package main

import (
	"fmt"
	"net"
	"time"
)

func GetUpstream() (string, error) {
	dnsServers := []string{
		"8.8.8.8:53",        // Google Public DNS
		"1.1.1.1:53",        // Cloudflare DNS
		"208.67.222.222:53", // OpenDNS
		"9.9.9.9:53",        // Quad9 DNS
		"8.26.56.26:53",     // Comodo Secure DNS
		"4.2.2.1:53",        // Level3 DNS
		"77.88.8.8:53",      // Yandex DNS Basic
	}

	for _, server := range dnsServers {
		if isDNSServerReachable(server) {
			return server, nil
		}
	}

	return "", fmt.Errorf("no DNS servers are reachable")
}

func isDNSServerReachable(server string) bool {
	conn, err := net.DialTimeout("udp", server, 2*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}