# Learned Memory 37: GitMap PAS Formula, Windows Credential Cache Safety, Ignore Management Suite, CPAR, and Split-DB Cache Engine

Date: 2026-10-01
Related Specs:
- [02-spec/21-app/195-pas-formula-fix-ignores-cpar-and-repo-cache.md](../../../02-spec/21-app/195-pas-formula-fix-ignores-cpar-and-repo-cache.md)
- [02-spec/22-app-issues/58-windows-git-cache-credential-store-failure-rca.md](../../../02-spec/22-app-issues/58-windows-git-cache-credential-store-failure-rca.md)
- [.ai-memory/issues/24-windows-git-cache-credential-store-failure-rca.md](../../issues/24-windows-git-cache-credential-store-failure-rca.md)
- [.ai-memory/plans/pending/58-pas-formula-fix-ignores-cpar-and-repo-cache.md](../../plans/pending/58-pas-formula-fix-ignores-cpar-and-repo-cache.md)
- [.ai-memory/spec/commands/10-pas-formula-ignore-suite-cpar-and-cache.md](../../spec/commands/10-pas-formula-ignore-suite-cpar-and-cache.md)

## 1. Architectural Lessons & Principles

### 1.1 The GitMap PAS Formula (Pull-All-SSH Standard)
- Local host execution MUST be direct and in-process. Wrapping local execution in SSH loops introduces unnecessary latency, key validation issues, and loses terminal color/progress bar richness.
- Remote nodes are executed in parallel using bounded worker concurrency (`A=2`, `H=2`).
- Structured audit logs and error tables must be written locally to allow inspection via `gitmap see errors ssh` (`ses`) and `gitmap see history`.

### 1.2 Windows Credential Store Daemon Absence
- Git for Windows does NOT support the `cache` credential store daemon because Windows lacks UNIX domain sockets out-of-the-box.
- Subprocesses invoking git on Windows will immediately abort with:
  `fatal: Can not use the 'cache' credential store on Windows due to lack of UNIX socket support in Git for Windows.`
- Standardize on `gitutil.BuildSafeGitEnv()`: on Windows, omit `GCM_CREDENTIAL_STORE=cache` while retaining `GCM_NO_PERSIST=1`, `GCM_INTERACTIVE=never`, and `GIT_TERMINAL_PROMPT=0`. On Linux/macOS, `cache` is permitted.

### 1.3 Pull-First Decoupling
- Routine batch operations (`gitmap pa`, `gitmap pas`) must prioritize pulling immediately without halting on blocking interactive prompts upfront.
- Maintenance inspections (e.g. `.gitignore` issues or resume task prompts) must run post-pull or asynchronously.

### 1.4 GitIgnore Idempotent Sanitization
- Naive line appending multiplies duplicate patterns over repeated runs.
- The `DeduplicateAndSanitizeGitignore` engine preserves comment blocks and structure, deduplicates pattern lines, and automatically guarantees `.gitmap/backup/` and `.antigravity_resume_task.json` are properly ignored.

### 1.5 Repository Split-DB Cache Engine
- Partitioning file lines across folder slug databases (`.gitmap/cache/<slug>.db`) prevents massive SQLite table contention and enables zero-latency concurrent search.
- Skipping files >200KB and binary assets keeps cache databases compact and performant.
