package usecases

import "github.com/so-install/internal/core/domain"

// supported maps distro ID to accepted version IDs.
var supported = map[string]map[string]bool{
	"debian": {"12": true, "13": true},
}

// DetectOSUseCase detects and validates the OS.
type DetectOSUseCase struct {
	detector domain.OSDetector
}

// NewDetectOSUseCase creates a new use case with the provided detector.
func NewDetectOSUseCase(detector domain.OSDetector) *DetectOSUseCase {
	return &DetectOSUseCase{detector: detector}
}

// Execute returns OSInfo if the OS is supported, or an OsNotSupportedError.
func (uc *DetectOSUseCase) Execute() (*domain.OSInfo, error) {
	info, err := uc.detector.Detect()
	if err != nil {
		return nil, err
	}
	versions, ok := supported[info.ID]
	if !ok || !versions[info.VersionID] {
		return nil, domain.OsNotSupportedError{OS: info.ID, Version: info.VersionID}
	}
	return info, nil
}
