# gitmap react-repos

List all detected React projects across tracked repositories.

## Alias

rr

## Usage

    gitmap react-repos [--json]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Output as structured JSON |

## Prerequisites

- Run `gitmap scan` first to detect projects (see scan.md)

## Examples

### Example 1: List all React projects

    gitmap react-repos

**Output:**

    REPO          PACKAGE             REACT VER  PATH
    web-app       @user/web-app       18.2.0     D:\repos\web-app
    docs-site     docs-site           18.2.0     D:\repos\docs-site
    admin-panel   @user/admin-panel   18.3.1     D:\repos\admin-panel
    3 React projects detected

### Example 2: JSON output

    gitmap rr --json

**Output:**

    [
      {"repo":"web-app","package":"@user/web-app","react_version":"18.2.0","path":"D:\\repos\\web-app"},
      {"repo":"docs-site","package":"docs-site","react_version":"18.2.0","path":"D:\\repos\\docs-site"},
      {"repo":"admin-panel","package":"@user/admin-panel","react_version":"18.3.1","path":"D:\\repos\\admin-panel"}
    ]

### Example 3: No React projects found

    gitmap react-repos

**Output:**

    No React projects detected.
    → Run 'gitmap scan' to detect projects in your repos.

## See Also

- [scan](scan.md) — Scan directories to detect projects
- [node-repos](node-repos.md) — List Node.js projects
- [go-repos](go-repos.md) — List Go projects
- [csharp-repos](csharp-repos.md) — List C# projects

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter react-repos
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
