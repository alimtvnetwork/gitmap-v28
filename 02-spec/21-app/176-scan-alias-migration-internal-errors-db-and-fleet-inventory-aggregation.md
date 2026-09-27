# Canonical Specification: Scan Alias Migration, Internal Errors DB, and Fleet Inventory Aggregation

## Specification Metadata
- **Spec ID:** 176-scan-alias-migration-internal-errors-db-and-fleet-inventory-aggregation
- **Status:** Active
- **Category:** Architecture & Split-DB / Fleet Telemetry
- **Created:** 2026-09-27
- **Target Release:** v6.355.3

---

## 1. User Request (Verbatim)

```text
Okay. So I want you to run Git map scan on the work directory and take the JSON file, gitmap.json file, which actually allows to circle back or, let's say, scan on this current machine. So current machine I can run. Oh, I do see the error. There is a error. Consider fixing some error as well. And also, I think we have discussed this in the past. So any error that happens that should be logged inside the SQLite DB, errors DB, the general error, and we can do a errors E or errors that would give the internal errors that has been saved. I don't think that you considered this, so that needs to be fixed, I do think. And okay, so the JSON file is saved properly. The first thing is to consider saving this JSON file, because this JSON file should contain the latest repo, because a lot of repos has been added. That's all right. But I want to have the other machines scan to be done from the current machine and get that JSON file here and move it to the repo secrets folder with their work areas. So the first one is the work-areas one, right? So that would go as 04-W1 machine. And then the 05W2 machine would go, then 06, the work W3 machine would go. I have given a format for W2. You can check. So similar to this, you need to just clone and find these files and put into those folders, create folders and keep it there. Can you please do that for me? Yeah. Do you understand the task, and is it doable for you? Confirm it.
```

---

## 2. Architecture & Design

### 2.1 Problem Analysis
1. **Alias Column Missing Failure:**
   - When running `gitmap scan`, `autoPopulateScanAliases` invokes `cmdinstall.EnsureTrackedRepoAliases(db)` which calls `store.CreateAliasWithDetails`.
   - `CreateAliasWithDetails` executes `INSERT INTO Alias (Alias, RepoId, IsPrimary, Source) VALUES (?, ?, ?, ?)`.
   - In older databases, the `Alias` table lacked `IsPrimary` and `Source`. Because `Setting[schema_version] == 31`, `db.Migrate()` short-circuited via `isSchemaUpToDate()`, skipping `db.migrateAliasColumns()`.
   - Solution: Bump `SchemaVersionCurrent = 32` and introduce defensive `ensureAliasColumns()` in `cli/store/alias.go` that inspects column existence via `PRAGMA table_info` and applies `ALTER TABLE` dynamically.

2. **Internal Errors Split-DB & Inspection Command:**
   - Previous behavior: Errors in `ExecWrapper` and `QueryWrapper` were only logged to `os.Stderr`. CLI errors were only dumped to `.gitmap/last_error.log`.
   - Solution: Implement a dedicated SQLite Split-DB `gitmap-errors.db` with table `InternalErrorLog`.
   - Provide `gitmap errors` and `gitmap e` CLI commands supporting listing, single-record detail view (`gitmap e <id>`), clearing (`gitmap e clear`), and JSON output (`gitmap e --json`).

3. **Fleet Inventory Aggregation in `D:\work\repo-secrets`:**
   - SSH nodes: `w1` (192.168.1.3), `w2` (192.168.1.7), `w3` (192.168.1.12).
   - Target folders:
     - `D:\work\repo-secrets\04-w1-machine\`
     - `D:\work\repo-secrets\05-w2-machine\` (existing reference)
     - `D:\work\repo-secrets\06-w3-machine\`
   - Required files: `gitmap.json`, `gitmap-ssh-nodes.json`, `gitmap-ssh.json`, `ooshutup10.cfg` (if present).

---

## 3. Data Contracts & Schema

### InternalErrorLog Table Schema
```sql
CREATE TABLE IF NOT EXISTS InternalErrorLog (
    InternalErrorLogId INTEGER PRIMARY KEY AUTOINCREMENT,
    ErrorCode          TEXT NOT NULL DEFAULT '',
    ErrorType          TEXT NOT NULL DEFAULT 'general',
    Command            TEXT NOT NULL DEFAULT '',
    Message            TEXT NOT NULL,
    Details            TEXT NULL,
    SourceFile         TEXT NULL,
    ContextJson        TEXT NULL,
    StackTrace         TEXT NULL,
    GitMapVersion      TEXT NOT NULL DEFAULT '',
    IsResolved         INTEGER NOT NULL DEFAULT 0,
    Notes              TEXT NULL,
    Comments           TEXT NULL,
    CreatedAt          TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS IdxInternalErrorLog_CreatedAt ON InternalErrorLog(CreatedAt DESC);
CREATE INDEX IF NOT EXISTS IdxInternalErrorLog_ErrorType ON InternalErrorLog(ErrorType);
CREATE INDEX IF NOT EXISTS IdxInternalErrorLog_ErrorCode ON InternalErrorLog(ErrorCode);
CREATE INDEX IF NOT EXISTS IdxInternalErrorLog_IsResolved ON InternalErrorLog(IsResolved);
```

---

## 4. Acceptance Criteria & Verification Gates

- [x] **AC-01:** `gitmap scan D:\work --output json` completes with ZERO `SQL logic error: table Alias has no column named IsPrimary` errors.
- [x] **AC-02:** `gitmap e` and `gitmap errors` display recorded internal errors or a clean state message.
- [x] **AC-03:** `gitmap e <id>` displays full detail card with code, type, stack trace, and context.
- [x] **AC-04:** `gitmap e clear` flushes all recorded errors from `gitmap-errors.db`.
- [x] **AC-05:** `D:\work\.gitmap\output\gitmap.json` contains fresh inventory of all local repositories.
- [ ] **AC-06:** Remote nodes W1, W2, and W3 are scanned via SSH and their inventories retrieved.
- [ ] **AC-07:** `D:\work\repo-secrets\04-w1-machine` and `D:\work\repo-secrets\06-w3-machine` are populated matching `05-w2-machine` format.
