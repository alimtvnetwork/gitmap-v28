# Subtask 01: Specification Suite & Database Registry Engine

- **Subtask ID**: `subtask-01-spec-and-registry-engine`
- **Parent Task**: `command-and-database-integration-instructions`
- **Specification Reference**:
  - `02-spec/21-app/command-and-database-integration-instructions/01-architecture-overview.md`
  - `02-spec/21-app/command-and-database-integration-instructions/03-alias-registry-and-storage.md`
- **Owned Files**:
  - `cli/cmddb/alias_types.go`
  - `cli/cmddb/alias_registry.go`
  - `cli/cmddb/cmddb_listing.go`
  - `cli/cmddb/cmddb_listing_test.go`
- **Assigned Worker**: Worker 01
- **Status**: PENDING

---

## 1. Objectives

1. Maintain modular specification documents in `02-spec/21-app/command-and-database-integration-instructions/` ensuring alignment with Go guidelines and CGO-free architecture.
2. Architect and implement the persistent `DatabaseAlias` table inside GitMap's master database (`gitmap.db`) adhering strictly to PascalCase table naming and singular conventions.
3. Construct the alias resolver engine in `cli/cmddb/alias_registry.go` to support:
   - Dynamic alias registration (`-as <alias>`) upon invoking database commands.
   - Transparent resolution of aliases to physical absolute paths.
   - Real-time disk inspection tracking `FileSizeBytes` and query activity counter `SearchHistoryCount`.
4. Implement the `gitmap db-listing` (alias `gitmap dbl`) CLI command with dual-mode presentation:
   - Rich terminal display using Lipgloss/termout tabular formatting.
   - Structured JSON envelope via the `--json` flag.
5. Provide 100% unit test coverage in `cli/cmddb/cmddb_listing_test.go` verifying registration persistence, path validation, and JSON schemas.

---

## 2. Key Implementation Interfaces

```go
type DatabaseAlias struct {
    DatabaseAliasId    int64      `json:"databaseAliasId"`
    AliasName          string     `json:"aliasName"`
    DatabasePath       string     `json:"databasePath"`
    DatabaseType       string     `json:"databaseType"`
    IsActive           bool       `json:"isActive"`
    FileSizeBytes      int64      `json:"fileSizeBytes"`
    SearchHistoryCount int64      `json:"searchHistoryCount"`
    LastAccessedAt     *time.Time `json:"lastAccessedAt,omitempty"`
}

type AliasRegistry interface {
    Register(aliasName, dbPath, dbType string) (*DatabaseAlias, error)
    Resolve(aliasOrPath string) (string, error)
    RecordAccess(aliasName string) error
    ListAll() ([]DatabaseAlias, error)
}
```

---

## 3. Verification Criteria

1. `go test ./cli/cmddb -run TestAliasRegistry -v` passes cleanly.
2. CLI verification: `gitmap db-listing --json` outputs valid JSON containing all system and user aliases.
3. Verifying non-existent database paths during registration returns clear validation errors without creating dangling registry entries.
