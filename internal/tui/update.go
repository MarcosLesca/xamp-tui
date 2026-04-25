package tui

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/charmbracelet/bubbletea"

	"xampp-tui/internal/models"
)

// InstallProgressChan canal global para progreso de instalación
var InstallProgressChan = make(chan InstallProgress, 10)

// InstallCompleteChan canal global para instalación completa
var InstallCompleteChan = make(chan InstallComplete, 1)

// mu protects the channels
var mu sync.Mutex

// Update es la función principal de actualización del modelo.
// Maneja los mensajes y devuelve el modelo actualizado.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// tea.KeyMsg maneja los eventos de teclado
	case tea.KeyMsg:
		return handleKeyMsg(m, msg)

	// tea.WindowSizeMsg maneja cambios de tamaño de ventana
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	// Custom: PollingTick para auto-refresh
	case PollingTick:
		if m.Screen == models.ScreenDashboard {
			// Check for install progress
			select {
			case progress, ok := <-InstallProgressChan:
				if ok {
					msgStr := fmt.Sprintf("[%d/%d] %s", progress.Step, progress.Total, progress.Message)
					if m.InstallLog != "" {
						m.InstallLog += "\n" + msgStr
					} else {
						m.InstallLog = msgStr
					}
				}
			default:
			}
			
			select {
			case complete, ok := <-InstallCompleteChan:
				if ok {
					m.Installing = false
					if complete.Err != nil {
						m.ErrMsg = fmt.Sprintf("Install error: %v", complete.Err)
						m.InstallLog = complete.Log
					} else {
						m.InstallLog = complete.Log
					}
					m.RefreshServices()
				}
			default:
			}
			
			// Also refresh services
			if err := m.RefreshServices(); err != nil {
				log.Printf("Error refreshing services: %v", err)
			}
		}
		// Reiniciar timer de polling
		m.StartPolling()
		return m, nil

	// Custom: ServiceAction para acciones de servicio
	case ServiceAction:
		return handleServiceAction(m, msg)

	// Custom: ErrorMsg para errores
	case ErrorMsg:
		m.ErrMsg = string(msg)
		return m, nil

	// Custom: InstallComplete cuando termina la instalación
	case InstallComplete:
		m.Installing = false
		if msg.Err != nil {
			m.ErrMsg = fmt.Sprintf("Install error: %v", msg.Err)
			m.InstallLog = msg.Log
		} else {
			m.InstallLog = msg.Log
		}
		m.RefreshServices()
		return m, nil

	// Custom: InstallProgress para mostrar progreso en tiempo real
	case InstallProgress:
		msgStr := fmt.Sprintf("[%d/%d] %s", msg.Step, msg.Total, msg.Message)
		if m.InstallLog != "" {
			m.InstallLog += "\n" + msgStr
		} else {
			m.InstallLog = msgStr
		}
		return m, nil

	// Custom: CheckInstallStatus para verificar si terminó la instalación
	case CheckInstallStatus:
		// Verificar si hay nuevos servicios
		if err := m.RefreshServices(); err == nil {
			// Si hay servicios, terminó
			if len(m.Services) > 0 {
				m.Installing = false
				m.InstallLog = "✓ Installation complete! Services detected."
			}
		}
		// Reiniciar polling
		m.StartPolling()
		return m, nil

	default:
		return m, nil
	}
}

// handleKeyMsg maneja los mensajes de teclado.
func handleKeyMsg(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	m.LastKey = key

	switch key {
	// q o Ctrl+C para salir
	case "q", "ctrl+c":
		m.Quitting = true
		return m, tea.Quit

	// Escape para ir atrás
	case "esc":
		return handleEscape(m)

	// Flechas para navegación
	case "up", "k":
		return handleUp(m)

	case "down", "j":
		return handleDown(m)

	// Enter para seleccionar
	case "enter":
		return handleEnter(m)

	// Space para acción
	case " ":
		return handleSpace(m)

	// Teclas numéricas para edición de puerto
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		return handlePortDigit(m, key)

	default:
		return m, nil
	}
}

// handleEscape maneja la tecla Escape (volver atrás).
func handleEscape(m Model) (tea.Model, tea.Cmd) {
	switch m.Screen {
	case models.ScreenWelcome:
		m.Quitting = true
		return m, tea.Quit

	case models.ScreenStackSelect:
		m.SetScreen(models.ScreenWelcome)
		m.ResetSelection()

	case models.ScreenInstall:
		m.SetScreen(models.ScreenStackSelect)
		m.ResetSelection()

	case models.ScreenConfig:
		// Si está editando, salir del modo edición
		if m.ConfigEditing {
			m.ConfigEditing = false
			m.SelectedMenuIndex = 1 // Volver a "Edit Options"
		} else {
			m.SetScreen(models.ScreenStackSelect)
		}
		m.ResetSelection()

	case models.ScreenDashboard:
		m.SetScreen(models.ScreenWelcome)
		m.ResetSelection()
		m.StopPolling()

	case models.ScreenDetails:
		m.Screen = models.ScreenDashboard
		m.SelectedMenuIndex = 0
		m.StartPolling()

	case models.ScreenLogs:
		m.Screen = models.ScreenDashboard
		m.ResetSelection()
		m.StartPolling()

	case models.ScreenDatabase:
		m.SetScreen(models.ScreenDashboard)
		m.ResetSelection()
		m.StartPolling()

	case models.ScreenPortEdit:
		m.Screen = models.ScreenDetails
		m.SelectedMenuIndex = 0

	default:
		m.SetScreen(models.ScreenWelcome)
		m.ResetSelection()
	}
	return m, nil
}

// handleEnter maneja la tecla Enter.
func handleEnter(m Model) (tea.Model, tea.Cmd) {
	switch m.Screen {
	case models.ScreenWelcome:
		// Índice 0 = Install Stack, Índice 1 = Manage Stack, Índice 2 = Quit
		switch m.SelectedMenuIndex {
		case 0:
			m.StackSelectMode = "install"
			m.SetScreen(models.ScreenStackSelect)
		case 1:
			m.StackSelectMode = "manage"
			m.SelectedMenuIndex = 0
			m.SetScreen(models.ScreenStackSelect)
		case 2:
			m.Quitting = true
			return m, tea.Quit
		}
		m.ResetSelection()

	case models.ScreenStackSelect:
		// Seleccionar stack - caso 3 = Back
		switch m.SelectedMenuIndex {
		case 0:
			m.Config.StackType = models.StackTypeLAMP
		case 1:
			m.Config.StackType = models.StackTypeLAMM
		case 2:
			m.Config.StackType = models.StackTypeLEPP
		case 3:
			m.SetScreen(models.ScreenWelcome)
			m.ResetSelection()
			return m, nil
		}

		// Según el modo, decidir qué hacer
		if m.StackSelectMode == "manage" {
			m.SetScreen(models.ScreenDashboard)
			m.InitServices()
			m.Services = filterServicesByStack(m.Services, m.Config.StackType)
			m.StartPolling()
		} else {
			m.SetScreen(models.ScreenConfig)
			if m.Config.StackType == models.StackTypeLAMP || m.Config.StackType == models.StackTypeLAMM {
				m.LAMMConfig = new(models.LAMMConfig)
				*m.LAMMConfig = models.DefaultLAMMConfig()
			}
		}
		m.ResetSelection()

	case models.ScreenInstall:
		if m.SelectedMenuIndex == 0 {
			m.Installing = true
			m.InstallLog = "Installing..."
			m.SetScreen(models.ScreenDashboard)
			m.InitServices()
			// Filter services by stack
			m.Services = filterServicesByStack(m.Services, m.Config.StackType)
			stackType := m.Config.StackType
			svcMgr := m.ServiceManager

			// Start installation in goroutine with progress
			go func() {
				_, err := svcMgr.InstallStackWithProgress(stackType, func(step, total int, message string) {
					InstallProgressChan <- InstallProgress{Step: step, Total: total, Message: message}
				})
				logMsg, _ := svcMgr.InstallStack(stackType)
				InstallCompleteChan <- InstallComplete{Log: logMsg, Err: err}
			}()

			m.StartPolling()
			m.ResetSelection()
		} else {
			m.SetScreen(models.ScreenStackSelect)
		}

	case models.ScreenConfig:
		switch m.SelectedMenuIndex {
		case 0: // Start Installation
			m.Installing = true
			m.InstallLog = "Installing..."
			m.SetScreen(models.ScreenDashboard)
			m.InitServices()
			// Filter services by stack
			m.Services = filterServicesByStack(m.Services, m.Config.StackType)
			stackType := m.Config.StackType
			svcMgr := m.ServiceManager
			
			// Start installation in goroutine with progress
			go func() {
				InstallProgressChan <- InstallProgress{Step: 1, Total: 4, Message: "Starting..."}
				logMsg, err := svcMgr.InstallStackWithProgress(stackType, func(step, total int, message string) {
					InstallProgressChan <- InstallProgress{Step: step, Total: total, Message: message}
				})
				InstallCompleteChan <- InstallComplete{Log: logMsg, Err: err}
			}()
			
			// Start polling to check install progress
			m.StartPolling()
			return m, nil
		case 1: // Edit Options
			m.ConfigEditing = true
			m.SelectedMenuIndex = 0
		case 2: // Back
			m.SetScreen(models.ScreenStackSelect)
		}
		m.ResetSelection()

	case models.ScreenDashboard:
		m.StopPolling()
		m.Screen = models.ScreenDetails
		m.SelectedMenuIndex = 0
		return m, nil

	case models.ScreenDetails:
		return handleDetailsAction(m)

	case models.ScreenLogs:
		m.SetScreen(models.ScreenDatabase)
		m.ResetSelection()

	case models.ScreenDatabase:
		return handleDBAction(m)

	case models.ScreenPortEdit:
		return handlePortEditEnter(m)

	default:
		m.SetScreen(models.ScreenWelcome)
		m.ResetSelection()
	}

	return m, nil
}

// handleUp maneja la flecha arriba.
func handleUp(m Model) (tea.Model, tea.Cmd) {
	switch m.Screen {
	case models.ScreenWelcome, models.ScreenStackSelect:
		if m.SelectedMenuIndex > 0 {
			m.SelectedMenuIndex--
		}

	case models.ScreenDashboard:
		if m.SelectedServiceIndex > 0 {
			m.SelectedServiceIndex--
		}

	case models.ScreenDetails, models.ScreenLogs:
		if m.SelectedMenuIndex > 0 {
			m.SelectedMenuIndex--
		}

	case models.ScreenDatabase:
		if m.SelectedMenuIndex > 0 {
			m.SelectedMenuIndex--
		}

	case models.ScreenConfig:
		if m.ConfigEditing {
			// Editing mode: navigate config options (0-5)
			if m.SelectedMenuIndex > 0 {
				m.SelectedMenuIndex--
			}
		} else {
			// Menu mode: navigate menu (0-2)
			if m.SelectedMenuIndex > 0 {
				m.SelectedMenuIndex--
			}
		}

	case models.ScreenPortEdit:
		// ↑ Borra último dígito
		if len(m.PortEditValue) > 0 {
			m.PortEditValue = m.PortEditValue[:len(m.PortEditValue)-1]
		}

	default:
		if m.SelectedMenuIndex > 0 {
			m.SelectedMenuIndex--
		}
	}

	return m, nil
}

// handleDown maneja la flecha abajo.
func handleDown(m Model) (tea.Model, tea.Cmd) {
	switch m.Screen {
	case models.ScreenConfig:
		if m.ConfigEditing {
			// Editing mode: navigate config options (max 5 for LAMM/LAMP)
			if m.SelectedMenuIndex < 5 {
				m.SelectedMenuIndex++
			}
		} else {
			// Menu mode: navigate menu (max 2)
			if m.SelectedMenuIndex < 2 {
				m.SelectedMenuIndex++
			}
		}
		return m, nil
	
	default:
		// Use the standard getMaxIndex
		maxIndex := getMaxIndex(m)
		
		switch m.Screen {
		case models.ScreenWelcome, models.ScreenStackSelect:
			if m.SelectedMenuIndex < maxIndex {
				m.SelectedMenuIndex++
			}

		case models.ScreenDashboard:
			if m.SelectedServiceIndex < maxIndex {
				m.SelectedServiceIndex++
			}

		case models.ScreenDetails, models.ScreenLogs:
			if m.SelectedMenuIndex < maxIndex {
				m.SelectedMenuIndex++
			}

		case models.ScreenDatabase:
			if m.SelectedMenuIndex < maxIndex {
				m.SelectedMenuIndex++
			}

		case models.ScreenInstall:
			if m.SelectedMenuIndex < maxIndex {
				m.SelectedMenuIndex++
			}

		case models.ScreenPortEdit:
			if len(m.PortEditValue) < 5 {
				m.PortEditValue += "0"
			}

		default:
			if m.SelectedMenuIndex < maxIndex {
				m.SelectedMenuIndex++
			}
		}
	}

	return m, nil
}

// handleSpace maneja la barra espaciadora.
func handleSpace(m Model) (tea.Model, tea.Cmd) {
	switch m.Screen {
	case models.ScreenDetails:
		// Ejecutar acción de control (start/stop/restart)
		return handleDetailsAction(m)

	case models.ScreenConfig:
		// En modo edición, Space toggla el valor booleano
		return handleConfigEdit(m)

	default:
		return m, nil
	}
}

// handleConfigEdit/edita valores en Config screen.
func handleConfigEdit(m Model) (tea.Model, tea.Cmd) {
	if m.LAMMConfig == nil || !m.ConfigEditing {
		return m, nil
	}

	switch m.SelectedMenuIndex {
	case 0:
		// MariaDB Root Password - por ahora no editable (requiere input)
		m.ErrMsg = "Password editing not implemented - use terminal"
	case 1:
		m.LAMMConfig.RemoveAnonymous = !m.LAMMConfig.RemoveAnonymous
	case 2:
		m.LAMMConfig.DisallowRemote = !m.LAMMConfig.DisallowRemote
	case 3:
		m.LAMMConfig.RemoveTestDB = !m.LAMMConfig.RemoveTestDB
	case 4:
		// Web Server toggle
		if m.LAMMConfig.WebServer == "apache2" {
			m.LAMMConfig.WebServer = "nginx"
		} else {
			m.LAMMConfig.WebServer = "apache2"
		}
	case 5:
		// phpMyAdmin Path toggle
		if m.LAMMConfig.PHPMyAdminPath == "/phpmyadmin" {
			m.LAMMConfig.PHPMyAdminPath = "/pma"
		} else {
			m.LAMMConfig.PHPMyAdminPath = "/phpmyadmin"
		}
	}

	return m, nil
}

// handleDetailsAction maneja las acciones en la pantalla de detalles.
func handleDetailsAction(m Model) (tea.Model, tea.Cmd) {
	svc := m.GetServiceByIndex(m.SelectedServiceIndex)
	if svc == nil {
		return m, nil
	}

	running := svc.Status == "running"
	available := svc.Status == "available"

	if running {
		// Menú: Stop(0), Restart(1), View Logs(2), Change Port(3), Back(4)
		switch m.SelectedMenuIndex {
		case 0: // Stop
			if err := m.ServiceManager.Stop(svc.Name); err != nil {
				m.ErrMsg = err.Error()
			}
		case 1: // Restart
			if err := m.ServiceManager.Restart(svc.Name); err != nil {
				m.ErrMsg = err.Error()
			}
		case 2: // View Logs
			m.SetScreen(models.ScreenLogs)
			m.ResetSelection()
		case 3: // Change Port
			m.Screen = models.ScreenPortEdit
			m.PortEditValue = ""
		case 4: // Back
			m.Screen = models.ScreenDashboard
			m.SelectedMenuIndex = 0
			m.StartPolling()
		}
	} else if available {
		// Menú: Install(0), View Logs(1), Change Port(2), Back(3)
		switch m.SelectedMenuIndex {
		case 0: // Install
			m.ErrMsg = "Instalar " + svc.Name + " desde Stack Select"
			m.SetScreen(models.ScreenStackSelect)
			m.SelectedMenuIndex = 0
		case 1: // View Logs - no disponible para servicios no instalados
			m.ErrMsg = "Logs no disponibles para servicios no instalados"
		case 2: // Change Port
			m.ErrMsg = "Cambiar puerto requiere instalar primero"
		case 3: // Back
			m.Screen = models.ScreenDashboard
			m.SelectedMenuIndex = 0
			m.StartPolling()
		}
	} else {
		// Menú: Start(0), View Logs(1), Change Port(2), Back(3)
		switch m.SelectedMenuIndex {
		case 0: // Start
			if err := m.ServiceManager.Start(svc.Name); err != nil {
				m.ErrMsg = err.Error()
			}
		case 1: // View Logs
			m.SetScreen(models.ScreenLogs)
			m.ResetSelection()
		case 2: // Change Port
			m.Screen = models.ScreenPortEdit
			m.PortEditValue = ""
		case 3: // Back
			m.Screen = models.ScreenDashboard
			m.SelectedMenuIndex = 0
			m.StartPolling()
		}
	}

	// Refresh services after action
	if err := m.RefreshServices(); err != nil {
		fmt.Printf("Error refreshing: %v\n", err)
	}

	return m, nil
}

// handleDBAction maneja las acciones en la pantalla de database.
func handleDBAction(m Model) (tea.Model, tea.Cmd) {
	// Por implementar: crear/eliminar bases de datos
	return m, nil
}

// handleServiceAction maneja acciones de servicio.
func handleServiceAction(m Model, msg ServiceAction) (tea.Model, tea.Cmd) {
	var err error

	switch msg.Action {
	case "start":
		err = m.ServiceManager.Start(msg.ServiceName)
	case "stop":
		err = m.ServiceManager.Stop(msg.ServiceName)
	case "restart":
		err = m.ServiceManager.Restart(msg.ServiceName)
	}

	if err != nil {
		m.ErrMsg = err.Error()
	}

	m.RefreshServices()
	return m, nil
}

// getMaxIndex retorna el índice máximo según la pantalla actual.
func getMaxIndex(m Model) int {
	switch m.Screen {
	case models.ScreenWelcome, models.ScreenStackSelect:
		return 3 // 4 opciones de menú (LAMP, LAMM, LEPP, Back)

	case models.ScreenDashboard:
		if len(m.Services) == 0 {
			return 0
		}
		return len(m.Services) - 1

	case models.ScreenDetails:
		svc := m.GetServiceByIndex(m.SelectedServiceIndex)
		if svc != nil && svc.Status == "running" {
			return 4 // Stop, Restart, View Logs, Change Port, Back
		}
		if svc != nil && svc.Status == "available" {
			return 1 // Install, Back
		}
		return 3 // Start, View Logs, Change Port, Back

	case models.ScreenLogs:
		return 2 // Apache Logs, Database Logs, Back

	case models.ScreenDatabase:
		return 2 // Create DB, Drop DB, Back

	case models.ScreenInstall:
		return 1 // Start Installation, Back (2 options - 1 = max index - pero el menú muestra solo 2)

	case models.ScreenConfig:
		return 2 // Start Installation, Edit Options, Back (but when editing, 6 options)

	case models.ScreenPortEdit:
		return 0 // Solo un "modo" de edición por teclado

	default:
		return 0
	}
}

// handlePortDigit maneja dígitos para edición de puerto.
func handlePortDigit(m Model, digit string) (tea.Model, tea.Cmd) {
	if m.Screen != models.ScreenPortEdit {
		return m, nil
	}

	// Máximo 5 dígitos (65535)
	if len(m.PortEditValue) < 5 {
		// Verificar que no pase de 65535
		newVal := m.PortEditValue + digit
		if len(newVal) <= 5 {
			m.PortEditValue = newVal
		}
	}

	return m, nil
}

// handleEnter en PortEdit: confirmar cambio de puerto.
func handlePortEditEnter(m Model) (tea.Model, tea.Cmd) {
	svc := m.GetServiceByIndex(m.SelectedServiceIndex)
	if svc == nil {
		return m, nil
	}

	portStr := m.PortEditValue
	if portStr == "" {
		// No cambiar - volver
		m.Screen = models.ScreenDetails
		return m, nil
	}

	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		m.ErrMsg = "Puerto inválido (1-65535)"
		m.Screen = models.ScreenDetails
		return m, nil
	}

	// Ejecutar cambio de puerto
	if _, err := m.ServiceManager.ChangePort(svc.Name, port); err != nil {
		m.ErrMsg = err.Error()
	}

	// Refrescar servicios
	m.RefreshServices()

	// Volver al dashboard
	m.Screen = models.ScreenDashboard
	m.StartPolling()

	return m, nil
}

// filterServicesByStack filtra los servicios para mostrar solo los del stack seleccionado.
func filterServicesByStack(services []models.Service, stackType models.StackType) []models.Service {
	var filtered []models.Service

	// Mapear stacks a servicios que contienen
	stackServices := map[models.StackType][]string{
		models.StackTypeLAMP: {"apache2", "mysql"},
		models.StackTypeLAMM: {"apache2", "mariadb"},
		models.StackTypeLEPP: {"nginx", "postgresql"},
	}

	servicesToShow, ok := stackServices[stackType]
	if !ok {
		return services
	}

	// Filtrar solo los servicios del stack
	for _, svc := range services {
		for _, name := range servicesToShow {
			if svc.Name == name {
				filtered = append(filtered, svc)
				break
			}
		}
	}

	return filtered
}