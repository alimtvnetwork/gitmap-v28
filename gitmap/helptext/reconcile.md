# gitmap reconcile

Reconcile failed or dirty repositories repo-by-repo after pull conflicts.

## Alias

recon

## Usage

```bash
gitmap reconcile [repo-name] [stash|wip|discard]
gitmap reconcile --all [stash|wip|discard]
```

## Description

Reconciles repositories that encountered merge conflicts, untracked file
collisions, or uncommitted working tree changes during `gitmap pull`.
Allows targeting specific repositories by name or resolving all pending
repositories in batch.

## Options

| Action | Description |
|--------|-------------|
| `stash` | Stash local modifications and untracked files (`-u`), pull remote, and restore stash |
| `wip` | Stage all changes, commit as temporary WIP, and rebase pull |
| `discard` | Reset working tree to HEAD, clean untracked files (`-fd`), and pull remote |

## Examples

### Example 1: View all pending repositories needing reconciliation

```bash
gitmap reconcile
```

**Output:**

```
  ── Pull Remediation Needed ── (2 repo(s))

  1. codelane (+12 untracked)
     gitmap reconcile codelane stash   (stash untracked, pull, re-apply)
     gitmap reconcile codelane wip     (commit WIP, pull --rebase)
     gitmap reconcile codelane discard (discard local changes, pull)

  2. atto-property (+11 untracked)
     gitmap reconcile atto-property stash
     gitmap reconcile atto-property wip
     gitmap reconcile atto-property discard
```

### Example 2: Reconcile a specific repository by name with stash

```bash
gitmap reconcile codelane stash
```

**Output:**

```
ℹ Applying Fix: Option 1 (Stash & Re-apply) on codelane
  Running: git -C "D:\wp-work\riseup-asia\codelane" stash -u && git -C "D:\wp-work\riseup-asia\codelane" pull && git -C "D:\wp-work\riseup-asia\codelane" stash pop

Saved working directory and index state WIP on main: cbf8072
Updating cbf8072..1ccac7c
Fast-forward
Dropped refs/stash@{0}

✓ Fix applied successfully on codelane
```

### Example 3: Reconcile all pending repositories with discard

```bash
gitmap reconcile --all discard
```

**Output:**

```
ℹ Reconciling 2 repository(ies) with action: discard

ℹ Applying Fix: Option 3 (Discard Local Changes) on codelane
✓ Fix applied successfully on codelane

ℹ Applying Fix: Option 3 (Discard Local Changes) on atto-property
✓ Fix applied successfully on atto-property
```

## See Also

- [pull](pull.md) — Pull tracked repositories concurrently
- [status](status.md) — Inspect repository working tree states
- [fix](fix.md) — Apply remediation to the last failed repository
