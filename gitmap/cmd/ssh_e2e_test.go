package cmd

import (
	"context"
	"testing"
)

func TestDispatchSSHLogin(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	args := []string{"login", "a@b"}
	
	err := dispatchSSH(ctx, args, nil)
	if err == nil {
		t.Fatalf("expected error due to context cancellation, got nil")
	}
}
