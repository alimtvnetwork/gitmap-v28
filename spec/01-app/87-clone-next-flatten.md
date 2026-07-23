# 87 — Clone-Next Flatten Mode (Default Behavior)

## Overview

As of v2.75.0, `gitmap-v28 clone-next` (`cn`) **flattens by default** — it
clones a versioned repository into a **consistent base-name folder** instead
of the version-suffixed folder. No `--flatten` flag is required.

This keeps a single, predictable local path across version iterations
(e.g., always `macro-ahk/` instead of `macro-ahk-v15/`, `macro-ahk-v16/`).

The behavior also enables **version tracking** in the gitmap-v28 database:
both the current active version and a full transition history.

---

## Command Syntax

```
gitmap-v28 cn <version-spec> [-f|--force] [--delete] [--keep] [--no-desktop] [--verbose]
```

| Flag | Description |
|------|-------------|
| `-f`, `--force` | Force flatten even when cwd IS the target folder. chdir to parent and remove cwd before cloning. Refuses the versioned-folder fallback. |
| `--delete` | Remove the current versioned folder after clone (when different from flattened path) |
| `--keep` | Keep current folder without prompting |
| `--no-desktop` | Skip GitHub Desktop registration |
| `--verbose` | Print detailed progress |

### Flag Interactions

| Flags Used | Behavior |
|------------|----------|
| (none) | Clone into `macro-ahk/`, replacing it if it exists. If cwd IS `macro-ahk/` (already flattened), Windows file lock prevents removal → falls back to `macro-ahk-vN/` with a warning. |
| `-f` / `--force` | Same as default, but if cwd IS the target folder, gitmap-v28 chdirs to parent first, removes the cwd, and clones into the flattened name. **Never falls back to a versioned folder.** Aborts with a clear error if removal still fails. |
| `--delete` | Clone into `macro-ahk/`, then delete the old versioned folder (if a different path) |

### `-f` Use Case (v3.50.0+)

Working continuously from one flattened folder across version bumps:

```
PS C:\repos\macro-ahk> gitmap-v28 cn v++ -f
  → Force-flatten: leaving D:\repos\macro-ahk to release lock...
  Removing existing macro-ahk for fresh clone...
  Cloning macro-ahk-v22 into macro-ahk (flattened)...
  ✓ Cloned macro-ahk-v22 into macro-ahk
  → Now in macro-ahk
```

Without `-f`, this same flow falls back to `macro-ahk-v22/` because the
shell holds an open handle on `macro-ahk/`. `-f` is the explicit
contract that the user accepts losing their cwd in exchange for a
guaranteed-flat layout.



---

## Folder Name Resolution

### Base Name Extraction

Strip the version suffix from the repo name (derived from remote URL):

```
macro-ahk-v15     → macro-ahk
my-tool-v2        → my-tool
project-v100      → project
some-repo         → some-repo  (no version suffix — unchanged)
```

**Rules:**

1. Match the pattern `-v<digits>` at the end of the repo name.
2. Strip the matched suffix to produce the base name.
3. If no version suffix is found, use the repo name as-is.
4. The base name must be non-empty after stripping.

### Regex Pattern

```
^(.+)-v(\d+)$
```

- Group 1: base name (e.g., `macro-ahk`)
- Group 2: version number as string (e.g., `15`)

### Target Folder

The target clone folder is the base name only, at the same parent
directory level as the current folder:

```
# Current: /projects/macro-ahk-v15/
# Target:  /projects/macro-ahk/
```

---

## Version Number Parsing

Two representations are stored for every version:

| Field | Type | Example | Source |
|-------|------|---------|--------|
| Version tag | `TEXT` | `v16` | Extracted from target repo name |
| Version number | `INTEGER` | `16` | Parsed integer from the tag |

### Parsing Rules

1. Extract digits from the version suffix: `v16` → `16`.
2. Parse as integer. If parsing fails, store `0` and log a warning.
3. The tag always retains the `v` prefix for display consistency.

---

## Database Schema Changes

### Repos Table — New Columns

Add two columns to the existing `Repos` table via idempotent migration:

```sql
ALTER TABLE Repos ADD COLUMN CurrentVersionTag TEXT DEFAULT '';
ALTER TABLE Repos ADD COLUMN CurrentVersionNum INTEGER DEFAULT 0;
```

| Column | Type | Default | Description |
|--------|------|---------|-------------|
| `CurrentVersionTag` | `TEXT` | `''` | Full version string, e.g., `v16` |
| `CurrentVersionNum` | `INTEGER` | `0` | Integer version number, e.g., `16` |

These columns are updated on every clone-next operation.

### New Table: `RepoVersionHistory`

```sql
CREATE TABLE IF NOT EXISTS RepoVersionHistory (
    Id              INTEGER PRIMARY KEY AUTOINCREMENT,
    RepoId          INTEGER NOT NULL REFERENCES Repos(Id) ON DELETE CASCADE,
    FromVersionTag  TEXT NOT NULL,
    FromVersionNum  INTEGER NOT NULL,
    ToVersionTag    TEXT NOT NULL,
    ToVersionNum    INTEGER NOT NULL,
    FlattenedPath   TEXT DEFAULT '',
    CreatedAt       TEXT DEFAULT CURRENT_TIMESTAMP
);
```

| Column | Type | Description |
|--------|------|-------------|
| `Id` | `INTEGER PK` | Auto-incrementing primary key |
| `RepoId` | `INTEGER FK` | References `Repos(Id)`, cascade delete |
| `FromVersionTag` | `TEXT` | Previous version tag (e.g., `v15`) |
| `FromVersionNum` | `INTEGER` | Previous version number (e.g., `15`) |
| `ToVersionTag` | `TEXT` | New version tag (e.g., `v16`) |
| `ToVersionNum` | `INTEGER` | New version number (e.g., `16`) |
| `FlattenedPath` | `TEXT` | Base name of the flattened folder |
| `CreatedAt` | `TEXT` | ISO 8601 timestamp of the transition |

### Migration

The migration uses the idempotent `ALTER TABLE ADD COLUMN` pattern
with `isDuplicateColumnError` to silently skip if columns already exist.
The `CREATE TABLE IF NOT EXISTS` handles the new table idempotently.

---

## Workflow

### Step-by-Step Execution

```
1. Parse flags (--delete, --keep, --no-desktop, --verbose)
2. Resolve target version (existing clone-next logic)
3. Extract base name from remote repo name (strip -vN suffix)
4. Compute target path = parent_dir / base_name
5. IF target path exists:
   a. Remove target folder entirely (no prompt)
6. git clone <target-url> <base-name-folder>
7. Update Repos row: CurrentVersionTag, CurrentVersionNum
8. INSERT into RepoVersionHistory (from -> to)
9. Register with GitHub Desktop (using flattened path)
10. IF --delete AND source folder != target folder:
    a. Remove the old versioned folder
11. Shell handoff: write flattened path to `$GITMAP_HANDOFF_FILE`
    (set by the shell wrapper function) so the parent shell cds to it.
```

---

## Shell Handoff

The flattened path is written to the sentinel file pointed to by
`GITMAP_HANDOFF_FILE` (exported by the `gitmap-v28` shell wrapper function
before invocation). After the binary exits, the wrapper reads that file
and `cd`s the parent shell to the flattened folder.

Without the wrapper installed, `GITMAP_HANDOFF_FILE` is unset and the
write becomes a silent no-op — the user's cwd is unchanged.

> History: prior to v3.103.0 the spec referenced
> `GITMAP_SHELL_HANDOFF` as an env var set by `os.Setenv`. That was a
> no-op (child cannot mutate parent env). v3.103.0 replaced it with
> the sentinel-file mechanism. See
> [.lovable/memory/features/shell-handoff-file.md](../../.lovable/memory/features/shell-handoff-file.md).

See the [navigation helper](31-cd.md) spec for the shell wrapper
mechanism.

---

## Error Handling

All errors follow the project's zero-swallow policy. Every failure must
be logged to `os.Stderr` using the standardized format.

| Scenario | Behavior |
|----------|----------|
| Target folder removal fails | Exit with error, do not clone |
| Git clone fails | Exit with error, do not update DB |
| DB update fails | Log error to stderr, do not exit (clone succeeded) |
| Version number parse fails | Log warning, store `0` as version number |
| No version suffix on current repo | Use repo name as-is, version = 1 |

---

## Examples

### Basic Clone-Next (Flattened)

```bash
# In /projects/macro-ahk-v15/
gitmap-v28 cn v+1

# Result:
#   Cloned macro-ahk-v16 into macro-ahk/
#   DB: CurrentVersionTag="v16", CurrentVersionNum=16
#   History: v15 -> v16
#   Shell navigates to /projects/macro-ahk/
```

### Clone-Next with Delete

```bash
# In /projects/macro-ahk-v15/
gitmap-v28 cn v+1 --delete

# Result:
#   Cloned macro-ahk-v16 into macro-ahk/
#   Deleted /projects/macro-ahk-v15/
#   DB and history updated as above
```

### Repeated Clone-Next (v16 -> v17)

```bash
# In /projects/macro-ahk/ (flattened from v16)
gitmap-v28 cn v+1

# Result:
#   Removes /projects/macro-ahk/
#   Clones macro-ahk-v17 into macro-ahk/
#   DB: CurrentVersionTag="v17", CurrentVersionNum=17
#   History: v16 -> v17
```

### No Version Suffix

```bash
# In /projects/some-tool/
gitmap-v28 cn v+1

# Result:
#   Removes /projects/some-tool/
#   Clones some-tool-v2 into some-tool/
#   DB: CurrentVersionTag="v2", CurrentVersionNum=2
#   History: v1 -> v2
```

---

## Viewing Version History

Use `gitmap-v28 version-history` (`vh`) to see all recorded transitions:

```bash
gitmap-v28 vh

# Output:
# Version history for D:\wp-work\riseup-asia\macro-ahk:
#
# FROM        TO          FOLDER                    TIMESTAMP
# v15         v16         macro-ahk                 2026-04-16T10:30:00Z
# v16         v17         macro-ahk                 2026-04-16T14:22:00Z
#
# 2 transition(s) recorded.
```

---

## Acceptance Criteria

1. **Flatten by default**: `clone-next` always clones into the base-name folder.
2. **Folder replacement**: If the base-name folder exists, it is removed
   before cloning (no prompt).
3. **Version suffix parsing**: Correctly strips `-v<digits>` suffix from
   repo names; handles missing suffix gracefully.
4. **DB current version**: `Repos.CurrentVersionTag` and
   `Repos.CurrentVersionNum` are updated on every clone-next.
5. **DB history**: A new `RepoVersionHistory` row is inserted for each
   transition, with correct from/to version data and flattened path.
6. **Shell handoff**: `GITMAP_SHELL_HANDOFF` is set to the flattened path.
7. **GitHub Desktop**: The flattened path is registered correctly.
8. **Error handling**: All failures logged to stderr; clone failure
   prevents DB update; DB failure does not roll back a successful clone.
9. **Idempotent migration**: Schema changes use `isDuplicateColumnError`
   and `CREATE TABLE IF NOT EXISTS`.
10. **Constants**: All messages and error strings are defined in the
    constants package — no magic strings in command logic.

---

## Component Mapping

| Component | File / Package | Responsibility |
|-----------|---------------|----------------|
| Command handler | `cmd/clonenext.go` | Flatten logic, clone, folder removal |
| Version history | `cmd/clonenexthistory.go` | Record transitions in DB |
| Flag parser | `cmd/clonenextflags.go` | Parse clone-next flags |
| Version parser | `clonenext/version.go` | Strip `-vN` suffix, resolve target |
| DB migration | `store/store.go` | Add columns to `Repos`, create `RepoVersionHistory` |
| DB write (current) | `store/version_history.go` | Update `CurrentVersionTag` and `CurrentVersionNum` |
| DB write (history) | `store/version_history.go` | Insert transition row |
| Version history CLI | `cmd/versionhistory.go` | Display version transitions |
| Constants (messages) | `constants/constants_version_history.go` | SQL, messages, error strings |
| Shell handoff | `cmd/clonenext.go` | Set `GITMAP_SHELL_HANDOFF` to flattened path |
| GitHub Desktop | `cmd/clonenext.go` | Register flattened path |

---

## Cross-References

| Document | Relevance |
|----------|-----------|
| [59-clone-next.md](59-clone-next.md) | Full clone-next command spec |
| [31-cd.md](31-cd.md) | Shell handoff via `GITMAP_SHELL_HANDOFF` |
| [13-release-data-model.md](./13-release-data-model.md) | Migration pattern reference |

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
