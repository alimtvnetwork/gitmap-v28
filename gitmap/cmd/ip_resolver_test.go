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
