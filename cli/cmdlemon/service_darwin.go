//go:build darwin

package cmdlemon

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	DarwinPlistName       = "com.gitmap.daemon.plist"
	DarwinLegacyPlistName = "com.gitmap.lemon.plist"
)

type DarwinService struct{}

func resolvePlatformService() PlatformService {
	return &DarwinService{}
}

func (s *DarwinService) GenerateConfig() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		execPath = "gitmap"
	}
	plistTmpl := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
<key>Label</key><string>com.gitmap.daemon</string>
<key>ProgramArguments</key><array><string>%s</string><string>daemon</string><string>run</string></array>
<key>RunAtLoad</key><true/><key>KeepAlive</key><true/>
</dict></plist>`

	return fmt.Sprintf(plistTmpl, execPath), nil
}

func (s *DarwinService) Install() error {
	plist, err := s.GenerateConfig()
	if err != nil {
		return err
	}

	dir := filepath.Join(os.Getenv("HOME"), "Library/LaunchAgents")
	_ = os.MkdirAll(dir, 0o755)
	target := filepath.Join(dir, DarwinPlistName)
	if err := os.WriteFile(target, []byte(plist), 0o644); err != nil {
		return err
	}

	return exec.Command("launchctl", "load", target).Run()
}

func (s *DarwinService) Uninstall() error {
	target := filepath.Join(os.Getenv("HOME"), "Library/LaunchAgents", DarwinPlistName)
	_ = exec.Command("launchctl", "unload", target).Run()

	return os.Remove(target)
}

func (s *DarwinService) Start() error {
	target := filepath.Join(os.Getenv("HOME"), "Library/LaunchAgents", DarwinPlistName)
	return exec.Command("launchctl", "start", target).Run()
}

func (s *DarwinService) Stop() error {
	target := filepath.Join(os.Getenv("HOME"), "Library/LaunchAgents", DarwinPlistName)
	return exec.Command("launchctl", "stop", target).Run()
}

func (s *DarwinService) Restart() error {
	_ = s.Stop()
	return s.Start()
}

func (s *DarwinService) Status() (*ServiceDetails, error) {
	out, err := exec.Command("launchctl", "list").CombinedOutput()
	statusType := ServiceStatusTypeStopped
	outStr := string(out)
	isRunning := err == nil && (strings.Contains(outStr, "com.gitmap.daemon") || strings.Contains(outStr, "com.gitmap.lemon"))
	if isRunning {
		statusType = ServiceStatusTypeRunning
	}
	return &ServiceDetails{
		Name:       DarwinPlistName,
		Platform:   "darwin",
		StatusType: statusType,
		IsEnabled:  statusType == ServiceStatusTypeRunning,
		ConfigPath: filepath.Join("Library/LaunchAgents", DarwinPlistName),
	}, nil
}
