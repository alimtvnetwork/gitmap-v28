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

const (
	probeStatusOnline   = "Online"
	probeStatusOffline  = "Offline"
	probeStatusTimedOut = "Timed out"
	probeTimeoutMarker  = "timeout"
)

func finishPingResult(res SSHPingResult) SSHPingResult {
	res.Status = probeStatus(res)

	return res
}

func rememberNodeProbe(alias string, res SSHPingResult) {
	nodeProbeMu.Lock()
	nodeProbes[alias] = nodeProbe{IsOnline: res.IsOnline, Status: probeStatus(res)}
	nodeProbeMu.Unlock()
}

func probeStatus(res SSHPingResult) string {
	if res.IsOnline {
		return probeStatusOnline
	}
	if strings.Contains(strings.ToLower(res.Error), probeTimeoutMarker) {
		return probeStatusTimedOut
	}

	return probeStatusOffline
}

func storedNodeProbe(alias string) (nodeProbe, bool) {
	nodeProbeMu.Lock()
	probe, isFound := lookupNodeProbe(alias)
	nodeProbeMu.Unlock()

	return probe, isFound
}

func lookupNodeProbe(alias string) (nodeProbe, bool) {
	probe, isFound := nodeProbes[alias]
	if isFound {
		return probe, true
	}

	return probeByFold(alias)
}

func probeByFold(alias string) (nodeProbe, bool) {
	for key, probe := range nodeProbes {
		if strings.EqualFold(key, alias) {
			return probe, true
		}
	}

	return nodeProbe{}, false
}

func applyStoredProbe(summary *NodeSummary) {
	probe, isFound := storedNodeProbe(summary.Alias)
	if !isFound {
		return
	}
	summary.IsOnline = probe.IsOnline
	summary.Status = probe.Status
}
