---
name: Fixture Stamp Autobump
description: GITMAP_FIXTURE_AUTOBUMP=1 + MustValidateBodyWithAutobump self-heal stale fixture-stamp markers; `make fixtures-bump RUN=...` runs the cycle and verifies
type: feature
---
# Fixture Stamp Autobump

`gitmap/fixtureversion` ships pure helpers (`BumpStampInBody`, `NextGeneration`, `ParseGenerationFromBody`, `FormatBumpSummary`) plus an env-gated file rewriter `MaybeAutoBumpFile`. Tests opt in by calling `MustValidateBodyWithAutobump(t, body, sourcePath, want)` instead of `MustValidateBody`.

Behavior:
- Without `GITMAP_FIXTURE_AUTOBUMP=1`: identical to `MustValidateBody` — t.Fatals on stale fixture with regenerate recipe.
- With `GITMAP_FIXTURE_AUTOBUMP=1`: rewrites the `// fixture-stamp:` marker's `generation=` (clamped to `want.MinGeneration`) inside `sourcePath` and lets the test pass for that run.

Gate is **strict equality** to `"1"` and lives only in `MaybeAutoBumpFile` — no other code path can mutate source.

Make targets:
- `make fixtures-bump RUN=<pattern> [PKG=./cmd/...]` — runs with `GITMAP_FIXTURE_AUTOBUMP=1`, then chains `fixtures-bump-verify` to re-run cleanly.
- `make fixtures-bump-verify RUN=<pattern>` — same pattern with the env var explicitly unset; mandatory second pass mirroring `goldens-verify`.

First consumer: `TestFixRepoRewriteV9ToV12Fixture` (rewrites `gitmap/cmd/fixrepo_rewrite_v9tov12_test.go`).

Files: `gitmap/fixtureversion/bump.go`, `bump_test.go`, `validate.go` (`MustValidateBodyWithAutobump`), `Makefile` (`fixtures-bump`, `fixtures-bump-verify`).
