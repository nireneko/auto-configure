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

func TestGetSteps(t *testing.T) {
	steps := domain.GetSteps()

	if len(steps) != 3 {
		t.Fatalf("Expected 3 steps, got %d", len(steps))
	}

	// Step 1: Browsers (Not critical)
	if steps[0].ID != "browsers" {
		t.Errorf("Step 1 ID = %s, want browsers", steps[0].ID)
	}
	if steps[0].Critical {
		t.Error("Step 1 (browsers) should NOT be critical")
	}
	if len(steps[0].Software) != 4 {
		t.Errorf("Step 1 should have 4 browsers, got %d", len(steps[0].Software))
	}

	// Step 2: Docker (Critical)
	if steps[1].ID != "docker" {
		t.Errorf("Step 2 ID = %s, want docker", steps[1].ID)
	}
	if !steps[1].Critical {
		t.Error("Step 2 (docker) SHOULD be critical")
	}

	// Step 3: DDEV (Not critical, but depends on Docker)
	if steps[2].ID != "ddev" {
		t.Errorf("Step 3 ID = %s, want ddev", steps[2].ID)
	}
}
