# Subtask 02: SQLite File Uniqueness and OS-Aware Indexing Guarantee

- **Task Code:** Task-02
- **Owner:** Worker 02
- **Status:** PENDING
- **Owned Files:**
  - `cli/repodb/repo_db.go`
  - `cli/repodb/repo_file.go`

## Instructions:
1. In `cli/repodb/repo_db.go`:
   - In `InitRepoSchema(ctx context.Context, db *sql.DB) error`:
     - Make `IdxRepoFile_RelativePath` OS-aware:
       - On Windows (`runtime.GOOS == "windows"`): `"CREATE UNIQUE INDEX IF NOT EXISTS IdxRepoFile_RelativePath ON RepoFile(RelativePath COLLATE NOCASE);"`
       - On Unix/other: `"CREATE UNIQUE INDEX IF NOT EXISTS IdxRepoFile_RelativePath ON RepoFile(RelativePath);"`
2. In `cli/repodb/repo_file.go`:
   - Implement `EnsureFileUniqueInDB(ctx context.Context, db *sql.DB, relPath string) (bool, error)`:
     - Check if SQLite already contains a file record with `relPath` (using `COLLATE NOCASE` on Windows, binary on Unix).
     - Return `true` if unique (count == 0), `false` if duplicate exists.
3. Verify `Upsert` method in `RepoFileDbRepo` continues to enforce `ON CONFLICT(RelativePath) DO UPDATE SET ...`.
4. Log actions in SQLite task DB before touching files, mark completed with evidence.
