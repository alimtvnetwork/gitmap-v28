# Repo Aliases — Shortcut Names for Repositories

## Overview

Repo aliases let users assign short, memorable names to any tracked
repository. Once aliased, any `gitmap-v28` command can target the repo by
alias instead of requiring `cd` or a full slug/path — enabling remote
execution from anywhere in the filesystem.

---

## Goals

1. **Short names** — replace long slugs or paths with concise aliases
   (e.g., `api`, `web`, `infra`).
2. **Run from anywhere** — execute any gitmap-v28 command against a repo
   using `-A <alias>` without changing directory.
3. **Auto-suggest** — during `scan`/`rescan`, suggest aliases for newly
   discovered repos based on their slug or repo name.
4. **Persistent** — aliases stored in a dedicated SQLite table, surviving
   across sessions and terminal restarts.
5. **Conflict-safe** — warn and prompt when an alias collision occurs.

---

## Command: `gitmap-v28 alias`

Alias: `a`

Manages repo alias definitions stored in SQLite.

### Subcommands

| Subcommand | Alias | Description |
|------------|-------|-------------|
| `set`      | —     | Create or update an alias for a repo |
| `remove`   | —     | Remove an alias |
| `list`     | —     | List all aliases |
| `show`     | —     | Show the repo linked to an alias |
| `suggest`  | —     | Auto-suggest aliases for unaliased repos |

### Usage

```
gitmap-v28 alias set <alias> <slug-or-path>
gitmap-v28 alias remove <alias>
gitmap-v28 alias list
gitmap-v28 alias show <alias>
gitmap-v28 alias suggest [--apply]
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| --apply | false | Automatically accept all suggestions (used with `suggest`) |

---

## Global Flag: `-A` / `--alias`

Any gitmap-v28 command that operates on a repository accepts `-A <alias>`
to resolve the target repo by its alias.

### Usage

```bash
# Pull a repo by alias
gitmap-v28 pull -A api

# Run a shell command in the aliased repo directory
gitmap-v28 exec -A web -- npm test

# Navigate to the aliased repo
gitmap-v28 cd -A infra

# Check status of an aliased repo
gitmap-v28 status -A api

# Watch an aliased repo
gitmap-v28 watch -A web
```

The `-A` flag resolves the alias to the repo's `AbsolutePath` in the
database, then executes the command as if run from that directory.

---

## Auto-Suggestion During Scan

When `gitmap-v28 scan` or `gitmap-v28 rescan` discovers new repositories, the
system suggests aliases based on:

1. **Repo name** — last segment of the clone URL (e.g., `gitmap-v28`).
2. **Slug suffix** — last segment of the slug if different from repo name.
3. **Short unique prefix** — shortest unique prefix across all repos.

### Suggestion Flow

```
  Scan complete. 3 new repos found.

  Suggested aliases:
    api-gateway  → api       Accept? (y/N): y
    web-frontend → web       Accept? (y/N): y
    shared-libs  → libs      Accept? (y/N): n

  ✓ Created 2 alias(es).
```

### Non-Interactive Mode

```bash
gitmap-v28 scan --apply-aliases    # Accept all suggestions automatically
gitmap-v28 alias suggest --apply   # Suggest for existing unaliased repos
```

---

## Conflict Handling

When a user creates an alias that already exists:

```
  ⚠ Alias "api" already points to github/user/old-api.
  → Reassign to github/user/new-api? (y/N):
```

- **y** → update the alias to point to the new repo.
- **n** → abort, alias unchanged.

During auto-suggestion, conflicting aliases are skipped with a warning:

```
  ⚠ Skipping "api" — alias already assigned to github/user/old-api.
```

---

## Data Model

### SQLite Table: `Aliases`

| Column    | Type | Constraints |
|-----------|------|-------------|
| Id        | TEXT | PRIMARY KEY |
| Alias     | TEXT | NOT NULL UNIQUE |
| RepoId    | TEXT | NOT NULL REFERENCES Repos(Id) ON DELETE CASCADE |
| CreatedAt | TEXT | DEFAULT CURRENT_TIMESTAMP |

### Model Struct

```go
// Alias links a short name to a repository.
type Alias struct {
    ID        string `json:"id"`
    Alias     string `json:"alias"`
    RepoID    string `json:"repoId"`
    CreatedAt string `json:"createdAt"`
}
```

---

## Alias Resolution

When `-A <alias>` is provided:

1. Query `Aliases` table for matching `Alias` value.
2. Join with `Repos` table to get `AbsolutePath`.
3. If not found → print error: `no alias found: <alias>`.
4. If found → set working directory context to `AbsolutePath`.
5. Execute the command as normal.

### Resolution Priority

If both `-A` and a slug/path argument are provided, `-A` takes
precedence with a warning:

```
  ⚠ Both alias and slug provided — using alias "api".
```

---

## Shell Completion Integration

The `-A` flag should support tab-completion for alias names. Add a
`--list-aliases` data flag to the `completion` system:

```bash
gitmap-v28 completion --list-aliases
# Output: api web infra libs
```

This feeds into Bash/Zsh/PowerShell completion generators.

---

## Interaction with Existing Commands

| Command | Behavior with `-A` |
|---------|-------------------|
| `cd`    | Navigate to aliased repo directory |
| `pull`  | Pull in aliased repo |
| `exec`  | Execute command in aliased repo directory |
| `status`| Show status of aliased repo |
| `watch` | Watch aliased repo for changes |
| `release`| Run release from aliased repo |
| `scan`  | No effect (scan operates on directories, not repos) |
| `group` | No effect (operates on group names) |

Commands that don't operate on a single repo (e.g., `scan`, `list`,
`group`) ignore the `-A` flag.

---

## Package Structure

| File | Responsibility |
|------|----------------|
| `cmd/alias.go` | Subcommand dispatch (set/remove/list/show/suggest) |
| `cmd/aliasops.go` | Subcommand implementation |
| `cmd/aliasresolve.go` | `-A` flag resolution logic |
| `store/alias.go` | Database CRUD for Aliases table |
| `model/alias.go` | Data struct |
| `constants/constants_alias.go` | Messages, SQL, flag descriptions |
| `helptext/alias.md` | Command help |

---

## Dry-Run and Verbose Support

When `--verbose` is active with `-A`:

```
  → Resolved alias "api" → /home/user/repos/api-gateway (slug: github/user/api-gateway)
```

---

## Error Handling

- Alias not found → `no alias found: <alias>` (exit 1).
- Repo deleted but alias remains → cascade delete handles cleanup.
- Empty alias name → `alias name cannot be empty` (exit 1).
- Alias contains spaces or special chars → `alias must be alphanumeric with hyphens` (exit 1).
- Alias name conflicts with a gitmap-v28 subcommand → `alias cannot shadow command: <name>` (exit 1).

---

## Acceptance Criteria

1. `gitmap-v28 alias set api github/user/api-gateway` creates an alias.
2. `gitmap-v28 alias list` displays all aliases with their target slugs.
3. `gitmap-v28 alias show api` shows the linked repo details.
4. `gitmap-v28 alias remove api` deletes the alias.
5. `gitmap-v28 pull -A api` pulls the repo from any working directory.
6. `gitmap-v28 exec -A web -- npm test` runs the command in the repo dir.
7. `gitmap-v28 cd -A infra` navigates to the repo directory.
8. `gitmap-v28 alias suggest` proposes aliases for unaliased repos.
9. `gitmap-v28 alias suggest --apply` auto-accepts all suggestions.
10. `gitmap-v28 scan` suggests aliases for newly discovered repos.
11. Duplicate alias → warn and prompt for reassignment.
12. Deleting a repo cascades to remove its alias.
13. Tab-completion works for alias names via `--list-aliases`.
14. Alias names cannot shadow existing gitmap-v28 subcommands.
