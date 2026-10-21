# Subtask 04: Storytelling README Casing Remediation & Repo-Cache Recording

- **Subtask Path**: `.ai-memory/plans/subtasks/gitmap-enhancement-and-issue-resolution/04-storytelling-readme-casing.md`
- **Parent Plan**: [.ai-memory/plans/gitmap-enhancement-and-issue-resolution.md](../../gitmap-enhancement-and-issue-resolution.md)
- **Specification Reference**: [02-spec/21-app/gitmap-enhancement-and-issue-resolution/02-component-spec.md](../../../02-spec/21-app/gitmap-enhancement-and-issue-resolution/02-component-spec.md)
- **Primary Owner**: Worker 02 (Casing & Fix Engine Specialist)
- **Target Files (Exclusive)**:
  - `cli/cmdfix/fix_casing.go` (new file)
  - `cli/cmdfix/fix_casing_test.go` (new file)
  - `cli/cmdfix/fix_cmd.go` (integration update)
  - `repo-cache/01-gitmap/remediation_recipes.json` (recipe recording)
- **Status**: READY
- **Target Release**: GitMap v6.590.0

---

## 1. Context & Objectives

In repositories such as `storytelling-v1`, developers face an intractable dirty working tree loop caused by Windows NTFS case-insensitivity combined with Git's default `core.ignorecase = true`:
1. **Case Collision / Duality**: The repository history tracks `README.md` (or both `story/README.md` and `story/readme.md`), while workspace standards mandate lowercase `readme.md`.
2. **Rename Deadlock**: Directly executing `git mv README.md readme.md` on Windows fails with `fatal: destination exists`. Discarding or stashing cannot resolve the conflict because checking out one entry overwrites the single physical NTFS disk inode, immediately causing Git status to flag the other entry as modified or deleted.
3. **Core Objectives**:
   - Provide automated detection of casing collisions (`DetectCasingDiscrepancies`).
   - Execute an atomic two-step rename (`README.md` -> `README.md.tmp-lcf` -> `readme.md`) via the proven routines in `cli/cmd/lowercasefix_git.go`.
   - Purge uppercase index duplicates (`git rm --cached -f`) when both entries exist in `.git/index`.
   - Wire the solution directly into `gitmap fix-all` (`fa`) in `cli/cmdfix/fix_cmd.go`.
   - Persist the executed remediation recipe into `repo-cache/01-gitmap/remediation_recipes.json` and local state (`~/.gitmap/output/last_remediation.json`).

---

## 2. Step-by-Step Implementation Plan

### Step 2.1: Implement `cli/cmdfix/fix_casing.go`

1. **Define Data Structures**:
   ```go
   package cmdfix

   type CasingDiscrepancy struct {
       RepoPath      string `json:"repoPath"`
       CurrentPath   string `json:"currentPath"`
       CanonicalPath string `json:"canonicalPath"`
       IsDualTracked bool   `json:"isDualTracked"`
   }

   type CasingFixResult struct {
       RepoPath      string `json:"repoPath"`
       IsFixed       bool   `json:"isFixed"`
       RecipeId      string `json:"recipeId"`
       ErrorMessage  string `json:"errorMessage,omitempty"`
   }
   ```

2. **Implement Detection (`DetectCasingDiscrepancies`)**:
   - Inspect Git index via `git ls-files` in `repoPath`.
   - Identify files matching case-insensitive `(?i)^readme\.md$` or subfolder readmes (`(?i)/readme\.md$`).
   - If the tracked or filesystem path contains uppercase characters (e.g. `README.md`), create a `CasingDiscrepancy` record.
   - Check if both lowercase and uppercase entries exist simultaneously (`IsDualTracked = true`).
   - Ensure function length <= 15 lines with guard clauses.

3. **Implement Atomic Remediation (`ExecuteCasingNormalization`)**:
   - If `IsDualTracked`:
     - Buffer current disk file contents to memory.
     - Execute `git rm --cached -f <uppercase_path>` to purge the uppercase entry from index.
     - Verify lowercase file on disk is present (restore from buffer if needed).
     - Execute `git add -- <canonical_path>`.
     - Execute `git commit -m "fix(case): normalize uppercase duplicate to canonical lowercase"`.
   - If single tracked uppercase file:
     - Execute two-step rename:
       - Step 1: `git mv <uppercase> <uppercase>.tmp-lcf`.
       - Step 2: `git mv <uppercase>.tmp-lcf <canonical_lowercase>`.
       - Rollback: if Step 2 fails, restore from `<uppercase>.tmp-lcf` back to `<uppercase>`.
       - Fallback: if `git mv` fails, execute two-step `os.Rename` followed by `git add <canonical_lowercase>`.
   - Ensure all sub-operations return direct errors without calling exit handlers.

4. **Implement Recipe Recording (`RecordCasingRemediationRecipe`)**:
   - Resolve `repo-cache` path via `special_repos_rc_heatmap.go:ResolveRepoCachePath()`.
   - Ensure target directory `01-gitmap/` exists.
   - Append structured JSON recipe record to `remediation_recipes.json`.
   - Also write to local state `~/.gitmap/output/last_remediation.json`.

---

### Step 2.2: Integrate into `cli/cmdfix/fix_cmd.go`

1. **Add Request Matcher**:
   ```go
   func isFixCasingRequest(args []string) bool {
       for _, arg := range args {
           low := strings.ToLower(arg)
           if low == "casing" || low == "case" || low == "readme-casing" {
               return true
           }
       }
       return false
   }
   ```

2. **Update `RunFix` & `runFixAll`**:
   - In `RunFix`: if `isFixCasingRequest(args)` is true, route to targeted casing remediation.
   - In `runFixAll`: before or alongside pull remediation, scan tracked repositories for casing discrepancies via `DetectCasingDiscrepancies`.
   - If found, invoke `ExecuteCasingNormalization` and record recipe.
   - Render informative terminal summary pill:
     `✔ Case Normalization: Renamed README.md -> readme.md in storytelling-v1`

---

### Step 2.3: Implement Unit & Integration Tests in `cli/cmdfix/fix_casing_test.go`

1. `TestDetectCasingDiscrepancies`: Verify detection correctly flags `README.md` and `story/README.md`.
2. `TestCaseNormalizationTwoStepRename`: Mock filesystem and Git execution to verify the two-step `.tmp-lcf` sequence and rollback on step 2 failure.
3. `TestCaseNormalizationDualIndexPurge`: Verify `--cached -f` purge when both entries exist in Git index.
4. `TestRecordCasingRemediationRecipe`: Verify JSON serialization and directory creation under `repo-cache`.

---

## 3. Coding Guidelines & Invariant Verification

- [ ] **Error Sovereignty**: Leaf functions return `error` or `*apperror.AppError` directly. Zero calls to `cliexit.HandleError` or `os.Exit`.
- [ ] **Boolean Naming**: All booleans prefixed with `is*` or `has*` only (`isFixed`, `isDualTracked`, `hasFailed`).
- [ ] **Function Sizing**: <= 15 lines max per function (<= 8 lines preferred).
- [ ] **Conditional Flattening**: Zero nested `if` statements. Guard clauses and early returns throughout.
- [ ] **Path Hygiene**: Strictly relative paths in specs, plans, and output formatters.
- [ ] **Worker Git Ban**: Subagents must not run git commands directly; orchestrator handles atomic gitmap commits.

---

## 4. Verification & Validation Commands

1. **Unit Test Suite**:
   ```bash
   go test -v ./cli/cmdfix/ -run TestCasing
   ```
2. **Dry Run Simulation on Target Repository**:
   ```bash
   gitmap fix casing --repo ../storytelling-v1 --dry-run
   ```
3. **Verify Recipe in Repo-Cache**:
   Inspect `repo-cache/01-gitmap/remediation_recipes.json` for proper JSON schema and verified execution steps.
