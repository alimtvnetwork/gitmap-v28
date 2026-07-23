# gitmap version-history

Show version transitions recorded for the current repository.

## Alias

vh

## Usage

    gitmap version-history [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --limit \<N\> | 0 (all) | Show only the last N transitions |
| --json | false | Output as JSON instead of table |

## Prerequisites

- Must be run inside a Git repository that has been cloned via `gitmap clone-next`
- Version transitions are only recorded when clone-next is used

## Examples

### Example 1: Show all version transitions

    gitmap vh

**Output:**

    Version history for D:\wp-work\riseup-asia\macro-ahk:

    FROM        TO          FOLDER                    TIMESTAMP
    v11         v12         macro-ahk                 2026-04-16T10:30:00Z
    v12         v15         macro-ahk                 2026-04-16T14:22:00Z
    v15         v16         macro-ahk                 2026-04-16T16:45:00Z

    3 transition(s) recorded.

### Example 2: Show last 2 transitions as JSON

    gitmap vh --limit 2 --json

**Output:**

    [
      {
        "id": 3,
        "repoId": 42,
        "fromVersionTag": "v15",
        "fromVersionNum": 15,
        "toVersionTag": "v16",
        "toVersionNum": 16,
        "flattenedPath": "macro-ahk",
        "createdAt": "2026-04-16T16:45:00Z"
      },
      {
        "id": 2,
        "repoId": 42,
        "fromVersionTag": "v12",
        "fromVersionNum": 12,
        "toVersionTag": "v15",
        "toVersionNum": 15,
        "flattenedPath": "macro-ahk",
        "createdAt": "2026-04-16T14:22:00Z"
      }
    ]

## See Also

- [clone-next](clone-next.md) — Clone next version of a repo (auto-flattened)
- [history](history.md) — Show command execution audit log

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter version-history
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
