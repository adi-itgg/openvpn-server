package usecase

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"server/dto"
	"server/pkg/network"
	"strings"

	"github.com/mitchellh/go-ps"
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
		fmt.Printf("Error reading file forti.log: %v\n", err)
	}

	f, err = os.ReadFile("/opt/app/forticonfig")
	if err == nil {
		host := ""
		port := ""
		content := string(f)
		for _, line := range strings.Split("\n", content) {
			if strings.Contains(line, "host") {
				host = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(line, " ", ""), "host=", ""))
			} else if strings.Contains(line, "port") {
				port = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(line, " ", ""), "port=", ""))
			}
		}
		data.Server = host + ":" + port
	} else {
		fmt.Printf("Error reading file forticonfig: %v\n", err)
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
		ip := network.GetIPFormDNS(host)
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
		return fmt.Errorf("invalid cookie")
	}

	fortiConfigPath := "/opt/app/forticonfig"
	newContent := readFileAndReplaceHostPort(fortiConfigPath, body.Host, body.Port)

	err := os.WriteFile(fortiConfigPath, []byte(newContent), 0644)
	if err != nil {
		fmt.Printf("Write file failed: %v\n", err)
		return err
	}

	err = os.WriteFile("/opt/app/forti-cookie.txt", []byte(body.Cookie), 0644)
	if err != nil {
		fmt.Printf("Write file failed: %v\n", err)
		return err
	}

	// stop current vpn
	processes, err := ps.Processes()
	if err != nil {
		fmt.Println("Error getting processes:", err)
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

	// start vpn
	command := `nohup openfortivpn --config /opt/app/forticonfig --cookie="` + body.Cookie + `" > /var/log/openvpn/forti.log 2>&1 &`
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
