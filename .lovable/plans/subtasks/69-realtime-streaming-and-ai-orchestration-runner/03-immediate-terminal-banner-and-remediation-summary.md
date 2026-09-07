# Subtask 03: Immediate Terminal Failure Banners & Post-Run AI Remediation Summary

## Objective
Update `format_failure_banner` to include `Working Dir`, `Env Overrides`, and `Stream Events`, and implement a post-run `AI AGENT REMEDIATION SUMMARY & LOG LOCATIONS` banner in `03-ai-scripts/06-cicd-local-runner.py`.

## Requirements
1. **Immediate Failure Banner Enhancements**:
   - Display `Command`, `Working Dir` (e.g. `.` or `gitmap`), and `Env Overrides`.
   - List paths to `errors.log`, `errors.json`, and `events.jsonl`.
   - Flush output immediately with ANSI formatting.
2. **Post-Execution AI Remediation Banner**:
   - If any gate failed or timed out, display an aggregated summary banner at suite conclusion:
     - Header with total failed gate count.
     - Exact relative paths to `errors.log`, `errors.json`, `summary.json`, `state.json`, `run.log`, and session folder.
     - Deduplicated list of suspect files across all failing gates.
     - Copy-pasteable targeted re-test commands for both runner filter (`--filter "<Gate Name>"`) and direct script execution.
     - Actionable step-by-step remediation instructions for AI agents.
3. **Report Integration**:
   - Wire `print_failure_report` into `handle_text_output` so the remediation banner prints right after the summary table.
4. **Coding Guidelines**:
   - All functions $\le 15$ lines.
   - Blank line before every return statement.
   - Affirmative booleans (`is_*`, `has_*`).
   - Zero swallowed exceptions.

## Target Files
- `03-ai-scripts/06-cicd-local-runner.py`
- `.lovable/ai-fix-scripts/06-cicd-local-runner.py`

## Acceptance Criteria
- [x] Terminal failure banner includes `Working Dir` and `Env Overrides`.
- [x] Post-run failure report prints the `AI AGENT REMEDIATION SUMMARY & LOG LOCATIONS` banner.
- [x] Suspect files are deduplicated and displayed clearly.
- [x] Single-gate re-test commands are generated for each failed gate.
