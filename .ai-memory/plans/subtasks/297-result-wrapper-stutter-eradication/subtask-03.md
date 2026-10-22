# Subtask 03: Modernize OSUser Portable and OSClean Cache Cleaning Functions

> **Subtask ID:** `SUBTASK-297-03`  
> **Parent Task:** `297-result-wrapper-stutter-eradication`  
> **Status:** Completed  
> **Target Files:**  
> - `cli/osuser/group_portable.go`  
> - `cli/osuser/portable.go`  
> - `cli/osclean/dev_clean.go`  
> - `cli/osclean/terminal_clean.go`  

---

## 1. Action Items
1. In `cli/osuser/group_portable.go`:
   - Replace `result.Result[PortableGroup]` with `PortableGroupWrap`.
   - Replace `result.Result[bool]` with `result.BoolResult`.
   - Replace `result.Result[int]` with `result.IntResult`.
2. In `cli/osuser/portable.go`:
   - Replace `result.Result[PortableUser]` with `PortableUserWrap`.
   - Replace `result.Result[bool]` with `result.BoolResult`.
   - Replace `result.Result[int]` with `result.IntResult`.
3. In `cli/osclean/dev_clean.go`:
   - Replace `result.Result[DevCleanSummary]` with `DevCleanSummaryWrap`.
4. In `cli/osclean/terminal_clean.go`:
   - Replace `result.Result[TerminalCleanSummary]` with `TerminalCleanSummaryWrap`.
5. Coding Standards:
   - Functions $\le 15$ lines, zero nested ifs, positive booleans, blank lines around returns.
