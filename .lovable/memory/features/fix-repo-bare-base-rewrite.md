---
name: fix-repo bare-base rewrite
description: Bare `{base}` rewrite is restricted to the v1→v2 transition only (current==2). At v3+ the bare token is left alone — it is almost always an unrelated identifier (binary, package, brand)
type: feature
---

`gitmap fix-repo` (v5.8.0+) extends its rewrite engine so that when
bumping a `-v1` repo to `-v2`, the sweep ALSO substitutes standalone
`{base}` tokens — not just `{base}-v1`. This handles the case where
the original repo shipped without a `-v1` suffix (e.g.
`alimtvnetwork/img-pdf` instead of `img-pdf-v1`), leaving downstream
references reading the bare name.

**v5.38.0 scope tightening — bare-base ONLY runs when `current == 2`.**
Before v5.38.0 the sweep ran whenever v1 was in the target span, so
`fix-repo` inside a `-v4` repo would rewrite every bare `gitmap`
mention to `gitmap-v27`, silently corrupting binary names, package
identifiers, brand strings, and unrelated repo URLs. The pre-versioned
origin only matters for the v1→v2 transition; at v3+ any bare token is
overwhelmingly NOT the old origin and must be preserved.

**Guard rules** (`isBareBaseBoundary` in `gitmap/cmd/fixrepo_rewrite.go`):
prev byte AND next byte must NOT be a "word char" — defined as ASCII
alnum, `_`, `-`, or `.`.

The bare-base pass runs ONLY when `n == 1 && current == 2` inside
`applyAllTargets`.

Tests: `gitmap/cmd/fixrepo_rewrite_barebase_test.go` —
`TestApplyAllTargets_BareBase_SkippedAtV3Plus` and
`TestApplyAllTargets_BareBase_SkippedAtV4WithV1InTargets` lock in the
v5.38.0 scope.
