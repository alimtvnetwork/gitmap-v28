# Completed Plan 104: OS Discovery, Enum Reusability, Cross-Platform Shell Runner, and Node OS Profiling

Spec Reference: [02-spec/21-app/156-which-os-cross-platform-shell-and-node-profiling/01-overview.md](../../../02-spec/21-app/156-which-os-cross-platform-shell-and-node-profiling/01-overview.md)
Execution Summary: Completed in 1 unified orchestrator cycle across 6 deliverables.
Consolidation Status: 100% Complete.

## Original User Request (Verbatim)

To log into the nodes, the first time when you log in, you should run the Git maps, which OS command that would actually reveal the OS type, OS group, OS version, everything, build version, OS version, everything. Okay? So that you can actually learn from the OS enum. I hope you have learned that and reused the code so that this can be displayed, and this can also be displayed as JSON flag. Okay? So when we are logging to different machines and Git map is there, if not, then Git map will be installed, and then Git map will tell us which type of OS that is, and that would be saved to the SQLite DB so that in future we know which OS requires which type of instruction to be executed. If we are running an installer, by default, it would skip the Unix version. If we're doing a POSIX install, then it would by default skip the Windows version. But also we can run a bash command in Windows using the Git bash, so remember that. And Git map would be powerful enough to run bash command in Windows and also the Ubuntu. Okay? So that would be the power. So if we run a bash, Git map bash or shell, so both would run using the Git bash. It would first make sure that Git is installed, and where that is, it would find that and put that information into the root DB. If the Git is not installed, then it will say, "Hey, you have to install the Git using the Git map SSH, Git compact, or Git install command." So it will suggest that which one the user wants to pick, and how they could install using SSH, that would also be suggested. So these are the things I think very important that you need to work on. I hope you understand. If you have any question, confusion, let me know. Once you apply this, then we can actually run smoothly any commands in the future. Let's say we want to run a PowerShell command, it would already know that by default, we're going to run the PowerShell command in Windows. But also we should have another force flag using the execute force all, so that would execute some command in all OS, because in Ubuntu or Unix-based system, you can also install PowerShell, and then again, it will be same thing. If the target is not found, it will say, "Do you like to install the PowerShell using the Git..." Sorry, Git map, whatever command that is, it would suggest that. Do you understand? Do you have any question or confusion? If not, then please create tasks, proceed with implementation. Okay?

Please make sure to do git pull before commiting or chnaging to files please do it

and train yourself on gitmap llm train and use gitmap aum search to find file sfaster check the commands please and complete the stuff, clear???????

---

## Consolidated Subtasks & Delivered Implementations

### Subtask 01: Which-OS Command and Canonical Enum Reusability
- **Traceability ID:** Task-02
- **Delivered Changes:**
  - Added OS Group constants (`OSGroupWindows`, `OSGroupUnix`, `OSGroupMac`, `OSGroupPOSIX`) to `cli/constants/os_targets.go`.
  - Defined `LocalOSProbe` and extended `OSInfoReport` with `OSGroup`, `BuildVersion`, `GitPath`, `BashPath`, `PowerShellPath`, `HasGit`, `HasBash`, `HasPowerShell` in `cli/cmdos/os_info_types.go`.
  - Updated `probeLocalPlatformOS` on Windows (`cli/cmdos/os_info_windows.go`) and POSIX (`cli/cmdos/os_info_other.go`) mapping canonical target constants.
  - Added tool discovery for Git, Bash, and PowerShell in `cli/cmdos/os_info.go`.
  - Registered `which-os`, `whichos`, `os-which` aliases in `cli/cmdos/os_info_cmd.go` and `cli/cmd/roottooling.go`.
- **Status:** Verified & Complete.

### Subtask 02: Cross-Platform Bash & Shell Command Runner with Git Detection
- **Traceability ID:** Task-03
- **Delivered Changes:**
  - Implemented `runBash` and `runShell` in `cli/cmd/bash_runner.go`.
  - Automatically locates Git; if found, records `git_path` into SQLite Split-DB (`gitmap.db` via `store.DB.SetSetting("git_path", gitPath)`).
  - On Windows, routes execution through Git Bash (`bash.exe` in Git directories).
  - On POSIX, executes native `bash` or `sh`.
  - If Git is not installed, prints structured recommendations (`gitmap install git`, `gitmap compact`, `gitmap ssh install git --target <node>`, or package manager commands).
  - Registered `bash`, `git-bash`, `shell`, `sh` in `cli/cmd/roottooling.go`.
- **Status:** Verified & Complete.

### Subtask 03: OS-Aware Installer Filtering and Force-All Execution Override
- **Traceability ID:** Task-04
- **Delivered Changes:**
  - Updated `executeOSInstall` in `cli/cmdinstaller/installer_os_cmds.go` to gate installer targets by host OS: skips Unix targets on Windows and Windows targets on POSIX by default.
  - Added `--force-all` / `-f` override flag parsing in installer commands and cluster exec options (`cli/cmdssh/cluster_exec_cmd.go`).
  - Added `install-win` command to installer subcommands in `cli/cmdinstaller/installer_os_cmds.go`.
- **Status:** Verified & Complete.

### Subtask 04: Node First-Login OS Profiling and SQLite DB Persistence
- **Traceability ID:** Task-05
- **Delivered Changes:**
  - Implemented `probeTargetWithWhichOS` and `probeRemoteGitmapWhichOS` in `cli/cmdssh/sshjoin_common.go` and `cli/cmdssh/sshjoin_enroll.go`.
  - On node login / SSH join, executes remote `gitmap which-os --json` to extract full OS telemetry.
  - If GitMap is absent on the remote node, prints install recommendation (`gitmap ssh deploy <alias>` / `gitmap ssh install gitmap -t <alias>`) and falls back to remote SSH discovery.
  - Persists node OS details to SQLite DB `SSHConnection` table.
- **Status:** Verified & Complete.

### Subtask 05: PowerShell Execution, OS Defaulting, and Unix Installation Prompts
- **Traceability ID:** Task-04
- **Delivered Changes:**
  - Implemented `runPowerShell` in `cli/cmd/powershell_runner.go` supporting `powershell`, `pwsh`, `ps`.
  - Defaults to Windows PowerShell; on Unix, checks for `pwsh` and prints actionable installation guidance if missing (`gitmap install powershell`, snap, or apt commands).
  - Enhanced `ExecPS` in `cli/cluster/exec_ps.go` to advise installing PowerShell via GitMap when absent on Unix nodes.
  - Registered `powershell`, `pwsh`, `ps` in `cli/cmd/roottooling.go`.
- **Status:** Verified & Complete.
