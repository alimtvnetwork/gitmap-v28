// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RunNodesCommits handles standard fleet commits without semantic prefixes.
func RunNodesCommits(args []string) error {
	return RunNodesCommitsWithAction("commits", args)
}

// RunNodesCommitPushFeature handles fleet feature commits (Feature: <msg>).
func RunNodesCommitPushFeature(args []string) error {
	return RunNodesCommitsWithAction("cpf", args)
}

// RunNodesCommitPushBug handles fleet bugfix commits (Bug: <msg>).
func RunNodesCommitPushBug(args []string) error {
	return RunNodesCommitsWithAction("cpb", args)
}

// RunNodesCommitPushRelease handles fleet release commits (Release: <msg>).
func RunNodesCommitPushRelease(args []string) error {
	return RunNodesCommitsWithAction("cpr", args)
}

// RunNodesCommitFix handles fleet repair and merge commits (Fix: <msg>).
func RunNodesCommitFix(args []string) error {
	return RunNodesCommitsWithAction("commit-fix", args)
}

// RunNodesCommitsAction is an alias for RunNodesCommitsWithAction.
func RunNodesCommitsAction(action string, args []string) error {
	return RunNodesCommitsWithAction(action, args)
}

// RunNodesCommitsWithAction executes remote semantic commits across fleet nodes.
func RunNodesCommitsWithAction(action string, args []string) error {
	if isCommitsHelp(args) {
		return printNodesCommitsHelp(action)
	}

	opts, errOpts := parseNodesCommitsOptions(action, args)
	if errOpts != nil {
		return errOpts
	}

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

	results := runFleetCommitsParallel(targets, opts)
	if opts.IsJSON {
		return emitNodesCommitsJSON(results, opts, len(targets))
	}

	renderNodesCommitsTable(results, action)
	printCommitsSummary(results, action, len(targets))
	return nil
}

func parseNodesCommitsOptions(action string, args []string) (NodesCommitsOptions, error) {
	opts := NodesCommitsOptions{
		FilterOpts:    ParseNodeFilterOptions(args),
		Action:        action,
		IsPushEnabled: true,
	}

	var positionals []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if isCommitsFlag(arg) {
			handleCommitsFlag(&opts, arg)
			continue
		}

		step := isFilterFlagStep(args, i)
		if step > 0 {
			i += (step - 1)
			continue
		}

		if !strings.HasPrefix(arg, "-") {
			positionals = append(positionals, arg)
		}
	}

	return finalizeCommitsOptions(opts, positionals)
}

func isCommitsFlag(arg string) bool {
	return arg == "--dry-run" || arg == "-n" || arg == "--no-push" || arg == "--json" || arg == "-j"
}

func handleCommitsFlag(opts *NodesCommitsOptions, arg string) {
	if arg == "--dry-run" || arg == "-n" {
		opts.IsDryRun = true
		return
	}

	if arg == "--no-push" {
		opts.IsPushEnabled = false
		return
	}

	if arg == "--json" || arg == "-j" {
		opts.IsJSON = true
		return
	}
}

func finalizeCommitsOptions(opts NodesCommitsOptions, positionals []string) (NodesCommitsOptions, error) {
	if len(positionals) > 0 {
		opts.TargetScope = positionals[0]
	}

	if len(positionals) > 1 {
		opts.CommitMessage = strings.Join(positionals[1:], " ")
	}

	if opts.TargetScope == "" {
		return opts, apperror.NewSimple("target repository name or 'all' is required", "E9081")
	}

	if opts.CommitMessage == "" && !opts.IsDryRun {
		return opts, apperror.NewSimple("commit message is required", "E9080")
	}

	return opts, nil
}

func runFleetCommitsParallel(targets []db.SSHConnection, opts NodesCommitsOptions) []NodeCommitResult {
	results := make([]NodeCommitResult, len(targets))
	var wg sync.WaitGroup

	for i, conn := range targets {
		wg.Add(1)
		go func(idx int, c db.SSHConnection) {
			defer wg.Done()
			results[idx] = executeSingleNodeCommits(c, opts)
		}(i, conn)
	}

	wg.Wait()
	return results
}

func executeSingleNodeCommits(conn db.SSHConnection, opts NodesCommitsOptions) NodeCommitResult {
	start := time.Now()
	res := NodeCommitResult{
		NodeAlias: conn.Alias,
		Host:      conn.IPAddress,
		OS:        conn.OS,
		Action:    opts.Action,
		IsDryRun:  opts.IsDryRun,
	}

	cmd, shell := resolveRemoteCommitsCommand(conn, opts)
	out, errExec := currentSSHExecutor.Execute(conn, cmd, shell)
	res.Latency = time.Since(start)
	res.LatencyMs = res.Latency.Milliseconds()

	if errExec != nil {
		return handleCommitsError(res, errExec)
	}

	res.IsNodeOnline = true
	res.IsSuccess = true
	payload := extractJSONPayload(out)
	res.Repositories, res.CommittedCount = parseRemoteCommitsOutcomes(payload, opts.CommitMessage, opts.IsPushEnabled)
	return res
}

func handleCommitsError(res NodeCommitResult, err error) NodeCommitResult {
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

func resolveRemoteCommitsCommand(conn db.SSHConnection, opts NodesCommitsOptions) (string, string) {
	if isWindowsNode(conn) {
		return buildRemotePowerShellCommitsCommand(opts)
	}

	return buildRemoteBashCommitsCommand(opts)
}

func buildRemoteBashCommitsCommand(opts NodesCommitsOptions) (string, string) {
	parts := []string{
		"gitmap", "sends", opts.Action,
		quoteBashArg(opts.TargetScope),
		quoteBashArg(opts.CommitMessage),
		"--json",
	}

	if opts.IsDryRun {
		parts = append(parts, "--dry-run")
	}

	if !opts.IsPushEnabled {
		parts = append(parts, "--no-push")
	}

	return strings.Join(parts, " "), "bash"
}

func buildRemotePowerShellCommitsCommand(opts NodesCommitsOptions) (string, string) {
	parts := []string{
		"gitmap", "sends", opts.Action,
		quotePowerShellArg(opts.TargetScope),
		quotePowerShellArg(opts.CommitMessage),
		"--json",
	}

	if opts.IsDryRun {
		parts = append(parts, "--dry-run")
	}

	if !opts.IsPushEnabled {
		parts = append(parts, "--no-push")
	}

	inner := strings.Join(parts, " ")
	return fmt.Sprintf(`powershell.exe -NoProfile -Command "%s"`, inner), "ps"
}

func quoteBashArg(s string) string {
	if s == "" {
		return "''"
	}

	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

func quotePowerShellArg(s string) string {
	if s == "" {
		return `\"\"`
	}

	escaped := strings.ReplaceAll(s, `"`, `\"`)
	return `\"` + escaped + `\"`
}

func extractJSONPayload(s string) string {
	trimmed := strings.TrimSpace(s)
	if len(trimmed) == 0 {
		return ""
	}

	objStart := strings.Index(trimmed, "{")
	objEnd := strings.LastIndex(trimmed, "}")
	arrStart := strings.Index(trimmed, "[")
	arrEnd := strings.LastIndex(trimmed, "]")

	hasObj := objStart >= 0 && objEnd > objStart
	hasArr := arrStart >= 0 && arrEnd > arrStart

	if hasObj && hasArr {
		return resolveEnclosingJSONBounds(trimmed, objStart, objEnd, arrStart, arrEnd)
	}

	if hasObj {
		return trimmed[objStart : objEnd+1]
	}

	if hasArr {
		return trimmed[arrStart : arrEnd+1]
	}

	return trimmed
}

func resolveEnclosingJSONBounds(trimmed string, objStart, objEnd, arrStart, arrEnd int) string {
	if objStart < arrStart && objEnd > arrEnd {
		return trimmed[objStart : objEnd+1]
	}

	if arrStart < objStart && arrEnd > objEnd {
		return trimmed[arrStart : arrEnd+1]
	}

	if objStart < arrStart {
		return trimmed[objStart : objEnd+1]
	}

	return trimmed[arrStart : arrEnd+1]
}

type rawRemoteCommitsEnvelope struct {
	Repositories   []rawRemoteRepoCommit `json:"repositories"`
	Results        []rawRemoteRepoCommit `json:"results"`
	CommittedCount int                   `json:"committed_count"`
	TotalCommitted int                   `json:"total_committed"`
	IsPushed       bool                  `json:"is_pushed"`
	ErrorMessage   string                `json:"error_message"`
}

type rawRemoteRepoCommit struct {
	RepoName     string `json:"repo_name"`
	Branch       string `json:"branch"`
	HeadSHA      string `json:"head_sha"`
	CommitHash   string `json:"commit_hash"`
	CommitMsg    string `json:"commit_msg"`
	IsSuccess    bool   `json:"is_success"`
	IsPushed     bool   `json:"is_pushed"`
	HasChanges   bool   `json:"has_changes"`
	FilesStaged  int    `json:"files_staged"`
	FilesChanged int    `json:"files_changed"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
}

func parseRemoteCommitsOutcomes(payload string, defaultMsg string, isPushEnabled bool) ([]RemoteRepoCommitOutcome, int) {
	if len(payload) == 0 {
		return nil, 0
	}

	rawItems, committed := extractRawCommitsItems(payload)
	var outcomes []RemoteRepoCommitOutcome

	for _, item := range rawItems {
		outcome := convertRawToCommitOutcome(item, defaultMsg, isPushEnabled)
		outcomes = append(outcomes, outcome)
	}

	if committed == 0 {
		committed = countSuccessfulCommits(outcomes)
	}

	return outcomes, committed
}

func extractRawCommitsItems(payload string) ([]rawRemoteRepoCommit, int) {
	items, committed := extractEnvelopeCommitsItems(payload)
	if len(items) > 0 {
		return items, committed
	}

	var arr []rawRemoteRepoCommit
	if errArr := json.Unmarshal([]byte(payload), &arr); errArr == nil {
		return arr, 0
	}

	return nil, 0
}

func extractEnvelopeCommitsItems(payload string) ([]rawRemoteRepoCommit, int) {
	var env rawRemoteCommitsEnvelope
	if errEnv := json.Unmarshal([]byte(payload), &env); errEnv != nil {
		return nil, 0
	}

	committed := resolveCommittedCount(env)
	if len(env.Repositories) > 0 {
		return env.Repositories, committed
	}

	return env.Results, committed
}

func resolveCommittedCount(env rawRemoteCommitsEnvelope) int {
	if env.CommittedCount > 0 {
		return env.CommittedCount
	}

	return env.TotalCommitted
}

func convertRawToCommitOutcome(item rawRemoteRepoCommit, defaultMsg string, isPushEnabled bool) RemoteRepoCommitOutcome {
	hash := item.CommitHash
	if hash == "" {
		hash = item.HeadSHA
	}

	msg := item.CommitMsg
	if msg == "" {
		msg = defaultMsg
	}

	files := item.FilesChanged
	if files == 0 {
		files = item.FilesStaged
	}

	hasChanges := item.HasChanges || files > 0 || hash != ""
	isSuccess := item.IsSuccess || (item.Status == "COMMITTED" || item.Status == "SUCCESS" || hash != "")
	isPushed := item.IsPushed || (isPushEnabled && isSuccess && hasChanges)

	return RemoteRepoCommitOutcome{
		RepoName:     item.RepoName,
		Branch:       item.Branch,
		CommitHash:   hash,
		CommitMsg:    msg,
		IsSuccess:    isSuccess,
		IsPushed:     isPushed,
		HasChanges:   hasChanges,
		FilesChanged: files,
		ErrorMessage: item.ErrorMessage,
	}
}

func countSuccessfulCommits(outcomes []RemoteRepoCommitOutcome) int {
	count := 0
	for _, o := range outcomes {
		if o.IsSuccess && o.HasChanges {
			count++
		}
	}

	return count
}

func emitNodesCommitsJSON(results []NodeCommitResult, opts NodesCommitsOptions, targetCount int) error {
	successfulNodes := 0
	totalUpdated := 0

	for _, r := range results {
		if r.IsSuccess {
			successfulNodes++
		}
		totalUpdated += r.CommittedCount
	}

	agg := NodesCommitsAggregated{
		Timestamp:         time.Now().UTC().Format(time.RFC3339),
		Action:            opts.Action,
		TargetRepo:        opts.TargetScope,
		CommitMessage:     opts.CommitMessage,
		IsDryRun:          opts.IsDryRun,
		TotalNodesTarget:  targetCount,
		SuccessfulNodes:   successfulNodes,
		TotalReposUpdated: totalUpdated,
		Results:           results,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if errEnc := enc.Encode(agg); errEnc != nil {
		return apperror.WrapSimple(errEnc, "encode nodes commits aggregated json")
	}

	return nil
}

func renderNodesCommitsTable(results []NodeCommitResult, action string) {
	columns := []termout.Column{
		{Title: "NODE ALIAS", Align: termout.AlignLeft, MinWidth: 14},
		{Title: "REPOSITORY", Align: termout.AlignLeft, MinWidth: 20},
		{Title: "ACTION", Align: termout.AlignLeft, MinWidth: 10},
		{Title: "COMMIT HASH", Align: termout.AlignLeft, MinWidth: 12},
		{Title: "STATUS", Align: termout.AlignLeft, MinWidth: 16},
		{Title: "DURATION", Align: termout.AlignRight, MinWidth: 10},
	}

	rows := buildCommitsRows(results, action)
	tableCfg := termout.TableConfig{
		Columns: columns,
		Rows:    rows,
	}

	fmt.Println()
	termout.PrintTable(tableCfg)
}

func buildCommitsRows(results []NodeCommitResult, action string) []termout.Row {
	var rows []termout.Row

	for _, res := range results {
		if !res.IsNodeOnline {
			rows = append(rows, termout.Row{
				Cells: []string{res.NodeAlias, "-", strings.ToUpper(action), "-", constants.ColorDim + "○ OFFLINE" + constants.ColorReset, "-"},
			})
			continue
		}

		if len(res.Repositories) == 0 {
			rows = append(rows, termout.Row{
				Cells: []string{res.NodeAlias, "(none)", strings.ToUpper(action), "-", constants.ColorDim + "○ NO REPOS" + constants.ColorReset, res.Latency.Round(time.Millisecond).String()},
			})
			continue
		}

		for _, repo := range res.Repositories {
			statusStr := formatCommitStatusBadge(repo.IsSuccess, repo.HasChanges, repo.IsPushed, res.IsDryRun)
			rows = append(rows, termout.Row{
				Cells: []string{
					res.NodeAlias,
					repo.RepoName,
					strings.ToUpper(action),
					shortHash(repo.CommitHash),
					statusStr,
					res.Latency.Round(time.Millisecond).String(),
				},
			})
		}
	}

	return rows
}

func formatCommitStatusBadge(isOK, hasChanges, isPushed, isDryRun bool) string {
	if isDryRun {
		return constants.ColorYellow + "● DRY-RUN" + constants.ColorReset
	}

	if !isOK {
		return constants.ColorRed + "▲ FAILED" + constants.ColorReset
	}

	if !hasChanges {
		return constants.ColorDim + "○ NO CHANGES" + constants.ColorReset
	}

	if isPushed {
		return constants.ColorGreen + "● PUSHED" + constants.ColorReset
	}

	return constants.ColorCyan + "● COMMITTED" + constants.ColorReset
}

func shortHash(sha string) string {
	if len(sha) > 8 {
		return sha[:8]
	}

	if sha == "" {
		return "-"
	}

	return sha
}

func printCommitsSummary(results []NodeCommitResult, action string, targetCount int) {
	successNodes := 0
	totalCommitted := 0

	for _, r := range results {
		if r.IsSuccess {
			successNodes++
		}
		totalCommitted += r.CommittedCount
	}

	prefix := constants.ColorGreen + "✓" + constants.ColorReset
	if successNodes < targetCount {
		prefix = constants.ColorYellow + "▲" + constants.ColorReset
	}

	fmt.Printf("\n  %s Completed %s across %d/%d node(s) (%d repositories updated)\n\n",
		prefix, action, successNodes, targetCount, totalCommitted)
}

func isCommitsHelp(args []string) bool {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			return true
		}
	}

	return false
}

func printNodesCommitsHelp(action string) error {
	fmt.Println()
	fmt.Println("  ┌─────────────────────────────────────────────────────────────────────────────────┐")
	fmt.Printf("  │  gitmap nodes %-66s│\n", fmt.Sprintf("%s <repoName|all> \"<msg>\" [flags]", action))
	fmt.Println("  │  Delegate semantic commits and pushes across remote cluster nodes via SSH       │")
	fmt.Println("  └─────────────────────────────────────────────────────────────────────────────────┘")
	fmt.Println()
	fmt.Println("  ACTIONS:")
	fmt.Println("    commits       Standard commit (no semantic prefix)")
	fmt.Println("    cpf           Feature commit (prepends 'Feature: ')")
	fmt.Println("    cpb           Bug fix commit (prepends 'Bug: ')")
	fmt.Println("    cpr           Release commit (prepends 'Release: ')")
	fmt.Println("    commit-fix    Merge and conflict fix commit (prepends 'Fix: ')")
	fmt.Println()
	fmt.Println("  TARGETS:")
	fmt.Println("    <repoName>    Specific repository directory or slug on target node(s)")
	fmt.Println("    all           All dirty repositories detected on target node(s)")
	fmt.Println()
	fmt.Println("  FLAGS:")
	fmt.Println("    -t, --target <alias|ip>    Target specific node alias or IP address")
	fmt.Println("    -e, --except <nodes>       Exclude comma-separated node aliases")
	fmt.Println("        --include <nodes>      Whitelist comma-separated node aliases")
	fmt.Println("        --include-main         Include main controller node")
	fmt.Println("        --open-only            Only dispatch to reachable nodes")
	fmt.Println("    -n, --dry-run              Simulate commit operations without staging/pushing")
	fmt.Println("        --no-push              Stage and commit on remote node without git push")
	fmt.Println("    -j, --json                 Output raw aggregated JSON telemetry")
	fmt.Println("    -h, --help                 Display this help menu")
	fmt.Println()
	fmt.Println("  EXAMPLES:")
	fmt.Println("    gitmap nodes cpf gitmap \"add nodes commit suite\"")
	fmt.Println("    gitmap nodes cpb all \"fix nil pointer exception in split-db indexer\"")
	fmt.Println("    gitmap nodes commit-fix my-service \"resolve merge conflict in main branch\"")
	fmt.Println("    gitmap nodes commits all \"sync latest submodules\" --dry-run -t u1")
	fmt.Println()
	return nil
}
