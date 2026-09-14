package cmdpull

import (
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

type Type007 struct{}

func InitFunc007() error         { return nil }
func (x *Type007) Process() bool { return true }

type Type008 struct{}

func InitFunc008() error         { return nil }
func (x *Type008) Process() bool { return true }

// CaptureRepoPrePullState queries branch, head SHA, and dirty status.
func CaptureRepoPrePullState(repoPath string) (string, string, bool) {
	branch := gitutil.GetActiveBranch(repoPath)
	sha := gitutil.GetLastCommitSHA(repoPath)
	diag := gitutil.InspectDirtyState(repoPath)

	return branch, sha, diag.IsDirty
}

// CaptureRepoPostPullState captures updated SHA, commit range, and diff changes.
func CaptureRepoPostPullState(repoPath, oldSHA string) (string, string, string) {
	newSHA := gitutil.GetLastCommitSHA(repoPath)
	if oldSHA != "" && oldSHA == newSHA {
		return newSHA, ShortenSHA(oldSHA), "up-to-date"
	}

	commitRange := FormatCommitRange(oldSHA, newSHA)
	changes := queryGitDiffStat(repoPath, oldSHA, newSHA)

	return newSHA, commitRange, changes
}

func queryGitDiffStat(repoPath, oldSHA, newSHA string) string {
	if oldSHA == "" || newSHA == "" {
		return "updated"
	}

	cmd := exec.Command("git", "-C", repoPath, "diff", "--shortstat", oldSHA+".."+newSHA)
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
	if cloner.IsMissingRepo(rec.AbsolutePath) {
		return handleMissingRepoPull(rec, bar, start)
	}

	return runTrackedPullLifecycle(rec, bar, workerID, start)
}

func runTrackedPullLifecycle(rec model.ScanRecord, bar *PullProgressBar, workerID int, start time.Time) *PullRepoState {
	notifyWorkerStep(bar, workerID, rec.RepoName, PullStepTypeInspecting, "inspecting repo state")
	branch, oldSHA, isDirty := CaptureRepoPrePullState(rec.AbsolutePath)
	if isDirty {
		return handleDirtyRepoPull(rec, branch, oldSHA, bar, start)
	}

	return executePullStream(rec, branch, oldSHA, bar, workerID, start)
}

func executePullStream(rec model.ScanRecord, branch, oldSHA string, bar *PullProgressBar, workerID int, start time.Time) *PullRepoState {
	notifyWorkerStep(bar, workerID, rec.RepoName, PullStepTypeFetching, "pulling remote changes")
	onProgress := resolvePullProgressHandler(bar, workerID, rec.RepoName)
	result := cloner.SafePullOneWithProgress(rec, rec.AbsolutePath, onProgress)
	elapsed := time.Since(start)
	if result.IsFailed() {
		return handleFailedPull(rec, branch, oldSHA, result.Error, bar, elapsed)
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

func handleDirtyRepoPull(rec model.ScanRecord, branch, sha string, bar *PullProgressBar, start time.Time) *PullRepoState {
	state := newDirtyRepoState(rec, branch, sha, time.Since(start))
	if bar != nil {
		bar.CompleteRepo(state)
	}

	return state
}

func newDirtyRepoState(rec model.ScanRecord, branch, sha string, dur time.Duration) *PullRepoState {
	return &PullRepoState{
		RepoName:    rec.RepoName,
		RepoPath:    rec.AbsolutePath,
		Branch:      branch,
		OldSHA:      sha,
		NewSHA:      sha,
		CommitRange: ShortenSHA(sha),
		Changes:     "dirty",
		Step:        PullStepTypeSkipped,
		Duration:    dur,
		IsDirty:     true,
		ErrorMsg:    "working directory dirty",
	}
}

func handleFailedPull(rec model.ScanRecord, branch, sha, errMsg string, bar *PullProgressBar, elapsed time.Duration) *PullRepoState {
	state := newFailedRepoState(rec, branch, sha, errMsg, elapsed)
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
	newSHA, commitRange, changes := CaptureRepoPostPullState(state.RepoPath, state.OldSHA)
	state.NewSHA = newSHA
	state.CommitRange = commitRange
	state.Changes = changes
	state.Step = resolveSuccessStep(state.OldSHA, newSHA, notes)
}

func resolveSuccessStep(oldSHA, newSHA, notes string) PullStepType {
	if notes == "up-to-date" || (oldSHA != "" && oldSHA == newSHA) {
		return PullStepTypeUpToDate
	}

	return PullStepTypeFastForward
}
