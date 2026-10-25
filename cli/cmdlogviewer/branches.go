package cmdlogviewer

import (
	"strings"
)

// QueryBranches returns a list of local branches and the currently checked-out branch.
func QueryBranches(repoRoot string) (BranchesResponse, error) {
	branches, err := listLocalBranches(repoRoot)
	if err != nil {
		return BranchesResponse{}, err
	}

	active := queryActiveBranch(repoRoot)

	return BranchesResponse{Branches: branches, ActiveBranch: active}, nil
}

func listLocalBranches(repoRoot string) ([]string, error) {
	out, err := runGitInDir(repoRoot, "for-each-ref", "--format=%(refname:short)", "refs/heads/")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(out, "\n")
	branches := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			branches = append(branches, trimmed)
		}
	}

	return branches, nil
}

func queryActiveBranch(repoRoot string) string {
	out, err := runGitInDir(repoRoot, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "main"
	}

	return strings.TrimSpace(out)
}

// ValidateBranch reports whether the branch name resolves to a valid git commit.
func ValidateBranch(repoRoot, branch string) bool {
	if strings.TrimSpace(branch) == "" {
		return false
	}

	_, err := runGitInDir(repoRoot, "rev-parse", "--verify", branch+"^{commit}")

	return err == nil
}
