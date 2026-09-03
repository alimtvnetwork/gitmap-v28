# Pipeline and Repository Split Database Architecture

## Purpose

Defines the isolated split SQLite database architecture for pipeline telemetry and repository-specific file index storage, decoupling active CI/CD execution records and search indexes from the master database (`gitmap.db`).

## Scoring & Compliance

- **Status:** Production-Ready
- **Health Score:** 100/100
- **Ambiguity Level:** None (All schemas, paths, and commands explicitly declared)
- **Database Standard:** PascalCase, Singular Table Names, `{TableName}Id INTEGER PRIMARY KEY AUTOINCREMENT`, SetMaxOpenConns(1), filepath.EvalSymlinks binary anchoring.

## Specification Inventory

- [`01-pipeline-split-db-architecture.md`](01-pipeline-split-db-architecture.md) — Isolated pipeline database schema, storage paths, and run/error lifecycle.
- [`02-repo-split-db-commands.md`](02-repo-split-db-commands.md) — Repository split DB commands (`gitmap repo db <status|log|errorlogs|clear|reset|optimize|help>`) and `RepoScanLog` table.
- [`03-unified-db-status-and-optimization.md`](03-unified-db-status-and-optimization.md) — Unified database status and cross-database vacuum optimization (`gitmap db <status|optimize>`).
- [`97-acceptance-criteria.md`](97-acceptance-criteria.md) — Gherkin Given/When/Then validation scenarios.
- [`99-consistency-report.md`](99-consistency-report.md) — Structural consistency audit and invariant verification.
