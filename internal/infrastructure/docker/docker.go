package docker

import (
	"fmt"

	"github.com/so-install/internal/core/domain"
)

const (
	dockerGPGURL    = "https://download.docker.com/linux/debian/gpg"
	dockerGPGPath   = "/etc/apt/keyrings/docker.asc"
	dockerRepoPath  = "/etc/apt/sources.list.d/docker.list"
	dockerRepoEntry = `echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/debian $(. /etc/os-release && echo $VERSION_CODENAME) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null`
)

// DockerInstaller installs Docker CE from the official Docker repository.
type DockerInstaller struct {
	executor domain.Executor
	sudoUser string
}

// NewDockerInstaller creates a new DockerInstaller.
// sudoUser is the real user to add to the docker group (typically os.Getenv("SUDO_USER")).
// If sudoUser is empty, the usermod step is skipped.
func NewDockerInstaller(executor domain.Executor, sudoUser string) *DockerInstaller {
	return &DockerInstaller{executor: executor, sudoUser: sudoUser}
}

var _ domain.SoftwareInstaller = (*DockerInstaller)(nil)

// ID returns the SoftwareID for Docker.
func (d *DockerInstaller) ID() domain.SoftwareID { return domain.Docker }

// IsInstalled checks if docker is already installed.
func (d *DockerInstaller) IsInstalled() (bool, error) {
	_, _, err := d.executor.Execute("docker", "version")
	return err == nil, nil
}

// Install installs Docker CE from the official Docker repository.
// Step 1 (remove conflicting packages) is best-effort — errors are ignored.
// All subsequent steps stop execution on failure.
func (d *DockerInstaller) Install() error {
	// Step 1: remove conflicting packages — ignore errors (may not be installed)
	d.executor.Execute("apt", "remove", "-y", "docker.io", "docker-doc", "docker-compose", "podman-docker", "containerd", "runc") //nolint:errcheck

	// Steps 2-11: sequential — stop on first failure
	steps := [][]string{
		{"apt", "update"},
		{"apt", "install", "-y", "ca-certificates", "curl"},
		{"install", "-m", "0755", "-d", "/etc/apt/keyrings"},
		{"curl", "-fsSL", dockerGPGURL, "-o", dockerGPGPath},
		{"chmod", "a+r", dockerGPGPath},
		{"sh", "-c", dockerRepoEntry},
		{"apt", "update"},
		{"apt", "install", "-y", "docker-ce", "docker-ce-cli", "containerd.io", "docker-buildx-plugin", "docker-compose-plugin"},
		{"systemctl", "enable", "docker"},
		{"systemctl", "enable", "containerd"},
		{"systemctl", "start", "docker"},
	}
	for _, step := range steps {
		_, stderr, err := d.executor.Execute(step[0], step[1:]...)
		if err != nil {
			return domain.WrapInstallError("docker", step[0], step[1:], "", stderr, err)
		}
	}

	// Step 12: add real user to docker group (skip if sudoUser is empty)
	if d.sudoUser == "" {
		fmt.Println("warning: SUDO_USER not set — skipping usermod. Run manually: usermod -aG docker <your-user>")
		return nil
	}
	_, stderr, err := d.executor.Execute("usermod", "-aG", "docker", d.sudoUser)
	if err != nil {
		return domain.WrapInstallError("docker", "usermod", []string{"-aG", "docker", d.sudoUser}, "", stderr, err)
	}
	return nil
}
