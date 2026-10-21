package gitutil

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// GenerateCaseNormalizationRecipe builds a remediation recipe to resolve case collision in the index.
func GenerateCaseNormalizationRecipe(repoPath string, col CaseCollision) RemediationRecipe {
	rawPath := CleanRepoPathRaw(repoPath)
	p := CleanRepoPath(repoPath)
	steps := buildCaseNormalizationSteps(rawPath, col)
	cmdStr := buildCaseNormalizationCommand(p, col)

	return RemediationRecipe{
		Title:       fmt.Sprintf("Normalize Case Collision (%s)", col.CanonicalPath),
		Description: fmt.Sprintf("Purge uppercase index collision and stage canonical %s", col.CanonicalPath),
		Command:     cmdStr,
		Steps:       steps,
	}
}

func buildCaseCommitMessage(canonicalPath string) string {
	return fmt.Sprintf("fix(case): normalize %s and purge uppercase index collision", canonicalPath)
}

func appendRmCachedSteps(steps []RemediationStep, rawPath string, col CaseCollision) []RemediationStep {
	for _, badPath := range col.CollidingPaths {
		if badPath != col.CanonicalPath {
			step := RemediationStep{Name: "git", Args: []string{"-C", rawPath, "rm", "--cached", "-f", badPath}}
			steps = append(steps, step)
		}
	}

	return steps
}

func buildCaseNormalizationSteps(rawPath string, col CaseCollision) []RemediationStep {
	var steps []RemediationStep
	steps = appendRmCachedSteps(steps, rawPath, col)
	steps = append(steps, RemediationStep{Name: "git", Args: []string{"-C", rawPath, "add", "--", col.CanonicalPath}})
	steps = append(steps, RemediationStep{Name: "git", Args: []string{"-C", rawPath, "commit", "-m", buildCaseCommitMessage(col.CanonicalPath)}})

	return steps
}

func buildCaseNormalizationCommand(p string, col CaseCollision) string {
	var cmds []string
	for _, badPath := range col.CollidingPaths {
		if badPath != col.CanonicalPath {
			cmds = append(cmds, fmt.Sprintf("git -C %s rm --cached -f %q", p, badPath))
		}
	}
	cmds = append(cmds, fmt.Sprintf("git -C %s add -- %q", p, col.CanonicalPath))
	msg := buildCaseCommitMessage(col.CanonicalPath)
	cmds = append(cmds, fmt.Sprintf("git -C %s commit -m %q", p, msg))

	return strings.Join(cmds, " && ")
}

// ExecuteCaseNormalization executes atomic index purge, stages canonical path, and commits.
func ExecuteCaseNormalization(repoPath string, col CaseCollision) error {
	diskPath := filepath.Join(repoPath, col.CanonicalPath)
	buf, hasBackup := readDiskFileBuffer(diskPath)

	if err := purgeCollidingIndexEntries(repoPath, col); err != nil {
		return apperror.WrapSimple(err, "purge case collision")
	}
	restoreDiskFileIfNeeded(diskPath, buf, hasBackup)

	return stageAndCommitCanonical(repoPath, col.CanonicalPath)
}

func readDiskFileBuffer(filePath string) ([]byte, bool) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, false
	}

	return data, true
}

func restoreDiskFileIfNeeded(filePath string, buf []byte, hasBackup bool) {
	if !hasBackup {
		return
	}
	_, err := os.Stat(filePath)
	if err != nil {
		_ = os.WriteFile(filePath, buf, 0644)
	}
}

func purgeCollidingIndexEntries(repoPath string, col CaseCollision) error {
	for _, badPath := range col.CollidingPaths {
		if err := purgeSingleCollidingPath(repoPath, col.CanonicalPath, badPath); err != nil {
			return err
		}
	}

	return nil
}

func purgeSingleCollidingPath(repoPath, canonicalPath, badPath string) error {
	if badPath == canonicalPath {
		return nil
	}

	cmd := exec.Command("git", "-C", repoPath, "rm", "--cached", "-f", badPath)
	cmd.Env = BuildSafeGitEnv()

	return cmd.Run()
}

func stageAndCommitCanonical(repoPath, canonicalPath string) error {
	addCmd := exec.Command("git", "-C", repoPath, "add", "--", canonicalPath)
	addCmd.Env = BuildSafeGitEnv()
	if err := addCmd.Run(); err != nil {
		return apperror.WrapSimple(err, "git add canonical path")
	}

	msg := buildCaseCommitMessage(canonicalPath)
	commitCmd := exec.Command("git", "-C", repoPath, "commit", "-m", msg)
	commitCmd.Env = BuildSafeGitEnv()
	if err := commitCmd.Run(); err != nil {
		return apperror.WrapSimple(err, "git commit canonical path")
	}

	return nil
}
