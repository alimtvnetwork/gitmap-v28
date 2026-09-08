package cmd

import (
	"strings"
	"testing"
)

func TestNginxStatus(t *testing.T) {
	err := runNginx([]string{})
	if err != nil {
		t.Errorf("expected runNginx with empty args to succeed, got %v", err)
	}
}

func TestNginxStatusSubcommand(t *testing.T) {
	err := runNginx([]string{"status"})
	if err != nil {
		t.Errorf("expected runNginx status to succeed, got %v", err)
	}

	errAlias := runNginx([]string{"st"})
	if errAlias != nil {
		t.Errorf("expected runNginx st to succeed, got %v", errAlias)
	}
}

func TestNginxUnknownSubcommand(t *testing.T) {
	err := runNginx([]string{"invalid-subcommand"})
	if err == nil {
		t.Errorf("expected error for invalid subcommand, got nil")
	}

	if !strings.Contains(err.Error(), "unknown nginx subcommand") {
		t.Errorf("expected unknown nginx subcommand error message, got %v", err)
	}
}

func TestNginxVHostDispatch(t *testing.T) {
	isHandled, err := dispatchNginxVHostSubcommand("vhost", []string{"unknown"})
	if !isHandled {
		t.Errorf("expected vhost to be handled by dispatchNginxVHostSubcommand")
	}

	if err == nil {
		t.Errorf("expected error for unknown vhost subcommand")
	}

	isHandledList, errList := dispatchNginxVHostSubcommand("list", []string{})
	if !isHandledList || errList != nil {
		t.Errorf("expected list to be handled and succeed, got %t, %v", isHandledList, errList)
	}
}

func TestNginxVHostAliases(t *testing.T) {
	isHandledCreate, _ := dispatchNginxVHostSubcommand("create", []string{"unknown-missing-args"})
	if !isHandledCreate {
		t.Errorf("expected create to be handled by dispatchNginxVHostSubcommand")
	}

	isHandledEnable, _ := dispatchNginxVHostSubcommand("enable", []string{})
	if !isHandledEnable {
		t.Errorf("expected enable to be handled by dispatchNginxVHostSubcommand")
	}

	isHandledDisable, _ := dispatchNginxVHostSubcommand("disable", []string{})
	if !isHandledDisable {
		t.Errorf("expected disable to be handled by dispatchNginxVHostSubcommand")
	}
}

func TestNginxDispatchTestAndReloadDryRun(t *testing.T) {
	errTest := runNginx([]string{"test", "--dry-run"})
	if errTest != nil {
		t.Errorf("expected test --dry-run to succeed, got %v", errTest)
	}

	errReload := runNginx([]string{"reload", "--dry-run"})
	if errReload != nil {
		t.Errorf("expected reload --dry-run to succeed, got %v", errReload)
	}
}
