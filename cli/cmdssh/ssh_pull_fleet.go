package cmdssh

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"golang.org/x/crypto/ssh"
)

type fleetPullRepoItem struct {
	RepoName        string `json:"repoName"`
	Status          string `json:"status"`
	Changes         string `json:"changes"`
	ErrorDetails    string `json:"errorDetails,omitempty"`
	RemediationHint string `json:"remediationHint,omitempty"`
}

type fleetPullJSONSummary struct {
	Total        int                 `json:"total"`
	PulledCount  int                 `json:"pulledCount"`
	SuccessCount int                 `json:"successCount"`
	FailedCount  int                 `json:"failedCount"`
	States       []fleetPullRepoItem `json:"states"`
	DurationMs   int64               `json:"durationMs"`
}

// FleetNodePullOutcome represents the pull result for a local or remote node.
type FleetNodePullOutcome struct {
	NodeName     string               `json:"node"`
	IP           string               `json:"ip"`
	IsLocal      bool                 `json:"isLocal"`
	Success      bool                 `json:"success"`
	IsSkipped    bool                 `json:"isSkipped,omitempty"`
	ErrorMsg     string               `json:"errorMsg,omitempty"`
	Summary      fleetPullJSONSummary `json:"summary,omitempty"`
	ActiveStates []fleetPullRepoItem  `json:"activeStates,omitempty"`
}

type fleetNodeLiveness struct {
	conn     db.SSHConnection
	isOnline bool
}

// RunSSHPullAllFleet executes pull-all across SSH fleet nodes and Local VM concurrently.
func RunSSHPullAllFleet(cleanArgs []string) error {
	isJSON := hasFleetJSONFlag(cleanArgs)
	conns, _ := fetchAllSSHConnections()
	probes := probeFleetLiveness(conns)
	if !isJSON {
		printFleetEnqueueBanner(probes)
	}
	outcomes := executeFleetPullAll(probes, cleanArgs)

	return renderFleetPullAllOutcomes(outcomes, isJSON)
}

func probeFleetLiveness(conns []db.SSHConnection) []fleetNodeLiveness {
	results := make([]fleetNodeLiveness, len(conns))
	var wg sync.WaitGroup
	for i, c := range conns {
		wg.Add(1)
		idx := i
		targetConn := c
		go func() {
			defer wg.Done()
			isOnline, _ := CheckConnLiveness(context.Background(), targetConn.IPAddress, 22, 0)
			results[idx] = fleetNodeLiveness{conn: targetConn, isOnline: isOnline}
		}()
	}
	wg.Wait()

	return results
}

func hasFleetJSONFlag(args []string) bool {
	for _, a := range args {
		if strings.EqualFold(a, "--json") || strings.EqualFold(a, "-json") {
			return true
		}
	}

	return false
}

func printFleetEnqueueBanner(probes []fleetNodeLiveness) {
	fmt.Printf("\n%s  Enqueuing 'pull-all' across SSH fleet:%s\n", constants.ColorCyan, constants.ColorReset)
	for _, p := range probes {
		c := p.conn
		if p.isOnline {
			fmt.Printf("    • Remote Node [%s] (%s): %sOnline → Enqueued (async)%s\n",
				c.Alias, c.IPAddress, constants.ColorGreen, constants.ColorReset)
			continue
		}
		fmt.Printf("    • Remote Node [%s] (%s): %sOffline (skipped, no task enqueued)%s\n",
			c.Alias, c.IPAddress, constants.ColorYellow, constants.ColorReset)
	}
	host := resolveLocalHostname()
	fmt.Printf("    • Current Machine [%s (127.0.0.1)]: %sRunning locally (direct execution, not enqueued)%s\n\n",
		host, constants.ColorGreen, constants.ColorReset)
}

func resolveLocalHostname() string {
	host, err := os.Hostname()
	if err != nil || host == "" {
		return "localhost"
	}
	return host
}

func executeFleetPullAll(probes []fleetNodeLiveness, cleanArgs []string) []FleetNodePullOutcome {
	var online []db.SSHConnection
	var skipped []FleetNodePullOutcome
	for _, p := range probes {
		if p.isOnline {
			online = append(online, p.conn)
			continue
		}
		skipped = append(skipped, buildSkippedFleetOutcome(p.conn.Alias, p.conn.IPAddress))
	}
	dispatched := dispatchOnlineFleetPull(online, cleanArgs)

	return append(dispatched, skipped...)
}

func buildSkippedFleetOutcome(name, ip string) FleetNodePullOutcome {
	return FleetNodePullOutcome{
		NodeName:  name,
		IP:        ip,
		IsLocal:   false,
		Success:   false,
		IsSkipped: true,
		ErrorMsg:  "offline (skipped, no task enqueued)",
	}
}

func dispatchOnlineFleetPull(conns []db.SSHConnection, cleanArgs []string) []FleetNodePullOutcome {
	total := len(conns) + 1
	outcomes := make([]FleetNodePullOutcome, total)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		outcomes[0] = executeLocalVMPull(cleanArgs)
	}()
	for idx, c := range conns {
		wg.Add(1)
		slot := idx + 1
		targetConn := c
		go func() {
			defer wg.Done()
			outcomes[slot] = executeRemoteNodePull(targetConn)
		}()
	}
	wg.Wait()

	return outcomes
}

// RunLocalPullAllJSONFn executes pull-all in-process and returns structured JSON output.
var RunLocalPullAllJSONFn func(args []string) (string, error)

func executeLocalVMPull(cleanArgs []string) FleetNodePullOutcome {
	if RunLocalPullAllJSONFn != nil {
		return executeLocalVMPullInProcess(cleanArgs)
	}
	return executeLocalVMPullSubprocess(cleanArgs)
}

func executeLocalVMPullInProcess(cleanArgs []string) FleetNodePullOutcome {
	out, err := RunLocalPullAllJSONFn(cleanArgs)
	if err != nil && len(strings.TrimSpace(out)) == 0 {
		return buildFailedFleetOutcome("Local VM", "127.0.0.1", true, err.Error())
	}
	return parseFleetPullOutcome("Local VM", "127.0.0.1", true, out)
}

func executeLocalVMPullSubprocess(cleanArgs []string) FleetNodePullOutcome {
	cmdName := resolveLocalGitmapExecutable()
	subArgs := []string{"pa", "--json"}
	cmd := exec.Command(cmdName, subArgs...)
	out, runErr := cmd.CombinedOutput()
	if runErr != nil && len(out) == 0 {
		return buildFailedFleetOutcome("Local VM", "127.0.0.1", true, runErr.Error())
	}

	return parseFleetPullOutcome("Local VM", "127.0.0.1", true, string(out))
}

func resolveLocalGitmapExecutable() string {
	exePath, err := os.Executable()
	if err == nil && exePath != "" {
		return exePath
	}
	return "gitmap"
}

func executeRemoteNodePull(c db.SSHConnection) FleetNodePullOutcome {
	isOnline, _ := CheckConnLiveness(context.Background(), c.IPAddress, 22, 0)
	if !isOnline {
		return buildFailedFleetOutcome(c.Alias, c.IPAddress, false, "offline: node unreachable")
	}
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)
	client, isConnected := connectSSHClient(c, header)
	if !isConnected {
		return buildFailedFleetOutcome(c.Alias, c.IPAddress, false, "connection or auth failed")
	}
	defer client.Close()
	ensureRemoteNodeGitmap(client, c)
	out, err := runRemotePullJSON(client, c)
	if err != nil && len(out) == 0 {
		return buildFailedFleetOutcome(c.Alias, c.IPAddress, false, err.Error())
	}

	return parseFleetPullOutcome(c.Alias, c.IPAddress, false, out)
}

func ensureRemoteNodeGitmap(client *ssh.Client, c db.SSHConnection) {
	osType := resolveTargetNodeOS(client, c)
	ver, isInstalled := queryNodeVersionViaSSH(client, osType)
	if !isInstalled {
		cmd := BuildGitmapInstallOneLiner(osType, "latest")
		_, _ = crypto.RunCommand(client, cmd, resolveRemoteShell(osType))
		return
	}
	if CompareSemverStrings(ver, "6.349.0") < 0 {
		updateTargetNodeGitmap(client, osType)
	}
}

func updateTargetNodeGitmap(client *ssh.Client, osType string) {
	cmd := resolveRemoteUpdateCommand(osType, "gitmap")
	_, _ = crypto.RunCommand(client, cmd, resolveRemoteShell(osType))
}

func resolveTargetNodeOS(client *ssh.Client, c db.SSHConnection) string {
	osType := c.OS
	if probed := probeRemoteOSType(client); probed != "" {
		return probed
	}
	return osType
}

const remotePullCmd = "gitmap pa --json --parallel 2"

func runRemotePullJSON(client *ssh.Client, c db.SSHConnection) (string, error) {
	out, err := crypto.RunCommand(client, remotePullCmd, "")
	if strings.Contains(out, "flag provided but not defined: -json") {
		return handleLegacyRemotePull(client, c)
	}
	if strings.Contains(out, "pending task already exists for pa") {
		return recoverRemotePendingTask(client, out)
	}
	return out, err
}

func handleLegacyRemotePull(client *ssh.Client, c db.SSHConnection) (string, error) {
	osType := resolveTargetNodeOS(client, c)
	updateTargetNodeGitmap(client, osType)
	return crypto.RunCommand(client, remotePullCmd, "")
}

func recoverRemotePendingTask(client *ssh.Client, rawOutput string) (string, error) {
	taskID := extractPendingTaskID(rawOutput)
	if taskID != "" {
		_, _ = crypto.RunCommand(client, "gitmap task cancel "+taskID, "")
	}
	return crypto.RunCommand(client, remotePullCmd, "")
}

func extractPendingTaskID(s string) string {
	lower := strings.ToLower(s)
	idx := strings.Index(lower, "(id ")
	if idx == -1 {
		return ""
	}
	sub := s[idx+4:]
	end := strings.Index(sub, ")")
	if end == -1 {
		return ""
	}
	return strings.TrimSpace(sub[:end])
}

func buildFailedFleetOutcome(name, ip string, isLocal bool, msg string) FleetNodePullOutcome {
	return FleetNodePullOutcome{
		NodeName: name,
		IP:       ip,
		IsLocal:  isLocal,
		Success:  false,
		ErrorMsg: msg,
	}
}

func parseFleetPullOutcome(name, ip string, isLocal bool, rawOutput string) FleetNodePullOutcome {
	jsonStr := extractJSONPayload(rawOutput)
	var summary fleetPullJSONSummary
	if err := json.Unmarshal([]byte(jsonStr), &summary); err != nil {
		return buildFailedFleetOutcome(name, ip, isLocal, "invalid JSON: "+strings.TrimSpace(rawOutput))
	}
	active := filterActiveFleetStates(summary.States)

	return FleetNodePullOutcome{
		NodeName:     name,
		IP:           ip,
		IsLocal:      isLocal,
		Success:      true,
		Summary:      summary,
		ActiveStates: active,
	}
}

func extractJSONPayload(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start >= 0 && end > start {
		return s[start : end+1]
	}

	return s
}

func filterActiveFleetStates(states []fleetPullRepoItem) []fleetPullRepoItem {
	var active []fleetPullRepoItem
	for _, s := range states {
		if isFleetRepoItemActive(s) {
			active = append(active, s)
		}
	}

	return active
}

func isFleetRepoItemActive(item fleetPullRepoItem) bool {
	status := strings.ToLower(item.Status)
	changes := strings.ToLower(item.Changes)
	if status == "up-to-date" || changes == "up-to-date" || status == "synced" || changes == "synced" {
		return false
	}
	if status == "failed" || status == "dirty" {
		return true
	}
	if strings.HasPrefix(item.Changes, "+") {
		return true
	}

	return status != "" && status != "up-to-date"
}

func renderFleetPullAllOutcomes(outcomes []FleetNodePullOutcome, isJSON bool) error {
	if isJSON {
		data, _ := json.MarshalIndent(map[string]any{"fleet": outcomes}, "", "  ")
		fmt.Println(string(data))

		return nil
	}
	for _, o := range outcomes {
		renderSingleFleetOutcome(o)
	}
	fmt.Println()

	return nil
}

func renderSingleFleetOutcome(o FleetNodePullOutcome) {
	if o.IsSkipped {
		return
	}
	if !o.Success {
		fmt.Printf("  %s▶ Node [%s] (IP: %s): %sfailed%s\n",
			constants.ColorYellow, o.NodeName, o.IP, constants.ColorRed, constants.ColorReset)
		fmt.Printf("      %s(%s)%s\n", constants.ColorDim, o.ErrorMsg, constants.ColorReset)

		return
	}
	activeCount := len(o.ActiveStates)
	upToDateCount := o.Summary.PulledCount - activeCount
	if upToDateCount < 0 {
		upToDateCount = 0
	}
	printNodeOutcomeHeader(o, activeCount, upToDateCount)
	if activeCount == 0 {
		fmt.Printf("      %s(all repositories are up-to-date)%s\n", constants.ColorDim, constants.ColorReset)

		return
	}
	renderActiveStatesList(o.ActiveStates)
}

func printNodeOutcomeHeader(o FleetNodePullOutcome, activeCount, upToDateCount int) {
	nodeTitle := fmt.Sprintf("[%s] (IP: %s)", o.NodeName, o.IP)
	if o.IsLocal {
		nodeTitle = fmt.Sprintf("Local VM (%s - localhost)", o.IP)
	}
	if activeCount > 0 {
		fmt.Printf("\n  %s▶ %s:%s %d pulled (%d active, %d up-to-date)\n",
			constants.ColorCyan, nodeTitle, constants.ColorReset,
			o.Summary.PulledCount, activeCount, upToDateCount)

		return
	}
	fmt.Printf("\n  %s▶ %s:%s %d pulled (all up-to-date)\n",
		constants.ColorCyan, nodeTitle, constants.ColorReset, o.Summary.PulledCount)
}

func renderActiveStatesList(states []fleetPullRepoItem) {
	colWidth := resolveFleetStatesColWidth(states)
	for _, s := range states {
		styled := styleFleetStatus(s.Status, s.Changes)
		fmt.Printf("      • %-*s  %s\n", colWidth, s.RepoName, styled)
		if s.ErrorDetails != "" {
			fmt.Printf("        %s↳ Reason: %s%s\n", constants.ColorDim, s.ErrorDetails, constants.ColorReset)
		}
		if s.RemediationHint != "" {
			fmt.Printf("        %s↳ Next Step: %s%s\n", constants.ColorCyan, s.RemediationHint, constants.ColorReset)
		}
	}
}

func resolveFleetStatesColWidth(states []fleetPullRepoItem) int {
	colWidth := 26
	for _, s := range states {
		if len(s.RepoName) > colWidth {
			colWidth = len(s.RepoName)
		}
	}
	return colWidth
}

func styleFleetStatus(status, changes string) string {
	if status == "failed" || changes == "failed" {
		return constants.ColorRed + "failed" + constants.ColorReset
	}
	if status == "dirty" || changes == "dirty" {
		return constants.ColorYellow + "dirty" + constants.ColorReset
	}
	if strings.HasPrefix(changes, "+") {
		return formatDiffStats(changes)
	}

	return changes
}

func formatDiffStats(s string) string {
	parts := strings.SplitN(s, " ", 2)
	diffPart := parts[0]
	slashIdx := strings.Index(diffPart, "/")
	if slashIdx == -1 {
		return constants.ColorGreen + s + constants.ColorReset
	}
	ins := diffPart[:slashIdx]
	del := diffPart[slashIdx:]
	colored := constants.ColorGreen + ins + constants.ColorReset + constants.ColorRed + del + constants.ColorReset
	if len(parts) > 1 {
		return colored + " " + constants.ColorDim + parts[1] + constants.ColorReset
	}

	return colored
}
