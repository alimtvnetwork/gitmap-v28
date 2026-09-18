# Subtask 02: Fix Payload Assembly, Full Error Logs Embedding & Full File Paths

## Objective
Audit and fix prompt assembly to ensure the failing pipeline error logs are always embedded in `active-agy-pipeline-fix-prompt.txt`, metrics calculate correct component sizes, and terminal output strictly prints full absolute file paths.

## Target Files
- `cli/cmdagy/agy_fix_pipeline.go`
- `cli/cmdagy/agy_fix_pipeline_assembly.go`
- `cli/cmdagy/agy_fix_pipeline_queue.go`

## Status
Completed: 2026-09-18

## Verification
- `active-agy-pipeline-fix-prompt.txt` strictly receives `primaryPayload` containing RCA header, target repo, workflow run ID, commit SHA, recent git commits, full failing pipeline error logs, and canonical 4-part RCA prompt template.
- `queued-agy-followup-prompt.txt` strictly receives `followupPrompt` containing `# Verification Check: Is It Fixed?`.
- Metric calculations fixed: `Pipeline Logs: len(logs)`, `Fix Prompt: len(promptContent)`, `Total Payload: len(primaryPayload)`.
- All terminal output file paths (`Saved Payload`, `Follow-up Verification Prompt Queued`, `Queue Ledger`) converted to full absolute paths via `toAbsPath(resolve...())`.
- Zero nested ifs; all files $\le 100$ lines, functions $\le 15$ lines.
- `go vet ./cmdagy/...`: PASS (code 0).
