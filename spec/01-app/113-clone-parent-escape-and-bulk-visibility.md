# Spec 113 — Parent-escape clone family + bulk visibility + prior-version privatize

Status: Draft (v5.61.0 target)
Owner: cmd/clone, cmd/clonenext, cmd/visibility

## 1. Problem

Two related rough edges:

1. **Locked-cwd deadlock for clone family.** When the user is *already inside*
   the flattened target folder and runs `gitmap clone <url>`, `gitmap cfr ...`,
   `gitmap cfrp ...`, or `gitmap cn v++`, the existing-folder removal fails on
   Windows (cwd is held open) and gitmap silently degrades to a versioned
   sibling folder (`foo-v40/`). Only `gitmap cn v++ -f` currently escapes.

2. **Stale public versions after release.** `gitmap cfrp <vN-url>` flips
   `vN` public but leaves `vN-1`, `vN-2`, … publicly visible. The user
   has to chase them down one by one with `make-private`.

## 2. Scope

### 2.1 Parent-escape (auto, no flag)

For **every** command in the clone family — `clone`, `clone-fix-repo` (`cfr`),
`clone-fix-repo-pub` (`cfrp`), `clone-next` (`cn`) — when:

- the resolved target folder already exists, **and**
- the current working directory *is* that folder (or a subdirectory of it),

the command MUST:

1. `os.Chdir` to the parent directory.
2. `os.RemoveAll` the target.
3. Proceed with the clone into the flattened base-name folder.
4. `os.Chdir` into the freshly cloned folder.
5. Emit `GITMAP_SHELL_HANDOFF` so the shell wrapper follows.

This is the default behavior — no `-f` / `--force` required. The existing
`-f` flag on `cn` becomes a no-op alias (kept for back-compat) and stops
being the only escape hatch.

`cn v++` (and `v+N`) MUST always write into the base-name folder
(`foo/`), never `foo-vN/`. The current fallback that creates
`foo-v3/` when the flatten path is blocked is removed; the
parent-escape path replaces it.

### 2.2 Bulk visibility

New positional form for both commands:

```
gitmap make-public  [<repo-or-url>] [<count>] [--yes] [--dry-run] [--verbose]
gitmap make-private [<repo-or-url>] [<count>] [--dry-run] [--verbose]
```

- `<repo-or-url>` — optional. Accepts a full URL
  (`https://github.com/owner/foo-v40`), an SSH shorthand
  (`git@github.com:owner/foo-v40.git`), or a bare base name
  (`foo`). If omitted, defaults to the current repo's origin.
- `<count>` — optional positive integer. When present, the command
  resolves the *N most recent versions* of the base repo on the
  provider (descending `vN, vN-1, … vN-count+1`) and flips each one.
  When omitted, behavior is unchanged: only the single repo is
  flipped.
- A single `--yes` covers the whole batch.
- `--dry-run` lists every slug it *would* flip; performs no API call.
- Failures on individual repos are reported but do not abort the batch
  (exit code reflects the worst failure observed).

### 2.3 cfrp prior-version privatize

After `cfrp <url-for-vN>` successfully clones, fixes, and publishes
`vN`, it MUST:

1. Probe the provider for `vN-1`, `vN-2`, … (stop at first 404 or
   gap of 2 consecutive misses).
2. Filter to those currently `public`.
3. If any are public:
   - **No `-y`:** print the list and prompt
     `Make N prior version(s) private? [y/N]:` (default no).
   - **With `-y`:** auto-confirm — privatize all of them.
4. Apply via the same provider CLI path used by `make-private`.
5. Report each result on its own line; non-fatal on individual
   failures.

The `-y` flag on `cfrp` already short-circuits the `make-public`
confirmation; this spec extends it to the prior-version step too.

## 3. Out of scope

- Other providers beyond GitHub / GitLab.
- Cross-org bulk flips (count search stays within the same owner).
- Rewriting the `-f` flag semantics on `cn` — it remains accepted.

## 4. Exit codes

Existing visibility exit codes (2–8) unchanged. Bulk mode returns the
**worst** non-zero exit observed across the batch; `0` only when every
slug succeeded (or was already in the target state).

## 5. Tests (end-to-end, table-driven)

- `TestCloneParentEscape_{Clone,CFR,CFRP,CNvPlusPlus}` — cwd == target
  folder; assert chdir-up, remove, re-clone, post-clone chdir.
- `TestCNvPlusPlus_NoVersionedFallback` — blocked removal must not
  create `foo-vN+1/`.
- `TestMakeVisibilityBulk_DryRun` — `make-public foo 3` lists three
  slugs without invoking `gh`.
- `TestCFRP_PriorVersionPrompt_{NoY,WithY}` — stubs provider probe;
  asserts prompt suppression with `-y`.
- `TestMakeVisibilityBulk_PartialFailure` — middle slug fails;
  command still processes the rest and exits with that slug's code.

## 6. Docs

- Update `gitmap/helptext/{clone,clone-fix-repo,clone-fix-repo-pub,clone-next,make-public,make-private}.md`.
- Update root `README.md` Visibility + Clone sections.
- Add CHANGELOG.md entry under v5.61.0.
