package cmdui

import (
	"strings"
	"sync"
)

type nodeProbe struct {
	IsOnline bool
	Status   string
}

var (
	nodeProbeMu sync.Mutex
	nodeProbes  = map[string]nodeProbe{}
)

func rememberNodeProbe(alias string, res SSHPingResult) {
	nodeProbeMu.Lock()
	nodeProbes[alias] = nodeProbe{IsOnline: res.IsOnline, Status: probeStatus(res)}
	nodeProbeMu.Unlock()
}

func probeStatus(res SSHPingResult) string {
	if res.IsOnline {
		return "Online"
	}
	if strings.Contains(strings.ToLower(res.Error), "timeout") {
		return "Timed out"
	}

	return "Offline"
}

func storedNodeProbe(alias string) (nodeProbe, bool) {
	nodeProbeMu.Lock()
	probe, isFound := nodeProbes[alias]
	nodeProbeMu.Unlock()

	return probe, isFound
}

func applyStoredProbe(summary *NodeSummary) {
	probe, isFound := storedNodeProbe(summary.Alias)
	if !isFound {
		return
	}
	summary.IsOnline = probe.IsOnline
	summary.Status = probe.Status
}
