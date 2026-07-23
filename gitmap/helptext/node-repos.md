# gitmap node-repos

List all detected Node.js projects across tracked repositories.

## Alias

nr

## Usage

    gitmap node-repos [--json]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Output as structured JSON |

## Prerequisites

- Run `gitmap scan` first to detect projects (see scan.md)

## Examples

### Example 1: List all Node.js projects

    gitmap node-repos

**Output:**

    REPO          PACKAGE             NODE VER  PATH
    web-app       @user/web-app       18.x      D:\repos\web-app
    docs-site     docs-site           20.x      D:\repos\docs-site
    landing-page  landing-page        18.x      D:\repos\landing-page
    admin-panel   @user/admin-panel   20.x      D:\repos\admin-panel
    4 Node.js projects detected

### Example 2: JSON output

    gitmap nr --json

**Output:**

    [
      {"repo":"web-app","package":"@user/web-app","node_version":"18.x","path":"D:\\repos\\web-app"},
      {"repo":"docs-site","package":"docs-site","node_version":"20.x","path":"D:\\repos\\docs-site"},
      {"repo":"landing-page","package":"landing-page","node_version":"18.x","path":"D:\\repos\\landing-page"}
    ]

### Example 3: No Node.js projects found

    gitmap node-repos

**Output:**

    No Node.js projects detected.
    → Run 'gitmap scan' to detect projects in your repos.

## See Also

- [scan](scan.md) — Scan directories to detect projects
- [react-repos](react-repos.md) — List React projects
- [go-repos](go-repos.md) — List Go projects
- [cpp-repos](cpp-repos.md) — List C++ projects

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter node-repos
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
