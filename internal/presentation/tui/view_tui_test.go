package tui_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/presentation/tui"
)

func TestModel_ViewMethods(t *testing.T) {
	installers := makeInstallers(nil, nil)
	m := tui.NewModel(installers, nil)

	m.View()

	m_update, _ := m.Update(tui.OSDetectedMsg{Info: &domain.OSInfo{ID: "debian", VersionID: "13"}})
	m = m_update.(tui.Model)
	m.View()

	// stateNvidiaConfig
	m.SetCursor(0)
	m_update, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})
	m = m_update.(tui.Model)
	m_update, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = m_update.(tui.Model)
	m.View()

	// stateTokenInput
	m_update, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = m_update.(tui.Model)
	// Find GitlabTokenConfig
	idx := -1
	for i, id := range m.VisibleSoftware() { if id == domain.GitlabTokenConfig { idx = i; break } }
	m.SetCursor(idx)
	m_update, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})
	m = m_update.(tui.Model)
	m_update, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = m_update.(tui.Model)
	m.View()

	// stateProgress
	m_update, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("t")})
	m = m_update.(tui.Model)
	m_update, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = m_update.(tui.Model)
	// We can't easily reach progress state here without mocking the pre-install cmd
}
