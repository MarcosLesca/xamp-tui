package service

import (
	"errors"
	"fmt"
	"testing"

	"xampp-tui/internal/models"
)

// MockServiceManager es un mock de ServiceManager para testing.
type MockServiceManager struct {
	services     []models.Service
	statusMap    map[string]string
	shouldError bool
}

// NewMockServiceManager creates a new MockServiceManager for testing.
func NewMockServiceManager() *MockServiceManager {
	return &MockServiceManager{
		services: []models.Service{
			{Name: "apache2", Status: "running", Version: "2.4.52", Port: 80, PID: 1234},
			{Name: "mysql", Status: "stopped", Version: "8.0.35", Port: 3306, PID: 0},
		},
		statusMap: map[string]string{
			"apache2": "running",
			"mysql":  "stopped",
		},
		shouldError: false,
	}
}

func (m *MockServiceManager) DetectServices() ([]models.Service, error) {
	if m.shouldError {
		return nil, errors.New("mock error")
	}
	return m.services, nil
}

func (m *MockServiceManager) GetStatus(serviceName string) (string, error) {
	if m.shouldError {
		return "", errors.New("mock error")
	}
	status, ok := m.statusMap[serviceName]
	if !ok {
		return "error", nil
	}
	return status, nil
}

func (m *MockServiceManager) Start(serviceName string) error {
	if m.shouldError {
		return errors.New("mock error")
	}
	if status, ok := m.statusMap[serviceName]; ok {
		m.statusMap[serviceName] = "running"
		// Update service
		for i := range m.services {
			if m.services[i].Name == serviceName {
				m.services[i].Status = "running"
				m.services[i].PID = 9999
			}
		}
		_ = status // silence unused variable warning
	}
	return nil
}

func (m *MockServiceManager) Stop(serviceName string) error {
	if m.shouldError {
		return errors.New("mock error")
	}
	if _, ok := m.statusMap[serviceName]; ok {
		m.statusMap[serviceName] = "stopped"
		// Update service
		for i := range m.services {
			if m.services[i].Name == serviceName {
				m.services[i].Status = "stopped"
				m.services[i].PID = 0
			}
		}
	}
	return nil
}

func (m *MockServiceManager) Restart(serviceName string) error {
	if m.shouldError {
		return errors.New("mock error")
	}
	// Restart is stop + start
	m.statusMap[serviceName] = "running"
	for i := range m.services {
		if m.services[i].Name == serviceName {
			m.services[i].Status = "running"
			m.services[i].PID = 9999
		}
	}
	return nil
}

func (m *MockServiceManager) GetVersion(serviceName string) (string, error) {
	if m.shouldError {
		return "", errors.New("mock error")
	}
	for _, svc := range m.services {
		if svc.Name == serviceName {
			return svc.Version, nil
		}
	}
	return "", nil
}

func (m *MockServiceManager) InstallStack(stackType models.StackType) (string, error) {
	if m.shouldError {
		return "", errors.New("mock error")
	}
	return "Mock installation complete", nil
}

func (m *MockServiceManager) InstallStackWithProgress(stackType models.StackType, onProgress func(step, total int, message string)) (string, error) {
	if m.shouldError {
		return "", errors.New("mock error")
	}
	onProgress(1, 4, "Starting...")
	onProgress(2, 4, "Installing packages...")
	onProgress(3, 4, "Enabling services...")
	onProgress(4, 4, "Starting services...")
	return "Mock installation complete", nil
}

func (m *MockServiceManager) HasRootAccess() bool {
	return true
}

func (m *MockServiceManager) ChangePort(serviceName string, newPort int) (string, error) {
	if m.shouldError {
		return "", errors.New("mock error")
	}
	return fmt.Sprintf("Port changed to %d for %s", newPort, serviceName), nil
}

// Ensure MockServiceManager implements ServiceManager
var _ ServiceManager = (*MockServiceManager)(nil)

func TestKnownServices(t *testing.T) {
	if len(KnownServices) == 0 {
		t.Error("KnownServices should not be empty")
	}

	// Check for expected services
	expected := []string{"apache2", "nginx", "mysql", "mariadb", "postgresql", "php-fpm"}
	for _, name := range expected {
		found := false
		for _, known := range KnownServices {
			if known == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected service %q not in KnownServices", name)
		}
	}
}

func TestServicePorts(t *testing.T) {
	tests := []struct {
		service string
		want    int
	}{
		{"apache2", 80},
		{"nginx", 80},
		{"mysql", 3306},
		{"mariadb", 3306},
		{"postgresql", 5432},
		{"php-fpm", 9000},
	}

	for _, tt := range tests {
		t.Run(tt.service, func(t *testing.T) {
			port, ok := ServicePorts[tt.service]
			if !ok {
				t.Errorf("ServicePorts[%q] not found", tt.service)
				return
			}
			if port != tt.want {
				t.Errorf("ServicePorts[%q] = %d, want %d", tt.service, port, tt.want)
			}
		})
	}
}

func TestMockDetectServices(t *testing.T) {
	mock := NewMockServiceManager()

	services, err := mock.DetectServices()
	if err != nil {
		t.Fatalf("DetectServices() error = %v", err)
	}

	if len(services) != 2 {
		t.Errorf("len(services) = %d, want %d", len(services), 2)
	}
}

func TestMockGetStatus(t *testing.T) {
	tests := []struct {
		name      string
		service   string
		want      string
		wantError bool
	}{
		{"running service", "apache2", "running", false},
		{"stopped service", "mysql", "stopped", false},
		{"unknown service", "unknown", "error", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockServiceManager()

			status, err := mock.GetStatus(tt.service)
			if (err != nil) != tt.wantError {
				t.Errorf("GetStatus() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if status != tt.want {
				t.Errorf("GetStatus() = %q, want %q", status, tt.want)
			}
		})
	}
}

func TestMockStart(t *testing.T) {
	mock := NewMockServiceManager()

	// Start mysql
	if err := mock.Start("mysql"); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Verify it's running
	status, _ := mock.GetStatus("mysql")
	if status != "running" {
		t.Errorf("mysql status = %q, want running", status)
	}
}

func TestMockStop(t *testing.T) {
	mock := NewMockServiceManager()

	// Stop apache2
	if err := mock.Stop("apache2"); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}

	// Verify it's stopped
	status, _ := mock.GetStatus("apache2")
	if status != "stopped" {
		t.Errorf("apache2 status = %q, want stopped", status)
	}
}

func TestMockRestart(t *testing.T) {
	mock := NewMockServiceManager()

	// Restart apache2
	if err := mock.Restart("apache2"); err != nil {
		t.Fatalf("Restart() error = %v", err)
	}

	// Verify it's still running
	status, _ := mock.GetStatus("apache2")
	if status != "running" {
		t.Errorf("apache2 status = %q, want running", status)
	}
}

func TestMockGetVersion(t *testing.T) {
	mock := NewMockServiceManager()

	version, err := mock.GetVersion("apache2")
	if err != nil {
		t.Fatalf("GetVersion() error = %v", err)
	}

	if version != "2.4.52" {
		t.Errorf("apache2 version = %q, want 2.4.52", version)
	}
}

func TestMockError(t *testing.T) {
	mock := NewMockServiceManager()
	mock.shouldError = true

	_, err := mock.DetectServices()
	if err == nil {
		t.Error("DetectServices() should return error when shouldError=true")
	}
}

func TestStatusConstants(t *testing.T) {
	if StatusRunning != "running" {
		t.Errorf("StatusRunning = %q, want 'running'", StatusRunning)
	}

	if StatusStopped != "stopped" {
		t.Errorf("StatusStopped = %q, want 'stopped'", StatusStopped)
	}

	if StatusError != "error" {
		t.Errorf("StatusError = %q, want 'error'", StatusError)
	}
}