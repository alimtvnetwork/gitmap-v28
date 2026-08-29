package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/release"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/verbose"
)

// runUpdate handles the "update" subcommand.
//
// As of v5.51.0 the default path downloads + runs the canonical remote
// installer (install.ps1 on Windows, install.sh elsewhere) so updates
// no longer require a source checkout. The installer itself performs
// the parallel `-v<N+i>` sibling-repo probe, so the latest published
// gitmap-vN release always wins.
//
// The legacy source-rebuild handoff is preserved behind --source-rebuild
// for users who maintain a local clone and want to ship in-tree changes.
func runUpdate() error {
	requireOnline()
	if hasFlag(constants.FlagSourceRebuild) == false && runUpdateRemoteInstall() == true {
		return nil
	}
	if hasFlag(constants.FlagSourceRebuild) == false {
		fmt.Fprint(os.Stderr, constants.MsgUpdateRemoteFallback)
	}
	repoPath, err := resolveRepoPath()
	if err != nil {
		return err
	}
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

// resolveRepoPath returns the repo path from --repo-path flag or embedded constant.
// If neither is available, it attempts to delegate to gitmap-updater.
func resolveRepoPath() (string, error) {
	for _, path := range []string{
		resolveRepoPathFromFlag(),
		resolveRepoPathFromEmbedded(),
		resolveRepoPathFromDB(),
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

	// Try to fall back to gitmap-updater for release-based update
	if tryUpdaterFallback() {
		cliexit.HandleError(nil, 0)
		return "", nil
	}

	return "", apperror.NewSimple("no repo path resolved", "E9024")
}

// tryUpdaterFallback looks for gitmap-updater on PATH and launches it.
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
	if errRun != nil && errors.As(errRun, &exitErr) == true {
		appErr := apperror.NewWithDetails(
			"cmd.update.updaterFallback",
			"E1132",
			fmt.Sprintf("gitmap-updater failed with exit code %d", exitErr.ExitCode()),
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

// runUpdateRunner is a hidden command that performs the real update work.
// After the binary update completes, it triggers a best-effort schema
// migration so the next CLI invocation never has to repair the database.
//
// Phase 3 of the handoff chain runs at the end via
// scheduleDeployedCleanupHandoff: the freshly-deployed gitmap binary
// (a different file than this handoff copy) is invoked detached to
// run `update-cleanup`. Only the deployed binary can safely remove
// the still-locked handoff copy and the just-renamed *.exe.old.
// See spec/08-generic-update/06-cleanup.md and
// spec/03-general/02f-self-update-orchestration.md.
func runUpdateRunner() error {
	repoPath, err := resolveRepoPath()
	if err != nil {
		return err
	}
	report := resolveReportErrors()

	currentVersion := constants.Version
	targetVersion := readTargetVersion(repoPath)

	initRunnerVerbose()
	fmt.Printf(constants.MsgUpdateStarting)
	fmt.Printf(constants.MsgUpdateRepoPath, repoPath)
	fmt.Printf(constants.MsgUpdateVersionCompare, currentVersion, targetVersion)
	executeUpdate(repoPath, report)
	runPostUpdateMigrate()
	report.summarize()
	printUpdateSummary(currentVersion, targetVersion, repoPath)
	scheduleDeployedCleanupHandoff()
	return nil
}

// getFlagValue returns the value following a flag like --repo-path <value>.
func getFlagValue(name string) string {
	args := os.Args[2:]
	for i, arg := range args {
		if arg == name && i+1 < len(args) {
			return args[i+1]
		}
	}

	return ""
}

// initRunnerVerbose initializes verbose logging if --verbose flag is present.
func initRunnerVerbose() {
	if hasFlag(constants.FlagVerbose) == false {
		return
	}
	log, err := verbose.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnVerboseLogFailed, err)
		return
	}
	defer log.Close()
	log.Log(constants.UpdateRunnerLogStart, constants.RepoPath)
}

// hasFlag checks if a flag is present in os.Args[2:].
func hasFlag(name string) bool {
	for _, arg := range os.Args[2:] {
		if arg == name {
			return true
		}
	}

	return false
}

// copyFile copies src to dst.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)

	return err
}

// readTargetVersion reads the target version from the repository.
func readTargetVersion(repoPath string) string {
	versionPath := filepath.Join(repoPath, constants.DefaultVersionFile)
	data, err := os.ReadFile(versionPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUpdateVersionRead, versionPath, err)
		return "unknown"
	}

	var parsed struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUpdateVersionRead, versionPath, err)
		return "unknown"
	}
	return parsed.Version
}

// printUpdateSummary outputs the final update summary.
func printUpdateSummary(oldVer, newVer, repoPath string) {
	fmt.Printf(constants.MsgUpdateSummaryDetail, oldVer, newVer, repoPath)
	fmt.Printf("  → Update: %s -> %s\n", oldVer, newVer)
	printUpdateChangelog()
	fmt.Printf("\n\n")
}

// printUpdateChangelog prints the last two notes from the changelog.
func printUpdateChangelog() {
	entries, err := release.ReadChangelog()
	hasEntries := err == nil && len(entries) > 0
	if hasEntries == false {
		return
	}

	for _, note := range getLastTwoNotes(entries[0].Notes) {
		fmt.Printf("    - %s\n", note)
	}
}

func getLastTwoNotes(notes []string) []string {
	count := len(notes)
	if count > 2 {
		return notes[count-2:]
	}
	return notes
}
