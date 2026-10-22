# Subtask Plan 04: Propagate Directory Scanning & I/O Errors in `cli/cmddb/db_doctor.go`

- **Subtask ID:** `subtask-04`
- **Subtask Name:** `04-db-doctor-io-error-propagation`
- **Parent Spec:** `02-spec/21-app/290-improvements-planning-and-execution/01-architecture-spec.md`
- **Assigned Agent Role:** Spec Author 01 / Implementation Worker 01
- **Target Files:**
  - `cli/cmddb/db_doctor.go`
- **Status:** Completed

---

## 1. Targeted Problem & Architectural Rationale

In `cli/cmddb/db_doctor.go`, database health inspection functions (`collectDoctorItems`, `statDoctorItem`, `dirDoctorSize`) silently swallow I/O and directory traversal errors:

1. In `dirDoctorSize(dir string) int64`:
   ```go
   entries, readErr := os.ReadDir(dir)
   if readErr != nil {
       return 0
   }
   ```
   When `os.ReadDir(dir)` fails due to permission denial, filesystem locks, or corruption, `dirDoctorSize` silently returns `0`. The doctor report falsely indicates an empty directory rather than an unreadable directory.
2. In `statDoctorItem(name, path string, isDir bool) doctorItem`:
   ```go
   info, statErr := os.Stat(path)
   if statErr != nil {
       return doctorItem{name: name, path: path, isDir: isDir, isFound: false}
   }
   ```
   Any `statErr` (including disk faults or permission errors) is collapsed into `isFound: false`. A database file that exists but cannot be accessed is wrongly reported as `MISSING`.
3. To uphold **Rule 1 (Error Return Sovereignty)** and **Rule 4 (Positive Booleans)**:
   - `dirDoctorSize` must return `(int64, error)`.
   - `doctorItem` must track access faults with positive flags (`hasError bool`, `errDesc string`).
   - `doctorStatus` must report `ERROR` when I/O operations fail, alerting the user to filesystem permission or corruption issues.

---

## 2. Component Boundaries & Code Signatures

### 2.1 Enhance `doctorItem` Struct
- **File:** `cli/cmddb/db_doctor.go`
- **Requirements:**
  - Add positive boolean `hasError bool`.
  - Add `errDesc string` to describe the captured fault.

```go
type doctorItem struct {
	name     string
	path     string
	size     int64
	isDir    bool
	isFound  bool
	hasError bool
	errDesc  string
	modified time.Time
}
```

### 2.2 Refactor `dirDoctorSize`
- **File:** `cli/cmddb/db_doctor.go`
- **Signature:** `func dirDoctorSize(dir string) (int64, error)`
- **Requirements:**
  - Propagate `readErr` from `os.ReadDir(dir)` directly instead of returning 0.
  - Return `(total, nil)` on success.
  - Function length <= 15 lines.

```go
func dirDoctorSize(dir string) (int64, error) {
	var total int64
	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		return 0, readErr
	}
	for _, e := range entries {
		info, infoErr := e.Info()
		if infoErr != nil {
			continue
		}
		total += info.Size()
	}

	return total, nil
}
```

### 2.3 Refactor `statDoctorItem`
- **File:** `cli/cmddb/db_doctor.go`
- **Signature:** `func statDoctorItem(name, path string, isDir bool) doctorItem`
- **Requirements:**
  - Differentiate between non-existent files (`os.IsNotExist(statErr)`) and access errors.
  - Call `dirDoctorSize(path)` and capture directory reading errors.
  - Set `hasError: true` and `errDesc: err.Error()` when I/O fails.
  - Keep function length <= 15 lines.

```go
func statDoctorItem(name, path string, isDir bool) doctorItem {
	info, statErr := os.Stat(path)
	if statErr != nil {
		isMissing := os.IsNotExist(statErr)
		return doctorItem{name: name, path: path, isDir: isDir, isFound: !isMissing, hasError: !isMissing, errDesc: statErr.Error()}
	}

	size := info.Size()
	if isDir {
		dirSize, dirErr := dirDoctorSize(path)
		if dirErr != nil {
			return doctorItem{name: name, path: path, isDir: true, isFound: true, hasError: true, errDesc: dirErr.Error(), modified: info.ModTime()}
		}
		size = dirSize
	}

	return doctorItem{name: name, path: path, size: size, isDir: isDir, isFound: true, modified: info.ModTime()}
}
```

### 2.4 Update `doctorStatus`
- **File:** `cli/cmddb/db_doctor.go`
- **Signature:** `func doctorStatus(it doctorItem) string`
- **Requirements:**
  - Return `"ERROR"` if `it.hasError`.
  - Return `"MISSING"` if `!it.isFound`.
  - Return `"STALE"` if modified duration exceeds `staleAfter`.
  - Return `"OK"` otherwise.
  - Function length <= 12 lines.

```go
func doctorStatus(it doctorItem) string {
	if it.hasError {
		return "ERROR"
	}
	if !it.isFound {
		return "MISSING"
	}
	if time.Since(it.modified) > staleAfter {
		return "STALE"
	}

	return "OK"
}
```

---

## 3. Implementation Verification Checklist

- [x] `doctorItem` struct has `hasError bool` and `errDesc string`.
- [x] `dirDoctorSize` returns `(int64, error)` and no longer swallows directory read failures.
- [x] `statDoctorItem` distinguishes `os.IsNotExist` from genuine I/O read/stat errors.
- [x] `doctorStatus` returns `"ERROR"` when `hasError` is true.
- [x] All functions <= 15 lines.
- [x] Positive booleans (`hasError`, `isFound`, `isDir`) enforced.
- [x] Blank line before `return` and after `}` preserved.
