package cmdpull

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunPullAllEfficient executes the efficient pull workflow skipping inactive repos.
func RunPullAllEfficient(args []string, isTableMode bool, invokedAlias string, isShortForm bool) error {
	isHelp := isHelpArg(args)
	if isHelp {
		cmdName := resolveEfficientHelpTopic(args, isTableMode, invokedAlias)
		checkHelp(cmdName, args)

		return nil
	}

	useSSH, useHTTPS, isTableFlag, isJSON, targetSSH, cleanArgs := extractEfficientFlags(args)
	if isTableFlag {
		isTableMode = true
	}
	if targetSSH != "" && RunRemoteSSHPullFn != nil {
		return RunRemoteSSHPullFn(targetSSH, isTableMode)
	}

	opts := buildEfficientOptions(cleanArgs, isTableMode, isJSON, invokedAlias, isShortForm, useSSH, useHTTPS, targetSSH)
	fullCmd := resolveEfficientFullCmdName(isTableMode)
	if !isJSON {
		PrintPullBanner(fullCmd, invokedAlias, isShortForm)
	}
	requireOnline()

	records := resolveAllTrackedRecords()
	if len(records) == 0 {
		return handleEmptyRecords(isJSON)
	}

	return processEfficientPullLifecycle(records, opts)
}

func handleEmptyRecords(isJSON bool) error {
	if isJSON {
		return renderJSONInactiveResults(0, nil)
	}
	printNothingToPull()
	return nil
}

func buildEfficientOptions(args []string, isTable, isJSON bool, alias string, isShort bool, ssh, https bool, target string) EfficientPullOptions {
	return EfficientPullOptions{
		IsTableMode:  isTable,
		IsJSON:       isJSON,
		InvokedAlias: alias,
		IsShortForm:  isShort,
		UseSSH:       ssh,
		UseHTTPS:     https,
		TargetSSH:    target,
		Args:         args,
	}
}

func resolveEfficientFullCmdName(isTable bool) string {
	if isTable {
		return "pull all-efficient-table"
	}

	return "pull all-efficient"
}

func isHelpArg(args []string) bool {
	for _, a := range args {
		lower := strings.ToLower(a)
		isHelp := lower == "--help" || lower == "-h" || lower == "help"
		if isHelp {
			return true
		}
	}

	return false
}

func resolveEfficientHelpTopic(args []string, isTableMode bool, invokedAlias string) string {
	lowerAlias := strings.ToLower(invokedAlias)
	hasPaet := strings.Contains(lowerAlias, "paet") || lowerAlias == "paet"
	if hasPaet {
		return "paet"
	}
	hasPae := strings.Contains(lowerAlias, "pae") || lowerAlias == "pae"
	if hasPae {
		return "pae"
	}

	isTable := isTableMode || hasTableFlag(args)
	if isTable {
		return "pull-all-efficient-table"
	}

	return "pull-all-efficient"
}

func hasTableFlag(args []string) bool {
	for _, a := range args {
		lower := strings.ToLower(a)
		isTable := lower == "--status" || lower == "--status-table" || lower == "-status" || lower == "--table"
		if isTable {
			return true
		}
	}

	return false
}

func extractEfficientFlags(args []string) (bool, bool, bool, bool, string, []string) {
	var useSSH, useHTTPS, isTable, isJSON bool
	var targetSSH string
	var rest []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		lower := strings.ToLower(a)
		if lower == "--status" || lower == "--status-table" || lower == "-status" || lower == "--table" {
			isTable = true
			continue
		}
		if lower == "--json" {
			isJSON = true
			continue
		}
		if (lower == "-t" || lower == "--target") && i+1 < len(args) {
			targetSSH = args[i+1]
			i++
			continue
		}
		if strings.HasPrefix(lower, "-t=") {
			targetSSH = a[3:]
			continue
		}
		if strings.HasPrefix(lower, "--target=") {
			targetSSH = a[9:]
			continue
		}
		if lower == "--ssh" || lower == "-ssh" || lower == "ssh" {
			useSSH = true
			continue
		}
		if lower == "--https" || lower == "-https" || lower == "https" {
			useHTTPS = true
			continue
		}
		rest = append(rest, a)
	}

	return useSSH, useHTTPS, isTable, isJSON, targetSSH, rest
}

func resolveAllTrackedRecords() []model.ScanRecord {
	db, err := store.OpenDefault()
	if err != nil {
		return nil
	}
	defer db.Close()

	records, err := db.ListRepos()
	if err != nil {
		return nil
	}

	return records
}

func printNothingToPull() {
	arrow := resolveSubArrow()
	fmt.Printf("    %s%s%s %snothing to pull: no tracked repositories found in database.%s\n\n",
		constants.ColorCyan, arrow, constants.ColorReset,
		constants.ColorDim, constants.ColorReset)
}

func processEfficientPullLifecycle(records []model.ScanRecord, opts EfficientPullOptions) error {
	partition, _ := PartitionRecordsByActivity(records)
	total := len(records)
	if len(partition.ActiveRecords) == 0 {
		return handleAllReposInactive(total, partition.InactiveRepos, opts.IsTableMode, opts.IsJSON)
	}

	if !opts.IsJSON {
		printEfficientPartitionNotice(total, len(partition.ActiveRecords), len(partition.InactiveRepos))
	}
	return executeActiveEfficientBatch(partition, opts, total)
}

func printEfficientPartitionNotice(total, activeCount, inactiveCount int) {
	arrow := resolveSubArrow()
	fmt.Printf("    %s%s%s resolved %d repo(s) (%d active, %d inactive skipped in 24h window)\n\n",
		constants.ColorCyan, arrow, constants.ColorReset,
		total, activeCount, inactiveCount)
}

func executeActiveEfficientBatch(partition EfficientPullPartition, opts EfficientPullOptions, total int) error {
	pullOpts := pullOptions{
		all:      true,
		useSSH:   opts.UseSSH,
		useHTTPS: opts.UseHTTPS,
	}
	maybeApplyTransportToRecords(partition.ActiveRecords, opts.UseSSH, opts.UseHTTPS)

	bar := NewPullProgressBar(len(partition.ActiveRecords), false, false)
	if !opts.IsJSON {
		bar.Start()
	}
	startTime := time.Now()
	executePull(partition.ActiveRecords, bar, pullOpts)
	if !opts.IsJSON {
		bar.Stop()
	}
	dur := time.Since(startTime)

	sortedStates := sortStatesAlphabetically(bar.States())
	if opts.IsJSON {
		return renderJSONEfficientResults(total, sortedStates, partition.InactiveRepos, dur)
	}
	renderEfficientResults(sortedStates, partition.InactiveRepos, opts.IsTableMode)
	recordEfficientPullTelemetry(total, sortedStates, partition.InactiveRepos, dur, opts.IsTableMode)

	return nil
}

func renderEfficientResults(states []*PullRepoState, inactive []InactiveRepoDetail, isTable bool) {
	if isTable {
		renderPullBatchResults(states)
	} else {
		renderConciseActiveResults(states)
	}

	renderEfficientPullSummary(states, inactive)
}

func renderConciseActiveResults(states []*PullRepoState) {
	fmt.Println()
	for _, s := range states {
		statusLabel := s.Changes
		if statusLabel == "" || statusLabel == "synced" {
			statusLabel = "up-to-date"
		}
		fmt.Printf("    • %-26s %s\n", s.RepoName, statusLabel)
	}
}

func renderEfficientPullSummary(states []*PullRepoState, inactive []InactiveRepoDetail) {
	fmt.Printf("\n  %s✓%s %sPull efficient complete:%s %d active pulled, %d inactive skipped\n",
		constants.ColorGreen, constants.ColorReset,
		constants.ColorBold, constants.ColorReset,
		len(states), len(inactive))

	if len(inactive) > 0 {
		printInactiveSkipSummary(inactive)
	}
}

func printInactiveSkipSummary(inactive []InactiveRepoDetail) {
	arrow := resolveSubArrow()
	names := make([]string, 0, len(inactive))
	for _, in := range inactive {
		names = append(names, in.RepoName)
	}
	fmt.Printf("    %s%s%s %sskipped inactive repos (0 changes over 3+ pulls in last 24h):%s\n",
		constants.ColorCyan, arrow, constants.ColorReset,
		constants.ColorDim, constants.ColorReset)
	fmt.Printf("      %s%s%s\n", constants.ColorDim, strings.Join(names, ", "), constants.ColorReset)
	fmt.Printf("    %sTo pull all repositories including inactive ones, run:%s %sgitmap pull all%s\n\n",
		constants.ColorDim, constants.ColorReset,
		constants.ColorBold, constants.ColorReset)
}

func handleAllReposInactive(total int, inactive []InactiveRepoDetail, isTable bool, isJSON bool) error {
	if isJSON {
		return renderJSONInactiveResults(total, inactive)
	}
	fmt.Printf("\n  %sℹ%s All %d repository(ies) are currently inactive (0 changes in last 3+ pulls within 24h).\n",
		constants.ColorCyan, constants.ColorReset, total)
	printInactiveSkipSummary(inactive)
	recordEfficientPullTelemetry(total, nil, inactive, 0, isTable)

	return nil
}

// PartitionRecordsByActivity separates scan records into active and inactive sets based on gitmap-pull.db history.
func PartitionRecordsByActivity(records []model.ScanRecord) (EfficientPullPartition, error) {
	db, err := store.OpenPullSplitDB()
	if err != nil {
		return EfficientPullPartition{ActiveRecords: records}, nil
	}
	defer db.Close()

	var part EfficientPullPartition
	for _, rec := range records {
		classifyRepoActivity(db, rec, &part)
	}

	return part, nil
}

func classifyRepoActivity(db *store.PullSplitDB, rec model.ScanRecord, part *EfficientPullPartition) {
	status, err := db.EvaluateRepoInactivity(rec.AbsolutePath, 3, 24)
	isInactive := err == nil && status.IsInactive
	if isInactive {
		part.InactiveRepos = append(part.InactiveRepos, InactiveRepoDetail{
			RepoName: rec.RepoName,
			RepoPath: rec.AbsolutePath,
			Reason:   status.Reason,
		})
		return
	}

	part.ActiveRecords = append(part.ActiveRecords, rec)
}

func recordEfficientPullTelemetry(total int, states []*PullRepoState, inactive []InactiveRepoDetail, dur time.Duration, isTable bool) {
	cwd, _ := os.Getwd()
	cmdType := resolveEfficientFullCmdName(isTable)
	successCount := 0
	failedCount := 0
	for _, s := range states {
		if s.ErrorMsg == "" && s.Step != PullStepTypeError {
			successCount++
		} else {
			failedCount++
		}
	}

	telemetry := PullSessionTelemetry{
		CommandType:   cmdType,
		WorkingDir:    cwd,
		TotalRepos:    total,
		PulledRepos:   len(states),
		SkippedRepos:  len(inactive),
		SuccessCount:  successCount,
		FailedCount:   failedCount,
		IsEfficient:   true,
		Duration:      dur,
		GitMapVersion: constants.Version,
	}

	repoRuns := buildCombinedRepoRuns(states, inactive)
	_ = recordTelemetryToDB(telemetry, repoRuns)
}

func buildCombinedRepoRuns(states []*PullRepoState, inactive []InactiveRepoDetail) []store.PullRepoRunRecord {
	runs := buildRepoRunRecords(states)
	for _, in := range inactive {
		runs = append(runs, store.PullRepoRunRecord{
			RepoPath:   in.RepoPath,
			RepoName:   in.RepoName,
			PullStatus: "skipped-inactive",
			IsActive:   false,
			HasChanges: false,
			DurationMs: 0,
			Notes:      in.Reason,
		})
	}

	return runs
}

func recordTelemetryToDB(telemetry PullSessionTelemetry, records []store.PullRepoRunRecord) error {
	db, err := store.OpenPullSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	runID, err := db.InsertPullRun(telemetry.ConvertToStoreRunRecord())
	if err != nil {
		return err
	}

	return db.InsertPullRepoRuns(runID, records)
}
