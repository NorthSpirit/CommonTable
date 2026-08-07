package discovery

import (
	"log"
	"net"
)

// PrintLANAddresses logs every non-loopback IPv4 address currently bound
// to this machine, on every server startup — so you can see at a glance
// (in console or the log file) which address to bookmark/type on a phone
// or PC, without needing to run ipconfig/ifconfig by hand each time the
// DHCP lease changes.
func PrintLANAddresses() {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Printf("discovery: could not list network interfaces: %v", err)
		return
	}

	log.Println("discovery: current LAN addresses —")

	found := false
	for _, iface := range interfaces {
		// Skip interfaces that are down or loopback-only.
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.To4() == nil {
				continue // IPv4 only — matches what the phone/discovery flow actually uses.
			}

			log.Printf("discovery:   %s -> %s", iface.Name, ip.String())
			found = true
		}
	}

	if !found {
		log.Println("discovery:   (no active IPv4 interfaces found)")
	}
}
