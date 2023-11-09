package utils

import (
	"net"
	"os"
)

func GetLocalIPv4() net.IP {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	var ip net.IP
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			ip = ipv4
		}
	}

	return ip
}
