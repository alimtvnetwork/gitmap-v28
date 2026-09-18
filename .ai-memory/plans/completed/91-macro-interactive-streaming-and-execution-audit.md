# 91-macro-interactive-streaming-and-execution-audit

## Overview
Comprehensive audit and enhancement of macro creation (`gitmap macro add`), interactive prompts, in-builder helpers, and live terminal streaming during macro execution (`gitmap macro run`). This plan documents all confirmed capabilities (real-time unbuffered terminal streaming, live command execution, in-builder inspection tools, universal path expansion) and addresses key gaps identified during the audit (stdin attachment in live builder, error feedback on live execution failure, local file priority in open launcher, and step timeout enforcement).

## Audit Findings & Baseline
1. **Interactive Macro Creation & Live Execution**:
   - `gitmap macro add <name> --exec` (or `-e`) enables live execution mode.
   - Inside the builder, `exec on` / `:exec on` enables live execution, and `exec off` / `exec` toggles it.
   - All entered commands run live with output streaming to `os.Stdout` and `os.Stderr`.
   - In-builder helpers: `ls`, `dir`, `pwd`, `pwd on/off`, `cd`, `mkdir`, `find`, `search`, `replace`, `+add`.
   - PWD banner `[PWD: <cwd>]` displays before each prompt.
   - Gaps: `cmd.Stdin` was not attached to `os.Stdin` in `executeLiveCommand`, causing interactive prompts to fail; live command exit errors were ignored without user feedback.
2. **Macro Execution & Live Streaming (`gitmap macro run`)**:
   - Uses `io.MultiWriter(os.Stdout, outBuf)` and `io.MultiWriter(os.Stderr, errBuf)` for unbuffered live streaming.
   - Step progress `[ 1/N] ➜ <cmd>`, directory transitions `➜ 📁 Directory: <dir>`, and timers `✔ ok (X.Xs)` provide real-time visibility.
   - Universal path expansion (`%TEMP%`, `//temp`, `~`, `$VAR`) works across steps via `ExpandPathAndEnv`.
   - Gaps: In `macro/open.go`, `parseURLTarget` was checked before local file existence, misidentifying local dot files (e.g. `sample.txt`, `readme.txt`) as URLs; `step.TimeoutSeconds` was not applied via `context.WithTimeout`; step completion lacked 2-space indentation.

## Subtask Breakdown
- `01-task-builder-live-stdin-and-error-feedback.md`: In `gitmap/cmd/macro_add_interactive.go`, wire `cmd.Stdin = os.Stdin` in `executeLiveCommand`, inspect `cmd.Run()` exit errors, and print clear alert feedback if a live command fails.
- `02-task-engine-open-local-files-and-timeout.md`: In `gitmap/macro/open.go`, check local file existence before URL detection; in `gitmap/macro/execute.go`, apply `step.TimeoutSeconds` via `context.WithTimeout`, honor `step.WorkingDir` precedence, and align terminal indentation.
- `03-task-verification-tests-and-quality-gates.md`: Add unit tests for live execution, local file opening, path expansion, and run all CI quality gates (`06-cicd-local-runner.py`).
