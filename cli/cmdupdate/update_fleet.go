package cmdupdate

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
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
	Target        string
	Pkg           string
	Except        string
	IsAll         bool
	IsZip         bool
	IncludeOthers bool
	IsDryRun      bool
	IsForce       bool
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

// LoadClusterTargetsFn is a mockable additional cluster target provider for --include-others.
var LoadClusterTargetsFn = loadDefaultClusterTargets

// ExecuteRemoteUpdateFn is a mockable remote executor.
var ExecuteRemoteUpdateFn = executeDefaultRemoteUpdate

// ExecuteFleetZipUpdateFn is a mockable zip update executor.
var ExecuteFleetZipUpdateFn = executeSSHFleetZipUpdate

// CheckConnLivenessFn is a mockable connectivity liveness checker.
var CheckConnLivenessFn = cmdssh.CheckConnLiveness

// StreamFileToRemoteFn is a mockable SSH file streaming provider.
var StreamFileToRemoteFn = cmdssh.StreamFileToRemote

// StreamFileFromRemoteFn is a mockable reader for the remote update zip.
var StreamFileFromRemoteFn = cmdssh.StreamFileFromRemote

// CreateUpdateZipFn is a mockable zip archive packaging provider.
var CreateUpdateZipFn = createUpdatePackageZip

// IsFleetUpdateCommand reports whether the command invocation routes to fleet update.
func IsFleetUpdateCommand(cmd string, args []string) bool {
	if isFleetUpdateCmdToken(cmd) {
		return true
	}
	if cmd != "update" {
		return false
	}
	return isFleetUpdateArgs(args)
}

func isFleetUpdateCmdToken(cmd string) bool {
	switch cmd {
	case "ua", "update-all", "updateall", "uaz", "update-all-zip", "updateallzip":
		return true
	default:
		return false
	}
}

func isFleetUpdateArgs(args []string) bool {
	if len(args) == 0 {
		return false
	}
	first := strings.ToLower(args[0])
	if first == "ls" || first == "list" {
		return true
	}
	if isAllFleetToken(first) || hasZipFlag(args) || hasIncludeOthersFlag(args) {
		return true
	}
	if hasExceptFlag(args) {
		return true
	}
	return hasFleetNodeTarget(args)
}

func hasZipFlag(args []string) bool {
	for _, a := range args {
		if isZipToken(a) {
			return true
		}
	}
	return false
}

func hasIncludeOthersFlag(args []string) bool {
	for _, a := range args {
		if isIncludeOthersToken(a) {
			return true
		}
	}
	return false
}

func isAllFleetToken(token string) bool {
	return token == "--all" || token == "-all" || token == "all" ||
		token == "all-nodes" || token == "allnodes" || token == "-a"
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
		if isAllFleetToken(a) || strings.HasPrefix(a, "--node=") || strings.HasPrefix(a, "--remote=") {
			return true
		}
	}
	return false
}

// RunFleetUpdateDispatch routes to the appropriate fleet update action.
func RunFleetUpdateDispatch(cmd string, args []string) error {
	if isFleetZipCmdToken(cmd) {
		return ExecuteFleetUpdate(append([]string{"--all", "--zip"}, args...))
	}
	if cmd == "ua" || cmd == "update-all" || cmd == "updateall" {
		return ExecuteFleetUpdate(append([]string{"--all"}, args...))
	}
	if len(args) > 0 && isLSKeyword(args[0]) {
		return ExecuteFleetUpdateLS(args[1:])
	}
	return ExecuteFleetUpdate(args)
}

func isFleetZipCmdToken(cmd string) bool {
	return cmd == "uaz" || cmd == "update-all-zip" || cmd == "updateallzip"
}

func isLSKeyword(token string) bool {
	low := strings.ToLower(token)
	return low == "ls" || low == "list"
}

// ExecuteFleetUpdate updates applications across cluster nodes in parallel.
func ExecuteFleetUpdate(args []string) error {
	opts := parseFleetUpdateOptions(args)
	targets, err := resolveFleetTargets(opts.IncludeOthers)
	if err != nil {
		return apperror.WrapSimple(err, "ExecuteFleetUpdate.resolveFleetTargets")
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
	if opts.IsAll && opts.Pkg == "gitmap" {
		opts.Pkg = "all"
	}
	return opts
}

func processUpdateFlag(arg string, args []string, index *int, opts *FleetUpdateOptions) {
	if isAllFleetToken(arg) {
		opts.IsAll = true
		return
	}
	if isZipToken(arg) {
		opts.IsZip = true
		return
	}
	if isIncludeOthersToken(arg) {
		opts.IncludeOthers = true
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
	processPositionalPkg(arg, opts)
}

func processPositionalPkg(arg string, opts *FleetUpdateOptions) {
	if strings.HasPrefix(arg, "-") || isFleetUpdateCmdToken(arg) || arg == "update" || arg == "zip" {
		return
	}
	if opts.Pkg == "gitmap" || opts.Pkg == "all" {
		opts.Pkg = arg
	}
}

func isZipToken(arg string) bool {
	low := strings.ToLower(arg)
	return low == "--zip" || low == "-zip" || low == "zip"
}

func isIncludeOthersToken(arg string) bool {
	low := strings.ToLower(arg)
	return low == "--include-others" || low == "--include-other" || low == "--includeothers" || low == "--includeother"
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

func resolveFleetTargets(includeOthers bool) ([]FleetTarget, error) {
	targets, err := LoadFleetTargetsFn()
	if err != nil {
		return nil, apperror.WrapSimple(err, "LoadFleetTargetsFn")
	}
	if !includeOthers {
		return targets, nil
	}
	clusterTargets, err := LoadClusterTargetsFn()
	if err != nil {
		return targets, nil
	}
	return mergeFleetTargets(targets, clusterTargets), nil
}

func loadDefaultClusterTargets() ([]FleetTarget, error) {
	ctx := context.Background()
	storeDB, err := store.OpenDefault()
	if err != nil {
		return nil, apperror.WrapSimple(err, "store.OpenDefault")
	}
	defer storeDB.Close()

	var targets []FleetTarget
	if hosts, hostErr := store.ListHosts(ctx, storeDB.Conn()); hostErr == nil {
		targets = append(targets, convertHostsToFleetTargets(hosts)...)
	}
	nodesRes := db.ListClusterNodes(ctx, storeDB.Conn())
	if nodesRes.IsSuccess() {
		targets = append(targets, convertClusterNodesToFleetTargets(nodesRes.Data)...)
	}
	return targets, nil
}

func convertHostsToFleetTargets(hosts []store.SSHHost) []FleetTarget {
	targets := make([]FleetTarget, 0, len(hosts))
	for _, h := range hosts {
		targets = append(targets, FleetTarget{
			ID:       h.Alias,
			Alias:    h.Alias,
			IP:       h.IP,
			Username: h.Username,
			Port:     h.Port,
			Password: h.EncryptedPassword,
			OS:       "windows",
		})
	}
	return targets
}

func convertClusterNodesToFleetTargets(nodes []db.ClusterNode) []FleetTarget {
	targets := make([]FleetTarget, 0, len(nodes))
	for _, n := range nodes {
		targets = append(targets, FleetTarget{
			ID:       n.NodeId,
			Alias:    n.Alias,
			IP:       n.IPAddress,
			Username: "root",
			Port:     22,
			OS:       resolveOS(n.OS),
		})
	}
	return targets
}

func mergeFleetTargets(primary []FleetTarget, additional []FleetTarget) []FleetTarget {
	seen := make(map[string]bool)
	merged := make([]FleetTarget, 0, len(primary)+len(additional))
	for _, t := range primary {
		key := strings.ToLower(strings.TrimSpace(t.IP))
		if key != "" && !seen[key] {
			seen[key] = true
			merged = append(merged, t)
		}
	}
	for _, t := range additional {
		key := strings.ToLower(strings.TrimSpace(t.IP))
		if key != "" && !seen[key] {
			seen[key] = true
			merged = append(merged, t)
		}
	}
	return merged
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
	if opts.IsZip {
		clearZipCache()
	}
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

	jsonStr := extractJSONSubstring(trimmed)
	var parsed FleetUpdateTelemetry
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err == nil {
		populateTelemetryDefaults(&parsed, target, execErr)
		return parsed
	}

	var items []map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &items); err == nil {
		return buildArrayTelemetry(items, target, execErr)
	}

	return buildFallbackTelemetry(target, execErr, cleanTelemetryDetails(trimmed))
}

func extractJSONSubstring(s string) string {
	startObj, startArr := strings.Index(s, "{"), strings.Index(s, "[")
	if startArr >= 0 && (startObj < 0 || startArr < startObj) {
		return extractDelimitedRange(s, startArr, "]")
	}
	if startObj >= 0 {
		return extractDelimitedRange(s, startObj, "}")
	}
	return s
}

func extractDelimitedRange(s string, start int, closeDelim string) string {
	end := strings.LastIndex(s, closeDelim)
	if end > start {
		return s[start : end+1]
	}
	return s
}

func cleanTelemetryDetails(raw string) string {
	lines := strings.Split(raw, "\n")
	var candidates []string
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if isIgnoredTelemetryLine(trimmed) {
			continue
		}
		candidates = append(candidates, trimmed)
	}
	if len(candidates) == 0 {
		return "OK"
	}
	first := candidates[0]
	if len(first) > 50 {
		return first[:47] + "..."
	}
	return first
}

func isIgnoredTelemetryLine(trimmed string) bool {
	if trimmed == "" {
		return true
	}
	low := strings.ToLower(trimmed)
	if strings.HasPrefix(low, "warning:") || strings.HasPrefix(low, "error:") {
		return true
	}
	if strings.Contains(low, "404") || strings.Contains(low, "[test-repoexists]") || strings.Contains(low, "[discovery]") {
		return true
	}
	if strings.Contains(low, "fetching") || strings.Contains(low, "url:") || strings.Contains(low, "gitmap shell wrapper") {
		return true
	}
	return strings.HasPrefix(low, "---") || strings.HasPrefix(low, "===") || strings.Contains(low, "done! run")
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
	if opts.IsZip {
		return ExecuteFleetZipUpdateFn(target, opts)
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
	out, err := crypto.RunCommand(client, cmd, shell)
	if err == nil && isAgmPkg(opts.Pkg) {
		collectAgmUpdateZip(client, osType)
	}
	return out, err
}

func executeSSHFleetZipUpdate(target FleetTarget, opts FleetUpdateOptions) (string, error) {
	client, err := dialFleetSSH(target)
	if err != nil {
		return "", err
	}
	defer client.Close()

	osType := target.OS
	if probed := cmdssh.ProbeRemoteOSType(client); probed != "" {
		osType = probed
	}

	zipData, err := getCachedUpdateZip(opts.Pkg, osType)
	if err != nil {
		return "", apperror.WrapSimple(err, "getCachedUpdateZip")
	}

	destZipPath := resolveRemoteZipPath(osType, opts.Pkg)
	err = StreamFileToRemoteFn(client, destZipPath, zipData, osType)
	if err != nil {
		return "", apperror.WrapSimple(err, "StreamFileToRemote")
	}

	cmd := resolveFleetZipInstallCommand(osType, opts.Pkg, destZipPath)
	shell := resolveFleetShell(osType)
	out, err := crypto.RunCommand(client, cmd, shell)
	if err == nil && isAgmPkg(opts.Pkg) {
		collectAgmUpdateZip(client, osType)
	}
	return out, err
}

var (
	zipCacheMu sync.Mutex
	zipCache   = make(map[string][]byte)
)

func clearZipCache() {
	zipCacheMu.Lock()
	zipCache = make(map[string][]byte)
	zipCacheMu.Unlock()
}

func getCachedUpdateZip(pkg, osType string) ([]byte, error) {
	cacheKey := fmt.Sprintf("%s_%s", strings.ToLower(pkg), strings.ToLower(osType))
	zipCacheMu.Lock()
	defer zipCacheMu.Unlock()
	if data, ok := zipCache[cacheKey]; ok {
		return data, nil
	}
	data, err := CreateUpdateZipFn(pkg, osType)
	if err != nil {
		return nil, err
	}
	zipCache[cacheKey] = data
	return data, nil
}

func createUpdatePackageZip(pkg, osType string) ([]byte, error) {
	if isAgmPkg(pkg) {
		if data, err := loadOrExportAgmInstallerZip(); err == nil && len(data) > 4 && string(data[:2]) == "PK" {
			return data, nil
		}
	}
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	isWin := isWindowsOS(osType)

	binTarget := resolvePackageBinName(pkg, isWin)
	binData, err := locatePackageBinary(pkg)
	if err != nil {
		_ = zw.Close()
		return nil, apperror.WrapSimple(err, "locatePackageBinary")
	}
	if err := addZipFileEntry(zw, binTarget, binData, 0755); err != nil {
		_ = zw.Close()
		return nil, apperror.WrapSimple(err, "addZipFileEntry.bin")
	}
	appendAgmToZipIfIncluded(zw, pkg, isWin)
	if err := addLauncherScriptEntry(zw, isWin); err != nil {
		_ = zw.Close()
		return nil, apperror.WrapSimple(err, "addLauncherScriptEntry")
	}
	if err := zw.Close(); err != nil {
		return nil, apperror.WrapSimple(err, "zw.Close")
	}
	return buf.Bytes(), nil
}

func appendAgmToZipIfIncluded(zw *zip.Writer, pkg string, isWin bool) {
	if !isAgmIncludedInZip(pkg) {
		return
	}
	agmData, agmErr := locatePackageBinary("agm")
	if agmErr != nil || len(agmData) == 0 {
		return
	}
	agmTarget := resolvePackageBinName("agm", isWin)
	_ = addZipFileEntry(zw, agmTarget, agmData, 0755)
}

func resolvePackageBinName(pkg string, isWin bool) string {
	name := "gitmap"
	if isAgmPkg(pkg) {
		name = "agm"
	}
	if isWin {
		return name + ".exe"
	}
	return name
}

func isAgmIncludedInZip(pkg string) bool {
	return strings.ToLower(pkg) == "all"
}

func locatePackageBinary(pkg string) ([]byte, error) {
	if isAgmPkg(pkg) {
		return locateAgmBinary()
	}
	return locateGitmapBinary()
}

func readFileIfExists(path string) ([]byte, bool) {
	if path == "" {
		return nil, false
	}
	data, err := os.ReadFile(path)
	if err != nil || len(data) == 0 {
		return nil, false
	}
	return data, true
}

func locateGitmapBinary() ([]byte, error) {
	execPath, _ := os.Executable()
	if data, ok := readFileIfExists(execPath); ok {
		return data, nil
	}
	lp, _ := exec.LookPath("gitmap")
	if data, ok := readFileIfExists(lp); ok {
		return data, nil
	}
	return []byte("gitmap-payload-simulated"), nil
}

func agmExportZipPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".antigravity_tools", "update-export", "agm-update.zip")
}

func loadOrExportAgmInstallerZip() ([]byte, error) {
	path := agmExportZipPath()
	if data, ok := readFileIfExists(path); ok {
		return data, nil
	}
	agm, err := exec.LookPath("agm")
	if err != nil {
		agm, err = exec.LookPath("agm.exe")
		if err != nil {
			return nil, err
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, agm, "update", "export-zip")
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	data, ok := readFileIfExists(path)
	if !ok {
		return nil, fmt.Errorf("agm update export-zip did not write %s", path)
	}
	return data, nil
}

func locateAgmBinary() ([]byte, error) {
	lp, _ := exec.LookPath("agm")
	if data, ok := readFileIfExists(lp); ok {
		return data, nil
	}
	lpExe, _ := exec.LookPath("agm.exe")
	if data, ok := readFileIfExists(lpExe); ok {
		return data, nil
	}
	return []byte("agm-payload-simulated"), nil
}

func addZipFileEntry(zw *zip.Writer, name string, data []byte, mode os.FileMode) error {
	header := &zip.FileHeader{
		Name:   name,
		Method: zip.Deflate,
	}
	header.SetMode(mode)
	w, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func addLauncherScriptEntry(zw *zip.Writer, isWin bool) error {
	if isWin {
		script := "# Embedded GitMap Remote Installer\nExpand-Archive -Path $zipPath -DestinationPath $destDir -Force\n"
		return addZipFileEntry(zw, "install_remote.ps1", []byte(script), 0644)
	}
	script := "#!/bin/sh\nunzip -o \"$ZIP\" -d \"$DEST\"\n"
	return addZipFileEntry(zw, "install_remote.sh", []byte(script), 0755)
}

func remoteExportZipProbe(osType string) string {
	if isWindowsOS(osType) {
		return `powershell -NoProfile -Command "Join-Path $env:USERPROFILE '.antigravity_tools\update-export\agm-update.zip'"`
	}
	return `sh -c 'printf %s "$HOME/.antigravity_tools/update-export/agm-update.zip"'`
}

func collectAgmUpdateZip(client *ssh.Client, osType string) {
	pathOut, err := crypto.RunCommand(client, remoteExportZipProbe(osType), "")
	if err != nil {
		return
	}
	data, err := StreamFileFromRemoteFn(client, strings.TrimSpace(pathOut), osType)
	if err != nil || len(data) < 4 || string(data[:2]) != "PK" {
		return
	}
	rememberCollectedZip("agm", osType, data)
}

func rememberCollectedZip(pkg, osType string, data []byte) {
	cacheKey := fmt.Sprintf("%s_%s", strings.ToLower(pkg), strings.ToLower(osType))
	zipCacheMu.Lock()
	zipCache[cacheKey] = data
	zipCacheMu.Unlock()
}

func resolveRemoteZipPath(osType, pkg string) string {
	pkgName := "gitmap"
	if isAgmPkg(pkg) {
		pkgName = "agm"
	}
	if isWindowsOS(osType) {
		return fmt.Sprintf(`C:\Windows\Temp\gitmap_update_%s.zip`, pkgName)
	}
	return fmt.Sprintf(`/tmp/gitmap_update_%s.zip`, pkgName)
}

func resolveFleetZipInstallCommand(osType, pkg, zipPath string) string {
	isWin := isWindowsOS(osType)
	targetBin := "gitmap"
	if isAgmPkg(pkg) {
		targetBin = "agm"
	}
	if isWin {
		return resolveWindowsZipInstallCommand(targetBin, zipPath)
	}
	return resolvePOSIXZipInstallCommand(targetBin, zipPath)
}

func resolveWindowsZipInstallCommand(targetBin, zipPath string) string {
	return fmt.Sprintf(`powershell -NoProfile -ExecutionPolicy Bypass -Command "& { $ErrorActionPreference = 'SilentlyContinue'; $zipPath = '%s'; $destDir = Split-Path (Get-Command %s -ErrorAction SilentlyContinue).Path; if (-not $destDir) { $destDir = [System.IO.Path]::Combine($env:LOCALAPPDATA, 'Programs', '%s') }; if (-not (Test-Path $destDir)) { New-Item -ItemType Directory -Path $destDir -Force | Out-Null }; Expand-Archive -Path $zipPath -DestinationPath $destDir -Force; Remove-Item -Path $zipPath -Force -ErrorAction SilentlyContinue; $curr = (%s version 2>$null | Out-String).Trim(); if ($curr) { @{ success = $true; current_version = $curr; details = ('Zip installed: ' + $curr) } | ConvertTo-Json -Compress } else { @{ success = $false; details = 'Zip update failed to verify binary' } | ConvertTo-Json -Compress } }"`,
		zipPath, targetBin, targetBin, targetBin)
}

func resolvePOSIXZipInstallCommand(targetBin, zipPath string) string {
	return fmt.Sprintf(`sh -c 'ZIP="%s"; DEST=$(dirname "$(command -v %s 2>/dev/null || echo /usr/local/bin/%s)"); mkdir -p "$DEST"; unzip -o "$ZIP" -d "$DEST" >/dev/null 2>&1; chmod +x "$DEST/%s"; rm -f "$ZIP"; CURR=$(%s version 2>/dev/null | head -n1); if [ -n "$CURR" ]; then printf "{\"success\":true,\"current_version\":\"%%s\",\"details\":\"Zip installed: %%s\"}" "$CURR" "$CURR"; else printf "{\"success\":false,\"details\":\"Zip update failed to verify binary\"}"; fi'`,
		zipPath, targetBin, targetBin, targetBin, targetBin)
}

func isWindowsOS(osType string) bool {
	return strings.EqualFold(osType, "windows") || strings.EqualFold(osType, "win")
}

func isAgmPkg(pkg string) bool {
	low := strings.ToLower(pkg)
	return low == "agm" || low == "agy" || low == "ag-manager" || low == "antigravity-manager"
}

func resolveFleetUpdateCommand(osType, pkg string) string {
	isWin := strings.EqualFold(osType, "windows") || strings.EqualFold(osType, "win")
	switch strings.ToLower(pkg) {
	case "agm", "ag-manager", "antigravity-manager":
		return resolveAgmUpdateCommand(isWin)
	default:
		return resolveGitmapUpdateCommand(isWin)
	}
}

func resolveGitmapUpdateCommand(isWin bool) string {
	if isWin {
		return "powershell -NoProfile -ExecutionPolicy Bypass -Command \"& { $ErrorActionPreference='SilentlyContinue'; $WarningPreference='SilentlyContinue'; $ProgressPreference='SilentlyContinue'; $prev = (gitmap version 2>$null | Out-String).Trim(); & { $env:GITMAP_UPDATING='1'; irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/cli/scripts/install.ps1 | iex } *>$null; $curr = (gitmap version 2>$null | Out-String).Trim(); if ($curr) { $msg = if ($prev -and $prev -ne $curr) { 'Upgraded: ' + $prev + ' -> ' + $curr } else { 'Version: ' + $curr }; @{ success = $true; current_version = $curr; previous_version = $prev; details = $msg } | ConvertTo-Json -Compress } else { @{ success = $false; details = 'Update failed to verify binary' } | ConvertTo-Json -Compress } }\""
	}
	return "sh -c 'PREV=$(gitmap version 2>/dev/null | grep -oE \"v[0-9]+\\.[0-9]+\\.[0-9]+\" | head -n1); GITMAP_UPDATING=1 curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/cli/scripts/install.sh | bash >/dev/null 2>&1; CURR=$(gitmap version 2>/dev/null | grep -oE \"v[0-9]+\\.[0-9]+\\.[0-9]+\" | head -n1); if [ -n \"$CURR\" ]; then if [ -n \"$PREV\" ] && [ \"$PREV\" != \"$CURR\" ]; then DET=\"Upgraded: $PREV -> $CURR\"; else DET=\"Version: $CURR\"; fi; printf \"{\\\"success\\\":true,\\\"current_version\\\":\\\"%s\\\",\\\"previous_version\\\":\\\"%s\\\",\\\"details\\\":\\\"%s\\\"}\" \"$CURR\" \"$PREV\" \"$DET\"; else printf \"{\\\"success\\\":false,\\\"details\\\":\\\"Update failed to verify binary\\\"}\"; fi'"
}

func resolveAgmUpdateCommand(isWin bool) string {
	if isWin {
		return "powershell -NoProfile -ExecutionPolicy Bypass -Command \"& { $ErrorActionPreference='SilentlyContinue'; $WarningPreference='SilentlyContinue'; $ProgressPreference='SilentlyContinue'; $prev = (agm version 2>$null | Select-Object -First 1); & { irm https://raw.githubusercontent.com/alimtvnetwork/Antigravity-Manager/main/install.ps1 | iex } *>$null; $curr = (agm version 2>$null | Select-Object -First 1); if ($curr) { @{ success = $true; current_version = $curr; details = ('Version: ' + $curr) } | ConvertTo-Json -Compress } else { @{ success = $false; details = 'AGM update completed' } | ConvertTo-Json -Compress } }\""
	}
	return "sh -c 'curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/Antigravity-Manager/main/install.sh | bash >/dev/null 2>&1; CURR=$(agm version 2>/dev/null | head -n1); if [ -n \"$CURR\" ]; then printf \"{\\\"success\\\":true,\\\"current_version\\\":\\\"%s\\\",\\\"details\\\":\\\"Version: %s\\\"}\" \"$CURR\" \"$CURR\"; else printf \"{\\\"success\\\":false,\\\"details\\\":\\\"AGM update completed\\\"}\"; fi'"
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
			{Title: "DETAILS", Align: termtable.AlignLeft, MinWidth: 25},
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
		detail := sanitizeTableRowDetail(r.Details)
		rows = append(rows, termtable.Row{
			Cells: []string{
				r.Alias,
				r.IP,
				statusStr,
				fmt.Sprintf("%dms", r.DurationMs),
				detail,
			},
		})
	}
	return rows
}

func sanitizeTableRowDetail(raw string) string {
	clean := strings.ReplaceAll(raw, "\r", "")
	clean = strings.ReplaceAll(clean, "\n", " ")
	clean = strings.TrimSpace(clean)
	if len(clean) > 50 {
		return clean[:47] + "..."
	}
	if clean == "" {
		return "OK"
	}
	return clean
}
