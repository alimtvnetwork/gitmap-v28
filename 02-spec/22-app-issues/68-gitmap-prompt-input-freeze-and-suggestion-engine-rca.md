# App Issue 68: GitMap Interactive Prompt Input Freeze and Suggestion Engine Collapse RCA

**Issue ID:** 68  
**Date:** 2026-10-05  
**Status:** Resolved  
**Affected Subsystems:** `cli/cmd` (`console_windows.go`, `rootsuggest.go`, `rootsuggest_calc.go`, `root_cobra_completion.go`), `cli/constants` (`constants_cd.go`, `constants_cd_shim.go`, `constants_pathsnippet.go`), `scripts` (`install.ps1`, `cli/scripts/install.ps1`)  
**Reference Commits:** `cli - fix prompt input freeze and suggestion engine`  

---

## 1. Symptom

Users and automated workflows reported three interrelated terminal execution and usability failures:

1. **Terminal Keystroke Freeze on Interactive Prompts:**
   - When GitMap prompts for user confirmation (`[y/N]`, `[Y/n]`) in commands such as `gitmap secret`, `gitmap clean`, or repository reset operations, user keystrokes in the console fail to register. The cursor appears unresponsive or deadlocked.
   - During ambiguous repository navigation via `gitmap cd <query>` (or `gcd <query>`), when multiple repositories match and GitMap invokes `promptCDPick()`, the selection list either prints and immediately hangs without accepting number/arrow selections, or instantly terminates with an unexpected `io.EOF`, picking the default option or crashing.
   - Master passphrase prompts (e.g., encryption vault unlocking) exhibit complete input capture failure on Windows terminal hosts.

2. **Typo Suggestions Completely Missing for Core Commands:**
   - When entering slight misspellings of foundational GitMap commands—such as `gitmap instlal` (intended `install`), `gitmap seach` (intended `search`), `gitmap fnd` (intended `find`), `gitmap commti` (intended `commit`), or `gitmap ap` (intended `apps`)—GitMap failed to suggest the intended command:
     ```text
     gitmap: unknown command "instlal"
       Run 'gitmap help' for usage.
     ```
   - Instead of offering the obvious Levenshtein distance-1 match (`install`), the CLI emitted either no suggestions or unrelated keyword matches, eroding developer experience.

3. **Shell Tab Autocompletion Truncation (9 Commands vs 512 Available):**
   - When invoking shell completion via `gitmap __complete ""` or pressing `<Tab>` at the root `gitmap ` prompt in PowerShell, Bash, or Zsh, only 9 commands were presented as completions.
   - GitMap provides over 512 commands, flags, and aliases across scanning, SSH fleet orchestration, release workflows, and AI integrations (defined in `completion.AllCommands()`), but 503 commands were missing from shell autocompletion.

---

## 2. Root Cause Analysis (RCA)

A rigorous forensic audit of the codebase identified three root causes:

### 2.1 Win32 Input Console Code Page Corruption (`cli/cmd/console_windows.go:35`)

In `cli/cmd/console_windows.go`, `initConsole()` initializes console parameters on Windows to ensure Unicode glyphs (`✓`, `→`, `⚠`, `▸`) render properly:

```go
26: func initConsole() {
27: 	kernel32 := syscall.NewLazyDLL("kernel32.dll")
28: 	setConsoleOutputCP := kernel32.NewProc("SetConsoleOutputCP")
29: 	setConsoleCP := kernel32.NewProc("SetConsoleCP")
30: 	getConsoleMode := kernel32.NewProc("GetConsoleMode")
31: 	setConsoleMode := kernel32.NewProc("SetConsoleMode")
32: 	getStdHandle := kernel32.NewProc("GetStdHandle")
33: 
34: 	_, _, _ = setConsoleOutputCP.Call(uintptr(consoleCodePageUTF8))
35: 	_, _, _ = setConsoleCP.Call(uintptr(consoleCodePageUTF8))
36: 
37: 	enableVTOnHandle(getStdHandle, getConsoleMode, setConsoleMode, consoleStdOutHandle)
38: 	enableVTOnHandle(getStdHandle, getConsoleMode, setConsoleMode, consoleStdErrHandle)
39: }
```

- **Forensic Line:** `cli/cmd/console_windows.go:35` calls `SetConsoleCP(65001)`.
- **Mechanism:** `SetConsoleOutputCP(65001)` correctly configures console *output* rendering. However, calling `SetConsoleCP(65001)` sets the console *input* code page (`STD_INPUT_HANDLE`) to UTF-8. In Windows conhost (`conhost.exe`) and ConPTY, setting input code page 65001 triggers an established Win32 kernel subsystem bug: standard console input reads (`ReadFile` / `ReadConsoleA` / `ReadConsoleW`) fail to translate multi-byte sequences properly.
- In Go, `bufio.Reader` and `os.Stdin.Read()` either:
  1. Immediately return 0 bytes read, which Go interprets as `io.EOF`, instantly skipping interactive prompts; or
  2. Deadlock inside conhost's input buffer state machine, dropping all keyboard scan codes and freezing the console prompt.

### 2.2 PowerShell Wrapper Pipe Redirection & TTY Detachment (`cli/constants/constants_cd.go:201, 237`)

In `cli/constants/constants_cd.go` (and identical implementations in `constants_cd_shim.go`, `constants_pathsnippet.go`, and installer scripts), the PowerShell wrapper function intercepts `cd` and `go` invocations:

```powershell
200:   $env:GITMAP_COMMAND_WRAPPER = "1"
201:   $dest = [string](& $real cd @args | Out-String)
202:   if ($LASTEXITCODE -ne 0) {
...
236:     $env:GITMAP_COMMAND_WRAPPER = "1"
237:     $dest = [string](& $real @args | Out-String)
238:     if ($LASTEXITCODE -ne 0) {
```

- **Forensic Lines:** `cli/constants/constants_cd.go:201` and `cli/constants/constants_cd.go:237` invoke `& $real cd @args | Out-String` and `& $real @args | Out-String`.
- **Mechanism:** In PowerShell, piping a native binary through `| Out-String` captures standard output (`STDOUT`) into an in-memory string stream and detaches standard input (`STDIN`) from the interactive console keyboard.
- When `gitmap cd` detects multiple ambiguous matches and invokes `promptCDPick()`, GitMap attempts to write menu prompts to stdout and read user selection from stdin. Because stdout is buffered in `Out-String` and stdin is disconnected from the interactive console host, the user cannot see the prompt in real time, and the input stream receives `io.EOF` or deadlocks on the pipeline buffer.
- In contrast, general `gitmap` invocations in lines 250–260 of `constants_cd.go` correctly use foreground execution `& $real @args` combined with the `GITMAP_HANDOFF_FILE` IPC mechanism. The `cd` and `go` wrapper branches bypassed this robust IPC mechanism in favor of `| Out-String`.

### 2.3 Suggestion Engine Disconnect and Cobra Completion Pruning

The suggestion and completion failure stems from two architectural gaps:

#### A. Typo Candidate Pool Isolation (`cli/cmd/rootsuggest.go:14-30`)

In `cli/cmd/rootsuggest.go`, `primaryTopCommands` maintained a manual, hardcoded list of 101 commands:

```go
14: var primaryTopCommands = []string{
15: 	"scan", "clone", "clone-only-missing", "com", "create", "clone-sync", "pull", "push", "pull-all",
...
29: 	"gitignore", "agm", "ta", "nodes", "ping", "ports", "port",
30: }
```

- Critical commands—including `install`, `uninstall`, `apps`, `commit`, `search`, `find`, `login`, `setup`, `cpf`, `cpb`, and `cpr`—were completely omitted from `primaryTopCommands`.
- In `cli/cmd/rootsuggest_calc.go:15-26`, `collectTopCommandCandidates()` only aggregated `primaryTopCommands`, `"run"`, `"run-until"`, and dynamic user macros. It completely failed to reference `completion.AllCommands()`, which contains the authoritative list of all 512 commands and aliases.
- As a result, the Levenshtein distance ranking in `rankCandidateCommands()` never evaluated the omitted commands, returning no matches for `gitmap instlal`, `gitmap seach`, or `gitmap fnd`.

#### B. Cobra Command Runnable Pruning (`cli/cmd/root_cobra_completion.go:214-222`)

In `cli/cmd/root_cobra_completion.go`, `populateRemainingCommands(root)` dynamically added commands from `completion.AllCommands()` to the root Cobra command:

```go
213: 	for _, cmdName := range completion.AllCommands() {
214: 		if existing[cmdName] || strings.TrimSpace(cmdName) == "" {
215: 			continue
216: 		}
217: 		desc := resolveCommandHelpShort(cmdName)
218: 		root.AddCommand(&cobra.Command{
219: 			Use:   cmdName,
220: 			Short: desc,
221: 		})
222: 		existing[cmdName] = true
223: 	}
```

- **Forensic Lines:** `cli/cmd/root_cobra_completion.go:218-221` creates `cobra.Command` instances with `Run: nil` and `RunE: nil`.
- **Mechanism:** In Cobra's architecture, `cmd.Runnable()` returns `c.Run != nil || c.RunE != nil`. When both are nil, `c.Runnable()` evaluates to `false`.
- When Cobra generates shell completion scripts or handles `gitmap __complete ""`, it checks `cmd.IsAvailableCommand()`. For leaf commands without subcommands, `IsAvailableCommand()` requires `cmd.Runnable() == true`. Because 503 commands lacked execution handlers, Cobra pruned them as non-runnable group headers, reducing the completion output to only 9 statically initialized commands.

---

## 3. Resolution

A three-part surgical remediation was engineered across the codebase:

### 3.1 Eradicate `SetConsoleCP(65001)` in `cli/cmd/console_windows.go`

Removed line 35 from `initConsole()`, leaving `SetConsoleOutputCP(65001)` and Virtual Terminal Processing intact:

```diff
--- a/cli/cmd/console_windows.go
+++ b/cli/cmd/console_windows.go
@@ -32,7 +32,6 @@ func initConsole() {
 	getStdHandle := kernel32.NewProc("GetStdHandle")
 
 	_, _, _ = setConsoleOutputCP.Call(uintptr(consoleCodePageUTF8))
-	_, _, _ = setConsoleCP.Call(uintptr(consoleCodePageUTF8))
 
 	enableVTOnHandle(getStdHandle, getConsoleMode, setConsoleMode, consoleStdOutHandle)
 	enableVTOnHandle(getStdHandle, getConsoleMode, setConsoleMode, consoleStdErrHandle)
```

- Output glyphs and ANSI color formatting continue to render via code page 65001 on stdout and stderr.
- The console input code page remains at system default, allowing `ReadConsole` and `os.Stdin` to receive interactive keystrokes without Win32 conhost deadlocks or spurious EOF signals.

### 3.2 Eliminate Pipe Redirection in PowerShell Wrappers

Updated `CDFuncPowerShell` in `cli/constants/constants_cd.go`, `constants_cd_shim.go`, `constants_pathsnippet.go`, and installer scripts (`scripts/install.ps1`, `cli/scripts/install.ps1`, and `$PROFILE`). Replaced `| Out-String` with foreground execution and `GITMAP_HANDOFF_FILE` IPC:

```powershell
function gcd {
  $real = Get-GitmapCommand
  if (-not $real) {
    Write-Error "gitmap executable not found"
    return
  }

  $handoff = [System.IO.Path]::Combine([System.IO.Path]::GetTempPath(), "gitmap-handoff-$([System.Guid]::NewGuid().ToString('N')).txt")
  try {
    $env:GITMAP_HANDOFF_FILE = $handoff
    $env:GITMAP_WRAPPER = "1"
    $env:GITMAP_COMMAND_WRAPPER = "1"
    & $real cd @args
    if ((Test-Path -LiteralPath $handoff) -and ((Get-Item -LiteralPath $handoff).Length -gt 0)) {
      $dest = [string](Get-Content -LiteralPath $handoff -Raw)
      $dest = $dest.Trim()
      if ($dest -and (Test-Path -LiteralPath ([string]$dest))) {
        Set-Location -LiteralPath ([string]$dest)
      }
    }
  } finally {
    Remove-Item -LiteralPath $handoff -ErrorAction SilentlyContinue
  }
}
```

- The executable runs directly in the terminal foreground with unrestricted access to `STDIN` and `STDOUT`.
- When interactive disambiguation (`promptCDPick`) or confirmations are displayed, users can type and select options seamlessly.
- Target directory navigation is communicated cleanly via the handoff file upon process exit.

### 3.3 Overhaul Suggestion Engine and Cobra Completion

#### A. Comprehensive Candidate Merging in `cli/cmd/rootsuggest_calc.go`

Refactored `collectTopCommandCandidates()` to merge `completion.AllCommands()`, `primaryTopCommands`, `"run"`, `"run-until"`, and user macros with deduplication:

```go
func collectTopCommandCandidates() []string {
	seen := make(map[string]bool)
	var candidates []string

	add := func(cmd string) {
		cmd = strings.TrimSpace(cmd)
		if cmd != "" && !seen[cmd] {
			seen[cmd] = true
			candidates = append(candidates, cmd)
		}
	}

	for _, c := range completion.AllCommands() {
		add(c)
	}
	for _, c := range primaryTopCommands {
		add(c)
	}
	add("run")
	add("run-until")

	macroList := macro.ListMacros()
	if macroList.IsSuccess() {
		for _, m := range macroList.Data {
			add(m.Name)
		}
	}

	return candidates
}
```

- Added `install`, `uninstall`, `apps`, `commit`, `search`, `find`, `login`, `setup` to `primaryTopCommands` in `cli/cmd/rootsuggest.go`.
- All 512 commands and aliases are now candidates for Levenshtein distance matching.

#### B. Enable Runnable State on Completion Commands in `cli/cmd/root_cobra_completion.go`

Added a stub `Run` handler to dynamically populated commands:

```go
func populateRemainingCommands(root *cobra.Command) {
	existing := make(map[string]bool)
	for _, c := range root.Commands() {
		existing[c.Name()] = true
		for _, alias := range c.Aliases {
			existing[alias] = true
		}
	}

	for _, cmdName := range completion.AllCommands() {
		if existing[cmdName] || strings.TrimSpace(cmdName) == "" {
			continue
		}
		desc := resolveCommandHelpShort(cmdName)
		root.AddCommand(&cobra.Command{
			Use:   cmdName,
			Short: desc,
			Run:   func(cmd *cobra.Command, args []string) {},
		})
		existing[cmdName] = true
	}
}
```

- Setting `Run: func(...) {}` ensures `cmd.Runnable() == true` and `cmd.IsAvailableCommand() == true`.
- All 512 commands are now emitted during shell tab completion.
- Configured dynamic argument completion via `ValidArgsFunction` for `cd`, `clone`, `apps`, and `install`.

---

## 4. Prevention & Learnings

To permanently prevent regressions across all supported platforms:

1. **Static Linter Gate for Console APIs:**
   - Incorporate a rule in `03-ai-scripts/05-guideline-autofixer.py` and `03-ai-scripts/check-forbidden-strings.py` that flags any occurrence of `SetConsoleCP` across the codebase. Only `SetConsoleOutputCP` is permitted.

2. **Shell Wrapper IPC Standard:**
   - All shell wrapper functions (PowerShell, Bash, Zsh, Fish) must never pipe interactive commands (`| Out-String`, `| cat`, `| tee`).
   - File-based IPC (`GITMAP_HANDOFF_FILE`) is the mandatory pattern for passing state between GitMap and parent shells.

3. **Suggestion Engine Unit Tests:**
   - Added unit tests in `cli/cmd/rootsuggest_test.go` asserting that `suggestTopLevelCommands()` correctly matches common misspellings:
     - `instlal` $\to$ `install`
     - `seach` $\to$ `search`
     - `fnd` $\to$ `find`
     - `commti` $\to$ `commit`
     - `ap` $\to$ `apps`
   - Asserted that `collectTopCommandCandidates()` contains at least 500 unique commands.

4. **Cobra Completion Invariant Tests:**
   - Added automated tests in `cli/cmd/root_cobra_completion_test.go` verifying that `GetRootCompletionCmd().Commands()` contains $\ge 500$ commands, and that all leaf commands satisfy `cmd.Runnable() == true`.

---

## 5. Verification Evidence

- **Console Input Verification:** Verified that `bufio.NewReader(os.Stdin).ReadString('\n')` and interactive `[y/N]` prompts accept keyboard input without hanging or receiving premature EOF on Windows ConHost and Windows Terminal.
- **PowerShell Wrapper Verification:** Tested `gitmap cd` with ambiguous query; verified selection menu renders interactively in real time, accepts numerical input, and successfully navigates via handoff file.
- **Suggestion Engine Verification:** Ran `suggestTopLevelCommands` against `instlal`, `seach`, `fnd`, `commti`, `ap`; verified each resolves to its expected primary command.
- **Completion Verification:** Executed `gitmap __complete ""` in PowerShell and Bash; verified output contains 512 commands and aliases.
