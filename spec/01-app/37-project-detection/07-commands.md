# Project Type Detection — Commands

## Integrated into Scan

Project detection runs automatically as part of `scan` and `rescan`.
No additional flags are needed. After repo discovery and record
building, the scan pipeline adds a project detection phase.

---

## Query Commands

| Command              | Alias | Description                        |
|----------------------|-------|------------------------------------|
| `gitmap-v28 go-repos`    | `gr`  | List repos containing Go projects  |
| `gitmap-v28 node-repos`  | `nr`  | List repos containing Node projects|
| `gitmap-v28 react-repos` | `rr`  | List repos containing React projects|
| `gitmap-v28 cpp-repos`   | `cr`  | List repos containing C++ projects |
| `gitmap-v28 csharp-repos`| `sr`  | List repos containing C# projects  |

---

## Query Command Output

Terminal output for each detected project:

```
  go  github.com/user/my-api
      Path: /home/user/repos/my-api
      Indicator: go.mod

  go  github.com/user/my-cli/tools/linter
      Path: /home/user/repos/my-cli/tools/linter
      Indicator: go.mod
```

---

## Query Command Flags

| Flag       | Default    | Description                       |
|------------|------------|-----------------------------------|
| `--json`   | false      | Output as JSON instead of terminal|
| `--count`  | false      | Print count only                  |

---

## Query Command Data Source

Query commands read from the SQLite database. If the database does
not exist, print: `"No database found. Run 'gitmap-v28 scan' first."`
