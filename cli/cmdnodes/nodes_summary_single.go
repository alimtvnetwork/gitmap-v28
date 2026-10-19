package cmdnodes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsummary"
	"github.com/pterm/pterm"
)

// RunNodesSummary delegates a single repository release summary across the fleet.
func RunNodesSummary(args []string) error {
	if len(args) == 0 {
		return cmdsummary.RunSummary(args)
	}

	targetRepo := args[0]
	if strings.HasPrefix(targetRepo, "-") {
		return cmdsummary.RunSummary(args)
	}

	cwd, _ := os.Getwd()
	localRecords := resolveWorkspaceReposHook(cwd)
	isLocal := false
	var localPath string

	if _, err := os.Stat(filepath.Join(targetRepo, ".git")); err == nil {
		isLocal = true
		localPath = targetRepo
	} else if _, err := os.Stat(filepath.Join(cwd, targetRepo, ".git")); err == nil {
		isLocal = true
		localPath = filepath.Join(cwd, targetRepo)
	} else {
		targetNorm := normalizeRepoKey(targetRepo)
		for _, r := range localRecords {
			if normalizeRepoKey(r.RepoName) == targetNorm ||
				normalizeRepoKey(r.Slug) == targetNorm ||
				normalizeRepoKey(filepath.Base(r.AbsolutePath)) == targetNorm ||
				normalizeRepoKey(r.HTTPSUrl) == targetNorm ||
				normalizeRepoKey(r.SSHUrl) == targetNorm ||
				normalizeRepoKey(r.DiscoveredURL) == targetNorm {
				isLocal = true
				if r.AbsolutePath != "" {
					localPath = r.AbsolutePath
				} else {
					localPath = r.RelativePath
				}
				break
			}
		}
	}

	if isLocal {
		localArgs := make([]string, len(args))
		copy(localArgs, args)
		if localPath != "" && targetRepo != localPath {
			localArgs[0] = localPath
		}
		return cmdsummary.RunSummary(localArgs)
	}

	conns, errConns := fetchAllSSHConnectionsHook()
	if errConns != nil || len(conns) == 0 {
		return cmdsummary.RunSummary(args)
	}

	pterm.Info.Printf("[GitMap Nodes] Querying fleet for repository '%s'...\n", targetRepo)
	for _, conn := range conns {
		cmd := fmt.Sprintf("gitmap summary %s", strings.Join(args, " "))
		shell := "bash"
		if isWindowsNode(conn) {
			cmd = fmt.Sprintf("gitmap.exe summary %s", strings.Join(args, " "))
			shell = "powershell"
		}
		out, errExec := currentSSHExecutor.Execute(conn, cmd, shell)
		if errExec == nil && !strings.Contains(out, "not a valid git repository") && !strings.Contains(out, "not found") {
			pterm.Success.Printf("[Node: %s (%s)]\n", conn.Alias, conn.IPAddress)
			fmt.Println(out)
			return nil
		}
	}

	return apperror.NewSimple(fmt.Sprintf("repository '%s' not found across local or fleet nodes", targetRepo), "E9080")
}
