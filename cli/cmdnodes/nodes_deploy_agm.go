// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

// DeployAGMResult captures telemetry for accounts deployment to a fleet node.
type DeployAGMResult struct {
	NodeAlias string        `json:"node_alias"`
	Host      string        `json:"host"`
	OS        string        `json:"os"`
	Accounts  int           `json:"accounts"`
	Status    string        `json:"status"`
	Latency   time.Duration `json:"latency"`
	LatencyMs int64         `json:"latency_ms"`
	Error     string        `json:"error,omitempty"`
}

// DeployAGMOptions stores configuration flags for fleet AGM accounts deployment.
type DeployAGMOptions struct {
	Target      string   `json:"target,omitempty"`
	Except      []string `json:"except,omitempty"`
	IncludeMain bool     `json:"include_main"`
	OpenOnly    bool     `json:"open_only"`
	DryRun      bool     `json:"dry_run"`
	IsJSON      bool     `json:"is_json"`
}

// ParseDeployAGMOptions parses CLI arguments into DeployAGMOptions.
func ParseDeployAGMOptions(args []string) (DeployAGMOptions, error) {
	var opts DeployAGMOptions
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if isDeployAGMHelpFlag(arg) {
			return opts, nil
		}
		consumed, err := parseSingleDeployAGMFlag(&opts, args, i)
		if err != nil {
			return opts, err
		}
		if consumed > 0 {
			i += (consumed - 1)
			continue
		}
		if !strings.HasPrefix(arg, "-") && opts.Target == "" {
			opts.Target = arg
		}
	}
	return opts, nil
}

func parseSingleDeployAGMFlag(opts *DeployAGMOptions, args []string, i int) (int, error) {
	arg := args[i]
	if arg == "--include-main" {
		opts.IncludeMain = true
		return 1, nil
	}
	if arg == "--open-only" {
		opts.OpenOnly = true
		return 1, nil
	}
	if arg == "--dry-run" || arg == "-n" {
		opts.DryRun = true
		return 1, nil
	}
	if arg == "--json" || arg == "-j" {
		opts.IsJSON = true
		return 1, nil
	}
	return parseParamDeployAGMFlag(opts, args, i)
}

func parseParamDeployAGMFlag(opts *DeployAGMOptions, args []string, i int) (int, error) {
	arg := args[i]
	if strings.HasPrefix(arg, "--target=") {
		opts.Target = strings.TrimPrefix(arg, "--target=")
		return 1, nil
	}
	if strings.HasPrefix(arg, "-t=") {
		opts.Target = strings.TrimPrefix(arg, "-t=")
		return 1, nil
	}
	if (arg == "--target" || arg == "-t") && i+1 < len(args) {
		opts.Target = args[i+1]
		return 2, nil
	}
	return parseExceptFlag(opts, args, i)
}

func parseExceptFlag(opts *DeployAGMOptions, args []string, i int) (int, error) {
	arg := args[i]
	if strings.HasPrefix(arg, "--except=") {
		appendExceptTokens(opts, strings.TrimPrefix(arg, "--except="))
		return 1, nil
	}
	if strings.HasPrefix(arg, "-e=") {
		appendExceptTokens(opts, strings.TrimPrefix(arg, "-e="))
		return 1, nil
	}
	if strings.HasPrefix(arg, "--exclude=") {
		appendExceptTokens(opts, strings.TrimPrefix(arg, "--exclude="))
		return 1, nil
	}
	if (arg == "--except" || arg == "-e" || arg == "--exclude") && i+1 < len(args) {
		appendExceptTokens(opts, args[i+1])
		return 2, nil
	}
	return 0, nil
}

func appendExceptTokens(opts *DeployAGMOptions, raw string) {
	parts := strings.Split(raw, ",")
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			opts.Except = append(opts.Except, trimmed)
		}
	}
}

func isDeployAGMHelpFlag(arg string) bool {
	return arg == "--help" || arg == "-h" || arg == "help"
}

// ResolveDeployAGMExceptList determines final except slice with default "main" inclusion.
func ResolveDeployAGMExceptList(opts DeployAGMOptions) []string {
	var list []string
	for _, e := range opts.Except {
		clean := strings.TrimSpace(e)
		if clean != "" {
			list = append(list, clean)
		}
	}
	if !opts.IncludeMain && !containsToken(list, "main") {
		list = append(list, "main")
	}
	return list
}

func containsToken(list []string, token string) bool {
	for _, item := range list {
		if strings.EqualFold(item, token) {
			return true
		}
	}
	return false
}

// FilterDeployAGMNodes filters fleet connections according to target, except, and local exclusions.
func FilterDeployAGMNodes(conns []db.SSHConnection, opts DeployAGMOptions) []db.SSHConnection {
	exceptList := ResolveDeployAGMExceptList(opts)
	var filtered []db.SSHConnection
	for _, c := range conns {
		if isLocalMachineConnection(c) {
			continue
		}
		if opts.Target != "" && !matchesTargetNode(c, opts.Target) {
			continue
		}
		if isNodeExcludedByExcept(c, exceptList) {
			continue
		}
		filtered = append(filtered, c)
	}
	return filtered
}

func matchesTargetNode(c db.SSHConnection, target string) bool {
	clean := strings.TrimSpace(target)
	if strings.EqualFold(c.Alias, clean) {
		return true
	}
	return strings.EqualFold(c.IPAddress, clean)
}

func isNodeExcludedByExcept(c db.SSHConnection, exceptList []string) bool {
	for _, ex := range exceptList {
		if strings.EqualFold(c.Alias, ex) {
			return true
		}
		if strings.EqualFold(c.IPAddress, ex) {
			return true
		}
	}
	return false
}

// ResolveLocalAGMToolsDir resolves the ~/.antigravity_tools directory.
func ResolveLocalAGMToolsDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", apperror.WrapSimple(err, "resolve user home directory")
	}
	return filepath.Join(homeDir, ".antigravity_tools"), nil
}

// PackageLocalAGMAccounts packages accounts, accounts.json, and user_tokens.db into a .tar.gz archive.
func PackageLocalAGMAccounts(baseDir string) ([]byte, int, error) {
	stat, err := os.Stat(baseDir)
	if err != nil || !stat.IsDir() {
		return nil, 0, apperror.NewNotFound("directory", "E9060", "directory "+baseDir+" not found")
	}

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	accCount := packageAccountsDir(tw, filepath.Join(baseDir, "accounts"))
	accJsonCount := packageAccountsJSON(tw, filepath.Join(baseDir, "accounts.json"))
	packageUserTokensDB(tw, filepath.Join(baseDir, "user_tokens.db"))

	_ = tw.Close()
	_ = gw.Close()

	totalAccounts := resolveTotalAccountCount(accCount, accJsonCount)
	if totalAccounts == 0 {
		return nil, 0, apperror.NewSimple("no Antigravity accounts found in "+baseDir+" (neither in accounts/ nor accounts.json)", "E9061")
	}

	return buf.Bytes(), totalAccounts, nil
}

func resolveTotalAccountCount(accDirCount, accJsonCount int) int {
	if accDirCount > 0 {
		return accDirCount
	}
	return accJsonCount
}

func packageAccountsDir(tw *tar.Writer, accountsDir string) int {
	stat, err := os.Stat(accountsDir)
	if err != nil || !stat.IsDir() {
		return 0
	}
	writeTarDirHeader(tw, "accounts/")
	entries, readErr := os.ReadDir(accountsDir)
	if readErr != nil {
		return 0
	}
	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(strings.ToLower(entry.Name()), ".json") {
			continue
		}
		fullPath := filepath.Join(accountsDir, entry.Name())
		if writeTarFileEntry(tw, "accounts/"+entry.Name(), fullPath) {
			count++
		}
	}
	return count
}

func packageAccountsJSON(tw *tar.Writer, jsonPath string) int {
	stat, err := os.Stat(jsonPath)
	if err != nil || stat.IsDir() {
		return 0
	}
	data, readErr := os.ReadFile(jsonPath)
	if readErr != nil {
		return 0
	}
	writeTarFileBytes(tw, "accounts.json", data)
	return countAccountsInJSON(data)
}

func packageUserTokensDB(tw *tar.Writer, tokensPath string) {
	stat, err := os.Stat(tokensPath)
	if err != nil || stat.IsDir() {
		return
	}
	data, readErr := os.ReadFile(tokensPath)
	if readErr != nil {
		return
	}
	writeTarFileBytes(tw, "user_tokens.db", data)
}

func writeTarDirHeader(tw *tar.Writer, dirName string) {
	hdr := &tar.Header{
		Name:     dirName,
		Mode:     0700,
		Typeflag: tar.TypeDir,
		ModTime:  time.Now(),
	}
	_ = tw.WriteHeader(hdr)
}

func writeTarFileEntry(tw *tar.Writer, tarName, filePath string) bool {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}
	return writeTarFileBytes(tw, tarName, data)
}

func writeTarFileBytes(tw *tar.Writer, tarName string, data []byte) bool {
	hdr := &tar.Header{
		Name:     tarName,
		Mode:     0600,
		Size:     int64(len(data)),
		Typeflag: tar.TypeReg,
		ModTime:  time.Now(),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return false
	}
	if _, err := tw.Write(data); err != nil {
		return false
	}
	return true
}

func countAccountsInJSON(data []byte) int {
	var payload struct {
		Accounts []any `json:"accounts"`
	}
	if err := json.Unmarshal(data, &payload); err == nil && len(payload.Accounts) > 0 {
		return len(payload.Accounts)
	}
	var arr []any
	if err := json.Unmarshal(data, &arr); err == nil {
		return len(arr)
	}
	return 0
}

func resolveAGMStagingPath(conn db.SSHConnection) string {
	if isWindowsNode(conn) {
		return `C:\Windows\Temp\agm_accounts_sync.tar.gz`
	}
	return "/tmp/agm_accounts_sync.tar.gz"
}

func resolveAGMExtractCommand(conn db.SSHConnection) (string, string) {
	if isWindowsNode(conn) {
		cmd := `powershell.exe -NoProfile -Command "$dest = Join-Path $env:USERPROFILE '.antigravity_tools'; if (-not (Test-Path $dest)) { New-Item -ItemType Directory -Force -Path $dest | Out-Null }; tar.exe -xzf 'C:\Windows\Temp\agm_accounts_sync.tar.gz' -C $dest; Remove-Item 'C:\Windows\Temp\agm_accounts_sync.tar.gz' -Force -ErrorAction SilentlyContinue"`
		return cmd, "ps"
	}
	cmd := `sh -c "mkdir -p ~/.antigravity_tools/accounts && tar -xzf /tmp/agm_accounts_sync.tar.gz -C ~/.antigravity_tools && rm -f /tmp/agm_accounts_sync.tar.gz && chmod 700 ~/.antigravity_tools ~/.antigravity_tools/accounts && chmod 600 ~/.antigravity_tools/accounts.json ~/.antigravity_tools/user_tokens.db 2>/dev/null || true; chmod 600 ~/.antigravity_tools/accounts/*.json 2>/dev/null || true"`
	return cmd, "sh"
}

func probeAGMNodeLiveness(ctx context.Context, conn db.SSHConnection, timeout time.Duration) bool {
	port := 22
	addr := net.JoinHostPort(conn.IPAddress, strconv.Itoa(port))
	d := net.Dialer{Timeout: timeout}
	rawConn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return false
	}
	_ = rawConn.Close()
	return true
}

func executeDeployAGMWorker(conn db.SSHConnection, opts DeployAGMOptions, tarData []byte, totalAccounts int) DeployAGMResult {
	start := time.Now()
	res := DeployAGMResult{
		NodeAlias: conn.Alias,
		Host:      conn.IPAddress,
		OS:        conn.OS,
	}
	if !probeAGMNodeLiveness(context.Background(), conn, 1500*time.Millisecond) {
		return buildOfflineAGMResult(res, start)
	}
	if opts.DryRun {
		return buildDryRunAGMResult(res, totalAccounts, start)
	}
	return runRemoteAGMDeploy(conn, res, tarData, totalAccounts, start)
}

func buildOfflineAGMResult(res DeployAGMResult, start time.Time) DeployAGMResult {
	res.Status = "OFFLINE"
	res.Accounts = 0
	res.Latency = time.Since(start)
	res.LatencyMs = res.Latency.Milliseconds()
	res.Error = "connection timed out or unreachable"
	return res
}

func buildDryRunAGMResult(res DeployAGMResult, totalAccounts int, start time.Time) DeployAGMResult {
	res.Status = "DRY-RUN"
	res.Accounts = totalAccounts
	res.Latency = time.Since(start)
	res.LatencyMs = res.Latency.Milliseconds()
	return res
}

func runRemoteAGMDeploy(conn db.SSHConnection, res DeployAGMResult, tarData []byte, totalAccounts int, start time.Time) DeployAGMResult {
	client, errConnect := cmdssh.ConnectSSHClientWithErr(conn)
	if errConnect != nil {
		return handleAGMConnectFailure(res, errConnect, start)
	}
	defer client.Close()
	return streamAndExtractAGM(client, conn, res, tarData, totalAccounts, start)
}

func handleAGMConnectFailure(res DeployAGMResult, errConnect error, start time.Time) DeployAGMResult {
	res.Latency = time.Since(start)
	res.LatencyMs = res.Latency.Milliseconds()
	if isOfflineError(errConnect.Error()) {
		res.Status = "OFFLINE"
		res.Error = errConnect.Error()
		return res
	}
	res.Status = "FAILED"
	res.Error = errConnect.Error()
	return res
}

func streamAndExtractAGM(client *ssh.Client, conn db.SSHConnection, res DeployAGMResult, tarData []byte, totalAccounts int, start time.Time) DeployAGMResult {
	stagingPath := resolveAGMStagingPath(conn)
	if errStream := cmdssh.StreamFileToRemote(client, stagingPath, tarData, conn.OS); errStream != nil {
		res.Status = "FAILED"
		res.Error = "stream failed: " + errStream.Error()
		res.Latency = time.Since(start)
		res.LatencyMs = res.Latency.Milliseconds()
		return res
	}
	return executeAGMExtraction(client, conn, res, totalAccounts, start)
}

func executeAGMExtraction(client *ssh.Client, conn db.SSHConnection, res DeployAGMResult, totalAccounts int, start time.Time) DeployAGMResult {
	cmd, shell := resolveAGMExtractCommand(conn)
	out, errExec := crypto.RunCommand(client, cmd, shell)
	res.Latency = time.Since(start)
	res.LatencyMs = res.Latency.Milliseconds()
	if errExec != nil {
		res.Status = "FAILED"
		res.Error = fmt.Sprintf("extraction failed: %s (%s)", errExec.Error(), strings.TrimSpace(out))
		return res
	}
	res.Status = "SUCCESS"
	res.Accounts = totalAccounts
	return res
}

// RunNodesDeployAGMAccounts orchestrates deployment of Antigravity Manager accounts across fleet nodes.
func RunNodesDeployAGMAccounts(args []string) error {
	opts, err := ParseDeployAGMOptions(args)
	if err != nil {
		return err
	}
	baseDir, errDir := ResolveLocalAGMToolsDir()
	if errDir != nil {
		return errDir
	}
	tarData, totalAccounts, errPkg := PackageLocalAGMAccounts(baseDir)
	if errPkg != nil {
		return errPkg
	}
	return executeFleetAGMDeployment(opts, tarData, totalAccounts)
}

func executeFleetAGMDeployment(opts DeployAGMOptions, tarData []byte, totalAccounts int) error {
	conns, err := cmdssh.FetchAllSSHConnections()
	if err != nil || len(conns) == 0 {
		return apperror.NewSimple("no registered fleet SSH nodes found", "E9062")
	}
	targets := FilterDeployAGMNodes(conns, opts)
	if len(targets) == 0 {
		return apperror.NewSimple("no candidate fleet nodes found matching filter criteria", "E9063")
	}
	results := deployToFleetTargets(targets, opts, tarData, totalAccounts)
	return outputAGMResults(results, opts, totalAccounts)
}

func deployToFleetTargets(targets []db.SSHConnection, opts DeployAGMOptions, tarData []byte, totalAccounts int) []DeployAGMResult {
	var results []DeployAGMResult
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, c := range targets {
		wg.Add(1)
		go func(conn db.SSHConnection) {
			defer wg.Done()
			res := executeDeployAGMWorker(conn, opts, tarData, totalAccounts)
			mu.Lock()
			results = append(results, res)
			mu.Unlock()
		}(c)
	}
	wg.Wait()
	return filterOpenOnlyResults(results, opts.OpenOnly)
}

func filterOpenOnlyResults(results []DeployAGMResult, openOnly bool) []DeployAGMResult {
	if !openOnly {
		return results
	}
	var out []DeployAGMResult
	for _, r := range results {
		if r.Status != "OFFLINE" {
			out = append(out, r)
		}
	}
	return out
}

func outputAGMResults(results []DeployAGMResult, opts DeployAGMOptions, totalAccounts int) error {
	if opts.IsJSON {
		return emitDeployAGMJSON(results)
	}
	renderDeployAGMTable(results)
	printDeployAGMSummary(results, totalAccounts, opts.DryRun)
	return nil
}

func emitDeployAGMJSON(results []DeployAGMResult) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}

func formatStatusBadge(status string) string {
	switch status {
	case "SUCCESS":
		return constants.ColorGreen + "● SUCCESS" + constants.ColorReset
	case "DRY-RUN":
		return constants.ColorYellow + "◌ DRY-RUN" + constants.ColorReset
	case "OFFLINE":
		return constants.ColorDim + "○ OFFLINE" + constants.ColorReset
	default:
		return constants.ColorRed + "▲ " + status + constants.ColorReset
	}
}

func renderDeployAGMTable(results []DeployAGMResult) {
	fmt.Println()
	fmt.Printf("  %s%-18s %-16s %-10s %-10s %-12s %-10s%s\n",
		constants.ColorCyan, "NODE ALIAS", "HOST", "OS", "ACCOUNTS", "STATUS", "LATENCY", constants.ColorReset)
	fmt.Println("  " + strings.Repeat("-", 80))
	for _, r := range results {
		statusStr := formatStatusBadge(r.Status)
		osStr := r.OS
		if osStr == "" {
			osStr = "linux"
		}
		fmt.Printf("  %-18s %-16s %-10s %-10d %-12s %-10v\n",
			r.NodeAlias, r.Host, osStr, r.Accounts, statusStr, r.Latency.Round(time.Millisecond))
	}
	fmt.Println()
}

func printDeployAGMSummary(results []DeployAGMResult, totalAccounts int, isDryRun bool) {
	successCount := 0
	offlineCount := 0
	for _, r := range results {
		if r.Status == "SUCCESS" || r.Status == "DRY-RUN" {
			successCount++
		}
		if r.Status == "OFFLINE" {
			offlineCount++
		}
	}
	prefix := constants.ColorGreen + "✓" + constants.ColorReset
	mode := "Deployed"
	if isDryRun {
		mode = "Simulated deployment of"
	}
	fmt.Printf("  %s %s %d account(s) across %d/%d node(s) (%d offline)\n\n",
		prefix, mode, totalAccounts, successCount, len(results), offlineCount)
}
