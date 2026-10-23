# RCA-110: `gitmap fix --all` Terminal Output Collision, Concurrency Tearing & Missing Repo Context

## 1. Symptom & Evidence
When executing `gitmap fix --all` with multiple pending repositories:
```text
asia-uploader\vendor\theseer\tokenizer\readme.md in wp-html-automate
• Case Normalization: Renamed wp-deploy/README.md -> wp-deploy\readme.md in wp-link-manager
• Case Normalization: Renamed wp-plugins/scripts/README.md -> wp-plugins\scripts\readme.md in wp-onboarding
Remediating 6 repository(ies) in parallel (20 workers) with action: stash

[Worker 6] Starting remediation on: seo-packages-v2
[Worker 3] Starting remediation on: story-telling-v1
[Worker 1] Starting remediation on: wp-html-automate
[Worker 2] Starting remediation on: alim.karim.profile
[Worker 4] Starting remediation on: gitmap-v28
[Worker 5] Starting remediation on: kubernetes-training-v1
Applying Fix: Option 1 (Stash & Re-apply) on story-telling-v1
Plan:    Temporarily save local changes (including untracked), pull latest remote commits, then re-apply

[1/3] → git stash -u ... [ Applying Fix: Option 1 (Stash & Re-apply) on kubernetes-training-v1
Plan:    Temporarily save local changes (including untracked), pull latest remote commits, then re-apply

[1/3] → git stash -u ... [ Applying Fix: Option 1 (Stash & Re-apply) on alim.karim.profile
Plan:    Temporarily save local changes (including untracked), pull latest remote commits, then re-apply

[1/3] → git stash -u ... 🗹 ok
[2/3] → git pull ... 🗹 ok
[2/3] → git pull ... 🗹 ok
```

### Visual & Functional Defects Identified:
1. **Concurrent Stdout Collision & Tearing:**
   Workers printed partial strings without newlines (`[1/3] → git stash -u ... `), causing subsequent worker headers (`Applying Fix: Option 1...`) to interleave on the exact same terminal line:
   `[1/3] → git stash -u ... [ Applying Fix: Option 1 (Stash & Re-apply) on kubernetes-training-v1`
2. **Missing Repository Context on Progress Lines:**
   Lines like `[1/3] → git stash -u ... ok` and `[2/3] → git pull ... ok` contained no repository name, making it impossible to determine which repository completed which step.
3. **Noisy Worker ID Clutter:**
   `[Worker 6] Starting remediation on: ...` and `[Worker 6] Completed (1/6): ...` exposed internal thread IDs that provided no value to the user.
4. **Repetitive 2-Line Plan Dumping:**
   `Applying Fix: Option 1 ...` and `Plan: Temporarily save local changes...` was printed redundantly for every single repository in parallel.
5. **Inflated Worker Count:**
   Reported `(20 workers)` even when only 6 repositories existed.
6. **Vendor Path Casing Pollution:**
   `reconcileAllCasingDiscrepancies` inspected all tracked files indiscriminately, attempting to normalize third-party vendored dependencies (`vendor/theseer/tokenizer/README.md`).

---

## 2. Root Cause Analysis
1. **Unbuffered & Unlocked Stdout in Parallel Execution:**
   In `cli/cmdfix/fix_execute.go`, `executeSingleStep` executed:
   `fmt.Printf("  [%d/%d] ➜ %s ... ", idx, total, stepCmd)`
   without a newline character (`\n`), expecting to append `ok` upon command completion. When 6 or 20 workers executed in parallel via `parallel.ExecuteParallel`, multiple goroutines invoked `fmt.Printf` concurrently without synchronization, causing terminal line clobbering and race conditions.
2. **Missing Repo Name in Step Formatter:**
   `executeSingleStep` formatted the command string (`formatStepCommand(step)`) without prepending `item.RepoName`.
3. **No Batch Mode Separation:**
   `executeFixRecipe` was written solely for single-repo sequential execution and lacked a batch-mode execution path (`executeFixRecipeBatch`) that executes steps quietly and yields a structured summary.
4. **Unclamped Concurrency Calculation:**
   `parallel.ResolveConcurrency(0)` returned CPU-based maximum (e.g. 20), which was printed directly without clamping to `itemCount`.
5. **Vendored Path Filtering Absence in Casing Inspector:**
   `inspectReadmeDuality` did not filter out `vendor/` or `node_modules/` paths from `git ls-files`.

---

## 3. Corrective Architecture & Implementation Plan
1. **Clean Batch Execution Pipeline (`executeFixRecipeBatch`):**
   Execute remediation steps silently for each repository in batch mode, capturing errors, benign outcomes, and measuring elapsed duration.
2. **Synchronized Tree Node Progress Reporter (`FixBatchReporter`):**
   Render completed repositories under `sync.Mutex` using GitMap standard tree connectors (`├── ` and `└── `), aligned repo columns, color status badges (`✔`/`✗`), human-friendly action summaries (`stashed -> pulled -> re-applied`), and step durations `(0.8s)`.
3. **Clamped Concurrency Header:**
   Display `min(concurrency, itemCount)` workers with clear invocation header:
   `→ gitmap fix --all (action: <action>)`
4. **Vendor Directory Exclusion in Casing Remediation:**
   Exclude `vendor/`, `node_modules/`, and `.git/` from `inspectReadmeDuality`.
5. **Actionable Post-Remediation Summary:**
   If any repository fails, print a structured failure card with the exact error and suggested remediation command.

---

## 4. Verification Gate
- Pass FastGate: `python 03-ai-scripts/50-fastgate.py`.
- Pass Relative Path Linter: `python linter-scripts/check-relative-paths.py`.
- Pass Newline Styling: `python linter-scripts/check-newline-styling.py --all`.
- Pass Unit Tests: `go test -v ./cli/cmdfix/...` (isolated via FastGate rules).
