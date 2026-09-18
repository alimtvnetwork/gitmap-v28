# RCA: Nested Ifs, SSH Help Exit Panic, Absolute Paths, and Gofmt Drift

## 1. Symptom

GitHub Actions CI run `#35354190330` for commit `b9d9023` on `main` failed across four distinct validation gates:

1. **Boolean & Enum Linter (`check-enum-and-boolean.py`)**:
   - `cli/cmdagy/agy_fix_pipeline_inject.go:35`: Nested 'if' detected (depth 2):
     `if logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err == nil`
   - `cli/cmdssh/ssh_exec_command.go:74`: Nested 'if' detected (depth 2):
     `if len(fields) > 1` inside `if len(args) == 1`

2. **Full Suite Guard (`go test ./...` in `cli/cmdssh`)**:
   - Panic in `TestDispatchPrimarySSH_MatchedSubcommands`:
     ```text
     panic: unexpected call to os.Exit(0) during test [recovered, repanicked]
     os.Exit(0x0)
     github.com/alimtvnetwork/gitmap-v28/cli/cliexit.Exit(0x0)
     github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.parseSEFlags(...)
     github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.runSSHExec(...)
     github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.dispatchPrimarySSH(...)
     github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.TestDispatchPrimarySSH_MatchedSubcommands(...)
     ```

3. **Relative Path Linter (`check-relative-paths.py`)**:
   - 12 hardcoded absolute repo paths (`D:\work\gitmap` and `d:\work\gitmap`) across:
     - `.ai-memory/plans/completed/210-ssh-multi-command-discovery-and-agy-terminal-verification.md:64`
     - `.ai-memory/plans/completed/211-ssh-multi-command-discovery-and-agy-terminal-verification.md:65`
     - `.ai-memory/plans/completed/206-pipeline-errors-agy-fix-and-multi-project-batch-hardening.md:54, 55`
     - `.ai-memory/plans/completed/207-pipeline-errors-agy-fix-comprehensive-verification-and-audit.md:5, 22, 27`
     - `.lovable/plans/completed/206-pipeline-errors-agy-fix-and-multi-project-batch-hardening.md:54, 55`
     - `.lovable/plans/completed/207-pipeline-errors-agy-fix-comprehensive-verification-and-audit.md:5, 22, 27`

4. **Lint Script Unit Tests (`.github/scripts/tests/test_ci_scripts.py`)**:
   - `test_gofmt_check_clean_repo` failed because recent Go changes in `cmdssh` were not formatted with `gofmt`.

## 2. Root Cause

1. **Nested Ifs**: In `cli/cmdagy/agy_fix_pipeline_inject.go`, an `os.OpenFile` check was enclosed directly within `if len(repoDir) > 0`, producing a nesting depth of 2. In `cli/cmdssh/ssh_exec_command.go:extractShellCommandArgs`, `if len(fields) > 1` was enclosed directly inside `if len(args) == 1`, violating the repository maximum nesting depth rule (depth <= 1).
2. **Process Exit in Unit Tests**: `parseSEFlags` in `cli/cmdssh/sshexec.go` directly invoked `cliexit.Exit(0)` when `hasHelpFlag(args)` was true. In `ssh_test.go:TestDispatchPrimarySSH_MatchedSubcommands`, `dispatchPrimarySSH` is invoked with `[]string{"--help"}` across all matched subcommands to verify sub-routing. Calling `cliexit.Exit(0)` triggered `os.Exit(0)`, which the Go standard test runner recovers and panics on.
3. **Absolute Windows Path Leakage**: Verbatim PowerShell terminal output (`PS D:\work\gitmap> ...`) and absolute file paths (`d:\work\gitmap\bin\gitmap.exe`) from local testing sessions were copied into completed plan documentation files in `.ai-memory/plans/completed/` and `.lovable/plans/completed/`.
4. **Gofmt Drift**: Files modified in Plans 212–214 under `cli/cmdssh/` were committed without running `gofmt -w`, resulting in minor whitespace and alignment drift detected by `.github/scripts/go-format-check.py`.

## 3. Resolution

1. **Flattened Conditionals**:
   - Extracted `attachInjectionLog` helper function in `cli/cmdagy/agy_fix_pipeline_inject.go` with early guard returns (`if len(repoDir) == 0`, `if err != nil`), keeping both functions under 15 lines and nesting depth 1.
   - Refactored `extractShellCommandArgs` in `cli/cmdssh/ssh_exec_command.go` using guard clauses (`if len(args) > 1`, `if len(args) == 0`), eliminating nesting.
2. **Graceful Help Flag Handling**:
   - Added `IsShowHelp bool` field to `seOptions` in `cli/cmdssh/sshexec.go`. `parseSEFlags` prints help and returns `seOptions{IsShowHelp: true}`. `runSSHExec` checks `opts.IsShowHelp` and returns `nil` cleanly.
   - Refactored `showJoinHelpAndExit` in `cli/cmdssh/sshjoin_cmd.go` to return `nil` without calling `cliexit.Exit(0)`.
   - Updated `checkHelp` and `checkSSHHelp` in `cli/cmdssh/helpers.go` to return `bool`, and updated `runSSH` in `cli/cmdssh/ssh.go` to return `nil` when help was displayed.
   - Removed unused `cliexit` imports in `cli/cmdssh/helpers.go` and `cli/cmdssh/sshjoin_cmd.go`.
3. **Path Sanitization**:
   - Sanitized all 12 absolute path occurrences in `.ai-memory/plans/completed/` and `.lovable/plans/completed/` to relative paths (`PS gitmap>`, `./bin/gitmap.exe`, `%LOCALAPPDATA%`).
4. **Go Formatting**:
   - Executed `gofmt -w cli/`, resolving all formatting drift across `cmdssh`.

## 4. Prevention & Learnings

- **Never call `cliexit.Exit(0)` or `os.Exit` inside subcommand handlers**: Command routing handlers must return `nil` or an error, allowing the top-level CLI runner to control process termination and allowing unit tests to execute cleanly without panicking.
- **Maintain Nesting Depth <= 1**: Always use early guard returns or decompose into small single-purpose helper functions (each <= 15 lines).
- **Sanitize Plan Evidence Before Commit**: When documenting command output in plans, always strip absolute drive letters and user paths (`D:\work\gitmap`, `C:\Users\...`).
- **Always Run `gofmt -w cli/` and `check-relative-paths.py` Pre-Commit**: Verify all linters pass locally before pushing.
