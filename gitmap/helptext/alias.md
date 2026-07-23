# gitmap alias

Assign short names to repositories for quick access from anywhere.

## Alias

a

## Usage

    gitmap alias <subcommand> [arguments]

## Subcommands

| Subcommand | Description |
|------------|-------------|
| set        | Create or update an alias for a repo |
| remove     | Remove an alias |
| list       | List all aliases |
| show       | Show the repo linked to an alias |
| suggest    | Auto-suggest aliases for unaliased repos |

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --apply | false | Auto-accept all suggestions (with suggest) |

## Prerequisites

- Run `gitmap scan` first to populate the database (see scan.md)

## Examples

### Example 1: Create an alias and use it

    gitmap alias set api github/user/api-gateway
    gitmap pull -A api

**Output:**

    ✓ Alias "api" → github/user/api-gateway

    Pulling api-gateway (main)...
    Already up to date.

### Example 2: Auto-suggest aliases for all unaliased repos

    gitmap alias suggest

**Output:**

    Suggesting aliases for 5 unaliased repos...
    api-gateway    → api       Accept? (y/N): y
    ✓ Alias "api" created
    web-frontend   → web       Accept? (y/N): y
    ✓ Alias "web" created
    billing-svc    → billing   Accept? (y/N): n
    (skipped)
    auth-service   → auth      Accept? (y/N): y
    ✓ Alias "auth" created
    ✓ 3 aliases created, 1 skipped

### Example 3: Auto-accept all suggestions

    gitmap alias suggest --apply

**Output:**

    ✓ api-gateway → api
    ✓ web-frontend → web
    ✓ billing-svc → billing
    ✓ auth-service → auth
    ✓ 4 aliases created

### Example 4: List all aliases with paths

    gitmap alias list

**Output:**

    ALIAS   REPO                          PATH
    api     github/user/api-gateway       D:\repos\api-gateway
    web     github/user/web-frontend      D:\repos\web-frontend
    auth    github/user/auth-service      D:\repos\auth-service
    3 aliases defined

### Example 5: Show details for a specific alias

    gitmap alias show api

**Output:**

    Alias: api
    Repo:  github/user/api-gateway
    Path:  D:\repos\api-gateway
    Branch: main

## See Also

- [cd](cd.md) — Navigate to a repository (supports -A flag)
- [exec](exec.md) — Run commands in a repository (supports -A flag)
- [pull](pull.md) — Pull a repository (supports -A flag)
- [list](list.md) — List tracked repositories

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter alias
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
