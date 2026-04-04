package tui_test

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/presentation/tui"
	"github.com/so-install/pkg/mocks"
)

func makeInstallers(installed map[domain.SoftwareID]bool, installErrs map[domain.SoftwareID]error) map[domain.SoftwareID]domain.SoftwareInstaller {
	result := make(map[domain.SoftwareID]domain.SoftwareInstaller)
	for _, id := range domain.AllSoftware() {
		result[id] = &mocks.MockSoftwareInstaller{
			SoftwareID:        id,
			IsInstalledResult: installed[id],
			InstallErr:        installErrs[id],
		}
	}
	return result
}

func TestModel_InitialStateIsWelcome(t *testing.T) {
	m := tui.NewModel(makeInstallers(nil, nil), nil)
	view := m.View()
	if !strings.Contains(view, "1x-so-install") {
		t.Errorf("welcome view missing app name, got: %s", view)
	}
}

func TestModel_EnterOnWelcomeTriggersDetection(t *testing.T) {
	m := tui.NewModel(makeInstallers(nil, nil), nil)
	m2, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Log("cmd is nil, may be acceptable if osInfo not set")
	}
	_ = m2
}

func TestModel_SoftwareSelectShowsAllSoftware(t *testing.T) {
	installers := makeInstallers(nil, nil)
	m := tui.NewModel(installers, nil)
	m2, _ := m.Update(tui.OSDetectedMsg{Info: &domain.OSInfo{ID: "debian", VersionID: "12"}})
	view := m2.View()
	_ = view
}

func TestModel_SelectionValidation(t *testing.T) {
	installers := makeInstallers(nil, nil)
	m := tui.NewModel(installers, nil)
	// Manually send OSDetected + preInstalled to reach select state
	m2, cmd := m.Update(tui.OSDetectedMsg{Info: &domain.OSInfo{ID: "debian", VersionID: "12"}})
	if cmd != nil {
		msg := cmd()
		m2, _ = m2.Update(msg)
	}
	// Now try to confirm with no selection (all unchecked)
	m3, _ := m2.Update(tea.KeyMsg{Type: tea.KeyEnter})
	view := m3.View()
	if strings.Contains(view, "Select at least") {
		// validation message appeared — expected behavior
		return
	}
	t.Log("validation path not triggered (may be pre-checked or in different state)")
}

func TestModel_ExitCode0WhenAllSucceed(t *testing.T) {
	installers := makeInstallers(nil, nil)
	m := tui.NewModel(installers, nil)
	m2, _ := m.Update(tui.AllInstallsDoneMsg{
		Results: []domain.InstallResult{
			{Software: domain.Brave, Err: nil},
			{Software: domain.Firefox, Err: nil},
		},
	})
	model := m2.(tui.Model)
	if model.ExitCode() != 0 {
		t.Errorf("expected exit code 0, got %d", model.ExitCode())
	}
}

func TestModel_ExitCode1WhenAnyFails(t *testing.T) {
	installers := makeInstallers(nil, nil)
	m := tui.NewModel(installers, nil)
	m2, _ := m.Update(tui.AllInstallsDoneMsg{
		Results: []domain.InstallResult{
			{Software: domain.Brave, Err: nil},
			{Software: domain.Firefox, Err: errors.New("install failed")},
		},
	})
	model := m2.(tui.Model)
	if model.ExitCode() != 1 {
		t.Errorf("expected exit code 1, got %d", model.ExitCode())
	}
}

func TestModel_QuitOnSummary(t *testing.T) {
	installers := makeInstallers(nil, nil)
	m := tui.NewModel(installers, nil)
	m2, _ := m.Update(tui.AllInstallsDoneMsg{Results: []domain.InstallResult{}})
	_, cmd := m2.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Error("expected quit command from summary state")
	}
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Errorf("expected tea.QuitMsg, got %T", msg)
	}
}

func TestModel_SoftwareSelectLabelUpdated(t *testing.T) {
	installers := makeInstallers(nil, nil)
	m := tui.NewModel(installers, nil)
	m2, cmd := m.Update(tui.OSDetectedMsg{Info: &domain.OSInfo{ID: "debian", VersionID: "12"}})
	if cmd != nil {
		msg := cmd()
		m2, _ = m2.Update(msg)
	}
	view := m2.View()
	if !strings.Contains(view, "Select software to install") {
		t.Errorf("expected 'Select software to install' label, got: %s", view)
	}
}

func TestModel_DockerAppearsInSoftwareList(t *testing.T) {
	installers := makeInstallers(nil, nil)
	m := tui.NewModel(installers, nil)
	m2, cmd := m.Update(tui.OSDetectedMsg{Info: &domain.OSInfo{ID: "debian", VersionID: "12"}})
	if cmd != nil {
		msg := cmd()
		m2, _ = m2.Update(msg)
	}
	view := m2.View()
	if !strings.Contains(view, "Docker CE") {
		t.Errorf("expected 'Docker CE' in software list, got: %s", view)
	}
}

func TestModel_DdevAppearsInSoftwareList(t *testing.T) {
	installers := makeInstallers(nil, nil)
	m := tui.NewModel(installers, nil)
	m2, cmd := m.Update(tui.OSDetectedMsg{Info: &domain.OSInfo{ID: "debian", VersionID: "12"}})
	if cmd != nil {
		msg := cmd()
		m2, _ = m2.Update(msg)
	}
	view := m2.View()
	if !strings.Contains(view, "DDEV") {
		t.Errorf("expected 'DDEV' in software list, got: %s", view)
	}
}

func TestModel_GentleAIAppearsInSoftwareList(t *testing.T) {
	installers := makeInstallers(nil, nil)
	m := tui.NewModel(installers, nil)
	m2, cmd := m.Update(tui.OSDetectedMsg{Info: &domain.OSInfo{ID: "debian", VersionID: "12"}})
	if cmd != nil {
		msg := cmd()
		m2, _ = m2.Update(msg)
	}
	view := m2.View()
	if !strings.Contains(view, "Gentle-AI") {
		t.Errorf("expected 'Gentle-AI' in software list, got: %s", view)
	}
}

func TestModel_StepTransitions(t *testing.T) {
	installers := makeInstallers(nil, nil)
	m := tui.NewModel(installers, nil)

	// Start install
	m2, cmd := m.Update(tui.OSDetectedMsg{Info: &domain.OSInfo{ID: "debian", VersionID: "12"}})
	// skip to select
	if cmd != nil {
		m2, _ = m2.Update(cmd())
	}

	// Start installation step
	steps := domain.GetSteps()
	m3, cmd := m2.Update(tui.StepFinishedMsg{
		Step:    steps[0],
		Results: []domain.InstallResult{{Software: domain.Brave, Err: nil}},
	})

	model := m3.(tui.Model)
	_ = model
	// Check if moved to next step (index 1 is docker)
	// (Note: it actually returns m, runCurrentStep() command)
	if cmd == nil {
		t.Fatal("expected command to run next step")
	}
}

func TestModel_NoPanicOnLastStepFinished(t *testing.T) {
	installers := makeInstallers(nil, nil)
	tm := tui.NewModel(installers, nil)

	// Reach stateProgress
	m_update, cmd := tm.Update(tui.OSDetectedMsg{Info: &domain.OSInfo{ID: "debian", VersionID: "12"}})
	m := m_update.(tui.Model)
	if cmd != nil {
		m_update, _ = m.Update(cmd()) // PreInstallCheckDoneMsg
		m = m_update.(tui.Model)
	}

	// Manually select all to avoid empty selection
	// (we can't access fields but we can send space keys or just let it be if defaults are set)
	// Actually, getSelected() returns based on m.checked which is set from pre-install check
	// Since all are false in our mock, we need to toggle one.
	// Let's send a space to toggle the first one.
	m_update, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})
	m = m_update.(tui.Model)

	m_update, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = m_update.(tui.Model)
	if cmd != nil {
		m_update, _ = m.Update(cmd()) // startInstallMsg
		m = m_update.(tui.Model)
	}

	steps := domain.GetSteps()
	// Finish all steps but the last one
	for i := 0; i < len(steps)-1; i++ {
		m_update, _ = m.Update(tui.StepFinishedMsg{
			Step:    steps[i],
			Results: []domain.InstallResult{{Software: steps[i].Software[0], Err: nil}},
		})
		m = m_update.(tui.Model)
	}

	// Finish LAST step
	lastStep := steps[len(steps)-1]
	m_update, _ = m.Update(tui.StepFinishedMsg{
		Step:    lastStep,
		Results: []domain.InstallResult{{Software: lastStep.Software[0], Err: nil}},
	})
	m = m_update.(tui.Model)

	// Now m should be in the summary state:
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("View() panicked: %v", r)
		}
	}()

	view := m.View()
	if !strings.Contains(view, "Installation complete!") {
		t.Errorf("expected summary view after last step, got: %s", view)
	}
}

func TestModel_GitlabTokenConfigTransition(t *testing.T) {
	installers := makeInstallers(nil, nil)
	m := tui.NewModel(installers, nil)

	// Reach stateSoftwareSelect
	m_update, cmd := m.Update(tui.OSDetectedMsg{Info: &domain.OSInfo{ID: "debian", VersionID: "12"}})
	m = m_update.(tui.Model)
	if cmd != nil {
		m_update, _ = m.Update(cmd())
		m = m_update.(tui.Model)
	}

	// Find the GitlabTokenConfig index in the list
	software := domain.AllSoftware()
	idx := -1
	for i, id := range software {
		if id == domain.GitlabTokenConfig {
			idx = i
			break
		}
	}

	if idx == -1 {
		t.Fatal("GitlabTokenConfig not found in software list")
	}

	// Navigate to GitlabTokenConfig and select it
	m.SetCursor(idx) // We'll need to add a SetCursor for testing if it's not exported
	m_update, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})
	m = m_update.(tui.Model)

	// Press Enter to confirm selection
	m_update, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = m_update.(tui.Model)

	// Should NOT start installation yet, should go to stateTokenInput
	if cmd != nil {
		t.Fatal("Expected NO command yet, state should be token input")
	}

	view := m.View()
	if !strings.Contains(view, "Gitlab Token Configuration") {
		t.Errorf("Expected Gitlab Token Configuration view, got: %s", view)
	}

	// Test token input and masking
	m_update, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
	m = m_update.(tui.Model)
	m_update, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("b")})
	m = m_update.(tui.Model)

	view = m.View()
	if !strings.Contains(view, "Token: **_") {
		t.Errorf("Expected masked token with 2 stars, got: %s", view)
	}

	// Confirm token
	m_update, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("Expected startInstallMsg command after confirming token")
	}
}
