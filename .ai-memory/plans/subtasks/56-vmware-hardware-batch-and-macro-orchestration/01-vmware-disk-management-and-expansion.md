# Subtask [01]: Virtual Disk Lifecycle & Capacity Expansion
Traceability ID: Task-56-01
Spec Reference: [02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md](../../../../02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md)
Target Files: [scripts/vmware/manage-vm.ps1, cli/cmdvm/vm_disk.go]

Action:
- Extend `scripts/vmware/manage-vm.ps1` with `-Action Add-Disk` and `-Action Expand-Disk`.
- Implement safe `.vmdk.bak` backup before disk expansion.
- Invoke `vmware-vdiskmanager.exe` if available, or fall back to `.vmx` descriptor mutation.
- Expose CLI binding `gitmap vm disk add` and `gitmap vm disk expand`.

Acceptance Criteria:
- Existing virtual disk configurations are preserved verbatim with zero key loss.
- Disk capacity changes update geometry attributes accurately.
