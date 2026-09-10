# Subtask 02: Go Corrupted Dirs Cleaner & Subcommand

## Objective
Author `gitmap/cmd/corrupted_dirs_cleaner.go` and `gitmap/cmd/corrupted_dirs_cmd.go`:
1. Implement `CleanCorruptedDirs(opts CleanOptions) (CleanResult, error)`:
   - Evaluates detected directories.
   - If current process working directory (`os.Getwd()`) is inside or equal to the corrupted directory:
     - Escapes CWD to safe target (`~/.local/bin` or `$HOME`) using `os.Chdir()`.
   - Recovers any binary/config assets into canonical `~/.local/bin` (or `/usr/local/bin` if root).
   - If `opts.IsDryRun`: reports planned actions without modifying filesystem.
   - Removes corrupted directories using `os.RemoveAll()`.
   - Returns structured `CleanResult` with removed paths, recovered files, and escaped paths.
2. Implement CLI subcommand `gitmap clean-corrupted`:
   - Flags: `--dry-run`, `--force`, `--json`.
   - Outputs colorized human-readable report or JSON format.
   - Register in command dispatch table in `gitmap/cmd/rootcore.go`.

## Files Affected
- `gitmap/cmd/corrupted_dirs_cleaner.go`
- `gitmap/cmd/corrupted_dirs_cmd.go`
- `gitmap/cmd/rootcore.go`
