# Completed Plan 108: Remote Node First-Time Login OS Profiling, Cross-Platform Git Bash Execution & Multi-OS Installer Filtering

Spec Reference: [02-spec/21-app/156-which-os-cross-platform-shell-and-node-profiling/01-overview.md](../../02-spec/21-app/156-which-os-cross-platform-shell-and-node-profiling/01-overview.md)  
Execution Summary: Verified and completed across all 6 core deliverables with live execution proof, robust JSON parsing, automatic Windows/Unix installer filtering, `--force-all` bypass, and Git/PowerShell installation guidance.

## User Request (Verbatim)

```text
is it done properly and released??


To log into the nodes, the first time when you log in, you should run the Git maps, which OS command that would actually reveal the OS type, OS group, OS version, everything, build version, OS version, everything. Okay? So that you can actually learn from the OS enum. I hope you have learned that and reused the code so that this can be displayed, and this can also be displayed as JSON flag. Okay? So when we are logging to different machines and Git map is there, if not, then Git map will be installed, and then Git map will tell us which type of OS that is, and that would be saved to the SQLite DB so that in future we know which OS requires which type of instruction to be executed. If we are running an installer, by default, it would skip the Unix version. If we're doing a POSIX install, then it would by default skip the Windows version. But also we can run a bash command in Windows using the Git bash, so remember that. And Git map would be powerful enough to run bash command in Windows and also the Ubuntu. Okay? So that would be the power. So if we run a bash, Git map bash or shell, so both would run using the Git bash. It would first make sure that Git is installed, and where that is, it would find that and put that information into the root DB. If the Git is not installed, then it will say, "Hey, you have to install the Git using the Git map SSH, Git compact, or Git install command." So it will suggest that which one the user wants to pick, and how they could install using SSH, that would also be suggested. So these are the things I think very important that you need to work on. I hope you understand. If you have any question, confusion, let me know. Once you apply this, then we can actually run smoothly any commands in the future. Let's say we want to run a PowerShell command, it would already know that by default, we're going to run the PowerShell command in Windows. But also we should have another force flag using the execute force all, so that would execute some command in all OS, because in Ubuntu or Unix-based system, you can also install PowerShell, and then again, it will be same thing. If the target is not found, it will say, "Do you like to install the PowerShell using the Git..." Sorry, Git map, whatever command that is, it would suggest that. Do you understand? Do you have any question or confusion? If not, then please create tasks, proceed with implementation. Okay?

Please make sure to do git pull before commiting or chnaging to files please do it

and train yourself on gitmap llm train and use gitmap aum search to find file sfaster check the commands please and complete the stuff, clear???????
```

---

## Consolidated Subtasks & Outcomes

### Subtask 01: Release Confirmation & Remote Repository Pull
- Verified release `v6.338.0` is active and pushed to `origin main`.
- Executed `git pull origin main` to verify up-to-date workspace status.
- Trained on `gitmap llm train` curriculum and updated official skill `.agents/skills/gitmap/SKILL.md`.

### Subtask 02: Remote Node First-Time Login OS Profiling & Auto-Bootstrapping
- **Target Files:** `cli/cmdssh/ssh_login_cmd.go`, `cli/cmdssh/ssh_login_profile.go`, `cli/cmdssh/sshjoin_common.go`
- **Delivered:**
  - `probeAndEnsureNodeProfile` checks if node was previously profiled with valid OS version and `FirstRunAt`.
  - On first login: checks if GitMap is installed on remote machine; if missing, auto-bootstraps GitMap via `bootstrapRemoteGitmap`.
  - Executes remote `gitmap which-os --json` to probe full system OS report.
  - Made `probeRemoteGitmapWhichOS` JSON unmarshaling resilient by trimming between first `{` and last `}` and adding fallback binary path probes for non-interactive SSH environments.

### Subtask 03: SQLite DB Schema & OS Enum Persistence
- **Target Files:** `cli/db/sshconnection.go`, `cli/cmdos/os_info_types.go`, `cli/cmdssh/ssh_login_profile.go`
- **Delivered:**
  - Persists `OSType`, `OSGroup`, `OSVersion`, `BuildVersion`, `FirstRunAt` directly into SQLite `SSHConnection` table.
  - Reuses canonical OS target constants (`OSTargetWin`, `OSTargetUbuntu`, `OSTargetDebian`, `OSTargetMac`, `OSTargetUnix`).
  - Enables subsequent cluster runs to dispatch commands using cached OS capabilities without redundant roundtrip discovery.

### Subtask 04: Cross-Platform Git Bash Execution (`gitmap bash` / `shell`) & Git Advisor
- **Target Files:** `cli/cmd/bash_runner.go`, `cli/cmd/roottooling.go`
- **Delivered:**
  - Implemented `gitmap bash` and `gitmap shell` commands that operate uniformly across Windows and Linux/Ubuntu.
  - Automatically locates Git; records `git_path` into SQLite Split-DB (`gitmap.db` via `store.DB.SetSetting("git_path", gitPath)`).
  - On Windows, routes execution through Git Bash (`C:\Program Files\Git\bin\bash.exe`, etc.).
  - On POSIX, executes native `bash` or `sh`.
  - If Git is not installed, prints structured recommendations:
    `Hey, you have to install the Git using the Git map SSH, Git compact, or Git install command.`
    suggesting local, compact, and remote SSH installation commands.

### Subtask 05: Multi-OS PowerShell Execution (`--force-all`), Installer Filtering & pwsh Advisor
- **Target Files:** `cli/cmdssh/ssh_install_exec.go`, `cli/cmd/powershell_runner.go`, `cli/cmdinstaller/installer_os_cmds.go`
- **Delivered:**
  - In `gitmap ssh install-exec`:
    - Automatically targets Windows (`os: "win"`) for `.exe`, `.msi`, `.bat`, `.cmd`, `.ps1` installers, skipping Unix machines by default.
    - Automatically targets Unix (`os: "unix"`) for `.sh`, `.bash`, `.deb`, `.rpm`, `.bin`, `.run` installers, skipping Windows machines by default.
    - Added `-f, --force-all` flag to deploy installers to all machines regardless of default OS target.
  - In `gitmap powershell` / `gitmap pwsh`:
    - Defaults PowerShell execution to Windows; on Unix/Ubuntu systems or under `--force-all`, probes for `pwsh` and outputs installation guidance if missing:
      `Do you like to install the PowerShell using GitMap? (gitmap install powershell)`
  - Added unit test `TestParseInstallExecArgs_DefaultOSAndForceAll` in `cli/cmdssh/ssh_exec_install_test.go`.
