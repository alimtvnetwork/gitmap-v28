# Subtask [05]: Single-File Enterprise VMware Automation PowerShell Script
Traceability ID: Task-05
Spec Reference: [02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md](../../../02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md)
Target Files: [D:/work/repo-secrets/vmware/manage-vm.ps1, scripts/vmware/manage-vm.ps1]
Action:
- Author comprehensive `manage-vm.ps1` in `D:/work/repo-secrets/vmware/manage-vm.ps1` and mirror to `scripts/vmware/manage-vm.ps1`.
- Implement robust `.vmx` configuration parser, validator, and modifier:
  - `-Action Scan`: Recursively find `.vmx` virtual machines in a target folder.
  - `-Action Inspect`: Show all virtual machine hardware properties (RAM, vCPUs, Disks, Network Adapters, MAC addresses).
  - `-Action Set-MAC`: Set static MAC address on adapter (e.g. `00:50:56:38:57:7B` from user screenshot) or generate new.
  - `-Action Set-Network`: Configure Bridged, NAT, HostOnly, or Custom virtual network (`VMnet`), connected and connect-at-power-on state.
  - `-Action Set-Hardware`: Modify RAM (`-MemoryMB`), Processors (`-vCPUs`), and 3D graphics display settings.
  - `-Action Manage-Disk`: List all virtual disks (NVMe, SCSI, SATA), test existence, repair via `vmware-vdiskmanager -R`, and expand disk size.
  - `-Action Manage-Printer`: Enable or remove virtual printer device.
  - `-Action Manage-Snapshot`: List snapshots, create snapshot, clone snapshot, and restore snapshot using `vmrun.exe` / `.vmsd`.
  - `-Action Optimize-VM`: Apply standard performance optimizations (disable memory trimming, paging tweaks).
- Error Handling & Diagnostics:
  - Comprehensive `try/catch` with timestamp, `$_.ScriptStackTrace`, and error origin.
  - `-CopyLog`: Automatically copy execution log to the Windows clipboard via `Set-Clipboard` for instant sharing.
  - `-Help` / `--help`: Formatted command reference with parameter definitions and copy-pasteable examples.
Acceptance Criteria:
- Script executes without syntax errors on Windows PowerShell 5.1 and PowerShell 7.
- `.vmx` modifications preserve all existing keys and formatting, creating `.vmx.bak` prior to modification.
- `-CopyLog` copies formatted output directly to clipboard.
Targeted Verification: [powershell -NoProfile -ExecutionPolicy Bypass -Command "& 'scripts/vmware/manage-vm.ps1' -Help"]
