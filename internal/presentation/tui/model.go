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
	stateSoftwareSelect
	stateProgress
	stateSummary
)

// Model is the root Bubbletea model for 1x-so-install.
type Model struct {
	state      appState
	osInfo     *domain.OSInfo
	installers map[domain.SoftwareID]domain.SoftwareInstaller
	selected   []domain.SoftwareID
	results    []domain.InstallResult
	exitCode   int
	width      int
	height     int

	// sub-model state
	cursor        int
	checked       map[domain.SoftwareID]bool
	preInstalled  map[domain.SoftwareID]bool
	validationErr string
	interrupted   bool

	// step-based state
	steps       []domain.InstallStep
	currentStep int
}

// NewModel creates the TUI model with the given installers.
func NewModel(installers map[domain.SoftwareID]domain.SoftwareInstaller) Model {
	return Model{
		state:        stateWelcome,
		installers:   installers,
		checked:      make(map[domain.SoftwareID]bool),
		preInstalled: make(map[domain.SoftwareID]bool),
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
		return m, m.checkInstalledSoftware()

	case preInstalledCheckDoneMsg:
		m.preInstalled = msg.results
		for id, installed := range msg.results {
			m.checked[id] = installed
		}
		m.state = stateSoftwareSelect
		return m, nil

	case startInstallMsg:
		m.state = stateProgress
		m.steps = domain.GetSteps()
		m.currentStep = 0
		return m, m.runCurrentStep()

	case StepFinishedMsg:
		m.results = append(m.results, msg.Results...)

		// Check if we should stop due to critical failure
		if msg.Step.Critical {
			for _, r := range msg.Results {
				if r.Err != nil {
					m.state = stateSummary
					return m, func() tea.Msg { return AllInstallsDoneMsg{Results: m.results} }
				}
			}
		}

		m.currentStep++
		if m.currentStep >= len(m.steps) {
			m.state = stateSummary
			return m, func() tea.Msg { return AllInstallsDoneMsg{Results: m.results} }
		}
		return m, m.runCurrentStep()

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

	case stateSoftwareSelect:
		software := domain.AllSoftware()
		switch msg.String() {
		case "up", "k":
			m.cursor = (m.cursor - 1 + len(software)) % len(software)
		case "down", "j":
			m.cursor = (m.cursor + 1) % len(software)
		case " ":
			id := software[m.cursor]
			m.checked[id] = !m.checked[id]
			m.validationErr = ""
		case "enter":
			sel := m.getSelected()
			if len(sel) == 0 {
				m.validationErr = "Select at least one item"
				return m, nil
			}
			// Prepend mandatory system prep steps
			m.selected = append([]domain.SoftwareID{domain.SystemUpdate, domain.BaseDeps}, sel...)
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
	case stateSoftwareSelect:
		return m.viewSoftwareSelect()
	case stateProgress:
		return m.viewProgress()
	case stateSummary:
		return m.viewSummary()
	}
	return ""
}

func (m Model) getSelected() []domain.SoftwareID {
	var sel []domain.SoftwareID
	for _, id := range domain.AllSoftware() {
		if m.checked[id] {
			sel = append(sel, id)
		}
	}
	return sel
}

// detectOSCmd emits osInfo (already set by main.go) as an OSDetectedMsg to trigger
// software pre-install check and transition to software select state.
func (m Model) detectOSCmd() tea.Cmd {
	info := m.osInfo
	return func() tea.Msg {
		return OSDetectedMsg{Info: info}
	}
}

type preInstalledCheckDoneMsg struct {
	results map[domain.SoftwareID]bool
}

type startInstallMsg struct{}

func (m Model) checkInstalledSoftware() tea.Cmd {
	installers := m.installers
	return func() tea.Msg {
		results := make(map[domain.SoftwareID]bool)
		for id, inst := range installers {
			installed, _ := inst.IsInstalled()
			results[id] = installed
		}
		return preInstalledCheckDoneMsg{results: results}
	}
}

func (m Model) runCurrentStep() tea.Cmd {
	step := m.steps[m.currentStep]

	// Find which software in this step is actually selected
	var stepSelected []domain.SoftwareID
	selectedMap := make(map[domain.SoftwareID]bool)
	for _, id := range m.selected {
		selectedMap[id] = true
	}
	for _, id := range step.Software {
		if selectedMap[id] {
			stepSelected = append(stepSelected, id)
		}
	}

	// If no software from this step is selected, skip to next step immediately
	if len(stepSelected) == 0 {
		return func() tea.Msg {
			return StepFinishedMsg{Step: step, Results: nil}
		}
	}

	installers := m.installers
	return func() tea.Msg {
		uc := usecases.NewInstallSoftwareUseCase(installers, time.Sleep)
		results := uc.Execute(stepSelected)
		return StepFinishedMsg{Step: step, Results: results}
	}
}

// --- Views ---

func (m Model) viewWelcome() string {
	out := "\n"
	out += "  1x-so-install\n"
	out += "  Post-installation OS configurator\n"
	if m.osInfo != nil {
		out += "\n  OS: " + m.osInfo.ID + " " + m.osInfo.VersionID
		if m.osInfo.DesktopEnvironment != "" && m.osInfo.DesktopEnvironment != domain.Other {
			out += " (" + string(m.osInfo.DesktopEnvironment) + ")"
		}
		out += "\n"
	}
	out += "\n  Press Enter to continue  •  q to quit\n"
	return out
}

func (m Model) viewSoftwareSelect() string {
	out := "\n  Select software to install:\n\n"
	for i, id := range domain.AllSoftware() {
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
	if m.currentStep >= len(m.steps) {
		return "\n  Finalizing installation summary...\n"
	}
	currentStep := m.steps[m.currentStep]
	out := "\n  Phase: " + currentStep.ID + "...\n\n"

	// Create map for easy lookup of results
	resultsMap := make(map[domain.SoftwareID]domain.InstallResult)
	for _, r := range m.results {
		resultsMap[r.Software] = r
	}

	for _, id := range m.selected {
		var status string
		r, done := resultsMap[id]
		if done {
			if r.Err != nil {
				status = "  [✗] " + id.DisplayName() + " — Failed"
			} else if r.AlreadyInstalled {
				status = "  [✓] " + id.DisplayName() + " — Already installed"
			} else {
				status = "  [✓] " + id.DisplayName() + " — Installed"
			}
		} else {
			// Check if it's currently installing (in the active step)
			inCurrentStep := false
			for _, sid := range currentStep.Software {
				if sid == id {
					inCurrentStep = true
					break
				}
			}
			if inCurrentStep {
				status = "  [~] " + id.DisplayName() + " — Installing..."
			} else {
				status = "  [ ] " + id.DisplayName()
			}
		}
		out += status + "\n"
	}
	out += "\n  Ctrl+C to abort\n"
	return out
}

func (m Model) viewSummary() string {
	out := "\n  Installation complete!\n\n"
	success, failed, skipped := 0, 0, 0

	// Track which ones were actually processed
	processed := make(map[domain.SoftwareID]bool)
	for _, r := range m.results {
		processed[r.Software] = true
		if r.Err != nil {
			out += "  [✗] " + r.Software.DisplayName() + " — Failed\n"
			failed++
		} else if r.AlreadyInstalled {
			out += "  [✓] " + r.Software.DisplayName() + " — Already installed\n"
			success++
		} else {
			out += "  [✓] " + r.Software.DisplayName() + " — Installed\n"
			success++
		}
	}

	// Any selected but not processed was skipped
	for _, id := range m.selected {
		if !processed[id] {
			out += "  [ ] " + id.DisplayName() + " — Skipped (dependency failed)\n"
			skipped++
		}
	}

	out += "\n"
	out += "  Installed: " + strconv.Itoa(success) + "  Failed: " + strconv.Itoa(failed)
	if skipped > 0 {
		out += "  Skipped: " + strconv.Itoa(skipped)
	}
	out += "\n"
	out += "\n  Press Enter or q to exit\n"
	return out
}
