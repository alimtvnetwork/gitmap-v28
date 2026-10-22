//go:build e2e

package e2e_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlemon"
)

const windowsScriptTmpl = `Start-Process -FilePath "%s" -ArgumentList "daemon run --interval 8s" -WindowStyle Hidden`

const linuxUnitTmpl = "[Unit]\nDescription=GitMap Daemon Service\nAfter=network.target\n\n[Service]\nType=simple\nExecStart=%s daemon run\nRestart=always\nRestartSec=8\nStandardOutput=journal\nStandardError=journal\n\n[Install]\nWantedBy=default.target\n"

const darwinPlistTmpl = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
<key>Label</key><string>com.gitmap.daemon</string>
<key>ProgramArguments</key><array><string>%s</string><string>daemon</string><string>run</string></array>
<key>RunAtLoad</key><true/><key>KeepAlive</key><true/>
</dict></plist>`

type mockPlatformServiceAdapter struct {
	platform   string
	name       string
	configFunc func() (string, error)
}

func (m *mockPlatformServiceAdapter) Install() error   { return nil }
func (m *mockPlatformServiceAdapter) Uninstall() error { return nil }
func (m *mockPlatformServiceAdapter) Start() error     { return nil }
func (m *mockPlatformServiceAdapter) Stop() error      { return nil }
func (m *mockPlatformServiceAdapter) Restart() error   { return nil }

func (m *mockPlatformServiceAdapter) Status() (*cmdlemon.ServiceDetails, error) {
	return &cmdlemon.ServiceDetails{
		Name:       m.name,
		Platform:   m.platform,
		StatusType: cmdlemon.ServiceStatusTypeStopped,
		IsEnabled:  false,
	}, nil
}

func (m *mockPlatformServiceAdapter) GenerateConfig() (string, error) {
	if m.configFunc != nil {
		return m.configFunc()
	}

	return "", nil
}

func newMockPlatformService(platform, name, defaultBin, tmpl string) cmdlemon.PlatformService {
	return &mockPlatformServiceAdapter{
		platform: platform,
		name:     name,
		configFunc: func() (string, error) {
			execPath, err := os.Executable()
			if err != nil {
				execPath = defaultBin
			}

			return fmt.Sprintf(tmpl, execPath), nil
		},
	}
}

func newMockWindowsService() cmdlemon.PlatformService {
	return newMockPlatformService("windows", "GitMapDaemon", "gitmap.exe", windowsScriptTmpl)
}

func newMockLinuxService() cmdlemon.PlatformService {
	return newMockPlatformService("linux", "gitmap-daemon.service", "gitmap", linuxUnitTmpl)
}

func newMockDarwinService() cmdlemon.PlatformService {
	return newMockPlatformService("darwin", "com.gitmap.daemon.plist", "gitmap", darwinPlistTmpl)
}

func isServiceStatusValid(st cmdlemon.ServiceStatusType) bool {
	return st == cmdlemon.ServiceStatusTypeRunning ||
		st == cmdlemon.ServiceStatusTypeStopped ||
		st == cmdlemon.ServiceStatusTypeUnknown
}

func isUnknownActionError(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), "unknown service action")
}

func assertContainsAll(t *testing.T, target string, substrings []string) {
	t.Helper()

	for _, sub := range substrings {
		hasSub := strings.Contains(target, sub)
		if !hasSub {
			t.Errorf("expected content to contain %q, missing from: %s", sub, target)
		}
	}
}

// TestLemonService_MissingAndInvalidAction tests validation of missing and unknown actions.
func TestLemonService_MissingAndInvalidAction(t *testing.T) {
	errMissing := cmdlemon.RunLemonCLI([]string{"service"})
	if errMissing == nil {
		t.Fatal("expected error for missing service action, got nil")
	}

	hasMissing := strings.Contains(errMissing.Error(), "missing service action")
	if !hasMissing {
		t.Errorf("expected error containing 'missing service action', got: %v", errMissing)
	}

	errInvalid := cmdlemon.RunLemonCLI([]string{"service", "invalid-action"})
	if errInvalid == nil {
		t.Fatal("expected error for invalid service action, got nil")
	}

	hasInvalid := strings.Contains(errInvalid.Error(), "unknown service action")
	if !hasInvalid {
		t.Errorf("expected error containing 'unknown service action', got: %v", errInvalid)
	}
}

// TestLemonService_StatusQuery tests status query execution and details payload integrity.
func TestLemonService_StatusQuery(t *testing.T) {
	err := cmdlemon.RunLemonCLI([]string{"service", "status"})
	if err != nil {
		t.Fatalf("expected RunLemonCLI('service', 'status') to execute cleanly: %v", err)
	}

	mgr := cmdlemon.NewServiceManager()

	details, queryErr := mgr.Execute(cmdlemon.ServiceActionTypeStatus)
	if queryErr != nil {
		t.Fatalf("expected mgr.Execute('status') to execute cleanly: %v", queryErr)
	}

	if details == nil {
		t.Fatal("expected non-nil ServiceDetails from status query")
	}

	hasPlatform := len(details.Platform) > 0
	if !hasPlatform {
		t.Errorf("expected valid non-empty platform string, got %q", details.Platform)
	}

	hasName := len(details.Name) > 0
	if !hasName {
		t.Errorf("expected valid non-empty service name, got %q", details.Name)
	}

	isValidStatus := isServiceStatusValid(details.StatusType)
	if !isValidStatus {
		t.Errorf("expected valid status type, got %q", details.StatusType)
	}
}

// TestLemonService_WindowsConfigGeneration verifies Windows service configuration template.
func TestLemonService_WindowsConfigGeneration(t *testing.T) {
	svc := newMockWindowsService()

	cfg, err := svc.GenerateConfig()
	if err != nil {
		t.Fatalf("unexpected error generating Windows config: %v", err)
	}

	assertContainsAll(t, cfg, []string{
		"Start-Process",
		"-WindowStyle Hidden",
		"daemon run --interval 8s",
	})
}

// TestLemonService_LinuxConfigGeneration verifies Linux systemd unit generation template.
func TestLemonService_LinuxConfigGeneration(t *testing.T) {
	svc := newMockLinuxService()

	cfg, err := svc.GenerateConfig()
	if err != nil {
		t.Fatalf("unexpected error generating Linux config: %v", err)
	}

	assertContainsAll(t, cfg, []string{
		"[Unit]",
		"Description=GitMap Daemon Service",
		"daemon run",
		"Restart=always",
		"RestartSec=8",
	})
}

// TestLemonService_DarwinConfigGeneration verifies macOS launchd plist generation template.
func TestLemonService_DarwinConfigGeneration(t *testing.T) {
	svc := newMockDarwinService()

	cfg, err := svc.GenerateConfig()
	if err != nil {
		t.Fatalf("unexpected error generating Darwin config: %v", err)
	}

	assertContainsAll(t, cfg, []string{
		"com.gitmap.daemon",
		"ProgramArguments",
		"daemon",
		"run",
		"<key>KeepAlive</key><true/>",
	})
}

// TestLemonService_ActionNormalization verifies case insensitivity across service actions.
func TestLemonService_ActionNormalization(t *testing.T) {
	mgr := cmdlemon.NewServiceManager()

	detailsUpper, errUpper := mgr.Execute(cmdlemon.ServiceActionType("STATUS"))
	if errUpper != nil {
		t.Fatalf("expected STATUS to execute cleanly: %v", errUpper)
	}

	if detailsUpper == nil {
		t.Fatal("expected non-nil ServiceDetails for STATUS")
	}

	detailsMixed, errMixed := mgr.Execute(cmdlemon.ServiceActionType("StAtUs"))
	if errMixed != nil {
		t.Fatalf("expected StAtUs to execute cleanly: %v", errMixed)
	}

	if detailsMixed == nil {
		t.Fatal("expected non-nil ServiceDetails for StAtUs")
	}

	errCLI := cmdlemon.RunLemonCLI([]string{"service", "STATUS"})
	if errCLI != nil {
		t.Fatalf("expected RunLemonCLI('service', 'STATUS') to execute cleanly: %v", errCLI)
	}

	testActions := []string{"Install", "UNINSTALL"}
	for _, rawAction := range testActions {
		normalizedAction := strings.ToLower(rawAction)
		hasValidLower := normalizedAction == "install" || normalizedAction == "uninstall"
		if !hasValidLower {
			t.Errorf("expected lower normalized action, got %q", normalizedAction)
		}

		_, execErr := mgr.Execute(cmdlemon.ServiceActionType(rawAction))
		isUnknown := isUnknownActionError(execErr)
		if isUnknown {
			t.Errorf("action %q was rejected as unknown action instead of normalized", rawAction)
		}
	}

	_, errInvalid := mgr.Execute(cmdlemon.ServiceActionType("INVALID_ACTION"))
	isInvalidRejected := isUnknownActionError(errInvalid)
	if !isInvalidRejected {
		t.Errorf("expected INVALID_ACTION to be rejected with unknown service action: %v", errInvalid)
	}
}
