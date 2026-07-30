package discovery

import (
	"encoding/json"
	"log"
	"net"
)

const (
	discoveryPort  = 41242
	discoveryMagic = "COMMONTABLE_DISCOVER"
)

type discoveryResponse struct {
	Service string `json:"service"`
	Port    int    `json:"port"`
}

// StartResponder listens for UDP broadcast discovery requests and replies
// directly to the sender with the server's identity and HTTP port. This
// lets clients on the LAN find the server without a hardcoded IP — the
// server's actual IP is taken from the incoming packet's source address,
// which the OS fills in correctly regardless of DHCP changes.
func StartResponder(httpPort int) {
	addr := net.UDPAddr{Port: discoveryPort, IP: net.IPv4zero}
	conn, err := net.ListenUDP("udp4", &addr)
	if err != nil {
		log.Printf("discovery: failed to start UDP listener: %v", err)
		return
	}

	log.Printf("discovery: listening on UDP :%d", discoveryPort)

	go func() {
		defer conn.Close()
		buf := make([]byte, 1024)
		for {
			n, remoteAddr, err := conn.ReadFromUDP(buf)
			if err != nil {
				continue
			}
			if string(buf[:n]) != discoveryMagic {
				continue
			}

			resp, err := json.Marshal(discoveryResponse{Service: "commontable", Port: httpPort})
			if err != nil {
				continue
			}
			_, _ = conn.WriteToUDP(resp, remoteAddr)
		}
	}()
}
