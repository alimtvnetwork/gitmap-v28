# Specification 367.02: GitMap Cat Query & Window Filtering Enhancements

- **Sequence ID:** `367`
- **Subtask ID:** `TASK-02`
- **Feature Slug:** `sync-script-versioning-and-gitmap-skill-enhancement - Task 02`

---

## 1. Objective

Upgrade `gitmap cat` to seamlessly handle combined window and query operations, eliminating the need for PowerShell `Get-Content ... -Tail N | Select-String` or Node/Python line inspection scripts.

---

## 2. Technical Contracts

### 2.1 Flag Parsing Parity in `cat_options.go`
- Add `-t` to `isTailFlag(arg)` alongside `--tail` and `-tail`.
- Add `-q` and `-s` support alongside `--query` and `--filter`.
- Support single line range notation `-L N` where `StartLine = N` and `EndLine = N`.

### 2.2 Combined Window & Query Evaluation in `cat_reader.go`
When both a window flag (`TailCount > 0`, `HeadCount > 0`, or `HasLineRange`) and a query filter (`Query != ""`) are provided:
1. Extract the candidate lines matching the requested window bounds.
2. Filter the windowed lines using `filterLinesByQuery`.
3. If `--tail` is specified with query and the user expects the last N matching lines from the entire file, evaluate accordingly:
   - Pattern A (`gitmap cat <file> -t 10 -q "pattern"`): Filters the entire file for "pattern" and returns the last 10 matching lines with original file line numbers!
   - Pattern B (`gitmap cat <file> -L 100:200 -q "pattern"`): Searches only within lines 100-200 for "pattern".
   - Pattern C (`gitmap cat <file> -n 10 -q "pattern"`): Returns the first 10 matching lines.

### 2.3 `gitmap py` Guard Expansion in `cli/cmdpy/py_cmd.go`
Detect and emit helpful terminal tips for:
- `open(...)` with `.read` or `readlines`
- `os.walk(...)` -> suggest `gitmap find "<pattern>" [dir]` or `gitmap lf [dir]`
- `time.sleep(...)` -> suggest background task scheduling or `gitmap watch`
