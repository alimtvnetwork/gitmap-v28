package cmd

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdip"
)

func TestIPResolver(t *testing.T) {
	resolver := cmdip.IPResolver{
		Cache:   make(map[string]string),
		Timeout: 5 * time.Second,
	}

	if resolver.Timeout != 5*time.Second {
		t.Errorf("expected Timeout to be 5s")
	}

	_, err := resolver.FetchLocalIP(context.Background())
	if err == nil {
		t.Errorf("expected an error from scaffolded FetchLocalIP")
	}
}

func TestGetLocalIP(t *testing.T) {
	ctx := context.Background()
	ip, err := cmdip.GetLocalIP(ctx, true, "")
	if err != nil {
		t.Skipf("cmdip.GetLocalIP returned error (maybe no network): %v", err)
	}

	if ip == "" {
		t.Errorf("expected non-empty IP")
	}

	if ip == "127.0.0.1" {
		t.Errorf("expected non-loopback IP, got 127.0.0.1")
	}

	matched, _ := regexp.MatchString(`^\d{1,3}(\.\d{1,3}){3}$`, ip)
	if !matched {
		t.Errorf("expected IPv4 format, got %s", ip)
	}
}
