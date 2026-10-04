# Subtask 78.2: Comprehensive Database Reset & Reseed Engine

## 1. Context & Objective
The user requires `gitmap reset` / `gitmap db reset` to completely wipe all database files across the database directory (`%LOCALAPPDATA%\gitmap-cli\data` and `.gitmap/`), report every removed file with exact bytes, and reseed clean schemas and default project types so that all previous state is wiped and fresh.

---

## 2. Technical Scope
- `cli/cmddb/cmddb_reset.go`
- `cli/cmd/reset.go`
- `cli/store/storage_inventory.go`

---

## 3. Remediation Checklist
- [ ] Scan and collect all `.db`, `-wal`, `-shm` database files across `%LOCALAPPDATA%\gitmap-cli\data` (master `gitmap.db`, `installation/`, `tasks/`, `pipeline/`, `automation/`, errors, and pull databases).
- [ ] Safely remove all matching files.
- [ ] Print table/list of removed databases with individual file sizes and grand total reclaimed.
- [ ] Re-run migrations (`db.Migrate()`) and project type seeds (`db.SeedProjectTypes()`).
- [ ] Verify execution with `gitmap reset --dry-run` and live `gitmap reset -y`.
