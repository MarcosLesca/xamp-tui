# XAMPP TUI

A terminal UI for managing XAMPP stack on Linux. Built with Bubble Tea.

## Features

- Interactive TUI for XAMPP stack management
- Install LAMP, LAMM, or LEPP stacks
- Start/Stop/Restart services (Apache, MySQL/MariaDB, PHP, FileZilla)
- View service status and logs
- Configure phpMyAdmin port
- Database management interface

## Supported Stacks

| Stack | Components |
|-------|------------|
| **LAMP** | Linux + Apache + MySQL + PHP |
| **LAMM** | Linux + Apache + MariaDB + PHP |
| **LEPP** | Linux + nginx + PostgreSQL + PHP |

## Installation

### From Homebrew Tap

```bash
brew tap MarcosLesca/homebrew-tap
brew install xampp-tui
```

### From Source

```bash
git clone git@github.com:MarcosLesca/xamp-tui.git
cd xampp-tui
go build -o xampp-tui
sudo mv xampp-tui /usr/local/bin/
```

### Quick Install

```bash
curl -sL https://github.com/MarcosLesca/xamp-tui/raw/main/main.go | sudo sh
```

## Usage

```bash
xampp-tui
```

### Navigation

| Key | Action |
|-----|--------|
| `↑` `↓` | Navigate menu |
| `Enter` | Select |
| `Space` | Toggle / Execute |
| `Esc` | Go back / Quit |

## Requirements

- Linux (tested on Ubuntu/Debian)
- Go 1.21+ (for building from source)
- Root privileges for service management

## License

MIT