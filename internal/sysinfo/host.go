package sysinfo

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// HostInfo contiene información del host.
type HostInfo struct {
	Hostname string `json:"hostname"` // Nombre del host
	IP      string `json:"ip"`      // Dirección IP primaria
}

// GetHostInfo obtiene información del host.
func GetHostInfo() (*HostInfo, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	ip, err := getPrimaryIP()
	if err != nil {
		ip = ""
	}

	return &HostInfo{
		Hostname: hostname,
		IP:      ip,
	}, nil
}

// getPrimaryIP obtiene la dirección IP primaria del sistema.
// Busca la IP de la interfaz que no sea loopback.
func getPrimaryIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("error al obtener interfaces: %w", err)
	}

	for _, iface := range interfaces {
		// Ignorar loopback
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// Ignorar down
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		// Obtener direcciones
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

			if ip == nil {
				continue
			}

			// Solo IPv4
			if ip4 := ip.To4(); ip4 != nil {
				return ip4.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no se encontró IP primaria")
}

// GetHostname obtiene solo el nombre del host.
func GetHostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("error al obtener hostname: %w", err)
	}
	return hostname, nil
}

// GetIP obtiene solo la IP primaria.
func GetIP() (string, error) {
	ip, err := getPrimaryIP()
	if err != nil {
		return "", fmt.Errorf("error al obtener IP: %w", err)
	}
	return ip, nil
}

// GetFQDN obtiene el nombre de dominio completo.
func GetFQDN() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	// Intentar resolver el hostname
	names, err := net.LookupAddr(hostname)
	if err != nil || len(names) == 0 {
		return hostname, nil
	}

	// Limpiar trailing dot
	return strings.TrimSuffix(names[0], "."), nil
}