# Specification 156: Data Contracts & Type Schemas

## 1. OS Enum Definitions (`cli/constants/os_targets.go`)

```go
package constants

const (
	// Target OS identifiers for installers & discovery
	OSTargetWin    = "win"
	OSTargetUbuntu = "ubuntu"
	OSTargetCentOS = "centos"
	OSTargetDebian = "debian"
	OSTargetFedora = "fedora"
	OSTargetArch   = "arch-linux"
	OSTargetMac    = "macos"
	OSTargetUnix   = "unix"
	OSTargetAll    = "all"

	// OS Group categories
	OSGroupWindows = "windows"
	OSGroupUnix    = "unix"
	OSGroupMac     = "macos"
	OSGroupPOSIX   = "posix"
)
```

## 2. OS Info Report Contract (`cli/cmdos/os_info_types.go`)

```go
package cmdos

type OSInfoReport struct {
	OSType             string `json:"osType"`
	OSGroup            string `json:"osGroup"`
	OSVersion          string `json:"osVersion"`
	BuildVersion       string `json:"buildVersion"`
	Architecture       string `json:"architecture"`
	Platform           string `json:"platform"`
	Hostname           string `json:"hostname"`
	NumCPU             int    `json:"numCpu"`
	Kernel             string `json:"kernel,omitempty"`
	GitPath            string `json:"gitPath,omitempty"`
	BashPath           string `json:"bashPath,omitempty"`
	PowerShellPath     string `json:"powerShellPath,omitempty"`
	HasGit             bool   `json:"hasGit"`
	HasBash            bool   `json:"hasBash"`
	HasPowerShell      bool   `json:"hasPowerShell"`
}
```

## 3. Database Schema Extensions (`ClusterNode` in `gitmap.db`)

```sql
-- Enhanced ClusterNode table schema
CREATE TABLE IF NOT EXISTS ClusterNode (
    NodeId TEXT PRIMARY KEY,
    Alias TEXT NOT NULL DEFAULT "",
    DisplayId INTEGER NOT NULL DEFAULT 0,
    IPAddress TEXT NOT NULL DEFAULT "",
    NodeRole TEXT NOT NULL DEFAULT "client",
    OS TEXT NOT NULL DEFAULT "windows",
    OSGroup TEXT NOT NULL DEFAULT "windows",
    OSVersion TEXT NOT NULL DEFAULT "",
    BuildVersion TEXT NOT NULL DEFAULT "",
    Architecture TEXT NOT NULL DEFAULT "amd64",
    GitPath TEXT NOT NULL DEFAULT "",
    BashPath TEXT NOT NULL DEFAULT "",
    JoinedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    LastHeartbeat TIMESTAMP,
    Status TEXT NOT NULL DEFAULT "online",
    PasswordHash TEXT,
    PackageManager TEXT
);
```

## 4. Root Database Key-Value Configuration

Key `system.git_path`: Stores discovered path to Git binary on host (e.g. `C:\Program Files\Git\cmd\git.exe` or `/usr/bin/git`).
Key `system.bash_path`: Stores discovered path to Bash binary on host (e.g. `C:\Program Files\Git\bin\bash.exe` or `/bin/bash`).
Key `system.os_type`: Normalized OS target identifier (`win`, `ubuntu`, `debian`, etc.).
Key `system.os_group`: Normalized OS group category (`windows`, `unix`, `posix`, `macos`).
