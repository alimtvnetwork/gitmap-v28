// Package cmdssh — ssh_deploy_bin.go automates binary deployment across remote fleet nodes.
package cmdssh

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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
		if (arg == "--file" || arg == "-f") && i+1 < len(args) {
			opts.filePath = args[i+1]
			i++
		} else if arg == "--dry-run" || arg == "-n" {
			opts.isDryRun = true
		} else if (arg == "--timeout" || arg == "-t") && i+1 < len(args) {
			opts.timeout = parseTimeoutFlag(args, i, opts.timeout)
			i++
		} else if !strings.HasPrefix(arg, "-") {
			positional = append(positional, arg)
		}
	}

	if len(positional) > 0 {
		opts.target = positional[0]
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
			`d:\work\gitmap\cli\gitmap.exe`,
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

func loadAllSSHConnections() ([]db.SSHConnection, error) {
	dbConn, err := store.OpenDefault()
	if err != nil {
		return nil, apperror.WrapSimple(err, "open connection database")
	}
	defer dbConn.Close()

	res := db.GetSSHConnections(dbConn.Context(), dbConn.SQL())
	if res.IsFailure() {
		return nil, apperror.NewExecutionError(fmt.Sprintf("load ssh nodes: %v", res.Err))
	}

	return res.Data, nil
}

func resolveDeployNodes(target string) ([]db.SSHConnection, error) {
	connections, err := loadAllSSHConnections()
	if err != nil {
		return nil, err
	}

	if target == "" || target == "all" || target == "--all" {
		return connections, nil
	}

	var matched []db.SSHConnection
	for _, c := range connections {
		if strings.EqualFold(c.Alias, target) || strings.EqualFold(c.IPAddress, target) {
			matched = append(matched, c)
		}
	}

	if len(matched) == 0 {
		return nil, apperror.NewNotFound("target_node", "E404", fmt.Sprintf("target node %q not found in fleet", target))
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

func streamAndVerify(client *ssh.Client, node db.SSHConnection, binPath, destPath string) bool {
	t0 := time.Now()
	if err := streamBinaryToRemote(client, destPath, binPath, isWindowsOS(node.OS)); err != nil {
		fmt.Printf("    %s✖ Upload error: %v%s\n", constants.ColorRed, err, constants.ColorReset)

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

	return session.Run(remoteCmd)
}

func buildStreamReceiverCmd(destPath string, isWin bool) string {
	if isWin {
		return fmt.Sprintf(
			`powershell -NoProfile -Command "$dest = '%s'; $tmp = $dest + '.tmp'; $dir = [System.IO.Path]::GetDirectoryName($dest); if ($dir -and !(Test-Path $dir)) { New-Item -ItemType Directory -Path $dir -Force | Out-Null }; $in = [System.Console]::OpenStandardInput(); $out = [System.IO.File]::OpenWrite($tmp); $in.CopyTo($out); $out.Close(); Move-Item -Path $tmp -Destination $dest -Force"`,
			destPath,
		)
	}

	return fmt.Sprintf(
		`sh -c "mkdir -p '$(dirname '%s')' && cat > '%s.tmp' && chmod +x '%s.tmp' && mv -f '%s.tmp' '%s'"`,
		destPath, destPath, destPath, destPath, destPath,
	)
}
