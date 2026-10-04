# RCA-54: Remote Windows Node Clone Failure (`'bash' is not recognized`) and Fleet Output Degradation

**Status:** Resolved
**Date:** 2026-09-30
**Affected Subsystem:** `cli/cmdnodes/`, Fleet Remote Execution (`nodes clone`, `nodes cfr`, `nodes cfrp`)
**Associated Spec:** [02-spec/21-app/191-nodes-clone-except-self-windows-runner-and-path.md](../21-app/191-nodes-clone-except-self-windows-runner-and-path.md)
**Visual Proof:** `assets/screenshots/nodes-clone-fleet-ui-and-windows-bash-failure.png`

---

## 1. Reproduction Steps

1. Enroll Windows node `w4` (`node-w4`) into the SSH fleet database.
2. Run fleet clone against any Git repository:
   ```powershell
   gitmap nodes clone https://github.com/alimtvnetwork/awansoft-v10
   ```
3. Observe output:
   - Local clone completes on the master host.
   - Node `w4` fails immediately (in 219ms) with error:
     ```text
     command execution failed: Process exited with status 1 (output: 'bash' is not recognized as an internal or external command,
     operable program or batch file.
     )
     ```
   - Terminal table wraps multiple raw lines of error text across the table, destroying column alignment.

---

## 2. Root Cause Analysis

### 2.1 Stale or Default OS Classification
- In `SSHConnection` SQLite table, node `w4` had `OS = 'linux'`.
- OpenSSH server on Windows launches the default login shell configured for the user (Windows `cmd.exe` or `powershell.exe`).
- When dispatching, `runRemoteExecOverSSH` evaluated:
  ```go
  isWin := strings.EqualFold(conn.OS, "windows") || strings.EqualFold(conn.OS, "win")
  shell := "bash"
  if isWin {
      shell = "ps"
  }
  ```
- Because `conn.OS` was `"linux"`, it forced `shell = "bash"`. Over SSH, `wrapBashCommand` sent `bash -c ...` to Windows `cmd.exe`, which failed with status 1 because `bash` was not installed.

### 2.2 Lack of Self-Healing Fallback
- `runRemoteExecOverSSH` did not inspect stderr or return values to check if the shell binary was missing.
- When `bash` was missing, it directly failed the node instead of attempting a PowerShell fallback or checking if the remote host is a Windows machine.

### 2.3 Terminal Output Formatting Fragility
- `formatDetails` in `nodes_clone_table.go` passed `r.Error` directly to `fmt.Fprintf` without stripping embedded newlines.
- When stderr contained line breaks, the table row collapsed across multiple terminal lines, corrupting visual presentation.

---

## 3. Code Fix

1. **Auto-Fallback to PowerShell:**
   - In `runRemoteExecOverSSH`, if execution with `bash` fails with `'bash' is not recognized` or `bash: command not found`:
     - Dynamically fallback to `shell = "ps"`.
     - Rebuild the command using Windows syntax (`Set-Location`).
     - Retry over SSH.
     - On success, update `SSHConnection.OS = 'windows'` in the database via `db.UpdateConnectionOS`.
2. **Work Directory Pre-Navigation:**
   - Both Windows (`Set-Location "<workDir>"`) and Unix (`cd "<workDir>"`) commands now prepend directory navigation prior to invoking `gitmap clone`, supporting both default work directories and user-specified `[target_path]` arguments.
3. **`except-self` Support:**
   - Introduced `--except-self` flag and `except-self` positional keyword, skipping local execution and reporting `○ skipped (except-self)` in the results table.
4. **Table Sanitization:**
   - Sanitized `formatDetails` to strip embedded newlines, extract core error descriptions, and maintain rigid column boundaries.

---

## 4. Prevention & Quality Guardrails

1. **Heuristic Windows Identification:**
   - Check `conn.OS == "windows"`, `conn.OS == "win"`, `conn.OSGroup == "windows"`, or `conn.Username == "Administrator"` (case-insensitive).
2. **Unit Test Coverage:**
   - Added test cases in `cli/cmdnodes/nodes_clone_test.go` verifying that `'bash' is not recognized` triggers PowerShell fallback, `--except-self` sets `opts.IsSkipLocal = true`, and custom paths are threaded into `opts.TargetDir`.
