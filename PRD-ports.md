# PRD — Cambio de Puertos

**Fecha**: 2026-04-21
**Proyecto**: xampp-tui
**Autor**: Marcos L.

---

## Feature — Cambiar Puerto de un Servicio

### Problema
Los servicios usan puertos por defecto (80, 3306, etc.). Si el usuario quiere usar otro puerto (ej: 8080 para Apache), no tiene forma de cambiarlo desde la TUI.

### User Value
- Flexibilidad para evitar conflictos de puertos
- UX completa: misma app para instalar, gestionar y configurar

### Scope (v1)

| Componente | Detalle |
|-----------|---------|
| Apache2 | Editar `/etc/apache2/ports.conf` — `Listen {puerto}` |
| Nginx | Editar `/etc/nginx/sites-available/default` — `listen {puerto}` |
| MySQL/MariaDB | Editar `/etc/mysql/mariadb.conf.d/*.cnf` — `port={puerto}` |
| PostgreSQL | Editar `/etc/postgresql/*/main/postgresql.conf` — `port={puerto}` |

### Flujo
```
Dashboard → Enter → Details → Navigate a "Change Port" → Enter
  → PortEdit screen: "Puerto: [____] (default: 80)"
  → Enter para confirmar, Esc para cancelar
  → Aplica cambio + restart del servicio → Volver a Details
```

### Modelo de datos
```go
// En models/service.go - agregar metodo
ChangePort(serviceName string, newPort int) error
```

### Archivos a tocar
| Archivo | Cambio |
|---------|--------|
| `internal/models/service.go` | + `ChangePort()` |
| `internal/service/interfaces.go` | + `ChangePort` en interfaz |
| `internal/service/linux.go` | Implementar ChangePort por servicio |
| `internal/tui/update.go` | + `ScreenPortEdit`, + handlers |
| `internal/tui/model.go` | + `EditPort string` |
| `internal/tui/screens/details.go` | + "Change Port" en menú |
| `new file: internal/tui/screens/portedit.go` | Pantalla de edición de puerto |

### Criterio de Done
- [ ] "Change Port" aparece en Details (para servicios suportados)
- [ ] Editar puerto funciona para Apache2
- [ ] Esc cancela y vuelve a Details sin cambios
- [ ] Puerto inválido (< 1 o > 65535) da error
- [ ] go test ./... pasa