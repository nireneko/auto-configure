# 1X-SO Install Orchestrator (1X-SO-INSTALL)

## Description
This project is a robust and deterministic CLI (Command Line Interface) tool designed for post-installation configuration of Debian-based operating systems. It provides an interactive User Interface (TUI) that guides the user through the selection and installation of various software packages and configurations in a modular and idempotent manner.

## Key Features
- **Interactive TUI:** Built on `bubbletea`, it offers a smooth and visually appealing terminal experience.
- **Environment Validation:** Automatic checking for root/sudo privileges and detection of Debian distribution and version.
- **Modular Installation:** Multiple selection of packages and configurations.
- **Idempotency:** Detects if a package is already installed to avoid unnecessary failures.
- **Robust Development:** Built following Clean Architecture principles and guided by TDD (Test Driven Development).

## What it Installs and Configures
The tool allows you to install and configure a wide range of software, categorized as follows:

### Browsers
- **Brave:** Installation from the official repository.
- **Google Chrome:** Download and installation of the official `.deb` package.
- **Firefox:** Official Mozilla repository (latest stable version).
- **Chromium:** Official Debian repository.

### Development & Tools
- **Docker:** Container engine.
- **DDEV:** Local development environment based on Docker.
- **Ollama:** Local execution of LLMs (Large Language Models).
- **VS Code (OpenCode):** Fork of VS Code with optimized configurations.
- **NVM & NPM:** Node.js version management and global packages.
- **Homebrew:** Alternative package manager.

### Infrastructure & System
- **APT:** System update and basic dependency management.
- **Flatpak:** Universal package management system.
- **OpenVPN:** VPN client configuration.
- **GitLab:** Token and credential configuration.
- **Desktop/Screen Lock:** Screen lock configuration and desktop preferences.

### AI & Agents
- **Gentle AI:** Installation and configuration of the Gemini CLI agent.

## System Requirements
To compile and run this project, you will need:
- **Golang:** Version 1.24.2 or higher.
- **Make:** For compilation and task management.
- **Operating System:** Debian-based distribution (Debian 12/13 recommended).
- **Privileges:** Run with `sudo` or as `root` user.

## Usage

### Compilation
To compile the project and generate the binary in the `bin/` folder:
```bash
make build
```

### Running
To compile (if necessary) and run the application with root privileges:
```bash
make run
```

Alternatively, if you already have the binary:
```bash
sudo ./bin/so-install
```

### Testing
To run the test suite:
```bash
make test
```

To perform linting checks:
```bash
make lint
```

## Cleaning
To remove generated binaries and coverage files:
```bash
make clean
```
