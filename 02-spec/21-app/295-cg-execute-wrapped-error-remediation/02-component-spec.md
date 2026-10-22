# Component Specification 295: Monadic Wrapped Result Architecture

## Status: APPROVED
- **Spec ID:** 295-02
- **Slug:** `295-cg-execute-wrapped-error-remediation`
- **Domain:** Component Design & Type Refactoring

---

## 1. Component Boundaries & Target Files

```
┌──────────────────────────────────────────────────────────────────────────────┐
│ cli/store/repo_repository.go                                                 │
│ • ListByBranch, ListDirty, ListByTag, ListUnpushed, FindPathBySlug           │
│ • queryScanRecords internal helper returning result.Result[[]model.ScanRecord]│
├──────────────────────────────────────────────────────────────────────────────┤
│ cli/store/repo_query_filter.go                                               │
│ • ListReposByTag, ListReposUnpushed delegators returning result.Result[...]  │
├──────────────────────────────────────────────────────────────────────────────┤
│ cli/store/repo.go                                                            │
│ • ListReposByBranch, ListReposDirtyStatus delegators returning result.Result │
├──────────────────────────────────────────────────────────────────────────────┤
│ cli/cmdpipeline/pipeline_history_render.go                                   │
│ • queryRepoPathFromConn consuming result.Result[string]                      │
└──────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Component Details

### 2.1 `cli/store/repo_repository.go`
```go
// Package store manages repository database queries and persistence.
package store

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// RepoRepository encapsulates repository database queries returning typed Result monads.
type RepoRepository struct {
	conn sqlExecutor
}

func NewRepoRepository(conn sqlExecutor) *RepoRepository {
	return &RepoRepository{conn: conn}
}

func (r *RepoRepository) ListByBranch(branch string) result.Result[[]model.ScanRecord] {
	clean := strings.TrimSpace(branch)
	if clean == "" {
		return r.queryScanRecords(constants.SQLSelectAllRepos)
	}

	return r.queryScanRecords(constants.SQLSelectReposByBranch, clean)
}

func (r *RepoRepository) ListDirty() result.Result[[]model.ScanRecord] {
	return r.queryScanRecords(constants.SQLSelectDirtyRepos)
}

func (r *RepoRepository) ListByTag(tag string) result.Result[[]model.ScanRecord] {
	clean := strings.TrimSpace(tag)
	if clean == "" {
		return r.queryScanRecords(constants.SQLSelectAllRepos)
	}

	pattern := "%" + clean + "%"

	return r.queryScanRecords(constants.SQLSelectReposByTag, pattern)
}

func (r *RepoRepository) ListUnpushed() result.Result[[]model.ScanRecord] {
	return r.queryScanRecords(constants.SQLSelectUnpushedRepos)
}

func (r *RepoRepository) FindPathBySlug(slug string) result.Result[string] {
	rq, ok := r.conn.(sqlRowQueryer)
	if !ok {
		return result.Fail[string](apperror.NewSimple("database executor does not support row query", "E_DB_EXECUTOR"))
	}

	parts := strings.Split(slug, "/")
	shortSlug := parts[len(parts)-1]
	row := QueryRowWrapper(rq, constants.SQLSelectRepoPathBySlug, slug, shortSlug)

	var absPath string
	scanErr := row.Scan(&absPath)
	if scanErr != nil {
		return result.Fail[string](apperror.WrapSimple(scanErr, "store.RepoRepository.FindPathBySlug"))
	}

	return result.Ok(absPath)
}

func (r *RepoRepository) queryScanRecords(query string, args ...any) result.Result[[]model.ScanRecord] {
	q, ok := r.conn.(sqlQueryer)
	if !ok {
		return result.Fail[[]model.ScanRecord](apperror.NewSimple("database executor does not support query", "E_DB_EXECUTOR"))
	}

	res := QueryWrapper(q, query, args...)
	rows, err := res.Destruct()
	if err != nil {
		return result.Fail[[]model.ScanRecord](apperror.WrapSimple(err, "store.RepoRepository.queryScanRecords"))
	}
	defer rows.Close()

	records, scanErr := scanRows(rows)
	if scanErr != nil {
		return result.Fail[[]model.ScanRecord](apperror.WrapSimple(scanErr, "store.RepoRepository.scanRows"))
	}

	return result.Ok(records)
}
```

### 2.2 `cli/store/repo_query_filter.go`
```go
package store

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// ListReposByTag returns all repositories matching a specific tag stored in metadata or notes.
func (db *DB) ListReposByTag(tag string) result.Result[[]model.ScanRecord] {
	return NewRepoRepository(db.conn).ListByTag(tag)
}

// ListReposUnpushed returns repositories flagged with unpushed commits in metadata or local tracking.
func (db *DB) ListReposUnpushed() result.Result[[]model.ScanRecord] {
	return NewRepoRepository(db.conn).ListUnpushed()
}
```

### 2.3 `cli/store/repo.go`
```go
// ListReposByBranch returns all repositories configured for the specified branch.
func (db *DB) ListReposByBranch(branch string) result.Result[[]model.ScanRecord] {
	return NewRepoRepository(db.conn).ListByBranch(branch)
}

// ListReposDirtyStatus returns repositories flagged as dirty or requiring synchronization.
func (db *DB) ListReposDirtyStatus() result.Result[[]model.ScanRecord] {
	return NewRepoRepository(db.conn).ListDirty()
}
```

### 2.4 `cli/cmdpipeline/pipeline_history_render.go`
```go
func queryRepoPathFromConn(db *store.DB, slug string) string {
	if db == nil {
		return ""
	}
	defer db.Close()

	repo := store.NewRepoRepository(db.Conn())
	pathRes := repo.FindPathBySlug(slug)
	if pathRes.HasError() {
		logRepoLookupFailure(slug, pathRes.Err)

		return ""
	}

	return pathRes.Value
}
```
