package cmdssh

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/termtable"
)

type remotePullRepoItem struct {
	RepoName string `json:"repoName"`
	Status   string `json:"status"`
	Changes  string `json:"changes"`
}

type remotePullSummary struct {
	Total         int                  `json:"total"`
	ActiveCount   int                  `json:"activeCount"`
	InactiveCount int                  `json:"inactiveCount"`
	States        []remotePullRepoItem `json:"states"`
	DurationMs    int64                `json:"durationMs"`
}

// RunSSHPullJSON executes pull all-efficient remotely via SSH, parses JSON responses, and renders native tables.
func RunSSHPullJSON(target string, isTable bool) error {
	normalizedTarget := normalizePullTarget(target)
	conns, err := loadTargetNodes(normalizedTarget)
	if err != nil {
		return err
	}
	opts := FleetParallelOptions{
		Target:   normalizedTarget,
		TaskName: "Pull All-Efficient (JSON)",
	}
	results := RunParallelFleetExecution(conns, opts, executeRemotePullWorker)
	renderFleetPullResponses(results, isTable)
	return nil
}

func normalizePullTarget(target string) string {
	low := strings.ToLower(target)
	if low == "" || low == "ssh" || low == "--ssh" || low == "all" {
		return "all-nodes"
	}
	return target
}

func executeRemotePullWorker(c db.SSHConnection) (string, error) {
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)
	if !checkRemoteNodeOnline(c.IPAddress, header) {
		return "", fmt.Errorf("node %s is offline", header)
	}
	client, isConnected := connectSSHClient(c, header)
	if !isConnected {
		return "", fmt.Errorf("connection/auth failed for %s", header)
	}
	defer client.Close()

	return crypto.RunCommand(client, "gitmap pull all-efficient --json", "")
}

func renderFleetPullResponses(results []FleetNodeResult, isTable bool) {
	for _, r := range results {
		if !r.Success {
			continue
		}
		renderSingleNodePullJSON(r, isTable)
	}
}

func renderSingleNodePullJSON(r FleetNodeResult, isTable bool) {
	var summary remotePullSummary
	if err := json.Unmarshal([]byte(strings.TrimSpace(r.Output)), &summary); err != nil {
		fmt.Printf("  [%s|%s] Raw output:\n%s\n", r.Alias, r.IP, r.Output)
		return
	}
	fmt.Printf("\n%s▶ Node [%s] (IP: %s):%s %d active, %d inactive skipped\n",
		constants.ColorCyan, r.Alias, r.IP, constants.ColorReset, summary.ActiveCount, summary.InactiveCount)

	if len(summary.States) == 0 {
		fmt.Println("  (No repositories tracked or updated)")
		return
	}
	renderRemotePullTable(summary.States)
}

func renderRemotePullTable(states []remotePullRepoItem) {
	var rows []termtable.Row
	for _, s := range states {
		rows = append(rows, termtable.Row{
			Cells: []string{s.RepoName, s.Status, s.Changes},
		})
	}
	cfg := termtable.TableConfig{
		Columns: []termtable.Column{
			{Title: "REPOSITORY", Align: termtable.AlignLeft, MinWidth: 25},
			{Title: "STATUS", Align: termtable.AlignLeft, MinWidth: 15},
			{Title: "CHANGES", Align: termtable.AlignLeft, MinWidth: 20},
		},
		Rows: rows,
	}
	termtable.PrintTable(cfg)
}
