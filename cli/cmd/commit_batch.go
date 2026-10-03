package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
)

// isAllCommitRequested checks whether "all", "--all", or "-a" was passed or "commit-all"/"ca" was invoked.
func isAllCommitRequested(args []string) bool {
	if isAllSubcommand() {
		return true
	}

	return hasAllFlag(args)
}

func isAllSubcommand() bool {
	if len(os.Args) <= 1 {
		return false
	}

	sub := os.Args[1]

	return sub == "commit-all" || sub == "ca"
}

func hasAllFlag(args []string) bool {
	for _, arg := range args {
		if isAllArg(arg) {
			return true
		}
	}

	return false
}

func isAllArg(arg string) bool {
	return arg == "all" || arg == "--all" || arg == "-a"
}

// stripAllFlags removes "all", "--all", "-a" flags. If empty, defaults to a chore commit message.
func stripAllFlags(args []string) []string {
	var remaining []string
	for _, arg := range args {
		if isAllArg(arg) {
			continue
		}

		remaining = append(remaining, arg)
	}

	if len(remaining) == 0 {
		return []string{"chore: commit pending changes"}
	}

	return remaining
}

// handleNonGitRepoCommit handles commit attempts when current directory is not a git repo.
func handleNonGitRepoCommit(cleanArgs []string, hasPush, isDryRun bool) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	childRepos, _ := fsutil.DiscoverChildGitRepos(cwd)
	isAll := isAllCommitRequested(cleanArgs)

	if isAll && len(childRepos) > 0 {
		return executeChildReposBatch(childRepos, cleanArgs, hasPush, isDryRun)
	}

	return handleNonGitRepoFallback()
}

func handleNonGitRepoFallback() error {
	printPaddedError("Not a git repository (or any of the parent directories).")
	printPaddedInfo("To commit across child repositories, run: gitmap commit all")

	return newNonGitRepoAbortError("not a git repository")
}

func executeChildReposBatch(childRepos []string, cleanArgs []string, hasPush, isDryRun bool) error {
	msgArgs := stripAllFlags(cleanArgs)
	commitMessage := strings.Join(msgArgs, " ")

	if isDryRun {
		return previewChildReposBatch(childRepos)
	}

	return commitChildReposBatch(childRepos, commitMessage, hasPush)
}

func hasChildRepoChanges(repoPath string) bool {
	out, err := exec.Command("git", "-C", repoPath, "status", "--porcelain").Output()
	if err != nil {
		return false
	}

	return len(strings.TrimSpace(string(out))) > 0
}

func previewChildReposBatch(childRepos []string) error {
	printPaddedInfo("Dry run: discovering dirty child repositories...")

	dirtyCount := previewDirtyRepos(childRepos)
	if dirtyCount == 0 {
		printPaddedSuccess("All %d child repositories are clean.", len(childRepos))
	}

	return nil
}

func previewDirtyRepos(childRepos []string) int {
	dirtyCount := 0
	for _, repo := range childRepos {
		if hasChildRepoChanges(repo) {
			dirtyCount++
			previewSingleChildRepo(repo)
		}
	}

	return dirtyCount
}

func previewSingleChildRepo(repo string) {
	repoName := filepath.Base(repo)
	printPaddedInfo("Child repo %s has pending changes:", repoName)
	_ = execGitPaddedFiltered("-C", repo, "status", "--short")
}

func commitChildReposBatch(childRepos []string, msg string, hasPush bool) error {
	committedCount := commitAllDirtyRepos(childRepos, msg, hasPush)
	if committedCount == 0 {
		printPaddedSuccess("All %d child repositories are clean, nothing to commit.", len(childRepos))

		return nil
	}

	printPaddedSuccess("Committed across %d child repositories.", committedCount)

	return nil
}

func commitAllDirtyRepos(childRepos []string, msg string, hasPush bool) int {
	committedCount := 0
	for _, repo := range childRepos {
		if hasChildRepoChanges(repo) {
			committedCount++
			commitSingleChildRepo(repo, msg, hasPush)
		}
	}

	return committedCount
}

func commitSingleChildRepo(repo, msg string, hasPush bool) {
	repoName := filepath.Base(repo)
	printPaddedInfo("Committing in child repo: %s", repoName)
	_ = execGitPaddedFiltered("-C", repo, "add", "-A")
	_ = execGitPaddedFiltered("-C", repo, "commit", "-m", msg)
	if hasPush {
		printPaddedInfo("Pushing child repo: %s", repoName)
		_ = execGitPaddedFiltered("-C", repo, "push")
	}
}
