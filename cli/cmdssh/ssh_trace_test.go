package cmdssh

import (
	"errors"
	"strings"
	"testing"
)

func TestSSHTrace_Recording(t *testing.T) {
	trace := BeginSSHTrace("test_enroll", "user@10.0.0.1:22", "user", "10.0.0.1", 22)
	if trace == nil {
		t.Fatalf("expected non-nil trace")
	}

	trace.AddStep("TCP Probe", "10.0.0.1:22 reachable", "SUCCESS", nil)
	trace.AddStep("Password Auth", "invalid credentials", "FAILED", errors.New("handshake failed"))
	trace.SetInternalError(errors.New("raw ssh handshake rejected"), "Check password")

	if len(trace.Steps) != 2 {
		t.Errorf("expected 2 steps, got %d", len(trace.Steps))
	}
	if trace.InternalError != "raw ssh handshake rejected" {
		t.Errorf("unexpected internal error: %s", trace.InternalError)
	}

	report := trace.FormatTerminalReport()
	if !strings.Contains(report, "TCP Probe") {
		t.Errorf("expected report to contain 'TCP Probe', got %s", report)
	}
	if !strings.Contains(report, "raw ssh handshake rejected") {
		t.Errorf("expected report to contain internal error, got %s", report)
	}
}

func TestSSHErrorLogsCLI_Empty(t *testing.T) {
	err := RunSSHErrorLogsCLI([]string{"--clear"})
	if err != nil {
		t.Fatalf("unexpected error on clear: %v", err)
	}

	err = RunSSHErrorLogsCLI([]string{})
	if err != nil {
		t.Fatalf("unexpected error on empty logs: %v", err)
	}
}
