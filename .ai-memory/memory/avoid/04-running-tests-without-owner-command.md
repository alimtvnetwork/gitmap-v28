# Avoid: Running Tests Without Owner Explicit Command

**Status:** 🚫 Blocked — Explicit Owner Authorization Required
**Recorded:** 2026-09-12

---

## Rule

🔴 **NEVER execute unit tests, test suites, or test jobs (`go test`, `pytest`, `npm test`, or test jobs in CI runners) during standard development tasks, prompt executions, coding guideline fixes, or refactoring loops unless explicitly commanded by the repository owner.**

---

## Forbidden Actions

- ❌ Executing `go test`, `pytest`, `npm test`, or `cargo test` in regular turns.
- ❌ Running `python 03-ai-scripts/06-cicd-local-runner.py` without the `--no-tests` (or `--skip-tests`) flag.
- ❌ Adding automated test runs or test assertions into non-release prompts or workflows.
- ❌ Re-enabling test suites without explicit owner authorization.

---

## Permitted & Mandatory Actions

- ✅ Run tests **ONLY** when the repository owner explicitly requests it in their prompt (e.g., "run tests", "execute unit tests", "fix failing tests").
- ✅ Mandatory test runs during explicit release workflows (e.g. `03-ai-scripts/29-release-orchestrator.py` or explicit release prompts) where 100% test passing is a required pre-release quality gate.
- ✅ Always use `--no-tests` when running `03-ai-scripts/06-cicd-local-runner.py` for standard quality gate verification.

---

## Rationale

Unit test suites across polyglot repositories can be time-consuming, resource-intensive, and prone to non-deterministic side effects during rapid iterative file refactoring. Decoupling standard development verification (linting, AST checks, structural integrity) from test execution keeps iterations fast and deterministic while reserving comprehensive test execution for release ceremonies and explicit debugging requests.
