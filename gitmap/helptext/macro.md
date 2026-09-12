# Macro

Record, replay, automate, and loop command sequences with environment variable expansion, sleep timers, AI error diagnostics, and structured JSON/YAML reporting.

## Aliases

`m`, `macros`, `retry`, `loop`

## Usage

    gitmap macro add <name> [cmd1...] [--desc <text>] [--tag <tag>] [--no-exec]
    gitmap macro edit <name> [--no-exec]
    gitmap macro run <name> [--json] [--yaml] [--file <path>] [--dry-run] [--verbose]
    gitmap macro run-until-succeed <name|"cmd"> [--sleep <sec>] [--max-retries <N>] [--backoff fixed|linear|exponential] [--ai]
    gitmap macro record <name>
    gitmap macro list [--json] [--yaml] [--file <path>]
    gitmap macro show <name> [--json] [--yaml] [--file <path>]
    gitmap macro rm <name>
    gitmap macro export <name|all> [-f <path>] [--format json|yaml|sqlite|zip] [--all] [-except <list>]
    gitmap macro import <file> [--format json|yaml|sqlite|zip] [--force] [--dry-run] [-except <list>]

## Subcommands

| Subcommand | Description |
|---|---|
| `add` (`create`, `new`) `<name> [steps...]` | Create a new macro from arguments or interactively with live execution |
| `edit` (`modify`) `<name>` | Interactively modify, insert, delete, or append steps with live command feedback |
| `run-until-succeed` (`retry`, `until-success`, `loop`) `<name\|"cmd">` | Execute macro or command repeatedly until success (with sleep, backoff & AI diagnostics) |
| `record` (`rec`) `<name>` | Start interactive recording session for shell commands |
| `run` (`exec`) `<name>` | Replay a recorded macro sequence with optional JSON/YAML export |
| `list` (`ls`) | List all saved macros |
| `show` `<name>` | View steps in a recorded macro |
| `rm` (`delete`) `<name>` | Delete a saved macro |
| `export` (`exp`, `dump`) `<name\|all>` | Export macro(s) to JSON, YAML, SQLite DB, or ZIP bundle |
| `import` (`imp`, `load`, `restore`) `<file>` | Import macro(s) safely with format auto-inference and overwrite guards |

## In-Builder Commands (Create & Edit Mode)

When creating (`gitmap macro add <name>`) or editing (`gitmap macro edit <name>`) interactively, every shell command entered executes live in the terminal with streamed output by default. The builder also supports dedicated helper commands:

| Command | Action |
|---|---|
| `cat <file>` / `view <file>` | Inspect file contents live in terminal |
| `touch <file>` | Create an empty file (auto-creates parent directories) |
| `mkfile <file> [content]` | Create a file with initial text content |
| `rmfile <file>` | Remove a file live during the session |
| `cpfile <src> <dst>` | Copy a file live during the session |
| `copy <text\|--file <path>>` | Copy text or file to OS clipboard and persistent memory buffer |
| `paste [--file <path>]` | Paste from memory buffer or OS clipboard to terminal or file |
| `explorer [path]` | Open native graphical file manager at path |
| `browse <url> [--chrome]` | Open URL in default browser or Google Chrome |
| `del <n>` / `rm <n>` | (Edit mode) Delete step `<n>` and re-index remaining steps |
| `replace <n> <cmd>` | (Edit mode) Replace step `<n>` with a new command |
| `insert <n> <cmd>` | (Edit mode) Insert a command at position `<n>` |
| `list` / `show` | Display the current step list |
| `exec on` / `exec off` | Toggle live execution of entered commands |
| `cd <directory>` | Change working directory live for the session |
| `pwd` / `pwd on` / `pwd off` | Display or toggle current working directory banner |
| `done` / `exit` / `quit` | Save macro steps and finalize |
| `cancel` / `abort` | Discard changes without saving |

## Flags

| Flag | Description |
|---|---|
| `--no-exec` | Disable live command execution during interactive creation or edit |
| `--sleep <duration>`, `--delay <duration>` | Sleep interval between retry attempts (default: `5s`) |
| `--max-retries <N>`, `-n <N>` | Maximum number of retry attempts before stopping (default: unlimited) |
| `--timeout <duration>` | Maximum overall timeout for retry loop (e.g. `10m`, `1h`) |
| `--backoff <strategy>` | Backoff strategy: `fixed` (default), `linear`, or `exponential` |
| `--ai` | Output AI-ready diagnostic report on failure |
| `--ai-file <path>`, `--error-file <path>` | Save AI failure prompt report to specified file |
| `--json` | Output execution report in formatted JSON |
| `--yaml`, `-y` | Output execution report in formatted YAML |
| `--file <path>`, `-o <path>` | Save execution report to file |
| `--dry-run` | Simulate macro execution or import without modifying disk |
| `--verbose`, `-v` | Show live command stdout and stderr |
| `--all` | Export all stored macros into a single bundle file or archive |
| `--sqlite`, `--db` | Export/import using portable SQLite database table format |
| `--zip` | Export/import using compressed ZIP archive with per-macro JSON files |
| `--force`, `--overwrite` | Overwrite existing saved macros during import |
| `-except <list>`, `--exclude <list>` | Comma-separated list of macro names to exclude from export/import |

## Examples

### 1. Interactive Macro Creation with File Operations & Live Output

```bash
# Start interactive builder
gitmap macro add setup-dev

# In builder:
# Step 1> mkdir -p configs
# Step 2> mkfile configs/app.env APP_ENV=development PORT=8080
# Step 3> cat configs/app.env
# Step 4> copy configs/app.env
# Step 5> gitmap explorer configs
# Step 6> done
```

### 2. Interactive Macro Editing

```bash
# Edit existing macro
gitmap macro edit setup-dev

# In editor:
# Edit [7]> list
# Edit [7]> replace 2 mkfile configs/app.env APP_ENV=staging PORT=3000
# Edit [7]> insert 3 cat configs/app.env
# Edit [8]> done
```

### 3. Command Line Add & Replay

```bash
# Add a multi-command macro directly
gitmap macro add build-test "go build -o gitmap.exe ." "go test ./..." --desc "Build and test"

# Replay saved macro
gitmap macro run build-test

# Direct native root-level execution
gitmap build-test
```

### 4. Retry Loop Until Success

```bash
# Retry until success with exponential backoff & AI diagnostics
gitmap macro run-until-succeed "npm run test" --sleep 2s --backoff exponential --ai
```

### 5. Multi-Format Macro Export & Import

```bash
# Export single macro to JSON or YAML
gitmap macro export setup-dev -f setup-dev.json
gitmap macro export setup-dev -f setup-dev.yaml --yaml

# Export all macros into a portable SQLite database
gitmap macro export --all -f macros.db --sqlite

# Export all macros into a compressed ZIP archive
gitmap macro export --all -f macros.zip --zip

# Export all except specific macros
gitmap macro export --all -f backup.json -except "test-run,debug-macro"

# Safely preview an import without writing to disk
gitmap macro import macros.db --dry-run

# Import macros with overwrite enabled
gitmap macro import macros.db --force
```
