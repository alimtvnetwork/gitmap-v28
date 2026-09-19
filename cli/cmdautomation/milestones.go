package cmdautomation

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// RunMilestones queries GitHub milestone issues/PRs and formats release notes.
func RunMilestones(opts MilestonesOptions) MilestonesResultMonad {
	start := time.Now()
	root := resolveAuditDir(opts.Dir)
	id := resolveMilestoneId(opts.MilestoneId)
	issues := queryMilestoneIssues(root, id)
	notes := formatMilestoneReleaseNotes(id, issues)
	isWritten := writeNotesIfRequested(opts.OutputPath, notes)
	isFailedWrite := opts.OutputPath != "" && isWritten == false
	if isFailedWrite {
		err := apperror.NewExecutionError("failed to write release notes to output path")
		return result.Fail[MilestonesResult](err)
	}
	res := buildMilestonesResult(id, issues, notes, time.Since(start))
	return result.Ok(res)
}

func resolveMilestoneId(id string) string {
	hasId := id != ""
	if hasId {
		return id
	}
	return "current"
}

func queryMilestoneIssues(root, milestoneId string) []MilestoneIssue {
	ghIssues := queryGhCli(milestoneId)
	hasGhIssues := len(ghIssues) > 0
	if hasGhIssues {
		return ghIssues
	}
	return scanLocalMilestoneItems(root, milestoneId)
}

func queryGhCli(milestoneId string) []MilestoneIssue {
	cmd := exec.Command("gh", "issue", "list", "--milestone", milestoneId, "--state", "all", "--json", "number,title,state,url")
	out, err := cmd.Output()
	hasErr := err != nil
	if hasErr {
		return nil
	}
	return parseGhIssuesJson(out)
}

func parseGhIssuesJson(data []byte) []MilestoneIssue {
	var raw []map[string]any
	hasErr := json.Unmarshal(data, &raw) != nil
	if hasErr {
		return nil
	}
	var issues []MilestoneIssue
	for _, item := range raw {
		issues = append(issues, convertRawIssue(item))
	}
	return issues
}

func convertRawIssue(item map[string]any) MilestoneIssue {
	num, _ := item["number"].(float64)
	title, _ := item["title"].(string)
	state, _ := item["state"].(string)
	url, _ := item["url"].(string)
	isClosed := strings.EqualFold(state, "closed")
	category := categorizeIssue(title)
	return MilestoneIssue{
		Number: int(num), Title: title, State: state,
		Category: category, Url: url, IsClosed: isClosed,
	}
}

func scanLocalMilestoneItems(root, milestoneId string) []MilestoneIssue {
	plansDir := filepath.Join(root, ".ai-memory", "plans", "completed")
	entries, err := os.ReadDir(plansDir)
	hasErr := err != nil
	if hasErr {
		return nil
	}
	return filterPlansForMilestone(entries, plansDir, milestoneId)
}

func filterPlansForMilestone(entries []os.DirEntry, dir, milestoneId string) []MilestoneIssue {
	var issues []MilestoneIssue
	for idx, entry := range entries {
		isMatch := isMilestonePlanMatch(entry.Name(), milestoneId)
		if isMatch {
			title := extractPlanTitle(filepath.Join(dir, entry.Name()))
			issues = append(issues, makeLocalIssue(idx+1, title, entry.Name()))
		}
	}
	return issues
}

func isMilestonePlanMatch(name, milestoneId string) bool {
	isCurrent := milestoneId == "current" || milestoneId == "all"
	if isCurrent {
		return strings.HasSuffix(name, ".md")
	}
	return strings.HasPrefix(name, milestoneId)
}

func extractPlanTitle(path string) string {
	content := readFileContent(path)
	reH1 := regexp.MustCompile(`(?m)^#\s+(.+)`)
	m := reH1.FindStringSubmatch(content)
	hasMatch := len(m) > 1
	if hasMatch {
		return strings.TrimSpace(m[1])
	}
	return filepath.Base(path)
}

func makeLocalIssue(num int, title, name string) MilestoneIssue {
	return MilestoneIssue{
		Number:   num,
		Title:    title,
		State:    "closed",
		Category: categorizeIssue(title),
		Url:      fmt.Sprintf(".ai-memory/plans/completed/%s", name),
		IsClosed: true,
	}
}

func categorizeIssue(title string) string {
	lower := strings.ToLower(title)
	isAdded := strings.Contains(lower, "feat") || strings.Contains(lower, "add")
	if isAdded {
		return "Added"
	}
	isFixed := strings.Contains(lower, "fix") || strings.Contains(lower, "bug")
	if isFixed {
		return "Fixed"
	}
	isRemoved := strings.Contains(lower, "remove") || strings.Contains(lower, "purge")
	if isRemoved {
		return "Removed"
	}
	return "Changed"
}

func formatMilestoneReleaseNotes(title string, issues []MilestoneIssue) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("## Milestone %s Release Notes\n\n", title))
	categories := []string{"Added", "Fixed", "Changed", "Removed"}
	for _, cat := range categories {
		appendCategorySection(&sb, cat, issues)
	}
	return sb.String()
}

func appendCategorySection(sb *strings.Builder, category string, issues []MilestoneIssue) {
	var catItems []MilestoneIssue
	for _, issue := range issues {
		isCatMatch := issue.Category == category && issue.IsClosed
		if isCatMatch {
			catItems = append(catItems, issue)
		}
	}
	hasItems := len(catItems) > 0
	if hasItems {
		sb.WriteString(fmt.Sprintf("### %s\n\n", category))
		for _, item := range catItems {
			sb.WriteString(fmt.Sprintf("- %s (#%d)\n", item.Title, item.Number))
		}
		sb.WriteString("\n")
	}
}

func writeNotesIfRequested(outPath, content string) bool {
	hasOutPath := outPath != ""
	if hasOutPath {
		dir := filepath.Dir(outPath)
		_ = os.MkdirAll(dir, 0755)
		return os.WriteFile(outPath, []byte(content), 0644) == nil
	}
	return true
}

func countClosedIssues(issues []MilestoneIssue) int {
	closed := 0
	for _, issue := range issues {
		if issue.IsClosed {
			closed++
		}
	}
	return closed
}

func buildMilestonesResult(id string, issues []MilestoneIssue, notes string, dur time.Duration) MilestonesResult {
	closedCount := countClosedIssues(issues)
	openCount := len(issues) - closedCount
	return MilestonesResult{
		MilestoneId: id, Title: fmt.Sprintf("Milestone %s", id),
		TotalItems: len(issues), ClosedCount: closedCount, OpenCount: openCount,
		Items: issues, ReleaseNotes: notes, Duration: dur, IsSuccess: true,
	}
}
