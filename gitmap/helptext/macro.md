# Macro

Record, replay, automate, and loop command sequences with environment variable expansion, sleep timers, AI error diagnostics, and structured JSON/YAML reporting.

## Aliases

`m`, `macros`, `retry`, `loop`

## Usage

    gitmap macro add <name> <cmd1> [cmd2...] [--desc <text>] [--tag <tag>]
    gitmap macro run <name> [--json] [--yaml] [--file <path>] [--dry-run] [--verbose]
    gitmap macro run-until-succeed <name|"cmd"> [--sleep <sec>] [--max-retries <N>] [--backoff fixed|linear|exponential] [--ai]
    gitmap macro record <name>
    gitmap macro list [--json] [--yaml] [--file <path>]
    gitmap macro show <name> [--json] [--yaml] [--file <path>]
    gitmap macro rm <name>

## Subcommands

| Subcommand | Description |
|---|---|
| `add` (`create`, `new`) `<name> <steps...>` | Create a new macro directly from command line arguments |
| `run-until-succeed` (`retry`, `until-success`, `loop`) `<name\|"cmd">` | Execute macro or command repeatedly until success (with sleep, backoff & AI diagnostics) |
| `record` (`rec`) `<name>` | Start interactive recording session for shell commands |
| `run` (`exec`) `<name>` | Replay a recorded macro sequence with optional JSON/YAML export |
| `list` (`ls`) | List all saved macros |
| `show` `<name>` | View steps in a recorded macro |
| `rm` (`delete`) `<name>` | Delete a saved macro |

## Flags

| Flag | Description |
|---|---|
| `--sleep <duration>`, `--delay <duration>` | Sleep interval between retry attempts (default: `5s`) |
| `--max-retries <N>`, `-n <N>` | Maximum number of retry attempts before stopping (default: unlimited) |
| `--timeout <duration>` | Maximum overall timeout for retry loop (e.g. `10m`, `1h`) |
| `--backoff <strategy>` | Backoff strategy: `fixed` (default), `linear`, or `exponential` |
| `--ai` | Output AI-ready diagnostic report on failure |
| `--ai-file <path>`, `--error-file <path>` | Save AI failure prompt report to specified file |
| `--json` | Output execution report in formatted JSON |
| `--yaml`, `-y` | Output execution report in formatted YAML |
| `--file <path>`, `-o <path>` | Save execution report to file |
| `--dry-run` | Simulate macro execution without invoking commands |
| `--verbose`, `-v` | Show live command stdout and stderr |

## Examples

```bash
# Add a new macro with multiple commands
gitmap macro add build-test "go build -o gitmap.exe ." "go test ./..." --desc "Build and test"

# Add a chained command macro
gitmap macro add deploy "npm run build && npm run test && npm run deploy"

# Run macro until it succeeds with 5s sleep interval
gitmap macro run-until-succeed build-test --sleep 5s

# Run arbitrary command until success with exponential backoff & AI error diagnostic report
gitmap macro run-until-succeed "go test -v ./..." --sleep 2s --backoff exponential --ai

# Run via top-level shortcut
gitmap retry "npm run build" --sleep 3s --max-retries 5 --ai-file .lovable/temp/error.md

# Replay a recorded macro
gitmap macro run deploy-temp
```
