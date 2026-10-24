// Package cmdpending — pending_commits_cmd.go inspects and classifies uncommitted changes and unpushed commits.
package cmdpending

import (
	"database/sql"
	"fmt"
	"os/exec"
	"slices"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
)

// gitCmdExecutor defines a function that executes a git command in a specific directory.
type gitCmdExecutor func(dir string, args ...string) (string, error)

var defaultPendingCommitsGitExecutor gitCmdExecutor = func(dir string, args ...string) (string, error) {
	cmdArgs := append([]string{"-C", dir}, args...)
	cmd := exec.Command("git", cmdArgs...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

var currentPendingCommitsGitExecutor = defaultPendingCommitsGitExecutor

// pendingCommitsSSHExecutor defines an interface for executing commands over SSH.
type pendingCommitsSSHExecutor interface {
	Execute(conn db.SSHConnection, cmd, shell string) (string, error)
}

type defaultPendingCommitsSSHRunner struct{}

func (r defaultPendingCommitsSSHRunner) Execute(conn db.SSHConnection, cmd, shell string) (string, error) {
	client, errConnect := cmdssh.ConnectSSHClientWithErr(conn)
	if errConnect != nil {
		return "", errConnect
	}
	defer client.Close()

	out, errExec := secrets.RunCommand(client, cmd, shell)
	if errExec != nil && shell == "bash" {
		out, errExec = secrets.RunCommand(client, cmd, "sh")
	}
	return out, errExec
}

var currentPendingCommitsSSHExecutor pendingCommitsSSHExecutor = defaultPendingCommitsSSHRunner{}

// RunPendingCommits executes discovery of uncommitted and unpushed changes.
func RunPendingCommits(args []string) error {
	if isPendingCommitsHelpRequested(args) {
		printPendingCommitsHelp()
		return nil
	}

	opts := parsePendingCommitsOptions(args)
	if opts.IsSSH {
		return runFleetSSHPendingCommits(opts)
	}

	return runLocalPendingCommits(opts)
}

func isPendingCommitsHelpRequested(args []string) bool {
	for _, a := range args {
		if a == "-h" || a == "--help" || a == "help" {
			return true
		}
	}
	return false
}

func parsePendingCommitsOptions(args []string) PendingCommitsOptions {
	opts := defaultPendingCommitsOptions()
	for i := 0; i < len(args); i++ {
		i = processPendingCommitArg(args, i, &opts)
	}
	normalizePendingCommitsOptions(&opts)
	return opts
}

func defaultPendingCommitsOptions() PendingCommitsOptions {
	return PendingCommitsOptions{
		SortMode:    "priority",
		DetailMode:  "summary",
		IsDirtyOnly: true,
	}
}

func normalizePendingCommitsOptions(opts *PendingCommitsOptions) {
	if opts.SortMode != "name" && opts.SortMode != "count" && opts.SortMode != "priority" {
		opts.SortMode = "priority"
	}
	if opts.DetailMode == "1" {
		opts.DetailMode = "1-by-1"
	}
}

func processPendingCommitArg(args []string, i int, opts *PendingCommitsOptions) int {
	arg := args[i]
	switch {
	case slices.Contains([]string{"-s", "--sort"}, arg) || strings.HasPrefix(arg, "--sort="):
		return parseSortFlag(args, i, opts)
	case slices.Contains([]string{"-d", "--detail"}, arg) || strings.HasPrefix(arg, "--detail="):
		return parseDetailFlag(args, i, opts)
	default:
		applyPendingBooleanOrTarget(arg, opts)
		return i
	}
}

func parseSortFlag(args []string, i int, opts *PendingCommitsOptions) int {
	arg := args[i]
	if strings.HasPrefix(arg, "--sort=") {
		opts.SortMode = strings.ToLower(strings.TrimPrefix(arg, "--sort="))
		return i
	}
	if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
		opts.SortMode = strings.ToLower(args[i+1])
		return i + 1
	}
	return i
}

func parseDetailFlag(args []string, i int, opts *PendingCommitsOptions) int {
	arg := args[i]
	if strings.HasPrefix(arg, "--detail=") {
		opts.DetailMode = strings.ToLower(strings.TrimPrefix(arg, "--detail="))
		return i
	}
	if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
		opts.DetailMode = strings.ToLower(args[i+1])
		return i + 1
	}
	opts.DetailMode = "all"
	return i
}

func applyPendingBooleanOrTarget(arg string, opts *PendingCommitsOptions) {
	if applyPendingScopeFlags(arg, opts) {
		return
	}
	applyPendingCacheAndTarget(arg, opts)
}

func applyPendingScopeFlags(arg string, opts *PendingCommitsOptions) bool {
	switch {
	case slices.Contains([]string{"-j", "--json"}, arg):
		opts.IsJSON = true
		return true
	case slices.Contains([]string{"--ssh", "-ssh"}, arg):
		opts.IsSSH = true
		return true
	case arg == "--dirty-only":
		opts.IsDirtyOnly, opts.IsAll = true, false
		return true
	case arg == "--all":
		opts.IsAll, opts.IsDirtyOnly = true, false
		return true
	}
	return false
}

func applyPendingCacheAndTarget(arg string, opts *PendingCommitsOptions) bool {
	switch {
	case arg == "--no-cache":
		opts.IsNoCache = true
		return true
	case arg == "--refresh":
		opts.IsRefresh = true
		return true
	case slices.Contains([]string{"--backup", "--serve-backup"}, arg):
		opts.IsBackup = true
		return true
	case !strings.HasPrefix(arg, "-"):
		opts.TargetRepo = arg
		return true
	}
	return false
}

func openOptionalCacheConn(isNoCache bool) *sql.DB {
	if isNoCache {
		return nil
	}
	conn, _ := OpenPendingCommitsCache("")
	return conn
}

func runLocalPendingCommits(opts PendingCommitsOptions) error {
	cacheConn := openOptionalCacheConn(opts.IsNoCache)
	if cacheConn != nil {
		defer cacheConn.Close()
	}

	backupConn, errBackup := OpenPendingCommitsBackup("")
	if errBackup == nil && backupConn != nil {
		defer backupConn.Close()
	}

	if opts.IsBackup {
		return handleBackupMode(backupConn, errBackup, opts)
	}

	inspected, errInspect := resolveAndInspectLocalReposWithDB(opts.TargetRepo, cacheConn, backupConn, opts)
	if errInspect != nil {
		return errInspect
	}
	payload := buildPendingCommitsPayload(inspected, opts)
	if opts.IsJSON {
		return emitPendingCommitsJSON(payload)
	}
	renderLocalPendingCommitsTerminal(payload, len(inspected), inspected, opts)
	return nil
}

func printPendingCommitsHelp() {
	cyan := constants.ColorCyan
	reset := constants.ColorReset
	bold := constants.ColorBold
	yellow := constants.ColorYellow

	fmt.Printf(`
%s┌──────────────────────────────────────────────────────────────────────────────┐
│                      GITMAP PENDING COMMITS HELP                             │
└──────────────────────────────────────────────────────────────────────────────┘%s

%sDESCRIPTION:%s
  Scans local workspace repositories or fleet cluster nodes over SSH to identify
  uncommitted working tree modifications and unpushed commits ahead of origin.

%sUSAGE:%s
  gitmap pending-commits [repoName|all] [flags]
  gitmap pc [repoName|all] [flags]

%sFLAGS:%s
  -s, --sort <mode>       Sort order: priority | name | count (default: priority)
  -d, --detail <mode>     Detail view: summary | all | 1 | none (default: summary)
      --ssh, -ssh         Aggregate pending commits across all registered cluster SSH nodes
  -j, --json              Output machine-readable JSON telemetry
      --dirty-only        Display only repositories with changes (default: true)
      --all               Include completely clean repositories in output
      --no-cache          Bypass status cache and perform live inspection
      --refresh           Force refresh cached pending commit status
      --backup            Serve status from secondary persistent backup DB
      --serve-backup      Serve status from secondary persistent backup DB (alias)
  -h, --help              Display this help menu

%sEXAMPLES:%s
  %sgitmap pc%s
      Summary overview of all dirty repositories in workspace.

  %sgitmap pc --sort=count --detail=all%s
      Detailed file diff breakdown sorted by total modification count.

  %sgitmap pc -ssh --json%s
      Aggregate pending commits across entire SSH fleet as JSON.
`,
		cyan, reset,
		bold, reset,
		bold, reset,
		bold, reset,
		bold, reset,
		yellow, reset,
		yellow, reset,
		yellow, reset,
	)
}
