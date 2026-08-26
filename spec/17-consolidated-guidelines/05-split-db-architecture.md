# Consolidated: Split Database Architecture — Complete Reference

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Purpose

This is the **standalone consolidated reference** for the Split Database Architecture. An AI reading only this file must be able to implement the full pattern — including DB registry, dynamic creation, lifecycle management, import/export, RBAC, and reset API — without consulting source specs.

**Source:** `spec/05-split-db-architecture/` (01-fundamentals + 02-features/01-05)

---

## CRITICAL: Naming Convention

**All field names use PascalCase. No underscores allowed.**

| ❌ Wrong | ✅ Correct |
|----------|-----------|
| `session_id` | `SessionId` |
| `created_at` | `CreatedAt` |
| `message_count` | `MessageCount` |

---

## 1. Pattern Summary

A hierarchical SQLite database organization where a **Root DB** manages metadata about child databases. Item-specific databases are created **dynamically on-demand** per entity. This enables data isolation, improved performance, logical organization, and easy import/export via zip files.

### Database Terminology

| Term | Meaning | Example Path |
|------|---------|--------------|
| **Root DB** | Global registry, settings, app list | `data/root.db` |
| **App DB** | Application-scoped metadata | `data/{appName}/search.db` |
| **Session DB** | Per-session isolated storage | `data/{appName}/ai/chat/001-{id}.db` |
| **Cache DB** | Cached results with TTL | `data/{appName}/rag/cache/001-{slug}.db` |
| **Document DB** | RAG chunks + embeddings | `data/{appName}/rag/documents/001-{id}.db` |

### Benefits

| Benefit | Description |
|---------|-------------|
| **Isolation** | Each entity has its own database — no table bloat |
| **Performance** | Smaller databases = faster queries, per-DB WAL |
| **Portability** | Zip any project folder for backup/migration |
| **Testability** | In-memory per domain in tests |
| **Cleanup** | Easy to archive/delete unused databases |
| **Debugging** | Inspect specific databases in isolation |

---

## 2. Hierarchical Structure

### 2-Layer (Simple)

```
data/
├── root.db                    # Root registry
└── {project-slug}/
    ├── config.db              # Project config
    ├── cache.db               # Project cache
    └── logs.db                # Project logs
```

### 3-Layer (Standard — Most Common)

```
data/
├── root.db                    # Root registry
├── {project-slug}/
│   ├── history/               # History databases
│   │   └── {file-slug}.db
│   ├── cache/                 # Cache databases
│   │   └── search-cache.db
│   ├── config/                # Config databases
│   │   └── settings.db
│   ├── chat/                  # Chat sessions
│   │   └── {session-id}.db
│   └── search/                # Search indices
│       └── {index-id}.db
└── {project-slug-2}/
    └── ...
```

### 4-Layer (Complex — With Categories)

```
data/
├── root.db
├── {project-slug}/
│   ├── ai/                    # AI category
│   │   ├── chat/              # Chat type
│   │   │   └── {session-id}.db
│   │   └── embeddings/        # Embeddings type
│   │       └── {model-id}.db
│   ├── workflow/              # Workflow category
│   │   ├── history/
│   │   │   └── {file-slug}.db
│   │   └── queue/
│   │       └── {queue-id}.db
│   └── search/                # Search category
│       ├── indices/
│       │   └── {index-id}.db
│       └── cache/
│           └── {query-hash}.db
```

### File Path Convention

| Component | Pattern | Example |
|-----------|---------|---------|
| Root DB | `{data}/root.db` | `data/root.db` |
| Project Folder | `{data}/{project-slug}/` | `data/my-project/` |
| Type Folder | `{project}/{type}/` | `data/my-project/history/` |
| Entity DB | `{type}/{entity-slug}.db` | `data/my-project/history/readme-md.db` |

---

## 3. Root Database Schema

### Table: Project

```sql
-- linter-waive: MISSING-DESC-001 reason="Consolidated split-DB example; canonical version in 05-split-db-architecture/"
CREATE TABLE Project (
    ProjectId INTEGER PRIMARY KEY AUTOINCREMENT,
    Slug TEXT UNIQUE NOT NULL,
    DisplayName TEXT NOT NULL,
    Path TEXT NOT NULL,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    Status TEXT DEFAULT 'active'
);

CREATE INDEX IdxProject_Slug ON Project(Slug);
CREATE INDEX IdxProject_Status ON Project(Status);
```

### Table: Database

```sql
-- linter-waive: MISSING-DESC-001 reason="Consolidated split-DB example; canonical version in 05-split-db-architecture/"
CREATE TABLE Database (
    DatabaseId INTEGER PRIMARY KEY AUTOINCREMENT,
    ProjectId INTEGER NOT NULL,
    Type TEXT NOT NULL,
    EntityId TEXT,
    Path TEXT NOT NULL,
    SizeBytes INTEGER DEFAULT 0,
    RecordCount INTEGER DEFAULT 0,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    LastAccessedAt DATETIME,
    Status TEXT DEFAULT 'active',
    FOREIGN KEY (ProjectId) REFERENCES Project(ProjectId)
);

CREATE INDEX IdxDatabase_ProjectId ON Database(ProjectId);
CREATE INDEX IdxDatabase_Type ON Database(Type);
CREATE INDEX IdxDatabase_EntityId ON Database(EntityId);
```

### Table: DatabaseStat

```sql
-- linter-waive: MISSING-DESC-001 reason="Consolidated split-DB example; canonical version in 05-split-db-architecture/"
CREATE TABLE DatabaseStat (
    DatabaseStatId INTEGER PRIMARY KEY AUTOINCREMENT,
    DatabaseId INTEGER NOT NULL,
    RecordedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    SizeBytes INTEGER,
    RecordCount INTEGER,
    QueryCount INTEGER DEFAULT 0,
    AvgQueryMs REAL,
    FOREIGN KEY (DatabaseId) REFERENCES Database(DatabaseId)
);
```

---

## 4. Database Types

| Type | Purpose | Retention |
|------|---------|-----------|
| `history` | Version history tracking | Permanent |
| `cache` | Cached data (search, API) | 7–30 days |
| `config` | Configuration/settings | Permanent |
| `search` | Search results | 30 days |
| `session` | User sessions | 24 hours |
| `analytics` | Usage analytics | 90 days |
| `logs` | Application logs | 14 days |
| `queue` | Job/task queues | Until processed |

---

## 5. Concurrency & Locking — WAL Mode

All databases use Write-Ahead Logging for concurrent access. **Every connection MUST set all three PRAGMAs:**

```go
func (m *DbManager) configureDb(db *sql.DB) error {
    _, err := db.Exec("PRAGMA journal_mode=WAL")
    if err != nil {
        return err
    }

    _, err = db.Exec("PRAGMA busy_timeout=5000")
    if err != nil {
        return err
    }

    _, err = db.Exec("PRAGMA foreign_keys=ON")

    return err
}
```

---

## 6. Connection Pooling & DbManager

### DbManager Struct

```go
type DbManager struct {
    rootDb   *sql.DB
    dataDir  string
    openDbs  map[string]*sql.DB
    mu       sync.RWMutex
    maxOpen  int           // Default: 50
    maxIdle  int           // Default: 2 per DB
    connLife time.Duration // Default: 1h
    logger   *DbLogger
}
```

### Constructor

```go
func NewDbManager(dataDir string) apperror.Result[*DbManager] {
    if err := pathutil.EnsureDir(dataDir, 0755); err != nil {
        return nil, apperror.Wrap(err, ErrDbDirCreate, "create data directory")
    }

    rootPath := filepath.Join(dataDir, "root.db")
    rootDb, err := sql.Open("sqlite3", rootPath)
    if err != nil {
        return nil, apperror.Wrap(err, ErrDbOpen, "open root database").WithPath(rootPath)
    }

    manager := &DbManager{
        rootDb:  rootDb,
        dataDir: dataDir,
        openDbs: make(map[string]*sql.DB),
    }

    if err := manager.initRootSchema(); err != nil {
        return nil, err
    }

    return manager, nil
}
```

### Dynamic Database Creation

```go
func (m *DbManager) GetOrCreateDb(
    projectSlug string,
    dbType string,
    entityId string,
) apperror.Result[*sql.DB] {
    m.mu.Lock()
    defer m.mu.Unlock()

    key := fmt.Sprintf("%s/%s/%s", projectSlug, dbType, entityId)
    if db, ok := m.openDbs[key]; ok {
        return db, nil
    }

    project, err := m.getOrCreateProject(projectSlug)
    if err != nil {
        return nil, err
    }

    dbPath := m.buildDbPath(projectSlug, dbType, entityId)
    fullPath := filepath.Join(m.dataDir, dbPath)

    if err := pathutil.EnsureDir(filepath.Dir(fullPath), 0755); err != nil {
        return nil, apperror.Wrap(err, ErrDbDirCreate, "create database directory")
    }

    db, err := sql.Open("sqlite3", fullPath)
    if err != nil {
        return nil, apperror.Wrap(err, ErrDbOpen, "open database").WithPath(fullPath)
    }

    m.configureDb(db)
    m.openDbs[key] = db
    m.updateLastAccessed(dbRecord.DatabaseId)

    return db, nil
}
```

### Close All

```go
func (m *DbManager) Close() error {
    m.mu.Lock()
    defer m.mu.Unlock()

    for _, db := range m.openDbs {
        db.Close()
    }
    m.openDbs = make(map[string]*sql.DB)

    return m.rootDb.Close()
}
```

---

## 7. Lifecycle Management

### Database Creation Flow

1. Check if project exists in root.db, create if not
2. Check if database record exists, create if not
3. Create directory structure if needed
4. Open SQLite database file
5. Initialize schema (caller responsibility)

### Archive Stale Databases

```go
func (m *DbManager) ArchiveStale(maxAge time.Duration) error {
    cutoff := time.Now().Add(-maxAge)

    _, err := m.rootDb.Exec(`
        UPDATE Database
        SET Status = 'archived', UpdatedAt = CURRENT_TIMESTAMP
        WHERE LastAccessedAt < ? AND Status = 'active'
    `, cutoff)

    return err
}
```

### Purge Archived

```go
func (m *DbManager) PurgeArchived(retention time.Duration) error {
    cutoff := time.Now().Add(-retention)

    rows, _ := m.rootDb.Query(`
        SELECT Path FROM Database
        WHERE Status = 'archived' AND UpdatedAt < ?
    `, cutoff)
    defer rows.Close()

    for rows.Next() {
        var path string
        rows.Scan(&path)
        pathutil.Remove(filepath.Join(m.dataDir, path))
    }

    _, err := m.rootDb.Exec(`
        DELETE FROM Database
        WHERE Status = 'archived' AND UpdatedAt < ?
    `, cutoff)

    return err
}
```

---

## 8. Backup & Recovery

### Incremental Backup

Uses SQLite backup API for consistency. Backups are timestamped: `{projectSlug}/{timestamp}/`.

```go
func (m *DbManager) BackupProject(projectSlug, backupDir string) error {
    dbs, err := m.ListDatabases(projectSlug)
    if err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102-150405")
    projectBackupDir := filepath.Join(backupDir, projectSlug, timestamp)
    pathutil.EnsureDir(projectBackupDir, 0755)

    for _, db := range dbs {
        srcPath := filepath.Join(m.dataDir, db.Path)
        dstPath := filepath.Join(projectBackupDir, filepath.Base(db.Path))

        if err := m.backupDb(srcPath, dstPath); err != nil {
            return apperror.Wrap(err, ErrDbBackupFailed, "backup database").
                WithContext("path", db.Path)
        }
    }

    return nil
}
```

### Point-in-Time Recovery

Close all open databases for the project → restore from backup → re-open on next access.

### Zip Export/Import

```
project-a/ (folder) → ZIP → project-a.zip → UNZIP → project-a/ (restored)
```

**Export:**
```go
func (m *DbManager) ExportProjectToZip(projectSlug, outputPath string) error {
    projectDir := filepath.Join(m.dataDir, projectSlug)

    zipFile, err := pathutil.Create(outputPath)
    if err != nil {
        return apperror.Wrap(err, ErrFsWrite, "create zip file").WithPath(outputPath)
    }
    defer zipFile.Close()

    zipWriter := zip.NewWriter(zipFile)
    defer zipWriter.Close()

    return filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
        if err != nil || info.IsDir() || stringutil.IsMissingSuffix(path, ".db") {
            return nil
        }

        relPath, _ := filepath.Rel(projectDir, path)
        writer, _ := zipWriter.Create(relPath)
        file, _ := pathutil.Open(path)
        defer file.Close()
        _, err = io.Copy(writer, file)

        return err
    })
}
```

**Selective Export (by type):**
```go
func (m *DbManager) ExportByType(
    projectSlug string,
    dbTypes []string,
    outputPath string,
) error {
    dbs, err := m.ListDatabases(projectSlug)
    if err != nil {
        return err
    }

    typeSet := make(map[string]bool)
    for _, t := range dbTypes {
        typeSet[t] = true
    }

    // Create zip with only matching types
    // ... zip creation logic filtering by typeSet[db.Type]
}
```

**Import:**
```go
func (m *DbManager) ImportProjectFromZip(
    zipPath string,
    projectSlug string,
    overwrite bool,
) error {
    // 1. Close any open databases for this project
    m.closeProjectDbs(projectSlug)
    // 2. Extract zip to project directory
    // 3. Register imported databases in root.db
    m.registerImportedDatabases(projectSlug)
}
```

---

## 9. User-Scoped Isolation

Three scoping levels for data isolation:

| Level | Path Pattern | Use Case |
|-------|-------------|----------|
| App-level | `data/{app}/users/{userId}/` | Independent users |
| Company-level | `data/{app}/companies/{company}/users/{userId}/` | Enterprise multi-tenant |
| Module-level | `data/{app}/{module}/users/{userId}/` | Per-module isolation |

### User Database Registry (in Root DB)

```sql
-- linter-waive: MISSING-DESC-001 reason="Consolidated split-DB example; canonical version in 05-split-db-architecture/"
CREATE TABLE UserDbRegistry (
    UserDbRegistryId INTEGER PRIMARY KEY AUTOINCREMENT,
    UserId INTEGER NOT NULL,
    CompanyId INTEGER,
    Category TEXT NOT NULL,
    EntityId TEXT NOT NULL,
    SequenceNum INTEGER NOT NULL,
    Path TEXT NOT NULL,
    SizeBytes INTEGER DEFAULT 0,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    LastAccessedAt DATETIME,
    ExpiresAt DATETIME,
    Status TEXT DEFAULT 'active',
    FOREIGN KEY (UserId) REFERENCES User(UserId)
);
```

---

## 10. RBAC via Casbin

Role-Based Access Control using Casbin with SQLite via GORM adapter.

### Scope Levels

| Level | Database Location | Use Case |
|-------|-------------------|----------|
| Root | `data/rbac.db` | Multi-tenant platform |
| App | `data/{app}/rbac.db` | Single-application |
| Company | `data/{app}/companies/{slug}/rbac.db` | Enterprise |

### Casbin Model (`rbac_model.conf`)

```ini
[request_definition]
r = Sub, Obj, Act

[policy_definition]
p = Sub, Obj, Act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.Sub, p.Sub) && keyMatch2(r.Obj, p.Obj) && r.Act == p.Act
```

### Seed Policies

```sql
INSERT INTO CasbinRule (Ptype, V0, V1, V2) VALUES
('p', 'admin', '*', '*'),
('p', 'manager', '/api/*', 'read'),
('p', 'manager', '/api/*', 'write'),
('p', 'viewer', '/api/*', 'read');

-- Role hierarchy
INSERT INTO CasbinRule (Ptype, V0, V1) VALUES
('g', 'admin', 'manager'),
('g', 'manager', 'editor'),
('g', 'editor', 'viewer');
```

### HTTP Middleware

```go
func RbacMiddleware(manager *rbac.RbacManager) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user := r.Context().Value("userId").(string)
            resource := r.URL.Path
            action := methodToAction(r.Method) // GET→read, POST→write, DELETE→delete

            allowed, err := manager.Enforce(user, resource, action)
            if err != nil || !allowed {
                http.Error(w, "Forbidden", http.StatusForbidden)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

---

## 11. 2-Step Reset API Standard

All CLIs implement a standardized safe-deletion flow with 5-minute TTL.

### Flow

1. **Request** — `POST /api/v1/reset/request` → returns `ResetId`, preview of affected items, 5-min expiry
2. **Confirm** — `POST /api/v1/reset/confirm` → executes deletion, returns count and freed bytes
3. **Cancel** — `POST /api/v1/reset/cancel` → cancels pending reset

### Reset Database Schema

```sql
-- linter-waive: MISSING-DESC-001 reason="Consolidated split-DB example; canonical version in 05-split-db-architecture/"
CREATE TABLE ResetRequest (
    ResetRequestId INTEGER PRIMARY KEY AUTOINCREMENT,
    ResetToken TEXT UNIQUE NOT NULL,
    Scope TEXT NOT NULL,
    RequestedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    ExpiresAt DATETIME NOT NULL,
    AffectedItems TEXT,
    Status TEXT DEFAULT 'pending',
    DeletedCount INTEGER,
    FreedBytes INTEGER
);
```

---

## 12. Structured Logging

All database operations are logged with structured context:

### Log Levels by Operation

| Operation | Success Level | Failure Level |
|-----------|---------------|---------------|
| GetOrCreateDb | INFO | ERROR |
| ListDatabases | DEBUG | WARN |
| Export/Import | INFO | ERROR |
| Backup | INFO | ERROR |
| Archive/Purge | INFO | ERROR |

### Log Format

```json
{
  "time": "2026-02-01T10:30:00Z",
  "level": "INFO",
  "msg": "Database ready",
  "project": "my-project",
  "type": "history",
  "entity": "readme-md",
  "duration_ms": 12,
  "cached": false
}
```

---

## 13. Usage Patterns

### History Database

```go
historyDb, err := manager.GetOrCreateDb("my-project", "history", "readme-md")
if err != nil {
    return err
}

_, err = historyDb.Exec(`
    CREATE TABLE IF NOT EXISTS Version (
        VersionId INTEGER PRIMARY KEY AUTOINCREMENT,
        Content TEXT NOT NULL,
        Hash TEXT NOT NULL,
        CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
        Author TEXT
    )
`)
```

### Cache Database

```go
cacheDb, err := manager.GetOrCreateDb("my-project", "cache", "search")
if err != nil {
    return err
}

_, err = cacheDb.Exec(`
    CREATE TABLE IF NOT EXISTS SearchCache (
        SearchCacheId INTEGER PRIMARY KEY AUTOINCREMENT,
        QueryHash TEXT UNIQUE NOT NULL,
        Results TEXT NOT NULL,
        CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
        ExpiresAt DATETIME NOT NULL
    )
`)
```

---

## 14. Applicable Projects

| Project | Usage |
|---------|-------|
| Spec Management | File history, search cache |
| GSearch CLI | Search results, cache |
| BRun CLI | Build artifacts, logs |
| AI Bridge | Conversation history, chat DBs |
| Nexus Flow | Workflow state, execution history |

---

## 15. Key Rules Summary

| # | Rule |
|---|------|
| 1 | **WAL mode** on every database — no exceptions |
| 2 | **Connection pooling** — max 50 open databases, 2 idle per DB |
| 3 | **Dynamic creation** — databases created on first access, registered in Root DB |
| 4 | **PascalCase columns** — zero tolerance for underscores |
| 5 | **Zip export** — any project folder can be zipped for backup/migration |
| 6 | **Foreign keys enabled** — `PRAGMA foreign_keys=ON` on every connection |
| 7 | **Busy timeout** — `PRAGMA busy_timeout=5000` prevents SQLITE_BUSY errors |
| 8 | **Archive stale** — databases not accessed in 30 days get archived |
| 9 | **Purge archived** — archived databases beyond retention period get deleted |
| 10 | **2-step reset** — all destructive operations require request + confirm with 5-min TTL |
| 11 | **Structured logging** — all operations logged with project/type/entity context |

---

*Consolidated split-db architecture — v3.2.0 — 2026-04-16*
