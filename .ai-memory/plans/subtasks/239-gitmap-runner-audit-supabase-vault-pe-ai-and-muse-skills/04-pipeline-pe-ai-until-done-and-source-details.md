# Subtask 04: Pipeline PE-AI, Until-Done, 2-Min Poll & Source Details

> **Subtask ID:** Subtask-04  
> **Parent Plan:** `.ai-memory/plans/239-gitmap-runner-audit-supabase-vault-pe-ai-and-muse-skills.md`  
> **Associated Specs:**  
> - `02-spec/21-app/239-gitmap-runner-audit-supabase-vault-pe-ai-and-muse-skills/01-architecture-spec.md`  
> - `02-spec/21-app/239-gitmap-runner-audit-supabase-vault-pe-ai-and-muse-skills/02-component-spec.md`  
> **Owned Files:**  
> - `cli/cmdpipeline/pipeline_flags.go`  
> - `cli/cmdpipeline/pipeline_ai.go`  
> - `cli/cmdpipeline/pipeline_logs.go`  
> - `cli/cmdpipeline/pipeline_dynamic_timeline.go`  
> - `cli/cmdpipeline/pipeline_cache_eval.go`  
> - `cli/cmdpipeline/pipeline_details.go`  
> - `cli/cmdpipeline/helpers.go`  

---

## 1. Objectives

- [x] 1. Add `--ai` flag to `cli/cmdpipeline/pipeline_flags.go`:
  - Add `HasAI bool` to `PipelineErrorFlags`.
  - Parse `--ai` in `parseCommonErrorFlags(args []string, flags *PipelineErrorFlags)`:
    `flags.HasAI = hasArgFlag(args, "--ai") || hasArgFlag(os.Args, "--ai")`.
- [x] 2. Suppress clipboard copying and clipboard notices when `--ai` is enabled:
  - In `cli/cmdpipeline/pipeline_logs.go`:
    - In `renderFailureTerminal(p)` and `renderCleanSuccessTerminal(p)`: Guard `copyReportToClipboard` behind `if !flags.HasAI`.
    - In `outputJSONErrorLogs(content string)` and `writeErrorLogsToDisk(params, content)`: Guard `writeClipboard(content)` behind `if !params.Flags.HasAI`.
  - In `cli/cmdpipeline/pipeline_history.go`: Guard `copyPositionalReportToClipboard` behind `if !hasAIFlag`.
  - In `cli/cmdpipeline/helpers.go`: Provide `shouldWriteClipboard(hasAI bool) bool`.
- [x] 3. Implement 2-minute polling interval in `gitmap pe -t` with early error abort:
  - In `cli/cmdpipeline/pipeline_dynamic_timeline.go`:
    - Configure default timeline polling interval to 2 minutes (`120s`).
    - At each polling iteration, inspect active run logs and stack traces. If errors or failures are detected in active log streams, immediately surface failure summary via `gitmap pe` and terminate execution without waiting for the timeout duration to elapse.
- [x] 4. Add `-ud` / `--until-done` flag:
  - In `cli/cmdpipeline/pipeline_flags.go`: Add `HasUntilDone bool` to `PipelineErrorFlags`, parsed from `-ud`, `--until-done`.
  - In `cli/cmdpipeline/pipeline_dynamic_timeline.go`: When `HasUntilDone` is true, run continuous polling loop until the full CI/CD run completes (all jobs succeed or any job fails), then display final `gitmap pe` status and diagnostics.
- [x] 5. Implement immediate Git hash retrieval from SQLite DB:
  - In `cli/cmdpipeline/pipeline_cache_eval.go`:
    - If the target commit SHA (or current HEAD commit SHA) is already stored in the local SQLite split-db and all workflow runs for that commit have completed status, pull immediately from the database without waiting or issuing network requests.
    - Exempt completed SHA hits from being bypassed when `-t` is provided.
- [x] 6. Enhance failure diagnostics with log source details:
  - In `cli/cmdpipeline/pipeline.go` and `cli/cmdpipeline/pipeline_logs.go`:
    - Enrich `SectionFailure` struct with `LoggerName string`, `WorkflowName string`, `RunUrl string`, `JobName string`, `StepName string`, and relative `SavedLogFile string`.
    - Render a dedicated `Source Details:` block in terminal failure output:
      ```text
        Source Details:
          Logger:    <LoggerName>
          Workflow:  <WorkflowName>
          Job/Step:  <JobName> / <StepName>
          URL:       <RunUrl>
          Log File:  <SavedLogFile>
      ```
- [x] 7. Verify compilation and linter compliance:
  - `go test -v ./cli/cmdpipeline/...`
  - `python linter-scripts/check-nested-ifs.py`
  - `python linter-scripts/check-enum-and-boolean.py`
  - `python linter-scripts/check-relative-paths.py`

---

## 2. Types & Data Structures

```go
package cmdpipeline

// PipelineErrorFlags additions
type PipelineErrorFlags struct {
	// ... existing fields ...
	HasAI        bool // Parsed from --ai
	HasUntilDone bool // Parsed from -ud, --until-done
}

// SectionFailure with source metadata
type SectionFailure struct {
	LoggerName     string   `json:"loggerName,omitempty"`
	WorkflowName   string   `json:"workflowName"`
	RunId          uint64   `json:"runId"`
	RunUrl         string   `json:"runUrl,omitempty"`
	JobName        string   `json:"jobName"`
	StepName       string   `json:"stepName"`
	FailureSummary string   `json:"failureSummary"`
	ErrorLines     []string `json:"errorLines"`
	Warnings       []string `json:"warnings,omitempty"`
	SavedLogFile   string   `json:"savedLogFile,omitempty"` // Relative path
	CreatedAt      string   `json:"createdAt,omitempty"`
	StackTrace     string   `json:"stackTrace,omitempty"`
}
```

---

## 3. Acceptance Criteria

- **AC-PE-AI-001 (Clipboard Suppression):** Running `gitmap pe --ai` emits complete terminal failure diagnostics and log extracts to stdout/stderr, but executes zero clipboard write operations (`writeClipboard`) and suppresses `📋 Copied...` notice messages.
- **AC-PE-2MIN-002 (2-Minute Early Error Abort):** When running `gitmap pe -t`, the watcher polls every 2 minutes. If errors/stack traces appear in the workflow logs during an active run, `gitmap` terminates the polling loop immediately, displays the error diagnostics, and exits with code 1 without waiting for timeout expiration.
- **AC-PE-UD-003 (Until-Done Polling):** Running `gitmap pe -ud` (or `gitmap pe --until-done`) watches continuously until all workflow jobs complete, then presents the final status report.
- **AC-PE-SHA-004 (Instant Hash Retrieval):** If the requested commit SHA is already stored in the local SQLite database and all runs for that commit are finished, `gitmap pe` retrieves the data instantly without initiating live GitHub API polling.
- **AC-PE-SRC-005 (Log Source Details):** Every failure reported by `gitmap pe` displays comprehensive log source details including logger/runner name, workflow name, run URL, job name, step name, and relative log path.

---

## 4. Verification Instructions

1. Run unit tests for pipeline flags, cache evaluation, and log rendering:
   ```bash
   go test -v ./cli/cmdpipeline/...
   ```
2. Test clipboard suppression:
   ```bash
   gitmap pe --ai --force
   ```
   Verify system clipboard contents remain unmodified and no clipboard notice appears.
3. Verify linter clean status:
   ```bash
   python linter-scripts/check-nested-ifs.py
   python linter-scripts/check-enum-and-boolean.py
   python linter-scripts/check-relative-paths.py
   ```
