# gitmap make-public

Make the current repository **public** on GitHub or GitLab.

```
gitmap make-public [<repo-or-url>] [<count>] [--yes] [--dry-run] [--verbose]
```

## Bulk form (v5.61.0+)

`gitmap make-public <count>` flips the **N most recent versions** of
the *current* repo (vN, vN-1, …, vN-count+1).

`gitmap make-public <repo-or-url> <count>` does the same against a
different base — `<repo-or-url>` accepts a full URL, an SSH
shorthand, or a bare base name (`foo` or `foo-v40`). The provider
and owner are inherited from the current repo's origin.

```
gitmap make-public 3                                   # flip current + 2 prior
gitmap make-public macro-ahk-v40 5                     # flip v40..v36
gitmap make-public https://github.com/me/foo-v9 2 --yes
```

A single `--yes` covers the whole batch. Individual failures are
non-fatal; the command exits with the worst code observed.


## What it does

1. Detects the provider (GitHub or GitLab) and the `owner/repo` slug
   from `git remote get-url origin`.
2. Verifies that the matching CLI (`gh` or `glab`) is on `PATH` and
   already authenticated — gitmap does not store any tokens.
3. Reads the current visibility. If the repo is **already public**,
   exits 0 with no changes.
4. Prompts for explicit `yes` confirmation before going private →
   public. Pass `--yes` (or `-y`) to skip the prompt for scripts/CI.
5. Runs `gh repo edit <slug> --visibility public
   --accept-visibility-change-consequences` (GitHub) or `glab repo
   edit <slug> --visibility public` (GitLab).
6. Re-reads visibility to verify the change actually took effect.

## Flags

| Flag | Behavior |
|------|----------|
| `--yes`, `-y` | Skip the "are you sure?" private→public confirmation. |
| `--dry-run` | Print the provider command that would run; do not invoke it. |
| `--verbose` | Echo every shell command to stderr before running it. |

## Examples

```
# Interactive (will prompt for confirmation)
gitmap make-public

# Non-interactive (CI / scripts)
gitmap make-public --yes

# Preview without touching the API
gitmap make-public --dry-run

# Debug auth or argv issues
gitmap make-public --yes --verbose
```

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | Success (or already public) |
| 2 | Not inside a git repository |
| 3 | No `origin` remote configured |
| 4 | Unsupported provider host, or unparseable owner/repo |
| 5 | Provider CLI missing, not authenticated, or apply failed |
| 6 | Bad flag |
| 7 | Confirmation required (re-run with `--yes`) |
| 8 | Verification failed (visibility did not change) |

## See also

- `gitmap make-private` — the opposite direction.
- `gh auth login` / `glab auth login` — authenticate the provider CLI.

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter make-public
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
