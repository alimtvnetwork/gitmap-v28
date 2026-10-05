# Subtask 05: Pull Error Logging, Stack Trace Capture & Diagnostic Hints

> **Parent Plan:** `219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging`  
> **Status:** READY  
> **Target Files:**  
> - `cli/store/pull_split_db_errors.go`  
> - `cli/cmdpull/pull_efficient.go`  
> - `cli/cmdpull/pull_db_sync.go`  
> - `cli/cmdpull/pull_error_logger.go` *(new helper)*  
> - `cli/cmdpullerror/pull_error_cmd.go`  
> - `cli/cmdpullerror/pull_error_render.go`  
> - `cli/store/pull_split_db_errors_test.go`  

---

## 1. Objectives

1. **RFC3339 DateTime Parsing Fix:**
   Resolve silent zero-time bugs in `cli/store/pull_split_db_errors.go` by replacing single-format datetime parsing with a resilient multi-layout parser supporting RFC3339, RFC3339Nano, ISO 8601, and SQLite datetime strings.
2. **NodeID & Stack Trace Population:**
   Ensure `NodeID` and `StackTrace` fields are rigorously populated when writing failure records to `gitmap-pull.db` via `recordPullErrorsToDB` (`pull_efficient.go`) and `syncPullFailuresToDB` (`pull_db_sync.go`).
3. **Structured Persistent File Logging:**
   Implement thread-safe JSONL structured error logging to `.gitmap/logs/pull-errors.log` so failures can be inspected and ingested independently of SQLite queries.
4. **Diagnostic Command Guidance:**
   Standardize terminal output to guide users and agents to inspect detailed stack traces using `gitmap pull-error <repo>` (or `gitmap pe`).
5. **Database Roundtrip Verification:**
   Verify through unit tests that `InsertPullError` and `QueryLatestPullErrors` preserve all fields (including `NodeID`, `StackTrace`, and RFC3339 timestamps) without data truncation or zero-time conversion.

---

## 2. Implementation Steps

### Step 2.1: Robust DateTime Parsing (`cli/store/pull_split_db_errors.go`)

1. Implement `parseFlexibleDBTimestamp`:
   ```go
   func parseFlexibleDBTimestamp(s string) time.Time {
       clean := strings.TrimSpace(s)
       if clean == "" {
           return time.Now().UTC()
       }
       layouts := []string{
           time.RFC3339Nano,
           time.RFC3339,
           "2006-01-02 15:04:05.999999999",
           "2006-01-02 15:04:05",
           "2006-01-02T15:04:05",
           "2006-01-02",
       }
       for _, layout := range layouts {
           if t, err := time.Parse(layout, clean); err == nil {
               return t.UTC()
           }
       }
       return time.Now().UTC()
   }
   ```
2. In `assignNullableFields`, replace:
   ```go
   rec.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdStr)
   ```
   with:
   ```go
   rec.CreatedAt = parseFlexibleDBTimestamp(createdStr)
   ```

### Step 2.2: Populate NodeID and StackTrace in Pull Split-DB

1. In `cli/cmdpull/pull_efficient.go` (`recordPullErrorsToDB`):
   - Resolve node ID dynamically:
     ```go
     nodeID := resolveLocalNodeIdentifier()
     ```
   - Populate `NodeID: nodeID` and `StackTrace: s.StackTrace` (or extract error stack details if empty).
2. In `cli/cmdpull/pull_db_sync.go` (`syncPullFailuresToDB`):
   - Populate:
     ```go
     rec := store.PullErrorRecord{
         RepoSlug:       state.RepoName,
         RepoPath:       state.RepoPath,
         NodeID:         resolveLocalNodeIdentifier(),
         NodeVersion:    constants.Version,
         ErrorType:      string(state.Step),
         ErrorText:      state.ErrorMsg,
         StackTrace:     state.StackTrace,
         RemediationCmd: ResolvePullRemediationHint(state),
         CreatedAt:      time.Now().UTC(),
     }
     ```

### Step 2.3: Structured JSONL File Logger (`cli/cmdpull/pull_error_logger.go`)

1. Author helper to append pull failures to `.gitmap/logs/pull-errors.log`:
   ```go
   func LogPullErrorToFile(rec store.PullErrorRecord) error {
       logDir := filepath.Join(".gitmap", "logs")
       if err := os.MkdirAll(logDir, 0755); err != nil {
           return err
       }
       logPath := filepath.Join(logDir, "pull-errors.log")
       f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
       if err != nil {
           return err
       }
       defer f.Close()

       data, err := json.Marshal(rec)
       if err != nil {
           return err
       }
       _, err = f.Write(append(data, '\n'))
       return err
   }
   ```
2. Invoke `LogPullErrorToFile` whenever a pull error record is inserted into SQLite.

### Step 2.4: Terminal Diagnostic Hint Standardization

1. Verify that `renderSingleFailedItem` outputs the terminal hint:
   ```text
   └── To inspect stack trace: gitmap pull-error <repo> (or: gitmap pe)
   ```
2. Verify that `gitmap pull-error <repo>` cleanly displays `Node ID`, `GitMap Version`, `Timestamp`, and `Stack Trace`.

### Step 2.5: Unit Test Suite Verification

1. Add tests in `cli/store/pull_split_db_errors_test.go`:
   - Test `parseFlexibleDBTimestamp` with:
     - RFC3339 string `"2026-10-05T02:44:45Z"`
     - RFC3339Nano string `"2026-10-05T02:44:45.123456789Z"`
     - SQLite format `"2026-10-05 02:44:45"`
     - Invalid/empty string fallback
   - Test DB insert and retrieval asserting `NodeID`, `StackTrace`, and `CreatedAt` match input exactly.

---

## 3. Acceptance Criteria

- [ ] `parseFlexibleDBTimestamp` handles RFC3339, RFC3339Nano, and SQLite datetime strings without zero-time fallbacks.
- [ ] `InsertPullError` persists `NodeID` and `StackTrace` correctly into `pull_errors` table in `gitmap-pull.db`.
- [ ] Pull failures are written to `.gitmap/logs/pull-errors.log` as valid single-line JSON records.
- [ ] Terminal failure cards render `└── To inspect stack trace: gitmap pull-error <repo> (or: gitmap pe)`.
- [ ] `go test ./cli/store/...` and `go test ./cli/cmdpullerror/...` pass with 100%.
