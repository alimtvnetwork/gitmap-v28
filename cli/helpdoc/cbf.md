# gitmap cbf

Stages all changes, commits with a "Bug: " prefix, and pushes to remote. First-class alias for `commit-push-bug`.

## Synopsis

```
gitmap cbf "<message>" [flags]
gitmap commit-push-bug "<message>" [flags]
gitmap cpb "<message>" [flags]
```

## Behavior

1. **Auto-Stage:** Automatically stages all tracked and untracked changes (`git add -A`).
2. **Semantic Prefix:** Injects `Bug: ` prefix into the commit message.
3. **Bracketed Module Support:** Recognizes `[module]: <item>` syntax and formats cleanly as `Bug: [module]: <item>`.
4. **Duplicate Prefix Trimming:** Safely strips duplicate `Bug: ` or `fix: ` prefixes if already typed by the user or AI agent.
5. **Immediate Push:** Executes `git push` with bound SSH keys and error handling.

## Examples

```bash
# Standard bracketed module syntax (recommended)
gitmap cbf "[db]: fix sqlite connection pool deadlock"

# Conventional type syntax preservation
gitmap cbf "fix([ssh]): handle missing private key"

# Legacy hyphen conversion (automatically converted to [module]: item)
gitmap cbf "cache - fix ttl expiration check"

# Scope commit to claimed task files only
gitmap cbf "[auth]: fix token expiration" --task 317
```
