package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/so-install/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestModel_Coverage(t *testing.T) {
	m := NewModel(map[domain.SoftwareID]domain.SoftwareInstaller{}, nil)

	assert.Nil(t, m.Init())
	assert.Equal(t, 0, m.ExitCode())

	osInfo := &domain.OSInfo{ID: "debian", VersionID: "12"}
	m.SetOSInfo(osInfo)
	assert.Equal(t, osInfo, m.osInfo)

	m.SetCursor(2)
	assert.Equal(t, 2, m.cursor)

	// Test Update WindowSizeMsg
	var updated tea.Model
	updated, _ = m.Update(tea.WindowSizeMsg{Width: 200, Height: 40})
	m = updated.(Model)
	assert.Equal(t, 100, m.width) // clamped
	assert.Equal(t, 40, m.height)

	// Test quit keys in welcome
	m.state = stateWelcome
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	assert.NotNil(t, cmd)

	// Test Token Input state
	m.state = stateTokenInput
	m.gitlabToken = ""
	// Type some characters
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("abc")})
	m = updated.(Model)
	assert.Equal(t, "abc", m.gitlabToken)
	// Backspace
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	m = updated.(Model)
	assert.Equal(t, "ab", m.gitlabToken)
	// Esc goes back to select
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(Model)
	assert.Equal(t, stateSoftwareSelect, m.state)

	// Test Progress state interrupt
	m.state = stateProgress
	updated, cmd = m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	m = updated.(Model)
	assert.True(t, m.interrupted)
	assert.Equal(t, 1, m.exitCode)
	assert.NotNil(t, cmd)

	// Test Summary state quit
	m.state = stateSummary
	_, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.NotNil(t, cmd)

	// Test Views
	m.state = stateWelcome
	assert.Contains(t, m.View(), "1x-so-install")

	m.state = stateSoftwareSelect
	assert.Contains(t, m.View(), "Select software to install:")

	m.state = stateTokenInput
	assert.Contains(t, m.View(), "Gitlab Token Configuration:")

	m.state = stateProgress
	assert.Contains(t, m.View(), "Finalizing installation summary") // if currentStep >= len(steps)

	m.state = stateSummary
	assert.Contains(t, m.View(), "Installation complete!")

	// Test StepFinishedMsg with critical failure
	m.state = stateProgress
	m.steps = []domain.InstallStep{{Critical: true}}
	m.currentStep = 0
	msg := StepFinishedMsg{
		Step: m.steps[0],
		Results: []domain.InstallResult{
			{Err: assert.AnError},
		},
	}
	updated, cmd = m.Update(msg)
	m = updated.(Model)
	assert.Equal(t, stateSummary, m.state)
}
