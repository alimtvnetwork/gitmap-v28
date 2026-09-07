# Subtask 01: Session Scaffolding, Multi-Session Directory Hierarchy & Atomic State Persistence

## Objective
Implement multi-session run directory management and crash-safe state persistence in `03-ai-scripts/06-cicd-local-runner.py`.

## Requirements
1. **Multi-Session Isolation**:
   - Create timestamped run directory: `.lovable/temp/cicd/runs/<timestamp>/`.
   - Maintain symlink / Windows junction / pointer at `.lovable/temp/cicd/latest/` pointing to the active session.
   - Maintain backward-compatible mirrors at `.lovable/temp/cicd/` root: `errors.log`, `errors.json`, `run.log`, `summary.json`, `state.json`.
2. **StateStore Machine**:
   - Track session metadata: `session_id`, `started_at`, `status`, `pid`, `batches`, `active_slots`, `gates`.
   - Gate states: `PENDING`, `RUNNING`, `PASSED`, `FAILED`, `TIMEOUT`, `SKIPPED`, `CRASHED`.
   - Atomic disk write via `write-tmp-flush-rename` (`os.replace`) to prevent file corruption.
3. **Coding Guidelines**:
   - All functions $\le 15$ lines.
   - Blank line before every return statement.
   - Affirmative booleans (`is_*`, `has_*`).
   - Zero swallowed exceptions.

## Target Files
- `03-ai-scripts/06-cicd-local-runner.py`
- `.lovable/ai-fix-scripts/06-cicd-local-runner.py`

## Acceptance Criteria
- [x] Session directory `.lovable/temp/cicd/runs/<timestamp>/` created on startup.
- [x] `state.json` written atomically with PID, batch, and gate statuses.
- [x] Root mirror files kept up to date for backwards compatibility.
