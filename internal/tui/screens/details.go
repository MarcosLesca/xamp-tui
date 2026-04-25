package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"xampp-tui/internal/models"
)

// DetailsView renderiza la pantalla de detalles de un servicio.
// selectedMenuIndex es el índice del menú seleccionado.
func DetailsView(svc *models.Service, width int, selectedMenuIndex int) string {
	styles := NewStyles()
	t := styles.Theme

	if svc == nil {
		return ErrorViewWidth("No service selected", width)
	}

	// Opciones de menú según el estado
	var menuOptions []string
	running := svc.Status == "running"
	available := svc.Status == "available"

	if running {
		menuOptions = []string{
			"Stop",
			"Restart",
			"View Logs",
			"Change Port",
			"Back",
		}
	} else if available {
		menuOptions = []string{
			"Install",
			"Back",
		}
	} else {
		menuOptions = []string{
			"Start",
			"View Logs",
			"Change Port",
			"Back",
		}
	}

	// Información del servicio
	info := renderServiceInfo(svc, t)

	// Menú con selección
	menu := renderMenuOptions(menuOptions, selectedMenuIndex, t)

	content := fmt.Sprintf(`  Service Details - %s

%s

%s`,
		strings.ToUpper(svc.Name),
		info,
		menu,
	)

	hints := []string{"↑↓", "Navigate", "Enter", "Execute", "Esc", "Back"}
	return BaseLayout(styles, width, content, hints)
}

// renderServiceInfo renderiza la información del servicio.
func renderServiceInfo(svc *models.Service, t *Theme) string {
	// Estado con color
	statusColor := t.Green
	if svc.Status == "stopped" {
		statusColor = t.Red
	} else if svc.Status == "error" {
		statusColor = t.Yellow
	}
	statusBadge := lipgloss.NewStyle().
		Foreground(statusColor).
		Background(t.Surface).
		Padding(0, 1).
		Render(fmt.Sprintf(" %s ", strings.ToUpper(svc.Status)))

	// Versión
	version := svc.Version
	if version == "" {
		version = "Not installed"
	}

	// Puerto
	port := fmt.Sprintf("%d", svc.Port)

	// PID
	pid := "N/A"
	if svc.PID > 0 {
		pid = fmt.Sprintf("%d", svc.PID)
	}

	// Uptime
	uptime := formatUptime(svc.Uptime)

	textStyle := lipgloss.NewStyle().Foreground(t.Text)

	info := fmt.Sprintf(`%s
%s
%s
%s
%s`,
		statusBadge,
		textStyle.Render("Version: "+version),
		textStyle.Render("Port: "+port),
		textStyle.Render("PID: "+pid),
		textStyle.Render("Uptime: "+uptime),
	)

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(t.Overlay).
		Padding(0, 1).
		Render(info)
}