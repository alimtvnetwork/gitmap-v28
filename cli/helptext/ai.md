# gitmap ai

Native AI scripts automation catalog, streaming execution, and repository autofix engine.

## Usage

```bash
gitmap ai [subcommand] [flags]
gitmap scripts [subcommand] [flags]
gitmap ai <token> [args...]
```

## Description

`gitmap ai` provides a unified CLI interface for managing, executing, and monitoring autonomous AI Python scripts located in `03-ai-scripts/`. It discovers available system Python runtimes, tracks script metadata and categories, streams unbuffered script output in real time, and exposes automated repository autofix suites for codebase hygiene.

## Subcommands

- `list` (aliases: `ls`, `catalog`): List registered AI automation scripts with status and categories.
- `run <token> [args...]` (alias: `exec`): Run a registered script by numeric token, slug, or alias with live streaming output.
- `fix [target] [args...]`: Run standardized repository autofix targets (`guidelines`, `newlines`, `paths`, `naming`, `encoding`, `spelling`, `all`).

## Flags

### List Flags (`gitmap ai list`)

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--category` | `-c` | `""` | Filter scripts by category (`audit`, `guidelines`, `ci`, `fix`, `spec`) |
| `--search` | `-s` | `""` | Search scripts by token, slug, or description substring |
| `--fix` | `-f` | `false` | Show only scripts supporting automated fix mode |
| `--json` | | `false` | Output script catalog as structured JSON |
| `--verbose` | `-v` | `false` | Display full script metadata, aliases, and file paths |

## Examples

```bash
# 1. List registered AI automation scripts in table format
gitmap ai list
gitmap scripts ls

# 2. Filter scripts by category or search keyword
gitmap ai list --category guidelines
gitmap ai list --search "format"

# 3. Output only scripts with fix mode in JSON
gitmap ai list --fix --json

# 4. Execute a script by numeric token, slug, or alias
gitmap ai run 01
gitmap ai run 06-cicd-local-runner.py
gitmap ai 01-ai-instruction-writer

# 5. Run standardized repository autofix targets
gitmap ai fix all
gitmap ai fix guidelines
gitmap ai fix newlines

# 6. Execute via scripts alias
gitmap scripts run 14
```

See also: `gitmap os ai-clean`, `gitmap pipeline-ai`, `gitmap agy`
