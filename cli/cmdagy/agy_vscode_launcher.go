package cmdagy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func launchVSCodeForLastProjects(projCount, promptLimit int) error {
	projects := discoverRecentCommitProjects(projCount)
	tempFile, genErr := generateVSCodeInspectionFile(projects, promptLimit)
	if genErr != nil {
		return genErr
	}
	fmt.Printf("\n  %s● Generated prompt history file:%s %s\n", constants.ColorCyan, constants.ColorReset, tempFile)
	if launchErr := spawnNonAdminVSCode(tempFile); launchErr != nil {
		fmt.Printf("  %s⚠ Could not launch VS Code directly: %v%s\n", constants.ColorYellow, launchErr, constants.ColorReset)
		fmt.Printf("    Open manually: code -n %s\n\n", tempFile)
		return nil
	}
	fmt.Printf("  %s✔ Launched new non-admin VS Code window (code -n).%s\n\n", constants.ColorGreen, constants.ColorReset)

	return nil
}

func discoverRecentCommitProjects(count int) []string {
	allPrompts := CollectAllPrompts()
	seen := make(map[string]bool)
	var projects []string
	for _, p := range allPrompts {
		if p.Workspace == "" || seen[p.Workspace] {
			continue
		}
		seen[p.Workspace] = true
		projects = append(projects, p.Workspace)
		if len(projects) >= count {
			break
		}
	}

	return projects
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
	sb.WriteString(fmt.Sprintf("## Project: %s\n\n", filepath.Base(proj)))
	sb.WriteString(fmt.Sprintf("Path: `%s`\n\n", proj))
	prompts := CollectPromptsForWorkspace(proj)
	if len(prompts) > promptLimit {
		prompts = prompts[:promptLimit]
	}
	for i, p := range prompts {
		sb.WriteString(fmt.Sprintf("### Prompt %d (%s)\n\n```text\n%s\n```\n\n", i+1, p.CreatedAt.Format("2006-01-02 15:04"), p.Content))
	}
	sb.WriteString("---\n\n")
}

func spawnNonAdminVSCode(filePath string) error {
	cmd := exec.Command("code", "-n", filePath)
	cmd.Env = filterAdminEnv(os.Environ())

	return cmd.Start()
}

func filterAdminEnv(env []string) []string {
	var filtered []string
	for _, e := range env {
		if !strings.HasPrefix(strings.ToUpper(e), "ELEVATED=") {
			filtered = append(filtered, e)
		}
	}

	return filtered
}
