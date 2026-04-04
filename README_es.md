# 1X-SO Install Orchestrator (1X-SO-INSTALL)

## Descripción
Este proyecto es una herramienta de CLI (Command Line Interface) robusta y determinista diseñada para la configuración post-instalación de sistemas operativos basados en Debian. Proporciona una interfaz de usuario interactiva (TUI) que guía al usuario a través de la selección e instalación de diversos paquetes y configuraciones de software de manera modular e idempotente.

## Características Principales
- **Interfaz TUI Interactiva:** Basada en `bubbletea`, ofrece una experiencia de usuario fluida y visualmente atractiva en la terminal.
- **Validación de Entorno:** Comprobación automática de privilegios de root/sudo y detección de la distribución y versión de Debian.
- **Instalación Modular:** Selección múltiple de paquetes y configuraciones.
- **Idempotencia:** Detecta si un paquete ya está instalado para evitar fallos innecesarios.
- **Logging:** Registro detallado de la ejecución en `so-install.log` para depuración y transparencia.
- **Desarrollo Robusto:** Construido siguiendo principios de Clean Architecture y con una cobertura de tests guiada por TDD.

## Qué instala y configura
La herramienta permite instalar y configurar una amplia gama de software, categorizado de la siguiente manera:

### Navegadores
- **Brave:** Instalación desde el repositorio oficial.
- **Google Chrome:** Descarga e instalación del `.deb` oficial.
- **Firefox:** Repositorio oficial de Mozilla (última versión estable).
- **Chromium:** Repositorio oficial de Debian.

### Desarrollo y Herramientas
- **Docker:** Motor de contenedores.
- **DDEV:** Entorno de desarrollo local basado en Docker.
- **Ollama:** Ejecución local de modelos de lenguaje (LLM).
- **VS Code (OpenCode):** Fork de VS Code con configuraciones optimizadas.
- **NVM & NPM:** Gestión de versiones de Node.js y paquetes globales.
- **Homebrew:** Gestor de paquetes alternativo.

### Infraestructura y Sistema
- **APT:** Actualización del sistema y gestión de dependencias básicas.
- **Flatpak:** Sistema de gestión de paquetes universales.
- **OpenVPN:** Configuración de cliente VPN.
- **GitLab:** Token y conflagración de credenciales.
- **Desktop/Screen Lock:** Configuración del bloqueo de pantalla y preferencias de escritorio.

### IA y Agentes
- **Gentle AI:** Instalación y configuración del agente Gemini CLI.

## Requisitos del Sistema
Para compilar y ejecutar este proyecto, necesitarás:
- **Golang:** Versión 1.24.2 o superior.
- **Make:** Para la gestión de tareas de compilación.
- **Sistema Operativo:** Distribución basada en Debian (Debian 12/13 recomendados).
- **Privilegios:** Ejecución con `sudo` o como usuario `root`.

## Uso

### Compilación
Para compilar el proyecto y generar el binario en la carpeta `bin/`:
```bash
make build
```

### Ejecución
Para compilar (si es necesario) y ejecutar la aplicación con privilegios de root:
```bash
make run
```

O directamente si ya tienes el binario:
```bash
sudo ./bin/so-install
```

### Tests
Para ejecutar la suite de pruebas:
```bash
make test
```

Para realizar comprobaciones de linting:
```bash
make lint
```

## Limpieza
Para eliminar los binarios generados y archivos de cobertura:
```bash
make clean
```

## Solución de problemas
Si el proceso de instalación parece congelarse o falla, consulta el archivo `so-install.log` en el directorio actual para obtener información detallada sobre los comandos ejecutados y los errores encontrados.
