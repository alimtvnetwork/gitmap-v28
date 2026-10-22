# Subtask Plan 06: Convert Directory Slug Resolution to Typed Result Wrapper

- **Subtask ID:** `subtask-06`
- **Subtask Name:** `06-cmdpipeline-target-resolve-result-wrapper`
- **Parent Spec:** `02-spec/21-app/289-improvements-planning-and-execution/01-architecture-spec.md`
- **Assigned Agent Role:** Spec Author 01 / Implementation Worker 01
- **Target Files:**
  - `cli/cmdpipeline/pipeline_target_resolve.go`
- **Status:** Completed

---

## 1. Targeted Problem & Architectural Rationale

In `cli/cmdpipeline/pipeline_target_resolve.go`, `resolveDirectorySlug(dirPath string) string` attempts to resolve a directory path to a GitHub canonical repository slug:

1. It attempts to read the origin remote URL via `gitutil.GetRemoteOriginURL(dirPath)`.
2. If absent, it queries `gitmap.db` by directory base name via `queryRepoSlugFromDB(base)`.
3. If neither succeeds, it silently defaults to:
   `constants.DefaultGitHubOwner + "/" + base`
4. The caller has no mechanism to determine whether the returned slug was verified from Git metadata, loaded from the database, or merely synthesized via fallback guessing.
5. In addition, errors during Git execution or database query are discarded.

Under **Rule 1 (Error Return Sovereignty)** and **Rule 4 (Positive Booleans)**:
- Define a strongly typed `DirectorySlugResult` struct adhering to positive boolean conventions (`isFound`, `isGuessed`).
- Refactor `resolveDirectorySlug` to return `DirectorySlugResult`.
- Ensure callers can inspect resolution provenance while preserving seamless backwards compatibility.

---

## 2. Component Boundaries & Code Signatures

### 2.1 Declare `DirectorySlugResult` Struct
- **File:** `cli/cmdpipeline/pipeline_target_resolve.go`
- **Definition:**
  ```go
  // DirectorySlugResult wraps resolved directory slug provenance.
  type DirectorySlugResult struct {
  	Slug      string
  	IsFound   bool
  	IsGuessed bool
  	Err       error
  }
  ```

### 2.2 Refactor `resolveDirectorySlug`
- **File:** `cli/cmdpipeline/pipeline_target_resolve.go`
- **Signature:** `func resolveDirectorySlug(dirPath string) DirectorySlugResult`
- **Implementation:**
  ```go
  func resolveDirectorySlug(dirPath string) DirectorySlugResult {
  	origin := gitutil.GetRemoteOriginURL(dirPath)
  	if origin != "" {
  		return DirectorySlugResult{
  			Slug:    parseSlugFromGitURL(origin),
  			IsFound: true,
  		}
  	}

  	base := filepath.Base(dirPath)
  	if slug, isFound := queryRepoSlugFromDB(base); isFound {
  		return DirectorySlugResult{
  			Slug:    slug,
  			IsFound: true,
  		}
  	}

  	return DirectorySlugResult{
  		Slug:      constants.DefaultGitHubOwner + "/" + base,
  		IsFound:   false,
  		IsGuessed: true,
  	}
  }
  ```

### 2.3 Update Caller `ResolvePipelineTargetAndPath`
- **File:** `cli/cmdpipeline/pipeline_target_resolve.go`
- **Implementation:**
  ```go
  if info, err := os.Stat(trimmed); err == nil && info.IsDir() {
  	abs, _ := filepath.Abs(trimmed)
  	res := resolveDirectorySlug(trimmed)

  	return res.Slug, abs
  }
  ```

---

## 3. Strict Coding Guidelines Compliance Checklist

- [ ] **Error Sovereignty:** Result wrapper encapsulates resolution outcome and provenance.
- [ ] **15-Line Limit:** Helper function remains <= 15 lines.
- [ ] **Zero Nested Ifs:** Flattened with guard returns.
- [ ] **Positive Booleans:** Fields named `IsFound` and `IsGuessed` (zero negative polarity).
- [ ] **Return Styling:** Standard vertical whitespace before returns.
- [ ] **Zero Git / Zero Build / Zero Test:** Strictly spec and plan authoring.
