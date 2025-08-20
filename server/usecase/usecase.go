package usecase

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"server/dto"
	"server/pkg/network"
	"strings"

	"github.com/mitchellh/go-ps"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewUsecase() *Usecase {
	return &Usecase{}
}

type Usecase struct {
}

func (u *Usecase) Status() (*dto.VPNStatusResponse, error) {
	data := dto.VPNStatusResponse{
		Active: false,
		Server: "",
		Logs:   "",
	}

	f, err := os.ReadFile("/var/log/openvpn/forti.log")
	if err == nil {
		content := string(f)
		data.Active = strings.Contains(content, "Tunnel is up and running") && !strings.Contains(content, "VPN disconnected")
		data.Logs = content
	} else {
		log.Err(err).Msg("Error reading file forti.log")
	}

	h, p := readFileHostPort("/opt/app/forticonfig")
	if h != "" && p != "" {
		data.Server = h + ":" + p
	}

	vs := os.Getenv("VPN_SERVERS")
	vpnServers := strings.Split(vs, ",")

	ipPortServers := make([]string, 0)

	for _, server := range vpnServers {
		sp := strings.Split(server, ":")
		host := sp[0]
		port := ""
		if len(sp) > 1 {
			port = sp[1]
		}
		ip := network.Ping(host)
		if port == "" {
			ipPortServers = append(ipPortServers, host+" ("+ip+")")
			continue
		}
		ipPortServers = append(ipPortServers, host+":"+port+" ("+ip+":"+port+")")
	}
	data.Servers = ipPortServers

	return &data, nil
}

func (u *Usecase) Activate(body *dto.VPNActivateRequest) error {
	if len(body.Cookie) <= 2000 || len(body.Cookie) > 3000 {
		return errors.New("invalid cookie")
	}

	fortiConfigPath := "/opt/app/forticonfig"
	newContent := readFileAndReplaceHostPort(fortiConfigPath, body.Host, body.Port)

	err := os.WriteFile(fortiConfigPath, []byte(newContent), 0644)
	if err != nil {
		log.Err(err).Msg("Write file failed forticonfig")
		return err
	}

	err = os.WriteFile("/opt/app/forti-cookie.txt", []byte(body.Cookie), 0644)
	if err != nil {
		log.Err(err).Msg("Write file failed forti-cookie.txt")
		return err
	}

	// stop current vpn
	processes, err := ps.Processes()
	if err != nil {
		log.Err(err).Msg("Error getting processes")
		return err
	}
	targetProcessName := "openfortivpn"
	for _, p := range processes {
		if p.Executable() == targetProcessName {
			p, err := os.FindProcess(p.Pid())
			if err == nil {
				_ = p.Kill()
			}
		}
	}

	extraOptions := " "

	if log.Logger.GetLevel() == zerolog.TraceLevel {
		extraOptions += "-v"
	}

	// start vpn
	command := `nohup openfortivpn` + extraOptions + ` --config /opt/app/forticonfig --cookie="` + body.Cookie + `" > /var/log/openvpn/forti.log 2>&1 &`
	cmd := exec.Command("/bin/sh", "-c", command)
	err = cmd.Start()
	if err != nil {
		log.Printf("Error starting vpn: %v", err)
		return err
	}

	return nil
}

func readFileAndReplaceHostPort(filename, host, port string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		return string(data)
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Err(err).Msg("Error opening file")
		return string(data)
	}
	defer file.Close()

	var result strings.Builder
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "host =") {
			line = "host = " + host
		} else if strings.HasPrefix(strings.TrimSpace(line), "port =") {
			line = "port = " + port
		}
		result.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		return string(data)
	}

	return result.String()
}

func readFileHostPort(filename string) (host string, port string) {
	host = ""
	port = ""

	file, err := os.Open(filename)
	if err != nil {
		log.Err(err).Str("filename", filename).Msg("Error opening file")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(strings.ReplaceAll(line, " ", ""))
		if strings.HasPrefix(line, "host=") {
			host = strings.ReplaceAll(line, "host=", "")
		} else if strings.HasPrefix(line, "port=") {
			port = strings.ReplaceAll(line, "port=", "")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Err(err).Msg("Error reading file scanner")
	}

	return
}
