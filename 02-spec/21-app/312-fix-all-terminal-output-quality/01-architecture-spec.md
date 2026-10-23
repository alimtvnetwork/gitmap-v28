# Plan 312 Architecture Specification: High-Quality Terminal Output for `gitmap fix --all`

## User Request (Verbatim)
```text
For gitmap fix --all have better quality terminal output please
```
Attached image: `media_1792723562507_90c913bb.png` displaying concurrent output collision, missing repo labels on steps, thread ID noise, and vendor path casing pollution.

---

## 1. Executive Summary & Design Principles

The goal of this specification is to elevate `gitmap fix --all` terminal UX to match GitMap's high-authority CLI standards:
1. **Zero Concurrency Collision:** All parallel worker output is coordinated through a thread-safe synchronized reporter. No partial lines (`... ` without `\n`) are ever emitted to stdout in parallel mode.
2. **Deterministic Tree Visualizer:** Completed repository remediations render as aligned tree nodes (`├── ` and `└── `) with color badges (`✔`/`✗`), aligned repository names, human-readable action summaries, and elapsed durations.
3. **Clamped Worker Diagnostics:** The concurrency header displays realistic worker counts (`min(concurrency, itemCount)`).
4. **Third-Party Dependency Hygiene:** Casing discrepancy detection strictly skips vendored and external dependency directories (`vendor/`, `node_modules/`).

---

## 2. Visual Output Specification

### 2.1 Batch Mode Execution (`gitmap fix --all`)
```text
  → gitmap fix --all (action: stash)
  Remediating 6 repository(ies) in parallel (6 workers)

  ├── [1/6] story-telling-v1          ✔ stashed -> pulled -> re-applied (0.8s)
  ├── [2/6] alim.karim.profile        ✔ stashed -> pulled -> re-applied (1.1s)
  ├── [3/6] wp-html-automate          ✔ stashed -> pulled -> re-applied (0.9s)
  ├── [4/6] kubernetes-training-v1    ✔ stashed -> pulled -> re-applied (1.2s)
  ├── [5/6] gitmap-v28                ✔ stashed -> pulled -> re-applied (0.7s)
  └── [6/6] seo-packages-v2           ✔ stashed -> pulled -> re-applied (1.0s)

  ✔ All 6 repository(ies) remediated successfully in 1.4s.
```

### 2.2 Partial Failure Handling
```text
  ├── [1/6] story-telling-v1          ✔ stashed -> pulled -> re-applied (0.8s)
  ├── [2/6] kubernetes-training-v1    ✗ git pull failed: merge conflict (1.2s)
  ...
  └── [6/6] seo-packages-v2           ✔ stashed -> pulled -> re-applied (1.0s)

  ⚠ Remediation completed with 1 failure(s) in 1.4s:
    • kubernetes-training-v1: git pull failed: merge conflict
      Suggested resolution: gitmap fix kubernetes-training-v1 resolve
```

### 2.3 Single-Repo Mode (`gitmap fix <repo> <action>`)
Maintains detailed step-by-step progress with repo context and clean timing:
```text
  Applying Fix: Option 1 (Stash & Re-apply) on story-telling-v1
  Plan:    Temporarily save local changes (including untracked), pull latest remote commits, then re-apply

  [1/3] ➜ git stash -u ... ✔ ok
  [2/3] ➜ git pull ... ✔ ok
  [3/3] ➜ git stash pop ... ✔ ok

  ✔ Fix applied successfully on story-telling-v1 (0.8s)
```

---

## 3. Disjoint Subtask Allocation

| Task-ID | Subtask | Owner | Owned Files | Status |
|---|---|---|---|---|
| Task-01 | Batch Remediation Reporter & Header | Worker 01 | `cli/cmdfix/fix_cmd.go` | PENDING |
| Task-02 | Silent Batch Recipe Execution & Action Summary | Worker 02 | `cli/cmdfix/fix_execute.go` | PENDING |
| Task-03 | Vendor Exclusion & Casing Output Polish | Worker 03 | `cli/cmdfix/fix_casing.go` | PENDING |

---

## 4. Coding & Structural Invariants
- Functions <= 15 lines.
- Positive booleans (`is*` or `has*` only).
- Blank line before return / after `}`.
- Relative paths only.
- FastGate pre-commit validation pass.
