# Subtask [03]: Multi-VM Discovery and Batch Operations
Traceability ID: Task-56-03
Spec Reference: [02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md](../../../../02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md)
Target Files: [cli/cmdvm/vm_batch.go, cli/cmdvm/vm_scan.go]

Action:
- Implement recursive multi-VM batch processing with directory scans (`--dir <path>`).
- Support filtering by VM name wildcard or network connection type.
- Execute batch hardware modifications in parallel or sequential modes with atomic file locking.

Acceptance Criteria:
- Batch runner modifies all matching `.vmx` files without leaving half-modified states.
- Clean progress bar and summary table emitted upon completion.
