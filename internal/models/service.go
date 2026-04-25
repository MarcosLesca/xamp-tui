package models

import (
	"time"
)

// Service representa un servicio de XAMPP (Apache, MySQL, PostgreSQL, etc.)
type Service struct {
	Name    string    `json:"name"`    // Nombre del servicio (mysql, apache2, postgresql, etc.)
	Status string    `json:"status"`  // Estado: running, stopped, error
	Version string   `json:"version"` // Versión instalada
	Port    int      `json:"port"`    // Puerto que usa el servicio
	Uptime time.Duration `json:"uptime"` // Tiempo ejecutándose
	PID    int      `json:"pid"`    // PID del proceso
}

// IsRunning verifica si el servicio está corriendo.
func (s *Service) IsRunning() bool {
	return s.Status == "running" && s.PID > 0
}

// Stopped devuelve true si el servicio está detenido.
func (s *Service) IsStopped() bool {
	return s.Status == "stopped"
}

// AllServices devuelve la lista de servicios típicos de XAMPP.
func AllServices(stackType StackType) []Service {
	switch stackType {
	case StackTypeLAMP:
		return []Service{
			{Name: "apache2", Status: "stopped", Version: "", Port: 80, Uptime: 0, PID: 0},
			{Name: "mysql", Status: "stopped", Version: "", Port: 3306, Uptime: 0, PID: 0},
		}
	case StackTypeLAMM:
		return []Service{
			{Name: "apache2", Status: "stopped", Version: "", Port: 80, Uptime: 0, PID: 0},
			{Name: "mariadb", Status: "stopped", Version: "", Port: 3306, Uptime: 0, PID: 0},
		}
	case StackTypeLEPP:
		return []Service{
			{Name: "apache2", Status: "stopped", Version: "", Port: 80, Uptime: 0, PID: 0},
			{Name: "postgresql", Status: "stopped", Version: "", Port: 5432, Uptime: 0, PID: 0},
		}
	default:
		return nil
	}
}