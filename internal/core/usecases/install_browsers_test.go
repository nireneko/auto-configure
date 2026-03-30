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

func TestInstallBrowsers_AlreadyInstalled(t *testing.T) {
	m := &mocks.MockBrowserInstaller{
		BrowserID:         domain.Brave,
		IsInstalledResult: true,
	}
	uc := usecases.NewInstallBrowsersUseCase(
		map[domain.BrowserID]domain.BrowserInstaller{domain.Brave: m},
		noSleep,
	)
	results := uc.Execute([]domain.BrowserID{domain.Brave})

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].AlreadyInstalled || results[0].Err != nil {
		t.Errorf("expected AlreadyInstalled=true, Err=nil; got %+v", results[0])
	}
	if m.InstallCalled {
		t.Error("Install should not have been called for already-installed browser")
	}
}

func TestInstallBrowsers_SuccessfulInstall(t *testing.T) {
	m := &mocks.MockBrowserInstaller{
		BrowserID:         domain.Firefox,
		IsInstalledResult: false,
		InstallErr:        nil,
	}
	uc := usecases.NewInstallBrowsersUseCase(
		map[domain.BrowserID]domain.BrowserInstaller{domain.Firefox: m},
		noSleep,
	)
	results := uc.Execute([]domain.BrowserID{domain.Firefox})

	if results[0].Err != nil || results[0].AlreadyInstalled {
		t.Errorf("expected success, got %+v", results[0])
	}
	if !m.InstallCalled {
		t.Error("Install should have been called")
	}
}

func TestInstallBrowsers_AptLockRetries(t *testing.T) {
	aptErr := domain.AptLockError{InstallError: domain.InstallError{Browser: "brave"}}
	installer := &countingInstaller{id: domain.Brave, err: aptErr}
	uc := usecases.NewInstallBrowsersUseCase(
		map[domain.BrowserID]domain.BrowserInstaller{domain.Brave: installer},
		noSleep,
	)
	results := uc.Execute([]domain.BrowserID{domain.Brave})

	if results[0].Err == nil {
		t.Fatal("expected AptLockError after exhausted retries")
	}
	if installer.callCount != 3 {
		t.Errorf("expected 3 install attempts, got %d", installer.callCount)
	}
}

func TestInstallBrowsers_SequentialOrder(t *testing.T) {
	order := []domain.BrowserID{}
	braveInst := &orderTracker{id: domain.Brave, order: &order}
	firefoxInst := &orderTracker{id: domain.Firefox, order: &order}

	uc := usecases.NewInstallBrowsersUseCase(
		map[domain.BrowserID]domain.BrowserInstaller{
			domain.Brave:   braveInst,
			domain.Firefox: firefoxInst,
		},
		noSleep,
	)
	uc.Execute([]domain.BrowserID{domain.Brave, domain.Firefox})

	if len(order) != 2 {
		t.Fatalf("expected 2 installs, got %d", len(order))
	}
	if order[0] != domain.Brave || order[1] != domain.Firefox {
		t.Errorf("wrong order: %v", order)
	}
}

func TestInstallBrowsers_OneFailureContinuesOthers(t *testing.T) {
	braveInst := &mocks.MockBrowserInstaller{
		BrowserID:  domain.Brave,
		InstallErr: errors.New("some error"),
	}
	firefoxInst := &mocks.MockBrowserInstaller{
		BrowserID: domain.Firefox,
	}

	uc := usecases.NewInstallBrowsersUseCase(
		map[domain.BrowserID]domain.BrowserInstaller{
			domain.Brave:   braveInst,
			domain.Firefox: firefoxInst,
		},
		noSleep,
	)
	results := uc.Execute([]domain.BrowserID{domain.Brave, domain.Firefox})

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
	id        domain.BrowserID
	err       error
	callCount int
}

func (c *countingInstaller) Install() error {
	c.callCount++
	return c.err
}
func (c *countingInstaller) IsInstalled() (bool, error) { return false, nil }
func (c *countingInstaller) ID() domain.BrowserID       { return c.id }

type orderTracker struct {
	id    domain.BrowserID
	order *[]domain.BrowserID
}

func (o *orderTracker) Install() error {
	*o.order = append(*o.order, o.id)
	return nil
}
func (o *orderTracker) IsInstalled() (bool, error) { return false, nil }
func (o *orderTracker) ID() domain.BrowserID       { return o.id }
