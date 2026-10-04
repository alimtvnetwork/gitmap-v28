# Subtask 77.1: DetectedProject Foreign Key Constraint RCA & Fix

- **Parent Plan:** [77-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca.md](../../pending/77-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca.md)
- **Spec Reference:** [02-spec/21-app/207-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca/01-architecture-spec.md](../../../../02-spec/21-app/207-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca/01-architecture-spec.md)
- **Status:** Pending
- **Target Area:** `cli/store/project.go`, `cli/store/projecttype.go`, `cli/store/repo_upsert.go`, `cli/scanner/`, `cli/detector/`, `cli/model/`

## Objective

Remediate SQLite `FOREIGN KEY constraint failed (787)` during project detection scanner upserts by guaranteeing parent record existence in `Repository` and `ProjectType` prior to inserting into `DetectedProject`, adding transactional consistency and defensive fallback queries.

## Root Cause Analysis (RCA)

1. SQLite enforces foreign keys (`PRAGMA foreign_keys = ON;`). In `DetectedProject`, `RepoId` references `Repository.Id` and `ProjectTypeId` references `ProjectType.Id`.
2. When scanning workspaces containing multiple projects or submodules, project detection can trigger before the parent `Repository` row transaction commits, or when clone validation fails, causing `RepoId` to be zero, invalid, or uncommitted in the current connection.
3. Newly detected project types or indicators might reference a `ProjectTypeId` that has not yet been seeded into the `ProjectType` reference table.
4. `UpsertDetectedProject` in `cli/store/project.go` attempts path-based fallback when `RepoID <= 0`, but does not verify that `RepoID` actually exists in the database before executing `INSERT INTO DetectedProject`.

## Implementation Details

1. **Pre-Seeding & Dynamic Type Guarantee:**
   - In `cli/store/projecttype.go`, verify `SeedProjectTypes` seeds all supported project types on schema initialization.
   - Add `EnsureProjectTypeExists(db, slug)` or verify `ProjectTypeId` against `ProjectType` table before upserting child projects.
2. **Defensive Foreign Key Validation:**
   - In `cli/store/project.go` (`UpsertDetectedProject`), verify that `p.RepoID` exists via `SELECT Id FROM Repository WHERE Id = ?`.
   - If missing, perform path lookup via `SelectRepoIDByPath(p.RepoPath)` and `SelectRepoIDByPath(p.AbsolutePath)`.
   - If `RepoID` is still absent, create or attach to a verified repository stub, or return a structured `apperror.AppError` (code `E2030`) rather than crashing with unhandled SQLite error 787.
3. **Scanner Pipeline Synchronization:**
   - Ensure the scanner commits repository records and verifies the returned primary key before dispatching detection routines for sub-projects.
   - Maintain strict coding standards: function line caps <= 15 lines, positive boolean variables, and explicit AppError wrapping.

## Verification & Acceptance Criteria

- Executing `gitmap scan` across multi-project and monorepo repositories completes without any `FOREIGN KEY constraint failed (787)` errors.
- Unit tests in `cli/store/` verify project upsert with valid parent `RepoId` and test defensive handling when parent repo is not yet committed.
