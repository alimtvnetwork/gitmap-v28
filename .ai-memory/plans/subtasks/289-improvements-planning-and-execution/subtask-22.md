# Subtask 22: Database Path Introspection & Accessor

> **Subtask ID:** Subtask-22  
> **Parent Plan:** `.ai-memory/plans/pending/289-improvements-planning-and-execution.md`  
> **Component Spec:** `02-spec/21-app/289-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/store/store.go`  
> **Status:** PENDING  

---

## 1. Objectives & Scope
1. Persist the resolved SQLite database file path onto the `DB` struct in `cli/store/store.go`.
2. Add a thread-safe, public getter method `GetDatabasePath() string` on `*DB`.
3. Eliminate duplicate path reconstruction logic across CLI diagnostics, backup helpers, and agent inspection tools.
4. Maintain nil safety when the method is invoked on an uninitialized or nil database pointer.

---

## 2. Detailed Technical Plan

### 2.1 Struct Modification in `cli/store/store.go`
- Update `DB` struct definition:
  ```go
  // DB wraps the SQLite database connection.
  type DB struct {
  	conn  *sql.DB
  	path  string
  	dbDir string
  	isMem bool
  }
  ```

### 2.2 Connection Initialization Update
- Update `openDBConnection` to store `dbPath`:
  ```go
  func openDBConnection(connPath, dbPath, dbDir string, isMem bool) (*DB, error) {
  	conn, err := sql.Open("sqlite", connPath)
  	if err != nil {
  		releaseLockIfNotMem(dbDir, isMem)

  		return nil, fmt.Errorf(constants.ErrDBOpen, dbPath, err)
  	}

  	if err := ConfigureSQLiteConn(conn); err != nil {
  		return closeOnConfigError(conn, dbDir, isMem, err)
  	}

  	return &DB{conn: conn, path: dbPath, dbDir: dbDir, isMem: isMem}, nil
  }
  ```

### 2.3 Getter Method Implementation
- Add `GetDatabasePath` getter on `*DB`:
  ```go
  // GetDatabasePath returns the filesystem path of the active SQLite database file.
  func (db *DB) GetDatabasePath() string {
  	if db == nil {
  		return ""
  	}

  	return db.path
  }
  ```

---

## 3. Strict Coding Guidelines Compliance
- **Rule 2 (Function Limit):** `GetDatabasePath` is 7 lines (<= 15 lines max).
- **Rule 3 (Zero Nested Ifs):** Single guard check for nil pointer.
- **Rule 6 (Vertical Whitespace):** Single blank line after nil guard and before return.

---

## 4. Verification Criteria
- [ ] `DB.path` is initialized on all non-memory and memory database opens.
- [ ] `db.GetDatabasePath()` returns the exact file path passed to `OpenAt` or resolved by `Open`.
- [ ] Nil pointer call `(*DB)(nil).GetDatabasePath()` returns `""` without panicking.
- [ ] Zero compiler builds and zero tests executed during planning.
