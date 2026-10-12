package cmdupdate

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// runUpdate handles the "update" subcommand.
//
// In accordance with 02-spec/14-update/02-self-update-overview.md and
// 02-spec/14-update/08-console-safe-handoff.md:
//  1. If a local source repository is resolved, build from repo using the
//     console-safe handoff mechanism (blocking cmd.Run(), zero popup windows,
//     same terminal session).
//  2. Otherwise fall back to the binary remote installer.
func runUpdate() error {
	requireOnline()

	// Priority 1: Source-based update via console-safe handoff
	repoPath, errRepo := resolveRepoPath()
	if errRepo == nil && len(repoPath) > 0 {
		return runSourceHandoffUpdate(repoPath)
	}

	// Priority 2: Binary-based remote installer update
	if runUpdateRemoteInstall() {
		return nil
	}

	if hasFlag(constants.FlagSourceRebuild) || isRepoPathPromptAllowed() {
		if prompted := promptRepoPath(); len(prompted) > 0 {
			saveRepoPathToDB(prompted)

			return runSourceHandoffUpdate(prompted)
		}
	}

	if errRepo != nil {
		return errRepo
	}

	return nil
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
// If neither is available, it attempts to delegate to cli-updater.
func resolveRepoPath() (string, error) {
	for _, path := range []string{
		resolveRepoPathFromFlag(),
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

	if hasFlag(constants.FlagSourceRebuild) {
		if prompted := promptRepoPath(); len(prompted) > 0 {
			saveRepoPathToDB(prompted)

			return prompted, nil
		}
	}

	// Try to fall back to cli-updater for release-based update
	if tryUpdaterFallback() {
		cliexit.HandleError(nil, 0)

		return "", nil
	}

	return "", apperror.NewSimple("no repo path resolved", "E9024")
}

// tryUpdaterFallback looks for cli-updater on PATH and launches it.
func tryUpdaterFallback() bool {
	updaterPath, err := exec.LookPath(constants.UpdaterBin)
	if err != nil {
		return false
	}

	fmt.Printf(constants.MsgUpdaterFallback, updaterPath)
	cmd := exec.Command(updaterPath, "run")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	errRun := cmd.Run()
	var exitErr *exec.ExitError
	if errRun != nil && errors.As(errRun, &exitErr) {
		appErr := apperror.NewWithDetails(
			"cmd.update.updaterFallback",
			"E1132",
			fmt.Sprintf("cli-updater failed with exit code %d", exitErr.ExitCode()),
			"cmd.update",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(appErr, exitErr.ExitCode())
	}

	if errRun != nil {
		return false
	}

	return true
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
