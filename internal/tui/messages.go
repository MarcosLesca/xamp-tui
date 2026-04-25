package tui

// PollingTick es un mensaje personalizado para el polling de servicios.
type PollingTick struct{}

// ServiceAction es un mensaje para ejecutar una acción en un servicio.
type ServiceAction struct {
	ServiceName string
	Action    string // "start", "stop", "restart"
}

// ErrorMsg es un mensaje de error.
type ErrorMsg string

// StackSelected es un mensaje cuando se selecciona un stack.
type StackSelected struct {
	StackType string
}

// ServiceSelected es un mensaje cuando se selecciona un servicio.
type ServiceSelected struct {
	ServiceName string
}

// InstallComplete es un mensaje cuando termina la instalación.
type InstallComplete struct {
	Log string
	Err error
}

// InstallProgress es un mensaje de progreso durante la instalación.
type InstallProgress struct {
	Step    int    // 1-4
	Total  int    // total de pasos
	Message string // mensaje del paso actual
}

// CheckInstallStatus es un mensaje para verificar el estado de instalación.
type CheckInstallStatus struct{}