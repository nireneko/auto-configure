# Proposal: 80% Test Coverage & Edge Cases

## Intent
Solucionar los errores actuales de compilación en los tests y aumentar la cobertura de los paquetes críticos (`openvpn`, `nvm`, `homebrew`) para que superen el 80%, asegurando que la herramienta maneje adecuadamente los edge cases (ej. fallos de red, permisos, comandos no encontrados).

## Scope

### In Scope
- Corregir errores de imports en `update_ids.go`.
- Remover variables sin uso en `cmd/so-install/main_test.go`.
- Arreglar asignación múltiple en `internal/presentation/tui/model_extra_test.go`.
- Añadir edge cases en tests para `openvpn`, `nvm` y `homebrew`.

### Out of Scope
- Refactorización mayor de la arquitectura.
- Tests E2E completos del ciclo de instalación.
- Aumentar la cobertura de los demás paquetes que ya superan el 80%.

## Approach
1. **Fijar Compilación**: Reparar sintaxis y validaciones de tipos en los tests fallidos y archivos inutilizados.
2. **Mocking y Edge Cases**: Utilizar `pkg/mocks` para simular fallos en la ejecución de shells, rutas inválidas y dependencias faltantes en las pruebas unitarias de `openvpn`, `nvm` y `homebrew`.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `update_ids.go` | Modified | Limpiar imports |
| `cmd/so-install/main_test.go` | Modified | Remover variable sin uso |
| `internal/presentation/tui/model_extra_test.go` | Modified | Corregir firma de `Update` |
| `internal/infrastructure/openvpn/openvpn_test.go` | Modified | Agregar casos límite |
| `internal/infrastructure/nvm/nvm_test.go` | Modified | Agregar casos límite |
| `internal/infrastructure/homebrew/homebrew_test.go` | Modified | Agregar casos límite |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Incompatibilidad de mocks | Medium | Asegurar que `pkg/mocks/executor.go` cubra todos los comandos shell esperados por los instaladores. |

## Rollback Plan
Revertir los commits relacionados con la adición de tests usando `git revert`.

## Dependencies
- `github.com/stretchr/testify` para aserciones y mocks.

## Success Criteria
- [ ] `go test ./...` compila y se ejecuta exitosamente.
- [ ] La cobertura global y por paquete (`openvpn`, `nvm`, `homebrew`) supera el 80%.