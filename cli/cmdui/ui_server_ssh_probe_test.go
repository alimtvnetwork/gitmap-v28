package cmdui

import "testing"

func TestProbeStatusTimeout(t *testing.T) {
	got := probeStatus(SSHPingResult{Error: "dial tcp 192.168.1.7:22: i/o timeout"})
	if got != "Timed out" {
		t.Fatalf("status = %s", got)
	}
}

func TestProbeStatusOnline(t *testing.T) {
	got := probeStatus(SSHPingResult{IsOnline: true})
	if got != "Online" {
		t.Fatalf("status = %s", got)
	}
}
