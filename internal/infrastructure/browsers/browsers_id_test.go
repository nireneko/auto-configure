package browsers_test

import (
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/browsers"
	"github.com/so-install/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestBrowsers_ID(t *testing.T) {
	exec := &mocks.MockExecutor{}
	
	brave := browsers.NewBraveInstaller(exec)
	assert.Equal(t, domain.Brave, brave.ID())

	chrome := browsers.NewChromeInstaller(exec)
	assert.Equal(t, domain.Chrome, chrome.ID())

	chromium := browsers.NewChromiumInstaller(exec)
	assert.Equal(t, domain.Chromium, chromium.ID())

	firefox := browsers.NewFirefoxInstaller(exec)
	assert.Equal(t, domain.Firefox, firefox.ID())
}
