# Spec: `gitmap gomod` — Go Module Path Rename

## Purpose

Rename a Go module path across an entire repository: updates `go.mod` and every `.go` file that imports the old module path. Wraps the operation in a safe branch workflow with backup and auto-merge.

## Command

```
gitmap gomod <new-module-path>
```

**Alias:** `gm`

## Behavior

### 1. Read Current Module Path

- Open `go.mod` in the current working directory.
- Parse the first `module` directive line (e.g., `module github.com/alimtvnetwork/core`).
- Extract the old module path: `github.com/alimtvnetwork/core`.
- If `go.mod` is not found or has no `module` line, print error and exit.

### 2. Validate Input

- `<new-module-path>` is required; if missing, print usage and exit.
- If old path equals new path, print "nothing to rename" and exit.
- Must be inside a Git work tree; if not, print error and exit.

### 3. Branch Workflow

Before any file modifications:

1. **Record current branch** — e.g., `main`.
2. **Derive slug** — sanitize `<new-module-path>` for branch names (replace `/` with `-`, strip special chars).
3. **Create backup branch** — `backup/before-replace-<slug>` from current HEAD. This branch is never modified; it exists purely as a restore point.
4. **Create feature branch** — `feature/replace-<slug>` from current HEAD. Check out this branch.

If either branch already exists, abort with an error message.

### 4. Replace Module Path

On the `feature/replace-<slug>` branch:

1. **Update `go.mod`** — Replace the `module <old-path>` line with `module <new-module-path>`.
2. **Update all `.go` files** — Recursively walk the repo directory (excluding `vendor/`, `.git/`, `node_modules/`). In every `.go` file, replace all occurrences of the old module path with the new one. This covers:
   - `import "old/path/pkg"` → `import "new/path/pkg"`
   - `import ( "old/path/pkg" )` → `import ( "new/path/pkg" )`
   - Any string literal containing the old path (rare but possible in `go:generate` directives).
3. **Run `go mod tidy`** — If `go` is on PATH, run `go mod tidy` to clean up. If it fails, print a warning but continue.

### 5. Commit

- Stage all changed files: `git add -A`
- Commit with message:
  ```
  refactor: rename go module path

  Old: github.com/alimtvnetwork/core
  New: x/y

  Replaced module directive in go.mod and all import paths
  across <N> .go files.
  ```
- `<N>` is the count of `.go` files that were modified.

### 6. Merge Back

1. **Checkout original branch** (e.g., `main`).
2. **Merge feature branch** — `git merge feature/replace-<slug> --no-ff -m "merge: module rename to <new-module-path>"`.
3. If merge fails (conflicts), print error and instruct user to resolve manually. Do not force.

### 7. Summary Output

Print a summary:

```
✔ Module path renamed
  Old: github.com/alimtvnetwork/core
  New: x/y
  Files updated: 47
  Backup branch: backup/before-replace-x-y
  Feature branch: feature/replace-x-y
  Merged into: main
```

## Flags

| Flag | Description |
|------|-------------|
| `--dry-run` | Preview changes without modifying any files or creating branches. Print what would be renamed and how many files would change. |
| `--no-merge` | Create branches and commit but do not merge back. Leave on the feature branch. |
| `--no-tidy` | Skip `go mod tidy` after replacement. |
| `--verbose` | Print each file path as it is modified. |
| `--ext <exts>` | Comma-separated file extensions to restrict replacement (e.g. `*.go,*.md,*.txt`). If omitted, all files in the repo are checked. `go.mod` is always updated regardless of this flag. |

## Edge Cases

| Scenario | Behavior |
|----------|----------|
| No `go.mod` in current directory | Error: `go.mod not found in current directory` |
| Not inside a Git repo | Error: `not inside a git repository` |
| Old path == new path | Info: `module path is already <path>, nothing to rename` |
| Branch already exists | Error: `branch <name> already exists, aborting` |
| No `.go` files contain old path | Warning: `no import paths found to replace (only go.mod updated)` |
| Dirty working tree | Error: `working tree has uncommitted changes, commit or stash first` |
| Merge conflict | Error: `merge conflict — resolve manually on <branch>` |
| `go mod tidy` fails | Warning: `go mod tidy failed: <error> (continuing)` |
| Vendor directory present | Skip `vendor/` during replacement (user should re-vendor after) |

## File Layout

| File | Responsibility |
|------|---------------|
| `cmd/gomod.go` | Flag parsing, orchestration, summary output |
| `cmd/gomodreplace.go` | File walking, path replacement logic |
| `cmd/gomodbranch.go` | Branch creation, merge, slug generation |

## Constants

Add to `constants/constants_cli.go`:
- `CmdGoMod = "gomod"`
- `CmdGoModAlias = "gm"`

Add to `constants/constants_messages.go`:
- Error/success format strings for all messages above.

## Dispatch

Add to `dispatchMisc` in `cmd/root.go`:
```go
if command == constants.CmdGoMod || command == constants.CmdGoModAlias {
    runGoMod(os.Args[2:])
    return true
}
```

## Help Text

Add to `constants/constants_cli.go`:
```
HelpGoMod = "  gomod (gm) <path>   Rename Go module path across repo with branch safety"
```

## Acceptance Criteria

1. Running `gitmap gomod "x/y"` in a repo with module `github.com/old/name` replaces all occurrences in `go.mod` and `.go` files.
2. Backup branch is created and left untouched.
3. Feature branch contains exactly one commit with the replacement.
4. Feature branch is merged back to the original branch.
5. `--dry-run` makes zero changes and prints accurate preview.
6. Dirty working tree is rejected before any work begins.
7. Existing branch names cause an abort, not an overwrite.
