// Package cmdagy — agy_wpr_actions.go executes WPR subcommands for backup, listing, control, and logs.
package cmdagy

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunWPRAll captures prompt snapshots across all projects and displays a consolidated summary.
func RunWPRAll(opts WPROptions) error {
	if opts.IsSSH {
		return AggregateSSHWPRAll(opts)
	}

	summaries, err := SnapshotAllRunningPrompts()
	if err != nil {
		return err
	}

	_ = RecordWPRActionToTaskHistory("all", "all_projects", opts)
	if opts.IsJSON {
		return printJSON(summaries)
	}

	RenderWPRAllSummaryTable(summaries)
	return nil
}

// RunWPRLs lists all scheduled and watched projects.
func RunWPRLs(opts WPROptions) error {
	if opts.IsSSH {
		return AggregateSSHWPRLs(opts)
	}

	st, isRunning := LoadWPRRuntimeStatus()
	projects := collectWatchedProjectStatuses(st.TargetSlugs, isRunning)
	if opts.IsJSON {
		return printJSON(projects)
	}

	RenderWPRProjectsTable(projects, st, isRunning)
	return nil
}

func collectWatchedProjectStatuses(targets []string, isWatching bool) []WPRProjectStatus {
	allProjects, _ := getAllProjects()
	var list []WPRProjectStatus

	for _, p := range filterNonRestrictedProjects(allProjects) {
		slug := store.SanitizeSlug(filepath.Base(p.GetPath()))
		status := buildSingleProjectStatus(p, slug, isWatching)
		list = append(list, status)
	}

	return list
}

func buildSingleProjectStatus(p AgyProject, slug string, isWatching bool) WPRProjectStatus {
	sum, hasSum := LoadWatchPromptsSummary(slug)
	promptStatus := "idle"
	wordCount := 0
	mediaCount := 0

	if hasSum {
		promptStatus = sum.PromptStatus
		wordCount = sum.WordCount
		mediaCount = sum.MediaCount
	}

	return WPRProjectStatus{
		RepoSlug:     slug,
		ProjectName:  p.Name,
		ProjectPath:  p.GetPath(),
		ProjectId:    p.ID,
		SequenceId:   p.ID,
		PromptStatus: promptStatus,
		WordCount:    wordCount,
		MediaCount:   mediaCount,
		IsRunning:    promptStatus == "running",
		IsWatching:   isWatching,
		UpdatedAt:    sum.UpdatedAt,
	}
}

// RunWPRStart registers targets and starts or triggers the monitoring loop.
func RunWPRStart(opts WPROptions) error {
	_ = RecordWPRActionToTaskHistory("start", opts.Target, opts)
	st, _ := LoadWPRRuntimeStatus()
	st.IsActive = true
	_ = SaveWPRRuntimeStatus(st)

	if opts.IsDryRun {
		fmt.Printf("  %s[dry-run] Watch loop previewed for target: %s (interval: %s)%s\n",
			constants.ColorDim, opts.Target, opts.IntervalStr, constants.ColorReset)
		return nil
	}

	fmt.Printf("  %s✔ Started watching prompts for %s (interval: %s)%s\n",
		constants.ColorGreen, opts.Target, opts.IntervalStr, constants.ColorReset)

	return ExecuteWPRWatchLoop(opts)
}

// RunWPRDisable pauses the watch loop without removing project data.
func RunWPRDisable(opts WPROptions) error {
	_ = RecordWPRActionToTaskHistory("disable", opts.Target, opts)
	st, _ := LoadWPRRuntimeStatus()
	st.IsActive = false
	_ = SaveWPRRuntimeStatus(st)

	fmt.Printf("  %s✔ Watch loop disabled (projects preserved in watch list)%s\n",
		constants.ColorGreen, constants.ColorReset)

	return nil
}

// RunWPRShutdown terminates the watch loop and cleans active processes.
func RunWPRShutdown(opts WPROptions) error {
	_ = RecordWPRActionToTaskHistory("shutdown", opts.Target, opts)
	st, _ := LoadWPRRuntimeStatus()
	st.IsActive = false
	_ = SaveWPRRuntimeStatus(st)

	terminateRunningIDE()
	fmt.Printf("  %s✔ Watch prompts running shut down and IDE closed%s\n",
		constants.ColorGreen, constants.ColorReset)

	return nil
}

// RunWPRRestart terminates and relaunches IDE, then replays prompts.
func RunWPRRestart(opts WPROptions) error {
	_ = RecordWPRActionToTaskHistory("restart", opts.Target, opts)
	target := opts.Target
	if target == "" || target == "all" {
		target = "1"
	}

	return RestartAndRerunProject(target, true, opts.IsDryRun, opts.PrefixTemplate)
}

// RunWPRSwitchAccount delegates account switching with prompt preservation.
func RunWPRSwitchAccount(opts WPROptions) error {
	_ = RecordWPRActionToTaskHistory("switch-account", opts.Email, opts)
	switchOpts := AccountSwitchOptions{
		ThresholdPct: 99,
		InstanceName: "wpr-switch",
		IsDryRun:     opts.IsDryRun,
		IsJSON:       opts.IsJSON,
	}

	return RunAccountSwitchWorkflow(switchOpts)
}

// RunWPRFastForward switches account and forces IDE restart with fast forward mode.
func RunWPRFastForward(opts WPROptions) error {
	_ = RecordWPRActionToTaskHistory("fast-forward", opts.Email, opts)
	fmt.Println("  • Fast-forward mode initiated: backing up prompts and restarting IDE...")
	_, _ = SnapshotAllRunningPrompts()

	cand := resolveTargetAccountCandidate(opts.Email)
	_ = delegateFastForwardSwitch(cand, "wpr-fast-forward")
	terminateRunningIDE()

	allProjects, _ := getAllProjects()
	if len(allProjects) > 0 {
		launchWorkspaceIDE(allProjects[0].GetPath())
	}

	fmt.Printf("  %s✔ Fast-forward completed for account: %s%s\n",
		constants.ColorGreen, cand.Email, constants.ColorReset)

	return nil
}

// RunWPRLogs prints execution events and recovery logs for a repo slug.
func RunWPRLogs(opts WPROptions) error {
	if opts.IsSSH {
		return AggregateSSHWPRLogs(opts)
	}

	slug := store.SanitizeSlug(opts.Target)
	db, err := store.OpenWatchPromptsSplitDB(slug)
	if err != nil {
		return err
	}
	defer db.Close()

	logs, err := db.ListWatchLogs(slug, 30)
	if err != nil {
		return err
	}

	if opts.IsJSON {
		return printJSON(logs)
	}

	RenderWPRLogsTable(slug, logs)
	return nil
}

// RunWPRStatus prints current watchdog state and remote machine identities.
func RunWPRStatus(opts WPROptions) error {
	if opts.IsSSH {
		return AggregateSSHWPRStatus(opts)
	}

	st, isRunning := LoadWPRRuntimeStatus()
	cachedMachines := LoadWPRMachineCache()
	if opts.IsJSON {
		payload := map[string]interface{}{
			"Runtime":  st,
			"Machines": cachedMachines,
		}
		return printJSON(payload)
	}

	RenderWPRStatusBox(st, isRunning, cachedMachines)
	return nil
}

// RunWPRRemove removes designated projects or all projects from the watch list.
func RunWPRRemove(opts WPROptions) error {
	_ = RecordWPRActionToTaskHistory("remove", opts.Target, opts)
	st, _ := LoadWPRRuntimeStatus()
	if strings.EqualFold(opts.Target, "all") || opts.Target == "*" {
		st.TargetSlugs = []string{}
	} else {
		st.TargetSlugs = filterRemovalTargets(st.TargetSlugs, opts.Target)
	}

	_ = SaveWPRRuntimeStatus(st)
	fmt.Printf("  %s✔ Removed %q from watch list (remaining: %d project(s))%s\n",
		constants.ColorGreen, opts.Target, len(st.TargetSlugs), constants.ColorReset)

	return nil
}

func filterRemovalTargets(existing []string, toRemove string) []string {
	var kept []string
	clean := strings.ToLower(strings.TrimSpace(toRemove))

	for _, item := range existing {
		if !strings.EqualFold(item, clean) && !strings.EqualFold(filepath.Base(item), clean) {
			kept = append(kept, item)
		}
	}

	return kept
}
