# Subtask 03: AI Analysis LLM Train, Export & Cleanup

> **Subtask ID:** Subtask-03  
> **Parent Plan:** `.ai-memory/plans/240-gitmap-mcp-ai-analysis-muse-installer-and-agm-fleet.md`  
> **Target Subsystems:** `cli/cmdai/`  
> **Owned Files:**  
> - `cli/cmdai/ai_train_types.go`  
> - `cli/cmdai/ai_train.go`  
> - `cli/cmdai/ai_export.go`  
> - `cli/cmdai/ai_clear.go`  
> - `cli/cmdai/ai_train_test.go`  

---

## 1. Concrete Objectives

1. **LLM Decision History Summarizer (`gitmap ai train` / `gitmap ai llm-train`):**
   - Query Split-DB SQLite (`ai_analysis.db`) to retrieve tasks, files inspected, line mutations, and rationale recorded by AI agents.
   - Synthesize decision sequences directly addressing: *"What was the previous decision-making? Why it took the decision-making, and what was the reason?"*
   - Implement formatting engines:
     - **Markdown Decision Chains (`--format markdown`):** Generates structured markdown summaries with chronological rationale per file and function.
     - **JSONL Dataset Pairs (`--format jsonl`):** Generates fine-tuning prompt/completion pairs or conversational messages (`system`, `user`, `assistant`) for model context injection.
   - Support streaming to stdout or writing directly to `--output <file>`.

2. **Portable Export Engine (`gitmap ai export`):**
   - Provide multi-format export to facilitate sharing analysis data across development nodes and CI workers:
     - `--format json`: Serializes full task trees and associated line telemetry into a single JSON schema.
     - `--format zip`: Bundles JSON metadata alongside referenced file content snippets for offline inspection.
     - `--format db`: Creates an atomic, read-consistent SQLite snapshot of `.gitmap/data/ai-analysis/<slug>/sql.db`.

3. **Analysis Ingestion Engine (`gitmap ai import`):**
   - Ingests exported JSON payloads or database snapshots into the current workspace's Split-DB.
   - Resolves or generates new task IDs on primary key collision using an `upsert` or `rename-on-conflict` strategy.

4. **Storage Pruning & Retention Governance (`gitmap ai clear` / `gitmap ai prune`):**
   - Prevent unbounded database growth across prolonged agent development sprints:
     - `--task <id>`: Deletes specific task and cascades to all child lines.
     - `--older-than <duration>`: Prunes records older than specified duration (e.g. `14d`, `30d`).
     - `--all --force`: Purges all task history and executes `VACUUM` on `sql.db` to reclaim disk space.

5. **CLI Command Binding:**
   - Register `trainCmd`, `exportCmd`, `importCmd`, and `clearCmd` under `AiCmd` (`gitmap ai`).
   - Provide command alias `gitmap ai llm-train` pointing to `gitmap ai train`.

---

## 2. Core Domain Types & Structs (`cli/cmdai/ai_train_types.go`)

```go
package cmdai

import "time"

// LlmTrainingFormat defines output format for LLM reasoning ingestion.
type LlmTrainingFormat string

const (
	LlmFormatMarkdown LlmTrainingFormat = "markdown"
	LlmFormatJsonl    LlmTrainingFormat = "jsonl"
)

// LlmMessageItem represents a single role turn in conversational training schemas.
type LlmMessageItem struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// LlmConversationRecord represents a JSONL fine-tuning record.
type LlmConversationRecord struct {
	Messages []LlmMessageItem `json:"messages"`
}

// LlmTrainOptions configures the reasoning extraction and formatting.
type LlmTrainOptions struct {
	TaskId    string
	Format    LlmTrainingFormat
	Limit     int
	Output    string
	RepoRoot  string
}

// ExportFormatType defines supported packaging formats for export.
type ExportFormatType string

const (
	ExportFormatJson ExportFormatType = "json"
	ExportFormatZip  ExportFormatType = "zip"
	ExportFormatDb   ExportFormatType = "db"
)

// AiExportBundle represents a serialized collection of tasks and lines.
type AiExportBundle struct {
	ExportedAt time.Time        `json:"exportedAt"`
	RepoSlug   string           `json:"repoSlug"`
	Tasks      []AiAnalysisTask `json:"tasks"`
	Lines      []AiAnalysisLine `json:"lines"`
}

// AiExportOptions configures export operations.
type AiExportOptions struct {
	TaskId   string
	Format   ExportFormatType
	Output   string
	RepoRoot string
}

// AiImportOptions configures import operations.
type AiImportOptions struct {
	InputPath string
	RepoRoot  string
}

// AiClearOptions configures pruning and truncation.
type AiClearOptions struct {
	TaskId    string
	OlderThan time.Duration
	All       bool
	Force     bool
	RepoRoot  string
}
```

---

## 3. Step-by-Step Implementation Plan

### Step 1: Decision Sequence Extractor

- In `cli/cmdai/ai_train.go`, implement `ExtractDecisionHistory(db *sql.DB, taskID string, limit int) ([]TaskDecisionChain, error)`:
  - Queries `AiAnalysisTask` and joins `AiAnalysisLine` ordered by `created_at ASC`.
  - Groups lines by file path and operation type.
  - Generates cohesive reasoning timelines per task.

### Step 2: Markdown & JSONL Formatters

- In `cli/cmdai/ai_train.go`, implement:
  - `FormatAsMarkdown(chains []TaskDecisionChain) string`:
    - Renders objective, files observed, line ranges, and decisions.
  - `FormatAsJsonl(chains []TaskDecisionChain) ([]byte, error)`:
    - Synthesizes user prompt asking why changes were made and assistant response detailing rationale.

### Step 3: Export Engine Implementation

- In `cli/cmdai/ai_export.go`, implement `ExecuteAiExport(opts AiExportOptions) error`:
  - `json`: Marshals `AiExportBundle` using indented JSON encoder.
  - `zip`: Creates archive containing `manifest.json`, `records.json`, and file snippets.
  - `db`: Executes SQLite backup API or file copy with WAL flush (`PRAGMA wal_checkpoint(TRUNCATE)`).

### Step 4: Import Engine Implementation

- In `cli/cmdai/ai_export.go`, implement `ExecuteAiImport(opts AiImportOptions) error`:
  - Detects input format (JSON or DB).
  - Validates schema compatibility.
  - Performs batch inserts into `AiAnalysisTask` and `AiAnalysisLine` inside an atomic transaction.

### Step 5: Pruning & Vacuum Engine

- In `cli/cmdai/ai_clear.go`, implement `ExecuteAiClear(opts AiClearOptions) (int64, error)`:
  - If `--task <id>` is provided: executes `DELETE FROM AiAnalysisTask WHERE task_id = ?`.
  - If `--older-than <dur>` is provided: calculates cutoff time and deletes matching tasks.
  - If `--all` is provided: executes `DELETE FROM AiAnalysisLine; DELETE FROM AiAnalysisTask;`.
  - If rows were deleted, executes `VACUUM;` to reclaim unused disk space.

### Step 6: CLI Command Wiring

- Create commands:
  - `gitmap ai train` (alias: `gitmap ai llm-train`)
  - `gitmap ai export`
  - `gitmap ai import`
  - `gitmap ai clear` (alias: `gitmap ai prune`)
- Bind flags `--task`, `--format`, `--output`, `--input`, `--older-than`, `--all`, `--force`.
- Wire into `AiCmd` in `cli/cmdai/ai_cmd.go`.

### Step 7: Automated Testing

- Create `cli/cmdai/ai_train_test.go`:
  - Seed test tasks and reasoning lines.
  - Verify Markdown and JSONL formatting outputs.
  - Test JSON export and subsequent import into a fresh database.
  - Test pruning with age threshold and table truncation with vacuum.

---

## 4. Acceptance Criteria

- [x] `gitmap ai train` formats historical reasoning into coherent Markdown summaries.
- [x] `gitmap ai train --format jsonl` produces valid JSONL records parseable by standard LLM loaders.
- [x] `gitmap ai export --format json --output out.json` writes a valid portable dataset.
- [x] `gitmap ai export --format zip --output out.zip` packages tasks and code snippets.
- [x] `gitmap ai import --input out.json` successfully restores tasks and lines into Split-DB.
- [x] `gitmap ai clear --task <id>` cleanly deletes the task and all linked reasoning lines.
- [x] `gitmap ai prune --older-than 30d` accurately removes expired records and runs `VACUUM`.
- [x] Zero absolute paths or `file:///` URIs exist in any owned file.
- [x] All unit tests in `cli/cmdai/ai_train_test.go` pass with 100% success rate.

---

## 5. Verification Commands

```bash
  # 1. Run unit tests for training, export, and cleanup
  go test -v ./cli/cmdai -run "TestAi(Train|Export|Clear).*"

  # 2. Generate Markdown LLM reasoning summary
  gitmap ai train --limit 5

  # 3. Generate JSONL dataset for fine-tuning
  gitmap ai train --format jsonl --output dataset.jsonl --limit 10

  # 4. Export analysis database to JSON
  gitmap ai export --format json --output ai_analysis_backup.json

  # 5. Export analysis database to ZIP
  gitmap ai export --format zip --output ai_analysis_bundle.zip

  # 6. Test pruning records older than 14 days
  gitmap ai prune --older-than 14d

  # 7. Test clearing specific task
  gitmap ai clear --task "test-subtask-01"

  # 8. Run repository relative paths linter
  python linter-scripts/check-relative-paths.py
```
