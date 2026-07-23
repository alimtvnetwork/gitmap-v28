# Spec 100: `gitmap clone-pick` (alias `cpk`)

> Status: Draft (target v3.153.0)
> Owner: cmd / clonepick package
> Related: spec 88 (clone-direct-url), spec 99 (CLI uniqueness CI guard),
> mem://features/clone-direct-url, mem://features/command-help-system

## 1. Summary

`gitmap clone-pick <repo-url> <relative-path>[,<relative-path>...]` performs
a **partial clone** of a Git repository: fetch only the requested
subdirectories/files into the current working directory using git's native
sparse-checkout machinery.

Every selection is auto-saved to a new SQLite table so the same selection
can be re-run with `--replay`.

The command intentionally does NOT add an alias clash with the existing
"interactive" command (`i`). Reserved short alias is `cpk`. The originally
proposed `ci`/`intra` was rejected because `ci` collides with CI/CD
muscle memory; `cpk` reads as "clone-pick" with no overlap.

## 2. Why a new command

| Existing command | Behaviour | Why insufficient |
|------------------|-----------|------------------|
| `clone <url>` | Full clone, optional folder name | Always full repo, no path filter |
| `clone-from <file>` | Clone many repos from a manifest | Manifest-driven, not path-subset |
| `clone-now <file>` | Re-clone from scan output | Round-trip cloner, no path filter |
| `clone-multi` | Multi-URL convenience | Still full clones |

None can fetch *part* of a single repo. Sparse-checkout is the right
primitive for "I just want `docs/` and `examples/foo.md` from this repo,
not the whole 800 MB monorepo."

## 3. Surface

```
gitmap clone-pick <repo-url> <paths> [flags]
gitmap cpk        <repo-url> <paths> [flags]
gitmap clone-pick --replay <id|name> [flags]
```

`<paths>` is a comma-separated list of repo-relative paths
(e.g. `docs,examples/foo.md,scripts/release-version.ps1`). Folders and
files both accepted. Leading `./` and trailing `/` are normalised away.

### 3.1 Flags

| Flag | Default | Purpose |
|------|---------|---------|
| `--ask` | false | Open interactive tree picker (see §6) before fetching |
| `--name <label>` | "" | Human label saved to DB for later `--replay <name>` |
| `--replay <id\|name>` | "" | Re-run a previously saved selection (no `<paths>` needed) |
| `--branch <name>` | "" | Pin branch (`--branch` on the underlying clone) |
| `--mode <https\|ssh>` | https | URL form to use when the input is shorthand `owner/repo` |
| `--depth <n>` | 1 | Shallow clone depth (`0` = full history) |
| `--cone` | true | Use sparse-checkout cone mode (faster, folder-only). Disable for file-level globs |
| `--dest <dir>` | "." | Destination directory; created if missing |
| `--keep-git` | true | Leave the `.git` dir in `<dest>`. `--keep-git=false` deletes `.git` after checkout (files-only mode) |
| `--dry-run` | false | Print plan + git commands without executing |
| `--quiet` | false | Suppress per-step progress on stderr |
| `--force` | false | Allow non-empty `<dest>` (default refuses to clobber) |

### 3.2 Exit codes

| Code | Meaning |
|------|---------|
| 0 | Success (or `--dry-run` rendered without error) |
| 1 | Runtime failure: git error, fs error, replay row not found |
| 2 | Bad CLI usage: missing `<repo-url>`, missing `<paths>` (when no `--replay`), invalid flag value |

## 4. Resolution of `<repo-url>`

Reuse the URL parser from `gitutil.CanonicalRepoID` (added with the export
extension) so:

- Full HTTPS / SSH URLs pass through verbatim
- `owner/repo` shorthand expands per `--mode` to:
  - https → `https://github.com/owner/repo.git`
  - ssh   → `git@github.com:owner/repo.git`
- `host/owner/repo` shorthand for non-GitHub hosts (gitlab, bitbucket)
  uses the same expansion table

URL canonicalisation happens once; the canonical form is what gets stored
in `CloneInteractiveSelection.RepoCanonicalId` so a later `--replay`
matches regardless of HTTPS↔SSH variation.

## 5. Sparse-checkout pipeline

Sparse checkout was selected over "shallow clone to /tmp + copy" because:

1. Single git invocation tree, no temp dir bookkeeping
2. Works with any host (no GitHub-API dependency)
3. Leaves a real `.git` dir so the user can `git pull`/`git switch`
4. `--cone` mode is O(matched paths), not O(repo size)

### 5.1 Steps

```bash
# 1. partial clone, no working tree
git clone --filter=blob:none --no-checkout \
  [--branch <name>] [--depth <n>] <url> <dest>

# 2. enable sparse-checkout
cd <dest>
git sparse-checkout init [--cone]   # cone unless --cone=false

# 3. set the path patterns (cone mode = folders only)
git sparse-checkout set <path1> <path2> ...

# 4. materialise the working tree
git checkout            # or `git read-tree -mu HEAD` when no branch given

# 5. optional: drop .git when --keep-git=false
rm -rf .git
```

### 5.2 Path-mode auto-detection

If any `<path>` contains `*`, `?`, `[`, or has `.<ext>` past a `/` we
silently flip `--cone=false` so non-cone (full pattern) sparse-checkout is
used. The user can override either direction with explicit `--cone` or
`--cone=false`.

### 5.3 Validation

- Reject empty path entries (`a,,b` → exit 2 with clear error)
- Reject absolute paths (`/etc/passwd` → exit 2)
- Reject `..` traversal (`../foo` → exit 2)
- Reject paths longer than 4096 bytes (sparse-checkout limit)

## 6. `--ask` interactive picker

Built with `bubbletea` (already vendored, used by `tui/browser.go`).

### 6.1 Source of the tree

Before launching the picker we must know what's in the repo. Two-step:

1. Run `git ls-tree -r --name-only HEAD` against a *partial* clone (`--filter=blob:none --depth=1 --no-checkout`) into a temp dir.
2. Parse the output into a tree-of-nodes structure rooted at `/`.

The temp clone is reused for the actual checkout (we just `mv` it to
`<dest>` and run `sparse-checkout set`) so we don't pay the clone twice.

### 6.2 UI

```
[ gitmap clone-pick — github.com/owner/repo ]

  ▸ [x] docs/                       (pre-selected from CLI)
    [ ] examples/
    [x] scripts/release-version.ps1 (pre-selected from CLI)
    [ ] src/
    [ ] tests/

  ↑/↓ move   →/enter expand   space toggle   a all   n none
  s save & clone   q quit
```

- User-supplied `<paths>` are pre-checked
- Toggling a folder toggles all descendants (cone-friendly)
- Expanding a folder reveals children; toggling a child auto-flips cone
  off and stores the file-level selection
- `s` writes selection to DB and proceeds to checkout
- `q` exits with code 130 (user-cancelled)

### 6.3 Defaults shipped with the picker

- All `.git/`, `node_modules/`, `vendor/`, `dist/`, `build/`,
  `__pycache__/` folders rendered greyed-out and pre-unchecked even if
  globbed. User can manually re-include them.

These defaults live in `constants.ClonePickAutoExclude` so they're
overridable via config (`clonePick.autoExclude` array).

## 7. Persistence

### 7.1 New table `CloneInteractiveSelection`

```sql
CREATE TABLE IF NOT EXISTS CloneInteractiveSelection (
    SelectionId       INTEGER PRIMARY KEY AUTOINCREMENT,
    Name              TEXT NOT NULL DEFAULT '',
    RepoCanonicalId   TEXT NOT NULL,         -- e.g. github.com/owner/repo
    RepoUrl           TEXT NOT NULL,         -- canonical URL used at clone time
    Mode              TEXT NOT NULL DEFAULT 'https',   -- https | ssh
    Branch            TEXT NOT NULL DEFAULT '',
    Depth             INTEGER NOT NULL DEFAULT 1,
    Cone              INTEGER NOT NULL DEFAULT 1,      -- 0/1
    KeepGit           INTEGER NOT NULL DEFAULT 1,      -- 0/1
    DestDir           TEXT NOT NULL DEFAULT '.',
    PathsCsv          TEXT NOT NULL,         -- normalised, sorted, comma-joined
    UsedAsk           INTEGER NOT NULL DEFAULT 0,      -- 0/1
    CreatedAt         TEXT DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_clonepick_repocanon
    ON CloneInteractiveSelection(RepoCanonicalId);
CREATE INDEX IF NOT EXISTS idx_clonepick_name
    ON CloneInteractiveSelection(Name) WHERE Name <> '';
```

- `Name` is optional but unique-non-empty across the table (enforced in
  the store layer, not as a UNIQUE constraint, because '' must repeat).
- No FK to `Repo(RepoId)` — the picked repo may not be in any local scan,
  and we don't want to require a prior `gitmap scan` run.

### 7.2 Insert policy

- Every successful run inserts exactly one row (auto-save, per the
  user's chosen behaviour).
- `--dry-run` does NOT insert.
- `--replay` does NOT insert a duplicate row; it updates the
  `CreatedAt` of the matched row to record the re-run time.

### 7.3 `--replay <id|name>` lookup

1. If the value parses as a positive integer → SELECT by `SelectionId`.
2. Else → SELECT by `Name = ?` (case-sensitive).
3. Zero results → exit 1 with `MsgClonePickReplayNotFound`.
4. >1 results (only possible if `Name = ''` was bypassed somehow) →
   exit 1 with `MsgClonePickReplayAmbiguous` and list candidate IDs.

## 8. Shell handoff

After a successful checkout (when `<dest>` was created or freshly
populated), call `WriteShellHandoff(dest)` so the wrapper `cd`s the user
into the cloned tree. Same pattern as `gitmap clone <url>`.

Skipped when `<dest> == "."` (already there) or `--dry-run`.

## 9. Help text

`gitmap/helptext/clone-pick.md` follows the standard 120-line limit and
3-8 line realistic simulations. Examples cover:

1. Single-folder pick: `gitmap cpk owner/repo docs`
2. Multi-path pick: `gitmap cpk owner/repo docs,examples,README.md`
3. Interactive pick: `gitmap cpk owner/repo --ask`
4. Save + replay: `gitmap cpk owner/repo docs --name docs-only` then
   `gitmap cpk --replay docs-only`
5. Files-only (no `.git`): `gitmap cpk owner/repo docs --keep-git=false`

## 10. Tests

| Layer | Test | File |
|-------|------|------|
| Parse | URL + paths normalisation | `clonepick/parse_test.go` |
| Parse | Reject empty / abs / `..` paths | `clonepick/parse_test.go` |
| Plan | Cone vs non-cone auto-detection | `clonepick/plan_test.go` |
| DB | Insert + lookup by id + by name | `store/cloneinteractiveselection_test.go` |
| DB | Replay updates CreatedAt, no duplicate row | `store/cloneinteractiveselection_test.go` |
| Cmd | Missing args → exit 2 | `cmd/clonepick_test.go` |
| Cmd | `--replay <unknown>` → exit 1 | `cmd/clonepick_test.go` |
| Cmd | `--dry-run` prints commands, no DB write | `cmd/clonepick_test.go` |
| Help | `clone-pick` registered in helptext registry | `helptext/coverage_test.go` |
| Marker | `// gitmap:cmd top-level` present on const | drift CI |

## 11. Future extensions (out of scope for v1)

- `gitmap clone-pick --update` to re-run sparse-checkout with the same
  selection but a fresh fetch.
- Group several selections into a `ClonePickProfile` for one-shot
  multi-repo subset clones.
- `clonepick.autoExclude` overridable per-repo via `.gitmap-pick.json`
  inside the cloned tree.

## 12. Open questions

None at spec time — all four ambiguity points were answered by the user
before drafting (command name, sparse-checkout, tree picker, auto-save).
