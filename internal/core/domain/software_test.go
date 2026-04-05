package domain_test

import (
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestSoftwareID_DisplayName(t *testing.T) {
	tests := []struct {
		id       domain.SoftwareID
		expected string
	}{
		{domain.SystemUpdate, "System Update"},
		{domain.BaseDeps, "Base Dependencies"},
		{domain.Brave, "Brave"},
		{domain.Firefox, "Firefox"},
		{domain.Chrome, "Google Chrome"},
		{domain.Chromium, "Chromium"},
		{domain.Docker, "Docker CE"},
		{domain.Ddev, "DDEV"},
		{domain.OpenVpn, "OpenVPN"},
		{domain.Nvm, "NVM & NPM"},
		{domain.Gemini, "Google Gemini CLI"},
		{domain.ClaudeCode, "Claude Code (Anthropic)"},
		{domain.Codex, "OpenAI Codex CLI"},
		{domain.Flatpak, "Flatpak"},
		{domain.Bitwarden, "Bitwarden"},
		{domain.Homebrew, "Homebrew"},
		{domain.GitlabTokenConfig, "Gitlab Token Configuration (Composer/NPM)"},
		{domain.ScreenLockConfig, "Screen Lock Configuration"},
		{domain.GentleAI, "Gentle-AI"},
		{domain.VsCode, "Visual Studio Code"},
		{domain.Cursor, "Cursor IDE"},
		{domain.Antigravity, "Google Antigravity"},
		{domain.SoftwareID("unknown-software"), "unknown-software"}, // Default fallback
	}

	for _, tt := range tests {
		t.Run(string(tt.id), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.id.DisplayName())
		})
	}
}

func TestGetSteps(t *testing.T) {
	steps := domain.GetSteps()
	assert.NotEmpty(t, steps)
}

func TestGetSteps_AiCli_ContainsTools(t *testing.T) {
	steps := domain.GetSteps()
	var aiCliStep *domain.InstallStep
	for i := range steps {
		if steps[i].ID == "ai-cli" {
			aiCliStep = &steps[i]
			break
		}
	}
	if aiCliStep == nil {
		t.Fatal("ai-cli step not found in GetSteps()")
	}
	assert.Contains(t, aiCliStep.Software, domain.Gemini)
	assert.Contains(t, aiCliStep.Software, domain.ClaudeCode)
	assert.Contains(t, aiCliStep.Software, domain.Codex)
}

func TestAllSoftware(t *testing.T) {
	all := domain.AllSoftware()
	assert.NotEmpty(t, all)
}

func TestAllSoftware_ContainsGentleAI(t *testing.T) {
	all := domain.AllSoftware()
	assert.Contains(t, all, domain.GentleAI)
}

func TestGetSteps_GentleAI_IsAfterAiCli(t *testing.T) {
	steps := domain.GetSteps()
	var aiCliIdx, gentleAiIdx int = -1, -1
	for i, step := range steps {
		if step.ID == "ai-cli" {
			aiCliIdx = i
		}
		if step.ID == "gentle-ai" {
			gentleAiIdx = i
		}
	}
	assert.NotEqual(t, -1, aiCliIdx, "ai-cli step not found")
	assert.NotEqual(t, -1, gentleAiIdx, "gentle-ai step not found")
	assert.Equal(t, aiCliIdx+1, gentleAiIdx, "gentle-ai step must be immediately after ai-cli")

	// Also check content
	var gentleAiStep domain.InstallStep
	for _, step := range steps {
		if step.ID == "gentle-ai" {
			gentleAiStep = step
			break
		}
	}
	assert.Contains(t, gentleAiStep.Software, domain.GentleAI)
	assert.Len(t, gentleAiStep.Software, 1, "gentle-ai step must contain only one item")
	assert.False(t, gentleAiStep.Critical, "gentle-ai step must not be critical")
}

func TestGetSteps_Ides_IsAfterGentleAi(t *testing.T) {
	steps := domain.GetSteps()
	var gentleAiIdx, idesIdx int = -1, -1
	for i, step := range steps {
		if step.ID == "gentle-ai" {
			gentleAiIdx = i
		}
		if step.ID == "ides" {
			idesIdx = i
		}
	}
	assert.NotEqual(t, -1, gentleAiIdx, "gentle-ai step not found")
	assert.NotEqual(t, -1, idesIdx, "ides step not found")
	assert.Equal(t, gentleAiIdx+1, idesIdx, "ides step must be immediately after gentle-ai")

	// Also check content
	var idesStep domain.InstallStep
	for _, step := range steps {
		if step.ID == "ides" {
			idesStep = step
			break
		}
	}
	assert.Contains(t, idesStep.Software, domain.VsCode)
	assert.Contains(t, idesStep.Software, domain.Cursor)
	assert.Contains(t, idesStep.Software, domain.Antigravity)
	assert.Len(t, idesStep.Software, 3, "ides step must contain three items")
	assert.False(t, idesStep.Critical, "ides step must not be critical")
}
