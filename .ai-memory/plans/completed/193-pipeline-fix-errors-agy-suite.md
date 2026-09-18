# Plan 193 (Consolidated): Pipeline Fix Errors AGY (AEF) Full Alias Suite, Direct Fix Routing, and Antigravity Queue Verification

## Execution Context & Lifecycle
- **Task Start**: Initiated to verify that `gitmap pipeline fix errors agy`, `gitmap pipeline-fix errors agy`, `gitmap pipeline agy errors fix (aef)`, `gitmap pipeline agy-errors-fix (aef)` are 100% complete and working; add direct routing for `gitmap fix agy`, `gitmap fix-agy`, `gitmap fix errors agy`, `gitmap fix agy errors`, and `gitmap agy errors fix` aliases; verify duplicate detection with `--force` / `-f` bypass; verify automated extraction of git commit history alongside pipeline error logs and canonical 4-part RCA prompt; and confirm staging of the follow-up verification prompt ("Is it fixed?") into the verification queue.
- **Workflow**: Continuous N-Step Self-Loop across 2-Agent Concurrency and Parallel Micro-Batches.
- **Total Steps/Loops**: 4 subtasks executed in parallel micro-batches.
- **Status**: 100% Completed & Verified.

---

## Consolidated Subtasks Ledger

### Subtask 01: Direct `fix agy` and `fix errors agy` Routing
- **Problem**: When users ran `gitmap fix agy` or `gitmap fix errors agy`, the command dispatched to `runFix` in `cli/cmd/fix_cmd.go`, which attempted to look up a repository named `"agy"` in remediation state, causing a "repository not found" error instead of delegating to Antigravity pipeline error resolution.
- **Implementation**:
  - `cli/cmd/fix_cmd.go` (modified): Added `isFixAgyRequest(args []string) bool` checking for `"agy"`, `"aef"`, or compound error tokens. When present, `runFix` immediately routes to `cmdagy.RunPipelineFixAgyCLI(args)`.
  - Refactored `runFix` and extracted `executeFixTarget` to ensure all functions strictly adhere to $\le 15$ lines.

### Subtask 02: Compound AGY Subcommand Normalization
- **Problem**: Invocations such as `gitmap agy errors fix` or `gitmap agy fix errors` were treated by Cobra as unknown subcommands because Cobra only inspects single-token subcommands by default.
- **Implementation**:
  - `cli/cmdagy/agy_cmd.go` (modified): Implemented `normalizeAgyArgs(args []string) []string`, `isCompoundAgyFix(args []string) bool`, `rewriteCompoundAgyFix`, and `stripAgyPrefix`.
  - Automatically rewrites compound patterns (`errors fix`, `fix errors`, `fix pipeline`, `aef`) into the canonical `fix-pipeline` subcommand while preserving all trailing arguments and flags (`--force`, `-f`, `--detailed`, `--dry-run`, `-p`).

### Subtask 03: Unit Testing All Command Syntaxes and Routing
- **Implementation**:
  - `cli/cmd/fix_cmd_test.go` (new): Implemented `TestIsFixAgyRequest_Variations` testing `fix agy`, `fix errors agy`, `fix agy errors`, `fix aef`, `fix agy-errors-fix`, regular repositories, and empty args.
  - `cli/cmdagy/agy_pipeline_fix_errors_test.go` (modified): Implemented `TestNormalizeAgyArgs_CompoundFixPhrases` testing `errors fix`, `fix errors`, `fix pipeline`, `aef`, and normal subcommands.

### Subtask 04: Verification, Linters, and Cleanliness
- **Verification**:
  - All 3 targeted linters passed with 0 violations:
    - `python linter-scripts/check-nested-ifs.py --changed-only` (PASS)
    - `python linter-scripts/check-boolean-guidelines.py --changed-only` (PASS)
    - `python linter-scripts/check-error-management.py --changed-only` (PASS)
  - Recorded 4 modified files to `.ai-memory/temp/recent-file-changes.json` under lock.

---

## Command Syntax & Feature Verification Matrix

| Syntax Variant | Subsystem Route | Status | Notes |
| :--- | :--- | :--- | :--- |
| `gitmap pipeline fix errors agy` | `cmdpipeline.runPipeline` $\to$ `PipelineAgyFixRunner` | Supported | Canonical multi-word pipeline command |
| `gitmap pipeline-fix errors agy` | `cmd/rootutility` $\to$ `cmdagy.RunPipelineFixAgyCLI` | Supported | Direct hyphenated utility command |
| `gitmap pipeline agy errors fix` | `cmdpipeline.runPipeline` $\to$ `PipelineAgyFixRunner` | Supported | Reversed word ordering |
| `gitmap pipeline agy-errors-fix` | `cmdpipeline.runPipeline` $\to$ `PipelineAgyFixRunner` | Supported | Hyphenated compound |
| `gitmap pipeline aef` | `cmdpipeline.runPipeline` $\to$ `PipelineAgyFixRunner` | Supported | Acronym shorthand alias |
| `gitmap aef` | `cmd/rootutility` $\to$ `cmdagy.RunPipelineFixAgyCLI` | Supported | Top-level shortcut alias |
| `gitmap fix agy` / `gitmap fix-agy` | `cli/cmd/fix_cmd.go` $\to$ `isFixAgyRequest` | Supported | Direct fix namespace routing |
| `gitmap fix errors agy` / `gitmap fix agy errors` | `cli/cmd/fix_cmd.go` $\to$ `isFixAgyRequest` | Supported | Direct fix error routing |
| `gitmap agy errors fix` / `gitmap agy fix errors` | `cli/cmdagy/agy_cmd.go` $\to$ `normalizeAgyArgs` | Supported | Compound agy subcommand normalization |
| Duplicate Detection | `ComputeErrorSignature` + `sent_agy_errors.json` | Supported | Prints warning and asks to use `--force` or `-f` |
| Force Bypass | `--force` / `-f` | Supported | Bypasses duplicate warning and updates records |
| Automated Context Extraction | `ExtractGitLog` (-n 5 --stat) + Error Report | Supported | Fully embedded before RCA prompt |
| Two-Prompt Feedback Queue | Primary RCA file + Queued "Is it fixed?" file | Supported | Active prompt copied to clipboard; follow-up staged |
