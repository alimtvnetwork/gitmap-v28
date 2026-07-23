# gitmap csharp-repos

List all detected C# projects across tracked repositories.

## Alias

csr

## Usage

    gitmap csharp-repos [--json]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Output as structured JSON |

## Prerequisites

- Run `gitmap scan` first to detect projects (see scan.md)

## Examples

### Example 1: List all C# projects

    gitmap csharp-repos

**Output:**

    REPO          SOLUTION           TARGET    PATH
    billing-svc   BillingSvc.sln     net8.0    D:\repos\billing-svc
    auth-api      AuthApi.sln        net7.0    D:\repos\auth-api
    web-portal    WebPortal.sln      net8.0    D:\repos\web-portal
    3 C# projects detected

### Example 2: JSON output

    gitmap csr --json

**Output:**

    [
      {"repo":"billing-svc","solution":"BillingSvc.sln","target":"net8.0","path":"D:\\repos\\billing-svc"},
      {"repo":"auth-api","solution":"AuthApi.sln","target":"net7.0","path":"D:\\repos\\auth-api"},
      {"repo":"web-portal","solution":"WebPortal.sln","target":"net8.0","path":"D:\\repos\\web-portal"}
    ]

### Example 3: No C# projects found

    gitmap csharp-repos

**Output:**

    No C# projects detected.
    → Run 'gitmap scan' to detect projects in your repos.

## See Also

- [scan](scan.md) — Scan directories to detect projects
- [cpp-repos](cpp-repos.md) — List C++ projects
- [go-repos](go-repos.md) — List Go projects
- [node-repos](node-repos.md) — List Node.js projects

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter csharp-repos
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
