# 74-pipeline-errorlogs-details-and-cross-platform-ci-fixes.md: Pipeline Error Logs Multi-Run Aggregation, Parsed Error Details & Cross-Platform CI Fixes

## 1. Executive Summary

This plan addresses runtime issues reported in pipeline failure diagnostics and CI/CD workflows for GitHub Actions:
1. **`gitmap pipeline errorlogs` / `error-logs` / `errors` Aggregation & Detailed Failure Extraction**:
   - Currently, `buildErrorLogsPayload` only queries the first failed run in `gh run list` and ignores other failed workflows from the same commit/push (e.g. `Release #34146129254` shadowed `Cross-Platform Build #34146127608` and `GoReleaser Smoke #34146127554`).
   - Query all failed runs matching the current commit/push/branch.
   - Parse raw `Job\tStep\tTimestamp\tLogText` lines into structured failures with job name, step name, failure message, root cause snippet (e.g. `--- FAIL: TestName` and error assertion lines), and URL.
   - In JSON view (`--json`): Emit structured `failedRuns` and `failedJobs` arrays with concatenated clean error details, avoiding raw unmarshaled escape soup.
   - In terminal view: Render multi-failure cards with job title, failed step, extracted assertion reasoning, and direct run URL.
2. **Fix `Cross-Platform Build` macOS Failure (`gitmap/power/driver_darwin.go` & `manager_test.go`)**:
   - `darwinDriver.GetStatus()` hardcodes `IsNeverSleep: false`.
   - `TestMockRunner_DriverInteraction` in `manager_test.go` was written with Windows `powercfg` mocks, failing on macOS (`Expected IsNeverSleep to be true`).
   - Implement `pmset -g` parsing in `darwinDriver` to detect `displaysleep 0` and `sleep 0` as `IsNeverSleep: true`.
   - Update `manager_test.go` with cross-platform mock runner support for macOS (`pmset`) and Linux (`gsettings`).
3. **Fix `GoReleaser Smoke (cfr cg)` Windows Failure (`codingguidelines_compat.go` & `goreleaser-smoke.yml`)**:
   - In `patchCGWindowsScriptFile`, PowerShell 5.1 in Windows runner failed with `Unexpected token 'local' in expression or statement` at `$TarballUrl = "(none — local archive)"` because the downloaded script lacks a UTF-8 BOM, causing Windows-1252 ANSI interpretation of `—` (`\xe2\x80\x94`).
   - Prepend UTF-8 BOM (`\xef\xbb\xbf`) when writing `install.ps1` in `patchCGWindowsScriptFile`.
   - In `.github/workflows/goreleaser-smoke.yml`: Print `$logText` (`Get-Content $log`) when `$rc -ne 0` so failure logs are visible directly in console and fetched by `gitmap pipeline errorlogs`.
4. **Fix `Release` Workflow Race Condition Collision (`.github/workflows/release.yml` & `29-release-orchestrator.py`)**:
   - `Release` workflow currently triggers on both push to `release/*` AND push to tags `v*`.
   - When `push_release_artifacts` pushed both branch `release/vX.Y.Z` and tag `vX.Y.Z` simultaneously, two concurrent runs executed `softprops/action-gh-release@v2`, colliding on asset uploads with `HttpError: Not Found - update-a-release-asset`.
   - Restrict `Release` workflow trigger exclusively to tags `v*` (or decouple branch push).
5. **Confirm Chrome Profile Import Status**:
   - Provide a clear report to the user confirming that Chrome Profile Import (`gitmap import-all`, ZIP discovery, and `import-check`) was resolved and verified, and explain that updating the Ubuntu environment to `v6.197.0` will deliver the fix.

## 2. Task-Specific Rules & Constraints

1. **Rule 1 (Strict Relative Paths):** Zero `file:///` URIs or machine-specific absolute paths in markdown files, artifacts, or code. All paths must be relative to repository root (`/`).
2. **Rule 2 (Coding Guidelines Compliance):** Functions <= 15 lines (prefer <= 8), blank line before every return statement, affirmative booleans (`is_*`, `has_*`), no magic numbers, and zero swallowed errors (`err != nil` must be handled or wrapped via `apperror`).
3. **Rule 3 (Strict 200-Line File Limit):** All new and modified Go source files must remain <= 200 lines. Decompose logic into focused single-responsibility files.
4. **Rule 4 (AST Parity & Constants):** All CLI verbs and help entries must be synchronized with `gitmap/constants/constants_cli.go`.
5. **Rule 5 (CI/CD Local Runner Validation):** Must run `python 03-ai-scripts/06-cicd-local-runner.py` exit 0 before release orchestration.

## 3. Subtasks Breakdown

- [01-pipeline-multi-run-and-parsed-failure-extraction.md](../subtasks/74-pipeline-errorlogs-details-and-cross-platform-ci-fixes/01-pipeline-multi-run-and-parsed-failure-extraction.md)
- [02-pipeline-errorlogs-json-and-terminal-rendering.md](../subtasks/74-pipeline-errorlogs-details-and-cross-platform-ci-fixes/02-pipeline-errorlogs-json-and-terminal-rendering.md)
- [03-fix-macos-power-driver-and-test.md](../subtasks/74-pipeline-errorlogs-details-and-cross-platform-ci-fixes/03-fix-macos-power-driver-and-test.md)
- [04-fix-windows-cg-compat-utf8-bom-and-smoke-logging.md](../subtasks/74-pipeline-errorlogs-details-and-cross-platform-ci-fixes/04-fix-windows-cg-compat-utf8-bom-and-smoke-logging.md)
- [05-prevent-release-workflow-branch-collision.md](../subtasks/74-pipeline-errorlogs-details-and-cross-platform-ci-fixes/05-prevent-release-workflow-branch-collision.md)
- [06-quality-gates-ci-validation-and-release.md](../subtasks/74-pipeline-errorlogs-details-and-cross-platform-ci-fixes/06-quality-gates-ci-validation-and-release.md)
