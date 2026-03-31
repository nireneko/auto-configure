# Design: 80% Test Coverage & Edge Cases

## Technical Approach
Implementar los mocks existentes de `pkg/mocks` (ej: `Executor`) en los tests de `openvpn`, `nvm` y `homebrew` para forzar fallos en los comandos shell y la red. Paralelamente, corregir los errores de compilación mediante limpieza directa de variables y asignaciones múltiples de `tea.Model, tea.Cmd` en el entorno de pruebas de la TUI.

## Architecture Decisions

### Decision: Inyección de Mocks en Tests de Infraestructura
**Choice**: Usar `pkg/mocks.NewMockExecutor()` para inyectar fallos específicos mediante funciones de tipo `On("Execute", ...).Return(..., errors.New(...))`.
**Alternatives considered**: Hacer testing real sobre el sistema con `os/exec`.
**Rationale**: Las pruebas con `os/exec` real son lentas, frágiles y dependientes del OS subyacente. Los mocks aseguran que el path de errores se verifique consistentemente.

## Data Flow
No hay cambios en el flujo de datos real. Solo en el flujo de pruebas:
    Test ──→ [Mock Executor: Error] ──→ Installer ──→ Return Error ──→ Test Assert (Pass)

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `update_ids.go` | Modify | Remover imports no usados (`os`, `path/filepath`, `strings`). |
| `cmd/so-install/main_test.go` | Modify | Remover variable `exitCode` o usarla en una aserción de log. |
| `internal/presentation/tui/model_extra_test.go` | Modify | Reemplazar `m.Update(...)` por `m, _ = m.Update(...)` |
| `internal/infrastructure/openvpn/openvpn_test.go` | Modify | Agregar test para error de descarga y error en `apt update`. |
| `internal/infrastructure/nvm/nvm_test.go` | Modify | Agregar test para error en `wget` o en la ejecución del script. |
| `internal/infrastructure/homebrew/homebrew_test.go` | Modify | Agregar test para cuando `brew` no existe. |

## Interfaces / Contracts
No hay nuevas interfaces. Se utilizará la interface `domain.Executor` y su mock correspondiente provisto por `testify/mock`.

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | Fixes compilación | `go test` sobre `internal/presentation/tui/` |
| Unit | Edge cases infrastructura | Simular errores en el Executor usando mocks y verificar que el Installer los propaga |

## Migration / Rollout
No migration required.

## Open Questions
- None.