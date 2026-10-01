# Learned Memory 37: GitMap PAS Formula, Windows Git Subprocess Environment, and GitIgnore Deduplication

Date: 2026-10-01
Related Specs:
- [02-spec/21-app/195-pas-formula-fix-ignores-cpar-and-repo-cache.md](../../../02-spec/21-app/195-pas-formula-fix-ignores-cpar-and-repo-cache.md)
- [02-spec/22-app-issues/58-windows-git-cache-credential-store-failure-rca.md](../../../02-spec/22-app-issues/58-windows-git-cache-credential-store-failure-rca.md)
- [.ai-memory/plans/completed/58-pas-formula-fix-ignores-cpar-and-repo-cache.md](../../plans/completed/58-pas-formula-fix-ignores-cpar-and-repo-cache.md)

## 1. Architectural Lessons & Principles

### 1.1 The GitMap PAS Formula
- **Direct In-Process Local Execution**: Always run local operations on `127.0.0.1` synchronously without SSH or background worker wrapping.
- **Bounded Remote Fleet Concurrency**: Limit parallel SSH node tasks to bounded workers (`A=2`, `H=2`) to prevent CPU and lock contention on fleet machines.
- **Immediate Offline Bypassing**: Check probe status and skip offline nodes without blocking or stalling other fleet operations.

### 1.2 Windows Git Subprocess Credential Store Safety
On Windows, `git-credential-cache` fails with:
`fatal: Can not use the 'cache' credential store on Windows due to lack of UNIX socket support in Git for Windows.`
- Centralize all Git command environment creation in `gitutil.BuildSafeGitEnv()`.
- Exclude `GCM_CREDENTIAL_STORE=cache` on Windows while keeping `GCM_NO_PERSIST=1`, `GCM_INTERACTIVE=never`, `GIT_TERMINAL_PROMPT=0`, and `GIT_ASKPASS=""`.

### 1.3 Pull-First Non-Blocking Workflow
- Never block `gitmap pa`, `gitmap pull-all`, or `gitmap pas` with pre-flight interactive confirmations.
- All sanitization, cleanup, or migration checks must occur post-execution or as non-blocking background tasks.

### 1.4 GitIgnore Deduplication Sanitization
- Avoid repeatedly appending identical comment headers and ignore lines to `.gitignore`.
- Run `gitignoreagm.DeduplicateAndSanitizeGitignore()` to prune duplicate rules and comments cleanly.
