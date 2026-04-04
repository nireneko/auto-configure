package tui_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/presentation/tui"
)

func TestModel_NvidiaConfigTransitions(t *testing.T) {
	installers := makeInstallers(nil, nil)
	m := tui.NewModel(installers, nil)

	m_update, _ := m.Update(tui.OSDetectedMsg{Info: &domain.OSInfo{ID: "debian", VersionID: "13"}})
	m = m_update.(tui.Model)

	software := m.VisibleSoftware()
	idx := -1
	for i, id := range software {
		if id == domain.NvidiaDrivers {
			idx = i
			break
		}
	}
	if idx == -1 { t.Fatal("NvidiaDrivers not found in visible list for Debian 13") }
	m.SetCursor(idx)
	m_update, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})
	m = m_update.(tui.Model)

	m_update, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = m_update.(tui.Model)

	// In stateNvidiaConfig
	m_update, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}) // down
	m = m_update.(tui.Model)
	m_update, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")}) // up
	m = m_update.(tui.Model)

	m_update, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = m_update.(tui.Model)

	m_update, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	m = m_update.(tui.Model)

	// Other keys
	m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
}
