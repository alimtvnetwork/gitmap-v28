# Plan 101: SSH Execution OS Filtering (`--except-os`), Multi-Command Sequencing, and Remote Installer Deployment (`gitmap ssh install-exec`)

> **Status:** `PENDING`  
> **Target Release:** `v6.328.0`  
> **Spec Reference:** [`02-spec/21-app/152-ssh-exec-os-filter-and-install-exec.md`](../../../02-spec/21-app/152-ssh-exec-os-filter-and-install-exec.md)

---

## 1. User Request (Verbatim)

```text
gitmap ssh exec cmd1,cmd2,cmd3 --except-os unix
gitmap ssh exec cmd1,cmd2,cmd3 --except-os win
gitmap ssh exec cmd1,cmd2,cmd3 --except-os ubuntu

also please make

gitmap ssh install-exec <give exe setup> --except id, ip, alias # this will install the setup in the notes, can you do it????

Okay. So what I want is that can you install installers using Git map to another machine? So let's say I have a setup or something, and I wanted to run it on a different machine and install it. Can you perform something like this or add a command that could install, let's say, EXE on different Windows OS? And also I want to execute. You can have commands. Set OS, the new command, or you can do accept Windows. We can do accept

release please bump the version
```

---

## 2. Granular Subtask Decomposition

- [ ] **SUBTASK-101-01**: Implement OS normalization, `--except-os` / `--exclude-os`, and `--os` / `--target-os` filtering functions in `cli/cmdssh/ssh_filter.go` or `cli/cmdssh/ssh_exec.go`.
- [ ] **SUBTASK-101-02**: Support multi-command splitting and sequential execution chaining in `gitmap ssh exec cmd1,cmd2,cmd3`.
- [ ] **SUBTASK-101-03**: Implement `gitmap ssh install-exec <setup-file> [args...]` in `cli/cmdssh/ssh_install_exec.go`, uploading local installer setup to remote `%TEMP%` or `/tmp` via SSH byte stream, executing silently with process wait, and logging status across fleet nodes.
- [ ] **SUBTASK-101-04**: Wire CLI commands and routing in `cli/cmdssh/ssh_cmd.go` and `cli/cmd/ssh_cmd.go`.
- [ ] **SUBTASK-101-05**: Author and run isolated temporary E2E test suite `cli/tests/e2e/ssh_exec_os_filter_and_install_exec_tempe2e_test.go` (`//go:build tempe2e`, `RUN_TEMP_E2E=1`).
- [ ] **SUBTASK-101-06**: Consolidate plan, bump version to `v6.328.0`, update changelogs, commit single atomic commit, tag, and push.
