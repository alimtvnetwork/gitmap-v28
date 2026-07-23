# Spec 115 — v6.0.0 Migration Guide

Status: PUBLISHED (implemented v6.0.0)
Version target: v6.0.0
Owners: committransfer package, release management

## Summary

This spec documents the single breaking behavioural change planned for
v6.0.0 and provides the migration path for users and downstream scripts
that rely on the v5.x default.

## Breaking change

### `committransfer.IncludeMerges` default flips from `false` → `true`

In v5.x (introduced v5.62.0), `gitmap commit-in`, `commit-out`,
`commit-left`, `commit-right`, and `commit-both` silently strip merge
commits from the source range because `Options.IncludeMerges` defaults to
the zero-value `false`. The CLI gained `--include-merges` / `--no-include-merges`
flags in v5.62.0 as an opt-in escape hatch.

In v6.0.0 the default flips to `true`. Merge commits are now preserved by
default; users who want the legacy strip behaviour must pass
`--no-include-merges` explicitly.

## Why this is breaking

Scripts and CI jobs that invoke commit-transfer commands without flags
will see a different replay plan when the source range contains merge
commits:

| Scenario | v5.x result | v6.0.0 result |
|---|---|---|
| Source has merge commits, no flags | Merges stripped silently | Merges preserved |
| Source has no merge commits | No change | No change |
| `--include-merges` passed | Merges preserved | Merges preserved (redundant) |
| `--no-include-merges` passed | Merges stripped | Merges stripped |

## Migration path

### If you want the new behaviour (merges preserved)

No action required. Remove any `--include-merges` flags you added as a
v5.x workaround; they become redundant but harmless.

### If you want the legacy behaviour (merges stripped)

Add `--no-include-merges` to every commit-transfer invocation:

```bash
# Before (v5.x — merges stripped by default)
gitmap commit-in ./source ./target

# After (v6.0.0 — explicit opt-out required)
gitmap commit-in ./source ./target --no-include-merges
```

### If you use the Go API directly

Set `Options.IncludeMerges = false` explicitly:

```go
opts := committransfer.Options{
    IncludeMerges: false, // now required for legacy behaviour
}
```

## Stderr notice inversion

| Version | Notice emitted when … |
|---|---|
| v5.62.0+ | `--no-merges` strips ≥1 merge AND `--include-merges` was **not** passed. Message: "N merge commits excluded; pass --include-merges to preserve them." |
| v6.0.0+ | `--no-merges` strips ≥1 merge AND `--no-include-merges` **was** passed. Message: "N merge commits excluded by --no-include-merges." |

The notice moves from an advisory ("you can opt in") to an informational
confirmation ("you explicitly opted out").

## Files to touch for the v6.0.0 flip

- `gitmap/committransfer/types.go` — update `Options.IncludeMerges` doc
- `gitmap/cmd/committransfer.go` — flip the flag default wiring
- `gitmap/committransfer/plan.go` — verify `countMergeExcluded` logic still
  reports correctly under the new default
- `gitmap/helptext/commit-in.md`, `commit-out.md` — update help text
- `CHANGELOG.md` — breaking-change section
- `README.md` — commit-transfer command table if it mentions merge handling

## Acceptance

- [x] `go test ./gitmap/committransfer/...` green with the flipped default.
- [x] `TestPrintPlanNoticeV6` asserts correct notice under both `--include-merges` and `--no-include-merges`.
- [x] `TestCommitTransferIncludeMergesDefault` + `TestCommitTransferIncludeMergesExplicit` assert CLI wiring.
- [x] Manual smoke: run `gitmap commit-in` on a source with merge commits,
  verify merges appear in the plan by default, verify
  `--no-include-merges` strips them and emits the inverted stderr notice.
- [x] Changelog entry under "Breaking changes" links back to this spec.

## Out of scope

- Removing the `--no-include-merges` / `--include-merges` flags. Both
  remain supported indefinitely; only the default changes.
- Any other breaking changes not explicitly listed above. If a future
  breaking change is identified, it must be added to this spec before
  the v6.0.0 tag is cut.
