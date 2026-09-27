# 167: Verification Gates & Acceptance Criteria

## 1. Acceptance Criteria

- **AC-01: Split-DB Schema & Lifecycle**:
  - `PromptBackupBatch`, `PromptBackupItem`, and `PromptRestoreLedger` tables are created in `data/backup-prompts/sql.db` (or custom `-file/-f`).
  - Active and enqueued prompts are captured with project name, path, ID, conversation ID, sequence ID, prompt text, word count, and status.
- **AC-02: Restore & Auto-Pruning**:
  - Restoring marks entries in `PromptRestoreLedger`.
  - When `--keep/-k` is passed, `IsKept` is set to 1 and entries are excluded from retention pruning.
  - Non-kept restored entries older than 24 hours (1 day) are pruned on subsequent backup/restore invocations, outputting the exact single-line notice `old data has been removed`.
  - Manual clean command (`running-prompts clean [--force]`) prunes restored entries or forces complete database purge.
- **AC-03: Running Prompts CLI & Truncation**:
  - `running-prompts ls` applies `--wordcount/--wc` default 100 words (or custom N), `--limit/-l` default 8 (or custom Y), and prints full prompt only when `--full` is supplied.
  - `--json` formats output as structured JSON.
  - `export` writes `.db` or `.json` (full prompts by default, or with `--wc`).
  - `import` ingests `.db` or `.json` and enqueues prompts.
- **AC-04: Running Projects & SSH Aggregation**:
  - `running-projects` identifies projects with active/enqueued prompts.
  - `--ssh` queries cluster nodes and produces JSON or table summaries.
- **AC-05: `fpug` and `sug` Automation**:
  - `fpug` / `finish-prompts-until-green` polls target projects at `-t` intervals (min 30s) until all pipelines are green.
  - `sug` / `shutdown-until-green` supports `ls`, `help`, `run`, `add-projects`, `rm`, `agy-running-projects`, `-t` flag (default 5m), and executes native OS shutdown commands on Windows, Ubuntu/Linux, and macOS upon green status.
- **AC-06: `gitmap pe` DB Cache Short-Circuit**:
  - Verifies local commit SHA in database before downloading logs; serves cached report immediately if commit matches.
  - Fix dispatch is isolated and only triggered when explicit fix subcommands are provided.

## 2. Quality & Linter Invariants
- Zero nested `if` statements (max depth 1 inside functions).
- Positive boolean identifiers only (`is*`, `has*`).
- Structured Go error returns (`*apperror.AppError`).
- Function sizing <= 8-15 lines.
