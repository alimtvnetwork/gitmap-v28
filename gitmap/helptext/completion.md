# gitmap completion

Generate or install shell tab-completion scripts for gitmap commands and repo names.

## Alias

cmp

## Usage

    gitmap completion <powershell|bash|zsh>

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --list-repos | — | Print repo slugs, one per line (for script use) |
| --list-groups | — | Print group names, one per line (for script use) |
| --list-commands | — | Print all command names, one per line (for script use) |
| --list-aliases | — | Print alias names, one per line (for script use) |
| --list-zip-groups | — | Print zip group names, one per line (for script use) |

## Prerequisites

- Run `gitmap scan` first for repo name completion
- Shell profile must be writable for auto-install via `gitmap setup`

## Examples

### Example 1: Generate PowerShell completion script

    gitmap completion powershell

**Output:**

    Register-ArgumentCompleter -CommandName gitmap -ScriptBlock {
        param($commandName, $wordToComplete, $cursorPosition)
        $commands = @('scan','clone','pull','status','exec','release',
                      'group','list','cd','watch','alias','bookmark',...)
        $commands | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
    }

### Example 2: Generate Bash completion script

    gitmap completion bash

**Output:**

    _gitmap_completions() {
        local cur="${COMP_WORDS[COMP_CWORD]}"
        local commands="scan clone pull status exec release group list cd watch"
        COMPREPLY=( $(compgen -W "$commands" -- "$cur") )
    }
    complete -F _gitmap_completions gitmap

### Example 3: List repo slugs for scripting

    gitmap completion --list-repos

**Output:**

    my-api
    web-app
    billing-svc
    auth-gateway
    shared-lib
    notification-svc

### Example 4: List all available commands

    gitmap completion --list-commands

**Output:**

    scan
    clone
    pull
    rescan
    status
    exec
    release
    release-branch
    release-pending
    ...

## See Also

- [setup](setup.md) — Auto-installs completions during setup
- [cd](cd.md) — Navigate to repos using tab-completed slugs
- [group](group.md) — Group names are also tab-completed

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter completion
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
