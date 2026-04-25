package screens

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"xampp-tui/internal/models"
)

// DatabaseView renderiza la pantalla de gestión de base de datos.
func DatabaseView(stackType models.StackType, width int) string {
	styles := NewStyles()
	t := styles.Theme

	// Opciones de menú
	menuOptions := []string{
		"List Databases",
		"Create Database",
		"Drop Database",
		"Back",
	}

	menu := renderMenuOptions(menuOptions, 0, t)

	// Obtener información de la DB
	dbInfo := renderDBInfo(stackType, t)

	content := fmt.Sprintf(`  Database Management

%s

%s`,
		dbInfo,
		menu,
	)

	hints := []string{"↑↓", "Navigate", "Enter", "Execute", "Esc", "Back"}
	return BaseLayout(styles, width, content, hints)
}

// renderDBInfo renderiza la información de la base de datos.
func renderDBInfo(stackType models.StackType, t *Theme) string {
	var dbType string
	var dbPort int

	switch stackType {
	case models.StackTypeLAMP:
		dbType = "MySQL/MariaDB"
		dbPort = 3306
	case models.StackTypeLEPP:
		dbType = "PostgreSQL"
		dbPort = 5432
	default:
		dbType = "None"
		dbPort = 0
	}

	textStyle := lipgloss.NewStyle().Foreground(t.Text)

	info := fmt.Sprintf(`%s
%s`,
		textStyle.Render("Database: "+dbType),
		textStyle.Render("Port: "+fmt.Sprintf("%d", dbPort)),
	)

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(t.Overlay).
		Padding(0, 1).
		Render(info)
}

// ListDatabasesCommand devuelve el comando para listar bases de datos.
func ListDatabasesCommand(stackType models.StackType) []string {
	switch stackType {
	case models.StackTypeLAMP:
		return []string{"mysql", "-u", "root", "-p", "-e", "SHOW DATABASES;"}
	case models.StackTypeLEPP:
		return []string{"psql", "-U", "postgres", "-c", "\\l"}
	default:
		return nil
	}
}

// CreateDatabaseCommand devuelve el comando para crear una base de datos.
func CreateDatabaseCommand(stackType models.StackType, dbName string) []string {
	switch stackType {
	case models.StackTypeLAMP:
		return []string{"mysql", "-u", "root", "-p", "-e", "CREATE DATABASE " + dbName + ";"}
	case models.StackTypeLEPP:
		return []string{"psql", "-U", "postgres", "-c", "CREATE DATABASE " + dbName + ";"}
	default:
		return nil
	}
}

// DropDatabaseCommand devuelve el comando para eliminar una base de datos.
func DropDatabaseCommand(stackType models.StackType, dbName string) []string {
	switch stackType {
	case models.StackTypeLAMP:
		return []string{"mysql", "-u", "root", "-p", "-e", "DROP DATABASE IF EXISTS " + dbName + ";"}
	case models.StackTypeLEPP:
		return []string{"psql", "-U", "postgres", "-c", "DROP DATABASE IF EXISTS " + dbName + ";"}
	default:
		return nil
	}
}