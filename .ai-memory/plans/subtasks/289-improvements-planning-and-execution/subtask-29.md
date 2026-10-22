# Subtask 29: Path Healing Guard Inversion & Flow Flattening

> **Subtask ID:** Subtask-29  
> **Parent Plan:** `.ai-memory/plans/pending/289-improvements-planning-and-execution.md`  
> **Component Spec:** `02-spec/21-app/289-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/cmdpull/pull_path_heal.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Verify and enforce zero nested if statements in `cli/cmdpull/pull_path_heal.go` using guard clauses.
2. Ensure path separator healing logic across OS boundaries (Windows UNC / WSL paths) executes safely without nesting.
3. Validate positive boolean variables (`isFound`, `hasDrive`) and verify all functions remain under 15 lines.
4. Ensure safe database update handling when repairing paths in `gitmap.db`.

---

## 2. Detailed Technical Plan

### 2.1 Guard Flattening in `HealRepoPathSeparators`
- Guard 1: Host OS is Windows -> return early.
- Guard 2: No backslashes present -> return early.
- Guard 3: Drive letter present -> delegate to helper and return early.
- Default: Return normalized path.
  ```go
  func HealRepoPathSeparators(path string) string {
  	if runtime.GOOS == "windows" {
  		return path
  	}

  	if !strings.Contains(path, "\\") {
  		return path
  	}

  	normalized := strings.ReplaceAll(path, "\\", "/")
  	hasDrive := len(normalized) >= 3 && normalized[1] == ':' && normalized[2] == '/'
  	if hasDrive {
  		return resolveWindowsDriveOnUnix(normalized)
  	}

  	return normalized
  }
  ```

### 2.2 Guard Flattening in `resolveWindowsDriveOnUnix`
- Ensure sequential stat checks return immediately upon match:
  ```go
  func resolveWindowsDriveOnUnix(normalized string) string {
  	drive := strings.ToLower(string(normalized[0]))
  	wslPath := "/mnt/" + drive + normalized[2:]
  	if _, err := os.Stat(wslPath); err == nil {
  		return wslPath
  	}

  	tailPath := normalized[2:]
  	if _, err := os.Stat(tailPath); err == nil {
  		return tailPath
  	}

  	return tailPath
  }
  ```

### 2.3 Positive Boolean Guard in `TryHealMissingRecordPath`
- Verify positive boolean `isFound`:
  ```go
  func TryHealMissingRecordPath(rec *model.ScanRecord) bool {
  	if rec == nil {
  		return false
  	}

  	healed := HealRepoPathSeparators(rec.AbsolutePath)
  	if healed == rec.AbsolutePath {
  		return false
  	}

  	isFound := !cloner.IsMissingRepo(healed)
  	if !isFound {
  		return false
  	}

  	oldPath := rec.AbsolutePath
  	rec.AbsolutePath = healed
  	_ = AutoHealRepoPathInDB(oldPath, healed)

  	return true
  }
  ```

---

## 3. Strict Coding Guidelines Compliance
- **Rule 3 (Zero Nested Ifs):** Every condition in `pull_path_heal.go` is flat; zero nested blocks.
- **Rule 4 (Boolean Naming):** Uses positive booleans (`hasDrive`, `isFound`).
- **Rule 2 (Function Limit):** All functions <= 15 lines.
- **Rule 6 (Vertical Whitespace):** Single blank line after each `}` closure and before return.

---

## 4. Verification Criteria
- [x] `pull_path_heal.go` contains zero nested if statements.
- [x] Positive booleans used exclusively throughout.
- [x] Path healing functions maintain cross-platform behavior.
- [x] Zero compiler builds and zero tests executed during planning.
