package sysinfo

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ResourcesInfo contiene información de recursos del sistema.
type ResourcesInfo struct {
	CPU    CPUInfo  `json:"cpu"`    // Información de CPU
	Memory MemoryInfo `json:"memory"` // Información de memoria
}

// CPUInfo contiene información de CPU.
type CPUInfo struct {
	Usage    float64 `json:"usage"`     // Porcentaje de uso (0-100)
	Cores    int     `json:"cores"`      // Número de cores
	Model    string  `json:"model"`     // Modelo del CPU
	MHz     float64 `json:"mhz"`       // Velocidad en MHz
	LoadAvg  []float64 `json:"load_avg"` // Carga del sistema (1min, 5min, 15min)
}

// MemoryInfo contiene información de memoria.
type MemoryInfo struct {
	Total       uint64  `json:"total"`        // Memoria total en bytes
	Used        uint64  `json:"used"`        // Memoria usada en bytes
	Free        uint64  `json:"free"`        // Memoria libre en bytes
	Available   uint64  `json:"available"`    // Memoria disponible en bytes
	UsagePercent float64 `json:"usage_percent"` // Porcentaje de uso
}

// GetResourcesInfo obtiene información de recursos del sistema.
func GetResourcesInfo() (*ResourcesInfo, error) {
	cpu, err := GetCPUInfo()
	if err != nil {
		return nil, err
	}

	mem, err := GetMemoryInfo()
	if err != nil {
		return nil, err
	}

	return &ResourcesInfo{
		CPU:    cpu,
		Memory: mem,
	}, nil
}

// GetCPUInfo obtiene información de CPU.
func GetCPUInfo() (CPUInfo, error) {
	var info CPUInfo

	// Obtener número de cores
	cores, err := getCPUCores()
	if err == nil {
		info.Cores = cores
	}

	// Obtener modelo
	model, err := getCPUModel()
	if err == nil {
		info.Model = model
	}

	// Obtener MHz
	mhz, err := getCPUMHz()
	if err == nil {
		info.MHz = mhz
	}

	// Obtener uso de CPU (lectura simple)
	usage, _ := getCPUUsage()
	info.Usage = usage

	// Obtener load average
	loadAvg, _ := getLoadAverage()
	info.LoadAvg = loadAvg

	return info, nil
}

// GetMemoryInfo obtiene información de memoria.
func GetMemoryInfo() (MemoryInfo, error) {
	var info MemoryInfo

	// Leer /proc/meminfo
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return info, fmt.Errorf("error al leer meminfo: %w", err)
	}

	// Parser simple de meminfo
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Formato: "MemTotal:  Kb"
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Extraer número (ignora "kB" o similar)
		re := regexp.MustCompile(`(\d+)`)
		matches := re.FindStringSubmatch(value)
		if len(matches) == 0 {
			continue
		}

		kb, err := strconv.ParseUint(matches[1], 10, 64)
		if err != nil {
			continue
		}

		bytes := kb * 1024 // Convertir a bytes

		switch key {
		case "MemTotal":
			info.Total = bytes
		case "MemUsed":
			info.Used = bytes
		case "MemFree":
			info.Free = bytes
		case "MemAvailable":
			info.Available = bytes
		case "Available":
			// A veces está en "Available"
			if info.Available == 0 {
				info.Available = bytes
			}
		}
	}

	// Calcular used si no está مباشر
	if info.Used == 0 && info.Total > 0 {
		info.Used = info.Total - info.Free
	}

	// Calcular available si no está مباشر
	if info.Available == 0 && info.Total > 0 {
		info.Available = info.Free
	}

	// Calcular porcentaje
	if info.Total > 0 {
		info.UsagePercent = float64(info.Used) / float64(info.Total) * 100
	}

	return info, nil
}

// getCPUCores obtiene el número de cores de CPU.
func getCPUCores() (int, error) {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return 0, err
	}

	// Contarprocessor
	re := regexp.MustCompile(`(?m)^processor\s*:`)
	matches := re.FindAllStringIndex(string(data), -1)
	if matches == nil {
		return 0, fmt.Errorf("no se encontró información de processors")
	}

	return len(matches), nil
}

// getCPUModel obtiene el modelo del CPU.
func getCPUModel() (string, error) {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "", err
	}

	// Buscar "model name"
	re := regexp.MustCompile(`(?m)^model name\s*:\s*(.+)`)
	matches := re.FindStringSubmatch(string(data))
	if len(matches) < 2 {
		return "", fmt.Errorf("no se encontró model name")
	}

	return strings.TrimSpace(matches[1]), nil
}

// getCPUMHz obtiene la velocidad del CPU en MHz.
func getCPUMHz() (float64, error) {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return 0, err
	}

	// Buscar "cpu MHz"
	re := regexp.MustCompile(`(?m)^cpu MHz\s*:\s*([\d.]+)`)
	matches := re.FindStringSubmatch(string(data))
	if len(matches) < 2 {
		return 0, fmt.Errorf("no se encontró cpu MHz")
	}

	mhz, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, err
	}

	return mhz, nil
}

// getCPUUsage obtiene el porcentaje de uso de CPU.
// Esta es una implementación simple que retorna 0 para la primera lectura.
// En uso real, se necesitaría un segundo sample.
func getCPUUsage() (float64, error) {
	// Implementación simple basada en /proc/stat
	// Para uso real, se necesitan 2 lecturas con intervalo
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, err
	}

	// Primera línea: cpu  user nice system idle iowait irq softirq...
	re := regexp.MustCompile(`(?m)^cpu\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)`)
	matches := re.FindStringSubmatch(string(data))
	if len(matches) < 5 {
		return 0, fmt.Errorf("formato de stat inválido")
	}

	user, _ := strconv.ParseUint(matches[1], 10, 64)
	nice, _ := strconv.ParseUint(matches[2], 10, 64)
	sys, _ := strconv.ParseUint(matches[3], 10, 64)
	idle, _ := strconv.ParseUint(matches[4], 10, 64)
	iowait, _ := strconv.ParseUint(matches[5], 10, 64)

	total := user + nice + sys + idle + iowait
	if total == 0 {
		return 0, nil
	}

	used := user + nice + sys
	usage := float64(used) / float64(total) * 100

	return usage, nil
}

// getLoadAverage obtiene el load average del sistema.
func getLoadAverage() ([]float64, error) {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return nil, err
	}

	// Formato: "0.52 0.58 0.59 1/245 1234"
	re := regexp.MustCompile(`([\d.]+)\s+([\d.]+)\s+([\d.]+)`)
	matches := re.FindStringSubmatch(string(data))
	if len(matches) < 4 {
		return nil, fmt.Errorf("formato de loadavg inválido")
	}

	result := make([]float64, 3)
	for i := 1; i <= 3; i++ {
		val, err := strconv.ParseFloat(matches[i], 64)
		if err != nil {
			return nil, err
		}
		result[i-1] = val
	}

	return result, nil
}