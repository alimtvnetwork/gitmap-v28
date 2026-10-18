# Specification: Heatmap RC JSON Response & AI Expectation Tracking

## 1. Overview
The Heatmap Telemetry Subsystem enables autonomous AI agents (Antigravity, Muse, Cloud Code, Codex, Grok, Cursor) to log cognitive decisions and file modifications directly to isolated SQLite split-databases (`repo-cache/.heatmap/<prefix>-<ai>/<repo>/heatmap.db` and local repository `data/heatmap.db`).

This specification defines the structured JSON response envelope and bidirectional AI expectation lifecycle (`ExpectedOutcome`, `ActualOutcome`, `IsExpectationMet`, and derived `Status`), allowing AI agents to record intended outcomes before execution and retrospectively confirm whether execution succeeded or failed.

---

## 2. JSON Response Schema (`RcLogResponseEnvelope`)

When `--json` or `-j` (or `-F json` / `--format=json`) is passed to `gitmap rc heatmap log` or `gitmap rc heatmap update` (and root `gitmap heatmap log` / `gitmap heatmap update`), the CLI outputs a single contiguous JSON envelope:

```json
{
  "status": "success",
  "action_id": 2,
  "where_saved": {
    "db_path": "D:/work/repo-cache/.heatmap/01-antigravity/gitmap-v28/heatmap.db",
    "relative_db_path": "../repo-cache/.heatmap/01-antigravity/gitmap-v28/heatmap.db",
    "repo_secrets_url": "https://github.com/alimtvnetwork/repo-secrets.git"
  },
  "db_info": {
    "driver": "sqlite",
    "table": "HeatmapActionLog",
    "ai_name": "antigravity",
    "repo_slug": "gitmap-v28",
    "path": "D:/work/repo-cache/.heatmap/01-antigravity/gitmap-v28/heatmap.db"
  },
  "what_saved": {
    "action_id": 2,
    "parent_log_id": null,
    "task_slug": "task-01",
    "repo_name": "gitmap-v28",
    "action_type": "EDIT",
    "file_path": "cli/heatmap/heatmap_rc_cmd.go",
    "start_line": 0,
    "end_line": 0,
    "code_snippet": "",
    "command_text": "",
    "reasoning_state": "EXPLORATION",
    "reasoning_notes": "Added structured JSON envelope",
    "execution_depth": 0,
    "execution_order": 1
  },
  "ai_expectation": {
    "expected_outcome": "Enable AI to log expectation and retrieve structured JSON envelope",
    "actual_outcome": "JSON response validated successfully",
    "is_expectation_met": true,
    "status": "MET"
  }
}
```

### 2.1 Section Breakdown
1. **`where_saved`**:
   - `db_path`: Canonical physical filesystem path to the target SQLite database.
   - `relative_db_path`: Relative filesystem path from the current working directory to the target database file.
   - `repo_secrets_url`: Upstream remote origin URL for the companion `repo-secrets` vault (resolved via `SpecialReposSplitDB` or local git config).
2. **`db_info`**:
   - `driver`: Database engine driver (`sqlite`).
   - `table`: Primary telemetry table (`HeatmapActionLog`).
   - `ai_name`: Active AI engine slug (`antigravity`, `muse`, etc.).
   - `repo_slug`: Sanitized target repository name.
   - `path`: Database file path.
3. **`what_saved`**:
   - Reflects the exact database row fields inserted into `HeatmapActionLog`.
4. **`ai_expectation`**:
   - `expected_outcome`: String description of what the AI intended to achieve beforehand.
   - `actual_outcome`: String description of the observed outcome after action execution.
   - `is_expectation_met`: Boolean flag indicating if expectation was met (`true` or `false`).
   - `status`: Lifecycle indicator (`PENDING`, `MET`, or `FAILED`).

---

## 3. Bidirectional AI Expectation Lifecycle

```
[Phase 1: Pre-Execution Log]
AI runs: gitmap rc heatmap log --expected "Fix nil pointer in config parser" --json
  Result: Action #42 recorded.
  Expectation Status: PENDING (is_expectation_met = true, actual_outcome = "")

[Phase 2: Execution]
AI performs edit and executes unit tests: go test ./...

[Phase 3: Post-Execution Confirmation]
Case A (Success):
  AI runs: gitmap rc heatmap update --id 42 --actual "Tests passed in 0.2s" --met true --json
  Result: Status transitions to MET.

Case B (Failure):
  AI runs: gitmap rc heatmap update --id 42 --actual "Test failed: index out of range" --failed --json
  Result: Status transitions to FAILED (is_expectation_met = false).
```

---

## 4. CLI Subcommand & Flag Mapping

### 4.1 Logging (`log`)
- Flags:
  - `--expected`, `-e`, `--expectation`: Description of planned outcome.
  - `--actual`, `-o`, `--outcome`: Description of actual outcome (if immediately known).
  - `--met`: Boolean (`true` / `false`).
  - `--json`, `-j`, `-F json`: Request structured JSON response.

### 4.2 Updating (`update` / `outcome`)
- Command forms:
  - `gitmap rc heatmap update --id <ID> [flags]`
  - `gitmap rc heatmap outcome --id <ID> [flags]`
  - `gitmap heatmap update --id <ID> [flags]`
- Supported resolution for `--id`:
  - `--id <N>`, `--action-id <N>`, or positional first argument `<N>`.
- Supported satisfaction flags:
  - Success: `--met true`, `--success`, `--status success`, `--status met`, `--status passed`.
  - Failure: `--met false`, `--failed`, `--fail`, `--status fail`, `--status failed`.
