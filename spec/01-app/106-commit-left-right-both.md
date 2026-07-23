# 106 — `commit-left` / `commit-right` / `commit-both`

**Status:** IMPLEMENTED — `commit-right` shipped in v3.76.0; `commit-left` and `commit-both` shipped in v3.102.0 (Phases 2 + 3) reusing the same Plan/Replay primitives via `committransfer.runOneDirection`.
**Companion family:** `mv` / `merge-both` / `merge-left` / `merge-right` (see `97-move-and-merge.md`)
**Related:** `24-amend-author.md` (commit metadata rewriting), `61-refactor-autocommit.md` (auto-commit primitives)

## 1. Purpose

The merge-* family transfers **file state** between two repo endpoints. The
commit-* family transfers **the history of how that file state was reached** —
i.e. it replays the source side's commit timeline onto the target side as a
sequence of fresh commits, one per source commit, preserving the order,
authorship intent, and (cleaned) commit messages.

Use cases:

- "I prototyped in `repo-A`, now I want the same step-by-step history in
  `repo-B` without the prototype's noise commits."
- "Two forks diverged with parallel work; I want both forks to share a single
  unified, chronologically-ordered history."
- "I want to rebase a feature from `wp-onboarding-v12` onto
  `wp-onboarding-v13` as a series of clean conventional commits."

## 2. Command surface

| Command | Alias | Direction | Writes commits to |
|---|---|---|---|
| `gitmap commit-left LEFT RIGHT [flags]`  | `cl` | RIGHT → LEFT  | LEFT  |
| `gitmap commit-right LEFT RIGHT [flags]` | `cr` | LEFT → RIGHT  | RIGHT |
| `gitmap commit-both LEFT RIGHT [flags]`  | `cb` | bidirectional | LEFT and RIGHT |

> **Naming mirror:** `commit-left` writes to LEFT (source = RIGHT), exactly
> like `merge-left` writes files to LEFT. The "-left" suffix always names
> **the destination**, never the source. This is intentional — it keeps the
> mental model identical to the merge-* family.

LEFT and RIGHT use the same endpoint syntax as the merge-* family — either a
local folder path or an `https://` / `git@` URL with optional `:branch`
suffix. Resolution rules (clone-if-missing, `--force-folder` to replace a
folder whose origin doesn't match, etc.) are inherited verbatim from
`movemerge.Endpoint`.

## 3. Commit selection — "all commits since divergence"

For each (source, target) pair the command:

1. Resolves both endpoints into working folders (same code path as merge-*).
2. Runs `git merge-base <source-HEAD> <target-HEAD>` to find the divergence
   point. If no shared base exists (truly unrelated histories), the source's
   entire reachable history is used.
3. Lists commits with
   `git rev-list --reverse --no-merges <base>..<source-HEAD>`. The
   `--reverse` flag guarantees oldest-first replay; `--no-merges` is on by
   default but can be disabled with `--include-merges` (see flags).
4. Filters out commits whose cleaned message would be empty or matches the
   "drop" patterns (see §6).

The replay set is computed **once per direction** before any side effects. A
preview is printed up front:

```
[commit-right] replaying 7 commits from LEFT onto RIGHT:
  [1/7] a3f2c1d  feat: add OAuth flow
  [2/7] b7e4a9f  fix: handle expired tokens
  ...
[commit-right] proceed? [y/N]
```

`-y` / `--yes` skips the prompt.

## 4. Replay mechanism — manual reconstruct

Per design decision, replay does **not** use `git cherry-pick` or
`format-patch`/`am`. Instead, each source commit becomes:

1. `git checkout <commit-sha>` (detached HEAD) on the source working dir.
2. File-level snapshot copied into the target working dir using the existing
   `movemerge` copy primitives (respecting `--include-vcs` /
   `--include-node-modules`).
3. `git add -A && git commit -m "<cleaned-message>"` on the target side.
4. Source working dir is restored to its original branch HEAD when the
   replay loop finishes (or aborts).

Rationale: this isolates the operation from cross-history concerns
(cherry-pick conflicts, divergent ancestors, sign-off tags). The cost is that
the **resulting tree on the target after each step is exactly the snapshot of
that source commit**, even if the target had unrelated files. Those unrelated
target-only files are preserved between commits — the manual reconstruct
copies *changed* files from source, it does not delete target-only files
unless `--mirror` is passed (see flags).

> **Important caveat (acknowledged by design):** because we are doing manual
> file-state reconstruction rather than git-native cherry-pick, the target's
> resulting tree at each step will not byte-match the source commit's tree
> when the target carries files that the source never had. This is accepted —
> the goal is "the same human-readable evolution," not a tree-hash-equivalent
> mirror.

## 5. `commit-both` — two sequential passes (default) or `--interleave`

> **Implementation note (v3.104.0):** sequential remains the default
> because it gives deterministic per-side summaries and avoids
> mid-run merge-base drift. The author-date interleave variant
> originally drafted in this section now ships behind `--interleave`.

### 5.1 Default — sequential

1. **Pass 1 — LEFT → RIGHT.** Build plan from LEFT, preview, prompt
   (unless `-y` / `--dry-run`), replay onto RIGHT, push.
2. **Pass 2 — RIGHT → LEFT.** Build a *fresh* plan from RIGHT (so
   LEFT's just-replayed commits are excluded by the new merge-base),
   preview, prompt, replay onto LEFT, push.
3. If Pass 1 fails the run aborts before Pass 2 — partial commit-both
   is worse than half-done because the second direction's merge-base
   would have shifted.

Each pass labels its log lines with a directional suffix
(`(left→right)` / `(right→left)`) so commit-both output is
visually attributable to a specific direction.

Edge cases:

- Pass 2's plan automatically excludes commits replayed by Pass 1
  because `git merge-base` will now point past them.
- If LEFT and RIGHT have a commit with the same cleaned message +
  same author date + same diff, the provenance footer's
  `AlreadyReplayed` check skips it on the second pass.

### 5.2 `--interleave` (v3.104.0+)

1. Build BOTH directional plans up front.
2. Merge the two commit lists into a single stream sorted by
   `AuthorAt` (stable sort; LEFT-side wins exact ties).
3. Print the unified plan with each step labelled `L→R` or `R→L`.
4. Single confirmation prompt.
5. Walk the stream and replay each commit onto its **opposite** side
   in chronological order.
6. After the stream, push each side that received commits and print
   a per-side summary.

Tradeoffs:

- ✓ Closer to "what actually happened first" across both sides.
- ✓ One prompt instead of two.
- ⚠ First per-commit failure aborts mid-stream — leaves whichever side
  was being written in a partial state. Use `--dry-run` first.
- ⚠ No per-direction merge-base re-computation between commits, so
  later interleaved commits may re-touch files that an opposite-direction
  commit just modified. Conflicts surface immediately as a failure.

`--interleave` is rejected (exit 2) for `commit-left` / `commit-right`.

Implementation: `committransfer.RunBothInterleaved` in
`gitmap/committransfer/interleave.go`. Sort invariant + tie-breaking
pinned by `TestBuildInterleavedStream*` in
`gitmap/committransfer/interleave_test.go`.

## 6. Commit-message normalization pipeline

Every replayed commit's message passes through this pipeline in order. Each
stage is independently configurable; defaults are conservative.

### 6.1 Drop filter (whole-commit skip)

A source commit is **skipped entirely** (not replayed) if its subject line
matches any pattern in the drop list. Defaults:

```
^Merge branch
^Merge pull request
^Revert "
^fixup!
^squash!
^WIP$
```

Skipped commits are reported in the summary but do not produce a target
commit. Override with `--no-drop-merges` etc. or the config knob below.

### 6.2 Strip rules (regex-based prefix/suffix removal)

Configurable via `gitmap/data/config.json`:

```json
{
  "commitTransfer": {
    "stripPatterns": [
      "^\\[WIP\\]\\s*",
      "^(JIRA|TICKET)-\\d+:\\s*",
      "\\s*\\(#\\d+\\)$"
    ],
    "dropPatterns": [
      "^Merge branch",
      "^fixup!",
      "^squash!"
    ],
    "conventionalCommit": true,
    "provenance": true
  }
}
```

CLI flags override config:

| Flag | Effect |
|---|---|
| `--strip <regex>` (repeatable) | Append a strip pattern for this run |
| `--no-strip` | Disable all strip patterns (config + flags) |
| `--drop <regex>` (repeatable) | Append a drop pattern for this run |
| `--no-drop` | Disable all drop patterns (replay every commit) |

Strip patterns run **before** conventional-commit normalization so that
"`[WIP] feat: add login`" → "`feat: add login`" → kept as-is by §6.3.

### 6.3 Conventional-commit normalization

When enabled (`conventionalCommit: true` in config, or `--conventional` on
the CLI; disable per-run with `--no-conventional`), the cleaned subject is
inspected:

- If it already starts with `<type>(<scope>)?: ` (one of `feat`, `fix`,
  `chore`, `docs`, `refactor`, `test`, `build`, `ci`, `perf`, `style`,
  `revert`), it is kept verbatim.
- Otherwise the diff is inspected to infer a type:
  - Only `*.md`, `docs/**` changes → `docs:`
  - Only `*_test.go`, `*.test.ts` etc. → `test:`
  - Only `Makefile`, `.github/**`, `.golangci.yml` → `ci:` or `build:`
    (CI files vs build files distinguished by path prefix)
  - Only formatting/whitespace deltas → `style:`
  - Default fallback → `chore:`
- Bug-fix heuristic: subject starts with "Fix", "Bugfix", "Hotfix"
  (case-insensitive) → `fix:`
- Feature heuristic: subject starts with "Add", "Introduce", "Implement"
  → `feat:`

The original subject (sans heuristic prefix) becomes the conventional-commit
subject after the colon. Body is preserved unchanged.

### 6.4 Provenance footer

When enabled (default `true`), every replayed commit gets a trailing footer
appended after a blank line:

```
gitmap-replay: from <source-display-name> <short-sha>
gitmap-replay-cmd: commit-right
gitmap-replay-at: 2026-04-21T19:55:00+08:00
```

Disable with `--no-provenance` or `"provenance": false` in config. The
footer is parseable so future tooling can detect/avoid double-replay
(`gitmap commit-right` should refuse to replay a commit that already carries
a `gitmap-replay:` footer pointing at the same source — see §10).

### 6.5 Empty-after-cleanup guard

If the cleaned message is empty (whole subject stripped), the commit is
skipped and reported in the summary as `cleaned-empty`.

## 7. Conflict handling

The manual-reconstruct mechanism does not produce git-merge conflicts at the
git layer (every step is a fresh commit on top of the target's HEAD).
However, **file-level** conflicts can still occur during the snapshot copy
when `--mirror` is **not** set and a target-only file would be overwritten
by a source-side file with the same path but different content.

The same `--prefer-left` / `--prefer-right` / `--prefer-newer` /
`--prefer-skip` policy used by merge-* applies here, with one rename:

- `--prefer-source` (alias `--prefer-from`) — source-side wins
- `--prefer-target` (alias `--prefer-to`) — target-side wins

The legacy `--prefer-left`/`--prefer-right` flags still work and resolve to
source/target based on command direction (e.g. `commit-right` →
`--prefer-left` == `--prefer-source`).

## 8. Flag reference

Inherits the full merge-* flag set unchanged, plus:

| Flag | Default | Description |
|---|---|---|
| `--mirror` | false | Delete target files not present in the source commit (true file-state mirror; closer to cherry-pick semantics) |
| `--include-merges` | false | Include `git rev-list` merge commits in the replay set |
| `--limit N` | 0 (no limit) | Replay at most N source commits (oldest first) |
| `--since <sha\|date>` | (auto: merge-base) | Override the divergence base |
| `--strip <regex>` | (config) | Add a strip pattern (repeatable) |
| `--no-strip` | false | Disable all strip patterns |
| `--drop <regex>` | (config) | Add a drop pattern (repeatable) |
| `--no-drop` | false | Replay every commit (disable drop filter) |
| `--conventional` | (config) | Force conventional-commit normalization on |
| `--no-conventional` | false | Disable conventional-commit normalization |
| `--provenance` | true | Append provenance footer |
| `--no-provenance` | false | Skip provenance footer |
| `--prefer-source` | false | Source side wins file conflicts |
| `--prefer-target` | false | Target side wins file conflicts |
| `--dry-run` | false | Print the full plan + cleaned messages; perform no writes |

## 9. Endpoint commit/push behavior

Inherits §5 of `97-move-and-merge.md` with one extension:

- After the replay loop completes successfully on a URL endpoint, gitmap
  pushes the *new commits* with `git push origin <branch>`.
- `--no-push` skips the push.
- `--no-commit` is honored at the **replay** layer too: when set, gitmap
  performs the file copies and stages them, but creates **no commits at
  all** on the target — useful for "dry-real" testing where you want the
  resulting working tree but want to inspect/edit before committing.

## 10. Idempotence and re-runs

The provenance footer (§6.4) is the idempotence anchor. Before replaying a
commit, gitmap:

1. Greps the target branch's last 200 commits (configurable) for
   `gitmap-replay: from <source-display-name> <short-sha>`.
2. If a match is found and `--force-replay` is **not** set, the commit is
   reported as `already-replayed` and skipped.
3. With `--force-replay`, the commit is replayed again and a new commit
   with a fresh timestamp is appended.

This makes `gitmap commit-right LEFT RIGHT` safe to re-run after adding a
few new commits to LEFT — only the genuinely-new commits will be appended
to RIGHT.

## 11. Failure modes and recovery

| Failure | Behavior |
|---|---|
| File-copy I/O error mid-replay | Abort with error; target HEAD is left at the last successfully-replayed commit. Re-running picks up where it left off (per §10). |
| `git commit` fails (e.g. nothing to commit because cleaned diff is empty) | Skip with reason `empty-diff`; continue to next source commit. |
| `git push` fails after successful replay | Replay is preserved; user gets the same `MsgAutoCommitSyncRetry`/abort path as `release/autocommitgit.go`. |
| User `Ctrl-C` during replay | Source working dir is restored to its original branch (deferred cleanup); target keeps any commits already made. |
| Source endpoint URL unreachable | Fail before any target writes (resolution happens up front). |

## 12. Output and logging

Log prefixes follow the merge-* pattern:

- `[commit-left]`, `[commit-right]`, `[commit-both]`

Per-commit lines:

```
[commit-right] [3/7] a3f2c1d → b91f4ce  feat: add OAuth flow
[commit-right] [4/7] b7e4a9f → c2d8e1a  fix: handle expired tokens  (was: "WIP fix tokens")
[commit-right] [5/7] e5fa12b → -        skipped: drop-pattern "^Merge branch"
```

Final summary:

```
[commit-right] done: replayed 5, skipped 2 (1 drop-pattern, 1 already-replayed)
[commit-right] pushed 5 commits to origin/main
```

## 13. Constants surface (planned)

A new `gitmap/constants/constants_committransfer.go` file will own:

- `CmdCommitLeft`, `CmdCommitLeftA` (`"commit-left"`, `"cl"`)
- `CmdCommitRight`, `CmdCommitRgtA` (`"commit-right"`, `"cr"`)
- `CmdCommitBoth`, `CmdCommitBothA` (`"commit-both"`, `"cb"`)
- `LogPrefixCommitLeft`, `LogPrefixCommitRight`, `LogPrefixCommitBoth`
- `FlagCTMirror`, `FlagCTIncludeMerges`, `FlagCTLimit`, `FlagCTSince`,
  `FlagCTStrip`, `FlagCTNoStrip`, `FlagCTDrop`, `FlagCTNoDrop`,
  `FlagCTConventional`, `FlagCTNoConventional`, `FlagCTProvenance`,
  `FlagCTNoProvenance`, `FlagCTPreferSource`, `FlagCTPreferTarget`,
  `FlagCTForceReplay`
- All user-facing message templates (`MsgCTReplayPlanFmt`,
  `MsgCTReplayStepFmt`, `MsgCTReplaySkipFmt`, `MsgCTSummaryFmt`,
  `ErrCTReplayFailedFmt`, `ErrCTSourceCheckoutFmt`, …)

Marked `// gitmap:cmd top-level` so the completion generator picks up
`commit-left`, `commit-right`, `commit-both`, `cl`, `cr`, `cb` automatically.

## 14. Package layout (planned)

Mirrors `movemerge/`:

```
gitmap/committransfer/
  types.go         # ReplaySpec, ReplayPlan, ReplayResult, MessagePolicy
  resolve.go       # endpoint resolution (delegates to movemerge.ResolveEndpoint)
  plan.go          # build replay set: merge-base → rev-list → drop filter
  message.go       # strip pipeline + conventional normalizer + provenance
  replay.go        # per-commit checkout/copy/commit loop
  conflict.go      # file-level conflict resolution (reuses movemerge.PreferPolicy)
  log.go           # [commit-*] prefixed structured output
  push.go          # final-push helper (reuses release/autocommitgit primitives)
```

Dispatcher: extend `gitmap/cmd/dispatchmovemerge.go` (or a new
`dispatchcommittransfer.go`) so `cmd/root.go` routes the three new commands.

## 15. Testing strategy

- Unit tests for the message pipeline with a table of (input, config) →
  expected cleaned message.
- Unit tests for the drop filter (default patterns + custom regex).
- Integration tests in `gitmap/tests/cmd_test/committransfer_test.go`
  using local-folder endpoints (no network):
  1. `commit-right` from a 5-commit source onto an empty target →
     verify 5 cleaned commits land in order.
  2. Re-run after appending 2 new source commits → verify only 2 new
     commits appended (idempotence via provenance footer).
  3. `commit-both` with 3 LEFT-only and 2 RIGHT-only commits → verify
     both sides end with the same 5-commit suffix in author-date order.
  4. Drop-pattern test: source containing a `Merge branch` commit →
     verify it's skipped and reported.
  5. Conventional-commit normalization: source commit "Add login form"
     → target commit subject "feat: Add login form".

## 16. Help text (planned files)

- `gitmap/helptext/commit-left.md`
- `gitmap/helptext/commit-right.md`
- `gitmap/helptext/commit-both.md`

Following the existing helptext template (Alias, Usage, Flags table,
Prerequisites, Examples, Exit Codes, Notes, See Also).

## 17. Out of scope (explicitly deferred)

- LLM-based commit-message rewriting (e.g. summarize a chatty WIP message
  into a clean conventional commit). The §6.3 heuristic stays rule-based.
- Cross-history rebase semantics (true cherry-pick with `-x`,
  three-way merges). The manual-reconstruct mechanism is intentionally
  simpler and lossier for this first pass.
- Signed commits / GPG signing of replayed commits. Inherits whatever
  `git commit` does in the target repo (so `commit.gpgsign=true` in the
  target's git config will Just Work, but gitmap doesn't manage keys).
- A `--squash` flag that collapses the entire replay set into one commit.
  Easy follow-up but not part of the v1 surface.

## 18. Implementation phasing (when this lands)

1. **Phase 1 — commit-right only.** Single direction is enough to validate
   the pipeline end-to-end. Ship constants, helptext, dispatcher,
   committransfer package, integration tests.
2. **Phase 2 — commit-left.** Trivial flip of source/target wiring once
   Phase 1 is solid.
3. **Phase 3 — commit-both.** Adds the interleave-by-timestamp planner
   and idempotence-across-both-sides handling.
4. **Phase 4 — config knobs.** Wire `commitTransfer` block in
   `gitmap/data/config.json` and the three-layer merge pattern (defaults
   → config → CLI flags) per the project's standard config-pattern.
