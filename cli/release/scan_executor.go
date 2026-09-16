package release

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
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
	isFound, err := isRefPresent(repoDir, "refs/heads/"+branchName)
	if err != nil {
		return apperror.WrapSimple(err, "processBranch")
	}

	if isFound {
		action.IsBranchSkipped = true

		return nil
	}

	return createBranch(repoDir, branchName, commit.Hash, action)
}

func processTag(repoDir string, commit ParsedCommit, action *ScanCommitAction) error {
	isFound, err := isRefPresent(repoDir, "refs/tags/"+commit.Version)
	if err != nil {
		return apperror.WrapSimple(err, "processTag")
	}

	if isFound {
		action.IsTagSkipped = true

		return nil
	}

	return createTag(repoDir, commit.Version, commit.Hash, action)
}
