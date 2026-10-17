# gitmap resolve

Resolve merge and rebase conflicts across repositories using IDE 3-way merge editors, remote/local checkout strategies, or safe abort.

## Alias

res

## Usage

```bash
gitmap resolve [repo] [--theirs | --ours | --ide | --abort | --rerere]
```

## Description

Resolves active merge, rebase, or cherry-pick conflicts on a specific repository or the current working directory.
When conflicts arise during `gitmap pull`, `gitmap reconcile`, or manual rebases, `gitmap resolve` provides instant automated resolutions:

- **`--theirs`**: Automatically accepts incoming (remote/upstream) versions of all conflicted files and stages them.
- **`--ours`**: Automatically preserves current (local) versions of all conflicted files and stages them.
- **`--ide`**: Launches your preferred installed IDE (VS Code, Cursor, Windsurf, or VSCodium) with the 3-Way Merge Editor walkthrough.
- **`--abort`**: Safely aborts any in-progress merge, rebase, or cherry-pick and returns the repository to a clean state.
- **`--rerere`**: Enables Git recorded resolution (`rerere.enabled` and `rerere.autoupdate`) to automatically resolve recurring conflicts.

## Options

| Flag | Short | Action | Description |
|------|-------|--------|-------------|
| `--theirs` | `-t` | Accept Remote | Overwrite conflicted files with incoming remote version (`git checkout --theirs`) and stage them |
| `--ours` | `-o` | Accept Local | Keep current local version (`git checkout --ours`) and stage them |
| `--ide` | `-i`, `--code` | Launch IDE | Launch VS Code / Cursor / Windsurf / VSCodium in 3-Way Merge Editor |
| `--abort` | `-a` | Abort | Abort active merge (`git merge --abort`) or rebase (`git rebase --abort`) |
| `--rerere` | | Enable rerere | Configure `git config rerere.enabled true` and `rerere.autoupdate true` |

## VS Code 3-Way Merge Editor Walkthrough

When resolving complex conflicts interactively:

1. Run `gitmap resolve [repo] --ide` (or `gitmap fix [repo] ide`).
2. In VS Code / Cursor, open the Source Control panel (`Ctrl+Shift+G` or `Cmd+Shift+G` on macOS).
3. Under **Merge Changes**, click on the conflicted file, or click the **Resolve in Merge Editor** button at the bottom right of the editor.
4. The 3-Way Merge Editor displays:
   - **Incoming (Remote)**: incoming branch changes.
   - **Current (Local)**: your local branch changes.
   - **Result (Bottom)**: the editable output file with conflict markers resolved.
5. Click **Accept Incoming** or **Accept Current** checkboxes, or manually edit the Result pane.
6. Press `Ctrl+S` (`Cmd+S`) to save. VS Code automatically removes conflict markers and stages the file.
7. Finalize and push with GitMap:
   ```bash
   gitmap cpf "resolve merge conflicts"
   ```

## Examples

### Example 1: Launch IDE Merge Editor for current repository

```bash
gitmap resolve --ide
```

**Output:**

```
✓ Launched VS Code for D:\work\my-service

=================================================================
  VS Code 3-Way Merge Editor Walkthrough
=================================================================
  Conflicted Files (2):
    • src/auth/login.ts
    • package.json

  1. Press Ctrl+Shift+G (Cmd+Shift+G on macOS) to open Source Control
  2. Click a conflicted file, then click 'Resolve in Merge Editor'
  3. Review the 3-Way split:
     - Incoming (Remote): Changes from incoming branch
     - Current (Local):   Your current branch changes
     - Result (Bottom):   Final merged output
  4. Click 'Accept Incoming', 'Accept Current', or combine changes
  5. Save the file (Ctrl+S / Cmd+S) — Git marks it resolved
  6. Stage & commit with GitMap:
     gitmap cpf "resolve merge conflicts"
=================================================================
```

### Example 2: Accept incoming remote changes for a specific repository

```bash
gitmap resolve my-service --theirs
```

**Output:**

```
✓ Accepted incoming (theirs) changes on 2 conflicted file(s) in D:\work\my-service
  Next: Run 'gitmap cpf "resolve merge conflicts"' to commit and push.
```

### Example 3: Accept current local changes

```bash
gitmap resolve my-service --ours
```

**Output:**

```
✓ Accepted current (ours) changes on 2 conflicted file(s) in D:\work\my-service
  Next: Run 'gitmap cpf "resolve merge conflicts"' to commit and push.
```

### Example 4: Safely abort an in-progress merge or rebase

```bash
gitmap resolve my-service --abort
```

**Output:**

```
✓ Merge aborted successfully in D:\work\my-service
```

### Example 5: Enable git rerere for automatic replay of conflict resolutions

```bash
gitmap resolve my-service --rerere
```

**Output:**

```
✓ Git rerere (Reuse Recorded Resolution) enabled for D:\work\my-service
```

## See Also

- [conflicts](conflicts.md) — Inspect all conflicted and unmerged files across repositories
- [reconcile](reconcile.md) — Reconcile dirty repositories after pull
- [fix](fix.md) — Remediation walkthrough for sync failures
- [pull-commit-push](pull-commit-push.md) — Commit and push once conflicts are resolved
