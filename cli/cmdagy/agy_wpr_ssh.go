// Package cmdagy — agy_wpr_ssh.go coordinates SSH aggregation and remote telemetry for WPR.
package cmdagy

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// AggregateSSHWPRAll snapshots and queries WPR status across all SSH fleet nodes.
func AggregateSSHWPRAll(opts WPROptions) error {
	local, _ := SnapshotAllRunningPrompts()
	remote := queryFleetWPRAll()
	all := make([]store.WatchPromptsSummary, 0, len(local)+len(remote))
	all = append(all, local...)
	all = append(all, remote...)

	if opts.IsJSON {
		return printJSON(all)
	}

	RenderWPRAllSummaryTable(all)
	return nil
}

func queryFleetWPRAll() []store.WatchPromptsSummary {
	if SSHConnectionsFetcher == nil {
		return nil
	}

	conns, err := SSHConnectionsFetcher()
	if err != nil || len(conns) == 0 {
		return nil
	}

	return fetchNodesWPRAllParallel(conns)
}

func fetchNodesWPRAllParallel(conns []db.SSHConnection) []store.WatchPromptsSummary {
	var aggregated []store.WatchPromptsSummary
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, c := range conns {
		wg.Add(1)
		go func(conn db.SSHConnection) {
			defer wg.Done()
			records := querySingleNodeWPRAll(conn)
			mu.Lock()
			aggregated = append(aggregated, records...)
			mu.Unlock()
		}(c)
	}

	wg.Wait()
	return aggregated
}

func querySingleNodeWPRAll(c db.SSHConnection) []store.WatchPromptsSummary {
	client, err := dialSSHNodeClient(c)
	if err != nil {
		UpdateWPRMachineCache(c, c.Alias, "offline")
		return nil
	}
	defer client.Close()

	UpdateWPRMachineCache(c, c.Alias, "online")
	cmdStr := "gitmap agy wpr all --json"
	out, runErr := secrets.RunCommand(client, cmdStr, resolveNodeShell(c.OS))
	if runErr != nil {
		return nil
	}

	return parseRemoteWPRSummaries(c, out)
}

func parseRemoteWPRSummaries(c db.SSHConnection, out string) []store.WatchPromptsSummary {
	clean := extractJSONArrayFromOutput(out)
	if len(clean) == 0 {
		return nil
	}

	var list []store.WatchPromptsSummary
	if json.Unmarshal([]byte(clean), &list) != nil {
		return nil
	}

	return list
}

// AggregateSSHWPRLs queries watched projects list across all SSH fleet nodes.
func AggregateSSHWPRLs(opts WPROptions) error {
	st, isRunning := LoadWPRRuntimeStatus()
	localProjects := collectWatchedProjectStatuses(st.TargetSlugs, isRunning)
	remoteProjects := queryFleetWPRProjects()
	all := make([]WPRProjectStatus, 0, len(localProjects)+len(remoteProjects))
	all = append(all, localProjects...)
	all = append(all, remoteProjects...)

	if opts.IsJSON {
		return printJSON(all)
	}

	RenderWPRProjectsTable(all, st, isRunning)
	return nil
}

func queryFleetWPRProjects() []WPRProjectStatus {
	if SSHConnectionsFetcher == nil {
		return nil
	}

	conns, err := SSHConnectionsFetcher()
	if err != nil || len(conns) == 0 {
		return nil
	}

	var aggregated []WPRProjectStatus
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, c := range conns {
		wg.Add(1)
		go func(conn db.SSHConnection) {
			defer wg.Done()
			records := querySingleNodeWPRProjects(conn)
			mu.Lock()
			aggregated = append(aggregated, records...)
			mu.Unlock()
		}(c)
	}

	wg.Wait()
	return aggregated
}

func querySingleNodeWPRProjects(c db.SSHConnection) []WPRProjectStatus {
	client, err := dialSSHNodeClient(c)
	if err != nil {
		UpdateWPRMachineCache(c, c.Alias, "offline")
		return nil
	}
	defer client.Close()

	UpdateWPRMachineCache(c, c.Alias, "online")
	cmdStr := "gitmap agy wpr ls --json"
	out, runErr := secrets.RunCommand(client, cmdStr, resolveNodeShell(c.OS))
	if runErr != nil {
		return nil
	}

	return parseRemoteWPRProjects(c, out)
}

func parseRemoteWPRProjects(c db.SSHConnection, out string) []WPRProjectStatus {
	clean := extractJSONArrayFromOutput(out)
	if len(clean) == 0 {
		return nil
	}

	var list []WPRProjectStatus
	if json.Unmarshal([]byte(clean), &list) != nil {
		return nil
	}

	for i := range list {
		list[i].NodeAlias = c.Alias
		list[i].MachineName = c.Alias
	}

	return list
}

// AggregateSSHWPRStatus aggregates live status across all SSH fleet nodes.
func AggregateSSHWPRStatus(opts WPROptions) error {
	_ = queryFleetWPRProjects()
	st, isRunning := LoadWPRRuntimeStatus()
	cachedMachines := LoadWPRMachineCache()

	if opts.IsJSON {
		payload := map[string]interface{}{
			"Runtime":  st,
			"Machines": cachedMachines,
		}
		return printJSON(payload)
	}

	RenderWPRStatusBox(st, isRunning, cachedMachines)
	return nil
}

// AggregateSSHWPRLogs retrieves logs from remote SSH fleet nodes.
func AggregateSSHWPRLogs(opts WPROptions) error {
	if SSHConnectionsFetcher == nil {
		return RunWPRLogs(opts)
	}

	conns, err := SSHConnectionsFetcher()
	if err != nil || len(conns) == 0 {
		return RunWPRLogs(opts)
	}

	for _, c := range conns {
		printRemoteNodeLogs(c, opts)
	}

	return nil
}

func printRemoteNodeLogs(c db.SSHConnection, opts WPROptions) {
	client, err := dialSSHNodeClient(c)
	if err != nil {
		return
	}
	defer client.Close()

	cmdStr := fmt.Sprintf("gitmap agy wpr logs %s --json", opts.Target)
	out, runErr := secrets.RunCommand(client, cmdStr, resolveNodeShell(c.OS))
	if runErr == nil && len(out) > 0 {
		fmt.Printf("[%s] %s\n", c.Alias, out)
	}
}
