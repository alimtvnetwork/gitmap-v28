package cmdssh

import (
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
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
	client, err := connectTargetNode(c, header)
	if err != nil {
		return "", err
	}

	defer client.Close()
	shellType := resolveRemoteShell(resolveTargetNodeOSFast(client, c))

	return secrets.RunCommand(client, "gitmap pull all-efficient --json", shellType)
}

func connectTargetNode(c db.SSHConnection, header string) (*ssh.Client, error) {
	if !checkRemoteNodeOnline(c.IPAddress, header) {
		return nil, apperror.NewExecutionError(fmt.Sprintf("node %s is offline", header))
	}

	client, err := connectSSHClientWithErr(c, header)
	if err != nil {
		return nil, apperror.NewExecutionError(fmt.Sprintf("connection/auth failed for %s: %v", header, err))
	}

	return client, nil
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
	summary, err := parseRemotePullSummary(r.Output)
	if err != nil {
		fmt.Printf("  [%s|%s] Raw output:\n%s\n", r.Alias, r.IP, r.Output)

		return
	}

	displayNodePullSummary(r, summary)
}

func parseRemotePullSummary(raw string) (*remotePullSummary, error) {
	var summary remotePullSummary
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &summary); err != nil {
		return nil, err
	}

	return &summary, nil
}

func displayNodePullSummary(r FleetNodeResult, summary *remotePullSummary) {
	fmt.Printf("\n%s▶ Node [%s] (IP: %s):%s %d active, %d inactive skipped\n",
		constants.ColorCyan, r.Alias, r.IP, constants.ColorReset, summary.ActiveCount, summary.InactiveCount)
	if len(summary.States) == 0 {
		fmt.Println("  (No repositories tracked or updated)")

		return
	}

	renderRemotePullTable(summary.States)
}

func renderRemotePullTable(states []remotePullRepoItem) {
	cfg := termout.TableConfig{
		Columns: []termout.Column{
			{Title: "REPOSITORY", Align: termout.AlignLeft, MinWidth: 25},
			{Title: "STATUS", Align: termout.AlignLeft, MinWidth: 15},
			{Title: "CHANGES", Align: termout.AlignLeft, MinWidth: 20},
		},
		Rows: buildRemotePullRows(states),
	}
	termout.PrintTable(cfg)
}

func buildRemotePullRows(states []remotePullRepoItem) []termout.Row {
	var rows []termout.Row
	for _, s := range states {
		rows = append(rows, termout.Row{
			Cells: []string{s.RepoName, s.Status, s.Changes},
		})
	}

	return rows
}
