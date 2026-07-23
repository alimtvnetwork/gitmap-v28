# audit-legacy

Scan the workspace for forbidden legacy strings and exit non-zero on any
hit. Built as a regression guard so future remixes / rename commits
(e.g. `gitmap-v28` → `gitmap-v28`) can't silently leave stale references <!-- gitmap-legacy-ref-allow -->
behind.

## Synopsis

```
gitmap audit-legacy [--patterns <csv>] [--path <dir>] [--json] [--report[=<file>]] [--diffs]
gitmap audit         [--patterns <csv>] [--path <dir>] [--json] [--report[=<file>]] [--diffs]
gitmap al            [--patterns <csv>] [--path <dir>] [--json] [--report[=<file>]] [--diffs]
```

## Defaults

- `--patterns` defaults to `gitmap-v[567]\b` — catches every old
  versioned-repo reference.
- `--path` defaults to the current working directory.
- `--report` without a value writes to `.gitmap/audit-legacy-report.md`.
  Pass `--report=path/to/file.md` to override.
- `--diffs` (requires `--report`) writes one unified-diff artifact per
  offending file under `<reportDir>/diffs/<flat-path>.diff` and links
  each one from the report's "Counts by file" + "Per-file diffs"
  sections. Apply with `patch -p0 < <file>.diff` from the repo root.
- Skips `.git`, `node_modules`, `dist`, `build`, `bin`, `.next`,
  `.gitmap`, `vendor`, `coverage`, plus binary file extensions
  (images, archives, executables, fonts, sqlite).

## Exit codes

| Code | Meaning                                  |
|------|------------------------------------------|
| 0    | No matches — workspace is clean          |
| 1    | One or more matches found (regression!)  |
| 2    | Bad flags / regex / walk error           |

## Examples

```
# Default scan from repo root
gitmap audit-legacy

# Add custom patterns (comma-separated regexes)
gitmap audit-legacy --patterns "gitmap-v[567]\b,old-org-name"

# Machine-readable output for CI
gitmap al --json > legacy-report.json

# Scope the scan to a subtree
gitmap audit-legacy --path ./src

# Full audit + per-file unified diffs as downloadable artifacts
gitmap audit --report --diffs
```

## CI usage

Add to your pre-merge or release workflow:

```yaml
- name: Guard against legacy refs
  run: gitmap audit-legacy --json
```

A non-zero exit fails the job and the JSON report names every
offending file + line.

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter audit-legacy
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
