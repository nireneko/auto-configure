package domain_test

import (
	"testing"

	"github.com/so-install/internal/core/domain"
)

func TestSoftwareID_DisplayName(t *testing.T) {
	tests := []struct {
		id   domain.SoftwareID
		want string
	}{
		{domain.Brave, "Brave"},
		{domain.Firefox, "Firefox"},
		{domain.Chrome, "Google Chrome"},
		{domain.Chromium, "Chromium"},
		{domain.Docker, "Docker CE"},
		{domain.SoftwareID("ddev"), "DDEV"}, // RED: ddev doesn't exist as constant yet
	}
	for _, tt := range tests {
		t.Run(string(tt.id), func(t *testing.T) {
			if got := tt.id.DisplayName(); got != tt.want {
				t.Errorf("SoftwareID.DisplayName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllSoftware(t *testing.T) {
	all := domain.AllSoftware()
	foundDdev := false
	foundDocker := false
	dockerIdx := -1
	ddevIdx := -1

	for i, id := range all {
		if id == domain.Docker {
			foundDocker = true
			dockerIdx = i
		}
		if id == domain.SoftwareID("ddev") {
			foundDdev = true
			ddevIdx = i
		}
	}

	if !foundDdev {
		t.Error("AllSoftware() should include ddev")
	}
	if !foundDocker {
		t.Error("AllSoftware() should include docker")
	}
	if foundDocker && foundDdev && ddevIdx <= dockerIdx {
		t.Errorf("ddev (index %d) should come after docker (index %d)", ddevIdx, dockerIdx)
	}
}
