package cmdnodes

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/pterm/pterm"
)

type remoteDiscovery struct {
	conn       db.SSHConnection
	uniqueURLs []string
	isOnline   bool
	err        error
}

func discoverRemoteRepositories(conns []db.SSHConnection, localURLMap, localSlugMap map[string]bool) []remoteDiscovery {
	var discoveredRemotes []remoteDiscovery
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, c := range conns {
		wg.Add(1)
		go func(conn db.SSHConnection) {
			defer wg.Done()

			cmd := "gitmap list --json"
			shell := "bash"
			if isWindowsNode(conn) {
				cmd = "gitmap.exe list --json"
				shell = "powershell"
			}

			out, errExec := currentSSHExecutor.Execute(conn, cmd, shell)
			if errExec != nil {
				mu.Lock()
				discoveredRemotes = append(discoveredRemotes, remoteDiscovery{
					conn:     conn,
					isOnline: false,
					err:      errExec,
				})
				mu.Unlock()
				return
			}

			manifest := parseRemoteManifest(out)
			var uniqueURLs []string
			for _, item := range manifest {
				isAlreadyLocal := (item.RemoteURL != "" && localURLMap[normalizeRepoKey(item.RemoteURL)]) ||
					(item.Slug != "" && localSlugMap[normalizeRepoKey(item.Slug)]) ||
					(item.RepoName != "" && localSlugMap[normalizeRepoKey(item.RepoName)])
				if isAlreadyLocal {
					continue
				}
				uniqueURLs = append(uniqueURLs, resolveManifestItemIdentifier(item))
			}

			mu.Lock()
			discoveredRemotes = append(discoveredRemotes, remoteDiscovery{
				conn:       conn,
				uniqueURLs: uniqueURLs,
				isOnline:   true,
			})
			mu.Unlock()
		}(c)
	}

	wg.Wait()
	return discoveredRemotes
}

func executeRemoteSummaries(discoveredRemotes []remoteDiscovery, opts nodesSummaryOptions) []NodeSummaryResult {
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

		remoteCmd, shell := buildRemoteSummaryCommand(d.conn, opts)
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

func buildOfflineNodeResult(conn db.SSHConnection, err error) NodeSummaryResult {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return NodeSummaryResult{
		NodeAlias:       conn.Alias,
		Host:            conn.IPAddress,
		IsOnline:        false,
		IsLocalMaster:   false,
		UniqueRepoCount: 0,
		Account:         conn.Username,
		Instances:       0,
		ErrorMessage:    errMsg,
	}
}

func buildErrorNodeResult(conn db.SSHConnection, uniqueCount int, err error) NodeSummaryResult {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return NodeSummaryResult{
		NodeAlias:       conn.Alias,
		Host:            conn.IPAddress,
		IsOnline:        true,
		IsLocalMaster:   false,
		UniqueRepoCount: uniqueCount,
		Account:         conn.Username,
		Instances:       1,
		ErrorMessage:    errMsg,
	}
}

func appendSuccessfulRemoteResult(results []NodeSummaryResult, conn db.SSHConnection, uniqueCount int, out string, isJSON bool) []NodeSummaryResult {
	if isJSON {
		return append(results, NodeSummaryResult{
			NodeAlias:       conn.Alias,
			Host:            conn.IPAddress,
			IsOnline:        true,
			IsLocalMaster:   false,
			UniqueRepoCount: uniqueCount,
			Account:         conn.Username,
			Instances:       1,
			Payload:         parsePayloadOrRaw(out),
			RawOutput:       out,
		})
	}
	fmt.Println(out)
	pterm.Success.Printf("✓ %s completed\n", conn.Alias)
	return results
}

func buildRemoteSummaryCommand(conn db.SSHConnection, opts nodesSummaryOptions) (string, string) {
	remoteCmd := fmt.Sprintf("gitmap fs %d", opts.releasesCount)
	if opts.withPE {
		remoteCmd = fmt.Sprintf("gitmap fs+pe %d", opts.releasesCount)
	}
	if opts.isJSON {
		remoteCmd += " --json"
	}
	if opts.forceAll {
		remoteCmd += " --force-all"
	}

	shell := "bash"
	if isWindowsNode(conn) {
		remoteCmd = "gitmap.exe " + remoteCmd
		shell = "powershell"
	}

	return remoteCmd, shell
}

func renderFleetEnvelopeJSON(conns []db.SSHConnection, reachableCount, localReposCount, delegatedCount int, localPayload interface{}, remoteResults []NodeSummaryResult) error {
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
