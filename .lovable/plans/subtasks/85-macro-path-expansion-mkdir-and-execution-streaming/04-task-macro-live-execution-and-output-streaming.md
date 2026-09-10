# Subtask 04: Macro Live Execution & Real-Time Output Streaming

## 1. Objective
Enable real-time stdout and stderr output streaming during macro step execution (`macro.Execute`) and provide live execution feedback in the interactive macro builder (`gitmap macro add`).

## 2. Target Files
- `gitmap/macro/execute.go`
- `gitmap/macro/execute_test.go`
- `gitmap/cmd/macro_add_interactive.go`
- `gitmap/cmd/macro_add_helpers.go`

## 3. Requirements
- In `macro/execute.go`:
  - Wire `cmd.Stdout` and `cmd.Stderr` via `io.MultiWriter` so output streams to terminal in real time while also buffering for `ExecutionReport` (unless structured output `--json`/`--yaml` is requested).
  - On step failure, display clear error diagnostic output so compiler, installer, and command errors are immediately visible.
- In `cmd/macro_add_interactive.go`:
  - When user types non-helper commands at `Step N> <cmd>` (e.g. `gitmap mkdir -p test`, `gitmap install vscode`), execute them live in the session working directory using `macro.RecordLive` or dual execution, stream output to terminal, and record the step.
  - Provide an interactive toggle `exec on` / `exec off` (default `exec on`).
  - When `cd <dir>` is run, execute directory change and stage the step into `*steps`.
- All Go functions <= 15 lines with blank lines before every return.
