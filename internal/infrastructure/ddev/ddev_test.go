package ddev_test

import (
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/ddev"
	"github.com/so-install/pkg/mocks"
)

func TestDdevInstaller_Install_HappyPath(t *testing.T) {
	executor := &mocks.MockExecutor{}
	installer := ddev.NewDdevInstaller(executor)

	err := installer.Install()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedSteps := []struct {
		name string
		args []string
	}{
		{"apt-get", []string{"update"}},
		{"apt-get", []string{"install", "-y", "curl"}},
		{"install", []string{"-m", "0755", "-d", "/etc/apt/keyrings"}},
		{"sh", []string{"-c", "curl -fsSL https://pkg.ddev.com/apt/gpg.key | gpg --dearmor | sudo tee /etc/apt/keyrings/ddev.gpg > /dev/null"}},
		{"chmod", []string{"a+r", "/etc/apt/keyrings/ddev.gpg"}},
		{"sh", []string{"-c", `echo "deb [signed-by=/etc/apt/keyrings/ddev.gpg] https://pkg.ddev.com/apt/ * *" | sudo tee /etc/apt/sources.list.d/ddev.list >/dev/null`}},
		{"apt-get", []string{"update"}},
		{"apt-get", []string{"install", "-y", "ddev"}},
		{"mkcert", []string{"-install"}},
	}

	if len(executor.Calls) != len(expectedSteps) {
		t.Fatalf("Expected %d calls, got %d", len(expectedSteps), len(executor.Calls))
	}

	for i, step := range expectedSteps {
		if executor.Calls[i].Name != step.name {
			t.Errorf("Step %d: expected command %s, got %s", i+1, step.name, executor.Calls[i].Name)
		}
		if len(executor.Calls[i].Args) != len(step.args) {
			t.Errorf("Step %d: expected %d args, got %d", i+1, len(step.args), len(executor.Calls[i].Args))
			continue
		}
		for j, arg := range step.args {
			if executor.Calls[i].Args[j] != arg {
				t.Errorf("Step %d, arg %d: expected %s, got %s", i+1, j, arg, executor.Calls[i].Args[j])
			}
		}
	}
}

func TestDdevInstaller_IsInstalled(t *testing.T) {
	t.Run("Already installed", func(t *testing.T) {
		executor := &mocks.MockExecutor{}
		executor.AddResponse("v1.23.0", "", nil)
		installer := ddev.NewDdevInstaller(executor)

		installed, err := installer.IsInstalled()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !installed {
			t.Error("Expected installed to be true")
		}
		if len(executor.Calls) != 1 || executor.Calls[0].Name != "ddev" || executor.Calls[0].Args[0] != "--version" {
			t.Errorf("Unexpected executor calls: %v", executor.Calls)
		}
	})

	t.Run("Not installed", func(t *testing.T) {
		executor := &mocks.MockExecutor{}
		executor.AddResponse("", "not found", domain.InstallError{Software: "ddev", ExitCode: 127})
		installer := ddev.NewDdevInstaller(executor)

		installed, err := installer.IsInstalled()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if installed {
			t.Error("Expected installed to be false")
		}
	})
}
