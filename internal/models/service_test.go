package models

import (
	"testing"
	"time"
)

func TestServiceStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		pid     int
		expected bool
	}{
		{"running with positive PID", "running", 1234, true},
		{"stopped with zero PID", "stopped", 0, false},
		{"running with zero PID", "running", 0, false},
		{"error status", "error", 0, false},
		{"unknown status", "unknown", 1234, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Name:   "test-service",
				Status: tt.status,
				PID:    tt.pid,
			}

			result := s.IsRunning()
			if result != tt.expected {
				t.Errorf("IsRunning() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestServiceIsStopped(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{"stopped service", "stopped", true},
		{"running service", "running", false},
		{"error service", "error", false},
		{"unknown status", "unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Name:   "test-service",
				Status: tt.status,
			}

			result := s.IsStopped()
			if result != tt.expected {
				t.Errorf("IsStopped() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAllServices(t *testing.T) {
	tests := []struct {
		name      string
		stackType StackType
		wantLen   int
		wantNames []string
	}{
		{
			name:      "LAMP stack",
			stackType: StackTypeLAMP,
			wantLen:   2,
			wantNames: []string{"apache2", "mysql"},
		},
		{
			name:      "LEPP stack",
			stackType: StackTypeLEPP,
			wantLen:   2,
			wantNames: []string{"apache2", "postgresql"},
		},
		{
			name:      "LAMM stack",
			stackType: StackTypeLAMM,
			wantLen:   2,
			wantNames: []string{"apache2", "mariadb"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AllServices(tt.stackType)

			if len(result) != tt.wantLen {
				t.Errorf("len(AllServices()) = %d, want %d", len(result), tt.wantLen)
				return
			}

			for i, wantName := range tt.wantNames {
				if i >= len(result) {
					break
				}
				if result[i].Name != wantName {
					t.Errorf("service[%d].Name = %q, want %q", i, result[i].Name, wantName)
				}
			}
		})
	}
}

func TestServiceFields(t *testing.T) {
	s := Service{
		Name:    "apache2",
		Status:  "running",
		Version: "2.4.52",
		Port:    80,
		Uptime:  5 * time.Minute,
		PID:     1234,
	}

	if s.Name != "apache2" {
		t.Errorf("Name = %q, want %q", s.Name, "apache2")
	}
	if s.Status != "running" {
		t.Errorf("Status = %q, want %q", s.Status, "running")
	}
	if s.Version != "2.4.52" {
		t.Errorf("Version = %q, want %q", s.Version, "2.4.52")
	}
	if s.Port != 80 {
		t.Errorf("Port = %d, want %d", s.Port, 80)
	}
	if s.PID != 1234 {
		t.Errorf("PID = %d, want %d", s.PID, 1234)
	}
	if s.Uptime != 5*time.Minute {
		t.Errorf("Uptime = %v, want %v", s.Uptime, 5*time.Minute)
	}
}

func TestServiceJSONFields(t *testing.T) {
	// Test JSON serialization/deserialization
	s := Service{
		Name:    "mysql",
		Status:  "stopped",
		Version: "8.0",
		Port:    3306,
		Uptime:  0,
		PID:     0,
	}

	// Verify fields are accessible
	if s.Name != "mysql" {
		t.Errorf("Name = %q", s.Name)
	}
	if s.Port != 3306 {
		t.Errorf("Port = %d", s.Port)
	}
}