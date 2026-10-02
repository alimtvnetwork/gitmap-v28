package cmdpipeline

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

// HistoryAIOptions holds configuration for gitmap pe history-ai.
type HistoryAIOptions struct {
	Limit        int
	Repo         string
	IsJSON       bool
	IsDetailed   bool
	FilePath     string
	TempFileName string
}

// HistoryAIPayload represents the structured error dossier for AI ingestion.
type HistoryAIPayload struct {
	Repo         string            `json:"repo"`
	GeneratedAt  string            `json:"generatedAt"`
	CommitCount  int               `json:"commitCount"`
	TotalErrors  int               `json:"totalErrors"`
	MarkdownFile string            `json:"markdownFile"`
	JSONFile     string            `json:"jsonFile"`
	Commits      []HistoryAICommit `json:"commits"`
}

// HistoryAICommit captures failure telemetry for a single commit.
type HistoryAICommit struct {
	Sha             string                `json:"sha"`
	ShortSha        string                `json:"shortSha"`
	Branch          string                `json:"branch"`
	WorkflowName    string                `json:"workflowName"`
	RunId           uint64                `json:"runId"`
	RunUrl          string                `json:"runUrl"`
	CreatedAt       string                `json:"createdAt"`
	DurationSeconds int                   `json:"durationSeconds"`
	FailedSteps     []HistoryAIFailedStep `json:"failedSteps"`
	StackTrace      string                `json:"stackTrace,omitempty"`
}

// HistoryAIFailedStep records discrete step failure detail.
type HistoryAIFailedStep struct {
	JobName        string   `json:"jobName"`
	StepName       string   `json:"stepName"`
	FailureSummary string   `json:"failureSummary"`
	ErrorLines     []string `json:"errorLines"`
}

// IsHistoryAISubcmd checks whether an argument matches the history-ai command name.
func IsHistoryAISubcmd(arg string) bool {
	lower := strings.ToLower(strings.TrimSpace(arg))
	return lower == "history-ai" || lower == "hai" || lower == "history_ai"
}

// HasHistoryAISubcmd reports whether args invoke the history-ai command.
func HasHistoryAISubcmd(args []string) bool {
	for _, a := range args {
		if IsHistoryAISubcmd(a) {
			return true
		}
	}
	return false
}

// HandlePipelineHistoryAI handles the gitmap pe history-ai [N] command invocation.
func HandlePipelineHistoryAI(args []string) (bool, error) {
	if !HasHistoryAISubcmd(args) {
		return false, nil
	}
	if hasArgFlag(args, "--help") || hasArgFlag(args, "-h") {
		printHistoryAIHelp()
		return true, nil
	}
	return true, executePipelineHistoryAI(args)
}

func executePipelineHistoryAI(args []string) error {
	opts := parseHistoryAIOptions(args)
	db, err := pipelinedb.OpenPipelineSplitDb(opts.Repo)
	if err != nil {
		return err
	}
	defer db.Close()

	runs := fetchHistoricalFailedRuns(db, opts.Repo, opts.Limit)
	commits := groupRunsByFailingCommits(db, opts.Repo, runs, opts.Limit)
	payload := buildHistoryAIPayload(opts.Repo, commits)
	mdPath, jsonPath := resolveHistoryAIOutputPaths(opts)

	if err := writeHistoryAIFiles(&payload, mdPath, jsonPath); err != nil {
		return err
	}
	return outputHistoryAIResults(payload, opts)
}

func parseHistoryAIOptions(args []string) HistoryAIOptions {
	flags := ParsePipelineErrorFlags(args)
	repo := flags.RepoTarget
	if len(repo) == 0 {
		repo = resolveCurrentRepoSlug()
	}
	return HistoryAIOptions{
		Limit:        extractHistoryAILimit(args),
		Repo:         repo,
		IsJSON:       flags.IsJSON,
		IsDetailed:   flags.IsDetailed,
		FilePath:     flags.FilePath,
		TempFileName: flags.TempFileName,
	}
}

func extractNextHistoryAILimit(args []string, i int, a string) (int, bool) {
	if !IsHistoryAISubcmd(a) || i+1 >= len(args) {
		return 0, false
	}
	return tryExtractLimitArg(args[i+1])
}

func extractHistoryAILimit(args []string) int {
	for i, a := range args {
		if val, isOk := extractNextHistoryAILimit(args, i, a); isOk {
			return val
		}
		if val, isOk := tryExtractLimitArg(a); isOk {
			return val
		}
	}
	return 5
}

func tryExtractLimitArg(arg string) (int, bool) {
	val, err := strconv.Atoi(arg)
	if err == nil && val > 0 && val <= 50 {
		return val, true
	}
	return 0, false
}

func fetchHistoricalFailedRuns(db *pipelinedb.PipelineSplitDb, repo string, limit int) []pipelinedb.PipelineRunRecord {
	ensureHistoryRunsPopulated(db, repo)
	queryLimit := limit * 4
	if queryLimit < 20 {
		queryLimit = 20
	}
	res := db.QueryLastNFailedRuns(queryLimit)
	if res.IsSuccess() {
		return res.Data
	}
	return nil
}

func ensureHistoryRunsPopulated(db *pipelinedb.PipelineSplitDb, repo string) {
	res := db.QueryLastNFailedRuns(5)
	if res.IsSuccess() && len(res.Data) > 0 {
		return
	}
	runs := queryWorkflowRuns(repo)
	RecordFetchedRunsToSplitDb(repo, runs)
}

func groupRunsByFailingCommits(db *pipelinedb.PipelineSplitDb, repo string, runs []pipelinedb.PipelineRunRecord, limit int) []HistoryAICommit {
	selected := selectUniqueFailingRuns(runs, limit)
	if len(selected) == 0 {
		return nil
	}
	return parallelBuildHistoryCommits(db, repo, selected)
}

func selectUniqueFailingRuns(runs []pipelinedb.PipelineRunRecord, limit int) []pipelinedb.PipelineRunRecord {
	var selected []pipelinedb.PipelineRunRecord
	seen := make(map[string]bool)
	for _, r := range runs {
		if seen[r.Sha] {
			continue
		}
		seen[r.Sha] = true
		selected = append(selected, r)
		if len(selected) >= limit {
			break
		}
	}
	return selected
}

func parallelBuildHistoryCommits(db *pipelinedb.PipelineSplitDb, repo string, runs []pipelinedb.PipelineRunRecord) []HistoryAICommit {
	commits := make([]HistoryAICommit, len(runs))
	var wg sync.WaitGroup
	for i, r := range runs {
		wg.Add(1)
		go func(idx int, record pipelinedb.PipelineRunRecord) {
			defer wg.Done()
			commits[idx] = buildHistoryAICommit(db, repo, record)
		}(i, r)
	}
	wg.Wait()
	return commits
}

func buildHistoryAICommit(db *pipelinedb.PipelineSplitDb, repo string, r pipelinedb.PipelineRunRecord) HistoryAICommit {
	shortSha := r.Sha
	if len(shortSha) > 7 {
		shortSha = shortSha[:7]
	}
	steps := collectFailedStepsForRun(db, repo, r)
	return HistoryAICommit{
		Sha:             r.Sha,
		ShortSha:        shortSha,
		Branch:          r.Branch,
		WorkflowName:    r.WorkflowName,
		RunId:           r.RunId,
		RunUrl:          r.RunUrl,
		CreatedAt:       r.CreatedAt,
		DurationSeconds: r.DurationSeconds,
		FailedSteps:     steps,
		StackTrace:      resolveRunStackTrace(repo, r.RunId),
	}
}

func collectFailedStepsForRun(db *pipelinedb.PipelineSplitDb, repo string, r pipelinedb.PipelineRunRecord) []HistoryAIFailedStep {
	if steps := parseErrorsFromDb(db, r.RunId); len(steps) > 0 {
		return steps
	}
	if steps := parseErrorsFromLogFallback(repo, r.RunId); len(steps) > 0 {
		return steps
	}
	return parseErrorsFromRunOrFetch(db, repo, r)
}

func parseErrorsFromDb(db *pipelinedb.PipelineSplitDb, runId uint64) []HistoryAIFailedStep {
	if db == nil {
		return nil
	}
	if steps := parseErrorsFromDetailedDb(db, runId); len(steps) > 0 {
		return steps
	}
	return parseErrorsFromErrorLogDb(db, runId)
}

func parseErrorsFromRunOrFetch(db *pipelinedb.PipelineSplitDb, repo string, r pipelinedb.PipelineRunRecord) []HistoryAIFailedStep {
	rawLogs := queryFailedRunLogs(repo, r.RunId)
	if len(rawLogs) == 0 {
		return nil
	}
	ghJobs := queryRunJobs(repo, r.RunId)
	jobs := CorrelateFailedJobsWithRunJobs(rawLogs, ghJobs)
	if db != nil && len(jobs) > 0 {
		ghRun := ghRunItem{
			DatabaseId: r.RunId, Name: r.WorkflowName,
			HeadSha: r.Sha, HeadBranch: r.Branch,
		}
		saveParsedFailedJobs(db, repo, ghRun, jobs, rawLogs)
	}
	return convertFailedJobsToHistorySteps(jobs, rawLogs)
}

func convertFailedJobsToHistorySteps(jobs []FailedJobItem, rawLogs string) []HistoryAIFailedStep {
	var steps []HistoryAIFailedStep
	for _, j := range jobs {
		steps = append(steps, HistoryAIFailedStep{
			JobName:        j.JobName,
			StepName:       j.StepName,
			FailureSummary: j.FailureSummary,
			ErrorLines:     j.ErrorLines,
		})
	}
	if len(steps) > 0 {
		return steps
	}

	lines := cleanAndFilterErrorLines(rawLogs)
	if len(lines) > 0 {
		steps = append(steps, HistoryAIFailedStep{
			JobName:        "Workflow Execution",
			StepName:       "Failed Step",
			FailureSummary: "Failure diagnostic extracted from run logs",
			ErrorLines:     lines,
		})
	}

	return steps
}

func parseErrorsFromDetailedDb(db *pipelinedb.PipelineSplitDb, runId uint64) []HistoryAIFailedStep {
	res := db.QueryDetailedErrorLogsByRunId(runId)
	if res.IsFailure() || len(res.Data) == 0 {
		return nil
	}
	var steps []HistoryAIFailedStep
	for _, rec := range res.Data {
		lines := cleanAndFilterErrorLines(rec.RawLogs)
		steps = append(steps, HistoryAIFailedStep{
			JobName:        rec.WorkflowName,
			StepName:       rec.StepName,
			FailureSummary: rec.ErrorText,
			ErrorLines:     lines,
		})
	}
	return steps
}

func parseErrorsFromErrorLogDb(db *pipelinedb.PipelineSplitDb, runId uint64) []HistoryAIFailedStep {
	res := db.QueryErrorLogsByRunId(runId)
	if res.IsFailure() || len(res.Data) == 0 {
		return nil
	}
	var steps []HistoryAIFailedStep
	for _, rec := range res.Data {
		lines := cleanAndFilterErrorLines(rec.RawLogs)
		steps = append(steps, HistoryAIFailedStep{
			JobName:        rec.WorkflowName,
			StepName:       rec.StepName,
			FailureSummary: rec.ErrorText,
			ErrorLines:     lines,
		})
	}
	return steps
}

func parseErrorsFromLogFallback(repo string, runId uint64) []HistoryAIFailedStep {
	rawLog, ok := readCachedPipelineLogForRepo(repo, runId)
	if !ok || len(rawLog) == 0 {
		return nil
	}
	lines := cleanAndFilterErrorLines(rawLog)
	if len(lines) == 0 {
		return nil
	}
	return []HistoryAIFailedStep{
		{
			JobName:        "Workflow Run",
			StepName:       "Failed Step",
			FailureSummary: "Step failure extracted from run logs",
			ErrorLines:     lines,
		},
	}
}

func resolveRunStackTrace(repo string, runId uint64) string {
	rawLog, ok := readCachedPipelineLogForRepo(repo, runId)
	if !ok || len(rawLog) == 0 {
		return ""
	}
	return extractStackTraceFromLog(rawLog)
}

func cleanAndFilterErrorLines(raw string) []string {
	if len(raw) == 0 {
		return nil
	}
	lines := strings.Split(raw, "\n")
	var filtered []string
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if len(trimmed) == 0 || !isActionableErrorLine(trimmed) {
			continue
		}
		filtered = append(filtered, trimmed)
		if len(filtered) >= 10 {
			break
		}
	}

	return filtered
}

func isActionableErrorLine(line string) bool {
	lower := strings.ToLower(line)
	return strings.Contains(lower, "error") || strings.Contains(lower, "fail") ||
		strings.Contains(lower, "fatal") || strings.Contains(lower, "panic") ||
		strings.Contains(lower, "exit status") || strings.Contains(lower, "violation")
}

func buildHistoryAIPayload(repo string, commits []HistoryAICommit) HistoryAIPayload {
	totalErrors := 0
	for _, c := range commits {
		totalErrors += len(c.FailedSteps)
	}
	return HistoryAIPayload{
		Repo:        repo,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		CommitCount: len(commits),
		TotalErrors: totalErrors,
		Commits:     commits,
	}
}

func resolveHistoryAIOutputPaths(opts HistoryAIOptions) (string, string) {
	if len(opts.FilePath) > 0 {
		return opts.FilePath, opts.FilePath + ".json"
	}
	if len(opts.TempFileName) > 0 {
		dir := resolveTempDir()
		base := strings.TrimSuffix(opts.TempFileName, filepath.Ext(opts.TempFileName))
		return filepath.Join(dir, base+".md"), filepath.Join(dir, base+".json")
	}
	baseDir := filepath.Join(".", ".ai-memory", "pipeline-ai")
	return filepath.Join(baseDir, "history-errors.md"), filepath.Join(baseDir, "history-errors.json")
}

func writeHistoryAIFiles(p *HistoryAIPayload, mdPath, jsonPath string) error {
	p.MarkdownFile = filepath.ToSlash(mdPath)
	p.JSONFile = filepath.ToSlash(jsonPath)

	mdContent := generateHistoryAIMarkdown(*p)
	if err := writeContentToFile(mdPath, mdContent); err != nil {
		return err
	}
	jsonData, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return writeContentToFile(jsonPath, string(jsonData))
}

func generateHistoryAIMarkdown(p HistoryAIPayload) string {
	var sb strings.Builder
	appendMarkdownHeader(&sb, p)
	appendMarkdownTaxonomy(&sb, p)
	for i, c := range p.Commits {
		appendCommitMarkdownSection(&sb, i+1, c)
	}
	appendMarkdownAIPreventionDirectives(&sb)
	return sb.String()
}

func appendMarkdownHeader(sb *strings.Builder, p HistoryAIPayload) {
	sb.WriteString("# Pipeline Historical Errors & AI Anti-Mistake Training Dossier\n\n")
	sb.WriteString(fmt.Sprintf("- **Repository:** `%s`\n", p.Repo))
	sb.WriteString(fmt.Sprintf("- **Generated At:** %s\n", p.GeneratedAt))
	sb.WriteString(fmt.Sprintf("- **Failing Commits Analyzed:** %d\n", p.CommitCount))
	sb.WriteString(fmt.Sprintf("- **Total Failure Diagnostics:** %d\n\n", p.TotalErrors))
	sb.WriteString("> **Purpose:** This dossier exposes all historical pipeline failures across recent commits\n")
	sb.WriteString("> to train AI coding agents so they NEVER repeat these architectural, linting, or testing mistakes.\n\n")
}

func appendMarkdownTaxonomy(sb *strings.Builder, p HistoryAIPayload) {
	sb.WriteString("## Executive Failure Taxonomy\n\n")
	sb.WriteString("| Commit | Workflow | Failing Step(s) | Primary Issue Category |\n")
	sb.WriteString("|--------|----------|-----------------|------------------------|\n")
	for _, c := range p.Commits {
		stepName := "General Workflow"
		if len(c.FailedSteps) > 0 {
			stepName = c.FailedSteps[0].StepName
		}
		category := categorizeFailure(stepName)
		sb.WriteString(fmt.Sprintf("| `%s` | %s | %s | %s |\n", c.ShortSha, c.WorkflowName, stepName, category))
	}
	sb.WriteString("\n")
}

func categorizeFailure(stepName string) string {
	lower := strings.ToLower(stepName)
	if strings.Contains(lower, "lint") || strings.Contains(lower, "gocritic") {
		return "Code Quality / Linter"
	}
	if strings.Contains(lower, "policy") || strings.Contains(lower, "guard") {
		return "Coding Guideline Violation"
	}
	if strings.Contains(lower, "build") {
		return "Cross-Platform Build"
	}
	if strings.Contains(lower, "test") {
		return "Unit / Integration Test"
	}
	return "Pipeline Execution"
}

func appendCommitMarkdownSection(sb *strings.Builder, idx int, c HistoryAICommit) {
	sb.WriteString(fmt.Sprintf("### [%d] Commit `%s` (Branch: `%s`)\n\n", idx, c.ShortSha, c.Branch))
	sb.WriteString(fmt.Sprintf("- **Workflow:** %s\n", c.WorkflowName))
	sb.WriteString(fmt.Sprintf("- **Run ID:** [%d](%s)\n", c.RunId, c.RunUrl))
	sb.WriteString(fmt.Sprintf("- **Recorded At:** %s (Duration: %ds)\n\n", c.CreatedAt, c.DurationSeconds))

	for _, step := range c.FailedSteps {
		appendStepMarkdownDetails(sb, step)
	}
	if len(c.StackTrace) > 0 {
		sb.WriteString("#### Stack Trace\n```text\n")
		sb.WriteString(c.StackTrace)
		sb.WriteString("\n```\n\n")
	}
}

func appendStepMarkdownDetails(sb *strings.Builder, s HistoryAIFailedStep) {
	sb.WriteString(fmt.Sprintf("#### Job: `%s` | Step: `%s`\n", s.JobName, s.StepName))
	sb.WriteString(fmt.Sprintf("- **Summary:** %s\n", s.FailureSummary))
	if len(s.ErrorLines) > 0 {
		sb.WriteString("```text\n")
		for _, l := range s.ErrorLines {
			sb.WriteString(l)
			sb.WriteString("\n")
		}
		sb.WriteString("```\n")
	}
	sb.WriteString("\n")
}

func appendMarkdownAIPreventionDirectives(sb *strings.Builder) {
	sb.WriteString("## AI Agent Directives to Prevent Mistakes\n\n")
	sb.WriteString("1. **Pre-Commit Verification:** Run `python 03-ai-scripts/05-guideline-autofixer.py <dir>` before staging code.\n")
	sb.WriteString("2. **Strict Sizing Rules:** Keep Go functions $\\le 15$ lines and files $\\le 100$ lines (Rule R14).\n")
	sb.WriteString("3. **Positive Boolean Hygiene:** Use affirmative identifiers (`isSuccess`, `hasErrors`, `isReady`); reject double negatives.\n")
	sb.WriteString("4. **Universal Error Handling:** Always wrap and return errors using `*appfault.AppError` / `apperror.WrapSimple`.\n")
	sb.WriteString("5. **Cross-Platform Safety:** Ensure all path operations use `filepath.ToSlash` and tests do NOT rely on unmocked OS commands.\n\n")
}

func outputHistoryAIResults(p HistoryAIPayload, opts HistoryAIOptions) error {
	if !opts.IsJSON {
		renderHistoryAITerminal(p, opts.IsDetailed)
		return nil
	}

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}

func renderHistoryAITerminal(p HistoryAIPayload, isDetailed bool) {
	fmt.Printf("\n  %s● Pipeline AI Historical Error Dossier (Last %d Commits) [%s]%s\n\n",
		constants.ColorCyan, p.CommitCount, p.Repo, constants.ColorReset)
	printHistoryAITable(p.Commits)
	fmt.Printf("  %s✓%s AI Error Dossier: %s\n", constants.ColorGreen, constants.ColorReset, p.MarkdownFile)
	fmt.Printf("  %s✓%s Training JSON:    %s\n", constants.ColorGreen, constants.ColorReset, p.JSONFile)
	fmt.Printf("\n  %sTip:%s Feed this file into AI context to avoid repeating past pipeline failures.\n\n",
		constants.ColorYellow, constants.ColorReset)
}

func printHistoryAITable(commits []HistoryAICommit) {
	fmt.Printf("    %-8s  %-10s  %-24s  %s\n", "Commit", "Branch", "Workflow", "Failing Step")
	fmt.Printf("    %-8s  %-10s  %-24s  %s\n", "------", "------", "--------", "------------")
	for i, c := range commits {
		printHistoryAITableRow(i+1, c)
	}
	fmt.Println()
}

func printHistoryAITableRow(idx int, c HistoryAICommit) {
	stepName := "General Failure"
	if len(c.FailedSteps) > 0 {
		stepName = c.FailedSteps[0].StepName
	}
	if len(stepName) > 40 {
		stepName = stepName[:37] + "..."
	}
	fmt.Printf("    %-8s  %-10s  %-24s  %s\n", c.ShortSha, c.Branch, c.WorkflowName, stepName)
}

func printHistoryAIHelp() {
	fmt.Println(constants.ColorCyan + "Usage: gitmap pe history-ai [N] [flags]" + constants.ColorReset)
	fmt.Println("Extracts historical CI/CD pipeline errors across recent commits to train AI agents.")
	fmt.Println("\nArguments:")
	fmt.Println("  [N]                 Number of historical failing commits to extract (default: 5)")
	fmt.Println("\nFlags:")
	fmt.Println("  --repo <slug>       Target repository slug (default: current repository)")
	fmt.Println("  --file <path>       Output Markdown file path (default: .ai-memory/pipeline-ai/history-errors.md)")
	fmt.Println("  --tempfile <name>   Write dossier to .ai-memory/temp/<name>")
	fmt.Println("  --json              Output structured JSON to stdout")
	fmt.Println("  -v, --detailed      Show expanded failure step logs in output")
	fmt.Println("  -h, --help          Show this history-ai help menu")
}
