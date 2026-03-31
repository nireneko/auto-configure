# Tasks: 80% Test Coverage & Edge Cases

## Phase 1: Fix Compilation Errors
- [x] 1.1 Eliminar imports no usados en `update_ids.go` (`os`, `path/filepath`, `strings`).
- [x] 1.2 Remover o usar correctamente la variable `exitCode` en `cmd/so-install/main_test.go`.
- [x] 1.3 Corregir llamadas a `m.Update` en `internal/presentation/tui/model_extra_test.go` para usar asignación múltiple `m, cmd = m.Update(...)`.

## Phase 2: OpenVPN Tests (Edge Cases)
- [x] 2.1 Agregar test en `internal/infrastructure/openvpn/openvpn_test.go` simulando error de red (ej. fallo en wget de keyrings) usando `MockExecutor`.
- [x] 2.2 Agregar test simulando fallo en `apt update` al instalar OpenVPN.

## Phase 3: NVM Tests (Edge Cases)
- [x] 3.1 Agregar test en `internal/infrastructure/nvm/nvm_test.go` simulando error en la descarga del script de instalación.
- [x] 3.2 Agregar test simulando fallo en la ejecución de la carga de NVM en el profile (`source ~/.nvm/nvm.sh`).

## Phase 4: Homebrew Tests (Edge Cases)
- [x] 4.1 Agregar test en `internal/infrastructure/homebrew/homebrew_test.go` para verificar comportamiento cuando `brew` no está en el sistema y la instalación inicial falla.

## Phase 5: Verification
- [x] 5.1 Ejecutar `go test ./... -coverprofile=coverage.out` y confirmar compilación limpia.
- [x] 5.2 Ejecutar `go tool cover -func=coverage.out` y verificar que openvpn, nvm y homebrew superan 80%.