package service

import "xampp-tui/internal/models"

// ServiceManager define la interfaz para gestionar servicios del sistema.
// Implementaciones concretas deben manejar servicios como apache2, mysql, postgresql, etc.
type ServiceManager interface {
	// DetectServices detecta los servicios instalados en el sistema.
	// Retorna una lista de servicios encontrados con su información básica.
	DetectServices() ([]models.Service, error)

	// GetStatus obtiene el estado actual de un servicio específico.
	// Retorna el estado: running, stopped, error
	GetStatus(serviceName string) (string, error)

	// Start inicia un servicio.
	Start(serviceName string) error

	// Stop detiene un servicio.
	Stop(serviceName string) error

	// Restart reinicia un servicio.
	Restart(serviceName string) error

	// GetVersion obtiene la versión de un servicio.
	// Retorna string con la versión o empty si no está instalado.
	GetVersion(serviceName string) (string, error)

	// InstallStack instala el stack seleccionado (LAMP o LEPP).
	// Retorna los comandos ejecutados y cualquier error.
	InstallStack(stackType models.StackType) (string, error)

	// HasRootAccess verifica si tiene permisos de root.
	HasRootAccess() bool

	// ChangePort cambia el puerto de un servicio.
	// Retorna los comandos ejecutados y cualquier error.
	ChangePort(serviceName string, newPort int) (string, error)
}

// StatusRunning es el estado "running".
const StatusRunning = "running"

// StatusStopped es el estado "stopped".
const StatusStopped = "stopped"

// StatusError es el estado "error".
const StatusError = "error"

// KnownServices define los servicios que xampp-tui puede gestionar.
// NOTA: "phpmyadmin" es especial - no es un servicio systemd, es una app web.
var KnownServices = []string{
	"apache2",
	"nginx",
	"mysql",
	"mariadb",
	"postgresql",
	"php-fpm",
	"phpmyadmin",
}

// ServicePorts mapea servicios a sus puertos por defecto.
var ServicePorts = map[string]int{
	"apache2":    80,
	"nginx":     80,
	"mysql":     3306,
	"mariadb":   3306,
	"postgresql": 5432,
	"php-fpm":   9000,
	"phpmyadmin": 80, // Same as Apache - acceso via web
}