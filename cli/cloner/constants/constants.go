package constants

// Static immutable git CLI arguments and environment variables for cloner.
var (
	// GitPullFastForwardArgs defines args for safe fast-forward git pull.
	GitPullFastForwardArgs = []string{"pull", "--progress", "--ff-only", "--autostash"}

	// GitPullAutoMergeArgs defines args for auto-merge fallback git pull.
	GitPullAutoMergeArgs = []string{"pull", "--progress", "--no-rebase", "--no-edit", "-X", "diff-algorithm=histogram", "--autostash"}

	// GitStatusPorcelainArgs defines args for checking working tree porcelain status.
	GitStatusPorcelainArgs = []string{"status", "--porcelain"}
)

const (
	// EnvGitLFSNoSmudge skips LFS smudge filter during git operations.
	EnvGitLFSNoSmudge = "GIT_LFS_SKIP_SMUDGE=1"

	// EnvGitSSHCommand configures SSH with accept-new host key checking.
	EnvGitSSHCommand = "GIT_SSH_COMMAND=ssh -o StrictHostKeyChecking=accept-new"
)
