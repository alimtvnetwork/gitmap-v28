# gitmap make-all-public-except-latest

Same as `gitmap make-all-public`, but **flips the newest `-vN` sibling
of every base group to the OPPOSITE visibility** (private). Useful when
you want every prior version of a project public but keep only the
current/latest cut private.

```
gitmap make-all-public-except-latest <owner-or-url> <patterns> \
    [-Y|--yes] [--verbose] [--parallel=N] [--cache-ttl=SECONDS]
```

## What `--except-latest` does

For every base group with 2+ versioned siblings:

- All non-latest versions       → flipped to **public** (the requested target)
- Highest `-vN` per base group  → flipped to **private** (the inverted target)
- Repos with no `-vN` suffix    → flipped to **public** (passed through)
- Bases with only one versioned entry → flipped to **public** (nothing to invert)

| Repo name        | Base group  | Version | Final visibility |
|------------------|-------------|---------|------------------|
| `myapp-v25`      | `myapp`     | 25      | public           |
| `myapp-v26`      | `myapp`     | 26 ✅   | **private** (inverted) |
| `myapp-v26-rc1`  | (no match)  | —       | public           |
| `tooling`        | (no match)  | —       | public           |

## Flags

| Flag | Behavior |
|------|----------|
| `-Y`, `--yes` | Skip the interactive confirmation prompt. |
| `--verbose` | Echo every shell command to stderr before running it. |
| `--parallel=N` | Apply N repos concurrently (default 8, max 32). |
| `--cache-ttl=N` | Override the owner repo-list cache TTL (seconds; 0 disables). |

## Examples

```
# All myapp-vN public, but myapp-v<latest> goes private
gitmap make-all-public-except-latest alice "myapp-v*" -Y

# Same for two base groups at once, parallel 16
gitmap make-all-public-except-latest alice "demo-v*,api-v*" --parallel=16 -Y

# Uppercase shorthand
gitmap MAPUBXL alice "demo-v*"
```

## See also

- `gitmap make-all-public` — flip every match, no inversion.
- `gitmap make-all-private-except-latest` — mirror command (all private, latest public).
- `gitmap MAPUBXL` — uppercase shorthand for this command.
