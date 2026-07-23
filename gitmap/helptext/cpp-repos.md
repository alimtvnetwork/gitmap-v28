# gitmap cpp-repos

List all detected C++ projects across tracked repositories.

## Alias

cr

## Usage

    gitmap cpp-repos [--json]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Output as structured JSON |

## Prerequisites

- Run `gitmap scan` first to detect projects (see scan.md)

## Examples

### Example 1: List all C++ projects

    gitmap cpp-repos

**Output:**

    REPO          BUILD SYSTEM  STANDARD  PATH
    game-engine   CMake         C++20     D:\repos\game-engine
    codec-lib     Makefile      C++17     D:\repos\codec-lib
    renderer      CMake         C++20     D:\repos\renderer
    3 C++ projects detected

### Example 2: JSON output

    gitmap cr --json

**Output:**

    [
      {"repo":"game-engine","build_system":"CMake","standard":"C++20","path":"D:\\repos\\game-engine"},
      {"repo":"codec-lib","build_system":"Makefile","standard":"C++17","path":"D:\\repos\\codec-lib"}
    ]

### Example 3: No C++ projects found

    gitmap cpp-repos

**Output:**

    No C++ projects detected.
    → Run 'gitmap scan' to detect projects in your repos.

## See Also

- [scan](scan.md) — Scan directories to detect projects
- [csharp-repos](csharp-repos.md) — List C# projects
- [go-repos](go-repos.md) — List Go projects
- [node-repos](node-repos.md) — List Node.js projects

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter cpp-repos
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
