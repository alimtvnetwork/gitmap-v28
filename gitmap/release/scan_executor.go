package release

import (
	"os/exec"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

type ParsedCommit struct {
	Hash    string
	Message string
	Version string
}

type ScanCommitAction struct {
	CommitHash      string
	Version         string
	IsBranchCreated bool
	IsBranchSkipped bool
	IsTagCreated    bool
	IsTagSkipped    bool
}

func ExecuteCommitActions(repoDir string, commits []ParsedCommit) ([]ScanCommitAction, error) {
	var actions []ScanCommitAction
	for _, commit := range commits {
		action, err := processCommit(repoDir, commit)
		if err != nil {
			return nil, apperror.Wrap(err, "ExecuteCommitActions", map[string]any{"hash": commit.Hash})
		}
		actions = append(actions, action)
	}
	return actions, nil
}

func processCommit(repoDir string, commit ParsedCommit) (ScanCommitAction, error) {
	action := ScanCommitAction{
		CommitHash: commit.Hash,
		Version:    commit.Version,
	}
	if err := processBranch(repoDir, commit, &action); err != nil {
		return action, apperror.WrapSimple(err, "processCommit")
	}
	if err := processTag(repoDir, commit, &action); err != nil {
		return action, apperror.WrapSimple(err, "processCommit")
	}
	return action, nil
}

func processBranch(repoDir string, commit ParsedCommit, action *ScanCommitAction) error {
	branchName := "release/" + commit.Version
	isExists, err := isRefExists(repoDir, "refs/heads/"+branchName)
	if err != nil {
		return apperror.WrapSimple(err, "processBranch")
	}
	if isExists {
		action.IsBranchSkipped = true
		return nil
	}
	return createBranch(repoDir, branchName, commit.Hash, action)
}

func createBranch(repoDir, branchName, hash string, action *ScanCommitAction) error {
	cmd := exec.Command("git", "branch", branchName, hash)
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		return apperror.Wrap(err, "createBranch", map[string]any{"branch": branchName})
	}
	action.IsBranchCreated = true
	return nil
}

func processTag(repoDir string, commit ParsedCommit, action *ScanCommitAction) error {
	isExists, err := isRefExists(repoDir, "refs/tags/"+commit.Version)
	if err != nil {
		return apperror.WrapSimple(err, "processTag")
	}
	if isExists {
		action.IsTagSkipped = true
		return nil
	}
	return createTag(repoDir, commit.Version, commit.Hash, action)
}

func createTag(repoDir, tagName, hash string, action *ScanCommitAction) error {
	cmd := exec.Command("git", "tag", tagName, hash)
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		return apperror.Wrap(err, "createTag", map[string]any{"tag": tagName})
	}
	action.IsTagCreated = true
	return nil
}

func isRefExists(repoDir, refPath string) (bool, error) {
	cmd := exec.Command("git", "show-ref", "--verify", "--quiet", refPath)
	cmd.Dir = repoDir
	err := cmd.Run()
	if err == nil {
		return true, nil
	}
	if _, isExit := err.(*exec.ExitError); isExit {
		return false, nil
	}
	return false, apperror.Wrap(err, "isRefExists", map[string]any{"ref": refPath})
}
