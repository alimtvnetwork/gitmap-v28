# gitmap fix

Remediate dirty repositories and merge conflicts with batch or interactive prompt walkthrough.

## Usage

```bash
gitmap fix
gitmap fix ls
gitmap fix list
gitmap fix all [stash|wip|discard|1|2|3]
gitmap fix --prompt
gitmap fix [repo-name] [stash|wip|discard|1|2|3]
```

## Description

`gitmap fix` resolves repositories that have uncommitted local changes, untracked files,
or conflicts preventing clean synchronization. It offers four execution modes:

1. **List Repositories with Issues (`gitmap fix ls` / `gitmap fix list`)**: Live-scans all
   tracked repositories across the workspace, detects dirty files, merge conflicts, behind/ahead
   branch desync, unpopped stashes, and git locks, and displays an aligned diagnostic table with
   suggested remediation commands.
2. **Batch All Mode (`gitmap fix all`)**: Automatically applies the remediation strategy
   (defaulting to `stash` & re-apply) across all pending repositories in sequence.
3. **Interactive Prompt Walkthrough (`gitmap fix --prompt` / `gitmap fix -p`)**: Steps through each
   pending project one by one, displaying its exact pending changes and modified/untracked files
   with colored status tags, allowing individual selection, batch continuation, or skipping.
4. **Targeted Mode (`gitmap fix <repo> [1|2|3]`)**: Directly executes a remediation strategy
   on a specific repository.

## Remediation Strategies

| Option | Name | Action | Description |
|--------|------|--------|-------------|
| `1` | Stash & Re-apply | `stash` | Stashes local changes (`-u`), pulls remote changes, then restores the stash |
| `2` | Commit WIP | `wip` | Stages all changes, commits as temporary WIP commit, then pulls with `--rebase` |
| `3` | Discard Local | `discard` / `clean` | Permanently discards changes (`reset --hard` & `clean -fd`), then pulls |

## Interactive Walkthrough Controls

When running in interactive mode (`gitmap fix`, `gitmap fix --prompt`, or `gitmap fix all --prompt`),
the following choices are available at each repository prompt:

- `1` or `stash`: Apply Stash & Re-apply to this repository.
- `2` or `wip`: Apply Commit WIP to this repository.
- `3` or `discard`: Discard local changes and pull.
- `s` or `skip`: Skip this repository without modifications.
- `a` or `all`: Apply stash to this repository and all remaining pending repositories.
- `q` or `quit`: Abort interactive walkthrough immediately.

## Examples

### Example 1: Batch fix all pending repositories with default stash strategy

```bash
gitmap fix all
```

**Output:**

```
ℹ Remediating 3 repository(ies) with action: stash

ℹ Applying Fix: Option 1 (Stash & Re-apply) on icon-coding-guidelines
  Plan:    Temporarily save local changes (including untracked), pull latest remote commits, then re-apply

  [1/3] ➜ git stash -u ... ✔ ok
  [2/3] ➜ git pull ... ✔ ok
  [3/3] ➜ git stash pop ... ✔ ok

✓ Fix applied successfully on icon-coding-guidelines

ℹ Applying Fix: Option 1 (Stash & Re-apply) on bsrm-presentation-hiltrax-v4
  Plan:    Temporarily save local changes (including untracked), pull latest remote commits, then re-apply

  [1/3] ➜ git stash -u ... ✔ ok
  [2/3] ➜ git pull ... ✔ ok
  [3/3] ➜ git stash pop ... ✔ ok

✓ Fix applied successfully on bsrm-presentation-hiltrax-v4

✓ All 3 repository(ies) remediated successfully.
```

### Example 2: Interactive step-by-step prompt walkthrough showcasing files

```bash
gitmap fix --prompt
```

**Output:**

```
[1/3] icon-coding-guidelines (+1 staged, +1 modified, +6 untracked)
    Path: D:\work\icon-coding-guidelines
    Pending Changes (8 files):
      • [staged]    package.json
      • [modified]  src/index.ts
      • [untracked] temp.txt
      ... and 5 more files
    Remediation Options:
      [1] Stash & Re-apply   Stash local changes (-u), pull remote, then pop stash
      [2] Commit WIP         Commit all modified/untracked files, then pull --rebase
      [3] Discard Local      Permanently discard changes (reset --hard & clean -fd), pull
      [s] Skip               Skip this repository for now
      [a] Apply to All       Apply stash to this and all remaining repositories
      [q] Quit               Exit interactive prompt
  Pick [1=stash, 2=wip, 3=discard, s=skip, a=all-stash, q=quit]: 1
```

### Example 3: Batch fix with explicit WIP commit strategy

```bash
gitmap fix all wip
```

### Example 4: Scan and list all repositories with issues

```bash
gitmap fix ls
```

**Output:**

```
ℹ Repositories Requiring Fix / Remediation (2):

  #  Repository             Status / Issues                   Branch      Suggested Fix
  ─  ─────────────────────  ────────────────────────────────  ──────────  ───────────────────────────────
  1  icon-coding-guideline  dirty: +1 staged, +2 modified     main        gitmap fix icon-coding-guideline 1
  2  riseup-asia            behind (3 commits)                main        gitmap fix riseup-asia 1

  Remediation Commands:
    gitmap fix <repo> [1|2|3]        (1=stash, 2=wip, 3=discard)
    gitmap fix all                   (apply stash to all repositories)
    gitmap fix all [1|2|3]           (apply specific strategy to all)
    gitmap fix --prompt              (step-by-step interactive walkthrough)
```

## See Also

- [reconcile](reconcile.md) — Reconcile failed repositories
- [pull](pull.md) — Pull tracked repositories
- [status](status.md) — Inspect working tree status
