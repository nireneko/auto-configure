package usecases

import "github.com/so-install/internal/core/domain"

// CheckPrivilegesUseCase validates that the process runs as root or via sudo.
type CheckPrivilegesUseCase struct {
	uidFn func() int
	envFn func(string) string
}

// NewCheckPrivilegesUseCase creates a new use case. In production, pass os.Getuid and os.Getenv.
func NewCheckPrivilegesUseCase(uidFn func() int, envFn func(string) string) *CheckPrivilegesUseCase {
	return &CheckPrivilegesUseCase{uidFn: uidFn, envFn: envFn}
}

// Execute returns nil if running as root or via sudo, otherwise PrivilegeError.
func (uc *CheckPrivilegesUseCase) Execute() error {
	if uc.uidFn() == 0 {
		return nil
	}
	if uc.envFn("SUDO_UID") != "" {
		return nil
	}
	return domain.PrivilegeError{}
}
