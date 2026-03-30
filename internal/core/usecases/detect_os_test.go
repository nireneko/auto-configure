package usecases_test

import (
	"errors"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/core/usecases"
	"github.com/so-install/pkg/mocks"
)

func TestDetectOS_Debian12(t *testing.T) {
	det := &mocks.MockOSDetector{ReturnID: "debian", ReturnVersionID: "12"}
	uc := usecases.NewDetectOSUseCase(det)
	info, err := uc.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.ID != "debian" || info.VersionID != "12" {
		t.Errorf("wrong info: %+v", info)
	}
}

func TestDetectOS_Debian13(t *testing.T) {
	det := &mocks.MockOSDetector{ReturnID: "debian", ReturnVersionID: "13"}
	uc := usecases.NewDetectOSUseCase(det)
	_, err := uc.Execute()
	if err != nil {
		t.Fatalf("unexpected error for debian 13: %v", err)
	}
}

func TestDetectOS_Ubuntu_Rejected(t *testing.T) {
	det := &mocks.MockOSDetector{ReturnID: "ubuntu", ReturnVersionID: "24.04"}
	uc := usecases.NewDetectOSUseCase(det)
	_, err := uc.Execute()
	if err == nil {
		t.Fatal("expected OsNotSupportedError, got nil")
	}
	var nsErr domain.OsNotSupportedError
	if !errors.As(err, &nsErr) {
		t.Errorf("expected OsNotSupportedError, got %T", err)
	}
}

func TestDetectOS_Debian11_Rejected(t *testing.T) {
	det := &mocks.MockOSDetector{ReturnID: "debian", ReturnVersionID: "11"}
	uc := usecases.NewDetectOSUseCase(det)
	_, err := uc.Execute()
	if err == nil {
		t.Fatal("expected OsNotSupportedError for debian 11, got nil")
	}
}

func TestDetectOS_DetectorError(t *testing.T) {
	det := &mocks.MockOSDetector{ReturnErr: errors.New("file not found")}
	uc := usecases.NewDetectOSUseCase(det)
	_, err := uc.Execute()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
