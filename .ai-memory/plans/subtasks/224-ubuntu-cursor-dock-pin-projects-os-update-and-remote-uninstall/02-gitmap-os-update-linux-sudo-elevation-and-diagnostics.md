# Subtask Plan 02: GitMap OS Update Linux Sudo Elevation & Diagnostic Reporting

- **Spec Reference:** [02-spec/21-app/224-ubuntu-cursor-dock-pin-projects-os-update-and-remote-uninstall/01-architecture-spec.md](../../../../02-spec/21-app/224-ubuntu-cursor-dock-pin-projects-os-update-and-remote-uninstall/01-architecture-spec.md)
- **Status:** Queued
- **Target Area:** `cli/cmdos/os_update_engine.go`, `cli/cmdos/os_update_cmd.go`, `cli/cmdos/os_update_types.go`, `cli/cmdos/os_update_test.go`

---

## 1. Objective

Resolve package index update failures on Linux (`apt-get update` exiting with status 100) when executed by non-root users:
1. Detect non-root execution (`os.Geteuid() != 0`) on Linux (`runtime.GOOS == "linux"`).
2. Auto-elevate privileged package managers (`apt`, `dnf`, `pacman`) using `sudo -n` when `sudo` is available in `PATH`.
3. Capture execution stdout and stderr during command execution.
4. Enhance `printUpdateSummary` to display descriptive diagnostic error outputs and actionable remediation tips instead of a bare `✖ FAILED`.
5. Maintain test isolation adhering to repository guidelines (no real OS updates or power/system state modifications during tests).

---

## 2. Implementation Details

### Step 1: Update Toolchain Types & Privilege Classification (`cli/cmdos/os_update_types.go`)
1. Extend `UpdateToolchain`:
   ```go
   type UpdateToolchain struct {
       Name         string
       Binary       string
       UpdateArgs   []string
       UpgradeArgs  []string
       RequiresSudo bool
   }
   ```
2. In `discoverUpdateToolchains()`, mark `RequiresSudo: true` for:
   - `apt` (`apt-get`)
   - `dnf`
   - `pacman`
   User/isolated package managers (`brew`, `flatpak`, `snap`, `winget`) remain `RequiresSudo: false`.

### Step 2: Linux Sudo Auto-Elevation Engine (`cli/cmdos/os_update_engine.go`)
1. Add helper `isLinuxNonRoot() bool`:
   ```go
   func isLinuxNonRoot() bool {
       return runtime.GOOS == "linux" && os.Geteuid() != 0
   }
   ```
2. Add helper `hasSudoBinary() bool`:
   ```go
   func hasSudoBinary() bool {
       _, err := exec.LookPath("sudo")
       return err == nil
   }
   ```
3. In `executeToolchain(tc UpdateToolchain, isUpgrade bool, isDryRun bool) UpdateResult`:
   - Determine if elevation applies:
     `shouldElevate := tc.RequiresSudo && isLinuxNonRoot() && hasSudoBinary()`
   - In dry-run mode (`isDryRun == true`):
     - If `shouldElevate`: print `  • [dry-run] sudo -n %s %s\n`
     - Else: print `  • [dry-run] %s %s\n`
   - In actual execution:
     - If `shouldElevate`:
       ```go
       elevatedArgs := append([]string{"-n", tc.Binary}, args...)
       cmd = exec.Command("sudo", elevatedArgs...)
       ```
     - Else:
       ```go
       cmd = exec.Command(tc.Binary, args...)
       ```
   - Run `out, err := cmd.CombinedOutput()`.
   - Populate `UpdateResult{Name: tc.Name, Success: err == nil, Output: string(out), Error: err}`.

### Step 3: Diagnostic Reporting & Output Rendering (`cli/cmdos/os_update_cmd.go`)
1. Refactor `printUpdateSummary(results []UpdateResult)`:
   - Render each result:
     - On success: `  • %-10s : ✔ OK\n`
     - On failure: `  • %-10s : ✖ FAILED\n`
   - For all failed results, render an indented diagnostics section:
     ```text
     ▶ Diagnostic Details for Failed Managers:
       [apt] Exit Error: exit status 100
       [apt] Output:
         E: Could not open lock file /var/lib/apt/lists/lock - open (13: Permission denied)
         E: Unable to lock directory /var/lib/apt/lists/
       [apt] Tip: Ensure the user has sudo permissions or configure passwordless sudo in /etc/sudoers.d/
     ```
   - Trim output and cap to the most relevant 10 diagnostic lines to prevent terminal buffer flooding.

### Step 4: Unit Test Coverage (`cli/cmdos/os_update_test.go`)
1. Add test `TestDiscoverUpdateToolchains_RequiresSudo`:
   - Assert `RequiresSudo` is true for `apt`, `dnf`, `pacman`.
   - Assert `RequiresSudo` is false for `brew`, `flatpak`, `snap`, `winget`.
2. Add test `TestPrintUpdateSummary_Diagnostics`:
   - Construct simulated failed `UpdateResult` with mock output.
   - Capture stdout and verify diagnostic tips and error messages are rendered.
3. Adhere to `02-spec/02-coding-guidelines/` isolating test execution with mocks/dry-runs.

---

## 3. Verification Commands

```bash
# 1. Run unit tests for cmdos package
go test -v ./cli/cmdos/... -run "TestDiscover|TestUpdate|TestPrint"

# 2. Verify dry-run behavior on Linux
./gitmap os update --dry-run

# 3. Check code style and guideline compliance
python 03-ai-scripts/05-guideline-autofixer.py cli/cmdos --check-only
```
