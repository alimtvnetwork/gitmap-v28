# Project Type Detection — Go Metadata

## Purpose

Go projects have additional metadata beyond basic detection. This
metadata is stored in a separate table because it is disjoint from
the core `DetectedProjects` record and only applies to Go projects.

---

## Collected Metadata

| Field            | Source              | Notes                              |
|------------------|---------------------|------------------------------------|
| GoModPath        | `go.mod` location   | Absolute path to `go.mod` file     |
| GoSumPath        | `go.sum` location   | Absolute path to `go.sum` file     |
| ModuleName       | `go.mod` content    | Parsed from `module` directive     |
| GoVersion        | `go.mod` content    | Parsed from `go` directive         |

---

## Runnable File Detection

### What Is a Runnable File

In Go, a runnable (executable) is any `main.go` file inside a
`cmd/` directory structure. The standard Go project layout uses:

```
cmd/
├── server/
│   └── main.go       ← runnable
├── cli/
│   └── main.go       ← runnable
└── worker/
    └── main.go       ← runnable
```

### Detection Rules

1. Starting from the detected Go project root, look for a `cmd/`
   directory.
2. For each immediate subdirectory inside `cmd/`, check for a
   `main.go` file at any depth (typically `cmd/<name>/main.go`
   or `cmd/<name>/main/main.go`).
3. Each `main.go` found is recorded as a runnable file.
4. The runnable name is derived from the `cmd/` subdirectory name
   (e.g., `cmd/server/main.go` → runnable name `server`).

### Alternative Runnable Patterns

| Pattern                   | Detection                         |
|---------------------------|-----------------------------------|
| `cmd/<name>/main.go`      | Standard layout                   |
| `cmd/<name>/main/main.go` | Nested main package               |
| `main.go` at project root | Root-level executable             |

A `main.go` at the project root (same level as `go.mod`) is also
recorded as a runnable with the name derived from the module name
or directory name.

---

## GoProjectMetadata Table

| Column            | Type    | Constraints                                        | Notes                        |
|-------------------|---------|----------------------------------------------------|------------------------------|
| Id                | TEXT    | PRIMARY KEY                                        | UUID                         |
| DetectedProjectId | TEXT    | NOT NULL, FK → DetectedProjects(Id) ON DELETE CASCADE | Link to parent project    |
| GoModPath         | TEXT    | NOT NULL                                           | Absolute path to go.mod      |
| GoSumPath         | TEXT    | DEFAULT ''                                         | Absolute path to go.sum      |
| ModuleName        | TEXT    | NOT NULL                                           | e.g., `github.com/user/repo` |
| GoVersion         | TEXT    | DEFAULT ''                                         | e.g., `1.22`                 |

**Unique constraint:** `(DetectedProjectId)` — one metadata row per
detected Go project.

---

## GoRunnableFiles Table

| Column            | Type    | Constraints                                            | Notes                     |
|-------------------|---------|--------------------------------------------------------|---------------------------|
| Id                | TEXT    | PRIMARY KEY                                            | UUID                      |
| GoMetadataId      | TEXT    | NOT NULL, FK → GoProjectMetadata(Id) ON DELETE CASCADE | Link to Go metadata       |
| RunnableName      | TEXT    | NOT NULL                                               | e.g., `server`, `cli`     |
| FilePath          | TEXT    | NOT NULL                                               | Absolute path to main.go  |
| RelativePath      | TEXT    | NOT NULL                                               | Path relative to project  |

**Unique constraint:** `(GoMetadataId, RelativePath)` — one entry
per runnable per project.

---

## SQL Statements

### Create GoProjectMetadata

```sql
CREATE TABLE IF NOT EXISTS GoProjectMetadata (
    Id                TEXT PRIMARY KEY,
    DetectedProjectId TEXT NOT NULL UNIQUE
        REFERENCES DetectedProjects(Id) ON DELETE CASCADE,
    GoModPath         TEXT NOT NULL,
    GoSumPath         TEXT DEFAULT '',
    ModuleName        TEXT NOT NULL,
    GoVersion         TEXT DEFAULT ''
)
```

### Create GoRunnableFiles

```sql
CREATE TABLE IF NOT EXISTS GoRunnableFiles (
    Id                TEXT PRIMARY KEY,
    GoMetadataId      TEXT NOT NULL
        REFERENCES GoProjectMetadata(Id) ON DELETE CASCADE,
    RunnableName      TEXT NOT NULL,
    FilePath          TEXT NOT NULL,
    RelativePath      TEXT NOT NULL,
    UNIQUE(GoMetadataId, RelativePath)
)
```

### Upsert GoProjectMetadata

```sql
INSERT INTO GoProjectMetadata (Id, DetectedProjectId, GoModPath,
    GoSumPath, ModuleName, GoVersion)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(DetectedProjectId) DO UPDATE SET
    GoModPath=excluded.GoModPath,
    GoSumPath=excluded.GoSumPath,
    ModuleName=excluded.ModuleName,
    GoVersion=excluded.GoVersion
```

### Upsert GoRunnableFiles

```sql
INSERT INTO GoRunnableFiles (Id, GoMetadataId, RunnableName,
    FilePath, RelativePath)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT(GoMetadataId, RelativePath) DO UPDATE SET
    RunnableName=excluded.RunnableName,
    FilePath=excluded.FilePath
```

### Query Go Metadata

```sql
SELECT gm.Id, gm.DetectedProjectId, gm.GoModPath, gm.GoSumPath,
    gm.ModuleName, gm.GoVersion
FROM GoProjectMetadata gm
WHERE gm.DetectedProjectId = ?
```

### Query Go Runnables

```sql
SELECT gr.Id, gr.GoMetadataId, gr.RunnableName, gr.FilePath,
    gr.RelativePath
FROM GoRunnableFiles gr
WHERE gr.GoMetadataId = ?
ORDER BY gr.RunnableName
```

---

## JSON Output Extension

When writing `go-projects.json`, each Go project record includes
additional fields:

```json
{
  "id": "uuid",
  "repoId": "uuid",
  "repoName": "my-api",
  "projectType": "go",
  "projectName": "github.com/user/my-api",
  "absolutePath": "/home/user/repos/my-api",
  "repoPath": "/home/user/repos/my-api",
  "relativePath": ".",
  "primaryIndicator": "go.mod",
  "detectedAt": "2026-03-11T09:54:00Z",
  "goMetadata": {
    "goModPath": "/home/user/repos/my-api/go.mod",
    "goSumPath": "/home/user/repos/my-api/go.sum",
    "moduleName": "github.com/user/my-api",
    "goVersion": "1.22",
    "runnables": [
      {
        "name": "server",
        "filePath": "/home/user/repos/my-api/cmd/server/main.go",
        "relativePath": "cmd/server/main.go"
      },
      {
        "name": "worker",
        "filePath": "/home/user/repos/my-api/cmd/worker/main.go",
        "relativePath": "cmd/worker/main.go"
      }
    ]
  }
}
```
