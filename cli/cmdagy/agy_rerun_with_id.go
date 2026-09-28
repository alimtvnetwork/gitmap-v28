// Package cmdagy — agy_rerun_with_id.go implements `gitmap agy rerun-with-id` (`rwi`) and `gitmap agy rerun-with-convid` (`rwc`).
package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RerunWithIDResult models the structured response of `rwi` and `rwc`.
type RerunWithIDResult struct {
	SequenceID     string `json:"sequenceId"`
	ProjectID      string `json:"projectId"`
	ProjectAlias   string `json:"projectAlias"`
	ProjectPath    string `json:"projectPath"`
	ConversationID string `json:"conversationId"`
	PromptSource   string `json:"promptSource"`
	PromptPreview  string `json:"promptPreview"`
	IsDryRun       bool   `json:"isDryRun"`
	QueueFile      string `json:"queueFile,omitempty"`
}

var (
	rwiPromptFlag string
	rwiJSONFlag   bool
	rwiDryRunFlag bool
)

// AgyRerunWithIDCmd injects or replays a prompt into a target project and conversation.
var AgyRerunWithIDCmd = &cobra.Command{
	Use:     "rerun-with-id <project-id|alias|sequence|path> <convid> [prompt-or-file]",
	Aliases: []string{"rwi"},
	Short:   "Inject/rerun prompt by project ID/alias/sequence (#1) and conversation ID (or P1)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunRerunWithIDCLI(args)
	},
}

// AgyRerunWithConvIDCmd injects or replays a prompt using only a short sequence ID (P1) or conversation ID.
var AgyRerunWithConvIDCmd = &cobra.Command{
	Use:     "rerun-with-convid <convid|short-seq-id> [prompt-or-file]",
	Aliases: []string{"rwc", "rerun-with-prompt-id", "rwp"},
	Short:   "Inject/rerun prompt by conversation ID or 24h sequence ID (e.g. P1, #1)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunRerunWithConvIDCLI(args)
	},
}

func init() {
	AgyRerunWithIDCmd.Flags().StringVarP(&rwiPromptFlag, "prompt", "p", "", "Named prompt template from 01-prompts/ (e.g. read, fix-pipeline)")
	AgyRerunWithIDCmd.Flags().BoolVar(&rwiJSONFlag, "json", false, "Output result in JSON format")
	AgyRerunWithIDCmd.Flags().BoolVarP(&rwiDryRunFlag, "dry-run", "d", false, "Preview target and prompt without injecting")

	AgyRerunWithConvIDCmd.Flags().StringVarP(&rwiPromptFlag, "prompt", "p", "", "Named prompt template from 01-prompts/")
	AgyRerunWithConvIDCmd.Flags().BoolVar(&rwiJSONFlag, "json", false, "Output result in JSON format")
	AgyRerunWithConvIDCmd.Flags().BoolVarP(&rwiDryRunFlag, "dry-run", "d", false, "Preview target and prompt without injecting")
}

// RunRerunWithIDCLI executes `gitmap agy rerun-with-id` (`rwi`).
func RunRerunWithIDCLI(args []string) error {
	if len(args) < 2 {
		return apperror.NewSimple("usage: gitmap agy rerun-with-id <project-id|alias|seq|path> <convid> [prompt-or-file] [-p <name>]", "E9042")
	}
	projRec, err := resolveRwiProjectTarget(args[0])
	if err != nil {
		return err
	}
	convID, seqID := resolveRwiConvTarget(args[1], projRec)
	rawPromptArg := strings.Join(args[2:], " ")
	promptText, source := resolveRwiPromptContent(rawPromptArg, rwiPromptFlag, projRec.ProjectPath, projRec.PromptSnippet)
	return dispatchAndRenderRwi(projRec, convID, seqID, promptText, source)
}

// RunRerunWithConvIDCLI executes `gitmap agy rerun-with-convid` (`rwc` / `rwp`).
func RunRerunWithConvIDCLI(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("usage: gitmap agy rerun-with-convid <convid|short-seq-id> [prompt-or-file] [-p <name>]", "E9043")
	}
	projRec, convID, seqID, err := resolveRwcSingleToken(args[0])
	if err != nil {
		return err
	}
	rawPromptArg := strings.Join(args[1:], " ")
	promptText, source := resolveRwiPromptContent(rawPromptArg, rwiPromptFlag, projRec.ProjectPath, projRec.PromptSnippet)
	return dispatchAndRenderRwi(projRec, convID, seqID, promptText, source)
}

func resolveRwiProjectTarget(token string) (*store.AgySequenceRecord, error) {
	if cached, err := store.ResolveSequenceEntry(token); err == nil && cached != nil {
		return cached, nil
	}
	if checkDirExists(token) {
		return buildRecordFromFolder(token)
	}
	projects, _ := getAllProjects()
	matched, err := findProjectByFlexibleTarget(projects, token)
	if err != nil {
		return nil, apperror.WrapSimple(err, "resolve project target")
	}
	return &store.AgySequenceRecord{
		SeqId:        "1",
		ProjectId:    matched.ID,
		ProjectAlias: strings.ToLower(matched.Name),
		ProjectPath:  matched.GetPath(),
	}, nil
}

func buildRecordFromFolder(rawPath string) (*store.AgySequenceRecord, error) {
	abs, err := filepath.Abs(rawPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "resolve folder path")
	}
	clean := findEnclosingGitRoot(abs)
	name := filepath.Base(clean)
	return &store.AgySequenceRecord{
		SeqId:          "1",
		ProjectId:      deriveOrLookupProjectID(clean),
		ProjectAlias:   strings.ToLower(name),
		ProjectPath:    clean,
		ConversationId: resolveWorkspaceConvID(clean),
	}, nil
}

func resolveRwiConvTarget(convToken string, projRec *store.AgySequenceRecord) (string, string) {
	if cached, err := store.ResolveSequenceEntry(convToken); err == nil && cached != nil && cached.ConversationId != "" {
		return cached.ConversationId, cached.SeqId
	}
	clean := strings.TrimPrefix(strings.TrimSpace(convToken), "#")
	if clean == "" {
		return projRec.ConversationId, projRec.SeqId
	}
	return clean, projRec.SeqId
}

func resolveRwcSingleToken(token string) (*store.AgySequenceRecord, string, string, error) {
	if cached, err := store.ResolveSequenceEntry(token); err == nil && cached != nil {
		convID := selectOrFallbackConvID(cached.ConversationId, cached.ProjectPath)
		return cached, convID, cached.SeqId, nil
	}
	ws := resolveConvWorkspaceFromBrain(token)
	if ws != "" {
		rec, err := buildRecordFromFolder(ws)
		return rec, token, "P1", err
	}
	rec, err := resolveRwiProjectTarget(token)
	if err != nil {
		return nil, "", "", err
	}
	return rec, resolveWorkspaceConvID(rec.ProjectPath), rec.SeqId, nil
}

func selectOrFallbackConvID(convID, projectPath string) string {
	if convID != "" {
		return convID
	}
	return resolveWorkspaceConvID(projectPath)
}

func resolveConvWorkspaceFromBrain(convPrefix string) string {
	clean := strings.ToLower(strings.TrimSpace(convPrefix))
	for _, p := range CollectAllPrompts() {
		if strings.HasPrefix(strings.ToLower(p.ConvID), clean) && p.Workspace != "" {
			return p.Workspace
		}
	}
	return ""
}

func resolveRwiPromptContent(rawArg, namedTemplate, projectPath, fallbackSnippet string) (string, string) {
	if strings.TrimSpace(namedTemplate) != "" {
		return loadNamedPromptTemplate(namedTemplate, projectPath), "template:" + namedTemplate
	}
	trimmed := strings.TrimSpace(rawArg)
	if fileContent, isFile := tryReadPromptFileCandidate(trimmed, projectPath); isFile {
		return fileContent, "file:" + trimmed
	}
	if trimmed != "" {
		return trimmed, "literal"
	}
	if strings.TrimSpace(fallbackSnippet) != "" {
		return fallbackSnippet, "cache-snippet"
	}
	return defaultReadMemoryPrompt, "default-read"
}

func loadNamedPromptTemplate(name, projectPath string) string {
	clean := strings.TrimSuffix(strings.TrimSpace(name), ".md")
	home, _ := os.UserHomeDir()
	candidates := []string{
		filepath.Join(projectPath, "01-prompts", clean+".md"),
		filepath.Join(projectPath, ".ai-memory", "prompts", clean+".md"),
		filepath.Join("01-prompts", clean+".md"),
		filepath.Join(home, ".gitmap", "prompts", clean+".md"),
	}
	for _, c := range candidates {
		data, err := os.ReadFile(c)
		if err == nil && len(strings.TrimSpace(string(data))) > 0 {
			return strings.TrimSpace(string(data))
		}
	}
	return fmt.Sprintf("Execute canonical prompt workflow: %s", clean)
}

func tryReadPromptFileCandidate(candidate, projectPath string) (string, bool) {
	if candidate == "" {
		return "", false
	}
	paths := []string{candidate, filepath.Join(projectPath, candidate)}
	for _, p := range paths {
		if content, isRead := readValidPromptFile(p); isRead {
			return content, true
		}
	}
	return "", false
}

func readValidPromptFile(path string) (string, bool) {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return "", false
	}
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		return "", false
	}
	return strings.TrimSpace(string(data)), true
}

func dispatchAndRenderRwi(proj *store.AgySequenceRecord, convID, seqID, promptText, source string) error {
	if convID == "" {
		convID = resolveWorkspaceConvID(proj.ProjectPath)
	}
	res := &RerunWithIDResult{
		SequenceID:     seqID,
		ProjectID:      proj.ProjectId,
		ProjectAlias:   proj.ProjectAlias,
		ProjectPath:    proj.ProjectPath,
		ConversationID: convID,
		PromptSource:   source,
		PromptPreview:  CompactWords(promptText, 40),
		IsDryRun:       rwiDryRunFlag,
	}
	if !rwiDryRunFlag {
		title := fmt.Sprintf("Rerun [%s | %s]", proj.ProjectAlias, shortProjectId(convID))
		qFile, _ := writeReadPromptToProjectQueue(proj.ProjectPath, title, promptText)
		res.QueueFile = qFile
		_ = InjectAgyPrompt(proj.ProjectPath, promptText, title, "rerun_with_id", false)
	}
	return renderRwiResult(res)
}

func renderRwiResult(res *RerunWithIDResult) error {
	if rwiJSONFlag {
		return printRwiJSON(res)
	}
	printRwiConfirmationBox(res)
	return nil
}

func printRwiJSON(res *RerunWithIDResult) error {
	data, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal rwi result")
	}
	fmt.Println(string(data))
	return nil
}

func printRwiConfirmationBox(res *RerunWithIDResult) {
	fmt.Printf("\n  %s╔══ ANTIGRAVITY RERUN-WITH-ID DISPATCH ═══════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  │ Target Bracket: [%s | %s | Seq: %s]\n", shortProjectId(res.ProjectID), shortProjectId(res.ConversationID), res.SequenceID)
	fmt.Printf("  │ Project Alias:  %s (%s)\n", res.ProjectAlias, res.ProjectID)
	fmt.Printf("  │ Project Path:   %s\n", res.ProjectPath)
	fmt.Printf("  │ Conversation:   %s\n", res.ConversationID)
	fmt.Printf("  │ Prompt Source:  %s\n", res.PromptSource)
	fmt.Printf("  │ Prompt Preview: %s\n", res.PromptPreview)
	fmt.Printf("  %s╚═════════════════════════════════════════════════════════════════════╝%s\n\n", constants.ColorCyan, constants.ColorReset)
}
