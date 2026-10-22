# Subtask 04: Eradicate Result Stutter Across CmdAutomation Types & Database Queries

> **Subtask ID:** `SUBTASK-297-04`  
> **Parent Task:** `297-result-wrapper-stutter-eradication`  
> **Status:** Completed  
> **Target Files:**  
> - `cli/cmdautomation/phase1_types.go`  
> - `cli/cmdautomation/phase2_types.go`  
> - `cli/cmdautomation/phase3_types.go`  
> - `cli/cmdautomation/phase4_types.go`  
> - `cli/cmdautomation/phase5_types.go`  
> - `cli/cmdautomation/phase6_types.go`  
> - `cli/cmdautomation/types.go`  
> - `cli/cmdautomation/runtime_db.go`  
> - `cli/cmdautomation/worker_pool.go`  

---

## 1. Action Items
1. In `phase1_types.go` through `phase6_types.go` and `types.go`:
   - Replace all `result.Result[...]` type aliases with canonical `result.Wrap[...]`.
2. In `runtime_db.go`:
   - Replace `result.Result[*sql.DB]` with `result.Wrap[*sql.DB]`.
   - Replace `result.Result[RuntimeRecord]` with `RuntimeRecordWrap` (or `result.Wrap[RuntimeRecord]`).
3. In `worker_pool.go`:
   - Replace `result.Result[FileContext]` with `result.Wrap[FileContext]`.
4. Coding Standards:
   - Functions $\le 15$ lines, zero nested ifs, positive booleans.
