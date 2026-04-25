package screens

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"xampp-tui/internal/models"
)

// LogsView renderiza la pantalla de logs de un servicio.
// selectedMenuIndex es el índice del menú seleccionado.
func LogsView(svc *models.Service, width int, selectedMenuIndex int) string {
	styles := NewStyles()
	t := styles.Theme

	if svc == nil {
		return ErrorViewWidth("No service selected", width)
	}

	// Obtener logs del servicio
	logs := getServiceLogs(svc.Name)

	// Opciones de servicio
	serviceOptions := []string{
		"Apache/Nginx Logs",
		"Database Logs",
		"Back",
	}

	menu := renderMenuOptions(serviceOptions, selectedMenuIndex, t)

	content := fmt.Sprintf(`  Logs - %s

%s

%s`,
		strings.ToUpper(svc.Name),
		lipgloss.NewStyle().Foreground(t.Subtext).Render(logs),
		menu,
	)

	hints := []string{"↑↓", "Navigate", "Esc", "Back"}
	return BaseLayout(styles, width, content, hints)
}

// getServiceLogs obtiene los logs de un servicio.
func getServiceLogs(serviceName string) string {
	var logPaths = map[string][]string{
		"apache2":    {"/var/log/apache2/error.log", "/var/log/apache2/access.log"},
		"nginx":      {"/var/log/nginx/error.log", "/var/log/nginx/access.log"},
		"mysql":      {"/var/log/mysql/error.log", "/var/log/mysql/general.log"},
		"mariadb":    {"/var/log/mysql/error.log", "/var/log/mysql/general.log"},
		"postgresql": {"/var/log/postgresql/postgresql-*.log"},
	}

	paths, ok := logPaths[serviceName]
	if !ok {
		return "No logs available for this service"
	}

	var logs []string
	for _, path := range paths {
		if strings.Contains(path, "*") {
			logs = append(logs, fmt.Sprintf("  %s", path))
		} else {
			if _, err := os.Stat(path); err == nil {
				content, err := readLastLines(path, 20)
				if err == nil {
					logs = append(logs, fmt.Sprintf("  %s:\n%s", path, content))
				}
			}
		}
	}

	if len(logs) == 0 {
		return "No logs found"
	}

	return strings.Join(logs, "\n\n")
}

// readLastLines lee las últimas n líneas de un archivo.
func readLastLines(path string, n int) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}

	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	result := make([]string, len(lines))
	for i, line := range lines {
		timestamp := time.Now().Format("15:04:05")
		result[i] = fmt.Sprintf("    %s | %s", timestamp, line)
	}

	return strings.Join(result, "\n"), nil
}

// renderMenuOptions renderiza las opciones del menú.
func renderMenuOptions(options []string, selected int, t *Theme) string {
	var lines []string

	for i, opt := range options {
		// El símbolo ▶ tiene ancho 1
		symbolWidth := 1

		if i == selected {
			prefix := lipgloss.NewStyle().Foreground(t.Green).Bold(true).Render("▶")
			lines = append(lines, prefix+" "+lipgloss.NewStyle().Foreground(t.Green).Bold(true).Render(opt))
		} else {
			spaces := strings.Repeat(" ", symbolWidth)
			lines = append(lines, spaces+" "+lipgloss.NewStyle().Foreground(t.Text).Render(opt))
		}
	}

	return strings.Join(lines, "\n")
}