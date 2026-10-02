# Subtask 65.6: .gitignore Policy Engine: Retain .gitmap/ & Exclude Only .gitmap/backup/

> **Parent Plan:** [65-gitmap-update-all-zip-and-fixes.md](../../pending/65-gitmap-update-all-zip-and-fixes.md)  
> **Tracking Spec:** [02-component-spec.md](../../../../02-spec/21-app/65-gitmap-update-all-zip-and-fixes/02-component-spec.md)  
> **Status:** Pending  
> **Primary File Targets:** `cli/cmdpull/pull.go`, `cli/cmdignore/ignore_groups.go`, `cli/cmdignore/fix_ignore.go`  

---

## 1. Objective

Refactor GitMap's `.gitignore` policy, pull verification, and automated cleanup engines:
1. Ensure `.gitmap/` is retained and allowed to be tracked in git repository trees.
2. Only `.gitmap/backup/` (and ephemeral resume task JSONs) is ignored by default.
3. Remove mandatory forced injection of `.gitmap/` into `.gitignore`.
4. Eliminate duplicate ignore patterns by normalizing lines (trimming both leading and trailing slashes).
5. Remove obsolete `isGitmapDevelopmentRepo` logic in pull inspection.

---

## 2. Implementation Scope

### 2.1 Pull Tracked Path Inspection (`cli/cmdpull/pull.go`)
- **`buildDefaultTrackedPathsToCheck(repoDir string)`:**
  - Remove `.gitmap/` from default paths.
  - Return exclusively `.gitmap/backup/` along with the standard resume task JSON file candidates.
- **Remove Obsolete Development Repo Guard:**
  - Remove `isGitmapDevelopmentRepo(repoDir string)` helper function so all repositories follow the same consistent rule.

### 2.2 Ignore Groups & Immutable Rules (`cli/cmdignore/ignore_groups.go`)
- **`ImmutableDefaultRules`:**
  - Update slice to contain only `".gitmap/backup/"` (remove `".gitmap/"`).
- **`EnsureImmutableRules(patterns []string)`:**
  - Ensure presence of `".gitmap/backup/"` only.
  - Do not automatically prepend or enforce `".gitmap/"`.

### 2.3 Ignore Inspection & Deduplication (`cli/cmdignore/fix_ignore.go`)
- **`analyzeGitignoreData`:**
  - Stop flagging missing `.gitmap/` as an issue.
  - Optionally check for missing `.gitmap/backup/`.
- **`assembleCleanedGitignore`:**
  - When appending default backup persistence, add only `".gitmap/backup/"`.
  - Never add `".gitmap/"`.
- **Deduplication Normalization (`processIgnoreLine`):**
  - Normalize line patterns using `strings.Trim(trimmed, "/ \t\r\n")` to recognize that `dist/`, `/dist`, and `dist` represent identical patterns.
  - Increment duplicate counter and prune duplicate occurrences.

---

## 3. Coding Guidelines & Constraints
- All boolean conditions and flags must be positive.
- Functions must remain small (<= 15 lines per function) and modular.
- Do not remove or alter legitimate user-defined ignore rules.
- Maintain existing Split-DB caching and performance characteristics.

---

## 4. Verification Steps
- Run tests in `cli/cmdignore`: `go test ./cli/cmdignore/...`.
- Run tests in `cli/cmdpull`: `go test ./cli/cmdpull/...`.
- Verify with a test `.gitignore` containing `/dist` and `dist/` that only one entry is preserved and duplicate count is reported correctly.
- Verify that a repository containing tracked files inside `.gitmap/` (e.g., `.gitmap/config.json`) is not flagged as an ignore violation during `gitmap pull` or `gitmap fix-ignore-all`.
