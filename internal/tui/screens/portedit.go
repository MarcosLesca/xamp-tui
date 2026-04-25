package screens

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"xampp-tui/internal/models"
)

// PortEditView renderiza la pantalla de edición de puerto.
func PortEditView(svc *models.Service, width int, portValue string) string {
	styles := NewStyles()
	t := styles.Theme

	if svc == nil {
		return ErrorViewWidth("No service selected", width)
	}

	defaultPort := strconv.Itoa(svc.Port)
	if portValue == "" {
		portValue = defaultPort
	}

	// Preview del valor actual
	preview := portValue
	if portValue == "" {
		preview = "_"
	}

	previewStyle := lipgloss.NewStyle().
		Foreground(t.Yellow).
		Width(width - 8)

	portStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(t.Blue).
		Foreground(t.Text).
		Width(width - 8).
		Padding(0, 1)

	portLine := fmt.Sprintf("Puerto: %s", previewStyle.Render(preview))

	textStyle := lipgloss.NewStyle().Foreground(t.Subtext)

	content := fmt.Sprintf(`  Change Port - %s

%s

  Puerto actual: %s

%s

  %s`,
		strings.ToUpper(svc.Name),
		portStyle.Render(portLine),
		defaultPort,
		textStyle.Render("↑↓ Editar valor   Enter: Confirmar   Esc: Cancelar"),
		lipgloss.NewStyle().Foreground(t.Red).Render(" Puerto inválido (1-65535)"),
	)

	hints := []string{"↑↓/0-9", "Edit Port", "Enter", "Confirm", "Esc", "Cancel"}
	return BaseLayout(styles, width, content, hints)
}

// PortEditResultView muestra el resultado del cambio de puerto.
func PortEditResultView(svc *models.Service, width int, portValue string, success bool, message string) string {
	styles := NewStyles()
	t := styles.Theme

	if svc == nil {
		return ErrorViewWidth("No service selected", width)
	}

	if success {
		msg := lipgloss.NewStyle().
			Foreground(t.Green).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(t.Green).
			Width(width - 8).
			Padding(0, 1).
			Render(fmt.Sprintf(" Puerto cambiado a %s ", portValue))
		content := fmt.Sprintf(`  Port Changed - %s

%s

  %s

  Presiona cualquier tecla para volver`,
			strings.ToUpper(svc.Name),
			msg,
			lipgloss.NewStyle().Foreground(t.Subtext).Render(message),
		)
		return BaseLayout(styles, width, content, []string{"Enter", "Back"})
	}

	msg := lipgloss.NewStyle().
		Foreground(t.Red).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(t.Red).
		Width(width - 8).
		Padding(0, 1).
		Render(fmt.Sprintf(" Error: %s ", message))
	content := fmt.Sprintf(`  Port Error - %s

%s

  %s`,
		strings.ToUpper(svc.Name),
		msg,
		lipgloss.NewStyle().Foreground(t.Subtext).Render("Presiona cualquier tecla para volver"),
	)
	return BaseLayout(styles, width, content, []string{"Enter", "Back"})
}