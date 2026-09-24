# Specification 156: OS Discovery, Enum Reusability, Cross-Platform Shell Runner, and Node OS Profiling

## User Request (Verbatim)

To log into the nodes, the first time when you log in, you should run the Git maps, which OS command that would actually reveal the OS type, OS group, OS version, everything, build version, OS version, everything. Okay? So that you can actually learn from the OS enum. I hope you have learned that and reused the code so that this can be displayed, and this can also be displayed as JSON flag. Okay? So when we are logging to different machines and Git map is there, if not, then Git map will be installed, and then Git map will tell us which type of OS that is, and that would be saved to the SQLite DB so that in future we know which OS requires which type of instruction to be executed. If we are running an installer, by default, it would skip the Unix version. If we're doing a POSIX install, then it would by default skip the Windows version. But also we can run a bash command in Windows using the Git bash, so remember that. And Git map would be powerful enough to run bash command in Windows and also the Ubuntu. Okay? So that would be the power. So if we run a bash, Git map bash or shell, so both would run using the Git bash. It would first make sure that Git is installed, and where that is, it would find that and put that information into the root DB. If the Git is not installed, then it will say, "Hey, you have to install the Git using the Git map SSH, Git compact, or Git install command." So it will suggest that which one the user wants to pick, and how they could install using SSH, that would also be suggested. So these are the things I think very important that you need to work on. I hope you understand. If you have any question, confusion, let me know. Once you apply this, then we can actually run smoothly any commands in the future. Let's say we want to run a PowerShell command, it would already know that by default, we're going to run the PowerShell command in Windows. But also we should have another force flag using the execute force all, so that would execute some command in all OS, because in Ubuntu or Unix-based system, you can also install PowerShell, and then again, it will be same thing. If the target is not found, it will say, "Do you like to install the PowerShell using the Git..." Sorry, Git map, whatever command that is, it would suggest that. Do you understand? Do you have any question or confusion? If not, then please create tasks, proceed with implementation. Okay?

Please make sure to do git pull before commiting or chnaging to files please do it

and train yourself on gitmap llm train and use gitmap aum search to find file sfaster check the commands please and complete the stuff, clear???????

---

## 1. Domain Architecture & System Overview

This specification establishes a unified, cross-platform OS detection, command runner, and node profiling framework across GitMap:

1. **`which-os` CLI & Canonical OS Enum**:
   - Provides `gitmap which-os` (with aliases `whichos`, `os-which`, and `os-info`) to probe, classify, and format host OS identity.
   - Reuses canonical OS constants (`OSTargetWin`, `OSTargetUbuntu`, `OSTargetDebian`, `OSTargetMac`, `OSTargetUnix`, `OSTargetAll` from `cli/constants/os_targets.go`).
   - Reports `OSType`, `OSGroup`, `OSVersion`, `BuildVersion`, `Architecture`, `Platform`, `Hostname`, `NumCPU`, `Kernel`, `GitPath`, `BashPath`, and `PowerShellPath`.
   - Supports `--json` flag for machine-readable JSON output.

2. **Cross-Platform Bash & Shell Command Execution**:
   - Implements `gitmap bash` and `gitmap shell` commands that operate uniformly across Windows and Linux/Ubuntu.
   - On Windows, automatically discovers Git Bash (`C:\Program Files\Git\bin\bash.exe`, `C:\Program Files\Git\usr\bin\bash.exe`, etc.) and routes execution through it.
   - On Linux/macOS, invokes standard system `bash` or `sh`.
   - Probes Git availability: if present, records Git binary path into the root database (`gitmap.db`).
   - If Git is absent, halts with clear actionable guidance suggesting `gitmap install git`, `gitmap ssh install git`, `gitmap compact`, or remote SSH package manager commands.

3. **OS-Aware Installer Filtering & Force-All Override**:
   - Default execution filtering: skips Unix/POSIX installers on Windows hosts, and skips Windows installers on Unix/POSIX hosts.
   - Introduces `--force-all` (`-f`) flag for execution overrides, enabling scripts and commands to execute across all target platforms unconditionally.
   - PowerShell handling: defaults PowerShell execution to Windows; on Unix/Ubuntu systems or under `--force-all`, probes for `pwsh` and outputs installation guidance if missing.

4. **Node Initial Login OS Profiling & SQLite Persistence**:
   - Upon first connecting to a remote node via SSH or cluster join, triggers automated OS detection.
   - If GitMap binary exists on the remote node, leverages `gitmap which-os --json` directly.
   - If GitMap is missing, provides install instructions and falls back to remote SSH discovery commands.
   - Persists OS profile attributes (`OSGroup`, `BuildVersion`, `GitPath`, `PackageManager`, etc.) into `ClusterNode` / `SSHHost` tables in SQLite Split-DB (`gitmap.db`).
   - Enables subsequent cluster runs to dispatch commands using cached OS capabilities without redundant roundtrip discovery.
