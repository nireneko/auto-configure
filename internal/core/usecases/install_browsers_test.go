package usecases_test

import (
	"errors"
	"testing"
	"time"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/core/usecases"
	"github.com/so-install/pkg/mocks"
)

func noSleep(_ time.Duration) {}

func TestInstallSoftware_AlreadyInstalled(t *testing.T) {
	m := &mocks.MockSoftwareInstaller{
		SoftwareID:        domain.Brave,
		IsInstalledResult: true,
	}
	uc := usecases.NewInstallSoftwareUseCase(
		map[domain.SoftwareID]domain.SoftwareInstaller{domain.Brave: m},
		noSleep,
	)
	results := uc.Execute([]domain.SoftwareID{domain.Brave})

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].AlreadyInstalled || results[0].Err != nil {
		t.Errorf("expected AlreadyInstalled=true, Err=nil; got %+v", results[0])
	}
	if m.InstallCalled {
		t.Error("Install should not have been called for already-installed software")
	}
}

func TestInstallSoftware_SuccessfulInstall(t *testing.T) {
	m := &mocks.MockSoftwareInstaller{
		SoftwareID:        domain.Firefox,
		IsInstalledResult: false,
		InstallErr:        nil,
	}
	uc := usecases.NewInstallSoftwareUseCase(
		map[domain.SoftwareID]domain.SoftwareInstaller{domain.Firefox: m},
		noSleep,
	)
	results := uc.Execute([]domain.SoftwareID{domain.Firefox})

	if results[0].Err != nil || results[0].AlreadyInstalled {
		t.Errorf("expected success, got %+v", results[0])
	}
	if !m.InstallCalled {
		t.Error("Install should have been called")
	}
}

func TestInstallSoftware_AptLockRetries(t *testing.T) {
	aptErr := domain.AptLockError{InstallError: domain.InstallError{Software: "brave"}}
	installer := &countingInstaller{id: domain.Brave, err: aptErr}
	uc := usecases.NewInstallSoftwareUseCase(
		map[domain.SoftwareID]domain.SoftwareInstaller{domain.Brave: installer},
		noSleep,
	)
	results := uc.Execute([]domain.SoftwareID{domain.Brave})

	if results[0].Err == nil {
		t.Fatal("expected AptLockError after exhausted retries")
	}
	if installer.callCount != 3 {
		t.Errorf("expected 3 install attempts, got %d", installer.callCount)
	}
}

func TestInstallSoftware_SequentialOrder(t *testing.T) {
	order := []domain.SoftwareID{}
	braveInst := &orderTracker{id: domain.Brave, order: &order}
	firefoxInst := &orderTracker{id: domain.Firefox, order: &order}

	uc := usecases.NewInstallSoftwareUseCase(
		map[domain.SoftwareID]domain.SoftwareInstaller{
			domain.Brave:   braveInst,
			domain.Firefox: firefoxInst,
		},
		noSleep,
	)
	uc.Execute([]domain.SoftwareID{domain.Brave, domain.Firefox})

	if len(order) != 2 {
		t.Fatalf("expected 2 installs, got %d", len(order))
	}
	if order[0] != domain.Brave || order[1] != domain.Firefox {
		t.Errorf("wrong order: %v", order)
	}
}

func TestInstallSoftware_OneFailureContinuesOthers(t *testing.T) {
	braveInst := &mocks.MockSoftwareInstaller{
		SoftwareID: domain.Brave,
		InstallErr: errors.New("some error"),
	}
	firefoxInst := &mocks.MockSoftwareInstaller{
		SoftwareID: domain.Firefox,
	}

	uc := usecases.NewInstallSoftwareUseCase(
		map[domain.SoftwareID]domain.SoftwareInstaller{
			domain.Brave:   braveInst,
			domain.Firefox: firefoxInst,
		},
		noSleep,
	)
	results := uc.Execute([]domain.SoftwareID{domain.Brave, domain.Firefox})

	if results[0].Err == nil {
		t.Error("expected Brave to fail")
	}
	if results[1].Err != nil {
		t.Error("expected Firefox to succeed even after Brave failed")
	}
	if !firefoxInst.InstallCalled {
		t.Error("Firefox Install should have been called despite Brave failure")
	}
}

// Test helpers

type countingInstaller struct {
	id        domain.SoftwareID
	err       error
	callCount int
}

func (c *countingInstaller) Install() error {
	c.callCount++
	return c.err
}
func (c *countingInstaller) IsInstalled() (bool, error) { return false, nil }
func (c *countingInstaller) ID() domain.SoftwareID      { return c.id }

type orderTracker struct {
	id    domain.SoftwareID
	order *[]domain.SoftwareID
}

func (o *orderTracker) Install() error {
	*o.order = append(*o.order, o.id)
	return nil
}
func (o *orderTracker) IsInstalled() (bool, error) { return false, nil }
func (o *orderTracker) ID() domain.SoftwareID      { return o.id }
