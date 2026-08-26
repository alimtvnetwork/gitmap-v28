# Split DB Architecture: User-Scoped Isolation

**Version:** 3.2.0  
**Created:** 2026-03-09  
**Status:** Active  
**Parent:** [00-overview.md](../00-overview.md)

---

## Overview

This document defines the **user-scoped isolation pattern** for the Split DB architecture. This pattern enables per-user data isolation within applications, complementing the existing company-scoped isolation pattern.

---

## Scoping Hierarchy

The Split DB architecture supports three scoping levels:

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                           SCOPING HIERARCHY                                          │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                      │
│   ┌─────────────────┐                                                               │
│   │   ROOT LEVEL    │  data/root.db                                                 │
│   │   (Global)      │  data/{appName}/                                              │
│   └────────┬────────┘                                                               │
│            │                                                                         │
│   ┌────────┴────────┐                                                               │
│   │                  │                                                               │
│   ▼                  ▼                                                               │
│ ┌─────────────┐   ┌─────────────┐                                                   │
│ │  COMPANY    │   │    USER     │                                                   │
│ │   SCOPE     │   │   SCOPE     │                                                   │
│ │             │   │             │                                                   │
│ │ companies/  │   │  users/     │                                                   │
│ │ {company}/  │   │  {userId}/  │                                                   │
│ └──────┬──────┘   └──────┬──────┘                                                   │
│        │                  │                                                          │
│        ▼                  ▼                                                          │
│ ┌─────────────┐   ┌─────────────┐                                                   │
│ │  USER SCOPE │   │  SESSION    │                                                   │
│ │  (Company)  │   │  SCOPE      │                                                   │
│ │             │   │             │                                                   │
│ │ {company}/  │   │ sessions/   │                                                   │
│ │ users/      │   │ {session}/  │                                                   │
│ │ {userId}/   │   │             │                                                   │
│ └─────────────┘   └─────────────┘                                                   │
│                                                                                      │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

---

## Directory Structure Patterns

### Pattern 1: App-Level User Isolation

For applications where users are independent (not grouped by company):

```
data/
├── root.db                                    # Global registry
├── {appName}/
│   ├── app.db                                 # App metadata
│   │
│   └── users/
│       ├── {userId-1}/
│       │   ├── settings.db                    # User preferences
│       │   ├── sessions/
│       │   │   └── {sessionId}.db             # User sessions
│       │   ├── history/
│       │   │   └── {historyId}.db             # User activity history
│       │   └── data/
│       │       └── {dataType}/{entityId}.db   # User-owned data
│       │
│       ├── {userId-2}/
│       │   └── ...
│       └── ...
```

### Pattern 2: Company + User Isolation (Enterprise)

For multi-tenant applications with company hierarchy:

```
data/
├── root.db                                    # Global registry
├── {appName}/
│   ├── app.db                                 # App metadata
│   │
│   └── companies/
│       ├── {companySlug-1}/
│       │   ├── company.db                     # Company metadata
│       │   │
│       │   └── users/
│       │       ├── {userId-1}/
│       │       │   ├── settings.db            # User preferences
│       │       │   ├── sessions/
│       │       │   │   └── {sessionId}.db
│       │       │   └── data/
│       │       │       └── ...
│       │       │
│       │       └── {userId-2}/
│       │           └── ...
│       │
│       └── {companySlug-2}/
│           └── ...
```

### Pattern 3: Module-Based User Isolation

For applications with distinct modules:

```
data/
├── root.db                                    # Global registry
├── {appName}/
│   │
│   ├── chat/
│   │   └── users/
│   │       └── {userId}/
│   │           └── {sessionId}.db
│   │
│   ├── rag/
│   │   └── users/
│   │       └── {userId}/
│   │           └── documents/
│   │               └── {docId}.db
│   │
│   └── seo/
│       └── companies/
│           └── {companySlug}/
│               └── users/
│                   └── {userId}/
│                       └── jobs/
│                           └── {jobId}.db
```

---

## Database Schema

### User Registry Table (in Root or App DB)

```sql
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE User (
    UserId INTEGER PRIMARY KEY AUTOINCREMENT,
    ExternalId TEXT UNIQUE,                    -- External auth provider ID
    Username TEXT UNIQUE NOT NULL,
    Email TEXT UNIQUE,
    DisplayName TEXT,
    AvatarUrl TEXT,
    CompanyId INTEGER,                         -- NULL for app-level users
    Status TEXT DEFAULT 'active',              -- active, suspended, deleted
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    LastActiveAt DATETIME,
    FOREIGN KEY (CompanyId) REFERENCES Company(CompanyId)
);

CREATE INDEX IdxUser_ExternalId ON User(ExternalId);
CREATE INDEX IdxUser_CompanyId ON User(CompanyId);
CREATE INDEX IdxUser_Status ON User(Status);
```

### User Settings Table (in User's settings.db)

```sql
CREATE TABLE Setting (
    SettingId INTEGER PRIMARY KEY AUTOINCREMENT,
    Key TEXT UNIQUE NOT NULL,
    Value TEXT NOT NULL,
    ValueType TEXT DEFAULT 'string',           -- string, int, bool, json
    Source TEXT DEFAULT 'user',                -- seed, user, admin
    Description TEXT,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Default user settings
INSERT INTO Settings (Key, Value, ValueType, Source) VALUES
('Theme', 'system', 'string', 'seed'),
('Language', 'en', 'string', 'seed'),
('Notifications.Email', 'true', 'bool', 'seed'),
('Notifications.Push', 'true', 'bool', 'seed'),
('Privacy.ShareAnalytics', 'false', 'bool', 'seed');
```

### User Database Registry (in Root DB)

```sql
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE UserDbRegistry (
    UserDbRegistryId INTEGER PRIMARY KEY AUTOINCREMENT,
    UserId INTEGER NOT NULL,
    CompanyId INTEGER,                         -- NULL for app-level
    Category TEXT NOT NULL,                    -- settings, sessions, history, data
    SubCategory TEXT,                          -- Type within category
    EntityId TEXT NOT NULL,
    SequenceNum INTEGER NOT NULL,
    Path TEXT NOT NULL,
    SizeBytes INTEGER DEFAULT 0,
    RecordCount INTEGER DEFAULT 0,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    LastAccessedAt DATETIME,
    ExpiresAt DATETIME,                        -- For session cleanup
    Status TEXT DEFAULT 'active',
    FOREIGN KEY (UserId) REFERENCES User(UserId)
);

CREATE INDEX IdxUserDbRegistry_UserId ON UserDbRegistry(UserId);
CREATE INDEX IdxUserDbRegistry_Category ON UserDbRegistry(Category);
CREATE INDEX IdxUserDbRegistry_ExpiresAt ON UserDbRegistry(ExpiresAt);
```

---

## Go Implementation

### User Database Manager

```go
package userdb

import (
    // ALLOWED: database infrastructure layer — manages user-scoped SQLite databases with dynamic routing
    "database/sql"
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"

    _ "github.com/mattn/go-sqlite3"
)

// ScopeLevel defines the isolation level
type ScopeLevel string

const (
    ScopeLevelApp     ScopeLevel = "app"
    ScopeLevelCompany ScopeLevel = "company"
)

// UserDbManager manages user-scoped databases
type UserDbManager struct {
    rootDb  *sql.DB
    dataDir string
    scope   ScopeLevel
    openDbs map[string]*sql.DB
    mu      sync.RWMutex
}

// UserDbConfig defines configuration for user database manager
type UserDbConfig struct {
    DataDir string
    AppName string
    Scope   ScopeLevel
}

// NewUserDbManager creates a new user database manager
func NewUserDbManager(cfg UserDbConfig) apperror.Result[*UserDbManager] {
    rootPath := filepath.Join(cfg.DataDir, "root.db")
    rootDb, err := sql.Open("sqlite3", rootPath)
    if err != nil {
        return nil, apperror.Wrap(
            err,
            ErrDbOpen,
            "open root database",
        ).WithPath(rootPath)
    }

    manager := &UserDbManager{
        rootDb:  rootDb,
        dataDir: cfg.DataDir,
        scope:   cfg.Scope,
        openDbs: make(map[string]*sql.DB),
    }

    return manager, nil
}

// GetUserDb returns a database for a specific user
func (m *UserDbManager) GetUserDb(appName, userId, category, entityId string, companySlug ...string) apperror.Result[*sql.DB] {
    m.mu.Lock()
    defer m.mu.Unlock()

    // Build path based on scope
    var dbPath string
    if m.scope == ScopeLevelCompany && len(companySlug) > 0 {
        dbPath = m.buildCompanyUserPath(appName, companySlug[0], userId, category, entityId)
    } else {
        dbPath = m.buildAppUserPath(appName, userId, category, entityId)
    }

    // Check if already open
    key := dbPath
    if db, ok := m.openDbs[key]; ok {
        return db, nil
    }

    // Create directory structure
    fullPath := filepath.Join(m.dataDir, dbPath)
    if err := pathutil.EnsureDir(filepath.Dir(fullPath), 0755); err != nil {
        return nil, apperror.Wrap(
            err,
            ErrDbDirCreate,
            "create user database directory",
        ).WithContext("dir", filepath.Dir(fullPath))
    }

    // Open database
    db, err := sql.Open("sqlite3", fullPath)
    if err != nil {
        return nil, apperror.Wrap(
            err,
            ErrDbOpen,
            "open user database",
        ).WithPath(fullPath)
    }

    // Configure database
    m.configureDb(db)

    // Register in root DB
    m.registerUserDb(userId, category, entityId, dbPath)

    m.openDbs[key] = db
    return db, nil
}

func (m *UserDbManager) buildAppUserPath(appName, userId, category, entityId string) string {
    return filepath.Join(appName, "users", userId, category, entityId+".db")
}

func (m *UserDbManager) buildCompanyUserPath(appName, companySlug, userId, category, entityId string) string {
    return filepath.Join(appName, "companies", companySlug, "users", userId, category, entityId+".db")
}

func (m *UserDbManager) configureDb(db *sql.DB) {
    db.Exec("PRAGMA journal_mode=WAL")
    db.Exec("PRAGMA busy_timeout=5000")
    db.Exec("PRAGMA foreign_keys=ON")
}

func (m *UserDbManager) registerUserDb(userId, category, entityId, path string) error {
    _, err := m.rootDb.Exec(`
        INSERT OR REPLACE INTO UserDbRegistry 
        (UserId, Category, EntityId, Path, SequenceNum, UpdatedAt)
        VALUES (?, ?, ?, ?, 1, CURRENT_TIMESTAMP)
    `, userId, category, entityId, path)

    return err
}

// GetUserSettings returns the user's settings database
func (m *UserDbManager) GetUserSettings(appName, userId string, companySlug ...string) apperror.Result[*sql.DB] {
    return m.GetUserDb(appName, userId, "settings", "config", companySlug...)
}

// GetUserSession returns a specific session database
func (m *UserDbManager) GetUserSession(appName, userId, sessionId string, companySlug ...string) apperror.Result[*sql.DB] {
    return m.GetUserDb(appName, userId, "sessions", sessionId, companySlug...)
}

// ListUserDatabases returns all databases for a user
func (m *UserDbManager) ListUserDatabases(userId string) apperror.Result[[]UserDatabase] {
    rows, err := m.rootDb.Query(`
        SELECT UserDbRegistryId, UserId, Category, EntityId, Path, SizeBytes, CreatedAt, LastAccessedAt
        FROM UserDbRegistry
        WHERE UserId = ? AND Status = 'active'
        ORDER BY Category, CreatedAt DESC
    `, userId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var dbs []UserDatabase
    for rows.Next() {
        var db UserDatabase
        if err := rows.Scan(
            &db.UserDbRegistryId, &db.UserId, &db.Category, &db.EntityId,
            &db.Path, &db.SizeBytes, &db.CreatedAt, &db.LastAccessedAt,
        ); err != nil {
            return nil, err
        }
        dbs = append(dbs, db)
    }
    return dbs, nil
}

// DeleteUserData removes all databases for a user
func (m *UserDbManager) DeleteUserData(appName, userId string, companySlug ...string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    // Build user directory path
    var userDir string
    if m.scope == ScopeLevelCompany && len(companySlug) > 0 {
        userDir = filepath.Join(m.dataDir, appName, "companies", companySlug[0], "users", userId)
    } else {
        userDir = filepath.Join(m.dataDir, appName, "users", userId)
    }

    // Close any open databases for this user
    for key, db := range m.openDbs {
        if filepath.HasPrefix(key, userDir) {
            db.Close()
            delete(m.openDbs, key)
        }
    }

    // Remove from registry
    _, err := m.rootDb.Exec(`
        UPDATE UserDbRegistry SET Status = 'deleted', UpdatedAt = CURRENT_TIMESTAMP
        WHERE UserId = ?
    `, userId)
    if err != nil {
        return err
    }

    // Remove directory
    return pathutil.RemoveAll(userDir)
}

// UserDatabase represents a user database entry
type UserDatabase struct {
    UserDbRegistryId int64
    UserId          string
    Category         string
    EntityId         string
    Path             string
    SizeBytes        int64
    CreatedAt        time.Time
    LastAccessedAt   *time.Time
}

// Close closes all open databases
func (m *UserDbManager) Close() error {
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

### App-Level User Isolation

```go
// Create manager for app-level user isolation
manager, err := userdb.NewUserDbManager(userdb.UserDbConfig{
    DataDir: "./data",
    AppName: "myapp",
    Scope:   userdb.ScopeLevelApp,
})
if err != nil {
    log.Fatal(err)
}
defer manager.Close()

// Get user's settings database
settingsDb, err := manager.GetUserSettings("myapp", "user_123")
if err != nil {
    log.Fatal(err)
}

// Initialize settings schema
settingsDb.Exec(`
    CREATE TABLE IF NOT EXISTS Setting (
        Key TEXT PRIMARY KEY,
        Value TEXT NOT NULL,
        UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
    )
`)

// Save user preference
settingsDb.Exec(`INSERT OR REPLACE INTO Settings (Key, Value) VALUES (?, ?)`,
    "Theme", "dark")
```

### Company + User Isolation

```go
// Create manager for company-scoped user isolation
manager, err := userdb.NewUserDbManager(userdb.UserDbConfig{
    DataDir: "./data",
    AppName: "enterprise",
    Scope:   userdb.ScopeLevelCompany,
})
if err != nil {
    log.Fatal(err)
}
defer manager.Close()

// Get user's session database within a company
sessionDb, err := manager.GetUserSession("enterprise", "user_456", "session_abc", "acme-corp")
if err != nil {
    log.Fatal(err)
}

// Use the session database
sessionDb.Exec(`
    CREATE TABLE IF NOT EXISTS SessionData (
        Key TEXT PRIMARY KEY,
        Value TEXT NOT NULL,
        ExpiresAt DATETIME
    )
`)
```

---

## Session Management

### Session Database Schema

```sql
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE SessionMeta (
    SessionMetaId INTEGER PRIMARY KEY AUTOINCREMENT,
    SessionId TEXT UNIQUE NOT NULL,
    UserId TEXT NOT NULL,
    DeviceInfo TEXT,                           -- JSON: browser, OS, IP
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    LastActiveAt DATETIME,
    ExpiresAt DATETIME NOT NULL,
    Status TEXT DEFAULT 'active'               -- active, expired, revoked
);

-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE SessionActivity (
    SessionActivityId INTEGER PRIMARY KEY AUTOINCREMENT,
    Action TEXT NOT NULL,
    Resource TEXT,
    Metadata TEXT,                             -- JSON
    Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IdxSessionActivityTimestamp ON SessionActivity(Timestamp DESC);
```

### Session Cleanup

```go
// CleanupExpiredSessions removes expired session databases
func (m *UserDbManager) CleanupExpiredSessions() error {
    rows, err := m.rootDb.Query(`
        SELECT Path FROM UserDbRegistry
        WHERE Category = 'sessions' 
        AND ExpiresAt < CURRENT_TIMESTAMP 
        AND Status = 'active'
    `)
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var path string
        rows.Scan(&path)
        
        // Close if open
        m.mu.Lock()
        if db, ok := m.openDbs[path]; ok {
            db.Close()
            delete(m.openDbs, path)
        }
        m.mu.Unlock()
        
        // Delete file
        pathutil.Remove(filepath.Join(m.dataDir, path))
    }

    // Update registry
    _, err = m.rootDb.Exec(`
        UPDATE UserDbRegistry 
        SET Status = 'expired', UpdatedAt = CURRENT_TIMESTAMP
        WHERE Category = 'sessions' 
        AND ExpiresAt < CURRENT_TIMESTAMP 
        AND Status = 'active'
    `)
    return err
}
```

---

## Privacy & GDPR Compliance

### User Data Export

```go
// ExportUserData creates a GDPR-compliant data export
func (m *UserDbManager) ExportUserData(userId, outputPath string) error {
    dbs, err := m.ListUserDatabases(userId)
    if err != nil {
        return err
    }

    // Create zip with all user data
    zipFile, _ := pathutil.Create(outputPath)
    defer zipFile.Close()
    zipWriter := zip.NewWriter(zipFile)
    defer zipWriter.Close()

    for _, db := range dbs {
        fullPath := filepath.Join(m.dataDir, db.Path)
        writer, _ := zipWriter.Create(db.Path)
        file, _ := pathutil.Open(fullPath)
        io.Copy(writer, file)
        file.Close()
    }

    return nil
}
```

### User Data Deletion (Right to be Forgotten)

```go
// DeleteAllUserData permanently removes all user data
func (m *UserDbManager) DeleteAllUserData(appName, userId string, companySlug ...string) error {
    // Export before deletion (audit trail)
    exportPath := filepath.Join(m.dataDir, "exports", userId+"-final.zip")
    m.ExportUserData(userId, exportPath)

    // Delete all user databases
    return m.DeleteUserData(appName, userId, companySlug...)
}
```

---

## References

- [00-overview.md](../00-overview.md) - Split DB Architecture Overview
- [01-cli-examples.md](./01-cli-examples.md) - CLI-specific Examples
- [04-rbac-casbin.md](./04-rbac-casbin.md) - RBAC with Casbin
