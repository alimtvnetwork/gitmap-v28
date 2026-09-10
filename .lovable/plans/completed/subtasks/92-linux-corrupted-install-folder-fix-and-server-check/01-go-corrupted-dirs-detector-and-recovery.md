# Subtask 01: Go Corrupted Dirs Detector & Asset Recovery Engine

## Objective
Author `gitmap/cmd/corrupted_dirs_constants.go`, `gitmap/cmd/corrupted_dirs_detector.go`, and `gitmap/cmd/corrupted_dirs_recovery.go`:
1. Define constants for detection patterns:
   - `CorruptedPatternQuickInstaller = "quick installer"`
   - `CorruptedPatternInstaller = "gitmap installer"`
   - `CorruptedPatternDefault = "Default:"`
   - `CorruptedPatternPrompt = "Choose install folder"`
   - `CorruptedPatternPath = "Install path"`
   - `CorruptedTildeDir = "~"`
   - `AnsiEscapePrefix = "\x1b["`
2. Implement `DetectCorruptedDirs() ([]CorruptedDirInfo, error)`:
   - Scans candidates: `os.UserHomeDir()`, `os.Getwd()`, `os.TempDir()`, `/tmp`.
   - `isCorruptedDirName(name string) bool`:
     - Checks literal `"~"`
     - Checks ANSI escape sequences (`\x1b[`, `\033[`)
     - Checks newlines/carriage returns (`\n`, `\r`)
     - Checks installer prompt leakage
   - `isProtectedPath(path string, homeDir string) bool`:
     - Blacklists `/`, `$HOME`, `/usr`, `/usr/local`, `/usr/local/bin`, `/bin`, `/tmp`, `.`
     - Requires `filepath.Base(path) == "~"` and `path != homeDir` for tilde
3. Implement `RecoverCorruptedDirAssets(info CorruptedDirInfo, targetDir string) error`:
   - Inspects corrupted directory for `gitmap`, `gitmap-cli`, or `.gitmap` assets.
   - If found and `targetDir/gitmap` does not exist, copies with `0755` permissions to `targetDir` (`~/.local/bin`).
   - If current running executable (`os.Executable()`) is inside corrupted dir, safely copies out before unlinking.

## Files Affected
- `gitmap/cmd/corrupted_dirs_constants.go`
- `gitmap/cmd/corrupted_dirs_detector.go`
- `gitmap/cmd/corrupted_dirs_recovery.go`
