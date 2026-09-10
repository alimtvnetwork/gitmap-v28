# Subtask 04: Python Purge History Script Hardening & Line Budget Optimization

## Objective
Refactor `03-ai-scripts/30-purge-history.py`:
- Bring file length below 200 lines (currently 209 lines).
- Bring functions below 15 lines (e.g. break down `purge_history` and `restore_history`).
- Eliminate non-Windows `AttributeError: module 'ctypes' has no attribute 'windll'` by guarding `ctypes.windll` with `sys.platform == "win32"`.
- Replace `shell=True` with structured argument lists in `subprocess.run` to prevent shell injection and path quoting discrepancies on Windows.
- Normalize path globs with `Path(pattern).as_posix()`.
- Update boolean naming to affirmative prefixes (`is_auto_confirm`, `is_restore`, `should_check`).

## Files Affected
- `03-ai-scripts/30-purge-history.py`
