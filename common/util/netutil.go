package util

import (
	"math/rand"
	"net"

	"starmesh/common/model"
)

// NormalizeIP converts IPv6 loopback and IPv4-mapped IPv6 to IPv4
func NormalizeIP(ip string) string {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return ip // Invalid IP, return as-is
	}

	// Handle IPv6 loopback (::1)
	if parsedIP.IsLoopback() {
		return "127.0.0.1"
	}

	// If IPv6-mapped IPv4 (e.g., ::ffff:192.168.0.1), extract IPv4
	if ipv4 := parsedIP.To4(); ipv4 != nil {
		return ipv4.String()
	}

	return ip
}

// GetOutboundIP returns the preferred outbound IP of this machine
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80") // 외부 주소 사용 (구글 DNS)
	if err != nil {
		return "unknown"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// selectRandomPeers picks up to maxCount random peers from the given peer map
func SelectRandomPeers(peerMap map[string]model.Peer, maxCount int) []model.Peer {
	// Prepare slice
	allPeers := make([]model.Peer, 0, len(peerMap))
	for _, p := range peerMap {
		allPeers = append(allPeers, p)
	}

	// Random shuffle (Go 1.20+: no need to seed)
	rand.Shuffle(len(allPeers), func(i, j int) {
		allPeers[i], allPeers[j] = allPeers[j], allPeers[i]
	})

	// Return up to maxCount
	if len(allPeers) > maxCount {
		return allPeers[:maxCount]
	}
	return allPeers
}
