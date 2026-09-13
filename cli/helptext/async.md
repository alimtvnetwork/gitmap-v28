# async

Execute CLI commands asynchronously in the background with detached process tracking.

## Aliases

`bg`, `spawn`

## Usage

    gitmap async <command> [args...]

## Description

The `async` command launches long-running tasks, scripts, or build processes in the background without blocking the interactive terminal. Process status, exit codes, and output logs are persisted in repository-scoped runtime records.

## Examples

```bash
# Launch a background build command
gitmap async npm run build

# Spawn an asynchronous pipeline runner
gitmap async gitmap pipeline run all
```
