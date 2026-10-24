# commit-push-refactor

Stages all changes, formats commit message with uppercase `REFACTOR(MODULE): <summary>` prefix, and pushes to remote tracking branch.

## Synopsis

```bash
gitmap commit-push-refactor "(module): description" [flags]
gitmap cfr "(module): description" [flags]
```

## Behavior

1. **Auto-Stage:** Automatically stages all tracked and untracked changes (`git add -A`).
2. **Uppercase Conventional Format:** Formats as `REFACTOR(MODULE): description` or `REFACTOR: description`.
3. **No Third Brackets:** Replaces any square brackets with standard round parentheses.
4. **Duplicate Prefix Trimming:** Safely strips duplicate `refactor:` or `REFACTOR:` prefixes.
5. **Direct Push:** Executes `git push` using the folder-specific bound SSH key if configured.

## Aliases

- `gitmap cfr "(module): description"`

## Examples

### Refactor with Module Scope

```bash
gitmap cfr "(auth): decouple token verification from handler"
```

### Refactor with Bare Module

```bash
gitmap cfr "cache: hoist compiled regular expressions"
```

### Dry Run Preview

```bash
gitmap cfr "(db): optimize query index" --dry-run
```

## See Also

- `gitmap hd commit-format` — Uppercase Conventional Commit formatting guide.
- `gitmap hd cpf` — Feature commit and push (`FEAT(MODULE): ...`).
- `gitmap hd cbf` — Bug fix commit and push (`FIX(MODULE): ...`).
