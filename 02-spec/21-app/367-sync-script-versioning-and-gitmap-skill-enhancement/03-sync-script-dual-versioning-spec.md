# Specification 367.03: Sync Script Dual Versioning Architecture

- **Sequence ID:** `367`
- **Subtask ID:** `TASK-03`
- **Feature Slug:** `sync-script-versioning-and-gitmap-skill-enhancement - Task 03`

---

## 1. Objective

Ensure that the multi-repository synchronization toolchain provides two strictly differentiated execution paths across both `gitmap` and `coding-guidelines`:
1. **V7 Execution Path:**
   - Dedicated script: `03-ai-scripts/38-sync-prompts-skills-scripts-v7.py` and `.agents/scripts/38-sync-prompts-skills-scripts-v7.py`.
   - CLI flags: `--v7` or `--mode v7`.
   - Enforces `validate_v7_restrictions`:
     - Checks target timezone (`--tz` or local).
     - Checks calendar weekday vs weekend:
       - Weekday (Monday–Friday): Commits strictly restricted to morning window (< 08:00 AM). Blocked if at or after 08:00 AM unless `--auto-window` or `--schedule` is supplied.
       - Weekend (Saturday–Sunday): Unconstrained 24-hour commits permitted.
     - Checks AI bot identity attribution in Git configuration.
2. **V6 Execution Path:**
   - Base script: `03-ai-scripts/38-sync-prompts-skills-scripts.py` and `.agents/scripts/38-sync-prompts-skills-scripts.py`.
   - CLI flag: `--mode v6` (default).
   - Runs unrestricted without calendar restriction checks.

---

## 2. Invariants
- Dual scripts must be maintained in sync between `gitmap` and `coding-guidelines`.
- Both `.agents/scripts/` and `03-ai-scripts/` must contain the identical scripts.
