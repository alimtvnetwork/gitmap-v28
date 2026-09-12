# Subtask 02: Update Gitmap Prompts and Skills

> **Parent Plan:** [129-smart-test-runner-and-eta-sleep-sync.md](../../pending/129-smart-test-runner-and-eta-sleep-sync.md)
> **Status:** completed

## Requirements
1. Update prompts in `01-prompts/16-ci-cd/`:
   - Document `.lovable/temp/failures/` failure log isolation.
   - Document 100% silent passing tests rule (zero console lines, zero disk files).
   - Document dual-queue concurrency: slow queue (4w x 2 tests) and fast queue (4w x 4 tests in 100-test chunks).
   - Document dynamic ETA sleep protocol using `.lovable/temp/runner-eta.json` (AI sleeps instead of looping, re-checks remaining ETA if still running).
   - Document repository-relative code-to-test mapping and configurable slow test threshold (`4.0s`).
2. Update prompts in `01-prompts/14-execute/` and `01-prompts/15-cg-execute/`:
   - Incorporate the smart runner parameters, failure folder, and ETA sleep protocol into execution and QA workflows.
3. Update skills in `.agents/skills/`:
   - `autonomous-qa-and-testing/skill.md`: update with smart test runner, failure isolation, dual-queue, and ETA sleep protocol.
   - `ci-cd-fix/skill.md` & `ci-cd-create/skill.md`: reference the new test runner rules.
   - `smart-test-runner-and-inventory/skill.md`: create new dedicated skill documenting the complete runner and inventory contract.

## Acceptance Criteria
- [x] All CI/CD and execution prompts updated with the new test runner and sleep protocol specifications.
- [x] Skills updated and new `smart-test-runner-and-inventory` skill created.
- [x] No hardcoded absolute paths in any modified or new files.
