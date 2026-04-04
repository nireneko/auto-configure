package nvidia_test

import (
	"errors"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/nvidia"
	"github.com/so-install/pkg/mocks"
)

func TestNvidiaInstaller_IsInstalled_ProprietaryFound(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", nil)
	inst := nvidia.NewNvidiaInstaller(m, &domain.OSInfo{})
	got, err := inst.IsInstalled()
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	if !got { t.Error("expected true") }
}

func TestNvidiaInstaller_IsInstalled_FreeFound(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", errors.New("not found"))
	m.AddResponse("", "", nil)
	inst := nvidia.NewNvidiaInstaller(m, &domain.OSInfo{})
	got, err := inst.IsInstalled()
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	if !got { t.Error("expected true") }
}

func TestNvidiaInstaller_Install_ErrorWithoutSetOptions(t *testing.T) {
	inst := nvidia.NewNvidiaInstaller(&mocks.MockExecutor{}, &domain.OSInfo{})
	err := inst.Install()
	if err == nil { t.Fatal("expected error") }
}

func TestNvidiaInstaller_Install_Free(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", errors.New("not found"))
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	inst := nvidia.NewNvidiaInstaller(m, &domain.OSInfo{})
	inst.SetOptions(domain.NvidiaFree, false)
	if err := inst.Install(); err != nil { t.Fatalf("unexpected error: %v", err) }
}

func TestNvidiaInstaller_Install_ProprietaryDebian_Debian_Wayland_CUDA(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	m.AddResponse("6.1.0-31-amd64\n", "", nil)
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	inst := nvidia.NewNvidiaInstaller(m, &domain.OSInfo{ID: "debian", IsWayland: true})
	inst.SetOptions(domain.NvidiaProprietaryDebian, true)
	if err := inst.Install(); err != nil { t.Fatalf("unexpected error: %v", err) }
}

func TestNvidiaInstaller_Install_ProprietaryDebian_NonDebian(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", nil)
	inst := nvidia.NewNvidiaInstaller(m, &domain.OSInfo{ID: "ubuntu"})
	inst.SetOptions(domain.NvidiaProprietaryDebian, false)
	if err := inst.Install(); err != nil { t.Fatalf("unexpected error: %v", err) }
}

func TestNvidiaInstaller_Install_ProprietaryNvidia_Wayland_CUDA(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	m.AddResponse("6.1.0-31-amd64\n", "", nil)
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	inst := nvidia.NewNvidiaInstaller(m, &domain.OSInfo{ID: "debian", IsWayland: true})
	inst.SetOptions(domain.NvidiaProprietaryNvidia, true)
	if err := inst.Install(); err != nil { t.Fatalf("unexpected error: %v", err) }
}
