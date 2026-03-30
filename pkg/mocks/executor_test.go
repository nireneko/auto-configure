package mocks

import (
	"errors"
	"testing"
)

func TestMockExecutor_ResponsesConsumedInOrder(t *testing.T) {
	m := &MockExecutor{}
	m.AddResponse("out1", "err1", nil)
	m.AddResponse("out2", "err2", errors.New("fail"))

	s, e, err := m.Execute("cmd", "a")
	if s != "out1" || e != "err1" || err != nil {
		t.Fatalf("first response wrong: got %q %q %v", s, e, err)
	}
	s, e, err = m.Execute("cmd", "b")
	if s != "out2" || e != "err2" || err == nil {
		t.Fatalf("second response wrong: got %q %q %v", s, e, err)
	}
}

func TestMockExecutor_CallsRecorded(t *testing.T) {
	m := &MockExecutor{}
	m.Execute("apt", "install", "-y", "brave-browser")
	m.Execute("curl", "-fsSLo", "/tmp/key.gpg")

	if len(m.Calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(m.Calls))
	}
	if m.Calls[0].Name != "apt" {
		t.Errorf("first call: expected apt, got %s", m.Calls[0].Name)
	}
	if m.Calls[1].Name != "curl" {
		t.Errorf("second call: expected curl, got %s", m.Calls[1].Name)
	}
}

func TestMockExecutor_ExhaustedQueueReturnsDefault(t *testing.T) {
	m := &MockExecutor{}
	m.SetDefault("default-out", "", nil)

	s, _, _ := m.Execute("anything")
	if s != "default-out" {
		t.Errorf("expected default-out, got %s", s)
	}
}
