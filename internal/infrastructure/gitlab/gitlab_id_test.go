package gitlab

import (
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/pkg/mocks"
)

func TestGitlabTokenConfigurator_ID(t *testing.T) {
	configurator := NewGitlabTokenConfigurator(&mocks.MockExecutor{})
	if configurator.ID() != domain.GitlabTokenConfig { t.Errorf("Expected ID %v, got %v", domain.GitlabTokenConfig, configurator.ID()) }
}

func TestGitlabTokenConfigurator_IsInstalled(t *testing.T) {
	configurator := NewGitlabTokenConfigurator(&mocks.MockExecutor{})
	installed, _ := configurator.IsInstalled()
	if installed { t.Error("Expected false") }
}

func TestGitlabTokenConfigurator_Install_NoToken(t *testing.T) {
	configurator := NewGitlabTokenConfigurator(&mocks.MockExecutor{})
	err := configurator.Install()
	if err == nil { t.Error("Expected error when token is not set") }
}
