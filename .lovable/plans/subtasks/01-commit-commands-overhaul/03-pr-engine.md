# Subtask 3: PR Engine Integration

1. When --pr is active, intercept the commit application logic.
2. If ll, every applied commit is pushed to a new branch, a PR is created via gh pr create, merged, and deleted.
3. If 	ags or
elease, only commits associated with tags/releases trigger the PR flow.
4. Implement safe fallbacks or error handling using AppError.
