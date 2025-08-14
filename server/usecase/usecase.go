package usecase

import (
	"fmt"
	"github.com/mitchellh/go-ps"
	"log"
	"os"
	"os/exec"
	"server/dto"
	"strings"
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

	return &data, nil
}

func (u *Usecase) Activate(body *dto.VPNActivateRequest) error {
	f, err := os.ReadFile("/opt/app/forticonfig")
	if err != nil {
		return err
	}

	content := string(f)
	newContent := ""
	for _, line := range strings.Split("\n", content) {
		if strings.Contains(line, "host") {
			line = "host = " + body.Host
		} else if strings.Contains(line, "port") {
			line = "port = " + body.Port
		}
		newContent += line + "\n"
	}

	err = os.WriteFile("/opt/app/forticonfig", []byte(content), 0644)
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
