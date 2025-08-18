package network

import (
	"os/exec"
	"regexp"
	"runtime"
)

func Ping(host string) string {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", "-w", "1", host)
	} else {
		cmd = exec.Command("ping", "-c", "1", "-W", "1", host)
	}

	out, _ := cmd.Output()

	output := string(out)

	var re *regexp.Regexp
	if runtime.GOOS == "windows" {
		re = regexp.MustCompile(`\[(\d+\.\d+\.\d+\.\d+)]`)
	} else {
		re = regexp.MustCompile(`\((\d+\.\d+\.\d+\.\d+)\)`)
	}

	match := re.FindStringSubmatch(output)
	if len(match) > 1 {
		return match[1]
	}

	return ""
}
