# What to Read: GitMap Reading & Onboarding Guide

> **Primary Documentation:** [readme.md](readme.md)  
> **Autonomous Agent Memory:** [.ai-memory/what-to-read.md](.ai-memory/what-to-read.md)  
> **Canonical Version:** [version.json](version.json) (`v6.511.0`)  

This document serves as the top-level navigational reading sequence for developers and autonomous AI agents working in the GitMap repository. To eliminate documentation drift and redundancy, root [readme.md](readme.md) is the authoritative front door for human and agent interaction, while [.ai-memory/what-to-read.md](.ai-memory/what-to-read.md) contains deep agent-specific execution rules and historical changelogs.

---

## Recommended Ingestion Sequence

To establish complete context before executing code modifications, follow this phased reading order:

### Phase 1: Identity & Core Architecture
1. **[version.json](version.json)**: Canonical source of truth for versioning, repository slugs, and metadata.
2. **[readme.md](readme.md)**: Executive overview, core capabilities, quick install, and top-level architecture.

### Phase 2: Memory & Autonomous Context
1. **[.ai-memory/what-to-read.md](.ai-memory/what-to-read.md)**: Agent instructions, recent memory writes, and mandatory pre-flight checks.
2. **[.ai-memory/overview.md](.ai-memory/overview.md)**: High-level architectural blueprint, split-DB subsystems, and active conventions.
3. **[.ai-memory/plans/](.ai-memory/plans/)**: Active task ledgers, pending milestones, and subtask breakdown directories.

### Phase 3: Engineering Specifications & Standards
1. **[02-spec/01-spec-authoring-guide/](02-spec/01-spec-authoring-guide/)**: Spec authoring standards, required headings, and verification checklists.
2. **[02-spec/02-coding-guidelines/](02-spec/02-coding-guidelines/)**: Zero-nesting conventions, positive boolean naming, error wrapping, and immutability rules.
3. **[02-spec/03-error-manage/](02-spec/03-error-manage/)**: Universal error envelope, error codes, and failure isolation patterns.

### Phase 4: Command Taxonomy & Engine Manuals
1. **[docs/commands/readme.md](docs/commands/readme.md)**: Full index of all GitMap CLI commands and categories.
2. **[docs/commands/cloning-architecture.md](docs/commands/cloning-architecture.md)**: Architectural guide mapping all 5 cloning engines and concurrency controls.
3. **[docs/benchmarks/benchmark.md](docs/benchmarks/benchmark.md)**: Polyglot performance benchmarks and latency comparison tables.
