---
name: autonomous-qa-and-testing
description: Autonomously run test suites, verify quality gates, and prevent regressions across polyglot stacks.
---

# Autonomous QA and Testing

> [!IMPORTANT]
> **Owner Command Required:** Running unit tests and test suites (`go test`, `npm run test`, `pytest`) is strictly prohibited in standard development turns. Tests may ONLY be run when explicitly commanded by the repository owner, or as part of the mandatory pre-release quality gate during a release ceremony. In routine turns, the full CI/CD runner is banned; run targeted single-file linters instead. The full runner (`06-cicd-local-runner.py`) may ONLY run when explicitly commanded by the repository owner.

Executes comprehensive testing, linting verification, and quality gate validation.

## Checks
1. **TypeScript / React:** `npm run test`, `npm run lint`
2. **Go:** `go test ./...`
3. **Targeted Verification (Routine):** Run targeted linters on modified files (full `06-cicd-local-runner.py` is strictly banned in routine turns).
4. **Full Runner (Owner Explicit Command Only):** `python 03-ai-scripts/06-cicd-local-runner.py`
5. **Pre-Release Full Gates (Release Ceremony Only):** `python 03-ai-scripts/06-cicd-local-runner.py --run-tests`
4. **Git Hygiene:** Verify no un-ignored test dumps or binaries via `git status`.
