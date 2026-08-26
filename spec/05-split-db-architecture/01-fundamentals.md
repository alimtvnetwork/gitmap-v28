# Split Database Architecture — Fundamentals

**Version:** 3.4.0
**Updated:** 2026-05-06
**Parent:** [00-overview.md](./00-overview.md)

> **v3.4.0 — DB convention propagation (Phase 2):** Every tier in this architecture inherits the global DB conventions:
> - All `*At` timestamp columns MUST be `INTEGER` (epoch seconds, UTC) per [`spec/04-database-conventions/01-naming-conventions.md`](../04-database-conventions/01-naming-conventions.md) **Rule 7.1 v2**. TEXT/ISO-8601 storage is forbidden for new tier-local schemas.
> - Every enum-like reference table in any tier MUST follow the canonical `(Id, Code, Label)` shape per **Rule 13** ([`spec/04-database-conventions/02-schema-design.md`](../04-database-conventions/02-schema-design.md) §6.5). The legacy `{Table}Code` / `{Table}Label` column-name style is removed.
> - These rules apply uniformly to Root, Settings, App, Session, Cache, and Document tiers.

> **Downstream tier consumers (FU-6 back-link, added 2026-05-04):** This file defines the canonical 6-tier model (Root / Settings / App / Session / Cache / Document). Service-level specs that consume only a subset MUST publish a reconciliation file naming which tiers they use. Current consumers:
> - **Main / Worker service** — see [`spec/19-main-worker-service/11-split-db-tier-reconciliation.md`](../19-main-worker-service/11-split-db-tier-reconciliation.md). Main = 3 tiers (Root/Settings/Session); Worker = 4 tiers (Root/Settings/App/Session). Cache + Document unused in v1.0.

---

## Database Terminology

| Term | Meaning | Example Path |
|------|---------|--------------|
| **Root DB** | Global registry, settings, app list | `data/aibridge.db` |
| **Settings DB** | Configuration (seeded + user) | Inside Root DB |
| **App DB** | Application-scoped metadata | `data/{appName}/search.db` |
| **Session DB** | Per-session isolated storage | `data/{appName}/ai/chat/001-{id}.db` |
| **Cache DB** | Cached results with TTL | `data/{appName}/rag/cache/search/001-{slug}.db` |
| **Document DB** | RAG chunks + embeddings | `data/{appName}/rag/documents/001-{id}.db` |

---

## Core Concepts

---

## Hierarchical Structure Examples

### 2-Layer Structure (Simple)

```
data/
├── root.db                              # Root registry
└── {project-slug}/
    ├── config.db                        # Project config
    ├── cache.db                         # Project cache
    └── logs.db                          # Project logs
```

### 3-Layer Structure (Standard - Most Common)

```
data/
├── root.db                              # Root registry database
├── {project-slug}/
│   ├── history/                         # History databases folder
│   │   ├── {file-slug}.db               # Per-file history
│   │   ├── {file-slug-2}.db
│   │   └── ...
│   ├── cache/                           # Cache databases folder
│   │   └── search-cache.db
│   ├── config/                          # Config databases
│   │   └── settings.db
│   ├── chat/                            # Chat session databases
│   │   ├── {session-id}.db
│   │   └── ...
│   ├── voice/                           # Voice recording databases
│   │   └── {recording-id}.db
│   └── search/                          # Search index databases
│       └── {index-id}.db
└── {project-slug-2}/
    └── ...
```

### 4-Layer Structure (Complex - With Categories)

```
data/
├── root.db                              # Root registry database
├── {project-slug}/
│   ├── ai/                              # AI category
│   │   ├── chat/                        # Chat type
│   │   │   ├── {session-id}.db
│   │   │   └── ...
│   │   ├── embeddings/                  # Embeddings type
│   │   │   └── {model-id}.db
│   │   └── prompts/                     # Prompts type
│   │       └── {template-id}.db
│   ├── workflow/                        # Workflow category
│   │   ├── history/                     # History type
│   │   │   └── {file-slug}.db
│   │   └── queue/                       # Queue type
│   │       └── {queue-id}.db
│   └── search/                          # Search category
│       ├── indices/                     # Indices type
│       │   └── {index-id}.db
│       └── cache/                       # Cache type
│           └── {query-hash}.db
└── {project-slug-2}/
    └── ...
```

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                         SPLIT DATABASE ARCHITECTURE (v2.0)                           │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                      │
│                              ┌─────────────────┐                                     │
│                              │    ROOT.DB      │                                     │
│                              │   (Registry +   │                                     │
│                              │    Logging)     │                                     │
│                              └────────┬────────┘                                     │
│                                       │                                              │
│           ┌───────────────────────────┼───────────────────────────┐                  │
│           │                           │                           │                  │
│           ▼                           ▼                           ▼                  │
│  ┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐            │
│  │  PROJECT-A/     │       │  PROJECT-B/     │       │  PROJECT-C/     │            │
│  └────────┬────────┘       └────────┬────────┘       └────────┬────────┘            │
│           │                         │                         │                     │
│   ┌───────┴───────┐                ...                       ...                    │
│   │       │       │                                                                 │
│   ▼       ▼       ▼                                                                 │
│  ai/   workflow/ search/     ← Categories (optional 4-layer)                        │
│   │                                                                                 │
│   ├── chat/                  ← Types                                                │
│   │    │                                                                            │
│   │    ├── session-001.db    ← Entity DBs                                           │
│   │    └── session-002.db                                                           │
│   │                                                                                 │
│   └── embeddings/                                                                   │
│        └── gpt-4.db                                                                 │
│                                                                                      │
│  ┌──────────────────────────────────────────────────────────────────────────────┐   │
│  │                           IMPORT / EXPORT                                     │   │
│  │  ┌─────────────┐    ZIP     ┌─────────────┐    UNZIP    ┌─────────────┐      │   │
│  │  │ project-a/  │ ────────►  │ project-a   │ ──────────► │ project-a/  │      │   │
│  │  │ (folder)    │            │ .zip        │             │ (restored)  │      │   │
│  │  └─────────────┘            └─────────────┘             └─────────────┘      │   │
│  └──────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                      │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

---

## Root Database Schema

### Table: Project

```sql
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE Project (
    ProjectId INTEGER PRIMARY KEY AUTOINCREMENT,
    Slug TEXT UNIQUE NOT NULL,
    DisplayName TEXT NOT NULL,
    Path TEXT NOT NULL,                    -- Relative path to project folder
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    Status TEXT DEFAULT 'active'           -- active, archived, deleted
);

CREATE INDEX IdxProject_Slug ON Project(Slug);
CREATE INDEX IdxProject_Status ON Project(Status);
```

### Table: Database

```sql
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE Database (
    DatabaseId INTEGER PRIMARY KEY AUTOINCREMENT,
    ProjectId INTEGER NOT NULL,
    Type TEXT NOT NULL,                    -- history, cache, config, search, etc.
    EntityId TEXT,                         -- File slug, search ID, etc.
    Path TEXT NOT NULL,                    -- Relative path to .db file
    SizeBytes INTEGER DEFAULT 0,
    RecordCount INTEGER DEFAULT 0,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    LastAccessedAt DATETIME,
    Status TEXT DEFAULT 'active',          -- active, archived, deleted
    FOREIGN KEY (ProjectId) REFERENCES Project(ProjectId)
);

CREATE INDEX IdxDatabase_ProjectId ON Database(ProjectId);
CREATE INDEX IdxDatabase_Type ON Database(Type);
CREATE INDEX IdxDatabase_EntityId ON Database(EntityId);
```

### Table: DatabaseStat

```sql
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
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

CREATE INDEX IdxDatabaseStat_DatabaseId ON DatabaseStat(DatabaseId);
CREATE INDEX IdxDatabaseStat_RecordedAt ON DatabaseStat(RecordedAt);
```

---

## Database Types

| Type | Purpose | Entity Example | Retention |
|------|---------|----------------|-----------|
| `history` | Version history tracking | File slug | Permanent |
| `cache` | Cached data (search, API) | Cache type | 7-30 days |
| `config` | Configuration/settings | - | Permanent |
| `search` | Search results | Search ID | 30 days |
| `session` | User sessions | Session ID | 24 hours |
| `analytics` | Usage analytics | - | 90 days |
| `logs` | Application logs | Date | 14 days |
| `queue` | Job/task queues | Queue name | Until processed |

---

## Concurrency & Locking

### SQLite WAL Mode

All databases use Write-Ahead Logging for concurrent access:

```go
func (m *DbManager) configureDb(db *sql.DB) error {
    // Enable WAL mode for concurrent reads
    _, err := db.Exec("PRAGMA journal_mode=WAL")
    if err != nil {
        return err
    }
    
    // Set busy timeout to avoid SQLITE_BUSY errors
    _, err = db.Exec("PRAGMA busy_timeout=5000")
    if err != nil {
        return err
    }
    
    // Enable foreign keys
    _, err = db.Exec("PRAGMA foreign_keys=ON")
    return err
}
```

### Connection Pooling

```go
type DbManager struct {
    rootDb    *sql.DB
    dataDir   string
    openDbs   map[string]*sql.DB
    mu        sync.RWMutex
    maxOpen   int           // Max open databases (default: 50)
    maxIdle   int           // Max idle connections per DB (default: 2)
    connLife  time.Duration // Max connection lifetime (default: 1h)
}

func (m *DbManager) getDb(key string) (*sql.DB, bool) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    db, ok := m.openDbs[key]
    return db, ok
}
```

---

## Backup & Recovery

### Incremental Backup

```go
// BackupProject creates a backup of all databases for a project
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
        
        // Use SQLite backup API for consistency
        if err := m.backupDb(srcPath, dstPath); err != nil {
            return apperror.Wrap(
                err,
                ErrDbBackupFailed,
                "backup database",
            ).WithContext("path", db.Path)
        }
    }
    
    return nil
}
```

### Point-in-Time Recovery

```go
// RestoreProject restores databases from a backup
func (m *DbManager) RestoreProject(projectSlug, backupPath string) error {
    // Close all open databases for this project
    m.closeProjectDbs(projectSlug)
    
    // Restore from backup
    return filepath.Walk(backupPath, func(path string, info os.FileInfo, err error) error {
        if err != nil || info.IsDir() || stringutil.IsMissingSuffix(path, ".db") {
            return err
        }
        
        relPath := strings.TrimPrefix(path, backupPath)
        dstPath := filepath.Join(m.dataDir, projectSlug, relPath)
        
        return copyFile(path, dstPath)
    })
}
```

---

## Go Implementation

### DbManager Interface

```go
package splitdb

import (
    // ALLOWED: database infrastructure layer — manages raw SQLite connection pooling and dynamic DB routing
    "database/sql"
    "fmt"
    "path/filepath"
    "sync"
    "time"
    
    _ "github.com/mattn/go-sqlite3"
    "pkg/pathutil"
)

type DbManager struct {
    rootDb    *sql.DB
    dataDir   string
    openDbs   map[string]*sql.DB
    mu        sync.RWMutex
}

type Project struct {
    ProjectId  int64
    Slug        string
    DisplayName string
    Path        string
    Status      string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type Database struct {
    DatabaseId  int64
    ProjectId   int64
    Type         string
    EntityId     string
    Path         string
    SizeBytes    int64
    RecordCount  int64
    Status       string
    CreatedAt    time.Time
    UpdatedAt    time.Time
    LastAccessed *time.Time
}

// NewDbManager creates a new split database manager
func NewDbManager(dataDir string) apperror.Result[*DbManager] {
    if err := pathutil.EnsureDir(dataDir, 0755); err != nil {
        return nil, apperror.Wrap(
            err,
            ErrDbDirCreate,
            "create data directory",
        ).WithContext("dir", dataDir)
    }
    
    rootPath := filepath.Join(dataDir, "root.db")
    rootDb, err := sql.Open("sqlite3", rootPath)
    if err != nil {
        return nil, apperror.Wrap(
            err,
            ErrDbOpen,
            "open root database",
        ).WithPath(rootPath)
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

// GetOrCreateDb returns a database, creating it if it doesn't exist
func (m *DbManager) GetOrCreateDb(projectSlug, dbType, entityId string) apperror.Result[*sql.DB] {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    // Check if already open
    key := fmt.Sprintf("%s/%s/%s", projectSlug, dbType, entityId)
    if db, ok := m.openDbs[key]; ok {
        return db, nil
    }
    
    // Ensure project exists
    project, err := m.getOrCreateProject(projectSlug)
    if err != nil {
        return nil, err
    }
    
    // Get or create database record
    dbPath := m.buildDbPath(projectSlug, dbType, entityId)
    dbRecord, err := m.getOrCreateDatabase(project.ProjectId, dbType, entityId, dbPath)
    if err != nil {
        return nil, err
    }
    
    // Ensure directory exists
    dir := filepath.Dir(filepath.Join(m.dataDir, dbRecord.Path))
    if err := pathutil.EnsureDir(dir, 0755); err != nil {
        return nil, apperror.Wrap(
            err,
            ErrDbDirCreate,
            "create database directory",
        ).WithContext("dir", dir)
    }
    
    // Open the database
    fullPath := filepath.Join(m.dataDir, dbRecord.Path)
    db, err := sql.Open("sqlite3", fullPath)
    if err != nil {
        return nil, apperror.Wrap(
            err,
            ErrDbOpen,
            "open database",
        ).WithPath(fullPath)
    }
    
    m.openDbs[key] = db
    
    // Update last accessed
    m.updateLastAccessed(dbRecord.DatabaseId)
    
    return db, nil
}

// ListDatabases returns all databases for a project
func (m *DbManager) ListDatabases(projectSlug string) apperror.Result[[]Database] {
    query := `
        SELECT D.DatabaseId, D.ProjectId, D.Type, D.EntityId, D.Path, 
               D.SizeBytes, D.RecordCount, D.Status, D.CreatedAt, D.UpdatedAt
        FROM Database D
        JOIN Project P ON D.ProjectId = P.ProjectId
        WHERE P.Slug = ? AND D.Status = 'active'
    `
    
    rows, err := m.rootDb.Query(query, projectSlug)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var dbs []Database
    for rows.Next() {
        var db Database
        if err := rows.Scan(
            &db.DatabaseId, &db.ProjectId, &db.Type, &db.EntityId, &db.Path,
            &db.SizeBytes, &db.RecordCount, &db.Status, &db.CreatedAt, &db.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        dbs = append(dbs, db)
    }
    
    return dbs, nil
}

// Close closes all open databases
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

## Usage Examples

### History Database Pattern

```go
// Get history database for a specific file
historyDb, err := manager.GetOrCreateDb("my-project", "history", "readme-md")
if err != nil {
    return err
}

// Create history table if not exists
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

### Cache Database Pattern

```go
// Get cache database for search results
cacheDb, err := manager.GetOrCreateDb("my-project", "cache", "search")
if err != nil {
    return err
}

// Create cache table if not exists
_, err = cacheDb.Exec(`
    CREATE TABLE IF NOT EXISTS SearchCache (
        SearchCacheId INTEGER PRIMARY KEY AUTOINCREMENT,
        QueryHash TEXT UNIQUE NOT NULL,
        Results TEXT NOT NULL,          -- JSON encoded
        CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
        ExpiresAt DATETIME NOT NULL
    )
`)
```

---

## File Path Convention

| Component | Pattern | Example |
|-----------|---------|---------|
| Root DB | `{data}/root.db` | `data/root.db` |
| Project Folder | `{data}/{project-slug}/` | `data/my-project/` |
| Type Folder | `{project}/{type}/` | `data/my-project/history/` |
| Entity DB | `{type}/{entity-slug}.db` | `data/my-project/history/readme-md.db` |

### Slug Generation

```go
func GenerateSlug(name string) string {
    // Convert to lowercase
    slug := strings.ToLower(name)
    // Replace spaces and special chars with hyphens
    slug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")
    // Remove leading/trailing hyphens
    slug = strings.Trim(slug, "-")
    return slug
}
```

---

## Lifecycle Management

### Database Creation

1. Check if project exists in root.db, create if not
2. Check if database record exists, create if not
3. Create directory structure if needed
4. Open SQLite database file
5. Initialize schema (caller responsibility)

### Database Cleanup

```go
// Archive databases not accessed in 30 days
func (m *DbManager) ArchiveStale(maxAge time.Duration) error {
    cutoff := time.Now().Add(-maxAge)
    
    _, err := m.rootDb.Exec(`
        UPDATE Database 
        SET Status = 'archived', UpdatedAt = CURRENT_TIMESTAMP
        WHERE LastAccessedAt < ? AND Status = 'active'
    `, cutoff)
    
    return err
}

// Delete archived databases older than retention period
func (m *DbManager) PurgeArchived(retention time.Duration) error {
    cutoff := time.Now().Add(-retention)
    
    // Get databases to delete
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
    
    // Remove records
    _, err := m.rootDb.Exec(`
        DELETE FROM Database 
        WHERE Status = 'archived' AND UpdatedAt < ?
    `, cutoff)
    
    return err
}
```

---

## Import / Export (Zip Files)

### Export Project to Zip

```go
// ExportProjectToZip creates a zip file of all project databases
func (m *DbManager) ExportProjectToZip(projectSlug, outputPath string) error {
    m.logger.Info("Starting export", "project", projectSlug, "output", outputPath)
    
    projectDir := filepath.Join(m.dataDir, projectSlug)
    if pathutil.IsDirMissing(projectDir) {
        return apperror.New(
            ErrProjectNotFound,
            "project not found",
        ).WithContext("project", projectSlug)
    }
    
    // Create zip file
    zipFile, err := pathutil.Create(outputPath)
    if err != nil {
        return apperror.Wrap(
            err,
            ErrFsWrite,
            "create zip file",
        ).WithPath(outputPath)
    }
    defer zipFile.Close()
    
    zipWriter := zip.NewWriter(zipFile)
    defer zipWriter.Close()
    
    // Walk project directory and add all .db files
    err = filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            m.logger.Warn("Skip file due to error", "path", path, "error", err)
            return nil // Continue walking
        }
        
        if info.IsDir() || stringutil.IsMissingSuffix(path, ".db") {
            return nil
        }
        
        // Get relative path within project
        relPath, _ := filepath.Rel(projectDir, path)
        
        m.logger.Debug("Adding to zip", "file", relPath, "size", info.Size())
        
        // Create zip entry
        writer, err := zipWriter.Create(relPath)
        if err != nil {
            return err
        }
        
        // Copy file content
        file, err := pathutil.Open(path)
        if err != nil {
            return err
        }
        defer file.Close()
        
        _, err = io.Copy(writer, file)
        return err
    })
    
    if err != nil {
        return apperror.Wrap(
            err,
            ErrExportFailed,
            "export project to zip",
        ).WithContext("project", projectSlug)
    }
    
    m.logger.Info("Export complete", "project", projectSlug, "output", outputPath)
    return nil
}
```

### Import Project from Zip

```go
// ImportProjectFromZip imports databases from a zip file
func (m *DbManager) ImportProjectFromZip(zipPath, projectSlug string, overwrite bool) error {
    m.logger.Info("Starting import", "zip", zipPath, "project", projectSlug, "overwrite", overwrite)
    
    // Open zip file
    reader, err := zip.OpenReader(zipPath)
    if err != nil {
        return apperror.Wrap(
            err,
            ErrFsRead,
            "open zip file",
        ).WithPath(zipPath)
    }
    defer reader.Close()
    
    projectDir := filepath.Join(m.dataDir, projectSlug)
    
    // Check if project exists
    isProjectExists := pathutil.IsDir(projectDir)
    isReadOnly := !overwrite
    isProjectConflict := isProjectExists && isReadOnly

    if isProjectConflict {
        return apperror.FailNew[ImportResult](
            errors.ErrFsConflict,
            "project already exists; use overwrite=true to replace",
        )
    }
    
    // Close any open databases for this project
    m.closeProjectDbs(projectSlug)
    
    // Create project directory
    if err := pathutil.EnsureDir(projectDir, 0755); err != nil {
        return err
    }
    
    // Extract files
    for _, file := range reader.File {
        if file.FileInfo().IsDir() {
            continue
        }
        
        destPath := filepath.Join(projectDir, file.Name)
        
        m.logger.Debug("Extracting", "file", file.Name, "size", file.UncompressedSize64)
        
        // Create directory structure
        if err := pathutil.EnsureDir(filepath.Dir(destPath), 0755); err != nil {
            return err
        }
        
        // Extract file
        if err := m.extractZipFile(file, destPath); err != nil {
            return apperror.Wrap(
                err,
                ErrImportFailed,
                "extract zip entry",
            ).WithContext("file", file.Name)
        }
    }
    
    // Register databases in root.db
    if err := m.registerImportedDatabases(projectSlug); err != nil {
        m.logger.Warn("Failed to register databases", "error", err)
    }
    
    m.logger.Info("Import complete", "project", projectSlug, "files", len(reader.File))
    return nil
}

func (m *DbManager) extractZipFile(file *zip.File, destPath string) error {
    src, err := file.Open()
    if err != nil {
        return err
    }
    defer src.Close()
    
    dst, err := pathutil.Create(destPath)
    if err != nil {
        return err
    }
    defer dst.Close()
    
    _, err = io.Copy(dst, src)
    return err
}
```

### Selective Export (By Type/Category)

```go
// ExportByType exports only specific database types
func (m *DbManager) ExportByType(projectSlug string, dbTypes []string, outputPath string) error {
    m.logger.Info("Selective export", "project", projectSlug, "types", dbTypes)
    
    // Filter databases by type
    dbs, err := m.ListDatabases(projectSlug)
    if err != nil {
        return err
    }
    
    typeSet := make(map[string]bool)
    for _, t := range dbTypes {
        typeSet[t] = true
    }
    
    // Create zip with only matching types
    zipFile, err := pathutil.Create(outputPath)
    if err != nil {
        return err
    }
    defer zipFile.Close()
    
    zipWriter := zip.NewWriter(zipFile)
    defer zipWriter.Close()
    
    for _, db := range dbs {
        if !typeSet[db.Type] {
            continue
        }
        
        m.logger.Debug("Including", "type", db.Type, "path", db.Path)
        
        fullPath := filepath.Join(m.dataDir, db.Path)
        relPath := strings.TrimPrefix(db.Path, projectSlug+"/")
        
        writer, _ := zipWriter.Create(relPath)
        file, _ := pathutil.Open(fullPath)
        io.Copy(writer, file)
        file.Close()
    }
    
    return nil
}
```

---

## Logging System

### Structured Logging

All database operations are logged with structured context for debugging and audit:

```go
type DbLogger struct {
    logger *slog.Logger
}

func NewDbLogger(output io.Writer) *DbLogger {
    return &DbLogger{
        logger: slog.New(slog.NewJsonHandler(output, &slog.HandlerOptions{
            Level: slog.LevelDebug,
        })),
    }
}

func (l *DbLogger) Info(msg string, args ...any) {
    l.logger.Info(msg, args...)
}

func (l *DbLogger) Debug(msg string, args ...any) {
    l.logger.Debug(msg, args...)
}

func (l *DbLogger) Warn(msg string, args ...any) {
    l.logger.Warn(msg, args...)
}

func (l *DbLogger) Error(msg string, args ...any) {
    l.logger.Error(msg, args...)
}
```

### Operation Logging

```go
// GetOrCreateDb with logging
func (m *DbManager) GetOrCreateDb(projectSlug, dbType, entityId string) apperror.Result[*sql.DB] {
    startTime := time.Now()
    
    m.logger.Debug("GetOrCreateDb called",
        "project", projectSlug,
        "type", dbType,
        "entity", entityId,
    )
    
    // ... existing logic ...
    
    db, err := m.doGetOrCreateDb(projectSlug, dbType, entityId)
    
    duration := time.Since(startTime)
    
    if err != nil {
        m.logger.Error("GetOrCreateDb failed",
            "project", projectSlug,
            "type", dbType,
            "entity", entityId,
            "error", err,
            "duration_ms", duration.Milliseconds(),
        )
        return nil, err
    }
    
    m.logger.Info("Database ready",
        "project", projectSlug,
        "type", dbType,
        "entity", entityId,
        "duration_ms", duration.Milliseconds(),
        "cached", cached,
    )
    
    return db, nil
}
```

### Log Levels by Operation

| Operation | Success Level | Failure Level |
|-----------|---------------|---------------|
| GetOrCreateDb | INFO | ERROR |
| ListDatabases | DEBUG | WARN |
| Export | INFO | ERROR |
| Import | INFO | ERROR |
| Backup | INFO | ERROR |
| Archive/Purge | INFO | ERROR |
| Query Stats | DEBUG | WARN |

### Log Format Examples

```json
// Successful operation
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

// Export operation
{
  "time": "2026-02-01T10:35:00Z",
  "level": "INFO",
  "msg": "Export complete",
  "project": "my-project",
  "output": "/backups/my-project-2026-02-01.zip",
  "files_count": 15,
  "total_size_bytes": 1048576,
  "duration_ms": 250
}

// Error with context
{
  "time": "2026-02-01T10:40:00Z",
  "level": "ERROR",
  "msg": "Import failed",
  "zip": "/imports/invalid.zip",
  "error": "zip: not a valid zip file",
  "project": "new-project"
}
```

---

## Benefits

| Benefit | Description |
|---------|-------------|
| **Isolation** | Each entity has its own database, preventing table bloat |
| **Performance** | Smaller databases = faster queries |
| **Scalability** | Add new entities without affecting existing ones |
| **Backup** | Backup individual databases or folders |
| **Cleanup** | Easy to archive/delete unused databases |
| **Debugging** | Inspect specific databases in isolation |
| **Portability** | Easy import/export via zip files |
| **Auditability** | Structured logging for all operations |

---

## Applicable Projects

This pattern is used by:

| Project | Usage |
|---------|-------|
| Spec Management | File history, search cache |
| GSearch CLI | Search results, cache |
| BRun CLI | Build artifacts, logs |
| AI Bridge | Conversation history, chat DBs |
| Nexus Flow | Workflow state, execution history |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Seedable Config | [Seedable Config Architecture](../06-seedable-config-architecture/00-overview.md) |

---

*This pattern ensures consistent database organization across all projects.*
