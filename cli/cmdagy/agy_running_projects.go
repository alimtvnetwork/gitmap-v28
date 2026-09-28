package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/spf13/cobra"
)

// RunningProjectRecord models a project hosting active or queued prompts.
type RunningProjectRecord struct {
	ProjectName      string `json:"projectName"`
	ProjectPath      string `json:"projectPath"`
	ProjectId        string `json:"projectId,omitempty"`
	Node             string `json:"node,omitempty"`
	Host             string `json:"host,omitempty"`
	HasActivePrompt  bool   `json:"hasActivePrompt"`
	HasQueuedPrompt  bool   `json:"hasQueuedPrompt"`
	ActivePromptFile string `json:"activePromptFile,omitempty"`
	QueuedPromptFile string `json:"queuedPromptFile,omitempty"`
	PromptPreview    string `json:"promptPreview,omitempty"`
	QueuedCount      int    `json:"queuedCount,omitempty"`
	Status           string `json:"status"`
	UpdatedAt        string `json:"updatedAt,omitempty"`
}

var (
	runningProjectsJSON bool
	runningProjectsSSH  bool
	runningProjectsFile string
)

// AgyRunningProjectsCmd inspects projects with active or enqueued prompts.
var AgyRunningProjectsCmd = &cobra.Command{
	Use:     "running-projects [ls]",
	Aliases: []string{"runningprojects", "rp"},
	Short:   "List projects hosting active or queued Antigravity prompts",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunRunningProjectsCLI(args)
	},
}

func init() {
	AgyRunningProjectsCmd.Flags().BoolVar(&runningProjectsJSON, "json", false, "Output results in JSON format")
	AgyRunningProjectsCmd.Flags().BoolVar(&runningProjectsSSH, "ssh", false, "Aggregate running projects across SSH cluster nodes")
	AgyRunningProjectsCmd.Flags().StringVarP(&runningProjectsFile, "file", "f", "", "Export results to file (.json or .db)")
	_ = cobra.MarkFlagFilename(AgyRunningProjectsCmd.Flags(), "file")
	AgyRunningProjectsCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return []string{"ls\tList projects with active or queued prompts"}, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	AgyCmd.AddCommand(AgyRunningProjectsCmd)
}

// RunRunningProjectsCLI handles running-projects command execution.
func RunRunningProjectsCLI(args []string) error {
	if checkRunningProjectsHelp(args) {
		printRunningProjectsHelp()
		return nil
	}
	projects, err := DiscoverRunningProjects()
	if err != nil {
		return err
	}
	if runningProjectsSSH {
		return AggregateSSHRunningProjects(projects, runningProjectsJSON, runningProjectsFile)
	}
	return outputRunningProjects(projects, runningProjectsJSON, runningProjectsFile)
}

func checkRunningProjectsHelp(args []string) bool {
	for _, a := range args {
		sub := strings.ToLower(a)
		if sub == "help" || sub == "--help" || sub == "-h" {
			return true
		}
	}
	return false
}

func printRunningProjectsHelp() {
	fmt.Println()
	fmt.Printf("  %sAntigravity Running Projects Discovery%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("  Usage:")
	fmt.Println("    gitmap agy running-projects [ls] [--json] [-f <file>] [--ssh]")
	fmt.Println("    gitmap agy running-projects help")
	fmt.Println()
	fmt.Println("  Flags:")
	fmt.Println("    --json          Output in structured JSON format")
	fmt.Println("    --ssh           Query and aggregate projects across cluster SSH fleet")
	fmt.Println("    -f, --file      Export results to output file (.json or .db)")
	fmt.Println()
}

// DiscoverRunningProjects scans candidate workspaces for active or queued prompts.
func DiscoverRunningProjects() ([]RunningProjectRecord, error) {
	wsMap := collectCandidateWorkspaces()
	var list []RunningProjectRecord
	for ws, name := range wsMap {
		rec, isRunning := inspectWorkspaceRunningStatus(ws, name)
		if isRunning {
			list = append(list, rec)
		}
	}
	return list, nil
}

func inspectWorkspaceRunningStatus(ws, name string) (RunningProjectRecord, bool) {
	activePrompt, activeFile, hasActive := checkActivePromptFile(ws)
	qSummary, hasQueue := inspectSingleWorkspaceQueue(ws, name)
	if !hasActive && !hasQueue {
		return RunningProjectRecord{}, false
	}
	return buildRunningRecord(ws, name, activePrompt, activeFile, hasActive, qSummary, hasQueue), true
}

func buildRunningRecord(ws, name, activePrompt, activeFile string, hasActive bool, q AgyWorkspaceQueueSummary, hasQueue bool) RunningProjectRecord {
	return RunningProjectRecord{
		ProjectName:      name,
		ProjectPath:      ws,
		HasActivePrompt:  hasActive,
		HasQueuedPrompt:  hasQueue,
		ActivePromptFile: activeFile,
		QueuedPromptFile: q.QueueFile,
		PromptPreview:    resolvePromptPreview(activePrompt, q),
		QueuedCount:      q.TotalQueued,
		Status:           resolveRunningStatus(hasActive, hasQueue),
		UpdatedAt:        time.Now().UTC().Format(time.RFC3339),
	}
}

func checkActivePromptFile(ws string) (string, string, bool) {
	paths := []string{
		filepath.Join(ws, ".ai-memory", "temp", "active-agy-pipeline-fix-prompt.txt"),
		filepath.Join(ws, "active-agy-pipeline-fix-prompt.txt"),
	}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		trimmed := strings.TrimSpace(string(data))
		if err == nil && len(trimmed) > 0 {
			return trimmed, p, true
		}
	}
	return "", "", false
}

func resolveRunningStatus(hasActive, hasQueue bool) string {
	if hasActive && hasQueue {
		return "RUNNING+QUEUED"
	}
	if hasActive {
		return "RUNNING"
	}
	return "QUEUED"
}

func resolvePromptPreview(active string, q AgyWorkspaceQueueSummary) string {
	if len(active) > 0 {
		return CompactWords(active, 12)
	}
	if len(q.QueuedItems) > 0 {
		return CompactWords(q.QueuedItems[0].Prompt, 12)
	}
	return ""
}

func outputRunningProjects(projects []RunningProjectRecord, isJSON bool, filePath string) error {
	if filePath != "" {
		return writeRunningProjectsToFile(projects, filePath)
	}
	if isJSON {
		return renderRunningProjectsJSON(projects)
	}
	RenderRunningProjectsTable(projects)
	return nil
}

func renderRunningProjectsJSON(projects []RunningProjectRecord) error {
	data, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal running projects JSON")
	}
	fmt.Println(string(data))
	return nil
}

func writeRunningProjectsToFile(projects []RunningProjectRecord, path string) error {
	data, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal running projects")
	}
	if writeErr := os.WriteFile(path, data, 0644); writeErr != nil {
		return apperror.WrapSimple(writeErr, "write running projects file")
	}
	fmt.Printf("  ✔ Exported %d running project(s) to %s\n", len(projects), path)
	return nil
}

// RenderRunningProjectsTable prints a formatted terminal table of running projects.
func RenderRunningProjectsTable(projects []RunningProjectRecord) {
	fmt.Println()
	fmt.Printf("  %s%s ANTIGRAVITY RUNNING PROJECTS (%d active) %s%s\n",
		constants.ColorCyan, "╔════", len(projects), "════╗", constants.ColorReset)
	fmt.Printf("  %s%-20s  %-16s  %-8s  %-35s  %s%s\n",
		constants.ColorWhite, "PROJECT", "STATUS", "QUEUED", "ACTIVE PROMPT PREVIEW", "PATH", constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorDim, strings.Repeat("─", 95), constants.ColorReset)
	if len(projects) == 0 {
		fmt.Printf("  %sNo active or queued Antigravity projects found.%s\n\n", constants.ColorDim, constants.ColorReset)
		return
	}
	for _, p := range projects {
		printRunningProjectRow(p)
	}
	fmt.Println()
}

func printRunningProjectRow(p RunningProjectRecord) {
	color := resolveRunningColor(p.HasActivePrompt)
	name := truncateRunningStr(p.ProjectName, 18)
	preview := truncateRunningStr(p.PromptPreview, 33)
	pathTrunc := truncateRunningPath(p.ProjectPath, 25)
	fmt.Printf("  %-20s  %s%-16s%s  %-8d  %-35s  %s\n",
		name, color, p.Status, constants.ColorReset, p.QueuedCount, preview, pathTrunc)
}

func resolveRunningColor(hasActive bool) string {
	if hasActive {
		return constants.ColorGreen + "\033[1m"
	}
	return constants.ColorYellow
}

func truncateRunningStr(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}

func truncateRunningPath(p string, maxLen int) string {
	if len(p) > maxLen {
		return "..." + p[len(p)-maxLen+3:]
	}
	return p
}
