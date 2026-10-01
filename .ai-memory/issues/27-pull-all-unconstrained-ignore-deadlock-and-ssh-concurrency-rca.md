# RCA-27: Pull-All Unconstrained Ignore Subprocess Deadlock and SSH Concurrency Multiplication

Spec Reference: [02-spec/22-app-issues/61-pull-all-unconstrained-ignore-deadlock-and-ssh-concurrency-rca.md](../../02-spec/22-app-issues/61-pull-all-unconstrained-ignore-deadlock-and-ssh-concurrency-rca.md)  
Parent Spec: [02-spec/21-app/201-pull-ignore-concurrency-split-db-cache-and-deadlock-prevention.md](../../02-spec/21-app/201-pull-ignore-concurrency-split-db-cache-and-deadlock-prevention.md)  
Plan Reference: [.ai-memory/plans/pending/64-pull-ignore-concurrency-split-db-cache-and-deadlock-prevention.md](../plans/pending/64-pull-ignore-concurrency-split-db-cache-and-deadlock-prevention.md)

---

## 1. Symptom

During batch repository operations (`gitmap pa` / `gitmap pull-all`) across 62 repositories:
1. The execution hangs indefinitely at 48% progress (30/62 repos).
2. The operating system task manager and terminal telemetry show up to 64 `git.exe` processes accumulating concurrently (`git.exe[*64]: 13720`), exhausting CPU cycles and file handles.
3. The terminal interface freezes and stops receiving heartbeat ticker updates.
4. Over SSH sessions, high concurrency causes SSH transport drops, socket timeouts, and remote VM freezing.

Telemetry screenshot: `assets/screenshots/64-pull-ignore-concurrency-split-db-cache-and-deadlock-prevention-01.png`

---

## 2. Root Cause

Unconstrained asynchronous `.gitignore` audits spawned up to 8 parallel workers firing ~370 git subprocesses simultaneously with the primary pull worker pool, git subprocess invocations lacked `context.WithTimeout` causing them to block indefinitely on index file lock contention, and remote SSH sessions inherited local high-concurrency settings.

---

## 3. Resolution

1. **Throttled Low-Priority Queue:**
   Clamp background ignore worker pool to max 1 worker (`CalculateIgnoreWorkersForPull() = 1`) and insert 10ms cooperative pauses between repository audits.
2. **Split-DB Ignore Caching:**
   Implement SQLite Split-DB cache at `.gitmap/data/gitignore/cache/sql.db` (`gitignore_repo_cache` table) with configurable TTL (default 24h). Warm repositories trigger 0 git subprocesses during pull.
3. **Fail-Safe Subprocess Timeouts:**
   Introduce `gitutil.ExecGitWithTimeout(15*time.Second, dir, args...)` with `context.WithTimeout`, ensuring any stalled process is killed cleanly.
4. **SSH Single-Hand Clamping:**
   Detect SSH sessions via `cloneconcurrency.IsSSHSession()` or `--ssh` and clamp execution to 1 worker, 1 hand (`--w 1 --hand 1`, `parallel = 1`).
5. **Configuration & Management CLI:**
   Add `gitmap ignore config`, `gitmap ignore cache`, and `gitmap ignore cache clear` commands.

---

## 4. Prevention & Learnings

- **Decouple Maintenance from Primary Workflows:** Background auditing tasks must run in low-priority, throttled queues and never contend with high-priority primary operations like `git pull`.
- **Mandatory Subprocess Context Timeouts:** All external process invocations (`git`, `powershell`, `sh`) must have explicit timeouts to eliminate zombie process accumulation.
- **Cache Expensive Audits:** Idempotent, high-cost disk audits must be cached in dedicated Split-DB tables with configurable TTLs.
- **Environment-Aware Concurrency:** Remote SSH sessions must automatically clamp to single-worker, single-hand execution to protect constrained remote environments.
