package cmdai

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var (
	trainCmd = &cobra.Command{
		Use:     "train",
		Aliases: []string{"llm-train", "reconstruct"},
		Short:   "Extract AI decision history and rationale for LLM context or training",
		RunE:    runTrainCmd,
	}

	trainFlags LlmTrainOptions
)

func runTrainCmd(cmd *cobra.Command, args []string) error {
	db, err := OpenAiAnalysisSplitDB(trainFlags.RepoRoot)
	if err != nil {
		return err
	}
	defer db.Close()

	return RunAiTrain(db, trainFlags)
}

// RunAiTrain extracts decision chains and renders them according to format options.
func RunAiTrain(db *sql.DB, opts LlmTrainOptions) error {
	chains, extractErr := ExtractDecisionHistory(db, opts.TaskId, opts.Limit)
	if extractErr != nil {
		return extractErr
	}

	data, formatErr := formatDecisionChains(chains, opts.Format)
	if formatErr != nil {
		return formatErr
	}

	return writeOrPrintOutput(data, opts.Output)
}

func formatDecisionChains(chains []TaskDecisionChain, format LlmTrainingFormat) ([]byte, error) {
	cleanFormat := strings.ToLower(string(format))
	isJsonl := cleanFormat == string(LlmFormatJsonl)
	if isJsonl {
		return FormatAsJsonl(chains)
	}

	renderedMd := FormatAsMarkdown(chains)

	return []byte(renderedMd), nil
}

func writeOrPrintOutput(data []byte, outputPath string) error {
	hasOutput := strings.TrimSpace(outputPath) != ""
	if !hasOutput {
		fmt.Println(string(data))

		return nil
	}

	dir := filepath.Dir(outputPath)
	mkdirErr := os.MkdirAll(dir, 0755)
	if mkdirErr != nil {
		return apperror.WrapSimple(mkdirErr, "cmdai.writeOrPrintOutput.mkdir")
	}

	writeErr := os.WriteFile(outputPath, data, 0644)
	if writeErr != nil {
		return apperror.WrapSimple(writeErr, "cmdai.writeOrPrintOutput.writeFile")
	}

	fmt.Printf("Wrote AI decision history to %q (%d bytes)\n", outputPath, len(data))

	return nil
}

// ExtractDecisionHistory retrieves chronological decision sequences from the database.
func ExtractDecisionHistory(db *sql.DB, taskID string, limit int) ([]TaskDecisionChain, error) {
	hasTaskID := strings.TrimSpace(taskID) != ""
	if hasTaskID {
		return extractSingleTaskChain(db, strings.TrimSpace(taskID))
	}

	return extractMultipleTaskChains(db, limit)
}

func extractSingleTaskChain(db *sql.DB, taskID string) ([]TaskDecisionChain, error) {
	task, err := GetAiAnalysisTask(db, taskID)
	if err != nil {
		return nil, err
	}

	lines, linesErr := GetAiAnalysisLinesForTask(db, taskID)
	if linesErr != nil {
		return nil, linesErr
	}

	return []TaskDecisionChain{{Task: *task, Lines: lines}}, nil
}

func extractMultipleTaskChains(db *sql.DB, limit int) ([]TaskDecisionChain, error) {
	hasInvalidLimit := limit <= 0
	if hasInvalidLimit {
		limit = 10
	}

	tasks, err := ListAiAnalysisTasks(db, limit, 0)
	if err != nil {
		return nil, err
	}

	chains := make([]TaskDecisionChain, 0, len(tasks))
	for _, t := range tasks {
		lines, linesErr := GetAiAnalysisLinesForTask(db, t.TaskId)
		if linesErr != nil {
			return nil, linesErr
		}
		chains = append(chains, TaskDecisionChain{Task: t, Lines: lines})
	}

	return chains, nil
}

// FormatAsMarkdown transforms decision chains into human-readable markdown.
func FormatAsMarkdown(chains []TaskDecisionChain) string {
	var sb strings.Builder
	sb.WriteString("# AI Decision History & Rationale Training Log\n\n")

	hasChains := len(chains) > 0
	if !hasChains {
		sb.WriteString("*No recorded AI decision sessions available.*\n")

		return sb.String()
	}

	for _, chain := range chains {
		renderMarkdownTask(&sb, chain.Task)
		renderMarkdownLines(&sb, chain.Lines)
	}

	return sb.String()
}

func renderMarkdownTask(sb *strings.Builder, task AiAnalysisTask) {
	sb.WriteString(fmt.Sprintf("## Task: %s (`%s`)\n\n", task.TaskId, task.Status))
	sb.WriteString(fmt.Sprintf("- **Objective:** %s\n", task.TaskDescription))
	sb.WriteString(fmt.Sprintf("- **Model:** `%s`\n", task.ModelName))
	sb.WriteString(fmt.Sprintf("- **Started:** %s\n", task.StartedAt.Format("2006-01-02 15:04:05 UTC")))
	sb.WriteString(fmt.Sprintf("- **Files Impacted:** Read: %d | Modified: %d | Deleted: %d\n",
		task.TotalFilesRead, task.TotalFilesModified, task.TotalFilesDeleted))
	hasSummary := strings.TrimSpace(task.ReasoningSummary) != ""
	if hasSummary {
		sb.WriteString(fmt.Sprintf("- **Reasoning Summary:** %s\n", task.ReasoningSummary))
	}
	sb.WriteString("\n### Chronological File Decisions & Rationale:\n\n")
}

func renderMarkdownLines(sb *strings.Builder, lines []AiAnalysisLine) {
	hasLines := len(lines) > 0
	if !hasLines {
		sb.WriteString("*(No granular file line decisions recorded for this task)*\n\n")

		return
	}

	for i, l := range lines {
		lineSpan := fmt.Sprintf("%d-%d", l.StartLine, l.EndLine)
		isSingleLine := l.StartLine == l.EndLine
		if isSingleLine {
			lineSpan = fmt.Sprintf("%d", l.StartLine)
		}

		sb.WriteString(fmt.Sprintf("%d. **[%s]** `%s` (lines %s):\n", i+1, strings.ToUpper(string(l.OperationType)), l.FilePath, lineSpan))
		sb.WriteString(fmt.Sprintf("   - **Rationale:** %s\n", l.Reasoning))
		hasSnippet := strings.TrimSpace(l.ContentSnippet) != ""
		if hasSnippet {
			sb.WriteString(fmt.Sprintf("   - **Snippet:** `%s`\n", l.ContentSnippet))
		}
	}
	sb.WriteString("\n---\n\n")
}

// FormatAsJsonl produces JSON Lines suitable for LLM fine-tuning or few-shot injection.
func FormatAsJsonl(chains []TaskDecisionChain) ([]byte, error) {
	var buf bytes.Buffer
	for _, chain := range chains {
		rec := buildConversationRecord(chain)
		lineBytes, err := json.Marshal(rec)
		if err != nil {
			return nil, apperror.WrapSimple(err, "cmdai.FormatAsJsonl.marshal")
		}
		buf.Write(lineBytes)
		buf.WriteByte('\n')
	}

	return buf.Bytes(), nil
}

func buildConversationRecord(chain TaskDecisionChain) LlmConversationRecord {
	userPrompt := fmt.Sprintf("What decisions were made during task %q (%s)?", chain.Task.TaskId, chain.Task.TaskDescription)
	assistantReply := buildAssistantReply(chain)

	return LlmConversationRecord{
		Messages: []LlmMessageItem{
			{Role: "system", Content: "You are an autonomous AI coding assistant providing chronological decision-making rationale."},
			{Role: "user", Content: userPrompt},
			{Role: "assistant", Content: assistantReply},
		},
	}
}

func buildAssistantReply(chain TaskDecisionChain) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("During task %q, the objective was: %s\n", chain.Task.TaskId, chain.Task.TaskDescription))
	hasSummary := strings.TrimSpace(chain.Task.ReasoningSummary) != ""
	if hasSummary {
		sb.WriteString(fmt.Sprintf("Outcome Summary: %s\n", chain.Task.ReasoningSummary))
	}
	sb.WriteString("Decision sequence:\n")

	for i, l := range chain.Lines {
		sb.WriteString(fmt.Sprintf("%d. [%s] %s (lines %d-%d): %s\n",
			i+1, l.OperationType, l.FilePath, l.StartLine, l.EndLine, l.Reasoning))
	}

	return sb.String()
}

func initTrainCommands() {
	trainCmd.Flags().StringVarP(&trainFlags.TaskId, "task", "t", "", "Filter reasoning by specific task ID")
	trainCmd.Flags().StringVarP((*string)(&trainFlags.Format), "format", "f", "markdown", "Output format: markdown, jsonl")
	trainCmd.Flags().IntVarP(&trainFlags.Limit, "limit", "l", 10, "Maximum number of tasks to inspect")
	trainCmd.Flags().StringVarP(&trainFlags.Output, "output", "o", "", "Write formatted output to file")

	AiCmd.AddCommand(trainCmd)
}
