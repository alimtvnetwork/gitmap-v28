# Subtask [03]: Universal Audit Task Logging and Polished Terminal History View
Traceability ID: Task-03
Spec Reference: [02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md](../../../02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md)
Target Files: [cli/cmdtask/task_history_cmd.go, cli/cmdtask/task_audit.go, cli/cmdmacro/macro_cmd.go, cli/cmdssh/ssh_history_tasks.go, cli/store/tasks_split_db.go]
Action:
- Implement a reusable task audit interceptor `RecordAuditTask(section, action, target, forwardPayload, status)` that writes to `TaskHistory` in `tasks_root` Split-DB.
- Wire audit recording into:
  - Macro: `macro add`, `macro edit`, `macro rm`, `macro run`
  - SSH: `ssh connect`, `ssh clone`, `ssh reclone`, `ssh join`
  - Installer: `installer create`, `installer run`
- Upgrade terminal display view in `cli/cmdtask/task_history_cmd.go`:
  - Format with termtable border rules, cyan headers, colored status badges (`[DONE]`, `[FAIL]`, `[PENDING]`), and human-readable timestamps.
  - Support section filtering (`gitmap task history --section macro`) and limit/offset pagination.
Acceptance Criteria:
- Running macros, SSH operations, or installers creates a persistent entry in `TaskHistory`.
- `gitmap task history` displays the history in a polished, readable terminal table.
Targeted Verification: [python 03-ai-scripts/05-guideline-autofixer.py cli/cmdtask/task_history_cmd.go]
