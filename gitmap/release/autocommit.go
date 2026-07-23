package release

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/verbose"
)

// AutoCommitResult describes what happened during auto-commit.
type AutoCommitResult struct {
	Committed bool
	AllFiles  bool
	Message   string
}

// AutoCommit inspects working tree changes after returning to the original branch.
// If only .gitmap/release/ files changed, it commits and pushes silently.
// If other files also changed, it prompts the user (or auto-confirms with yes=true).
// On decline, it commits only .gitmap/release/.
func AutoCommit(version string, dryRun, yes bool) AutoCommitResult {
	fmt.Print(constants.MsgAutoCommitScanning)

	if verbose.IsEnabled() {
		verbose.Get().Log("autocommit: starting for %s (dry-run=%v)", version, dryRun)
	}

	if dryRun {
		fmt.Print(constants.MsgAutoCommitDryRun)

		return AutoCommitResult{}
	}

	changed := listChangedFiles()
	if len(changed) == 0 {
		fmt.Print(constants.MsgAutoCommitNone)

		if verbose.IsEnabled() {
			verbose.Get().Log("autocommit: no changed files detected")
		}

		return AutoCommitResult{}
	}

	releaseFiles, otherFiles := classifyFiles(changed)

	if verbose.IsEnabled() {
		verbose.Get().Log("autocommit: %d release file(s), %d other file(s)", len(releaseFiles), len(otherFiles))
	}

	commitMsg := fmt.Sprintf(constants.AutoCommitMsgFmt, version)

	if len(otherFiles) == 0 {
		return commitReleaseOnly(releaseFiles, commitMsg)
	}

	return promptAndCommit(releaseFiles, otherFiles, commitMsg, yes)
}

// listChangedFiles returns all modified/untracked files in the working tree.
func listChangedFiles() []string {
	cmd := exec.Command(constants.GitBin, constants.GitStatus, constants.GitStatusShort)
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	return parsePorcelainOutput(string(out))
}

// parsePorcelainOutput extracts file paths from git status --porcelain output.
func parsePorcelainOutput(output string) []string {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var files []string

	for _, line := range lines {
		if len(line) < 4 {
			continue
		}

		path := strings.TrimSpace(line[3:])
		if len(path) > 0 {
			files = append(files, path)
		}
	}

	return files
}

// classifyFiles separates .gitmap/release/ files (and legacy .release/ files) from everything else.
func classifyFiles(files []string) (releaseFiles, otherFiles []string) {
	for _, f := range files {
		if strings.HasPrefix(f, constants.DefaultReleaseDir+"/") || f == constants.DefaultReleaseDir ||
			strings.HasPrefix(f, constants.LegacyReleaseDir+"/") || f == constants.LegacyReleaseDir {
			releaseFiles = append(releaseFiles, f)
		} else {
			otherFiles = append(otherFiles, f)
		}
	}

	return releaseFiles, otherFiles
}

// commitReleaseOnly stages and commits only .gitmap/release/ files.
func commitReleaseOnly(files []string, msg string) AutoCommitResult {
	err := stageFiles(files)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAutoCommitFailed, err)

		return AutoCommitResult{}
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("autocommit: staged %d release file(s)", len(files))
	}

	err = commitStaged(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAutoCommitFailed, err)

		return AutoCommitResult{}
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("autocommit: committed release-only: %s", msg)
	}

	fmt.Printf(constants.MsgAutoCommitReleaseOnly, msg)

	err = pushCurrentBranch()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAutoCommitPush, err)

		return AutoCommitResult{Committed: true, Message: msg}
	}

	branch, branchErr := CurrentBranchName()
	if branchErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not determine current branch: %v\n", branchErr)
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("autocommit: pushed to %s", branch)
	}

	fmt.Printf(constants.MsgAutoCommitPushed, branch)

	return AutoCommitResult{Committed: true, Message: msg}
}

// promptAndCommit shows changed files and asks the user whether to commit all.
// If yes is true, it skips the interactive prompt and commits all changes.
func promptAndCommit(releaseFiles, otherFiles []string, msg string, yes bool) AutoCommitResult {
	fmt.Print(constants.MsgAutoCommitPrompt)

	for _, f := range otherFiles {
		fmt.Printf(constants.MsgAutoCommitFile, f)
	}

	if yes {
		fmt.Print(constants.MsgAutoCommitAutoYes)

		return commitAll(msg)
	}

	fmt.Print(constants.MsgAutoCommitAsk)

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return commitReleaseOnly(releaseFiles, msg)
	}

	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	if answer == "y" || answer == "yes" {
		return commitAll(msg)
	}

	if len(releaseFiles) > 0 {
		result := commitReleaseOnly(releaseFiles, msg)
		fmt.Printf(constants.MsgAutoCommitPartial, msg)

		return result
	}

	return AutoCommitResult{}
}
