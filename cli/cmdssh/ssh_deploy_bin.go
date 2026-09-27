// Package cmdssh — ssh_deploy_bin.go automates binary deployment across remote fleet nodes.
package cmdssh

import (
	"bytes"
	"context"
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
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type deployBinOptions struct {
	target   string
	filePath string
	isDryRun bool
	timeout  time.Duration
}

// RunSSHDeployBinCLI executes one-liner binary deployment to target fleet nodes.
func RunSSHDeployBinCLI(args []string) error {
	opts := parseDeployBinFlags(args)
	if hasFlag(args, "--help") || hasFlag(args, "-h") || opts.target == "help" {
		RenderDeployBinHelp()

		return nil
	}

	binPath, binSize, err := resolveDeployBinary(opts.filePath)
	if err != nil {
		return err
	}

	nodes, err := resolveDeployNodes(opts.target)
	if err != nil {
		return err
	}

	return executeDeployPipeline(nodes, binPath, binSize, opts)
}

func parseDeployBinFlags(args []string) deployBinOptions {
	opts := deployBinOptions{timeout: 60 * time.Second}
	var positional []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case (arg == "--file" || arg == "-f") && i+1 < len(args):
			opts.filePath = args[i+1]
			i++
		case arg == "--dry-run" || arg == "-n":
			opts.isDryRun = true
		case (arg == "--timeout" || arg == "-t") && i+1 < len(args):
			opts.timeout = parseTimeoutFlag(args, i, opts.timeout)
			i++
		case !strings.HasPrefix(arg, "-"):
			positional = append(positional, arg)
		}
	}

	if len(positional) > 0 {
		opts.target = strings.Join(positional, ",")
	}

	return opts
}

func parseTimeoutFlag(args []string, i int, fallback time.Duration) time.Duration {
	if i+1 >= len(args) {
		return fallback
	}

	d, err := time.ParseDuration(args[i+1])
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
			`C:\Users\Administrator\AppData\Local\gitmap-cli\gitmap.exe`,
		})
	}

	for _, path := range candidates {
		if path == "" {
			continue
		}
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
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
		if strings.EqualFold(c.Alias, alias) || (ip != "" && strings.EqualFold(c.IPAddress, ip)) {
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
		if strings.EqualFold(c.Alias, token) || strings.EqualFold(c.IPAddress, token) || strings.EqualFold(c.IPAddress, cleanIP) {
			return c, true
		}
	}
	return db.SSHConnection{}, false
}

func resolveTargetToken(conns []db.SSHConnection, hosts []store.SSHHost, token string) (db.SSHConnection, bool) {
	if conn, isMatch := matchConnectionByAliasOrIP(conns, token); isMatch {
		return conn, true
	}
	if conn, isSeq := matchConnectionBySeq(conns, hosts, token); isSeq {
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
	return low == "" || low == "all" || low == "--all" || low == "all-nodes"
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
	var matched []db.SSHConnection
	seen := make(map[string]bool)

	for _, tok := range tokens {
		conn, isFound := resolveTargetToken(conns, hosts, tok)
		if isFound && !seen[conn.Alias] {
			seen[conn.Alias] = true
			matched = append(matched, conn)
		}
	}
	if len(matched) == 0 {
		return nil, apperror.NewNotFound("target_node", "E404", fmt.Sprintf("target %q not found in fleet (use: alias, ip, seq 1..%d, or id)", target, len(conns)))
	}
	return matched, nil
}

func executeDeployPipeline(nodes []db.SSHConnection, binPath string, binSize int64, opts deployBinOptions) error {
	printDeployStartHeader(nodes, binPath, binSize, opts.isDryRun)
	start := time.Now()
	successCount := 0

	for _, node := range nodes {
		if deploySingleNode(node, binPath, binSize, opts) {
			successCount++
		}
	}

	dur := time.Since(start).Round(time.Millisecond)
	fmt.Printf("\n  %s✓ Fleet binary deployment complete:%s %d/%d nodes updated (%s)\n\n",
		constants.ColorGreen, constants.ColorReset, successCount, len(nodes), dur)

	return nil
}

func printDeployStartHeader(nodes []db.SSHConnection, binPath string, binSize int64, isDryRun bool) {
	sizeMB := float64(binSize) / (1024 * 1024)
	mode := "Active"
	if isDryRun {
		mode = "Dry-Run"
	}

	fmt.Printf("\n%s🚀 GitMap Fleet Binary Deployment (%s)%s\n", constants.ColorCyan, mode, constants.ColorReset)
	fmt.Printf("  • Local Binary:  %s (%.1f MB)\n", binPath, sizeMB)
	fmt.Printf("  • Target Nodes:  %d node(s)\n\n", len(nodes))
}

func deploySingleNode(node db.SSHConnection, binPath string, binSize int64, opts deployBinOptions) bool {
	destPath := resolveBinaryInstallPath(node)
	fmt.Printf("  ▸ [%s|%s] %s -> %s\n", node.Alias, node.IPAddress, filepath.Base(binPath), destPath)
	if opts.isDryRun {
		fmt.Printf("    %s[dry-run] Skipped transfer%s\n", constants.ColorDim, constants.ColorReset)
		return true
	}
	if isOnline, reason := CheckConnLiveness(context.Background(), node.IPAddress, 22, 1*time.Second); !isOnline {
		fmt.Printf("    %s✖ Node offline (%s)%s\n", constants.ColorYellow, reason, constants.ColorReset)
		return false
	}
	client, isConnected := connectSSHClient(node)
	if !isConnected {
		fmt.Printf("    %s✖ Failed to connect via SSH%s\n", constants.ColorRed, constants.ColorReset)
		return false
	}
	defer client.Close()
	return streamAndVerify(client, node, binPath, destPath)
}

func resolveBinaryInstallPath(node db.SSHConnection) string {
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
	return strings.Contains(msg, "used by another process") ||
		strings.Contains(msg, "sharing violation") ||
		strings.Contains(msg, "text file busy") ||
		strings.Contains(msg, "process cannot access") ||
		strings.Contains(msg, "access is denied") ||
		strings.Contains(msg, "process change")
}

func printProcessChangeException(alias string) {
	fmt.Printf("    %s✖ [E_PROCESS_CHANGE_REQUIRED] Cannot deploy gitmap directly with active process on %s%s\n",
		constants.ColorRed, alias, constants.ColorReset)
	fmt.Printf("      %sIn-place binary replacement is blocked because gitmap is actively running.%s\n",
		constants.ColorYellow, constants.ColorReset)
	fmt.Printf("      %sProcess change required: Clone to stage binary -> delegate call to cloned exe -> replace target.%s\n",
		constants.ColorDim, constants.ColorReset)
}

func handleUploadFailure(err error, alias string) {
	if isProcessChangeRequiredError(err) {
		printProcessChangeException(alias)
		return
	}
	fmt.Printf("    %s✖ Upload error: %v%s\n", constants.ColorRed, err, constants.ColorReset)
}

func streamAndVerify(client *ssh.Client, node db.SSHConnection, binPath, destPath string) bool {
	t0 := time.Now()
	if err := streamBinaryToRemote(client, destPath, binPath, isWindowsOS(node.OS)); err != nil {
		handleUploadFailure(err, node.Alias)
		return false
	}

	elapsed := time.Since(t0).Round(time.Millisecond)
	fmt.Printf("    %s✓ Uploaded in %s%s\n", constants.ColorGreen, elapsed, constants.ColorReset)

	verOut, err := crypto.RunCommand(client, "gitmap version", "")
	if err == nil && strings.TrimSpace(verOut) != "" {
		fmt.Printf("    %s✓ Verified remote version: %s%s\n", constants.ColorGreen, strings.TrimSpace(verOut), constants.ColorReset)
	}

	return true
}

func streamBinaryToRemote(client *ssh.Client, destPath, localPath string, isWin bool) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("open stdin: %w", err)
	}

	var stderrBuf bytes.Buffer
	session.Stderr = &stderrBuf

	remoteCmd := buildStreamReceiverCmd(destPath, isWin)
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("open local file: %w", err)
	}
	defer file.Close()

	go func() {
		defer stdin.Close()
		_, _ = io.Copy(stdin, file)
	}()

	runErr := session.Run(remoteCmd)
	if runErr != nil {
		return fmt.Errorf("%w (remote stderr: %s)", runErr, strings.TrimSpace(stderrBuf.String()))
	}
	return nil
}

func buildStreamReceiverCmd(destPath string, isWin bool) string {
	if isWin {
		return fmt.Sprintf(
			`powershell -NoProfile -Command "$dest = '%s'; $cloned = $dest + '-cloned.exe'; $dir = [System.IO.Path]::GetDirectoryName($dest); if ($dir -and -not (Test-Path $dir)) { New-Item -ItemType Directory -Path $dir -Force | Out-Null }; $in = [System.Console]::OpenStandardInput(); $out = [System.IO.File]::OpenWrite($cloned); $in.CopyTo($out); $out.Close(); Stop-Process -Name gitmap -Force -ErrorAction SilentlyContinue; Start-Sleep -Milliseconds 200; if (Test-Path $dest) { $old = $dest + '.old'; Remove-Item $old -Force -ErrorAction SilentlyContinue; Move-Item -Path $dest -Destination $old -Force -ErrorAction SilentlyContinue }; Move-Item -Path $cloned -Destination $dest -Force"`,
			destPath,
		)
	}

	return fmt.Sprintf(
		`sh -c "dir=$(dirname '%s'); mkdir -p \"$dir\"; cloned='%s-cloned'; cat > \"$cloned\"; chmod +x \"$cloned\"; pkill -f gitmap 2>/dev/null || true; mv -f \"$cloned\" '%s'"`,
		destPath, destPath, destPath,
	)
}
