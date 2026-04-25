package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"xampp-tui/internal/models"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg == nil {
		t.Fatal("Default() returned nil")
	}

	if cfg.StackType != models.StackTypeLAMP {
		t.Errorf("StackType = %v, want %v", cfg.StackType, models.StackTypeLAMP)
	}

	if cfg.Theme != "dark" {
		t.Errorf("Theme = %q, want %q", cfg.Theme, "dark")
	}

	if cfg.Port != 8080 {
		t.Errorf("Port = %d, want %d", cfg.Port, 8080)
	}

	if cfg.AutoStart != false {
		t.Errorf("AutoStart = %v, want %v", cfg.AutoStart, false)
	}
}

func TestLoadNotExists(t *testing.T) {
	// Use a non-existent path
	cfg, err := Load("/nonexistent/path/config.json")

	if err != nil {
		t.Errorf("Load() error = %v", err)
	}

	if cfg == nil {
		t.Error("Load() should return Default config for non-existent file")
	}
}

func TestLoadExists(t *testing.T) {
	// Create temp config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Write test config
	testCfg := &Config{
		StackType:  models.StackTypeLEPP,
		Theme:     "light",
		Port:      9090,
		LogPath:   "/tmp/logs",
		DataPath:  "/tmp/data",
		AutoStart: true,
	}

	data, err := json.Marshal(testCfg)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	// Load it
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.StackType != models.StackTypeLEPP {
		t.Errorf("StackType = %v, want %v", loaded.StackType, models.StackTypeLEPP)
	}

	if loaded.Theme != "light" {
		t.Errorf("Theme = %q, want %q", loaded.Theme, "light")
	}

	if loaded.Port != 9090 {
		t.Errorf("Port = %d, want %d", loaded.Port, 9090)
	}

	if loaded.AutoStart != true {
		t.Errorf("AutoStart = %v, want %v", loaded.AutoStart, true)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Write invalid JSON
	if err := os.WriteFile(configPath, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	// Load should fail
	_, err := Load(configPath)
	if err == nil {
		t.Error("Load() should return error for invalid JSON")
	}
}

func TestSave(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	cfg := &Config{
		StackType:  models.StackTypeLEPP,
		Theme:     "dark",
		Port:      8888,
		LogPath:   "/var/log/xampp",
		DataPath:  "/var/lib/xampp",
		AutoStart: true,
	}

	if err := Save(cfg, configPath); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("config file not created: %v", err)
	}

	// Verify content
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.StackType != cfg.StackType {
		t.Errorf("StackType = %v, want %v", loaded.StackType, cfg.StackType)
	}

	if loaded.Port != cfg.Port {
		t.Errorf("Port = %d, want %d", loaded.Port, cfg.Port)
	}
}

func TestSaveCreatesDir(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "subdir", "config.json")

	cfg := Default()

	if err := Save(cfg, configPath); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("config file not created: %v", err)
	}
}

func TestGetConfigPath(t *testing.T) {
	path := GetConfigPath()

	if path == "" {
		t.Error("GetConfigPath() returned empty string")
	}

	// Should contain .config/xampp-tui/config.json
	if filepath.Base(path) != "config.json" {
		t.Errorf("GetConfigPath() = %q, want file named config.json", path)
	}
}

func TestConfigRoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	original := &Config{
		StackType:  models.StackTypeLAMP,
		Theme:     "dark",
		Port:      8080,
		LogPath:   "/home/user/.config/xampp-tui/logs",
		DataPath:  "/home/user/.local/share/xampp-tui",
		AutoStart: false,
	}

	// Save
	if err := Save(original, configPath); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Load
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify
	if loaded.StackType != original.StackType {
		t.Errorf("StackType = %v, want %v", loaded.StackType, original.StackType)
	}

	if loaded.Theme != original.Theme {
		t.Errorf("Theme = %q, want %q", loaded.Theme, original.Theme)
	}

	if loaded.Port != original.Port {
		t.Errorf("Port = %d, want %d", loaded.Port, original.Port)
	}
}

func TestConfigDefaultPorts(t *testing.T) {
	cfg := Default()

	// Verify port is set
	if cfg.Port == 0 {
		t.Error("Default port should not be 0")
	}

	if cfg.Port < 1024 || cfg.Port > 65535 {
		t.Errorf("Port = %d, should be in valid range (1024-65535)", cfg.Port)
	}
}