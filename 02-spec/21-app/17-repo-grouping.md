# Repo Grouping

## Overview

Groups allow users to organize tracked repositories into named
collections for batch operations. Groups are stored in the SQLite
database alongside repo records (see [16-database.md](./16-database.md)).

---

## Data Model

### groups Table

| Column      | Type    | Constraints     | Notes                        |
|-------------|---------|-----------------|------------------------------|
| id          | TEXT    | PRIMARY KEY     | UUID                         |
| name        | TEXT    | NOT NULL UNIQUE | Group display name           |
| description | TEXT    | DEFAULT ''      | Optional description         |
| color       | TEXT    | DEFAULT ''      | Terminal color (e.g. "green") |
| created_at  | TEXT    | DEFAULT CURRENT_TIMESTAMP |                    |

### group_repos Table (Join)

| Column   | Type | Constraints                              |
|----------|------|------------------------------------------|
| group_id | TEXT | NOT NULL, FK → groups(id) ON DELETE CASCADE |
| repo_id  | TEXT | NOT NULL, FK → repos(id) ON DELETE CASCADE  |
| PRIMARY KEY | | (group_id, repo_id)                      |

### Group Model (model/group.go)

```go
type Group struct {
    ID          string
    Name        string
    Description string
    Color       string
    CreatedAt   string
}
```

---

## CLI Commands

### `gitmap list` (alias: `ls`)

Show all tracked repos with slugs.

```
SLUG                 REPO NAME
──────────────────────────────────────────
my-api               My API
my-api               My API (personal)
dashboard            Dashboard
auth-service         Auth Service
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--group` | `-g` | Filter by group name |
| `--verbose` | `-V` | Show full paths and URLs |

### `gitmap group create <name>` (alias: `g create`)

Create a new group.

```
gitmap group create backend
gitmap group create backend --description "All backend services"
gitmap group create backend --color cyan
```

### `gitmap group add <group> <slug...>` (alias: `g add`)

Add one or more repos to a group by slug.

```
gitmap group add backend my-api auth-service
```

If a slug is duplicated, disambiguation triggers (interactive or
`slug@path`).

### `gitmap group remove <group> <slug...>` (alias: `g rm`)

Remove repos from a group.

```
gitmap group remove backend my-api
```

### `gitmap group list` (alias: `g ls`)

List all groups with repo counts.

```
GROUP           REPOS   DESCRIPTION
──────────────────────────────────────────
backend         3       All backend services
frontend        2       UI applications
```

### `gitmap group show <name>` (alias: `g show`)

Show repos in a specific group.

```
Group: backend (3 repos)
  my-api           /home/user/work/my-api
  auth-service     /home/user/work/auth-service
  gateway          /home/user/work/gateway
```

### `gitmap group delete <name>` (alias: `g del`)

Delete a group (does not delete repos, only the grouping).

---

## Batch Operations on Groups

All existing repo-level commands support a `--group` (`-g`) flag and
an `--all` flag to target repos in bulk:

| Command | `--group` Example | `--all` Example |
|---------|-------------------|-----------------|
| `pull` | `gitmap pull --group backend` | `gitmap pull --all` |
| `exec` | `gitmap exec --group backend "git fetch --all"` | `gitmap exec --all "git fetch --all"` |
| `status` | `gitmap status --group backend` | `gitmap status --all` |
| `release` | `gitmap release --group backend` | — |
| `clone` | `gitmap clone json --group backend` | — |

### Selective Multi-Repo

Select specific repos by slug (comma-separated or repeated flag):

```
gitmap pull my-api,auth-service
gitmap pull my-api auth-service
gitmap status my-api auth-service
```

---

## Package Structure (Grouping)

### New Files

| File | Contents |
|------|----------|
| `store/group.go` | Group CRUD (create, list, add/remove repos) |
| `model/group.go` | Group and GroupRepo structs |
| `cmd/list.go` | `list` command handler |
| `cmd/group.go` | `group` command routing |
| `cmd/groupcreate.go` | `group create` handler |
| `cmd/groupadd.go` | `group add` handler |
| `cmd/groupremove.go` | `group remove` handler |
| `cmd/grouplist.go` | `group list` handler |
| `cmd/groupshow.go` | `group show` handler |
| `cmd/groupdelete.go` | `group delete` handler |

### Updated Files

| File | Change |
|------|--------|
| `cmd/root.go` | Register `list`, `group` commands |
| `cmd/pull.go` | Add `--group`, `--all` flags |
| `cmd/exec.go` | Add `--group`, `--all` flags |
| `cmd/status.go` | Add `--group`, `--all` flags |
| `constants/constants_cli.go` | New command names, aliases, help text |

---

## Error Handling (Grouping)

| Scenario | Behavior |
|----------|----------|
| Group not found | `"No group found: %s"` |
| Repo already in group | Silent no-op |
| Group name already exists | `"Group already exists: %s"` |

---

## Constraints

- All string literals in `constants` package.
- All files under 200 lines.
- All functions 8–15 lines.
- Positive conditions only (no negation).
- Blank line before `return`.
