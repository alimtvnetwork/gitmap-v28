# Plan 101: SSH Execution OS Filtering (`--except-os`), Multi-Command Sequencing, and Remote Installer Deployment (`gitmap ssh install-exec`)

> **Status:** `COMPLETED`  
> **Release Version:** `v6.329.0`  
> **Spec Reference:** [`02-spec/21-app/152-ssh-exec-os-filter-and-install-exec.md`](../../../02-spec/21-app/152-ssh-exec-os-filter-and-install-exec.md)

---

## 1. Completed Subtasks

- [x] **SUBTASK-101-01**: Implemented `MatchesOSToken`, `IsConnectionOSExcluded`, `IsConnectionOSIncluded`, and `FilterSSHConnectionsByOS` in `cli/cmdssh/ssh_filter.go`, allowing `--except-os unix`, `--except-os win`, `--except-os ubuntu`, `--except-os darwin`, and `--os` filters across all SSH operations.
- [x] **SUBTASK-101-02**: Implemented `extractSECustomFlags` in `cli/cmdssh/sshexec.go` and `splitCommandsRespectingQuotes` in `cli/cmdssh/ssh_exec_command.go`. Correctly parses and chains multi-commands (`cmd1,cmd2,cmd3`) on target machines with OS-specific sequential operators (`;` for PowerShell on Windows, `&&` for bash/sh on Unix/Linux).
- [x] **SUBTASK-101-03**: Implemented `cli/cmdssh/ssh_install_exec.go` (`gitmap ssh install-exec <setup-file> [args...]` / `in-exec` / `setup-exec`) streaming local installer setups over pure SSH byte buffers to remote `%TEMP%` (Windows) or `/tmp` (Linux), executing silently with process wait, capturing exit codes, and outputting formatted results.
- [x] **SUBTASK-101-04**: Wired `install-exec` into `dispatchPackageSSH` in `cli/cmdssh/ssh.go`, `cli/cmdssh/exports.go`, and top-level utility dispatch in `cli/cmd/rootutility.go`.
- [x] **SUBTASK-101-05**: Authored and verified isolated temporary E2E test suite `cli/tests/e2e/ssh_exec_os_filter_and_install_exec_tempe2e_test.go` (`//go:build tempe2e`, `RUN_TEMP_E2E=1`) passing 100% of tests.
- [x] **SUBTASK-101-06**: Consolidated plan, bumped version to `v6.328.0`, updated changelogs, performed atomic commit, tagged, and pushed.
