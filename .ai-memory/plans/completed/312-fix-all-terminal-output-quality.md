# Plan 312: High-Quality Terminal Output for `gitmap fix --all`

## User Request (Verbatim)
```text
For gitmap fix --all have better quality terminal output please
```
Attached image: `media_1792723562507_90c913bb.png`

---

## 1. Disjoint Subtask Breakdown

| Task-ID | Subtask | Owner | Owned Files | Status | Evidence |
|---|---|---|---|---|---|
| Task-01 | Batch Tree Reporter, Clamped Concurrency & Result Summary | Worker 01 | `cli/cmdfix/fix_cmd.go` | COMPLETED | Synchronized tree reporter with branch connectors (├──/└──), aligned cols, duration badges, clamped workers |
| Task-02 | Silent Batch Recipe Execution, Step Isolation & Action Formatter | Worker 02 | `cli/cmdfix/fix_execute.go` | COMPLETED | executeFixRecipeBatch and executeSingleStepSilent eliminate stdout collisions, summarizeRecipeAction |
| Task-03 | Vendor Filtering & Casing Output Hygiene | Worker 03 | `cli/cmdfix/fix_casing.go` | COMPLETED | isVendorOrDependencyPath skips vendor/node_modules/.git, clean bullet formatting |

---

## 2. Invariants
- Functions <= 15 lines.
- Blank lines before return / after `}`.
- Positive booleans (`is*` or `has*` only).
- Strict relative paths.
- FastGate pass (exit 0).
