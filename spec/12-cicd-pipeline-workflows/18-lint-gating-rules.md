# Lint Gating Rules & Reusable CI Guard Architecture

> **/goal** Provide strict, automated lint gating across all language toolchains using baseline diffing and non-blocking grandfathered rules.
> **/learn** Read the modular guard specifications under `03-reusable-ci-guards/` before configuring repository lint pipelines.

## 🎯 Actionable CI/CD & Agent Checklist

1. [ ] `/goal` Apply baseline diff linting to block new violations while grandfathering existing debt.
2. [ ] `/learn` Read `02-spec/12-cicd-pipeline-workflows/03-reusable-ci-guards/05-baseline-diff-lint-gate.md` for algorithm details.
3. [ ] `/goal` Ensure all guard scripts emit standard `::error file=...,line=...::` GitHub annotations.
4. [ ] `/learn` Run all local verification linters via `python 03-ai-scripts/06-cicd-local-runner.py`.

---

## Modular Guard Execution Checklist

All concrete guard implementations and patterns are located under [`03-reusable-ci-guards/`](./03-reusable-ci-guards/01-index.md). Follow this numbered checklist:

1. [ ] `/learn` [02-forbidden-name-guard.md](./03-reusable-ci-guards/02-forbidden-name-guard.md) — Block collision-prone helper names in flat-namespace packages.
2. [ ] `/learn` [03-grandfather-baseline-naming.md](./03-reusable-ci-guards/03-grandfather-baseline-naming.md) — Enforce naming conventions on new identifiers only.
3. [ ] `/learn` [04-cross-file-collision-audit.md](./03-reusable-ci-guards/04-cross-file-collision-audit.md) — Detect duplicate and case-insensitive identifier collisions.
4. [ ] `/learn` [05-baseline-diff-lint-gate.md](./03-reusable-ci-guards/05-baseline-diff-lint-gate.md) — Fail build only on new lint findings vs cached baseline.
5. [ ] `/learn` [06-actionable-lint-suggestions.md](./03-reusable-ci-guards/06-actionable-lint-suggestions.md) — Generate PR comments mapping findings to fix templates.
6. [ ] `/learn` [07-matrix-test-aggregator.md](./03-reusable-ci-guards/07-matrix-test-aggregator.md) — Combine matrix-job test outputs into a single report.
7. [ ] `/learn` [08-shared-cli-wrapper.md](./03-reusable-ci-guards/08-shared-cli-wrapper.md) — Unified CLI wrapper dispatching to all guards.
8. [ ] `/learn` [09-config-schema.md](./03-reusable-ci-guards/09-config-schema.md) — Unified YAML configuration schema for guards.
9. [ ] `/learn` [10-workflow-templates.md](./03-reusable-ci-guards/10-workflow-templates.md) — Reusable GitHub Actions composite actions and starters.
10. [ ] `/learn` [11-changelog-awk-integration.md](./03-reusable-ci-guards/11-changelog-awk-integration.md) — AWK-free JSON-based changelog parser guard.
11. [ ] `/learn` [12-strict-enum-enforcement.md](./03-reusable-ci-guards/12-strict-enum-enforcement.md) — Strict validation of Type suffixes and Enum objects.
12. [ ] `/learn` [13-query-wrapper-python-ts.md](./03-reusable-ci-guards/13-query-wrapper-python-ts.md) — Cross-language validation of typed query wrappers.
