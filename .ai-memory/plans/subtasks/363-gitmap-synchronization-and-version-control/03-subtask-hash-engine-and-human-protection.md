# Subtask 03: Non-Destructive Hash Verification & Human Conflict Algorithm

- **Parent Plan:** [363-gitmap-synchronization-and-version-control](../../pending/363-gitmap-synchronization-and-version-control.md)
- **Subtask ID:** `TASK-03`
- **Slug:** `gitmap-synchronization-and-version-control - Task 03`
- **Target Package:** `cli/cmdsync/`, `cli/cmdcg/`, `pkg/hashengine/`
- **Status:** `pending`
- **Assigned Worker:** Security & Sync Subagent

---

## 1. Goal

Implement the non-destructive 3-way hash verification engine and human modification conflict protection algorithm. Ensure that automated fleet synchronizations never silently overwrite manual developer changes, local branch customizations, or uncommitted work.

---

## 2. In Scope

1. **3-Way Hash State Machine (`hash_verifier.go`):**
   - For every managed asset, evaluate three discrete hashes:
     - `CanonicalHash`: Digest of the authoritative asset in canonical source repo.
     - `BaseHash`: Digest recorded in target repo's `version.json` at the previous sync.
     - `DiskHash`: Digest of the actual file currently residing on the target filesystem.
   - Classification rules:
     - `IN_SYNC`: `DiskHash == CanonicalHash` (No action required).
     - `OUTDATED`: `DiskHash == BaseHash` AND `DiskHash != CanonicalHash` (Pristine target, safe to update automatically).
     - `HUMAN_MODIFIED`: `DiskHash != BaseHash` AND `DiskHash != CanonicalHash` (File was edited locally by a human or custom branch).
     - `NEW_ASSET`: `BaseHash == ""` AND file absent on disk (New canonical asset to install).
     - `ORPHANED_LOCAL`: `CanonicalHash == ""` AND file present on disk (Asset deprecated upstream).

2. **Human Conflict Protection Protocol (`conflict_guard.go`):**
   - Strict Default Policy: If state is `HUMAN_MODIFIED`, abort modification for that file.
   - Reporting: Add entry to `ConflictDocket` containing relative path, modification timestamp, and diff summary.
   - Backup Handling (`--backup` flag): Write non-destructive `<path>.bak.<timestamp>` sidecar before any action.
   - Force Override (`--force` flag): Allow deliberate override when explicitly commanded by the operator, logging the override to GitMap's audit trail.
   - Incoming Patch Mode (`--sidecar` flag): Write the new canonical version as `<path>.incoming` to allow manual visual inspection.

3. **Batch Sync Verification Orchestrator (`sync_resolver.go`):**
   - Run comparison across all discovered files in memory before performing any filesystem mutations.
   - Calculate summary metrics: total clean updates, skipped human modifications, newly added assets, and orphan detections.
   - Return strongly-typed `Result[VerificationBatch]` container.

---

## 3. Out of Scope

- CLI flag parsing and root command wiring (deferred to Task 04).
- Terminal card UI formatting and color rendering (deferred to Task 07).
- Version releases or Git commits.

---

## 4. Architectural Invariants

- **Zero Silent Overwrite:** It is an architectural failure if any file with local human modifications is overwritten without `--force`.
- **Positive Booleans:** `isCleanUpdate`, `isHumanModified`, `hasConflict`, `hasBackupRequested`.
- **Leaf Error Sovereignty:** Disk read/stat failures bubble upward via `*apperror.AppError`.
- **Function Line Limit:** Complex state evaluation decomposed into discrete predicate functions <= 15 lines.

---

## 5. Verification Checklist

- [ ] 3-way hash evaluator accurately identifies `IN_SYNC`, `OUTDATED`, and `HUMAN_MODIFIED` states.
- [ ] Modified local files are strictly preserved by default.
- [ ] `--backup` generates valid `<file>.bak` sidecar files before modifications.
- [ ] `--force` successfully updates human-modified files while logging clear audit warnings.
- [ ] Verification batch accurately tabulates summary statistics before disk writes.
