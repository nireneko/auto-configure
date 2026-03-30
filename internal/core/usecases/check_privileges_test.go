package usecases_test

import (
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/core/usecases"
)

func TestCheckPrivileges_RootUser(t *testing.T) {
	uc := usecases.NewCheckPrivilegesUseCase(
		func() int { return 0 },
		func(string) string { return "" },
	)
	if err := uc.Execute(); err != nil {
		t.Errorf("expected nil error for root user, got %v", err)
	}
}

func TestCheckPrivileges_SudoUser(t *testing.T) {
	uc := usecases.NewCheckPrivilegesUseCase(
		func() int { return 1000 },
		func(key string) string {
			if key == "SUDO_UID" {
				return "1000"
			}
			return ""
		},
	)
	if err := uc.Execute(); err != nil {
		t.Errorf("expected nil error for sudo user, got %v", err)
	}
}

func TestCheckPrivileges_Unprivileged(t *testing.T) {
	uc := usecases.NewCheckPrivilegesUseCase(
		func() int { return 1000 },
		func(string) string { return "" },
	)
	err := uc.Execute()
	if err == nil {
		t.Fatal("expected PrivilegeError, got nil")
	}
	if _, ok := err.(domain.PrivilegeError); !ok {
		t.Errorf("expected PrivilegeError, got %T: %v", err, err)
	}
}
