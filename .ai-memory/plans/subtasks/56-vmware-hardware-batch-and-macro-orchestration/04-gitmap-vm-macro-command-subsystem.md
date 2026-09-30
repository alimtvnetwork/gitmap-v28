# Subtask [04]: GitMap CLI Subsystem & Macro Orchestration
Traceability ID: Task-56-04
Spec Reference: [02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md](../../../../02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md)
Target Files: [cli/cmd/root.go, cli/cmdvm/vm_root.go, cli/macro/execute.go]

Action:
- Register `gitmap vm` root command family in CLI dispatcher.
- Implement argument forwarders and quote normalization for macro runner compatibility.
- Ensure all `gitmap vm` actions log audit entries into `TaskHistory` in Split-DB (`tasks_root`).

Acceptance Criteria:
- `gitmap macro add vm_setup "gitmap vm net set --vm myvm.vmx --type bridged"` registers and executes cleanly.
- `gitmap task history --section vm` lists all executed actions with timestamp and status.
