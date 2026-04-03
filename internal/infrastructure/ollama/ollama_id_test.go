package ollama_test

import (
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/ollama"
	"github.com/so-install/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestOllama_ID(t *testing.T) {
	exec := &mocks.MockExecutor{}
	installer := ollama.NewOllamaInstaller(exec)
	assert.Equal(t, domain.Ollama, installer.ID())
}
