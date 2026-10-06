# Issue Domain 04: SQLite and Database Locking Issues

- **Domain:** Database Concurrency, Foreign Keys, and Split-DB
- **Status:** Consolidated Problem & Resolution Matrix

## 1. Database Locking under Worker Concurrency
- **Symptoms:** `sqlite3: database is locked` (code 5) during high-concurrency mass pull operations.
- **Root Cause:** Multiple goroutines attempted simultaneous writes on a single SQLite handle without WAL mode.
- **Resolution:** Enabled WAL mode (`PRAGMA journal_mode=WAL;`), set busy timeout to 5000ms, and implemented channel-based batching.

## 2. Foreign Key Constraint Violations on DB Reset
- **Symptoms:** `gitmap db reset` crashed with foreign key constraint errors across `DetectedProjects`.
- **Root Cause:** Child rows remained referenced when parent entities were truncated.
- **Resolution:** Added `PRAGMA foreign_keys = ON;` and structured cascading deletion order.
