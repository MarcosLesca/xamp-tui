package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"xampp-tui/internal/models"
)

// InstallView renderiza la pantalla de instalación.
// selectedIndex es la opción seleccionada del menú.
// canInstall indica si tiene permisos para instalar.
func InstallView(stackType models.StackType, width int, selectedIndex int, canInstall bool) string {
	styles := NewStyles()
	t := styles.Theme

	// Instrucciones según tenga permisos
	var instructions string
	instructStyle := lipgloss.NewStyle().Foreground(t.Subtext)
	if canInstall {
		instructions = "Select 'Start Installation' to begin"
	} else {
		instructions = "Run in terminal: sudo ./xampp-tui\nOr run commands manually, then reopen."
		instructStyle = lipgloss.NewStyle().Foreground(t.Yellow).Bold(true)
	}

	// Determinar los comandos según el stack
	var installCmds []string
	var installDesc []string

	switch stackType {
	case models.StackTypeLAMP:
		installCmds = []string{
			"sudo apt-get update",
			"sudo apt-get install -y apache2 mysql-server php libapache2-mod-php php-mysql",
			"sudo systemctl enable apache2 mysql",
			"sudo systemctl start apache2 mysql",
			"sudo mysql_secure_installation",
		}
		installDesc = []string{
			"Update package lists",
			"Install Apache, MySQL Server, PHP",
			"Enable services on boot",
			"Start Apache & MySQL services",
			"Secure MySQL installation (set root password)",
		}

	case models.StackTypeLAMM:
		installCmds = []string{
			"sudo apt-get update",
			"sudo apt-get install -y apache2 mariadb-server php libapache2-mod-php php-mysql phpmyadmin",
			"sudo systemctl enable apache2 mariadb",
			"sudo systemctl start apache2 mariadb",
			"sudo mysql_secure_installation",
			"sudo cp -rn /usr/share/phpmyadmin /var/www/html/phpmyadmin 2>/dev/null || sudo ln -sf /usr/share/phpmyadmin /var/www/html/phpmyadmin",
		}
		installDesc = []string{
			"Update package lists",
			"Install Apache, MariaDB Server, PHP, phpMyAdmin",
			"Enable services on boot",
			"Start Apache & MariaDB services",
			"Secure MariaDB installation (set root password)",
			"Configure phpMyAdmin web access",
		}

	case models.StackTypeLEPP:
		installCmds = []string{
			"sudo apt-get update",
			"sudo apt-get install -y nginx postgresql php-fpm php-pgsql",
			"sudo systemctl enable nginx postgresql",
			"sudo systemctl start nginx postgresql",
		}
		installDesc = []string{
			"Update package lists",
			"Install Nginx, PostgreSQL, PHP",
			"Enable services on boot",
			"Start services",
		}
	}

	// Opciones del menú
	menuOptions := []string{
		"Start Installation",
		"Back",
	}

	commands := renderInstallCommands(installCmds, installDesc, t, selectedIndex == 0)
	menu := renderMenuOptions(menuOptions, selectedIndex, t)

	content := fmt.Sprintf(`  Installation

  Stack: %s

%s

%s

  %s

  Use arrows to select, Enter to confirm
`,
		stackTypeToString(stackType),
		commands,
		menu,
		instructStyle.Render(instructions),
	)

	hints := []string{"↑↓", "Navigate", "Enter", "Select", "Esc", "Back"}
	return BaseLayout(styles, width, content, hints)
}

// stackTypeToString convierte el stack a string
func stackTypeToString(stackType models.StackType) string {
	switch stackType {
	case models.StackTypeLAMP:
		return "LAMP (Apache + MySQL + PHP)"
	case models.StackTypeLAMM:
		return "LAMM (Apache + MariaDB + phpMyAdmin)"
	case models.StackTypeLEPP:
		return "LEPP (Nginx + PostgreSQL + PHP)"
	default:
		return "Unknown"
	}
}

// renderInstallCommands renderiza los comandos de instalación.
func renderInstallCommands(cmds []string, descs []string, t *Theme, isFirstSelected bool) string {
	var lines []string

	lines = append(lines, lipgloss.NewStyle().Foreground(t.Subtext).Render("Commands:"))

	// El símbolo ▶ tiene ancho 1
	symbolWidth := 1

	for i, cmd := range cmds {
		if isFirstSelected && i == 0 {
			prefix := lipgloss.NewStyle().Foreground(t.Green).Render("▶")
			lines = append(lines, prefix+" "+lipgloss.NewStyle().Foreground(t.Text).Render(cmd))
		} else {
			spaces := strings.Repeat(" ", symbolWidth)
			lines = append(lines, spaces+" "+lipgloss.NewStyle().Foreground(t.Text).Render(cmd))
		}
		lines = append(lines, "  "+lipgloss.NewStyle().Foreground(t.Subtext).Render(descs[i]))
	}

	return strings.Join(lines, "\n")
}
