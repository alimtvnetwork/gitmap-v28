# Subtask Plan 03: Dev Tools Remote & Local Uninstall Subsystem with --node Delegation

- **Spec Reference:** [02-spec/21-app/224-ubuntu-cursor-dock-pin-projects-os-update-and-remote-uninstall/02-component-and-cli-spec.md](../../../../02-spec/21-app/224-ubuntu-cursor-dock-pin-projects-os-update-and-remote-uninstall/02-component-and-cli-spec.md)
- **Status:** Queued
- **Target Area:** `cli/cmd/uninstall.go`, `cli/cmdinstall/uninstall_tools.go`, `cli/cmdcursor/cursor_install.go`, `cli/cmdcursor/cursor_cmd.go`, `cli/constants/constants_install.go`

---

## 1. Objective

Architect and implement the expanded **Dev Tools Uninstallation Subsystem** in GitMap CLI:
1. Support uninstallation of the 5 targeted developer tools: PowerShell (`pwsh` / `powershell`), Antigravity Manager (`agm` / `antigravity-manager`), Vim (`vim`), Visual Studio Code (`vscode` / `code`), and Cursor IDE (`cursor`).
2. Implement `--node <alias>` remote fleet delegation using `cmdssh.RunSSHExec`, enabling developers to trigger tool removal across remote nodes (e.g. `u1`) from their central terminal.
3. Implement `--dry-run` (`-d`, `-n`) simulation mode to print precise uninstallation plans and command lines without deleting packages or files, satisfying the non-destructive requirement.
4. Support `--force` (`-f`, `-y`) for non-interactive execution and `--purge` (`-a`, `--all`) for full removal of configuration directories, caches, shortcuts, and dock favorites.
5. Guarantee strict safety guards preventing accidental modification or deletion of working repositories (`d:\work`, `/home/a/git-work`).

---

## 2. Implementation Details

### Step 1: Flag Binding & Data Structures in `cli/cmd/uninstall.go`
1. Extend `uninstallFlags` struct in `cli/cmd/uninstall.go`:
   ```go
   type uninstallFlags struct {
       isDryRun        bool
       isForce         bool
       isPurge         bool
       isPurgeWebView2 bool
       backupPath      string
       targetNode      string
   }
   ```
2. Update flag parsing functions:
   - In `bindCoreUninstallFlags`:
     * Bind `--node` / `-node` to `f.targetNode` (string).
     * Bind `--dry-run` and `-d` to `f.isDryRun` (boolean).
     * Retain `-n` as boolean dry-run alias only when followed by no value or non-node string, ensuring unambiguous token parsing.
     * Bind `--force`, `-f`, `-y`, `--yes` to `f.isForce`.
     * Bind `--purge`, `-a`, `--all` to `f.isPurge`.
3. Update `hasPositionalToolArg(args)` to properly skip `--node <value>` pairs when determining the presence of positional tool arguments.

### Step 2: Remote Fleet Delegation Interceptor
1. In `executeUninstallFlow(args []string)`:
   - Check if `flags.targetNode != ""`.
   - If a target node is specified:
     * Construct sanitized remote command line:
       ```go
       remoteCmd := fmt.Sprintf("gitmap uninstall %s", tool)
       if flags.isDryRun {
           remoteCmd += " --dry-run"
       }
       if flags.isForce {
           remoteCmd += " --force"
       }
       if flags.isPurge {
           remoteCmd += " --purge"
       }
       ```
     * Dispatch to fleet execution:
       ```go
       fmt.Printf("● Delegating uninstallation of '%s' to remote node: %s\n", tool, flags.targetNode)
       return cmdssh.RunSSHExec([]string{flags.targetNode, remoteCmd})
       ```
     * Return early upon successful dispatch.

### Step 3: Tool Dispatch & Specialized Handlers
1. In `processToolUninstall(tool string, flags *uninstallFlags)`:
   - Add specialized dispatcher `dispatchDevToolUninstall(tool, canonical, flags)`:
     ```go
     if handled, err := dispatchDevToolUninstall(tool, canonical, flags); handled {
         return err
     }
     ```
2. Implement `cli/cmdinstall/uninstall_tools.go`:
   - **`UninstallPwsh(flags)`:**
     * Linux:
       - If dry-run: print planned `apt-get remove powershell` (or `purge`), unlinking `/usr/local/bin/pwsh`, removing `~/.config/powershell`.
       - If live: execute `sudo -n apt-get remove/purge -y powershell`, delete wrappers, clean config if purge requested.
     * Windows:
       - If dry-run: print `winget uninstall Microsoft.PowerShell` or `choco uninstall powershell-core`.
       - If live: execute package manager uninstall.
   - **`UninstallAgm(flags)`:**
     * Linux:
       - If dry-run: print removal of `/usr/local/bin/agm`, `~/.local/bin/agm`, desktop entries, `~/.config/antigravity`.
       - If live: remove wrappers, remove desktop entries, purge `~/.config/antigravity` if purge requested.
     * Windows:
       - Delegate to `cmdinstall.RunAGMUninstall(flags.isPurge, flags.isForce, flags.isDryRun)`.
   - **`UninstallVim(flags)`:**
     * Linux:
       - If dry-run: print planned `apt-get remove/purge -y vim`, removing `~/.vim`, `~/.vimrc`.
       - If live: execute `sudo -n apt-get remove/purge -y vim`, clean configs if purge requested.
     * Windows:
       - Winget/Chocolatey uninstaller invocation.
   - **`UninstallVSCode(flags)`:**
     * Linux:
       - If dry-run: print planned `apt-get remove/purge -y code`, removing desktop entries, `/usr/local/bin/code`, `~/.config/Code`.
       - If live: execute `sudo -n apt-get remove/purge -y code`, remove wrappers and desktop launchers, purge config if requested.
     * Windows:
       - Winget `Microsoft.VisualStudioCode` uninstallation.
   - **`UninstallCursor(flags)`:**
     * Linux:
       - If dry-run: print planned removal of `/opt/cursor`, `~/.local/share/cursor`, `/usr/local/bin/cursor`, `~/.local/bin/cursor`, desktop entries, branding icons, unpinning from GNOME dock favorites (`favorite-apps`), and purging `~/.config/Cursor`.
       - If live: execute file deletions, call DBus helper to unpin `'cursor.desktop'` from `org.gnome.shell favorite-apps`, purge config if requested.
     * Windows:
       - Winget `Anysphere.Cursor` removal or AppData directory cleanup.

### Step 4: Dry-Run Simulation Engine
1. Provide a unified `renderUninstallDryRun(tool string, plan UninstallPlan)` formatter in `cli/cmdinstall/`:
   ```text
   ● [DryRun] Simulated Uninstallation for '<tool>' (OS: <runtime.GOOS>):
     - Package Removal: <command or N/A>
     - Binary Wrappers: <paths to unlink>
     - Desktop Launchers: <paths to remove>
     - System Icons: <paths to remove>
     - Desktop Dock: <dock modifications>
     - Configuration: <paths to purge if --purge>
   ✔ Dry-run completed. No packages or files were modified.
   ```
2. Ensure exit code is 0 on dry run.

### Step 5: Work Directory Protection Guard
1. Add validation check `isProtectedWorkDirectory(targetPath string) bool`:
   - Checks if target path is equal to, or an ancestor/descendant of:
     * Current working directory (`.`)
     * `d:\work` or `/home/a/git-work`
     * Active Git repositories containing `.git`
2. If a protected path is detected, abort immediately with an error and refuse execution.

---

## 3. Verification & Testing Protocol

1. **Local Command Registration Test:**
   - Execute: `gitmap uninstall --help`
   - Verify flags `--node`, `--dry-run`, `--force`, `--purge` are displayed.
2. **Local Dry-Run Testing across all 5 Tools:**
   - `gitmap uninstall pwsh --dry-run`
   - `gitmap uninstall agm --dry-run`
   - `gitmap uninstall vim --dry-run`
   - `gitmap uninstall vscode --dry-run`
   - `gitmap uninstall cursor --dry-run`
   - Verify zero files deleted and 0 exit code.
3. **Remote Dry-Run Testing on Node `u1`:**
   - `gitmap uninstall cursor --node u1 --dry-run`
   - Verify output matches expected remote uninstallation plan.
4. **Unit Tests:**
   - Add unit test coverage in `cli/cmd/uninstall_unit_test.go` and `cli/cmdinstall/uninstall_tools_test.go` verifying flag parsing and plan construction.

---

## 4. Acceptance Criteria

- [ ] `gitmap uninstall` parses `--node <alias>`, `--dry-run`, `--force`, `--purge`.
- [ ] Positional tool argument parsing correctly resolves aliases for `pwsh`, `agm`, `vim`, `vscode`, `cursor`.
- [ ] Passing `--node <alias>` cleanly delegates to `cmdssh.RunSSHExec` with reconstructed arguments.
- [ ] Passing `--dry-run` outputs full execution plan without touching filesystem or packages.
- [ ] Workspace directories (`d:\work`, `/home/a/git-work`) are protected by assertion guards.
