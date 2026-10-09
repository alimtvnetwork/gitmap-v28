package cmdssh

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
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
	version  string
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
	renderErr := renderFleetPullAllOutcomes(outcomes, isJSON)
	if renderErr != nil {
		return renderErr
	}
	failures := countFleetFailures(outcomes)
	if failures > 0 {
		return apperror.NewSimple(fmt.Sprintf("SSH fleet pull-all had %d node failure(s)", failures), "E_SSH_FLEET_FAIL")
	}

	return nil
}

func countFleetFailures(outcomes []FleetNodePullOutcome) int {
	failures := 0
	for _, o := range outcomes {
		if hasOutcomeFailed(o) {
			failures++
		}
	}
	return failures
}

func hasOutcomeFailed(o FleetNodePullOutcome) bool {
	if o.IsSkipped {
		return false
	}
	return !o.Success
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
			ver := "unknown"
			if isOnline {
				ver = queryFleetNodeVersion(targetConn)
			}
			if ver == "" {
				ver = "unknown"
			}
			results[idx] = fleetNodeLiveness{conn: targetConn, isOnline: isOnline, version: ver}
		}()
	}
	wg.Wait()

	return results
}

func queryFleetNodeVersion(c db.SSHConnection) string {
	if ver := tryRESTNodeVersion(c.IPAddress); ver != "" {
		return ver
	}
	if ver := tryQuickSSHProbe(c); ver != "" {
		return ver
	}
	if c.BuildVersion != "" {
		return cleanFleetVersionString(c.BuildVersion)
	}
	return "unknown"
}

func tryRESTNodeVersion(ip string) string {
	endpoints := []string{
		fmt.Sprintf("http://%s:49152/api/v1/version", ip),
		fmt.Sprintf("http://%s:49152/api/v1/status", ip),
	}
	client := http.Client{Timeout: 1500 * time.Millisecond}
	for _, u := range endpoints {
		resp, err := client.Get(u)
		if err != nil {
			continue
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			continue
		}
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 4096))
		resp.Body.Close()
		if readErr != nil {
			continue
		}
		var statusResp struct {
			Version string `json:"version"`
		}
		if err := json.Unmarshal(body, &statusResp); err == nil && statusResp.Version != "" {
			return cleanFleetVersionString(statusResp.Version)
		}
		trimmed := strings.TrimSpace(string(body))
		if trimmed != "" && !strings.HasPrefix(trimmed, "<") && !strings.HasPrefix(trimmed, "{") {
			return cleanFleetVersionString(trimmed)
		}
	}
	return ""
}

func tryQuickSSHProbe(c db.SSHConnection) string {
	type probeResult struct {
		ver string
	}
	ch := make(chan probeResult, 1)
	go func() {
		client, err := dialNodeWithFallback(c, "")
		if err != nil || client == nil {
			ch <- probeResult{ver: ""}
			return
		}
		defer client.Close()
		osType := resolveTargetNodeOS(client, c)
		ver, isInstalled := queryNodeVersionViaSSH(client, osType)
		if isInstalled && ver != "" {
			ch <- probeResult{ver: cleanFleetVersionString(ver)}
			return
		}
		ch <- probeResult{ver: ""}
	}()

	select {
	case res := <-ch:
		return res.ver
	case <-time.After(1500 * time.Millisecond):
		return ""
	}
}

func cleanFleetVersionString(v string) string {
	s := strings.TrimSpace(v)
	if idx := strings.IndexAny(s, "\r\n"); idx != -1 {
		s = strings.TrimSpace(s[:idx])
	}
	s = strings.TrimPrefix(s, "gitmap")
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "version")
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "v")
	return strings.TrimSpace(s)
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
		if !p.isOnline {
			fmt.Printf("    • Remote Node [%s] (%s): %sOffline (skipped, no task enqueued)%s\n",
				c.Alias, c.IPAddress, constants.ColorYellow, constants.ColorReset)
			continue
		}
		verStr := p.version
		if verStr != "" && verStr != "unknown" && !strings.HasPrefix(verStr, "v") {
			verStr = "v" + verStr
		}
		fmt.Printf("    • Remote Node [%s] (%s) [GitMap %s]: %sOnline → Enqueued (async)%s\n",
			c.Alias, c.IPAddress, verStr, constants.ColorGreen, constants.ColorReset)
	}
	host := resolveLocalHostname()
	localVer := constants.Version
	if !strings.HasPrefix(localVer, "v") {
		localVer = "v" + localVer
	}
	fmt.Printf("    • Current Machine [%s (127.0.0.1)] [GitMap %s]: %sRunning locally (direct execution, not enqueued)%s\n\n",
		host, localVer, constants.ColorGreen, constants.ColorReset)
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
	concurrencyFlags := ParsePASConcurrencyFlags(cleanArgs)
	dispatched := dispatchOnlineFleetPull(online, cleanArgs, concurrencyFlags)

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

const (
	// MaxPASWorkers is the maximum concurrency for remote nodes under the PAS Formula.
	MaxPASWorkers = 2
	// MaxAsyncOpsPerNode is the maximum parallel operations per remote node.
	MaxAsyncOpsPerNode = 2
)

func resolvePASWorkerCount() int {
	if isHighCPUPressure() {
		return 1
	}
	return MaxPASWorkers
}

func isHighCPUPressure() bool {
	if os.Getenv("GITMAP_HIGH_CPU") == "1" {
		return true
	}
	return runtime.NumCPU() <= 2
}

func dispatchOnlineFleetPull(conns []db.SSHConnection, cleanArgs, concurrencyFlags []string) []FleetNodePullOutcome {
	total := len(conns) + 1
	outcomes := make([]FleetNodePullOutcome, total)
	var wg sync.WaitGroup
	dispatchLocalVMPullAsync(&wg, outcomes, cleanArgs, concurrencyFlags)
	dispatchRemoteFleetNodes(&wg, outcomes, conns, concurrencyFlags)

	stopHeartbeat := startFleetPullHeartbeat(len(conns), hasFleetJSONFlag(cleanArgs))
	wg.Wait()
	stopHeartbeat()

	return outcomes
}

func startFleetPullHeartbeat(nodeCount int, isJSON bool) func() {
	if isJSON {
		return func() {}
	}
	stopCh := make(chan struct{})
	doneCh := make(chan struct{})
	start := time.Now()
	go func() {
		defer close(doneCh)
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-stopCh:
				return
			case now := <-ticker.C:
				elapsed := now.Sub(start)
				if elapsed >= 30*time.Second {
					fmt.Printf("    %s[⏳ %ds elapsed] Pulling across SSH fleet (%d online nodes, 1 local)...%s\n",
						constants.ColorDim, int(elapsed.Seconds()), nodeCount, constants.ColorReset)
				}
			}
		}
	}()
	return func() {
		close(stopCh)
		<-doneCh
	}
}

func dispatchLocalVMPullAsync(wg *sync.WaitGroup, outcomes []FleetNodePullOutcome, cleanArgs, concurrencyFlags []string) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		outcomes[0] = executeLocalVMPull(cleanArgs, concurrencyFlags)
	}()
}

func dispatchRemoteFleetNodes(wg *sync.WaitGroup, outcomes []FleetNodePullOutcome, conns []db.SSHConnection, concurrencyFlags []string) {
	workerCount := resolvePASWorkerCount()
	sem := make(chan struct{}, workerCount)
	for idx, c := range conns {
		wg.Add(1)
		slot := idx + 1
		targetConn := c
		go executeBoundedRemotePull(wg, sem, outcomes, slot, targetConn, concurrencyFlags)
	}
}

func executeBoundedRemotePull(wg *sync.WaitGroup, sem chan struct{}, outcomes []FleetNodePullOutcome, slot int, c db.SSHConnection, concurrencyFlags []string) {
	defer wg.Done()
	sem <- struct{}{}
	defer func() { <-sem }()
	outcomes[slot] = executeRemoteNodePull(c, concurrencyFlags)
}

// RunLocalPullAllJSONFn executes pull-all in-process and returns structured JSON output.
var RunLocalPullAllJSONFn func(args []string) (string, error)

func executeLocalVMPull(cleanArgs, concurrencyFlags []string) FleetNodePullOutcome {
	if RunLocalPullAllJSONFn != nil {
		return executeLocalVMPullInProcess(cleanArgs)
	}
	return executeLocalVMPullSubprocess(cleanArgs, concurrencyFlags)
}

func executeLocalVMPullInProcess(cleanArgs []string) FleetNodePullOutcome {
	out, err := RunLocalPullAllJSONFn(cleanArgs)
	if err != nil && len(strings.TrimSpace(out)) == 0 {
		return buildFailedFleetOutcome("Local VM", "127.0.0.1", true, err.Error())
	}
	return parseFleetPullOutcome("Local VM", "127.0.0.1", true, out)
}

func executeLocalVMPullSubprocess(cleanArgs, concurrencyFlags []string) FleetNodePullOutcome {
	cmdName := resolveLocalGitmapExecutable()
	subArgs := []string{"pa", "--json"}
	if len(concurrencyFlags) > 0 {
		subArgs = append(subArgs, concurrencyFlags...)
	}
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

func executeRemoteNodePull(c db.SSHConnection, concurrencyFlags []string) FleetNodePullOutcome {
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
	remoteCmd := buildRemoteFleetPullCmd(concurrencyFlags)
	out, err := runRemotePullJSON(client, c, remoteCmd)
	if err != nil && len(out) == 0 {
		return buildFailedFleetOutcome(c.Alias, c.IPAddress, false, err.Error())
	}

	return parseFleetPullOutcome(c.Alias, c.IPAddress, false, out)
}

func buildRemoteFleetPullCmd(concurrencyFlags []string) string {
	cmd := "gitmap pa --json"
	if len(concurrencyFlags) > 0 {
		cmd += " " + strings.Join(concurrencyFlags, " ")
	} else {
		cmd += " --parallel 2"
	}
	return cmd
}

func ensureRemoteNodeGitmap(client *ssh.Client, c db.SSHConnection) {
	osType := resolveTargetNodeOS(client, c)
	ver, isInstalled := queryNodeVersionViaSSH(client, osType)
	if !isInstalled {
		cmd := BuildGitmapInstallOneLiner(osType, "latest")
		_, _ = secrets.RunCommand(client, cmd, resolveRemoteShell(osType))
		return
	}
	if CompareSemverStrings(ver, "6.349.0") < 0 {
		updateTargetNodeGitmap(client, osType)
	}
}

func updateTargetNodeGitmap(client *ssh.Client, osType string) {
	cmd := resolveRemoteUpdateCommand(osType, "gitmap")
	_, _ = secrets.RunCommand(client, cmd, resolveRemoteShell(osType))
}

func resolveTargetNodeOS(client *ssh.Client, c db.SSHConnection) string {
	osType := c.OS
	if probed := probeRemoteOSType(client); probed != "" {
		return probed
	}
	return osType
}

func runRemotePullJSON(client *ssh.Client, c db.SSHConnection, remoteCmd string) (string, error) {
	out, err := secrets.RunCommand(client, remoteCmd, "")
	if strings.Contains(out, "flag provided but not defined: -json") {
		return handleLegacyRemotePull(client, c, remoteCmd)
	}
	if strings.Contains(out, "pending task already exists for pa") {
		return recoverRemotePendingTask(client, out, remoteCmd)
	}
	return out, err
}

func handleLegacyRemotePull(client *ssh.Client, c db.SSHConnection, remoteCmd string) (string, error) {
	osType := resolveTargetNodeOS(client, c)
	updateTargetNodeGitmap(client, osType)
	return secrets.RunCommand(client, remoteCmd, "")
}

func recoverRemotePendingTask(client *ssh.Client, rawOutput string, remoteCmd string) (string, error) {
	taskID := extractPendingTaskID(rawOutput)
	if taskID != "" {
		_, _ = secrets.RunCommand(client, "gitmap task cancel "+taskID, "")
	}
	return secrets.RunCommand(client, remoteCmd, "")
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
	logFleetFailure(isLocal, msg, ip)
	return FleetNodePullOutcome{
		NodeName: name,
		IP:       ip,
		IsLocal:  isLocal,
		Success:  false,
		ErrorMsg: msg,
	}
}

func logFleetFailure(isLocal bool, msg, ip string) {
	errType := "SSH_DELEGATION_ERROR"
	if isLocal {
		errType = "LOCAL_PULL_ERROR"
	}
	store.LogInternalError("PULL_ALL_SSH", errType, msg, ip, "")
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
