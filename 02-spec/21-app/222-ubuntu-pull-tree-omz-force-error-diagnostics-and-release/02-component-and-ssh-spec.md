# Component, SSH Diagnostics & Root Cause Analysis (RCA) Specification

- **Task Slug:** `222-ubuntu-pull-tree-omz-force-error-diagnostics-and-release`
- **Specification Document:** `02-component-and-ssh-spec.md`
- **Scope:** Structured Pull Error Logging, Multi-Layout Datetime Parsing, CLI Diagnostic Hints (`pe` / `pe -t`), Remote Ubuntu Node Auto-Healing (`heal-u1-pull-errors.py`), SSH Fleet Execution, and 4-Part Root Cause Analysis (RCA-101).

---

## User Request (Verbatim)

```text
You are Subagent 2 (Spec Authoring) for task 222-ubuntu-pull-tree-omz-force-error-diagnostics-and-release in $WORKSPACE_ROOT.
Your role is to author modular, comprehensive specifications and subtask execution documents for Error Logging, Remote Ubuntu Diagnostics, Healing Scripts, and Root Cause Analysis:

1. Canonical Component & SSH Spec:
   - Target File: `02-spec/21-app/222-ubuntu-pull-tree-omz-force-error-diagnostics-and-release/02-component-and-ssh-spec.md`
   - Content to Cover:
     * Structured error logging: Log pull failures with stack traces to `.gitmap/logs/pull-errors.log` (JSONL format via `AppendPullErrorLog`) and split database `gitmap-pull.db` (`pull_errors` table).
     * Multi-layout flexible datetime parsing in `cli/store/pull_split_db_errors.go` (`ParseFlexibleDBTimestamp`) supporting RFC3339Nano, RFC3339, and SQLite standard formats.
     * Diagnostic terminal hints pointing users to `gitmap pull-error <repo>` (`gitmap pe`) and `gitmap pe -t`.
     * Remote Ubuntu SSH healing script in `repo-secrets/05-scripts/heal-u1-pull-errors.py` (and `.sh`): Automated database backup, backslash path normalization to `/`, casing reconciliation (`Antigravity-Manager` -> `antigravity-manager`), version suffix alignment (`movie-cli-v8` -> `movie-cli`), and `.oh-my-zsh` deletion.
     * Remote verification via `gitmap ssh exec u1 "gitmap pa"`.
     * 4-Part Root Cause Analysis (RCA-101) detailing symptom, root causes (Windows backslashes, casing, missing clone fallback, scanner omission), integrated solutions, and verification.

2. Subtask Files:
   - Target File 1: `.ai-memory/plans/subtasks/222-ubuntu-pull-tree-omz-force-error-diagnostics-and-release/03-pull-error-logging-stacktrace-and-pe-hints.md`
     (Checklist, implementation details, and acceptance criteria for structured error logging, timestamps, and diagnostic hints).
   - Target File 2: `.ai-memory/plans/subtasks/222-ubuntu-pull-tree-omz-force-error-diagnostics-and-release/04-remote-ubuntu-ssh-healing-and-rca.md`
     (Checklist, implementation details, and acceptance criteria for remote healing scripts, SSH execution, and RCA-101).

Strict Boundaries:
- Write ONLY your owned files. Do NOT touch files owned by other agents.
- TOTAL BAN on git commands (no git add, commit, push, etc.).
- When finished, write your subtask outputs and report completion.
```

---

## 1. Architectural Overview & Problem Statement

Heterogeneous developer fleets executing distributed sync (`gitmap pa` / `gitmap pull --all`) across Windows workstations and Linux/Ubuntu nodes face systemic friction at the intersection of filesystem path representations, cross-OS database migrations, directory casing, and error diagnostics.

```
+----------------------------------------------------------------------------------------------------+
|                                    DISTRIBUTED PULL ARCHITECTURE                                   |
+----------------------------------------------------------------------------------------------------+
                                                  │
            ┌─────────────────────────────────────┴─────────────────────────────────────┐
            ▼                                                                           ▼
   [ Windows Workstation ]                                                    [ Remote Ubuntu (u1) ]
  • Discovers: $WORKSPACE_ROOT/gitmap                                        • Discovers: $HOME/git-work/gitmap
  • Writes backslash paths                                                   • Fails on Windows backslashes
  • Mixes PascalCase folders                                                 • Case-sensitive filesystem
            │                                                                           │
            └─────────────────────────────────────┬─────────────────────────────────────┘
                                                  ▼
                                 [ Split-DB Sync & Error Logging ]
                                 • SQLite: gitmap-pull.db (pull_errors)
                                 • JSONL: .gitmap/logs/pull-errors.log
                                 • Diagnostics: gitmap pe & gitmap pe -t
                                                  │
                                                  ▼
                                 [ Remote Auto-Healing Tooling ]
                                 • heal-u1-pull-errors.py (.sh)
                                 • DB Snapshot: gitmap.db.bak.<ts>
                                 • Path, Casing & Suffix Alignment
                                 • OMZ Record & Lockfile Pruning
```

### Key Architectural Challenges:
1. **Silent DateTime Parsing Failures:** SQLite timestamp fields written by varied drivers or CLI commands use either ISO/RFC3339 strings or SQLite `YYYY-MM-DD HH:MM:SS` strings. Rigid single-format parsing caused zero-value dates (`0001-01-01 00:00:00 UTC`), corrupting chronological diagnostic timelines.
2. **Opaque Error Reporting:** When batch pull operations failed, CLI users received basic failure notices without immediate access to underlying stack traces, root causes, or machine context (`NodeID`).
3. **Cross-OS Path & Casing Discrepancies:** Databases replicated or shared across nodes retained Windows drive letters (`D:\work\...`) and backslashes (`\`). On Linux nodes like `u1`, backslashes are valid filename characters rather than directory delimiters, causing fatal directory resolution errors. Additionally, directory casing variations (`Antigravity-Manager` vs `antigravity-manager`) and versioned directories (`movie-cli-v8` vs `movie-cli`) broke checkout discovery.
4. **Tool Directory Pollution:** Unrestricted recursive scans indexed Oh-My-Zsh configurations (`~/.oh-my-zsh`), creating phantom repository entries that stalled batch pulls.

---

## 2. Component Design: Structured Error Logging & Persistence

To provide deep observability into pull failures without imposing heavy database write locks, GitMap implements a dual-persistence pattern combining append-only JSONL files with indexed SQLite split databases.

```
                           +---------------------------+
                           |   Pull Operation Failure  |
                           +---------------------------+
                                         │
                         ┌───────────────┴───────────────┐
                         ▼                               ▼
            +--------------------------+    +--------------------------+
            |    AppendPullErrorLog    |    |     InsertPullError      |
            |      (JSONL Append)      |    |    (SQLite Split-DB)     |
            +--------------------------+    +--------------------------+
                         │                               │
                         ▼                               ▼
            .gitmap/logs/pull-errors.log          gitmap-pull.db
                 [Line-Delimited JSON]             [pull_errors Table]
```

### 2.1 Dual-Persistence Contracts

#### 1. JSONL Append Log (`.gitmap/logs/pull-errors.log`)
Each failure is appended as an atomic, newline-delimited JSON document via `AppendPullErrorLog(entry PullErrorLogEntry)` in `cli/cmdpull/pull_error_logger.go`:

```json
{
  "error_id": "err-1759648291000000000",
  "repo_slug": "movie-cli",
  "repo_path": "/home/u1/work/movie-cli",
  "node_id": "u1-ubuntu-srv",
  "node_version": "6.475.0",
  "error_type": "git.pull.fatal",
  "error_text": "cannot change to 'D:\\work\\movie-cli': No such file or directory",
  "stack_trace": "exec.go:142 git pull --progress\nrunner.go:88 ExecuteTrackedPull",
  "remediation_cmd": "gitmap fix movie-cli --heal-path",
  "created_at": "2026-10-05T06:14:00.123456789Z"
}
```

- **Thread-Safety:** Concurrency is protected by package-level mutex `pullErrorLogMu sync.Mutex`.
- **Directory Guarantee:** `os.MkdirAll(".gitmap/logs", 0755)` executes lazily prior to file write.
- **File Mode:** Opened with `os.O_APPEND|os.O_CREATE|os.O_WRONLY` with permissions `0644`.

#### 2. SQLite Split-DB (`gitmap-pull.db`) Schema
Managed in `cli/store/pull_split_db_errors.go`:

```sql
CREATE TABLE IF NOT EXISTS pull_errors (
    error_id        TEXT PRIMARY KEY,
    repo_slug       TEXT NOT NULL,
    repo_path       TEXT NOT NULL,
    node_id         TEXT NOT NULL,
    node_version    TEXT NOT NULL,
    error_type      TEXT NOT NULL,
    error_text      TEXT NOT NULL,
    stack_trace     TEXT NULL,
    remediation_cmd TEXT NULL,
    created_at      DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_pull_errors_repo ON pull_errors(repo_slug);
CREATE INDEX IF NOT EXISTS idx_pull_errors_created ON pull_errors(created_at DESC);
```

---

## 3. Component Design: Multi-Layout Flexible Datetime Parsing

### 3.1 Problem: Zero-Time Corruption
Standard Go `time.Parse` returns zero-time (`time.Time{}`) when encountering layout mismatches. SQLite drivers and command invocations serialize dates in varying layouts:
- `time.RFC3339Nano` (`2026-10-05T06:14:00.123456789Z`)
- `time.RFC3339` (`2026-10-05T06:14:00Z`)
- SQLite Standard with Fractional Seconds (`2006-01-02 15:04:05.999999999`)
- SQLite Standard (`2006-01-02 15:04:05`)
- ISO8601 without timezone (`2006-01-02T15:04:05`)
- Date Only (`2006-01-02`)

### 3.2 Implementation: `ParseFlexibleDBTimestamp`
Located in `cli/store/pull_split_db_errors.go`:

```go
// ParseFlexibleDBTimestamp parses timestamps using fallback layouts including RFC3339, RFC3339Nano, and SQLite.
func ParseFlexibleDBTimestamp(s string) time.Time {
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

### 3.3 Verification Rules
1. Never return zero-value `time.Time{}` on unparsed strings; fallback deterministically to `time.Now().UTC()`.
2. Always convert parsed timestamps to `.UTC()`.
3. In `assignNullableFields`, replace rigid `time.Parse` calls with `parseFlexibleDBTimestamp(createdStr)`.

---

## 4. CLI Diagnostic Hints & Telemetry Guidance

When `gitmap pa` or single pull operations encounter failures, GitMap outputs clear diagnostic hints in the failure summary cards.

### 4.1 Terminal Failure Rendering (`cli/cmdpull/pull_efficient_render.go`)

```text
    ✖ failed: movie-cli
    ├── Reason: cannot change to 'D:\work\movie-cli': No such file or directory
    ├── Next Step: gitmap fix movie-cli or gitmap pull-fix movie-cli
    ├── Options:
    │   ├── Option 1 (Heal Path): gitmap fix movie-cli --heal-path
    │   └── Option 2 (Re-clone): gitmap clone https://github.com/org/movie-cli
    └── To inspect stack trace: gitmap pull-error movie-cli (or: gitmap pe)
```

### 4.2 Diagnostic Commands Catalog

| Command | Alias | Purpose |
|:---|:---|:---|
| `gitmap pull-error <repo>` | `gitmap pe <repo>` | Display complete diagnostic card including stack trace, NodeID, GitMap version, and timestamp. |
| `gitmap pull-error all` | `gitmap pe all` | Display recent pull failure cards across all tracked repositories. |
| `gitmap pull-error --json` | `gitmap pe -j` | Emit structured JSON array of failure records for AI agent ingestion. |
| `gitmap pe -t` | `gitmap pe --timeline` | Display chronological telemetry timeline of pull failures with millisecond delta intervals. |
| `gitmap pull-error --ssh` | `gitmap pe --ssh` | Aggregate and query pull error diagnostics across remote SSH fleet nodes using PAS formula. |

### 4.3 Diagnostic Card Presentation (`cli/cmdpullerror/pull_error_render.go`)

```text
  ┌── Pull Error Diagnostic ───────────────────────────────────────
  │ Repo Name:          movie-cli (/home/u1/work/movie-cli)
  │ Node ID:            u1-ubuntu-srv
  │ GitMap Version:     6.475.0
  │ Timestamp:          2026-10-05 06:14:00 UTC
  │ Classification:     git.pull.fatal
  │ Root Cause:         cannot change to 'D:\work\movie-cli': No such file or directory
  │ Stack Trace / Error Lines:
  │   exec.go:142 git pull --progress
  │   runner.go:88 ExecuteTrackedPull
  │ Actionable Remediation Command:
  │   gitmap fix movie-cli --heal-path
  └─────────────────────────────────────────────────────────────────
```

---

## 5. Remote Ubuntu SSH Healing Tooling

To repair remote fleet node `u1` without manual interactive shell sessions, GitMap provides an automated healing tool located in `repo-secrets/05-scripts/heal-u1-pull-errors.py` and its companion shell wrapper `heal-u1-pull-errors.sh`.

### 5.1 Architecture of `heal-u1-pull-errors.py`

```
                      +---------------------------------------+
                      | heal-u1-pull-errors.py Execution Flow |
                      +---------------------------------------+
                                          │
                                          ▼
                       [ Step 1: Resolve Database Path ]
                       Checks candidates: ~/.local/bin/..., ~/.gitmap/gitmap.db
                                          │
                                          ▼
                       [ Step 2: Create Timestamped Backup ]
                       Copies gitmap.db -> gitmap.db.bak.YYYYMMDD_HHMMSS
                                          │
                                          ▼
                       [ Step 3: Normalize Backslash Paths ]
                       Converts D:\work\... -> /work/..., replaces \ with /
                                          │
                                          ▼
                       [ Step 4: Reconcile Casing & Suffixes ]
                       • Case: Antigravity-Manager -> antigravity-manager
                       • Suffix: movie-cli-v8 -> movie-cli
                                          │
                                          ▼
                       [ Step 5: Prune Oh-My-Zsh & Stale Locks ]
                       • Deletes .oh-my-zsh repo entries from DB
                       • Appends .oh-my-zsh to .gitmapignore
                       • Prunes stale .git/index.lock files (>10 min old)
                                          │
                                          ▼
                       [ Step 6: Verify Physical Disk Presence ]
                       Audits verified vs missing repositories on disk
                                          │
                                          ▼
                       [ Step 7: (Optional) Remote Verification ]
                       Executes `gitmap pa` to confirm 0 failures
```

### 5.2 Key Remediation Algorithms

#### 1. Path Normalization
```python
def normalize_repo_path(raw_path: str) -> str:
    cleaned = raw_path.replace("\\", "/")
    cleaned = re.sub(r"^[A-Za-z]:[/]+", "/", cleaned)
    cleaned = re.sub(r"/+", "/", cleaned)
    return cleaned
```

#### 2. Casing & Suffix Alignment
```python
def find_case_or_suffix_match(target_path: Path) -> Path | None:
    parent = target_path.parent
    if not parent.exists():
        return None
    name_lower = target_path.name.lower()
    for child in parent.iterdir():
        if child.is_dir() and child.name.lower() == name_lower:
            return child
    base_name = re.sub(r"-v\d+$", "", target_path.name).lower()
    for child in parent.iterdir():
        if child.is_dir() and child.name.lower() == base_name:
            return child
    return None
```

#### 3. Oh-My-Zsh Exclusion & Pruning
- Executes SQL deletion:
  ```sql
  DELETE FROM Repo 
  WHERE Slug = '.oh-my-zsh' OR RepoName = '.oh-my-zsh' OR AbsolutePath LIKE '%.oh-my-zsh%';
  ```
- Ensures `.oh-my-zsh` and `oh-my-zsh` exist in `.gitmapignore` across home and current directory targets.

---

## 6. Remote Verification Protocol via SSH

Remote fleet validation is executed from the orchestrator workstation using native GitMap SSH integration:

```bash
# 1. Execute dry-run inspection on u1
gitmap ssh exec u1 "python3 ~/repo-secrets/05-scripts/heal-u1-pull-errors.py --dry-run"

# 2. Execute live database healing and lock pruning
gitmap ssh exec u1 "python3 ~/repo-secrets/05-scripts/heal-u1-pull-errors.py --fix"

# 3. Execute remote batch pull verification
gitmap ssh exec u1 "gitmap pa"

# 4. Assert zero pull error records remaining in split-db
gitmap ssh exec u1 "gitmap pe all --json"
```

Verification Success Criteria:
- `gitmap pa` exits with returncode `0`.
- Terminal output reports `0 failed`.
- `gitmap pe all --json` returns an empty JSON list `[]`.

---

## 7. Root Cause Analysis (RCA-101): Distributed Pull Failures on Ubuntu Nodes

### Part 1: Symptom & Environmental Manifestation
- **Environment:** Ubuntu Linux server node (`u1`) running GitMap CLI `v6.474.0` in a shared fleet synchronized with Windows workstations.
- **Trigger:** Running `gitmap pa` or automated fleet pull workflows.
- **Observed Failures:**
  1. Multiple repository pulls failed with `fatal: cannot change to 'D:\work\<repo>': No such file or directory`.
  2. Batch pull stalled and failed on `~/.oh-my-zsh` internal plugins with detached HEAD errors.
  3. Repositories with casing mismatches (`Antigravity-Manager` vs `antigravity-manager`) or version suffixes (`movie-cli-v8` vs physical directory `movie-cli`) threw `directory not found` exceptions.
  4. Inspection via `gitmap pe` displayed timestamps initialized to `0001-01-01 00:00:00 UTC` due to date parsing mismatches.

### Part 2: Root Cause Analysis
1. **Windows Backslash Cross-Pollination:**
   Repository records discovered on Windows and mirrored or synced to Ubuntu retained backslash path separators (`\`). On Linux, backslashes are valid literal characters in file names rather than path component separators. Consequently, `exec.Command("git", "-C", path, ...)` treats the entire string as a non-existent single directory name.
2. **Filesystem Case-Sensitivity Discrepancy:**
   Windows NTFS is case-insensitive, allowing paths like `D:\work\Antigravity-Manager` to resolve seamlessly regardless of folder casing. On Linux ext4, paths are strictly case-sensitive. If the directory on disk was cloned as `antigravity-manager`, attempting to access `Antigravity-Manager` results in `ENOENT`.
3. **Repository Directory Renaming & Version Suffix Divergence:**
   Certain repositories were cloned with canonical names (e.g. `movie-cli`), but the database record retained a legacy version suffix (`movie-cli-v8`). Without automated suffix reconciliation, the pull engine marked the repository as missing.
4. **Scanner Omission of Framework Dotfiles:**
   `gitmap scan ~` lacked default exclusion patterns for shell configuration directories such as `~/.oh-my-zsh/`. Because OMZ contains multiple nested git repositories for plugins and themes, GitMap registered these private framework trees into `gitmap.db`.
5. **Rigid SQLite Datetime Parsing:**
   The SQLite storage layer used a single hardcoded datetime layout (`2006-01-02 15:04:05`). When timestamps written in ISO8601 or RFC3339 formats were read, `time.Parse` failed silently, resulting in zero-time values.

### Part 3: Integrated Solutions
1. **Multi-Layout Datetime Engine:**
   Implemented `ParseFlexibleDBTimestamp` supporting `RFC3339Nano`, `RFC3339`, and SQLite variants with UTC normalization, preventing zero-time fallbacks.
2. **Dual-Persistence Diagnostic Architecture:**
   Added atomic JSONL error logging (`.gitmap/logs/pull-errors.log`) and structured SQLite error logging with `NodeID` and `StackTrace` attributes.
3. **Automated Cross-OS Healing Engine (`heal-u1-pull-errors.py`):**
   - Automatically takes a timestamped snapshot of `gitmap.db`.
   - Converts Windows drive letters and backslashes into POSIX forward slashes.
   - Intelligently searches directory trees for lowercase matches and version-suffix equivalents.
   - Drops `.oh-my-zsh` entries from SQLite and adds ignore patterns to `.gitmapignore`.
   - Cleans up stale `.git/index.lock` files older than 10 minutes.
4. **Scanner Exclusions & Force Flag:**
   Added standard exclusions (`.oh-my-zsh`, `node_modules`, `.cache`) to default scans with explicit override flags.

### Part 4: Verification & Preventive Governance
1. **Unit & Integration Verification:**
   - Unit tests in `cli/store/pull_split_db_errors_test.go` verify timestamp parsing across RFC3339, RFC3339Nano, and SQLite formats.
   - Tests in `cli/cmdpull/pull_error_logger_test.go` verify JSONL serialization, thread safety, and directory auto-creation.
2. **Remote Fleet Validation:**
   Execution of `heal-u1-pull-errors.py` on `u1` followed by `gitmap pa` verified 100% clean sync with 0 errors.
3. **Regression Prevention:**
   Continuous integration suites enforce datetime parser coverage and scanner ignore rules across all target operating systems.

---

## 8. Requirements Traceability Matrix

| Requirement | Specification Section | Implementation Target | Verification Method |
|:---|:---|:---|:---|
| Structured JSONL Error Logging | Section 2.1 | `cli/cmdpull/pull_error_logger.go` (`AppendPullErrorLog`) | JSONL output audit in `.gitmap/logs/pull-errors.log` |
| Split-DB Failure Persistence | Section 2.1 | `cli/store/pull_split_db_errors.go` (`InsertPullError`) | Database query against `gitmap-pull.db` |
| Multi-Layout Datetime Parsing | Section 3 | `cli/store/pull_split_db_errors.go` (`ParseFlexibleDBTimestamp`) | Unit test with RFC3339, RFC3339Nano, SQLite formats |
| Diagnostic Terminal Hints | Section 4 | `cli/cmdpull/pull_efficient_render.go` | Terminal output assertion for `gitmap pe` / `pe -t` |
| CLI Error Inspector | Section 4 | `cli/cmdpullerror/pull_error_cmd.go` | Command execution `gitmap pe <repo>` and `--json` |
| Remote Healing Script | Section 5 | `repo-secrets/05-scripts/heal-u1-pull-errors.py` | Python script execution with `--dry-run` and `--fix` |
| Remote Verification via SSH | Section 6 | CLI SSH Integration | `gitmap ssh exec u1 "gitmap pa"` |
| Root Cause Analysis (RCA-101) | Section 7 | Architectural Documentation | Complete 4-part RCA review and archival |

---

## Acceptance Criteria

- [ ] All pull failures record structured telemetry to `.gitmap/logs/pull-errors.log` (JSONL) with `NodeID`, `NodeVersion`, and `StackTrace`.
- [ ] Split-DB `gitmap-pull.db` (`pull_errors` table) persists error entries with thread-safe execution and indexed lookups.
- [ ] `ParseFlexibleDBTimestamp` parses dates across `RFC3339Nano`, `RFC3339`, ISO8601, and SQLite formats without falling back to zero-time.
- [ ] Terminal failure summaries render nested box-drawing connectors (`├──`, `└──`, `│   `) with `Diagnostic: To inspect stack trace: gitmap pull-error <repo> (or: gitmap pe)`.
- [ ] Interactive dual hints for missing directories return Option 1 (`gitmap clone <repo>`) and Option 2 (`gitmap rm --db-only <repo>`).
- [ ] Remote node healing script `repo-secrets/05-scripts/heal-u1-pull-errors.py` creates database backup, normalizes backslashes to `/`, reconciles casing and suffixes, and purges `.oh-my-zsh` from registry.
- [ ] Remote SSH execution `gitmap ssh exec u1 "gitmap pa"` succeeds with 0 failures and 0 missing repositories.

