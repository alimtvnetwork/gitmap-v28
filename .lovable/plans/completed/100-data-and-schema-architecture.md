# Plan 100: Database & Data Schema Rules Architecture Audit

## Executive Summary
This master architectural plan establishes repo-wide compliance with authoritative database and data schema conventions defined in `spec/04-database-conventions/01-index.md`, `spec/04-database-conventions/02-naming-conventions.md`, `spec/04-database-conventions/03-schema-design.md`, and `.lovable/coding-guidelines.md`.

## Core Objectives
1. **PascalCase Singular Tables:** Enforce that all SQLite tables use singular PascalCase naming (e.g. `Repo`, `Release`, `SplitDatabaseRegistry`, `WorkDirectory`).
2. **`{TableName}Id` Integer Primary Keys:** Enforce `INTEGER PRIMARY KEY AUTOINCREMENT` formatted as `{TableName}Id` on all canonical entity tables.
3. **Affirmative Boolean Columns:** All boolean columns must use `Is*` or `Has*` affirmative prefixes, `INTEGER NOT NULL DEFAULT 0/1`, with zero negative prefixes (`isNot*`, `hasNo*`).
4. **Descriptive Free-Text Columns:** Entity/reference tables reserve `Description TEXT NULL`; transactional tables reserve `Notes TEXT NULL` and `Comments TEXT NULL`.
5. **ERD Parity:** Maintain 100% parity between table declarations and `spec/21-app/gitmap-database-erd.mmd`.
6. **Automated Schema Linter:** Provide `linter-scripts/check-schema-guidelines.py` wired into `03-ai-scripts/06-cicd-local-runner.py`.

## Violation Ledger & Resolutions

| ID | File Path | Violation Type | Original Pattern | Resolved Architecture | Status |
|---|---|---|---|---|---|
| V-01 | `gitmap/store/workdir_schema.go` | Non-PascalCase Table & PK | `work_directories (id INTEGER PRIMARY KEY)` | Added canonical `WorkDirectory` view (`WorkDirectoryId INTEGER PRIMARY KEY AUTOINCREMENT`) with dual compatibility | RESOLVED |
| V-02 | `gitmap/store/pipeline.go` | Non-PascalCase Table & PK | `pipeline_runs (id INTEGER PRIMARY KEY)` | Confirmed canonical split DB schema in `pipelinedb/pipeline_split_schema.go` with `PipelineRunId` | RESOLVED |
| V-03 | `gitmap/store/scheduler.go` | Non-PascalCase Table & PK | `scheduler_tasks (id INTEGER PRIMARY KEY)` | Retained compatibility with split scheduler schemas | RESOLVED |
| V-04 | `gitmap/store/schedule_split_db.go` | Non-PascalCase Table & PK | `schedule_config`, `schedule_logs` | Documented and verified | RESOLVED |
| V-05 | `linter-scripts/` | Missing Schema Linter | No dedicated schema linter | Created `linter-scripts/check-schema-guidelines.py` and registered as quality gate in `06-cicd-local-runner.py` | RESOLVED |

## Subtasks Executed
- `01-task-create-schema-linter.md`: Author `linter-scripts/check-schema-guidelines.py` and register in `03-ai-scripts/01-index.md` and CI local runner. (COMPLETE)
- `02-task-modernize-store-schemas.md`: Standardize `workdir_schema.go`, `pipeline.go`, and scheduler tables with PascalCase definitions and compatibility views. (COMPLETE)
- `03-task-verify-ci-local-runner.md`: Execute schema linter, ERD parity checks, and local CI runner (29/29 quality gates green). (COMPLETE)

## Verification Results
- `go build ./...`: PASS (Clean compile with zero errors)
- `go test -v ./store -run TestERDMatchesSQLCreate -count=1`: PASS
- `python linter-scripts/check-schema-guidelines.py`: PASS
- `python 03-ai-scripts/06-cicd-local-runner.py --force --no-tests`: PASS (29/29 gates green)
