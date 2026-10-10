// Package cmdnodes coordinates deduplicated full summaries and pipeline errors across fleet nodes.
package cmdnodes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpending"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsummary"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/pterm/pterm"
)

var (
	fetchAllSSHConnectionsHook = cmdssh.FetchAllSSHConnections
	resolveWorkspaceReposHook  = cmdpending.ResolveWorkspaceRepositories
)

// RemoteRepoManifestItem models a remote repository discovered during the discovery handshake.
type RemoteRepoManifestItem struct {
	RepoName  string `json:"repoName"`
	Slug      string `json:"slug"`
	RemoteURL string `json:"remoteUrl"`
	Path      string `json:"path"`
}

// FleetNodesEnvelope models root JSON envelope for distributed commands.
type FleetNodesEnvelope struct {
	Attributes FleetNodesAttributes `json:"attributes"`
	Data       FleetNodesData       `json:"data"`
}

// FleetNodesAttributes models metadata attributes of the fleet nodes command run.
type FleetNodesAttributes struct {
	GeneratedAt     string `json:"generatedAt"`
	TotalNodes      int    `json:"totalNodes"`
	ReachableNodes  int    `json:"reachableNodes"`
	LocalReposCount int    `json:"localReposCount"`
	DelegatedCount  int    `json:"delegatedCount"`
}

// FleetNodesData models the data payload containing local and remote results.
type FleetNodesData struct {
	LocalNode   NodeSummaryResult   `json:"localNode"`
	RemoteNodes []NodeSummaryResult `json:"remoteNodes"`
}

// NodeSummaryResult aggregates execution result for one fleet node.
type NodeSummaryResult struct {
	NodeAlias       string      `json:"nodeAlias"`
	Host            string      `json:"host"`
	IsOnline        bool        `json:"isOnline"`
	IsLocalMaster   bool        `json:"isLocalMaster"`
	UniqueRepoCount int         `json:"uniqueRepoCount"`
	Payload         interface{} `json:"payload,omitempty"`
	RawOutput       string      `json:"rawOutput,omitempty"`
	ErrorMessage    string      `json:"errorMessage,omitempty"`
}

type nodesSummaryOptions struct {
	isJSON        bool
	withPE        bool
	forceAll      bool
	verbose       bool
	releasesCount int
}

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
		if !opts.isJSON {
			pterm.Warning.Println("No remote fleet SSH nodes found; executing full summary locally.")
		}

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

		if !opts.isJSON {
			return cmdsummary.RunFullSummary(localArgs)
		}

		captured, _ := captureStdout(func() error {
			return cmdsummary.RunFullSummary(localArgs)
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
					Payload:         localPayload,
				},
				RemoteNodes: []NodeSummaryResult{},
			},
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(envelope)
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
				if !isAlreadyLocal {
					ident := item.RemoteURL
					if ident == "" {
						ident = item.Slug
					}
					if ident == "" {
						ident = item.RepoName
					}
					uniqueURLs = append(uniqueURLs, ident)
				}
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

	var localPayload interface{}
	if opts.isJSON {
		captured, _ := captureStdout(func() error {
			return cmdsummary.RunFullSummary(localArgs)
		})
		cleanJSON := extractJSONPayload(captured)
		if err := json.Unmarshal([]byte(cleanJSON), &localPayload); err != nil {
			localPayload = captured
		}
	} else {
		fmt.Println()
		pterm.Bold.Println("[Node: Local Machine]")
		_ = cmdsummary.RunFullSummary(localArgs)
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
				ErrorMessage:    d.err.Error(),
			})
			continue
		}

		if len(d.uniqueURLs) == 0 {
			if opts.isJSON {
				remoteResults = append(remoteResults, NodeSummaryResult{
					NodeAlias:       d.conn.Alias,
					Host:            d.conn.IPAddress,
					IsOnline:        true,
					IsLocalMaster:   false,
					UniqueRepoCount: 0,
				})
			}
			continue
		}

		if !opts.isJSON {
			fmt.Println()
			pterm.Bold.Printf("[Node: %s (%s) - Delegated %d Unique Repositories]\n",
				d.conn.Alias, d.conn.IPAddress, len(d.uniqueURLs))
		}

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
		if isWindowsNode(d.conn) {
			remoteCmd = "gitmap.exe " + remoteCmd
			shell = "powershell"
		}

		out, errExec := currentSSHExecutor.Execute(d.conn, remoteCmd, shell)
		if errExec != nil {
			if !opts.isJSON {
				pterm.Error.Printf("Failed to execute on node %s: %v\n", d.conn.Alias, errExec)
			}
			remoteResults = append(remoteResults, NodeSummaryResult{
				NodeAlias:       d.conn.Alias,
				Host:            d.conn.IPAddress,
				IsOnline:        true,
				IsLocalMaster:   false,
				UniqueRepoCount: len(d.uniqueURLs),
				ErrorMessage:    errExec.Error(),
			})
			continue
		}

		if opts.isJSON {
			var remPayload interface{}
			cleanRem := extractJSONPayload(out)
			if err := json.Unmarshal([]byte(cleanRem), &remPayload); err != nil {
				remPayload = out
			}
			remoteResults = append(remoteResults, NodeSummaryResult{
				NodeAlias:       d.conn.Alias,
				Host:            d.conn.IPAddress,
				IsOnline:        true,
				IsLocalMaster:   false,
				UniqueRepoCount: len(d.uniqueURLs),
				Payload:         remPayload,
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

// RunNodesSummary delegates a single repository release summary across the fleet.
func RunNodesSummary(args []string) error {
	if len(args) == 0 {
		return cmdsummary.RunSummary(args)
	}

	targetRepo := args[0]
	if strings.HasPrefix(targetRepo, "-") {
		return cmdsummary.RunSummary(args)
	}

	cwd, _ := os.Getwd()
	localRecords := resolveWorkspaceReposHook(cwd)
	isLocal := false
	var localPath string

	if _, err := os.Stat(filepath.Join(targetRepo, ".git")); err == nil {
		isLocal = true
		localPath = targetRepo
	} else if _, err := os.Stat(filepath.Join(cwd, targetRepo, ".git")); err == nil {
		isLocal = true
		localPath = filepath.Join(cwd, targetRepo)
	} else {
		targetNorm := normalizeRepoKey(targetRepo)
		for _, r := range localRecords {
			if normalizeRepoKey(r.RepoName) == targetNorm ||
				normalizeRepoKey(r.Slug) == targetNorm ||
				normalizeRepoKey(filepath.Base(r.AbsolutePath)) == targetNorm ||
				normalizeRepoKey(r.HTTPSUrl) == targetNorm ||
				normalizeRepoKey(r.SSHUrl) == targetNorm ||
				normalizeRepoKey(r.DiscoveredURL) == targetNorm {
				isLocal = true
				if r.AbsolutePath != "" {
					localPath = r.AbsolutePath
				} else {
					localPath = r.RelativePath
				}
				break
			}
		}
	}

	if isLocal {
		localArgs := make([]string, len(args))
		copy(localArgs, args)
		if localPath != "" && targetRepo != localPath {
			localArgs[0] = localPath
		}
		return cmdsummary.RunSummary(localArgs)
	}

	conns, errConns := fetchAllSSHConnectionsHook()
	if errConns != nil || len(conns) == 0 {
		return cmdsummary.RunSummary(args)
	}

	pterm.Info.Printf("[GitMap Nodes] Querying fleet for repository '%s'...\n", targetRepo)
	for _, conn := range conns {
		cmd := fmt.Sprintf("gitmap summary %s", strings.Join(args, " "))
		shell := "bash"
		if isWindowsNode(conn) {
			cmd = fmt.Sprintf("gitmap.exe summary %s", strings.Join(args, " "))
			shell = "powershell"
		}
		out, errExec := currentSSHExecutor.Execute(conn, cmd, shell)
		if errExec == nil && !strings.Contains(out, "not a valid git repository") && !strings.Contains(out, "not found") {
			pterm.Success.Printf("[Node: %s (%s)]\n", conn.Alias, conn.IPAddress)
			fmt.Println(out)
			return nil
		}
	}

	return apperror.NewSimple(fmt.Sprintf("repository '%s' not found across local or fleet nodes", targetRepo), "E9080")
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
		if !opts.isJSON {
			pterm.Warning.Println("No remote fleet SSH nodes found; executing pipeline error discovery locally.")
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
					Payload:         localPayload,
				},
				RemoteNodes: []NodeSummaryResult{},
			},
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(envelope)
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
				if !isAlreadyLocal {
					ident := item.RemoteURL
					if ident == "" {
						ident = item.Slug
					}
					if ident == "" {
						ident = item.RepoName
					}
					uniqueURLs = append(uniqueURLs, ident)
				}
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
		cleanJSON := extractJSONPayload(captured)
		if err := json.Unmarshal([]byte(cleanJSON), &localPayload); err != nil {
			localPayload = captured
		}
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
				ErrorMessage:    d.err.Error(),
			})
			continue
		}

		if len(d.uniqueURLs) == 0 {
			if opts.isJSON {
				remoteResults = append(remoteResults, NodeSummaryResult{
					NodeAlias:       d.conn.Alias,
					Host:            d.conn.IPAddress,
					IsOnline:        true,
					IsLocalMaster:   false,
					UniqueRepoCount: 0,
				})
			}
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
			if !opts.isJSON {
				pterm.Error.Printf("Failed to inspect node %s: %v\n", d.conn.Alias, errExec)
			}
			remoteResults = append(remoteResults, NodeSummaryResult{
				NodeAlias:       d.conn.Alias,
				Host:            d.conn.IPAddress,
				IsOnline:        true,
				IsLocalMaster:   false,
				UniqueRepoCount: len(d.uniqueURLs),
				ErrorMessage:    errExec.Error(),
			})
			continue
		}

		if opts.isJSON {
			var remPayload interface{}
			cleanRem := extractJSONPayload(out)
			if err := json.Unmarshal([]byte(cleanRem), &remPayload); err != nil {
				remPayload = out
			}
			remoteResults = append(remoteResults, NodeSummaryResult{
				NodeAlias:       d.conn.Alias,
				Host:            d.conn.IPAddress,
				IsOnline:        true,
				IsLocalMaster:   false,
				UniqueRepoCount: len(d.uniqueURLs),
				Payload:         remPayload,
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

type rawRemoteManifestItem struct {
	RepoName     string `json:"repoName"`
	RepoNameAlt  string `json:"repo_name"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	RemoteURL    string `json:"remoteUrl"`
	RemoteURLAlt string `json:"remote_url"`
	URL          string `json:"url"`
	Path         string `json:"path"`
	PathAlt      string `json:"absolute_path"`
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func parseRemoteManifest(raw string) []RemoteRepoManifestItem {
	var rawItems []rawRemoteManifestItem
	clean := extractJSONPayload(raw)
	if err := json.Unmarshal([]byte(clean), &rawItems); err == nil && len(rawItems) > 0 {
		var items []RemoteRepoManifestItem
		for _, ri := range rawItems {
			name := firstNonEmpty(ri.RepoName, ri.RepoNameAlt, ri.Name)
			slug := ri.Slug
			if slug == "" {
				slug = name
			}
			if name == "" {
				name = slug
			}
			url := firstNonEmpty(ri.RemoteURL, ri.RemoteURLAlt, ri.URL)
			path := firstNonEmpty(ri.Path, ri.PathAlt)

			if name != "" || slug != "" || url != "" {
				items = append(items, RemoteRepoManifestItem{
					RepoName:  name,
					Slug:      slug,
					RemoteURL: url,
					Path:      path,
				})
			}
		}
		if len(items) > 0 {
			return items
		}
	}

	var genericList []map[string]interface{}
	if err := json.Unmarshal([]byte(clean), &genericList); err != nil {
		var env map[string]interface{}
		if errEnv := json.Unmarshal([]byte(clean), &env); errEnv == nil {
			if dataSlice, ok := env["data"].([]interface{}); ok {
				for _, el := range dataSlice {
					if m, ok := el.(map[string]interface{}); ok {
						genericList = append(genericList, m)
					}
				}
			} else if itemSlice, ok := env["items"].([]interface{}); ok {
				for _, el := range itemSlice {
					if m, ok := el.(map[string]interface{}); ok {
						genericList = append(genericList, m)
					}
				}
			}
		}
	}

	var items []RemoteRepoManifestItem
	for _, m := range genericList {
		item := RemoteRepoManifestItem{}
		if v, ok := m["repoName"].(string); ok {
			item.RepoName = v
		} else if v, ok := m["repo_name"].(string); ok {
			item.RepoName = v
		} else if v, ok := m["name"].(string); ok {
			item.RepoName = v
		}
		if v, ok := m["slug"].(string); ok {
			item.Slug = v
		}
		if v, ok := m["remoteUrl"].(string); ok {
			item.RemoteURL = v
		} else if v, ok := m["remote_url"].(string); ok {
			item.RemoteURL = v
		} else if v, ok := m["url"].(string); ok {
			item.RemoteURL = v
		}
		if v, ok := m["path"].(string); ok {
			item.Path = v
		} else if v, ok := m["absolute_path"].(string); ok {
			item.Path = v
		}
		if item.Slug == "" {
			item.Slug = item.RepoName
		}
		if item.RepoName == "" {
			item.RepoName = item.Slug
		}
		if item.RepoName != "" || item.Slug != "" || item.RemoteURL != "" {
			items = append(items, item)
		}
	}

	return items
}
