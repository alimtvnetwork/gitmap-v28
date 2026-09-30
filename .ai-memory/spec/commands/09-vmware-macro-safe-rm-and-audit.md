# Command Specification: VMware Automation, Macro Idempotent Removal & Audit History

Spec Reference: [02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md](../../../02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md)

## 1. Commands Added & Enhanced

### 1.1 VMware Automation Engine (`manage-vm.ps1`)
Single-file PowerShell automation engine located in `repo-secrets/vmware/manage-vm.ps1` (and `scripts/vmware/manage-vm.ps1`).

```powershell
# Scan for VMs
.\manage-vm.ps1 -Action Scan -VMPath "D:\VMs"

# Inspect VM Configuration
.\manage-vm.ps1 -Action Inspect -VMPath "D:\VMs\WinServer\WinServer.vmx"

# Change MAC Address (Static or Generated)
.\manage-vm.ps1 -Action Set-MAC -VMPath "D:\VMs\WinServer\WinServer.vmx" -MAC "00:50:56:38:57:7B" -Adapter ethernet0

# Network Configuration (Bridged, NAT, HostOnly, Custom)
.\manage-vm.ps1 -Action Set-Network -VMPath "D:\VMs\Ubuntu\Ubuntu.vmx" -NetworkType bridged -Connected $true -ConnectAtPowerOn $true

# Adjust Hardware (RAM, vCPUs, 3D)
.\manage-vm.ps1 -Action Set-Hardware -VMPath "D:\VMs\WinServer\WinServer.vmx" -MemoryGB 8 -vCPUs 20

# Disks (List, Repair, Expand)
.\manage-vm.ps1 -Action Manage-Disk -VMPath "D:\VMs\WinServer\WinServer.vmx" -DiskAction List
.\manage-vm.ps1 -Action Manage-Disk -VMPath "D:\VMs\WinServer\WinServer.vmx" -DiskAction Repair

# Printer Device
.\manage-vm.ps1 -Action Manage-Printer -VMPath "D:\VMs\WinServer\WinServer.vmx" -PrinterEnabled $false

# Snapshots (List, Create, Revert, Clone)
.\manage-vm.ps1 -Action Manage-Snapshot -VMPath "D:\VMs\WinServer\WinServer.vmx" -SnapshotAction List
.\manage-vm.ps1 -Action Manage-Snapshot -VMPath "D:\VMs\WinServer\WinServer.vmx" -SnapshotAction Create -SnapshotName "Clean-Baseline"

# Copy Execution Log to Clipboard
.\manage-vm.ps1 -Action <Action> -CopyLog
```

### 1.2 Macro Idempotent Removal & CLI Safe-Rm (`gitmap safe-rm`)
- **CLI Command:** `gitmap safe-rm <path...> [--force]` (alias `gitmap rm-safe`). Idempotently removes files or directories, exiting 0 even if target does not exist.
- **Macro Step & Live Runner:** When macros or live interactive runners execute commands starting with `rm`, `rmdir`, `Remove-Item`, `rd`, or `del` on Windows PowerShell, they are automatically transformed into idempotent scripts checking `Test-Path -LiteralPath` before deleting.
- Deletion of non-existent files or directories returns exit code 0 rather than terminating with `ItemNotFoundException`.

### 1.3 Interactive Macro Editor Polish (`gitmap macro edit <name>`)
- Shell commands typed in interactive edit mode (such as `mkdir -p test`) are appended directly into macro steps.
- In-builder transient helpers are restricted to colon-prefixed commands (`:ls`, `:cat`), `cd`, `pwd`, or `help`.

### 1.4 Task Audit History View (`gitmap task history`)
- Displays execution audit records across `macro`, `ssh`, and `installer` subsystems from split-DB `TaskHistory` table.
- Filter by subsystem: `gitmap task history --section macro|ssh|installer`.
- Pagination: `gitmap task history [limit] [offset]`.
