package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// WelcomeView renderiza la pantalla de bienvenida.
// hasServices indica si ya hay servicios instalados.
func WelcomeView(selectedIndex int, width int, hasServices bool) string {
	styles := NewStyles()
	t := styles.Theme

	// SIEMPRE mostrar ambas opciones para que el usuario elija
	options := []string{
		"Install Stack",
		"Manage Stack",
		"Quit",
	}

	// Construir contenido con estilos del tema
	menu := renderMenu(options, selectedIndex, t)

	content := fmt.Sprintf(`  Welcome to XAMPP-TUI

  Manage your local web server stack (Apache, MySQL, PostgreSQL, Nginx)
  with a simple and intuitive interface.

%s`,
		menu,
	)

	hints := []string{"↑↓", "Navigate", "Enter", "Select", "q", "Quit"}
	return BaseLayout(styles, width, content, hints)
}

// renderMenu renderiza un menú de opciones.
func renderMenu(options []string, selected int, t *Theme) string {
	var lines []string

	// El símbolo ▶ tiene ancho 1
	symbolWidth := 1

	for i, opt := range options {
		if i == selected {
			// Seleccionado: símbolo verde + espacio + texto
			prefix := lipgloss.NewStyle().Foreground(t.Green).Bold(true).Render("▶")
			lines = append(lines, prefix+" "+lipgloss.NewStyle().Foreground(t.Green).Bold(true).Render(opt))
		} else {
			// No seleccionado: espacios en blanco (ancho del símbolo) + espacio + texto
			spaces := strings.Repeat(" ", symbolWidth)
			lines = append(lines, spaces+" "+lipgloss.NewStyle().Foreground(t.Text).Render(opt))
		}
	}

	return strings.Join(lines, "\n")
}

// ErrorView renderiza una pantalla de error.
func ErrorView(msg string) string {
	return ErrorViewWidth(msg, 80)
}

// ErrorViewWidth renderiza una pantalla de error con ancho específico.
func ErrorViewWidth(msg string, width int) string {
	styles := NewStyles()

	content := fmt.Sprintf(`  Error: %s

  Press any key to continue
`, msg)

	hints := []string{"any", "Continue"}
	return BaseLayout(styles, width, content, hints)
}