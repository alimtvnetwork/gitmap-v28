// Package cmdnodes coordinates deduplicated full summaries and pipeline errors across fleet nodes.
package cmdnodes

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpending"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsummary"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/pterm/pterm"
)

// RemoteRepoManifestItem models a remote repository discovered during the discovery handshake.
type RemoteRepoManifestItem struct {
	RepoName  string `json:"repoName"`
	Slug      string `json:"slug"`
	RemoteURL string `json:"remoteUrl"`
	Path      string `json:"path"`
}

// NodeSummaryResult aggregates execution result for one fleet node.
type NodeSummaryResult struct {
	NodeAlias       string `json:"nodeAlias"`
	Host            string `json:"host"`
	IsOnline        bool   `json:"isOnline"`
	IsLocalMaster   bool   `json:"isLocalMaster"`
	UniqueRepoCount int    `json:"uniqueRepoCount"`
	RawOutput       string `json:"rawOutput,omitempty"`
	ErrorMessage    string `json:"errorMessage,omitempty"`
}

// RunNodesFullSummary coordinates distributed full summary execution with local precedence deduplication.
func RunNodesFullSummary(args []string) error {
	isJSON := false
	withPE := false
	for _, a := range args {
		lower := strings.ToLower(a)
		if lower == "--json" || lower == "-j" {
			isJSON = true
		}
		if lower == "+pe" || lower == "--pe" || strings.Contains(lower, "fspe") {
			withPE = true
		}
	}

	pterm.Info.Println("[GitMap Nodes] Initializing fleet handshake and local repository cataloging...")

	cwd, _ := os.Getwd()
	localRecords := cmdpending.ResolveWorkspaceRepositories(cwd)
	localURLMap := make(map[string]bool)
	localSlugMap := make(map[string]bool)
	for _, rec := range localRecords {
		if rec.HTTPSUrl != "" {
			localURLMap[strings.ToLower(rec.HTTPSUrl)] = true
		}
		if rec.SSHUrl != "" {
			localURLMap[strings.ToLower(rec.SSHUrl)] = true
		}
		if rec.DiscoveredURL != "" {
			localURLMap[strings.ToLower(rec.DiscoveredURL)] = true
		}
		if rec.Slug != "" {
			localSlugMap[strings.ToLower(rec.Slug)] = true
		}
		if rec.RepoName != "" {
			localSlugMap[strings.ToLower(rec.RepoName)] = true
		}
	}

	conns, errConns := cmdssh.FetchAllSSHConnections()
	if errConns != nil || len(conns) == 0 {
		pterm.Warning.Println("No remote fleet SSH nodes found; executing full summary locally.")
		return cmdsummary.RunFullSummary(args)
	}

	pterm.Info.Printf("[Phase 1] Cataloged %d repositories on local master machine.\n", len(localRecords))
	pterm.Info.Printf("[Phase 2] Discovering remote repositories across %d fleet nodes...\n", len(conns))

	type remoteDiscovery struct {
		conn       db.SSHConnection
		uniqueURLs []string
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
				return
			}

			manifest := parseRemoteManifest(out)
			var uniqueURLs []string
			for _, item := range manifest {
				isAlreadyLocal := localURLMap[strings.ToLower(item.RemoteURL)] ||
					localSlugMap[strings.ToLower(item.Slug)] ||
					localSlugMap[strings.ToLower(item.RepoName)]
				if !isAlreadyLocal {
					uniqueURLs = append(uniqueURLs, item.RemoteURL)
				}
			}

			mu.Lock()
			discoveredRemotes = append(discoveredRemotes, remoteDiscovery{
				conn:       conn,
				uniqueURLs: uniqueURLs,
			})
			mu.Unlock()
		}(c)
	}

	wg.Wait()

	totalDelegated := 0
	for _, d := range discoveredRemotes {
		totalDelegated += len(d.uniqueURLs)
	}
	pterm.Success.Printf("[Phase 3] Deduplication applied: %d local repos evaluated locally; %d remote-unique repos delegated.\n",
		len(localRecords), totalDelegated)

	fmt.Println()
	pterm.Bold.Println("[Node: Local Machine]")
	_ = cmdsummary.RunFullSummary(args)

	for _, d := range discoveredRemotes {
		if len(d.uniqueURLs) == 0 {
			continue
		}
		fmt.Println()
		pterm.Bold.Printf("[Node: %s (%s) - Delegated %d Unique Repositories]\n",
			d.conn.Alias, d.conn.IPAddress, len(d.uniqueURLs))

		remoteCmd := "gitmap fs 3"
		if withPE {
			remoteCmd = "gitmap fs+pe"
		}
		if isJSON {
			remoteCmd += " --json"
		}
		shell := "bash"
		if isWindowsNode(d.conn) {
			remoteCmd = "gitmap.exe " + remoteCmd
			shell = "powershell"
		}

		out, errExec := currentSSHExecutor.Execute(d.conn, remoteCmd, shell)
		if errExec != nil {
			pterm.Error.Printf("Failed to execute on node %s: %v\n", d.conn.Alias, errExec)
			continue
		}
		fmt.Println(out)
	}

	return nil
}

// RunNodesSummary delegates a single repository release summary across the fleet.
func RunNodesSummary(args []string) error {
	if len(args) == 0 {
		return cmdsummary.RunSummary(args)
	}

	targetRepo := args[0]
	cwd, _ := os.Getwd()
	localRecords := cmdpending.ResolveWorkspaceRepositories(cwd)
	isLocal := false
	for _, r := range localRecords {
		if strings.EqualFold(r.RepoName, targetRepo) || strings.EqualFold(r.Slug, targetRepo) {
			isLocal = true
			break
		}
	}

	if isLocal {
		return cmdsummary.RunSummary(args)
	}

	conns, errConns := cmdssh.FetchAllSSHConnections()
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
		if errExec == nil && !strings.Contains(out, "not a valid git repository") {
			pterm.Success.Printf("[Node: %s (%s)]\n", conn.Alias, conn.IPAddress)
			fmt.Println(out)
			return nil
		}
	}

	return apperror.NewSimple(fmt.Sprintf("repository '%s' not found across local or fleet nodes", targetRepo), "E9080")
}

// RunNodesPipelineErrorsAll coordinates distributed pe all execution across fleet nodes.
func RunNodesPipelineErrorsAll(args []string) error {
	pterm.Info.Println("[GitMap Nodes] Executing distributed CI/CD pipeline error discovery...")
	fmt.Println("[Node: Local Machine]")
	_ = cmdsummary.RunPipelineErrorsAll(args)

	conns, errConns := cmdssh.FetchAllSSHConnections()
	if errConns != nil || len(conns) == 0 {
		return nil
	}

	for _, conn := range conns {
		fmt.Println()
		pterm.Bold.Printf("[Node: %s (%s)]\n", conn.Alias, conn.IPAddress)
		cmd := "gitmap pe all"
		for _, a := range args {
			if a == "--force-all" || a == "--json" {
				cmd += " " + a
			}
		}
		shell := "bash"
		if isWindowsNode(conn) {
			cmd = "gitmap.exe " + cmd
			shell = "powershell"
		}
		out, errExec := currentSSHExecutor.Execute(conn, cmd, shell)
		if errExec != nil {
			pterm.Error.Printf("Failed to inspect node %s: %v\n", conn.Alias, errExec)
			continue
		}
		fmt.Println(out)
	}

	return nil
}

func parseRemoteManifest(raw string) []RemoteRepoManifestItem {
	var items []RemoteRepoManifestItem
	_ = json.Unmarshal([]byte(extractJSONPayload(raw)), &items)
	return items
}
