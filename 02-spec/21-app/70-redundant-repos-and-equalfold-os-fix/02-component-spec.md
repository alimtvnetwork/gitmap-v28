# Component Specification: EqualFold String Efficiency & SQLite Ingestion Uniqueness

## User Request (Verbatim)

```text
Do you have the code to improve the optimize the redundant repos that has been created? That's it. Do you have the code to do that? Confirm it. And also, I've seen in the code you are trying to compare the code in lowercase. Try not to do that. the specific method strings unfold, which is a lot more faster. Try to use that when you are comparing two strings without the case sensitivity. So try to apply that in terms of check-in, because in Unix, the different paths, different cases actually mean same thing. So you can only check ignoring the path in Windows. So you need to understand which OS you are in. So I think that is a bug we need to fix, and again, make a release and make sure you also optimize in the next pool if there is a redundancy. Okay? You can find the redundancy from your SQLite database, not from file system. Okay, so in future, when you add new files, you make sure that it is unique. Do you understand?

# Actionable Items Must Follow Non-Negotiable

1. Confirm if the code to optimize redundant repositories exists.
2. Avoid comparing code in lowercase; use the `strings unfold` method for case-insensitive comparisons.
3. Identify and fix the OS-specific bug related to path case sensitivity.
4. Ensure redundancy checks are performed using the SQLite database.
5. Guarantee uniqueness of new files added in the future.

Must follow and spawn agent using 

@[.agents/skills/execute-parent-task-with-n-steps-v6]

## Additional Instructions

learn /learn if you have to learn something and /plan stuff before working please.
```

---

## 1. Scope & Architectural Purpose

This component specification defines the technical design, contracts, and quality standards for two core pillars of Task 70:
1. **Zero-Allocation `strings.EqualFold` Refactoring (Task-02)**: Eliminating heap allocations and performance bottlenecks caused by naive lowercase string conversions (`strings.ToLower(a) == strings.ToLower(b)` or `strings.ToLower(s) == "literal"`) across 30+ call sites throughout the GitMap codebase.
2. **SQLite Ingestion Uniqueness Guarantees (Task-04)**: Establishing bulletproof uniqueness and deterministic upsert rules for files (`RepoFile` in per-repo Split DBs) and repositories (`Repo` in `gitmap.db`) to ensure that future file additions and repository registrations are strictly idempotent and free from duplicate entries.

---

## 2. Module 1: Zero-Allocation `strings.EqualFold` Architecture

### 2.1 The Problem with Naive Lowercase Equality

In polyglot and high-throughput CLI engines like GitMap, comparing string values without case sensitivity is frequent (flags, subcommands, branch names, user prompts, JSON type keys, etc.). 

A widespread anti-pattern observed across `cli/` is:
```go
// ❌ ANTI-PATTERN: Heap allocations + eager full-string copying
if strings.ToLower(a) == strings.ToLower(b) { ... }
if strings.ToLower(input) == "yes" { ... }
```

#### Why This Is Detrimental:
1. **Dual Heap Allocation**: `strings.ToLower(s)` creates a new byte buffer, converts each rune, and allocates a new `string` on the Go heap for each operand. For two 64-byte paths or URLs, this forces two separate heap allocations (128 bytes + GC tracking overhead) for a single comparison.
2. **Lack of Short-Circuiting**: Even if strings differ on the very first character (e.g. comparing `"Windows"` to `"Linux"`), `strings.ToLower` traverses the entire length of both strings, converting every character before equality is evaluated.
3. **Garbage Collection Pressure**: In tight loops (traversing thousands of repository files, fleet nodes, or command history items), accumulating short-lived lowercase string allocations triggers frequent garbage collection pauses.

### 2.2 The Solution: `strings.EqualFold`

Go's standard library provides `strings.EqualFold(s, t string) bool`, which implements Unicode case-folding comparison:
```go
// ✅ RECOMMENDED: Zero allocations + early exit short-circuiting
if strings.EqualFold(a, b) { ... }
if strings.EqualFold(input, "yes") { ... }
```

#### Key Technical Benefits:
- **Zero Allocations (0 B/op, 0 allocs/op)**: `strings.EqualFold` scans the underlying byte slices directly without copying or allocating new memory.
- **Short-Circuit Evaluation**: It aborts immediately upon encountering the first non-equivalent rune, providing \(O(1)\) performance in mismatched cases.
- **Full ASCII & Simple Unicode Fold Support**: Correctly handles ASCII letters (`A-Z` to `a-z`) as well as Unicode case mappings without external dependencies.

### 2.3 Refactoring Pattern Transformation Matrix

The following table establishes the mandatory transformation patterns to apply across the repository:

| Context Pattern | Legacy Anti-Pattern | Modernized Zero-Alloc Pattern |
|:---|:---|:---|
| **Two Dynamic Strings** | `strings.ToLower(a) == strings.ToLower(b)` | `strings.EqualFold(a, b)` |
| **Dynamic String vs Literal** | `strings.ToLower(s) == "yes"` | `strings.EqualFold(s, "yes")` |
| **Literal vs Dynamic String** | `"all" == strings.ToLower(token)` | `strings.EqualFold(token, "all")` |
| **Negative Equality** | `strings.ToLower(s) != "folder"` | `!strings.EqualFold(s, "folder")` |
| **Whitespace Trimmed** | `strings.TrimSpace(strings.ToLower(s)) == "y"` | `strings.EqualFold(strings.TrimSpace(s), "y")` |
| **Prefix with Trim** | `strings.ToLower(strings.TrimPrefix(s, ".")) == ext` | `strings.EqualFold(strings.TrimPrefix(s, "."), ext)` |
| **Pre-lowercased Variable** | `low := strings.ToLower(s); if low == "a" \|\| low == "b"` | `if strings.EqualFold(s, "a") \|\| strings.EqualFold(s, "b")` |
| **Switch Statement Cases** | `low := strings.ToLower(arg); switch { case low == "a": ... }` | `switch { case strings.EqualFold(arg, "a"): ... }` |

### 2.4 Nuances & Strict Exclusions (When NOT to replace)

Engineers must exercise care and NOT replace `strings.ToLower` in the following scenarios:
1. **Map Key Normalization**: When storing or retrieving entries from a hash map where keys are canonically lowercase (e.g. `seen[strings.ToLower(key)] = true`). `EqualFold` does not produce a normalized hash key.
2. **Substrings and Contains**: `strings.Contains(strings.ToLower(s), target)` cannot be replaced by `EqualFold`. For case-insensitive substring searching, an explicit case-folding search helper or regex is required.
3. **SQL Parameters for Case-Sensitive Columns**: When binding query parameters to SQLite columns where database collation requires lowercased inputs (e.g., `WHERE LOWER(Slug) = ?`).
4. **Canonical Path Identity on Unix**: When generating unique path identifiers, Unix filesystems are case-sensitive. Lowercasing must never be applied to Unix paths without OS detection.

### 2.5 Catalog of Identified Refactoring Sites

A deep scan across `cli/` has cataloged the following target sites for `strings.EqualFold` conversion:

| Module / Package | File Path | Line Range | Original Comparison Pattern |
|:---|:---|:---|:---|
| `cloner` | `cli/cloner/runners.go` | 249 | `strings.ToLower(strings.TrimSpace(response)) == "y"` |
| `cmd` | `cli/cmd/profiles_ops.go` | 214 | `strings.ToLower(p.Name) == lower \|\| strings.ToLower(p.ID) == lower` |
| `cmd` | `cli/cmd/reverttxn.go` | 232 | `strings.TrimSpace(strings.ToLower(scanner.Text())) == "yes"` |
| `cmd` | `cli/cmd/reverttxn_lastn.go` | 101 | `strings.TrimSpace(strings.ToLower(scanner.Text())) == "yes"` |
| `cmd` | `cli/cmd/visibilityresolve.go` | 237 | `strings.TrimSpace(strings.ToLower(line)) != "yes"` |
| `commitin` | `cli/cmd/commitin/config_json.go` | 395-398 | `strings.ToLower(item.Category) == lowerRef \|\| ...` |
| `commitin` | `cli/cmd/commitin/message/weak.go` | 39 | `strings.ToLower(strings.TrimSpace(ww)) == w` |
| `folder` | `cli/cmd/folder/filter.go` | 103 | `strings.TrimPrefix(strings.ToLower(a), ".") == cleanExt` |
| `cmdagy` | `cli/cmdagy/agy_list_prompts.go` | 67 | `strings.ToLower(cleanArgs[0]) == "all"` |
| `cmdagy` | `cli/cmdagy/agy_pin_projects_lookup.go` | 52 | `strings.ToLower(p.Name) == low ...` |
| `cmdagy` | `cli/cmdagy/agy_pin_projects_resolve.go` | 57 | `strings.ToLower(p.Name) == low ...` |
| `cmdagy` | `cli/cmdagy/agy_rm_ops.go` | 40 | `strings.ToLower(args[0]) == "folder" \|\| strings.ToLower(args[0]) == "dir"` |
| `cmdagy` | `cli/cmdagy/agy_target_match.go` | 10 | `strings.ToLower(p.ID) == lower ...` |
| `cmdagy` | `cli/cmdagy/agy_wpr_cache.go` | 125 | `strings.ToLower(entry.Alias) == clean \|\| strings.ToLower(entry.MachineName) == clean` |
| `cmdai` | `cli/cmdai/ai_list.go` | 58 | `strings.ToLower(string(s.Category)) == clean` |
| `cmdchromeprofile`| `cli/cmdchromeprofile/chromeprofile_resolve.go` | 143 | `strings.ToLower(strings.TrimSpace(info.Name)) == want` |
| `cmdclone` | `cli/cmdclone/clonenext.go` | 389 | `strings.ToLower(strings.TrimSpace(answer)) == "y"` |
| `cmdclone` | `cli/cmdclone/clonenext.go` | 456 | `strings.ToLower(strings.TrimSpace(answer)) != "y"` |
| `cmdignore` | `cli/cmdignore/ignore_cmd.go` | 22 | `strings.ToLower(args[0]) == "set"` |
| `cmdignore` | `cli/cmdignore/ignore_cmd.go` | 45 | `strings.ToLower(args[0]) == "interval"` |
| `cmdignore` | `cli/cmdignore/ignore_cmd.go` | 220 | `strings.ToLower(a) == "--json"` |
| `cmdinstall` | `cli/cmdinstall/install_export.go` | 115 | `strings.ToLower(opts.Format) == "json"` |
| `cmdinstall` | `cli/cmdinstall/install_export.go` | 137 | `strings.ToLower(opts.Format) == "json"` |
| `cmdmacro` | `cli/cmdmacro/macro_deploy_ssh.go` | 106 | `strings.ToLower(m.Name) == targetLower` |
| `cmdos` | `cli/cmdos/os_ai_clean.go` | 81 | `strings.ToLower(args[0]) == target` |
| `cmdos` | `cli/cmdos/os_dev_clean.go` | 118 | `strings.ToLower(strings.TrimSpace(text)) == "yes"` |
| `cmdos` | `cli/cmdos/terminal_clean_cmd.go` | 108 | `strings.ToLower(strings.TrimSpace(text)) == "yes"` |
| `cmdssh` | `cli/cmdssh/sshdelete.go` | 48 | `strings.TrimSpace(strings.ToLower(input)) != "y"` |
| `cmdssh` | `cli/cmdssh/sshgen.go` | 273 | `strings.TrimSpace(strings.ToLower(input)) == "y"` |
| `cmdupdate` | `cli/cmdupdate/update_fleet.go` | 912 | `strings.ToLower(pkg) == "all"` |
| `cmdupdate` | `cli/cmdupdate/update_fleet.go` | 1019 | `low == "agm" \|\| low == "agy" ...` |
| `cmdvscode` | `cli/cmdvscode/vscode_cmd.go` | 47 | `sub == "ls" \|\| sub == "list" ...` |
| `cmdwinutil` | `cli/cmdwinutil/winutil_cmd.go` | 71-77 | `low == "--dry-run" \|\| low == "-n" ...` |
| `committransfer` | `cli/committransfer/prclean/pr_clean.go` | 82 | `resp == "y" \|\| resp == "yes"` |
| `gitignoreagm` | `cli/gitignoreagm/cli_prompt.go` | 122 | `ans == "" \|\| ans == "y" \|\| ans == "yes"` |
| `jsonenvelope` | `cli/jsonenvelope/registry.go` | 123 | `strings.ToLower(d.Type) == norm` |
| `macro` | `cli/macro/open.go` | 28 | `strings.ToLower(fields[0]) != "open"` |
| `macro` | `cli/macro/open.go` | 94 | `lower == "chrome" \|\| lower == "google-chrome" ...` |
| `macro` | `cli/macro/record.go` | 79-87 | `lower == "stop" \|\| lower == "exit" ...` |
| `macro` | `cli/macro/record_undo.go` | 80 | `ans == "y" \|\| ans == "yes"` |
| `osuser` | `cli/osuser/kill.go` | 57 | `strings.ToLower(curr.Username) == target` |
| `release` | `cli/release/autocommit.go` | 176 | `answer == "y" \|\| answer == "yes"` |
| `release` | `cli/release/workflowvalidate.go` | 129 | `answer != "y" && answer != "yes"` |

---

## 3. Module 2: SQLite Ingestion Uniqueness Guarantees

### 3.1 Overview of GitMap SQLite Architecture

GitMap uses a multi-tier SQLite architecture:
1. **Root Database (`gitmap.db`)**: Tracks overall workspace state, settings, and discovered git repositories (`Repo` table).
2. **Per-Repository Split Database (`repo_search/<slug>-<repoId>.db`)**: Stores indexed file manifests, contents, search caches, and sequence numbers (`RepoFile`, `FileSequence`, `SearchCache` tables).

When new repositories or files are ingested into SQLite, the database engine MUST guarantee uniqueness and enforce deterministic upsert semantics without unhandled constraint violations or duplicate records.

```
+-------------------------------------------------------------------------------+
|                             Ingestion Pipeline                                |
+-------------------------------------------------------------------------------+
                                      |
         +----------------------------+----------------------------+
         |                                                         |
         v                                                         v
+-------------------------------+         +-------------------------------------+
|   Repo Ingestion (gitmap.db)   |         | RepoFile Ingestion (Split DB)       |
+-------------------------------+         +-------------------------------------+
| Table: Repo                   |         | Table: RepoFile                     |
| Primary Key: RepoId           |         | Primary Key: RepoFileId             |
| Unique Index: AbsolutePath    |         | Unique Constraint: RelativePath     |
|   (COLLATE NOCASE)            |         | Index: IdxRepoFile_RelativePath     |
| Strategy: UPSERT              |         | Strategy: UPSERT                    |
| ON CONFLICT(AbsolutePath)     |         | ON CONFLICT(RelativePath)           |
| DO UPDATE SET ...             |         | DO UPDATE SET ...                   |
+-------------------------------+         +-------------------------------------+
```

---

### 3.2 Guarantee 1: File Ingestion Uniqueness (`RepoFile`)

#### 3.2.1 Schema Definition & Unique Index
In each repository's split database (`cli/repodb/repo_db.go`), `RepoFile` represents tracked and indexed source files.
The schema enforces uniqueness on `RelativePath`:

```sql
CREATE TABLE IF NOT EXISTS RepoFile (
    RepoFileId   INTEGER PRIMARY KEY AUTOINCREMENT,
    RelativePath TEXT NOT NULL UNIQUE,
    AbsolutePath TEXT NOT NULL,
    Content      TEXT,
    IsBig        INTEGER NOT NULL,
    WriteTime    INTEGER NOT NULL,
    CreatedAt    INTEGER NOT NULL,
    UpdatedAt    INTEGER NOT NULL
);

-- Explicit uniqueness index for high-performance lookup and conflict resolution
CREATE UNIQUE INDEX IF NOT EXISTS IdxRepoFile_RelativePath ON RepoFile(RelativePath);
```

#### 3.2.2 Conflict Resolution & Deterministic Upsert Rules
Blind `INSERT INTO RepoFile` without conflict resolution causes `UNIQUE constraint failed: RepoFile.RelativePath` whenever a file is re-scanned or modified.

All insertions into `RepoFile` MUST use SQLite 3.24+ `ON CONFLICT(RelativePath) DO UPDATE`:
```sql
INSERT INTO RepoFile (
    RelativePath, AbsolutePath, Content, IsBig, WriteTime, CreatedAt, UpdatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(RelativePath) DO UPDATE SET
    AbsolutePath = excluded.AbsolutePath,
    Content      = excluded.Content,
    IsBig        = excluded.IsBig,
    WriteTime    = excluded.WriteTime,
    UpdatedAt    = excluded.UpdatedAt;
```

#### 3.2.3 Typed Repository Contract (`cli/repodb/repo_file.go`)
Currently, `RepoFileDbRepo` only exposes `Insert(...)` which executes a plain `INSERT`.
To guarantee uniqueness across all ingestion code, `RepoFileDbRepo` must provide a typed `Upsert` method:

```go
// Upsert inserts a new RepoFile record or updates existing content on RelativePath conflict.
func (r *RepoFileDbRepo) Upsert(ctx context.Context, item *RepoFile) dbengine.RowsAffectedResult {
    query := `INSERT INTO RepoFile (
        RelativePath, AbsolutePath, Content, IsBig, WriteTime, CreatedAt, UpdatedAt
    ) VALUES (?, ?, ?, ?, ?, ?, ?)
    ON CONFLICT(RelativePath) DO UPDATE SET
        AbsolutePath = excluded.AbsolutePath,
        Content      = excluded.Content,
        IsBig        = excluded.IsBig,
        WriteTime    = excluded.WriteTime,
        UpdatedAt    = excluded.UpdatedAt;`

    now := time.Now().Unix()
    createdAt := item.CreatedAt
    if createdAt == 0 {
        createdAt = now
    }
    updatedAt := now

    return r.db.ExecRowsAffected(
        ctx, query,
        item.RelativePath, item.AbsolutePath, item.Content,
        item.IsBig, item.WriteTime, createdAt, updatedAt,
    )
}
```

#### 3.2.4 Walker Idempotency (`cli/indexer/walker.go`)
The background filesystem indexer (`Walker.Walk`) streams discovered files into a worker pool.
- `inspectAndQueueFile`: Checks `cachedTimes` loaded from `SELECT RelativePath, WriteTime FROM RepoFile`.
- `upsertRepoFile`: Executes `sqlUpsertRepoFile`.
- **Guarantee**: Even if two concurrent worker routines process the same relative path or sequential scans re-index existing files, no duplicate row can ever be created. The row is updated in place, updating `WriteTime`, `Content`, and `UpdatedAt`.

---

### 3.3 Guarantee 2: Repository Ingestion Uniqueness (`Repo`)

#### 3.3.1 Schema & Collation Constraints
The root database `gitmap.db` stores all discovered and cloned repositories in table `Repo` (`cli/constants/constants_store.go`).

```sql
CREATE TABLE IF NOT EXISTS Repo (
    RepoId                INTEGER PRIMARY KEY AUTOINCREMENT,
    Slug                  TEXT NOT NULL,
    RepoName              TEXT NOT NULL,
    HttpsUrl              TEXT NOT NULL,
    SshUrl                TEXT NOT NULL,
    Branch                TEXT NOT NULL,
    RelativePath          TEXT NOT NULL,
    AbsolutePath          TEXT NOT NULL,
    CloneInstruction      TEXT NOT NULL,
    Notes                 TEXT DEFAULT '',
    LastInjectedDesktopAt TEXT DEFAULT '',
    LastInjectedVSCodeAt  TEXT DEFAULT '',
    LastClonedAt          TEXT DEFAULT '',
    IdentifiedTransport   TEXT NOT NULL DEFAULT '',
    CreatedAt             TEXT DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt             TEXT DEFAULT CURRENT_TIMESTAMP
);

-- Case-insensitive unique index on AbsolutePath prevents mixed-case duplicates on Windows
CREATE UNIQUE INDEX IF NOT EXISTS IdxRepo_AbsolutePath ON Repo(AbsolutePath COLLATE NOCASE);
```

#### 3.3.2 Pre-Ingestion Storage Path Normalization
Before any repository path enters SQL statements, it MUST be passed through `store.NormalizeStoragePath(pathStr)`:
- Cleans relative navigation components (`.`, `..`).
- Normalizes path separators to OS standard.
- On Windows, capitalizes drive letters (e.g. `./repo` -> `./repo`) to eliminate case variations.

#### 3.3.3 Upsert Strategy (`constants.SQLUpsertRepoByPath`)
Repository additions are executed via `db.UpsertRepos(records)`:
```sql
INSERT INTO Repo (
    Slug, RepoName, HttpsUrl, SshUrl, Branch,
    RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(AbsolutePath) DO UPDATE SET
    Slug                = excluded.Slug,
    RepoName            = excluded.RepoName,
    HttpsUrl            = excluded.HttpsUrl,
    SshUrl              = excluded.SshUrl,
    Branch              = excluded.Branch,
    RelativePath        = excluded.RelativePath,
    CloneInstruction    = excluded.CloneInstruction,
    Notes               = excluded.Notes,
    IdentifiedTransport = excluded.IdentifiedTransport,
    UpdatedAt           = CURRENT_TIMESTAMP;
```

#### 3.3.4 Remote URL Redundancy Guard
While `AbsolutePath` is unique at the filesystem level, repositories with the same remote URL (e.g. cloned into multiple folders) represent logical redundancies.
- During future ingestion (scans or clones), GitMap checks `gitmap.db` using:
  ```sql
  SELECT RepoId, AbsolutePath FROM Repo 
  WHERE HttpsUrl = ? OR SshUrl = ? OR Slug = ?;
  ```
- If a match exists, GitMap logs the association rather than duplicating database records blindly, adhering to the SQLite-first redundancy optimization rule.

---

## 4. Verification and Acceptance Criteria

### 4.1 Zero-Alloc `strings.EqualFold` Acceptance Criteria
1. **Zero Naive Lowercase Comparisons**: Static analysis confirms that no `strings.ToLower(a) == strings.ToLower(b)` or `strings.ToLower(s) == "literal"` remains in user prompt/flag/command matching routines across `cli/cmd`, `cli/cmddb`, `cli/cmdclone`, `cli/cmdssh`, `cli/release`, `cli/osclean`.
2. **Zero Additional Heap Allocations**: Benchmark tests comparing `BenchmarkEqualFoldVsToLower` confirm 0 B/op and 0 allocs/op for `strings.EqualFold`.
3. **Logic Invariance**: All existing unit tests pass without behavioral regressions.

### 4.2 SQLite Uniqueness Acceptance Criteria
1. **Duplicate File Ingestion Test**: Ingesting the same relative path 100 times into `RepoFile` via `Upsert` or `Walker` produces exactly 1 row in the SQLite database, with updated `UpdatedAt` timestamp.
2. **Duplicate Repo Ingestion Test**: Ingesting the same absolute path with different casing (e.g. `D:\Work\repo` vs `./repo`) into `Repo` resolves to a single row due to `IdxRepo_AbsolutePath (COLLATE NOCASE)`.
3. **Foreign Key and Cascading Integrity**: Updating or deduplicating records leaves no orphaned entries in `GroupRepo` or `Release`.
