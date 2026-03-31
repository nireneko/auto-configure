package gitlab

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/so-install/internal/core/domain"
)

// GitlabTokenConfigurator configures Gitlab tokens for Composer and NPM.
type GitlabTokenConfigurator struct {
	executor domain.Executor
	token    string
	homeDir  string
}

// NewGitlabTokenConfigurator creates a new GitlabTokenConfigurator.
func NewGitlabTokenConfigurator(executor domain.Executor) *GitlabTokenConfigurator {
	return &GitlabTokenConfigurator{
		executor: executor,
		homeDir:  domain.GetActualHome(),
	}
}

var _ domain.SoftwareInstaller = (*GitlabTokenConfigurator)(nil)

// ID returns the SoftwareID.
func (g *GitlabTokenConfigurator) ID() domain.SoftwareID {
	return domain.GitlabTokenConfig
}

// SetToken sets the Gitlab personal access token.
func (g *GitlabTokenConfigurator) SetToken(token string) {
	g.token = token
}

// SetHomeDir overrides the home directory (useful for testing).
func (g *GitlabTokenConfigurator) SetHomeDir(homeDir string) {
	g.homeDir = homeDir
}

// IsInstalled always returns false as this is a configuration task.
// Or we could check if the token is already present in the files.
func (g *GitlabTokenConfigurator) IsInstalled() (bool, error) {
	return false, nil
}

// Install performs the configuration of Composer and NPM.
func (g *GitlabTokenConfigurator) Install() error {
	if g.token == "" {
		return fmt.Errorf("gitlab token is not set")
	}

	if err := g.configureComposer(); err != nil {
		return err
	}

	if err := g.configureNpm(); err != nil {
		return err
	}

	return nil
}

func (g *GitlabTokenConfigurator) configureComposer() error {
	composerDir := filepath.Join(g.homeDir, ".composer")
	authFile := filepath.Join(composerDir, "auth.json")

	if err := os.MkdirAll(composerDir, 0755); err != nil {
		return fmt.Errorf("failed to create composer directory: %w", err)
	}

	data := make(map[string]interface{})
	if _, err := os.Stat(authFile); err == nil {
		content, err := os.ReadFile(authFile)
		if err == nil {
			_ = json.Unmarshal(content, &data)
		}
	}

	gitlabTokens, ok := data["gitlab-token"].(map[string]interface{})
	if !ok {
		gitlabTokens = make(map[string]interface{})
	}
	gitlabTokens["gitlab.com"] = g.token
	data["gitlab-token"] = gitlabTokens

	newContent, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal composer auth.json: %w", err)
	}

	if err := os.WriteFile(authFile, newContent, 0600); err != nil {
		return fmt.Errorf("failed to write composer auth.json: %w", err)
	}

	return nil
}

func (g *GitlabTokenConfigurator) configureNpm() error {
	npmrcFile := filepath.Join(g.homeDir, ".npmrc")
	configLine := fmt.Sprintf("//gitlab.com/api/v4/packages/npm/:_authToken=%s", g.token)

	var lines []string
	if _, err := os.Stat(npmrcFile); err == nil {
		content, err := os.ReadFile(npmrcFile)
		if err == nil {
			lines = strings.Split(string(content), "\n")
		}
	}

	found := false
	for i, line := range lines {
		if strings.Contains(line, "//gitlab.com/api/v4/packages/npm/:_authToken=") {
			lines[i] = configLine
			found = true
			break
		}
	}

	if !found {
		lines = append(lines, configLine)
	}

	newContent := strings.Join(lines, "\n")
	if !strings.HasSuffix(newContent, "\n") {
		newContent += "\n"
	}

	if err := os.WriteFile(npmrcFile, []byte(newContent), 0600); err != nil {
		return fmt.Errorf("failed to write .npmrc: %w", err)
	}

	return nil
}
