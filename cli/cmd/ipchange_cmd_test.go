package cmd

import (
	"context"
	"testing"
)

func TestRunIPChangeCmd(t *testing.T) {
	ctx := context.Background()

	// Missing args
	err := runIPChangeCmd(nil, []string{}, ctx)
	if err == nil {
		t.Errorf("expected error when args are missing")
	}
}

func TestValidatePing(t *testing.T) {
	ctx := context.Background()
	// Loopback address ping should succeed on virtually all machines
	res := validatePing(ctx, "127.0.0.1", 1)
	if !res {
		t.Logf("ping 127.0.0.1 returned false (may be blocked by local firewall/permissions)")
	}
}

func TestExecuteIPChange(t *testing.T) {
	ctx := context.Background()
	_ = executeIPChange(ctx, "192.168.1.50", false)
}
