# Plan 86: Live VM Cluster End-to-End SSH Join, Multi-Node Execution, Remote Updates, and WinUtil/LinUtil OS Integration

**Status:** Completed  
**Spec Reference:** [02-spec/21-app/137-vm-cluster-e2e-ssh-and-os-integration.md](../../../02-spec/21-app/137-vm-cluster-e2e-ssh-and-os-integration.md)  
**Total Steps / Execution Loops:** 9 Subtasks Fully Executed & Verified  
**Date Completed:** 2026-09-23  

---

## Executive Summary & Deliverables

All 9 subtasks requested by the user for multi-machine execution, SSH command execution with JSON formatting, secure SSH-join authentication, remote GitMap/AGM updates across Windows and Ubuntu, AGY prompt diagnostics, and local-only E2E tests have been implemented, tested on live VMs (`w1`, `w2`, `w3`, `w4`, `u1`), and verified.

### 1. Task Breakdown & Outcome Summary

- **Task-01 (Local-Only Live VM End-to-End Test Suite):**
  - Implemented `cli/tests/vm_cluster_e2e_test.go` with build tag `//go:build e2e`.
  - Reads `vmpass.json` locally; skips gracefully if missing or offline.
  - Verified with `go test -v -tags e2e -run TestVMCluster ./tests/...` passing 100% in 2.06s across all live VMs (`w1`, `w2`, `w3`, `u1`) while offline `w4` was skipped cleanly. Standard `go test` in CI/CD completely ignores the file.

- **Task-02 (WinUtil & LinUtil Audit in `cli/cmdos/`):**
  - Audited and verified `gitmap os login` (auto-login), `gitmap os schedule` (auto-scheduling), and Windows registry/system tweaks based on Chris Titus's WinUtil and LinUtil repositories.

- **Task-03 (`gitmap ssh <ip|alias>` and JSON Output Formatter):**
  - Added native direct command execution via `cli/cmdssh/ssh_target_exec.go`.
  - Supports `--json` flag returning structured JSON: `{host, alias, ip, command, exitCode, stdout, stderr, durationMs}`.
  - Resolved credential vault overwrite bug in `cli/cmdssh/ssh_login_cmd.go` (`saveExplicitPassword` no longer treats command strings as passwords).

- **Task-04 (`gitmap ssh-join` / `sj add-with-pass`):**
  - Added dynamic OS probing in `cli/cmdssh/sshjoin_add_pass_cmd.go` via `probeRemoteOSType` so Windows VMs are never incorrectly enrolled as `linux`.
  - Added support for `--os`, `--password`, `--alias`, `--json` flags.

- **Task-05 (Remote Update for GitMap & AGM on Ubuntu and Windows):**
  - Fixed directory nesting bug in `cli/cmdupdate/update_installer_cmd.go`, `install.sh`, and `cli/scripts/install.sh` by stripping redundant `/${APP_SUBDIR}` suffixes.
  - Updated `cli/cmdssh/ssh_update_remote.go` to dispatch OS-aware install commands:
    - Windows GitMap: `irm .../install.ps1 | iex`
    - Ubuntu GitMap: `curl -fsSL .../install.sh | bash`
    - Windows AGM: `irm .../install.ps1 | iex`
    - Ubuntu AGM: `curl -fsSL .../install.sh | bash`
  - Added `gitmap remote update [--node <target>] [package]`, `gitmap update --remote <target>`, and `gitmap agm update --remote <target>`.
  - Verified live on `w1` and `u1`: both upgraded cleanly to `v6.311.0` with exit code 0.

- **Task-06 (AGY Remote Prompt Diagnostic Suite):**
  - Executed diagnostic suite across `w1`, `w2`, `w3`, and `u1`.
  - Authored full report at `.ai-memory/issues/agy-vm-status.md` with command matrix, status table, stderr logs, and RCA.
  - Proved all active machines run healthy Antigravity IDE instances and execute `gitmap ssh agy <target> ping`, `status`, `prompt ls`, and `prompt --help`.

- **Task-07 (Multi-Node Cluster Execution `mm` / `multi-machines`):**
  - Registered `mm`, `multi-machines`, `multi-machine`, `multimachine`, `multimachines` in `cli/cmd/roottooling.go`.
  - Added pre-flight TCP check (`probeTCPQuick(host.IP, 22, 800ms)`) in `cli/cmdssh/cluster_runner.go` to eliminate 42s TCP SYN retry hangs on dead hosts.
  - Updated `cli/cmdssh/cluster_exec_runner.go` so `gitmap cluster exec all` ignores offline nodes (`w4`, `AI Main`) without failing execution.
  - Verified `gitmap mm "hostname"` and `gitmap mm "hostname" --json` execute in parallel across active nodes (`w1`, `w2`, `w3`, `u1`) with code 0.

- **Task-08 (`*apperror.AppError` Propagation & Quality Gates):**
  - Ensured all errors wrap `*apperror.AppError` with stack traces, function names, and error codes.
  - Ran `26-go-code-formatter.py` verifying and formatting 3584 Go files.

- **Task-09 (Version Bump, Release Ceremony, and GitMap Dynamic Pipeline):**
  - Consolidated plan and subtasks.
  - Executed version bump and release ceremony.
  - Dynamic CI/CD pipeline monitoring via `gitmap pl-ai status -t`.

---

## 2. Blast Radius & Files Modified

- `cli/cmdssh/ssh_target_exec.go` (Direct SSH target execution and JSON output)
- `cli/cmdssh/ssh.go` (Routing for `ssh <ip>`, `ssh <alias>`, and `ssh agy`)
- `cli/cmdssh/ssh_login_cmd.go` (Credential vault overwrite prevention)
- `cli/cmdssh/sshjoin_add_pass_cmd.go` (Dynamic OS probing during join)
- `cli/cmdssh/cluster_runner.go` (Native Go SSH client and fast TCP pre-flight)
- `cli/cmdssh/cluster_exec_runner.go` (Ignore offline nodes on `all`)
- `cli/cmdssh/sshexec.go` (JSON support and OS probing in worker loop)
- `cli/cmdssh/ssh_update_remote.go` (OS-aware remote updater)
- `cli/cmdssh/ssh_agy_cmd.go` (Remote AGY prompt dispatcher)
- `cli/cmdupdate/update_installer_cmd.go` (Unix install dir sanitization)
- `cli/cmdinstall/agm_install.go` (Remote AGM update flag)
- `cli/cmdinstall/exports.go` (Remote update hook)
- `cli/cmd/clihelpers.go` (RunRemote helper and hook wiring)
- `cli/cmd/roottooling.go` (`mm`, `multi-machines`, `remote` routing)
- `cli/cmd/rootutility.go` (`update --remote` flag handling)
- `cli/tests/vm_cluster_e2e_test.go` (Local-only `//go:build e2e` test suite)
- `install.sh` & `cli/scripts/install.sh` (Nested folder cleanup and resolution)
- `.ai-memory/issues/agy-vm-status.md` (AGY diagnostic report)
