package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"xampp-tui/internal/models"
)

// Config representa la configuración de la aplicación.
type Config struct {
	StackType  models.StackType    `json:"stack_type"`  // Tipo de stack: LAMP o LEPP
	Theme     string             `json:"theme"`      // Tema visual
	Port      int                `json:"port"`       // Puerto para servidor integrado
	LogPath   string            `json:"log_path"`  // Ruta de logs
	DataPath  string            `json:"data_path"` // Ruta de datos
	AutoStart bool               `json:"auto_start"` // Iniciar servicios automáticamente
}

// Default devuelve la configuración por defecto.
func Default() *Config {
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".config", "xampp-tui")
	dataDir := filepath.Join(home, ".local", "share", "xampp-tui")

	return &Config{
		StackType:  models.StackTypeLAMP,
		Theme:     "dark",
		Port:      8080,
		LogPath:   filepath.Join(configDir, "logs"),
		DataPath:  dataDir,
		AutoStart: false,
	}
}

// Load carga la configuración desde el archivo.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		// Si no existe, devolver config por defecto
		if os.IsNotExist(err) {
			return Default(), nil
		}
		return nil, err
	}

	cfg := &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Save guarda la configuración en el archivo.
func Save(cfg *Config, path string) error {
	// Crear directorio si no existe
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// GetConfigPath devuelve la ruta del archivo de configuración.
func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "xampp-tui", "config.json")
}