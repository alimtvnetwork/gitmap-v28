//go:build e2e

package e2e_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlemon"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// TestGitmapDaemon_SettingsToggle verifies programmatic enablement toggling in SQLite settings.
func TestGitmapDaemon_SettingsToggle(t *testing.T) {
	initialState := cmdlemon.IsDaemonEnabled()
	t.Cleanup(func() {
		_ = cmdlemon.SetDaemonEnabled(initialState)
	})

	errDisable := cmdlemon.SetDaemonEnabled(false)
	if errDisable != nil {
		t.Fatalf("failed to set daemon enabled false: %v", errDisable)
	}

	isDisabled := !cmdlemon.IsDaemonEnabled()
	if !isDisabled {
		t.Errorf("expected IsDaemonEnabled() to be false after disabling")
	}

	errEnable := cmdlemon.SetDaemonEnabled(true)
	if errEnable != nil {
		t.Fatalf("failed to set daemon enabled true: %v", errEnable)
	}

	isEnabled := cmdlemon.IsDaemonEnabled()
	if !isEnabled {
		t.Errorf("expected IsDaemonEnabled() to be true after enabling")
	}
}

// TestGitmapDaemon_CLISettingsToggle verifies CLI enable/disable subcommand operations.
func TestGitmapDaemon_CLISettingsToggle(t *testing.T) {
	initialState := cmdlemon.IsDaemonEnabled()
	t.Cleanup(func() {
		_ = cmdlemon.SetDaemonEnabled(initialState)
	})

	disableOut := captureStdout(t, func() {
		if err := cmdlemon.RunLemonCLI([]string{"disable"}); err != nil {
			t.Errorf("expected disable subcommand to succeed, got: %v", err)
		}
	})

	hasDisabledNotice := strings.Contains(disableOut, "DISABLED")
	if !hasDisabledNotice {
		t.Errorf("expected output to contain 'DISABLED', got: %s", disableOut)
	}

	isDisabled := !cmdlemon.IsDaemonEnabled()
	if !isDisabled {
		t.Errorf("expected daemon to be disabled in settings")
	}

	enableOut := captureStdout(t, func() {
		if err := cmdlemon.RunLemonCLI([]string{"enable"}); err != nil {
			t.Errorf("expected enable subcommand to succeed, got: %v", err)
		}
	})

	hasEnabledNotice := strings.Contains(enableOut, "ENABLED")
	if !hasEnabledNotice {
		t.Errorf("expected output to contain 'ENABLED', got: %s", enableOut)
	}

	isEnabled := cmdlemon.IsDaemonEnabled()
	if !isEnabled {
		t.Errorf("expected daemon to be enabled in settings")
	}
}

// TestGitmapDaemon_ExecutionSuppressionWhenDisabled verifies daemon run skips execution when disabled.
func TestGitmapDaemon_ExecutionSuppressionWhenDisabled(t *testing.T) {
	initialState := cmdlemon.IsDaemonEnabled()
	t.Cleanup(func() {
		_ = cmdlemon.SetDaemonEnabled(initialState)
	})

	_ = cmdlemon.SetDaemonEnabled(false)

	isPollInvoked := false
	cfg := cmdlemon.DaemonConfig{
		Interval: 5 * time.Second,
		IsOnce:   false,
		PollFunc: func(ctx context.Context) error {
			isPollInvoked = true
			return nil
		},
	}

	daemon := cmdlemon.NewLemonDaemon(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	runErr := daemon.Run(ctx)
	if runErr != nil {
		t.Fatalf("expected daemon.Run to return cleanly when disabled, got: %v", runErr)
	}

	if isPollInvoked {
		t.Errorf("expected poll function NOT to be invoked when daemon is disabled in settings")
	}
}

// TestGitmapDaemon_StatusOutputReflectsSettings verifies daemon status reports enabled state.
func TestGitmapDaemon_StatusOutputReflectsSettings(t *testing.T) {
	initialState := cmdlemon.IsDaemonEnabled()
	t.Cleanup(func() {
		_ = cmdlemon.SetDaemonEnabled(initialState)
	})

	_ = cmdlemon.SetDaemonEnabled(true)
	enabledOut := captureStdout(t, func() {
		_ = cmdlemon.RunLemonCLI([]string{"status"})
	})

	hasEnabledLine := strings.Contains(enabledOut, "Enabled:      true")
	if !hasEnabledLine {
		t.Errorf("expected status output to include 'Enabled:      true', got: %s", enabledOut)
	}

	_ = cmdlemon.SetDaemonEnabled(false)
	disabledOut := captureStdout(t, func() {
		_ = cmdlemon.RunLemonCLI([]string{"status"})
	})

	hasDisabledLine := strings.Contains(disabledOut, "Enabled:      false")
	if !hasDisabledLine {
		t.Errorf("expected status output to include 'Enabled:      false', got: %s", disabledOut)
	}
}

// TestGitmapDaemon_CurrentPlatformServiceConfig verifies real service config templates on host OS.
func TestGitmapDaemon_CurrentPlatformServiceConfig(t *testing.T) {
	mgr := cmdlemon.NewServiceManager()
	details, err := mgr.Execute(cmdlemon.ServiceActionTypeStatus)
	if err != nil {
		t.Fatalf("unexpected error querying service status: %v", err)
	}
	if details == nil {
		t.Fatal("expected non-nil service details")
	}
	if !strings.Contains(strings.ToLower(details.Name), "daemon") && !strings.Contains(strings.ToLower(details.Name), "lemon") {
		t.Errorf("expected service name to contain daemon or lemon, got %q", details.Name)
	}
}

// TestGitmapDaemon_SettingsKeyAliases verifies fallback alias keys like gitmap-daemon.enabled.
func TestGitmapDaemon_SettingsKeyAliases(t *testing.T) {
	initialState := cmdlemon.IsDaemonEnabled()
	t.Cleanup(func() {
		_ = cmdlemon.SetDaemonEnabled(initialState)
	})

	dbConn, err := store.OpenDefault()
	if err != nil {
		t.Fatalf("failed to open default DB: %v", err)
	}
	defer dbConn.Close()

	_ = dbConn.DeleteSetting(cmdlemon.DaemonSettingKey)

	_ = dbConn.SetSetting("gitmap-daemon.enabled", "false")
	if cmdlemon.IsDaemonEnabled() {
		t.Errorf("expected IsDaemonEnabled() to be false when alias gitmap-daemon.enabled=false")
	}

	_ = dbConn.SetSetting("gitmap-daemon.enabled", "true")
	if !cmdlemon.IsDaemonEnabled() {
		t.Errorf("expected IsDaemonEnabled() to be true when alias gitmap-daemon.enabled=true")
	}

	_ = dbConn.DeleteSetting("gitmap-daemon.enabled")
}
