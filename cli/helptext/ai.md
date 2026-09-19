# gitmap ai

Native AI scripts automation catalog, streaming execution, scaffolding, and repository autofix engine.

## Usage

```bash
gitmap ai [subcommand] [flags]
gitmap scripts [subcommand] [flags]
gitmap ai <token> [args...]
```

## Description

`gitmap ai` provides a unified CLI interface for managing, executing, creating, and monitoring autonomous AI Python scripts located in `03-ai-scripts/`. It discovers available system Python runtimes, tracks script metadata and categories, dynamically scans new scripts, scaffolds compliant script archetypes, and exposes automated repository autofix suites for codebase hygiene.

## Subcommands

- `list` (aliases: `ls`, `catalog`): List registered and dynamic AI automation scripts.
- `run <token> [args...]` (alias: `exec`): Run a script by token, slug, or filename with live streaming.
- `create <name> [flags]` (aliases: `new`, `scaffold`, `gen`): Scaffold a new AI script skeleton.
- `fix [target] [args...]`: Run standardized repository autofix targets (`all`, `guidelines`, `newlines`).

## Flags

### List Flags (`gitmap ai list`)

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--category` | `-c` | `""` | Filter scripts by category (`core`, `formatting`, `guideline`, `audit`) |
| `--search` | `-s` | `""` | Search scripts by token, slug, or description substring |
| `--fix` | `-f` | `false` | Show only scripts supporting automated fix mode |
| `--json` | | `false` | Output script catalog as structured JSON |
| `--verbose` | `-v` | `false` | Display full script metadata, aliases, and file paths |

### Create Flags (`gitmap ai create`)

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--type` | `-t` | `"linter"` | Archetype (`linter`, `fixer`, `auditor`, `generator`, `util`) |
| `--desc` | `-d` | `""` | Description of script intent and functionality |
| `--parallel` | `-p` | `false` | Include multi-worker `ThreadPoolExecutor` scaffold |
| `--fix` | | `false` | Include `--fix` modification flags and handler |
| `--dry-run` | | `false` | Preview generated script without writing to disk |
| `--force` | `-f` | `false` | Overwrite existing script file if already present |

## Examples

```bash
# 1. List registered AI automation scripts in table format
gitmap ai list
gitmap scripts ls

# 2. Filter scripts by category or search keyword
gitmap ai list --category guidelines
gitmap ai list --search "format"

# 3. Scaffold a new linter with parallel workers
gitmap ai create import-order-checker --type linter --parallel

# 4. Preview script scaffolding in dry-run mode
gitmap ai new schema-validator --type auditor --dry-run

# 5. Scaffold a custom autofixer script
gitmap ai scaffold comment-stripper --type fixer --fix

# 6. Execute a script by numeric token, slug, or alias
gitmap ai run 01
gitmap ai run 06-cicd-local-runner.py
gitmap ai 01-ai-instruction-writer

# 7. Run standardized repository autofix targets
gitmap ai fix all
gitmap ai fix guidelines
gitmap ai fix newlines
```

See also: `gitmap os ai-clean`, `gitmap pipeline-ai`, `gitmap agy`
