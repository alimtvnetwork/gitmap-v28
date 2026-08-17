package cluster

import (
	"context"
)

// ExecGitPull runs gitmap pull --all on the remote node via shell.
func ExecGitPull(ctx context.Context, node ClusterNode) (string, string, int, error) {
	return ExecCmd(ctx, node, "gitmap pull --all")
}

// ExecGitPush runs gitmap push --all on the remote node via shell.
func ExecGitPush(ctx context.Context, node ClusterNode) (string, string, int, error) {
	return ExecCmd(ctx, node, "gitmap push --all")
}

// ExecGitCommit runs gitmap commit --all on the remote node via shell.
func ExecGitCommit(ctx context.Context, node ClusterNode) (string, string, int, error) {
	return ExecCmd(ctx, node, "gitmap commit --all")
}

// ExecGitStatus runs gitmap status --all on the remote node via shell.
func ExecGitStatus(ctx context.Context, node ClusterNode) (string, string, int, error) {
	return ExecCmd(ctx, node, "gitmap status --all")
}
