package cmdnodes

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsummary"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/pterm/pterm"
)

type nodesPEAllOptions struct {
	isJSON   bool
	forceAll bool
	verbose  bool
}

func parseNodesPEAllOptions(args []string) nodesPEAllOptions {
	var opts nodesPEAllOptions

	for _, a := range args {
		lower := strings.ToLower(a)
		if lower == "--json" || lower == "-j" {
			opts.isJSON = true
			continue
		}

		if lower == "--force-all" || lower == "--all" || lower == "-a" {
			opts.forceAll = true
			continue
		}

		if lower == "--verbose" || lower == "-v" {
			opts.verbose = true
			continue
		}
	}

	return opts
}

func buildLocalPEAllArgs(opts nodesPEAllOptions) []string {
	var localArgs []string
	if opts.isJSON {
		localArgs = append(localArgs, "--json")
	}
	if opts.forceAll {
		localArgs = append(localArgs, "--force-all")
	}
	if opts.verbose {
		localArgs = append(localArgs, "--verbose")
	}
	return localArgs
}

func executeLocalPEAllFallback(opts nodesPEAllOptions, localRecords []model.ScanRecord) error {
	if !opts.isJSON {
		pterm.Warning.Println("No remote fleet SSH nodes found; executing pipeline error discovery locally.")
	}

	localArgs := buildLocalPEAllArgs(opts)
	if !opts.isJSON {
		return cmdsummary.RunPipelineErrorsAll(localArgs)
	}

	captured, _ := captureStdout(func() error {
		return cmdsummary.RunPipelineErrorsAll(localArgs)
	})

	var localPayload interface{}
	cleanJSON := extractJSONPayload(captured)
	if err := json.Unmarshal([]byte(cleanJSON), &localPayload); err != nil {
		localPayload = captured
	}

	envelope := FleetNodesEnvelope{
		Attributes: FleetNodesAttributes{
			GeneratedAt:     time.Now().UTC().Format(time.RFC3339),
			TotalNodes:      1,
			ReachableNodes:  1,
			LocalReposCount: len(localRecords),
			DelegatedCount:  0,
		},
		Data: FleetNodesData{
			LocalNode: NodeSummaryResult{
				NodeAlias:       "local-master",
				Host:            "localhost",
				IsOnline:        true,
				IsLocalMaster:   true,
				UniqueRepoCount: len(localRecords),
				Account:         resolveLocalMasterAccount(),
				Instances:       1,
				Payload:         localPayload,
			},
			RemoteNodes: []NodeSummaryResult{},
		},
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(envelope)
}

// RunNodesPipelineErrorsAll coordinates distributed pe all execution across fleet nodes.
func RunNodesPipelineErrorsAll(args []string) error {
	opts := parseNodesPEAllOptions(args)

	if !opts.isJSON {
		pterm.Info.Println("[GitMap Nodes] Initializing fleet handshake and local repository cataloging...")
	}

	cwd, _ := os.Getwd()
	localRecords := resolveWorkspaceReposHook(cwd)
	localURLMap := make(map[string]bool)
	localSlugMap := make(map[string]bool)

	for _, rec := range localRecords {
		if rec.HTTPSUrl != "" {
			localURLMap[normalizeRepoKey(rec.HTTPSUrl)] = true
		}
		if rec.SSHUrl != "" {
			localURLMap[normalizeRepoKey(rec.SSHUrl)] = true
		}
		if rec.DiscoveredURL != "" {
			localURLMap[normalizeRepoKey(rec.DiscoveredURL)] = true
		}
		if rec.Slug != "" {
			localSlugMap[normalizeRepoKey(rec.Slug)] = true
		}
		if rec.RepoName != "" {
			localSlugMap[normalizeRepoKey(rec.RepoName)] = true
		}
	}

	conns, errConns := fetchAllSSHConnectionsHook()
	if errConns != nil || len(conns) == 0 {
		return executeLocalPEAllFallback(opts, localRecords)
	}

	if !opts.isJSON {
		pterm.Info.Printf("[Phase 1] Cataloged %d repositories on local master machine.\n", len(localRecords))
		pterm.Info.Printf("[Phase 2] Discovering remote repositories across %d fleet nodes...\n", len(conns))
	}

	type remoteDiscovery struct {
		conn       db.SSHConnection
		uniqueURLs []string
		isOnline   bool
		err        error
	}

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

	totalDelegated := 0
	reachableCount := 0
	for _, d := range discoveredRemotes {
		if d.isOnline {
			reachableCount++
		}
		totalDelegated += len(d.uniqueURLs)
	}

	if !opts.isJSON {
		pterm.Success.Printf("[Phase 3] Deduplication applied: %d local repos evaluated locally; %d remote-unique repos delegated.\n",
			len(localRecords), totalDelegated)
	}

	var localArgs []string
	if opts.isJSON {
		localArgs = append(localArgs, "--json")
	}
	if opts.forceAll {
		localArgs = append(localArgs, "--force-all")
	}
	if opts.verbose {
		localArgs = append(localArgs, "--verbose")
	}

	var localPayload interface{}
	if opts.isJSON {
		captured, _ := captureStdout(func() error {
			return cmdsummary.RunPipelineErrorsAll(localArgs)
		})
		localPayload = parsePayloadOrRaw(captured)
	} else {
		fmt.Println()
		pterm.Bold.Println("[Node: Local Machine]")
		_ = cmdsummary.RunPipelineErrorsAll(localArgs)
	}

	var remoteResults []NodeSummaryResult
	for _, d := range discoveredRemotes {
		if !d.isOnline {
			remoteResults = append(remoteResults, NodeSummaryResult{
				NodeAlias:       d.conn.Alias,
				Host:            d.conn.IPAddress,
				IsOnline:        false,
				IsLocalMaster:   false,
				UniqueRepoCount: 0,
				Account:         d.conn.Username,
				Instances:       0,
				ErrorMessage:    d.err.Error(),
			})
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

		remoteCmd := "gitmap pe all"
		if opts.forceAll {
			remoteCmd += " --force-all"
		}
		if opts.isJSON {
			remoteCmd += " --json"
		}

		shell := "bash"
		if isWindowsNode(d.conn) {
			remoteCmd = "gitmap.exe " + remoteCmd
			shell = "powershell"
		}

		out, errExec := currentSSHExecutor.Execute(d.conn, remoteCmd, shell)
		if errExec != nil {
			printNodeExecError(d.conn.Alias, errExec, opts.isJSON)
			remoteResults = append(remoteResults, NodeSummaryResult{
				NodeAlias:       d.conn.Alias,
				Host:            d.conn.IPAddress,
				IsOnline:        true,
				IsLocalMaster:   false,
				UniqueRepoCount: len(d.uniqueURLs),
				Account:         d.conn.Username,
				Instances:       1,
				ErrorMessage:    errExec.Error(),
			})
			continue
		}

		if opts.isJSON {
			remoteResults = append(remoteResults, NodeSummaryResult{
				NodeAlias:       d.conn.Alias,
				Host:            d.conn.IPAddress,
				IsOnline:        true,
				IsLocalMaster:   false,
				UniqueRepoCount: len(d.uniqueURLs),
				Account:         d.conn.Username,
				Instances:       1,
				Payload:         parsePayloadOrRaw(out),
				RawOutput:       out,
			})
		} else {
			fmt.Println(out)
			pterm.Success.Printf("✓ %s completed\n", d.conn.Alias)
		}
	}

	if opts.isJSON {
		envelope := FleetNodesEnvelope{
			Attributes: FleetNodesAttributes{
				GeneratedAt:     time.Now().UTC().Format(time.RFC3339),
				TotalNodes:      len(conns) + 1,
				ReachableNodes:  reachableCount + 1,
				LocalReposCount: len(localRecords),
				DelegatedCount:  totalDelegated,
			},
			Data: FleetNodesData{
				LocalNode: NodeSummaryResult{
					NodeAlias:       "local-master",
					Host:            "localhost",
					IsOnline:        true,
					IsLocalMaster:   true,
					UniqueRepoCount: len(localRecords),
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

	return nil
}
