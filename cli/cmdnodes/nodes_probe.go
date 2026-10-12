package cmdnodes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func probeUnifiedNodesLive(ctx context.Context, nodes []UnifiedFleetNode, hosts []store.SSHHost) []UnifiedFleetNode {
	if len(nodes) == 0 {
		return nodes
	}

	tempHosts := make([]store.SSHHost, len(nodes))
	for i, n := range nodes {
		tempHosts[i] = store.SSHHost{
			Alias:             n.Alias,
			IP:                n.Host,
			Port:              n.Port,
			Username:          n.User,
			EncryptedPassword: lookupHostPassword(n.Host, hosts),
		}
	}

	probed := cmdssh.ProbeHostsLiveStatus(ctx, tempHosts)

	for i := range nodes {
		nodes[i].Status = probed[i].Status
		nodes[i].IsOnline = strings.Contains(probed[i].Status, "ready") || strings.Contains(probed[i].Status, "●")
	}

	return nodes
}

func lookupHostPassword(ip string, hosts []store.SSHHost) string {
	for _, h := range hosts {
		if h.IP == ip && h.EncryptedPassword != "" {
			return h.EncryptedPassword
		}
	}
	return ""
}

func applyNodesFilter(nodes []UnifiedFleetNode, opts nodesFilterOptions) []UnifiedFleetNode {
	var filtered []UnifiedFleetNode
	for _, n := range nodes {
		if !matchesTargetFilter(n, opts.targetFilter) {
			continue
		}
		if !matchesSubsystemFilter(n, opts) {
			continue
		}
		filtered = append(filtered, n)
	}
	return filtered
}

func matchesTargetFilter(n UnifiedFleetNode, target string) bool {
	if target == "" {
		return true
	}
	low := strings.ToLower(target)
	return strings.Contains(strings.ToLower(n.Alias), low) || strings.Contains(strings.ToLower(n.Host), low)
}

func matchesSubsystemFilter(n UnifiedFleetNode, opts nodesFilterOptions) bool {
	hasSpecificFilter := opts.filterSSH || opts.filterClust || opts.filterSC
	if !hasSpecificFilter {
		return true
	}
	if opts.filterSSH && hasSubsystem(n, "SSH") {
		return true
	}
	if opts.filterClust && hasSubsystem(n, "Cluster") {
		return true
	}
	if opts.filterSC && hasSubsystem(n, "SC") {
		return true
	}
	return false
}

func hasSubsystem(n UnifiedFleetNode, sys string) bool {
	for _, s := range n.Subsystems {
		if strings.EqualFold(s, sys) {
			return true
		}
	}
	return false
}

func printNodesJSON(nodes []UnifiedFleetNode) error {
	data, err := json.MarshalIndent(nodes, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "printNodesJSON")
	}
	fmt.Println(string(data))
	return nil
}
