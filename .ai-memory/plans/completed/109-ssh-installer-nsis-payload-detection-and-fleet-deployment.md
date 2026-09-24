# Plan 109: SSH Installer NSIS Payload Detection, Unattended Execution & Fleet Deployment

- **Status:** Completed
- **Date:** 2026-09-25
- **Related Spec:** Spec 158 (02-spec/21-app/158-ssh-installer-payload-nsis-auto-detection-and-fleet-deployment.md)
- **Tracking RCA:** Issue 44 (02-spec/22-app-issues/44-ssh-install-exec-nsis-silent-hanging-and-empty-registry-resolution-rca.md)

---

## 1. Objectives

1. Identify and fix root cause of SSH installer hanging / failing on Windows nodes (`Antigravity.Manager.Tools_4.70.0_x64-setup.exe`).
2. Resolve empty SSH host registry behavior when executed outside repository root (e.g. `Downloads`).
3. Ensure automated NSIS signature detection (`NullsoftInst`) and injection of `/S` silent flags.
4. Default remote SSH executions to unattended silent mode while preserving CLI overrides.
5. Eliminate unmocked network dial latency in unit tests, fixing CI/CD 10-minute timeout failure.
6. Verify live deployment end-to-end on cluster nodes (`w1`, `w2`, `w3`).

---

## 2. Execution Log

- [x] Enrolled cluster nodes `w1` (`192.168.1.3`), `w2` (`192.168.1.7`), `w3` (`192.168.1.12`) into central `gitmap.db` using credentials from `vmpass.json` via `gitmap sjc`.
- [x] Identified that `Antigravity.Manager.Tools_4.70.0_x64-setup.exe` is an NSIS installer requiring `/S`, not Inno Setup (`/VERYSILENT`).
- [x] Implemented `BuildRemoteInstallerExecCmdWithPayload` and `resolveDefaultSilentInstallerArgs` in `cli/cmdssh/ssh_install_exec.go`.
- [x] Defaulted `SSHInstallExecOptions.IsSilent = true` in `ParseInstallExecArgs`.
- [x] Flattened all nested `if` statements across `lowercasefix_commit.go`, `lowercasefix_ops.go`, `lowercasefix_report.go`, and `ssh_install_exec.go`.
- [x] Fixed raw rune conversion in `cli/tests/e2e/agy_rerun_restart_e2e_test.go`.
- [x] Added mock test hooks `SetConnectProbeClientForTesting`, `SetAutoTrustTargetHostForTesting`, and `SetCryptoConnectWithKeyForTesting`.
- [x] Pinned `Timeout: 4 * time.Second` in `cli/crypto/ssh_client.go` to prevent unbounded network hanging.
- [x] Verified unit tests and isolated temporary E2E tests (`TestTempE2E_SSH*`).
- [x] Tested live fleet deployment from `C:\Users\Administrator\Downloads>` on `w1`, `w2`, `w3` with 100% success (`✔ INSTALLED (0)` in 2.7s - 3.8s).
- [x] Authored canonical Spec 158 and RCA 44.

---

## 3. Verified Outcomes

- `gitmap ssh install-exec "Z:\VmDownloads\04. SharedSoft\Antigravity.Manager.Tools_4.70.0_x64-setup.exe"` runs silently and succeeds across `w1`, `w2`, `w3` with zero user flags needed.
- `go test ./cmdssh` suite reduced from 600s+ timeout to under 42s.
- `python linter-scripts/check-nested-ifs.py`: PASS (0 violations).
- `python linter-scripts/check-enum-and-boolean.py`: PASS (0 violations).
- `python linter-scripts/check-relative-paths.py`: PASS (0 violations across 7703 files).
- `python linter-scripts/check-error-management.py`: PASS (0 violations across 3768 files).
