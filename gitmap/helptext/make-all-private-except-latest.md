# gitmap make-all-private-except-latest

Same as `gitmap make-all-private`, but **flips the newest `-vN` sibling
of every base group to the OPPOSITE visibility** (public). Useful when
you want every prior version of a project hidden but keep the
current/latest cut visible.

```
gitmap make-all-private-except-latest <owner-or-url> <patterns> \
    [-Y|--yes] [--verbose] [--parallel=N] [--cache-ttl=SECONDS]
```

## What `--except-latest` does

For every base group with 2+ versioned siblings:

- All non-latest versions       → flipped to **private** (the requested target)
- Highest `-vN` per base group  → flipped to **public** (the inverted target)
- Repos with no `-vN` suffix    → flipped to **private** (passed through)
- Bases with only one versioned entry → flipped to **private** (nothing to invert)

| Repo name        | Base group  | Version | Final visibility |
|------------------|-------------|---------|------------------|
| `myapp-v25`      | `myapp`     | 25      | private          |
| `myapp-v26`      | `myapp`     | 26 ✅   | **public** (inverted) |
| `myapp-v26-rc1`  | (no match)  | —       | private          |
| `tooling`        | (no match)  | —       | private          |

## Flags

| Flag | Behavior |
|------|----------|
| `-Y`, `--yes` | Skip the interactive confirmation prompt. |
| `--verbose` | Echo every shell command to stderr before running it. |
| `--parallel=N` | Apply N repos concurrently (default 8, max 32). |
| `--cache-ttl=N` | Override the owner repo-list cache TTL (seconds; 0 disables). |

## Examples

```
# Archive every prior version privately, keep latest public
gitmap make-all-private-except-latest alice "myapp-v*" -Y

# Two base groups, parallel 16
gitmap make-all-private-except-latest alice "demo-v*,api-v*" --parallel=16 -Y

# Uppercase shorthand
gitmap MAPRIXL alice "demo-v*"
```

## See also

- `gitmap make-all-private` — flip every match, no inversion.
- `gitmap make-all-public-except-latest` — mirror command (all public, latest private).
- `gitmap MAPRIXL` — uppercase shorthand for this command.
