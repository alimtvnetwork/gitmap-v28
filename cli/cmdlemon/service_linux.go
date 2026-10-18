//go:build linux

package cmdlemon

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const LinuxUnitName = "gitmap-lemon.service"

type LinuxService struct{}

func resolvePlatformService() PlatformService {
	return &LinuxService{}
}

func (s *LinuxService) GenerateConfig() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		execPath = "gitmap"
	}
	unitTmpl := "[Unit]\nDescription=GitMap Lemon Service\nAfter=network.target\n\n[Service]\nType=simple\nExecStart=%s lemon run\nRestart=always\nRestartSec=8\nStandardOutput=journal\nStandardError=journal\n\n[Install]\nWantedBy=default.target\n"

	return fmt.Sprintf(unitTmpl, execPath), nil
}

func (s *LinuxService) Install() error {
	unit, err := s.GenerateConfig()
	if err != nil {
		return err
	}

	dir := filepath.Join(os.Getenv("HOME"), ".config/systemd/user")
	_ = os.MkdirAll(dir, 0o755)
	target := filepath.Join(dir, LinuxUnitName)
	if err := os.WriteFile(target, []byte(unit), 0o644); err != nil {
		return err
	}

	_ = exec.Command("systemctl", "--user", "daemon-reload").Run()

	return exec.Command("systemctl", "--user", "enable", LinuxUnitName).Run()
}

func (s *LinuxService) Uninstall() error {
	_ = exec.Command("systemctl", "--user", "stop", LinuxUnitName).Run()
	_ = exec.Command("systemctl", "--user", "disable", LinuxUnitName).Run()
	target := filepath.Join(os.Getenv("HOME"), ".config/systemd/user", LinuxUnitName)
	_ = os.Remove(target)

	return exec.Command("systemctl", "--user", "daemon-reload").Run()
}

func (s *LinuxService) Start() error {
	return exec.Command("systemctl", "--user", "start", LinuxUnitName).Run()
}

func (s *LinuxService) Stop() error {
	return exec.Command("systemctl", "--user", "stop", LinuxUnitName).Run()
}

func (s *LinuxService) Restart() error {
	return exec.Command("systemctl", "--user", "restart", LinuxUnitName).Run()
}

func (s *LinuxService) Status() (*ServiceDetails, error) {
	out, err := exec.Command("systemctl", "--user", "is-active", LinuxUnitName).Output()
	state := strings.TrimSpace(string(out))
	statusType := ServiceStatusTypeStopped
	if err == nil && state == "active" {
		statusType = ServiceStatusTypeRunning
	}

	return &ServiceDetails{
		Name:       LinuxUnitName,
		Platform:   "linux",
		StatusType: statusType,
		IsEnabled:  statusType == ServiceStatusTypeRunning,
		ConfigPath: filepath.Join(".config/systemd/user", LinuxUnitName),
	}, nil
}
