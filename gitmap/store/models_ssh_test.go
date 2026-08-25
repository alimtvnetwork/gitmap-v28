package store

import (
	"testing"
	"time"
)

func TestSSHHost(t *testing.T) {
	now := time.Now()
	host := SSHHost{
		ID:        "host-123",
		Alias:     "my-server",
		IP:        "192.168.1.10",
		Username:  "admin",
		CreatedAt: now,
	}

	if host.ID != "host-123" {
		t.Errorf("expected ID host-123, got %s", host.ID)
	}
	if host.Alias != "my-server" {
		t.Errorf("expected Alias my-server, got %s", host.Alias)
	}
	if host.IP != "192.168.1.10" {
		t.Errorf("expected IP 192.168.1.10, got %s", host.IP)
	}
	if host.Username != "admin" {
		t.Errorf("expected Username admin, got %s", host.Username)
	}
	if host.CreatedAt != now {
		t.Errorf("expected CreatedAt %v, got %v", now, host.CreatedAt)
	}
}
