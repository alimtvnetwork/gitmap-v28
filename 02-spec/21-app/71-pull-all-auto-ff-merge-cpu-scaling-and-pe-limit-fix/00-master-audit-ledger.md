# Master Audit Ledger: Task 71 - Pull-All Auto FF Merge, CPU Auto-Scaling, Batch Fix, and PE Flag Fix

**Version:** 1.0.0  
**Status:** In Progress (Phase 1 Planning & Discovery)  
**Parent Task ID:** 1  
**Database:** `.ai-memory/temp-agents/71-71-pull-all-auto-ff-merge/agent-task.db`  

---

## 1. User Request (Verbatim)

```text
PS C:\Users\Administrator> gitmap pe -l 10

  ✖ Repository not found: "10"

  Did you mean:
    • gitmap pe awansoft-v10
    • gitmap pe enum-v10

gitmap pe: execute failed: target repository "10" not found
Stack Trace:
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.resolveErrorStackTrace (cmd/root.go:292)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.handleGlobalError (cmd/root.go:221)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.handleDispatchResult (cmd/root.go:461)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.dispatch (cmd/root.go:497)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.runDispatch (cmd/root.go:145)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.Run (cmd/root.go:112)
    at main.main (cli/main.go:7)
PS C:\Users\Administrator> gitmap pe -limit 10

  ✖ Repository not found: "10"

  Did you mean:
    • gitmap pe awansoft-v10
    • gitmap pe enum-v10

gitmap pe: execute failed: target repository "10" not found
Stack Trace:
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.resolveErrorStackTrace (cmd/root.go:292)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.handleGlobalError (cmd/root.go:221)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.handleDispatchResult (cmd/root.go:461)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.dispatch (cmd/root.go:497)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.runDispatch (cmd/root.go:145)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.Run (cmd/root.go:112)
    at main.main (cli/main.go:7)

Hi there. Can you please try to fix this error? And also when we do the git map PA, that means pull all, the spec should have given. It should actually try to fix all this if the fast forward is required. Okay? Don't do the rebase, try to do the fast forward and merge. Now, in case that user wants to go with all of these failure to fix rather than singular single, it should also have one command that actually shows to start with all of these failed projects to fix. So that is what I want, one command. And also at the end, after you resolve the ignore, you should also prompt the user, do they want it to resolve all of these all together or one by one, just like ignore. And usually by default, it will try to do the fast forward without the rebase, and it should fix it. And for example, the last one, you said unknown repository. Okay? That actually says repo cache. It's actually repo cache. Why you have unknown? So try to figure this out, what is the root cause of it, because from the current error, it's really hard to figure it out. And this happens in this machine, so you try to use the git map ignores, other stuff to see the stack trace and find the root cause and try to fix it. Try to find the root cause of it, right? The root cause of it, try to find the solution, and then you finally do a release bump, and then finally you release it. Do you understand? And can we increase the performance of the git pull all a little bit automatically if the CPU pressure is not that much? Again, try to have this automation automatically, so it will extend the power if the CPU is not that much in use. Okay, so it will auto scale. So there should be two types of configuration when CPU is not using that much, when the CPU is in too much pressure. Two types of configuration needs to be there, but also at the same time, user can actually change their configuration. User can do a PA space help to see these type of configuration and flags that they could use. Try to understand the issue, try to fix the issue, and then solve the problem. Is it clear?
```

---

## 2. Discrete Technical Deliverables (Traceable IDs)

1. **Task-01: Fix `gitmap pe -l/-limit/--limit` Flag Parsing & Extraction**
   - Support `-l`, `-limit`, `--limit`, `-n`, `--lines` flags in `gitmap pe` / `gitmap pipeline-errors`.
   - Ensure numbers like `10` are recognized as the line/entry limit rather than positional repository slugs.

2. **Task-02: Investigate and Resolve Unknown Repo-Cache Issue in Pull-All**
   - Root-cause why repository paths in `repo-cache` or `./repo...` resolve as slug/name "unknown".
   - Fix name/slug resolution fallback to derive repository names from path basename or remote URL.

3. **Task-03: Implement Auto Fast-Forward Merge Fallback (No Rebase) in Pull-All**
   - Automatically attempt fast-forward merge (`git merge --ff-only` or clean non-ff standard merge) without rebase.
   - Heal divergent branches cleanly during pull without failing with "safe-pull failed".

4. **Task-04: Batch Failure Resolution Command & Interactive Prompt**
   - Single command to fix all failed projects from the last pull batch: `gitmap pull-fix` / `gitmap pa --fix`.
   - Interactive prompt at the end of `gitmap pa` when failures occur: resolve all together or one by one.

5. **Task-05: CPU-Aware Concurrency Auto-Scaling for Pull-All**
   - Dynamically sample CPU usage / core availability and auto-scale worker concurrency (Low Pressure vs High Pressure).
   - Provide configuration presets and expose flags in `gitmap pa --help`.

6. **Task-06: Verification, Golden Tests, Version Bump & Minor Release Ceremony**
   - Run unit tests, guideline autofixer, nested if linters, bump minor version, tag, merge, and publish release.

---

## 3. Disjoint Subtask Ownership Matrix

| Subtask ID | Task Code | Title | Assigned Agent | Owned File Scope | Status |
| :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | Task-01 | Fix `gitmap pe` flag parsing (`-l`, `-limit`) | Worker 01 | `cli/cmdpipeline/cmd_pipeline_errors.go`, `cli/cmd/pipeline.go` | PENDING |
| 2 | Task-02 | Resolve "unknown" repo-cache issue | Worker 01 | `cli/cmdpull/pull_inventory.go`, `cli/cmdpull/pull_cache.go`, `cli/scanner/cache.go` | PENDING |
| 3 | Task-03 | Auto fast-forward merge in pull-all | Worker 02 | `cli/cmdpull/pull_safe.go`, `cli/cmdpull/pull_exec.go`, `cli/cmdpull/pull_merge.go` | PENDING |
| 4 | Task-04 | Batch failure resolution & prompt | Worker 02 | `cli/cmdpull/pull_failures.go`, `cli/cmdpull/pull_interactive.go`, `cli/cmd/pull_fix.go` | PENDING |
| 5 | Task-05 | CPU-aware concurrency auto-scaling | Worker 01 | `cli/cmdpull/pull_autoscale.go`, `cli/cmdpull/pull_cpu.go`, `cli/cmdpull/pull_config.go` | PENDING |
| 6 | Task-06 | Verification, Bump & Release Ceremony | Lead Orchestrator | `version.json`, `package.json`, `cli/constants/constants.go`, docs | PENDING |
