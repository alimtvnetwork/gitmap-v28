# gitmap exec

Run a git command across all tracked repositories.

## Alias

x

## Usage

    gitmap exec <git-args...>

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| -A, --alias \<name\> | — | Target a repo by its alias |

All other arguments are passed directly to git.

## Prerequisites

- Run `gitmap scan` first to populate the database (see scan.md)

## Examples

### Example 1: Fetch and prune across all repos

    gitmap exec fetch --prune

**Output:**

    [my-api] git fetch --prune...
      - [deleted] origin/feature/old-branch
      done
    [web-app] git fetch --prune... done
    [billing-svc] git fetch --prune... done
    ✓ 3 repos processed

### Example 2: Check remote URLs for all repos

    gitmap x remote -v

**Output:**

    [my-api]
      origin  https://github.com/user/my-api.git (fetch)
      origin  https://github.com/user/my-api.git (push)
    [web-app]
      origin  https://github.com/user/web-app.git (fetch)
      origin  https://github.com/user/web-app.git (push)
    ✓ 2 repos processed

### Example 3: Run git log across all repos

    gitmap exec log --oneline -3

**Output:**

    [my-api]
      abc1234 Fix auth middleware
      def5678 Add rate limiting
      ghi9012 Update dependencies
    [web-app]
      jkl3456 Redesign dashboard
      mno7890 Fix responsive layout
      pqr1234 Add dark mode toggle
    ✓ 2 repos processed

### Example 4: Run a command on a specific repo by alias

    gitmap exec -A api status --short

**Output:**

    [my-api]
      M  src/handler.go
      ?? src/new-file.go
    ✓ 1 repo processed

## See Also

- [scan](scan.md) — Scan directories to populate the database
- [pull](pull.md) — Pull repos (built-in alternative to exec fetch)
- [status](status.md) — View repo statuses
- [alias](alias.md) — Manage repo aliases

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter exec
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
