# Env — Cross-Platform Environment Variable Management

## Overview

The `env` command provides persistent, cross-platform management of
environment variables and PATH entries. Changes survive terminal
restarts and system reboots by writing directly to the OS-appropriate
configuration store.

---

## How It Works

### Windows
- Uses `os/exec` to call `setx` for User-level variables.
- With `--system` flag, calls `setx /M` (requires elevated prompt).
- PATH modifications update the registry `Environment` key.

### Linux / macOS
- Appends `export VAR=value` to the user's shell profile.
- Detects active shell: `~/.bashrc` (Bash), `~/.zshrc` (Zsh), `~/.profile` (fallback).
- Deduplicates entries before writing.
- PATH modifications append to the `export PATH=...` line or add a new one.

---

## Commands

### `gitmap-v28 env` (alias: `ev`)

Manage environment variables across platforms.

```bash
gitmap-v28 env set NAME value
gitmap-v28 env set NAME "value with spaces"
gitmap-v28 env path add /usr/local/go/bin
gitmap-v28 env path add "C:\Program Files\Go\bin"
gitmap-v28 env path list
gitmap-v28 env path remove /old/path
gitmap-v28 env list
gitmap-v28 env get NAME
gitmap-v28 env delete NAME
```

---

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `set`      | Set an environment variable (persists across sessions) |
| `get`      | Read current value of a variable |
| `delete`   | Remove an environment variable |
| `list`     | List all gitmap-managed variables |
| `path add` | Add a directory to the system PATH |
| `path remove` | Remove a directory from the system PATH |
| `path list` | List all PATH entries |

---

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--system` | `false` | Target system-level variables (Windows: requires admin) |
| `--shell <name>` | auto-detect | Force target shell profile (bash, zsh, fish) |
| `--verbose` | `false` | Show file modifications and registry calls |
| `--dry-run` | `false` | Preview changes without applying |

---

## Tracking File

gitmap-v28 tracks which variables it has set in `.gitmap/env-registry.json`:

```json
{
  "variables": [
    {
      "name": "GOPATH",
      "value": "C:\\Users\\alim\\go",
      "scope": "user",
      "set_at": "2026-04-04T10:00:00Z"
    }
  ],
  "paths": [
    {
      "entry": "C:\\Program Files\\Go\\bin",
      "added_at": "2026-04-04T10:01:00Z"
    }
  ]
}
```

---

## Cross-Platform Details

### Windows — User Level (default)
```
setx GOPATH "C:\Users\alim\go"
```
- Modifies `HKCU\Environment` registry key.
- Effective on next terminal session.

### Windows — System Level (`--system`)
```
setx /M GOPATH "C:\Users\alim\go"
```
- Modifies `HKLM\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`.
- Requires elevated (admin) prompt.

### Linux — Bash
```bash
# Appended to ~/.bashrc
export GOPATH="/home/alim/go"
export PATH="$PATH:/usr/local/go/bin"
```

### Linux / macOS — Zsh
```bash
# Appended to ~/.zshrc
export GOPATH="/home/alim/go"
export PATH="$PATH:/usr/local/go/bin"
```

### macOS — Fallback
```bash
# Appended to ~/.profile if no .bashrc or .zshrc exists
export GOPATH="/Users/alim/go"
```

---

## PATH Management

### Adding a Path
1. Read current PATH from OS/profile.
2. Check for duplicates (case-insensitive on Windows).
3. If not present, append to PATH.
4. Write to appropriate store (registry / profile file).
5. Record in env-registry.json.

### Removing a Path
1. Read current PATH.
2. Filter out the target entry.
3. Write updated PATH back.
4. Remove from env-registry.json.

---

## Error Handling

| Scenario | Behavior |
|----------|----------|
| `--system` without admin (Windows) | Exit with "Run as Administrator" message |
| Shell profile not found | Create the file with a header comment |
| Duplicate PATH entry | Skip with info message |
| Variable name contains spaces | Reject with error |
| Empty value | Allow (sets empty variable) |

---

## File Layout

| File | Purpose |
|------|---------|
| `constants/constants_env.go` | Command names, defaults, messages |
| `cmd/env.go` | Subcommand routing and flag parsing |
| `cmd/envset.go` | Set/get/delete variable logic |
| `cmd/envpath.go` | PATH add/remove/list logic |
| `cmd/envplatform_windows.go` | Windows registry and setx calls |
| `cmd/envplatform_unix.go` | Unix profile file operations |
| `model/envregistry.go` | Registry struct and JSON serialization |

---

## Constraints

- Variable names must be alphanumeric + underscore (validated).
- PATH entries are validated as existing directories (warn if not).
- Profile file modifications are atomic (write to temp, rename).
- All files under 200 lines, all functions 8-15 lines.
- No `!` or `!=` in conditions; use positive logic.

---

## Examples

```bash
# Set a variable
gitmap-v28 env set GOPATH "C:\Users\alim\go"

# Set system-level variable (Windows admin)
gitmap-v28 env set JAVA_HOME "C:\Java\jdk-21" --system

# Add to PATH
gitmap-v28 env path add "C:\Program Files\Go\bin"
gitmap-v28 env path add /usr/local/go/bin

# List managed variables
gitmap-v28 env list

# Preview changes
gitmap-v28 env set NODE_ENV production --dry-run

# Remove a PATH entry
gitmap-v28 env path remove /old/path
```

---

## See Also

- setup — shell profile and completion installation
- doctor — verify environment configuration
- install — automated tool installation
