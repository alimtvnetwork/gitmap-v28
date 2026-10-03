# Reusable CI Guards — AI-Implementation Guide

> **/goal** Master and implement language-agnostic CI guards for baseline diffing, naming conventions, and lint gating.
> **/learn** Read the modular guard specifications in this directory to implement zero-dependency CI quality gates across any language.

## 🎯 Actionable CI/CD & Agent Checklist

1. [ ] `/goal` Block collision-prone identifiers and enforce positive naming conventions.
2. [ ] `/learn` Read `02-spec/12-cicd-pipeline-workflows/03-reusable-ci-guards/05-baseline-diff-lint-gate.md` to prevent new lint errors.
3. [ ] `/goal` Ensure every guard emits standard `::error file=...,line=...::` annotations.
4. [ ] `/learn` Verify local guard executions via `python 03-ai-scripts/06-cicd-local-runner.py`.

---

**Version:** 2.0.0
**Updated:** 2026-08-30
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## Pattern Inventory

| # | File | Pattern | One-Line Summary |
|---|------|---------|------------------|
| 02 | [02-forbidden-name-guard.md](./02-forbidden-name-guard.md) | Forbidden-name guard | Block collision-prone helper names in flat-namespace packages |
| 03 | [03-grandfather-baseline-naming.md](./03-grandfather-baseline-naming.md) | Grandfather-baseline naming | Enforce naming convention on **new** identifiers only |
| 04 | [04-cross-file-collision-audit.md](./04-cross-file-collision-audit.md) | Cross-file collision audit | Detect duplicate / case-insensitive identifier collisions |
| 05 | [05-baseline-diff-lint-gate.md](./05-baseline-diff-lint-gate.md) | Baseline-diff lint gate | Fail build only on **new** lint findings vs cached baseline |
| 06 | [06-actionable-lint-suggestions.md](./06-actionable-lint-suggestions.md) | Actionable lint suggestions | PR comment mapping each new finding to a fix template |
| 07 | [07-matrix-test-aggregator.md](./07-matrix-test-aggregator.md) | Matrix test aggregator | Combine matrix-job test outputs into one copy-paste report |
| 08 | [08-shared-cli-wrapper.md](./08-shared-cli-wrapper.md) | Shared CLI wrapper | Unified `--phase check\|lint\|test\|all` entry point dispatching to all guards |
| 09 | [09-config-schema.md](./09-config-schema.md) | Unified config schema | Single `ci-guards.yaml` parameterizes every guard |
| 10 | [10-workflow-templates.md](./10-workflow-templates.md) | GitHub Actions templates | Composite action + reusable workflow + 4 language starters |
| 11 | [11-changelog-awk-integration.md](./11-changelog-awk-integration.md) | Changelog parser guard | AWK-free JSON-based changelog extraction in CI |
| 12 | [12-strict-enum-enforcement.md](./12-strict-enum-enforcement.md) | Enum convention guard | Strict validation of Type suffixes and Enum objects |
| 13 | [13-query-wrapper-python-ts.md](./13-query-wrapper-python-ts.md) | Query wrapper guard | Validation of cross-language typed query wrappers |
| 99 | [99-ai-implementation-guide.md](./99-ai-implementation-guide.md) | AI handoff | How an AI should select, configure, and ship these guards |
