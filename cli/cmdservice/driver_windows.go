package cmdservice

import (
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type WindowsServiceDriver struct{}

// ResolveServiceDriver returns the Windows service controller driver.
func ResolveServiceDriver() ServiceDriver {
	return &WindowsServiceDriver{}
}

func (d *WindowsServiceDriver) ListServices() ([]ServiceInfo, error) {
	cmd := exec.Command("sc.exe", "query", "state=", "all")
	out, err := cmd.Output()
	if err != nil {
		return nil, apperror.WrapSimple(err, "sc.exe query state= all")
	}

	return parseWindowsServiceList(string(out)), nil
}

func parseWindowsServiceList(output string) []ServiceInfo {
	lines := strings.Split(output, "\n")
	var services []ServiceInfo
	var current ServiceInfo
	for _, l := range lines {
		t := strings.TrimSpace(l)
		if strings.HasPrefix(t, "SERVICE_NAME:") {
			current.Name = strings.TrimSpace(strings.TrimPrefix(t, "SERVICE_NAME:"))
		}
		if strings.HasPrefix(t, "STATE") && current.Name != "" {
			current.Status = parseWindowsServiceState(t)
			current.IsEnabled = current.Status == "RUNNING"
			services = append(services, current)
			current = ServiceInfo{}
		}
	}

	return services
}

func parseWindowsServiceState(stateLine string) string {
	parts := strings.Fields(stateLine)
	if len(parts) >= 4 {
		return parts[3]
	}

	return "UNKNOWN"
}

func (d *WindowsServiceDriver) GetService(name string) (*ServiceInfo, error) {
	cmd := exec.Command("sc.exe", "query", name)
	out, err := cmd.Output()
	if err != nil {
		return nil, apperror.WrapSimple(err, "sc.exe query "+name)
	}

	state := parseWindowsQueryOutput(string(out))

	return &ServiceInfo{Name: name, Status: state, IsEnabled: state == "RUNNING"}, nil
}

func parseWindowsQueryOutput(out string) string {
	lines := strings.Split(out, "\n")
	for _, l := range lines {
		t := strings.TrimSpace(l)
		if strings.HasPrefix(t, "STATE") {
			return parseWindowsServiceState(t)
		}
	}

	return "UNKNOWN"
}

func (d *WindowsServiceDriver) StartService(name string) error {
	cmd := exec.Command("sc.exe", "start", name)
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "sc.exe start "+name)
	}

	return nil
}

func (d *WindowsServiceDriver) StopService(name string) error {
	cmd := exec.Command("sc.exe", "stop", name)
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "sc.exe stop "+name)
	}

	return nil
}

func (d *WindowsServiceDriver) CreateService(name, execPath, description string) error {
	cmd := exec.Command("sc.exe", "create", name, "binPath=", execPath, "start=", "auto")
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "sc.exe create "+name)
	}

	return nil
}

func (d *WindowsServiceDriver) RemoveService(name string) error {
	_ = d.StopService(name)
	cmd := exec.Command("sc.exe", "delete", name)
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "sc.exe delete "+name)
	}

	return nil
}
