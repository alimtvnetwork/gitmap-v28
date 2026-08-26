package cmd

import (
	"testing"
)

func TestSSHMultiParser(t *testing.T) {
	raw := "192.168.1.1, 192.168.1.2 192.168.1.3"
	ips := ParseMultiIPList(raw)
	if len(ips) != 3 {
		t.Fatalf("expected 3 IPs, got %d", len(ips))
	}
	if ips[0] != "192.168.1.1" || ips[1] != "192.168.1.2" || ips[2] != "192.168.1.3" {
		t.Fatalf("unexpected IPs parsed: %v", ips)
	}
}
