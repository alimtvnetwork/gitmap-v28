# Specification 16 - Part 8: Database Schemas & Migrations

## 1. Migration Overview

To support default work directory tracking, interactive macro recordings, and cluster node diagnostics, the SQLite database layer (`gitmap/db/`) defines the following unified tables.

## 2. Table Definitions

### 2.1 Scan Folders & Work Directories (`scan_folders`)
```sql
CREATE TABLE IF NOT EXISTS scan_folders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path TEXT UNIQUE NOT NULL,
    alias TEXT,
    is_default INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_scanned_at TIMESTAMP,
    repo_count INTEGER NOT NULL DEFAULT 0
);
```

### 2.2 Interactive Macros (`macros` & `macro_steps`)
```sql
CREATE TABLE IF NOT EXISTS macros (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    total_steps INTEGER NOT NULL DEFAULT 0,
    tags TEXT
);

CREATE TABLE IF NOT EXISTS macro_steps (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    macro_id INTEGER NOT NULL,
    step_num INTEGER NOT NULL,
    command_line TEXT NOT NULL,
    working_dir TEXT,
    continue_on_error INTEGER NOT NULL DEFAULT 0,
    timeout_seconds INTEGER NOT NULL DEFAULT 300,
    FOREIGN KEY(macro_id) REFERENCES macros(id) ON DELETE CASCADE
);
```

### 2.3 Cluster Nodes & Connectivity (`cluster_nodes`)
```sql
CREATE TABLE IF NOT EXISTS cluster_nodes (
    id TEXT PRIMARY KEY,
    display_id INTEGER UNIQUE NOT NULL,
    alias TEXT,
    ip TEXT NOT NULL,
    port INTEGER NOT NULL DEFAULT 9999,
    role TEXT NOT NULL DEFAULT 'client', -- 'server' or 'client'
    os TEXT NOT NULL DEFAULT 'windows',
    is_active INTEGER NOT NULL DEFAULT 1,
    status TEXT NOT NULL DEFAULT 'online', -- 'online', 'offline', 'unreachable'
    last_heartbeat TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    auth_token TEXT,
    password_hash TEXT
);
```
