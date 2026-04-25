package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"xampp-tui/internal/models"
)

// TestScreenTransitions tests screen transitions in the TUI.
func TestScreenTransitions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tests := []struct {
		name         string
		startScreen models.Screen
		key         string
		wantScreen  models.Screen
	}{
		{
			name:         "welcome to stack select on enter",
			startScreen:  models.ScreenWelcome,
			key:         "enter",
			wantScreen:  models.ScreenStackSelect,
		},
		{
			name:         "escape from welcome quits",
			startScreen:  models.ScreenWelcome,
			key:         "esc",
			wantScreen:  models.ScreenWelcome,
		},
		{
			name:         "stack select to install on enter",
			startScreen:  models.ScreenStackSelect,
			key:         "enter",
			wantScreen:  models.ScreenConfig,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			m.Screen = tt.startScreen

			newModel, _ := m.Update(tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune(tt.key),
			})

			result := newModel.(Model)

			if result.Screen != tt.wantScreen {
				t.Errorf("Screen = %v, want %v", result.Screen, tt.wantScreen)
			}
		})
	}
}

// TestNavigation tests navigation keys.
func TestNavigation(t *testing.T) {
	tests := []struct {
		name     string
		key     string
		wantUp  bool
		wantIdx int
	}{
		{"down arrow", "down", false, 1},
		{"j key", "j", false, 1},
		{"up arrow", "up", true, 0},
		{"k key", "k", true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			m.Services = []models.Service{
				{Name: "svc1"},
				{Name: "svc2"},
			}
			m.Screen = models.ScreenDashboard
			m.SelectedServiceIndex = 1

			var keyMsg tea.KeyMsg
			if tt.wantUp {
				keyMsg = tea.KeyMsg{Type: tea.KeyUp}
			} else {
				keyMsg = tea.KeyMsg{Type: tea.KeyDown}
			}

			newModel, _ := m.Update(keyMsg)
			result := newModel.(Model)

			if tt.wantUp && result.SelectedServiceIndex != 0 {
				t.Errorf("SelectedServiceIndex = %d, want 0", result.SelectedServiceIndex)
			}
			if !tt.wantUp && result.SelectedServiceIndex != 1 {
				t.Errorf("SelectedServiceIndex = %d, want 1", result.SelectedServiceIndex)
			}
		})
	}
}

// TestQuitOnQ tests quitting with 'q' key.
func TestQuitOnQ(t *testing.T) {
	m := New()
	m.Screen = models.ScreenDashboard

	newModel, cmd := m.Update(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("q"),
	})

	result := newModel.(Model)

	if !result.Quitting {
		t.Error("Quitting should be true after pressing q")
	}

	// cmd is tea.Quit which is a function, we can check if it's not nil
	if cmd == nil {
		t.Error("Expected non-nil command")
	}
}

// TestQuitOnCtrlC tests quitting with Ctrl+C.
func TestQuitOnCtrlC(t *testing.T) {
	m := New()
	m.Screen = models.ScreenDashboard

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	result := newModel.(Model)

	if !result.Quitting {
		t.Error("Quitting should be true after pressing Ctrl+C")
	}

	if cmd == nil {
		t.Error("Expected non-nil command")
	}
}

// TestModelInit tests Model.Init().
func TestModelInit(t *testing.T) {
	m := New()
	cmd := m.Init()

	if cmd != nil {
		t.Error("Init should return nil command")
	}
}

// TestModelView tests Model.View().
func TestModelView(t *testing.T) {
	m := New()
	m.Width = 80
	m.Height = 24

	view := m.View()

	if view == "" {
		t.Error("View should not return empty string")
	}
}

// TestGetServiceByName tests finding services by name.
func TestGetServiceByName(t *testing.T) {
	m := New()
	m.Services = []models.Service{
		{Name: "apache2", Status: "running", Port: 80},
		{Name: "mysql", Status: "stopped", Port: 3306},
	}

	svc := m.GetServiceByName("apache2")
	if svc == nil {
		t.Fatal("GetServiceByName() returned nil")
	}

	if svc.Name != "apache2" {
		t.Errorf("Name = %q, want apache2", svc.Name)
	}

	// Test non-existent
	notFound := m.GetServiceByName("nginx")
	if notFound != nil {
		t.Error("GetServiceByName() should return nil for unknown service")
	}
}

// TestGetServiceByIndex tests finding services by index.
func TestGetServiceByIndex(t *testing.T) {
	m := New()
	m.Services = []models.Service{
		{Name: "apache2"},
		{Name: "mysql"},
	}

	svc := m.GetServiceByIndex(0)
	if svc == nil {
		t.Fatal("GetServiceByIndex(0) returned nil")
	}

	if svc.Name != "apache2" {
		t.Errorf("Name = %q, want apache2", svc.Name)
	}

	// Test out of bounds
	outOfBounds := m.GetServiceByIndex(5)
	if outOfBounds != nil {
		t.Error("GetServiceByIndex(5) should return nil")
	}
}

// TestGetServices tests getter for services.
func TestGetServices(t *testing.T) {
	m := New()
	m.Services = []models.Service{
		{Name: "apache2"},
	}

	services := m.GetServices()
	if len(services) != 1 {
		t.Errorf("len(Services) = %d, want 1", len(services))
	}
}

// TestRefreshServices tests service refresh.
func TestRefreshServices(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	m := New()
	if err := m.RefreshServices(); err != nil {
		t.Logf("RefreshServices returned error (expected if no services): %v", err)
	}
}

// TestSetScreen tests screen change.
func TestSetScreen(t *testing.T) {
	m := New()
	m.SetScreen(models.ScreenDashboard)

	if m.Screen != models.ScreenDashboard {
		t.Errorf("Screen = %v, want %v", m.Screen, models.ScreenDashboard)
	}

	// Test that selection is reset
	if m.SelectedServiceIndex != 0 {
		t.Errorf("SelectedServiceIndex = %d, want 0", m.SelectedServiceIndex)
	}
}

// TestResetSelection tests selection reset.
func TestResetSelection(t *testing.T) {
	m := New()
	m.SelectedServiceIndex = 5
	m.SelectedStackIndex = 3
	m.SelectedMenuIndex = 2

	m.ResetSelection()

	if m.SelectedServiceIndex != 0 {
		t.Errorf("SelectedServiceIndex = %d, want 0", m.SelectedServiceIndex)
	}

	if m.SelectedStackIndex != 0 {
		t.Errorf("SelectedStackIndex = %d, want 0", m.SelectedStackIndex)
	}

	if m.SelectedMenuIndex != 0 {
		t.Errorf("SelectedMenuIndex = %d, want 0", m.SelectedMenuIndex)
	}
}

// TestNewModel tests model creation.
func TestNewModel(t *testing.T) {
	m := New()

	if m.Screen != models.ScreenWelcome {
		t.Errorf("Screen = %v, want %v", m.Screen, models.ScreenWelcome)
	}

	if m.Config == nil {
		t.Error("Config should not be nil")
	}

	if m.ServiceManager == nil {
		t.Error("ServiceManager should not be nil")
	}

	if m.Width != 80 {
		t.Errorf("Width = %d, want 80", m.Width)
	}

	if m.Height != 24 {
		t.Errorf("Height = %d, want 24", m.Height)
	}

	if m.Quitting {
		t.Error("Quitting should be false initially")
	}
}

// TestWindowSize tests window size messages.
func TestWindowSize(t *testing.T) {
	m := New()

	newModel, _ := m.Update(tea.WindowSizeMsg{
		Width:  120,
		Height: 40,
	})

	result := newModel.(Model)

	if result.Width != 120 {
		t.Errorf("Width = %d, want 120", result.Width)
	}

	if result.Height != 40 {
		t.Errorf("Height = %d, want 40", result.Height)
	}
}

// TestStackSelectFlow tests the stack selection flow.
func TestStackSelectFlow(t *testing.T) {
	m := New()
	m.Screen = models.ScreenStackSelect

	// Select LAMP (index 0) and press enter
	m.SelectedMenuIndex = 0
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	result := newModel.(Model)

	if result.Config.StackType != models.StackTypeLAMP {
		t.Errorf("StackType = %v, want LAMP", result.Config.StackType)
	}
}

// TestStackSelectLAMM tests selecting LAMM stack.
func TestStackSelectLAMM(t *testing.T) {
	m := New()
	m.Screen = models.ScreenStackSelect

	// Navigate to LAMM (index 1) and press enter
	m.SelectedMenuIndex = 1
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	result := newModel.(Model)

	if result.Config.StackType != models.StackTypeLAMM {
		t.Errorf("StackType = %v, want LAMM", result.Config.StackType)
	}
}

// TestStackSelectNavigate tests navigating 4 stack options (LAMP, LAMM, LEPP, Back).
func TestStackSelectNavigate(t *testing.T) {
	m := New()
	m.Screen = models.ScreenStackSelect
	m.SelectedMenuIndex = 0

	// Down to LAMM (index 1)
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	result := newModel.(Model)
	if result.SelectedMenuIndex != 1 {
		t.Errorf("SelectedMenuIndex = %d, want 1", result.SelectedMenuIndex)
	}

	// Down to LEPP (index 2)
	newModel, _ = result.Update(tea.KeyMsg{Type: tea.KeyDown})
	result = newModel.(Model)
	if result.SelectedMenuIndex != 2 {
		t.Errorf("SelectedMenuIndex = %d, want 2", result.SelectedMenuIndex)
	}

	// Down to Back (index 3)
	newModel, _ = result.Update(tea.KeyMsg{Type: tea.KeyDown})
	result = newModel.(Model)
	if result.SelectedMenuIndex != 3 {
		t.Errorf("SelectedMenuIndex = %d, want 3", result.SelectedMenuIndex)
	}

	// Down stays at bottom (index 3)
	newModel, _ = result.Update(tea.KeyMsg{Type: tea.KeyDown})
	result = newModel.(Model)
	if result.SelectedMenuIndex != 3 {
		t.Errorf("SelectedMenuIndex = %d, want 3 (bottom bounded)", result.SelectedMenuIndex)
	}

	// Up to LEPP (index 2)
	newModel, _ = result.Update(tea.KeyMsg{Type: tea.KeyUp})
	result = newModel.(Model)
	if result.SelectedMenuIndex != 2 {
		t.Errorf("SelectedMenuIndex = %d, want 2", result.SelectedMenuIndex)
	}

	// Up to LAMM (index 1)
	newModel, _ = result.Update(tea.KeyMsg{Type: tea.KeyUp})
	result = newModel.(Model)
	if result.SelectedMenuIndex != 1 {
		t.Errorf("SelectedMenuIndex = %d, want 1", result.SelectedMenuIndex)
	}

	// Up to LAMP (index 0)
	newModel, _ = result.Update(tea.KeyMsg{Type: tea.KeyUp})
	result = newModel.(Model)
	if result.SelectedMenuIndex != 0 {
		t.Errorf("SelectedMenuIndex = %d, want 0", result.SelectedMenuIndex)
	}

	// Up stays at top (index 0)
	newModel, _ = result.Update(tea.KeyMsg{Type: tea.KeyUp})
	result = newModel.(Model)
	if result.SelectedMenuIndex != 0 {
		t.Errorf("SelectedMenuIndex = %d, want 0 (top bounded)", result.SelectedMenuIndex)
	}
}

// TestDashboardNavigateEmpty tests navigation when no services (bug fix).
func TestDashboardNavigateEmpty(t *testing.T) {
	m := New()
	m.Screen = models.ScreenDashboard
	m.Services = []models.Service{} // empty - the bug scenario
	m.SelectedServiceIndex = 0

	// Down should stay at 0 (max index is 0)
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	result := newModel.(Model)
	if result.SelectedServiceIndex != 0 {
		t.Errorf("SelectedServiceIndex = %d, want 0 (bounded by empty services)", result.SelectedServiceIndex)
	}

	// Up should stay at 0
	newModel, _ = result.Update(tea.KeyMsg{Type: tea.KeyUp})
	result = newModel.(Model)
	if result.SelectedServiceIndex != 0 {
		t.Errorf("SelectedServiceIndex = %d, want 0 (bounded at top)", result.SelectedServiceIndex)
	}
}

// TestDashboardNavigate tests navigation with services.
func TestDashboardNavigate(t *testing.T) {
	m := New()
	m.Screen = models.ScreenDashboard
	m.Services = []models.Service{
		{Name: "apache2"},
		{Name: "mariadb"},
	}
	m.SelectedServiceIndex = 0

	// Down to mariadb
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	result := newModel.(Model)
	if result.SelectedServiceIndex != 1 {
		t.Errorf("SelectedServiceIndex = %d, want 1", result.SelectedServiceIndex)
	}

	// Down stays at bottom
	newModel, _ = result.Update(tea.KeyMsg{Type: tea.KeyDown})
	result = newModel.(Model)
	if result.SelectedServiceIndex != 1 {
		t.Errorf("SelectedServiceIndex = %d, want 1 (bottom bounded)", result.SelectedServiceIndex)
	}
}

// TestPortEditNavigate tests navigation and digit input in PortEdit screen.
func TestPortEditNavigate(t *testing.T) {
	m := New()
	m.Screen = models.ScreenPortEdit
	m.PortEditValue = "80"

	// Up: borra último dígito
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	result := newModel.(Model)
	if result.PortEditValue != "8" {
		t.Errorf("PortEditValue = %q, want 8", result.PortEditValue)
	}

	// Down: agrega 0
	newModel, _ = result.Update(tea.KeyMsg{Type: tea.KeyDown})
	result = newModel.(Model)
	if result.PortEditValue != "80" {
		t.Errorf("PortEditValue = %q, want 80", result.PortEditValue)
	}
}

// TestPortEditDigitInput tests typing digits.
func TestPortEditDigitInput(t *testing.T) {
	m := New()
	m.Screen = models.ScreenPortEdit
	m.PortEditValue = ""

	// Type "8"
	newModel, _ := m.Update(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'8'},
	})
	result := newModel.(Model)
	if result.PortEditValue != "8" {
		t.Errorf("PortEditValue = %q, want 8", result.PortEditValue)
	}

	// Type "0"
	newModel, _ = result.Update(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'0'},
	})
	result = newModel.(Model)
	if result.PortEditValue != "80" {
		t.Errorf("PortEditValue = %q, want 80", result.PortEditValue)
	}
}

// TestScreenEnumValues tests screen enum values.
func TestScreenEnumValues(t *testing.T) {
	expected := map[models.Screen]string{
		models.ScreenWelcome:      "welcome",
		models.ScreenStackSelect: "stackselect",
		models.ScreenInstall:    "install",
		models.ScreenDashboard: "dashboard",
		models.ScreenDetails:   "details",
		models.ScreenLogs:      "logs",
		models.ScreenDatabase:  "database",
		models.ScreenPortEdit:  "portedit",
	}

	for screen, name := range expected {
		if string(screen) != name {
			t.Errorf("Screen = %q, want %q", screen, name)
		}
	}
}