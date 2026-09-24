# Plan 110: SSH Installer Hermetic Test Isolation & Stale Process Guard

## Metadata
- Spec Reference: [02-spec/21-app/159-ssh-installer-hermetic-test-isolation-and-stale-process-guard.md](../../../02-spec/21-app/159-ssh-installer-hermetic-test-isolation-and-stale-process-guard.md)
- RCA Reference: [02-spec/22-app-issues/45-ssh-install-exec-stale-process-file-lock-and-test-db-wiping-rca.md](../../../02-spec/22-app-issues/45-ssh-install-exec-stale-process-file-lock-and-test-db-wiping-rca.md)
- Status: Completed
- Duration / Cycles: 1 cycle (Parent Task N-Step Loop v2.5.0)

## Overview & Scope
Resolved two major production fleet installer defects:
1. Unit tests in `cli/cmdssh` mutating or erasing registered SSH nodes in the central database (`AppData/Local/gitmap-cli/data/gitmap.db`).
2. SSH installer payload streaming failing with `✖ UPLOAD FAILED 0ms` or permission denied due to stale installer instances running in Session 0 holding open file locks.

## Outcomes & Verification
- Created `cli/cmdssh/main_test.go` providing package-level ephemeral SQLite test isolation.
- Wrapped `TestRunSSHJoinCLI_Validation` in `withMockSSHDB`.
- Added `killStaleInstallerProcess` before streaming to terminate lingering setup processes.
- Enhanced `BuildRemoteInstallerExecCmdWithPayload` with unconditional silent switch injection (`/S` for NSIS).
- Switched Windows remote directory creation to native `cmd.exe /c if not exist ... mkdir ...`.
- Live fleet verification on `w1`, `w2`, `w3`: 100% success (`✔ INSTALLED (0)` in ~2.7s each).
- Clean `gitmap pe` pipeline verification: 100% green.
