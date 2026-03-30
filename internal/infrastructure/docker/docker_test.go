package docker_test

import (
	"errors"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/docker"
	"github.com/so-install/pkg/mocks"
)

const (
	repoStep = `echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/debian $(. /etc/os-release && echo $VERSION_CODENAME) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null`
)

func addSuccessResponses(m *mocks.MockExecutor, n int) {
	for i := 0; i < n; i++ {
		m.AddResponse("", "", nil)
	}
}

func TestDockerInstaller_ID(t *testing.T) {
	m := &mocks.MockExecutor{}
	inst := docker.NewDockerInstaller(m, "alice")
	if inst.ID() != domain.Docker {
		t.Errorf("expected domain.Docker, got %v", inst.ID())
	}
}

func TestDockerInstaller_IsInstalled_True(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("Docker version 24.0.0", "", nil)
	inst := docker.NewDockerInstaller(m, "")
	got, err := inst.IsInstalled()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got {
		t.Error("expected true when docker version exits 0")
	}
	if len(m.Calls) != 1 || m.Calls[0].Name != "docker" {
		t.Errorf("expected 'docker' call, got %v", m.Calls)
	}
}

func TestDockerInstaller_IsInstalled_False(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", errors.New("exit 1"))
	inst := docker.NewDockerInstaller(m, "")
	got, err := inst.IsInstalled()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got {
		t.Error("expected false when docker version exits non-zero")
	}
}

func TestDockerInstaller_Install_HappyPath_WithSudoUser(t *testing.T) {
	m := &mocks.MockExecutor{}
	// 13 steps: step1(remove) + steps 2-12 + step13(usermod)
	addSuccessResponses(m, 13)

	inst := docker.NewDockerInstaller(m, "alice")
	err := inst.Install()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Calls) != 13 {
		t.Fatalf("expected 13 calls, got %d: %v", len(m.Calls), m.Calls)
	}

	// Step 1: remove conflicting packages (apt remove)
	assertCall(t, m.Calls[0], "apt", "remove", "-y", "docker.io", "docker-doc", "docker-compose", "podman-docker", "containerd", "runc")
	// Step 2: apt update
	assertCall(t, m.Calls[1], "apt", "update")
	// Step 3: install prerequisites
	assertCall(t, m.Calls[2], "apt", "install", "-y", "ca-certificates", "curl")
	// Step 4: create keyrings dir
	assertCall(t, m.Calls[3], "install", "-m", "0755", "-d", "/etc/apt/keyrings")
	// Step 5: download GPG key
	assertCall(t, m.Calls[4], "curl", "-fsSL", "https://download.docker.com/linux/debian/gpg", "-o", "/etc/apt/keyrings/docker.asc")
	// Step 6: set permissions
	assertCall(t, m.Calls[5], "chmod", "a+r", "/etc/apt/keyrings/docker.asc")
	// Step 7: add repo
	assertCall(t, m.Calls[6], "sh", "-c", repoStep)
	// Step 8: apt update
	assertCall(t, m.Calls[7], "apt", "update")
	// Step 9: install docker packages
	assertCall(t, m.Calls[8], "apt", "install", "-y", "docker-ce", "docker-ce-cli", "containerd.io", "docker-buildx-plugin", "docker-compose-plugin")
	// Step 10: enable docker
	assertCall(t, m.Calls[9], "systemctl", "enable", "docker")
	// Step 11: enable containerd
	assertCall(t, m.Calls[10], "systemctl", "enable", "containerd")
	// Step 12: start docker
	assertCall(t, m.Calls[11], "systemctl", "start", "docker")
	// Step 13: usermod (sudoUser = "alice")
	assertCall(t, m.Calls[12], "usermod", "-aG", "docker", "alice")
}

func TestDockerInstaller_Install_HappyPath_NoSudoUser(t *testing.T) {
	m := &mocks.MockExecutor{}
	// 12 steps: step1(remove) + steps 2-12, no usermod
	addSuccessResponses(m, 12)

	inst := docker.NewDockerInstaller(m, "")
	err := inst.Install()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Calls) != 12 {
		t.Fatalf("expected 12 calls (no usermod), got %d: %v", len(m.Calls), m.Calls)
	}
	// Verify usermod was NOT called
	for _, call := range m.Calls {
		if call.Name == "usermod" {
			t.Error("usermod must not be called when sudoUser is empty")
		}
	}
}

func TestDockerInstaller_Install_Step1ErrorIgnored(t *testing.T) {
	m := &mocks.MockExecutor{}
	// Step 1 fails (packages not installed) — must be ignored
	m.AddResponse("", "E: Unable to locate package docker.io", errors.New("exit 100"))
	// Steps 2-12 succeed (no usermod — empty sudoUser)
	addSuccessResponses(m, 11)

	inst := docker.NewDockerInstaller(m, "")
	err := inst.Install()
	if err != nil {
		t.Fatalf("step 1 error must be ignored, got: %v", err)
	}
	if len(m.Calls) != 12 {
		t.Fatalf("expected 12 calls total, got %d", len(m.Calls))
	}
}

func TestDockerInstaller_Install_MidStepFailureStopsExecution(t *testing.T) {
	m := &mocks.MockExecutor{}
	// Step 1: ok (remove)
	m.AddResponse("", "", nil)
	// Step 2: ok (apt update)
	m.AddResponse("", "", nil)
	// Step 3: ok (install prerequisites)
	m.AddResponse("", "", nil)
	// Step 4: ok (create keyrings dir)
	m.AddResponse("", "", nil)
	// Step 5: FAIL (curl GPG key)
	m.AddResponse("", "curl: failed to connect", errors.New("exit 7"))

	inst := docker.NewDockerInstaller(m, "alice")
	err := inst.Install()
	if err == nil {
		t.Fatal("expected error when step 5 fails")
	}
	var instErr domain.InstallError
	if !errors.As(err, &instErr) {
		t.Errorf("expected domain.InstallError, got %T: %v", err, err)
	}
	if len(m.Calls) != 5 {
		t.Errorf("expected execution to stop at step 5, got %d calls", len(m.Calls))
	}
}

func TestDockerInstaller_Install_AptLockReturnsAptLockError(t *testing.T) {
	m := &mocks.MockExecutor{}
	// Steps 1-8 succeed
	addSuccessResponses(m, 8)
	// Step 9: apt install returns APT lock error
	m.AddResponse("", "E: Could not get lock /var/lib/dpkg/lock", errors.New("exit 100"))

	inst := docker.NewDockerInstaller(m, "")
	err := inst.Install()
	if err == nil {
		t.Fatal("expected error")
	}
	var aptErr domain.AptLockError
	if !errors.As(err, &aptErr) {
		t.Errorf("expected AptLockError, got %T: %v", err, err)
	}
}

// assertCall verifies a specific executor call by command name and arguments.
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
