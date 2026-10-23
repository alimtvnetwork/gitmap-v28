# Spec 313: GitMap Execution and Optimization Instructions

- **Domain:** Distributed Architecture, Remote Fleet SSH, Split-DB Concurrency Gate, Encrypted Key Vault, Agent & Heatmap JSON Telemetry, Repository Scan Governance
- **Sequence:** `313`
- **Slug:** `gitmap-execution-and-optimization-instructions`
- **Status:** `active`
- **Parent Plan:** [.ai-memory/plans/pending/313-gitmap-execution-and-optimization-instructions.md](../../../.ai-memory/plans/pending/313-gitmap-execution-and-optimization-instructions.md)

---

## Documents in this Specification Module

1. [01-architecture-spec.md](01-architecture-spec.md): High-level system architecture, Split-DB `gitmap-pull.db` design, Windows OpenSSH quote-escaping analysis, RSA-OAEP Key Vault, and Agent/Heatmap JSON contracts.
2. [02-component-spec.md](02-component-spec.md): Component-level Go interfaces, structs, database schemas (`PullTask`, `PullItemHistory`, `SshKey`), CLI subcommands (`scan disable-repos`, `scan enable-repos`), REST API endpoints, and React UI controls.
3. [03-migration-and-verification-spec.md](03-migration-and-verification-spec.md): Step-by-step rollout plan, database migrations, 7-day retention verification, targeted quality gate checklist, and acceptance criteria.
