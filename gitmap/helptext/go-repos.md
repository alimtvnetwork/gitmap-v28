# gitmap go-repos

List all detected Go projects across tracked repositories.

## Alias

gr

## Usage

    gitmap go-repos [--json]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Output as structured JSON |

## Prerequisites

- Run `gitmap scan` first to detect projects (see scan.md)

## Examples

### Example 1: List all Go projects

    gitmap go-repos

**Output:**

    REPO          MODULE                            GO VER  PATH
    my-api        github.com/user/my-api            1.22    D:\repos\my-api
    shared-lib    github.com/user/shared-lib        1.21    D:\repos\shared-lib
    gitmap        github.com/alimtvnetwork/gitmap-v28/gitmap            1.22    D:\repos\gitmap
    auth-service  github.com/user/auth-service      1.22    D:\repos\auth-service
    4 Go projects detected

### Example 2: JSON output

    gitmap gr --json

**Output:**

    [
      {
        "repo": "my-api",
        "module": "github.com/user/my-api",
        "go_version": "1.22",
        "path": "D:\\repos\\my-api"
      },
      {
        "repo": "shared-lib",
        "module": "github.com/user/shared-lib",
        "go_version": "1.21",
        "path": "D:\\repos\\shared-lib"
      }
    ]

### Example 3: No Go projects found

    gitmap go-repos

**Output:**

    No Go projects detected.
    → Run 'gitmap scan' to detect projects in your repos.

## See Also

- [scan](scan.md) — Scan directories to detect projects
- [node-repos](node-repos.md) — List Node.js projects
- [react-repos](react-repos.md) — List React projects
- [csharp-repos](csharp-repos.md) — List C# projects
- [gomod](gomod.md) — Rename Go module paths

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter go-repos
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
