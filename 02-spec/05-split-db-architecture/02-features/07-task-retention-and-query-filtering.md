# Task Retention & Dynamic Query Filtering Specification

## 1. Overview

This specification standardizes automated retention policies and dynamic query filtering across the SQLite task database subsystem (`04-code/golang/pkg/applogger/sqlitelogger`).

## 2. Retention Management

### 2.1 Time-Based Pruning

`SplitDBManager.PruneTasks(maxAge time.Duration) (int, *appfault.AppError)`:
- Discovers all task databases in `tasksDir`.
- Compares each file's modification time against `cutoff := time.Now().Add(-maxAge)`.
- If older than cutoff:
  1. Closes any cached connection in `m.taskDbs`.
  2. Safely removes the `.db` file as well as associated `-wal` and `-shm` temporary files.
  3. Increments pruned count.

### 2.2 Capacity-Based Pruning

`SplitDBManager.PruneTaskCount(maxDbs int) (int, *appfault.AppError)`:
- If discovered database count exceeds `maxDbs`:
  1. Sorts task databases by modification time ascending (oldest first).
  2. Deletes oldest database files until total database count equals `maxDbs`.

## 3. Dynamic Query Filtering

### 3.1 Filter Specification

The `FilterOptions` struct controls log retrieval:
```go
type FilterOptions struct {
	Level     string
	Limit     int
	Offset    int
	StartTime string
	EndTime   string
}
```

### 3.2 Dynamic SQL Generation

`queryLogs(db *sql.DB, filter FilterOptions)` generates an parameterized SQL query:
- Base: `SELECT id, task_id, timestamp, level, message, caller, fields_json, stack_trace, duration_ms, status FROM logs WHERE 1=1`
- If `Level` provided: `AND level = ?`
- If `StartTime` provided: `AND timestamp >= ?`
- If `EndTime` provided: `AND timestamp <= ?`
- Ordering & Pagination: `ORDER BY id DESC LIMIT ? OFFSET ?`
- Default limit is 100 if unspecified or negative.
