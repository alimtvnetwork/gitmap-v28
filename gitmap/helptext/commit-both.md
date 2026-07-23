# gitmap commit-both

> **Status (v3.104.0):** sequential default + optional `--interleave`
> author-date variant. Sequential gives deterministic, auditable
> per-side summaries; interleave is closer to "what actually happened
> first" but aborts mid-stream on first error.

Bidirectional commit replay: each side ends up with the union of both
sides' commit timelines, applied either in two ordered passes
(default) or in author-date-merged order (`--interleave`).

## Alias

cmb

> Spec §13 reserved `cb`. `cb` is currently free, but the family uses
> `cmb` for visual consistency with `cml` / `cmr`. The long-form
> `commit-both` always works.

## Usage

    gitmap commit-both LEFT RIGHT [flags]
    gitmap commit-both LEFT RIGHT --interleave

## Algorithm — sequential (default)

1. **Pass 1 — LEFT → RIGHT.** Build plan from LEFT, preview, prompt
   (unless `-y` / `--dry-run`), replay onto RIGHT, push.
2. **Pass 2 — RIGHT → LEFT.** Now that RIGHT carries LEFT's commits
   too, build a fresh plan from RIGHT (so LEFT's just-replayed commits
   are excluded by the merge-base), preview, prompt, replay onto LEFT,
   push.
3. If Pass 1 fails the run aborts before Pass 2 — partial commit-both
   is worse than half-done because the second direction's merge-base
   would have shifted.

Each pass labels its log lines with a directional suffix
(`(left→right)` / `(right→left)`) so commit-both output is
visually attributable.

## Algorithm — `--interleave` (v3.104.0+)

1. Build BOTH directional plans up front (two `BuildPlan` calls).
2. Merge the two commit lists into a **single stream sorted by
   AuthorAt** (stable sort; LEFT-side wins on exact ties).
3. Print the unified plan with each step labelled `L→R` or `R→L`.
4. Single confirmation prompt (unless `-y` / `--dry-run`).
5. Walk the stream and replay each commit onto its **opposite** side
   in chronological order.
6. After the stream finishes, push each side that received commits
   and print a per-side summary.

Tradeoffs vs sequential:

- Faithful to original interleaved-history intent.
- One prompt instead of two.
- First per-commit failure aborts mid-stream — leaves whichever side
  was being written in a partial state. Use `--dry-run` to audit first.
- No per-direction merge-base re-computation between commits, so a
  commit replayed mid-stream might re-touch files that the just-prior
  opposite-direction commit had also touched.

## Flags

Same set as `commit-right` (see [commit-right.md](commit-right.md)),
plus:

| Flag           | Effect                                            |
|----------------|---------------------------------------------------|
| `--interleave` | Switch from sequential to author-date stream      |

`--interleave` is only valid for `commit-both`. Passing it to
`commit-left` or `commit-right` exits with code 2.

## Examples

    # Sequential (default)
    gitmap commit-both ./repo-A ./repo-B

    # Author-date interleave with dry-run audit first
    gitmap commit-both ./repo-A ./repo-B --interleave --dry-run

Sequential output skeleton:

    [commit-both] (left→right) replaying 3 commits from ./repo-A onto ./repo-B
    [commit-both] (left→right) [1/3] a3f2c1d  feat: add OAuth flow
    ...
    [commit-both] (right→left) replaying 2 commits from ./repo-B onto ./repo-A
    ...

Interleave output skeleton:

    [commit-both] interleave plan: 5 commits in author-date order
    [commit-both] [1/5] L→R  a3f2c1d  feat: add OAuth flow
    [commit-both] [2/5] R→L  b7e4a9f  fix: typo
    [commit-both] [3/5] L→R  c44d1ac  refactor: extract handler
    [commit-both] [4/5] R→L  d9e2510  docs: update README
    [commit-both] [5/5] L→R  e1afb20  test: add coverage

## See Also

- [commit-left](commit-left.md), [commit-right](commit-right.md) — single-direction siblings
- [merge-both](merge-both.md) — file-state mirror (no commit replay)
- spec/01-app/106-commit-left-right-both.md — full design

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter commit-both
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
