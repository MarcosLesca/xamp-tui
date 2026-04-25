# Skill Registry

## Overview
This file maps project contexts to available AI agent skills. Skills are auto-loaded when relevant code or task context is detected.

## Project
- **Name**: xampp-tui
- **Stack**: Go 1.22.2, Bubbletea TUI
- **Type**: CLI/TUI Application

---

## User Skills

| Trigger | Skill | When to Use |
|---------|-------|-------------|
| Go tests, testing | go-testing | Writing Go tests, using teatest for TUI |
| SDD workflow | sdd-init | Initialize SDD context |
| SDD workflow | sdd-propose | Create change proposal |
| SDD workflow | sdd-spec | Write specifications |
| SDD workflow | sdd-design | Technical design |
| SDD workflow | sdd-tasks | Task breakdown |
| SDD workflow | sdd-apply | Implement code |
| SDD workflow | sdd-verify | Verify implementation |
| SDD workflow | sdd-archive | Archive completed change |
| New skill creation | skill-creator | Creating new AI skills |

---

## Skill Paths (for reference)

### User-Level
- `~/.claude/skills/`
- `~/.config/opencode/skills/`

### Project-Level
- `.claude/skills/` (xampp-tui)
- `.agent/skills/`

---

## Active Triggers

This project activates:
- **go-testing**: When writing Go test files (`*_test.go`)
- **sdd-***: For any SDD workflow command
- **skill-creator**: If user requests new skill creation