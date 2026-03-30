package tui

import (
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/core/usecases"
)

type appState int

const (
	stateWelcome appState = iota
	stateBrowserSelect
	stateProgress
	stateSummary
)

// Model is the root Bubbletea model for 1x-so-install.
type Model struct {
	state      appState
	osInfo     *domain.OSInfo
	installers map[domain.BrowserID]domain.BrowserInstaller
	selected   []domain.BrowserID
	results    []domain.InstallResult
	exitCode   int
	width      int
	height     int

	// sub-model state
	cursor         int
	checked        map[domain.BrowserID]bool
	preInstalled   map[domain.BrowserID]bool
	validationErr  string
	currentInstall int
	interrupted    bool
}

// NewModel creates the TUI model with the given installers.
func NewModel(installers map[domain.BrowserID]domain.BrowserInstaller) Model {
	return Model{
		state:        stateWelcome,
		installers:   installers,
		checked:      make(map[domain.BrowserID]bool),
		preInstalled: make(map[domain.BrowserID]bool),
		width:        80,
		height:       24,
	}
}

// SetOSInfo sets the OS information on the model (called from main before TUI launch).
func (m *Model) SetOSInfo(info *domain.OSInfo) {
	m.osInfo = info
}

// ExitCode returns the process exit code (0 = success, 1 = failure).
func (m Model) ExitCode() int {
	return m.exitCode
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		if m.width > 100 {
			m.width = 100
		}
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case OSDetectedMsg:
		m.osInfo = msg.Info
		return m, m.checkInstalledBrowsers()

	case preInstalledCheckDoneMsg:
		m.preInstalled = msg.results
		for id, installed := range msg.results {
			m.checked[id] = installed
		}
		m.state = stateBrowserSelect
		return m, nil

	case startInstallMsg:
		m.state = stateProgress
		m.currentInstall = 0
		return m, m.runInstallations()

	case InstallProgressMsg:
		m.results = append(m.results, msg.Result)
		m.currentInstall++
		return m, nil

	case AllInstallsDoneMsg:
		m.results = msg.Results
		m.state = stateSummary
		for _, r := range m.results {
			if r.Err != nil {
				m.exitCode = 1
				break
			}
		}
		return m, nil
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case stateWelcome:
		switch msg.String() {
		case "enter", " ":
			return m, m.detectOSCmd()
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case stateBrowserSelect:
		browsers := domain.AllBrowsers()
		switch msg.String() {
		case "up", "k":
			m.cursor = (m.cursor - 1 + len(browsers)) % len(browsers)
		case "down", "j":
			m.cursor = (m.cursor + 1) % len(browsers)
		case " ":
			id := browsers[m.cursor]
			m.checked[id] = !m.checked[id]
			m.validationErr = ""
		case "enter":
			sel := m.getSelected()
			if len(sel) == 0 {
				m.validationErr = "Select at least one browser"
				return m, nil
			}
			m.selected = sel
			return m, func() tea.Msg { return startInstallMsg{} }
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case stateProgress:
		switch msg.String() {
		case "ctrl+c":
			m.interrupted = true
			m.exitCode = 1
			return m, tea.Quit
		}

	case stateSummary:
		switch msg.String() {
		case "q", "enter", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	switch m.state {
	case stateWelcome:
		return m.viewWelcome()
	case stateBrowserSelect:
		return m.viewBrowserSelect()
	case stateProgress:
		return m.viewProgress()
	case stateSummary:
		return m.viewSummary()
	}
	return ""
}

func (m Model) getSelected() []domain.BrowserID {
	var sel []domain.BrowserID
	for _, id := range domain.AllBrowsers() {
		if m.checked[id] {
			sel = append(sel, id)
		}
	}
	return sel
}

// detectOSCmd emits osInfo (already set by main.go) as an OSDetectedMsg to trigger
// browser pre-install check and transition to browser select state.
func (m Model) detectOSCmd() tea.Cmd {
	info := m.osInfo
	return func() tea.Msg {
		return OSDetectedMsg{Info: info}
	}
}

type preInstalledCheckDoneMsg struct {
	results map[domain.BrowserID]bool
}

type startInstallMsg struct{}

func (m Model) checkInstalledBrowsers() tea.Cmd {
	installers := m.installers
	return func() tea.Msg {
		results := make(map[domain.BrowserID]bool)
		for id, inst := range installers {
			installed, _ := inst.IsInstalled()
			results[id] = installed
		}
		return preInstalledCheckDoneMsg{results: results}
	}
}

func (m Model) runInstallations() tea.Cmd {
	selected := m.selected
	installers := m.installers
	return func() tea.Msg {
		uc := usecases.NewInstallBrowsersUseCase(installers, time.Sleep)
		results := uc.Execute(selected)
		return AllInstallsDoneMsg{Results: results}
	}
}

// --- Views ---

func (m Model) viewWelcome() string {
	out := "\n"
	out += "  1x-so-install\n"
	out += "  Post-installation OS configurator\n"
	if m.osInfo != nil {
		out += "\n  OS: " + m.osInfo.ID + " " + m.osInfo.VersionID + "\n"
	}
	out += "\n  Press Enter to continue  •  q to quit\n"
	return out
}

func (m Model) viewBrowserSelect() string {
	out := "\n  Select browsers to install:\n\n"
	for i, id := range domain.AllBrowsers() {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}
		checkbox := "[ ]"
		if m.checked[id] {
			checkbox = "[x]"
		}
		label := id.DisplayName()
		if m.preInstalled[id] {
			label += " (installed)"
		}
		out += cursor + checkbox + " " + label + "\n"
	}
	if m.validationErr != "" {
		out += "\n  ! " + m.validationErr + "\n"
	}
	out += "\n  Space: toggle  •  Enter: confirm  •  q: quit\n"
	return out
}

func (m Model) viewProgress() string {
	out := "\n  Installing browsers...\n\n"
	for i, id := range m.selected {
		var status string
		if i < len(m.results) {
			r := m.results[i]
			if r.Err != nil {
				status = "  [✗] " + id.DisplayName() + " — Failed"
			} else if r.AlreadyInstalled {
				status = "  [✓] " + id.DisplayName() + " — Already installed"
			} else {
				status = "  [✓] " + id.DisplayName() + " — Installed"
			}
		} else if i == len(m.results) {
			status = "  [~] " + id.DisplayName() + " — Installing..."
		} else {
			status = "  [ ] " + id.DisplayName()
		}
		out += status + "\n"
	}
	out += "\n  Ctrl+C to abort\n"
	return out
}

func (m Model) viewSummary() string {
	out := "\n  Installation complete!\n\n"
	success, failed := 0, 0
	for _, r := range m.results {
		if r.Err != nil {
			out += "  [✗] " + r.Browser.DisplayName() + " — Failed\n"
			failed++
		} else if r.AlreadyInstalled {
			out += "  [✓] " + r.Browser.DisplayName() + " — Already installed\n"
			success++
		} else {
			out += "  [✓] " + r.Browser.DisplayName() + " — Installed\n"
			success++
		}
	}
	out += "\n"
	out += "  Installed: " + strconv.Itoa(success) + "  Failed: " + strconv.Itoa(failed) + "\n"
	out += "\n  Press Enter or q to exit\n"
	return out
}
