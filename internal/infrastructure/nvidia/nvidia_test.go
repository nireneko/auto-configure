package nvidia_test

import (
	"errors"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/nvidia"
	"github.com/so-install/pkg/mocks"
)

// --- IsInstalled ---

func TestNvidiaInstaller_IsInstalled_ProprietaryFound(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", nil) // dpkg -s nvidia-driver → found
	inst := nvidia.NewNvidiaInstaller(m, &domain.OSInfo{})
	got, err := inst.IsInstalled()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got {
		t.Error("expected true when nvidia-driver is installed")
	}
	if len(m.Calls) != 1 {
		t.Errorf("expected 1 call, got %d", len(m.Calls))
	}
}

func TestNvidiaInstaller_IsInstalled_FreeFound(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", errors.New("not found")) // dpkg nvidia-driver → not found
	m.AddResponse("", "", nil)                     // dpkg firmware-nvidia-graphics → found
	inst := nvidia.NewNvidiaInstaller(m, &domain.OSInfo{})
	got, err := inst.IsInstalled()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got {
		t.Error("expected true when firmware-nvidia-graphics is installed")
	}
}

func TestNvidiaInstaller_IsInstalled_NoneFound(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", errors.New("not found"))
	m.AddResponse("", "", errors.New("not found"))
	inst := nvidia.NewNvidiaInstaller(m, &domain.OSInfo{})
	got, err := inst.IsInstalled()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got {
		t.Error("expected false when neither package is installed")
	}
}

// --- Install: error without SetOptions ---

func TestNvidiaInstaller_Install_ErrorWithoutSetOptions(t *testing.T) {
	inst := nvidia.NewNvidiaInstaller(&mocks.MockExecutor{}, &domain.OSInfo{})
	err := inst.Install()
	if err == nil {
		t.Fatal("expected error when SetOptions not called")
	}
	if err.Error() != "nvidia driver type not configured" {
		t.Errorf("unexpected error message: %v", err)
	}
}

// --- Install: Free driver ---

func TestNvidiaInstaller_Install_Free_SourcesAlreadyHaveNonFree(t *testing.T) {
	m := &mocks.MockExecutor{}
	// grep finds non-free → no sed needed
	m.AddResponse("", "", nil) // grep -q non-free
	m.AddResponse("", "", nil) // apt-get update
	m.AddResponse("", "", nil) // apt-get install firmware-nvidia-graphics

	inst := nvidia.NewNvidiaInstaller(m, &domain.OSInfo{})
	inst.SetOptions(domain.NvidiaFree, false)

	if err := inst.Install(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Calls) != 3 {
		t.Fatalf("expected 3 calls, got %d: %v", len(m.Calls), m.Calls)
	}
	assertCall(t, m.Calls[0], "grep", "-q", "non-free", "/etc/apt/sources.list.d/debian.sources")
	assertCall(t, m.Calls[1], "apt-get", "update")
	assertCall(t, m.Calls[2], "apt-get", "install", "-y", "firmware-nvidia-graphics")
}

func TestNvidiaInstaller_Install_Free_SourcesNeedUpdate(t *testing.T) {
	m := &mocks.MockExecutor{}
	// grep not found → sed needed
	m.AddResponse("", "", errors.New("not found")) // grep
	m.AddResponse("", "", nil)                     // sed
	m.AddResponse("", "", nil)                     // apt-get update
	m.AddResponse("", "", nil)                     // apt-get install

	inst := nvidia.NewNvidiaInstaller(m, &domain.OSInfo{})
	inst.SetOptions(domain.NvidiaFree, false)

	if err := inst.Install(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Calls) != 4 {
		t.Fatalf("expected 4 calls, got %d", len(m.Calls))
	}
	assertCall(t, m.Calls[1], "sed", "-i",
		"s/^Components: main$/Components: main contrib non-free non-free-firmware/",
		"/etc/apt/sources.list.d/debian.sources")
}

func TestNvidiaInstaller_Install_Free_NoCUDA(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.SetDefault("", "", nil)

	inst := nvidia.NewNvidiaInstaller(m, &domain.OSInfo{})
	inst.SetOptions(domain.NvidiaFree, true) // installCUDA=true must be ignored

	if err := inst.Install(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Verify no cuda-related call was made
	for _, c := range m.Calls {
		if c.Name == "apt-get" {
			for _, arg := range c.Args {
				if arg == "nvidia-cuda-dev" || arg == "nvidia-cuda-toolkit" || arg == "cuda-toolkit" {
					t.Errorf("unexpected CUDA install call in free driver mode: %v", c.Args)
				}
			}
		}
	}
}

// --- Install: Proprietary Debian ---

func TestNvidiaInstaller_Install_ProprietaryDebian_HappyPath(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", nil)             // grep → sources ok
	m.AddResponse("", "", nil)             // apt-get update
	m.AddResponse("6.1.0-31-amd64\n", "", nil) // uname -r
	m.AddResponse("", "", nil)             // apt-get install linux-headers
	m.AddResponse("", "", nil)             // apt-get install nvidia-kernel-dkms nvidia-driver
	m.AddResponse("", "", nil)             // systemctl enable nvidia-suspend
	m.AddResponse("", "", nil)             // systemctl enable nvidia-hibernate
	m.AddResponse("", "", nil)             // systemctl enable nvidia-resume
	// No wayland, no CUDA

	inst := nvidia.NewNvidiaInstaller(m, &domain.OSInfo{IsWayland: false})
	inst.SetOptions(domain.NvidiaProprietaryDebian, false)

	if err := inst.Install(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Calls) != 8 {
		t.Fatalf("expected 8 calls, got %d: %v", len(m.Calls), m.Calls)
	}
	assertCall(t, m.Calls[2], "uname", "-r")
	assertCall(t, m.Calls[3], "apt-get", "install", "-y", "linux-headers-6.1.0-31-amd64")
	assertCall(t, m.Calls[4], "apt-get", "install", "-y", "nvidia-kernel-dkms", "nvidia-driver")
	assertCall(t, m.Calls[5], "systemctl", "enable", "nvidia-suspend.service")
	assertCall(t, m.Calls[6], "systemctl", "enable", "nvidia-hibernate.service")
	assertCall(t, m.Calls[7], "systemctl", "enable", "nvidia-resume.service")
}

func TestNvidiaInstaller_Install_ProprietaryDebian_WithCUDA(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", nil)             // grep
	m.AddResponse("", "", nil)             // apt-get update
	m.AddResponse("6.1.0-31-amd64\n", "", nil) // uname -r
	m.AddResponse("", "", nil)             // headers
	m.AddResponse("", "", nil)             // nvidia driver
	m.AddResponse("", "", nil)             // systemctl ×3
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)             // apt-get install CUDA

	inst := nvidia.NewNvidiaInstaller(m, &domain.OSInfo{IsWayland: false})
	inst.SetOptions(domain.NvidiaProprietaryDebian, true)

	if err := inst.Install(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lastCall := m.Calls[len(m.Calls)-1]
	assertCall(t, lastCall, "apt-get", "install", "-y", "nvidia-cuda-dev", "nvidia-cuda-toolkit")
}

// --- Install: Proprietary Nvidia (official repo) ---

func TestNvidiaInstaller_Install_ProprietaryNvidia_HappyPath(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", nil)             // apt-get install prereqs
	m.AddResponse("", "", nil)             // sh -c curl | gpg
	m.AddResponse("", "", nil)             // sh -c printf | tee sources
	m.AddResponse("", "", nil)             // apt-get update
	m.AddResponse("6.1.0-31-amd64\n", "", nil) // uname -r
	m.AddResponse("", "", nil)             // apt-get install headers
	m.AddResponse("", "", nil)             // apt-get install cuda-drivers
	// No wayland, no CUDA

	inst := nvidia.NewNvidiaInstaller(m, &domain.OSInfo{IsWayland: false})
	inst.SetOptions(domain.NvidiaProprietaryNvidia, false)

	if err := inst.Install(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Calls) != 7 {
		t.Fatalf("expected 7 calls, got %d: %v", len(m.Calls), m.Calls)
	}
	assertCall(t, m.Calls[0], "apt-get", "install", "-y", "ca-certificates", "curl", "gpg")
	assertCall(t, m.Calls[4], "uname", "-r")
	assertCall(t, m.Calls[5], "apt-get", "install", "-y", "linux-headers-6.1.0-31-amd64")
	assertCall(t, m.Calls[6], "apt-get", "install", "-y", "cuda-drivers")
}

func TestNvidiaInstaller_Install_ProprietaryNvidia_WithCUDA(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", nil)             // prereqs
	m.AddResponse("", "", nil)             // gpg key
	m.AddResponse("", "", nil)             // sources
	m.AddResponse("", "", nil)             // apt-get update
	m.AddResponse("6.1.0-31-amd64\n", "", nil) // uname -r
	m.AddResponse("", "", nil)             // headers
	m.AddResponse("", "", nil)             // cuda-drivers
	m.AddResponse("", "", nil)             // cuda-toolkit

	inst := nvidia.NewNvidiaInstaller(m, &domain.OSInfo{IsWayland: false})
	inst.SetOptions(domain.NvidiaProprietaryNvidia, true)

	if err := inst.Install(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lastCall := m.Calls[len(m.Calls)-1]
	assertCall(t, lastCall, "apt-get", "install", "-y", "cuda-toolkit")
}

// --- Wayland config ---

func TestNvidiaInstaller_Install_Wayland_AppliedForProprietary(t *testing.T) {
	m := &mocks.MockExecutor{}
	// ProprietaryDebian + wayland
	m.AddResponse("", "", nil)             // grep
	m.AddResponse("", "", nil)             // apt-get update
	m.AddResponse("6.1.0-31-amd64\n", "", nil) // uname -r
	m.AddResponse("", "", nil)             // headers
	m.AddResponse("", "", nil)             // nvidia driver
	m.AddResponse("", "", nil)             // systemctl ×3
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)             // sh -c printf | tee modprobe
	m.AddResponse("", "", nil)             // update-initramfs
	m.AddResponse("", "", nil)             // update-grub

	inst := nvidia.NewNvidiaInstaller(m, &domain.OSInfo{IsWayland: true})
	inst.SetOptions(domain.NvidiaProprietaryDebian, false)

	if err := inst.Install(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Verify modprobe, update-initramfs and update-grub were called
	names := make([]string, len(m.Calls))
	for i, c := range m.Calls {
		names[i] = c.Name
	}
	found := func(name string) bool {
		for _, n := range names {
			if n == name {
				return true
			}
		}
		return false
	}
	if !found("sh") {
		t.Error("expected sh -c call for modprobe config")
	}
	if !found("update-initramfs") {
		t.Error("expected update-initramfs call")
	}
	if !found("update-grub") {
		t.Error("expected update-grub call")
	}
}

func TestNvidiaInstaller_Install_Wayland_NotAppliedForFree(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.SetDefault("", "", nil)

	inst := nvidia.NewNvidiaInstaller(m, &domain.OSInfo{IsWayland: true})
	inst.SetOptions(domain.NvidiaFree, false)

	if err := inst.Install(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, c := range m.Calls {
		if c.Name == "update-initramfs" || c.Name == "update-grub" {
			t.Errorf("unexpected %s call for free driver + wayland", c.Name)
		}
	}
}

// --- helpers ---

func assertCall(t *testing.T, call mocks.ExecutorCall, name string, args ...string) {
	t.Helper()
	if call.Name != name {
		t.Errorf("expected command %q, got %q", name, call.Name)
	}
	if len(call.Args) != len(args) {
		t.Errorf("command %q: expected args %v, got %v", name, args, call.Args)
		return
	}
	for i, a := range args {
		if call.Args[i] != a {
			t.Errorf("command %q arg[%d]: expected %q, got %q", name, i, a, call.Args[i])
		}
	}
}
