package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/so-install/internal/core/domain"
	"github.com/so-install/pkg/mocks"
)

func TestModel_MandatoryStepsPrepended(t *testing.T) {
	installers := make(map[domain.SoftwareID]domain.SoftwareInstaller)
	for _, id := range domain.AllSoftware() {
		installers[id] = &mocks.MockSoftwareInstaller{SoftwareID: id}
	}
	installers[domain.SystemUpdate] = &mocks.MockSoftwareInstaller{SoftwareID: domain.SystemUpdate}
	installers[domain.BaseDeps] = &mocks.MockSoftwareInstaller{SoftwareID: domain.BaseDeps}

	m := NewModel(installers, nil)

	// 1. Move to SoftwareSelect state
	m_update, _ := m.Update(OSDetectedMsg{Info: &domain.OSInfo{ID: "debian", VersionID: "12"}})
	m_update, _ = m_update.Update(preInstalledCheckDoneMsg{}) 

	// 2. Select first item (Brave)
	m_update, _ = m_update.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})

	// 3. Confirm selection (Enter)
	m_update, _ = m_update.Update(tea.KeyMsg{Type: tea.KeyEnter})
	
	// 4. Start installation (startInstallMsg)
	m_update, _ = m_update.Update(startInstallMsg{})
	
	m_final := m_update.(Model)
	view := m_final.View()

	if !strings.Contains(view, "System Update") {
		t.Errorf("Progress view should contain 'System Update', got: %s", view)
	}
	if !strings.Contains(view, "Base Dependencies") {
		t.Errorf("Progress view should contain 'Base Dependencies', got: %s", view)
	}
}
