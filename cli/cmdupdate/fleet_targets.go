package cmdupdate

import (
	"context"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func filterFleetTargets(targets []FleetTarget, targetFilter string, exclusions map[string]bool) ([]FleetTarget, int) {
	var filtered []FleetTarget
	excludedCount := 0
	targetLower := strings.ToLower(targetFilter)
	for _, t := range targets {
		if isFleetTargetMatch(t, targetLower) == false {
			continue
		}
		if isFleetNodeExcluded(t, exclusions) {
			excludedCount++
			continue
		}
		filtered = append(filtered, t)
	}
	return filtered, excludedCount
}

func isFleetTargetMatch(t FleetTarget, target string) bool {
	if target == "" || target == "all" || target == "all-nodes" {
		return true
	}
	if strings.EqualFold(t.Alias, target) {
		return true
	}
	if strings.EqualFold(t.IP, target) {
		return true
	}
	return strings.EqualFold(t.ID, target)
}

func isFleetNodeExcluded(t FleetTarget, exclusions map[string]bool) bool {
	if len(exclusions) == 0 {
		return false
	}
	if exclusions[strings.ToLower(t.ID)] {
		return true
	}
	if exclusions[strings.ToLower(t.Alias)] {
		return true
	}
	if exclusions[strings.ToLower(t.IP)] {
		return true
	}
	userHost := fmt.Sprintf("%s@%s", t.Username, t.IP)
	return exclusions[strings.ToLower(userHost)]
}

func loadDefaultFleetTargets() ([]FleetTarget, error) {
	conns, err := cmdssh.FetchAllSSHConnections()
	if err != nil {
		return nil, apperror.WrapSimple(err, "cmdssh.FetchAllSSHConnections")
	}
	return convertConnectionsToFleetTargets(conns), nil
}

func resolveFleetTargets(includeOthers bool) ([]FleetTarget, error) {
	targets, err := LoadFleetTargetsFn()
	if err != nil {
		return nil, apperror.WrapSimple(err, "LoadFleetTargetsFn")
	}
	if !includeOthers {
		return targets, nil
	}
	clusterTargets, err := LoadClusterTargetsFn()
	if err != nil {
		return targets, nil
	}
	return mergeFleetTargets(targets, clusterTargets), nil
}

func loadDefaultClusterTargets() ([]FleetTarget, error) {
	ctx := context.Background()
	storeDB, err := store.OpenDefault()
	if err != nil {
		return nil, apperror.WrapSimple(err, "store.OpenDefault")
	}
	defer storeDB.Close()

	var targets []FleetTarget
	if hosts, hostErr := store.ListHosts(ctx, storeDB.Conn()); hostErr == nil {
		targets = append(targets, convertHostsToFleetTargets(hosts)...)
	}
	nodesRes := db.ListClusterNodes(ctx, storeDB.Conn())
	if nodesRes.IsSuccess() {
		targets = append(targets, convertClusterNodesToFleetTargets(nodesRes.Data)...)
	}
	return targets, nil
}

func convertHostsToFleetTargets(hosts []store.SSHHost) []FleetTarget {
	targets := make([]FleetTarget, 0, len(hosts))
	for _, h := range hosts {
		targets = append(targets, FleetTarget{
			ID:       h.Alias,
			Alias:    h.Alias,
			IP:       h.IP,
			Username: h.Username,
			Port:     h.Port,
			Password: h.EncryptedPassword,
			OS:       "windows",
		})
	}
	return targets
}

func convertClusterNodesToFleetTargets(nodes []db.ClusterNode) []FleetTarget {
	targets := make([]FleetTarget, 0, len(nodes))
	for _, n := range nodes {
		targets = append(targets, FleetTarget{
			ID:       n.NodeId,
			Alias:    n.Alias,
			IP:       n.IPAddress,
			Username: "root",
			Port:     22,
			OS:       resolveOS(n.OS),
		})
	}
	return targets
}

func mergeFleetTargets(primary []FleetTarget, additional []FleetTarget) []FleetTarget {
	seen := make(map[string]bool)
	merged := make([]FleetTarget, 0, len(primary)+len(additional))
	for _, t := range primary {
		key := strings.ToLower(strings.TrimSpace(t.IP))
		if key != "" && !seen[key] {
			seen[key] = true
			merged = append(merged, t)
		}
	}
	for _, t := range additional {
		key := strings.ToLower(strings.TrimSpace(t.IP))
		if key != "" && !seen[key] {
			seen[key] = true
			merged = append(merged, t)
		}
	}
	return merged
}

func convertConnectionsToFleetTargets(conns []db.SSHConnection) []FleetTarget {
	seen := make(map[string]bool)
	var targets []FleetTarget

	for _, c := range conns {
		key := strings.ToLower(c.IPAddress)
		if seen[key] {
			continue
		}
		seen[key] = true
		targets = append(targets, FleetTarget{
			ID:       c.Alias,
			Alias:    c.Alias,
			IP:       c.IPAddress,
			Username: c.Username,
			Port:     22,
			Password: c.EncryptedPassword,
			KeyPath:  c.KeyPath,
			OS:       resolveOS(c.OS),
		})
	}
	return targets
}

func resolveOS(osName string) string {
	if osName != "" {
		return osName
	}
	return "windows"
}
