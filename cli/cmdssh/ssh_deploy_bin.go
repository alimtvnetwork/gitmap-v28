// Package cmdssh — ssh_deploy_bin.go automates binary deployment across remote fleet nodes.
package cmdssh

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type deployBinOptions struct {
	target     string
	filePath   string
	destPath   string
	isDryRun   bool
	hasForce   bool
	hasInstall bool
	timeout    time.Duration
}

// RunSSHDeployBinCLI executes one-liner binary deployment to target fleet nodes.
func RunSSHDeployBinCLI(args []string) error {
	opts := parseDeployBinFlags(args)
	if isDeployBinHelpRequested(args, opts.target) {
		RenderDeployBinHelp()

		return nil
	}

	return executeDeployCLIWorkflow(opts)
}

func isDeployBinHelpRequested(args []string, target string) bool {
	isHelp := hasFlag(args, "--help") || hasFlag(args, "-h") || target == "help"

	return isHelp
}

func executeDeployCLIWorkflow(opts deployBinOptions) error {
	binPath, binSize, binFormat, err := resolveDeployCandidate(opts.filePath)
	if err != nil {
		return err
	}

	nodes, err := resolveDeployNodes(opts.target)
	if err != nil {
		return err
	}

	return executeDeployPipeline(nodes, binPath, binSize, binFormat, opts)
}

func resolveDeployCandidate(custom string) (string, int64, BinaryFormatType, error) {
	binPath, binSize, err := resolveDeployBinary(custom)
	if err != nil {
		return "", 0, BinaryFormatUnknown, err
	}

	binFormat, appErr := DetectBinaryFormat(binPath)
	if appErr != nil {
		return "", 0, BinaryFormatUnknown, appErr
	}

	return binPath, binSize, binFormat, nil
}

func parseDeployBinFlags(args []string) deployBinOptions {
	opts := deployBinOptions{timeout: 60 * time.Second, hasInstall: true}
	var positional []string

	for i := 0; i < len(args); i++ {
		i, positional = processDeployBinArg(args, i, &opts, positional)
	}

	if len(positional) > 0 {
		opts.target = strings.Join(positional, ",")
	}

	return opts
}

func processDeployBinArg(args []string, i int, opts *deployBinOptions, pos []string) (int, []string) {
	arg := args[i]
	if isDeployBinModifierFlag(arg) {
		applyDeployBinModifier(arg, opts)

		return i, pos
	}

	return processDeployBinPositional(args, i, opts, pos)
}

func processDeployBinPositional(args []string, i int, opts *deployBinOptions, pos []string) (int, []string) {
	arg := args[i]
	if isDeployBinValueFlag(arg) && i+1 < len(args) {
		applyDeployBinValue(arg, args[i+1], opts)

		return i + 1, pos
	}

	if !strings.HasPrefix(arg, "-") {
		pos = append(pos, arg)
	}

	return i, pos
}

func isDeployBinModifierFlag(arg string) bool {
	isMod := argMatches(arg, "--dry-run", "-n") ||
		argMatches(arg, "--force", "--force-copy") ||
		argMatches(arg, "--install") ||
		argMatches(arg, "--no-install")

	return isMod
}

func applyDeployBinModifier(arg string, opts *deployBinOptions) {
	if argMatches(arg, "--dry-run", "-n") {
		opts.isDryRun = true
	}

	if argMatches(arg, "--force", "--force-copy") {
		opts.hasForce = true
	}

	applyDeployBinInstallFlag(arg, opts)
}

func applyDeployBinInstallFlag(arg string, opts *deployBinOptions) {
	if argMatches(arg, "--install") {
		opts.hasInstall = true
	}

	if argMatches(arg, "--no-install") {
		opts.hasInstall = false
	}
}

func isDeployBinValueFlag(arg string) bool {
	isVal := argMatches(arg, "--file", "-f") ||
		argMatches(arg, "--dest", "-d") ||
		argMatches(arg, "--timeout", "-t")

	return isVal
}

func applyDeployBinValue(arg, val string, opts *deployBinOptions) {
	if argMatches(arg, "--file", "-f") {
		opts.filePath = val
	}

	if argMatches(arg, "--dest", "-d") {
		opts.destPath = val
	}

	if argMatches(arg, "--timeout", "-t") {
		opts.timeout = parseDurationFlag(val, opts.timeout)
	}
}

func parseDurationFlag(val string, fallback time.Duration) time.Duration {
	d, err := time.ParseDuration(val)
	if err != nil {
		return fallback
	}

	return d
}

func resolveDeployBinary(custom string) (string, int64, error) {
	candidates := []string{custom}
	if custom == "" {
		candidates = appendExecutableCandidate([]string{
			filepath.Join("cli", "gitmap.exe"),
			filepath.Join("cli", "gitmap"),
		})
	}

	return findFirstExistingBinary(candidates)
}

func findFirstExistingBinary(candidates []string) (string, int64, error) {
	for _, path := range candidates {
		if path == "" {
			continue
		}

		info, err := os.Stat(path)
		if err == nil && !info.IsDir() {
			return path, info.Size(), nil
		}
	}

	return "", 0, apperror.NewNotFound("deploy_binary", "E404", "could not locate gitmap binary for deployment")
}

func appendExecutableCandidate(candidates []string) []string {
	exe, err := os.Executable()
	if err != nil {
		return candidates
	}

	return append(candidates, exe)
}

func loadFleetInventory() ([]db.SSHConnection, []store.SSHHost, error) {
	dbConn, err := store.OpenDefault()
	if err != nil {
		return nil, nil, apperror.WrapSimple(err, "open connection database")
	}

	defer dbConn.Close()

	return queryFleetRecords(dbConn)
}

func queryFleetRecords(dbConn *store.DB) ([]db.SSHConnection, []store.SSHHost, error) {
	ctx, sqlDB := dbConn.Context(), dbConn.SQL()
	res := db.GetSSHConnections(ctx, sqlDB)
	if res.IsFailure() {
		return nil, nil, apperror.NewExecutionError(fmt.Sprintf("load ssh nodes: %v", res.Err))
	}

	hosts, _ := store.ListHosts(ctx, sqlDB)

	return res.Data, hosts, nil
}

func findConnByAliasOrIP(conns []db.SSHConnection, alias, ip string) (db.SSHConnection, bool) {
	for _, c := range conns {
		isAliasMatch := strings.EqualFold(c.Alias, alias)
		isIPMatch := ip != "" && strings.EqualFold(c.IPAddress, ip)
		if isAliasMatch || isIPMatch {
			return c, true
		}
	}

	return db.SSHConnection{}, false
}

func matchConnectionBySeq(conns []db.SSHConnection, hosts []store.SSHHost, token string) (db.SSHConnection, bool) {
	clean := strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(token, "seq="), "seq:"), "#")
	idx, err := strconv.Atoi(clean)
	if err != nil || idx < 1 {
		return db.SSHConnection{}, false
	}

	return matchSeqInHostsOrConns(conns, hosts, idx)
}

func matchSeqInHostsOrConns(conns []db.SSHConnection, hosts []store.SSHHost, idx int) (db.SSHConnection, bool) {
	if len(hosts) > 0 && idx <= len(hosts) {
		h := hosts[idx-1]

		return findConnByAliasOrIP(conns, h.Alias, h.IP)
	}

	if idx <= len(conns) {
		return conns[idx-1], true
	}

	return db.SSHConnection{}, false
}

func matchConnectionByID(conns []db.SSHConnection, hosts []store.SSHHost, token string) (db.SSHConnection, bool) {
	for _, h := range hosts {
		if strings.EqualFold(h.ID, token) {
			return findConnByAliasOrIP(conns, h.Alias, h.IP)
		}
	}

	return db.SSHConnection{}, false
}

func matchConnectionByAliasOrIP(conns []db.SSHConnection, token string) (db.SSHConnection, bool) {
	cleanIP := strings.Split(token, ":")[0]
	for _, c := range conns {
		isAliasMatch := strings.EqualFold(c.Alias, token)
		isIPMatch := strings.EqualFold(c.IPAddress, token) || strings.EqualFold(c.IPAddress, cleanIP)
		if isAliasMatch || isIPMatch {
			return c, true
		}
	}

	return db.SSHConnection{}, false
}

func resolveTargetToken(conns []db.SSHConnection, hosts []store.SSHHost, token string) (db.SSHConnection, bool) {
	conn, isMatch := matchConnectionByAliasOrIP(conns, token)
	if isMatch {
		return conn, true
	}

	conn, isSeq := matchConnectionBySeq(conns, hosts, token)
	if isSeq {
		return conn, true
	}

	return matchConnectionByID(conns, hosts, token)
}

func splitTargetTokens(target string) []string {
	parts := strings.Split(target, ",")
	var tokens []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			tokens = append(tokens, trimmed)
		}
	}

	return tokens
}

func isDeployAllTarget(target string) bool {
	low := strings.ToLower(strings.TrimSpace(target))
	isAll := low == "" || low == "all" || low == "--all" || low == "all-nodes"

	return isAll
}

func resolveDeployNodes(target string) ([]db.SSHConnection, error) {
	conns, hosts, err := loadFleetInventory()
	if err != nil {
		return nil, err
	}

	if isDeployAllTarget(target) {
		return conns, nil
	}

	return filterDeployTargets(conns, hosts, target)
}

func filterDeployTargets(conns []db.SSHConnection, hosts []store.SSHHost, target string) ([]db.SSHConnection, error) {
	tokens := splitTargetTokens(target)
	matched := matchTargetTokens(conns, hosts, tokens)
	if len(matched) == 0 {
		return nil, apperror.NewNotFound("target_node", "E404", fmt.Sprintf("target %q not found in fleet", target))
	}

	return matched, nil
}

func matchTargetTokens(conns []db.SSHConnection, hosts []store.SSHHost, tokens []string) []db.SSHConnection {
	var matched []db.SSHConnection
	seen := make(map[string]bool)
	for _, tok := range tokens {
		conn, isFound := resolveTargetToken(conns, hosts, tok)
		if isFound && !seen[conn.Alias] {
			seen[conn.Alias] = true
			matched = append(matched, conn)
		}
	}

	return matched
}

func executeDeployPipeline(nodes []db.SSHConnection, binPath string, binSize int64, binFormat BinaryFormatType, opts deployBinOptions) error {
	printDeployStartHeader(nodes, binPath, binSize, binFormat, opts.isDryRun)
	start := time.Now()
	livenessMap := ProbeNodesLivenessConcurrent(nodes, 1200*time.Millisecond)
	docket := NewDeployFeedbackDocket(binPath, binFormat, len(nodes))

	for _, node := range nodes {
		deploySingleNodeToDocket(node, binPath, binSize, binFormat, opts, livenessMap, docket)
	}

	docket.TotalTime = time.Since(start).Round(time.Millisecond)
	RenderDeployFeedbackDocket(docket)

	return nil
}

func printDeployStartHeader(nodes []db.SSHConnection, binPath string, binSize int64, format BinaryFormatType, isDryRun bool) {
	sizeMB := float64(binSize) / (1024 * 1024)
	mode := "Active"
	if isDryRun {
		mode = "Dry-Run"
	}

	fmt.Printf("\n%s🚀 GitMap Fleet Binary Deployment (%s)%s\n", constants.ColorCyan, mode, constants.ColorReset)
	fmt.Printf("  • Local Binary:  %s (%.1f MB, %s)\n", binPath, sizeMB, FormatBinaryTypeName(format))
	fmt.Printf("  • Target Nodes:  %d node(s)\n\n", len(nodes))
}

func deploySingleNodeToDocket(node db.SSHConnection, binPath string, binSize int64, format BinaryFormatType, opts deployBinOptions, liveness map[string]NodeLivenessResult, d *DeployFeedbackDocket) {
	destPath := resolveBinaryInstallPath(node, opts.destPath)
	fmt.Printf("  ▸ [%s|%s] %s -> %s\n", node.Alias, node.IPAddress, filepath.Base(binPath), destPath)
	if isNodeDeploySkipped(node, format, opts, liveness, d) {
		return
	}

	executeNodeTransfer(node, binPath, destPath, opts, d)
}

func isNodeDeploySkipped(node db.SSHConnection, format BinaryFormatType, opts deployBinOptions, liveness map[string]NodeLivenessResult, d *DeployFeedbackDocket) bool {
	if isNodeOfflineSkip(node, liveness, d) {
		return true
	}

	if isNodeCompatSkip(node, format, opts.hasForce, d) {
		return true
	}

	return isDryRunSkip(node, opts.isDryRun, d)
}

func isNodeOfflineSkip(node db.SSHConnection, liveness map[string]NodeLivenessResult, d *DeployFeedbackDocket) bool {
	live := liveness[node.Alias]
	if live.IsOnline {
		return false
	}

	fmt.Printf("    %s✖ Node offline (%s: %s)%s\n", constants.ColorYellow, SkipReasonNodeOffline, live.Reason, constants.ColorReset)
	d.AddSkippedOffline(makeOfflineOutcome(node, live.Reason))

	return true
}

func makeOfflineOutcome(node db.SSHConnection, reason string) DeployNodeOutcome {
	return DeployNodeOutcome{
		NodeName:   node.Alias,
		NodeIP:     node.IPAddress,
		NodeOS:     node.OS,
		Status:     "SKIPPED",
		SkipReason: SkipReasonNodeOffline,
		Detail:     reason,
	}
}

func isNodeCompatSkip(node db.SSHConnection, format BinaryFormatType, hasForce bool, d *DeployFeedbackDocket) bool {
	targetOS := ResolveTargetNodeOS(node, nil)
	isCompatible, reason := CheckBinaryCompatibility(format, targetOS)
	if isCompatible || hasForce {
		logCompatOverrideNotice(isCompatible, hasForce)

		return false
	}

	rationale := FormatIncompatibleRationale(node.Alias, targetOS, format)
	fmt.Printf("    %s⊘ %s (%s)%s\n", constants.ColorYellow, SkipReasonOSIncompatible, rationale, constants.ColorReset)
	d.AddSkippedOS(makeCompatOutcome(node, targetOS, rationale, reason))

	return true
}

func logCompatOverrideNotice(isCompatible, hasForce bool) {
	if !isCompatible && hasForce {
		fmt.Printf("    %s⚠ Force override: proceeding with transfer%s\n", constants.ColorYellow, constants.ColorReset)
	}
}

func makeCompatOutcome(node db.SSHConnection, targetOS, rationale, reason string) DeployNodeOutcome {
	return DeployNodeOutcome{
		NodeName:   node.Alias,
		NodeIP:     node.IPAddress,
		NodeOS:     targetOS,
		Status:     "SKIPPED",
		SkipReason: SkipReasonOSIncompatible,
		Detail:     rationale,
	}
}

func isDryRunSkip(node db.SSHConnection, isDryRun bool, d *DeployFeedbackDocket) bool {
	if !isDryRun {
		return false
	}

	fmt.Printf("    %s[dry-run] Skipped transfer%s\n", constants.ColorDim, constants.ColorReset)
	d.AddSuccess(makeDryRunOutcome(node))

	return true
}

func makeDryRunOutcome(node db.SSHConnection) DeployNodeOutcome {
	return DeployNodeOutcome{
		NodeName: node.Alias,
		NodeIP:   node.IPAddress,
		NodeOS:   node.OS,
		Status:   "SUCCESS",
		Detail:   "[dry-run] Simulated transfer",
	}
}

func executeNodeTransfer(node db.SSHConnection, binPath, destPath string, opts deployBinOptions, d *DeployFeedbackDocket) {
	client, isConnected := connectSSHClient(node)
	if !isConnected {
		fmt.Printf("    %s✖ Failed to connect via SSH%s\n", constants.ColorRed, constants.ColorReset)
		d.AddFailure(DeployNodeOutcome{NodeName: node.Alias, NodeIP: node.IPAddress, NodeOS: node.OS, Status: "FAILED", Detail: "SSH connection failed"})

		return
	}

	defer client.Close()

	streamAndRecord(client, node, binPath, destPath, opts, d)
}

func streamAndRecord(client *ssh.Client, node db.SSHConnection, binPath, destPath string, opts deployBinOptions, d *DeployFeedbackDocket) {
	t0 := time.Now()
	err := streamBinaryToRemote(client, destPath, binPath, isWindowsOS(node.OS), opts.hasInstall)
	if err != nil {
		handleUploadFailure(err, node.Alias)
		d.AddFailure(DeployNodeOutcome{NodeName: node.Alias, NodeIP: node.IPAddress, NodeOS: node.OS, Status: "FAILED", Detail: err.Error()})

		return
	}

	elapsed := time.Since(t0).Round(time.Millisecond)
	fmt.Printf("    %s✓ Uploaded in %s%s\n", constants.ColorGreen, elapsed, constants.ColorReset)
	ver := verifyRemoteVersion(client, destPath)
	d.AddSuccess(DeployNodeOutcome{NodeName: node.Alias, NodeIP: node.IPAddress, NodeOS: node.OS, Status: "SUCCESS", Version: ver, Duration: elapsed})
}

func verifyRemoteVersion(client *ssh.Client, destPath string) string {
	cmd := fmt.Sprintf("%s version || gitmap version", destPath)
	verOut, err := secrets.RunCommand(client, cmd, "")
	ver := strings.TrimSpace(verOut)
	if err == nil && ver != "" {
		fmt.Printf("    %s✓ Verified remote version: %s%s\n", constants.ColorGreen, ver, constants.ColorReset)

		return ver
	}

	return ""
}

func resolveBinaryInstallPath(node db.SSHConnection, customDest string) string {
	if customDest != "" {
		return customDest
	}

	if isWindowsOS(node.OS) {
		return fmt.Sprintf("C:/Users/%s/AppData/Local/gitmap-cli/gitmap.exe", resolveNodeUser(node))
	}

	return "/usr/local/bin/gitmap"
}

func resolveNodeUser(node db.SSHConnection) string {
	if node.Username != "" {
		return node.Username
	}

	return "Administrator"
}

func isProcessChangeRequiredError(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())
	hasBusy := strings.Contains(msg, "used by another process") ||
		strings.Contains(msg, "sharing violation") ||
		strings.Contains(msg, "text file busy") ||
		strings.Contains(msg, "process cannot access") ||
		strings.Contains(msg, "access is denied") ||
		strings.Contains(msg, "process change")

	return hasBusy
}

func printProcessChangeException(alias string, err error) {
	fmt.Printf("    %s✖ [E_PROCESS_CHANGE_REQUIRED] Cannot deploy gitmap directly with active process on %s%s\n",
		constants.ColorRed, alias, constants.ColorReset)
	fmt.Printf("      %sIn-place binary replacement is blocked because gitmap is actively running.%s\n",
		constants.ColorYellow, constants.ColorReset)
	fmt.Printf("      %sProcess change required: Clone to stage binary -> delegate call to cloned exe -> replace target.%s\n",
		constants.ColorDim, constants.ColorReset)
	if err != nil {
		fmt.Printf("      %sRaw error: %v%s\n", constants.ColorDim, err, constants.ColorReset)
	}
}

func handleUploadFailure(err error, alias string) {
	if isProcessChangeRequiredError(err) {
		printProcessChangeException(alias, err)

		return
	}

	fmt.Printf("    %s✖ Upload error: %v%s\n", constants.ColorRed, err, constants.ColorReset)
}

func streamBinaryToRemote(client *ssh.Client, destPath, localPath string, isWin bool, hasInstall bool) error {
	session, stdin, err := openSessionPipes(client)
	if err != nil {
		return err
	}

	defer session.Close()

	var stderrBuf bytes.Buffer
	session.Stderr = &stderrBuf
	remoteCmd := buildStreamReceiverCmd(destPath, isWin, hasInstall)

	return pipeFileToSession(session, stdin, localPath, remoteCmd, &stderrBuf)
}

func openSessionPipes(client *ssh.Client) (*ssh.Session, io.WriteCloser, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, nil, apperror.WrapSimple(err, "create session")
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		session.Close()

		return nil, nil, apperror.WrapSimple(err, "open stdin")
	}

	return session, stdin, nil
}

func pipeFileToSession(session *ssh.Session, stdin io.WriteCloser, localPath, remoteCmd string, stderrBuf *bytes.Buffer) error {
	file, err := os.Open(localPath)
	if err != nil {
		return apperror.WrapSimple(err, "open local file")
	}

	defer file.Close()

	go copyFileToStdin(stdin, file)

	return runRemoteSession(session, remoteCmd, stderrBuf)
}

func copyFileToStdin(stdin io.WriteCloser, file *os.File) {
	defer stdin.Close()
	_, _ = io.Copy(stdin, file)
}

func runRemoteSession(session *ssh.Session, remoteCmd string, stderrBuf *bytes.Buffer) error {
	runErr := session.Run(remoteCmd)
	if runErr != nil {
		return apperror.WrapSimple(runErr, "remote stream command failed").WithContext("output", strings.TrimSpace(stderrBuf.String()))
	}

	return nil
}

func buildStreamReceiverCmd(destPath string, isWin bool, hasInstall bool) string {
	if isWin {
		return buildWindowsReceiverCmd(destPath)
	}

	return buildLinuxReceiverCmd(destPath)
}

func buildWindowsReceiverCmd(destPath string) string {
	psScript := "$dest = '%s'; $cloned = $dest + '-' + [System.Guid]::NewGuid().ToString('N') + '.tmp.exe'; " +
		"$dir = [System.IO.Path]::GetDirectoryName($dest); " +
		"if ($dir -and -not (Test-Path $dir)) { New-Item -ItemType Directory -Path $dir -Force | Out-Null }; " +
		"$in = [System.Console]::OpenStandardInput(); $out = [System.IO.File]::OpenWrite($cloned); $in.CopyTo($out); $out.Close(); " +
		"Unblock-File -Path $cloned -ErrorAction SilentlyContinue; " +
		"Stop-Process -Name gitmap -Force -ErrorAction SilentlyContinue; Start-Sleep -Milliseconds 200; " +
		"if (Test-Path $dest) { $old = $dest + '.old'; Remove-Item $old -Force -ErrorAction SilentlyContinue; Move-Item -Path $dest -Destination $old -Force -ErrorAction SilentlyContinue }; " +
		"Move-Item -Path $cloned -Destination $dest -Force; " +
		"Unblock-File -Path $dest -ErrorAction SilentlyContinue"

	return fmt.Sprintf(`powershell -NoProfile -Command "%s"`, fmt.Sprintf(psScript, destPath))
}

func buildLinuxReceiverCmd(destPath string) string {
	shScript := `dir=$(dirname '%s'); mkdir -p "$dir"; cloned='%s-cloned'; cat > "$cloned"; chmod 755 "$cloned"; pkill -f gitmap 2>/dev/null || true; mv -f "$cloned" '%s'; chmod 755 '%s'`

	return fmt.Sprintf(`sh -c "`+shScript+`"`, destPath, destPath, destPath, destPath)
}
