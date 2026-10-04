# Subtask 216.1: Terminal Prompt Input Freeze Remediation & Direct Handoff Wrapper

- **Parent Plan:** `pending/216-gitmap-prompting-freeze-and-suggestion-engine-fix.md`
- **Spec Reference:** [01-architecture-spec.md](../../../../02-spec/21-app/216-gitmap-prompting-freeze-and-suggestion-engine-fix/01-architecture-spec.md)
- **Status:** Pending
- **Target Area:** `cli/cmd/console_windows.go`, `cli/constants/constants_cd.go`, `cli/constants/constants_cd_shim.go`, `cli/constants/constants_pathsnippet.go`, `cli/scripts/install.ps1`, `install.ps1`, `cli/cmd/cdops.go`, PowerShell Profile (`$PROFILE`)

---

## 1. Objective

Eradicate the Windows console interactive prompt input freeze during `gitmap cd` picker prompts (`promptCDPick`) and confirmation dialogs (`[y/N]`) by:
1. Eliminating the bug-inducing `SetConsoleCP(65001)` system call from `cli/cmd/console_windows.go` while preserving UTF-8 output (`SetConsoleOutputCP(65001)`).
2. Refactoring all PowerShell wrappers (`gitmap`, `gcd`, `Invoke-GitmapAndSetLocation`, and shims) to eliminate pipeline interception (`| Out-String`), switching entirely to the direct execution `GITMAP_HANDOFF_FILE` IPC pattern.
3. Updating the active user PowerShell `$PROFILE` with the modernized, unpiped wrapper.
4. Verifying `readCDSelection` in `cli/cmd/cdops.go` operates cleanly with live keyboard input from `os.Stdin`.

---

## 2. Root Cause Analysis (RCA)

### 2.1 Win32 Console Input Handle Corruption
- In `cli/cmd/console_windows.go`, `initConsole()` calls `SetConsoleCP(65001)` to set the console input code page to UTF-8.
- On Windows (specifically ConHost and ConPTY), setting input code page to 65001 introduces a well-documented bug where `ReadFile` or `ReadConsoleA` on `STD_INPUT_HANDLE` fails, returns 0 bytes prematurely, or hangs indefinitely waiting for multi-byte sequences that never arrive when Enter/Return is pressed.
- Go's `os.Stdin.Read` and `bufio.Scanner` rely on Win32 `ReadFile` on `GetStdHandle(STD_INPUT_HANDLE)`. Under code page 65001, `readCDSelection` in `cli/cmd/cdops.go` freezes or crashes upon reaching `scanner.Scan()`.

### 2.2 PowerShell Pipeline STDIN Interception
- In `cli/constants/constants_cd.go` (`CDFuncPowerShell`), `constants_cd_shim.go`, `constants_pathsnippet.go`, and install scripts, `cd` and `go` commands are intercepted via:
  ```powershell
  $dest = [string](& $real @args | Out-String)
  ```
- Piping a native external executable into `| Out-String` redirects `StandardOutput` and detaches `StandardInput` from the interactive console keyboard.
- When `promptCDPick` renders on `os.Stderr` and attempts to read from `os.Stdin`, no keystrokes from the terminal keyboard are delivered to the process. The process hangs waiting for input from a detached pipe.

---

## 3. Step-by-Step Implementation Details

### Step 3.1: Eradicate `SetConsoleCP(65001)` in `cli/cmd/console_windows.go`
1. Open `cli/cmd/console_windows.go`.
2. Remove the lazy procedure binding:
   ```go
   // REMOVE:
   setConsoleCP := kernel32.NewProc("SetConsoleCP")
   ```
3. Remove the procedure invocation:
   ```go
   // REMOVE:
   _, _, _ = setConsoleCP.Call(uintptr(consoleCodePageUTF8))
   ```
4. Strictly retain:
   - `setConsoleOutputCP.Call(uintptr(consoleCodePageUTF8))`
   - `enableVTOnHandle(..., consoleStdOutHandle)`
   - `enableVTOnHandle(..., consoleStdErrHandle)`
5. Verify zero remaining occurrences of `setConsoleCP` or `SetConsoleCP` across the codebase.

### Step 3.2: Unify PowerShell Wrapper onto Direct `GITMAP_HANDOFF_FILE` IPC
1. In `cli/constants/constants_cd.go` (`CDFuncPowerShell`):
   - Eradicate the special-case `if ($args.Count -gt 0 -and ($args[0] -eq 'cd' -or $args[0] -eq 'go'))` block that piped via `| Out-String`.
   - Update `function gitmap` to execute all commands via the unified `GITMAP_HANDOFF_FILE` temporary file pattern using direct execution `& $real @args`.
   - Update `function gcd` to delegate cleanly to `gitmap cd @args` or invoke `& $real cd @args` using the same handoff mechanism.
2. In `cli/constants/constants_pathsnippet.go` (`PathSnippetPwshFmt`):
   - Update `Invoke-GitmapAndSetLocation`: eliminate the `| Out-String` branch.
   - Route all invocations through `& $real @GitMapArgs` with `GITMAP_HANDOFF_FILE` exported.
3. In `cli/constants/constants_cd_shim.go` (`PowerShellShimTemplateFmt`):
   - Remove the `| Out-String` branch; unify to direct execution with `$handoff`.
4. In `cli/scripts/install.ps1` and `install.ps1`:
   - Update the generated PowerShell wrapper functions to match the direct handoff structure.

### Step 3.3: Update Active PowerShell `$PROFILE`
1. Locate the active user PowerShell profile at:
   - `$PROFILE` (PowerShell 7+: `$HOME\Documents\PowerShell\Microsoft.PowerShell_profile.ps1`)
   - Windows PowerShell 5.1: `$HOME\Documents\WindowsPowerShell\Microsoft.PowerShell_profile.ps1`
2. Replace any legacy `# gitmap command wrapper v1` or `# gitmap shell wrapper v2` block containing `| Out-String` with the updated direct execution wrapper.
3. Verify that dot-sourcing `$PROFILE` loads without syntax errors.

### Step 3.4: Verify `readCDSelection` in `cli/cmd/cdops.go`
1. Review `readCDSelection`:
   - Verify `scanner := bufio.NewScanner(os.Stdin)` handles empty string (defaults to option 1).
   - Verify invalid integer input displays `invalid selection` and returns error without crashing.
   - Verify `WriteShellHandoff(selectedPath)` is called upon successful resolution.

---

## 4. Verification & Acceptance Criteria

- **No `SetConsoleCP`:** `grep -rn "SetConsoleCP" cli/` yields 0 matches in Go source files.
- **UTF-8 Output Intact:** Multi-byte glyphs (`✓`, `→`, `⚠`, `▸`) and Lipgloss colors continue to render without mojibake.
- **Zero Piped Wrapper Invocations:** Zero occurrences of `| Out-String` in PowerShell wrapper definitions across `cli/constants/`, `cli/scripts/`, `install.ps1`, and `$PROFILE`.
- **Responsive Interactive Prompts:** Executing `gitmap cd <ambiguous>` immediately presents the numerical picker list, responds to keyboard number input and Enter key, saves default choice, writes to `GITMAP_HANDOFF_FILE`, and changes the active PowerShell terminal directory via `Set-Location`.
