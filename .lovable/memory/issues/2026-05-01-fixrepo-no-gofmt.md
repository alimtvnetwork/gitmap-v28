# Root Cause Analysis — fix-repo / fix-repo.ps1 leave Go files un-gofmt'd

**Date:** 2026-05-01
**Reporter:** user
**Severity:** CI-blocker (recurring)
**Affected versions:** all `fix-repo` releases up to and including v4.6.0

## Symptom

After running `gitmap fix-repo --all` (or the legacy `fix-repo.ps1`) to bump
`{base}-vN` → `{base}-vN+1`, CI fails the gofmt check on `_test.go` files
that contain column-aligned map literals. Most recent failure (v4.6.0):

```
The following .go files are not gofmt-clean:
  cmd/replaceversionparse_test.go
  constants/cmd_constants_test.go
```

## Root Cause

`fix-repo` is a **token rewriter**. It edits the bytes of every matching
`{base}-vN` (and `{base}/vN`) occurrence in place. It does **not** re-run
`gofmt` on `.go` files after the rewrite.

Go's `gofmt` aligns map-literal keys/values into columns based on the
**widest key** in each contiguous block. When fix-repo changes the width of
even a single key — e.g. `gitmap-v28` (10 chars) → `gitmap-v28` (10 chars,
same width, fine) but `gitmap-v28` (9) → `gitmap-v28` (10, **wider**) — the
spacing of every other line in that block becomes one space short, and
`gofmt -l` flags the whole file.

The same class of breakage happens whenever a human edits a constants block
(e.g. adding `CmdCloneFixRepoPubAlias`, the longest key in its block) and
forgets to re-pad the older shorter keys. fix-repo just makes it
deterministic and frequent.

`fix-repo.ps1` has the same gap — it's a regex `-replace` over file bytes
with no language-aware post-processing.

## Why It Keeps Hitting Us

1. fix-repo touches every release (it's the whole point of the command).
2. Map-literal alignment is **block-scoped** — a single token change can
   silently shift every other line in the block.
3. Local devs running `go test ./...` don't see it; only the CI gofmt
   gate (`gofmt -l .`) does.
4. The fix-repo unit tests cover token correctness, not file formatting,
   so the regression is invisible until CI.

## Fix (this commit)

Manual gofmt-equivalent edits applied to the two failing files:

- `gitmap/cmd/replaceversionparse_test.go` — re-aligned both map literals
  in `TestSlugFromRemote` and `TestRemoteSlugRegex`.
- `gitmap/constants/cmd_constants_test.go` — re-aligned the
  Audit/inject/replace/regoldens/templates block so every key in the
  contiguous group pads to the width of the longest one
  (`CmdCloneFixRepoPubAlias`, 23 chars).

## Permanent Fix (next release — tracked separately)

`fix-repo` and `fix-repo.ps1` MUST run `gofmt -w` over every `.go` file
they actually rewrote, immediately after the rewrite phase, before the
"changed N files" summary.

Acceptance:
- After `gitmap fix-repo --all`, `gofmt -l .` returns empty on a clean
  worktree.
- An end-to-end test in `gitmap/cmd/fixrepo_test.go` writes a
  column-aligned map literal containing `{base}-v9`, runs fix-repo to
  bump to `v12`, and asserts `gofmt -l` returns empty.
- `fix-repo.ps1` shells out to `gofmt -w <file>` for each `.go` file it
  rewrote, behind a `Get-Command gofmt -ErrorAction SilentlyContinue`
  guard so it no-ops cleanly when Go isn't installed (and prints a
  warning telling the user to run `gofmt -w .` themselves).

## Related Memory

- Core rule: `Go v1.24.13. golangci-lint pinned to v1.64.8`
- mem://features/fix-repo-command — must be updated to mention
  the post-rewrite gofmt step once implemented.

## Follow-up (v6.80.1)

Diagnostics and tuning shipped to complement the chunker:

- `gitmap doctor fix-repo` — probes gofmt on PATH, gofmt executability,
  measured Windows argv cap vs configured budget, and a chunker
  invariant self-test. `--json` and `--budget N` supported.
- `gitmap fix-repo --dry-run` now prints a per-batch table with
  cmdLen and percent-of-budget, tagging NEAR-LIMIT (≥90%) and
  OVER-LIMIT (≥100%) batches.
- `gitmap fix-repo --verbose` prints batch header, per-batch
  start/done lines, and a rolling ETA.
- `--gofmt-max-cmd-len N` overrides the compiled-in 30,000-char
  default. Floor 512.

Spec: `spec/01-app/118-fix-repo-gofmt-tuning.md`.
