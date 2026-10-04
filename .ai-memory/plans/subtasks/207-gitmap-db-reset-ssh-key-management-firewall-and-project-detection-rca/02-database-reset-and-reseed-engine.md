# Subtask 77.2: Comprehensive Database Reset & Reseed Engine

- **Parent Plan:** [77-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca.md](../../pending/77-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca.md)
- **Spec Reference:** [02-spec/21-app/207-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca/01-architecture-spec.md](../../../../02-spec/21-app/207-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca/01-architecture-spec.md)
- **Status:** Pending
- **Target Area:** `cli/cmd/dbreset.go`, `cli/cmddb/cmddb_reset.go`, `cli/cmddb/exports.go`, `cli/cmd/rootcore.go`, `cli/store/store.go`

## Objective

Deliver a complete database purge and baseline reseed engine via `gitmap reset` and `gitmap db reset` that safely closes all active SQLite connections, deletes all database and WAL/SHM companion files across `.gitmap/`, outputs human-readable per-file deletion metrics, and re-executes migrations and baseline seeds.

## Architecture & Requirements

1. **CLI Commands & Flags:**
   - Primary: `gitmap reset`
   - Aliases: `gitmap db reset`, `gitmap database reset`
   - Flags:
     - `-y`, `--yes`, `--confirm`: Confirm deletion without interactive prompt.
     - `--dry-run`, `-n`: Preview all target databases and space that would be reclaimed without deleting anything.
     - `--rescan`: Automatically invoke `gitmap scan` after schemas are reseeded.
2. **Safe Discovery & Purge:**
   - Discover all SQLite files: `gitmap.db`, `installation.db`, `repodb/*.db`, `templates.db`, `pipeline.db`, and companion files (`-wal`, `-shm`, `.wal`, `.shm`).
   - Close active connection pools before deletion to prevent `EBUSY` / Windows file lock errors.
   - Remove files and log each purged database:
     `✔ Removed database: .gitmap/gitmap.db (1.2 MB)`
   - Output summary line: `✔ Total N database(s) purged (X MB reclaimed).`
3. **Clean Reseed Pipeline:**
   - Reopen connection pool using `store.OpenDefault()`.
   - Run baseline migrations (`db.Migrate()`).
   - Reseed core lookup dictionaries (`SeedProjectTypes()`, default settings).
   - Print confirmation: `✔ All databases reseeded with clean baseline schemas.`

## Verification & Acceptance Criteria

- Executing `gitmap reset --dry-run` displays candidate database paths and estimated reclaimed bytes without file deletion.
- Executing `gitmap reset -y` deletes all SQLite database files and outputs clean reseed confirmation.
- Running `gitmap reset` in an interactive terminal prompts `Are you sure you want to reset all databases? [y/N]: ` before proceeding.
- Function line counts strictly adhere to <= 15 lines.
