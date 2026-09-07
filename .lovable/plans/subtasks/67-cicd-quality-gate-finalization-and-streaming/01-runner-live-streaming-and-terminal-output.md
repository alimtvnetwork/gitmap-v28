# Subtask 01: Runner Live Streaming, Telemetry & Immediate Terminal Failure Output

## Objective
Refactor `03-ai-scripts/06-cicd-local-runner.py` and `.lovable/ai-fix-scripts/06-cicd-local-runner.py` to:
1. Append all gate outcomes to `CICD_RUN_LOG` (`.lovable/temp/cicd/run.log`).
2. Continuously update `CICD_SUMMARY_JSON` (`.lovable/temp/cicd/summary.json`) with active counts (`total_gates`, `passed_gates`, `failed_gates`, `remaining_gates`, `status`).
3. Add `extract_failing_files` to identify suspected source files from failure stderr/stdout and show them in terminal banners.
4. Protect JSON stdout mode by avoiding printing ANSI banners when `--json` is set.
5. Replace swallowed `except Exception: pass` blocks with typed exceptions and logging.
6. Remove unused imports and align docstring usage paths.

## Files
- `03-ai-scripts/06-cicd-local-runner.py`
- `.lovable/ai-fix-scripts/06-cicd-local-runner.py`

## Acceptance Criteria
- [ ] `run.log` receives timestamps and entries for every completed gate.
- [ ] `summary.json` reflects live execution progress and terminates in `"status": "completed"` or `"failed"`.
- [ ] Fails immediately display full stack trace and suspect files in terminal.
- [ ] `--json` produces valid JSON output without banner pollution.
