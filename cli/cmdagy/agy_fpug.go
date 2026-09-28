package cmdagy

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/spf13/cobra"
)

var fpugIntervalStr string

// AgyFPUGCmd continuously monitors target projects until prompts complete and pipelines turn green.
var AgyFPUGCmd = &cobra.Command{
	Use:     "finish-prompts-until-green [targets...]",
	Aliases: []string{"fpug"},
	Short:   "Monitor target projects until prompts finish and CI/CD pipelines turn green",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunFPUGCLI(args)
	},
}

func init() {
	AgyFPUGCmd.Flags().StringVarP(&fpugIntervalStr, "time", "t", "30s", "Polling interval (minimum 30s, e.g. 30s, 1m, 5m)")
	AgyFPUGCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return []string{
				"running-projects\tAll currently running Antigravity projects",
			}, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	AgyCmd.AddCommand(AgyFPUGCmd)
}

// RunFPUGCLI orchestrates the finish-prompts-until-green watcher loop.
func RunFPUGCLI(args []string) error {
	if checkFPUGHelp(args) {
		printFPUGHelp()
		return nil
	}
	targets, err := resolveFPUGTargets(args)
	if err != nil {
		return err
	}
	interval := parseFPUGInterval(fpugIntervalStr)
	return executeFPUGLoop(targets, interval)
}

func checkFPUGHelp(args []string) bool {
	if len(args) == 0 {
		return false
	}
	sub := strings.ToLower(args[0])
	return sub == "help" || sub == "--help" || sub == "-h"
}

func printFPUGHelp() {
	fmt.Println()
	fmt.Printf("  %sAntigravity Finish Prompts Until Green (fpug)%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("  Usage:")
	fmt.Println("    gitmap agy finish-prompts-until-green <targets...> [-t <duration>]")
	fmt.Println("    gitmap agy fpug running-projects [-t 5m]")
	fmt.Println("    gitmap agy fpug help")
	fmt.Println()
	fmt.Println("  Flags:")
	fmt.Println("    -t, --time      Check interval (default 30s, minimum 30s)")
	fmt.Println()
}

func parseFPUGInterval(raw string) time.Duration {
	d, err := time.ParseDuration(raw)
	if err != nil || d < 30*time.Second {
		return 30 * time.Second
	}
	return d
}

type fpugTarget struct {
	Name string
	Path string
	Repo string
}

func resolveFPUGTargets(args []string) ([]fpugTarget, error) {
	if hasRunningProjectsTarget(args) {
		return resolveRunningProjectsTargets()
	}
	return resolveExplicitTargets(args)
}

func hasRunningProjectsTarget(args []string) bool {
	for _, a := range args {
		clean := strings.Trim(strings.ToLower(a), ",")
		if clean == "running-projects" || clean == "runningprojects" {
			return true
		}
	}
	return false
}

func resolveRunningProjectsTargets() ([]fpugTarget, error) {
	running, err := DiscoverRunningProjects()
	if err != nil {
		return nil, err
	}
	var targets []fpugTarget
	for _, r := range running {
		targets = append(targets, fpugTarget{
			Name: r.ProjectName,
			Path: r.ProjectPath,
			Repo: resolveTargetRepoSlug(r.ProjectPath, r.ProjectName),
		})
	}
	return targets, nil
}

func resolveExplicitTargets(args []string) ([]fpugTarget, error) {
	projects, _ := loadActiveSortedProjects()
	matched, err := ResolveAgyProjectTargets(args, projects)
	if err != nil || len(matched) == 0 {
		return resolveFallbackTargets(args)
	}
	var targets []fpugTarget
	for _, p := range matched {
		targets = append(targets, fpugTarget{
			Name: p.Name,
			Path: p.GetPath(),
			Repo: resolveTargetRepoSlug(p.GetPath(), p.Name),
		})
	}
	return targets, nil
}

func resolveFallbackTargets(args []string) ([]fpugTarget, error) {
	var targets []fpugTarget
	for _, arg := range args {
		clean := strings.Trim(arg, ",")
		if clean == "" || strings.HasPrefix(clean, "-") {
			continue
		}
		targets = append(targets, fpugTarget{
			Name: clean,
			Path: clean,
			Repo: clean,
		})
	}
	return targets, nil
}

func resolveTargetRepoSlug(path, name string) string {
	if len(path) == 0 {
		return name
	}
	base := filepath.Base(path)
	if len(base) > 0 && base != "." {
		return base
	}
	return name
}

func executeFPUGLoop(targets []fpugTarget, interval time.Duration) error {
	fmt.Printf("\n  %s[FPUG]%s Monitoring %d project(s) until green (interval: %v)...\n\n",
		constants.ColorCyan, constants.ColorReset, len(targets), interval)
	for {
		pendingCount := checkTargetsStatus(targets)
		if pendingCount == 0 {
			printFPUGSuccess(len(targets))
			return nil
		}
		fmt.Printf("  %s⏳ [FPUG] %d project(s) still waiting to turn green. Next check in %v...%s\n",
			constants.ColorYellow, pendingCount, interval, constants.ColorReset)
		time.Sleep(interval)
	}
}

func checkTargetsStatus(targets []fpugTarget) int {
	pending := 0
	for _, t := range targets {
		isDone := verifyProjectGreenAndIdle(t)
		if !isDone {
			pending++
		}
	}
	return pending
}

func verifyProjectGreenAndIdle(t fpugTarget) bool {
	isQueueEmpty := checkWorkspaceQueueEmpty(t.Path)
	payload, _, hasFailures := cmdpipeline.FetchPipelineErrorReportWithMeta(t.Repo, false)
	isPipelineGreen := !hasFailures && !payload.IsRunning
	return isQueueEmpty && isPipelineGreen
}

func checkWorkspaceQueueEmpty(ws string) bool {
	if ws == "" {
		return true
	}
	activePrompt, _, hasActive := checkActivePromptFile(ws)
	if hasActive && len(activePrompt) > 0 {
		return false
	}
	summary, isFound := inspectSingleWorkspaceQueue(ws, filepath.Base(ws))
	if isFound && summary.TotalQueued > 0 {
		return false
	}
	return true
}

func printFPUGSuccess(total int) {
	fmt.Println()
	fmt.Printf("  %s✔ All %d monitored project(s) have finished prompts and green pipelines!%s\n\n",
		constants.ColorGreen+"\033[1m", total, constants.ColorReset)
}
