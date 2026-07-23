# Temp Release Command

## Overview

The `temp-release` command (`tr`) creates lightweight, temporary release branches from recent commits without creating tags. It enables rapid experimentation — spin up multiple candidate releases, review them, then clean up the ones you don't need.

**Key Difference from `release`:** No tags are created. No metadata (`.gitmap/release/` JSON) is written. No auto-commit step. Branches use the `temp-release/` prefix instead of `release/`.

## Commands

| Command | Alias | Description |
|---------|-------|-------------|
| `temp-release` | `tr` | Create temp-release branches from recent commits |
| `temp-release list` | `tr list` | List all temp-release branches |
| `temp-release remove` | `tr remove` | Remove temp-release branches |

## Branch Creation

### Syntax

```
gitmap-v28 temp-release <count> <version-pattern> [-s <start>]
gitmap-v28 tr <count> <version-pattern> [-s <start>]
```

### Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `count` | Yes | Number of recent commits to create branches for (1–50) |
| `version-pattern` | Yes | Version string with `$` placeholders for sequence number |
| `-s`, `--start` | No | Starting sequence number (default: auto-increment from last temp-release) |

### Version Pattern

The `$` placeholder determines zero-padded digit width:
- `$$` → 2 digits (`05`, `12`, `99`)
- `$$$` → 3 digits (`005`, `012`, `099`)
- `$$$$` → 4 digits (`0005`, `0012`, `0099`)

The pattern **must** contain at least one `$` sequence.

### Examples

```bash
# Create 10 branches from last 10 commits, starting at sequence 5
gitmap-v28 tr 10 v1.$$ -s 5
# Creates: temp-release/v1.05 (oldest) through temp-release/v1.14 (newest)

# Create 1 branch, auto-increment from last temp-release
gitmap-v28 tr 1 v1.$$
# If last temp-release was v1.43, creates: temp-release/v1.44

# Create 5 branches with 3-digit padding
gitmap-v28 tr 5 v2.1.$$$ -s 1
# Creates: temp-release/v2.1.001 through temp-release/v2.1.005
```

### Branch Naming

Branches use the prefix `temp-release/`:

```
temp-release/v1.05
temp-release/v1.06
...
```

### Commit Ordering

Commits are ordered from oldest to newest. The **oldest** commit in the range gets the **lowest** sequence number:

```
Commit history (git log order, newest first):
  abc1234  ← most recent  → highest sequence number
  def5678
  ...
  xyz9999  ← 10th back    → lowest sequence number (start)
```

### No Branch Switching

Branches are created using `git branch <name> <sha>` — this does NOT switch the current working branch. All branches are created from their respective commit SHAs without checkout.

### Push to Remote

After all branches are created locally, they are pushed to `origin` in a single batch:

```
git push origin temp-release/v1.05 temp-release/v1.06 ... temp-release/v1.14
```

## Auto-Increment

When `-s` is not provided, the system determines the next sequence number:

1. Query the `TempReleases` table for the last used sequence number matching the version prefix (e.g., `v1.` for pattern `v1.$$`).
2. If no DB record exists, scan remote branches matching `temp-release/v1.*` and parse the highest number.
3. If nothing found, start at 1.

The resolved start number is printed before branch creation:

```
  → Starting sequence: 44 (auto-detected from last temp-release)
```

## Listing Temp Releases

### Syntax

```
gitmap-v28 tr list [--json]
```

### Output (Terminal)

```
  Temp-release branches (6):

  temp-release/v1.05  abc1234  Fix login bug           2025-03-20
  temp-release/v1.06  def5678  Add dashboard widget    2025-03-20
  temp-release/v1.07  789abcd  Update dependencies     2025-03-21
  ...
```

Shows: branch name, short SHA, commit message (truncated), and commit date.

### Output (JSON)

```json
[
  {
    "branch": "temp-release/v1.05",
    "commit": "abc1234def5678...",
    "message": "Fix login bug",
    "date": "2025-03-20T14:30:00Z"
  }
]
```

## Removing Temp Releases

### Syntax

```bash
# Remove a single branch
gitmap-v28 tr remove v1.05

# Remove a range (inclusive)
gitmap-v28 tr remove v1.05 to v1.10

# Remove all temp-release branches
gitmap-v28 tr remove all
```

### Confirmation

All remove operations require interactive confirmation:

**Single:**
```
  Remove temp-release/v1.05? (y/N): y
  ✓ Removed temp-release/v1.05 (local + remote)
```

**Range:**
```
  Remove 6 temp-release branches:
    temp-release/v1.05
    temp-release/v1.06
    temp-release/v1.07
    temp-release/v1.08
    temp-release/v1.09
    temp-release/v1.10
  Proceed? (y/N): y
  ✓ Removed 6 temp-release branches (local + remote)
```

**All:**
```
  Remove ALL 12 temp-release branches:
    temp-release/v1.05
    temp-release/v1.06
    ...
    temp-release/v1.16
  Proceed? (y/N): y
  ✓ Removed 12 temp-release branches (local + remote)
```

### Remote Cleanup

Branches are deleted both locally and from the remote:

```bash
git branch -D temp-release/v1.05
git push origin --delete temp-release/v1.05
```

For batch operations, remote deletions use a single push:

```bash
git push origin --delete temp-release/v1.05 temp-release/v1.06 ...
```

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-s`, `--start` | int | auto | Starting sequence number |
| `--dry-run` | bool | false | Preview branch names without creating |
| `--json` | bool | false | JSON output for `list` subcommand |
| `--verbose` | bool | false | Detailed logging |

## Database Schema

### Table: `TempReleases`

Tracks temp-release branch metadata for auto-increment and listing:

```sql
CREATE TABLE IF NOT EXISTS TempReleases (
    Id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    Branch TEXT NOT NULL UNIQUE,
    VersionPrefix TEXT NOT NULL DEFAULT '',
    SequenceNumber INTEGER NOT NULL DEFAULT 0,
    Commit TEXT NOT NULL DEFAULT '',
    CommitMessage TEXT NOT NULL DEFAULT '',
    CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

- `Branch`: full branch name, e.g. `temp-release/v1.05`
- `VersionPrefix`: the prefix before the `$$`, e.g. `v1.`
- `SequenceNumber`: the numeric value substituted for `$$`, e.g. `5`
- `Commit`: full commit SHA
- `CommitMessage`: first line of commit message

Records are inserted on creation and deleted on removal.

## Error Handling

| Scenario | Behavior |
|----------|----------|
| `count` < 1 or > 50 | Print error, exit 1 |
| Pattern has no `$$` placeholder | Print error, exit 1 |
| Fewer commits than `count` | Print warning, create for available commits only |
| Branch already exists | Skip with warning, continue with remaining |
| Remote push fails | Print error for failed branches, continue with others |
| Range `remove` with missing branches | Skip missing with warning, remove existing ones |
| Sequence overflow (e.g., 100 for `$$`) | Print error: "sequence 100 exceeds 2-digit format" |

## Dry Run

```bash
gitmap-v28 tr 10 v1.$$ -s 5 --dry-run
```

```
  Dry-run: would create 10 temp-release branches:
    temp-release/v1.05  abc1234  Fix login bug
    temp-release/v1.06  def5678  Add dashboard widget
    ...
    temp-release/v1.14  999aaab  Bump version
```

## Implementation Notes

### Package Structure

| File | Responsibility |
|------|----------------|
| `cmd/temprelease.go` | CLI entry, flag parsing, dispatch to subcommands |
| `cmd/tempreleaseops.go` | Create, list, remove logic |
| `release/temprelease.go` | Git operations (branch create/delete, commit listing) |
| `store/temprelease.go` | DB CRUD for `TempReleases` table |
| `constants/constants_temprelease.go` | All string constants |

### Constants

```go
// constants/constants_temprelease.go
const (
    CmdTempRelease      = "temp-release"
    CmdTempReleaseShort = "tr"

    TempReleaseBranchPrefix = "temp-release/"
    TempReleasePlaceholder  = "$$"

    TempReleaseMaxCount     = 50
)
```

### Git Operations (No Checkout)

Branch creation without switching:
```go
// git branch temp-release/v1.05 abc1234
exec.Command("git", "branch", branchName, commitSHA)
```

Batch push:
```go
// git push origin branch1 branch2 branch3 ...
args := append([]string{"push", "origin"}, branchNames...)
exec.Command("git", args...)
```

Batch delete (remote):
```go
// git push origin --delete branch1 branch2 ...
args := append([]string{"push", "origin", "--delete"}, branchNames...)
exec.Command("git", args...)
```

Local delete:
```go
// git branch -D branch1
exec.Command("git", "branch", "-D", branchName)
```

### Dispatch Registration

Add to `cmd/root.go` in `dispatchMisc`:

```go
if command == constants.CmdTempRelease || command == constants.CmdTempReleaseShort {
    runTempRelease(os.Args[2:])
    return true
}
```

## CLI Help Entry

```
  temp-release (tr) <count> <pattern> [-s N]  Create temp branches from recent commits
```

## Acceptance Criteria

1. `gitmap-v28 tr 10 v1.$$ -s 5` creates 10 branches (`temp-release/v1.05`–`temp-release/v1.14`) from the last 10 commits, pushed to origin, without switching branches.
2. `gitmap-v28 tr 1 v1.$$` auto-detects the next sequence number from DB/remote and creates one branch.
3. `gitmap-v28 tr list` shows all temp-release branches with SHA, message, and date.
4. `gitmap-v28 tr list --json` outputs structured JSON.
5. `gitmap-v28 tr remove v1.05` prompts confirmation, then deletes local + remote branch.
6. `gitmap-v28 tr remove v1.05 to v1.10` prompts with full list, then deletes all in range.
7. `gitmap-v28 tr remove all` prompts with all branch names, then deletes everything.
8. `gitmap-v28 tr 10 v1.$$ --dry-run` shows preview without creating branches.
9. `$$$` produces 3-digit zero-padded numbers; `$$$$` produces 4-digit.
10. Sequence overflow (e.g., 100 for `$$`) produces a clear error.
11. Existing branches are skipped with a warning during creation.
12. DB records are inserted on create and deleted on remove.
13. No tags are created. No `.gitmap/release/` metadata is written.

## Related Commands

- [release](./12-release-command.md) — full release with tags and metadata
- [release-branch](./12-release-command.md) — complete release from existing branch
- [release-pending](./12-release-command.md) — release all pending branches
- [prune](./51-prune.md) — clean up stale branches
