package cmdservice

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type LinuxServiceDriver struct{}

func (d *LinuxServiceDriver) ListServices() ([]ServiceInfo, error) {
	cmd := exec.Command("systemctl", "list-unit-files", "--type=service", "--no-legend", "--no-pager")
	out, err := cmd.Output()
	if err != nil {
		return nil, apperror.WrapSimple(err, "systemctl list-unit-files")
	}

	return parseLinuxServiceList(string(out)), nil
}

func parseLinuxServiceList(output string) []ServiceInfo {
	lines := strings.Split(output, "\n")
	var services []ServiceInfo
	for _, l := range lines {
		f := strings.Fields(l)
		if len(f) >= 2 {
			name := strings.TrimSuffix(f[0], ".service")
			isEnabled := strings.EqualFold(f[1], "enabled")
			services = append(services, ServiceInfo{Name: name, IsEnabled: isEnabled, Status: f[1]})
		}
	}

	return services
}

func (d *LinuxServiceDriver) GetService(name string) (*ServiceInfo, error) {
	cmd := exec.Command("systemctl", "is-active", name)
	out, _ := cmd.Output()
	status := strings.TrimSpace(string(out))

	enCmd := exec.Command("systemctl", "is-enabled", name)
	enOut, _ := enCmd.Output()
	isEnabled := strings.TrimSpace(string(enOut)) == "enabled"

	return &ServiceInfo{Name: name, Status: status, IsEnabled: isEnabled}, nil
}

func (d *LinuxServiceDriver) StartService(name string) error {
	cmd := exec.Command("systemctl", "start", name)
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "systemctl start "+name)
	}

	return nil
}

func (d *LinuxServiceDriver) StopService(name string) error {
	cmd := exec.Command("systemctl", "stop", name)
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "systemctl stop "+name)
	}

	return nil
}

func (d *LinuxServiceDriver) CreateService(name, execPath, description string) error {
	unitContent := fmt.Sprintf("[Unit]\nDescription=%s\nAfter=network.target\n\n[Service]\nExecStart=%s\nRestart=always\n\n[Install]\nWantedBy=multi-user.target\n", description, execPath)
	unitPath := filepath.Join("/etc/systemd/system", name+".service")
	if err := os.WriteFile(unitPath, []byte(unitContent), 0644); err != nil {
		return apperror.WrapSimple(err, "write systemd unit")
	}

	_ = exec.Command("systemctl", "daemon-reload").Run()
	_ = exec.Command("systemctl", "enable", name).Run()

	return nil
}

func (d *LinuxServiceDriver) RemoveService(name string) error {
	_ = d.StopService(name)
	_ = exec.Command("systemctl", "disable", name).Run()
	unitPath := filepath.Join("/etc/systemd/system", name+".service")
	_ = os.Remove(unitPath)
	_ = exec.Command("systemctl", "daemon-reload").Run()

	return nil
}
