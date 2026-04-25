# PRD — phpMyAdmin+MariaDB en LAMP + Bug Selección Dashboard

**Fecha**: 2026-04-21
**Proyecto**: xampp-tui
**Autor**: Marcos L.

---

## Bug — Selección Dashboard Rota

### Problema
Cuando no hay servicios detectados, la navegación con `↑↓` en Dashboard no funciona. Solo se puede seleccionar la primera opción (índice 0).

### Root Cause
`getMaxIndex()` en `internal/tui/update.go:410`:

```go
case models.ScreenDashboard:
    return len(m.Services) - 1   // → -1 cuando len == 0
```

`handleDown()` compara `0 < -1` → **false** → no avanza.

### Fix
```go
case models.ScreenDashboard:
    if len(m.Services) == 0 {
        return 0
    }
    return len(m.Services) - 1
```

### Archivos
| File | Cambio |
|------|--------|
| `internal/tui/update.go` | Fix `getMaxIndex()` línea ~410 |

---

## Feature — LAMP con MariaDB + phpMyAdmin

### Problema
El stack LAMP actual usa `mysql-server`. Se necesita opción con **MariaDB** (drop-in replacement) + **phpMyAdmin** para gestión visual de bases de datos.

### User Value
- MariaDB es 100% compatible con MySQL, más performante, mantenido por comunidad
- phpMyAdmin permite administrar DBs sin CLI (útil para devs que vienen de XAMPP/WAMP)

### Decisión de Diseño
**Nueva opción de stack**: `LAMM` (Linux + Apache + **M**ariaDB + **M**ySQL/phpMyAdmin)

> Nombre: `LAMM` para diferenciarlo de `LAMP` (MySQL) y `LEPP` (PostgreSQL)

### Scope

| Componente | Incluye |
|-----------|---------|
| Servidor web | Apache 2 |
| Base de datos | MariaDB (no MySQL) |
| Lenguaje | PHP |
| Admin DB | phpMyAdmin |

### Paquetes a instalar
```
apache2, mariadb-server, php, libapache2-mod-php, php-mysql, phpmyadmin
```

### Servicios a gestionar
```
apache2, mariadb
```

### Cambios en código
1. **`internal/models/stacktype.go`** — agregar `StackTypeLAMM`
2. **`internal/service/interfaces.go`** — opcional: actualizar comentarios
3. **`internal/service/linux.go`** — `InstallStack()`: case `LAMM` con paquetes y servicios correctos
4. **`internal/tui/screens/stackselect.go`** — agregar "LAMM Stack" + descripción "Apache + MariaDB + phpMyAdmin"
5. **`internal/tui/screens/install.go`** — agregar case para `LAMM` con comandos
6. **`internal/tui/update.go`** — `getMaxIndex()`: `ScreenStackSelect` retorna 2 (3 opciones)

### Modelo de datos
```go
const StackTypeLAMM StackType = "LAMM"
```

### Criterio de Done
- [ ] Bug fix: navegación funciona con 0, 1, 2+ servicios
- [ ] Nueva opción "LAMM Stack" aparece en StackSelect
- [ ] Enter en LAMM → pantalla de instalación con comandos correctos
- [ ] `go test ./...` pasa