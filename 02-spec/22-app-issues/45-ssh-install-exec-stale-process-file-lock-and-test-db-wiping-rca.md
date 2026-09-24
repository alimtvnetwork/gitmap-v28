# Issue 45 RCA: SSH Install-Exec Stale Process File Lock & Unit Test DB Erasure

## 1. Reproduction

1. **Test-Induced Registry Erasure:**
   - Execute `go test ./cli/cmdssh/...`.
   - Run `gitmap ssh install-exec .\Antigravity.Manager.Tools_4.70.0_x64-setup.exe --dry-run`.
   - Result: `No matching SSH machines to deploy setup to (filtered by except: "", except-os: "", os: "win")`.
   - Cause: `TestRunSSHJoinCLI_Validation` invoked `runSSHJoinCLI([]string{"rm"})` without `withMockSSHDB`, falling back to `openDB()` which opened the real `AppData\Local\gitmap-cli\data\gitmap.db` and cleared all registered nodes.

2. **Session 0 Process File Lock (`✖ UPLOAD FAILED 0ms`):**
   - Execute an unattended installer that previously launched interactively in Windows OpenSSH Session 0.
   - Re-attempt deployment: `gitmap ssh install-exec .\Antigravity.Manager.Tools_4.70.0_x64-setup.exe`.
   - Result: `streamTarToRemote.Wait: Can't unlink already-existing object: Permission denied`, followed by `streamDirectToRemote.Write: EOF` (`0ms`).
   - Cause: The previous instance (e.g. PID 12576) remained alive in Session 0, holding an exclusive read/execute lock on `C:\Windows\Temp\Antigravity.Manager.Tools_4.70.0_x64-setup.exe`.

## 2. Root Cause Analysis

- **Registry Erasure:** `cli/cmdssh` lacked a package-level `TestMain` to redirect global store openers (`openSSHDBFunc` and `openSSHDB`). Any unmocked test calling CLI removal commands directly operated on production/user database state.
- **File Lock:** `StreamFileToRemote` attempted to write bytes to the target file path without validating whether a prior process with the same name was still running on the remote node.
- **PowerShell Quote Stripping:** `powershell -NoProfile -Command "..."` through Windows OpenSSH stripped double quotes, throwing `The string is missing the terminator: "` during remote directory verification.

## 3. Corrective Implementation

1. **Package-Level Hermetic Test Isolation (`cli/cmdssh/main_test.go`):**
   - Implemented `TestMain(m *testing.M)` configuring `store.SetBinaryDataDirForTesting` to an ephemeral directory with initialized SQLite tables.
   - Explicitly wrapped `TestRunSSHJoinCLI_Validation` in `withMockSSHDB`.

2. **Stale Process Termination (`killStaleInstallerProcess`):**
   - Added `killStaleInstallerProcess(client, fileName, c.OS)` prior to streaming.
   - Invokes `cmd.exe /c taskkill /F /IM "<base>*" 2>nul` on Windows and `pkill -9 -f` on Unix to release file locks.

3. **Guaranteed Silent Switch Injection (`injectSilentFlag`):**
   - Refactored `BuildRemoteInstallerExecCmdWithPayload` to check `containsSilentArg(args)`.
   - Automatically prepends or injects default silent switches (`/S` for NSIS, `/VERYSILENT` for Inno) whenever `isSilent` is true.

4. **Native Remote Directory Bootstrap (`ensureRemoteDir`):**
   - Switched from fragile PowerShell commands to `cmd.exe /c if not exist "%s" mkdir "%s"`.

## 4. Verification & Prevention

- Ran `go test -count=1 ./cli/cmdssh/...`: all unit tests passed in 33.9s.
- Verified central `gitmap.db` retention: `ssh_hosts` maintained all 3 cluster nodes (`w1`, `w2`, `w3`) through test runs without modification.
- Executed live installation across all 3 nodes:
  - `w1` (192.168.1.3:22): `✔ INSTALLED (0)` in 2680ms
  - `w2` (192.168.1.7:22): `✔ INSTALLED (0)` in 2727ms
  - `w3` (192.168.1.12:22): `✔ INSTALLED (0)` in 2743ms
