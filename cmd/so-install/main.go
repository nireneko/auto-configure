package main

import (
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/core/usecases"
	"github.com/so-install/internal/infrastructure/antigravity"
	"github.com/so-install/internal/infrastructure/apt"
	"github.com/so-install/internal/infrastructure/browsers"
	"github.com/so-install/internal/infrastructure/cursor"
	"github.com/so-install/internal/infrastructure/ddev"
	"github.com/so-install/internal/infrastructure/desktop"
	"github.com/so-install/internal/infrastructure/docker"
	"github.com/so-install/internal/infrastructure/flatpak"
	"github.com/so-install/internal/infrastructure/gentleai"
	"github.com/so-install/internal/infrastructure/gitlab"
	"github.com/so-install/internal/infrastructure/homebrew"
	"github.com/so-install/internal/infrastructure/npm"
	"github.com/so-install/internal/infrastructure/nvidia"
	"github.com/so-install/internal/infrastructure/nvm"
	"github.com/so-install/internal/infrastructure/ollama"
	"github.com/so-install/internal/infrastructure/opencode"
	"github.com/so-install/internal/infrastructure/openvpn"
	"github.com/so-install/internal/infrastructure/osrelease"
	"github.com/so-install/internal/infrastructure/shell"
	"github.com/so-install/internal/infrastructure/vscode"
	"github.com/so-install/internal/infrastructure/logging"
	"github.com/so-install/internal/presentation/tui"
	)
var (
	osGetuid   = os.Getuid
	osGetenv   = os.Getenv
	runProgram = func(p *tea.Program) (tea.Model, error) { return p.Run() }
	newDetector = func() domain.OSDetector { return osrelease.NewDefaultDetector() }
)

func Run(args []string, out io.Writer, errOut io.Writer) int {
        // Setup logging
        logFile, err := os.OpenFile("so-install.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
        var logger domain.Logger
        if err != nil {
                fmt.Fprintf(errOut, "Warning: could not open log file: %s. Logging disabled.\n", err.Error())
                logger = domain.NoopLogger{}
        } else {
                defer logFile.Close()
                logger = logging.NewFileLogger(logFile)
        }

        logger.Info("application starting")

        privUC := usecases.NewCheckPrivilegesUseCase(osGetuid, osGetenv)
        if err := privUC.Execute(); err != nil {
                logger.Error("privilege check failed", "err", err)
                fmt.Fprintf(errOut, "Error: %s\n", err.Error())
                fmt.Fprintf(errOut, "Please run with: sudo 1x-so-install\n")
                return 1
        }

        detector := newDetector()
        osUC := usecases.NewDetectOSUseCase(detector)
        osInfo, err := osUC.Execute()
        if err != nil {
                logger.Error("OS detection failed", "err", err)
                fmt.Fprintf(errOut, "Error: %s\n", err.Error())
                return 1
        }

        executor := shell.NewShellExecutor(logger)
        installerMap := map[domain.SoftwareID]domain.SoftwareInstaller{
                domain.SystemUpdate:      apt.NewAptUpdateInstaller(executor),
                domain.BaseDeps:          apt.NewBaseDepsInstaller(executor),
                domain.Brave:             browsers.NewBraveInstaller(executor),
                domain.Firefox:           browsers.NewFirefoxInstaller(executor),
                domain.Chrome:            browsers.NewChromeInstaller(executor),
                domain.Chromium:          browsers.NewChromiumInstaller(executor),
                domain.Docker:            docker.NewDockerInstaller(executor, osGetenv("SUDO_USER")),
                domain.Ddev:              ddev.NewDdevInstaller(executor),
                domain.OpenVpn:           openvpn.NewOpenVpnInstaller(executor, osInfo),
                domain.Nvm:               nvm.NewNvmInstaller(executor),
                domain.Gemini:            npm.NewNpmInstaller(executor, "@google/gemini-cli", "gemini", domain.Gemini),
                domain.ClaudeCode:        npm.NewNpmInstaller(executor, "@anthropic-ai/claude-code", "claude", domain.ClaudeCode),
                domain.Codex:             npm.NewNpmInstaller(executor, "@openai/codex", "codex", domain.Codex),
                domain.Ollama:            ollama.NewOllamaInstaller(executor),
                domain.OpenCode:          opencode.NewOpenCodeInstaller(executor),
                domain.GentleAI:          gentleai.NewGentleAIInstaller(executor),
                domain.VsCode:            vscode.NewVsCodeInstaller(executor),
                domain.Cursor:            cursor.NewCursorInstaller(executor),
                domain.Antigravity:       antigravity.NewAntigravityInstaller(executor),
                domain.Flatpak:           flatpak.NewFlatpakInstaller(executor, detector),
                domain.NvidiaDrivers:     nvidia.NewNvidiaInstaller(executor, osInfo),
                domain.Bitwarden:         flatpak.NewFlatpakAppInstaller(executor, "com.bitwarden.desktop", domain.Bitwarden),
                domain.Homebrew:          homebrew.NewHomebrewInstaller(executor),
                domain.GitlabTokenConfig: gitlab.NewGitlabTokenConfigurator(executor),
                domain.ScreenLockConfig:  desktop.NewScreenLockInstaller(executor, detector),
        }

        model := tui.NewModel(installerMap, logger)

	model.SetOSInfo(osInfo)

	p := tea.NewProgram(model, tea.WithAltScreen())
	finalModel, err := runProgram(p)
	if err != nil {
		fmt.Fprintf(errOut, "TUI error: %s\n", err.Error())
		return 1
	}

	if m, ok := finalModel.(tui.Model); ok {
		return m.ExitCode()
	}
	return 0
}

func main() {
	os.Exit(Run(os.Args, os.Stdout, os.Stderr))
}
