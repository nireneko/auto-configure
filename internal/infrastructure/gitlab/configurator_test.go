package gitlab

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/so-install/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGitlabTokenConfigurator_Install(t *testing.T) {
	t.Run("should configure composer and npm with token", func(t *testing.T) {
		tempHome := t.TempDir()
		executor := &mocks.MockExecutor{}
		configurator := NewGitlabTokenConfigurator(executor)
		configurator.SetHomeDir(tempHome)
		configurator.SetToken("my-secret-token")

		err := configurator.Install()
		require.NoError(t, err)

		// Verify Composer auth.json
		authFile := filepath.Join(tempHome, ".composer", "auth.json")
		content, err := os.ReadFile(authFile)
		require.NoError(t, err)

		var data map[string]interface{}
		err = json.Unmarshal(content, &data)
		require.NoError(t, err)

		gitlabTokens := data["gitlab-token"].(map[string]interface{})
		assert.Equal(t, "my-secret-token", gitlabTokens["gitlab.com"])

		// Verify NPM .npmrc
		npmrcFile := filepath.Join(tempHome, ".npmrc")
		npmrcContent, err := os.ReadFile(npmrcFile)
		require.NoError(t, err)
		assert.Contains(t, string(npmrcContent), "//gitlab.com/api/v4/packages/npm/:_authToken=my-secret-token")
	})

	t.Run("should update existing configurations", func(t *testing.T) {
		tempHome := t.TempDir()
		
		// Pre-create composer auth.json
		composerDir := filepath.Join(tempHome, ".composer")
		require.NoError(t, os.MkdirAll(composerDir, 0755))
		oldAuth := `{"gitlab-token": {"gitlab.com": "old-token", "other.com": "other-token"}}`
		require.NoError(t, os.WriteFile(filepath.Join(composerDir, "auth.json"), []byte(oldAuth), 0600))

		// Pre-create .npmrc
		oldNpmrc := "registry=https://registry.npmjs.org/\n//gitlab.com/api/v4/packages/npm/:_authToken=old-token\n"
		require.NoError(t, os.WriteFile(filepath.Join(tempHome, ".npmrc"), []byte(oldNpmrc), 0600))

		executor := &mocks.MockExecutor{}
		configurator := NewGitlabTokenConfigurator(executor)
		configurator.SetHomeDir(tempHome)
		configurator.SetToken("new-token")

		err := configurator.Install()
		require.NoError(t, err)

		// Verify Composer auth.json updated
		authFile := filepath.Join(tempHome, ".composer", "auth.json")
		content, err := os.ReadFile(authFile)
		require.NoError(t, err)

		var data map[string]interface{}
		json.Unmarshal(content, &data)
		gitlabTokens := data["gitlab-token"].(map[string]interface{})
		assert.Equal(t, "new-token", gitlabTokens["gitlab.com"])
		assert.Equal(t, "other-token", gitlabTokens["other.com"])

		// Verify NPM .npmrc updated
		npmrcFile := filepath.Join(tempHome, ".npmrc")
		npmrcContent, err := os.ReadFile(npmrcFile)
		require.NoError(t, err)
		assert.Contains(t, string(npmrcContent), "//gitlab.com/api/v4/packages/npm/:_authToken=new-token")
		assert.Contains(t, string(npmrcContent), "registry=https://registry.npmjs.org/")
		// Ensure only one gitlab token line exists
		assert.Equal(t, 1, strings.Count(string(npmrcContent), "_authToken="))
	})
}
