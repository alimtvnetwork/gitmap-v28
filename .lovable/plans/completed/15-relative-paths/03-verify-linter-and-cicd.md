# Subtask 03: Verify Linter and CI/CD Local Runner

## Description

Run `python linter-scripts/check-relative-paths.py` to assert zero absolute paths remain, wire the check into `.lovable/ai-fix-scripts/03-cicd-local-runner.py`, and run the full local runner suite.

## Citations

- [Coding Guidelines](spec/02-coding-guidelines/06-ai-optimization/05-citation-requirement.md)
- [Coding Guidelines Master](.lovable/coding-guidelines/coding-guidelines.md)

## Acceptance Criteria

- [ ] `check-relative-paths.py` reports exit code 0.
- [ ] `03-cicd-local-runner.py` reports exit code 0 across all quality gates.
