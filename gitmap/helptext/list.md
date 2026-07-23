# gitmap list

Show all tracked repositories with their slugs and paths.

## Alias

ls

## Usage

    gitmap list [--group <name>] [--verbose]
    gitmap ls go              List only Go projects
    gitmap ls node            List only Node.js projects
    gitmap ls react           List only React projects
    gitmap ls cpp             List only C++ projects
    gitmap ls csharp          List only C# projects
    gitmap ls groups          List all groups

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --group \<name\> | — | Filter to a specific group |
| --verbose | false | Show full paths and metadata |

## Prerequisites

- Run `gitmap scan` first to populate the database (see scan.md)

## Examples

### Example 1: List all tracked repos

    gitmap list

**Output:**

    REPO             PATH
    my-api           D:\wp-work\repos\my-api
    web-app          D:\wp-work\repos\web-app
    billing-svc      D:\wp-work\repos\billing-svc
    auth-gateway     D:\wp-work\repos\auth-gateway
    shared-lib       D:\wp-work\repos\shared-lib
    5 repos tracked

### Example 2: List only Go projects with metadata

    gitmap ls go --verbose

**Output:**

    REPO          MODULE                          GO      PATH
    my-api        github.com/user/my-api          1.22    D:\wp-work\repos\my-api
    shared-lib    github.com/user/shared-lib      1.21    D:\wp-work\repos\shared-lib
    gitmap        github.com/alimtvnetwork/gitmap-v28/gitmap          1.22    D:\wp-work\repos\gitmap
    3 Go projects detected

### Example 3: List all groups with member counts

    gitmap ls groups

**Output:**

    GROUP           REPOS   DESCRIPTION
    backend         5       All backend microservices
    frontend        3       React frontend applications
    infra           2       Infrastructure and tooling
    3 groups defined

### Example 4: List repos in a specific group

    gitmap list --group backend

**Output:**

    REPO             PATH
    billing-svc      D:\wp-work\repos\billing-svc
    auth-gateway     D:\wp-work\repos\auth-gateway
    payments-api     D:\wp-work\repos\payments-api
    3 repos in group 'backend'

## See Also

- [cd](cd.md) — Navigate to a tracked repo
- [group](group.md) — Manage repo groups
- [scan](scan.md) — Scan directories to populate the database
- [status](status.md) — View repo statuses

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter list
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
