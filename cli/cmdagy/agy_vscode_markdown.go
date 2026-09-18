package cmdagy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func discoverRecentCommitProjects(count int) []string {
	allPrompts := CollectAllPrompts()
	seen := make(map[string]bool)
	var projects []string
	for _, p := range allPrompts {
		if appendUniqueWorkspace(&projects, seen, p.Workspace, count) {
			break
		}
	}

	return projects
}

func appendUniqueWorkspace(projs *[]string, seen map[string]bool, ws string, count int) bool {
	if ws == "" || seen[ws] {
		return false
	}
	seen[ws] = true
	*projs = append(*projs, ws)

	return len(*projs) >= count
}

func generateVSCodeInspectionFile(projects []string, promptLimit int) (string, error) {
	tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("gitmap-agy-prompts-%d.md", time.Now().Unix()))
	var sb strings.Builder
	sb.WriteString("# AGY Prompt History & Project Inspection\n\n")
	sb.WriteString(fmt.Sprintf("> Generated: %s\n\n", time.Now().Format(time.RFC1123)))
	for _, proj := range projects {
		appendProjectPromptsMarkdown(&sb, proj, promptLimit)
	}

	return tempPath, os.WriteFile(tempPath, []byte(sb.String()), 0644)
}

func appendProjectPromptsMarkdown(sb *strings.Builder, proj string, promptLimit int) {
	sb.WriteString(fmt.Sprintf("## Project: %s\n\nPath: `%s`\n\n", filepath.Base(proj), proj))
	if commits := fetchProjectRecentCommits(proj); commits != "" {
		sb.WriteString(fmt.Sprintf("### Recent Commits\n\n```text\n%s\n```\n\n", commits))
	}
	prompts := CollectPromptsForWorkspace(proj)
	appendPromptsList(sb, prompts, promptLimit)
	sb.WriteString("---\n\n")
}

func fetchProjectRecentCommits(proj string) string {
	cmd := exec.Command("git", "-C", proj, "log", "-n", "5", "--oneline")
	out, err := cmd.Output()
	if err != nil || len(out) == 0 {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func appendPromptsList(sb *strings.Builder, prompts []AgyPromptEntry, limit int) {
	if len(prompts) > limit {
		prompts = prompts[:limit]
	}
	for i, p := range prompts {
		dateStr := p.CreatedAt.Format("2006-01-02 15:04")
		sb.WriteString(fmt.Sprintf("### Prompt %d (%s)\n\n```text\n%s\n```\n\n", i+1, dateStr, p.Content))
	}
}
