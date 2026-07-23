# Commit Author Override & Amend Command

## Overview

Two related features for controlling Git commit authorship:

1. **SEO-write author flags** — set custom author name/email on each commit during `seo-write`.
2. **`gitmap amend` command** — rewrite author name/email on existing commits (all or from a specific commit onwards).

---

## Feature 1: SEO-Write Author Flags

### New Flags

| Flag              | Description                              | Default              |
|-------------------|------------------------------------------|----------------------|
| `--author-name`   | Git author name for commits              | (current git config) |
| `--author-email`  | Git author email for commits             | (current git config) |

### Behavior

- When provided, each `git commit` in the seo-write loop uses:
  ```
  git commit -m "message" --author="Name <email>"
  ```
- If only `--author-name` is provided without `--author-email`, use the name with the current git config email.
- If only `--author-email` is provided without `--author-name`, use the current git config name with the provided email.
- Dry-run mode (`--dry-run`) should display the author that would be used.

### Examples

```bash
# SEO-write with custom author
gitmap sw --url example.com --service Plumbing --area London \
  --author-name "John Smith" --author-email "john@example.com"

# Only override name (email stays from git config)
gitmap sw --url example.com --service SEO --area Bristol \
  --author-name "Marketing Bot"
```

---

## Feature 2: `gitmap amend` Command

### Synopsis

```
gitmap amend [commit-hash] --name <name> --email <email> [--branch <branch>]
```

Alias: `gitmap am`

The **commit hash** (SHA) is always the **first positional argument** (before any flags).
If omitted, all commits on the target branch are rewritten.

### Branch Resolution

- `--branch <name>` — operate on a specific branch (checks it out first).
- No `--branch` — operates on the **current branch** (whatever `HEAD` points to).

### Modes

#### Mode 1: All Commits on Branch

```bash
gitmap amend --name "New Name" --email "new@email.com"
gitmap amend --branch develop --name "New Name" --email "new@email.com"
```

Rewrites **every commit** on the target branch.

#### Mode 2: From a Specific Commit Onwards

```bash
gitmap amend a1b2c3d --name "New Name" --email "new@email.com"
gitmap amend a1b2c3d --branch main --name "New Name" --email "new@email.com"
```

The SHA is the **first argument**. Rewrites all commits **from `a1b2c3d` (inclusive) to HEAD** of the target branch. Commits before `a1b2c3d` are left untouched.

#### Mode 3: Single Commit (HEAD only)

```bash
gitmap amend HEAD --name "New Name" --email "new@email.com"
```

Amends only the most recent commit on the current (or specified) branch.

### Argument Order

```
gitmap amend [SHA] [--flags...]
              ^
              first positional arg = commit hash (optional)
              everything else = named flags
```

### Flags

| Flag                | Description                              | Required |
|---------------------|------------------------------------------|----------|
| `--name <name>`     | New author name                          | Yes (at least one of name/email) |
| `--email <email>`   | New author email                         | Yes (at least one of name/email) |
| `--branch <branch>` | Target branch (default: current branch)  | No       |
| `--dry-run`         | Preview which commits would be amended   | No       |
| `--force-push`      | Auto-run `git push --force-with-lease` after amend | No |

### Implementation Approach

1. If `--branch` is provided, run `git checkout <branch>` first.
2. Resolve the commit range (all, from SHA, or HEAD).
3. Execute `git filter-branch`:

```bash
# All commits on branch
git filter-branch -f --env-filter '
  export GIT_AUTHOR_NAME="New Name"
  export GIT_AUTHOR_EMAIL="new@email.com"
  export GIT_COMMITTER_NAME="New Name"
  export GIT_COMMITTER_EMAIL="new@email.com"
' -- HEAD

# From specific SHA onwards
git filter-branch -f --env-filter '
  export GIT_AUTHOR_NAME="New Name"
  export GIT_AUTHOR_EMAIL="new@email.com"
  export GIT_COMMITTER_NAME="New Name"
  export GIT_COMMITTER_EMAIL="new@email.com"
' <commit-hash>^..HEAD
```

4. If `--branch` was used, switch back to the original branch.
5. Write an audit record to `.gitmap/` and persist to the `Amendments` database table.

### Audit Trail (`.gitmap/` folder)

Every amend operation writes a JSON file to `.gitmap/amendments/` in the repository root:

**File naming**: `.gitmap/amendments/amend-<timestamp>.json`

Example: `.gitmap/amendments/amend-2026-03-09T14-30-00.json`

**JSON structure**:

```json
{
  "timestamp": "2026-03-09T14:30:00Z",
  "branch": "develop",
  "fromCommit": "a1b2c3d",
  "toCommit": "9z8y7x6",
  "totalCommits": 15,
  "previousAuthor": {
    "name": "Old Name",
    "email": "old@email.com"
  },
  "newAuthor": {
    "name": "New Name",
    "email": "new@email.com"
  },
  "mode": "range",
  "forcePushed": false,
  "commits": [
    { "sha": "a1b2c3d", "message": "Fix login page" },
    { "sha": "def5678", "message": "Add dashboard" }
  ]
}
```

**`mode` values**: `"all"` (every commit), `"range"` (from SHA onwards), `"head"` (single HEAD commit).

The `.gitmap/` folder is created automatically if it doesn't exist. These files serve as a local audit log that can be reviewed, diffed, or committed.

### Database Persistence (`Amendments` table)

Each amend operation is also persisted to the SQLite database for queryable history.

**Table schema** (PascalCase, following project convention):

```sql
CREATE TABLE IF NOT EXISTS Amendments (
    Id            TEXT PRIMARY KEY,
    Branch        TEXT NOT NULL,
    FromCommit    TEXT NOT NULL,
    ToCommit      TEXT NOT NULL,
    TotalCommits  INTEGER NOT NULL,
    PreviousName  TEXT DEFAULT '',
    PreviousEmail TEXT DEFAULT '',
    NewName       TEXT DEFAULT '',
    NewEmail      TEXT DEFAULT '',
    Mode          TEXT NOT NULL,
    ForcePushed   INTEGER DEFAULT 0,
    CreatedAt     TEXT DEFAULT CURRENT_TIMESTAMP
)
```

**SQL operations**:

```sql
-- Insert
INSERT INTO Amendments (Id, Branch, FromCommit, ToCommit, TotalCommits, PreviousName, PreviousEmail, NewName, NewEmail, Mode, ForcePushed)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)

-- Query all (most recent first)
SELECT Id, Branch, FromCommit, ToCommit, TotalCommits, PreviousName, PreviousEmail, NewName, NewEmail, Mode, ForcePushed, CreatedAt
FROM Amendments ORDER BY CreatedAt DESC

-- Query by branch
SELECT Id, Branch, FromCommit, ToCommit, TotalCommits, PreviousName, PreviousEmail, NewName, NewEmail, Mode, ForcePushed, CreatedAt
FROM Amendments WHERE Branch = ? ORDER BY CreatedAt DESC
```

The table is added to the `Migrate()` call in `store.go` and reset in `db-reset`.

### Safety

- **Warning prompt**: Before executing, print a warning that this rewrites history and requires force-push. Proceed automatically (no interactive prompt — follows project convention).
- **Backup ref**: Git automatically creates `refs/original/` backup refs during filter-branch.
- **Dry-run**: List all commits that would be affected with their current author and the new author. No audit file or DB record is written in dry-run mode.

### Terminal Output

```
amend: rewriting 15 commits from abc1234..HEAD (branch: develop)
  author: "Old Name <old@email.com>" -> "New Name <new@email.com>"

  [1/15] abc1234 - Fix login page
  [2/15] def5678 - Add dashboard
  ...
  [15/15] 9z8y7x6 - Update README

Done: 15 commits amended
  Audit log: .gitmap/amendments/amend-2026-03-09T14-30-00.json
  Database:  1 record saved to Amendments table
Warning: Run 'git push --force-with-lease' to update the remote
```

### Examples

```bash
# Amend all commits on current branch
gitmap amend --name "John Smith" --email "john@company.com"
gitmap am --name "John Smith" --email "john@company.com"

# Amend all commits on a specific branch
gitmap amend --branch develop --name "John Smith" --email "john@company.com"

# Amend from a specific SHA onwards (SHA is first positional arg)
gitmap amend a1b2c3d --name "John Smith" --email "john@company.com"

# Amend from SHA on a specific branch
gitmap amend a1b2c3d --branch main --name "John Smith" --email "john@company.com"

# Amend only HEAD
gitmap amend HEAD --name "John Smith" --email "john@company.com"

# Preview what would change (dry-run, no audit saved)
gitmap amend --name "John Smith" --email "john@company.com" --dry-run
gitmap amend a1b2c3d --branch develop --name "John Smith" --dry-run

# Amend and auto force-push
gitmap amend a1b2c3d --name "John Smith" --email "john@company.com" --force-push

# Only change email (keep existing author name)
gitmap amend --email "newemail@company.com"

# Only change name on a specific branch
gitmap amend --branch feature/auth --name "CI Bot"
```

---

## File Layout

| File | Purpose |
|------|---------|
| `constants/constants_amend.go` | Command/flag/message/SQL constants |
| `cmd/amend.go` | Flag parsing, orchestration |
| `cmd/amendexec.go` | Git filter-branch execution and output |
| `cmd/amendaudit.go` | Audit JSON file writing and DB persistence |
| `store/amendment.go` | Amendment CRUD operations |
| `model/amendment.go` | AmendmentRecord struct |

SEO-write changes modify existing files:
- `constants/constants_seo.go` — add `FlagSEOAuthorName`, `FlagSEOAuthorEmail`
- `cmd/seowrite.go` — add fields to `seoWriteFlags`
- `cmd/seowriteloop.go` — pass author to `gitCommit`

Store migration changes:
- `constants/constants_store.go` — add `TableAmendments`, `SQLCreateAmendments`, SQL operations
- `store/store.go` — add `SQLCreateAmendments` to `Migrate()`, `SQLDropAmendments` to `Reset()`

---

## Database Updates

### New Table

| Table | Purpose |
|-------|---------|
| `Amendments` | Audit log of author rewrite operations |

### Columns

| Column | Type | Description |
|--------|------|-------------|
| `Id` | TEXT PK | UUID |
| `Branch` | TEXT | Target branch name |
| `FromCommit` | TEXT | First commit SHA in range (or first commit if all) |
| `ToCommit` | TEXT | Last commit SHA (HEAD at time of amend) |
| `TotalCommits` | INTEGER | Number of commits rewritten |
| `PreviousName` | TEXT | Original author name |
| `PreviousEmail` | TEXT | Original author email |
| `NewName` | TEXT | New author name |
| `NewEmail` | TEXT | New author email |
| `Mode` | TEXT | `all`, `range`, or `head` |
| `ForcePushed` | INTEGER | 0 or 1 |
| `CreatedAt` | TEXT | Timestamp |

---

## CLI Interface Updates

### Command Table Addition

| Command | Alias | Description |
|---------|-------|-------------|
| `amend [hash]` | `am` | Rewrite author name/email on commits |

### Dispatch

Add to `dispatchMisc` in `root.go`.

---

## Acceptance Criteria

- [ ] `gitmap sw --author-name "Bot" --author-email "bot@co.com"` sets author on each commit
- [ ] `gitmap amend --name "X" --email "x@y.com"` rewrites all commits on current branch
- [ ] `gitmap amend --branch develop --name "X" --email "x@y.com"` rewrites all commits on develop
- [ ] `gitmap amend abc123 --name "X" --email "x@y.com"` rewrites from abc123 to HEAD
- [ ] `gitmap amend abc123 --branch main --name "X"` rewrites from abc123 on main branch
- [ ] `gitmap amend HEAD --name "X" --email "x@y.com"` amends only the latest commit
- [ ] `--dry-run` shows affected commits without modifying or saving anything
- [ ] `--force-push` runs `git push --force-with-lease` after amend
- [ ] At least one of `--name` or `--email` is required
- [ ] Terminal output shows progress per commit and target branch
- [ ] When `--branch` is used, switches back to original branch after completion
- [ ] Audit JSON written to `.gitmap/amendments/amend-<timestamp>.json` after each operation
- [ ] Amendment record persisted to `Amendments` SQLite table
- [ ] `db-reset` clears the `Amendments` table along with other tables
