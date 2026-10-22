# Subtask 06: Modernize Utility Packages Type Aliases to result.Wrap

> **Subtask ID:** `SUBTASK-297-06`  
> **Parent Task:** `297-result-wrapper-stutter-eradication`  
> **Status:** Completed  
> **Target Files:**  
> - `cli/cmddownload/types.go`  
> - `cli/cmdupdate/types.go`  
> - `cli/cmdzsh/types.go`  
> - `cli/lazyregex/types.go`  
> - `cli/macro/types.go`  
> - `cli/netip/types.go`  
> - `cli/osfix/types.go`  
> - `cli/movemerge/types.go`  
> - `cli/cmdschedule/types.go`  

---

## 1. Action Items
1. In each targeted types file, replace `result.Result[...]` with canonical `result.Wrap[...]`.
2. Preserve backward-compatibility aliases if required.
3. Keep all files $\le 300$ lines.
