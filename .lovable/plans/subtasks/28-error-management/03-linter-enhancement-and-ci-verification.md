# Subtask 03 - Linter Enhancement & CI Verification

## Parent Specification

[28-error-management-audit.md](.lovable/plans/pending/28-error-management-audit.md)

## Acceptance Criteria & Requirements

- Update `linter-scripts/check-error-management.py` to add check for dual handling: calling `cliexit.HandleError` inside functions returning `error` followed by `return nil`.
- Run `python linter-scripts/check-error-management.py` to verify 0 violations.
- Run `python 03-ai-scripts/06-cicd-local-runner.py` to verify all 13 CI/CD quality gates exit with code 0.
- Update `.lovable/plans/index.md` upon completion.
