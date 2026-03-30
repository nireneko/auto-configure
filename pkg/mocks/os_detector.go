package mocks

import "github.com/so-install/internal/core/domain"

// MockOSDetector returns a configurable OSInfo or error.
type MockOSDetector struct {
	ReturnID        string
	ReturnVersionID string
	ReturnErr       error
}

var _ domain.OSDetector = (*MockOSDetector)(nil)

func (m *MockOSDetector) Detect() (*domain.OSInfo, error) {
	if m.ReturnErr != nil {
		return nil, m.ReturnErr
	}
	return &domain.OSInfo{ID: m.ReturnID, VersionID: m.ReturnVersionID}, nil
}
