package nvidia

import (
	"fmt"
	"strings"

	"github.com/so-install/internal/core/domain"
)

const (
	nvidiaCUDARepoKey     = "https://developer.download.nvidia.com/compute/cuda/repos/debian13/x86_64/8793F200.pub"
	nvidiaCUDARepoKeyPath = "/usr/share/keyrings/nvidia-drivers.gpg"
	nvidiaCUDASourcesPath = "/etc/apt/sources.list.d/nvidia-drivers.sources"
	nvidiaModprobePath    = "/etc/modprobe.d/nvidia-drm.conf"
	debianSourcesPath     = "/etc/apt/sources.list.d/debian.sources"

	// Raw \n sequences — interpreted by printf in shell.
	nvidiaCUDASourcesContent = `Types: deb\nURIs: https://developer.download.nvidia.com/compute/cuda/repos/debian13/x86_64/\nSuites: /\nComponents:\nSigned-By: /usr/share/keyrings/nvidia-drivers.gpg\n`
	nvidiaModprobeContent    = `options nvidia-drm modeset=1\noptions nvidia NVreg_PreserveVideoMemoryAllocations=1\n`
)

// NvidiaInstaller installs Nvidia GPU drivers on Debian 13.
type NvidiaInstaller struct {
	executor    domain.Executor
	osInfo      *domain.OSInfo
	driverType  domain.NvidiaDriverType
	installCUDA bool
	configured  bool
}

// NewNvidiaInstaller creates a new NvidiaInstaller.
func NewNvidiaInstaller(executor domain.Executor, osInfo *domain.OSInfo) *NvidiaInstaller {
	return &NvidiaInstaller{
		executor: executor,
		osInfo:   osInfo,
	}
}

// SetOptions configures the driver type and CUDA preference before installation.
func (n *NvidiaInstaller) SetOptions(driverType domain.NvidiaDriverType, installCUDA bool) {
	n.driverType = driverType
	n.installCUDA = installCUDA
	n.configured = true
}

var _ domain.SoftwareInstaller = (*NvidiaInstaller)(nil)

// ID returns the SoftwareID for Nvidia drivers.
func (n *NvidiaInstaller) ID() domain.SoftwareID { return domain.NvidiaDrivers }

// IsInstalled checks if any Nvidia driver package is already installed.
func (n *NvidiaInstaller) IsInstalled() (bool, error) {
	_, _, err := n.executor.Execute("dpkg", "-s", "nvidia-driver")
	if err == nil {
		return true, nil
	}
	_, _, err = n.executor.Execute("dpkg", "-s", "firmware-nvidia-graphics")
	return err == nil, nil
}

// Install installs the configured Nvidia driver.
func (n *NvidiaInstaller) Install() error {
	if !n.configured {
		return fmt.Errorf("nvidia driver type not configured")
	}
	switch n.driverType {
	case domain.NvidiaFree:
		return n.installFree()
	case domain.NvidiaProprietaryDebian:
		return n.installProprietaryDebian()
	case domain.NvidiaProprietaryNvidia:
		return n.installProprietaryNvidia()
	default:
		return fmt.Errorf("unknown nvidia driver type: %s", n.driverType)
	}
}

func (n *NvidiaInstaller) installFree() error {
	if err := n.enableNonFreeSources(); err != nil {
		return err
	}
	return n.runSteps([][]string{
		{"apt-get", "update"},
		{"apt-get", "install", "-y", "firmware-nvidia-graphics"},
	})
}

func (n *NvidiaInstaller) installProprietaryDebian() error {
	if err := n.enableNonFreeSources(); err != nil {
		return err
	}
	if err := n.runSteps([][]string{{"apt-get", "update"}}); err != nil {
		return err
	}
	if err := n.installKernelHeaders(); err != nil {
		return err
	}
	if err := n.runSteps([][]string{
		{"apt-get", "install", "-y", "nvidia-kernel-dkms", "nvidia-driver"},
		{"systemctl", "enable", "nvidia-suspend.service"},
		{"systemctl", "enable", "nvidia-hibernate.service"},
		{"systemctl", "enable", "nvidia-resume.service"},
	}); err != nil {
		return err
	}
	if err := n.applyWaylandConfig(); err != nil {
		return err
	}
	if n.installCUDA {
		return n.runSteps([][]string{
			{"apt-get", "install", "-y", "nvidia-cuda-dev", "nvidia-cuda-toolkit"},
		})
	}
	return nil
}

func (n *NvidiaInstaller) installProprietaryNvidia() error {
	if err := n.runSteps([][]string{
		{"apt-get", "install", "-y", "ca-certificates", "curl", "gpg"},
		{"sh", "-c", "curl -fsSL " + nvidiaCUDARepoKey + " | gpg --dearmor -o " + nvidiaCUDARepoKeyPath},
		{"sh", "-c", "printf '" + nvidiaCUDASourcesContent + "' | tee " + nvidiaCUDASourcesPath},
		{"apt-get", "update"},
	}); err != nil {
		return err
	}
	if err := n.installKernelHeaders(); err != nil {
		return err
	}
	if err := n.runSteps([][]string{
		{"apt-get", "install", "-y", "cuda-drivers"},
	}); err != nil {
		return err
	}
	if err := n.applyWaylandConfig(); err != nil {
		return err
	}
	if n.installCUDA {
		return n.runSteps([][]string{
			{"apt-get", "install", "-y", "cuda-toolkit"},
		})
	}
	return nil
}

func (n *NvidiaInstaller) enableNonFreeSources() error {
	_, _, err := n.executor.Execute("grep", "-q", "non-free", debianSourcesPath)
	if err != nil {
		// non-free not present — add it
		_, stderr, serr := n.executor.Execute("sed", "-i",
			"s/^Components: main$/Components: main contrib non-free non-free-firmware/",
			debianSourcesPath)
		if serr != nil {
			return domain.WrapInstallError("nvidia-drivers", "sed", nil, "", stderr, serr)
		}
	}
	return nil
}

func (n *NvidiaInstaller) installKernelHeaders() error {
	stdout, _, err := n.executor.Execute("uname", "-r")
	if err != nil {
		return fmt.Errorf("failed to detect kernel version: %w", err)
	}
	version := strings.TrimSpace(stdout)
	_, stderr, err := n.executor.Execute("apt-get", "install", "-y", "linux-headers-"+version)
	if err != nil {
		return domain.WrapInstallError("nvidia-drivers", "apt-get", []string{"install", "-y", "linux-headers-" + version}, "", stderr, err)
	}
	return nil
}

func (n *NvidiaInstaller) applyWaylandConfig() error {
	if n.osInfo == nil || !n.osInfo.IsWayland {
		return nil
	}
	return n.runSteps([][]string{
		{"sh", "-c", "printf '" + nvidiaModprobeContent + "' | tee " + nvidiaModprobePath},
		{"update-initramfs", "-u", "-k", "all"},
		{"update-grub"},
	})
}

func (n *NvidiaInstaller) runSteps(steps [][]string) error {
	for _, step := range steps {
		_, stderr, err := n.executor.Execute(step[0], step[1:]...)
		if err != nil {
			return domain.WrapInstallError("nvidia-drivers", step[0], step[1:], "", stderr, err)
		}
	}
	return nil
}
