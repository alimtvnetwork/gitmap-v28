# Subtask [02]: Macro Export, Import, and Multi-Node Fleet Deployment
Traceability ID: Task-02
Spec Reference: [02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md](../../../02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md)
Target Files: [cli/cmdmacro/macro_export.go, cli/cmdmacro/macro_import.go, cli/cmdmacro/macro_deploy_ssh.go, cli/cmdssh/ssh_deploy_router.go]
Action:
- Verify and refine `runMacroExport` to ensure `gitmap macro export <name> [file.json]` cleanly exports single macros or `--all` to structured JSON/YAML.
- Verify and refine `runMacroImport` to ensure `gitmap macro import <file.json>` imports macros with deduplication and overwrite confirmation.
- Connect `gitmap macro deploy <name>` and `gitmap deploy macro <name>` to distribute macros to remote SSH fleet nodes via `ExecuteMacroDeploySSH`.
Acceptance Criteria:
- `gitmap macro export <name>` produces valid JSON on disk or stdout.
- `gitmap macro import <file.json>` restores the macro into the SQLite store.
- Fleet deployment command formats payload cleanly for remote execution.
Targeted Verification: [python 03-ai-scripts/05-guideline-autofixer.py cli/cmdmacro/macro_export.go]
