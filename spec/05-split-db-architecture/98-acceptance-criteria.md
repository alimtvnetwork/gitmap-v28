# Split DB Architecture: Acceptance Criteria

**Version:** 3.2.0  
**Status:** Active  
**Updated:** 2026-04-16  
**Format:** GIVEN/WHEN/THEN (E2E-test-ready)

---

## Database Hierarchy (00-overview.md)

### SD-01: Root DB Initialization

**GIVEN** a CLI tool starts for the first time with no existing data directory  
**WHEN** the initialization sequence runs  
**THEN** the root database file is created at `data/{toolname}.db`  
**AND** all required tables (apps, settings, migrations) are auto-migrated via GORM  
**AND** WAL mode is enabled on the SQLite connection

**Edge Cases:**
- **GIVEN** the `data/` directory does not exist **WHEN** initialization runs **THEN** the directory is created with 0755 permissions before creating the DB
- **GIVEN** the root DB file exists but is corrupted **WHEN** initialization runs **THEN** the error is reported with error code and the tool exits cleanly without overwriting
- **GIVEN** the filesystem is read-only **WHEN** initialization runs **THEN** a clear error "Cannot create database: filesystem is read-only" is returned

### SD-02: App DB Dynamic Creation

**GIVEN** a root DB exists with the apps table  
**WHEN** a new app/project is registered (e.g., `brun init myapp`)  
**THEN** a new directory `data/{appName}/` is created  
**AND** an app-specific database file is created (e.g., `search.db`)  
**AND** the app is registered in the root DB's apps table with its metadata

**Edge Cases:**
- **GIVEN** an app with the same name already exists **WHEN** registration is attempted **THEN** a duplicate error is returned without modifying existing data
- **GIVEN** the app name contains invalid filesystem characters **WHEN** registration is attempted **THEN** a validation error is returned listing the invalid characters

### SD-03: Item DB Creation

**GIVEN** an app DB exists  
**WHEN** a new item (e.g., search result set, build run) is created  
**THEN** an item-specific DB file is created at `data/{appName}/{type}/{seq}-{slug}.db`  
**AND** the item is registered in the app DB's items table

**Edge Cases:**
- **GIVEN** the slug exceeds 50 characters **WHEN** the item DB is created **THEN** the slug is truncated to 50 characters with a hash suffix for uniqueness
- **GIVEN** disk space is insufficient **WHEN** creation is attempted **THEN** a storage error is returned with available and required space
- **GIVEN** two concurrent goroutines create items with the same sequence number **WHEN** both write to the items table **THEN** the UNIQUE constraint prevents duplicates and one goroutine retries with the next sequence number

---

## Reset API Standard (02-reset-api-standard.md)

### RA-01: Reset Request (Step 1)

**GIVEN** the CLI is running and the API is accessible  
**WHEN** a POST request is sent to `/api/v1/reset/request` with `{ "Scope": "all" }`  
**THEN** a response is returned with `ResetId`, `Scope`, `ExpiresAt` (now + 5 minutes), and `AffectedItems` preview  
**AND** the reset token is stored in memory with a 5-minute TTL

**Edge Cases:**
- **GIVEN** an invalid scope value is provided **WHEN** the request is sent **THEN** a 400 error is returned listing valid scope options
- **GIVEN** a reset request is already pending **WHEN** a new request is sent **THEN** the previous request is invalidated and a new token is issued

### RA-02: Reset Confirmation (Step 2)

**GIVEN** a valid reset token exists from Step 1  
**WHEN** a POST request is sent to `/api/v1/reset/confirm` with the `ResetId`  
**THEN** all data matching the scope is deleted  
**AND** the response includes `Status: "completed"`, `DeletedDatabases` count, and `Duration`

**Edge Cases:**
- **GIVEN** the reset token has expired (>5 minutes) **WHEN** confirmation is attempted **THEN** a 410 Gone error is returned with message "Reset token expired, please request again"
- **GIVEN** an invalid/unknown `ResetId` is provided **WHEN** confirmation is attempted **THEN** a 404 error is returned
- **GIVEN** scope is `all` **WHEN** reset is confirmed **THEN** all app databases and the root DB tables are wiped but the root DB file and schema remain intact
- **GIVEN** a database file is locked by another process during reset **WHEN** deletion is attempted **THEN** the locked file is skipped, reported in the response, and a partial completion status is returned

### RA-03: Per-Module Reset

**GIVEN** the CLI supports multiple modules (e.g., search, cache, settings)  
**WHEN** a reset request is sent with `{ "Scope": "cache" }`  
**THEN** only cache-related databases and tables are included in `AffectedItems`  
**AND** confirmation only deletes cache data, leaving other modules intact

---

## RBAC with Casbin (04-rbac-casbin.md)

### RB-01: Policy Enforcement

**GIVEN** RBAC is enabled with a Casbin policy model  
**WHEN** a user with role `viewer` attempts a `DELETE` operation  
**THEN** the request is denied with a 403 Forbidden response  
**AND** the denial is logged with user, role, resource, and action

**Edge Cases:**
- **GIVEN** no policy file exists **WHEN** the server starts **THEN** a default deny-all policy is loaded and a warning is logged
- **GIVEN** the policy file is modified at runtime **WHEN** the hot-reload interval triggers **THEN** the updated policies take effect without server restart

### RB-02: Role Assignment

**GIVEN** an admin user is authenticated  
**WHEN** they assign role `editor` to another user via API  
**THEN** the role is persisted in the database  
**AND** subsequent requests from that user are evaluated against the new role

---

## User-Scoped Isolation (05-user-scoped-isolation.md)

### US-01: Data Isolation

**GIVEN** multi-user mode is enabled  
**WHEN** user A creates data  
**THEN** user B cannot access user A's data through any API endpoint  
**AND** user A's databases are stored in `data/users/{userId}/`

**Edge Cases:**
- **GIVEN** user A's session token is used by an attacker **WHEN** attempting to access user B's path **THEN** path traversal is blocked and a security event is logged
- **GIVEN** a user is deleted **WHEN** cleanup runs **THEN** all user-scoped databases and directories are removed

---

*Wave 6 — Batch 1 (Patched): Split DB acceptance criteria with added edge cases for read-only filesystem, concurrent item creation race conditions, and locked DB files during reset.*
