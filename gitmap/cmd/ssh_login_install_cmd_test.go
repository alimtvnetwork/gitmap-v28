package cmd

import (
	"context"
	"testing"
)

func Test_getInstallPayload(t *testing.T) {
	ctx := context.Background()
	payload, err := getInstallPayload(ctx, "target", "v1.0.0")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expected := "curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/refs/heads/main/install.sh?version=v1.0.0 | bash"
	if payload != expected {
		t.Errorf("expected %q, got %q", expected, payload)
	}
}
