# Subtask 02: Modernize PRDB Connection & Operations to Monadic Result Containers

> **Subtask ID:** `SUBTASK-297-02`  
> **Parent Task:** `297-result-wrapper-stutter-eradication`  
> **Status:** Completed  
> **Target Files:**  
> - `cli/prdb/pr_split_conn.go`  
> - `cli/prdb/pr_split_ops.go`  

---

## 1. Action Items
1. In `cli/prdb/pr_split_conn.go`:
   - Replace `result.Result[*PrSplitDb]` with `PrSplitDbWrap` across `OpenPrSplitDb`, `OpenPrSplitDbAt`, `openPrSplitDbWithSlug`, and `initPrSplitConn`.
2. In `cli/prdb/pr_split_ops.go`:
   - Replace `result.Result[int64]` with `result.Int64Result` across `CreatePullRequest` and `AddPrRelease`.
   - Replace `result.Result[*PullRequestRecord]` with `PullRequestWrap` in `scanPullRequestRecord` and `GetPullRequestByNumber`.
   - Replace `result.Result[bool]` with `result.BoolResult` in `UpdatePullRequestStatus`, `UpsertPrBranch`, and `MarkPrBranchDeleted`.
   - Replace `result.Result[PrBranchRecord]` with `PrBranchRecordWrap` in `scanSinglePrBranch`.
3. Coding Standards:
   - Functions $\le 15$ lines, zero nested ifs, positive booleans, blank line before returns.
