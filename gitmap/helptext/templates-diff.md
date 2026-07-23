# templates diff

Show what `gitmap add ignore <lang>` or `gitmap add attributes <lang>`
would change in the current repo, **without writing to disk**.

## Usage

```
gitmap templates diff --lang <name> [--kind ignore|attributes] [--cwd <path>]
gitmap tpl td --lang <name>                                  # short form
```

## Flags

| Flag | Required | Default | Purpose |
|------|----------|---------|---------|
| `--lang <name>` | yes | — | Language to diff (e.g. `go`, `node`, `python`). |
| `--kind <ignore\|attributes>` | no | both | Restrict to one kind. |
| `--cwd <path>` | no | `.` | Repo directory. The target file is `<cwd>/.gitignore` or `<cwd>/.gitattributes`. |

## Exit codes

Mirrors standard `diff(1)` so it slots straight into shell pipelines:

| Code | Meaning |
|------|---------|
| 0    | No changes — the on-disk gitmap block already matches the template. |
| 1    | Differences found — the block is missing, or its body differs from the template. |
| 2    | Error — invalid flag value, unknown language, or I/O failure. |

## Output

The diff is **block-scoped**: only the gitmap-managed marker block
(`# >>> gitmap:<kind>/<lang> >>>` … `# <<< gitmap:<kind>/<lang> <<<`)
participates. Hand edits OUTSIDE the block are invisible — same
contract as `add` itself.

Each hunk starts with a banner line (`@@ gitmap:<kind>/<lang> @@`),
followed by `-` lines for content currently on disk and `+` lines for
content the template would write.

A no-change case prints a single human-readable summary line per
kind:

    no changes for ignore/go in /repo/.gitignore

## Examples

### Confirm `add ignore go` would be a no-op

    gitmap templates diff --lang go --kind ignore
    # → "no changes for ignore/go in /repo/.gitignore" (exit 0)

### See what would change for a brand-new repo

    gitmap templates diff --lang node
    # → @@ gitmap:ignore/node @@
    #   +node_modules/
    #   +*.log
    #   @@ gitmap:attributes/node @@
    #   +*.js  text eol=lf
    #   ...                                                (exit 1)

### Use as a pre-commit guard

    if ! gitmap tpl td --lang go --kind ignore > /dev/null; then
      echo "gitignore drift detected — run 'gitmap add ignore go'"
      exit 1
    fi

## See also

- [add ignore](add-ignore.md)
- [add attributes](add-attributes.md)
- [templates init](templates.md)

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter templates-diff
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
