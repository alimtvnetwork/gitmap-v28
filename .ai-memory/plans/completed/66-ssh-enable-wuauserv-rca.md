# Plan 66: OpenSSH Server Capability Installation Crash When Windows Update (wuauserv) Service is Disabled RCA & Fix

**Status:** Completed  
**Date:** 2026-10-03  
**Target:** `cli/cmdssh/ssh_daemon_enable.go`, `cli/cmdssh/ssh_daemon_enable_test.go`  
**Reference Issue:** `02-spec/22-app-issues/66-ssh-enable-wuauserv-disabled-capability-crash-rca.md`  

## Objectives & Deliverables
1. Fix OpenSSH Server capability detection on Windows by checking `sshd.exe` directly on disk and inspecting `Get-Service -Name sshd`.
2. Wrap `Add-WindowsCapability` in a PowerShell script that inspects `wuauserv`, temporarily enables/starts it if disabled, and cleanly stops/restores it to disabled in a `finally` block.
3. Suppress raw Go runtime stack traces when installation fails, providing clean, formatted diagnostic feedback with actionable next steps using `apperror.ErrorTypeAbort` with `reported: true`.
4. Add comprehensive unit tests in `cli/cmdssh/ssh_daemon_enable_test.go`.
5. Verify with local targeted linters and perform minor release bump.
