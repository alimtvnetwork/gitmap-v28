package store

import (
	"testing"
	"time"
)

func TestSSHHistory(t *testing.T) {
	now := time.Now()
	h := SSHHistory{
		ID:       "test-id",
		HostIP:   "192.168.1.1",
		JoinedAt: now,
		User:     "admin",
	}

	if h.ID != "test-id" {
		t.Errorf("expected ID test-id, got %s", h.ID)
	}
	if h.HostIP != "192.168.1.1" {
		t.Errorf("expected HostIP 192.168.1.1, got %s", h.HostIP)
	}
	if !h.JoinedAt.Equal(now) {
		t.Errorf("expected JoinedAt to match %v, got %v", now, h.JoinedAt)
	}
	if h.User != "admin" {
		t.Errorf("expected User admin, got %s", h.User)
	}
}
