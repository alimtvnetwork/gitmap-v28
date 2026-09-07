# Plan 72: CI/CD Pipeline Remote Log Fetching, Caching & Failure Remediation

## Executive Summary
This plan addresses the remote CI/CD pipeline error log retrieval issues in `gitmap pipeline errors` (5s timeout deadline exceeded, swallowed errors, lack of persistence, interactive prompts) and resolves the root causes of failing CI/CD runs (run `#34115544541`, `#34115544334`, `#34115544327`, `#34115544315`):
1. **Pipeline Error Logs & Caching**:
   - Adaptive/increased timeout (60s) for log downloads in `gitmap/cmd/pipeline_query.go`.
   - Detailed error reporting on failure (stdout/stderr) rather than static fallback strings.
   - Persistence and caching to `.gitmap/pipeline/<runId>.log` and `.gitmap/pipeline/<runId>.json` (configurable via SQLite setting `pipeline.dir`).
   - Non-interactive default output by removing unsolicited prompts unless `--fix` is passed.
   - Full flag parity (`--json`).
2. **Process Lock Reentrancy & WorkDir Test Contention**:
   - Add same-process reentrancy ref counting in `gitmap/store/lock.go` so recursive DB openings (same PID) do not deadlock or time out with `ErrLockHeld`.
   - Clean up DB handle lifetimes in `gitmap/cmd/workdir_test.go`.
3. **Windows Coding Guidelines Installer & Smoke Workflow**:
   - Download and patch upstream syntax bugs in `error-manage-install.ps1` (`$oldFile:`, `$destPath:`, `$targetVersionFile:`) in `gitmap/cmd/codingguidelines.go` before invoking `pwsh`/`powershell`.
   - Fix pipeline output leaking in `Invoke-CfrCg` within `.github/workflows/goreleaser-smoke.yml`.
4. **Validation & Verification**:
   - Pass all unit tests, AST parity checks, nested-if linter, and all 33 quality gates locally via `03-ai-scripts/06-cicd-local-runner.py`.

## Subtasks
- [01-pipeline-error-logs-and-caching.md](../subtasks/72-pipeline-error-logs-caching-and-cicd-fixes/01-pipeline-error-logs-and-caching.md)
- [02-lock-reentrancy-and-workdir-tests.md](../subtasks/72-pipeline-error-logs-caching-and-cicd-fixes/02-lock-reentrancy-and-workdir-tests.md)
- [03-windows-coding-guidelines-and-smoke.md](../subtasks/72-pipeline-error-logs-caching-and-cicd-fixes/03-windows-coding-guidelines-and-smoke.md)
- [04-quality-gates-and-verification.md](../subtasks/72-pipeline-error-logs-caching-and-cicd-fixes/04-quality-gates-and-verification.md)
- [05-ubuntu-os-fix-link-command-and-ui-help.md](../subtasks/72-pipeline-error-logs-caching-and-cicd-fixes/05-ubuntu-os-fix-link-command-and-ui-help.md)
