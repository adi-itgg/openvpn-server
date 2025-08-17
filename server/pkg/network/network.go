package network

import (
	"fmt"
	"net"
)

func GetIPFormDNS(domain string) string {
	host := "google.com"

	ips, err := net.LookupIP(host)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	for _, ip := range ips {
		if ip.To4() != nil {
			return ip.String()
		}
	}
	return ""
}
