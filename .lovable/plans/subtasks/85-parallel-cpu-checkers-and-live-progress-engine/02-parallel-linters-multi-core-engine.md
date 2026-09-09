# Subtask 02: Parallel Linters Multi-Core Engine

## Objective
Refactor repository linters (`linter-scripts/check-nested-ifs.py` and `linter-scripts/check-enum-and-boolean.py`) to gather all target files upfront and scan them in parallel using `concurrent.futures.ThreadPoolExecutor(max_workers=os.cpu_count())` with live progress percentages.

## Files to Touch
- `linter-scripts/check-nested-ifs.py`
- `linter-scripts/check-enum-and-boolean.py`

## Implementation Steps
1. In `linter-scripts/check-nested-ifs.py`:
   - Pre-gather all target files (`.go`, `.ts`, `.tsx`, `.js`, `.jsx`, `.py`, `.php`) into a list upfront.
   - Use `ThreadPoolExecutor(max_workers=os.cpu_count() or 16)` to execute `scan_file` concurrently across all cores.
   - Maintain a thread-safe completed counter and emit live progress percentage every 50 files or 0.2s:
     `Scanning: [ 650/2500 ] 26% | 16 workers`
   - Aggregate violations and exit 0 or 1.
2. In `linter-scripts/check-enum-and-boolean.py`:
   - Pre-gather all target files into a list upfront.
   - Use `ThreadPoolExecutor(max_workers=os.cpu_count() or 16)` to execute `check_file` concurrently.
   - Emit live progress percentage.
   - Aggregate violations and exit 0 or 1.
3. Test execution and verify 10x+ speedup and high CPU utilization.
