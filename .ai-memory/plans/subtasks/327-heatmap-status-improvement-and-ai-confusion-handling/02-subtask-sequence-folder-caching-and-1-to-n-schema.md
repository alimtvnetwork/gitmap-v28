# Subtask 02: Sequence Folder Caching & 1-to-N AI Confusion Schema

## Status: IN_PROGRESS
- **Subtask Slug:** `02-subtask-sequence-folder-caching-and-1-to-n-schema`
- **Parent Plan:** `.ai-memory/plans/pending/327-heatmap-status-improvement-and-ai-confusion-handling.md`
- **Architecture Spec:** `02-spec/21-app/327-heatmap-status-improvement-and-ai-confusion-handling/01-architecture-spec.md`
- **Target Subsystem:** `cli/store/`
- **Primary Source Files:**
  - `cli/store/repo_cache_sequence_store.go`
  - `cli/store/ai_confusion_types.go`
  - `cli/store/ai_confusion_schema.go`
  - `cli/store/ai_confusion_store.go`

---

## 1. Context & Scope

GitMap requires high-velocity access to repository-specific caches and persistent AI confusion episodes stored in `repo-cache`. Currently:
1. `repo-cache` organizes repositories under sequence-prefixed folders (e.g. `01-gitmap`, `02-scripts-fixer`). Every operation needing the repo cache directory previously walked the entire filesystem tree to discover the folder name matching the repository slug, resulting in repeated disk scans.
2. AI confusion episodes previously stored file references as flattened, denormalized strings. This prevented agents from efficiently querying episodes touching a specific file path or inspecting distinct line ranges and code snippets across multiple affected files.

This subtask implements the data access layer:
1. Root SQLite database caching of repo-cache sequence folders (`RepoCacheSequenceIndex`), providing $O(1)$ mapping from `RepoSlug` to `XX-<repo>` folders in `repo-cache`.
2. Normalized 1-to-N database schema and store routines for `AiConfusionFile` records linked to parent `AiConfusion` episodes with cascading deletion.

---

## 2. Technical Architecture & Schemas

### 2.1 Root SQLite Table: `RepoCacheSequenceIndex`

Created in GitMap's root SQLite database (via `store.OpenInstallationDB()` or `store.OpenRootHeatmapDB()`):

```sql
CREATE TABLE IF NOT EXISTS RepoCacheSequenceIndex (
    RepoCacheSequenceIndexId INTEGER PRIMARY KEY AUTOINCREMENT,
    RepoSlug                 TEXT NOT NULL UNIQUE,
    SequenceFolder           TEXT NOT NULL,
    ResolvedPath             TEXT NOT NULL,
    LastVerifiedAt           TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CreatedAt                TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt                TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS IdxRepoCacheSeq_Slug ON RepoCacheSequenceIndex(RepoSlug);
CREATE INDEX IF NOT EXISTS IdxRepoCacheSeq_Folder ON RepoCacheSequenceIndex(SequenceFolder);
```

### 2.2 Sequence Cache Store API (`cli/store/repo_cache_sequence_store.go`)

The store provides four core functions:

1. `GetCachedSequenceFolder(db *sql.DB, repoSlug string) (string, bool, error)`
   - Queries `RepoCacheSequenceIndex` for `RepoSlug`.
   - Returns cached `SequenceFolder` and `hasMatch = true` if present.
2. `SetCachedSequenceFolder(db *sql.DB, repoSlug, sequenceFolder, resolvedPath string) error`
   - Idempotently upserts `RepoSlug`, `SequenceFolder`, and `ResolvedPath`.
3. `ScanAndCacheRepoCacheSequenceFolders(db *sql.DB, repoCacheDir string) error`
   - Scans directory entries in `repoCacheDir` matching sequence pattern `^(\d{2,3})-(.+)$`.
   - Upserts all discovered mappings into the SQLite index in a single transaction.
4. `ResolveRepoCacheSequenceFolder(db *sql.DB, repoCacheDir, repoSlug string) (string, error)`
   - Checks the SQLite cache first.
   - On miss, falls back to targeted directory matching, caches the result, and returns the resolved path.

### 2.3 Normalized 1-to-N AI Confusion Schema (`cli/store/ai_confusion_schema.go`)

```sql
CREATE TABLE IF NOT EXISTS AiConfusionFile (
    AiConfusionFileId INTEGER PRIMARY KEY AUTOINCREMENT,
    AiConfusionId     INTEGER NOT NULL,
    ConfusionUuid     TEXT NOT NULL,
    FilePath          TEXT NOT NULL,
    FileType          TEXT NOT NULL DEFAULT 'SOURCE',
    LineNumbers       TEXT NOT NULL DEFAULT '',
    CodeSnippet       TEXT NULL,
    WhyModified       TEXT NOT NULL DEFAULT '',
    CreatedAt         TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (AiConfusionId) REFERENCES AiConfusion(AiConfusionId) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS IdxAiConfusionFile_ConfusionId ON AiConfusionFile(AiConfusionId);
CREATE INDEX IF NOT EXISTS IdxAiConfusionFile_FilePath ON AiConfusionFile(FilePath);
```

### 2.4 Store Persistence Extensions (`cli/store/ai_confusion_store.go`)

1. `InsertConfusionFiles(tx *sql.Tx, confusionId int64, uuid string, files []AiConfusionFileRecord) error`
   - Batch-inserts child records inside the parent transaction.
2. `GetConfusionFilesByConfusionId(confusionId int64) ([]AiConfusionFileRecord, error)`
   - Fetches all referenced files associated with a given `AiConfusionId`.
3. `GetConfusionByFilePath(filePath string) ([]AiConfusionRecord, error)`
   - Performs an inner join between `AiConfusion` and `AiConfusionFile` matching `f.FilePath = ?`.
   - Enriches each parent record with its complete set of `ReferencedFiles`.

---

## 3. Implementation Checklist

- [x] **Phase 1: Types & Schema Foundation**
  - [x] Define `AiConfusionFileRecord` in `cli/store/ai_confusion_types.go`.
  - [x] Add `ReferencedFiles []AiConfusionFileRecord` to `AiConfusionRecord`.
  - [x] Define `AiConfusionFileInput` and add `ReferencedFiles []AiConfusionFileInput` to `AiConfusionJSONInput`.
  - [x] Add `SQLCreateAiConfusionFileTable` and index DDL to `cli/store/ai_confusion_schema.go`.
  - [x] Register DDL statements in `InitAiConfusionSchema`.
- [ ] **Phase 2: Sequence Folder SQLite Cache**
  - [ ] Implement `RepoCacheSequenceIndex` DDL and initialization in `cli/store/repo_cache_sequence_store.go`.
  - [ ] Implement `GetCachedSequenceFolder`, `SetCachedSequenceFolder`, and `ScanAndCacheRepoCacheSequenceFolders`.
  - [ ] Integrate sequence folder cache with `ResolveAiConfusionDbPath`.
- [ ] **Phase 3: Store Persistence & Query Methods**
  - [ ] Update `InsertAiConfusion` in `cli/store/ai_confusion_store.go` to persist `ReferencedFiles` within a transaction.
  - [ ] Implement `GetConfusionFilesByConfusionId` in `cli/store/ai_confusion_store.go`.
  - [ ] Implement `GetConfusionByFilePath` to enable relative file path resolution for `view`.

---

## 4. Strict Coding Rules & Invariants

1. **Leaf Error Sovereignty:** All store functions return direct `error` or `*apperror.AppError`. Zero exit calls.
2. **15-Line Function Limit:** Decompose complex multi-column SQL scans and transaction loops into specialized helper functions $\le 15$ lines.
3. **Zero Nested Ifs:** Flatten all conditional branching with early returns and guard clauses.
4. **Positive Booleans:** Boolean identifiers must start with `is*` or `has*` only.
5. **Whitespace Discipline:** Blank line before `if`, blank line after `}`, and blank line before `return`.
6. **No-Build / No-Test Rule:** No builds or tests executed during editing turns.

---

## 5. Acceptance Criteria

- [x] `AiConfusionFileRecord` and `AiConfusionFileInput` defined and integrated in `cli/store/ai_confusion_types.go`.
- [x] `AiConfusionFile` table and indexes initialized via `InitAiConfusionSchema` in `cli/store/ai_confusion_schema.go`.
- [ ] `RepoCacheSequenceIndex` table and querying routines operational in `cli/store/repo_cache_sequence_store.go`.
- [ ] Child file records persisted and retrieved alongside parent `AiConfusionRecord` episodes.
- [ ] All file paths are strictly relative to repository root.
