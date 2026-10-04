# App Issue 67: DetectedProject Foreign Key Constraint (787) During Repository Scan RCA

**Issue ID:** 67  
**Date:** 2026-10-04  
**Status:** Resolved  
**Affected Subsystem:** `cli/cmdscan` (`scanprojects.go`), `cli/store` (`project.go`, `repo.go`, `store.go`)  
**Reference Commits:** `aa2dcfa00c85daeceacccb3471dc469ac4618857`, `b54d1b576e6c0df2008172dfbb7bbfe1e06a503b`  

---

## 1. Reproduction

When executing `gitmap scan <directory>` across multi-repo workspaces (for example, 37 repositories with 53 detected projects):

```text
gitmap: json: wrote 37 record(s), 0 validation issue(s)
gitmap: json: wrote 37 record(s), 0 validation issue(s)
  • Cache       last-scan.json

■ Database
────────────────────────────────────────────
  ✔ 37 repositories upserted into database
  • Tagged 37 repo(s) with scan folder #1
  ✔ Auto-generated 79 repository alias(es)
  ↪ background probe queued for 37 repo(s) (workers=3)

■ Project Detection
────────────────────────────────────────────
  [nav] Detected 53 project(s) across 35 repo(s)
  - go-projects.json       15 record(s)
  - react-projects.json    31 record(s)
  - node-projects.json     7 record(s)
[QueryWrapper Error]: exec failed: constraint failed: FOREIGN KEY constraint failed (787)
query: INSERT INTO DetectedProject
        (RepoId, ProjectTypeId, ProjectName, AbsolutePath, RepoPath, RelativePath, PrimaryIndicator)
        VALUES (?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(RepoId, ProjectTypeId, RelativePath) DO UPDATE SET
                ProjectName=excluded.ProjectName,
                AbsolutePath=excluded.AbsolutePath,
                RepoPath=excluded.RepoPath,
                PrimaryIndicator=excluded.PrimaryIndicator,
                DetectedAt=CURRENT_TIMESTAMP
failed to upsert detected project: constraint failed: FOREIGN KEY constraint failed (787)
... (8 errors total)
  [OK] Saved 45 detected project(s) to database
```

---

## 2. 4-Part Root Cause Analysis (RCA)

### 1. Symptom
During project detection upserts, 8 out of 53 detected projects fail with `[QueryWrapper Error]: exec failed: constraint failed: FOREIGN KEY constraint failed (787)` on `INSERT INTO DetectedProject`.

### 2. Root Cause
1. **Windows Path Separator Mismatch in Prefix Matcher:** In `cli/store/repo.go`, `SelectRepoIDByPath` executed `SELECT RepoId FROM Repo WHERE ? LIKE (AbsolutePath || '%')`. On Windows, `Repo.AbsolutePath` stored paths with backslashes (`\`), while detected projects and subtrees were normalized with forward slashes (`/`), causing the `LIKE` prefix match to fail and resolve `RepoId = 0`.
2. **Missing Pre-Flight Foreign Key Existence Check:** In `cli/store/project.go` and `cli/cmdscan/scanprojects.go`, if `p.RepoID > 0` was present (e.g. from memory index mismatch or pre-reset database IDs), the code bypassed path lookups and directly executed `INSERT INTO DetectedProject` without confirming `RepoExists(p.RepoID)`.
3. **Orphan Child Table References Across Migrations:** Child tables retained stale `RepoId` rows across database resets and re-creations.

### 3. Resolution
1. **Slash-Agnostic SQL Path Resolution:** In `cli/store/repo.go`, updated `SelectRepoIDByPath` with `queryRepoIDExact` and `queryRepoIDPrefix` using `REPLACE(AbsolutePath, '\\', '/') = ? COLLATE NOCASE` and `? LIKE (REPLACE(AbsolutePath, '\\', '/') || '%')`.
2. **Pre-Flight FK Validation:** In `cli/store/project.go`, introduced `validateProjectForeignKeys` checking `RepoExists(p.RepoID)` and `ensureProjectTypeExists(p.ProjectTypeID)`. In `cli/cmdscan/scanprojects.go`, added pre-insert guard:
   ```go
   if r.Project.RepoID <= 0 || !db.RepoExists(r.Project.RepoID) {
       return false
   }
   ```
3. **Orphan Purging:** In `cli/store/store.go`, added `purgeOrphanRepoReferences()` during `db.Migrate()`.

### 4. Prevention & Learnings
- Always verify parent record existence (`RepoExists`) before foreign key inserts; never assume in-memory IDs are persisted.
- Normalize path separators (`REPLACE(..., '\\', '/')`) in SQLite SQL queries when supporting Windows and POSIX.
- Purge orphan foreign key references during migration to avoid downstream cascading constraint violations.

---

## 3. Verification

- `go test ./store/... -v`: All tests passed (`TestSelectRepoIDByPath_SlashAgnostic`, `TestRepoExists`, `TestUpsertDetectedProject_ForeignKeyValidation`, `TestPurgeOrphanRepoReferences`).
- Live scan run `gitmap scan D:\work --compact`: Upserted 77 repos, detected 95 projects across 67 repos, and cleanly saved all 95 projects with 0 errors.
