# Plan 192 (Consolidated): Pipeline Fix Errors AGY (AEF), Duplicate Detection (--force), and Antigravity Queue Dispatch

## Execution Context & Lifecycle
- **Task Start**: Initiated to implement `gitmap pipeline fix errors agy`, `gitmap pipeline-fix errors agy`, `gitmap pipeline agy errors fix (aef)`, and `gitmap pipeline agy-errors-fix (aef)` aliases; deduplicate sent error diagnostics with `--force` / `-f` bypass; automatically capture recent git commit logs alongside pipeline error logs and the embedded canonical RCA fix prompt; dispatch the primary fix prompt to Antigravity IDE (via temp file, clipboard, and CLI) and stage a secondary verification check ("Is it fixed?") into the verification queue.
- **Workflow**: Continuous N-Step Self-Loop across 2-Agent Concurrency and Parallel Micro-Batches.
- **Total Steps/Loops**: 5 subtasks executed in parallel micro-batches.
- **Status**: 100% Completed & Verified.

---

## Consolidated Subtasks Ledger

### Subtask 01: Pipeline AGY Fix Command Routing & Aliases
- **Problem**: Users require flexible command patterns to trigger the Antigravity pipeline error fix flow: `gitmap pipeline fix errors agy`, `gitmap pipeline-fix errors agy`, `gitmap pipeline agy errors fix`, `gitmap pipeline agy-errors-fix`, `gitmap pipeline aef`, top-level `gitmap aef`, `gitmap pipeline-fix`, and `gitmap fix-agy`.
- **Implementation**:
  - `cli/cmdpipeline/pipeline_fix_agy_runner.go` (new): Implemented `IsPipelineFixAgyArgs(args []string) bool` and injectable runner `PipelineAgyFixRunner` detecting `fix errors agy`, `agy errors fix`, `aef`, and compound subcommand patterns while avoiding circular dependencies between packages.
  - `cli/cmdpipeline/pipeline.go` (modified): Updated `runPipeline` entrypoint to route to `PipelineAgyFixRunner` when matching args are encountered.
  - `cli/cmd/clihelpers.go` (modified): Wired `cmdpipeline.PipelineAgyFixRunner = cmdagy.RunPipelineFixAgyCLI` in `init()`.
  - `cli/cmd/rootutility.go` (modified): Registered `pipeline-fix`, `fix-pipeline`, `aef`, `agy-errors-fix`, `fix-agy` in `utilityPipelineEntries`.
  - `cli/cmd/root.go` (modified): Added `pipeline-fix`, `aef`, `agy-errors-fix`, `fix-agy` in root `executeAndAudit` routing.
  - `cli/cmdagy/agy_cmd.go` (modified): Expanded `isFixPipelineAlias` to recognize `aef`, `agy-errors-fix`, `errors-fix`, `fix-errors`, and `fix-agy`.

### Subtask 02: Deduplication State Tracking & `--force` / `-f` Flag
- **Problem**: Re-sending identical pipeline errors creates redundant AI agent tasks. Users must be alerted if the same errors were already dispatched unless explicitly requested with `--force` / `-f`.
- **Implementation**:
  - `cli/cmdagy/agy_fix_pipeline_dedup.go` (new): Implemented `ComputeErrorSignature(repo string, runID uint64, sha, errorLogs string) (string, string)` producing fingerprints (`repo:runId` or `repo:sha:hash`), `LoadSentAgyErrorsStore(path)`, `SaveSentAgyErrorsStore(path, store)`, and `CheckSentErrorDuplicate(sig, store, isForce)`.
  - Stored persistent dispatch history in `.gitmap/pipeline/sent_agy_errors.json` (with fallback to `.lovable/temp/sent_agy_errors.json`).
  - If identical errors were already dispatched without `--force` / `-f`, GitMap prints:
    `⚠ Pipeline errors for run #<id> (commit <sha>) have already been sent to Antigravity!`
    `  Previously sent at: <timestamp> (dispatch count: <N>)`
    `  Do you want to send again? Use --force or -f to send again.`
  - Passing `--force` or `-f` bypasses deduplication, increments dispatch count, and updates timestamp.

### Subtask 03: Prompt Assembly (Fix with RCA & Automated Git Log)
- **Problem**: Fixing pipeline bugs requires context on recent code changes leading up to the failure, in addition to error diagnostics and the canonical 4-part RCA workflow prompt.
- **Implementation**:
  - `cli/cmdagy/agy_fix_pipeline_assembly.go` (new): Implemented `ExtractGitLog(repoDir, 5)` running `git log -n 5 --stat --no-merges` (with fallback to oneline), `FormatRCAHeader`, `FormatGitLogSection`, `FormatPipelineErrorSection`, `LoadCanonicalRcaPrompt`, and `AssembleRcaFixPayload`.
  - Combines top-level `/goal` directive, recent git commit logs, failing pipeline error output, and canonical RCA prompt (`01-prompts/07-bug-fix/01-fix-with-rca.md` or `01-prompts/16-ci-cd/01-ci-cd-fix.md` / `04-ci-cd-fix-with-release.md`) with clean two-line gaps (`\n\n`).

### Subtask 04: Dual Prompt Antigravity Dispatch & Follow-up Verification Queue
- **Problem**: Autonomous debugging requires a two-step feedback loop: first dispatching the fix prompt, followed by a secondary prompt checking whether the issue is resolved ("Is it fixed?").
- **Implementation**:
  - `cli/cmdagy/agy_fix_pipeline_queue.go` (new): Implemented `BuildVerificationFollowupPrompt(repo, runId, sha)` creating the "Is it fixed?" verification check prompt, and `StageVerificationFollowupPrompt` saving to `.lovable/temp/queued-agy-followup-prompt.txt` and `.lovable/temp/agy-prompt-queue.json`.
  - `cli/cmdagy/agy_fix_pipeline.go` (modified): Refactored `RunPipelineFixAgyCLI` to orchestrate primary payload persistence to `.lovable/temp/active-agy-pipeline-fix-prompt.txt`, copy to OS clipboard via `clipboard.WriteAll`, and stage the secondary verification prompt into the queue.
  - Detects local Antigravity CLI binary (`agy.exe` / `agy`) and reports CLI status and resume command (`agy -c`).

### Subtask 05: Help Text Documentation & Hermetic Tests
- **Problem**: Documentation and test coverage needed updating for all new syntax variants, flags, and queuing mechanics.
- **Implementation**:
  - `cli/helptext/agy-fix-pipeline.md` (modified): Documented `gitmap pipeline fix errors agy`, `gitmap pipeline aef`, `gitmap aef`, `--force` / `-f`, and two-prompt architecture.
  - `cli/helptext/pipeline.md` (modified): Added `fix errors agy (aef)` and `--force` flag in subcommands and shortcuts tables.
  - `cli/cmdpipeline/pipeline_fix_agy_test.go` (new): Comprehensive unit tests verifying all command syntax variations and runner dispatch.
  - `cli/cmdagy/agy_pipeline_fix_errors_test.go` (new): Hermetic unit tests verifying flag parsing, deduplication store persistence, signature calculation, prompt assembly, and queue updating.

---

## Verification & Cleanliness Checks
- All functions strictly adhere to <= 15 lines (target <= 8 lines).
- Affirmative booleans only (`is*`, `has*`).
- Zero nested ifs (nesting depth <= 1).
- Universal AppError wrapping across all error paths.
- Targeted linters passed with 0 violations (`check-nested-ifs.py`, `check-boolean-guidelines.py`, `check-error-management.py`).
