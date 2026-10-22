package cmdnodes

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func probeUnifiedNodesLive(ctx context.Context, nodes []UnifiedFleetNode, hosts []store.SSHHost, isCacheAllowed ...bool) []UnifiedFleetNode {
	isNodesEmpty := len(nodes) == 0
	if isNodesEmpty {
		return nodes
	}

	useCache := true
	hasParam := len(isCacheAllowed) > 0
	if hasParam {
		useCache = isCacheAllowed[0]
	}

	isSkipEnvSet := os.Getenv("GITMAP_SKIP_NODE_CACHE") != ""
	isCachingEnabled := useCache && !isSkipEnvSet

	toCache := make(map[string]NodeLivenessCacheEntry, len(nodes))
	uncachedIndices := partitionCachedNodes(nodes, isCachingEnabled, toCache)

	isAllCached := len(uncachedIndices) == 0
	if isAllCached {
		return nodes
	}

	probeAndRecordUncached(ctx, nodes, hosts, uncachedIndices, toCache)
	SetCachedNodesLiveness(toCache)
	return nodes
}

func partitionCachedNodes(nodes []UnifiedFleetNode, isCachingEnabled bool, toCache map[string]NodeLivenessCacheEntry) []int {
	var uncached []int
	for i := range nodes {
		isNodeCached := checkAndApplyCachedLiveness(&nodes[i], isCachingEnabled, toCache)
		if isNodeCached {
			continue
		}
		uncached = append(uncached, i)
	}
	return uncached
}

func checkAndApplyCachedLiveness(node *UnifiedFleetNode, isCachingEnabled bool, toCache map[string]NodeLivenessCacheEntry) bool {
	isBypass := !isCachingEnabled
	if isBypass {
		return false
	}

	entry, isCached := GetCachedNodeLiveness(node.Host, node.Alias)
	isEntryMissing := !isCached
	if isEntryMissing {
		return false
	}

	node.Status = entry.Status
	node.IsOnline = entry.IsOnline
	cacheKey := resolveLivenessCacheKey(node.Host, node.Alias)
	toCache[cacheKey] = entry
	return true
}

func probeAndRecordUncached(ctx context.Context, nodes []UnifiedFleetNode, hosts []store.SSHHost, uncachedIndices []int, toCache map[string]NodeLivenessCacheEntry) {
	tempHosts := make([]store.SSHHost, len(uncachedIndices))
	for k, idx := range uncachedIndices {
		n := nodes[idx]
		tempHosts[k] = store.SSHHost{
			Alias:             n.Alias,
			IP:                n.Host,
			Port:              n.Port,
			Username:          n.User,
			EncryptedPassword: lookupHostPassword(n.Host, hosts),
		}
	}

	probed := cmdssh.ProbeHostsLiveStatus(ctx, tempHosts)
	nowUnix := time.Now().Unix()

	for k, idx := range uncachedIndices {
		st := probed[k].Status
		isOnline := strings.Contains(st, "ready") || strings.Contains(st, "●")
		nodes[idx].Status = st
		nodes[idx].IsOnline = isOnline

		cacheEntry := NodeLivenessCacheEntry{
			Alias:        nodes[idx].Alias,
			Host:         nodes[idx].Host,
			Port:         nodes[idx].Port,
			Status:       st,
			IsOnline:     isOnline,
			ProbedAtUnix: nowUnix,
		}
		cacheKey := resolveLivenessCacheKey(nodes[idx].Host, nodes[idx].Alias)
		toCache[cacheKey] = cacheEntry
		aliasKey := strings.ToLower(strings.TrimSpace(nodes[idx].Alias))
		hasDistinctAlias := aliasKey != "" && aliasKey != "-" && aliasKey != cacheKey
		if hasDistinctAlias {
			toCache[aliasKey] = cacheEntry
		}
	}
}

func lookupHostPassword(ip string, hosts []store.SSHHost) string {
	for _, h := range hosts {
		isIPMatch := h.IP == ip
		hasPass := h.EncryptedPassword != ""
		if isIPMatch && hasPass {
			return h.EncryptedPassword
		}
	}
	return ""
}

func applyNodesFilter(nodes []UnifiedFleetNode, opts nodesFilterOptions) []UnifiedFleetNode {
	var filtered []UnifiedFleetNode
	for _, n := range nodes {
		isNodeOffline := !n.IsOnline
		isOfflineOmitted := opts.openOnly && isNodeOffline
		if isOfflineOmitted {
			continue
		}
		hasTargetMatch := matchesTargetFilter(n, opts.targetFilter)
		isTargetMismatch := !hasTargetMatch
		if isTargetMismatch {
			continue
		}
		hasSubsystemMatch := matchesSubsystemFilter(n, opts)
		isSubsystemMismatch := !hasSubsystemMatch
		if isSubsystemMismatch {
			continue
		}
		filtered = append(filtered, n)
	}
	return filtered
}

func matchesTargetFilter(n UnifiedFleetNode, target string) bool {
	isTargetEmpty := target == ""
	if isTargetEmpty {
		return true
	}
	low := strings.ToLower(target)
	isAliasMatch := strings.Contains(strings.ToLower(n.Alias), low)
	isHostMatch := strings.Contains(strings.ToLower(n.Host), low)
	return isAliasMatch || isHostMatch
}

func matchesSubsystemFilter(n UnifiedFleetNode, opts nodesFilterOptions) bool {
	hasSpecificFilter := opts.filterSSH || opts.filterClust || opts.filterSC
	isUnfiltered := !hasSpecificFilter
	if isUnfiltered {
		return true
	}
	isSSHMatch := opts.filterSSH && hasSubsystem(n, "SSH")
	if isSSHMatch {
		return true
	}
	isClusterMatch := opts.filterClust && hasSubsystem(n, "Cluster")
	if isClusterMatch {
		return true
	}
	isSCMatch := opts.filterSC && hasSubsystem(n, "SC")
	if isSCMatch {
		return true
	}
	return false
}

func hasSubsystem(n UnifiedFleetNode, sys string) bool {
	for _, s := range n.Subsystems {
		isEqual := strings.EqualFold(s, sys)
		if isEqual {
			return true
		}
	}
	return false
}

func printNodesJSON(nodes []UnifiedFleetNode) error {
	data, err := json.MarshalIndent(nodes, "", "  ")
	isMarshalErr := err != nil
	if isMarshalErr {
		return apperror.WrapSimple(err, "printNodesJSON")
	}
	fmt.Println(string(data))
	return nil
}
