# gitmap var

Persistent variable engine for global and command-scoped configuration,
dynamic path expansion (`$VAR`, `${VAR}`), and OS environment export.

## Aliases

variable

## Usage

    gitmap var set <key> <value> [--scope <scope>]
    gitmap var get <key> [--scope <scope>]
    gitmap var ls [--scope <scope>]
    gitmap var rm <key> [--scope <scope>]
    gitmap var export

## Scopes

- `global` (default): Available across all gitmap commands.
- `<scope>` (e.g. `agy`, `cluster`): Scoped to a specific command family.

## Expansion

Variables are expanded in paths, prompts, and CLI arguments:
`$WORKSPACE/prompts` expands to the stored value.

## Examples

```bash
# Set global variable
gitmap var set MY_REPO /path/to/myrepo

# Set scoped variable for AGY
gitmap var set PROMPT_DIR /path/to/prompts --scope agy

# List variables
gitmap var ls
```

