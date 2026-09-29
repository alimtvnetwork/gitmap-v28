// Package cmdagy — agy_wpr_deploy.go handles cross-machine database transfer and remote task scheduling.
package cmdagy

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunWPRDeploy executes parallel deployment of WPR database and task queue across SSH nodes.
func RunWPRDeploy(args []string, opts WPROptions) error {
	targetAlias, targetProject := parseDeployArgs(args)
	if isDeployHelpNeeded(targetAlias, targetProject) {
		printDeploySuggestions()
		return nil
	}

	nodes, err := resolveTargetDeployNodes(targetAlias)
	if err != nil {
		printDeploySuggestions()
		return err
	}

	summary := executeAsyncFleetDeployment(nodes, targetProject, opts)
	if opts.IsJSON {
		return printJSON(summary)
	}

	RenderWPRDeploySummary(summary)
	return nil
}

func parseDeployArgs(args []string) (string, string) {
	alias := ""
	proj := ""
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		alias = args[0]
	}
	if len(args) > 1 && !strings.HasPrefix(args[1], "-") {
		proj = args[1]
	}
	return alias, proj
}

func isDeployHelpNeeded(alias, proj string) bool {
	return alias == "" || alias == "help" || proj == ""
}

func printDeploySuggestions() {
	fmt.Printf("\n  %sℹ GitMap WPR Deploy Usage:%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("    gitmap wpr deploy <alias|all> <project-name|all>")
	fmt.Println()
	fmt.Println("  Available Fleet Node Aliases:")
	for _, m := range LoadWPRMachineCache() {
		fmt.Printf("    • %-12s (%s | %s)\n", m.Alias, m.MachineName, m.IPAddress)
	}
	fmt.Println()
	fmt.Println("  Available Projects:")
	projects, _ := getAllProjects()
	for i, p := range filterNonRestrictedProjects(projects) {
		if i < 8 {
			fmt.Printf("    • %s\n", p.Name)
		}
	}
	fmt.Println()
}

func resolveTargetDeployNodes(alias string) ([]db.SSHConnection, error) {
	if SSHConnectionsFetcher == nil {
		return nil, apperror.NewNotFound("ssh_fleet", "E404", "SSH connection fetcher not initialized")
	}

	conns, err := SSHConnectionsFetcher()
	if err != nil || len(conns) == 0 {
		return nil, apperror.NewNotFound("ssh_fleet", "E404", "no remote SSH connections registered")
	}

	if alias == "all" || alias == "*" {
		return conns, nil
	}

	for _, c := range conns {
		if strings.EqualFold(c.Alias, alias) || c.IPAddress == alias {
			return []db.SSHConnection{c}, nil
		}
	}

	return nil, apperror.NewNotFound("node_alias", "E404", fmt.Sprintf("target node %q not found", alias))
}

func executeAsyncFleetDeployment(nodes []db.SSHConnection, targetProject string, opts WPROptions) WPRDeploySummary {
	start := time.Now()
	var wg sync.WaitGroup
	var mu sync.Mutex
	var results []WPRDeployResult

	for _, n := range nodes {
		wg.Add(1)
		go func(conn db.SSHConnection) {
			defer wg.Done()
			res := deploySingleNodeWPR(conn, targetProject, opts)
			mu.Lock()
			results = append(results, res)
			mu.Unlock()
		}(n)
	}

	wg.Wait()
	successCount := countSuccessfulDeploys(results)
	return WPRDeploySummary{
		Status:        "completed",
		TargetAlias:   opts.Target,
		TargetProject: targetProject,
		TotalNodes:    len(nodes),
		SuccessCount:  successCount,
		Results:       results,
		DurationMs:    time.Since(start).Milliseconds(),
	}
}

func countSuccessfulDeploys(results []WPRDeployResult) int {
	count := 0
	for _, r := range results {
		if r.IsCopied && r.IsTaskEnqueued {
			count++
		}
	}
	return count
}

func deploySingleNodeWPR(c db.SSHConnection, targetProject string, opts WPROptions) WPRDeployResult {
	start := time.Now()
	slug := store.SanitizeSlug(targetProject)
	srcDb := store.ResolveWatchPromptsDbPath(slug)
	ensureLocalSnapshotReady(srcDb, slug)

	client, err := dialSSHNodeClient(c)
	if err != nil {
		return buildDeployErrorResult(c, slug, err, start)
	}
	defer client.Close()

	remoteDest := resolveRemoteDbPath(c, slug)
	streamErr := streamFileToRemote(client, remoteDest, srcDb, isWindowsOS(c.OS))
	if streamErr != nil {
		return buildDeployErrorResult(c, slug, streamErr, start)
	}

	remoteStatus := enqueueRemoteWPRTask(client, c, slug)
	return buildDeploySuccessResult(c, slug, remoteDest, remoteStatus, start)
}

func ensureLocalSnapshotReady(srcDb, slug string) {
	if !isFileExisting(srcDb) {
		_, _ = SnapshotAllRunningPrompts()
	}
}

func resolveRemoteDbPath(c db.SSHConnection, slug string) string {
	if isWindowsOS(c.OS) {
		user := resolveNodeUser(c)
		return fmt.Sprintf("C:/Users/%s/AppData/Local/gitmap-cli/watch-prompts/%s/sql.db", user, slug)
	}

	return fmt.Sprintf("/tmp/gitmap/watch-prompts/%s/sql.db", slug)
}

func enqueueRemoteWPRTask(client *ssh.Client, c db.SSHConnection, slug string) bool {
	cmdStr := fmt.Sprintf("gitmap agy wpr start %s --once --json", slug)
	out, err := crypto.RunCommand(client, cmdStr, resolveNodeShell(c.OS))
	if err != nil {
		return false
	}

	return len(out) > 0
}

func streamFileToRemote(client *ssh.Client, destPath, localPath string, isWin bool) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}

	var stderrBuf bytes.Buffer
	session.Stderr = &stderrBuf
	remoteCmd := buildStreamReceiverCmd(destPath, isWin)
	file, openErr := os.Open(localPath)
	if openErr != nil {
		return openErr
	}
	defer file.Close()

	go func() {
		defer stdin.Close()
		_, _ = io.Copy(stdin, file)
	}()

	return session.Run(remoteCmd)
}

func buildDeployErrorResult(c db.SSHConnection, slug string, err error, start time.Time) WPRDeployResult {
	return WPRDeployResult{
		NodeAlias:    c.Alias,
		MachineName:  c.Alias,
		IPAddress:    c.IPAddress,
		RepoSlug:     slug,
		IsCopied:     false,
		ErrorMessage: err.Error(),
		DurationMs:   time.Since(start).Milliseconds(),
	}
}

func buildDeploySuccessResult(c db.SSHConnection, slug, dest string, isEnqueued bool, start time.Time) WPRDeployResult {
	return WPRDeployResult{
		NodeAlias:      c.Alias,
		MachineName:    c.Alias,
		IPAddress:      c.IPAddress,
		RepoSlug:       slug,
		RemotePath:     dest,
		IsCopied:       true,
		IsTaskEnqueued: isEnqueued,
		IsRunning:      true,
		DurationMs:     time.Since(start).Milliseconds(),
	}
}

func isWindowsOS(osName string) bool {
	lower := strings.ToLower(osName)
	return strings.Contains(lower, "win")
}

func resolveNodeUser(c db.SSHConnection) string {
	clean := strings.TrimSpace(c.Username)
	if len(clean) > 0 {
		return clean
	}

	return "Administrator"
}

func isFileExisting(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

func buildStreamReceiverCmd(destPath string, isWin bool) string {
	if isWin {
		dir := filepath.ToSlash(filepath.Dir(destPath))
		return fmt.Sprintf("cmd.exe /c if not exist \"%s\" mkdir \"%s\" && powershell.exe -NoProfile -Command \"$in=[System.Console]::OpenStandardInput(); $out=[System.IO.File]::Create('%s'); $in.CopyTo($out); $out.Close()\"", dir, dir, destPath)
	}

	dir := filepath.Dir(destPath)
	return fmt.Sprintf("mkdir -p '%s' && cat > '%s'", dir, destPath)
}
