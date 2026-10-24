package cmdnodes

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/pterm/pterm"
)

func executePEAllRemoteSummaries(discoveredRemotes []remoteDiscovery, opts nodesPEAllOptions) []NodeSummaryResult {
	var remoteResults []NodeSummaryResult
	for _, d := range discoveredRemotes {
		if !d.isOnline {
			remoteResults = append(remoteResults, buildOfflineNodeResult(d.conn, d.err))
			continue
		}

		if len(d.uniqueURLs) == 0 {
			remoteResults = appendEmptyOnlineNodeResult(remoteResults, d.conn, opts.isJSON)
			continue
		}

		if !opts.isJSON {
			fmt.Println()
			pterm.Bold.Printf("[Node: %s (%s) - Delegated %d Unique Repositories]\n",
				d.conn.Alias, d.conn.IPAddress, len(d.uniqueURLs))
		}

		remoteCmd, shell := buildRemotePEAllCommand(d.conn, opts)
		out, errExec := currentSSHExecutor.Execute(d.conn, remoteCmd, shell)
		if errExec != nil {
			printNodeExecError(d.conn.Alias, errExec, opts.isJSON)
			remoteResults = append(remoteResults, buildErrorNodeResult(d.conn, len(d.uniqueURLs), errExec))
			continue
		}

		remoteResults = appendSuccessfulRemoteResult(remoteResults, d.conn, len(d.uniqueURLs), out, opts.isJSON)
	}

	return remoteResults
}

func buildRemotePEAllCommand(conn db.SSHConnection, opts nodesPEAllOptions) (string, string) {
	remoteCmd := "gitmap pe all"
	if opts.forceAll {
		remoteCmd += " --force-all"
	}
	if opts.isJSON {
		remoteCmd += " --json"
	}

	shell := "bash"
	if isWindowsNode(conn) {
		remoteCmd = "gitmap.exe " + remoteCmd
		shell = "powershell"
	}

	return remoteCmd, shell
}

func renderPEAllFleetEnvelopeJSON(conns []db.SSHConnection, reachableCount, localReposCount, delegatedCount int, localPayload interface{}, remoteResults []NodeSummaryResult) error {
	envelope := FleetNodesEnvelope{
		Attributes: FleetNodesAttributes{
			GeneratedAt:     time.Now().UTC().Format(time.RFC3339),
			TotalNodes:      len(conns) + 1,
			ReachableNodes:  reachableCount + 1,
			LocalReposCount: localReposCount,
			DelegatedCount:  delegatedCount,
		},
		Data: FleetNodesData{
			LocalNode: NodeSummaryResult{
				NodeAlias:       "local-master",
				Host:            constants.Localhost,
				IsOnline:        true,
				IsLocalMaster:   true,
				UniqueRepoCount: localReposCount,
				Account:         resolveLocalMasterAccount(),
				Instances:       1,
				Payload:         localPayload,
			},
			RemoteNodes: remoteResults,
		},
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(envelope)
}
