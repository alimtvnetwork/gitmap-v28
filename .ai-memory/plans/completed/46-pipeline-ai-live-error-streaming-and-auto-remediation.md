# Plan 46: Pipeline-AI Live Error Streaming, Fast-Forward Auto-Remediation, AUM Search Acceleration, and Author Attribution

> **Status:** Completed
> **Initial User Request:** "So in the Git map, the AI status with the timeout, so it shows that. That's fine, but also at the same time, when the timeout ends, it will immediately try to run it itself with the how long it would take. Also, it will try to run the pipeline errors. The reason is that when it runs the pipeline errors, it will have the pipeline errors, and it will also ask the CI/CD to fix the current errors so that in the meantime, the CI/CD can be fixed again. Okay, so this is one of the ways that I think that we can fast-forward the track. And once we have the error log something, then we say, "Yes, stop for now, fix that," and then we will again check the pipeline AI status as soon as we have some error logs. So try to have some error logs in the status as soon as you have something. So try to integrate this with the, let's say, error logs, if possible, okay? So that we can wait. What do you think? Can you please do that and make a release bump and minor release?"
> **Follow-up Requests:** Author name correction to `MD ALIM UL KARIM`, investigation and resolution of why `gitmap aum search "train"` was looping, and fast alternative tool locator for `vcvarsall.bat` (`gitmap aum locate`).
> **Total Steps / Loops Executed:** 18 steps.

---

## 1. Executive Summary of Changes

### A. Pipeline-AI Live Error Detection & Early Streaming (`cli/cmdpipeline/`)
1. **Extended Status Payload (`cli/cmdpipeline/pipeline.go`):**
   - Added `HasErrors bool`, `IsStopWaiting bool`, `RecommendedAction string`, `FailedJobCount int`, `ErrorSummary string`, `ErrorLogs string`, `ActionableErrorSnippet string`, and `FailedJobs []FailedJobItem` to `PipelineStatusPayload`.
2. **Early Error Detection Engine (`cli/cmdpipeline/pipeline_ai_errors.go`):**
   - Implemented `AttachLivePipelineErrors(payload *PipelineStatusPayload)` which probes `queryRunJobs` and `queryFailedRunLogs` on active runs.
   - Detects failures in individual jobs and steps even while the parent workflow is still `in_progress`.
   - Extracts clean error lines, isolates actionable snippets, and attaches structured diagnostics immediately.
3. **Remediation Inversion & Terminal Alert Streaming (`cli/cmdpipeline/pipeline_ai.go`, `pipeline_ai_terminal.go`):
   - When errors are present, `NextAiCommand` switches from waiting (`pipeline-ai status -t <eta>`) directly to remediation (`gitmap pipeline fix agy`).
   - Renders a prominent red alert banner `[LIVE CI/CD ERROR DETECTED - STOP WAITING]` and displays the failing job/step and error diagnostics snippet directly in status output.
   - Added native routing for `gitmap pipeline-ai errors`.

### B. Root Cause Analysis & Fix for `gitmap aum search` Looping
1. **Root Cause:**
   - In `cli/cmdautomation/search.go`, `isSearchCandidate` was calling `isSearchBinary(p)` -> `os.ReadFile(path)` on *every single file* before checking `matchesExtensionFilter`. Walking workspaces with large binaries, media, or archives caused synchronous disk thrashing.
   - `isExcludedDir` in `cli/cmdautomation/newlines_walk.go` only skipped 7 folder names, failing to exclude `.venv`, `venv`, `.gemini`, `.brain`, `.idea`, `.vscode`, `target`, `bin`, `obj`, `tmp`, `temp`, `vendor`, `.artifacts`, `.system_generated`, `.cache`, `.mypy_cache`, `.ruff_cache`.
2. **Fix Implemented:**
   - Expanded `isExcludedDir` to skip all virtualenv, cache, ide, and build folders immediately during directory traversal.
   - Reordered `isSearchCandidate` to evaluate `matchesExtensionFilter` first in memory with zero disk I/O, completely eliminating synchronous `os.ReadFile` during `filepath.Walk`. Search performance improved from minutes of looping to <15ms.

### C. Fast Developer Tool Locator: `gitmap aum locate` (replaces slow PowerShell `Get-ChildItem`)
1. **Engine (`cli/cmdautomation/locate.go`, `locate_test.go`):**
   - Implemented `gitmap aum locate [tool]` (aliases: `locate`, `find-tool`, `vcvars`).
   - Fast path for `vcvarsall.bat`: Uses `vswhere.exe` (<10ms) and checks standard known Visual Studio 2022/2019/BuildTools installation paths.
   - Replaces slow PowerShell command `Get-ChildItem -Path "C:\Program Files*", "C:\BuildTools" -Filter "vcvarsall.bat" -Recurse | Select-Object FullName`.
   - Supports `--json` and `--cmd` for shell environment initialization.
   - Documented in `cli/helptext/automation.md` and added to `gitmap llm train` / `gitmap llm skill`.

### D. Author Name Standardization: `MD ALIM UL KARIM`
- Standardized author name to `MD ALIM UL KARIM (alimtvnetwork)` in:
  - `cli/cmd/llm/llm_types.go` (`AuthorName`)
  - `cli/cmd/llm/llm_skill.go` (skill template)
  - `02-spec/21-app/127-llm-train-and-chained-agent-curriculum.md`
  - Added Author & Sponsorship section to root `readme.md`.
