# gitmap group

Manage repository groups and activate a group for batch operations.

## Alias

g

## Usage

    gitmap group <create|add|remove|list|show|delete|pull|status|exec|clear> [args]
    gitmap g <name>           Activate a group
    gitmap g                  Show active group
    gitmap g pull             Pull repos in active group
    gitmap g status           Show status for active group
    gitmap g exec <args>      Run git across active group
    gitmap g clear            Clear active group

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --desc \<text\> | — | Group description (for create) |
| --color \<hex\> | — | Group color (for create) |

## Prerequisites

- Run `gitmap scan` first to populate the database (see scan.md)

## Examples

### Example 1: Create a group and add repos

    gitmap group create backend --desc "All backend services"
    gitmap group add backend my-api billing-svc auth-gateway

**Output:**

    ✓ Group 'backend' created
    ✓ Added 3 repos to group 'backend'

### Example 2: Activate a group and pull

    gitmap g backend
    gitmap g pull

**Output:**

    Active group set: backend (3 repos)

    Pulling 3 repos in group 'backend'...
    [1/3] my-api (main)... Already up to date.
    [2/3] billing-svc (main)... updated (2 commits)
    [3/3] auth-gateway (main)... Already up to date.
    ✓ 3 repos pulled (1 updated, 2 up to date)

### Example 3: Show group members

    gitmap group show backend

**Output:**

    Group: backend
    Description: All backend services
    Repos (3):
      my-api           D:\wp-work\repos\my-api
      billing-svc      D:\wp-work\repos\billing-svc
      auth-gateway     D:\wp-work\repos\auth-gateway

### Example 4: Run git status across active group

    gitmap g exec status --short

**Output:**

    [my-api] (clean)
    [billing-svc] M  src/handler.go
    [auth-gateway] (clean)
    ✓ 3 repos processed

### Example 5: List all groups

    gitmap group list

**Output:**

    GROUP           REPOS   DESCRIPTION
    backend         3       All backend services
    frontend        2       React frontend apps
    2 groups defined

## See Also

- [list](list.md) — View all tracked repos
- [multi-group](multi-group.md) — Select multiple groups
- [pull](pull.md) — Pull repos by group
- [status](status.md) — View status by group

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter group
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
