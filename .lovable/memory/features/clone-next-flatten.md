---
name: Clone-next flatten mode
description: Clone-next flattens by default; -f / --force flag forces flatten even when cwd IS the target folder by chdir-to-parent before remove
type: feature
---

## Feature: Default Flatten for `clone-next`

### Status: ✅ Implemented (v2.75.0); `-f` force flag added v3.50.0

### Behavior (Default — No Flag Required)

When `gitmap cn v++` is used on a repo like `macro-ahk-v15`:

1. Target clone folder is `macro-ahk/` (base name only, version suffix stripped)
2. If `macro-ahk/` already exists, remove it entirely first (no prompt)
3. Clone the target repo (e.g., `macro-ahk-v16`) into `macro-ahk/`
4. The remote URL still points to `macro-ahk-v16` on GitHub — only the local folder name is flattened
5. Update the database with the new version info
6. Record version transition in `RepoVersionHistory`

### Force-Flatten (`-f` / `--force`) — v3.50.0+

Solves the "already-flattened cwd" deadlock:

- **Scenario:** user is inside `macro-ahk/` (the flattened folder from v21) and runs `gitmap cn v++ -f`. Without `-f`, Windows holds the cwd open → `os.RemoveAll` fails → falls back to `macro-ahk-v22/`.
- **With `-f`:** before the existence check, gitmap `os.Chdir`s to the parent dir, then removes the now-unlocked cwd, then clones `macro-ahk-v22` into `macro-ahk/`. After clone, `os.Chdir`s into the new `macro-ahk/` and emits `GITMAP_SHELL_HANDOFF` so the shell wrapper follows.
- **Refuses the versioned-folder fallback.** If removal still fails (locked by some other process), `-f` aborts with `ErrCloneNextForceFailed` instead of silently degrading to `macro-ahk-v22/`. This is the explicit contract: user asked for a flat layout and gets either that or a clear error.
- Compatible with `--delete` / `--keep`: those still control disposition of the *prior versioned* folder when cwd is on disk under a versioned name (the older "v15 → flatten" upgrade flow).

### Key Code

- Flag: `CloneNextFlags.Force` in `gitmap/cmd/clonenextflags.go` (long `--force`, short `-f`).
- Force-handling helper: `forceReleaseLockOnCwd` in `gitmap/cmd/clonenext.go` — chdirs to parent ONLY when `Force && cwd == targetPath`.
- Fallback gate: when `Force` is true, the code path that sets `flattenedFolder = fallbackFolder` is replaced by `os.Exit(1)` after `ErrCloneNextForceFailed`.

### `gitmap clone <url>` Auto-Flatten

When cloning a versioned URL without a custom folder name:
- `gitmap clone https://github.com/user/wp-onboarding-v13` → clones into `wp-onboarding/`
- `gitmap clone https://github.com/user/wp-onboarding-v13 my-folder` → clones into `my-folder/` (no flatten)

### Database Schema

#### Repos table — version columns
- `CurrentVersionTag TEXT DEFAULT ''` — e.g., "v16"
- `CurrentVersionNum INTEGER DEFAULT 0` — e.g., 16

#### `RepoVersionHistory` table
Tracks every version transition with from/to tags, numbers, flattened path, and timestamp.

### Related Commands
- `gitmap version-history` (`vh`) — Display all version transitions for the current repo
