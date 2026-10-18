package cmdpull

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func InitFunc007() error { return nil }
func InitFunc008() error { return nil }

// CaptureRepoPrePullState queries branch, head SHA, and dirty status.
func CaptureRepoPrePullState(repoPath string) (string, string, bool) {
	cleanPath := strings.TrimSpace(repoPath)
	if cleanPath == "" {
		return "", "-", false
	}

	branch := gitutil.GetActiveBranch(cleanPath)
	sha := gitutil.GetLastCommitSHA(cleanPath)

	return branch, sha, false
}

// CaptureRepoPostPullState captures updated SHA, commit range, and diff changes.
func CaptureRepoPostPullState(repoPath, oldSHA string) (string, string, string) {
	cleanPath := strings.TrimSpace(repoPath)
	if cleanPath == "" {
		return "-", "-", "updated"
	}

	newSHA := gitutil.GetLastCommitSHA(cleanPath)
	hasValidOld := oldSHA != "" && oldSHA != "-"
	hasValidNew := newSHA != "" && newSHA != "-"
	hasIdenticalSHA := hasValidOld && hasValidNew && oldSHA == newSHA
	if hasIdenticalSHA {
		return newSHA, ShortenSHA(oldSHA), "up-to-date"
	}

	commitRange := FormatCommitRange(oldSHA, newSHA)
	changes := queryGitDiffStat(cleanPath, oldSHA, newSHA)

	return newSHA, commitRange, changes
}

func queryGitDiffStat(repoPath, oldSHA, newSHA string) string {
	cleanPath := strings.TrimSpace(repoPath)
	hasEmptyPath := cleanPath == ""
	hasInvalidOld := oldSHA == "" || oldSHA == "-"
	hasInvalidNew := newSHA == "" || newSHA == "-"
	hasIdenticalSHA := oldSHA == newSHA
	if hasEmptyPath || hasInvalidOld || hasInvalidNew || hasIdenticalSHA {
		return "updated"
	}

	fi, err := os.Stat(cleanPath)
	isDir := err == nil && fi.IsDir()
	if !isDir {
		return "updated"
	}

	cmd := exec.Command("git", "-C", cleanPath, "diff", "--shortstat", oldSHA+".."+newSHA)
	out, err := cmd.Output()
	if err != nil {
		return "updated"
	}

	return formatDiffStatOutput(string(out))
}

func formatDiffStatOutput(raw string) string {
	stat := FormatDiffStat(strings.TrimSpace(raw))
	if stat == "" || stat == "synced" {
		return "updated"
	}

	return stat
}

// ExecuteTrackedPull executes pull on a scan record reporting state to bar.
func ExecuteTrackedPull(rec model.ScanRecord, bar *PullProgressBar) *PullRepoState {
	return ExecuteTrackedPullWithWorker(rec, bar, -1)
}

// ExecuteTrackedPullWithWorker executes pull associating progress with a worker ID.
func ExecuteTrackedPullWithWorker(rec model.ScanRecord, bar *PullProgressBar, workerID int) *PullRepoState {
	start := time.Now()
	if !cloner.IsMissingRepo(rec.AbsolutePath) {
		return runTrackedPullLifecycle(rec, bar, workerID, start)
	}

	if TryHealMissingRecordPath(&rec) {
		return runTrackedPullLifecycle(rec, bar, workerID, start)
	}

	return handleMissingRepoPull(rec, bar, start)
}

func runTrackedPullLifecycle(rec model.ScanRecord, bar *PullProgressBar, workerID int, start time.Time) *PullRepoState {
	sanitizeScanRecordIdentity(&rec)
	notifyWorkerStep(bar, workerID, rec.RepoName, PullStepTypeInspecting, "inspecting repo state")
	if !isRemoteConfigured(rec) {
		return handleNoRemoteRepoPull(rec, bar, start)
	}

	branch, oldSHA, isPrePullDirty := CaptureRepoPrePullState(rec.AbsolutePath)

	// Attempt git pull --autostash directly instead of prematurely skipping dirty trees.
	return executePullStreamWithAutostash(rec, branch, oldSHA, isPrePullDirty, bar, workerID, start)
}

func isRemoteConfigured(rec model.ScanRecord) bool {
	if len(rec.HTTPSUrl) > 0 || len(rec.SSHUrl) > 0 {
		return true
	}

	return hasRemote(rec.AbsolutePath)
}

func hasRemote(repoPath string) bool {
	remoteURL, err := gitutil.RemoteURL(repoPath)
	if err != nil || len(strings.TrimSpace(remoteURL)) == 0 {
		return false
	}

	return true
}

func handleNoRemoteRepoPull(rec model.ScanRecord, bar *PullProgressBar, start time.Time) *PullRepoState {
	state := &PullRepoState{
		RepoName: rec.RepoName,
		RepoPath: rec.AbsolutePath,
		Changes:  "local-only",
		Step:     PullStepTypeSkipped,
		Duration: time.Since(start),
		ErrorMsg: "no remote configured",
	}
	if bar != nil {
		bar.CompleteRepo(state)
	}

	return state
}

func sanitizeScanRecordIdentity(rec *model.ScanRecord) {
	if rec.RepoName != "" && rec.RepoName != "unknown" {
		return
	}
	clean := filepath.Clean(strings.TrimRight(rec.AbsolutePath, "/\\"))
	base := filepath.Base(clean)
	if base == "" || base == "." || base == "/" || base == "\\" {
		return
	}
	rec.RepoName = base
	if rec.Slug == "" || rec.Slug == "unknown" {
		rec.Slug = strings.ToLower(base)
	}
}

func executePullStreamWithAutostash(rec model.ScanRecord, branch, oldSHA string, isPrePullDirty bool, bar *PullProgressBar, workerID int, start time.Time) *PullRepoState {
	notifyWorkerStep(bar, workerID, rec.RepoName, PullStepTypeFetching, "pulling remote changes")
	onProgress := resolvePullProgressHandler(bar, workerID, rec.RepoName)
	result := cloner.SafePullOneWithProgress(rec, rec.AbsolutePath, onProgress)
	elapsed := time.Since(start)
	if result.IsFailed() {
		return handleFailedAutostashPull(rec, branch, oldSHA, isPrePullDirty, result.Error, bar, elapsed)
	}

	return handleSuccessfulPull(rec, branch, oldSHA, result.Notes, bar, elapsed)
}

func resolvePullProgressHandler(bar *PullProgressBar, workerID int, repoName string) func(string) {
	if bar == nil {
		return nil
	}

	return MakePullProgressHandler(bar, workerID, repoName)
}

func notifyWorkerStep(bar *PullProgressBar, workerID int, repoName string, step PullStepType, desc string) {
	if bar == nil {
		return
	}
	isWorkerActive := workerID >= 0
	if isWorkerActive {
		bar.UpdateWorkerProgress(workerID, step, desc, StepMilestonePercent(step))
		return
	}
	bar.UpdateRepoStep(repoName, step, desc)
}

func handleMissingRepoPull(rec model.ScanRecord, bar *PullProgressBar, start time.Time) *PullRepoState {
	state := &PullRepoState{
		RepoName: rec.RepoName,
		RepoPath: rec.AbsolutePath,
		Step:     PullStepTypeSkipped,
		Duration: time.Since(start),
		ErrorMsg: "missing repository directory",
	}
	if bar != nil {
		bar.CompleteRepo(state)
	}

	return state
}

func handleFailedAutostashPull(rec model.ScanRecord, branch, sha string, isPrePullDirty bool, errMsg string, bar *PullProgressBar, elapsed time.Duration) *PullRepoState {
	if !isPrePullDirty {
		diag := gitutil.InspectDirtyState(rec.AbsolutePath)
		isPrePullDirty = diag.IsDirty
	}
	state := newFailedRepoState(rec, branch, sha, errMsg, elapsed)
	state.IsDirty = isPrePullDirty
	if bar != nil {
		bar.CompleteRepo(state)
	}

	return state
}

func newFailedRepoState(rec model.ScanRecord, branch, sha, errMsg string, dur time.Duration) *PullRepoState {
	return &PullRepoState{
		RepoName:    rec.RepoName,
		RepoPath:    rec.AbsolutePath,
		Branch:      branch,
		OldSHA:      sha,
		NewSHA:      sha,
		CommitRange: ShortenSHA(sha),
		Changes:     "failed",
		Step:        PullStepTypeError,
		Duration:    dur,
		ErrorMsg:    errMsg,
	}
}

func handleSuccessfulPull(rec model.ScanRecord, branch, oldSHA, notes string, bar *PullProgressBar, elapsed time.Duration) *PullRepoState {
	state := initSuccessRepoState(rec, branch, oldSHA, elapsed)
	finalizeSuccessRepoState(state, notes)
	if bar != nil {
		bar.CompleteRepo(state)
	}

	return state
}

func initSuccessRepoState(rec model.ScanRecord, branch, oldSHA string, dur time.Duration) *PullRepoState {
	return &PullRepoState{
		RepoName: rec.RepoName,
		RepoPath: rec.AbsolutePath,
		Branch:   branch,
		OldSHA:   oldSHA,
		Duration: dur,
	}
}

func finalizeSuccessRepoState(state *PullRepoState, notes string) {
	if notes == "auto-merged" {
		state.IsAutoMerged = true
	}
	if notes == "up-to-date" {
		state.NewSHA = state.OldSHA
		state.CommitRange = ShortenSHA(state.OldSHA)
		state.Changes = "up-to-date"
		state.Step = PullStepTypeUpToDate
		return
	}
	newSHA, commitRange, changes := CaptureRepoPostPullState(state.RepoPath, state.OldSHA)
	state.NewSHA = newSHA
	state.CommitRange = commitRange
	state.Changes = changes
	state.Step = resolveSuccessStep(state.OldSHA, newSHA, notes)
}

func resolveSuccessStep(oldSHA, newSHA, notes string) PullStepType {
	if notes == "auto-merged" {
		return PullStepTypeMerging
	}
	if notes == "up-to-date" || (oldSHA != "" && oldSHA == newSHA) {
		return PullStepTypeUpToDate
	}

	return PullStepTypeFastForward
}
