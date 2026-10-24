# gitmap cbf

Stages all changes, commits with an uppercase `FIX(MODULE): ` prefix, and pushes to remote. First-class alias for `commit-push-bug`.

## Synopsis

```bash
gitmap cbf "<message>" [flags]
gitmap commit-push-bug "<message>" [flags]
gitmap cpb "<message>" [flags]
```

## Behavior

1. **Auto-Stage:** Automatically stages all tracked and untracked changes (`git add -A`).
2. **Semantic Prefix:** Injects uppercase `FIX:` or `FIX(MODULE):` prefix into the commit message.
3. **Module Scope Support:** Formats module scopes cleanly into uppercase parentheses `(MODULE)`. Square brackets are automatically converted.
4. **Duplicate Prefix Trimming:** Safely strips duplicate `Bug: `, `fix: `, or `FIX:` prefixes if already typed.
5. **Immediate Push:** Executes `git push` with bound SSH keys and error handling.

## Aliases

- `gitmap cbf "(module): description"`
- `gitmap cpb "(module): description"`

## Examples

### Bug Fix with Module Scope

```bash
gitmap cbf "(db): fix sqlite connection pool deadlock"
```

### Bug Fix with Bare Module

```bash
gitmap cbf "ssh: handle missing private key"
```

### Dry Run Bug Fix

```bash
gitmap cbf "(scanner): handle nil pointer when DB is empty" --dry-run
```
