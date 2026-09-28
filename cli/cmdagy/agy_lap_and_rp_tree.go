// Package cmdagy — agy_lap_and_rp_tree.go implements `gitmap agy rp prompts ls` and `gitmap agy last-active-projects` (`lap`).
package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// LapConversationItem represents a single conversation and its truncated prompt inside a project tree.
type LapConversationItem struct {
	ConversationID string   `json:"conversationId"`
	PromptSeqID    string   `json:"promptSeqId"`
	Prompts        []string `json:"prompts"`
	UpdatedAt      string   `json:"updatedAt"`
}

// LapProjectTreeRecord models a project and its active/recent conversation prompts in the last N hours.
type LapProjectTreeRecord struct {
	SequenceID    string                `json:"sequenceId"`
	ProjectID     string                `json:"projectId"`
	ProjectAlias  string                `json:"projectAlias"`
	ProjectName   string                `json:"projectName"`
	ProjectPath   string                `json:"projectPath"`
	Status        string                `json:"status"`
	LastActiveAt  string                `json:"lastActiveAt"`
	Conversations []LapConversationItem `json:"conversations"`
}

var (
	lapLimitFlag     int
	lapOffsetFlag    int
	lapPageFlag      int
	lapWordCountFlag int
	lapJSONFlag      bool
	lapFileFlag      string
)

// AgyLastActiveProjectsCmd lists projects with communication in the last N hours as a tree view.
var AgyLastActiveProjectsCmd = &cobra.Command{
	Use:     "last-active-projects [N] [ls|help]",
	Aliases: []string{"lap", "active-projects"},
	Short:   "List projects with communication in the last N hours (default 24h) with conversation tree",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunLastActiveProjectsCLI(args)
	},
}

func init() {
	AgyLastActiveProjectsCmd.Flags().IntVarP(&lapLimitFlag, "limit", "l", 10, "Maximum number of projects to display (default 10)")
	AgyLastActiveProjectsCmd.Flags().IntVar(&lapOffsetFlag, "offset", 0, "Number of projects to skip")
	AgyLastActiveProjectsCmd.Flags().IntVar(&lapOffsetFlag, "skip", 0, "Number of projects to skip (alias for --offset)")
	AgyLastActiveProjectsCmd.Flags().IntVarP(&lapPageFlag, "page", "p", 1, "Page number (1-indexed, uses --limit)")
	AgyLastActiveProjectsCmd.Flags().IntVar(&lapWordCountFlag, "wordcount", 200, "Maximum word count per prompt sub-item (default 200)")
	AgyLastActiveProjectsCmd.Flags().IntVar(&lapWordCountFlag, "wc", 200, "Maximum word count per prompt sub-item (default 200)")
	AgyLastActiveProjectsCmd.Flags().BoolVar(&lapJSONFlag, "json", false, "Output results in JSON format")
	AgyLastActiveProjectsCmd.Flags().StringVarP(&lapFileFlag, "file", "f", "", "Write JSON output to file path")
	AgyLastActiveProjectsCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		RenderLastActiveProjectsHelp()
	})
}

// RunRunningProjectsPromptsTreeCLI renders `gitmap agy rp prompts ls` as a tree with 200-word prompts.
func RunRunningProjectsPromptsTreeCLI(args []string) error {
	if checkRunningProjectsHelp(args) {
		printRunningProjectsHelp()
		return nil
	}
	projects, err := DiscoverRunningProjects()
	if err != nil {
		return err
	}
	wc := resolveEffectiveWordCount(runningProjectsWordCount)
	trees := convertRunningToTreeRecords(projects, wc)
	return outputTreeRecords(trees, len(trees), len(trees), 0, wc, runningProjectsJSON, runningProjectsFile, "RUNNING PROJECTS PROMPTS TREE")
}

func convertRunningToTreeRecords(projects []RunningProjectRecord, wc int) []LapProjectTreeRecord {
	var out []LapProjectTreeRecord
	for _, p := range projects {
		snippet := CompactWords(p.FullPrompt, wc)
		if snippet == "" {
			snippet = "(active session — awaiting prompt output)"
		}
		conv := LapConversationItem{
			ConversationID: p.ConversationId,
			PromptSeqID:    p.PromptSeqId,
			Prompts:        []string{snippet},
			UpdatedAt:      p.UpdatedAt,
		}
		out = append(out, LapProjectTreeRecord{
			SequenceID: p.SequenceId, ProjectID: p.ProjectId, ProjectAlias: p.ProjectAlias,
			ProjectName: p.ProjectName, ProjectPath: p.ProjectPath, Status: p.Status,
			LastActiveAt: p.UpdatedAt, Conversations: []LapConversationItem{conv},
		})
	}
	return out
}

// RunLastActiveProjectsCLI executes `gitmap agy last-active-projects` (`lap`).
func RunLastActiveProjectsCLI(args []string) error {
	if checkRunningProjectsHelp(args) {
		RenderLastActiveProjectsHelp()
		return nil
	}
	hours := parseLapHoursArg(args)
	wc := resolveEffectiveWordCount(lapWordCountFlag)
	allTrees := discoverLastActiveProjectTrees(hours, wc)
	sliced, offset, limit := paginateLapTrees(allTrees, lapLimitFlag, lapOffsetFlag, lapPageFlag)
	header := fmt.Sprintf("LAST ACTIVE PROJECTS (Last %dh — %d total)", hours, len(allTrees))
	return outputTreeRecords(sliced, len(allTrees), limit, offset, wc, lapJSONFlag, lapFileFlag, header)
}

func parseLapHoursArg(args []string) int {
	for _, a := range args {
		n, err := strconv.Atoi(strings.TrimSpace(a))
		if err == nil && n > 0 {
			return n
		}
	}
	return 24
}

func resolveEffectiveWordCount(wc int) int {
	if wc <= 0 {
		return 200
	}
	return wc
}

func paginateLapTrees(all []LapProjectTreeRecord, limit, offset, page int) ([]LapProjectTreeRecord, int, int) {
	if limit <= 0 {
		limit = 10
	}
	if page > 1 {
		offset = (page - 1) * limit
	}
	if offset < 0 || offset >= len(all) {
		return nil, offset, limit
	}
	end := offset + limit
	if end > len(all) {
		end = len(all)
	}
	return all[offset:end], offset, limit
}

func discoverLastActiveProjectTrees(hours, wc int) []LapProjectTreeRecord {
	cutoff := time.Now().Add(-time.Duration(hours) * time.Hour)
	allPrompts := CollectAllPrompts()
	projects, _ := getAllProjects()
	running, _ := DiscoverRunningProjects()
	rawTrees := assembleActiveTrees(projects, running, allPrompts, cutoff, wc)
	return assignAndPersistLapSequences(rawTrees)
}

func assembleActiveTrees(projects []AgyProject, running []RunningProjectRecord, prompts []AgyPromptEntry, cutoff time.Time, wc int) []LapProjectTreeRecord {
	byPath := make(map[string]*LapProjectTreeRecord)
	seedFromRunningProjects(byPath, running, wc)
	seedFromBrainPrompts(byPath, prompts, cutoff, wc)
	seedFromRecentProjects(byPath, projects, cutoff, wc)
	return sortLapTreeMap(byPath)
}

func seedFromRunningProjects(byPath map[string]*LapProjectTreeRecord, running []RunningProjectRecord, wc int) {
	for _, r := range running {
		key := strings.ToLower(filepath.Clean(r.ProjectPath))
		conv := LapConversationItem{
			ConversationID: r.ConversationId,
			Prompts:        []string{CompactWords(r.FullPrompt, wc)},
			UpdatedAt:      r.UpdatedAt,
		}
		byPath[key] = &LapProjectTreeRecord{
			ProjectID: r.ProjectId, ProjectAlias: r.ProjectAlias, ProjectName: r.ProjectName,
			ProjectPath: r.ProjectPath, Status: r.Status, LastActiveAt: r.UpdatedAt,
			Conversations: []LapConversationItem{conv},
		}
	}
}

func seedFromBrainPrompts(byPath map[string]*LapProjectTreeRecord, prompts []AgyPromptEntry, cutoff time.Time, wc int) {
	for _, p := range prompts {
		if p.CreatedAt.Before(cutoff) || strings.TrimSpace(p.Workspace) == "" {
			continue
		}
		addPromptEntryToTreeMap(byPath, p, wc)
	}
}

func addPromptEntryToTreeMap(byPath map[string]*LapProjectTreeRecord, p AgyPromptEntry, wc int) {
	cleanPath := filepath.Clean(p.Workspace)
	key := strings.ToLower(cleanPath)
	rec, isFound := byPath[key]
	if !isFound {
		rec = newLapRecordForPath(cleanPath, p.CreatedAt.UTC().Format(time.RFC3339))
		byPath[key] = rec
	}
	appendPromptToConversation(rec, p.ConvID, CompactWords(p.Content, wc), p.CreatedAt.UTC().Format(time.RFC3339))
}

func newLapRecordForPath(cleanPath, ts string) *LapProjectTreeRecord {
	name := filepath.Base(cleanPath)
	return &LapProjectTreeRecord{
		ProjectID:    deriveOrLookupProjectID(cleanPath),
		ProjectAlias: strings.ToLower(name),
		ProjectName:  name,
		ProjectPath:  cleanPath,
		Status:       "ACTIVE-24H",
		LastActiveAt: ts,
	}
}

func appendPromptToConversation(rec *LapProjectTreeRecord, convID, text, ts string) {
	if convID == "" {
		convID = "default-conv"
	}
	for i := range rec.Conversations {
		if rec.Conversations[i].ConversationID == convID {
			rec.Conversations[i].Prompts = addPromptBounded(rec.Conversations[i].Prompts, text, 2)
			return
		}
	}
	rec.Conversations = append(rec.Conversations, LapConversationItem{
		ConversationID: convID, Prompts: []string{text}, UpdatedAt: ts,
	})
}

func addPromptBounded(prompts []string, text string, maxItems int) []string {
	if len(prompts) < maxItems {
		return append(prompts, text)
	}
	return prompts
}

func seedFromRecentProjects(byPath map[string]*LapProjectTreeRecord, projects []AgyProject, cutoff time.Time, wc int) {
	for _, p := range projects {
		ws := p.GetPath()
		if ws == "" || !isProjectUpdatedAfter(p.UpdatedAt, cutoff) {
			continue
		}
		ensureRecentProjectInMap(byPath, p, ws, wc)
	}
}

func isProjectUpdatedAfter(updatedAt string, cutoff time.Time) bool {
	ts, isValid := parseTimestampString(updatedAt)
	if !isValid {
		return false
	}
	return ts.After(cutoff)
}

func ensureRecentProjectInMap(byPath map[string]*LapProjectTreeRecord, p AgyProject, ws string, wc int) {
	clean := filepath.Clean(ws)
	key := strings.ToLower(clean)
	if _, hasExisting := byPath[key]; hasExisting {
		return
	}
	convID := resolveWorkspaceConvID(clean)
	conv := LapConversationItem{
		ConversationID: convID,
		Prompts:        []string{CompactWords("Workspace active within lookback window", wc)},
		UpdatedAt:      p.UpdatedAt,
	}
	byPath[key] = &LapProjectTreeRecord{
		ProjectID: p.ID, ProjectAlias: strings.ToLower(p.Name), ProjectName: p.Name,
		ProjectPath: clean, Status: "RECENT", LastActiveAt: p.UpdatedAt,
		Conversations: []LapConversationItem{conv},
	}
}

func sortLapTreeMap(byPath map[string]*LapProjectTreeRecord) []LapProjectTreeRecord {
	var list []LapProjectTreeRecord
	for _, v := range byPath {
		list = append(list, *v)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].LastActiveAt > list[j].LastActiveAt
	})
	return list
}

func assignAndPersistLapSequences(list []LapProjectTreeRecord) []LapProjectTreeRecord {
	var cacheEntries []store.AgySequenceRecord
	now := time.Now().Unix()
	promptCounter := 0
	for i := range list {
		seqNum := i + 1
		list[i].SequenceID = strconv.Itoa(seqNum)
		cacheEntries = appendLapProjectCache(cacheEntries, &list[i], seqNum, &promptCounter, now)
	}
	_, _ = store.EnsureAndGetSequenceCache(cacheEntries, false)
	return list
}

func appendLapProjectCache(entries []store.AgySequenceRecord, rec *LapProjectTreeRecord, seqNum int, pCount *int, now int64) []store.AgySequenceRecord {
	firstConv := resolveFirstConvID(rec.Conversations)
	entries = append(entries, store.AgySequenceRecord{
		SeqId: rec.SequenceID, SeqNum: seqNum, EntryType: "project",
		ProjectId: rec.ProjectID, ProjectAlias: rec.ProjectAlias, ProjectPath: rec.ProjectPath,
		ConversationId: firstConv, CreatedAt: now, ExpiresAt: now + store.SequenceCacheTTLSeconds,
	})
	for cIdx := range rec.Conversations {
		*pCount++
		pSeq := fmt.Sprintf("P%d", *pCount)
		rec.Conversations[cIdx].PromptSeqID = pSeq
		snippet := firstPromptOrEmpty(rec.Conversations[cIdx].Prompts)
		entries = append(entries, store.AgySequenceRecord{
			SeqId: pSeq, SeqNum: *pCount, EntryType: "prompt",
			ProjectId: rec.ProjectID, ProjectAlias: rec.ProjectAlias, ProjectPath: rec.ProjectPath,
			ConversationId: rec.Conversations[cIdx].ConversationID, PromptSnippet: snippet,
			CreatedAt: now, ExpiresAt: now + store.SequenceCacheTTLSeconds,
		})
	}
	return entries
}

func resolveFirstConvID(convs []LapConversationItem) string {
	if len(convs) > 0 {
		return convs[0].ConversationID
	}
	return "default-conv"
}

func firstPromptOrEmpty(prompts []string) string {
	if len(prompts) > 0 {
		return prompts[0]
	}
	return ""
}

func outputTreeRecords(records []LapProjectTreeRecord, totalCount, limit, offset, wc int, isJSON bool, filePath, header string) error {
	if filePath != "" {
		return writeLapTreeToFile(records, filePath)
	}
	if isJSON {
		return renderLapTreeJSON(records)
	}
	renderLapTreeTerminal(records, totalCount, limit, offset, wc, header)
	return nil
}

func writeLapTreeToFile(records []LapProjectTreeRecord, path string) error {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal active projects tree")
	}
	if writeErr := os.WriteFile(path, data, 0644); writeErr != nil {
		return apperror.WrapSimple(writeErr, "write active projects tree file")
	}
	fmt.Printf("  ✔ Saved %d project tree record(s) to %s\n", len(records), path)
	return nil
}

func renderLapTreeJSON(records []LapProjectTreeRecord) error {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal tree JSON")
	}
	fmt.Println(string(data))
	return nil
}

func renderLapTreeTerminal(records []LapProjectTreeRecord, totalCount, limit, offset, wc int, header string) {
	fmt.Printf("\n  %s╔══ %s ══╗%s\n", constants.ColorCyan, header, constants.ColorReset)
	if totalCount > limit {
		fmt.Printf("  %s[info] Showing %d of %d active projects (max limit %d — use --limit or --page to view more)%s\n",
			constants.ColorYellow, len(records), totalCount, limit, constants.ColorReset)
	}
	if len(records) == 0 {
		fmt.Printf("  %sNo matching active projects found in lookback window (offset=%d).%s\n\n", constants.ColorDim, offset, constants.ColorReset)
		return
	}
	for _, r := range records {
		renderSingleLapProjectTree(r, wc)
	}
	fmt.Println()
}

func renderSingleLapProjectTree(r LapProjectTreeRecord, wc int) {
	fmt.Printf("\n  %s[Seq: #%s]%s [ProjectID: %s] [Alias: %s] [Path: %s]\n",
		constants.ColorGreen, r.SequenceID, constants.ColorReset, r.ProjectID, r.ProjectAlias, r.ProjectPath)
	for _, c := range r.Conversations {
		fmt.Printf("    ├─ %s[ProjectID: %s | ConvID: %s | Seq: %s]%s\n",
			constants.ColorCyan, shortProjectId(r.ProjectID), c.ConversationID, c.PromptSeqID, constants.ColorReset)
		renderConversationPromptBranches(c.Prompts, wc)
	}
}

func renderConversationPromptBranches(prompts []string, wc int) {
	for idx, p := range prompts {
		fmt.Printf("    │  └─ Prompt #%d (%d words max): %s\n", idx+1, wc, p)
	}
}

// RenderLastActiveProjectsHelp prints the boxed help menu for `gitmap agy last-active-projects` (`lap`).
func RenderLastActiveProjectsHelp() {
	fmt.Println()
	fmt.Printf("  %s╔══ ANTIGRAVITY LAST-ACTIVE-PROJECTS (LAP) HELP ══════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("  │ Usage:                                                              │")
	fmt.Println("  │   gitmap agy last-active-projects [N] [ls|help] [flags]             │")
	fmt.Println("  │   gitmap agy lap [N] [--limit Y] [--offset Z] [--page P] [--wc T]   │")
	fmt.Println("  │                                                                     │")
	fmt.Println("  │ Parameters & Flags:                                                 │")
	fmt.Println("  │   N                  Lookback hours (default: 24, e.g. 12, 8)       │")
	fmt.Println("  │   -l, --limit Y      Max projects to list (default: 10)             │")
	fmt.Println("  │   --offset, --skip Z Skip Z projects (default: 0)                   │")
	fmt.Println("  │   -p, --page P       Page number (computes offset = (P-1)*limit)    │")
	fmt.Println("  │   --wc, --wordcount  Words per prompt in tree view (default: 200)   │")
	fmt.Println("  │   --json             Output tree in JSON format                     │")
	fmt.Println("  │   -f, --file <path>  Save JSON output to file                       │")
	fmt.Printf("  %s╚═════════════════════════════════════════════════════════════════════╝%s\n\n", constants.ColorCyan, constants.ColorReset)
}
