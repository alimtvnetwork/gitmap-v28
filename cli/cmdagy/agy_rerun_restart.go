package cmdagy

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/lipgloss"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type rerunRestartPlan struct {
	Project        AgyProject
	SequenceNumber int
	ConvID         string
	WorkspacePath  string
	PromptText     string
	MediaCount     int
	MediaPaths     []string
	OldPID         int
	IsRestarted    bool
}

// RestartAndRerunProject resolves a target project, extracts its prompt and pictures, restarts the IDE, and replays.
func RestartAndRerunProject(target string, isRestart, isDryRun bool, templateID string) error {
	projects, err := loadSortedProjects()
	if err != nil {
		return err
	}

	plan, planErr := buildRerunRestartPlan(projects, target, templateID)
	if planErr != nil {
		return planErr
	}

	renderRerunPlanHeader(plan)
	if isDryRun {
		fmt.Println("  • [dry-run] Preview complete. Skipping IDE restart and prompt injection.")
		return nil
	}

	if isRestart {
		executeIDERestart(&plan)
	}

	return dispatchRerunPayload(plan)
}

func loadSortedProjects() ([]AgyProject, error) {
	dirPath, err := getProjectsDirPath()
	if err != nil {
		return nil, apperror.WrapSimple(err, "projects dir")
	}

	projects, err := loadAllAgyProjects(dirPath)
	if err != nil || len(projects) == 0 {
		return nil, apperror.NewSimple("no Antigravity projects configured in ~/.gemini/config/projects", "E9030")
	}

	sortAgyProjects(projects, "recent")

	return projects, nil
}

func buildRerunRestartPlan(projects []AgyProject, target, templateID string) (rerunRestartPlan, error) {
	proj, seq := resolveProjectFromTarget(projects, target)
	ws := proj.GetPath()
	if ws == "" {
		return rerunRestartPlan{}, fmt.Errorf("project %q has no valid workspace folder", proj.Name)
	}

	promptEntry, convID := extractLastPromptForProject(ws)
	payload := composeRerunPayload(promptEntry, templateID)
	mediaPaths := extractMediaURIs(promptEntry.Media)

	return rerunRestartPlan{
		Project:        proj,
		SequenceNumber: seq,
		ConvID:         convID,
		WorkspacePath:  ws,
		PromptText:     payload,
		MediaCount:     len(promptEntry.Media),
		MediaPaths:     mediaPaths,
	}, nil
}

func resolveProjectFromTarget(projects []AgyProject, target string) (AgyProject, int) {
	trimmed := strings.TrimSpace(target)
	if trimmed == "" {
		return projects[0], 1
	}

	seq, err := strconv.Atoi(trimmed)
	if err == nil && seq >= 1 && seq <= len(projects) {
		return projects[seq-1], seq
	}

	matched, matchErr := ResolveAgyProjectTargets([]string{target}, projects)
	if matchErr == nil && len(matched) > 0 {
		return matched[0], findProjectSequence(projects, matched[0].ID)
	}

	return projects[0], 1
}

func findProjectSequence(projects []AgyProject, id string) int {
	for i, p := range projects {
		if p.ID == id {
			return i + 1
		}
	}

	return 1
}

func extractLastPromptForProject(workspace string) (AgyPromptEntry, string) {
	entry, convID, isFound := tryExtractMatchingConvPrompt(workspace)
	if isFound {
		return entry, convID
	}

	all := CollectPromptsForWorkspace(workspace)
	if len(all) > 0 {
		return all[0], all[0].ConvID
	}

	return AgyPromptEntry{Content: "Rerun previous task and verify completion."}, "default"
}

func tryExtractMatchingConvPrompt(workspace string) (AgyPromptEntry, string, bool) {
	conv, err := SelectMatchingConversation(workspace)
	if err != nil || conv.ID == "" {
		return AgyPromptEntry{}, "", false
	}

	brainDir, bErr := GetBrainLogsDirPath()
	if bErr != nil {
		return AgyPromptEntry{}, "", false
	}

	prompts := readConvPrompts(brainDir, conv.ID)
	if len(prompts) == 0 {
		return AgyPromptEntry{}, "", false
	}

	return prompts[len(prompts)-1], conv.ID, true
}

func composeRerunPayload(promptEntry AgyPromptEntry, templateID string) string {
	tplContent := resolveRerunTemplate(templateID)
	basePrompt := promptEntry.Content
	if tplContent != "" && !strings.Contains(basePrompt, tplContent) {
		basePrompt = tplContent + "\n\n" + basePrompt
	}

	return FormatPromptWithMedia(basePrompt, promptEntry.Media)
}

func extractMediaURIs(media []rawTranscriptMedia) []string {
	var paths []string
	for _, m := range media {
		clean := strings.TrimPrefix(m.URI, "file:///")
		paths = append(paths, filepath.FromSlash(clean))
	}

	return paths
}

func executeIDERestart(plan *rerunRestartPlan) {
	fmt.Println("  [1/2] Terminating active Antigravity IDE process...")
	plan.OldPID = terminateRunningIDE()
	time.Sleep(500 * time.Millisecond)

	fmt.Printf("  [2/2] Launching Antigravity IDE in %s...\n", plan.WorkspacePath)
	ideRes := ResolveAntigravityIDE()
	if ideRes.IsSuccess() {
		_ = launchAntigravityProcess(ideRes.Value, plan.WorkspacePath)
		plan.IsRestarted = true
		time.Sleep(800 * time.Millisecond)
	}
}

func terminateRunningIDE() int {
	procRes := DetectRunningAntigravityIDE()
	if procRes.IsFailure() {
		fmt.Println("    • No running Antigravity IDE process detected.")
		return 0
	}

	pid := procRes.Value.PID
	killIDEProcessByPID(pid)

	return pid
}

func killIDEProcessByPID(pid int) {
	if runtime.GOOS == "windows" {
		_ = exec.Command("taskkill", "/F", "/PID", strconv.Itoa(pid)).Run()
		_ = exec.Command("taskkill", "/F", "/IM", "Antigravity.exe").Run()
		fmt.Printf("    ✓ Terminated Antigravity IDE (PID: %d)\n", pid)
		return
	}

	_ = exec.Command("kill", "-9", strconv.Itoa(pid)).Run()
	fmt.Printf("    ✓ Terminated Antigravity IDE (PID: %d)\n", pid)
}

func dispatchRerunPayload(plan rerunRestartPlan) error {
	title := fmt.Sprintf("Rerun: %s (#%d)", plan.Project.Name, plan.SequenceNumber)
	targetPrompt := filepath.Join(plan.WorkspacePath, activeAgyPromptRelativePath)
	writePromptFile(targetPrompt, plan.PromptText)
	_ = clipboard.WriteAll(plan.PromptText)

	ideProcRes := DetectRunningAntigravityIDE()
	pid := resolveActiveOrZeroPID(ideProcRes)

	res := DispatchPromptToAntigravity(plan.WorkspacePath, targetPrompt, title, plan.PromptText, pid)
	renderRerunResult(res, plan)

	return nil
}

func renderRerunPlanHeader(plan rerunRestartPlan) {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#bd93f9"))
	fmt.Println(headerStyle.Render("\nAntigravity IDE Rerun & Restart Suite:"))
	fmt.Printf("  • Target Project:   #%d - %s\n", plan.SequenceNumber, plan.Project.Name)
	fmt.Printf("  • Workspace Path:   %s\n", plan.WorkspacePath)
	fmt.Printf("  • Conversation:     %s\n", plan.ConvID)
	if plan.MediaCount > 0 {
		fmt.Printf("  • Pictures Attached: %d media item(s)\n", plan.MediaCount)
		for i, mp := range plan.MediaPaths {
			fmt.Printf("      [%d] %s\n", i+1, mp)
		}
	}
	fmt.Println()
}

func renderRerunResult(res AgyInjectionResult, plan rerunRestartPlan) {
	successMark := lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b")).Render("✓")
	fmt.Println()
	if res.IsSuccess {
		fmt.Printf("  %s %s\n", successMark, res.Message)
	} else {
		fmt.Printf("  • Prompt staged in %s (copied to clipboard)\n", plan.WorkspacePath)
	}
	fmt.Printf("  %s Completed prompt replay for project #%d (%s)!\n\n", successMark, plan.SequenceNumber, plan.Project.Name)
}
