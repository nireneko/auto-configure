## Exploration: 80% Test Coverage & Edge Cases

### Current State
El proyecto actualmente tiene la mayoría de los paquetes por encima del 80% de cobertura, pero sufre de errores de compilación que impiden la ejecución global de los tests.
Hay varios paquetes que tienen fallas o no cumplen el objetivo del 80% de cobertura:
- `openvpn` (54.2%)
- `nvm` (76.4%)
- `homebrew` (77.3%)

Además, los errores de compilación son:
- `update_ids.go`: importaciones no usadas (`os`, `path/filepath`, `strings`).
- `cmd/so-install/main_test.go`: variable `exitCode` declarada pero no usada.
- `internal/presentation/tui/model_extra_test.go`: error al procesar el retorno de `m.Update` en contextos de un solo valor (esperando `tea.Model, tea.Cmd` en vez de un solo valor).

### Affected Areas
- `update_ids.go` — limpieza de imports no usados.
- `cmd/so-install/main_test.go` — eliminar declaración de `exitCode` si no se usa.
- `internal/presentation/tui/model_extra_test.go` — ajustar las llamadas a `m.Update(msg)` para recibir múltiples valores (ej. `m, _ = m.Update(msg)`).
- `internal/infrastructure/openvpn/openvpn_test.go` — tests insuficientes, faltan edge cases.
- `internal/infrastructure/nvm/nvm_test.go` — tests insuficientes, faltan edge cases.
- `internal/infrastructure/homebrew/homebrew_test.go` — tests insuficientes, faltan edge cases.

### Approaches
1. **Corregir compilación y agregar tests a paquetes críticos** — Limpiar imports y variables sin uso. Corregir las llamadas a Update de Bubbletea en el test de TUI. Luego iterar por los paquetes de openvpn, nvm y homebrew añadiendo edge cases (rutas que no existen, errores de ejecución de comandos, falta de binarios, etc.).
   - Pros: Alcanza la meta requerida e incorpora resiliencia en casos límite.
   - Cons: Puede ser tedioso verificar cada edge case manual.
   - Effort: Low/Medium

### Recommendation
Recomiendo proceder con la corrección de errores de sintaxis y type-checking de forma prioritaria, y después agregar los edge cases para las implementaciones de openvpn, nvm y homebrew, incrementando de esta manera la cobertura global por encima del 80% de forma segura.

### Risks
- Los edge cases en implementaciones del sistema (como openvpn, ddev, nvm) podrían requerir un uso extenso de `pkg/mocks` que hay que configurar de forma particular, lo que podría tomar más tiempo.

### Ready for Proposal
Yes — the project is ready to formally propose these changes.