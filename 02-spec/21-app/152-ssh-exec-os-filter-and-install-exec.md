# Spec 152: SSH Execution OS Filtering (`--except-os`), Multi-Command Sequencing, and Remote Installer Deployment (`gitmap ssh install-exec`)

> **Spec ID:** `SPEC-152`  
> **Version:** `v6.328.0`  
> **Status:** Draft / Active  
> **Date:** 2026-09-24  

---

## 0. User Request (Verbatim) & Actionable Deliverables

```text
gitmap ssh exec cmd1,cmd2,cmd3 --except-os unix
gitmap ssh exec cmd1,cmd2,cmd3 --except-os win
gitmap ssh exec cmd1,cmd2,cmd3 --except-os ubuntu

also please make

gitmap ssh install-exec <give exe setup> --except id, ip, alias # this will install the setup in the notes, can you do it????

Okay. So what I want is that can you install installers using Git map to another machine? So let's say I have a setup or something, and I wanted to run it on a different machine and install it. Can you perform something like this or add a command that could install, let's say, EXE on different Windows OS? And also I want to execute. You can have commands. Set OS, the new command, or you can do accept Windows. We can do accept

release please bump the version
```

### Extracted Actionable Deliverables:
1. **OS Exclusion & Inclusion Filtering on SSH Exec (`--except-os`, `--os`)**:
   - Support `--except-os <os>` / `--exclude-os <os>` / `--skip-os <os>`:
     - `unix`: Filters out all non-Windows OSes (`linux`, `ubuntu`, `debian`, `centos`, `darwin`, `mac`, `macos`).
     - `win`, `windows`: Filters out all Windows machines.
     - `ubuntu`, `linux`, `debian`: Filters out Linux/Ubuntu distributions.
     - `darwin`, `mac`, `macos`: Filters out macOS nodes.
   - Support target OS filtering `--os <os>` / `--target-os <os>` to target only specific operating systems.
2. **Multi-Command Sequential Splitting**:
   - `gitmap ssh exec cmd1,cmd2,cmd3`: Automatically splits comma-separated commands (or standard chained commands) and runs them sequentially per machine, reporting execution output per command.
3. **Remote Installer Execution Engine (`gitmap ssh install-exec`)**:
   - Syntax: `gitmap ssh install-exec <local-setup-path> [installer-args...] [flags]`
   - Aliases: `ssh in-exec`, `ssh setup`, `ssh install-run`.
   - Streaming: Copies the setup file using low-level SSH byte-buffer streaming into a temporary folder (`%TEMP%` on Windows, `/tmp` on Unix/Linux).
   - Execution:
     - On Windows: Launches via PowerShell with `Start-Process -FilePath '<destPath>' -ArgumentList '<args>' -Wait -PassThru` (or headless silent switches), records exit code.
     - On Linux/macOS: Enforces `chmod +x` and invokes `./setup.sh` or `./setup.bin` with provided arguments, capturing exit status.
   - Filters: Honors `--except id,ip,alias`, `--except-os <os>`, and `--os <os>`.
4. **Temporary E2E Test Suite**:
   - `cli/tests/e2e/ssh_exec_os_filter_and_install_exec_tempe2e_test.go` guarded by `//go:build tempe2e` and `RUN_TEMP_E2E=1`.
5. **SemVer Version Bump & Git Release**:
   - Bump to `v6.328.0`, update changelogs, single atomic commit, tag `v6.328.0`, and push to `origin main`.

---

## 1. Technical Architecture & Protocols

### 1.1 OS Filtering Protocol

In `cli/cmdssh/`, every registered `SSHConnection` has an `OS` field (`windows`, `linux`, `darwin`).
We introduce a standard OS classifier:
```go
func MatchesOS(nodeOS, targetFilter string) bool
func IsOSExcluded(nodeOS, exceptFilter string) bool
```
- When `exceptFilter` is `unix`: matches any node where `nodeOS != "windows"`.
- When `exceptFilter` is `win` or `windows`: matches any node where `nodeOS == "windows"`.
- When `exceptFilter` is `linux`, `ubuntu`, `debian`, `centos`: matches any node where `nodeOS == "linux"`.
- When `exceptFilter` is `darwin`, `mac`, `macos`: matches any node where `nodeOS == "darwin"`.

### 1.2 Multi-Command Sequencing

In `gitmap ssh exec <commands>`:
If the command string contains commas (and is not an escaped single command), GitMap parses the list into individual steps:
`[cmd1, cmd2, cmd3]`
On Windows remote hosts, commands are chained using `cmd1 ; cmd2 ; cmd3` in PowerShell or `cmd1 && cmd2 && cmd3` in CMD.
On Unix remote hosts, commands are chained with `cmd1 && cmd2 && cmd3`.

### 1.3 Remote Installer Deployment Protocol (`gitmap ssh install-exec`)

1. **Source Discovery**: Reads the local installer file (e.g. `setup.exe`, `install.msi`, `agent.sh`).
2. **Fleet Filtering**: Applies node target filters, `--except id,ip,alias`, and `--except-os <os>`.
3. **Payload Streaming**: Streams file data over SSH channel into remote temp directory.
4. **Remote Execution Command Construction**:
   - Windows:
     ```powershell
     powershell -NoProfile -Command "$p = '<destPath>'; $proc = Start-Process -FilePath $p -ArgumentList '<args>' -Wait -PassThru; exit $proc.ExitCode"
     ```
   - Unix/Linux:
     ```bash
     chmod +x '<destPath>' && '<destPath>' <args>
     ```
5. **Result Aggregation**: Displays execution outcome table across all target nodes.
