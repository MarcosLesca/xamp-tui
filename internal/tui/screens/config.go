package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"xampp-tui/internal/models"
)

// ConfigView renderiza la pantalla de configuración antes de instalar.
func ConfigView(stackType models.StackType, width int, selectedIndex int, config *models.LAMMConfig, editing bool) string {
	styles := NewStyles()
	t := styles.Theme

	if config == nil {
		c := models.DefaultLAMMConfig()
		config = &c
	}

	// Opciones de configuración segun stack
	var configItems []string

	dbName := "MariaDB"
	if stackType == models.StackTypeLAMP {
		dbName = "MySQL"
	}

	configItems = []string{
		fmt.Sprintf("%s Root Password: %s", dbName, renderPasswordField(config.MariaDBRootPassword)),
		fmt.Sprintf("Remove Anonymous Users: %s", boolToYesNo(config.RemoveAnonymous)),
		fmt.Sprintf("Disallow Root Remote Login: %s", boolToYesNo(config.DisallowRemote)),
		fmt.Sprintf("Remove Test Database: %s", boolToYesNo(config.RemoveTestDB)),
		fmt.Sprintf("Web Server: %s", config.WebServer),
		fmt.Sprintf("phpMyAdmin Path: %s", config.PHPMyAdminPath),
	}

	// Menu actions
	menuOptions := []string{
		"Start Installation",
		"Edit Options",
		"Back",
	}

	// Render config items
	var configLines []string
	for i, item := range configItems {
		prefix := "  "
		if editing && i == selectedIndex {
			// Highlight selected config option when editing
			prefix = lipgloss.NewStyle().Foreground(t.Yellow).Render("▶")
		}
		configLines = append(configLines, prefix+" "+lipgloss.NewStyle().Foreground(t.Text).Render(item))
	}
	configSection := strings.Join(configLines, "\n")

	// Render menu (with selection) - only when not editing
	menu := ""
	if !editing {
		menu = renderMenu(menuOptions, selectedIndex, t)
	} else {
		// When editing, show instructions
		menu = lipgloss.NewStyle().Foreground(t.Yellow).Render("▶ Press SPACE to toggle, UP/DOWN to navigate, ESC to exit edit mode")
	}

	// Header
	header := "Review and edit configuration before installation:"
	if editing {
		header = "EDIT MODE: Press SPACE to toggle values"
	}

	content := fmt.Sprintf(`  Configuration - %s

%s

%s

%s`,
		stackTypeToString(stackType),
		lipgloss.NewStyle().Foreground(t.Subtext).Render(header),
		configSection,
		menu,
	)

	hint1, hint2 := "Enter", "Select"
	if editing {
		hint1, hint2 = "Space", "Toggle"
	}
	
	hints := []string{"↑↓", "Navigate", hint1, hint2, "Esc", "Back"}
	if editing {
		hints = []string{"↑↓", "Navigate", "Space", "Toggle", "Esc", "Exit"}
	}
	
	return BaseLayout(styles, width, content, hints)
}

// renderPasswordField muestra asteriscos si hay password.
func renderPasswordField(pwd string) string {
	if pwd == "" {
		return "(not set - will prompt during install)"
	}
	return strings.Repeat("●", len(pwd))
}

// boolToYesNo convierte booleano a string.
func boolToYesNo(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}