package cmdclone

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvscode"
	"github.com/alimtvnetwork/gitmap-v28/cli/errreport"
	"github.com/alimtvnetwork/gitmap-v28/cli/helptext"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/vscodepm"
)

var versionPattern = regexp.MustCompile(`^(v?\d+\.\d+\.\d+([.\-+].+)?|v\+\+|v\+\d+)$`)

// Delegate hooks
var (
	CreatePendingTaskFn            func(typeName, targetPath, workDir, sourceCmd, cmdArgs string) (int64, *store.DB)
	CompletePendingTaskFn          func(db *store.DB, taskID int64)
	FailPendingTaskFn              func(db *store.DB, taskID int64, reason string)
	RequireOnlineFn                func()
	CheckHelpFn                    func(string, []string)
	WriteShellHandoffFn            func(string)
	EscapeCwdIfInsideFn            func(string) (string, error)
	FinalizeErrorReportFn          func(*errreport.Collector, bool)
	RunCodingGuidelinesInstallFn   func(string) error
	CommitCodingGuidelinesFn       func(string, bool, bool) error
	RunCFRPPriorVersionPrivatizeFn func(string, bool) error
	RunGitHubDesktopOptimizeFn     func([]string) error
	ResolveEndpointStringFn        func(string) string
	ResolveReleaseAliasPathFn      func(string) string
	RunStatusFn                    func([]string) error
)

func runStatus(args []string) error {
	if RunStatusFn != nil {
		return RunStatusFn(args)
	}
	return nil
}

func openDB() (*store.DB, error) {
	return store.OpenDefault()
}

func closeTaskDB(db *store.DB) {
	if db != nil {
		_ = db.Close()
	}
}

func buildCommandArgs(args []string) string {
	return strings.Join(args, " ")
}

func createPendingTask(typeName, targetPath, workDir, sourceCmd, cmdArgs string) (int64, *store.DB) {
	if CreatePendingTaskFn != nil {
		return CreatePendingTaskFn(typeName, targetPath, workDir, sourceCmd, cmdArgs)
	}
	return 0, nil
}

func completePendingTask(db *store.DB, taskID int64) {
	if CompletePendingTaskFn != nil {
		CompletePendingTaskFn(db, taskID)
	}
}

func failPendingTask(db *store.DB, taskID int64, reason string) {
	if FailPendingTaskFn != nil {
		FailPendingTaskFn(db, taskID, reason)
	}
}

func requireOnline() {
	if RequireOnlineFn != nil {
		RequireOnlineFn()
	}
}

func checkHelp(command string, args []string) {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			helptext.Print(command)
			cliexit.Exit(0)
		}
	}
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}

func WriteShellHandoff(targetPath string) {
	if WriteShellHandoffFn != nil {
		WriteShellHandoffFn(targetPath)
	}
}

func escapeCwdIfInside(target string) (string, error) {
	if EscapeCwdIfInsideFn != nil {
		return EscapeCwdIfInsideFn(target)
	}
	return "", nil
}

func finalizeErrorReport(c *errreport.Collector, quiet bool) {
	if FinalizeErrorReportFn != nil {
		FinalizeErrorReportFn(c, quiet)
	}
}

type CodingGuidelinesOpts struct {
	WorkingDir string
}

type CGCommitOpts struct {
	WorkingDir   string
	IsSkipCommit bool
	IsSkipPush   bool
}

func RunCodingGuidelinesInstall(opts CodingGuidelinesOpts) error {
	if RunCodingGuidelinesInstallFn != nil {
		return RunCodingGuidelinesInstallFn(opts.WorkingDir)
	}
	return nil
}

func CommitCodingGuidelines(opts CGCommitOpts) error {
	if CommitCodingGuidelinesFn != nil {
		return CommitCodingGuidelinesFn(opts.WorkingDir, opts.IsSkipCommit, opts.IsSkipPush)
	}
	return nil
}

func runCFRPPriorVersionPrivatize(absPath string, autoYes bool) error {
	if RunCFRPPriorVersionPrivatizeFn != nil {
		return RunCFRPPriorVersionPrivatizeFn(absPath, autoYes)
	}
	return nil
}

func runGitHubDesktopOptimize(args []string) error {
	if RunGitHubDesktopOptimizeFn != nil {
		return RunGitHubDesktopOptimizeFn(args)
	}
	return nil
}

func resolveEndpointString(raw string) string {
	if ResolveEndpointStringFn != nil {
		return ResolveEndpointStringFn(raw)
	}
	return raw
}

func resolveReleaseAliasPath(alias string) string {
	if ResolveReleaseAliasPathFn != nil {
		return ResolveReleaseAliasPathFn(alias)
	}
	return ""
}

func syncClonedReposToVSCodePM(pairs []vscodepm.Pair, skip bool) {
	cmdvscode.SyncClonedReposToVSCodePM(pairs, skip)
}

func syncSingleClonedRepoToVSCodePM(absPath, repoName string, skip bool) {
	cmdvscode.SyncSingleClonedRepoToVSCodePM(absPath, repoName, skip)
}

func buildClonePMPair(absPath, repoName string) vscodepm.Pair {
	return cmdvscode.BuildClonePMPair(absPath, repoName)
}

func printVSCodeOptimizeResult(s vscodepm.OptimizeSummary, isDryRun bool) {
	cmdvscode.PrintVSCodeOptimizeResult(s, isDryRun)
}

func argsTail() []string {
	if len(os.Args) > 2 {
		return os.Args[2:]
	}
	return nil
}

func looksLikeVersion(s string) bool {
	return versionPattern.MatchString(s)
}

func extractPositionalArgs(args []string) []string {
	out := make([]string, 0, len(args))
	skipNext := false
	for _, a := range args {
		if skipNext {
			skipNext = false
			continue
		}
		if len(a) > 0 && a[0] == '-' {
			continue
		}
		out = append(out, a)
	}
	return out
}

func extractFlagArgs(args []string) []string {
	out := make([]string, 0, len(args))
	for _, a := range args {
		if len(a) > 0 && a[0] == '-' {
			out = append(out, a)
		}
	}
	return out
}

func reorderFlagsBeforeArgs(args []string) []string {
	flags := make([]string, 0, len(args))
	positional := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
		} else {
			positional = append(positional, arg)
		}
	}

	return append(flags, positional...)
}

func expandTilde(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	if path == "~" {
		return home
	}
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~\\") {
		return filepath.Join(home, path[2:])
	}
	return path
}

func isGitRepo(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, ".git"))
	return err == nil
}

func currentOriginURL(dir string) (string, error) {
	out, err := exec.Command("git", "-C", dir, "config", "--get", "remote.origin.url").Output()
	if err != nil {
		return "", fmt.Errorf("read remote.origin.url in %s: %w", dir, err)
	}

	return strings.TrimSpace(string(out)), nil
}
