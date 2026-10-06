# AI Memory — Authoritative Domain Knowledge Base

Welcome to the consolidated AI Memory repository for **GitMap**. This directory contains the authoritative, non-redundant architectural references, governance invariants, and operational guidelines governing all development workflows.

All previous fragmented memory notes, legacy audit reports, and temporary logs have been deeply compacted into these 9 core documents.

---

## Authoritative Domain References

| # | Reference Document | Domain Scope & Architectural Coverage |
| :---: | :--- | :--- |
| **00** | [00-project-governance-and-invariants.md](00-project-governance-and-invariants.md) | Non-negotiable architectural invariants: positive booleans (`is*`, `has*`), 100% relative paths, zero-build execution during maintenance, total git command ban for subagents, and append-only prohibited registry. |
| **01** | [01-cli-architecture-and-contracts.md](01-cli-architecture-and-contracts.md) | Cobra command tree, typo suggestions, JSON Envelope V2, Lipgloss ANSI styling, auto-aligned `termtable`, unbuffered `StreamWriter` interface, and macro automation engine. |
| **02** | [02-scanner-projects-and-deduplication.md](02-scanner-projects-and-deduplication.md) | High-speed filesystem traversal, polyglot project heuristics (Go, Python, TS, Rust, PHP, C#), fast file reader, `.gitmapignore` parsing, and zero-allocation deduplication. |
| **03** | [03-git-operations-commit-and-pull.md](03-git-operations-commit-and-pull.md) | Flat commit mechanics (`gitmap commit`, `cm`), automatic staging, append-only history invariant, PAS worker concurrency formula ($\min(\text{cores}, 8)$), and pull-all OMZ ignore handling. |
| **04** | [04-fleet-nodes-ssh-and-credentials.md](04-fleet-nodes-ssh-and-credentials.md) | Fleet nodes topology, ICMP/TCP port 22 concurrent liveness probing (1500ms timeout), RSA-OAEP encrypted credential vault (`credentials.vault`), masked password input, and remote clone engine. |
| **05** | [05-antigravity-and-ide-ecosystem.md](05-antigravity-and-ide-ecosystem.md) | Google Antigravity (AGY) agent architecture, prompt manager, decision logs, multi-IDE synchronization (VS Code, Cursor, AGY), autonomous pipeline fix injection, and Ubuntu workstation governance. |
| **06** | [06-database-engine-and-split-storage.md](06-database-engine-and-split-storage.md) | Three-tier SQLite Split-DB (`gitmap.db`, `installation.db`, `repodb/pipeline.db`), singular PascalCase schemas, integer Tesla ID primary keys (`{Table}Id`), reentrant `TxLockManager`, WAL mode, and safe row scanners. |
| **07** | [07-pipeline-diagnostics-and-telemetry.md](07-pipeline-diagnostics-and-telemetry.md) | CI/CD pipeline telemetry (`gitmap pe`, `gitmap pea`), bounded stack trace extraction (5 leading + 20 trailing lines), heatmap telemetry, dynamic ETA sleep sync (`runner-eta.json`), and 4-part RCA framework. |
| **08** | [08-distribution-installers-and-release.md](08-distribution-installers-and-release.md) | Cross-platform execution runners (`run.ps1`, `run.sh`), `run.config.json` driver, NSIS Windows installer, Linux archives (`.tar.gz`), centralized `version.json`, and untouchable CI release pipeline invariant. |

---

## Governance & Hygiene Rules

1. **Strict File Limit:** This directory contains **strictly 10 files** (the 9 domain references above + this `readme.md`). Do not create loose files, session dumps, or uncurated scratch notes in this directory.
2. **Authoritative Cross-Linking:** All specifications in `02-spec/` and subtask plans in `.ai-memory/plans/` must reference these documents directly via relative links.
3. **Continuous Compaction:** When new operational learnings or architectural changes are ratified, update the corresponding domain document directly rather than creating incremental notes.
