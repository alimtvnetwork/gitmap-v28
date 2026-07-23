# gitmap cd

Navigate to a tracked repository directory using its slug or an interactive picker.

## Alias

go

## Setup

Run `gitmap setup` to auto-install the `gcd` shell function. After that,
use `gcd <repo-name>` to change directory directly.

If you prefer to install it manually, add this to your shell profile:

**PowerShell:** `function gcd { Set-Location (gitmap cd @args) }`
**Bash/Zsh:** `gcd() { cd "$(gitmap cd "$@")" ; }`

## Usage

    gcd <repo-name>
    gitmap cd <repo-name|repos> [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| -A, --alias \<name\> | — | Navigate to a repo by its alias |
| --group \<name\> | — | Filter picker to a specific group |
| --pick | false | Force interactive picker |
| --default \<slug\> | — | Set or clear the default repo |

## Prerequisites

- Run `gitmap scan` first to populate the database (see scan.md)

## Examples

### Example 1: Navigate to a repo by slug

    gitmap cd my-api

**Output:**

    Changed directory to D:\wp-work\repos\my-api
    Branch: main | Status: clean | 0 ahead / 0 behind

### Example 2: Interactive repo picker

    gitmap cd repos

**Output:**

     #  REPO             BRANCH     PATH
     1. my-api           main       D:\wp-work\repos\my-api
     2. web-app          develop    D:\wp-work\repos\web-app
     3. billing-svc      main       D:\wp-work\repos\billing-svc
     4. auth-gateway     main       D:\wp-work\repos\auth-gateway
     5. shared-lib       main       D:\wp-work\repos\shared-lib
    Enter number (1-5): _

### Example 3: Pick from a specific group

    gitmap cd repos --group backend

**Output:**

     #  REPO             BRANCH     PATH
     1. billing-svc      main       D:\wp-work\repos\billing-svc
     2. auth-gateway     main       D:\wp-work\repos\auth-gateway
     3. payments-api     develop    D:\wp-work\repos\payments-api
    Enter number (1-3): _

### Example 4: Navigate by alias

    gitmap cd -A api

**Output:**

    Changed directory to D:\wp-work\repos\my-api
    Branch: main | Status: clean

### Example 5: Set a default repo

    gitmap cd --default my-api

**Output:**

    ✓ Default repo set: my-api
    (Run 'gitmap cd' without args to jump here)

## See Also

- [setup](setup.md) — Auto-install gcd and shell completions
- [list](list.md) — View all tracked repos and paths
- [group](group.md) — Manage repo groups for filtered navigation
- [scan](scan.md) — Scan directories to populate the database
- [alias](alias.md) — Manage repo aliases

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter cd
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
