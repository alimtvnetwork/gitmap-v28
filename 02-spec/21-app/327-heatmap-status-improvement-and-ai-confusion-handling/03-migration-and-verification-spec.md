# Migration & Verification Spec: Heatmap Status Improvement and AI Confusion Handling

- **Task Sequence:** 327
- **Task Slug:** `heatmap-status-improvement-and-ai-confusion-handling`
- **Spec Path:** `02-spec/21-app/327-heatmap-status-improvement-and-ai-confusion-handling/`
- **Specification Document:** `03-migration-and-verification-spec.md`
- **Status:** `RATIFIED`
- **Created At:** 2026-10-25

---

## 1. Executive Summary & Problem Scope

This specification defines the migration protocols, database verification, pull sanitization architecture, and release verification standards for Task 327.

### 1.1 Root Cause Analysis: Commit Messages Treated as Failed Repositories (Screenshots 2 & 4)

In recent test and production runs, the pull remediation engine and post-pull summary reported phantom repository failures such as:
```text
D:\work\COMMIT - support cfr verb in commit formatter and eliminate magic strings across cli
```
These strings originated when CLI commands, test harnesses, or subagent tasks passed raw semantic commit strings or descriptive task directives where a repository slug or path was expected. When downstream pull workers failed to access these paths as Git repositories, the error was logged directly into the `PullErrorRecord` table of `gitmap-pull.db`.

Subsequent invocations of `RunBatchPullFix` and `convertRecordsToFailures`:
1. Ingested these malformed entries as legitimate tracked repository failures.
2. Formatted them into remediation cards and interactive fix menus.
3. Passed them into `ClassifyNonRepoFolder`, which erroneously attempted to suggest `git init` or `gitmap rm` for synthetic paths.

### 1.2 Core Remediation Directives

To eliminate this confusion, Task 327 introduces:
1. **Candidate Validation Guard (`isValidRepoNameCandidate`):** A strict gate rejecting any string containing newlines, the commit separator `" - "`, conventional commit prefixes (`COMMIT `, `FEAT `, `FIX `, `Feature: `, `Bug: `, `Release: `, `Chore: `), or exceeding 64 characters.
2. **Directory & Remote Validation in Remediation Queries:** Skipping pull error records whose directory does not exist on disk AND has no recorded remote URL.
3. **Non-Repo Classifier Guard (`ClassifyNonRepoFolder`):** Refusing to classify any path or name as a remediable repository if it fails `isValidRepoNameCandidate`.
4. **Pull Batch Boundary Verification (`pull_batch.go`):** Filtering failure records, unresolved queues, and remediation items before persisting or presenting them to the user.
5. **Output Collapsing & Guidance:** Collapsing up-to-date repository listings into a concise single-line count and displaying a clear tip command for granular inspection.

---

## 2. Database Schema Migration & Verification

### 2.1 Root SQLite `RepoCacheSequenceIndex`

To optimize `repo-cache` folder resolution without repetitive filesystem directory walks, the root SQLite database (`gitmap.db`) introduces the `RepoCacheSequenceIndex` table:

```sql
CREATE TABLE IF NOT EXISTS repo_cache_sequence_index (
    repo_slug TEXT PRIMARY KEY,
    sequence_folder TEXT NOT NULL,
    sequence_num INTEGER NOT NULL,
    cache_path TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_repo_cache_seq_folder ON repo_cache_sequence_index(sequence_folder);
```

#### Migration Invariants:
- `repo_slug`: Lowercase, normalized repository slug (e.g. `gitmap-v28`).
- `sequence_folder`: Canonical numbered directory in `repo-cache` (e.g. `01-gitmap-v28`).
- `sequence_num`: Extracted integer sequence prefix (`1`, `2`, ...).
- Fast lookup: `LookupRepoSequenceFolder(slug)` retrieves cached directory mapping in $O(1)$ time; falls back to directory walk and populates cache if missing.

### 2.2 Normalized 1-to-N AI Confusion Schema

The per-repository AI confusion SQLite database located at `repo-cache/XX-<repo>/ai-confusion.db` implements a normalized 1-to-N structure:

```sql
CREATE TABLE IF NOT EXISTS ai_confusion (
    id TEXT PRIMARY KEY,
    agent_slug TEXT NOT NULL,
    task_id TEXT NOT NULL,
    subtask_id TEXT NOT NULL,
    spec_path TEXT NOT NULL,
    reasoning TEXT NOT NULL,
    confusion_state TEXT NOT NULL,
    question_state TEXT NOT NULL,
    is_completed INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ai_confusion_file (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    confusion_id TEXT NOT NULL,
    file_path TEXT NOT NULL,
    start_line INTEGER NOT NULL DEFAULT 0,
    end_line INTEGER NOT NULL DEFAULT 0,
    code_snippet TEXT NOT NULL DEFAULT '',
    FOREIGN KEY(confusion_id) REFERENCES ai_confusion(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_ai_confusion_task ON ai_confusion(task_id, subtask_id);
CREATE INDEX IF NOT EXISTS idx_ai_confusion_file_path ON ai_confusion_file(file_path);
```

---

## 3. Pull Sanitization Architecture

### 3.1 Candidate Validation Specifications

The helper `isValidRepoNameCandidate(name string) bool` enforces the following mathematical rules:

```
                    +---------------------------+
                    |        Input String       |
                    +-------------+-------------+
                                  |
                                  v
                    +---------------------------+
                    | Trim whitespace; len == 0 | ----> INVALID (false)
                    | or len > 64 chars?        |
                    +-------------+-------------+
                                  | No
                                  v
                    +---------------------------+
                    | Contains '\n', '\r',      | ----> INVALID (false)
                    | or " - "?                 |
                    +-------------+-------------+
                                  | No
                                  v
                    +---------------------------+
                    | Starts with disallowed    |
                    | prefix: COMMIT , FEAT ,   | ----> INVALID (false)
                    | FIX , Feature: , Bug: ,   |
                    | Release: , Chore: ?       |
                    +-------------+-------------+
                                  | No
                                  v
                    +---------------------------+
                    |      VALID (true)         |
                    +---------------------------+
```

### 3.2 Leaf Error Sovereignty & Error Bubbling

All sanitization functions act as pure filters or return typed errors:
- No `os.Exit` or CLI termination inside leaf functions.
- If a query encounters database corruption, wrapped `*apperror.AppError` instances bubble up to command dispatchers.

---

## 4. Verification & Testing Matrix

| Component | Target Location | Verification Method | Acceptance Criteria |
| :--- | :--- | :--- | :--- |
| Candidate Filter | `cli/cmdpull/pull_remediation_query.go` | Unit validation of strings with commit prefixes, `" - "`, and $>64$ runes | Rejects `"COMMIT - ..."`, accepts valid slugs like `"gitmap"` |
| Query Conversion | `cli/cmdpull/pull_remediation_query.go` | Ingest records with non-existent paths & no remote URLs | Skips phantom records; returns only authentic repositories |
| Non-Repo Classifier | `cli/cloner/pulldiag.go` | Pass commit strings and non-repo paths to `ClassifyNonRepoFolder` | Returns empty `NonRepoDiagnosis{}` (`IsNonRepoFolder == false`) |
| Pull Batch Protection | `cli/cmdpull/pull_batch.go` | Batch pull execution with synthetic state failures | Prevents malformed slugs from entering remediation queue |
| Output Collapsing | `cli/cmdpull/pull_efficient_render.go` | Execute pull on repos with 0 changes | Shows collapsed `+ N up-to-date repositories (unchanged)` with tip |
| Sequence Cache | `cli/store/repo_cache_sequence_store.go` | Insert and query `RepoCacheSequenceIndex` | Instant resolution of `XX-<repo>` folder names from cache |
| 1-to-N Confusion | `cli/store/ai_confusion_store.go` | Insert confusion record with multiple files and line ranges | All files linked to parent confusion ID with cascade delete |

---

## 5. Release Ceremony Standards (v6.611.0)

Upon successful integration of all Task 327 components:
1. Version synchronization across `package.json`, `version.json`, and constants.
2. Comprehensive release notes in `changelog.md` and `.ai-memory/release/v6.611.0.md`.
3. Dedicated semantic commit using `gitmap cpr "RELEASE - bump version to v6.611.0"`.
