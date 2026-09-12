package cmdclone

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/clonenext"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/desktop"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/lockcheck"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/verbose"
)

// runCloneNext handles the "clone-next" subcommand.
func runCloneNext(args []string) error {
	if tryFolderArgCloneNext(args) || tryCrossDirCloneNext(args) {
		return nil
	}

	checkHelp("clone-next", args)
	cnFlags := parseCloneNextFlags(args)
	configureCloneNextFlags(cnFlags)

	log, errVerbose := initVerboseSafe(cnFlags.Verbose)
	if errVerbose != nil {
		fmt.Fprintf(os.Stderr, constants.WarnVerboseLogFailed, errVerbose)
	}

	if log != nil {
		defer log.Close()
	}

	return dispatchCloneNext(cnFlags)
}

func configureCloneNextFlags(cnFlags CloneNextFlags) {
	setCmdFaithfulVerify(cnFlags.VerifyCmdFaithful)
	setCmdFaithfulExitOnMismatch(cnFlags.VerifyCmdFaithfulExitOnMismatch)
	setCmdPrintArgv(cnFlags.PrintCloneArgv)
}

func dispatchCloneNext(cnFlags CloneNextFlags) error {
	isBatch := isBatchRunnable(cnFlags, currentWorkingDir())
	if isBatch {
		return handleCloneNextBatch(cnFlags)
	}

	if len(cnFlags.VersionArg) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrCloneNextUsage)

		return apperror.NewValidationError(constants.ErrCloneNextUsage)
	}

	requireOnline()
	applySSHKey(cnFlags.SSHKeyName)

	return executeSingleCloneNext(cnFlags)
}

func handleCloneNextBatch(cnFlags CloneNextFlags) error {
	if cnFlags.DryRun {
		previewDryRunBatch(cnFlags.CSVPath, cnFlags.All)

		return nil
	}

	runCloneNextBatch(cnFlags.CSVPath, cnFlags.All, cnFlags.MaxConcurrency, cnFlags.NoProgress, cnFlags.ReportErrors)
	maybeExitOnCmdFaithfulMismatch()

	return nil
}

func executeSingleCloneNext(cnFlags CloneNextFlags) error {
	cwd, cwdErr := resolveCurrentRepoCwd()
	if cwdErr != nil {
		return cwdErr
	}

	remoteURL, err := gitutil.RemoteURL(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneNextNoRemote, err)

		return apperror.WrapSimple(err, fmt.Sprintf(constants.ErrCloneNextNoRemote, err))
	}

	return executeCloneNextPipeline(cnFlags, cwd, remoteURL)
}

func resolveCurrentRepoCwd() (string, *apperror.AppError) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneNextCwd, err)

		return "", apperror.WrapSimple(err, fmt.Sprintf(constants.ErrCloneNextCwd, err))
	}

	return escapeToRepoRoot(cwd)
}

type cloneNextTargetParams struct {
	repoName        string
	currentFolder   string
	flattenedFolder string
	targetName      string
	targetURL       string
	targetPath      string
	currentVersion  int
	targetVersion   int
}

func resolveCloneNextParams(cnFlags CloneNextFlags, cwd, remoteURL string) (cloneNextTargetParams, *apperror.AppError) {
	currentFolder := filepath.Base(cwd)
	parentDir := filepath.Dir(cwd)
	repoName := extractRepoName(remoteURL)
	parsed := clonenext.ParseRepoName(repoName)
	targetVersion, err := clonenext.ResolveTarget(parsed, cnFlags.VersionArg)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneNextBadVersion, err)

		return cloneNextTargetParams{}, apperror.WrapSimple(err, fmt.Sprintf(constants.ErrCloneNextBadVersion, err))
	}

	targetName := clonenext.TargetRepoName(parsed.BaseName, targetVersion)
	targetURL := clonenext.ReplaceRepoInURL(remoteURL, repoName, targetName)
	flattenedFolder := parsed.BaseName
	targetPath := filepath.Join(parentDir, flattenedFolder)

	return cloneNextTargetParams{
		repoName:        repoName,
		currentFolder:   currentFolder,
		flattenedFolder: flattenedFolder,
		targetName:      targetName,
		targetURL:       targetURL,
		targetPath:      targetPath,
		currentVersion:  parsed.CurrentVersion,
		targetVersion:   targetVersion,
	}, nil
}

func executeCloneNextPipeline(cnFlags CloneNextFlags, cwd, remoteURL string) error {
	p, err := resolveCloneNextParams(cnFlags, cwd, remoteURL)
	if err != nil {
		return err
	}

	maybePrintCloneNextTermBlock(cnFlags, p.targetName, currentBranch(cwd), remoteURL, p.targetURL, p.targetPath)
	if handleCloneNextDryRun(cnFlags.DryRun, p.targetURL, p.targetPath) {
		return nil
	}

	return executeCloneNextSteps(cnFlags, cwd, remoteURL, p)
}

func executeCloneNextSteps(cnFlags CloneNextFlags, cwd, remoteURL string, p cloneNextTargetParams) error {
	if cnFlags.Force {
		fmt.Printf(constants.MsgCNStagePrepare, p.currentFolder, p.flattenedFolder)
	}

	if escErr := prepareCloneNextTarget(p.targetPath, p.flattenedFolder); escErr != nil {
		return escErr
	}

	return performCloneNextOperation(cnFlags, cwd, remoteURL, p)
}

func prepareCloneNextTarget(targetPath, flattenedFolder string) *apperror.AppError {
	if _, escErr := escapeCwdIfInside(targetPath); escErr != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneNextForceFailed, flattenedFolder, escErr, flattenedFolder)

		return apperror.WrapSimple(escErr, fmt.Sprintf(constants.ErrCloneNextForceFailed, flattenedFolder, escErr, flattenedFolder))
	}

	return removeExistingTargetFolder(targetPath, flattenedFolder)
}

func performCloneNextOperation(cnFlags CloneNextFlags, cwd, remoteURL string, p cloneNextTargetParams) error {
	remoteErr := handleCreateRemote(CreateRemoteParams{
		IsCreateRemote: cnFlags.CreateRemote,
		RemoteURL:      remoteURL,
		TargetName:     p.targetName,
	})
	if remoteErr != nil {
		return remoteErr
	}

	return executeCloneAndFinalize(cnFlags, cwd, p)
}

func executeCloneAndFinalize(cnFlags CloneNextFlags, cwd string, p cloneNextTargetParams) error {
	if cnFlags.Force {
		fmt.Printf(constants.MsgCNStageClone, p.targetName)
	}

	fmt.Printf(constants.MsgFlattenCloning, p.targetName, p.flattenedFolder)
	cloneResult := runGitClone(p.targetURL, p.targetPath)
	if !cloneResult {
		fmt.Fprintf(os.Stderr, constants.ErrCloneNextFailed, p.targetName)

		return apperror.NewExecutionError(fmt.Sprintf(constants.ErrCloneNextFailed, p.targetName))
	}

	fmt.Printf(constants.MsgFlattenDone, p.targetName, p.flattenedFolder)

	finalizeCloneNext(cnFlags, cwd, p)
	maybeExitOnCmdFaithfulMismatch()

	return nil
}

func finalizeCloneNext(cnFlags CloneNextFlags, cwd string, p cloneNextTargetParams) {
	if cnFlags.Force {
		fmt.Printf(constants.MsgCNStageFinalize)
	}

	recordVersionHistory(p.targetPath, p.currentVersion, p.targetVersion, p.flattenedFolder)
	if !cnFlags.NoDesktop {
		registerCloneNextDesktop(p.targetName, p.targetPath)
	}

	if p.currentFolder != p.flattenedFolder {
		keep := cnFlags.Keep || cnFlags.Force
		handleCloneNextRemoval(p.currentFolder, cwd, p.targetPath, cnFlags.Delete, keep)
	}

	WriteShellHandoff(p.targetPath)
	openInVSCode(p.targetPath)
	syncSingleClonedRepoToVSCodePM(p.targetPath, p.flattenedFolder, cnFlags.NoVSCodeSync)
	if cnFlags.Force {
		fmt.Printf(constants.MsgCNDone, p.flattenedFolder)
	}
}

func initVerboseSafe(isVerbose bool) (*verbose.Logger, error) {
	if !isVerbose {
		return nil, nil
	}

	return verbose.Init()
}

func escapeToRepoRoot(cwd string) (string, *apperror.AppError) {
	root, rootErr := gitutil.RepoRoot(cwd)
	if rootErr != nil || root == "" || root == cwd {
		return cwd, nil
	}

	chErr := os.Chdir(root)
	if chErr != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneNextCwd, chErr)

		return cwd, apperror.WrapSimple(chErr, fmt.Sprintf(constants.ErrCloneNextCwd, chErr))
	}

	fmt.Printf("cn: escaped to repo root: %s\n", root)

	return root, nil
}

func removeExistingTargetFolder(targetPath string, flattenedFolder string) *apperror.AppError {
	_, statErr := os.Stat(targetPath)
	if statErr != nil {
		return nil
	}

	fmt.Printf(constants.MsgFlattenRemoving, flattenedFolder)

	if !removeFolderWithLockCheck(flattenedFolder, targetPath) {
		msg := fmt.Sprintf(constants.ErrCloneNextForceFailed, flattenedFolder, "removal aborted or failed after lock check", flattenedFolder)
		fmt.Fprint(os.Stderr, msg)

		return apperror.NewExecutionError(msg)
	}

	return nil
}

// CreateRemoteParams encapsulates parameters for remote creation.
type CreateRemoteParams struct {
	IsCreateRemote bool
	RemoteURL      string
	TargetName     string
}

func handleCreateRemote(params CreateRemoteParams) *apperror.AppError {
	if !params.IsCreateRemote {
		return nil
	}

	owner, _, parseErr := clonenext.ParseOwnerRepo(params.RemoteURL)
	if parseErr != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneNextRemoteParse, parseErr)

		return apperror.WrapSimple(parseErr, fmt.Sprintf(constants.ErrCloneNextRemoteParse, parseErr))
	}

	return ensureAndCreateRemote(owner, params.TargetName)
}

func ensureAndCreateRemote(owner, targetName string) *apperror.AppError {
	exists, checkErr := clonenext.RepoExists(owner, targetName)
	if checkErr != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneNextRepoCheck, checkErr)

		return apperror.WrapSimple(checkErr, fmt.Sprintf(constants.ErrCloneNextRepoCheck, checkErr))
	}

	if exists {
		return nil
	}

	fmt.Printf(constants.MsgCloneNextCreating, targetName)
	createErr := clonenext.CreateRepo(owner, targetName, true)
	if createErr != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneNextRepoCreate, targetName, createErr)

		return apperror.WrapSimple(createErr, fmt.Sprintf(constants.ErrCloneNextRepoCreate, targetName, createErr))
	}

	fmt.Printf(constants.MsgCloneNextCreated, targetName)

	return nil
}

// extractRepoName extracts the repository name from a remote URL.
func extractRepoName(remoteURL string) string {
	name := remoteURL
	name = strings.TrimSuffix(name, ".git")
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		name = name[idx+1:]
	}

	if idx := strings.LastIndex(name, ":"); idx >= 0 {
		name = name[idx+1:]
	}

	return name
}

// runGitClone executes git clone via the unified colorful runner so
// clone-next shares the same header, spinner, timing, failure panel,
// and --dry-run semantics as `clone` / `cfr` / `cfrp`.
func runGitClone(url, dest string) bool {
	return runCloneCommandPretty(url, dest) == nil
}

// registerCloneNextDesktop registers the cloned repo with GitHub Desktop.
func registerCloneNextDesktop(name, absPath string) {
	records := []model.ScanRecord{{
		RepoName:     name,
		AbsolutePath: absPath,
	}}
	result := desktop.AddRepos(records)
	if result.Added > 0 {
		fmt.Printf(constants.MsgCloneNextDesktop, name)
	}
}

// handleCloneNextRemoval manages removal of the current version folder.
// It changes to the parent directory first to release file locks on Windows.
func handleCloneNextRemoval(folderName, fullPath, targetPath string, deleteFlag, keepFlag bool) {
	if keepFlag {
		return
	}

	parentDir := filepath.Dir(fullPath)
	if chErr := os.Chdir(parentDir); chErr != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not cd to %s: %v\n", parentDir, chErr)
	}

	removed := promptOrDeleteFolder(folderName, fullPath, deleteFlag)
	handlePostRemovalChdir(removed, targetPath)
}

func promptOrDeleteFolder(folderName, fullPath string, deleteFlag bool) bool {
	if deleteFlag {
		return removeFolderWithLockCheck(folderName, fullPath)
	}

	fmt.Printf(constants.MsgCloneNextRemovePrompt, folderName)
	var answer string
	_, _ = fmt.Scanln(&answer)
	shouldRemove := strings.ToLower(strings.TrimSpace(answer)) == "y"
	if shouldRemove {
		return removeFolderWithLockCheck(folderName, fullPath)
	}

	return false
}

func handlePostRemovalChdir(removed bool, targetPath string) {
	if !removed {
		return
	}

	chErr := os.Chdir(targetPath)
	if chErr != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not cd to %s: %v\n", targetPath, chErr)

		return
	}

	fmt.Printf(constants.MsgCloneNextMovedTo, filepath.Base(targetPath))
}

// removeFolderWithLockCheck attempts to remove a directory, and if it fails,
// scans for locking processes and offers to terminate them before retrying.
// All removal attempts are tracked as pending tasks in the database.
func removeFolderWithLockCheck(name, path string) bool {
	taskID, db := createPendingTask(constants.TaskTypeDelete, path, "", constants.CmdCloneNext, "")
	if db != nil {
		defer db.Close()
	}

	if os.RemoveAll(path) == nil {
		fmt.Printf(constants.MsgCloneNextRemoved, name)
		completePendingTask(db, taskID)

		return true
	}

	return handleLockedFolderRemoval(db, taskID, name, path)
}

func handleLockedFolderRemoval(db *store.DB, taskID int64, name, path string) bool {
	fmt.Printf(constants.MsgLockCheckScanning, name)
	procs, scanErr := lockcheck.FindLockingProcesses(path)
	if scanErr != nil {
		fmt.Fprintf(os.Stderr, constants.WarnLockCheckScanFailed, scanErr)
		failPendingTask(db, taskID, fmt.Sprintf(constants.ReasonLockScanFailed, scanErr))

		return false
	}

	if len(procs) == 0 {
		fmt.Print(constants.MsgLockCheckNoneFound)
		failPendingTask(db, taskID, constants.ReasonNoLockingProcs)

		return false
	}

	return confirmAndKillProcs(db, taskID, name, path, procs)
}

func confirmAndKillProcs(db *store.DB, taskID int64, name, path string, procs []lockcheck.LockingProcess) bool {
	fmt.Printf(constants.MsgLockCheckFound, lockcheck.FormatProcessList(procs))
	fmt.Print(constants.MsgLockCheckKillPrompt)
	var answer string
	_, _ = fmt.Scanln(&answer)
	if strings.ToLower(strings.TrimSpace(answer)) != "y" {
		failPendingTask(db, taskID, constants.ReasonUserDeclined)

		return false
	}

	killLockingProcesses(procs)
	time.Sleep(500 * time.Millisecond)

	return retryFolderRemoval(db, taskID, name, path)
}

func killLockingProcesses(procs []lockcheck.LockingProcess) {
	for _, p := range procs {
		fmt.Printf(constants.MsgLockCheckKilling, p.Name, p.PID)
		if err := lockcheck.KillProcess(p.PID); err != nil {
			fmt.Fprintf(os.Stderr, constants.WarnLockCheckKillFailed, p.Name, p.PID, err)
		} else {
			fmt.Printf(constants.MsgLockCheckKilled, p.Name)
		}
	}
}

func retryFolderRemoval(db *store.DB, taskID int64, name, path string) bool {
	fmt.Print(constants.MsgLockCheckRetrying)
	if err := os.RemoveAll(path); err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnCloneNextRemoveFailed, name, err)
		failPendingTask(db, taskID, fmt.Sprintf(constants.ReasonRetryFailed, err))

		return false
	}

	fmt.Printf(constants.MsgCloneNextRemoved, name)
	completePendingTask(db, taskID)

	return true
}
