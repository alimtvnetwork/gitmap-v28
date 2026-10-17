// Package cloner re-clones repos from structured files.
package cloner

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	// ErrMergeConflict indicates a git merge conflict encountered during auto-merge fallback.
	ErrMergeConflict = apperror.NewWithDetails(
		"pull_merge",
		"E_MERGE_CONFLICT",
		"git merge conflict detected: auto-merge pull aborted",
		"safe_pull",
		apperror.ErrorTypeExecution,
		apperror.SeverityError,
		nil,
	)
)

// IsMergeConflict checks whether the error represents a merge conflict.
func IsMergeConflict(err error) bool {
	return isMergeConflict(err)
}

func isMergeConflict(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrMergeConflict) {
		return true
	}

	return strings.Contains(err.Error(), "merge conflict") || strings.Contains(err.Error(), "E_MERGE_CONFLICT")
}

// IsMergeConflictOutput returns true if git output indicates merge conflicts.
func IsMergeConflictOutput(output string) bool {
	return isMergeConflictOutput(output)
}

func isMergeConflictOutput(output string) bool {
	lower := strings.ToLower(output)

	return strings.Contains(lower, "conflict") ||
		strings.Contains(lower, "automatic merge failed") ||
		strings.Contains(lower, "fix conflicts and then commit") ||
		strings.Contains(lower, "you have unmerged files")
}

// ExtractConflictedFilesFromOutput parses git conflict output lines to identify conflicted files.
func ExtractConflictedFilesFromOutput(output string) []string {
	var files []string
	seen := make(map[string]bool)

	lines := strings.Split(output, "\n")
	for _, rawLine := range lines {
		line := strings.TrimRight(rawLine, "\r ")
		file := extractConflictFileName(line)
		file = strings.Trim(strings.TrimSpace(file), "\"'")
		if file != "" && !seen[file] {
			seen[file] = true
			files = append(files, file)
		}
	}

	return files
}

func extractConflictFileName(line string) string {
	if strings.HasPrefix(line, "CONFLICT (content): Merge conflict in ") {
		return strings.TrimPrefix(line, "CONFLICT (content): Merge conflict in ")
	}
	if strings.HasPrefix(line, "CONFLICT (add/add): Merge conflict in ") {
		return strings.TrimPrefix(line, "CONFLICT (add/add): Merge conflict in ")
	}
	if strings.Contains(line, ": Merge conflict in ") {
		return parseArbitraryConflictPrefix(line)
	}

	return ""
}

func parseArbitraryConflictPrefix(line string) string {
	parts := strings.SplitN(line, ": Merge conflict in ", 2)
	if len(parts) != 2 || !strings.HasPrefix(strings.TrimSpace(parts[0]), "CONFLICT") {
		return ""
	}

	return parts[1]
}

// ExtractConflictedFilesFromRepo finds unmerged files in the repository using git status --porcelain.
func ExtractConflictedFilesFromRepo(repoDir string) []string {
	cmd := exec.Command(constants.GitBin, constants.GitDirFlag, repoDir, "status", "--porcelain")
	cmd.Env = buildSafePullEnv()
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	return parseUnmergedFilesFromPorcelain(string(out))
}

func parseUnmergedFilesFromPorcelain(output string) []string {
	var files []string
	seen := make(map[string]bool)

	lines := strings.Split(output, "\n")
	for _, rawLine := range lines {
		line := strings.TrimRight(rawLine, "\r")
		filePath, ok := parseUnmergedPorcelainLine(line)
		if ok && !seen[filePath] {
			seen[filePath] = true
			files = append(files, filePath)
		}
	}

	return files
}

func parseUnmergedPorcelainLine(line string) (string, bool) {
	if len(line) < 4 || line[2] != ' ' || !isUnmergedStatusCode(line[:2]) {
		return "", false
	}
	filePath := strings.TrimSpace(line[3:])
	idx := strings.Index(filePath, " -> ")
	if idx != -1 {
		filePath = strings.TrimSpace(filePath[idx+4:])
	}
	cleaned := strings.Trim(filePath, "\"")

	return cleaned, cleaned != ""
}

func isUnmergedStatusCode(code string) bool {
	switch code {
	case "UU", "AA", "UD", "DU", "UA", "AU", "DD":
		return true
	default:
		return len(code) == 2 && (code[0] == 'U' || code[1] == 'U')
	}
}

// SafeAbortMerge safely aborts any in-progress merge or resets merge state to protect the working tree.
func SafeAbortMerge(repoDir string) bool {
	abort := exec.Command(constants.GitBin, constants.GitDirFlag, repoDir, "merge", "--abort")
	abort.Env = buildSafePullEnv()
	if err := abort.Run(); err == nil {
		return true
	}
	reset := exec.Command(constants.GitBin, constants.GitDirFlag, repoDir, "reset", "--merge")
	reset.Env = buildSafePullEnv()

	return reset.Run() == nil
}

// SafeAbortRebase safely aborts any in-progress rebase to protect the working tree.
func SafeAbortRebase(repoDir string) bool {
	abort := exec.Command(constants.GitBin, constants.GitDirFlag, repoDir, "rebase", "--abort")
	abort.Env = buildSafePullEnv()

	return abort.Run() == nil
}
