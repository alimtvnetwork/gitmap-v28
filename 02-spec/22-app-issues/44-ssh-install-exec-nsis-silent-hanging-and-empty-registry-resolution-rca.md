# Issue 44: Root Cause Analysis (RCA) - SSH Installer Execution Failure, NSIS Session 0 GUI Hanging & Unmocked Network Dial Timeouts

- **Status:** Resolved
- **Date:** 2026-09-25
- **Severity:** High
- **Component:** `cli/cmdssh`, `cli/crypto`
- **Related Spec:** Spec 158 (02-spec/21-app/158-ssh-installer-payload-nsis-auto-detection-and-fleet-deployment.md)

---

## Part 1: Symptom & Impact

### 1.1 Symptoms
1. **Empty Fleet Discovery in Ad-Hoc Directories:**
   Running `gitmap ssh install-exec .\Antigravity.Manager.Tools_4.70.0_x64-setup.exe --dry-run` from `C:\Users\Administrator\Downloads>` resulted in:
   `No matching SSH machines to deploy setup to (filtered by except: "", except-os: "", os: "win").`
   The CLI failed to communicate whether any SSH nodes were enrolled in the registry.
2. **Interactive Hanging on Remote Windows Nodes:**
   When targeting active Windows nodes, `Antigravity.Manager.Tools_4.70.0_x64-setup.exe` uploaded successfully, but the installer execution command hung indefinitely or timed out after 20+ seconds with exit code `4294967295`. Inspection of the remote node revealed `Antigravity.Manager.Tools_4.70.0_x64-setup.exe` stuck in Session 0 waiting for user interaction.
3. **CI/CD Quality Gate Failures & Test Timeout Alarm:**
   - Commit `46cc7fa` failed on GitHub Actions in `Boolean & Enum Linter` and `Nested If Linter` due to nested `if` statements in `lowercasefix_*.go` and `ssh_install_exec.go`.
   - `Cross-Platform Build` test suite timed out after 600s (`panic: test timed out after 10m0s`), blocking at `TestRunSSHLogin_DirectIPWithPassword`.

---

## Part 2: Grounded Root Cause Analysis

### 2.1 NSIS Parameter Mismatch & Missing Silent Default
- `BuildRemoteInstallerExecCmd` previously assumed all Windows `.exe` installers accept Inno Setup arguments (`/VERYSILENT /SUPPRESSMSGBOXES /NORESTART /SP-`).
- `Antigravity.Manager.Tools_4.70.0_x64-setup.exe` is built with NSIS (Nullsoft Scriptable Install System), signed with `NullsoftInst`.
- NSIS rejects `/VERYSILENT` and requires `/S` for silent mode.
- In addition, `SSHInstallExecOptions.IsSilent` was initialized to `false` unless explicitly commanded via `-s / --silent`.
- When invoked without arguments over SSH, NSIS launched its GUI installer in Session 0 without desktop access, blocking on window creation until terminated.

### 2.2 Split-DB Registry Enrolment
- `fetchAllSSHConnections()` queries `C:\Users\Administrator\AppData\Local\gitmap-cli\data\gitmap.db`.
- The database tables `ssh_hosts` and `SSHConnection` were unpopulated, resulting in 0 connections. The CLI did not distinguish between a filter mismatch and an empty registry.

### 2.3 Unmocked Network Dials in Unit Tests
- `executeSSHLoginWithPassword` called `autoTrustTargetHost` and `probeAndEnsureNodeProfile` unconditionally.
- In unit tests with mock IPs (e.g. `192.168.1.88`, `10.0.0.5`), `autoTrustTargetHost` and `connectProbeClient` called `net.DialTimeout` and `ssh.Dial`.
- Without a test hook, each non-routable IP caused TCP SYN retransmissions and 4-second to 2-minute timeouts, totaling over 600 seconds on CI runners.

---

## Part 3: Corrective & Preventative Actions

1. **Payload-Aware Installer Command Builder:**
   Created `BuildRemoteInstallerExecCmdWithPayload` and `resolveDefaultSilentInstallerArgs`. Inspects payload for `NullsoftInst`, `Inno Setup`, and `WixBundle` signatures. Automatically injects `/S` for NSIS.
2. **Default Silent Execution for Remote SSH Sessions:**
   `ParseInstallExecArgs` now defaults `IsSilent: true`. Added `--no-silent`, `--gui`, and `--interactive` for manual overrides.
3. **Registry Guidance on 0 Nodes:**
   Added `renderNoMatchingMachinesMessage` to detect when 0 nodes are enrolled and advise running `gitmap sjc` or `gitmap sj add`.
4. **Test Isolation Hooks:**
   - Introduced `SetConnectProbeClientForTesting`, `SetAutoTrustTargetHostForTesting`, and `SetCryptoConnectWithKeyForTesting`.
   - Set `Timeout: 4 * time.Second` on `ssh.ClientConfig` in `crypto.ConnectWithPassword` and `crypto.ConnectWithKey`.
5. **Coding Guidelines & Linter Compliance:**
   Flattened all nested `if` statements in `lowercasefix_commit.go`, `lowercasefix_ops.go`, `lowercasefix_report.go`, and `ssh_install_exec.go`. Fixed rune numerical cast in `agy_rerun_restart_e2e_test.go`.

---

## Part 4: Verification & Regression Matrix

1. **Local Linter Verification:**
   - `python linter-scripts/check-nested-ifs.py`: PASS (0 violations).
   - `python linter-scripts/check-enum-and-boolean.py`: PASS (0 violations).
   - `python linter-scripts/check-enum-guidelines.py`: PASS (0 violations).
   - `python linter-scripts/check-relative-paths.py`: PASS (0 violations across 7703 files).
   - `python linter-scripts/check-error-management.py`: PASS (0 violations across 3768 files).
2. **Unit & Fast Suite Verification:**
   - `TestBuildRemoteInstallerExecCmd_PayloadDetection`: PASS.
   - `TestParseInstallExecArgs_IsSilentDefault`: PASS.
   - `TestRunSSHLogin_*`: 0.16s (previously 20.3s).
   - `TestSSHClient`: 0.01s (previously 4.0s).
   - `TestConnectWithDefaultKey_Unreachable`: 0.00s (previously 16.0s).
3. **Live Fleet Verification:**
   - Enrolled cluster nodes `w1` (`192.168.1.3`), `w2` (`192.168.1.7`), `w3` (`192.168.1.12`) via `gitmap sjc`.
   - Tested deployment of `Antigravity.Manager.Tools_4.70.0_x64-setup.exe` from `C:\Users\Administrator\Downloads>`:
     - Node `w1`: `✔ INSTALLED (0)` in 2736ms.
     - Node `w2`: `✔ INSTALLED (0)` in 2720ms.
     - Node `w3`: `✔ INSTALLED (0)` in 3813ms.
