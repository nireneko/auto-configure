# PRD: 1X-SO Install Orchestrator (1X-SO-INSTALL)

## 1. Introducción y Propósito
El objetivo de este proyecto es desarrollar una herramienta de CLI (Command Line Interface) robusta y determinista para la configuración post-instalación de sistemas operativos basados en Debian. La aplicación guiará al usuario a través de una TUI interactiva para seleccionar paquetes y configuraciones de forma modular.

## 2. Requerimientos Funcionales

### 2.1. Validación de Entorno (Pre-check)
- **Privilegios:** La aplicación debe abortar inmediatamente si no se ejecuta con privilegios de `root` o `sudo`.
- **Detección de OS:** Identificar la distribución (Debian) y la versión específica (12, 13, etc.) utilizando `/etc/os-release`.

### 2.2. Interfaz de Usuario (TUI)
- Basada en el framework **Bubble Tea** (The Elm Architecture).
- Selección múltiple de opciones (Checkboxes).
- Flujo secuencial de pasos (Wizards).
- Feedback visual del progreso de instalación.

### 2.3. Módulo de Navegadores (Fase 1)
El sistema debe permitir la instalación de los siguientes navegadores desde sus fuentes oficiales:

- **Brave:** Repositorio oficial (GPG key + sources list).
- **Google Chrome:** Descarga del `.deb` oficial de Google e instalación vía `apt`.
- **Firefox (Official):** Utilizando el repositorio APT de Mozilla (preferido sobre la versión ESR de Debian para tener la última versión estable).
- **Chromium:** Instalación directa desde los repositorios oficiales de Debian.

### 2.4. Motor de Ejecución
- Capacidad para ejecutar comandos de sistema de forma segura.
- Manejo de logs de salida para cada comando.
- Idempotencia: Si un paquete ya está instalado, el sistema debe informarlo y no fallar.

## 3. Requerimientos Técnicos

- **Lenguaje:** Golang.
- **TUI Framework:** `charmbracelet/bubbletea`, `charmbracelet/lipgloss` para estilos.
- **Testing Strategy:** 
  - **100% TDD:** Todo el código debe ser guiado por pruebas.
  - **Abstracción de Shell:** Uso de interfaces para el ejecutor de comandos (`Executor` interface) para permitir mocks en tests unitarios.
  - **Inyección de Dependencias:** Para facilitar la testabilidad de los estados de la TUI.

## 4. Estrategia de Instalación de Navegadores (Detalle Técnico)

### Google Chrome
1. Descargar: `wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb -P /tmp/`
2. Instalar: `apt install /tmp/google-chrome-stable_current_amd64.deb`

### Firefox (Mozilla Official Repo)
1. Importar GPG: `wget -q https://packages.mozilla.org/apt/repo-signing-key.gpg -O- | tee /etc/apt/keyrings/packages.mozilla.org.gpg > /dev/null`
2. Añadir Repo: `echo "deb [signed-by=/etc/apt/keyrings/packages.mozilla.org.gpg] https://packages.mozilla.org/apt mozilla main" | tee /etc/apt/sources.list.d/mozilla.list`
3. Priorizar Repo Mozilla: Configurar `/etc/apt/preferences.d/mozilla` para asegurar que el paquete de Mozilla tenga prioridad.
4. Instalar: `apt update && apt install firefox`

### Chromium
1. Instalar: `apt update && apt install chromium`

### Brave (Referencia confirmada)
1. GPG: `curl -fsSLo /usr/share/keyrings/brave-browser-archive-keyring.gpg https://brave-browser-apt-release.s3.brave.com/brave-browser-archive-keyring.gpg`
2. Repo: `curl -fsSLo /etc/apt/sources.list.d/brave-browser-release.sources https://brave-browser-apt-release.s3.brave.com/brave-browser.sources`
3. Instalar: `apt update && apt install brave-browser`

## 5. Roadmap
1. **Sprint 0:** Estructura base, detección de OS y validación de Root (TDD).
2. **Sprint 1:** Implementación de la TUI base (Bubble Tea).
3. **Sprint 2:** Módulo de instalación de Navegadores (Lógica de Shell).
4. **Sprint 3:** Integración TUI + Engine y Feedback Visual.
