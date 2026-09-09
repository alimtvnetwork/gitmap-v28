# Subtask 01: Parallel Go Code Formatter and Checker

## Objective
Refactor `03-ai-scripts/26-go-code-formatter.py` and `.github/scripts/go-format-check.py` to list all Go files upfront using `03-ai-scripts/02-shared-engine.py` (`stream_directory_files`) or fast recursive file gathering, distribute files into parallel batches across all CPU cores (`os.cpu_count()`), and output live progress percentages.

## Files to Touch
- `03-ai-scripts/26-go-code-formatter.py`
- `.github/scripts/go-format-check.py`

## Implementation Steps
1. In `03-ai-scripts/26-go-code-formatter.py`:
   - Pre-gather all `.go` files into `target_files` list upfront using `stream_directory_files(repo_root, extensions=[".go"])`.
   - Chunk `target_files` into batches (e.g. 20-50 files per chunk, avoiding Windows argv 8191-char limit while minimizing subprocess overhead).
   - Use `concurrent.futures.ThreadPoolExecutor(max_workers=os.cpu_count())` to run `gofmt -w` concurrently on file chunks.
   - Track completed files with a thread-safe counter. Print live progressing percentage (e.g. `[150/1500] 10% ... [1500/1500] 100%`).
2. In `.github/scripts/go-format-check.py`:
   - Replace single-threaded `gofmt -l .` with parallel chunked check:
     - Pre-gather all `.go` files using `os.walk` or `stream_directory_files`.
     - Run `gofmt -l` in parallel chunks across `os.cpu_count()` worker threads.
     - Collect unformatted files and print live percentage progress.
     - When auto-formatting, apply `gofmt -w` in parallel chunks.
3. Verify with dry-run and formatting runs.
