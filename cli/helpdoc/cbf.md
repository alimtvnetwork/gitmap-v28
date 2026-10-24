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
3. **Uppercase Module Support:** Recognizes module token and formats cleanly as `Bug: MODULE - <item>`.
4. **Duplicate Prefix Trimming:** Safely strips duplicate `Bug: ` or `fix: ` prefixes if already typed by the user or AI agent.
5. **Immediate Push:** Executes `git push` with bound SSH keys and error handling.

## Examples

```bash
# Standard uppercase module syntax (recommended)
gitmap cbf "DB - fix sqlite connection pool deadlock"

# Conventional type syntax preservation
gitmap cbf "fix(ssh): handle missing private key"

# Legacy hyphen conversion
gitmap cbf "cache - fix ttl expiration check"

# Scope commit to claimed task files only
gitmap cbf "AUTH - fix token expiration" --task 317
```
