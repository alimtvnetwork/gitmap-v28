# Root Cause Analysis — fix-repo silently desynced regex-capture expectations

**Date:** 2026-05-01
**Reporter:** user
**Severity:** test failure (CI-blocker), data-corruption class
**Affected versions:** every release where fix-repo bumped a version
crossing a digit-width boundary (v9→v10, v99→v100, …)

## Symptom

Test failure in `gitmap/cmd`:

```
TestRemoteSlugRegex
    regex "gitmap-v28" -> base="gitmap" num="12", want base="gitmap" num="9"
```

The map key `"gitmap-v28"` was expecting captured num `"9"` — internally
inconsistent.

## Root Cause

`fix-repo` (Go binary AND `fix-repo.ps1`) is a **literal-token rewriter**.
It scans every tracked file for the pattern `{base}-v{N}` (with negative
lookahead `(?!\d)`) and rewrites the matched substring to
`{base}-v{current}`. That is the entire contract.

But test fixtures often store the **digit alone**, separated from the
`-v` prefix, as the expected output of a regex capture:

```go
"gitmap-v28": {true, "gitmap", "9"},
//          input → matches?, expected base, expected num
```

When fix-repo bumped this repo from v9 → v12, it correctly rewrote
the **map key** (`"gitmap-v28"` matches `{base}-vN` → becomes
`"gitmap-v28"`) but **did NOT touch the bare `"9"` in the value**,
because `"9"` alone doesn't match the `{base}-vN` token shape. The
test now asserts that the regex extracting the version number from
`"gitmap-v28"` returns `"9"`, which is internally contradictory.

The test stayed green for years because v1 → v2, v2 → v3, … v8 → v9
all had the same digit width as the captured num. The first
**width-changing bump** (v9 → v10, then v9 → v12 in subsequent fixes)
exposed the latent inconsistency.

`fix-repo.ps1` has the same gap — same regex, same blindness to
non-`-v`-prefixed digits.

## Why It Keeps Hitting Us

1. Test fixtures intentionally pair a `{base}-vN` input with the
   captured `N` as a separate string — that's the contract under test.
2. fix-repo operates token-by-token; it has no semantic understanding
   of "this string `"9"` over here is meant to mirror the `9` in the
   key over there".
3. The breakage is **silent until a digit-width crossing**: every
   v1→v2 style bump preserves the bug invisibly.
4. Same root mechanism as the gofmt-alignment bug (also produced by
   fix-repo's lack of post-rewrite normalization), but worse because
   gofmt failures are **mechanical** (column padding) while this one
   is **semantic** (test now lies about what the regex does).

## Fix (this commit)

`gitmap/cmd/replaceversionparse_test.go` line 37:
`{true, "gitmap", "9"}` → `{true, "gitmap", "12"}`.

## Permanent Fix (next release — tracked separately)

This class of bug cannot be solved by fix-repo alone — fix-repo has
no way to know that `"9"` is semantically tied to `"gitmap-v28"`. Two
defenses are needed in combination:

1. **Test convention**: when a test pairs a `{base}-vN` key with a
   captured-version expectation, derive the expectation from the key
   at runtime instead of hard-coding it. e.g.:

   ```go
   for in, want := range cases {
       parts := remoteSlugRe.FindStringSubmatch(in)
       if !want.matches { … continue }
       wantNum := strings.TrimPrefix(strings.Split(in, "-v")[1], "")
       if parts[2] != wantNum { … }
   }
   ```

   This makes the test self-consistent under any version bump.

2. **fix-repo audit pass**: add a `--strict-go-tests` mode that, after
   rewriting, runs `go test -run TestRemote -count=1 -short ./...` on
   any package whose `_test.go` files were touched, and reports any
   failures. This catches regressions caused by fix-repo even when
   the test convention from (1) wasn't applied.

Acceptance:
- After `gitmap fix-repo --all`, `go test ./...` is green on the
  gitmap repo without manual edits.
- A regression test in `gitmap/cmd/fixrepo_test.go` writes a fixture
  with a `{base}-v9` key paired with `"9"`, runs fix-repo to bump to
  v12, and asserts the test it just modified still compiles AND
  passes — failing the harness loudly when (1) is violated.

## Related Memory

- mem://features/fix-repo-command — already lists the gofmt gap; add
  this digit-capture gap once a permanent fix lands.
- Core rule: `Strict code style: <200 lines/file, <15 lines/func` —
  unchanged.

## Companion Issue

- `.lovable/memory/issues/2026-05-01-fixrepo-no-gofmt.md` — sibling
  bug where fix-repo leaves Go files un-gofmt'd. Same root mechanism
  (token rewriter without post-pass normalization), different symptom.
