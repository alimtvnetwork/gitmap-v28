package cmdservice

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type DarwinServiceDriver struct{}

func (d *DarwinServiceDriver) ListServices() ([]ServiceInfo, error) {
	cmd := exec.Command("launchctl", "list")
	out, err := cmd.Output()
	if err != nil {
		return nil, apperror.WrapSimple(err, "launchctl list")
	}

	return parseDarwinServiceList(string(out)), nil
}

func parseDarwinServiceList(output string) []ServiceInfo {
	lines := strings.Split(output, "\n")
	var services []ServiceInfo
	for i, l := range lines {
		if i == 0 {
			continue
		}
		f := strings.Fields(l)
		if len(f) >= 3 {
			services = append(services, ServiceInfo{Name: f[2], Status: f[1], IsEnabled: true})
		}
	}

	return services
}

func (d *DarwinServiceDriver) GetService(name string) (*ServiceInfo, error) {
	cmd := exec.Command("launchctl", "list", name)
	out, err := cmd.Output()
	if err != nil {
		return &ServiceInfo{Name: name, Status: "inactive", IsEnabled: false}, nil
	}

	return &ServiceInfo{Name: name, Status: "running", IsEnabled: len(out) > 0}, nil
}

func (d *DarwinServiceDriver) StartService(name string) error {
	cmd := exec.Command("launchctl", "start", name)
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "launchctl start "+name)
	}

	return nil
}

func (d *DarwinServiceDriver) StopService(name string) error {
	cmd := exec.Command("launchctl", "stop", name)
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "launchctl stop "+name)
	}

	return nil
}

func (d *DarwinServiceDriver) CreateService(name, execPath, description string) error {
	home, _ := os.UserHomeDir()
	plistDir := filepath.Join(home, "Library", "LaunchAgents")
	_ = os.MkdirAll(plistDir, 0755)
	plistContent := fmt.Sprintf("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<plist version=\"1.0\">\n<dict>\n<key>Label</key>\n<string>%s</string>\n<key>ProgramArguments</key>\n<array>\n<string>%s</string>\n</array>\n<key>RunAtLoad</key>\n<true/>\n</dict>\n</plist>", name, execPath)
	plistPath := filepath.Join(plistDir, name+".plist")
	if err := os.WriteFile(plistPath, []byte(plistContent), 0644); err != nil {
		return apperror.WrapSimple(err, "write launchd plist")
	}

	_ = exec.Command("launchctl", "load", plistPath).Run()

	return nil
}

func (d *DarwinServiceDriver) RemoveService(name string) error {
	home, _ := os.UserHomeDir()
	plistPath := filepath.Join(home, "Library", "LaunchAgents", name+".plist")
	_ = exec.Command("launchctl", "unload", plistPath).Run()
	_ = os.Remove(plistPath)

	return nil
}
