# Subtask 75.6: Audit Trail, Automated Runbook & Upstream AI Instruction Prompt

- **Parent Plan:** [75-gitmap-u1-ubuntu-agm-fleet-integration.md](../../pending/75-gitmap-u1-ubuntu-agm-fleet-integration.md)
- **Spec Reference:** [02-spec/21-app/205-gitmap-u1-ubuntu-agm-fleet-integration/01-architecture-spec.md](../../../../02-spec/21-app/205-gitmap-u1-ubuntu-agm-fleet-integration/01-architecture-spec.md)
- **Status:** Pending
- **Target Area:** `docs/`, `scripts/`

## Objective
Maintain a complete audit trail of commands and architectural decisions, provide an automated migration script for future Linux machines, and generate a standalone AI instruction prompt for upstream AGM developers to implement native headless Linux keyring support.

## Implementation Details
1. Create `docs/migration-windows-to-ubuntu-agm.md` recording all executed commands, terminal outputs, benefits, and architectural rationale.
2. Author `scripts/run-migration.sh` enabling 1-command reproducible migration on any fresh Ubuntu/Debian node.
3. Author `docs/agm-linux-account-switching-fix.prompt.md` containing full context, RCA, and detailed instructions for another AI agent to implement permanent native Linux fixes in the `Antigravity-Manager` repository.
