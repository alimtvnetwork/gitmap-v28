package cmdupdate

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

const FlagFromSource = "--from-source"

func isSourceRebuildRequested() bool {
	for _, arg := range os.Args[1:] {
		if slices.Contains([]string{constants.FlagSourceRebuild, FlagFromSource}, arg) {
			return true
		}
	}

	return false
}

// runUpdate handles the "update" subcommand.
//
// By default, gitmap update performs a direct pre-compiled binary update
// from GitHub Releases without requiring Go, Git, or source checkouts.
// Source-based rebuild via local repo is ONLY executed if --source-rebuild
// or --from-source is explicitly supplied.
func runUpdate() error {
	requireOnline()

	opts := ParseUpdateOptions(os.Args[1:])

	// Priority 1 (Explicit Opt-In Only): Source rebuild via console-safe handoff
	if isSourceRebuildRequested() {
		return runSourceRebuildUpdate()
	}

	// Priority 2 (Default): Pre-compiled binary update from GitHub Releases
	return runPureBinaryUpdate(opts)
}

// runSourceRebuildUpdate performs the explicit opt-in source rebuild flow.
func runSourceRebuildUpdate() error {
	repoPath, errRepo := resolveRepoPath()
	if errRepo == nil && len(repoPath) > 0 {
		return runSourceHandoffUpdate(repoPath)
	}

	if prompted := promptSourceRepoPath(); len(prompted) > 0 {
		saveRepoPathToDB(prompted)

		return runSourceHandoffUpdate(prompted)
	}

	if errRepo != nil {
		return errRepo
	}

	return apperror.NewSimple("no repo path resolved for source rebuild", "E9024")
}

// promptSourceRepoPath prompts for a repo path only when prompting is allowed.
func promptSourceRepoPath() string {
	if !isRepoPathPromptAllowed() {
		return ""
	}

	return promptRepoPath()
}

func runSourceHandoffUpdate(repoPath string) error {
	report := resolveReportErrors()
	report.announce()

	selfPath, err := os.Executable()
	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.update.executable",
			"E1131",
			"failed to determine current executable path",
			"cmd.update",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(appErr, 1)

		return nil
	}

	copyPath := createHandoffCopy(selfPath)
	fmt.Printf(constants.MsgUpdateActive, selfPath, copyPath)
	launchHandoff(copyPath, repoPath, report)

	return nil
}

// resolveRepoPath returns the repo path from flags, embedded constant, DB, cwd, or registered repos.
// It is strictly guarded so it only searches if explicit source rebuild or --repo-path is requested.
func resolveRepoPath() (string, error) {
	if flagPath := resolveRepoPathFromFlag(); len(flagPath) > 0 {
		return flagPath, nil
	}

	if !isSourceRebuildRequested() {
		return "", apperror.NewSimple("source repo resolution only allowed with --source-rebuild or --from-source", "E9025")
	}

	for _, path := range []string{
		resolveRepoPathFromEmbedded(),
		resolveRepoPathFromDB(),
		resolveRepoPathFromCwd(),
		resolveRepoPathFromRepoTable(),
		resolveRepoPathFromWorkDirs(),
	} {
		if len(path) > 0 {
			saveRepoPathToDB(path)

			return path, nil
		}
	}

	if prompted := promptRepoPath(); len(prompted) > 0 {
		saveRepoPathToDB(prompted)

		return prompted, nil
	}

	return "", apperror.NewSimple("no repo path resolved", "E9024")
}

// createHandoffCopy creates a temporary copy of the binary for handoff.
func createHandoffCopy(selfPath string) string {
	nameFmt := constants.UpdateCopyFmtExe
	if runtime.GOOS != "windows" {
		nameFmt = constants.UpdateCopyFmtUnix
	}

	name := fmt.Sprintf(nameFmt, os.Getpid())
	copyPath := filepath.Join(filepath.Dir(selfPath), name)

	if copyFile(selfPath, copyPath) == nil {
		makeExecutable(copyPath)

		return copyPath
	}

	fallbackPath := filepath.Join(os.TempDir(), name)
	if err := copyFile(selfPath, fallbackPath); err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.update.copyFallback",
			"E1133",
			"failed to create fallback handoff copy",
			"cmd.update",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			map[string]any{"fallbackPath": fallbackPath},
		)
		cliexit.HandleError(appErr, 1)
	}

	makeExecutable(fallbackPath)

	return fallbackPath
}

// makeExecutable sets executable permission on Unix systems.
func makeExecutable(path string) {
	if runtime.GOOS == "windows" {
		return
	}

	if err := os.Chmod(path, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not make %s executable: %v\n", path, err)
	}
}

// launchHandoff runs the handoff binary with update-runner command.
func launchHandoff(copyPath, repoPath string, report reportErrorsConfig) {
	copyArgs := []string{constants.CmdUpdateRunner}
	if hasFlag(constants.FlagVerbose) {
		copyArgs = append(copyArgs, constants.FlagVerbose)
	}

	if isDebugWindowsRequested() {
		copyArgs = append(copyArgs, constants.FlagDebugWindows)
	}

	copyArgs = append(copyArgs, constants.FlagRepoPath, repoPath)
	copyArgs = report.applyToHandoffArgs(copyArgs)

	cmd := exec.Command(copyPath, copyArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	env := report.applyToEnv(os.Environ())
	if isDebugWindowsRequested() {
		env = append(env, constants.EnvDebugWindows+"=1")
	}

	cmd.Env = env
	dumpDebugWindowsHeader("phase-2 handoff (active → copy)")
	dumpDebugWindowsHandoff("phase-2-copy", copyPath,
		append([]string{copyPath}, copyArgs...))
	dumpDebugWindowsFooter()
	if err := cmd.Run(); err != nil {
		handleHandoffError(err)
	}
}

// handleHandoffError exits with the handoff process exit code if available.
func handleHandoffError(err error) {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		appErr := apperror.NewWithDetails(
			"cmd.update.handoffExit",
			"E1134",
			fmt.Sprintf("handoff process exited with code %d", exitErr.ExitCode()),
			"cmd.update",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(appErr, exitErr.ExitCode())
	}

	appErr := apperror.WrapWithDetails(
		err,
		"cmd.update.handoffRun",
		"E1135",
		"failed to run handoff update process",
		"cmd.update",
		apperror.ErrorTypeExecution,
		apperror.SeverityError,
		nil,
	)
	cliexit.HandleError(appErr, 1)
}
