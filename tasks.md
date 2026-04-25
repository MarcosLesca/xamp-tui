# Tasks: xampp-tui

## Phase 1: Infrastructure & Project Setup
- [x] 1.1 Initialize Go module: already done (go.mod exists)
- [x] 1.2 Add dependencies: bubbletea, lipgloss
- [x] 1.3 Create project directory structure
- [x] 1.4 Create `internal/config/config.go`

## Phase 2: Core Models
- [x] 2.1 Create `internal/models/service.go`: Service struct
- [x] 2.2 Create `internal/models/stacktype.go`: StackType enum
- [x] 2.3 Create `internal/models/screen.go`: Screen enum
- [x] 2.4 Complete Config struct in config.go

## Phase 3: Service Layer
- [x] 3.1 Create `internal/service/interfaces.go`: ServiceManager interface with Detect, Start, Stop, Restart, GetStatus methods
- [x] 3.2 Create `internal/service/linux.go`: ServiceManager implementation using systemctl commands
- [x] 3.3 Implement service detection (apache2, nginx, mysql, postgresql, php) via `which`/`dpkg`
- [x] 3.4 Implement service control (start/stop/restart) via systemctl
- [x] 3.5 Implement status polling with 5s interval

## Phase 4: System Info
- [x] 4.1 Create `internal/sysinfo/host.go`: Hostname, IP address retrieval
- [x] 4.2 Create `internal/sysinfo/resources.go`: CPU, RAM usage retrieval
- [x] 4.3 Create `internal/sysinfo/ports.go`: Port checking for services (80, 3306, 5432)

## Phase 5: TUI Core
- [x] 5.1 Create `internal/tui/model.go`: Main Model struct with state machine
- [x] 5.2 Create `internal/tui/update.go`: Update function handling msgs
- [x] 5.3 Create `internal/tui/view.go`: Base View function with header/footer (integrated in screens)
- [x] 5.4 Implement Screen transitions based on state

## Phase 6: Screen Implementations
- [x] 6.1 Create `internal/tui/screens/welcome.go`: Welcome screen with "Choose Stack" option
- [x] 6.2 Create `internal/tui/screens/stackselect.go`: StackSelect screen (LAMP vs LEPP options)
- [x] 6.3 Create `internal/tui/screens/install.go`: Installation screen with apt-get commands
- [x] 6.4 Create `internal/tui/screens/dashboard.go`: Dashboard with service status grid + auto-refresh
- [x] 6.5 Create `internal/tui/screens/details.go`: Details screen with version, uptime, ports, controls
- [x] 6.6 Create `internal/tui/screens/logs.go`: Log viewer for Apache/Nginx/MySQL/PostgreSQL
- [x] 6.7 Create `internal/tui/screens/database.go`: DB management (list/create/delete)

## Phase 7: Wiring & Integration
- [x] 7.1 Review `cmd/main.go`: Entry point with bubbletea.NewProgram (already created in phase 6)
- [x] 7.2 Wire Config to TUI model
- [x] 7.3 Wire ServiceManager to Dashboard screen
- [x] 7.4 Implement config persistence (`~/.config/xampp-tui/config.json`)

## Phase 8: Testing
- [x] 8.1 Write unit tests for `internal/models/service.go`
- [x] 8.2 Write unit tests for `internal/config/config.go` (mock filesystem)
- [x] 8.3 Write mock tests for `internal/service/interfaces.go` (mock ServiceManager)
- [x] 8.4 Write integration tests for Screen transitions
- [x] 8.5 Run `go test ./...` to verify all passing