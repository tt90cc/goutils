package utils

import (
	"net"
	"net/http"
	"strings"
)

// 检查 IP 是否为内网访问
func IsPrivateIP(r *http.Request) bool {
	clientIP := GetClientIP(r)
	parsedIP := net.ParseIP(clientIP)
	if parsedIP == nil {
		return false
	}

	privateBlocks := []*net.IPNet{
		{IP: net.ParseIP("127.0.0.0"), Mask: net.CIDRMask(8, 32)},
		{IP: net.ParseIP("10.0.0.0"), Mask: net.CIDRMask(8, 32)},
		{IP: net.ParseIP("172.16.0.0"), Mask: net.CIDRMask(12, 32)},
		{IP: net.ParseIP("192.168.0.0"), Mask: net.CIDRMask(16, 32)},
	}

	for _, block := range privateBlocks {
		if block.Contains(parsedIP) {
			return true
		}
	}

	return false
}

func GetClientIP(r *http.Request) string {
	var clientIP string

	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ", ")
		clientIP = ips[0]
	} else if rip := r.Header.Get("X-Real-IP"); rip != "" {
		clientIP = rip
	} else {
		clientIP = strings.Split(r.RemoteAddr, ":")[0]
	}

	return clientIP
}
