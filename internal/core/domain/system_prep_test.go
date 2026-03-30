package domain_test

import (
	"testing"

	"github.com/so-install/internal/core/domain"
)

func TestSystemPrep_GetSteps(t *testing.T) {
	steps := domain.GetSteps()

	// 1. Verify that first step is system-prep and is critical
	if len(steps) < 1 {
		t.Fatalf("Expected at least 1 step, got %d", len(steps))
	}

	firstStep := steps[0]
	if firstStep.ID != "system-prep" {
		t.Errorf("First step ID = %s, want system-prep", firstStep.ID)
	}
	if !firstStep.Critical {
		t.Error("System-prep step SHOULD be critical")
	}

	// 2. Verify software in system-prep step
	expectedSoftware := []domain.SoftwareID{domain.SystemUpdate, domain.BaseDeps}
	if len(firstStep.Software) != len(expectedSoftware) {
		t.Errorf("System-prep step software count = %d, want %d", len(firstStep.Software), len(expectedSoftware))
	}

	for i, id := range expectedSoftware {
		if i < len(firstStep.Software) && firstStep.Software[i] != id {
			t.Errorf("System-prep software at index %d = %s, want %s", i, firstStep.Software[i], id)
		}
	}
}

func TestSystemPrep_AllSoftware(t *testing.T) {
	all := domain.AllSoftware()
	for _, id := range all {
		if id == domain.SystemUpdate || id == domain.BaseDeps {
			t.Errorf("AllSoftware() SHOULD NOT include %s", id)
		}
	}
}

func TestSystemPrep_DisplayName(t *testing.T) {
	tests := []struct {
		id   domain.SoftwareID
		want string
	}{
		{domain.SystemUpdate, "System Update"},
		{domain.BaseDeps, "Base Dependencies"},
	}

	for _, tt := range tests {
		t.Run(string(tt.id), func(t *testing.T) {
			if got := tt.id.DisplayName(); got != tt.want {
				t.Errorf("SoftwareID.DisplayName() = %v, want %v", got, tt.want)
			}
		})
	}
}
