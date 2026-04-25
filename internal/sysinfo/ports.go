package sysinfo

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// PortInfo contiene información de un puerto.
type PortInfo struct {
	Port    int    `json:"port"`     // Número de puerto
	Status string `json:"status"`   // "open" o "closed"
	Service string `json:"service"` // Servicio típicamente asociado
}

// GetPortsInfo obtiene el estado de los puertos comunes de XAMPP.
func GetPortsInfo() ([]PortInfo, error) {
	commonPorts := []int{80, 443, 3306, 5432, 9000, 9200, 8080, 8443}

	var results []PortInfo
	for _, port := range commonPorts {
		status, err := checkPort(port)
		if err != nil {
			status = "closed"
		}

		service := getServiceName(port)
		results = append(results, PortInfo{
			Port:    port,
			Status:  status,
			Service: service,
		})
	}

	return results, nil
}

// checkPort verifica si un puerto está abierto.
// Intenta conectar al puerto en localhost.
func checkPort(port int) (string, error) {
	// Intentar escuchar en el puerto
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		// Si falla, el puerto podría estar en uso
		// Verificar conectando
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return "closed", nil
		}
		conn.Close()
		return "open", nil
	}
	ln.Close()
	return "open", nil
}

// checkPortByFile verifica si un puerto está en uso leyendo /proc/net/tcp.
func checkPortByFile(port int) (bool, error) {
	data, err := os.ReadFile("/proc/net/tcp")
	if err != nil {
		return false, err
	}

	// Puerto en hex
	hexPort := fmt.Sprintf("%04X", port)

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if !strings.Contains(line, hexPort) {
			continue
		}

		//Formato: sl local_address rem_address st ...
		parts := strings.Fields(line)
		if len(parts) < 4 {
			continue
		}

		// Estado 0A = LISTEN, 01 = ESTABLISHED
		state := parts[3]
		if state == "0A" || state == "01" || state == "06" || state == "07" || state == "08" {
			return true, nil
		}
	}

	return false, nil
}

// IsPortOpen verifica si un puerto específico está abierto.
func IsPortOpen(port int) (bool, error) {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		// Puerto en uso
		return true, nil
	}
	ln.Close()
	return false, nil
}

// IsPortInUse verifica si un puerto específico está en uso.
// Retorna el PID del proceso que lo usa si está disponible.
func IsPortInUse(port int) (bool, int, error) {
	// Intentar conectar
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return false, 0, nil
	}
	conn.Close()

	// Puerto está en uso, buscar PID
	pid, err := getPIDByPort(port)
	return true, pid, err
}

// getPIDByPort intenta obtener el PID de un proceso que escucha en un puerto.
// Esta es una aproximación simple.
func getPIDByPort(port int) (int, error) {
	// Buscar en /proc todos los directorios numéricos
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return 0, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Solo directorios numéricos (PIDs)
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		// Buscar en fd
		fdPath := fmt.Sprintf("/proc/%d/fd", pid)
		fds, err := os.ReadDir(fdPath)
		if err != nil {
			continue
		}

		for _, fd := range fds {
			// Leer el symlink
			link, err := os.Readlink(fdPath + "/" + fd.Name())
			if err != nil {
				continue
			}

			// Buscar el puerto en el link
			if strings.Contains(link, fmt.Sprintf("socket:[%d]", port)) {
				return pid, nil
			}

			// También buscar formato de pipe
			if strings.Contains(link, fmt.Sprintf("pipe:[%d]", port)) {
				return pid, nil
			}
		}
	}

	return 0, fmt.Errorf("no se encontró proceso")
}

// getServiceName devuelve el nombre del servicio típicamente asociado a un puerto.
func getServiceName(port int) string {
	services := map[int]string{
		80:   "nginx/apache",
		443:  "nginx/apache-ssl",
		3306: "mysql",
		5432: "postgresql",
		9000: "php-fpm",
		9200: "elasticsearch",
		8080: "nginx/apache-alt",
		8443: "nginx/apache-ssl-alt",
	}

	name, ok := services[port]
	if !ok {
		return "unknown"
	}
	return name
}

// GetPortInfo obtiene información de un puerto específico.
func GetPortInfo(port int) (*PortInfo, error) {
	open, err := IsPortOpen(port)
	if err != nil {
		return nil, err
	}

	status := "closed"
	if open {
		status = "open"
	}

	return &PortInfo{
		Port:    port,
		Status:  status,
		Service: getServiceName(port),
	}, nil
}

// GetOpenPorts retorna una lista de puertos abiertos comunes.
func GetOpenPorts() ([]int, error) {
	ports, err := GetPortsInfo()
	if err != nil {
		return nil, err
	}

	var open []int
	for _, p := range ports {
		if p.Status == "open" {
			open = append(open, p.Port)
		}
	}

	return open, nil
}