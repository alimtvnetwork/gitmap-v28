package cmdos

import (
	"testing"
)

func TestKnownDNSProviders(t *testing.T) {
	cf, ok := KnownDNSProviders["cloudflare"]
	if !ok {
		t.Fatalf("expected cloudflare provider to be registered")
	}

	if cf.Primary != "1.1.1.1" {
		t.Errorf("expected 1.1.1.1, got %s", cf.Primary)
	}

	goog, ok := KnownDNSProviders["google"]
	if !ok {
		t.Fatalf("expected google provider to be registered")
	}

	if goog.Primary != "8.8.8.8" {
		t.Errorf("expected 8.8.8.8, got %s", goog.Primary)
	}
}

func TestRunOSDNSRouting(t *testing.T) {
	if err := runOSDNSCommand([]string{"help"}); err != nil {
		t.Errorf("expected nil error for help, got: %v", err)
	}

	err := runOSDNSCommand([]string{"invalid-provider-xyz"})
	if err == nil {
		t.Errorf("expected error for invalid provider, got nil")
	}
}
