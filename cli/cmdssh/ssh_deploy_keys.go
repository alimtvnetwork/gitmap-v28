// Package cmdssh — ssh_deploy_keys.go orchestrates mesh public key collection and distribution.
package cmdssh

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

// RunSSHDeployKeysCLI gathers all node public keys, deduplicates them, and deploys to authorized_keys.
func RunSSHDeployKeysCLI(args []string) error {
	except, isDryRun, isJSON := parseDeployKeysFlags(args)
	conns, _ := fetchAllSSHConnections()
	filtered := FilterSSHConnectionsByExcept(conns, except)

	summary := DeployKeysSummary{
		NodesTargeted: len(filtered),
		IsDryRun:      isDryRun,
	}
	allKeys := collectAllMeshPublicKeys(filtered, &summary)
	uniqueKeys := deduplicatePublicKeys(allKeys)
	summary.UniqueKeysIdentified = len(uniqueKeys)

	executeKeysDeployment(filtered, uniqueKeys, &summary)
	return renderDeployKeysSummary(summary, except, isJSON)
}

func collectAllMeshPublicKeys(targets []db.SSHConnection, summary *DeployKeysSummary) []string {
	localKeys := gatherLocalPublicKeys()
	summary.LocalKeysCollected = len(localKeys)
	allKeys := append([]string{}, localKeys...)

	for _, c := range targets {
		keys, err := gatherRemotePublicKeys(c)
		if err != nil {
			continue
		}
		summary.RemoteKeysCollected += len(keys)
		allKeys = append(allKeys, keys...)
	}
	return allKeys
}

func executeKeysDeployment(targets []db.SSHConnection, uniqueKeys []string, summary *DeployKeysSummary) {
	for _, c := range targets {
		res, err := deployKeysToRemoteNode(c, uniqueKeys, summary.IsDryRun)
		if err == nil {
			summary.NodesSucceeded++
		}
		summary.NodeResults = append(summary.NodeResults, res)
	}
}

func parseDeployKeysFlags(args []string) (string, bool, bool) {
	var exceptParts []string
	isDryRun, isJSON := false, false
	for i := 0; i < len(args); i++ {
		low := strings.ToLower(args[i])
		if low == "keys" || low == "k" || low == "deploy" || low == "ssh" || low == "all" {
			continue
		}
		if low == "--dry-run" || low == "-n" {
			isDryRun = true
			continue
		}
		if low == "--json" {
			isJSON = true
			continue
		}
		if isExceptOrExcepFlag(low) {
			for i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				exceptParts = append(exceptParts, args[i+1])
				i++
			}
			continue
		}
		if strings.HasPrefix(low, "--except=") || strings.HasPrefix(low, "--excep=") || strings.HasPrefix(low, "--accept=") || strings.HasPrefix(low, "--exclude=") {
			idx := strings.IndexByte(args[i], '=')
			exceptParts = append(exceptParts, args[i][idx+1:])
		}
	}
	return strings.Join(exceptParts, ","), isDryRun, isJSON
}
