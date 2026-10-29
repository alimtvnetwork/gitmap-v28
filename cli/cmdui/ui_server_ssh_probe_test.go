package cmdui

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func TestProbeStatusTimeout(t *testing.T) {
	got := probeStatus(SSHPingResult{Error: "dial tcp 192.168.1.7:22: i/o timeout"})
	if got != "Timed out" {
		t.Fatalf("status = %s", got)
	}
}

func TestProbeStatusOnline(t *testing.T) {
	got := probeStatus(SSHPingResult{IsOnline: true})
	if got != probeStatusOnline {
		t.Fatalf("status = %s", got)
	}
}

func TestSummaryKeepsPing(t *testing.T) {
	t.Cleanup(clearProbeAlias("w1"))
	rememberNodeProbe("w1", SSHPingResult{IsOnline: true, Success: true})
	summary := firstSummary(t)

	if summary.Status != probeStatusOnline || !summary.IsOnline {
		t.Fatalf("status=%s online=%v", summary.Status, summary.IsOnline)
	}
}

func TestSummaryMatchesProbeFold(t *testing.T) {
	t.Cleanup(clearProbeAlias("W1"))
	rememberNodeProbe("W1", SSHPingResult{IsOnline: true, Success: true})
	summary := firstSummary(t)

	if summary.Status != probeStatusOnline {
		t.Fatalf("status=%s", summary.Status)
	}
}

func firstSummary(t *testing.T) NodeSummary {
	t.Helper()
	summaries := buildNodeSummaries([]db.SSHConnection{{
		Alias:     "w1",
		IPAddress: "192.168.1.3",
		Username:  "Administrator",
		OS:        "win",
	}})
	if len(summaries) != 1 {
		t.Fatalf("summaries=%d", len(summaries))
	}

	return summaries[0]
}

func clearProbeAlias(alias string) func() {
	return func() {
		nodeProbeMu.Lock()
		delete(nodeProbes, alias)
		nodeProbeMu.Unlock()
	}
}
