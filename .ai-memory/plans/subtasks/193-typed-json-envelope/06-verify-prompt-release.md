# Subtask 193.6: Verification, AI/CI Prompt, Minor Bump & Green CI/CD
Traceability ID: Task-06
Spec Reference: [02-spec/21-app/183-typed-json-envelope-and-format-inspection.md](../../../02-spec/21-app/183-typed-json-envelope-and-format-inspection.md)
Target Files: 01-prompts/24-verify-typed-json-envelope-and-which-format.md, 03-ai-scripts/42-clean-test-and-build-caches.py, 03-ai-scripts/29-release-orchestrator.py
Action: Run cache clean script `42-clean-test-and-build-caches.py`. Author verification prompt in `01-prompts/`. Consolidate plan and subtasks. Bump minor version to v6.404.0. Verify all GitHub Actions workflows are green.
Acceptance Criteria:
1. Cache cleanup executed.
2. Verification prompt authored.
3. Minor release v6.404.0 pushed.
4. GitHub Actions CI/CD workflows are 100% green.
Targeted Verification: `gh run list --limit 6`
