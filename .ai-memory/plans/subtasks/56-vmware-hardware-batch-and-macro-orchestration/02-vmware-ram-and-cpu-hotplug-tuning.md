# Subtask [02]: Dynamic RAM, CPU & Hotplug Resource Management
Traceability ID: Task-56-02
Spec Reference: [02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md](../../../../02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md)
Target Files: [scripts/vmware/manage-vm.ps1, cli/cmdvm/vm_hardware.go]

Action:
- Implement core topology tuning in `manage-vm.ps1` (`cpuid.coresPerSocket`, `numvcpus`, `memsize`).
- Add memory and CPU hot-add flags (`mem.hotadd = "TRUE"`, `vcpu.hotadd = "TRUE"`).
- Provide host boundary guard to reject RAM exceeding physical memory limits.

Acceptance Criteria:
- RAM values between 512 MB and 128 GB are validated.
- CPU core/socket calculation matches VMware Workstation constraints.
