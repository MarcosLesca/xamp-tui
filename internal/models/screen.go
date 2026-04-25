package models

// Screen representa la pantalla actual de la TUI.
type Screen string

const (
	ScreenWelcome    Screen = "welcome"    // Pantalla de bienvenida
	ScreenStackSelect Screen = "stackselect" // Selección de stack
	ScreenInstall    Screen = "install"   // Instalación
	ScreenConfig     Screen = "config"    // Configuración antes de instalar
	ScreenDashboard  Screen = "dashboard" // Panel principal
	ScreenDetails    Screen = "details"   // Detalles del servicio
	ScreenLogs      Screen = "logs"      // Ver logs
	ScreenDatabase  Screen = "database"  // Gestión de base de datos
	ScreenPortEdit Screen = "portedit" // Edición de puerto
)

// IsValid verifica si el Screen es válido.
func (s Screen) IsValid() bool {
	switch s {
	case ScreenWelcome, ScreenStackSelect, ScreenInstall, ScreenConfig, ScreenDashboard, ScreenDetails, ScreenLogs, ScreenDatabase, ScreenPortEdit:
		return true
	default:
		return false
	}
}

// AllScreens devuelve todos los valores de Screen válidos.
func AllScreens() []Screen {
	return []Screen{
		ScreenWelcome,
		ScreenStackSelect,
		ScreenInstall,
		ScreenConfig,
		ScreenDashboard,
		ScreenDetails,
		ScreenLogs,
		ScreenDatabase,
		ScreenPortEdit,
	}
}