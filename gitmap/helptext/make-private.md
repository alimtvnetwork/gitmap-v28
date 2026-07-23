# gitmap make-private

Make the current repository **private** on GitHub or GitLab.

```
gitmap make-private [<repo-or-url>] [<count>] [--dry-run] [--verbose]
```

## Bulk form (v5.61.0+)

`gitmap make-private <count>` flips the **N most recent versions** of
the *current* repo. `gitmap make-private <repo-or-url> <count>`
targets a different base. See `gitmap make-public --help` for the
full bulk semantics — they are identical except no confirmation
prompt is ever shown (private is the safe direction).


## What it does

1. Detects the provider (GitHub or GitLab) and the `owner/repo` slug
   from `git remote get-url origin`.
2. Verifies that the matching CLI (`gh` or `glab`) is on `PATH` and
   already authenticated — gitmap does not store any tokens.
3. Reads the current visibility. If the repo is **already private**,
   exits 0 with no changes.
4. Runs `gh repo edit <slug> --visibility private
   --accept-visibility-change-consequences` (GitHub) or `glab repo
   edit <slug> --visibility private` (GitLab).
5. Re-reads visibility to verify the change actually took effect.

> No confirmation prompt is shown for `make-private` — hiding a
> public repo is the safe direction (and reversible). The
> confirmation lives on `make-public`, where exposure happens.

## Flags

| Flag | Behavior |
|------|----------|
| `--dry-run` | Print the provider command that would run; do not invoke it. |
| `--verbose` | Echo every shell command to stderr before running it. |

> `--yes` / `-y` is accepted but has no effect (kept for symmetry
> with `make-public` so the same script works for both).

## Examples

```
# Standard
gitmap make-private

# Preview without touching the API
gitmap make-private --dry-run

# Debug auth or argv issues
gitmap make-private --verbose
```

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | Success (or already private) |
| 2 | Not inside a git repository |
| 3 | No `origin` remote configured |
| 4 | Unsupported provider host, or unparseable owner/repo |
| 5 | Provider CLI missing, not authenticated, or apply failed |
| 6 | Bad flag |
| 8 | Verification failed (visibility did not change) |

## See also

- `gitmap make-public` — the opposite direction (with confirmation).
- `gh auth login` / `glab auth login` — authenticate the provider CLI.

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter make-private
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
