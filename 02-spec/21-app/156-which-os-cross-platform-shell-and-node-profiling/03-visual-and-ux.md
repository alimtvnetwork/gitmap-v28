# Specification 156: Visual Presentation & CLI UX Layouts

## 1. `gitmap which-os` Standard Output

```text
FIELD               VALUE
OS Type             win
OS Group            windows
OS Version          Windows 11 Pro 23H2 (Build 22631)
Build Version       22631
Architecture        amd64
Platform            windows/amd64
Hostname            WORKSTATION-01
CPUs                16
Git Path            C:\Program Files\Git\cmd\git.exe (installed)
Bash Path           C:\Program Files\Git\bin\bash.exe (Git Bash)
PowerShell Path     C:\Program Files\PowerShell\7\pwsh.exe
```

## 2. `gitmap which-os --json` JSON Output

```json
{
  "osType": "win",
  "osGroup": "windows",
  "osVersion": "Windows 11 Pro 23H2 (Build 22631)",
  "buildVersion": "22631",
  "architecture": "amd64",
  "platform": "windows/amd64",
  "hostname": "WORKSTATION-01",
  "numCpu": 16,
  "gitPath": "C:\\Program Files\\Git\\cmd\\git.exe",
  "bashPath": "C:\\Program Files\\Git\\bin\\bash.exe",
  "powerShellPath": "C:\\Program Files\\PowerShell\\7\\pwsh.exe",
  "hasGit": true,
  "hasBash": true,
  "hasPowerShell": true
}
```

## 3. Git Missing Warning & Suggestion UX

When `gitmap bash` or `gitmap shell` runs on a machine where Git cannot be found:

```text
✗ Git is not installed or could not be located on this system.

GitMap requires Git to execute bash and manage repositories.
Please install Git using one of the following commands:

  Local Installation:
    ● gitmap install git
    ● gitmap compact

  Remote SSH Installation:
    ● gitmap ssh install git --target <node-alias>
    ● gitmap ssh exec <node-alias> "sudo apt-get update && sudo apt-get install -y git"   [Ubuntu/Debian]
    ● gitmap ssh exec <node-alias> "winget install --id Git.Git -e --source winget"       [Windows]
```

## 4. PowerShell Missing on Unix Suggestion UX

When PowerShell is executed on a Unix machine without `pwsh`:

```text
✗ PowerShell (pwsh) is not installed on this Unix node.

Would you like to install PowerShell using GitMap?
  ● Local:  gitmap install powershell
  ● Remote: gitmap ssh exec <node-alias> "sudo snap install powershell --classic"
  ● Apt:    gitmap ssh exec <node-alias> "sudo apt-get install -y powershell"
```
