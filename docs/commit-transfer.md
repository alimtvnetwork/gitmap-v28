# Commit Transfer Family — `commit-left` / `commit-right` / `commit-both`

> Replay one repo's commit history onto another as a fresh, cleaned, idempotent
> sequence. Spec: [`spec/01-app/106-commit-left-right-both.md`](../spec/01-app/106-commit-left-right-both.md).

## TL;DR

| Command | Alias | Direction | Status |
|---------|-------|-----------|--------|
| `gitmap commit-right` | `cmr` | LEFT → RIGHT (writes on RIGHT) | **LIVE** (v3.76.0+) |
| `gitmap commit-left`  | `cml` | RIGHT → LEFT (writes on LEFT)  | scaffold (Phase 2) |
| `gitmap commit-both`  | `cmb` | both, interleaved by author date | scaffold (Phase 3) |

> The suffix names the **destination**, exactly like `merge-left` writes to
> LEFT. Spec §13 reserved `cl` / `cr` / `cb`, but those are taken
> top-level (`cl` → changelog, `cr` → cpp-repos), so the disambiguated
> three-letter aliases `cml` / `cmr` / `cmb` are the canonical short forms.
> The long-form names always work.

## When to use this vs. `merge-*`

- **`merge-*`** = file-state mirror. Copies bytes between two folders. One
  resulting commit on the target. Loses source history.
- **`commit-*`** = commit-history replay. Walks every source commit since
  the divergence base and lands one cleaned commit per source commit on the
  target. Preserves authorship + dates; appends a `gitmap-replay:` footer
  so re-runs are idempotent.

If you want the target's `git log` to *look like* the source's log
(minus drops/strips), use `commit-*`. If you only care about the working
tree, use `merge-*`.

## Required inputs

```
gitmap <command> LEFT RIGHT [flags]
```

- **`LEFT` / `RIGHT`** — endpoints in any form the merge-* family accepts:
  - local path: `./repo-A`, `../upstream`
  - Git URL: `https://github.com/owner/name.git`, `git@github.com:owner/name.git`
  - alias: `my-api` (resolves via `gitmap as`)
- Both endpoints must be valid git working trees by the time the engine
  runs. URL endpoints are auto-cloned through the same resolver `merge-*`
  uses (`movemerge.ResolveEndpoint`).
- The source side (LEFT for `commit-right`) must have a clean working
  tree — replay does `git checkout` on it.

## How `commit-right` works (the live one)

1. **Resolve endpoints** → working dirs for LEFT and RIGHT.
2. **Build plan** (no writes):
   - `git merge-base LEFT/HEAD RIGHT/HEAD` (or `--since` override) → divergence base.
   - `git rev-list --reverse base..HEAD` on LEFT → the replay set (oldest first).
   - For each source commit: read subject/body/author, run the §6
     **message pipeline** (`drop` → `strip` → `conventional` → `provenance`),
     check the target's recent log for an existing `gitmap-replay:` footer
     pointing at the same source SHA → mark `already-replayed` if found.
3. **Print the plan** + ask for confirmation (skip with `-y`).
4. **Replay loop** — for each non-skipped commit:
   - `git checkout <sha>` on LEFT (detached HEAD)
   - file-snapshot copy LEFT → RIGHT (skips `.git/`, `node_modules/`,
     plus user `--strip` patterns; `--mirror` deletes target-only files)
   - `git add -A && git commit` on RIGHT with the cleaned message,
     preserving `GIT_AUTHOR_NAME` / `_EMAIL` / `_DATE` from the source
5. **Restore source HEAD** (also fires on Ctrl-C — see `signal.go`).
6. **Push** unless `--no-push`.

## Invocation cheatsheet

```bash
# Live: replay LEFT's commits onto RIGHT, with confirmation prompt.
gitmap commit-right ./repo-A ./repo-B
gitmap cmr         ./repo-A ./repo-B          # alias

# Skip prompt + dry-run (prints the plan, no writes).
gitmap cmr ./repo-A ./repo-B --dry-run -y

# Cap the replay set + override the divergence base.
gitmap cmr ./repo-A ./repo-B --limit 20 --since 2026-01-01

# Add a strip regex (repeatable) and disable provenance footer.
gitmap cmr ./A ./B --strip '\(#\d+\)$' --strip '^\[WIP\]\s*' --no-provenance

# True mirror (delete target-only files) + force re-replay over existing footers.
gitmap cmr ./A ./B --mirror --force-replay

# Snapshot + stage but don't commit (inspect before letting it commit).
gitmap cmr ./A ./B --no-commit

# Scaffold: prints "not yet implemented — see spec 106" and exits 2.
gitmap commit-left  ./A ./B
gitmap commit-both  ./A ./B
```

## Key flags (full table: `gitmap help commit-right`)

| Flag | Default | Effect |
|------|---------|--------|
| `--dry-run` | off | Print the plan; no writes. |
| `-y` / `--yes` | off | Skip the confirmation prompt. |
| `--limit N` | 0 (no cap) | Replay at most N source commits, oldest first. |
| `--since <ref\|date>` | merge-base | Override the divergence base. |
| `--mirror` | off | Delete target-only files during snapshot copy. |
| `--include-merges` | off | Include merge commits in the replay set. |
| `--strip <re>` | (none) | Regex stripped from the source subject (repeatable). |
| `--no-strip` | — | Disable all `--strip` patterns. |
| `--drop <re>` | spec §6.1 defaults | Skip commits whose subject matches (repeatable). |
| `--no-drop` | — | Replay every commit (disable drop filter). |
| `--conventional` / `--no-conventional` | on | Force `feat:`/`fix:`/`chore:` prefix. |
| `--provenance` / `--no-provenance` | on | Append `gitmap-replay:` footer (drives idempotence). |
| `--force-replay` | off | Replay even commits that already carry a `gitmap-replay:` footer. |
| `--no-commit` | off | Snapshot + stage but skip the commit. |
| `--no-push` | off | Stop after the local commit (skip `git push`). |

> **Negation toggles** (`--no-conventional`, `--no-provenance`, `--no-drop`)
> use Go 1.21 `BoolFunc` so they don't need a value, and order on the
> command line wins (last write to the flag is what the engine sees).

## Idempotence model

The provenance footer is the contract:

```
gitmap-replay: from repo-A a3f2c1d
gitmap-replay-cmd: commit-right
gitmap-replay-at: 2026-04-23T12:00:00Z
```

On every run, the planner reads the target's last 200 commit messages
and skips any source commit whose `(SourceDisplayName, ShortSHA)` pair
already appears in a footer. Use `--force-replay` to bypass this check
when you want to re-land cleaned-up versions of already-replayed commits.

## Interrupt safety

`Ctrl-C` (SIGINT) and SIGTERM during a `commit-right` run trigger a
best-effort `git checkout <originalRef>` on the source working dir
before the process exits with `128+signo`. You will not be left on a
detached HEAD pointing at some intermediate replay SHA. (See
`gitmap/committransfer/signal.go`.)

## See also

- `gitmap help commit-right` — full flag table + examples in the CLI
- [`merge-both` / `merge-left` / `merge-right`](../helptext/merge-both.md) — file-state mirrors (no commit replay)
- [`spec/01-app/106-commit-left-right-both.md`](../spec/01-app/106-commit-left-right-both.md) — full design doc, message pipeline (§6), and phasing plan (§18)
