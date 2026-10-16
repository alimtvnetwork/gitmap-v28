package cmdpull

import (
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunPullAllEfficient executes the efficient pull workflow skipping inactive repos.
func RunPullAllEfficient(args []string, isTableMode bool, invokedAlias string, isShortForm bool) error {
	if isHelpArg(args) {
		checkHelp(resolveEfficientHelpTopic(args, isTableMode, invokedAlias), args)
		return nil
	}
	if hasSSHFleetFlag(args) {
		return handleSSHFleetPullAll(args)
	}

	useSSH, useHTTPS, isTableFlag, isJSON, targetSSH, windowHours, cleanArgs := extractEfficientFlags(args)
	if isTableFlag {
		isTableMode = true
	}
	if targetSSH != "" && RunRemoteSSHPullFn != nil {
		return RunRemoteSSHPullFn(targetSSH, isTableMode)
	}

	effectiveWindow := resolveEfficientWindowHours(windowHours)
	opts := buildEfficientOptions(cleanArgs, isTableMode, isJSON, invokedAlias, isShortForm, useSSH, useHTTPS, targetSSH, effectiveWindow)
	if !isJSON {
		PrintPullBanner(resolveEfficientFullCmdName(isTableMode), invokedAlias, isShortForm)
	}
	requireOnline()
	maybeExecuteEfficientProbe(args, isJSON)

	records := resolveAllTrackedRecords()
	if len(records) == 0 {
		return handleEmptyRecords(isJSON)
	}
	return processEfficientPullLifecycle(records, opts)
}

func buildEfficientOptions(args []string, isTable, isJSON bool, alias string, isShort bool, ssh, https bool, target string, windowHours ...int) EfficientPullOptions {
	effectiveWindow := constants.DefaultPullActivityWindowHours
	if len(windowHours) > 0 && windowHours[0] > 0 {
		effectiveWindow = windowHours[0]
	}
	return EfficientPullOptions{
		IsTableMode: isTable, IsJSON: isJSON, InvokedAlias: alias, IsShortForm: isShort,
		UseSSH: ssh, UseHTTPS: https, TargetSSH: target, WindowHours: effectiveWindow, Args: args,
	}
}

func resolveEfficientWindowHours(cliWindow int) int {
	if cliWindow > 0 {
		return cliWindow
	}
	db, err := store.OpenDefault()
	if err != nil {
		return constants.DefaultPullActivityWindowHours
	}
	defer db.Close()
	if val := queryDBWindowSetting(db, constants.SettingPullActivityWindowHours); val > 0 {
		return val
	}
	if legacyVal := queryDBWindowSetting(db, "pull_activity_window_hours"); legacyVal > 0 {
		return legacyVal
	}
	return constants.DefaultPullActivityWindowHours
}

func queryDBWindowSetting(db *store.DB, key string) int {
	val, err := db.GetSetting(key)
	if err != nil || val == "" {
		return 0
	}
	parsed, parseErr := strconv.Atoi(val)
	if parseErr != nil || parsed <= 0 {
		return 0
	}
	return parsed
}

func isHelpArg(args []string) bool {
	for _, a := range args {
		switch strings.ToLower(a) {
		case "--help", "-h", "help":
			return true
		}
	}
	return false
}

func resolveEfficientHelpTopic(args []string, isTableMode bool, invokedAlias string) string {
	low := strings.ToLower(invokedAlias)
	switch {
	case strings.Contains(low, "paet"):
		return "paet"
	case strings.Contains(low, "pae"):
		return "pae"
	case isTableMode || hasTableFlag(args):
		return "pull-all-efficient-table"
	default:
		return "pull-all-efficient"
	}
}

func hasTableFlag(args []string) bool {
	for _, a := range args {
		switch strings.ToLower(a) {
		case "--status", "--status-table", "-status", "--table":
			return true
		}
	}
	return false
}

func parseWindowArg(lower, a string, args []string, idx int) (int, int, bool) {
	isEq := strings.HasPrefix(lower, "--hours=") || strings.HasPrefix(lower, "-hours=") ||
		strings.HasPrefix(lower, "--window=") || strings.HasPrefix(lower, "-window=")
	if isEq {
		parts := strings.SplitN(a, "=", 2)
		parsed, _ := strconv.Atoi(parts[1])
		return parsed, idx, parsed > 0
	}
	isKey := lower == "--hours" || lower == "-hours" || lower == "--window" || lower == "-window"
	if isKey && idx+1 < len(args) {
		parsed, _ := strconv.Atoi(args[idx+1])
		return parsed, idx + 1, parsed > 0
	}
	return 0, idx, false
}

func extractEfficientFlags(args []string) (bool, bool, bool, bool, string, int, []string) {
	var useSSH, useHTTPS, isTable, isJSON bool
	var targetSSH string
	var windowHours int
	var rest []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		lower := strings.ToLower(a)
		parsedWindow, nextIdx, isWindow := parseWindowArg(lower, a, args, i)
		if isWindow {
			windowHours, i = parsedWindow, nextIdx
			continue
		}
		switch {
		case lower == "--status" || lower == "--status-table" || lower == "-status" || lower == "--table":
			isTable = true
		case lower == "--json":
			isJSON = true
		case isProbeFlag(lower):
			continue
		case (lower == "-t" || lower == "--target") && i+1 < len(args):
			targetSSH, i = args[i+1], i+1
		case strings.HasPrefix(lower, "-t="):
			targetSSH = a[3:]
		case strings.HasPrefix(lower, "--target="):
			targetSSH = a[9:]
		case lower == "--ssh" || lower == "-ssh" || lower == "ssh":
			useSSH = true
		case lower == "--https" || lower == "-https" || lower == "https":
			useHTTPS = true
		default:
			rest = append(rest, a)
		}
	}
	return useSSH, useHTTPS, isTable, isJSON, targetSSH, windowHours, rest
}

func resolveAllTrackedRecords() []model.ScanRecord {
	db, err := store.OpenDefault()
	if err != nil {
		return nil
	}
	defer db.Close()
	_, _ = OptimizeRedundantRepos(db, false)
	records, err := db.ListRepos()
	if err != nil {
		return nil
	}
	return deduplicateTrackedRecords(records)
}

func deduplicateTrackedRecords(records []model.ScanRecord) []model.ScanRecord {
	seen := make(map[string]bool, len(records))
	unique := make([]model.ScanRecord, 0, len(records))
	for _, r := range records {
		canonical := fsutil.CanonicalPathKey(r.AbsolutePath)
		isNew := r.AbsolutePath != "" && !seen[canonical]
		if isNew {
			seen[canonical] = true
			unique = append(unique, r)
		}
	}
	return unique
}

func maybeExecuteEfficientProbe(args []string, isJSON bool) {
	if hasProbeFlag(args) {
		executePullProbeHook(pullOptions{isJSON: isJSON, yes: hasYesFlag(args)})
	}
}

func isProbeFlag(lower string) bool {
	return lower == "--probe" || lower == "-probe" || lower == "--probe-repos" || lower == "-probe-repos"
}

func hasProbeFlag(args []string) bool {
	for _, a := range args {
		if isProbeFlag(strings.ToLower(a)) {
			return true
		}
	}
	return false
}

func hasYesFlag(args []string) bool {
	for _, a := range args {
		switch strings.ToLower(a) {
		case "-y", "--yes", "-yes":
			return true
		}
	}
	return false
}

func processEfficientPullLifecycle(records []model.ScanRecord, opts EfficientPullOptions) error {
	partition, _ := PartitionRecordsByActivity(records, opts.WindowHours)
	total := len(records)
	if len(partition.ActiveRecords) == 0 {
		return handleAllReposInactive(total, partition.InactiveRepos, opts.IsTableMode, opts.IsJSON, opts.WindowHours)
	}
	if !opts.IsJSON {
		printEfficientPartitionNotice(total, len(partition.ActiveRecords), len(partition.InactiveRepos), opts.WindowHours)
	}
	return executeActiveEfficientBatch(partition, opts, total, records)
}

func executeActiveEfficientBatch(partition EfficientPullPartition, opts EfficientPullOptions, total int, allRecords []model.ScanRecord) error {
	bar := NewPullProgressBar(len(partition.ActiveRecords), false, false)
	if !opts.IsJSON {
		bar.Start()
		defer bar.Stop()
	}
	startTime := time.Now()
	stopHeartbeat := startActiveBatchHeartbeat(bar, opts.IsJSON)
	defer stopHeartbeat()

	ignoreHandle := launchThrottledActiveIgnoreScan(allRecords)
	executePull(partition.ActiveRecords, bar, pullOptions{all: true, useSSH: opts.UseSSH, useHTTPS: opts.UseHTTPS})
	collectAndRemediateEfficientIgnore(ignoreHandle, opts.Args, opts.IsJSON)

	return handleEfficientBatchFinish(total, bar, partition, allRecords, opts, time.Since(startTime))
}

// PartitionRecordsByActivity separates scan records into active and inactive sets based on gitmap-pull.db history.
func PartitionRecordsByActivity(records []model.ScanRecord, windowHours ...int) (EfficientPullPartition, error) {
	effectiveWindow := constants.DefaultPullActivityWindowHours
	if len(windowHours) > 0 && windowHours[0] > 0 {
		effectiveWindow = windowHours[0]
	}
	db, err := store.OpenPullSplitDB()
	if err != nil {
		return EfficientPullPartition{ActiveRecords: records}, nil
	}
	defer db.Close()

	var part EfficientPullPartition
	for _, rec := range records {
		classifyRepoActivity(db, rec, &part, effectiveWindow)
	}
	return part, nil
}

func classifyRepoActivity(db *store.PullSplitDB, rec model.ScanRecord, part *EfficientPullPartition, windowHours ...int) {
	effectiveWindow := store.DefaultActivityWindowHours
	if len(windowHours) > 0 && windowHours[0] > 0 {
		effectiveWindow = windowHours[0]
	}
	status, err := db.EvaluateRepoActivityStatus(rec.AbsolutePath, effectiveWindow, store.DefaultFreshnessCooldownMinutes)
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
