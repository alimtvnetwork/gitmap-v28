# Reusable CI Guards — AI-Implementation Guide

**Version:** 1.0.0
**Updated:** 2026-04-21
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## Keywords

`reusable-ci` · `ci-guards` · `baseline-diff` · `name-collision` · `grandfather-baseline` · `lint-suggest` · `test-summary` · `ai-portable`

---

## Scoring

| Criterion | Status |
|-----------|--------|
| `00-overview.md` present | ✅ |
| AI Confidence assigned | ✅ |
| Ambiguity assigned | ✅ |
| Keywords present | ✅ |
| Scoring table present | ✅ |

---

## Purpose

This sub-section captures **six language-agnostic CI guard patterns** that
were originally implemented in Bash and Python under `.github/scripts/` of a
Go monorepo. Each pattern solves a problem that recurs in **every**
non-trivial repository — flat-namespace collisions, lint-debt creep,
inconsistent naming, opaque test failures — regardless of programming
language or build system.

The intent is that an AI assistant (or human engineer) reading this folder
can re-implement any of the six guards for **any repository in any
language** without ever seeing the original Go-specific source. Each
pattern file is structured as:

1. **Problem statement** — what failure mode the guard prevents
2. **Algorithm** — abstract, language-agnostic pseudocode
3. **Inputs / outputs / exit codes** — the contract
4. **Configuration surface** — what must be parameterized per repo
5. **Adaptations** — concrete worked examples for Go, Node/TypeScript,
   Python, and Rust
6. **Failure modes** — known pitfalls captured from production use

---

## Pattern Inventory

| # | File | Pattern | One-Line Summary |
|---|------|---------|------------------|
| 01 | [01-forbidden-name-guard.md](./01-forbidden-name-guard.md) | Forbidden-name guard | Block collision-prone helper names in flat-namespace packages |
| 02 | [02-grandfather-baseline-naming.md](./02-grandfather-baseline-naming.md) | Grandfather-baseline naming | Enforce naming convention on **new** identifiers only |
| 03 | [03-cross-file-collision-audit.md](./03-cross-file-collision-audit.md) | Cross-file collision audit | Detect duplicate / case-insensitive identifier collisions |
| 04 | [04-baseline-diff-lint-gate.md](./04-baseline-diff-lint-gate.md) | Baseline-diff lint gate | Fail build only on **new** lint findings vs cached baseline |
| 05 | [05-actionable-lint-suggestions.md](./05-actionable-lint-suggestions.md) | Actionable lint suggestions | PR comment mapping each new finding to a fix template |
| 06 | [06-matrix-test-aggregator.md](./06-matrix-test-aggregator.md) | Matrix test aggregator | Combine matrix-job test outputs into one copy-paste report |
| 07 | [07-shared-cli-wrapper.md](./07-shared-cli-wrapper.md) | Shared CLI wrapper | Unified `--phase check\|lint\|test\|all` entry point dispatching to all six guards |
| 08 | [08-config-schema.md](./08-config-schema.md) | Unified config schema | Single `ci-guards.yaml` parameterizes every guard; loader emits env vars consumed by Pattern 07 |
| 09 | [09-workflow-templates.md](./09-workflow-templates.md) | GitHub Actions templates | Composite action + reusable workflow + 4 language starters wrapping Patterns 07/08 |
| 99 | [99-ai-implementation-guide.md](./99-ai-implementation-guide.md) | AI handoff | How an AI should select, configure, and ship these guards |

---

## Why These Six?

| Pattern | Recurring Pain It Solves |
|---------|--------------------------|
| Forbidden-name guard | Two contributors independently invent `runOne()` in the same flat package → silent build break in CI |
| Grandfather baseline | Style policy change ("all constants must use `Cmd*` prefix") would require renaming 2,700 identifiers in one PR |
| Collision audit | Compiler errors ("redeclared in this block") only catch exact matches; case-only differences and intra-file dupes slip through code review |
| Baseline-diff lint gate | Adopting a new linter on a mature repo produces 5,000 findings; team can't fix them all but must not add new ones |
| Lint suggestions | Contributors see `[gocritic] paramTypeCombine` and have no idea what to change; suggestion bot turns each warning into a copy-paste diff |
| Test aggregator | A 12-shard test matrix produces 12 separate logs; reviewer wants one report listing every failing test and its assertion message |

---

## Common Design Principles

All six patterns share these properties — copy them when extending:

1. **Single binary, single concern.** Each script does one thing and
   exits with a meaningful status code (`0` clean, `1` violation,
   `2` tool error).
2. **String-literal aware.** Source-text scanners must skip code inside
   raw strings, regular strings, and comments — otherwise `WHERE` inside
   an SQL constant gets reported as a top-level identifier.
3. **Portable runtime.** Bash + POSIX `awk` + Python ≥ 3.10. No `pip
   install`, no Node, no compiled tools required on the runner.
4. **Annotation-friendly output.** Every error line uses the
   `::error file=…,line=…::` GitHub Actions annotation format so
   findings appear inline in the PR diff view.
5. **Idempotent.** Re-running the guard with the same inputs produces
   the same exit code and the same output. No timestamps, no random IDs.
6. **Fail-soft on missing baselines.** A baseline file that is absent or
   empty triggers a "seeding" mode that warns but exits `0`; the next
   run has something to diff against.

---

## Cross-References

- [Shared Conventions](../01-shared-conventions.md) — Pinning, triggers, permissions
- [CI Pipeline (Browser Extension)](../01-browser-extension-deploy/01-ci-pipeline.md)
- [CI Pipeline (Go Binary)](../02-go-binary-deploy/01-ci-pipeline.md)
- [Coding Guidelines (Cross-Language)](../../02-coding-guidelines/01-cross-language/15-master-coding-guidelines/00-overview.md)

---

*Reusable CI guards overview — v1.0.0 — 2026-04-21*