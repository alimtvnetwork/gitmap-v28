package cmdos

import (
	"testing"
)

func TestParseAutoLoginArgs(t *testing.T) {
	args := []string{"-u", "admin", "-d", "CORP", "-p", "pass123"}
	cfg := parseAutoLoginArgs(args)

	if cfg.Username != "admin" {
		t.Errorf("expected username admin, got %s", cfg.Username)
	}
	if cfg.Domain != "CORP" {
		t.Errorf("expected domain CORP, got %s", cfg.Domain)
	}
	if cfg.Password != "pass123" {
		t.Errorf("expected password pass123, got %s", cfg.Password)
	}
	if !cfg.IsEnabled {
		t.Error("expected IsEnabled to be true")
	}
}

func TestRunOSAutoLoginRouting(t *testing.T) {
	mock := &mockAutoLoginEngine{}
	prev := GetAutoLoginEngine()
	defer SetAutoLoginEngine(prev)
	SetAutoLoginEngine(mock)

	if err := runOSAutoLogin([]string{"help"}); err != nil {
		t.Errorf("expected nil error on help, got %v", err)
	}
	if err := runOSAutoLogin([]string{"disable"}); err != nil {
		t.Errorf("expected nil error on disable, got %v", err)
	}
	if !mock.isDisabled {
		t.Error("expected mock engine to be marked disabled")
	}

	err := runOSAutoLogin([]string{"enable", "-u", "bob", "-p", "pwd"})
	if err != nil {
		t.Errorf("expected nil error on enable, got %v", err)
	}
	if mock.lastConfig.Username != "bob" {
		t.Errorf("expected bob, got %s", mock.lastConfig.Username)
	}
}
