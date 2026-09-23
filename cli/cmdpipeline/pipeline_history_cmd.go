package cmdpipeline

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// PipelineHistoryFlags holds options for gitmap pipeline history.
type PipelineHistoryFlags struct {
	Limit      int
	Repo       string
	Branch     string
	IsDetailed bool
	IsJSON     bool
}

func handlePipelineHistory(args []string) error {
	if hasArgFlag(args, "--help") || hasArgFlag(args, "-h") {
		printPipelineHistoryHelp()

		return nil
	}

	flags := parsePipelineHistoryFlags(args)
	repo := resolveHistoryRepo(flags.Repo)
	runs := queryWorkflowRuns(repo)
	groups := GroupRunsByCommit(runs)

	return renderPipelineHistory(repo, groups, flags)
}

func printPipelineHistoryHelp() {
	fmt.Println(constants.ColorCyan + "Usage: gitmap pipeline history [flags]" + constants.ColorReset)
	fmt.Println("Aliases: hist, h")
	fmt.Println("\nFlags:")
	fmt.Println("  -n, --limit <N>     Number of commits to display (default: 5)")
	fmt.Println("  -v, --detailed      Show expanded job/step tree details")
	fmt.Println("  --branch <branch>   Filter history by branch name")
	fmt.Println("  --repo <slug>       Specify repository slug")
	fmt.Println("  --json              Output structured JSON")
}

func parsePipelineHistoryFlags(args []string) PipelineHistoryFlags {
	return PipelineHistoryFlags{
		Limit:      parseHistoryLimit(args),
		Repo:       extractFlagVal(args, "--repo"),
		Branch:     extractFlagVal(args, "--branch"),
		IsDetailed: hasDetailedArg(args),
		IsJSON:     hasArgFlag(args, "--json"),
	}
}

func parseHistoryLimit(args []string) int {
	for i, a := range args {
		if val, isOk := tryExtractLimitVal(args, i, a); isOk {
			return val
		}
	}

	return 5
}

func tryExtractLimitVal(args []string, i int, a string) (int, bool) {
	isLimitFlag := a == "-n" || a == "--limit"
	hasValue := i+1 < len(args)
	if isLimitFlag && hasValue {
		val, err := strconv.Atoi(args[i+1])
		return parsePositiveInt(val, err)
	}

	return 0, false
}

func parsePositiveInt(val int, err error) (int, bool) {
	if err == nil && val > 0 {
		return val, true
	}

	return 0, false
}

func resolveHistoryRepo(flagRepo string) string {
	if len(flagRepo) > 0 {
		return flagRepo
	}

	return resolveCurrentRepoSlug()
}

func renderPipelineHistory(repo string, groups []CommitPipelineGroup, flags PipelineHistoryFlags) error {
	filtered := filterGroupsByBranch(groups, flags.Branch)
	limited := capCommitGroups(filtered, flags.Limit)
	if flags.IsJSON {
		return printHistoryJSON(limited)
	}

	renderHistoryTerminalTree(repo, limited, flags.IsDetailed)

	return nil
}

func filterGroupsByBranch(groups []CommitPipelineGroup, branch string) []CommitPipelineGroup {
	if len(branch) == 0 {
		return groups
	}

	out := make([]CommitPipelineGroup, 0, len(groups))
	for _, g := range groups {
		if strings.EqualFold(g.HeadBranch, branch) {
			out = append(out, g)
		}
	}

	return out
}

func printHistoryJSON(groups []CommitPipelineGroup) error {
	data, err := json.MarshalIndent(groups, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}

func renderHistoryTerminalTree(repo string, groups []CommitPipelineGroup, isDetailed bool) {
	if len(groups) == 0 {
		fmt.Printf("No pipeline runs found for %s.\n", repo)

		return
	}

	printHistoryTreeHeader(repo, len(groups))
	for i, g := range groups {
		renderCommitGroupNode(g, i, isDetailed)
	}

	renderHistoryFooter(groups, repo)
}

func printHistoryTreeHeader(repo string, count int) {
	fmt.Printf("\n  %s● Pipeline Commit History (Last %d Commits) [%s]:%s\n\n",
		constants.ColorCyan, count, repo, constants.ColorReset)
}

func renderCommitGroupNode(g CommitPipelineGroup, index int, isDetailed bool) {
	offsetLabel := formatOffsetLabel(index)
	badge := formatCommitGroupBadge(g.Conclusion, g.Status)
	shortSha := truncateHistoryStr(g.HeadSha, 7)
	timeAgo := formatRunTimestamp(g.CreatedAt)
	fmt.Printf("  %-8s %-7s  %-12s  %s  %-10s\n",
		offsetLabel, shortSha, truncateHistoryStr(g.HeadBranch, 12), badge, timeAgo)
	renderWorkflowsSubtree(g.Workflows, isDetailed)
	fmt.Println()
}

func formatOffsetLabel(index int) string {
	if index == 0 {
		return "[latest]"
	}

	return fmt.Sprintf("[-%d]", index)
}

func formatCommitGroupBadge(conclusion, status string) string {
	if isFailingConclusion(conclusion) {
		return constants.ColorRed + "✖ FAIL   " + constants.ColorReset
	}
	switch conclusion {
	case "success":
		return constants.ColorGreen + "● PASS   " + constants.ColorReset
	}
	if status == "in_progress" || status == "queued" {
		return constants.ColorYellow + "● RUNNING" + constants.ColorReset
	}

	return "◌ " + conclusion
}

func renderWorkflowsSubtree(workflows []CommitWorkflowItem, isDetailed bool) {
	total := len(workflows)
	for i, wf := range workflows {
		isLast := i == total-1
		renderWorkflowBranch(wf, isLast)
	}
}

func renderWorkflowBranch(wf CommitWorkflowItem, isLast bool) {
	connector := "├──"
	if isLast {
		connector = "└──"
	}
	badge := formatWorkflowTreeBadge(wf.Conclusion, wf.Status)
	dur := formatDurationSeconds(wf.Duration)
	dots := computeTreeDots(wf.Name, len(dur))
	fmt.Printf("    %s %s (#%d) %s %s (%s)\n",
		connector, wf.Name, wf.DatabaseId, dots, badge, dur)
}

func formatWorkflowTreeBadge(conclusion, status string) string {
	if isFailingConclusion(conclusion) {
		return constants.ColorRed + "FAIL" + constants.ColorReset
	}
	switch conclusion {
	case "success":
		return constants.ColorGreen + "PASS" + constants.ColorReset
	}
	if status == "in_progress" || status == "queued" {
		return constants.ColorYellow + "RUNNING" + constants.ColorReset
	}

	return conclusion
}

func computeTreeDots(name string, durLen int) string {
	target := 42 - len(name) - durLen
	if target < 4 {
		target = 4
	}

	return strings.Repeat(".", target)
}

func renderHistoryFooter(groups []CommitPipelineGroup, repo string) {
	stats := calculateHistoryStats(groups)
	fmt.Println("  ────────────────────────────────────────────────────────────────────────")
	fmt.Printf("  Summary: %d commits | %d Passed | %d Failed | %d In-Progress\n",
		stats.Total, stats.Passed, stats.Failed, stats.Running)
	fmt.Println("  Inspect errors:  gitmap pipeline errors -1 (or -2, -3)")
	fmt.Println("  View logs:       gitmap pipeline logs -1 (or -2, -3, <sha>)")
	printPipelineDbFooter(repo)
	fmt.Println()
}

type historySummaryStats struct {
	Total   int
	Passed  int
	Failed  int
	Running int
}

func calculateHistoryStats(groups []CommitPipelineGroup) historySummaryStats {
	var s historySummaryStats
	s.Total = len(groups)
	for _, g := range groups {
		tallyGroupStat(&s, g.Conclusion)
	}

	return s
}

func tallyGroupStat(s *historySummaryStats, conclusion string) {
	switch conclusion {
	case "success":
		s.Passed++
	case "failure":
		s.Failed++
	default:
		s.Running++
	}
}

func printPipelineDbFooter(repo string) {
	summary := GetPipelineDbSummary(repo)
	if summary != "" {
		fmt.Printf("  Pipeline DB:     %s\n", summary)
	}
}
