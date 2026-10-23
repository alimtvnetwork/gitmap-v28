package cmdcommit

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagent"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// hasStagedChangesCP checks whether there are any staged or unstaged changes in the working tree.
func hasStagedChangesCP() (bool, error) {
	out, err := execGitOutputCP("status", "--porcelain")
	if err != nil {
		return false, err
	}

	hasChanges := len(strings.TrimSpace(out)) > 0

	return hasChanges, nil
}

// countUnpushedCommitsCP returns the count of unpushed commits ahead of upstream remote tracking branch.
func countUnpushedCommitsCP() int {
	branch, errBranch := getCurrentBranchName()
	if errBranch != nil || !hasRemoteTracking(branch) {
		return 0
	}

	return parseUnpushedCount()
}

func parseUnpushedCount() int {
	out, errCount := execGitOutputCP("rev-list", "--count", "@{u}..HEAD")
	if errCount != nil {
		return 0
	}

	count, errParse := strconv.Atoi(strings.TrimSpace(out))
	if errParse != nil {
		return 0
	}

	return count
}

func stageAndCheckChanges(taskSlug string) (bool, *apperror.AppError) {
	printPaddedInfo("Staging all changes...")

	if errAdd := executeStageAllStep(taskSlug); errAdd != nil {
		return false, errAdd
	}

	hasChanges, errStatus := hasStagedChangesCP()
	if errStatus != nil {
		return false, apperror.WrapSimple(errStatus, "check git status failed:")
	}

	return hasChanges, nil
}

func executeStageAllStep(taskSlug string) *apperror.AppError {
	startTime := time.Now()
	files := scopedStageFiles(taskSlug)
	if len(files) == 0 {
		return stageAllFallback(taskSlug, startTime)
	}

	return stageClaimedFiles(taskSlug, files, startTime)
}

func stageAllFallback(taskSlug string, startTime time.Time) *apperror.AppError {
	if strings.TrimSpace(taskSlug) != "" {
		printPaddedInfo("No file claims for task %s; staging all changes...", taskSlug)
	}
	err := execGitPaddedFiltered("add", "-A")

	return recordAndWrapStage(err, "git add -A", startTime)
}

func stageClaimedFiles(taskSlug string, files []string, startTime time.Time) *apperror.AppError {
	printPaddedInfo("Staging %d task-claimed file(s) for task %s...", len(files), taskSlug)
	normal, forced := partitionStageFiles(files)
	err := stagePartitionFiles(normal, forced)
	label := fmt.Sprintf("git add -- (%d claimed files)", len(files))

	return recordAndWrapStage(err, label, startTime)
}

func stagePartitionFiles(normal, forced []string) error {
	var err error
	if len(normal) > 0 {
		gitArgs := append([]string{"add", "--"}, normal...)
		err = execGitPaddedFiltered(gitArgs...)
	}
	if err != nil || len(forced) == 0 {
		return err
	}
	forceArgs := append([]string{"add", "-f", "--"}, forced...)

	return execGitPaddedFiltered(forceArgs...)
}

func recordAndWrapStage(err error, label string, startTime time.Time) *apperror.AppError {
	durMs := time.Since(startTime).Milliseconds()
	recordGitCommandDirect("git add", label, ExtractGitExitCode(err), durMs)
	if err != nil {
		return apperror.WrapSimple(err, "git add failed:")
	}

	return nil
}

// partitionStageFiles splits claimed files for staging: ignored-but-tracked
// files need `git add -f`; ignored-and-untracked files are skipped (gitignore
// is respected); everything else stages normally. A single ignored path must
// never poison the whole stage.
func partitionStageFiles(files []string) (normal []string, forced []string) {
	ignored := ignoredPaths(files)
	trackedIgnored := trackedPaths(ignored)
	trackedSet := make(map[string]bool, len(trackedIgnored))
	for _, t := range trackedIgnored {
		trackedSet[t] = true
	}
	ignoredSet := make(map[string]bool, len(ignored))
	for _, g := range ignored {
		ignoredSet[g] = true
	}
	for _, f := range files {
		if !ignoredSet[f] {
			normal = append(normal, f)
		} else if trackedSet[f] {
			forced = append(forced, f)
		}
	}

	return normal, forced
}

func runCheckIgnoreCommand(files []string) ([]byte, error) {
	c := exec.Command("git", "check-ignore", "--no-index", "--stdin")
	stdin, err := c.StdinPipe()
	if err != nil {
		return nil, err
	}
	go func() {
		defer stdin.Close()
		for _, f := range files {
			fmt.Fprintln(stdin, f)
		}
	}()

	return c.Output()
}

// ignoredPaths returns the subset of files ignored by .gitignore.
// --no-index is required: without it, check-ignore skips tracked files even
// though `git add` still refuses them without -f.
func ignoredPaths(files []string) []string {
	if len(files) == 0 {
		return nil
	}
	out, err := runCheckIgnoreCommand(files)
	if err != nil {
		return nil
	}
	var ignored []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if strings.TrimSpace(line) != "" {
			ignored = append(ignored, strings.TrimSpace(line))
		}
	}

	return ignored
}

// trackedPaths returns the subset of files already tracked by git.
func trackedPaths(files []string) []string {
	if len(files) == 0 {
		return nil
	}
	args := append([]string{"ls-files", "--"}, files...)
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	var tracked []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if strings.TrimSpace(line) != "" {
			tracked = append(tracked, strings.TrimSpace(line))
		}
	}

	return tracked
}

// scopedStageFiles returns the union of a task's claimed files for scoped
// staging. Empty when no task context or no claims (caller falls back to -A).
func scopedStageFiles(taskSlug string) []string {
	if strings.TrimSpace(taskSlug) == "" {
		return nil
	}
	files, err := cmdagent.TaskClaimedFiles(taskSlug)
	if err != nil || len(files) == 0 {
		return nil
	}

	return files
}

func ExtractGitExitCode(err error) int {
	if err == nil {
		return 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}

	return 1
}

func recordGitCommandDirect(cmdName, cmdLine string, exitCode int, durationMs int64) {
	store.RecordCommandSafely(cmdLine, cmdName, exitCode, durationMs)
}

func RecordGitCommandHistory(gitArgs []string, exitCode int, durationMs int64) {
	if len(gitArgs) == 0 {
		recordGitCommandDirect("git", "git", exitCode, durationMs)
		return
	}
	cmdName := "git " + gitArgs[0]
	cmdLine := "git " + strings.Join(gitArgs, " ")
	recordGitCommandDirect(cmdName, cmdLine, exitCode, durationMs)
}
