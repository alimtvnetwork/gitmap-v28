// Package cmdnodes coordinates deduplicated full summaries and pipeline errors across fleet nodes.
package cmdnodes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpending"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsummary"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/pterm/pterm"
)

var (
	fetchAllSSHConnectionsHook = cmdssh.FetchAllSSHConnections
	resolveWorkspaceReposHook  = cmdpending.ResolveWorkspaceRepositories
)

func parseNodesSummaryOptions(args []string) nodesSummaryOptions {
	opts := nodesSummaryOptions{
		releasesCount: 3,
	}

	for _, a := range args {
		lower := strings.ToLower(a)
		if lower == "--json" || lower == "-j" {
			opts.isJSON = true
			continue
		}

		if lower == "+pe" || lower == "--pe" || lower == "pe" ||
			strings.Contains(lower, "+pe") || strings.Contains(lower, "fspe") {
			opts.withPE = true
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

		if n, err := strconv.Atoi(a); err == nil && n > 0 {
			opts.releasesCount = n
			continue
		}
	}

	return opts
}

func normalizeRepoKey(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.TrimSuffix(s, ".git")
	s = strings.TrimSuffix(s, "/")
	return s
}

func captureStdout(fn func() error) (string, error) {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", fn()
	}
	os.Stdout = w

	outChan := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outChan <- buf.String()
	}()

	var execErr error
	func() {
		defer func() {
			_ = recover()
		}()
		execErr = fn()
	}()

	_ = w.Close()
	os.Stdout = oldStdout
	captured := <-outChan
	_ = r.Close()

	return captured, execErr
}

func buildLocalSummaryArgs(opts nodesSummaryOptions) []string {
	var localArgs []string
	if opts.withPE {
		localArgs = append(localArgs, "+pe")
	}
	if opts.isJSON {
		localArgs = append(localArgs, "--json")
	}
	if opts.releasesCount > 0 {
		localArgs = append(localArgs, strconv.Itoa(opts.releasesCount))
	}
	if opts.forceAll {
		localArgs = append(localArgs, "--force-all")
	}
	if opts.verbose {
		localArgs = append(localArgs, "--verbose")
	}
	return localArgs
}

func executeLocalSummaryFallback(opts nodesSummaryOptions, localRecords []model.ScanRecord) error {
	if !opts.isJSON {
		pterm.Warning.Println("No remote fleet SSH nodes found; executing full summary locally.")
	}

	localArgs := buildLocalSummaryArgs(opts)
	if !opts.isJSON {
		return cmdsummary.RunFullSummary(localArgs)
	}

	captured, _ := captureStdout(func() error {
		return cmdsummary.RunFullSummary(localArgs)
	})

	localPayload := parsePayloadOrRaw(captured)

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
				Host:            constants.Localhost,
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

// RunNodesFullSummary coordinates distributed full summary execution with local precedence deduplication.
func RunNodesFullSummary(args []string) error {
	opts := parseNodesSummaryOptions(args)

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
		return executeLocalSummaryFallback(opts, localRecords)
	}

	if !opts.isJSON {
		pterm.Info.Printf("[Phase 1] Cataloged %d repositories on local master machine.\n", len(localRecords))
		pterm.Info.Printf("[Phase 2] Discovering remote repositories across %d fleet nodes...\n", len(conns))
	}

	discoveredRemotes := discoverRemoteRepositories(conns, localURLMap, localSlugMap)

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

	localArgs := buildLocalSummaryArgs(opts)
	var localPayload interface{}
	if opts.isJSON {
		captured, _ := captureStdout(func() error {
			return cmdsummary.RunFullSummary(localArgs)
		})
		localPayload = parsePayloadOrRaw(captured)
	} else {
		fmt.Println()
		pterm.Bold.Println("[Node: Local Machine]")
		_ = cmdsummary.RunFullSummary(localArgs)
	}

	remoteResults := executeRemoteSummaries(discoveredRemotes, opts)

	if opts.isJSON {
		return renderFleetEnvelopeJSON(conns, reachableCount, len(localRecords), totalDelegated, localPayload, remoteResults)
	}

	return nil
}
