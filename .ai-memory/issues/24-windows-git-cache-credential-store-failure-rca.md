# 24 — Windows Git Cache Credential Store Failure, Pull-First Workflow, and GitIgnore Duplication RCA

## Status: Resolved
- **Issue ID:** RCA-58 / Issue 24
- **Spec Reference:** [02-spec/21-app/195-pas-formula-fix-ignores-cpar-and-repo-cache.md](../../02-spec/21-app/195-pas-formula-fix-ignores-cpar-and-repo-cache.md)
- **Target Subsystems:** `cli/cloner`, `cli/cmdpull`, `cli/cmdclone`, `cli/gitignoreagm`, `cli/cmdignore`
- **Date:** 2026-10-01

---

## 1. Reproduction & Symptoms

### 1.1 Fatal Credential Store Abort on Windows Remote Nodes
When executing `gitmap pa --ssh` or remote pulls across Windows nodes (`w1`, `w2`, `w3`, `w4`), git operations failed immediately with:
```text
fatal: Can not use the 'cache' credential store on Windows due to lack of UNIX socket support in Git for Windows.
```
As a result, all pulls failed across Windows fleet nodes.

### 1.2 Upfront Interactive Prompt Blocking Pull Workflow
When running `gitmap pa` or `gitmap pull-all`, GitMap stopped execution before initiating any pulls:
```text
Found .antigravity_resume_task.json in 1 repo(s)
Untrack, delete, and ignore .antigravity_resume_task.json? [Y/n]:
```
This broke automated unattended batch pulls, pipeline scripts, and fleet sync operations.

### 1.3 Duplicate .gitignore Entries Accumulation
Repositories repeatedly appended `.antigravity_resume_task.json` and section comments into `.gitignore`:
```gitignore
# Antigravity Task Persistence
.antigravity_resume_task.json
antigravity-resume_task.json
.antigravity-resume_task.json

# Antigravity Task Persistence
.antigravity_resume_task.json
antigravity-resume_task.json
.antigravity-resume_task.json
```
Lacking a sanitizing deduplicator, multiple runs multiplied identical lines.

---

## 2. Root Cause Analysis (4-Part RCA)

### 2.1 Direct Cause
- `GCM_CREDENTIAL_STORE=cache` was unconditionally injected into `buildSafePullEnv()`, `runGitPullCommand()`, and clone environments.
- On Windows, `git-credential-cache` requires UNIX domain sockets, which are absent in Git for Windows.
- `checkAgmResumeTaskBeforePull` was placed prior to batch pull execution in `executePullBatchLifecycle` and `processEfficientPullLifecycle`.
- `.gitignore` writers checked only line presence without cleaning up pre-existing duplicated blocks.

### 2.2 Indirect Cause
- `GCM_CREDENTIAL_STORE=cache` was added during credential contention remediation without platform gating.
- GitIgnore maintenance logic was coupled directly into the pre-flight phase of `gitmap pa` instead of running asynchronously or post-pull.

### 2.3 Environmental Context
- Fleet nodes run mixed Windows Server 2022 and Windows 11 with standard Git for Windows 2.4x installations lacking Unix domain socket daemon emulation.

### 2.4 Systemic Vulnerability
- Absence of a centralized, OS-aware Git subprocess environment helper (`buildGitSubprocessEnv()`) across `cloner`, `cmdpull`, and `cmdclone`.

---

## 3. Corrective & Preventive Actions

1. **OS-Aware Git Subprocess Environment**:
   - Centralize environment construction in `gitutil.BuildSafeGitEnv()`.
   - On Windows, set `GCM_NO_PERSIST=1`, `GCM_INTERACTIVE=never`, `GIT_TERMINAL_PROMPT=0`, but omit `GCM_CREDENTIAL_STORE=cache`.
2. **Pull-First Mandate**:
   - Defer `.gitignore` auditing until after repo pulls complete.
   - Run ignore scans concurrently or present non-blocking summaries at the conclusion of pulls.
3. **Formalize GitMap PAS Formula**:
   - Establish spec standard for `gitmap pas`: local direct in-process execution with bounded remote SSH fleet fanout.
4. **Idempotent GitIgnore Sanitizer**:
   - Clean, deduplicate, and sort entries while preserving comments and structure in `.gitignore`.
