package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/spf13/cobra"
)

// RunningProjectRecord models a project hosting active or queued prompts.
type RunningProjectRecord struct {
	SequenceId       string `json:"sequenceId,omitempty"`
	ProjectName      string `json:"projectName"`
	ProjectAlias     string `json:"projectAlias,omitempty"`
	ProjectPath      string `json:"projectPath"`
	ProjectId        string `json:"projectId,omitempty"`
	ConversationId   string `json:"conversationId,omitempty"`
	PromptSeqId      string `json:"promptSeqId,omitempty"`
	Node             string `json:"node,omitempty"`
	Host             string `json:"host,omitempty"`
	HasActivePrompt  bool   `json:"hasActivePrompt"`
	HasQueuedPrompt  bool   `json:"hasQueuedPrompt"`
	ActivePromptFile string `json:"activePromptFile,omitempty"`
	QueuedPromptFile string `json:"queuedPromptFile,omitempty"`
	PromptPreview    string `json:"promptPreview,omitempty"`
	FullPrompt       string `json:"fullPrompt,omitempty"`
	QueuedCount      int    `json:"queuedCount,omitempty"`
	Status           string `json:"status"`
	UpdatedAt        string `json:"updatedAt,omitempty"`
}

var (
	runningProjectsJSON      bool
	runningProjectsSSH       bool
	runningProjectsFile      string
	runningProjectsWordCount int
)

// AgyRunningProjectsCmd inspects projects with active or enqueued prompts.
var AgyRunningProjectsCmd = &cobra.Command{
	Use:                "running-projects [ls|prompts ls]",
	Aliases:            []string{"runningprojects", "rp"},
	Short:              "List projects hosting active or queued Antigravity prompts",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunRunningProjectsCLI(args)
	},
}

func init() {
	AgyRunningProjectsCmd.Flags().BoolVar(&runningProjectsJSON, "json", false, "Output results in JSON format")
	AgyRunningProjectsCmd.Flags().BoolVar(&runningProjectsSSH, "ssh", false, "Aggregate running projects across SSH cluster nodes")
	AgyRunningProjectsCmd.Flags().StringVarP(&runningProjectsFile, "file", "f", "", "Export results to file (.json or .db)")
	AgyRunningProjectsCmd.Flags().IntVar(&runningProjectsWordCount, "wordcount", 200, "Maximum words per prompt in tree view")
	AgyRunningProjectsCmd.Flags().IntVar(&runningProjectsWordCount, "wc", 200, "Maximum words per prompt in tree view")
	_ = cobra.MarkFlagFilename(AgyRunningProjectsCmd.Flags(), "file")
	AgyCmd.AddCommand(AgyRunningProjectsCmd)
}

// RunRunningProjectsCLI handles running-projects command execution.
func RunRunningProjectsCLI(args []string) error {
	parseRunningProjectsFlags(args)
	if isPromptsSubcommand(args) {
		return RunRunningProjectsPromptsTreeCLI(args[1:])
	}
	if checkRunningProjectsHelp(args) {
		printRunningProjectsHelp()
		return nil
	}
	projects, err := DiscoverRunningProjects()
	if err != nil {
		return err
	}
	return dispatchRunningProjectsOutput(projects)
}

func dispatchRunningProjectsOutput(projects []RunningProjectRecord) error {
	if runningProjectsSSH {
		return AggregateSSHRunningProjects(projects, runningProjectsJSON, runningProjectsFile)
	}
	return outputRunningProjects(projects, runningProjectsJSON, runningProjectsFile)
}

func parseRunningProjectsFlags(args []string) {
	runningProjectsJSON, runningProjectsSSH = false, false
	runningProjectsFile, runningProjectsWordCount = "", 200
	for i := 0; i < len(args); i++ {
		i += stepRunningProjectsFlag(args, i)
	}
}

func stepRunningProjectsFlag(args []string, i int) int {
	tok := strings.TrimSpace(args[i])
	if tok == "--json" || tok == "-j" || tok == "-json" {
		runningProjectsJSON = true
		return 0
	}
	if tok == "--ssh" || tok == "-s" || tok == "-ssh" {
		runningProjectsSSH = true
		return 0
	}
	return stepRunningProjectsValueFlag(args, i, tok)
}

func stepRunningProjectsValueFlag(args []string, i int, tok string) int {
	if i+1 >= len(args) {
		return 0
	}
	if tok == "--file" || tok == "-f" || tok == "-file" {
		runningProjectsFile = args[i+1]
		return 1
	}
	if isWordCountFlag(tok) {
		runningProjectsWordCount = parseIntOrDefault(args[i+1], 200)
		return 1
	}
	return 0
}

func isWordCountFlag(tok string) bool {
	return tok == "--wc" || tok == "-wc" || tok == "--wordcount" || tok == "-wordcount" || tok == "--ww" || tok == "-ww"
}

func isPromptsSubcommand(args []string) bool {
	if len(args) == 0 {
		return false
	}
	first := strings.ToLower(strings.TrimSpace(args[0]))
	return first == "prompts" || first == "prompt"
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
	fmt.Println("    gitmap agy running-projects prompts ls [--wc 200] [--json] [-f <file>]")
	fmt.Println("    gitmap agy running-projects help")
	fmt.Println()
	fmt.Println("  Flags:")
	fmt.Println("    --json          Output in structured JSON format")
	fmt.Println("    --ssh           Query and aggregate projects across cluster SSH fleet")
	fmt.Println("    -f, --file      Export results to output file (.json or .db)")
	fmt.Println("    --wc, --wordcount  Word limit for prompt sub-items (default 200)")
	fmt.Println()
}

// DiscoverRunningProjects scans candidate workspaces for active or queued prompts and caches 24h sequences.
func DiscoverRunningProjects() ([]RunningProjectRecord, error) {
	wsMap := collectCandidateWorkspaces()
	var list []RunningProjectRecord
	for ws, name := range wsMap {
		rec, isRunning := inspectWorkspaceRunningStatus(ws, name)
		if isRunning {
			list = append(list, rec)
		}
	}
	list = enrichAndCacheRunningProjects(list)
	return list, nil
}

func enrichAndCacheRunningProjects(list []RunningProjectRecord) []RunningProjectRecord {
	var cacheEntries []store.AgySequenceRecord
	now := time.Now().Unix()
	for i := range list {
		seqNum := i + 1
		list[i].SequenceId = strconv.Itoa(seqNum)
		list[i].PromptSeqId = fmt.Sprintf("P%d", seqNum)
		cacheEntries = appendRunningCachePair(cacheEntries, list[i], seqNum, now)
	}
	_, _ = store.EnsureAndGetSequenceCache(cacheEntries, false)
	return list
}

func appendRunningCachePair(entries []store.AgySequenceRecord, r RunningProjectRecord, seqNum int, now int64) []store.AgySequenceRecord {
	projEntry := store.AgySequenceRecord{
		SeqId: r.SequenceId, SeqNum: seqNum, EntryType: "project",
		ProjectId: r.ProjectId, ProjectAlias: r.ProjectAlias, ProjectPath: r.ProjectPath,
		ConversationId: r.ConversationId, PromptSnippet: r.PromptPreview,
		CreatedAt: now, ExpiresAt: now + store.SequenceCacheTTLSeconds,
	}
	promptEntry := store.AgySequenceRecord{
		SeqId: r.PromptSeqId, SeqNum: seqNum, EntryType: "prompt",
		ProjectId: r.ProjectId, ProjectAlias: r.ProjectAlias, ProjectPath: r.ProjectPath,
		ConversationId: r.ConversationId, PromptSnippet: r.PromptPreview,
		CreatedAt: now, ExpiresAt: now + store.SequenceCacheTTLSeconds,
	}
	return append(entries, projEntry, promptEntry)
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
	projID := deriveOrLookupProjectID(ws)
	alias := strings.ToLower(filepath.Base(filepath.Clean(ws)))
	convID := resolveWorkspaceConvID(ws)
	fullText := resolveFullPromptText(activePrompt, q)
	return RunningProjectRecord{
		ProjectName:      name,
		ProjectAlias:     alias,
		ProjectPath:      ws,
		ProjectId:        projID,
		ConversationId:   convID,
		HasActivePrompt:  hasActive,
		HasQueuedPrompt:  hasQueue,
		ActivePromptFile: activeFile,
		QueuedPromptFile: q.QueueFile,
		PromptPreview:    CompactWords(fullText, 12),
		FullPrompt:       fullText,
		QueuedCount:      q.TotalQueued,
		Status:           resolveRunningStatus(hasActive, hasQueue),
		UpdatedAt:        time.Now().UTC().Format(time.RFC3339),
	}
}

func resolveWorkspaceConvID(ws string) string {
	conv, err := SelectMatchingConversation(ws)
	if err == nil && conv.ID != "" {
		return conv.ID
	}
	return "default-conv"
}

func resolveFullPromptText(active string, q AgyWorkspaceQueueSummary) string {
	if len(active) > 0 {
		return active
	}
	if len(q.QueuedItems) > 0 {
		return q.QueuedItems[0].Prompt
	}
	return ""
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

// RenderRunningProjectsTable prints a formatted terminal table of running projects with 24h sequence IDs.
func RenderRunningProjectsTable(projects []RunningProjectRecord) {
	fmt.Println()
	fmt.Printf("  %s%s ANTIGRAVITY RUNNING PROJECTS (%d active) %s%s\n",
		constants.ColorCyan, "╔════", len(projects), "════╗", constants.ColorReset)
	fmt.Printf("  %s%-5s  %-10s  %-16s  %-14s  %-28s  %s%s\n",
		constants.ColorWhite, "SEQ", "PROJ ID", "ALIAS", "STATUS", "[PROJ_ID | CONV_ID | SEQ]", "PATH", constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorDim, strings.Repeat("─", 105), constants.ColorReset)
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
	seqLabel := "#" + p.SequenceId
	shortID := shortProjectId(p.ProjectId)
	alias := truncateRunningStr(p.ProjectAlias, 15)
	bracket := fmt.Sprintf("[%s | %s | Seq:%s]", shortID, shortProjectId(p.ConversationId), p.PromptSeqId)
	fmt.Printf("  %-5s  %-10s  %-16s  %s%-14s%s  %-28s  %s\n",
		seqLabel, shortID, alias, color, p.Status, constants.ColorReset, bracket, p.ProjectPath)
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
