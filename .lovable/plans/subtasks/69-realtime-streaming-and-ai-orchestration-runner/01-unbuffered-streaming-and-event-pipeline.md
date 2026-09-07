# Subtask 01: Unbuffered Synchronized Streaming, Event Dispatcher & Windows File Lock Resilience

## Objective
Implement unbuffered immediate disk streaming with `os.fsync`, line-delimited `events.jsonl` and `changelog.log` pipelines, Windows file lock retry mechanisms, and early `summary.json` initialization in `03-ai-scripts/06-cicd-local-runner.py`.

## Requirements
1. **Unbuffered Streaming**:
   - Create `direct_append_sync(file_path: Path, content: str)` with a thread lock, immediate `fh.flush()`, and `os.fsync(fh.fileno())`.
   - Wire `direct_append_sync` into `write_failure_markdown` (`errors.log`), `append_run_log_entry` (`run.log`), and event dispatchers.
2. **Event Dispatcher (`events.jsonl` & `changelog.log`)**:
   - Define constants: `CICD_EVENTS_JSONL = CICD_TEMP_DIR / "events.jsonl"`, `CICD_CHANGELOG_LOG = CICD_TEMP_DIR / "changelog.log"`.
   - Implement `emit_telemetry_event(event_type: str, session_dir: Path | None, payload: dict)` appending to both root and session files.
   - Emit events on: `run_started`, `gate_started`, `gate_cached`, `gate_passed`, `gate_failed`, `gate_timeout`, `run_completed`.
3. **Windows Sharing Violation Resilience**:
   - Enhance `atomic_write_text` with an exponential backoff loop (up to 5 retries) catching `PermissionError` [WinError 32] caused by concurrent tailing agents holding read handles on Windows, with a direct overwrite fallback.
4. **Early `summary.json` Initialization**:
   - In `init_files`, initialize `summary.json` with `"status": "running"`, `total_gates`, and `remaining_gates` so agents never inspect stale data from prior runs.
5. **Coding Guidelines**:
   - All functions $\le 15$ lines.
   - Blank line before every return statement.
   - Affirmative booleans (`is_*`, `has_*`).
   - Zero swallowed exceptions.

## Target Files
- `03-ai-scripts/06-cicd-local-runner.py`
- `.lovable/ai-fix-scripts/06-cicd-local-runner.py`

## Acceptance Criteria
- [x] Every line written to `errors.log`, `run.log`, `events.jsonl`, and `changelog.log` is flushed with `os.fsync`.
- [x] `events.jsonl` contains valid NDJSON records for each gate lifecycle transition.
- [x] `summary.json` is initialized immediately upon runner startup with `"status": "running"`.
- [x] Windows sharing violations on JSON files are retried and recovered cleanly.
