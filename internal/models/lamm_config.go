package models

// LAMMConfig define la configuración para instalación LAMM.
type LAMMConfig struct {
	// MariaDB
	MariaDBRootPassword string
	RemoveAnonymous   bool
	DisallowRemote   bool
	RemoveTestDB     bool
	ReloadPrivileges bool

	// phpMyAdmin
	WebServer      string // apache2 o nginx
	PHPMyAdminPath string // /phpmyadmin, /phpmyadmin, vacío= raíz
}

// DefaultLAMMConfig retorna configuración por defecto.
func DefaultLAMMConfig() LAMMConfig {
	return LAMMConfig{
		MariaDBRootPassword: "",
		RemoveAnonymous:      true,
		DisallowRemote:     true,
		RemoveTestDB:        true,
		ReloadPrivileges:     true,
		WebServer:           "apache2",
		PHPMyAdminPath:      "/phpmyadmin",
	}
}