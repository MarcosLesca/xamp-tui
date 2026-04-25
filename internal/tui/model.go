package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"xampp-tui/internal/config"
	"xampp-tui/internal/models"
	"xampp-tui/internal/service"
	"xampp-tui/internal/tui/screens"
)

// Model es el modelo principal de la aplicación TUI.
// Maneja el estado de la aplicación y la máquina de estados.
type Model struct {
	// Screen representa la pantalla actual
	Screen models.Screen

	// Services contiene los servicios detectados
	Services []models.Service

	// Config es la configuración de la aplicación
	Config *config.Config

	// ServiceManager gestiona los servicios del sistema
	ServiceManager service.ServiceManager

	// SelectedServiceIndex es el índice del servicio seleccionado (para navegación)
	SelectedServiceIndex int

	// SelectedStackIndex es el índice del stack seleccionado
	SelectedStackIndex int

	// SelectedMenuIndex es el índice del menú seleccionado
	SelectedMenuIndex int

	// LastKey fue la última tecla presionada
	LastKey string

	// Width y Height almacenan el tamaño de la ventana
	Width  int
	Height int

	// pollingTimer maneja el auto-refresh del dashboard
	pollingTimer *time.Timer

	// QuitChannel para señal de salida
	Quitting bool

	// ErrMsg contiene el último error
	ErrMsg string

// Installing indica si está en proceso de instalación
	Installing bool

	// InstallLog contiene el log de instalación
	InstallLog string

	// PortEditValue contiene el valor del puerto siendo editado
	PortEditValue string

	// PortEditResult guarda el resultado del último cambio de puerto
	PortEditResult struct {
		Success bool
		Message string
	}

	// StackSelectMode: "install" o "manage" - para saber qué flujo traen a StackSelect
	StackSelectMode string

	// LAMMConfig guarda la configuración para instalación LAMM
	LAMMConfig *models.LAMMConfig

	// ConfigEditing indica si estamos editando la configuración
	ConfigEditing bool
}

// Init implementa tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// View implementa tea.Model.
func (m Model) View() string {
	return renderView(m)
}

// renderView interna
func renderView(m Model) string {
	switch m.Screen {
	case models.ScreenWelcome:
		// Verificar si hay servicios INSTALADOS (running, stopped, o installed), no solo "available"
		hasInstalledServices := false
		for _, svc := range m.Services {
			if svc.Status == "running" || svc.Status == "stopped" || svc.Status == "installed" {
				hasInstalledServices = true
				break
			}
		}
		return screens.WelcomeView(m.SelectedMenuIndex, m.Width, hasInstalledServices)

	case models.ScreenStackSelect:
		return screens.StackSelectView(m.SelectedMenuIndex, m.Width)

	case models.ScreenInstall:
		// Verificar si tiene permisos de root
		canInstall := m.HasRootAccess()
		return screens.InstallView(m.Config.StackType, m.Width, m.SelectedMenuIndex, canInstall)

	case models.ScreenConfig:
		// Configuración antes de instalar
		lammConfig := m.LAMMConfig
		if lammConfig == nil {
			c := models.DefaultLAMMConfig()
			lammConfig = &c
		}
		return screens.ConfigView(m.Config.StackType, m.Width, m.SelectedMenuIndex, lammConfig, m.ConfigEditing)

	case models.ScreenDashboard:
		return screens.DashboardView(m.Services, m.Width, m.Installing, m.InstallLog, m.SelectedServiceIndex)

	case models.ScreenDetails:
		svc := m.GetServiceByIndex(m.SelectedServiceIndex)
		return screens.DetailsView(svc, m.Width, m.SelectedMenuIndex)

	case models.ScreenLogs:
		svc := m.GetServiceByIndex(m.SelectedServiceIndex)
		return screens.LogsView(svc, m.Width, m.SelectedMenuIndex)

	case models.ScreenDatabase:
		return screens.DatabaseView(m.Config.StackType, m.Width)

	case models.ScreenPortEdit:
		svc := m.GetServiceByIndex(m.SelectedServiceIndex)
		return screens.PortEditView(svc, m.Width, m.PortEditValue)

	default:
		return screens.ErrorView("Pantalla desconocida")
	}
}

// New crea un nuevo Modelo con valores por defecto.
func New() Model {
	return Model{
		Screen:               models.ScreenWelcome,
		Services:            []models.Service{},
		Config:               config.Default(),
		ServiceManager:       service.NewLinuxServiceManager(),
		SelectedServiceIndex: 0,
		SelectedStackIndex:   0,
		SelectedMenuIndex:    0,
		LastKey:              "",
		Width:                80,
		Height:               24,
		pollingTimer:         nil,
		Quitting:             false,
		ErrMsg:               "",
	}
}

// InitServices detecta los servicios instalados en el sistema.
func (m *Model) InitServices() error {
	svcs, err := m.ServiceManager.DetectServices()
	if err != nil {
		return err
	}
	m.Services = svcs
	// Asegurar que el índice sea válido
	if m.SelectedServiceIndex >= len(m.Services) {
		m.SelectedServiceIndex = 0
	}
	return nil
}

// GetServices retorna la lista de servicios.
func (m *Model) GetServices() []models.Service {
	return m.Services
}

// GetServiceByName retorna un servicio por su nombre.
func (m *Model) GetServiceByName(name string) *models.Service {
	for i := range m.Services {
		if m.Services[i].Name == name {
			return &m.Services[i]
		}
	}
	return nil
}

// RefreshServices actualiza el estado de los servicios.
func (m *Model) RefreshServices() error {
	return m.InitServices()
}

// SetScreen cambia a una pantalla específica.
func (m *Model) SetScreen(screen models.Screen) {
	m.Screen = screen
	m.SelectedMenuIndex = 0
	m.SelectedServiceIndex = 0
}

// StartPolling inicia el auto-refresh de servicios.
func (m *Model) StartPolling() {
	if m.pollingTimer != nil {
		return
	}
	m.pollingTimer = time.NewTimer(5 * time.Second)
}

// StopPolling detiene el auto-refresh.
func (m *Model) StopPolling() {
	if m.pollingTimer != nil {
		m.pollingTimer.Stop()
		m.pollingTimer = nil
	}
}

// GetPollingChannel retorna el canal de timer para polling.
func (m *Model) GetPollingChannel() <-chan time.Time {
	if m.pollingTimer == nil {
		return nil
	}
	return m.pollingTimer.C
}

// ResetSelection resetea la selección al iniciar una pantalla.
func (m *Model) ResetSelection() {
	m.SelectedServiceIndex = 0
	m.SelectedStackIndex = 0
	m.SelectedMenuIndex = 0
}

// GetServiceByIndex retorna el servicio en el índice dado.
func (m *Model) GetServiceByIndex(index int) *models.Service {
	if index < 0 || index >= len(m.Services) {
		return nil
	}
	return &m.Services[index]
}

// HasRootAccess verifica si tiene permisos para ejecutar comandos como root.
func (m *Model) HasRootAccess() bool {
	return m.ServiceManager.HasRootAccess()
}