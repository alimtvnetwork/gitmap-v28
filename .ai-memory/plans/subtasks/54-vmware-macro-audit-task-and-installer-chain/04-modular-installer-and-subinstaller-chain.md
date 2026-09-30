# Subtask [04]: OS-Specific Modular Installer and Chained Sub-Installer Engine
Traceability ID: Task-04
Spec Reference: [02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md](../../../02-spec/21-app/190-vmware-macro-audit-task-and-installer-chain.md)
Target Files: [cli/cmdinstall/install_archive_deploy.go, cli/cmdinstall/installer_types.go]
Action:
- Introduce OS targeting in installer recipes (`Windows`, `Unix`, `Ubuntu`, `CentOS`) with appropriate package manager hooks.
- Support subcommand chaining: define sub-installers that execute downstream recipes or subcommands in dependency order.
- Provide data models for chained execution and CLI command chaining.
Acceptance Criteria:
- Installer definitions support target OS filtering and sub-installer dependencies.
- Sub-installers execute sequentially or report clear dependency errors.
Targeted Verification: [python 03-ai-scripts/05-guideline-autofixer.py cli/cmdinstall/install_archive_deploy.go]
