# gitmap conflicts

Inspect and list all currently conflicted / unmerged files with tree view and quick-action commands.

## Alias

cf, conflict

## Usage

```bash
gitmap conflicts [repo]
```

## Description

Scans the current repository, a specified repository by name/path, or all tracked repositories in the workspace for active merge, rebase, or cherry-pick conflicts.
Outputs a clean ASCII tree view of all conflicted files along with instant one-liner quick actions to resolve with IDE 3-way merge editor, `--theirs`, `--ours`, or `--abort`.

When run without arguments:
- If inside a git repository, inspects the current repository.
- If run at workspace root, scans all registered repositories in the database and lists any repos with pending conflicts.

## Examples

### Example 1: Check conflicts in current repository

```bash
gitmap conflicts
```

**Output:**

```
📁 Repository: my-service (D:\work\my-service)
   Conflicted Files (2):
   ├── src/auth/login.ts
   └── package.json

⚡ Quick Actions:
   gitmap resolve my-service --ide      Launch VS Code / IDE 3-Way Merge Editor
   gitmap resolve my-service --theirs   Accept incoming remote changes
   gitmap resolve my-service --ours     Accept local current changes
   gitmap resolve my-service --abort    Abort merge or rebase
   gitmap cpf "resolve conflicts"       Stage, commit, and push once resolved
```

### Example 2: Check conflicts in a specific repository

```bash
gitmap conflicts auth-backend
```

**Output:**

```
📁 Repository: auth-backend (D:\work\auth-backend)
   Conflicted Files (1):
   └── internal/token/jwt.go

⚡ Quick Actions:
   gitmap resolve auth-backend --ide      Launch VS Code / IDE 3-Way Merge Editor
   gitmap resolve auth-backend --theirs   Accept incoming remote changes
   gitmap resolve auth-backend --ours     Accept local current changes
   gitmap resolve auth-backend --abort    Abort merge or rebase
   gitmap cpf "resolve conflicts"         Stage, commit, and push once resolved
```

### Example 3: When no conflicts exist

```bash
gitmap conflicts
```

**Output:**

```
✓ No merge conflicts detected in current repository.
```

## See Also

- [resolve](resolve.md) — Resolve merge conflicts via IDE, `--theirs`, `--ours`, or `--abort`
- [reconcile](reconcile.md) — Reconcile dirty repositories after pull
- [fix](fix.md) — Remediation walkthrough for sync failures
- [pull-commit-push](pull-commit-push.md) — Commit and push once conflicts are resolved
