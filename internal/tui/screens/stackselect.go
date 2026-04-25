package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// StackSelectView renderiza la pantalla de selección de stack.
func StackSelectView(selectedIndex int, width int) string {
	styles := NewStyles()
	t := styles.Theme

	// Opciones de stack
	options := []string{
		"LAMP Stack",
		"LAMM Stack",
		"LEPP Stack",
		"Back",
	}

	// Descripciones
	descriptions := []string{
		"Apache + MySQL + PHP",
		"Apache + MariaDB + phpMyAdmin",
		"Nginx + PostgreSQL + PHP",
		"Volver al menú principal",
	}

	menu := renderStackMenu(options, descriptions, selectedIndex, t)

	content := fmt.Sprintf(`  Choose Your Stack

  Select which web server stack you want to manage:

%s`,
		menu,
	)

	hints := []string{"↑↓", "Navigate", "Enter", "Select", "Esc", "Back"}
	return BaseLayout(styles, width, content, hints)
}

// renderStackMenu renderiza las opciones de stack con su descripción.
func renderStackMenu(options []string, descriptions []string, selected int, t *Theme) string {
	var lines []string

	// El símbolo ▶ tiene ancho 1
	symbolWidth := 1

	for i, opt := range options {
		if i == selected {
			prefix := lipgloss.NewStyle().Foreground(t.Green).Bold(true).Render("▶")
			lines = append(lines, prefix+" "+lipgloss.NewStyle().Foreground(t.Green).Bold(true).Render(opt))
			lines = append(lines, "  "+lipgloss.NewStyle().Foreground(t.Subtext).Render(descriptions[i]))
		} else {
			spaces := strings.Repeat(" ", symbolWidth)
			lines = append(lines, spaces+" "+lipgloss.NewStyle().Foreground(t.Text).Render(opt))
			lines = append(lines, "  "+lipgloss.NewStyle().Foreground(t.Subtext).Render(descriptions[i]))
		}
		if i < len(options)-1 {
			lines = append(lines, "")
		}
	}

	return strings.Join(lines, "\n")
}