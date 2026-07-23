# fix-repo paired-literal desync (sibling digit not rewritten)

**Date:** 2026-05-02
**Severity:** High — silent test corruption after every version bump
**Status:** Test fixtures fixed; fix-repo audit hardened (PowerShell + shell)
**Related:** v4.12.0 digit-capture gap (issues/2026-05-01-fixrepo-digit-capture-desync.md) — this is the **second variant** of the same root class.

## Symptoms

After bumping v12 → v13 and running `fix-repo --all`:

```
--- FAIL: TestBuildAuditNeedlesWidthCrossing
    replaceaudit_test.go:63: needle[6] = "gitmap-v28", want "gitmap-v28"
--- FAIL: TestRemoteSlugRegex
    replaceversionparse_test.go:57: regex "gitmap-v28" -> base="gitmap" num="13", want base="gitmap" num="12"
```

Both failures present the same shape: the `{base}-vN` token half of a paired literal got rewritten, but the **bare digit half** (`"12"` in a sibling slice element / `12` in a `[]int` literal) did not.

## Root Cause

`fix-repo` only rewrites `{base}-vN` tokens (and the slash form `{base}/vN`). It is intentionally narrow — it cannot rewrite arbitrary integer literals because doing so would produce thousands of false positives in normal code (port numbers, indexes, bit widths…).

But test fixtures in `gitmap/cmd/replaceaudit_test.go` and `gitmap/cmd/replaceversionparse_test.go` paired a `{base}-vN` string with a sibling digit literal that conceptually represents the **same** version:

```go
// replaceversionparse_test.go (BROKEN)
"gitmap-v28": {true, "gitmap", "12"},   // ← "12" was the captured num for v12
```

```go
// replaceaudit_test.go (BROKEN)
got := buildAuditNeedles("gitmap", []int{8, 9, 10, 12})  // ← 12 stays
want := []string{
    ...
    "gitmap-v28", "gitmap/v12",  // ← only "gitmap-v28" got bumped
}
```

When fix-repo rewrites `gitmap-v28` → `gitmap-v28`, the sibling `"12"` / `12` remains, leaving an internally inconsistent test that fails on the next `go test` run.

This is the **same root class** as the v4.12.0 fix:
- v4.12.0 fix was for a `[]int` digit capture next to a `{base}-vN` map key.
- This 2026-05-02 fix is for the **mirror case**: a stringified digit (`"12"`) in a struct value, and `[]int` literals seeded with the previous `current`.

## The standing rule (now enforced)

> Any test that pairs a `{base}-vN` string literal with a sibling integer/string representing the same version MUST derive the sibling from the same `int` it formats the token from. Hard-coding both is forbidden.

Examples of correct patterns:

```go
// CORRECT — both halves derived from the same int
const current = 12
slug := fmt.Sprintf("gitmap-v%d", current+1)
want := struct{ ok bool; base, num string }{true, "gitmap", strconv.Itoa(current)}
cases := map[string]want{slug: want}
```

```go
// CORRECT — single source of truth
targets := []int{8, 9, 10, 12}
needleWants := []string{}
for _, t := range targets {
    needleWants = append(needleWants,
        fmt.Sprintf("gitmap-v%d", t),
        fmt.Sprintf("gitmap/v%d", t))
}
```

## Audit hardening (fix-repo PowerShell + shell)

`scripts/fix-repo/Rewrite-Engine.ps1` and `scripts/fix-repo/rewrite.sh` now run a **post-rewrite paired-literal audit** on every modified `*_test.go` file:

For every line that contains a freshly-bumped `{base}-v{Current}` token, the audit greps the **same line plus the next 2 lines** for a quoted-string or bare-int literal matching `^|[^v0-9]({Current-1})($|[^0-9])`. A match is reported with file:line and the audit exits non-zero (exit 10 = `FixRepoExitPairedLiteral`).

This catches the desync in CI before it lands. The audit is gated on `*_test.go` only — production code legitimately uses small integers next to module paths.

## Files touched (this fix)

- `gitmap/cmd/replaceversionparse_test.go` — derive `num` from the same int as the `-vN` token.
- `gitmap/cmd/replaceaudit_test.go` — derive needle list from `targets` slice via `fmt.Sprintf`.
- `scripts/fix-repo/Rewrite-Engine.ps1` — added `Invoke-PairedLiteralAudit`.
- `scripts/fix-repo/rewrite.sh` — added `paired_literal_audit`.
- `gitmap/constants/constants_fixrepo.go` — `FixRepoExitPairedLiteral = 10` (if not already defined).

## Why the Go-native rewriter was not the right place

The Go-native rewriter (`gitmap fix-repo`) already has `--strict` mode which runs `go test` on touched packages. That **does** catch this failure — strict mode is the long-term safety net. But:

1. Many users / CI flows still call the PowerShell / shell scripts directly (legacy entry points kept for cross-platform parity).
2. `--strict` runs `go test` which is expensive; the paired-literal audit is a 2-line regex grep that runs in microseconds.
3. The audit gives a precise diagnostic (`paired-literal desync at file:line: 'gitmap-v28' followed by sibling '12'`) instead of a downstream test failure.

So we fix it at both layers: strict-mode tests + cheap paired-literal audit.
