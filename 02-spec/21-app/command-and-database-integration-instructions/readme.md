# Command and Database Integration Instructions

- **Module:** `02-spec/21-app/command-and-database-integration-instructions`
- **Scope:** Native Go SQLite integration, Split-DB routing, CLI command contracts (`db-read`, `db-add`, `db-edit`, `db-find`, `db-view`, `db-listing`)
- **Status:** APPROVED
- **Version:** 1.0.0
- **Target Release:** Minor Bump

---

## 1. Overview & Directory Index

This specification suite defines the native Go SQLite persistence, registry, search, and introspection ecosystem within the GitMap CLI binary. These commands supersede legacy ad-hoc Python one-liners, establishing safe, typed, and low-latency database access across all Split-DB tiers.

| Specification Document | Scope & Focus | Status |
|---|---|:---:|
| **[`01-architecture-overview.md`](01-architecture-overview.md)** | Core architecture, deprecation of Python one-liners, Split-DB compatibility, pure-Go `modernc.org/sqlite` driver, single-writer invariants, and `PRAGMA query_only=ON` safety. | APPROVED |
| **[`02-cli-command-contracts.md`](02-cli-command-contracts.md)** | CLI contracts for `gitmap db-read` (`dbr`), `gitmap db-add` (`dba`), and `gitmap db-edit` (`dbe`), argument tables, error codes (20–28), and JSON envelopes. | APPROVED |
| **[`03-alias-registry-and-storage.md`](03-alias-registry-and-storage.md)** | SQLite Database Registry architecture, persistent alias-to-path mapping in `gitmap.db`, file size detection, and `gitmap db-listing` (`dbl`). | APPROVED |
| **[`04-view-and-join-engine.md`](04-view-and-join-engine.md)** | Specification for `gitmap db-find` across columns and `gitmap db-view` with automatic join inference, alias support, and projections. | APPROVED |
| **[`05-antigravity-verification-and-testing.md`](05-antigravity-verification-and-testing.md)** | End-to-end verification protocol, live Antigravity DB tests, deadlock regression test suite (`MaxOpenConns=1`), and CI integration. | APPROVED |

---

## 2. Architectural Pillars

1. **Native Pure-Go Binary Execution:** All SQLite operations execute inside the statically compiled `gitmap` CLI binary via `modernc.org/sqlite`, completely removing Python runtime requirements.
2. **Strict Read-Only Enforcement:** Read commands (`gitmap db-read` / `dbr`) enforce engine-level safety through `PRAGMA query_only = ON;` and AST token filtering, eliminating accidental data mutation.
3. **Split-DB Alias Resolution:** Universal `-as` / `--db` alias routing provides immediate access to all database tiers (`master`, `install`, `pipeline`, `heatmap`, `tasks`) as well as custom filesystem paths.
4. **Single-Writer Concurrency:** Mutating commands (`db-add`, `db-edit`) enforce `db.SetMaxOpenConns(1)`, `PRAGMA journal_mode = WAL`, and `BEGIN IMMEDIATE` transactions to prevent database lock contention.
5. **Machine & Human Interfaces:** Dual output pipeline supporting interactive Bubble Tea / lipgloss tables for human operators and standardized JSON envelopes (`--json`) for autonomous AI agents.
