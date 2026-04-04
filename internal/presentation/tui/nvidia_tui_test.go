package tui_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/presentation/tui"
	"github.com/so-install/pkg/mocks"
)

func makeInstallersWithNvidia() map[domain.SoftwareID]domain.SoftwareInstaller {
	result := make(map[domain.SoftwareID]domain.SoftwareInstaller)
	for _, id := range domain.AllSoftware() {
		result[id] = &mocks.MockSoftwareInstaller{SoftwareID: id}
	}
	result[domain.SystemUpdate] = &mocks.MockSoftwareInstaller{SoftwareID: domain.SystemUpdate}
	result[domain.BaseDeps] = &mocks.MockSoftwareInstaller{SoftwareID: domain.BaseDeps}
	return result
}

// reachSoftwareSelect drives the model from welcome to the software select screen
// using the given osInfo.
func reachSoftwareSelect(osInfo *domain.OSInfo) tui.Model {
	m := tui.NewModel(makeInstallersWithNvidia())
	m.SetOSInfo(osInfo)
	updated, cmd := m.Update(tui.OSDetectedMsg{Info: osInfo})
	if cmd != nil {
		updated, _ = updated.Update(cmd())
	}
	return updated.(tui.Model)
}

// TestNvidiaDrivers_HiddenOnNonDebian13 verifies that NvidiaDrivers does not
// appear in the software list when the OS is not Debian 13.
func TestNvidiaDrivers_HiddenOnNonDebian13(t *testing.T) {
	cases := []struct {
		id        string
		versionID string
	}{
		{"debian", "12"},
		{"ubuntu", "22.04"},
		{"debian", ""},
	}
	for _, tc := range cases {
		m := reachSoftwareSelect(&domain.OSInfo{ID: tc.id, VersionID: tc.versionID})
		view := m.View()
		if strings.Contains(view, "Nvidia Drivers") {
			t.Errorf("OS %s %s: expected Nvidia Drivers to be hidden, but it appeared in view",
				tc.id, tc.versionID)
		}
	}
}

// TestNvidiaDrivers_ShownOnDebian13 verifies that NvidiaDrivers appears
// in the software list when the OS is Debian 13.
func TestNvidiaDrivers_ShownOnDebian13(t *testing.T) {
	m := reachSoftwareSelect(&domain.OSInfo{ID: "debian", VersionID: "13"})
	view := m.View()
	if !strings.Contains(view, "Nvidia Drivers") {
		t.Errorf("expected Nvidia Drivers to appear on Debian 13, got:\n%s", view)
	}
}

// TestNvidiaDrivers_FreePath verifies the flow:
// softwareSelect → stateNvidiaConfig (sub-step 0) → startInstallMsg
// when the user selects the Free driver type (skips CUDA sub-step).
func TestNvidiaDrivers_FreePath(t *testing.T) {
	m := reachSoftwareSelect(&domain.OSInfo{ID: "debian", VersionID: "13"})

	// NvidiaDrivers is the last entry in the visible list on Debian 13.
	all := domain.AllSoftware()
	nvidiaIdx := -1
	for i, id := range all {
		if id == domain.NvidiaDrivers {
			nvidiaIdx = i
			break
		}
	}
	if nvidiaIdx == -1 {
		t.Fatal("NvidiaDrivers not found in AllSoftware()")
	}

	// Navigate cursor to NvidiaDrivers and toggle it.
	m.SetCursor(nvidiaIdx)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})
	m = updated.(tui.Model)

	// Confirm selection → should go to stateNvidiaConfig
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(tui.Model)
	if cmd != nil {
		t.Fatal("expected no cmd immediately after Enter on softwareSelect when Nvidia is selected")
	}

	view := m.View()
	if !strings.Contains(view, "Nvidia Driver Configuration") {
		t.Fatalf("expected Nvidia config view, got:\n%s", view)
	}
	if !strings.Contains(view, "Free") {
		t.Fatalf("expected driver type options, got:\n%s", view)
	}

	// Press Enter on option 0 (Free) → should emit startInstallMsg (or nvidiaConfigDoneMsg)
	updated, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected cmd after selecting Free driver")
	}
	// cmd is non-nil — it will deliver startInstallMsg to the Update loop
}

// TestNvidiaDrivers_ProprietaryPath verifies the flow:
// softwareSelect → stateNvidiaConfig sub-step 0 → sub-step 1 (CUDA) → startInstallMsg.
func TestNvidiaDrivers_ProprietaryPath(t *testing.T) {
	m := reachSoftwareSelect(&domain.OSInfo{ID: "debian", VersionID: "13"})

	all := domain.AllSoftware()
	nvidiaIdx := -1
	for i, id := range all {
		if id == domain.NvidiaDrivers {
			nvidiaIdx = i
			break
		}
	}
	m.SetCursor(nvidiaIdx)

	// Toggle NvidiaDrivers
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})
	m = updated.(tui.Model)

	// Confirm → stateNvidiaConfig
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(tui.Model)
	if cmd != nil {
		t.Fatal("unexpected cmd on Enter in softwareSelect")
	}

	// Move cursor down once (to ProprietaryDebian) and press Enter
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	m = updated.(tui.Model)
	updated, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(tui.Model)
	if cmd != nil {
		t.Fatal("expected no cmd after selecting ProprietaryDebian — should show CUDA prompt")
	}

	view := m.View()
	if !strings.Contains(view, "CUDA") {
		t.Fatalf("expected CUDA prompt after selecting proprietary driver, got:\n%s", view)
	}

	// Press N for no CUDA → startInstallMsg
	updated, cmd = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	if cmd == nil {
		t.Fatal("expected startInstallMsg cmd after answering CUDA prompt")
	}
}

// TestNvidiaDrivers_EscGoesBack verifies Esc in stateNvidiaConfig returns to
// softwareSelect.
func TestNvidiaDrivers_EscGoesBack(t *testing.T) {
	m := reachSoftwareSelect(&domain.OSInfo{ID: "debian", VersionID: "13"})

	all := domain.AllSoftware()
	for i, id := range all {
		if id == domain.NvidiaDrivers {
			m.SetCursor(i)
			break
		}
	}

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})
	m = updated.(tui.Model)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(tui.Model)

	// Esc from sub-step 0 → back to softwareSelect
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(tui.Model)

	view := m.View()
	if !strings.Contains(view, "Select software to install") {
		t.Errorf("expected softwareSelect after Esc, got:\n%s", view)
	}
}
