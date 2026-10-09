// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RunNodesPendingCommits executes cluster-wide pending commits discovery across fleet nodes.
func RunNodesPendingCommits(args []string) error {
	if isPendingCommitsHelp(args) {
		return printNodesPendingCommitsHelp()
	}

	opts := parseNodesPendingCommitsOptions(args)
	conns, errFetch := cmdssh.FetchAllSSHConnections()
	if errFetch != nil {
		return apperror.WrapSimple(errFetch, "fetch fleet ssh connections")
	}

	if len(conns) == 0 {
		return apperror.NewSimple("no registered fleet SSH nodes found", "E9070")
	}

	targets := FilterFleetNodes(conns, opts.FilterOpts)
	if len(targets) == 0 {
		return apperror.NewSimple("no candidate fleet nodes found matching filter criteria", "E9071")
	}

	results := runFleetPendingCommitsParallel(targets, opts)
	if opts.IsJSON {
		return emitNodesPendingCommitsJSON(results, opts, len(targets))
	}

	renderNodesPendingCommitsTable(results)
	if opts.IsDetail {
		renderPendingCommitsDetailView(results)
	}

	printPendingCommitsSummary(results, len(targets))
	return nil
}

func parseNodesPendingCommitsOptions(args []string) NodesPendingCommitsOptions {
	opts := NodesPendingCommitsOptions{
		FilterOpts: ParseNodeFilterOptions(args),
		SortMode:   "priority",
	}

	for i := 0; i < len(args); i++ {
		step := parseSinglePendingFlag(&opts, args, i)
		if step > 0 {
			i += (step - 1)
		}
	}

	return opts
}

func parseSinglePendingFlag(opts *NodesPendingCommitsOptions, args []string, i int) int {
	arg := args[i]
	if arg == "--json" || arg == "-j" {
		opts.IsJSON = true
		return 1
	}

	if arg == "--detail" || arg == "-d" {
		opts.IsDetail = true
		return 1
	}

	if arg == "--dirty-only" {
		opts.IsDirtyOnly = true
		return 1
	}

	return parsePendingSortFlag(opts, args, i)
}

func parsePendingSortFlag(opts *NodesPendingCommitsOptions, args []string, i int) int {
	arg := args[i]
	if strings.HasPrefix(arg, "--sort=") {
		opts.SortMode = strings.TrimPrefix(arg, "--sort=")
		return 1
	}

	if strings.HasPrefix(arg, "-s=") {
		opts.SortMode = strings.TrimPrefix(arg, "-s=")
		return 1
	}

	if (arg == "--sort" || arg == "-s") && i+1 < len(args) {
		opts.SortMode = args[i+1]
		return 2
	}

	return 0
}

func runFleetPendingCommitsParallel(targets []db.SSHConnection, opts NodesPendingCommitsOptions) []NodePendingCommitsResult {
	results := make([]NodePendingCommitsResult, len(targets))
	var wg sync.WaitGroup

	for i, conn := range targets {
		wg.Add(1)
		go func(idx int, c db.SSHConnection) {
			defer wg.Done()
			results[idx] = executeSingleNodePendingCommits(c, opts)
		}(i, conn)
	}

	wg.Wait()
	return results
}

func executeSingleNodePendingCommits(conn db.SSHConnection, opts NodesPendingCommitsOptions) NodePendingCommitsResult {
	start := time.Now()
	res := NodePendingCommitsResult{
		NodeAlias: conn.Alias,
		Host:      conn.IPAddress,
		OS:        conn.OS,
	}

	cmd, shell := resolveRemotePendingCommitsCommand(conn, opts)
	out, errExec := currentSSHExecutor.Execute(conn, cmd, shell)
	res.Latency = time.Since(start)
	res.LatencyMs = res.Latency.Milliseconds()

	if errExec != nil {
		return handlePendingCommitsError(res, errExec)
	}

	res.IsNodeOnline = true
	res.IsSuccess = true
	payload := extractJSONPayload(out)
	res.PendingRepos = parseRemotePendingRepos(payload, opts.SortMode, opts.IsDirtyOnly)
	calculatePendingTotals(&res)
	return res
}

func handlePendingCommitsError(res NodePendingCommitsResult, err error) NodePendingCommitsResult {
	errStr := err.Error()
	res.ErrorMessage = errStr
	if isOfflineError(errStr) {
		res.IsNodeOnline = false
		return res
	}

	res.IsNodeOnline = true
	res.IsSuccess = false
	return res
}

func calculatePendingTotals(res *NodePendingCommitsResult) {
	for _, repo := range res.PendingRepos {
		res.TotalDirtyFiles += repo.DirtyFileCount
		res.TotalUnpushed += repo.UnpushedCount
	}
}

func resolveRemotePendingCommitsCommand(conn db.SSHConnection, opts NodesPendingCommitsOptions) (string, string) {
	if isWindowsNode(conn) {
		return buildRemoteWindowsPendingCommitsCommand(opts)
	}

	return buildRemotePendingCommitsCommand(opts)
}

func buildRemotePendingCommitsCommand(opts NodesPendingCommitsOptions) (string, string) {
	parts := []string{"gitmap", "pending-commits", "--json"}
	if opts.SortMode != "" && opts.SortMode != "priority" {
		parts = append(parts, "--sort="+opts.SortMode)
	}

	if opts.IsDetail {
		parts = append(parts, "--detail")
	}

	if opts.IsDirtyOnly {
		parts = append(parts, "--dirty-only")
	}

	return strings.Join(parts, " "), "bash"
}

func buildRemoteWindowsPendingCommitsCommand(opts NodesPendingCommitsOptions) (string, string) {
	parts := []string{"gitmap", "pending-commits", "--json"}
	if opts.SortMode != "" && opts.SortMode != "priority" {
		parts = append(parts, "--sort="+opts.SortMode)
	}

	if opts.IsDetail {
		parts = append(parts, "--detail")
	}

	if opts.IsDirtyOnly {
		parts = append(parts, "--dirty-only")
	}

	inner := strings.Join(parts, " ")
	return fmt.Sprintf(`powershell.exe -NoProfile -Command "%s"`, inner), "ps"
}

type rawRemotePendingEnvelope struct {
	PendingRepos []rawRemoteRepoItem `json:"pending_repos"`
	Repositories []rawRemoteRepoItem `json:"repositories"`
	Results      []rawRemoteRepoItem `json:"results"`
}

type rawRemoteRepoItem struct {
	RepoName           string   `json:"repo_name"`
	Name               string   `json:"name"`
	RelativePath       string   `json:"relative_path"`
	Branch             string   `json:"branch"`
	CurrentBranch      string   `json:"current_branch"`
	IsClean            bool     `json:"is_clean"`
	HasChanges         bool     `json:"has_changes"`
	IsDirty            bool     `json:"is_dirty"`
	HasUncommitted     bool     `json:"has_uncommitted"`
	DirtyFileCount     int      `json:"dirty_file_count"`
	UntrackedFiles     int      `json:"untracked_files_count"`
	ModifiedFiles      int      `json:"modified_files_count"`
	StagedFiles        int      `json:"staged_files_count"`
	HasUnpushedCommits bool     `json:"has_unpushed_commits"`
	HasUnpushed        bool     `json:"has_unpushed"`
	UnpushedCount      int      `json:"unpushed_count"`
	UnpushedCommits    int      `json:"unpushed_commits_count"`
	PriorityScore      int      `json:"priority_score"`
	ChangedFiles       []string `json:"changed_files"`
	PendingFiles       []string `json:"pending_files"`
}

func parseRemotePendingRepos(payload string, sortMode string, isDirtyOnly bool) []NodePendingRepoItem {
	if len(payload) == 0 {
		return nil
	}

	rawItems := extractRawRepoItems(payload)
	var repos []NodePendingRepoItem

	for _, item := range rawItems {
		repo := convertRawToPendingRepoItem(item)
		if isDirtyOnly && !repo.HasChanges && !repo.HasUnpushedCommits {
			continue
		}

		repos = append(repos, repo)
	}

	sortPendingRepos(repos, sortMode)
	return repos
}

func extractRawRepoItems(payload string) []rawRemoteRepoItem {
	items := extractEnvelopeRepoItems(payload)
	if len(items) > 0 {
		return items
	}

	var arr []rawRemoteRepoItem
	if errArr := json.Unmarshal([]byte(payload), &arr); errArr == nil {
		return arr
	}

	return nil
}

func extractEnvelopeRepoItems(payload string) []rawRemoteRepoItem {
	var env rawRemotePendingEnvelope
	if errEnv := json.Unmarshal([]byte(payload), &env); errEnv != nil {
		return nil
	}

	if len(env.PendingRepos) > 0 {
		return env.PendingRepos
	}

	if len(env.Repositories) > 0 {
		return env.Repositories
	}

	return env.Results
}

func convertRawToPendingRepoItem(item rawRemoteRepoItem) NodePendingRepoItem {
	name := resolveRawRepoName(item)
	branch := resolveRawRepoBranch(item)
	dirtyFiles := resolveRawDirtyCount(item)
	unpushed := resolveRawUnpushedCount(item)
	hasChanges := resolveRawHasChanges(item, dirtyFiles)
	hasUnpushed := resolveRawHasUnpushed(item, unpushed)
	isClean := !hasChanges && !hasUnpushed
	score := computeRepoPriorityScore(hasChanges, hasUnpushed, dirtyFiles, unpushed)
	files := resolveRawChangedFiles(item)

	return NodePendingRepoItem{
		RepoName:           name,
		RelativePath:       item.RelativePath,
		Branch:             branch,
		IsClean:            isClean,
		HasChanges:         hasChanges,
		DirtyFileCount:     dirtyFiles,
		HasUnpushedCommits: hasUnpushed,
		UnpushedCount:      unpushed,
		PriorityScore:      score,
		ChangedFiles:       files,
	}
}

func resolveRawRepoName(item rawRemoteRepoItem) string {
	if item.RepoName != "" {
		return item.RepoName
	}

	return item.Name
}

func resolveRawRepoBranch(item rawRemoteRepoItem) string {
	if item.Branch != "" {
		return item.Branch
	}

	return item.CurrentBranch
}

func resolveRawDirtyCount(item rawRemoteRepoItem) int {
	if item.DirtyFileCount > 0 {
		return item.DirtyFileCount
	}

	return item.UntrackedFiles + item.ModifiedFiles + item.StagedFiles
}

func resolveRawUnpushedCount(item rawRemoteRepoItem) int {
	if item.UnpushedCount > 0 {
		return item.UnpushedCount
	}

	return item.UnpushedCommits
}

func resolveRawHasChanges(item rawRemoteRepoItem, dirtyCount int) bool {
	return item.HasChanges || item.IsDirty || item.HasUncommitted || dirtyCount > 0
}

func resolveRawHasUnpushed(item rawRemoteRepoItem, unpushedCount int) bool {
	return item.HasUnpushedCommits || item.HasUnpushed || unpushedCount > 0
}

func resolveRawChangedFiles(item rawRemoteRepoItem) []string {
	if len(item.ChangedFiles) > 0 {
		return item.ChangedFiles
	}

	return item.PendingFiles
}

func computeRepoPriorityScore(hasChanges, hasUnpushed bool, dirtyCount, unpushedCount int) int {
	if hasChanges && hasUnpushed {
		return 100 + dirtyCount + unpushedCount
	}

	if hasChanges {
		return 50 + dirtyCount
	}

	if hasUnpushed {
		return 20 + unpushedCount
	}

	return 0
}

func sortPendingRepos(repos []NodePendingRepoItem, sortMode string) {
	mode := strings.ToLower(sortMode)
	if mode == "name" {
		sort.Slice(repos, func(i, j int) bool {
			return repos[i].RepoName < repos[j].RepoName
		})
		return
	}

	if mode == "count" {
		sort.Slice(repos, func(i, j int) bool {
			return repos[i].DirtyFileCount > repos[j].DirtyFileCount
		})
		return
	}

	sort.Slice(repos, func(i, j int) bool {
		if repos[i].PriorityScore != repos[j].PriorityScore {
			return repos[i].PriorityScore > repos[j].PriorityScore
		}
		return repos[i].RepoName < repos[j].RepoName
	})
}

func emitNodesPendingCommitsJSON(results []NodePendingCommitsResult, opts NodesPendingCommitsOptions, targetCount int) error {
	onlineCount := 0
	totalRepos := 0
	totalDirty := 0
	totalUnpushed := 0

	for _, r := range results {
		if r.IsNodeOnline {
			onlineCount++
		}
		totalRepos += len(r.PendingRepos)
		totalDirty += r.TotalDirtyFiles
		totalUnpushed += r.TotalUnpushed
	}

	agg := NodesPendingCommitsAggregated{
		Timestamp:          time.Now().UTC().Format(time.RFC3339),
		TotalNodesQueried:  targetCount,
		OnlineNodesCount:   onlineCount,
		TotalPendingRepos:  totalRepos,
		TotalDirtyFiles:    totalDirty,
		TotalUnpushedCount: totalUnpushed,
		SortMode:           opts.SortMode,
		Results:            results,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if errEnc := enc.Encode(agg); errEnc != nil {
		return apperror.WrapSimple(errEnc, "encode pending commits aggregated json")
	}

	return nil
}

func renderNodesPendingCommitsTable(results []NodePendingCommitsResult) {
	columns := []termout.Column{
		{Title: "NODE ALIAS", Align: termout.AlignLeft, MinWidth: 14},
		{Title: "REPOSITORY", Align: termout.AlignLeft, MinWidth: 22},
		{Title: "BRANCH", Align: termout.AlignLeft, MinWidth: 12},
		{Title: "DIRTY FILES", Align: termout.AlignRight, MinWidth: 12},
		{Title: "UNPUSHED", Align: termout.AlignRight, MinWidth: 10},
		{Title: "STATUS", Align: termout.AlignLeft, MinWidth: 14},
	}

	rows := buildPendingCommitsRows(results)
	tableCfg := termout.TableConfig{
		Columns: columns,
		Rows:    rows,
	}

	fmt.Println()
	termout.PrintTable(tableCfg)
}

func buildPendingCommitsRows(results []NodePendingCommitsResult) []termout.Row {
	var rows []termout.Row

	for _, res := range results {
		if !res.IsNodeOnline {
			rows = append(rows, termout.Row{
				Cells: []string{res.NodeAlias, "-", "-", "-", "-", constants.ColorDim + "○ OFFLINE" + constants.ColorReset},
			})
			continue
		}

		if len(res.PendingRepos) == 0 {
			rows = append(rows, termout.Row{
				Cells: []string{res.NodeAlias, "(all clean)", "-", "0", "0", constants.ColorGreen + "● SYNCED" + constants.ColorReset},
			})
			continue
		}

		for _, repo := range res.PendingRepos {
			statusBadge := formatPendingBadge(repo.DirtyFileCount, repo.UnpushedCount)
			rows = append(rows, termout.Row{
				Cells: []string{
					res.NodeAlias,
					repo.RepoName,
					repo.Branch,
					fmt.Sprintf("%d", repo.DirtyFileCount),
					fmt.Sprintf("%d", repo.UnpushedCount),
					statusBadge,
				},
			})
		}
	}

	return rows
}

func formatPendingBadge(dirtyCount, unpushedCount int) string {
	if dirtyCount > 0 && unpushedCount > 0 {
		return constants.ColorYellow + "▲ DIRTY+UNPUSHED" + constants.ColorReset
	}

	if dirtyCount > 0 {
		return constants.ColorYellow + "▲ DIRTY" + constants.ColorReset
	}

	if unpushedCount > 0 {
		return constants.ColorCyan + "↑ UNPUSHED" + constants.ColorReset
	}

	return constants.ColorGreen + "● SYNCED" + constants.ColorReset
}

func renderPendingCommitsDetailView(results []NodePendingCommitsResult) {
	fmt.Println()
	fmt.Println("  " + constants.ColorCyan + "Granular File Modifications Breakdown:" + constants.ColorReset)

	for _, res := range results {
		if !res.IsNodeOnline || len(res.PendingRepos) == 0 {
			continue
		}

		for _, repo := range res.PendingRepos {
			if len(repo.ChangedFiles) == 0 {
				continue
			}

			fmt.Printf("    • %s [%s/%s] (%d files):\n",
				repo.RepoName, res.NodeAlias, repo.Branch, len(repo.ChangedFiles))
			for _, file := range repo.ChangedFiles {
				fmt.Printf("        %s- %s%s\n", constants.ColorDim, file, constants.ColorReset)
			}
		}
	}
}

func printPendingCommitsSummary(results []NodePendingCommitsResult, targetCount int) {
	onlineCount := 0
	dirtyRepos := 0
	totalDirtyFiles := 0
	totalUnpushed := 0

	for _, r := range results {
		if r.IsNodeOnline {
			onlineCount++
		}
		for _, repo := range r.PendingRepos {
			if repo.HasChanges {
				dirtyRepos++
			}
			totalDirtyFiles += repo.DirtyFileCount
			totalUnpushed += repo.UnpushedCount
		}
	}

	prefix := constants.ColorGreen + "✓" + constants.ColorReset
	if dirtyRepos > 0 || totalUnpushed > 0 {
		prefix = constants.ColorYellow + "▲" + constants.ColorReset
	}

	fmt.Printf("\n  %s Completed pending-commits query across %d/%d node(s) (%d dirty repos, %d dirty files, %d unpushed)\n\n",
		prefix, onlineCount, targetCount, dirtyRepos, totalDirtyFiles, totalUnpushed)
}

func isPendingCommitsHelp(args []string) bool {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			return true
		}
	}

	return false
}

func printNodesPendingCommitsHelp() error {
	fmt.Println()
	fmt.Println("  ┌─────────────────────────────────────────────────────────────────────────────────┐")
	fmt.Println("  │  gitmap nodes pending-commits (alias: pc)                                       │")
	fmt.Println("  │  Inspect uncommitted and unpushed changes across all cluster fleet nodes via SSH│")
	fmt.Println("  └─────────────────────────────────────────────────────────────────────────────────┘")
	fmt.Println()
	fmt.Println("  USAGE:")
	fmt.Println("    gitmap nodes pending-commits [flags]")
	fmt.Println("    gitmap nodes pc [flags]")
	fmt.Println()
	fmt.Println("  FLAGS:")
	fmt.Println("    -t, --target <alias|ip>    Target specific node alias or IP address")
	fmt.Println("    -e, --except <nodes>       Exclude comma-separated node aliases")
	fmt.Println("        --include <nodes>      Whitelist comma-separated node aliases")
	fmt.Println("        --include-main         Include main controller node in query")
	fmt.Println("        --open-only            Only query nodes currently reachable")
	fmt.Println("    -s, --sort <mode>          Sort repositories: priority (default), name, count")
	fmt.Println("    -d, --detail               Show granular 1-by-1 file modifications")
	fmt.Println("        --dirty-only           Only show repositories with dirty changes")
	fmt.Println("    -j, --json                 Output raw aggregated JSON telemetry")
	fmt.Println("    -h, --help                 Display this help menu")
	fmt.Println()
	fmt.Println("  EXAMPLES:")
	fmt.Println("    gitmap nodes pc")
	fmt.Println("    gitmap nodes pc --sort=count --detail")
	fmt.Println("    gitmap nodes pc -t u1 --json")
	fmt.Println()
	return nil
}
