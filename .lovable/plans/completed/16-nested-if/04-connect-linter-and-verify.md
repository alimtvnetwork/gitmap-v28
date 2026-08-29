# Subtask 04: Connect Nested If Linter to CI/CD and Verify

## Description
Register `linter-scripts/check-nested-ifs.py` into `.github/workflows/ci.yml` and `.lovable/ai-fix-scripts/03-cicd-local-runner.py`. Verify all 22 quality gates pass.

## Citations
- [Braces and Nesting](spec/02-coding-guidelines/01-cross-language/04-code-style/01-braces-and-nesting.md)
- [Canonical Size Tier](spec/02-coding-guidelines/00-canonical-size-tier.md)
- [Master Coding Guidelines](.lovable/coding-guidelines/coding-guidelines.md)

## Acceptance Criteria
- [ ] `check-nested-ifs.py` passes with 0 violations.
- [ ] `03-cicd-local-runner.py` passes all quality gates with exit code 0.
