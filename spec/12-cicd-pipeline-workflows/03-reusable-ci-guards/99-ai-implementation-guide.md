# AI Implementation Guide — Reusable CI Guards

**Version:** 1.0.0
**Updated:** 2026-04-21

---

## Purpose

This document tells an AI assistant **how to select, configure, and
ship** the six guards documented in this folder when asked to "add CI
guards to repository X". It assumes the AI has read this folder and
nothing else about the source repository that originated these
patterns.

---

## Decision Tree

```
REQUEST: "Add CI guards to <repo>"

1. What is the primary language?
   ─ Go        → consider all 6
   ─ Node/TS   → consider 1, 2, 4, 5, 6 (skip 3 if no flat namespace)
   ─ Python    → consider 1, 2, 4, 5, 6
   ─ Rust      → consider 1, 2, 4, 5, 6
   ─ Polyglot  → consider per-language combinations

2. Does the repo have a flat-namespace package?
   (Go cmd/ pkg, Node barrel exports, Python __init__.py re-exports,
    Rust mod with many file includes)
   ─ Yes  → install Pattern 01 (forbidden-name guard)
           install Pattern 03 (collision audit) on the flat package

3. Does the repo enforce a naming convention with thousands of
   pre-existing violations?
   ─ Yes  → install Pattern 02 (grandfather baseline)

4. Does the repo have a linter producing >50 pre-existing findings?
   ─ Yes  → install Pattern 04 (baseline-diff gate)
           install Pattern 05 (PR-comment suggestions) for UX

5. Does the test suite use a matrix of >2 shards?
   ─ Yes  → install Pattern 06 (test aggregator)
```

---

## Implementation Order

Always implement in this order — each builds on the previous:

1. **Shared scaffolding** — create `.github/scripts/` with a `readme.md`
   listing every guard's purpose and exit codes
2. **Test aggregator (Pattern 06)** — gives reviewers immediate value
3. **Baseline-diff lint gate (Pattern 04)** — biggest CI quality win
4. **Lint suggestions (Pattern 05)** — wired to the same JSON as 04
5. **Forbidden-name guard (Pattern 01)** — cheap to add once cmd-style
   namespace is identified
6. **Collision audit (Pattern 03)** — pairs with Pattern 01
7. **Grandfather baseline (Pattern 02)** — last, because it requires
   committing the baseline file alongside

---

## Per-Pattern Configuration Checklist

When implementing each guard, the AI MUST resolve every variable in
the pattern's "Configuration Surface" table by reading the target
repo. Do not invent defaults silently — if a value is unclear, ask
the user.

| Pattern | Critical Vars to Resolve First |
|---------|--------------------------------|
| 01 | `SOURCE_DIR` (which dir is flat-namespace?), `FILE_GLOB`, `FORBIDDEN_PATTERNS` |
| 02 | `SOURCE_DIR`, `ALLOWED_PREFIX_REGEX`, decide where `BASELINE_FILE` lives |
| 03 | `FILE_GLOB`, `IDENT_RE`, `RAWSTRING_DELIM` for the language |
| 04 | Linter JSON shape, normalizer function, baseline cache key |
| 05 | `SUGGESTERS` table — needs an entry per linter on the project |
| 06 | `PASS_MARKER` / `FAIL_MARKER` for the test framework |

---

## Workflow Skeleton

Wire the guards into `.github/workflows/ci.yml` like this:

```yaml
name: ci
on:
  push: { branches: [main] }
  pull_request: { branches: [main] }

permissions:
  contents: read
  pull-requests: write   # only if Pattern 05 posts comments

concurrency:
  group: ci-${{ github.ref }}
  cancel-in-progress: ${{ !startsWith(github.ref, 'refs/heads/release/') }}

jobs:
  guards:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6
      - name: Forbidden-name guard
        run: bash .github/scripts/check-forbidden-names.sh
      - name: Collision audit
        run: python3 .github/scripts/check-collisions.py
      - name: Naming baseline
        run: bash .github/scripts/check-naming.sh

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6
      - name: Restore lint baseline
        uses: actions/cache@v4
        with:
          path: .cache/lint-baseline.json
          key: lint-baseline-${{ github.event.repository.default_branch }}-<linter-version>
      - name: Run linter
        run: <linter> --output=json > current.json || true
      - name: Diff vs baseline
        run: python3 .github/scripts/lint-diff.py --current current.json --baseline .cache/lint-baseline.json
      - name: Suggestions
        if: always()
        run: python3 .github/scripts/lint-suggest.py --current current.json --baseline .cache/lint-baseline.json --repo ${{ github.repository }} --sha ${{ github.event.pull_request.head.sha || github.sha }}

  test:
    strategy:
      matrix: { shard: [01, 02, 03, 04] }
    steps: # ... run tests, upload test-results-${{ matrix.shard }} ...

  summarize:
    needs: test
    if: always()
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6
      - uses: actions/download-artifact@v4
        with:
          path: ./test-artifacts
          pattern: test-results-*
      - run: bash .github/scripts/test-summary.sh ./test-artifacts
```

Apply the [shared conventions](../01-shared-conventions.md):
exact-version pinning, `working-directory:` instead of `cd`,
least-privilege permissions.

---

## Common Mistakes to Avoid

| Mistake | Why It's Wrong | Correct Approach |
|---------|---------------|------------------|
| Use `@latest` for actions or tools | Breaks reproducibility silently | Pin every version per [shared conventions](../01-shared-conventions.md) |
| Use 3-arg `match()` in awk | Fails on mawk (default GH runner) with no error | Use 2-arg `match()` + `RSTART` / `RLENGTH` |
| Skip string-literal awareness in extractor | False positives from SQL keywords | Track raw-string + quote state per file |
| Sort baseline in non-C locale | `comm -23` produces phantom diffs | `LC_ALL=C sort -u` everywhere |
| Swallow errors in guard scripts | Defeats the gate | Surface every failure with `::error::` annotation |
| Hardcode repo paths in scripts | Not portable across forks | Accept paths via env vars or `${1:-default}` |
| Post stacking PR comments | Noisy review history | Use a unique HTML sentinel + create-or-update-comment |
| Cancel in-progress release builds | Leaves partial releases | `cancel-in-progress: false` for release branches |

---

## Testing the Guards

Before declaring an integration done, the AI MUST:

1. Run each guard locally against the repo to confirm exit `0` (or
   document expected violations and add them to the appropriate
   baseline)
2. Trigger a deliberate violation (e.g. add `func runOne(){}`) and
   confirm the guard exits `1` with a useful annotation
3. Verify the workflow YAML passes `actionlint` (or any equivalent)
4. Document each guard in `.github/scripts/readme.md` with one-line
   summary, invocation, and exit-code legend

---

## Handoff Template

When the AI finishes the integration, hand off with a summary that
answers:

- Which patterns were installed (1–6)
- Which were intentionally skipped, and why
- Where each script lives in `.github/scripts/`
- Where each baseline file is stored (committed file vs. cache)
- How to regenerate baselines (`--regenerate-baseline` flag, cache
  invalidation key)
- Expected first-run behavior (seeding warnings, no gating)

---

## Cross-References

- [00-overview.md](./00-overview.md) — Pattern inventory
- [01-shared-conventions.md](../01-shared-conventions.md) — Pinning, triggers, permissions
- [Coding Guidelines (Cross-Language)](../../02-coding-guidelines/01-cross-language/15-master-coding-guidelines/00-overview.md)
- [Linters CI/CD Integration](../../02-coding-guidelines/06-cicd-integration/97-acceptance-criteria.md)

---

*AI implementation guide — v1.0.0 — 2026-04-21*