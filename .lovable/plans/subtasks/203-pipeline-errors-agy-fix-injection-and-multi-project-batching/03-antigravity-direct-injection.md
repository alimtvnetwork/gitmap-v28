# Subtask 03: Direct Antigravity (AGY) Injection & Dispatch

## Objective
Implement direct prompt injection and execution into Antigravity (`agy` CLI or IDE session), transitioning from passive clipboard copying to active agent dispatch.

## Target Files
- `cli/cmdagy/agy_fix_pipeline_inject.go` (new)
- `cli/cmdagy/agy_fix_pipeline.go`

## Status
Completed: 2026-09-18

## Verification
- Implemented `InjectAgyFixTask` in `cli/cmdagy/agy_fix_pipeline_inject.go`.
- Resolves Antigravity binary (`resolveAntigravityBinary()` finds `agy.exe` on PATH).
- Dispatches prompt to `agy -p "<directive referencing full absolute payload path>"` in background to prevent process hanging or command-line length limits.
- Supports `--no-inject` flag to bypass direct injection if desired.
- Clipboard copy is preserved as a seamless backup/fallback.
- File is 32 lines, functions $\le 15$ lines.
- `go vet ./cmdagy/...`: PASS (code 0).
