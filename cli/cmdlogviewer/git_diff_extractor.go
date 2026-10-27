package cmdlogviewer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdgit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// TriggerExternalDiff prepares file revisions and launches the external compare tool.
func TriggerExternalDiff(repoRoot, commitHash, filePath string) (string, error) {
	return TriggerExternalDiffWithTool(repoRoot, commitHash, filePath, "")
}

// TriggerExternalDiffWithTool prepares file revisions and launches with user preferred tool.
func TriggerExternalDiffWithTool(repoRoot, commitHash, filePath, preferredTool string) (string, error) {
	leftPath, rightPath, err := prepareDiffFiles(repoRoot, commitHash, filePath)
	if err != nil {
		return "", err
	}

	tool, err := cmdgit.FindCompareToolWithPreference(preferredTool)
	if err != nil {
		return "", err
	}

	err = cmdgit.LaunchToolWithPathsAsync(tool, []string{leftPath, rightPath})
	if err != nil {
		return tool.Name, err
	}

	return tool.Name, nil
}

func prepareDiffFiles(repoRoot, commitHash, filePath string) (string, string, error) {
	leftPath, err := resolveLeftDiffFile(repoRoot, commitHash, filePath)
	if err != nil {
		return "", "", err
	}

	rightPath, err := resolveRightDiffFile(repoRoot, commitHash, filePath)
	if err != nil {
		return "", "", err
	}

	return leftPath, rightPath, nil
}

func resolveLeftDiffFile(repoRoot, commitHash, filePath string) (string, error) {
	parentRef, hasParent := resolveParentCommit(repoRoot, commitHash)
	if !hasParent {
		return createTempDiffFile("initial_empty", filePath, nil)
	}

	data, err := readGitFileBytes(repoRoot, parentRef+":"+filePath)
	if err != nil {
		return createTempDiffFile("parent_empty", filePath, nil)
	}

	return createTempDiffFile(shortenHash(parentRef), filePath, data)
}

func resolveRightDiffFile(repoRoot, commitHash, filePath string) (string, error) {
	if isWorkingTreeRef(commitHash) {
		return filepath.Join(repoRoot, filePath), nil
	}

	data, err := readGitFileBytes(repoRoot, commitHash+":"+filePath)
	if err != nil {
		return createTempDiffFile("commit_deleted", filePath, nil)
	}

	return createTempDiffFile(shortenHash(commitHash), filePath, data)
}

func isWorkingTreeRef(ref string) bool {
	clean := strings.ToLower(strings.TrimSpace(ref))

	return clean == "" || clean == "working" || clean == "worktree" || clean == "index"
}

func resolveParentCommit(repoRoot, commitHash string) (string, bool) {
	out, err := runGitInDir(repoRoot, "rev-parse", "--verify", commitHash+"^")
	if err != nil {
		return "", false
	}

	return strings.TrimSpace(out), true
}

func shortenHash(hash string) string {
	if len(hash) > 8 {
		return hash[:8]
	}

	return hash
}

// TriggerTwoCommitDiff launches external diff comparing a file across two commits.
func TriggerTwoCommitDiff(repoRoot, hashA, hashB, filePath, preferredTool string) (string, error) {
	if strings.TrimSpace(filePath) == "" || filePath == "." {
		return triggerCommitsTreeCompare(repoRoot, hashA, hashB, preferredTool)
	}

	leftPath, err := extractBranchFile(repoRoot, hashA, filePath)
	if err != nil {
		return "", err
	}

	rightPath, err := extractBranchFile(repoRoot, hashB, filePath)
	if err != nil {
		return "", err
	}

	tool, err := cmdgit.FindCompareToolWithPreference(preferredTool)
	if err != nil {
		return "", err
	}

	err = cmdgit.LaunchToolWithPathsAsync(tool, []string{leftPath, rightPath})
	if err != nil {
		return tool.Name, err
	}

	return tool.Name, nil
}

func triggerCommitsTreeCompare(repoRoot, hashA, hashB, preferredTool string) (string, error) {
	tool, err := cmdgit.FindCompareToolWithPreference(preferredTool)
	if err != nil {
		return "", err
	}

	err = cmdgit.RunCompare([]string{hashA, hashB})

	return tool.Name, err
}

// TriggerBranchDiff compares a path or whole trees across two branches.
func TriggerBranchDiff(repoRoot, branchA, branchB, filePath string) (string, error) {
	if strings.TrimSpace(filePath) == "" || filePath == "." {
		return triggerTreeCompare(repoRoot, branchA, branchB)
	}

	leftPath, err := extractBranchFile(repoRoot, branchA, filePath)
	if err != nil {
		return "", err
	}

	rightPath, err := extractBranchFile(repoRoot, branchB, filePath)
	if err != nil {
		return "", err
	}

	return launchPathsComparison(leftPath, rightPath)
}

func triggerTreeCompare(repoRoot, branchA, branchB string) (string, error) {
	tool, err := cmdgit.FindCompareTool()
	if err != nil {
		return "", err
	}

	err = cmdgit.RunCompare([]string{branchA, branchB})

	return tool.Name, err
}

func extractBranchFile(repoRoot, branch, filePath string) (string, error) {
	data, err := readGitFileBytes(repoRoot, branch+":"+filePath)
	if err != nil {
		return createTempDiffFile(sanitizeRefName(branch), filePath, nil)
	}

	return createTempDiffFile(sanitizeRefName(branch), filePath, data)
}

func launchPathsComparison(leftPath, rightPath string) (string, error) {
	tool, err := cmdgit.FindCompareTool()
	if err != nil {
		return "", err
	}

	err = cmdgit.LaunchToolWithPathsAsync(tool, []string{leftPath, rightPath})
	if err != nil {
		return tool.Name, err
	}

	return tool.Name, nil
}

func readGitFileBytes(repoRoot, gitPathSpec string) ([]byte, error) {
	cmd := exec.Command(constants.GitBin, "show", gitPathSpec)
	if repoRoot != "" {
		cmd.Dir = repoRoot
	}

	out, err := cmd.Output()
	if err != nil {
		return nil, apperror.WrapSimple(err, "git show "+gitPathSpec)
	}

	return out, nil
}

func createTempDiffFile(label, origPath string, data []byte) (string, error) {
	ext := filepath.Ext(origPath)
	base := strings.TrimSuffix(filepath.Base(origPath), ext)
	fileName := fmt.Sprintf("gitmap_diff_%s_%s%s", label, base, ext)
	tmpPath := filepath.Join(os.TempDir(), fileName)

	err := os.WriteFile(tmpPath, data, 0o600)
	if err != nil {
		return "", apperror.WrapSimple(err, "write temp diff file")
	}

	return tmpPath, nil
}

func sanitizeRefName(ref string) string {
	replacer := strings.NewReplacer("/", "_", "\\", "_", ":", "_", "*", "_")

	return replacer.Replace(ref)
}
