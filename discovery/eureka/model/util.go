package model

import (
	"fmt"
	"github.com/google/uuid"
	"net"
	"os"
	"strings"
)

func asPort(port int, enabled bool) Port {
	return Port{
		Number:  port,
		Enabled: enabled,
	}
}

func toInstanceUrlPathToUrl(host string, port int, path string, defaultPath string) string {
	p := path
	if len(p) <= 0 {
		p = defaultPath
	}
	p = strings.TrimPrefix(p, "/")
	return fmt.Sprintf(nonSecureUrlPath, host, port, p)
}

func getLocalIP(address string) string {

	if address != "" {
		return os.ExpandEnv(address)
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	panic("Unable to determine local IP address (non loopback). Exiting.")
}

func generateID(fields ...string) string {
	return strings.Join(append(fields, uuid.New().String()), ":")
}
