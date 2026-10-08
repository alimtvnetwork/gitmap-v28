package cmdpull

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func resolveExplicitPullTargets(opts pullOptions) ([]model.ScanRecord, bool) {
	records := resolvePullTargets(opts.slug, opts.group, opts.all)
	if len(records) == 0 {
		warnPullTargetNotFoundUnlessJSON(opts)

		return nil, false
	}

	return records, true
}

func warnPullTargetNotFoundUnlessJSON(opts pullOptions) {
	if !opts.isJSON {
		handlePullTargetNotFound(opts)
	}
}

func resolvePullBatchRecords(opts pullOptions) ([]model.ScanRecord, bool) {
	if hasExplicitPullTarget(opts) {
		return resolveExplicitPullTargets(opts)
	}
	records := discoverCWDPullRecords()
	if len(records) == 0 {
		handleNoTrackedReposInCWD()

		return nil, false
	}

	return deduplicatePullRecords(records), true
}

func hasExplicitPullTarget(opts pullOptions) bool {
	return opts.slug != "" || opts.group != "" || opts.all || HasAlias()
}

func discoverCWDPullRecords() []model.ScanRecord {
	cwd, _ := os.Getwd()
	records := ResolvePullDirectoryTargets(cwd)
	if len(records) == 0 {
		return findChildrenOfCWD(cwd)
	}

	return records
}

func handleNoTrackedReposInCWD() {
	arrow := resolveSubArrow()
	fmt.Printf("    %s%s%s %snothing to pull: no tracked repositories found in or under this directory.%s\n\n",
		constants.ColorCyan, arrow, constants.ColorReset,
		constants.ColorDim, constants.ColorReset)
	cliexit.HandleError(nil, int(cliexit.ExitCodeNotFound))
}

func handleNonGitPull(cwd string, extraArgs []string) error {
	childRepos, err := fsutil.DiscoverChildGitRepos(cwd)
	if err == nil && len(childRepos) > 0 {
		return pullDiscoveredChildren(cwd, childRepos, extraArgs)
	}
	fmt.Fprintln(os.Stderr, "✗ not a git repository (run `gitmap pull` inside a repo)")
	cliexit.HandleValidationError(apperror.NewValidationError("not a git repository (run `gitmap pull` inside a repo)"))

	return nil
}

func pullDiscoveredChildren(cwd string, childRepos []string, extraArgs []string) error {
	arrow := resolveSubArrow()
	fmt.Printf("    %s%s%s Discovered %d child repositories in %s for pull:\n",
		constants.ColorCyan, arrow, constants.ColorReset, len(childRepos), cwd)
	for _, r := range childRepos {
		pullSingleDiscoveredChild(r, extraArgs)
	}

	return nil
}

func pullSingleDiscoveredChild(repoPath string, extraArgs []string) {
	fmt.Printf("      • %s\n", filepath.Base(repoPath))
	gitArgs := append([]string{"-C", repoPath, "pull"}, extraArgs...)
	cmd := exec.Command("git", gitArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}

func enqueueTaskQueue(action, target string) (string, *store.TasksSplitDB) {
	queueId := fmt.Sprintf("%s-%d", action, time.Now().UnixNano())
	tasksDB, err := store.OpenTasksRootSplitDB()
	hasErr := err != nil
	if hasErr {
		return "", nil
	}
	now := time.Now().UTC().Format(time.RFC3339)
	store.ExecWrapper(tasksDB.Conn(), "INSERT INTO TaskQueue (QueueId, Section, Action, Target, Status, CreatedAt, UpdatedAt) VALUES (?, 'pull', ?, ?, 'pending', ?, ?)", queueId, action, target, now, now)
	return queueId, tasksDB
}

func updateTaskQueue(db *store.TasksSplitDB, queueId, status string) {
	hasDB := db != nil
	if !hasDB {
		return
	}
	now := time.Now().UTC().Format(time.RFC3339)
	store.ExecWrapper(db.Conn(), "UPDATE TaskQueue SET Status = ?, UpdatedAt = ? WHERE QueueId = ?", status, now, queueId)
}

func closeTaskDB(db *store.TasksSplitDB) {
	hasDB := db != nil
	if hasDB {
		_ = db.Close()
	}
}
