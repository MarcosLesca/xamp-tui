package screens

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"xampp-tui/internal/models"
)

// DashboardView renderiza el panel principal con el estado de los servicios.
// selectedIndex es el índice del servicio seleccionado.
func DashboardView(services []models.Service, width int, installing bool, installLog string, selectedIndex int) string {
	styles := NewStyles()
	t := styles.Theme

	// Mostrar banner de instalación si está instalando
	installBanner := ""
	if installing {
		installBanner = renderInstallBanner(t)
	}

	// Siempre mostrar el log si hay
	if installLog != "" {
		// Si está instalando, mostrar el estado
		if installing {
			installBanner += lipgloss.NewStyle().
				Foreground(t.Yellow).
				Render(installLog) + "\n\n"
		} else {
			// Si terminó, mostrar el resultado
			installBanner += renderInstallLog(installLog, t)
		}
	}

	if len(services) == 0 && !installing {
		return DashboardEmptyView(width)
	}

	// Construir grid de servicios
	grid := renderServiceGrid(services, t, selectedIndex)

	// Resumen de estado
	summary := renderStatusSummary(services, t)

	content := fmt.Sprintf(`  Dashboard - Service Status

%s

%s

%s`,
		installBanner,
		summary,
		grid,
	)

	hints := []string{"↑↓", "Select", "Enter", "Details", "Esc", "Quit"}
	return BaseLayout(styles, width, content, hints)
}

// renderInstallBanner renderiza el banner de instalación en progreso.
func renderInstallBanner(t *Theme) string {
	style := lipgloss.NewStyle().
		Background(t.Yellow).
		Foreground(t.Base).
		Padding(0, 1).
		Bold(true)

	return style.Render(" ● INSTALLING... ") + " Esperá que pida tu password\n\n"
}

// renderInstallLog renderiza el log de instalación al final.
func renderInstallLog(log string, t *Theme) string {
	lines := strings.Split(log, "\n")
	var formatted []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			formatted = append(formatted, "  "+lipgloss.NewStyle().Foreground(t.Subtext).Render(line))
		}
	}
	return strings.Join(formatted, "\n") + "\n\n"
}

// DashboardEmptyView renderiza el dashboard cuando no hay servicios.
func DashboardEmptyView(width int) string {
	styles := NewStyles()

	content := `
  Dashboard

  No services detected. Please install a stack first.

  Go to StackSelect to choose and install your stack.
`

	hints := []string{"Esc", "Back"}
	return BaseLayout(styles, width, content, hints)
}

// renderServiceGrid renderiza el grid de tarjetas de servicios.
// selectedIndex es el índice del servicio seleccionado.
func renderServiceGrid(services []models.Service, t *Theme, selectedIndex int) string {
	var cards []string

	// Crear tarjetas una por línea
	for i := 0; i < len(services); i++ {
		card := renderServiceCard(services[i], i, 38, t, i == selectedIndex)
		cards = append(cards, card)
	}

	// Unir con newline
	return strings.Join(cards, "\n")
}

// renderServiceCard renderiza una tarjeta de servicio.
// selected indica si está seleccionada.
func renderServiceCard(svc models.Service, index int, width int, t *Theme, selected bool) string {
	// Borde según selección
	borderColor := t.Overlay
	if selected {
		borderColor = t.Yellow
	}

	// Nombre del servicio
	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(t.Mauve)
	nameContent := nameStyle.Render(strings.ToUpper(svc.Name))

	// Estado con color
	statusColor := t.Green
	if svc.Status == "stopped" {
		statusColor = t.Red
	} else if svc.Status == "error" {
		statusColor = t.Yellow
	} else if svc.Status == "available" {
		statusColor = t.Sky // Azul para "disponible para instalar"
	}
	statusStyle := lipgloss.NewStyle().Foreground(statusColor).Bold(true)
	statusContent := statusStyle.Render(strings.ToUpper(svc.Status))

	// Versión
	version := svc.Version
	if version == "" {
		version = "N/A"
	}

	// Puerto
	port := fmt.Sprintf("%d", svc.Port)

	// Uptime formateado
	uptime := formatUptime(svc.Uptime)

	// Construir tarjeta
	cardStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(width).
		Padding(0, 1)

	textStyle := lipgloss.NewStyle().Foreground(t.Text)

	cardContent := fmt.Sprintf("%s\n%s\n%s\nPort: %s\nUptime: %s",
		nameContent,
		statusContent,
		textStyle.Render("Version: "+version),
		textStyle.Render(port),
		textStyle.Render(uptime),
	)

	return cardStyle.Render(cardContent)
}

// renderStatusSummary renderiza el resumen de estado de servicios.
func renderStatusSummary(services []models.Service, t *Theme) string {
	running := 0
	stopped := 0
	available := 0
	failed := 0

	for _, svc := range services {
		switch svc.Status {
		case "running":
			running++
		case "stopped":
			stopped++
		case "available":
			available++
		case "error":
			failed++
		}
	}

	runningBadge := lipgloss.NewStyle().Foreground(t.Green).Render(fmt.Sprintf("● %d Running", running))
	stoppedBadge := lipgloss.NewStyle().Foreground(t.Red).Render(fmt.Sprintf("○ %d Stopped", stopped))
	availableBadge := lipgloss.NewStyle().Foreground(t.Sky).Render(fmt.Sprintf("○ %d Available", available))

	summary := fmt.Sprintf("  %s  %s  %s", runningBadge, stoppedBadge, availableBadge)
	if failed > 0 {
		summary += lipgloss.NewStyle().Foreground(t.Yellow).Render(fmt.Sprintf(" ⚠ %d Failed", failed))
	}

	return summary
}

// formatUptime formatea el uptime.
func formatUptime(d time.Duration) string {
	if d == 0 {
		return "N/A"
	}

	seconds := d.Seconds()
	if seconds < 60 {
		return fmt.Sprintf("%.0fs", seconds)
	}

	minutes := seconds / 60
	if minutes < 60 {
		return fmt.Sprintf("%.0fm", minutes)
	}

	hours := minutes / 60
	if hours < 24 {
		return fmt.Sprintf("%.1fh", hours)
	}

	days := hours / 24
	return fmt.Sprintf("%.1fd", days)
}