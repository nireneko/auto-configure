package mocks

import (
	"errors"
	"testing"
)

func TestMockOSDetector_ReturnsConfiguredValues(t *testing.T) {
	m := &MockOSDetector{ReturnID: "debian", ReturnVersionID: "12"}
	info, err := m.Detect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.ID != "debian" || info.VersionID != "12" {
		t.Errorf("wrong info: %+v", info)
	}
}

func TestMockOSDetector_ReturnsError(t *testing.T) {
	m := &MockOSDetector{ReturnErr: errors.New("boom")}
	_, err := m.Detect()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
