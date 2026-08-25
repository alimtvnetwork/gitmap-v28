package cmd

import (
	"context"
	"testing"
	"time"
)

func TestIPResolver(t *testing.T) {
	resolver := IPResolver{
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
	ip, err := GetLocalIP(ctx, true, "")
	if err != nil {
		t.Logf("GetLocalIP returned error: %v", err)
	} else {
		if ip == "" {
			t.Errorf("expected non-empty IP")
		} else {
			t.Logf("Found IP: %s", ip)
		}
	}
}
