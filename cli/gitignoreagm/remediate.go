package gitignoreagm

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const PrimaryIgnoreEntry = ".antigravity_resume_task.json"
const SecondaryIgnoreEntry = "antigravity-resume_task.json"

var targetResumeFiles = []string{
	".antigravity_resume_task.json",
	"antigravity-resume_task.json",
	"antigravity_resume_task.json",
	".antigravity-resume_task.json",
}

// RepoRemediationResult captures the actions taken on a single repository.
type RepoRemediationResult struct {
	RepoPath       string
	WasUntracked   bool
	WasFileDeleted bool
	WasIgnored     bool
	WasCommitted   bool
}

// IsGitRepository checks whether dir contains a valid .git directory or file.
func IsGitRepository(dir string) bool {
	if dir == "" {
		return false
	}
	_, err := os.Stat(filepath.Join(dir, ".git"))
	return err == nil
}

// HasUnignoredResumeTask returns true if a repo has the resume task file on disk or tracked in git.
func HasUnignoredResumeTask(repoDir string) bool {
	if !IsGitRepository(repoDir) {
		return false
	}
	if len(findTrackedResumeFiles(repoDir)) > 0 {
		return true
	}
	if hasAnyResumeFileOnDisk(repoDir) {
		return true
	}
	return false
}

func hasAnyResumeFileOnDisk(repoDir string) bool {
	for _, name := range targetResumeFiles {
		if fileExistsOnDisk(filepath.Join(repoDir, name)) {
			return true
		}
	}
	return false
}

func fileExistsOnDisk(p string) bool {
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}

func findTrackedResumeFiles(repoDir string) []string {
	cmdArgs := append([]string{"-C", repoDir, "ls-files", "--"}, targetResumeFiles...)
	out, err := exec.Command("git", cmdArgs...).Output()
	if err != nil {
		return nil
	}
	return parseNonEmptyLines(string(out))
}

func parseNonEmptyLines(raw string) []string {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// RemediateRepo removes resume task files from git tracking and disk, updates .gitignore, and commits.
func RemediateRepo(repoDir string, shouldCommit bool) (RepoRemediationResult, error) {
	res := RepoRemediationResult{RepoPath: repoDir}
	if !IsGitRepository(repoDir) {
		return res, apperror.NewSimple("not a git repository: "+repoDir, "E_NOT_GIT_REPO")
	}
	res.WasUntracked = untrackResumeFiles(repoDir)
	res.WasFileDeleted = deletePhysicalResumeFiles(repoDir)
	wasAdded, err := ensureGitignoreEntries(repoDir)
	if err != nil {
		return res, err
	}
	res.WasIgnored = wasAdded
	if shouldCommit && (res.WasUntracked || res.WasFileDeleted || res.WasIgnored) {
		res.WasCommitted = stageAndCommitGitignore(repoDir)
	}
	return res, nil
}

func untrackResumeFiles(repoDir string) bool {
	tracked := findTrackedResumeFiles(repoDir)
	if len(tracked) == 0 {
		return false
	}
	args := append([]string{"-C", repoDir, "rm", "--cached", "-f", "--ignore-unmatch", "--"}, tracked...)
	err := exec.Command("git", args...).Run()
	return err == nil
}

func deletePhysicalResumeFiles(repoDir string) bool {
	wasDeleted := false
	for _, name := range targetResumeFiles {
		fullPath := filepath.Join(repoDir, name)
		if removeIfPresent(fullPath) {
			wasDeleted = true
		}
	}
	return wasDeleted
}

func removeIfPresent(fullPath string) bool {
	if !fileExistsOnDisk(fullPath) {
		return false
	}
	return os.Remove(fullPath) == nil
}

func ensureGitignoreEntries(repoDir string) (bool, error) {
	ignorePath := filepath.Join(repoDir, ".gitignore")
	data, err := readGitignoreFile(ignorePath)
	if err != nil {
		return false, err
	}
	entries := []string{PrimaryIgnoreEntry, SecondaryIgnoreEntry}
	updated, changed := appendMissingIgnoreLines(string(data), entries)
	if !changed {
		return false, nil
	}
	writeErr := os.WriteFile(ignorePath, []byte(updated), 0o644)
	if writeErr != nil {
		return false, apperror.WrapSimple(writeErr, "failed writing .gitignore")
	}
	return true, nil
}

func readGitignoreFile(ignorePath string) ([]byte, error) {
	data, err := os.ReadFile(ignorePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, apperror.WrapSimple(err, "failed reading .gitignore")
	}
	return data, nil
}

func appendMissingIgnoreLines(content string, entries []string) (string, bool) {
	existing := buildExistingLineMap(content)
	var missing []string
	for _, entry := range entries {
		if !existing[entry] && !existing["/"+entry] {
			missing = append(missing, entry)
		}
	}
	if len(missing) == 0 {
		return content, false
	}
	prefix := ensureTrailingNewline(content)
	return prefix + strings.Join(missing, "\n") + "\n", true
}

func buildExistingLineMap(content string) map[string]bool {
	lines := strings.Split(content, "\n")
	existing := make(map[string]bool, len(lines))
	for _, line := range lines {
		existing[strings.TrimSpace(line)] = true
	}
	return existing
}

func ensureTrailingNewline(content string) string {
	if len(content) == 0 || strings.HasSuffix(content, "\n") {
		return content
	}
	return content + "\n"
}

func stageAndCommitGitignore(repoDir string) bool {
	_ = exec.Command("git", "-C", repoDir, "add", ".gitignore").Run()
	if !hasStagedGitChanges(repoDir) {
		return false
	}
	commitErr := exec.Command("git", "-C", repoDir, "commit", "-m", "Update .gitignore").Run()
	return commitErr == nil
}

func hasStagedGitChanges(repoDir string) bool {
	err := exec.Command("git", "-C", repoDir, "diff", "--cached", "--quiet").Run()
	return err != nil
}
