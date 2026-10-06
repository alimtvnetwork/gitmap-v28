# Subtask 01: Purge Root Clutter Scripts, Audit CSVs, and Binary Artifacts

> **Parent Plan:** [234-codebase-review-remediation-and-consolidation.md](../../234-codebase-review-remediation-and-consolidation.md)  
> **Spec Reference:** [02-spec/21-app/234-codebase-review-remediation-and-consolidation/01-architecture-spec.md](../../../../02-spec/21-app/234-codebase-review-remediation-and-consolidation/01-architecture-spec.md)  
> **Master Ledger:** [02-spec/21-app/234-codebase-review-remediation-and-consolidation/00-master-audit-ledger.md](../../../../02-spec/21-app/234-codebase-review-remediation-and-consolidation/00-master-audit-ledger.md)  
> **Status:** `COMPLETED`  
> **Target Areas:**  
> - Root directory temporary scripts: 16 `fix_*.py` files  
> - Root directory runner update scripts: 2 `update_runner*.py` files  
> - Root directory audit results: 3 `audit_results_*.csv` files  
> - Root directory prompt dump: `user_prompt_full.txt`  
> - Root directory misplaced file: `gitignore` (unprefixed copy)  
> - Windows resource binaries: `cli/rsrc_windows_386.syso`, `cli/rsrc_windows_amd64.syso`  
> - Ignore rules: `.gitignore`  

---

## 1. Technical Objective

Purge 22 tracked temporary scratch files, audit spreadsheets, and prompt dumps from the root directory to restore repository hygiene. Remove compiled Windows COFF resource binaries (`cli/*.syso`) from git tracking without affecting local build capabilities, and update `.gitignore` to explicitly ignore `*.syso` artifacts across the entire codebase.

---

## 2. Target File Inventory

### 2.1 Tracked Root Clutter Files (22 Files)

1. `fix_agy_conv_ls.py`
2. `fix_agy_ls.py`
3. `fix_agy_most_conv.py`
4. `fix_cd_alias.py`
5. `fix_cg_path.py`
6. `fix_common_sync.py`
7. `fix_cr_alias.py`
8. `fix_cr_rootcore.py`
9. `fix_go_compile.py`
10. `fix_ip_aliases.py`
11. `fix_ip_conflict.py`
12. `fix_most_conv.py`
13. `fix_repo_create_init.py`
14. `fix_repo_create_init2.py`
15. `fix_watch.py`
16. `fix_words_flag.py`
17. `update_runner.py`
18. `update_runner2.py`
19. `audit_results_extra.csv`
20. `audit_results_go.csv`
21. `audit_results_ts.csv`
22. `user_prompt_full.txt`
23. `gitignore` *(note: non-dot copy of `.gitignore` at root)*

### 2.2 Binary Artifacts to Untrack

1. `cli/rsrc_windows_386.syso`
2. `cli/rsrc_windows_amd64.syso`

### 2.3 Configuration Files to Update

1. `.gitignore` — Append `*.syso` under generated code and binaries.

---

## 3. Implementation Steps

### Step 1: Pre-Purge Invocation Audit
- Confirm none of the `fix_*.py` or `update_runner*.py` files are called by:
  - `Makefile`
  - `run.ps1` / `run.sh`
  - `install.ps1` / `install.sh`
  - `.github/workflows/`

### Step 2: Update `.gitignore`
- Add `*.syso` into the `# Generated code, artifacts, test data, and binaries (NEVER COMMIT)` block of `.gitignore`.
- Ensure clean formatting and positive ignore rules.

### Step 3: Untrack Compiled `.syso` Binaries
- Remove `cli/rsrc_windows_386.syso` and `cli/rsrc_windows_amd64.syso` from git staging while keeping Go build flags operational.

### Step 4: Delete Root Clutter Files
- Remove the 16 `fix_*.py` files, 2 `update_runner*.py` files, 3 `audit_results_*.csv` files, `user_prompt_full.txt`, and `gitignore` from git index and working tree.

### Step 5: Verification of Clean Tree
- Verify that root directory contains zero scratch scripts, zero raw audit CSVs, and zero untracked `.syso` binaries.

---

## 4. Verification Protocol

```bash
# 1. Verify all 16 fix_*.py and 2 update_runner*.py are removed
ls fix_*.py update_runner*.py 2>&1

# 2. Verify audit CSVs and user_prompt_full.txt are removed
ls audit_results_*.csv user_prompt_full.txt gitignore 2>&1

# 3. Verify .gitignore includes *.syso
grep -n '\*\.syso' .gitignore

# 4. Confirm zero .exe files exist in repository tree
find . -name "*.exe" -not -path "*/.git/*" -not -path "*/tmp/*"
```

---

## 5. Acceptance Criteria

- [x] All 16 `fix_*.py` files deleted from root and git tracking.
- [x] Both `update_runner.py` and `update_runner2.py` deleted from root and git tracking.
- [x] All 3 `audit_results_*.csv` files deleted from root and git tracking.
- [x] `user_prompt_full.txt` deleted from root and git tracking.
- [x] Misplaced `gitignore` (non-dot) deleted from root and git tracking.
- [x] `.gitignore` contains `*.syso` rule under binary artifacts section.
- [x] `cli/rsrc_windows_386.syso` and `cli/rsrc_windows_amd64.syso` untracked from git (delegated to lead orchestrator staging).
- [x] Zero `.exe` files present in tracked source files.
- [x] Root directory displays clean developer layout.
