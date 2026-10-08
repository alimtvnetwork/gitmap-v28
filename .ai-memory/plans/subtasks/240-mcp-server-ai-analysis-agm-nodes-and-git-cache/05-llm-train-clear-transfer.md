# Subtask 05: LLM Training Dataset Generation, Database Pruning & Cross-System Transfer

- **Parent Task:** `240-mcp-server-ai-analysis-agm-nodes-and-git-cache`
- **Subtask ID:** `05-llm-train-clear-transfer`
- **Spec Reference:** `02-spec/21-app/240-mcp-server-ai-analysis-agm-nodes-and-git-cache/01-architecture-spec.md` (Sections 2.3, 2.5, and 6)
- **Status:** `PENDING`
- **Assigned Subagent:** Implementation Subagent (Phase 2)

---

## 1. Objective & Scope

Implement the **LLM Fine-Tuning Dataset Generator, Split-DB Pruning Engine, and Cross-System Transfer Subsystem** for the GitMap AI Analysis platform.

### Core Capabilities
1. **LLM Training Dataset Extraction (`gitmap ai-analysis llm-train`):**
   - Harvests recorded AST diffs, line substitutions, architectural reasoning, and rule violations from `ai-analysis.db`.
   - Generates instruction-tuning datasets in standard JSONL (Instruction/Input/Output), Alpaca, or ShareGPT formats to train and fine-tune specialized LLM coding models directly on high-discipline developer decisions.
2. **Database Pruning & Vault Scrubber (`gitmap ai-analysis clear`):**
   - Prunes completed or aged tasks from `ai-analysis.db` based on retention windows (`--days <N>`).
   - Automatically unlinks and removes associated temporary file vaults in `<temp>/gitmap/removed/<taskId>/` to prevent disk bloat.
   - Provides safe `--dry-run` inspection before committing permanent deletions.
3. **Cross-System Transport Engine (`gitmap ai-analysis export` / `import`):**
   - Supports portable exports across three formats: JSON (normalized metadata with strict relative paths), ZIP (complete package containing DB and staged files), and SQLite (`VACUUM INTO` portable transfer DB).
   - Allows importing external task sets into local repositories while preserving UUID uniqueness.

---

## 2. Concrete Files to Create / Modify

| File | Nature | Purpose |
| :--- | :--- | :--- |
| `cli/cmdai/ai_analysis_llmtrain.go` | Create | Implements `gitmap ai-analysis llm-train` CLI command and dataset transformation pipelines. |
| `cli/cmdai/ai_analysis_clear.go` | Create | Implements `gitmap ai-analysis clear` CLI command with dry-run estimation and staging directory purge. |
| `cli/cmdai/ai_analysis_transfer.go` | Create | Implements `gitmap ai-analysis export` and `gitmap ai-analysis import` supporting JSON, ZIP, and SQLite formats. |
| `cli/cmdai/ai_analysis_llmtrain_test.go` | Create | Unit and integration tests validating dataset formatting, pruning filters, and round-trip export/import. |

---

## 3. Detailed Component Architecture

### 3.1 LLM Training Formatter (`cli/cmdai/ai_analysis_llmtrain.go`)
- **JSONL Record Format:**
  ```json
  {
    "id": "gitmap-task-001-line-12",
    "system": "You are an autonomous senior Go software engineer enforcing strict enterprise coding guidelines.",
    "instruction": "Refactor the following function to enforce guard clauses and eliminate nested conditions.",
    "input": "<original_code_lines>",
    "reasoning": "<architectural_reasoning_from_db>",
    "output": "<proposed_code_lines>",
    "rule_violation": "nested-if-and-guard-clauses",
    "file_path": "02-spec/21-app/..."
  }
  ```
- **Filter Flags:**
  - `--out <path>`: Output destination (default: `ai-analysis-dataset.jsonl`).
  - `--format <jsonl|alpaca|sharegpt>`: Target schema format.
  - `--min-lines <n>`: Exclude trivial whitespace or one-line changes.
  - `--category <cat>`: Filter by task category (e.g., `refactor`, `audit`).

### 3.2 Task Pruner (`cli/cmdai/ai_analysis_clear.go`)
- **Syntax:** `gitmap ai-analysis clear [--days <N>] [--force] [--dry-run]`
- **Workflow:**
  1. Queries `AiTask` for tasks where `CompletedAt > 0` and `CompletedAt < now - days*86400`.
  2. If `--dry-run`, computes count of tasks, files, lines, and disk space of staged vaults, then renders summary table and exits.
  3. If executing:
     - Iterates through tasks, deleting staging directories in `<temp>/gitmap/removed/<taskId>/`.
     - Deletes tasks from `AiTask` (triggering cascading delete across `AiTaskFile` and `AiTaskLine`).
     - Emits success summary envelope.

### 3.3 Transfer Subsystem (`cli/cmdai/ai_analysis_transfer.go`)
- **Export Syntax:** `gitmap ai-analysis export [--out <path>] [--format <json|zip|sqlite>]`
- **Import Syntax:** `gitmap ai-analysis import <file> [--dry-run] [--force]`
- Enforces strict relative paths: any paths in exported artifacts must be relative to the repository root.

---

## 4. Acceptance Criteria

- [ ] `gitmap ai-analysis llm-train` outputs valid JSONL containing all required instruction-tuning attributes.
- [ ] Exported datasets strictly omit absolute host paths.
- [ ] `gitmap ai-analysis clear` removes completed tasks beyond the day threshold and safely scrubs orphaned temp directories in `<temp>/gitmap/removed/<taskId>/`.
- [ ] `--dry-run` flag in `clear` provides accurate impact forecasting without mutating database or disk.
- [ ] Round-trip export and import verified across JSON, ZIP, and SQLite formats without primary key collisions.
- [ ] All unit and contract tests in `cli/cmdai/ai_analysis_llmtrain_test.go` pass 100%.

---

## 5. Verification Commands

```powershell
# 1. Run unit tests for LLM train export and pruning
go test -v ./cli/cmdai -run "TestAiAnalysisLlmTrain|TestAiAnalysisClear|TestAiAnalysisTransfer"

# 2. Test dry-run dataset export manually
go run . ai-analysis llm-train --out test-dataset.jsonl --format jsonl
```
