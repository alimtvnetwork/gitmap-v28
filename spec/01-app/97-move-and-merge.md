# Move & Merge — `gitmap-v28 mv` / `merge-both` / `merge-left` / `merge-right`

> **Status:** Draft — spec only, implementation pending.
> **Related specs:**
> - [05-cloner.md](05-cloner.md) — clone primitive reused for URL endpoints
> - [88-clone-direct-url.md](88-clone-direct-url.md) — direct-URL clone semantics this spec depends on
> - [96-clone-replace-existing-folder.md](96-clone-replace-existing-folder.md) — two-strategy folder replacement reused on `mv` overwrite
> - [02-cli-interface.md](02-cli-interface.md) — global flag conventions (`-y`, `-a`)

## Overview

A unified family of four commands that moves or merges files between
two endpoints, where each endpoint (`LEFT` and `RIGHT`) can be either a
local folder OR a remote Git URL. The CLI normalises both endpoints
into resolved local working folders, performs the file-level
operation, and (when an endpoint was a URL) commits and pushes the
result back to its remote.

```
gitmap-v28 mv          LEFT RIGHT     # move LEFT into RIGHT, delete LEFT
gitmap-v28 merge-both  LEFT RIGHT     # bidirectional fill: each side gains the other's missing files
gitmap-v28 merge-left  LEFT RIGHT     # take from RIGHT, put into LEFT
gitmap-v28 merge-right LEFT RIGHT     # take from LEFT,  put into RIGHT
```

`mv` is an alias for `move`. The terms **LEFT/RIGHT** are used in
help text and prompts for the merge commands (more intuitive than
FROM/TO when the operation is bidirectional). For `mv`, **FROM/TO**
is also accepted as a synonym (FROM = LEFT, TO = RIGHT).

---

## Endpoint Resolution

Each of the two positional arguments is classified once at command
start:

| Form | Detected as | Resolution |
|------|-------------|------------|
| Starts with `https://`, `http://`, `git@`, `ssh://` | URL | See **URL endpoint** below |
| Anything else | Folder path | Resolved relative to CWD; absolute path kept as-is |

Optional `:branch` suffix on a URL endpoint pins the branch:

```
gitmap-v28 mv https://github.com/owner/repo:develop ./local
```

Folder paths do not accept a branch suffix (the working folder is
already on whatever branch is checked out).

### URL endpoint

The URL is mapped to a candidate working folder name:
`<basename-of-url-without-.git>` placed in CWD.

1. **Folder does NOT exist:** clone the URL into that folder.
   Normal `gitmap-v28 clone <url>` semantics (see
   [88-clone-direct-url.md](88-clone-direct-url.md)).
2. **Folder DOES exist:**
   a. Read its `origin` remote (`git remote get-url origin`).
   b. If `origin` matches the requested URL (after normalising
      `https/ssh` and trailing `.git`): treat it as the working
      folder. Run `git pull --ff-only` first; abort with a clear
      error if pull fails.
   c. If `origin` does NOT match: abort with
      `error: folder '<name>' exists but its remote is '<other>',
      not '<requested>'. Pass --force-folder to overwrite, or rename it.`
   d. With `--force-folder`: invoke the
      [96-clone-replace-existing-folder.md](96-clone-replace-existing-folder.md)
      two-strategy replace flow.

### Folder endpoint

1. **Path exists and is a Git repo:** use as-is. Run `git pull
   --ff-only` only if `--pull` was passed (folder endpoints don't
   auto-pull, since they may be intentionally offline).
2. **Path exists and is NOT a Git repo:** still allowed for `mv`
   (treated as a plain folder of files). For `merge-*` commands
   this is also allowed; `.git` checks are simply skipped.
3. **Path does not exist:**
   - For LEFT in any command: error
     (`error: source '<path>' does not exist`).
   - For RIGHT in `mv`: created automatically (the move target).
   - For RIGHT in `merge-*`: error (merge requires both sides).

---

## Operation: `mv` (move)

Semantics: copy LEFT's contents (excluding `.git/`) into RIGHT,
then **delete the LEFT folder entirely**.

```
+--------+      copy files (no .git)      +--------+
| LEFT/  | -------------------------->    | RIGHT/ |
| (.git) |                                | (.git) |
+--------+      then: rm -rf LEFT/        +--------+
```

### Step-by-step

1. Resolve both endpoints (clone if URL, pull if matching).
2. If RIGHT did not exist, create it (and `git init` only if
   `--init` was passed; otherwise it stays a plain folder).
3. Copy every entry under LEFT to RIGHT, **excluding `.git/`**.
   Conflict policy: any file that already exists in RIGHT is
   overwritten without prompt. (Use `merge-right` for safer copy
   with prompts.)
4. Delete LEFT recursively, including its `.git/` if any.
5. If RIGHT originated from a URL: stage all changes, commit with
   `gitmap-v28 mv from <LEFT-display>`, and `git push`.
6. If LEFT originated from a URL, its working folder was just
   deleted; nothing to push (the URL repo on the remote is
   **NOT** deleted — only the local clone).

### Flags

| Flag | Meaning |
|------|---------|
| `--no-push` | When RIGHT is a URL, stop after the local commit (skip `git push`). |
| `--no-commit` | When RIGHT is a URL, copy files but skip both commit and push. |
| `--force-folder` | If a URL endpoint maps to a folder whose origin doesn't match, replace it via [96-clone-replace-existing-folder.md](96-clone-replace-existing-folder.md). |
| `--pull` | Force `git pull --ff-only` on a pre-existing folder endpoint (no auto-pull by default for folder endpoints). |
| `--init` | When RIGHT is created fresh (folder endpoint that didn't exist), also run `git init` in it. |
| `--dry-run` | Print every action; perform none. |

`mv` does NOT prompt — the answer "Move-and-delete-FROM" is the
documented destructive semantic. Use `merge-right` for the safer
copy-with-prompt variant.

---

## Operation: `merge-both`

Semantics: each side gains every file the other side has but it
does not. Files that exist on **only one** side are copied to the
other. Files that exist on **both** sides with **different content**
trigger a conflict prompt.

```
                   missing on RIGHT --> copied
   +--------+    <----------------------    +--------+
   | LEFT/  |                               | RIGHT/ |
   +--------+    ---------------------->    +--------+
                   missing on LEFT --> copied
                       conflicts -> prompt
```

Files identical on both sides are no-ops.

`.git/` is always excluded from comparison and copying.

---

## Operation: `merge-left`

Semantics: copy every file from RIGHT into LEFT. Missing files are
added; conflicts prompt. RIGHT is never modified.

```
   +--------+    <----------------------    +--------+
   | LEFT/  |    files from RIGHT           | RIGHT/ |
   +--------+    (with conflict prompt)     +--------+
                                            (untouched)
```

---

## Operation: `merge-right`

Semantics: copy every file from LEFT into RIGHT. Missing files are
added; conflicts prompt. LEFT is never modified.

```
   +--------+    ---------------------->    +--------+
   | LEFT/  |    files from LEFT            | RIGHT/ |
   +--------+    (with conflict prompt)     +--------+
   (untouched)
```

---

## Conflict Prompt (merge-* commands)

A "conflict" is a path that exists on both sides with different
content (byte-level compare; identical files are skipped silently).

For each conflict, print:

```
  conflict: docs/architecture.md
    LEFT  : 4.2 KB  modified 2026-04-12 09:14
    RIGHT : 5.1 KB  modified 2026-04-15 17:03
  [L]eft  [R]ight  [S]kip  [A]ll-left  [B]all-right  [Q]uit
  > _
```

| Key | Action |
|-----|--------|
| `L` | Take LEFT's version (overwrite RIGHT — for `merge-left` and `merge-both`; no-op for `merge-right` since LEFT is the source). |
| `R` | Take RIGHT's version (overwrite LEFT — for `merge-right` and `merge-both`; no-op for `merge-left`). |
| `S` | Skip this file; both sides keep their current copy. |
| `A` | All-Left: apply `L` to this and every remaining conflict. |
| `B` | All-Right: apply `R` to this and every remaining conflict. |
| `Q` | Quit immediately. Already-applied changes are kept; no rollback. |

`L`/`R` semantics in `merge-both`: whichever side's content is
chosen is written to the other side (so both sides end up
identical for that file).

### Bypass: `-y` / `-a`

When `-y` (or its long form `--yes`) or `-a` (`--accept-all`) is
passed, no prompt is shown. The default per command is **source
side wins**:

| Command | `-y` / `-a` default |
|---------|--------------------|
| `merge-right` | LEFT wins (LEFT is the source) |
| `merge-left`  | RIGHT wins (RIGHT is the source) |
| `merge-both`  | Newer mtime wins (no single source side) |

To override the default, pass one of the explicit prefer flags:

| Flag | Effect (with `-y`) |
|------|--------------------|
| `--prefer-left` | LEFT always wins on conflict |
| `--prefer-right` | RIGHT always wins on conflict |
| `--prefer-newer` | Newer mtime wins (default for `merge-both`) |
| `--prefer-skip` | Skip every conflict; only missing files are copied |

`-y`/`-a` and the `--prefer-*` flags also imply `--no-confirm` on
the post-merge commit (when an endpoint is a URL).

---

## URL-Side Commit & Push

When either endpoint was a URL, after the file operation completes
the CLI runs in that working folder:

```
git add -A
git commit -m "<command-specific message>"
git push
```

| Command | Commit message template |
|---------|-------------------------|
| `mv` | `gitmap-v28 mv from <LEFT-display>` |
| `merge-both` | `gitmap-v28 merge-both with <other-display>` |
| `merge-left` | `gitmap-v28 merge-left from <RIGHT-display>` |
| `merge-right` | `gitmap-v28 merge-right from <LEFT-display>` |

`<LEFT-display>` / `<RIGHT-display>` is the original argument the
user typed (URL or folder path), trimmed.

If `git push` fails (auth, non-fast-forward, network):

1. Log the full git error.
2. Print: `Push failed. Local commit is preserved at <sha>.
   Resolve manually or re-run with --no-push to skip.`
3. Exit non-zero.

`--no-push` stops after the local commit. `--no-commit` stops after
copying files (no commit, no push) — useful for dry inspection
before committing.

---

## Branch Selector

A trailing `:<branch>` on a URL endpoint pins both clone-checkout
and post-merge push to that branch:

```
gitmap-v28 mv          ./local https://github.com/owner/repo:release
gitmap-v28 merge-right ./local https://github.com/owner/repo:feature/x
```

If the branch does not exist on the remote, the URL endpoint is
created on a new branch and pushed with `git push --set-upstream`.
Folder endpoints ignore any `:branch` suffix (folders are checked
out at whatever branch they currently hold).

---

## Examples

```
# move local folder into another local folder, deleting source
gitmap-v28 mv ./gitmap-v28 ./gitmap-v28

# move local folder into a remote repo (clone, copy, commit, push)
gitmap-v28 mv ./gitmap-v28 https://github.com/alimtvnetwork/gitmap-v28

# move a remote repo's contents into a local folder
gitmap-v28 mv https://github.com/alimtvnetwork/gitmap-v28 ./another-folder

# move between two remote repos (clones both, copies, pushes RIGHT)
gitmap-v28 mv https://github.com/alimtvnetwork/gitmap-v28 \
         https://github.com/alimtvnetwork/gitmap-v28

# merge missing files only (identical or differing files prompt)
gitmap-v28 merge-both ./gitmap-v28 ./gitmap-v28

# merge with auto-accept: each side's source wins
gitmap-v28 merge-right ./gitmap-v28 https://github.com/alimtvnetwork/gitmap-v28 -y

# merge with explicit policy
gitmap-v28 merge-both ./gitmap-v28 https://github.com/alimtvnetwork/gitmap-v28 \
         -y --prefer-newer

# pin remote branch
gitmap-v28 merge-right ./local https://github.com/owner/repo:develop

# preview without writing anything
gitmap-v28 mv ./gitmap-v28 ./gitmap-v28 --dry-run
```

---

## Logging

Every command emits structured `[mv]` / `[merge-both]` /
`[merge-left]` / `[merge-right]` prefixed log lines:

```
  [mv] resolving LEFT  : ./gitmap-v28 (folder, exists)
  [mv] resolving RIGHT : https://github.com/alimtvnetwork/gitmap-v28
  [mv]   -> mapped to working folder: ./gitmap-v28
  [mv]   -> folder does not exist; cloning
  [mv]   -> clone OK (47 files, 1.2 MB)
  [mv] copying files LEFT -> RIGHT (excluding .git/) ...
  [mv]   copied 47 files
  [mv] deleting LEFT (./gitmap-v28) ...
  [mv]   deleted
  [mv] committing in RIGHT ...
  [mv]   commit 9a3c1e2 "gitmap-v28 mv from ./gitmap-v28"
  [mv] pushing RIGHT ...
  [mv]   push OK
  [mv] done
```

Per-file conflict resolutions in `merge-*` are logged as
`[merge-both]   conflict docs/x.md -> took LEFT`.

---

## Constraints

- The `.git/` folder is **never** copied, compared, or deleted by
  the file operation (only by the post-step `rm -rf LEFT` in `mv`).
- The CLI MUST refuse to act if LEFT and RIGHT resolve to the same
  working folder (after URL→folder mapping). Error:
  `error: LEFT and RIGHT resolve to the same folder: <path>`.
- The CLI MUST refuse to act if RIGHT is a strict ancestor or
  descendant of LEFT on disk (would cause infinite recursion or
  self-overwrite). Error: `error: RIGHT is nested inside LEFT (or
  vice versa)`.
- Symlinks inside LEFT are copied as symlinks (not followed).
- File mode bits are preserved (executable bit on Unix).
- Empty directories on the source are recreated on the destination.
- The default ignore list excludes `.git/`, `node_modules/`,
  `.gitmap/release-assets/`. Override with `--include-vcs` /
  `--include-node-modules`.

---

## Acceptance Checklist

- [x] `gitmap-v28 mv folder-a folder-b` moves contents and deletes folder-a.
- [x] `gitmap-v28 mv folder url` clones url, copies files, commits + pushes.
- [x] `gitmap-v28 mv url folder` clones url, copies files into folder, deletes the cloned working folder.
- [x] `gitmap-v28 mv url-a url-b` clones both, copies, pushes RIGHT, deletes LEFT clone.
- [x] `merge-both` copies missing files both ways and prompts on conflicts.
- [x] `merge-left` only writes into LEFT.
- [x] `merge-right` only writes into RIGHT.
- [x] Conflict prompt accepts L/R/S/A/B/Q with the documented effects.
- [x] `-y` / `-a` bypasses prompts using the per-command source-side default.
- [x] `--prefer-left` / `--prefer-right` / `--prefer-newer` / `--prefer-skip` override the bypass default.
- [x] `--dry-run` prints all actions and writes nothing.
- [x] `--no-push` and `--no-commit` are honoured for URL endpoints.
- [x] `:branch` suffix on a URL pins the checkout + push branch.
- [x] Same-folder and nested-folder protection trips before any file write.

> **Implementation:** v2.98.0 — `gitmap-v28/movemerge/` package, `cmd/move.go`, `cmd/merge.go`, `cmd/movemergeflags.go`, `cmd/dispatchmovemerge.go` wired in `cmd/root.go`. Helptext lives in `gitmap-v28/helptext/{mv,merge-both,merge-left,merge-right}.md` (added in v2.96.0).
