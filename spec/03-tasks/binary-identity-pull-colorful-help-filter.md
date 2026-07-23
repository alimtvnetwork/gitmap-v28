# Task: Binary identity footer, `gitmap pull`, colorful help, filter highlight

Five-step plan. Each step is executed on the user's `next` signal.

---

## Step 1 — Embed gitmap binary identity + render dual-identity footer

**Goal.** Every command's tail prints TWO clearly separated identity
blocks so the user always knows (a) which gitmap build they're using
and (b) which repo they're currently inside.

- Extend `gitmap/constants/constants.go` with build-time vars:
  `BuildCommit`, `BuildBranch`, `BuildRepo`, `BuildDate` (already have
  `Version`, `RepoPath`).
- Wire `-ldflags "-X ..."` injection in:
  - root `Makefile` build target (if present)
  - `scripts/build*.ps1` / `scripts/build*.sh`
  - GitHub Actions `cross-platform.yml` build step
  Fallback: when unset, read from embedded `.gitmap/release/latest.json`
  at runtime so dev builds still show *something*.
- New package `gitmap/identity/`:
  - `Self()` → struct {Version, Repo, Branch, Commit, CommitShort, BuiltAt}
  - `Current(cwd)` → same struct, sourced from `git rev-parse` in cwd
- New renderer `gitmap/render/identityfooter.go`:
  - Prints `gitmap` block (magenta header) then a **blank line + thin
    rule** then `current repo` block (cyan header).
  - Reuses the colour palette from `gitmap/glyphs` / `gitmap/theme`.
- Hook into the post-command footer path (same place the existing
  version/repo/branch lines from the screenshot are emitted).

**Deliverable.** Footer shows both blocks with visual distance; works
on a release build and on `go run .`.

---

## Step 2 — `gitmap pull` with `--ssh` / `--pub`

**Goal.** From any git repo, `gitmap pull` pulls the current repo
after rewriting the `origin` URL to SSH or HTTPS on demand.

- New `gitmap/cmd/pull.go` + `gitmap/pullcmd/` package.
- Flags:
  - `--ssh`   → rewrite origin to `git@host:owner/repo.git`
  - `--pub`   → rewrite origin to `https://host/owner/repo.git`
    (alias `--https`)
  - bare `gitmap pull` → no rewrite, just `git pull --ff-only` with
    safe-pull retry policy from `constants.SafePullRetry*`.
- URL rewriter reuses `gitmap/clonefrom/summary_scheme.go` classifier
  + a new `pullcmd/rewriteurl.go` mapper (full unit-test coverage).
- Register alias `pl`. Add `// gitmap:cmd top-level` marker.
- Help file `gitmap/helptext/pull.md` with examples.

**Deliverable.** `gitmap pull --pub` and `gitmap pull --ssh` round-trip
correctly on a test repo. Existing `safepull` helper reused.

---

## Step 3 — Colourful help: flags + command shorthand

**Goal.** `gitmap help` and `gitmap help <cmd>` become scannable at a
glance.

- Extend `gitmap/render/prettypost.go` (added in earlier turn):
  - Flag tokens (`--foo`, `-f`) → cyan/bold.
  - Default values (`(default: ...)`) → dim.
  - In top-level command list, **shorthand aliases in parens** (e.g.
    `clone (cl)`) → yellow/bold, command name → green/bold.
  - Section headers (`## Flags`, `## Examples`) → magenta underline.
- Snapshot tests under `gitmap/render/prettypost_test.go` pinning the
  ANSI sequences for: a flag line, an alias line, a header line.
- No content changes — purely the post-processor's regex table grows.

**Deliverable.** `gitmap help clone` visibly differentiates flags,
aliases, headers, defaults.

---

## Step 4 — `--filter <q>` everywhere + end-of-help highlight banner

**Goal.** Users can `gitmap help --filter ssh` to grep help, and the
matched section is **highlighted again at the very bottom** so it's
visible without scrolling back up.

- Confirm/extend the existing `--filter` (v5.43.0) on:
  - `gitmap help`            — already present, verify
  - `gitmap clone`           — filter the row table before clone
  - `gitmap list-*` family   — wire through `prettypost`
- New `gitmap/render/filterhighlight.go`:
  - Wraps matched substrings in inverse-video ANSI.
  - Appends a `── matches for "<q>" ──` block at end with the matched
    lines repeated (max 10), so the user sees them after scrolling.
- Unit tests for filter+highlight composition.

**Deliverable.** `gitmap help --filter ssh` highlights inline AND
prints a summary block at the bottom.

---

## Step 5 — Version bump, CHANGELOG, README pin

**Goal.** Ship the four feature steps as `v5.58.0`.

- `gitmap/constants/constants.go` → `Version = "5.58.0"`
- `src/constants/index.ts` → `VERSION = "v5.58.0"`
- `CHANGELOG.md` → new top entry covering steps 1–4.
- `src/data/changelog.ts` → matching entry.
- `README.md` → version matrix + pinned-version block updated.
- Verify `version-sync.test.ts` passes.

**Deliverable.** Single self-contained release commit, ready for
`gitmap pr v5.58.0`.

---

## Execution protocol

User says `next` → execute the next un-done step, then list remaining
steps. No step is started ahead of the signal.
