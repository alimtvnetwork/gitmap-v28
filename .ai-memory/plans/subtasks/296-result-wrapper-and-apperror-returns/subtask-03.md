# Subtask 03: Store Repository Monadic Container Refactoring

> **Subtask ID:** `SUBTASK-296-03`  
> **Parent Task:** `296-result-wrapper-and-apperror-returns`  
> **Status:** Pending  
> **Target Files:**  
> - `cli/store/repo.go` (modify)  
> **Related Architecture Spec:** `02-spec/21-app/296-result-wrapper-and-apperror-returns/01-architecture-spec.md`  

---

## 1. Objective

Refactor multi-value error tuple methods in `cli/store/repo.go` (`ListRepos`, `FindBySlug`, `FindOneBySlug`, `FindByPath`, `CountRepos`, `DeleteByPath`, `DeleteBySlug`) to return canonical monadic Result containers (`ScanRecordSliceWrap`, `ScanRecordWrap`, `RepoCountWrap`, `RowsAffectedWrap`). This eliminates the anti-pattern of multi-value error tuples `(T, error)` across the store subsystem, enforces single-return sovereignty, and enables safe pointer-attached predicate inspection.

---

## 2. Action Checklist

- [ ] 1. Refactor query methods in `cli/store/repo.go`:
  - [ ] `ListRepos() ScanRecordSliceWrap`:
    - Replaces `func (db *DB) ListRepos() ([]model.ScanRecord, error)`.
    - Query `constants.SQLSelectAllRepos`. On query error, return `result.FailSlice[model.ScanRecord](apperror.WrapSimple(err, "store.ListRepos.query"))`.
    - Scan records using `scanRows(rows)`. On scan error, return `result.FailSlice[model.ScanRecord](apperror.WrapSimple(scanErr, "store.ListRepos.scan"))`.
    - On success, return `result.OkSlice(records)`.
  - [ ] `FindBySlug(slug string) ScanRecordSliceWrap`:
    - Replaces `func (db *DB) FindBySlug(slug string) ([]model.ScanRecord, error)`.
    - Delegates to `findBySlugRunner(db.conn, slug)`.
    - On query error, return `result.FailSlice[model.ScanRecord](...)`.
    - On success, return `result.OkSlice(records)`.
  - [ ] `FindOneBySlug(slug string) ScanRecordWrap`:
    - Replaces `func (db *DB) FindOneBySlug(slug string) (*model.ScanRecord, error)`.
    - Call `db.FindBySlug(slug)`.
    - If `sliceRes.IsFailure()`, propagate error: `result.Fail[model.ScanRecord](sliceRes.AppError())`.
    - If `sliceRes.IsEmpty() || sliceRes.Count() == 0`, return `result.Fail[model.ScanRecord](NewRepoNotFoundError("FindOneBySlug", "Repo", slug))`.
    - On success, return `result.Ok(sliceRes.Data()[0])`.
  - [ ] `FindByPath(absPath string) ScanRecordSliceWrap`:
    - Replaces `func (db *DB) FindByPath(absPath string) ([]model.ScanRecord, error)`.
    - Normalize path with `NormalizeStoragePath(absPath)`.
    - Execute query `constants.SQLSelectRepoByPath`.
    - On error, return `result.FailSlice[model.ScanRecord](apperror.WrapSimple(err, "store.FindByPath"))`.
    - On success, return `result.OkSlice(records)`.
  - [ ] `CountRepos() RepoCountWrap`:
    - Replaces `func (db *DB) CountRepos() (int, error)`.
    - Execute `SELECT COUNT(*) FROM Repo`.
    - On query/scan error, return `result.Fail[int](apperror.WrapSimple(err, "store.CountRepos"))`.
    - On success, return `result.Ok(total)`.

- [ ] 2. Refactor mutation methods in `cli/store/repo.go`:
  - [ ] `DeleteByPath(absPath string) RowsAffectedWrap`:
    - Replaces `func (db *DB) DeleteByPath(absPath string) (int64, error)`.
    - Normalize path with `NormalizeStoragePath(absPath)`.
    - Execute `constants.SQLDeleteRepoByPath`.
    - On exec error: `result.Fail[int64](apperror.WrapSimple(err, "store.DeleteByPath.exec"))`.
    - On rows affected error: `result.Fail[int64](apperror.WrapSimple(rowsErr, "store.DeleteByPath.rows"))`.
    - On success: `result.Ok(rows)`.
  - [ ] `DeleteBySlug(slug string) RowsAffectedWrap`:
    - Replaces `func (db *DB) DeleteBySlug(slug string) (int64, error)`.
    - Execute `constants.SQLDeleteRepoBySlug`.
    - On exec error: `result.Fail[int64](apperror.WrapSimple(err, "store.DeleteBySlug.exec"))`.
    - On rows affected error: `result.Fail[int64](apperror.WrapSimple(rowsErr, "store.DeleteBySlug.rows"))`.
    - On success: `result.Ok(rows)`.

- [ ] 3. Refactor internal helpers in `cli/store/repo.go`:
  - [ ] Update `scanRows` signature to return `(ScanRecordSlice, *apperror.AppError)` for clean monadic assembly.
  - [ ] Update `findBySlugRunner(runner sqlQueryer, slug string) ScanRecordSliceWrap`.
  - [ ] Update `findFirstRepoRecord(repoPath string)` to inspect `db.FindByPath(repoPath)` with `IsFailure()` / `IsEmpty()`.

- [ ] 4. Coding Standard & Invariant Verification:
  - [ ] Use centralized aliases from `cli/store/types.go` (`ScanRecordSliceWrap`, `ScanRecordWrap`, `RepoCountWrap`, `RowsAffectedWrap`).
  - [ ] Never leak raw generic brackets like `result.Wrap[model.ScanRecord]` in method signatures.
  - [ ] Zero double-wrapping: Never call `FailWrap` (does not exist) or double-wrap `*apperror.AppError`.
  - [ ] Functions strictly <= 15 lines with 0 nested ifs.
  - [ ] Strictly relative Git paths only.
  - [ ] Zero compiler builds and zero tests (`go build` / `go test` strictly banned).
  - [ ] Zero git commands.

---

## 3. Concrete Code Specifications

### Method Implementations for `cli/store/repo.go`

```go
// DeleteByPath removes the repo row whose AbsolutePath matches.
// Returns the number of rows deleted encapsulated in RowsAffectedWrap.
func (db *DB) DeleteByPath(absPath string) RowsAffectedWrap {
	cleanPath := NormalizeStoragePath(absPath)
	res, err := ExecWrapper(db.conn, constants.SQLDeleteRepoByPath, cleanPath).Destruct()
	if err != nil {
		return result.Fail[int64](apperror.WrapSimple(err, "store.DeleteByPath"))
	}

	rows, rowsErr := res.RowsAffected()
	if rowsErr != nil {
		return result.Fail[int64](apperror.WrapSimple(rowsErr, "store.DeleteByPath.rows"))
	}

	return result.Ok(rows)
}

// DeleteBySlug removes every repo row whose Slug matches.
// Returns the number of rows deleted encapsulated in RowsAffectedWrap.
func (db *DB) DeleteBySlug(slug string) RowsAffectedWrap {
	res, err := ExecWrapper(db.conn, constants.SQLDeleteRepoBySlug, slug).Destruct()
	if err != nil {
		return result.Fail[int64](apperror.WrapSimple(err, "store.DeleteBySlug"))
	}

	rows, rowsErr := res.RowsAffected()
	if rowsErr != nil {
		return result.Fail[int64](apperror.WrapSimple(rowsErr, "store.DeleteBySlug.rows"))
	}

	return result.Ok(rows)
}

// ListRepos returns all tracked repositories ordered by slug.
func (db *DB) ListRepos() ScanRecordSliceWrap {
	rows, err := QueryWrapper(db.conn, constants.SQLSelectAllRepos).Destruct()
	if err != nil {
		return result.FailSlice[model.ScanRecord](apperror.WrapSimple(err, "store.ListRepos.query"))
	}

	defer rows.Close()

	records, scanErr := scanRows(rows)
	if scanErr != nil {
		return result.FailSlice[model.ScanRecord](scanErr)
	}

	return result.OkSlice(records)
}

// CountRepos returns the total number of repositories tracked in the database.
func (db *DB) CountRepos() RepoCountWrap {
	var total int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM Repo").Scan(&total)
	if err != nil {
		return result.Fail[int](apperror.WrapSimple(err, "store.CountRepos"))
	}

	return result.Ok(total)
}

func findBySlugRunner(runner sqlQueryer, slug string) ScanRecordSliceWrap {
	rows, err := QueryWrapper(runner, constants.SQLSelectRepoBySlug, slug).Destruct()
	if err != nil {
		return result.FailSlice[model.ScanRecord](apperror.WrapSimple(err, "store.findBySlugRunner.query"))
	}

	defer rows.Close()

	records, scanErr := scanRows(rows)
	if scanErr != nil {
		return result.FailSlice[model.ScanRecord](scanErr)
	}

	return result.OkSlice(records)
}

// FindBySlug returns all repos matching the given slug.
func (db *DB) FindBySlug(slug string) ScanRecordSliceWrap {
	return findBySlugRunner(db.conn, slug)
}

// FindOneBySlug returns the first repository matching the given slug.
func (db *DB) FindOneBySlug(slug string) ScanRecordWrap {
	sliceRes := db.FindBySlug(slug)
	if sliceRes.IsFailure() {
		return result.Fail[model.ScanRecord](sliceRes.AppError())
	}

	if sliceRes.IsEmpty() || sliceRes.Count() == 0 {
		return result.Fail[model.ScanRecord](NewRepoNotFoundError("FindOneBySlug", "Repo", slug))
	}

	return result.Ok(sliceRes.Data()[0])
}

// FindByPath returns the repo at the given absolute path.
func (db *DB) FindByPath(absPath string) ScanRecordSliceWrap {
	cleanPath := NormalizeStoragePath(absPath)
	rows, err := QueryWrapper(db.conn, constants.SQLSelectRepoByPath, cleanPath).Destruct()
	if err != nil {
		return result.FailSlice[model.ScanRecord](apperror.WrapSimple(err, "store.FindByPath.query"))
	}

	defer rows.Close()

	records, scanErr := scanRows(rows)
	if scanErr != nil {
		return result.FailSlice[model.ScanRecord](scanErr)
	}

	return result.OkSlice(records)
}
```

### Internal Helper Refactoring (`scanRows`)

```go
// scanRows reads ScanRecord values from query result rows.
func scanRows(rows interface {
	Next() bool
	Scan(dest ...any) error
}) (ScanRecordSlice, *apperror.AppError) {
	var results ScanRecordSlice

	for rows.Next() {
		r, err := scanOneRow(rows)
		if err != nil {
			return nil, apperror.WrapSimple(err, "store.scanRows.scanOneRow")
		}

		results = append(results, r)
	}

	return results, nil
}
```

---

## 4. Verification & Done Criteria

- [ ] All 7 methods in `cli/store/repo.go` return single-return monadic containers (`ScanRecordSliceWrap`, `ScanRecordWrap`, `RepoCountWrap`, `RowsAffectedWrap`).
- [ ] No multi-value error tuples `(T, error)` remain on target methods.
- [ ] No double-wrapping (`FailWrap` or nested `apperror.WrapSimple`).
- [ ] Zero compiler builds and zero tests executed.
- [ ] Zero git commands executed.
- [ ] Strictly relative Git paths used throughout.
