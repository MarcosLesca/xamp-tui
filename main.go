package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"xampp-tui/internal/config"
	"xampp-tui/internal/tui"
)

func main() {
	// Cargar configuración
	cfg, err := config.Load(config.GetConfigPath())
	if err != nil {
		log.Printf("Warning: could not load config: %v", err)
		cfg = config.Default()
	}

	// Crear modelo inicial
	model := tui.New()
	model.Config = cfg

	// Inicializar servicios
	if err := model.InitServices(); err != nil {
		log.Printf("Warning: could not detect services: %v", err)
	}

	// Crear programa bubbletea
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	)

	// Ejecutar
	if err := p.Start(); err != nil {
		log.Fatalf("Error running TUI: %v", err)
	}

	// Guardar configuración al salir
	if err := config.Save(cfg, config.GetConfigPath()); err != nil {
		log.Printf("Warning: could not save config on exit: %v", err)
	}

	// Asegurar que el directorio de datos existe
	if err := ensureDataDir(cfg.DataPath); err != nil {
		log.Printf("Warning: could not create data dir: %v", err)
	}
}

// ensureDataDir asegura que el directorio de datos existe.
func ensureDataDir(path string) error {
	return os.MkdirAll(path, 0755)
}