// Package cmdssh — ssh_deploy_cmd.go parses deploy CLI commands and flags.
package cmdssh

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// DeployOptions contains runtime arguments for smart deployment.
type DeployOptions struct {
	SubCmd        string
	TargetToken   string
	Conn          db.SSHConnection
	Host          store.SSHHost
	SourcePath    string
	DestPath      string
	RemoteWorkDir string
	IsOverwrite   bool
	IsSkip        bool
	IsSync        bool
	IsSyncRight   bool
	IsSyncLeft    bool
	IsJSON        bool
	IsInteractive bool
	Parallel      int
	IsDryRun      bool
}

type deployFlagState struct {
	isOverwrite bool
	isSkip      bool
	isSync      bool
	isSyncRight bool
	isSyncLeft  bool
	isJSON      bool
	isDryRun    bool
	parallel    int
}

// RunSSHDeployCLI parses and executes gitmap deploy / deploy-right / deploy-left.
func RunSSHDeployCLI(subCmd string, args []string) error {
	if isLegacyDeploySub(subCmd, args) {
		return RunSSHDeployRouterCLI(args)
	}
	if isHelpDeployRequest(args) {
		printDeployHelpText(subCmd)
		return nil
	}
	return runDeployExecution(subCmd, args)
}

func isLegacyDeploySub(subCmd string, args []string) bool {
	if subCmd != "deploy" || len(args) == 0 {
		return false
	}
	return isKnownLegacyDeployWord(strings.ToLower(args[0]))
}

func isKnownLegacyDeployWord(first string) bool {
	switch first {
	case "config", "config-ssh", "ssh-config", "ssh", "keys", "key", "k", "bin", "binary", "exe", "gitmap", "node-config", "nodeconfig", "nc", "nodes", "macro", "macros":
		return true
	default:
		return false
	}
}

func isHelpDeployRequest(args []string) bool {
	if len(args) == 0 {
		return true
	}
	return isHelpFlag(args[0])
}

func runDeployExecution(subCmd string, args []string) error {
	opts, err := parseDeployOptions(subCmd, args)
	if err != nil {
		return err
	}
	return ExecuteDeploy(opts)
}

func parseDeployOptions(subCmd string, args []string) (DeployOptions, error) {
	flags, posArgs := parseDeployFlags(args)
	if len(posArgs) < 3 {
		return DeployOptions{}, apperror.NewValidationError("missing arguments: expected <target> <source> <destination>")
	}
	targetToken, srcPath, destPath := posArgs[0], posArgs[1], posArgs[2]
	return buildDeployOptions(subCmd, targetToken, srcPath, destPath, flags)
}

func parseDeployFlags(args []string) (deployFlagState, []string) {
	state := deployFlagState{parallel: 4}
	var posArgs []string
	for i := 0; i < len(args); i++ {
		i = consumeDeployArg(args, i, &state, &posArgs)
	}
	return state, posArgs
}

func consumeDeployArg(args []string, i int, state *deployFlagState, posArgs *[]string) int {
	arg := args[i]
	if !strings.HasPrefix(arg, "-") {
		*posArgs = append(*posArgs, arg)
		return i
	}
	return applyDeployFlag(args, i, state)
}

func applyDeployFlag(args []string, i int, state *deployFlagState) int {
	arg := args[i]
	if tryApplyBooleanDeployFlag(arg, state) {
		return i
	}
	return applyValueDeployFlag(args, i, state)
}

func tryApplyBooleanDeployFlag(arg string, state *deployFlagState) bool {
	switch arg {
	case "--overwrite", "-o":
		state.isOverwrite = true
		return true
	case "--skip", "-s":
		state.isSkip = true
		return true
	case "--sync":
		state.isSync = true
		return true
	case "--sync-right":
		state.isSyncRight = true
		return true
	case "--sync-left":
		state.isSyncLeft = true
		return true
	case "--json", "-j":
		state.isJSON = true
		return true
	case "--dry-run", "-n":
		state.isDryRun = true
		return true
	default:
		return false
	}
}

func applyValueDeployFlag(args []string, i int, state *deployFlagState) int {
	arg := args[i]
	if val, ok := extractInlineFlagValue(arg, "--parallel="); ok {
		state.parallel = parseParallelValue(val)
		return i
	}
	if val, ok := extractInlineFlagValue(arg, "-p="); ok {
		state.parallel = parseParallelValue(val)
		return i
	}
	return applyNextValueFlag(args, i, state)
}

func applyNextValueFlag(args []string, i int, state *deployFlagState) int {
	arg := args[i]
	hasValueAhead := (arg == "--parallel" || arg == "-p") && i+1 < len(args)
	if hasValueAhead {
		state.parallel = parseParallelValue(args[i+1])
		return i + 1
	}
	return i
}

func parseParallelValue(raw string) int {
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return 4
	}
	if n > 16 {
		return 16
	}
	return n
}

func extractInlineFlagValue(arg, prefix string) (string, bool) {
	if strings.HasPrefix(arg, prefix) {
		return strings.TrimPrefix(arg, prefix), true
	}
	return "", false
}

func buildDeployOptions(subCmd, targetToken, srcPath, destPath string, state deployFlagState) (DeployOptions, error) {
	conn, host, err := resolveDeployTargetInfo(targetToken)
	if err != nil {
		return DeployOptions{}, err
	}
	localPath, err := resolveLocalDeploySource(srcPath)
	if err != nil {
		return DeployOptions{}, err
	}
	remotePath, remoteWorkDir := resolveRemoteDeployDest(destPath, conn.OS)
	state = applySubcommandModeDefaults(subCmd, state)
	return makeDeployOptions(subCmd, targetToken, conn, host, localPath, remotePath, remoteWorkDir, state), nil
}

func resolveDeployTargetInfo(targetToken string) (db.SSHConnection, store.SSHHost, error) {
	nodes, err := resolveDeployNodes(targetToken)
	if err != nil {
		return db.SSHConnection{}, store.SSHHost{}, err
	}
	if len(nodes) == 0 {
		return db.SSHConnection{}, store.SSHHost{}, apperror.NewNotFound("target_node", "E404", "node not found: "+targetToken)
	}
	conn := nodes[0]
	host := findDeployTargetHost(conn)
	return conn, host, nil
}

func findDeployTargetHost(conn db.SSHConnection) store.SSHHost {
	_, hosts, err := loadFleetInventory()
	if err != nil || len(hosts) == 0 {
		return fallbackTargetHost(conn)
	}
	return matchHostFromInventory(hosts, conn)
}

func matchHostFromInventory(hosts []store.SSHHost, conn db.SSHConnection) store.SSHHost {
	for _, h := range hosts {
		if strings.EqualFold(h.Alias, conn.Alias) || (conn.IPAddress != "" && strings.EqualFold(h.IP, conn.IPAddress)) {
			return h
		}
	}
	return fallbackTargetHost(conn)
}

func fallbackTargetHost(conn db.SSHConnection) store.SSHHost {
	return store.SSHHost{
		ID:    conn.Alias,
		Alias: conn.Alias,
		IP:    conn.IPAddress,
	}
}

func resolveLocalDeploySource(srcPath string) (string, error) {
	if filepath.IsAbs(srcPath) {
		return filepath.Clean(srcPath), nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return filepath.Clean(srcPath), nil
	}
	return filepath.Clean(filepath.Join(cwd, srcPath)), nil
}

func resolveRemoteDeployDest(destPath, osType string) (string, string) {
	isWin := strings.EqualFold(osType, "windows")
	workDir := defaultRemoteWorkDir(isWin)
	if isRemotePathAbsolute(destPath, isWin) {
		return normalizeRemoteSlash(destPath, isWin), workDir
	}
	return combineRemotePath(workDir, destPath, isWin), workDir
}

func defaultRemoteWorkDir(isWindows bool) string {
	if isWindows {
		return "D:/work"
	}
	return "~"
}

func isRemotePathAbsolute(path string, isWindows bool) bool {
	if isWindows {
		return isWindowsAbsolute(path)
	}
	return strings.HasPrefix(path, "/") || strings.HasPrefix(path, "~")
}

func isWindowsAbsolute(path string) bool {
	hasDrive := len(path) >= 2 && path[1] == ':'
	return hasDrive || strings.HasPrefix(path, "/") || strings.HasPrefix(path, "\\")
}

func combineRemotePath(workDir, destPath string, isWindows bool) string {
	combined := strings.TrimRight(workDir, "/\\") + "/" + strings.TrimLeft(destPath, "/\\")
	return normalizeRemoteSlash(combined, isWindows)
}

func normalizeRemoteSlash(path string, isWindows bool) string {
	if isWindows {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(path)
}

func applySubcommandModeDefaults(subCmd string, s deployFlagState) deployFlagState {
	if isAnyModeFlagSet(s) {
		return s
	}
	if subCmd == "deploy-right" {
		s.isSyncRight = true
		return s
	}
	if subCmd == "deploy-left" {
		s.isSyncLeft = true
		return s
	}
	return s
}

func isAnyModeFlagSet(s deployFlagState) bool {
	return s.isOverwrite || s.isSkip || s.isSync || s.isSyncRight || s.isSyncLeft
}

func makeDeployOptions(subCmd, targetToken string, conn db.SSHConnection, host store.SSHHost, localPath, remotePath, workDir string, s deployFlagState) DeployOptions {
	return DeployOptions{
		SubCmd:        subCmd,
		TargetToken:   targetToken,
		Conn:          conn,
		Host:          host,
		SourcePath:    localPath,
		DestPath:      remotePath,
		RemoteWorkDir: workDir,
		IsOverwrite:   s.isOverwrite,
		IsSkip:        s.isSkip,
		IsSync:        s.isSync,
		IsSyncRight:   s.isSyncRight,
		IsSyncLeft:    s.isSyncLeft,
		IsJSON:        s.isJSON,
		IsInteractive: !s.isJSON,
		Parallel:      s.parallel,
		IsDryRun:      s.isDryRun,
	}
}

func printDeployHelpText(subCmd string) {
	RenderDeployRichHelp()
}

// MakeDeployCobraCmd creates a cobra.Command for deploy completion and routing.
func MakeDeployCobraCmd(use string) *cobra.Command {
	cmd := &cobra.Command{
		Use:               use + " [target-node] [source-path] [dest-path]",
		Short:             "Smart file and folder deployment across SSH fleet with sync modes",
		ValidArgsFunction: DeployValidArgsFunction,
		RunE: func(c *cobra.Command, args []string) error {
			return RunSSHDeployCLI(use, args)
		},
	}
	cmd.Flags().BoolP("overwrite", "o", false, "Force overwrite remote files")
	cmd.Flags().BoolP("skip", "s", false, "Skip already existing files")
	cmd.Flags().Bool("sync", false, "Bidirectional sync based on modification time")
	cmd.Flags().Bool("sync-right", false, "Push newer files from local to remote only")
	cmd.Flags().Bool("sync-left", false, "Pull newer files from remote to local only")
	cmd.Flags().BoolP("json", "j", false, "Output deployment telemetry as JSON")
	cmd.Flags().BoolP("dry-run", "n", false, "Preview transfer plan without executing")
	cmd.Flags().IntP("parallel", "p", 4, "Number of concurrent transfer workers")
	return cmd
}
