---
name: pipeline-errors-agy-fix-injection-and-multi-project-batching
description: Autonomously implement and verify Antigravity (AGY) CI/CD pipeline error fix injection, embedding full failing logs, full file path terminal display, direct prompt dispatch into Antigravity IDE/CLI, queue management, and multi-project parallel batching with configurable project limits.
---

# Pipeline Errors AGY Fix Injection and Multi-Project Batching

## Overview
Autonomously diagnose, repair, and verify `gitmap pipeline errors agy fix` and `gitmap agy fix-pipeline` across single-repo and multi-project contexts:
1. **Full Error Logs Inside Text File Payload:** Ensure `active-agy-pipeline-fix-prompt.txt` embeds the full failing CI/CD error report, recent git commit history, and canonical 4-part RCA directives.
2. **Follow-up Verification Prompt Isolation:** Staged into `queued-agy-followup-prompt.txt` without overwriting the active fix payload.
3. **Queue Ledger Integrity:** Both active fix task and queued follow-up verification are registered in `agy-prompt-queue.json`.
4. **100% Full Absolute File Paths in Terminal:** In all terminal outputs, display strictly full absolute paths on disk for every generated file.
5. **Direct Antigravity Injection:** Launch `agy.exe` with `--dangerously-skip-permissions` in background detached mode, setting `cmd.Dir` to the target repository directory.
6. **Multi-Project Parallel Batching:** Concurrently scan pipeline errors across tracked projects with a default limit of 3 projects per run, persisting pagination cursor in `pipeline-fix-batch-cursor.json`.
7. **Comprehensive Documentation Parity:** Maintain full parity across CLI markdown helptext (`cli/helptext/pipeline.md`, `cli/helptext/agy-fix-pipeline.md`, `cli/helptext/agy.md`), Web UI command definitions (`src/data/commands.ts`), and command examples.

## Key Architecture & Components
- `cli/cmdagy/agy_fix_pipeline_assembly.go`: `AssembleRcaFixPayload` combines header, git log, error report, and RCA template.
- `cli/cmdagy/agy_fix_pipeline_persist.go`: Persists primary payload to `.ai-memory/temp/active-agy-pipeline-fix-prompt.txt`.
- `cli/cmdagy/agy_fix_pipeline_queue.go`: Persists queued verification prompt to `.ai-memory/temp/queued-agy-followup-prompt.txt` and manages `.ai-memory/temp/agy-prompt-queue.json`.
- `cli/cmdagy/agy_fix_pipeline_inject.go`: Dispatches direct injection to `agy.exe` with background detached process attributes.
- `cli/cmdagy/agy_proc_windows.go` & `agy_proc_other.go`: Detaches background processes on Windows using `DETACHED_PROCESS`.
- `cli/cmdagy/agy_fix_pipeline_render.go`: Ensures all file paths printed to the terminal are converted via `toAbsPath(filepath.Clean(...))`.
- `cli/cmdagy/agy_fix_pipeline_batch_scan.go` & `batch_cursor.go`: Handles multi-project candidate discovery, error inspection, and cursor pagination.
