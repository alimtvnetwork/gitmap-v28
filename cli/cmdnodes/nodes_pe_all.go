package cmdnodes

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsummary"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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

	localArgs := buildLocalPEAllArgs(opts)
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

	remoteResults := executePEAllRemoteSummaries(discoveredRemotes, opts)

	if opts.isJSON {
		return renderPEAllFleetEnvelopeJSON(conns, reachableCount, len(localRecords), totalDelegated, localPayload, remoteResults)
	}

	return nil
}
