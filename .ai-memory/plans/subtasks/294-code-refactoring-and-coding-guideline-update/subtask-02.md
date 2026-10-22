# Subtask 02: Repository Pattern Extraction & Pipeline SQL Decoupling

## Status: COMPLETED
- **Parent Plan**: `.ai-memory/plans/294-code-refactoring-and-coding-guideline-update.md`
- **Architecture Spec**: `02-spec/21-app/294-code-refactoring-and-coding-guideline-update/01-architecture-spec.md`
- **Owner**: Worker 02
- **Owned Files**:
  - `cli/constants/constants_store.go`
  - `cli/store/repo_repository.go`
  - `cli/store/repo.go`
  - `cli/store/repo_query_filter.go`
  - `cli/cmdpipeline/pipeline_history_render.go`

---

## 1. Context & Objectives

This subtask encapsulates the repository data access layer under the Repository Pattern, eliminates inline SQL duplicates across store modules, enforces typed `*apperror.AppError` returns, and decouples ad-hoc SQL execution from the pipeline history rendering subsystem.

### Key Objectives:
1. In `cli/constants/constants_store.go`: add centralized SQL constants `SQLSelectReposByBranch`, `SQLSelectDirtyRepos`, `SQLSelectReposByTag`, and `SQLSelectUnpushedRepos`.
2. In `cli/store/repo_repository.go`: create dedicated `RepoRepository` returning `*apperror.AppError`.
3. In `cli/store/repo.go` & `cli/store/repo_query_filter.go`: remove duplicated local SQL constants, delegate query execution to centralized constants and the repository abstraction, and preserve error return sovereignty.
4. In `cli/cmdpipeline/pipeline_history_render.go`: eliminate direct `sql.Open` connection cycles and raw SQL string `"SELECT AbsolutePath FROM Repo..."`, query through the store repository layer, handle errors without swallowing, and refactor `parseFirstVTag` to use `strutil` helpers.

---

## 2. Invariants & Guardrails

- **Function Length Limit:** Maximum 15 lines per function (<= 8 lines preferred).
- **Conditionals:** Zero nested `if` statements; apply guard clauses and early returns.
- **Boolean Naming:** Use positive prefixes `is*` or `has*` only. No negative booleans.
- **Error Return Sovereignty:** Leaf functions must return `*apperror.AppError` or standard `error` directly; zero `os.Exit` or `cliexit.HandleError` in leaf code.
- **Path Invariants:** Strictly relative Git paths only. No absolute paths or URI schemes.
- **Safety Gate:** Zero builds (`go build` prohibited) and zero tests (`go test` prohibited).

---

## 3. Step-by-Step Implementation Plan

### Step 1: Centralize SQL Constants in `cli/constants/constants_store.go`
- **Location:** `cli/constants/constants_store.go`
- **Additions:**
  ```go
  // SQL queries for repository filtering and status queries.
  const (
  	SQLSelectReposByBranch = "SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport FROM Repo WHERE Branch = ? COLLATE NOCASE ORDER BY Slug ASC"
  	SQLSelectDirtyRepos    = "SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport FROM Repo WHERE Notes LIKE '%dirty%' OR RepoName IN (SELECT RepoSlug FROM HeatmapLocalTracking WHERE IsDirty = 1) ORDER BY Slug ASC"
  	SQLSelectReposByTag    = "SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport FROM Repo WHERE Notes LIKE ? COLLATE NOCASE ORDER BY Slug ASC"
  	SQLSelectUnpushedRepos = "SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport FROM Repo WHERE (Notes LIKE '%unpushed%' OR Notes LIKE '%ahead%') OR RepoName IN (SELECT RepoSlug FROM HeatmapLocalTracking WHERE IsDirty = 1) ORDER BY Slug ASC"
  	SQLSelectRepoPathBySlug = "SELECT AbsolutePath FROM Repo WHERE Slug = ? OR Slug = ? LIMIT 1"
  )
  ```

### Step 2: Create `cli/store/repo_repository.go`
- **Location:** `cli/store/repo_repository.go`
- **Interface & Structure:**
  ```go
  package store

  import (
  	"database/sql"
  	"strings"

  	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
  	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
  	"github.com/alimtvnetwork/gitmap-v28/cli/model"
  )

  // RepoRepository encapsulates repository database queries returning typed AppErrors.
  type RepoRepository interface {
  	FindByBranch(branch string) ([]model.ScanRecord, *apperror.AppError)
  	FindDirty() ([]model.ScanRecord, *apperror.AppError)
  	FindByTag(tag string) ([]model.ScanRecord, *apperror.AppError)
  	FindUnpushed() ([]model.ScanRecord, *apperror.AppError)
  	FindPathBySlug(slug string) (string, *apperror.AppError)
  }

  // DefaultRepoRepository implements RepoRepository using a database connection.
  type DefaultRepoRepository struct {
  	conn *sql.DB
  }

  // NewRepoRepository constructs a new DefaultRepoRepository.
  func NewRepoRepository(conn *sql.DB) *DefaultRepoRepository {
  	return &DefaultRepoRepository{conn: conn}
  }
  ```
- **Implementation Methods (each <= 15 lines):**
  - `FindByBranch(branch string) ([]model.ScanRecord, *apperror.AppError)`:
    - Validates trimmed branch; if empty, delegates to `ListRepos()`.
    - Queries using `constants.SQLSelectReposByBranch`.
    - Wraps any database error as `apperror.WrapSimple(err, "store.FindByBranch")`.
  - `FindDirty() ([]model.ScanRecord, *apperror.AppError)`:
    - Queries using `constants.SQLSelectDirtyRepos`.
    - Wraps any database error as `apperror.WrapSimple(err, "store.FindDirty")`.
  - `FindByTag(tag string) ([]model.ScanRecord, *apperror.AppError)`:
    - Queries using `constants.SQLSelectReposByTag` with pattern `"%"+tag+"%"`.
    - Wraps any database error as `apperror.WrapSimple(err, "store.FindByTag")`.
  - `FindUnpushed() ([]model.ScanRecord, *apperror.AppError)`:
    - Queries using `constants.SQLSelectUnpushedRepos`.
    - Wraps any database error as `apperror.WrapSimple(err, "store.FindUnpushed")`.
  - `FindPathBySlug(slug string) (string, *apperror.AppError)`:
    - Computes `shortSlug` from trailing path element of `slug`.
    - Queries `constants.SQLSelectRepoPathBySlug`.
    - Returns `(absPath, nil)` on success; wraps `sql.ErrNoRows` or query failures cleanly.

### Step 3: Refactor `cli/store/repo.go` & `cli/store/repo_query_filter.go`
- **In `cli/store/repo.go`:**
  - Remove local `const sqlSelectReposByBranch` (line 128) and `const sqlSelectDirtyRepos` (line 147).
  - Update `ListReposByBranch` to query using `constants.SQLSelectReposByBranch` or delegate to `NewRepoRepository(db.conn).FindByBranch(branch)`.
  - Update `ListReposDirtyStatus` to query using `constants.SQLSelectDirtyRepos` or delegate to `NewRepoRepository(db.conn).FindDirty()`.
- **In `cli/store/repo_query_filter.go`:**
  - Remove local `const sqlSelectReposByTag` (line 12) and `const sqlSelectUnpushedRepos` (line 14).
  - Update `ListReposByTag` to query using `constants.SQLSelectReposByTag` or delegate to `NewRepoRepository(db.conn).FindByTag(tag)`.
  - Update `ListReposUnpushed` to query using `constants.SQLSelectUnpushedRepos` or delegate to `NewRepoRepository(db.conn).FindUnpushed()`.
  - Ensure leaf error sovereignty is maintained throughout.

### Step 4: Refactor `cli/cmdpipeline/pipeline_history_render.go`
- **Location:** `cli/cmdpipeline/pipeline_history_render.go`
- **Refactoring `findLocalRepoDirForSlug`:**
  - Replace direct `sql.Open("sqlite", ...)` with store connection access:
    ```go
    func findLocalRepoDirForSlug(slug string) string {
    	db, err := store.OpenDefault()
    	if err != nil {
    		return ""
    	}
    	defer db.Close()

    	repo := store.NewRepoRepository(db.Conn())
    	absPath, appErr := repo.FindPathBySlug(slug)
    	if appErr != nil {
    		return ""
    	}

    	return absPath
    }
    ```
  - Delete `queryRepoPathFromConn` function and eliminate the inline SQL query string.
- **Refactoring `parseFirstVTag`:**
  - Replace naïve newline splitting and index access:
    ```go
    func parseFirstVTag(output string) string {
    	trimmed := strings.TrimSpace(output)
    	if len(trimmed) == 0 {
    		return ""
    	}

    	lines := strings.Split(trimmed, "\n")
    	for _, line := range lines {
    		candidate := strings.TrimSpace(line)
    		hasVPrefix := strings.HasPrefix(strings.ToLower(candidate), "v")
    		if hasVPrefix {
    			return candidate
    		}
    	}

    	return ""
    }
    ```
  - Ensures clean guard clauses, positive booleans, and <= 15 lines.

---

## 4. Verification & Quality Checklist

- [x] `SQLSelectReposByBranch`, `SQLSelectDirtyRepos`, `SQLSelectReposByTag`, `SQLSelectUnpushedRepos`, and `SQLSelectRepoPathBySlug` added to `cli/constants/constants_store.go`.
- [x] `RepoRepository` interface and `DefaultRepoRepository` created in `cli/store/repo_repository.go` returning `*apperror.AppError`.
- [x] Duplicated SQL constants removed from `cli/store/repo.go` and `cli/store/repo_query_filter.go`.
- [x] Raw SQL and direct `sql.Open` eliminated from `cli/cmdpipeline/pipeline_history_render.go`.
- [x] `parseFirstVTag` refactored with clean guard clauses and positive booleans.
- [x] All functions <= 15 lines (<= 8 lines preferred).
- [x] Zero nested `if` statements.
- [x] Positive booleans (`is*`, `has*` only).
- [x] Strictly relative Git paths only.
- [x] Zero compiler builds and zero tests executed.
