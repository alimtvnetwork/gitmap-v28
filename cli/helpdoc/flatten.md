# gitmap repo flatten

Compress and squash bloated repository history into a single clean snapshot commit.

> **`flatten` vs `git reset --soft`**
>
> While `git reset --soft` squashes commits on the current branch, it leaves stale packfiles, loose objects, reflogs, and unreachable blobs inside `.git/objects/`.
> `gitmap repo flatten` physically stages a clean successor repository (`<repo>-flat-v1`), copies pure working tree files while stripping build caches, initializes a fresh Git history, archives the bloated tree to `<repo>-marked-for-deletion`, and swaps to the canonical name with instant one-command undo support.

## Synopsis

```bash
gitmap repo flatten                       # flatten repository in current working directory
gitmap repo flatten --dry-run             # inspect commit count, tree size, and staging plan without changes
gitmap repo flatten -y                    # skip confirmation prompt and squash history immediately
gitmap repo flatten --target ./my-repo    # target a specific local repository path
gitmap flatten                            # top-level utility alias
gitmap repo undo-flatten <repo>           # restore previous bloated repository from marked-for-deletion archive
```

## Workflow & Staging Architecture

When executed, the flattening engine runs a zero-risk 6-stage lifecycle:

1. **Pre-flight Inspection**:
   - Resolves the target repository directory (defaults to cwd or `--target`).
   - Verifies `.git` repository existence and calculates commit count (`git rev-list --count HEAD`).
   - Computes working tree size while filtering excluded cache/build artifacts.
2. **Interactive Confirmation**:
   - Displays a bold warning banner unless `-y` or `--yes` is specified:
     `WARNING: This operation squashes repository history into a single flat commit. It is irreversible unless restored via undo.`
   - In non-TTY environments (CI/CD scripts), `-y` is mandatory.
3. **Successor Directory Staging**:
   - Generates an incremental numbered staging directory: `<repo>-flat-v1` (or next free `-vN` index).
   - Copies all regular working tree files preserving permissions.
   - Strictly ignores bloated build folders:
     - `.git/` (clears old packfiles and commit bloat)
     - `node_modules/`
     - `dist/`
     - `temp/`
     - `.pnpm-store/`
4. **Clean Snapshot Commit**:
   - Initializes a fresh Git repository in the successor directory:
     `git init -b main`
   - Stages all copied working files:
     `git add -A`
   - Creates a single flat squashed commit:
     `commit -m "flat: initial squashed repository snapshot"`
5. **Atomic Archive & Swap**:
   - Escapes any active shell directory handles holding the target open.
   - Archives the bloated original repository to `<repo>-marked-for-deletion`.
   - Promotes `<repo>-flat-v1` to the canonical `<repo>` path.
6. **Undo Telemetry & Guidance**:
   - Displays success summary with the exact undo command:
     `gitmap repo undo-flatten <repo>`

## Flags

| Flag | Default | Description |
|---|---|---|
| `--target <dir>` | `.` | Target repository directory to flatten. Can also be supplied as a positional argument. |
| `--dry-run` | `false` | Preview repo path, commit count, working tree size, and projected successor directory name without writing any changes to disk. |
| `-y`, `--yes` | `false` | Skip interactive `[y/N]` confirmation prompt. Mandatory in automated headless CI environments. |

## Undo Guidance

If you need to restore previous commit history, tags, or branches after flattening:

```bash
gitmap repo undo-flatten <repo>
```

This restores the canonical `<repo>` directly from `<repo>-marked-for-deletion` and removes the flattened working tree.

## Examples

```bash
# 1. Preview flattening for current repository
gitmap repo flatten --dry-run

# 2. Flatten current repository with interactive confirmation
gitmap repo flatten

# 3. Flatten a specific sibling repository non-interactively
gitmap repo flatten --target ./story-telling-v1 -y

# 4. Undo flattening and restore original history
gitmap repo undo-flatten story-telling-v1
```

## Exit Codes

- `0` — Success or dry-run preview completed.
- `1` — Execution failure (invalid target directory, git command failure, or archive rename error).

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter flatten
```
