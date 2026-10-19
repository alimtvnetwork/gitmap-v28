//go:build windows

package cmdlemon

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	WindowsServiceName       = "GitMapDaemon"
	WindowsLegacyServiceName = "GitMapLemon"
)

type WindowsService struct{}

func resolvePlatformService() PlatformService {
	return &WindowsService{}
}

func (s *WindowsService) GenerateConfig() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		execPath = "gitmap.exe"
	}
	script := fmt.Sprintf(`Start-Process -FilePath "%s" -ArgumentList "daemon run --interval 8s" -WindowStyle Hidden`, execPath)
	return script, nil
}

func (s *WindowsService) Install() error {
	execPath, err := os.Executable()
	if err != nil {
		execPath = "gitmap.exe"
	}
	binPath := fmt.Sprintf(`\"%s\" daemon run`, execPath)
	cmd := exec.Command("sc.exe", "create", WindowsServiceName, "binPath=", binPath, "start=", "auto")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to register Windows service: %w (output: %s)", err, string(out))
	}
	return nil
}

func (s *WindowsService) Uninstall() error {
	_ = s.Stop()
	cmd := exec.Command("sc.exe", "delete", WindowsServiceName)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to delete Windows service: %w (output: %s)", err, string(out))
	}
	return nil
}

func (s *WindowsService) Start() error {
	return exec.Command("sc.exe", "start", WindowsServiceName).Run()
}

func (s *WindowsService) Stop() error {
	return exec.Command("sc.exe", "stop", WindowsServiceName).Run()
}

func (s *WindowsService) Restart() error {
	_ = s.Stop()
	return s.Start()
}

func (s *WindowsService) Status() (*ServiceDetails, error) {
	out, err := exec.Command("sc.exe", "query", WindowsServiceName).CombinedOutput()
	isRunning := err == nil && strings.Contains(string(out), "RUNNING")
	if isRunning {
		return makeWindowsServiceDetails(ServiceStatusTypeRunning), nil
	}

	outLegacy, errLegacy := exec.Command("sc.exe", "query", WindowsLegacyServiceName).CombinedOutput()
	isLegacyRunning := errLegacy == nil && strings.Contains(string(outLegacy), "RUNNING")
	if isLegacyRunning {
		return makeWindowsServiceDetails(ServiceStatusTypeRunning), nil
	}

	return makeWindowsServiceDetails(ServiceStatusTypeStopped), nil
}

func makeWindowsServiceDetails(statusType ServiceStatusType) *ServiceDetails {
	return &ServiceDetails{
		Name:       WindowsServiceName,
		Platform:   "windows",
		StatusType: statusType,
		IsEnabled:  statusType == ServiceStatusTypeRunning,
		ConfigPath: filepath.Join("System32", WindowsServiceName),
	}
}
