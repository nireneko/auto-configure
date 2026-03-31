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

func TestAllSoftware(t *testing.T) {
	all := domain.AllSoftware()
	assert.NotEmpty(t, all)
}
