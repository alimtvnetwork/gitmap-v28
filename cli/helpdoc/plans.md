# Plans

List and summarize the completed and pending architectural plans of the
current repository.

## Usage

    gitmap plans [N] [flags]
    gitmap plan [N] [flags]

## Description

Discovers plan files under `.ai-memory/plans/` and displays them in two
sections: completed plans and pending (active) plans. Each section is
truncated to the display limit (10 by default); a positional `N` or
`-n, --limit <N>` raises the limit, and `--all` shows everything.

## Flags

| Flag                | Description                                   |
|---------------------|-----------------------------------------------|
| `-n, --limit <N>`   | Limit number of plans shown per section       |
| `-a, --all`         | Show all completed and pending plans          |
| `-p, --pending`     | Show only pending / active plans              |
| `-c, --completed`   | Show only completed plans                     |
| `-j, --json`        | Emit machine-readable JSON to stdout          |
| `--yaml, --yml`     | Emit machine-readable YAML to stdout          |
| `-h, --help`        | Print usage help                              |

## Examples

### List plans with the default limit

    $ gitmap plans

### Show the first 20 plans per section

    $ gitmap plans 20

### Show only pending plans as JSON

    $ gitmap plans --pending --json

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter plans
```

The JSON schema is published at `02-spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
