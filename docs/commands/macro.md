# gitmap macro

Record, manage, edit, and replay multi-step command sequences as reusable macros, with live terminal execution, in-builder file operations, and support for direct root-level execution.

```bash
gitmap macro <subcommand> [arguments]
```

Aliases: `m`

---

## Direct Root-Level Macro Execution

Any macro saved in GitMap can be run **directly as a root-level command** by typing `gitmap <macro-name>`:

```bash
# Save a macro
gitmap macro add test-all "go test ./..." "npm run test"

# Run it directly as a native command!
gitmap test-all

# Pass flags directly to the macro
gitmap test-all --verbose
gitmap test-all --dry-run
```

---

## Subcommands & Root Utility Commands

| Subcommand | Root-Level Alias | Description |
|---|---|---|
| `gitmap macro list` | `gitmap macro-list`, `gitmap macro-ls` | List all saved macros and step counts (`--json`, `--yaml`) |
| `gitmap macro add <name> [steps...]` | `gitmap macro-add`, `gitmap macro-create` | Create a new macro with live terminal execution |
| `gitmap macro edit <name>` | `gitmap macro-edit`, `gitmap macro-modify` | Interactively edit, replace, insert, or delete steps in a macro |
| `gitmap macro run <name>` | `gitmap macro-run`, `gitmap execute` | Replay a macro with timing and step progress |
| `gitmap macro record <name>` | `gitmap macro-record`, `gitmap record` | Record an interactive shell session as a macro |
| `gitmap macro show <name>` | `gitmap macro-show` | Inspect the steps and parameters of a macro |
| `gitmap macro rm <name>` | `gitmap macro-rm`, `gitmap macro-del` | Delete a saved macro |
| `gitmap macro run-until-succeed <name>` | `gitmap retry`, `gitmap loop` | Retry macro until success with AI diagnostics |
| `gitmap macro export <name\|all>` | `gitmap macro-export`, `gitmap macro-exp` | Export macro(s) to JSON, YAML, SQLite DB, or ZIP bundle |
| `gitmap macro import <file>` | `gitmap macro-import`, `gitmap macro-imp` | Import macro(s) safely with format auto-inference and overwrite guards |

---

## In-Builder Commands (Create & Edit Mode)

When creating (`gitmap macro add <name>`) or editing (`gitmap macro edit <name>`) interactively, commands run live in your terminal with streamed output so you see exactly what happens in real time.

| Helper Command | Description |
|---|---|
| `cat <file>` / `view <file>` | Inspect file contents live in terminal |
| `touch <file>` | Create an empty file (auto-creates parent folders) |
| `mkfile <file> [content]` | Create a file with initial content |
| `rmfile <file>` | Remove a file live during the session |
| `cpfile <src> <dst>` | Copy a file live during the session |
| `copy <text\|--file <path>>` | Copy text or file to OS clipboard & persistent memory |
| `paste [--file <path>]` | Paste from memory buffer or OS clipboard |
| `explorer [path]` | Open native desktop file manager |
| `browse <url> [--chrome]` | Open URL in default browser or Google Chrome |
| `del <n>` / `rm <n>` | (Edit mode) Delete step `<n>` and re-index steps |
| `replace <n> <cmd>` | (Edit mode) Replace step `<n>` with a new command |
| `insert <n> <cmd>` | (Edit mode) Insert a command at position `<n>` |
| `list` / `show` | Display the current step list |
| `exec on` / `exec off` | Toggle live execution of entered commands |
| `cd <directory>` | Change working directory live for the session |
| `pwd` / `pwd on` / `pwd off` | Show or toggle current working directory |
| `done` / `exit` / `quit` | Save macro steps and exit |
| `cancel` / `abort` | Abort without saving |

---

## Comprehensive Examples

### 1. Interactive Macro Creation with Live Execution & File Ops

```bash
# Start interactive builder
gitmap macro add init-project

# Step 1> mkdir -p configs
# Step 2> mkfile configs/app.env "ENV=local PORT=8080"
# Step 3> cat configs/app.env
# Step 4> copy configs/app.env
# Step 5> gitmap explorer configs
# Step 6> gitmap browse http://localhost:8080 --chrome
# Step 7> done
```

### 2. Modifying an Existing Macro

```bash
# Modify steps in init-project
gitmap macro edit init-project

# Edit [8]> list
# Edit [8]> replace 2 mkfile configs/app.env "ENV=production PORT=443"
# Edit [8]> del 5
# Edit [7]> done
```

### 3. Running Macros & Passing Flags

```bash
# Replay via macro subcommand
gitmap macro run init-project

# Direct root shortcut execution
gitmap init-project --verbose
```

### 4. Exporting and Importing Macros

```bash
# Export single macro to JSON
gitmap macro export init-project -f init-project.json

# Export all macros into a portable SQLite database
gitmap macro export --all -f macros.db --sqlite

# Export all macros into a ZIP archive
gitmap macro export --all -f macros.zip --zip

# Preview importing macros without saving
gitmap macro import macros.db --dry-run

# Import macros and overwrite any existing macros with matching names
gitmap macro import macros.db --force
```
