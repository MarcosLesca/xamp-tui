package models

// StackType represents the type of XAMPP stack being managed.
type StackType string

const (
	StackTypeLAMP StackType = "LAMP" // Linux + Apache + MySQL + PHP
	StackTypeLAMM StackType = "LAMM" // Linux + Apache + MariaDB + phpMyAdmin
	StackTypeLEPP StackType = "LEPP" // Linux + Nginx + PHP + PostgreSQL
)

// IsValid checks if the StackType is valid.
func (s StackType) IsValid() bool {
	switch s {
	case StackTypeLAMP, StackTypeLAMM, StackTypeLEPP:
		return true
	default:
		return false
	}
}

// AllStackTypes returns all valid StackType values.
func AllStackTypes() []StackType {
	return []StackType{StackTypeLAMP, StackTypeLAMM, StackTypeLEPP}
}