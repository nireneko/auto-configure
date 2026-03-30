package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/core/usecases"
	"github.com/so-install/internal/infrastructure/browsers"
	"github.com/so-install/internal/infrastructure/ddev"
	"github.com/so-install/internal/infrastructure/docker"
	"github.com/so-install/internal/infrastructure/osrelease"
	"github.com/so-install/internal/infrastructure/shell"
	"github.com/so-install/internal/presentation/tui"
)

func main() {
	// 1. Validate privileges before anything else
	privUC := usecases.NewCheckPrivilegesUseCase(os.Getuid, os.Getenv)
	if err := privUC.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		fmt.Fprintf(os.Stderr, "Please run with: sudo 1x-so-install\n")
		os.Exit(1)
	}

	// 2. Detect OS
	detector := osrelease.NewDefaultDetector()
	osUC := usecases.NewDetectOSUseCase(detector)
	osInfo, err := osUC.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	// 3. Build concrete installers
	executor := shell.NewShellExecutor()
	installerMap := map[domain.SoftwareID]domain.SoftwareInstaller{
	        domain.Brave:    browsers.NewBraveInstaller(executor),
	        domain.Firefox:  browsers.NewFirefoxInstaller(executor),
	        domain.Chrome:   browsers.NewChromeInstaller(executor),
	        domain.Chromium: browsers.NewChromiumInstaller(executor),
	        domain.Docker:   docker.NewDockerInstaller(executor, os.Getenv("SUDO_USER")),
	        domain.Ddev:     ddev.NewDdevInstaller(executor),
	}
	// 4. Build TUI model and inject osInfo
	model := tui.NewModel(installerMap)
	model.SetOSInfo(osInfo)

	// 5. Run TUI
	p := tea.NewProgram(model, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %s\n", err.Error())
		os.Exit(1)
	}

	if m, ok := finalModel.(tui.Model); ok {
		os.Exit(m.ExitCode())
	}
}
