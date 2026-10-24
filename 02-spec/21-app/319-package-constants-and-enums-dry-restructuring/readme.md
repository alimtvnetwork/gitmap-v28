# 319 - Package Constants & Enums DRY Restructuring

## Overview
This specification defines the architectural design, directory conventions, and refactoring roadmap for centralizing constants and strongly-typed enums across 5 critical packages in the GitMap codebase:
1. `cli/cmdpull` - Pull modes, transport flags, progress indicators, status keys, diffstat indicators.
2. `cli/cloner` - Clone protocols, depth options, hooks, subprocess execution arguments.
3. `cli/cmdclone` - CLI options, clone strategy types, concurrency profiles, telemetry tokens.
4. `cli/store` - Split-DB table names, query operations, SQLite pragmas, transaction status codes.
5. `cli/cmdcommit` (5th Rollout Package) - Conventional commit verbs, module bracket formatters, validation flags.

## Document Index
- [01-architecture-spec.md](01-architecture-spec.md): Enum standards (`*Type` suffix, `IsValid()`, `String()`), constant layout, zero circular imports.
- [02-component-spec.md](02-component-spec.md): Deep domain mapping for each package, concrete enum variants, and constant identifiers.
- [03-migration-and-verification-spec.md](03-migration-and-verification-spec.md): Step-by-step migration protocol, backwards compatibility aliases, and AST verification.

## Execution Sequence
- **Phase 1 (Current Turn):** Architecture planning, discovery audit, specification authoring, and subtask ledger generation.
- **Phase 2 & 3 (Subsequent Turn upon user trigger `@[01-prompts/15-cg-execute/06-constants-and-enums.md]`):** Multi-agent tree swarm execution ($A = 3$, $H = 2$, Capacity = 6) applying the enums and constants across all 5 packages in bounded micro-batches.
