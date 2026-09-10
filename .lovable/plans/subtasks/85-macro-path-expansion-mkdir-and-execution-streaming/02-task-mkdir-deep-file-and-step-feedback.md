# Subtask 02: Enhanced Mkdir Engine (Deep Paths, File Creation, Step Feedback)

## 1. Objective
Overhaul `gitmap/cmd/mkdir.go` and `gitmap/cmd/mkdir_test.go` to support cross-platform deep directory creation, file creation (`-f`), multi-OS slash normalization, and progressive step-by-step console logging.

## 2. Target Files
- `gitmap/cmd/mkdir.go`
- `gitmap/cmd/mkdir_test.go`

## 3. Requirements
- Flag parsing: support `-p` / `--parents` anywhere in arguments, support `-f` / `--file` to create a file with its directory hierarchy, and support `-v` / `--verbose`.
- Path sanitization: normalize slashes (`//`, `\\`, `/`, `\`), expand environment variables and temp aliases using `macro.NormalizeTargetPath`.
- Hierarchy creation: create all ancestor directories and report progress for each directory created (`  [DIR]  created: ...`) or existing (`  [DIR]  exists: ...`).
- If `-f` / `--file` is set, create the file via `os.OpenFile(..., os.O_CREATE|os.O_WRONLY, 0644)` and report `  [FILE] created: ...`.
- Idempotency: does not fail if directory already exists; prints clear status.
- Print final success summary: `✔ Created: <path>`.
- All Go functions <= 15 lines with blank lines before every return.
- Unit tests covering:
  - Deep path creation with mixed slashes.
  - File creation (`-f`) with parent directory synthesis.
  - Idempotent re-run.
  - Flag position independence (`mkdir path -p` vs `mkdir -p path`).
  - Error handling (empty args, unwritable target).
