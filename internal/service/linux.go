package service

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"xampp-tui/internal/models"
)

// LinuxServiceManager es una implementación de ServiceManager para Linux.
// Utiliza systemctl para gestionar servicios del sistema.
type LinuxServiceManager struct {
	// pollInterval define el intervalo de polling para status.
	pollInterval time.Duration
}

// NewLinuxServiceManager crea un nuevo LinuxServiceManager.
func NewLinuxServiceManager() *LinuxServiceManager {
	return &LinuxServiceManager{
		pollInterval: 5 * time.Second,
	}
}

// DetectServices detecta los servicios instalados en el sistema.
// Si no están instalados, los muestra con estado "available" para poder instalar.
func (m *LinuxServiceManager) DetectServices() ([]models.Service, error) {
	var services []models.Service

	for _, name := range KnownServices {
		// Verificar si está instalado
		installed := isInstalled(name)

		// Obtener versión
		version, _ := m.GetVersion(name)
		if version == "" && installed {
			version = "unknown"
		}

		// Obtener estado actual
		var status string
		if installed {
			status, _ = m.GetStatus(name)
			if status == "error" {
				status = "stopped"
			}
		} else {
			status = "available"
		}

		// Obtener puerto
		port := ServicePorts[name]
		if port == 0 {
			port = 0
		}

		// Crear modelo de servicio
		svc := models.Service{
			Name:    name,
			Status:  status,
			Version: version,
			Port:    port,
			Uptime: 0,
			PID:    0,
		}

		// Si está corriendo, obtener PID y uptime
		if status == "running" {
			if pid, uptime, err := getProcessInfo(name); err == nil {
				svc.PID = pid
				svc.Uptime = uptime
			}
		}

		services = append(services, svc)
	}

	return services, nil
}

// GetStatus obtiene el estado actual de un servicio.
func (m *LinuxServiceManager) GetStatus(serviceName string) (string, error) {
	// Verificar si el servicio está instalado
	if !isInstalled(serviceName) {
		return "error", fmt.Errorf("servicio %s no instalado", serviceName)
	}

	// Usar systemctl is-active para verificar estado
	cmd := exec.Command("systemctl", "is-active", serviceName)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = nil

	_ = cmd.Run()
	output := strings.TrimSpace(out.String())

	switch output {
	case "active":
		return "running", nil
	case "inactive", "failed":
		return "stopped", nil
	default:
		// Puede ser que el servicio no tenga systemd
		// Intentar verificar con proceso
		return m.checkStatusByProcess(serviceName)
	}
}

// checkStatusByProcess verifica el estado buscando el proceso.
func (m *LinuxServiceManager) checkStatusByProcess(serviceName string) (string, error) {
	// phpMyAdmin especial: no es un servicio systemd
	//RETORNA "installed" si existe el directorio
	if serviceName == "phpmyadmin" {
		paths := []string{
			"/usr/share/phpmyadmin/index.php",
			"/var/www/html/phpmyadmin/index.php",
		}
		for _, p := range paths {
			if _, err := os.Stat(p); err == nil {
				return "installed", nil // Estado especial: instalado pero no es servicio
			}
		}
		return "stopped", nil
	}

	// Buscar proceso por nombre
	procPatterns := map[string][]string{
		"apache2":   {"apache2", "httpd"},
		"nginx":     {"nginx"},
		"mysql":     {"mysqld", "mariadb"},
		"mariadb":   {"mysqld", "mariadb"},
		"postgresql": {"postgres"},
		"php-fpm":  {"php-fpm"},
	}

	patterns, ok := procPatterns[serviceName]
	if !ok {
		return "stopped", nil
	}

	for _, pattern := range patterns {
		cmd := exec.Command("pgrep", "-x", pattern)
		if err := cmd.Run(); err == nil {
			return "running", nil
		}
	}

	return "stopped", nil
}

// Start inicia un servicio.
func (m *LinuxServiceManager) Start(serviceName string) error {
	if !isInstalled(serviceName) {
		return fmt.Errorf("servicio %s no instalado", serviceName)
	}

	cmd := exec.Command("systemctl", "start", serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error al iniciar %s: %w", serviceName, err)
	}

	return nil
}

// Stop detiene un servicio.
func (m *LinuxServiceManager) Stop(serviceName string) error {
	if !isInstalled(serviceName) {
		return fmt.Errorf("servicio %s no instalado", serviceName)
	}

	cmd := exec.Command("systemctl", "stop", serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error al detener %s: %w", serviceName, err)
	}

	return nil
}

// Restart reinicia un servicio.
func (m *LinuxServiceManager) Restart(serviceName string) error {
	if !isInstalled(serviceName) {
		return fmt.Errorf("servicio %s no instalado", serviceName)
	}

	cmd := exec.Command("systemctl", "restart", serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error al reiniciar %s: %w", serviceName, err)
	}

	return nil
}

// GetVersion obtiene la versión de un servicio.
func (m *LinuxServiceManager) GetVersion(serviceName string) (string, error) {
	if !isInstalled(serviceName) {
		return "", fmt.Errorf("servicio %s no instalado", serviceName)
	}

	// Comandos para obtener versión de cada servicio
	versionCmds := map[string][]string{
		"apache2":   {"apache2", "-v"},
		"nginx":     {"nginx", "-v"},
		"mysql":     {"mysql", "--version"},
		"mariadb":   {"mariadb", "--version"},
		"postgresql": {"postgres", "--version"},
		"php-fpm":   {"php-fpm", "-v"},
	}

	args, ok := versionCmds[serviceName]
	if !ok {
		return "", nil
	}

	cmd := exec.Command(args[0], args[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return "", nil
	}

	// Parsear versión del output
	output := out.String()
	version := parseVersion(output)

	return version, nil
}

// isInstalled verifica si un servicio está instalado en el sistema.
// Mapea nombres lógicos a nombres de paquetes reales.
func isInstalled(serviceName string) bool {
	// Mapeo de nombres lógicos a paquetes Debian
	packageMap := map[string]string{
		"apache2":    "apache2",
		"nginx":      "nginx",
		"mysql":      "mysql-server",
		"mariadb":    "mariadb-server",
		"postgresql":  "postgresql",
		"php-fpm":    "php-fpm",
		"phpmyadmin": "phpmyadmin",
	}

	pkg := serviceName
	if mapped, ok := packageMap[serviceName]; ok {
		pkg = mapped
	}

	// Verificar con dpkg
	cmd := exec.Command("dpkg", "-l", pkg)
	if err := cmd.Run(); err == nil {
		return true
	}

	// Verificar con which (para binaries que no son paquetes)
	binaryMap := map[string]string{
		"apache2":    "apache2",
		"nginx":     "nginx",
		"mysql":     "mysqld",
		"mariadb":   "mysqld",
		"postgresql": "postgres",
		"php-fpm":  "php-fpm",
		"phpmyadmin": "index.php", // phpMyAdmin es un archivo, no un binario
	}

	if bin, ok := binaryMap[serviceName]; ok {
		// Para phpMyAdmin, buscar el archivo index.php
		if serviceName == "phpmyadmin" {
			paths := []string{
				"/usr/share/phpmyadmin/index.php",
				"/var/www/html/phpmyadmin/index.php",
			}
			for _, p := range paths {
				if _, err := os.Stat(p); err == nil {
					return true
				}
			}
			return false
		}
		// Para otros, buscar binario
		cmd := exec.Command("which", bin)
		if err := cmd.Run(); err == nil {
			return true
		}
	}

	return false
}

// getProcessInfo obtiene el PID y uptime de un servicio.
// Intenta primero con systemctl, luego con pgrep.
func getProcessInfo(serviceName string) (int, time.Duration, error) {
	// Primero intentar con systemctl show para obtener el PID del main process
	cmd := exec.Command("systemctl", "show", serviceName, "-p", "MainPID", "--value")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = nil

	if err := cmd.Run(); err == nil {
		pidStr := strings.TrimSpace(out.String())
		if pidStr != "0" && pidStr != "" {
			pid, err := strconv.Atoi(pidStr)
			if err == nil && pid > 0 {
				uptime, _ := getProcessUptime(pid)
				return pid, uptime, nil
			}
		}
	}

	// Fallback: buscar proceso con pgrep
	procPatterns := map[string][]string{
		"apache2":   {"apache2", "httpd"},
		"nginx":     {"nginx"},
		"mysql":     {"mysqld", "mariadb"},
		"mariadb":   {"mysqld", "mariadb"},
		"postgresql": {"postgres"},
		"php-fpm":  {"php-fpm"},
	}

	patterns, ok := procPatterns[serviceName]
	if !ok {
		return 0, 0, fmt.Errorf("patrones no encontrados")
	}

	for _, pattern := range patterns {
		cmd := exec.Command("pgrep", "-x", pattern)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = nil

		if err := cmd.Run(); err != nil {
			continue
		}

		pidStr := strings.TrimSpace(out.String())
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue
		}

		// Obtener uptime del proceso
		uptime, err := getProcessUptime(pid)
		if err != nil {
			uptime = 0
		}

		return pid, uptime, nil
	}

	return 0, 0, fmt.Errorf("proceso no encontrado")
}

// getProcessUptime obtiene el uptime de un proceso dado su PID.
func getProcessUptime(pid int) (time.Duration, error) {
	// Leer inicio del proceso desde /proc/pid/stat
	statPath := fmt.Sprintf("/proc/%d/stat", pid)
	cmd := exec.Command("cat", statPath)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return 0, err
	}

	// El formato es: pid (comm) state ppid pgrp session tty...
	// El campo 21 (index 20 desde 1) es el tiempo de inicio en ticks
	output := out.String()
	fields := strings.Fields(output)
	if len(fields) < 22 {
		return 0, fmt.Errorf("formato de stat inválido")
	}

	startTime, err := strconv.ParseInt(fields[21], 10, 64)
	if err != nil {
		return 0, err
	}

	// Obtener uptime del sistema
	sysUptime, err := getSystemUptime()
	if err != nil {
		return 0, err
	}

	// Calcular uptime del proceso (en segundos)
	// startTime está en clock ticks desde el inicio del sistema
	clockTicks := float64(sysUptime)
	procUptime := clockTicks - float64(startTime)/100.0 // 100 Hz typical

	return time.Duration(procUptime) * time.Second, nil
}

// getSystemUptime obtiene el uptime del sistema en segundos.
func getSystemUptime() (float64, error) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}

	fields := strings.Fields(string(data))
	if len(fields) < 1 {
		return 0, fmt.Errorf("formato de uptime inválido")
	}

	uptime, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, err
	}

	return uptime, nil
}

// parseVersion extrae la versión de un string de output de comando.
func parseVersion(output string) string {
	// Buscar patrón de versión: "version X.Y.Z" o "vX.Y.Z"
	re := regexp.MustCompile(`(?:v|version\s*)(\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(output)

	if len(matches) > 1 {
		return matches[1]
	}

	// Si no encuentra patrón común, buscar cualquier número de versión
	re = regexp.MustCompile(`(\d+\.\d+\.\d+)`)
	matches = re.FindStringSubmatch(output)
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

// PollInterval devuelve el intervalo de polling.
func (m *LinuxServiceManager) PollInterval() time.Duration {
	return m.pollInterval
}

// HasRootAccess verifica si tiene permisos de root.
func (m *LinuxServiceManager) HasRootAccess() bool {
	return hasSudo() || hasPkexec() || hasDoas()
}

// InstallStack instala el stack seleccionado paso a paso para mostrar progreso.
// Retorna un log acumulativo de cada paso.
func (m *LinuxServiceManager) InstallStack(stackType models.StackType) (string, error) {
	var packages []string
	var services []string
	var enableCmds []string
	var startCmds []string

	switch stackType {
	case models.StackTypeLAMP:
		packages = []string{"apache2", "mysql-server", "php", "libapache2-mod-php", "php-mysql"}
		services = []string{"apache2", "mysql"}
	case models.StackTypeLAMM:
		packages = []string{"apache2", "mariadb-server", "php", "libapache2-mod-php", "php-mysql", "phpmyadmin"}
		services = []string{"apache2", "mariadb"}
	case models.StackTypeLEPP:
		packages = []string{"nginx", "postgresql", "php-fpm", "php-pgsql"}
		services = []string{"nginx", "postgresql"}
	default:
		return "", fmt.Errorf("stack desconocido: %v", stackType)
	}

	// Build commands arrays
	for _, svc := range services {
		enableCmds = append(enableCmds, fmt.Sprintf("systemctl enable %s", svc))
		startCmds = append(startCmds, fmt.Sprintf("systemctl start %s", svc))
	}

	var log strings.Builder
	log.WriteString("=== XAMPP-TUI Installation ===\n\n")

	// Step 1: apt update
	log.WriteString("[1/4] Updating package list...\n")
	var err error
	_, err = runCmdWithOutput("pkexec", "apt-get", "update", "-qq")
	if err != nil {
		out, err2 := runCmdWithOutput("sudo", "apt-get", "update", "-qq")
		if err2 != nil {
			return log.String(), fmt.Errorf("apt update failed: %v", err2)
		}
		log.WriteString(out)
		log.WriteString("✓ Done\n\n")
	}

	// Step 2: Install packages
	log.WriteString(fmt.Sprintf("[2/4] Installing packages: %s\n", strings.Join(packages, ", ")))
	args := []string{"apt-get", "install", "-y"}
	args = append(args, packages...)
	_, err = runCmdWithOutput("pkexec", args...)
	if err != nil {
		out, err2 := runCmdWithOutput("sudo", args...)
		if err2 != nil {
			return log.String(), fmt.Errorf("apt install failed: %v", err2)
		}
		log.WriteString(out)
		log.WriteString("✓ Done\n\n")
	}

	// Step 3: Enable services
	log.WriteString(fmt.Sprintf("[3/4] Enabling services: %s\n", strings.Join(services, ", ")))
	for _, cmd := range enableCmds {
		log.WriteString(fmt.Sprintf("  > %s\n", cmd))
		_, err = runCmdWithOutput("pkexec", "bash", "-c", cmd)
		if err != nil {
			_, err2 := runCmdWithOutput("sudo", "bash", "-c", cmd)
			if err2 != nil {
				log.WriteString(fmt.Sprintf("  ⚠ %v (continuing)\n", err2))
			}
		}
	}
	log.WriteString("✓ Done\n\n")

	// Step 4: Start services
	log.WriteString(fmt.Sprintf("[4/4] Starting services: %s\n", strings.Join(services, ", ")))
	for _, cmd := range startCmds {
		log.WriteString(fmt.Sprintf("  > %s\n", cmd))
		_, err = runCmdWithOutput("pkexec", "bash", "-c", cmd)
		if err != nil {
			_, err2 := runCmdWithOutput("sudo", "bash", "-c", cmd)
			if err2 != nil {
				log.WriteString(fmt.Sprintf("  ⚠ %v (continuing)\n", err2))
			}
		}
	}
	log.WriteString("✓ Done\n\n")

	log.WriteString("=== Installation Complete ===\n")
	log.WriteString(fmt.Sprintf("Stack: %s\n", stackType))
	log.WriteString("Services: " + strings.Join(services, ", ") + "\n")

	return log.String(), nil
}

// InstallStackWithProgress installs with callback using ONE password.
// All commands are bundled into a single script.
func (m *LinuxServiceManager) InstallStackWithProgress(stackType models.StackType, onProgress func(step, total int, message string)) (string, error) {
	var packages []string
	var services []string

	switch stackType {
	case models.StackTypeLAMP:
		packages = []string{"apache2", "mysql-server", "php", "libapache2-mod-php", "php-mysql"}
		services = []string{"apache2", "mysql"}
	case models.StackTypeLAMM:
		packages = []string{"apache2", "mariadb-server", "php", "libapache2-mod-php", "php-mysql", "phpmyadmin"}
		services = []string{"apache2", "mariadb"}
	case models.StackTypeLEPP:
		packages = []string{"nginx", "postgresql", "php-fpm", "php-pgsql"}
		services = []string{"nginx", "postgresql"}
	default:
		return "", fmt.Errorf("stack desconocido: %v", stackType)
	}

	var log strings.Builder
	total := 4

	// Build single script - DON'T exit on errors
	var script strings.Builder
	script.WriteString("#!/bin/bash\n")
	script.WriteString("# Don't exit on errors - continue even if something fails\n")

	script.WriteString("echo '[1/4] Updating package list'\n")
	script.WriteString("sudo apt-get update 2>&1\n")

	script.WriteString("echo '[2/4] Installing packages'\n")
	script.WriteString(fmt.Sprintf("sudo apt-get install -y %s 2>&1\n", strings.Join(packages, " ")))

	script.WriteString("echo '[3/4] Enabling services'\n")
	for _, svc := range services {
		script.WriteString(fmt.Sprintf("sudo systemctl enable %s 2>&1 || echo 'enable failed: %s'\n", svc, svc))
	}

	script.WriteString("echo '[4/4] Starting services'\n")
	for _, svc := range services {
		script.WriteString(fmt.Sprintf("sudo systemctl start %s 2>&1 || echo 'start failed: %s'\n", svc, svc))
	}

	script.WriteString("echo '[DONE]'\n")
	script.WriteString("echo 'Checking status...'\n")
	for _, svc := range services {
		script.WriteString(fmt.Sprintf("sudo systemctl is-active %s 2>&1 || echo '%s is not active'\n", svc, svc))
	}

	// Write script
	scriptPath := "/tmp/xampp-install.sh"
	if err := os.WriteFile(scriptPath, []byte(script.String()), 0755); err != nil {
		return "", fmt.Errorf("error creating script: %w", err)
	}
	defer os.Remove(scriptPath)

	var err error
	var output string

	onProgress(1, total, "Installing...")
	log.WriteString("[1/4] Starting installation...\n")

	// Run with pkexec for the WHOLE script - just ONE password
	output, err = runCmdWithOutput("pkexec", "bash", scriptPath)

	// ALWAYS show some progress, even on error
	onProgress(2, total, "Running...")
	log.WriteString(fmt.Sprintf("[2/4] Running... output: %s\n", output))

	if err != nil {
		// Fallback to sudo
		output2, err2 := runCmdWithOutput("sudo", "bash", scriptPath)
		onProgress(2, total, "Retrying with sudo...")
		
		if err2 != nil {
			onProgress(2, total, fmt.Sprintf("ERROR: %v", err2))
			log.WriteString(fmt.Sprintf("Install error: %v\nOutput:\n%v\n", err2, output2))
			return log.String(), fmt.Errorf("install failed: %v", err2)
		}
		output = output2
	}

	// Parse output
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		switch {
		case strings.Contains(line, "[1/4]"):
			onProgress(1, total, "Updating...")
			log.WriteString("[1/4] Updating...\n")
		case strings.Contains(line, "[2/4]"):
			onProgress(2, total, "Installing...")
			log.WriteString("[2/4] Installing...\n")
		case strings.Contains(line, "[3/4]"):
			onProgress(3, total, "Enabling...")
			log.WriteString("[3/4] Enabling...\n")
		case strings.Contains(line, "[4/4]"):
			onProgress(4, total, "Starting...")
			log.WriteString("[4/4] Starting...\n")
		case strings.Contains(line, "[DONE]"):
			onProgress(4, total, "Complete!")
			log.WriteString("[4/4] Complete!\n")
		case strings.Contains(line, "FAILED:"):
			log.WriteString(fmt.Sprintf("⚠ %s\n", line))
		}
	}

	onProgress(4, total, "✓ Complete")
	log.WriteString("=== Installation Complete ===\n")
	log.WriteString(fmt.Sprintf("Stack: %s\n", stackType))
	log.WriteString(fmt.Sprintf("Services: %s\n", strings.Join(services, ", ")))

	return log.String(), nil
}

// joinCmds une comandos con prefijo.
func joinCmds(prefix string, items []string) string {
	var cmds []string
	for _, item := range items {
		cmds = append(cmds, fmt.Sprintf("sudo %s %s 2>/dev/null || true", prefix, item))
	}
return strings.Join(cmds, "\n")
}

// detectElevator detecta qué método deelevation está disponible.
// Retorna "sudo-nopass", "pkexec", "doas", o "none".
func detectElevator() string {
	// 1. sudo sin password
	if hasSudo() {
		return "sudo-nopass"
	}

	// 2. pkexec (polkit)
	if hasPkexec() {
		return "pkexec"
	}

	// 3. doas
	if hasDoas() {
		return "doas"
	}

	return "none"
}

// hasSudo verifica si el usuario tiene permisos sudo sin contraseña.
func hasSudo() bool {
	cmd := exec.Command("sudo", "-n", "true")
	return cmd.Run() == nil
}

// hasPkexec verifica si pkexec está disponible y funciona.
func hasPkexec() bool {
	// pkexec requiere configuración de polkit, intentarlo directamente
	cmd := exec.Command("pkexec", "--version")
	return cmd.Run() == nil
}

// hasDoas verifica si doas está disponible.
func hasDoas() bool {
	cmd := exec.Command("doas", "echo", "test")
	return cmd.Run() == nil
}

// runElevated ejecuta un comando con elevation.
func runElevated(elevator string, cmdArgs ...string) (string, error) {
	switch elevator {
	case "sudo-nopass":
		// Intentar directamente
		out, err := runCmdWithOutput("sudo", cmdArgs...)
		if err == nil {
			return out, nil
		}
		// Si falló sin TTY, intentar con setsid
		allArgs := []string{"sudo"}
		allArgs = append(allArgs, cmdArgs...)
		return runCmdWithOutput("setsid", allArgs...)
	case "pkexec":
		return runCmdWithOutput("pkexec", cmdArgs...)
	case "doas":
		return runCmdWithOutput("doas", cmdArgs...)
	default:
		return "", fmt.Errorf("no hay método de elevation disponible")
	}
}

// runCmdCmd crea un comando para ejecutar.
func runCmd(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}

// runCmdWithOutput ejecuta un comando y retorna output+error.
func runCmdWithOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	output := stdout.String()
	if stderr.Len() > 0 {
		output += "\n" + stderr.String()
	}
	return output, err
}

// Ensure LinuxServiceManager implementa ServiceManager
var _ ServiceManager = (*LinuxServiceManager)(nil)

// ChangePort cambia el puerto de un servicio.
// Solo soporta servicios que usan archivos de configuración editables.
func (m *LinuxServiceManager) ChangePort(serviceName string, newPort int) (string, error) {
	if newPort < 1 || newPort > 65535 {
		return "", fmt.Errorf("puerto inválido: %d (debe estar entre 1 y 65535)", newPort)
	}

	if !isInstalled(serviceName) {
		return "", fmt.Errorf("servicio %s no instalado", serviceName)
	}

	switch serviceName {
	case "apache2":
		return changeApachePort(newPort)
	case "nginx":
		return changeNginxPort(newPort)
	case "mysql", "mariadb":
		return changeMySQLPort(serviceName, newPort)
	case "postgresql":
		return changePostgresPort(newPort)
	default:
		return "", fmt.Errorf("servicio %s no soporta cambio de puerto", serviceName)
	}
}

// changeApachePort cambia el puerto de Apache.
func changeApachePort(newPort int) (string, error) {
	portsConf := "/etc/apache2/ports.conf"

	if _, err := os.Stat(portsConf); os.IsNotExist(err) {
		// Crear el archivo si no existe
		content := fmt.Sprintf("Listen %d\n", newPort)
		if err := os.WriteFile(portsConf, []byte(content), 0644); err != nil {
			return "", fmt.Errorf("error creando %s: %w", portsConf, err)
		}
		return fmt.Sprintf("Puerto configurado en %s", portsConf), nil
	}

	// Editar Listen en ports.conf
	data, err := os.ReadFile(portsConf)
	if err != nil {
		return "", fmt.Errorf("error leyendo %s: %w", portsConf, err)
	}

	lines := strings.Split(string(data), "\n")
	var newLines []string
	found := false
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "Listen ") {
			newLines = append(newLines, fmt.Sprintf("Listen %d", newPort))
			found = true
		} else {
			newLines = append(newLines, line)
		}
	}

	if !found {
		newLines = append(newLines, fmt.Sprintf("Listen %d", newPort))
	}

	if err := os.WriteFile(portsConf, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		return "", fmt.Errorf("error escribiendo %s: %w", portsConf, err)
	}

	return fmt.Sprintf("Puerto Apache cambiado a %d en %s", newPort, portsConf), nil
}

// changeNginxPort cambia el puerto de Nginx.
func changeNginxPort(newPort int) (string, error) {
	siteConfig := "/etc/nginx/sites-available/default"

	if _, err := os.Stat(siteConfig); os.IsNotExist(err) {
		return "", fmt.Errorf("archivo de configuración de Nginx no encontrado")
	}

	data, err := os.ReadFile(siteConfig)
	if err != nil {
		return "", fmt.Errorf("error leyendo %s: %w", siteConfig, err)
	}

	lines := strings.Split(string(data), "\n")
	var newLines []string
	for _, line := range lines {
		if strings.Contains(line, "listen ") && strings.Contains(line, ";") {
			// Reemplazar puerto en líneas listen
			re := regexp.MustCompile(`(\slisten\s+)([0-9]+)`)
			newLine := re.ReplaceAllString(line, fmt.Sprintf("${1}%d", newPort))
			newLines = append(newLines, newLine)
		} else {
			newLines = append(newLines, line)
		}
	}

	if err := os.WriteFile(siteConfig, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		return "", fmt.Errorf("error escribiendo %s: %w", siteConfig, err)
	}

	return fmt.Sprintf("Puerto Nginx cambiado a %d en %s", newPort, siteConfig), nil
}

// changeMySQLPort cambia el puerto de MySQL/MariaDB.
func changeMySQLPort(serviceName string, newPort int) (string, error) {
	// Buscar archivo de configuración de MySQL/MariaDB
	configFiles := []string{
		"/etc/mysql/mariadb.conf.d/50-server.cnf",
		"/etc/mysql/mariadb.conf.d/50-default-cnf",
		"/etc/mysql/my.cnf",
		"/etc/my.cnf",
	}

	var configPath string
	for _, path := range configFiles {
		if _, err := os.Stat(path); err == nil {
			configPath = path
			break
		}
	}

	if configPath == "" {
		return "", fmt.Errorf("archivo de configuración de %s no encontrado", serviceName)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("error leyendo %s: %w", configPath, err)
	}

	lines := strings.Split(string(data), "\n")
	var newLines []string
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "port") &&
			(strings.Contains(line, "=") || strings.HasPrefix(strings.TrimSpace(line), "port")) {
			// Reemplazar línea de puerto
			re := regexp.MustCompile(`(port\s*[=]*\s*)([0-9]+)`)
			if re.MatchString(line) {
				newLine := re.ReplaceAllString(line, fmt.Sprintf("${1}%d", newPort))
				newLines = append(newLines, newLine)
			} else {
				newLines = append(newLines, line)
			}
		} else {
			newLines = append(newLines, line)
		}
	}

	if err := os.WriteFile(configPath, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		return "", fmt.Errorf("error escribiendo %s: %w", configPath, err)
	}

	return fmt.Sprintf("Puerto %s cambiado a %d en %s", serviceName, newPort, configPath), nil
}

// changePostgresPort cambia el puerto de PostgreSQL.
func changePostgresPort(newPort int) (string, error) {
	// Buscar directorio de datos de PostgreSQL
	versions := []string{"15", "14", "13", "12", "11"}
	dataDir := ""

	for _, ver := range versions {
		path := fmt.Sprintf("/etc/postgresql/%s/main/postgresql.conf", ver)
		if _, err := os.Stat(path); err == nil {
			dataDir = path
			break
		}
	}

	if dataDir == "" {
		return "", fmt.Errorf("postgresql.conf no encontrado")
	}

	data, err := os.ReadFile(dataDir)
	if err != nil {
		return "", fmt.Errorf("error leyendo %s: %w", dataDir, err)
	}

	lines := strings.Split(string(data), "\n")
	var newLines []string
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "port") {
			re := regexp.MustCompile(`(#?\s*port\s*[=]*\s*)([0-9]+)`)
			if re.MatchString(line) {
				newLine := re.ReplaceAllString(line, fmt.Sprintf("${1}%d", newPort))
				newLines = append(newLines, newLine)
			} else {
				newLines = append(newLines, line)
			}
		} else {
			newLines = append(newLines, line)
		}
	}

	if err := os.WriteFile(dataDir, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		return "", fmt.Errorf("error escribiendo %s: %w", dataDir, err)
	}

	return fmt.Sprintf("Puerto PostgreSQL cambiado a %d en %s", newPort, dataDir), nil
}

// SupportsPortChange retorna true si el servicio soporta cambio de puerto.
func SupportsPortChange(serviceName string) bool {
	switch serviceName {
	case "apache2", "nginx", "mysql", "mariadb", "postgresql":
		return true
	default:
		return false
	}
}