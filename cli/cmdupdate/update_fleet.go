package cmdupdate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/termtable"
	"golang.org/x/crypto/ssh"
)

// FleetTarget represents an active cluster node for remote fleet updates.
type FleetTarget struct {
	ID       string
	Alias    string
	IP       string
	Username string
	Port     int
	Password string
	KeyPath  string
	OS       string
}

// FleetUpdateOptions contains parsed parameters for fleet update.
type FleetUpdateOptions struct {
	Target   string
	Pkg      string
	Except   string
	IsAll    bool
	IsDryRun bool
	IsForce  bool
}

// FleetUpdateTelemetry is the structured JSON summary returned from a remote node.
type FleetUpdateTelemetry struct {
	NodeID          string   `json:"node_id,omitempty"`
	Alias           string   `json:"alias,omitempty"`
	IP              string   `json:"ip,omitempty"`
	Success         bool     `json:"success"`
	Updated         []string `json:"updated,omitempty"`
	Failed          []string `json:"failed,omitempty"`
	CurrentVersion  string   `json:"current_version,omitempty"`
	PreviousVersion string   `json:"previous_version,omitempty"`
	DurationMs      int64    `json:"duration_ms,omitempty"`
	Details         string   `json:"details,omitempty"`
}

// FleetUpdateNodeResult stores execution outcome on a single node.
type FleetUpdateNodeResult struct {
	Alias      string
	IP         string
	IsSuccess  bool
	IsOffline  bool
	Status     string
	Telemetry  FleetUpdateTelemetry
	DurationMs int64
	Details    string
	Error      error
}

// LoadFleetTargetsFn is a mockable target provider.
var LoadFleetTargetsFn = loadDefaultFleetTargets

// ExecuteRemoteUpdateFn is a mockable remote executor.
var ExecuteRemoteUpdateFn = executeDefaultRemoteUpdate

// CheckConnLivenessFn is a mockable connectivity liveness checker.
var CheckConnLivenessFn = cmdssh.CheckConnLiveness

// IsFleetUpdateCommand reports whether the command invocation routes to fleet update.
func IsFleetUpdateCommand(cmd string, args []string) bool {
	if cmd == "ua" {
		return true
	}
	if cmd != "update" {
		return false
	}
	return isFleetUpdateArgs(args)
}

func isFleetUpdateArgs(args []string) bool {
	if len(args) == 0 {
		return false
	}
	first := strings.ToLower(args[0])
	if first == "ls" || first == "list" {
		return true
	}
	if first == "--all" || first == "all" || first == "all-nodes" || first == "allnodes" || first == "-a" {
		return true
	}
	if hasExceptFlag(args) {
		return true
	}
	return hasFleetNodeTarget(args)
}

func hasExceptFlag(args []string) bool {
	for _, a := range args {
		if isExceptParam(a) {
			return true
		}
	}
	return false
}

func hasFleetNodeTarget(args []string) bool {
	for _, a := range args {
		if a == "--all" || a == "-a" || strings.HasPrefix(a, "--node=") || strings.HasPrefix(a, "--remote=") {
			return true
		}
	}
	return false
}

// RunFleetUpdateDispatch routes to the appropriate fleet update action.
func RunFleetUpdateDispatch(cmd string, args []string) error {
	if cmd == "ua" {
		return ExecuteFleetUpdate(append([]string{"--all"}, args...))
	}
	if len(args) > 0 && isLSKeyword(args[0]) {
		return ExecuteFleetUpdateLS(args[1:])
	}
	return ExecuteFleetUpdate(args)
}

func isLSKeyword(token string) bool {
	low := strings.ToLower(token)
	return low == "ls" || low == "list"
}

// ExecuteFleetUpdate updates applications across cluster nodes in parallel.
func ExecuteFleetUpdate(args []string) error {
	opts := parseFleetUpdateOptions(args)
	targets, err := LoadFleetTargetsFn()
	if err != nil {
		return apperror.WrapSimple(err, "ExecuteFleetUpdate.LoadFleetTargets")
	}

	exclusionSet := parseFleetExclusionSet(opts.Except)
	filteredTargets, excludedCount := filterFleetTargets(targets, opts.Target, exclusionSet)
	if len(filteredTargets) == 0 {
		printFleetNoTargetsBanner(opts.Pkg, excludedCount)
		return nil
	}

	results := executeParallelFleetUpdate(filteredTargets, opts)
	renderFleetUpdateSummary(results, opts.Pkg, excludedCount)
	return nil
}

func parseFleetUpdateOptions(args []string) FleetUpdateOptions {
	opts := FleetUpdateOptions{
		Pkg: "gitmap",
	}
	for i := 0; i < len(args); i++ {
		a := args[i]
		processUpdateFlag(a, args, &i, &opts)
	}
	if opts.IsAll {
		opts.Pkg = "all"
	}
	return opts
}

func processUpdateFlag(arg string, args []string, index *int, opts *FleetUpdateOptions) {
	if arg == "--all" || arg == "all" || arg == "-a" {
		opts.IsAll = true
		return
	}
	if isExceptParam(arg) {
		consumeExceptFlagUpdate(args, index, opts)
		return
	}
	if isTargetParam(arg) {
		consumeTargetFlagUpdate(args, index, opts)
		return
	}
	if arg == "--dry-run" {
		opts.IsDryRun = true
		return
	}
	if arg == "-f" || arg == "--force" {
		opts.IsForce = true
		return
	}
	isPositionalPkg := !strings.HasPrefix(arg, "-") && arg != "update" && arg != "ua"
	if isPositionalPkg && opts.Pkg == "gitmap" {
		opts.Pkg = arg
	}
}

func isExceptParam(arg string) bool {
	return arg == "--except" || arg == "--excep" || arg == "--exclude" ||
		strings.HasPrefix(arg, "--except=") || strings.HasPrefix(arg, "--excep=") || strings.HasPrefix(arg, "--exclude=")
}

func isTargetParam(arg string) bool {
	return arg == "-t" || arg == "--target" || strings.HasPrefix(arg, "--target=") || strings.HasPrefix(arg, "--node=") || strings.HasPrefix(arg, "--remote=")
}

func consumeExceptFlagUpdate(args []string, index *int, opts *FleetUpdateOptions) {
	arg := args[*index]
	if strings.Contains(arg, "=") {
		opts.Except = strings.SplitN(arg, "=", 2)[1]
		return
	}
	var tokens []string
	for *index+1 < len(args) {
		next := args[*index+1]
		isFlag := strings.HasPrefix(next, "-")
		if isFlag {
			break
		}
		*index++
		tokens = append(tokens, next)
	}
	opts.Except = strings.Join(tokens, ",")
}

func consumeTargetFlagUpdate(args []string, index *int, opts *FleetUpdateOptions) {
	arg := args[*index]
	if strings.Contains(arg, "=") {
		opts.Target = strings.SplitN(arg, "=", 2)[1]
		return
	}
	if *index+1 < len(args) {
		*index++
		opts.Target = args[*index]
	}
}

func parseFleetExclusionSet(rawExcept string) map[string]bool {
	set := make(map[string]bool)
	if rawExcept == "" {
		return set
	}
	parts := strings.Split(rawExcept, ",")
	for _, p := range parts {
		cleaned := strings.ToLower(strings.TrimSpace(p))
		if cleaned != "" {
			set[cleaned] = true
		}
	}
	return set
}

func filterFleetTargets(targets []FleetTarget, targetFilter string, exclusions map[string]bool) ([]FleetTarget, int) {
	var filtered []FleetTarget
	excludedCount := 0
	targetLower := strings.ToLower(targetFilter)
	for _, t := range targets {
		if isFleetTargetMatch(t, targetLower) == false {
			continue
		}
		if isFleetNodeExcluded(t, exclusions) {
			excludedCount++
			continue
		}
		filtered = append(filtered, t)
	}
	return filtered, excludedCount
}

func isFleetTargetMatch(t FleetTarget, target string) bool {
	if target == "" || target == "all" || target == "all-nodes" {
		return true
	}
	if strings.EqualFold(t.Alias, target) {
		return true
	}
	if strings.EqualFold(t.IP, target) {
		return true
	}
	return strings.EqualFold(t.ID, target)
}

func isFleetNodeExcluded(t FleetTarget, exclusions map[string]bool) bool {
	if len(exclusions) == 0 {
		return false
	}
	if exclusions[strings.ToLower(t.ID)] {
		return true
	}
	if exclusions[strings.ToLower(t.Alias)] {
		return true
	}
	if exclusions[strings.ToLower(t.IP)] {
		return true
	}
	userHost := fmt.Sprintf("%s@%s", t.Username, t.IP)
	return exclusions[strings.ToLower(userHost)]
}

func loadDefaultFleetTargets() ([]FleetTarget, error) {
	conns, err := cmdssh.FetchAllSSHConnections()
	if err != nil {
		return nil, apperror.WrapSimple(err, "cmdssh.FetchAllSSHConnections")
	}
	return convertConnectionsToFleetTargets(conns), nil
}

func convertConnectionsToFleetTargets(conns []db.SSHConnection) []FleetTarget {
	seen := make(map[string]bool)
	var targets []FleetTarget

	for _, c := range conns {
		key := strings.ToLower(c.IPAddress)
		if seen[key] {
			continue
		}
		seen[key] = true
		targets = append(targets, FleetTarget{
			ID:       c.Alias,
			Alias:    c.Alias,
			IP:       c.IPAddress,
			Username: c.Username,
			Port:     22,
			Password: c.EncryptedPassword,
			KeyPath:  c.KeyPath,
			OS:       resolveOS(c.OS),
		})
	}
	return targets
}

func resolveOS(osName string) string {
	if osName != "" {
		return osName
	}
	return "windows"
}

func executeParallelFleetUpdate(targets []FleetTarget, opts FleetUpdateOptions) []FleetUpdateNodeResult {
	results := make([]FleetUpdateNodeResult, len(targets))
	var wg sync.WaitGroup
	var mu sync.Mutex

	fmt.Printf("\n%s[FLEET UPDATE]%s Updating '%s' across %d cluster node(s) in parallel...\n\n",
		constants.ColorCyan, constants.ColorReset, opts.Pkg, len(targets))

	for idx, t := range targets {
		wg.Add(1)
		go func(i int, target FleetTarget) {
			defer wg.Done()
			res := runSingleFleetUpdate(target, opts)
			mu.Lock()
			results[i] = res
			mu.Unlock()
		}(idx, t)
	}
	wg.Wait()
	return results
}

func runSingleFleetUpdate(target FleetTarget, opts FleetUpdateOptions) FleetUpdateNodeResult {
	if opts.IsDryRun {
		return FleetUpdateNodeResult{
			Alias:      target.Alias,
			IP:         target.IP,
			IsSuccess:  true,
			Status:     "SUCCESS",
			DurationMs: 0,
			Details:    fmt.Sprintf("[DRY-RUN] Would update %s", opts.Pkg),
		}
	}

	isOnline, reason := CheckConnLivenessFn(context.Background(), target.IP, target.Port, 1000*time.Millisecond)
	if !isOnline {
		res := FleetUpdateNodeResult{
			Alias:      target.Alias,
			IP:         target.IP,
			IsSuccess:  false,
			IsOffline:  true,
			Status:     "OFFLINE",
			DurationMs: 15,
			Details:    fmt.Sprintf("machine is off or unreachable (%s)", reason),
		}
		printSingleFleetUpdateProgress(res)
		return res
	}

	start := time.Now()
	rawOutput, err := ExecuteRemoteUpdateFn(target, opts)
	dur := time.Since(start).Milliseconds()

	telemetry := ParseFleetUpdateTelemetry(rawOutput, target, err)
	isSuccess := err == nil && telemetry.Success
	status := "FAILED"
	if isSuccess {
		status = "SUCCESS"
	}

	details := telemetry.Details
	if details == "" {
		details = formatTelemetryDetails(telemetry)
	}

	res := FleetUpdateNodeResult{
		Alias:      target.Alias,
		IP:         target.IP,
		IsSuccess:  isSuccess,
		Status:     status,
		Telemetry:  telemetry,
		DurationMs: dur,
		Details:    details,
		Error:      err,
	}
	printSingleFleetUpdateProgress(res)
	return res
}

func formatTelemetryDetails(t FleetUpdateTelemetry) string {
	if len(t.Updated) > 0 {
		return fmt.Sprintf("Updated: %s", strings.Join(t.Updated, ", "))
	}
	if t.CurrentVersion != "" {
		return fmt.Sprintf("Version: %s", t.CurrentVersion)
	}
	return "OK"
}

func printSingleFleetUpdateProgress(res FleetUpdateNodeResult) {
	if res.IsSuccess {
		fmt.Printf("  %s✓%s [%s|%s] %s (%dms)\n",
			constants.ColorGreen, constants.ColorReset, res.Alias, res.IP, res.Details, res.DurationMs)
		return
	}
	if res.IsOffline {
		fmt.Printf("  %s●%s [%s|%s] OFFLINE: %s (%dms)\n",
			constants.ColorYellow, constants.ColorReset, res.Alias, res.IP, res.Details, res.DurationMs)
		return
	}
	fmt.Printf("  %s✖%s [%s|%s] FAILED: %s (%dms)\n",
		constants.ColorRed, constants.ColorReset, res.Alias, res.IP, res.Details, res.DurationMs)
}

// ParseFleetUpdateTelemetry parses JSON summary telemetry from remote execution.
func ParseFleetUpdateTelemetry(raw string, target FleetTarget, execErr error) FleetUpdateTelemetry {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return buildFallbackTelemetry(target, execErr, "Empty response")
	}

	var parsed FleetUpdateTelemetry
	if err := json.Unmarshal([]byte(trimmed), &parsed); err == nil {
		populateTelemetryDefaults(&parsed, target, execErr)
		return parsed
	}

	var items []map[string]any
	if err := json.Unmarshal([]byte(trimmed), &items); err == nil {
		return buildArrayTelemetry(items, target, execErr)
	}

	return buildFallbackTelemetry(target, execErr, trimmed)
}

func populateTelemetryDefaults(t *FleetUpdateTelemetry, target FleetTarget, execErr error) {
	if t.NodeID == "" {
		t.NodeID = target.ID
	}
	if t.Alias == "" {
		t.Alias = target.Alias
	}
	if t.IP == "" {
		t.IP = target.IP
	}
	if execErr != nil {
		t.Success = false
		t.Details = execErr.Error()
	}
}

func buildArrayTelemetry(items []map[string]any, target FleetTarget, execErr error) FleetUpdateTelemetry {
	var updated []string
	var failed []string
	for _, item := range items {
		name, _ := item["name"].(string)
		status, _ := item["status"].(string)
		if status == "updated" || status == "success" || status == "ok" {
			updated = append(updated, name)
			continue
		}
		failed = append(failed, name)
	}
	isSuccess := execErr == nil && len(failed) == 0
	return FleetUpdateTelemetry{
		NodeID:  target.ID,
		Alias:   target.Alias,
		IP:      target.IP,
		Success: isSuccess,
		Updated: updated,
		Failed:  failed,
		Details: fmt.Sprintf("%d updated, %d failed", len(updated), len(failed)),
	}
}

func buildFallbackTelemetry(target FleetTarget, execErr error, raw string) FleetUpdateTelemetry {
	isSuccess := execErr == nil
	details := raw
	if execErr != nil {
		details = execErr.Error()
	}
	return FleetUpdateTelemetry{
		NodeID:  target.ID,
		Alias:   target.Alias,
		IP:      target.IP,
		Success: isSuccess,
		Details: details,
	}
}

func executeDefaultRemoteUpdate(target FleetTarget, opts FleetUpdateOptions) (string, error) {
	if opts.IsDryRun {
		return `{"success": true, "details": "dry-run simulated"}`, nil
	}
	if out, ok := tryRestFleetUpdate(target, opts); ok {
		return out, nil
	}
	return executeSSHFleetUpdate(target, opts)
}

func tryRestFleetUpdate(target FleetTarget, opts FleetUpdateOptions) (string, bool) {
	url := fmt.Sprintf("http://%s:49152/api/v1/update", target.IP)
	payload := map[string]any{
		"package": opts.Pkg,
		"force":   opts.IsForce,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return "", false
	}
	client := http.Client{Timeout: 3000 * time.Millisecond}
	resp, reqErr := client.Post(url, "application/json", strings.NewReader(string(b)))
	if reqErr != nil {
		return "", false
	}
	defer resp.Body.Close()
	isOk := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated
	if !isOk {
		return "", false
	}
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return "", false
	}
	return string(body), true
}

func executeSSHFleetUpdate(target FleetTarget, opts FleetUpdateOptions) (string, error) {
	client, err := dialFleetSSH(target)
	if err != nil {
		return "", err
	}
	defer client.Close()

	osType := target.OS
	if probed := cmdssh.ProbeRemoteOSType(client); probed != "" {
		osType = probed
	}

	cmd := resolveFleetUpdateCommand(osType, opts.Pkg)
	shell := resolveFleetShell(osType)
	return crypto.RunCommand(client, cmd, shell)
}

func resolveFleetUpdateCommand(osType, pkg string) string {
	isWin := strings.EqualFold(osType, "windows") || strings.EqualFold(osType, "win")
	switch strings.ToLower(pkg) {
	case "agm", "ag-manager", "antigravity-manager":
		if isWin {
			return constants.AgManagerWindowsInstallCmd
		}
		return constants.AgManagerUnixInstallCmd
	default:
		if isWin {
			return "powershell -NoProfile -ExecutionPolicy Bypass -Command \"irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/cli/scripts/install.ps1 | iex\""
		}
		return "curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/cli/scripts/install.sh | bash"
	}
}

func resolveFleetShell(osType string) string {
	if strings.EqualFold(osType, "windows") || strings.EqualFold(osType, "win") {
		return ""
	}
	return "sh"
}

func dialFleetSSH(target FleetTarget) (*ssh.Client, error) {
	conn := db.SSHConnection{
		Alias:             target.Alias,
		IPAddress:         target.IP,
		Username:          target.Username,
		EncryptedPassword: target.Password,
		KeyPath:           target.KeyPath,
		OS:                target.OS,
	}
	client, isOk := cmdssh.ConnectSSHClient(conn, fmt.Sprintf("[%s|%s]", target.Alias, target.IP))
	if isOk && client != nil {
		return client, nil
	}
	if c, ok := tryDialFleetPassword(target); ok {
		return c, nil
	}
	if c, ok := tryDialFleetKey(target); ok {
		return c, nil
	}
	if c, ok := tryDialFleetCandidateKeys(target); ok {
		return c, nil
	}
	return nil, fmt.Errorf("ssh dial failed for %s@%s", target.Username, target.IP)
}

func tryDialFleetPassword(target FleetTarget) (*ssh.Client, bool) {
	if target.Password == "" {
		return nil, false
	}
	plain, err := crypto.DecryptStoredPassword(target.Password)
	if err != nil || plain == "" {
		plain = target.Password
	}
	c, connErr := crypto.ConnectWithPassword(target.IP, target.Username, plain)
	return c, connErr == nil
}

func tryDialFleetKey(target FleetTarget) (*ssh.Client, bool) {
	if target.KeyPath == "" {
		return nil, false
	}
	c, err := crypto.ConnectWithKey(target.IP, target.Username, target.KeyPath)
	return c, err == nil
}

func tryDialFleetCandidateKeys(target FleetTarget) (*ssh.Client, bool) {
	home, _ := os.UserHomeDir()
	candidateKeys := []string{
		fmt.Sprintf("%s/.ssh/id_ed25519", home),
		fmt.Sprintf("%s/.ssh/id_rsa", home),
	}
	for _, k := range candidateKeys {
		c, err := crypto.ConnectWithKey(target.IP, target.Username, k)
		if err == nil {
			return c, true
		}
	}
	return nil, false
}

func printFleetNoTargetsBanner(pkg string, excludedCount int) {
	fmt.Printf("\n%s[FLEET UPDATE]%s No matching active targets found for package '%s' (excluded: %d).\n\n",
		constants.ColorYellow, constants.ColorReset, pkg, excludedCount)
}

func renderFleetUpdateSummary(results []FleetUpdateNodeResult, pkg string, excludedCount int) {
	if len(results) == 0 {
		return
	}
	successCount, failCount, offlineCount := calculateFleetMetrics(results)
	fmt.Printf("\n%s================================================================================%s\n",
		constants.ColorCyan, constants.ColorReset)
	fmt.Printf(" %sSSH Fleet Update Summary [%s]:%s Total: %d | Succeeded: %d | Failed: %d | Offline: %d | Excluded: %d\n",
		constants.ColorBold, pkg, constants.ColorReset, len(results), successCount, failCount, offlineCount, excludedCount)
	fmt.Printf("%s--------------------------------------------------------------------------------%s\n",
		constants.ColorDim, constants.ColorReset)

	cfg := termtable.TableConfig{
		Columns: []termtable.Column{
			{Title: "ALIAS", Align: termtable.AlignLeft, MinWidth: 15},
			{Title: "IP", Align: termtable.AlignLeft, MinWidth: 16},
			{Title: "STATUS", Align: termtable.AlignLeft, MinWidth: 10},
			{Title: "DURATION", Align: termtable.AlignRight, MinWidth: 10},
			{Title: "TELEMETRY DETAILS", Align: termtable.AlignLeft, MinWidth: 30},
		},
		Rows: buildFleetUpdateRows(results),
	}
	termtable.PrintTable(cfg)
	fmt.Printf("%s================================================================================%s\n\n",
		constants.ColorCyan, constants.ColorReset)
}

func calculateFleetMetrics(results []FleetUpdateNodeResult) (int, int, int) {
	succeeded := 0
	failed := 0
	offline := 0
	for _, r := range results {
		if r.IsSuccess {
			succeeded++
			continue
		}
		if r.IsOffline {
			offline++
			continue
		}
		failed++
	}
	return succeeded, failed, offline
}

func buildFleetUpdateRows(results []FleetUpdateNodeResult) []termtable.Row {
	rows := make([]termtable.Row, 0, len(results))
	for _, r := range results {
		statusStr := constants.ColorRed + "FAILED" + constants.ColorReset
		if r.IsSuccess {
			statusStr = constants.ColorGreen + "SUCCESS" + constants.ColorReset
		} else if r.IsOffline {
			statusStr = constants.ColorYellow + "OFFLINE" + constants.ColorReset
		}
		rows = append(rows, termtable.Row{
			Cells: []string{
				r.Alias,
				r.IP,
				statusStr,
				fmt.Sprintf("%dms", r.DurationMs),
				r.Details,
			},
		})
	}
	return rows
}
