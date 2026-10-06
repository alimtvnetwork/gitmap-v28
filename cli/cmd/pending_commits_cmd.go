// Package cmd — pending_commits_cmd.go inspects and classifies uncommitted changes and unpushed commits.
package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// gitCmdExecutor defines a function that executes a git command in a specific directory.
type gitCmdExecutor func(dir string, args ...string) (string, error)

var defaultPendingCommitsGitExecutor gitCmdExecutor = func(dir string, args ...string) (string, error) {
	cmdArgs := append([]string{"-C", dir}, args...)
	cmd := exec.Command("git", cmdArgs...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

var currentPendingCommitsGitExecutor = defaultPendingCommitsGitExecutor

// pendingCommitsSSHExecutor defines an interface for executing commands over SSH.
type pendingCommitsSSHExecutor interface {
	Execute(conn db.SSHConnection, cmd, shell string) (string, error)
}

type defaultPendingCommitsSSHRunner struct{}

func (r defaultPendingCommitsSSHRunner) Execute(conn db.SSHConnection, cmd, shell string) (string, error) {
	client, errConnect := cmdssh.ConnectSSHClientWithErr(conn)
	if errConnect != nil {
		return "", errConnect
	}
	defer client.Close()

	out, errExec := crypto.RunCommand(client, cmd, shell)
	if errExec != nil && shell == "bash" {
		out, errExec = crypto.RunCommand(client, cmd, "sh")
	}
	return out, errExec
}

var currentPendingCommitsSSHExecutor pendingCommitsSSHExecutor = defaultPendingCommitsSSHRunner{}

// runPendingCommits is the CLI dispatch entry point.
func runPendingCommits(args []string) error {
	return RunPendingCommits(args)
}

// RunPendingCommits executes discovery of uncommitted and unpushed changes.
func RunPendingCommits(args []string) error {
	if isPendingCommitsHelpRequested(args) {
		printPendingCommitsHelp()
		return nil
	}

	opts := parsePendingCommitsOptions(args)
	if opts.IsSSH {
		return runFleetSSHPendingCommits(opts)
	}

	return runLocalPendingCommits(opts)
}

func isPendingCommitsHelpRequested(args []string) bool {
	for _, a := range args {
		if a == "-h" || a == "--help" || a == "help" {
			return true
		}
	}
	return false
}

func parsePendingCommitsOptions(args []string) PendingCommitsOptions {
	opts := PendingCommitsOptions{
		SortMode:    "priority",
		DetailMode:  "summary",
		IsDirtyOnly: true,
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-s" || arg == "--sort":
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				opts.SortMode = strings.ToLower(args[i+1])
				i++
			}
		case strings.HasPrefix(arg, "--sort="):
			opts.SortMode = strings.ToLower(strings.TrimPrefix(arg, "--sort="))
		case arg == "-d" || arg == "--detail":
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				opts.DetailMode = strings.ToLower(args[i+1])
				i++
			} else {
				opts.DetailMode = "all"
			}
		case strings.HasPrefix(arg, "--detail="):
			opts.DetailMode = strings.ToLower(strings.TrimPrefix(arg, "--detail="))
		case arg == "-j" || arg == "--json":
			opts.IsJSON = true
		case arg == "--ssh" || arg == "-ssh":
			opts.IsSSH = true
		case arg == "--dirty-only":
			opts.IsDirtyOnly = true
			opts.IsAll = false
		case arg == "--all":
			opts.IsAll = true
			opts.IsDirtyOnly = false
		case !strings.HasPrefix(arg, "-"):
			opts.TargetRepo = arg
		}
	}

	if opts.SortMode != "name" && opts.SortMode != "count" && opts.SortMode != "priority" {
		opts.SortMode = "priority"
	}
	if opts.DetailMode == "1" {
		opts.DetailMode = "1-by-1"
	}

	return opts
}

func runLocalPendingCommits(opts PendingCommitsOptions) error {
	cwd, errCwd := os.Getwd()
	if errCwd != nil {
		cwd = "."
	}

	records := resolveWorkspaceRepositories(cwd)
	var recordsToInspect []model.ScanRecord

	if opts.TargetRepo != "" && !strings.EqualFold(opts.TargetRepo, "all") {
		for _, r := range records {
			if strings.EqualFold(r.RepoName, opts.TargetRepo) ||
				strings.EqualFold(r.Slug, opts.TargetRepo) ||
				strings.EqualFold(r.RelativePath, opts.TargetRepo) {
				recordsToInspect = append(recordsToInspect, r)
			}
		}
		if len(recordsToInspect) == 0 {
			return apperror.NewSimple(fmt.Sprintf("repository '%s' not found", opts.TargetRepo), "E9002")
		}
	} else {
		recordsToInspect = records
	}

	inspected := inspectRepositoriesPendingCommits(recordsToInspect)
	totalScanned := len(inspected)
	dirtyCount := 0
	totalUncommitted := 0
	totalUnpushed := 0

	for _, item := range inspected {
		if item.IsDirty {
			dirtyCount++
		}
		totalUncommitted += item.UntrackedFilesCount + item.ModifiedFilesCount + item.StagedFilesCount
		totalUnpushed += item.UnpushedCommitsCount
	}

	var displayed []RepoPendingCommitRecord
	for _, item := range inspected {
		if opts.IsDirtyOnly {
			if item.IsDirty || item.HasUnpushed {
				displayed = append(displayed, item)
			}
		} else {
			displayed = append(displayed, item)
		}
	}

	sortPendingCommits(displayed, opts.SortMode)

	payload := PendingCommitsPayload{
		Timestamp:             time.Now().UTC(),
		TotalReposScanned:     totalScanned,
		TotalDirtyRepos:       dirtyCount,
		TotalUncommittedFiles: totalUncommitted,
		TotalUnpushedCommits:  totalUnpushed,
		SortMode:              opts.SortMode,
		DetailMode:            opts.DetailMode,
		Repositories:          displayed,
	}

	if opts.IsJSON {
		return emitPendingCommitsJSON(payload)
	}

	renderPendingCommitsTerminal(payload, totalScanned, inspected)
	if opts.DetailMode == "all" {
		renderPendingCommitsDetailsAll(displayed)
	} else if opts.DetailMode == "1-by-1" {
		renderPendingCommitsInteractive(displayed)
	}

	return nil
}

func resolveWorkspaceRepositories(cwd string) []model.ScanRecord {
	targets := ResolvePullDirectoryTargets(cwd)
	if len(targets) > 0 {
		return targets
	}

	gitPath := filepath.Join(cwd, ".git")
	if _, err := os.Stat(gitPath); err == nil {
		base := filepath.Base(cwd)
		return []model.ScanRecord{
			{
				RepoName:     base,
				Slug:         base,
				AbsolutePath: cwd,
				RelativePath: ".",
			},
		}
	}

	return nil
}

func inspectRepositoriesPendingCommits(records []model.ScanRecord) []RepoPendingCommitRecord {
	results := make([]RepoPendingCommitRecord, 0, len(records))
	for _, rec := range records {
		results = append(results, inspectSingleRepoPendingCommits(rec))
	}
	return results
}

func inspectSingleRepoPendingCommits(rec model.ScanRecord) RepoPendingCommitRecord {
	dir := rec.AbsolutePath
	if dir == "" {
		dir = rec.RelativePath
	}

	porcelain, _ := currentPendingCommitsGitExecutor(dir, "status", "--porcelain")
	untracked, modified, staged, pendingFiles := parsePorcelainStatusLines(porcelain)

	branch, _ := currentPendingCommitsGitExecutor(dir, "rev-parse", "--abbrev-ref", "HEAD")
	branch = strings.TrimSpace(branch)
	if branch == "" {
		branch = rec.Branch
	}
	if branch == "" {
		branch = "main"
	}

	upstream, errUpstream := currentPendingCommitsGitExecutor(dir, "rev-parse", "--abbrev-ref", "@{u}")
	hasUpstream := errUpstream == nil && strings.TrimSpace(upstream) != "" && !strings.Contains(strings.ToLower(upstream), "fatal")

	unpushedCount := 0
	var unpushedSHAs []string
	if hasUpstream {
		revCount, errCount := currentPendingCommitsGitExecutor(dir, "rev-list", "--count", "@{u}..HEAD")
		if errCount == nil {
			unpushedCount, _ = strconv.Atoi(strings.TrimSpace(revCount))
		}
		if unpushedCount > 0 {
			revLog, errLog := currentPendingCommitsGitExecutor(dir, "log", "--oneline", "-n", "10", "@{u}..HEAD")
			if errLog == nil {
				unpushedSHAs = parsePendingCommitLines(revLog)
			}
		}
	}

	isDirty := (untracked + modified + staged) > 0
	hasUncommitted := isDirty
	hasUnpushed := unpushedCount > 0
	isClean := !isDirty && !hasUnpushed

	relPath := rec.RelativePath
	if relPath == "" {
		relPath = rec.RepoName
	}

	return RepoPendingCommitRecord{
		RepoName:             rec.RepoName,
		RelativePath:         relPath,
		CurrentBranch:        branch,
		IsDirty:              isDirty,
		IsClean:              isClean,
		HasUncommitted:       hasUncommitted,
		HasUnpushed:          hasUnpushed,
		HasUpstream:          hasUpstream,
		UntrackedFilesCount:  untracked,
		ModifiedFilesCount:   modified,
		StagedFilesCount:     staged,
		UnpushedCommitsCount: unpushedCount,
		PendingFiles:         pendingFiles,
		UnpushedCommitSHAs:   unpushedSHAs,
	}
}

func parsePorcelainStatusLines(output string) (int, int, int, []string) {
	var untracked, modified, staged int
	var pendingFiles []string

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimRight(line, "\r")
		if len(line) == 0 {
			continue
		}
		pendingFiles = append(pendingFiles, strings.TrimSpace(line))
		if strings.HasPrefix(line, "??") {
			untracked++
			continue
		}
		if len(line) >= 2 {
			idxChar := line[0]
			wtChar := line[1]
			if idxChar != ' ' && idxChar != '?' {
				staged++
			}
			if wtChar != ' ' && wtChar != '?' {
				modified++
			}
		}
	}

	return untracked, modified, staged, pendingFiles
}

func parsePendingCommitLines(output string) []string {
	var shas []string
	lines := strings.Split(output, "\n")
	for _, l := range lines {
		t := strings.TrimSpace(l)
		if t != "" {
			shas = append(shas, t)
		}
	}
	return shas
}

func sortPendingCommits(records []RepoPendingCommitRecord, sortMode string) {
	switch sortMode {
	case "count":
		sort.SliceStable(records, func(i, j int) bool {
			countI := records[i].UntrackedFilesCount + records[i].ModifiedFilesCount + records[i].StagedFilesCount + records[i].UnpushedCommitsCount
			countJ := records[j].UntrackedFilesCount + records[j].ModifiedFilesCount + records[j].StagedFilesCount + records[j].UnpushedCommitsCount
			if countI != countJ {
				return countI > countJ
			}
			return strings.ToLower(records[i].RepoName) < strings.ToLower(records[j].RepoName)
		})
	case "name":
		sort.SliceStable(records, func(i, j int) bool {
			return strings.ToLower(records[i].RepoName) < strings.ToLower(records[j].RepoName)
		})
	default: // "priority"
		sort.SliceStable(records, func(i, j int) bool {
			pI := calculatePriorityRank(records[i])
			pJ := calculatePriorityRank(records[j])
			if pI != pJ {
				return pI < pJ
			}
			return strings.ToLower(records[i].RepoName) < strings.ToLower(records[j].RepoName)
		})
	}
}

func calculatePriorityRank(r RepoPendingCommitRecord) int {
	if r.IsDirty && r.HasUnpushed {
		return 1
	}
	if r.IsDirty {
		return 2
	}
	if r.HasUnpushed {
		return 3
	}
	return 4
}

func emitPendingCommitsJSON(payload PendingCommitsPayload) error {
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal pending commits json")
	}
	fmt.Println(string(b))
	return nil
}

func renderPendingCommitsTerminal(payload PendingCommitsPayload, totalScanned int, allInspected []RepoPendingCommitRecord) {
	border := constants.ColorCyan
	reset := constants.ColorReset
	bold := constants.ColorBold

	fmt.Println()
	fmt.Printf("  %s┌──────────────────────────────────────────────────────────────────────────────┐%s\n", border, reset)
	fmt.Printf("  %s│%s%s                    GITMAP PENDING COMMITS SUMMARY                           %s%s│%s\n", border, reset, bold, reset, border, reset)
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)
	fmt.Printf("  %s│%s Scanned: %-8s │ Dirty: %-8s │ Uncommitted: %-6s │ Unpushed: %-8s %s│%s\n",
		border, reset,
		fmt.Sprintf("%d repos", payload.TotalReposScanned),
		fmt.Sprintf("%d repos", payload.TotalDirtyRepos),
		fmt.Sprintf("%d files", payload.TotalUncommittedFiles),
		fmt.Sprintf("%d", payload.TotalUnpushedCommits),
		border, reset)
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)
	fmt.Printf("  %s│%s %-20s %-8s %-5s %7s %6s %6s %9s %7s  %s│%s\n",
		border, reset,
		"REPOSITORY", "BRANCH", "DIRTY", "UNTRACK", "MODIF", "STAGE", "UNPUSHED", "STATUS",
		border, reset)
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)

	cleanCount := 0
	for _, rec := range allInspected {
		if rec.IsClean {
			cleanCount++
		}
	}

	for _, rec := range payload.Repositories {
		dirtyLabel := "NO "
		statusLabel := constants.ColorGreen + "○ CLEAN" + reset
		if rec.IsDirty {
			dirtyLabel = constants.ColorYellow + "YES" + reset
			statusLabel = constants.ColorYellow + "● PEND " + reset
		} else if rec.HasUnpushed {
			statusLabel = constants.ColorYellow + "● PEND " + reset
		}

		repoDisplay := rec.RepoName
		if len(repoDisplay) > 20 {
			repoDisplay = repoDisplay[:19] + "…"
		}
		branchDisplay := rec.CurrentBranch
		if len(branchDisplay) > 8 {
			branchDisplay = branchDisplay[:7] + "…"
		}

		fmt.Printf("  %s│%s %-20s %-8s %-5s %7d %6d %6d %9d   %-7s %s│%s\n",
			border, reset,
			repoDisplay,
			branchDisplay,
			dirtyLabel,
			rec.UntrackedFilesCount,
			rec.ModifiedFilesCount,
			rec.StagedFilesCount,
			rec.UnpushedCommitsCount,
			statusLabel,
			border, reset)
	}

	if payload.DetailMode == "summary" && cleanCount > 0 && len(payload.Repositories) < totalScanned {
		fmt.Printf("  %s│%s %-20s %-8s %-5s %7d %6d %6d %9d   %-7s %s│%s\n",
			border, reset,
			fmt.Sprintf("%d clean repos", cleanCount),
			"-",
			"NO ",
			0, 0, 0, 0,
			constants.ColorGreen+"○ CLEAN"+reset,
			border, reset)
	}

	fmt.Printf("  %s└──────────────────────────────────────────────────────────────────────────────┘%s\n", border, reset)
}

func renderPendingCommitsDetailsAll(records []RepoPendingCommitRecord) {
	fmt.Println()
	fmt.Printf("  %sDetailed Pending Breakdown:%s\n", constants.ColorBold, constants.ColorReset)
	for _, rec := range records {
		if rec.IsClean {
			continue
		}
		fmt.Printf("\n  %s● %s%s (%s)\n", constants.ColorCyan, rec.RepoName, constants.ColorReset, rec.CurrentBranch)
		if len(rec.PendingFiles) > 0 {
			fmt.Printf("    %sModified Files (%d):%s\n", constants.ColorYellow, len(rec.PendingFiles), constants.ColorReset)
			for _, f := range rec.PendingFiles {
				fmt.Printf("      %s\n", f)
			}
		}
		if len(rec.UnpushedCommitSHAs) > 0 {
			fmt.Printf("    %sUnpushed Commits (%d):%s\n", constants.ColorYellow, len(rec.UnpushedCommitSHAs), constants.ColorReset)
			for _, sha := range rec.UnpushedCommitSHAs {
				fmt.Printf("      %s\n", sha)
			}
		}
	}
}

func renderPendingCommitsInteractive(records []RepoPendingCommitRecord) {
	var dirtyList []RepoPendingCommitRecord
	for _, r := range records {
		if !r.IsClean {
			dirtyList = append(dirtyList, r)
		}
	}

	if len(dirtyList) == 0 {
		fmt.Println("  All repositories are clean. Nothing to step through.")
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	for idx, rec := range dirtyList {
		fmt.Printf("\n  [%d/%d] Repository: %s%s%s (branch: %s)\n",
			idx+1, len(dirtyList), constants.ColorCyan, rec.RepoName, constants.ColorReset, rec.CurrentBranch)
		if len(rec.PendingFiles) > 0 {
			fmt.Printf("    Files (%d):\n", len(rec.PendingFiles))
			for _, f := range rec.PendingFiles {
				fmt.Printf("      %s\n", f)
			}
		}
		if len(rec.UnpushedCommitSHAs) > 0 {
			fmt.Printf("    Unpushed Commits (%d):\n", len(rec.UnpushedCommitSHAs))
			for _, sha := range rec.UnpushedCommitSHAs {
				fmt.Printf("      %s\n", sha)
			}
		}

		if idx < len(dirtyList)-1 {
			fmt.Printf("  Press [Enter] for next, or 'q' to quit: ")
			if !scanner.Scan() {
				break
			}
			input := strings.TrimSpace(scanner.Text())
			if strings.EqualFold(input, "q") || strings.EqualFold(input, "quit") {
				break
			}
		}
	}
}

func runFleetSSHPendingCommits(opts PendingCommitsOptions) error {
	conns, errFetch := cmdssh.FetchAllSSHConnections()
	if errFetch != nil {
		return apperror.WrapSimple(errFetch, "fetch fleet ssh connections")
	}

	if len(conns) == 0 {
		if opts.IsJSON {
			return emitPendingCommitsFleetJSON(NodesPendingCommitsPayload{
				Timestamp:    time.Now().UTC(),
				TotalNodes:   0,
				OnlineNodes:  0,
				FleetResults: []NodePendingCommitsRecord{},
			})
		}
		fmt.Printf("  %sNo registered SSH cluster nodes found.%s\n", constants.ColorYellow, constants.ColorReset)
		return nil
	}

	results := make([]NodePendingCommitsRecord, len(conns))
	var wg sync.WaitGroup

	for i, conn := range conns {
		wg.Add(1)
		go func(idx int, c db.SSHConnection) {
			defer wg.Done()
			results[idx] = querySingleNodePendingCommits(c, opts)
		}(i, conn)
	}

	wg.Wait()

	onlineCount := 0
	for _, r := range results {
		if r.IsOnline {
			onlineCount++
		}
	}

	fleetPayload := NodesPendingCommitsPayload{
		Timestamp:    time.Now().UTC(),
		TotalNodes:   len(conns),
		OnlineNodes:  onlineCount,
		FleetResults: results,
	}

	if opts.IsJSON {
		return emitPendingCommitsFleetJSON(fleetPayload)
	}

	renderPendingCommitsFleetTerminal(fleetPayload)
	return nil
}

func querySingleNodePendingCommits(conn db.SSHConnection, opts PendingCommitsOptions) NodePendingCommitsRecord {
	start := time.Now()
	rec := NodePendingCommitsRecord{
		NodeAlias: conn.Alias,
		Host:      conn.IPAddress,
		OSType:    conn.OS,
	}

	isWin := strings.EqualFold(conn.OS, "windows")
	cmdStr := "gitmap pending-commits --json"
	if opts.SortMode != "" && opts.SortMode != "priority" {
		cmdStr += " --sort=" + opts.SortMode
	}
	if opts.DetailMode == "all" {
		cmdStr += " --detail=all"
	}
	if opts.IsAll {
		cmdStr += " --all"
	}

	shell := "bash"
	if isWin {
		cmdStr = fmt.Sprintf(`powershell.exe -NoProfile -Command "%s"`, cmdStr)
		shell = "ps"
	}

	out, errExec := currentPendingCommitsSSHExecutor.Execute(conn, cmdStr, shell)
	rec.LatencyMs = time.Since(start).Milliseconds()

	if errExec != nil {
		rec.ErrorMessage = errExec.Error()
		if isOfflineErrorString(errExec.Error()) {
			rec.IsOnline = false
			rec.IsSuccess = false
			return rec
		}
		rec.IsOnline = true
		rec.IsSuccess = false
		return rec
	}

	rec.IsOnline = true
	rec.IsSuccess = true

	cleanJSON := extractJSONPayload(out)
	var nodePayload PendingCommitsPayload
	if errUnmarshal := json.Unmarshal([]byte(cleanJSON), &nodePayload); errUnmarshal == nil {
		rec.Payload = &nodePayload
	}

	return rec
}

func isOfflineErrorString(errStr string) bool {
	low := strings.ToLower(errStr)
	return strings.Contains(low, "network unreachable") ||
		strings.Contains(low, "connection refused") ||
		strings.Contains(low, "timeout") ||
		strings.Contains(low, "i/o timeout") ||
		strings.Contains(low, "offline")
}

func extractJSONPayload(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start >= 0 && end > start {
		return s[start : end+1]
	}
	return s
}

func emitPendingCommitsFleetJSON(payload NodesPendingCommitsPayload) error {
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal fleet pending commits json")
	}
	fmt.Println(string(b))
	return nil
}

func renderPendingCommitsFleetTerminal(fleet NodesPendingCommitsPayload) {
	border := constants.ColorCyan
	reset := constants.ColorReset
	bold := constants.ColorBold

	fmt.Println()
	fmt.Printf("  %s┌──────────────────────────────────────────────────────────────────────────────┐%s\n", border, reset)
	fmt.Printf("  %s│%s%s               GITMAP FLEET PENDING COMMITS (SSH CLUSTER)                     %s%s│%s\n", border, reset, bold, reset, border, reset)
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)
	fmt.Printf("  %s│%s Total Nodes: %-6d │ Online Nodes: %-6d                                   %s│%s\n",
		border, reset, fleet.TotalNodes, fleet.OnlineNodes, border, reset)
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)
	fmt.Printf("  %s│%s %-16s %-15s %-7s %-7s %-8s %-10s %7s %s│%s\n",
		border, reset,
		"NODE", "HOST", "STATUS", "DIRTY", "UNCOMMIT", "UNPUSHED", "LATENCY",
		border, reset)
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)

	for _, node := range fleet.FleetResults {
		statusLabel := constants.ColorRed + "OFFLINE" + reset
		dirtyStr := "-"
		uncommitStr := "-"
		unpushedStr := "-"

		if node.IsOnline {
			if node.IsSuccess && node.Payload != nil {
				statusLabel = constants.ColorGreen + "ONLINE " + reset
				dirtyStr = fmt.Sprintf("%d", node.Payload.TotalDirtyRepos)
				uncommitStr = fmt.Sprintf("%d", node.Payload.TotalUncommittedFiles)
				unpushedStr = fmt.Sprintf("%d", node.Payload.TotalUnpushedCommits)
			} else {
				statusLabel = constants.ColorYellow + "ERROR  " + reset
			}
		}

		fmt.Printf("  %s│%s %-16s %-15s %-7s %-7s %-8s %-10s %5dms %s│%s\n",
			border, reset,
			node.NodeAlias,
			node.Host,
			statusLabel,
			dirtyStr,
			uncommitStr,
			unpushedStr,
			node.LatencyMs,
			border, reset)
	}

	fmt.Printf("  %s└──────────────────────────────────────────────────────────────────────────────┘%s\n", border, reset)
}

func printPendingCommitsHelp() {
	cyan := constants.ColorCyan
	reset := constants.ColorReset
	bold := constants.ColorBold
	yellow := constants.ColorYellow

	fmt.Printf(`
%s┌──────────────────────────────────────────────────────────────────────────────┐
│                      GITMAP PENDING COMMITS HELP                             │
└──────────────────────────────────────────────────────────────────────────────┘%s

%sDESCRIPTION:%s
  Scans local workspace repositories or fleet cluster nodes over SSH to identify
  uncommitted working tree modifications and unpushed commits ahead of origin.

%sUSAGE:%s
  gitmap pending-commits [repoName|all] [flags]
  gitmap pc [repoName|all] [flags]

%sFLAGS:%s
  -s, --sort <mode>       Sort order: priority | name | count (default: priority)
  -d, --detail <mode>     Detail view: summary | all | 1 | none (default: summary)
      --ssh, -ssh         Aggregate pending commits across all registered cluster SSH nodes
  -j, --json              Output machine-readable JSON telemetry
      --dirty-only        Display only repositories with changes (default: true)
      --all               Include completely clean repositories in output
  -h, --help              Display this help menu

%sEXAMPLES:%s
  %sgitmap pc%s
      Summary overview of all dirty repositories in workspace.

  %sgitmap pc --sort=count --detail=all%s
      Detailed file diff breakdown sorted by total modification count.

  %sgitmap pc -ssh --json%s
      Aggregate pending commits across entire SSH fleet as JSON.
`,
		cyan, reset,
		bold, reset,
		bold, reset,
		bold, reset,
		bold, reset,
		yellow, reset,
		yellow, reset,
		yellow, reset,
	)
}
