# Shell Completion — Generic CLI Spec

> **Related specs:**
> - [03-subcommand-architecture.md](03-subcommand-architecture.md) — subcommand registry that completion draws from
> - [09-help-system.md](09-help-system.md) — help metadata used by completion generators
> - [04-flag-parsing.md](04-flag-parsing.md) — flags that completion must enumerate

## Overview

Shell tab-completion provides autocomplete for subcommand names, repo
slugs, group names, and per-command flags across PowerShell, Bash, and
Zsh. Completions use a hidden `completion` subcommand with `--list-*`
flags that query the database and print one value per line.

---

## Completion Subcommand

```
toolname completion <powershell|bash|zsh>
toolname completion --list-repos
toolname completion --list-groups
toolname completion --list-commands
```

### List Flag Behaviour

| Flag | Output |
|------|--------|
| `--list-repos` | One repo slug per line, no headers |
| `--list-groups` | One group name per line, no headers |
| `--list-commands` | One command name per line (includes aliases) |

All list outputs are plain text with no color or decoration.
They open the DB, query, print, and exit immediately.

---

## Completed Contexts

| User Typing | Values Offered |
|-------------|---------------|
| `toolname <tab>` | All subcommand names and aliases |
| `toolname cd <tab>` | Repo slugs + `repos`, `set-default`, `clear-default` |
| `toolname pull <tab>` | Repo slugs |
| `toolname exec --group <tab>` | Group names |
| `toolname group <tab>` | `create`, `add`, `remove`, `list`, `show`, `delete` |

---

## Shell Scripts

### PowerShell

Uses `Register-ArgumentCompleter`. Inspects `$commandAst.CommandElements`
to determine the active subcommand and decide what to complete.

### Bash

Uses `complete -F`. The function checks `COMP_WORDS` and `COMP_CWORD`
to route to `compgen -W` with the appropriate list source.

### Zsh

Uses `#compdef`. Routes via `words[2]` and calls `_describe` with
arrays populated from list flags.

---

## Setup Integration

The `setup` command auto-installs completions:

1. Detects current shell.
2. Generates the completion script via `Generate(shell)`.
3. Writes script to a well-known path.
4. Appends a source line to the user's shell profile.

Installation is idempotent — the source line check prevents duplicates.

| Shell | Script Path | Profile |
|-------|------------|---------|
| PowerShell | `$APPDATA/toolname/completions.ps1` | `$PROFILE` |
| Bash | `~/.local/share/toolname/completions.bash` | `~/.bashrc` |
| Zsh | `~/.local/share/toolname/completions.zsh` | `~/.zshrc` |

---

## File Layout

| File | Purpose |
|------|---------|
| `constants/constants_completion.go` | Shell names, list flags, messages |
| `cmd/completion.go` | Subcommand handler + list printers |
| `completion/completion.go` | `Generate()` + `AllCommands()` |
| `completion/powershell.go` | PowerShell script generator |
| `completion/bash.go` | Bash script generator |
| `completion/zsh.go` | Zsh script generator |
| `completion/install.go` | Profile detection + source-line writer |

---

## Constraints

- List outputs are one value per line, no color, no headers.
- Scripts call `toolname completion --list-*` at tab-time for dynamic data.
- PowerShell scripts use ASCII only (no em-dashes, no Unicode).
- All files under 200 lines, all functions 8–15 lines.
- Setup never duplicates source lines in profiles.
