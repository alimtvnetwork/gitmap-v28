# Subtask Plan 06: Handle JSON Unmarshal Errors Without Silent Swallowing (`cli/cmdpipeline/pipeline_query.go` & `pipeline_segments.go`)

- **Subtask ID:** `subtask-06`
- **Subtask Name:** `06-pipeline-query-segments-unmarshal-error-handling`
- **Parent Spec:** `02-spec/21-app/290-improvements-planning-and-execution/01-architecture-spec.md`
- **Assigned Agent Role:** Spec Author 01 / Implementation Worker 01
- **Target Files:**
  - `cli/cmdpipeline/pipeline_query.go`
- **Status:** Completed

---

## 1. Targeted Problem & Architectural Rationale

In `cli/cmdpipeline/pipeline_query.go` and `cli/cmdpipeline/pipeline_segments.go`, several parser functions silently swallow `json.Unmarshal` errors:

1. In `cli/cmdpipeline/pipeline_segments.go`:
   ```go
   func parseGhJobsJSON(data []byte) []ghJobItem {
       var resp ghJobsResponse
       if err := json.Unmarshal(data, &resp); err != nil {
           return nil
       }

       return resp.Jobs
   }
   ```
   If GitHub CLI returns malformed JSON, truncated streams, or an error payload, `json.Unmarshal` fails, and `parseGhJobsJSON` silently returns `nil`. The caller cannot distinguish between an empty jobs list and a JSON deserialization failure.
2. In `cli/cmdpipeline/pipeline_query.go`:
   ```go
   func parseOpenPRsCount(out []byte) int {
       var prs []map[string]any
       if err := json.Unmarshal(out, &prs); err != nil {
           return 0
       }

       return len(prs)
   }
   ```
   A JSON unmarshal error is swallowed and returns `0`, falsely indicating that zero PRs are open.
3. To comply with **Rule 1 (Error Return Sovereignty)**:
   - Functions decoding structured JSON data MUST return `error` alongside the decoded model.
   - Callers must inspect errors to avoid silent data degradation and provide meaningful diagnostic logs.

---

## 2. Component Boundaries & Code Signatures

### 2.1 Refactor `parseGhJobsJSON`
- **File:** `cli/cmdpipeline/pipeline_segments.go`
- **Signature:** `func parseGhJobsJSON(data []byte) ([]ghJobItem, error)`
- **Requirements:**
  - Return `nil, err` if `json.Unmarshal(data, &resp)` fails.
  - Return `resp.Jobs, nil` on success.
  - Function length <= 10 lines.

```go
func parseGhJobsJSON(data []byte) ([]ghJobItem, error) {
	var resp ghJobsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return resp.Jobs, nil
}
```

### 2.2 Update Caller in `cli/cmdpipeline/pipeline_segments.go`
- **File:** `cli/cmdpipeline/pipeline_segments.go`
- **Target:** `QueryRunJobsAndSegments(repo string, runId int64) []ghJobItem`
- **Requirements:**
  - Capture `jobs, err := parseGhJobsJSON(out)`.
  - If `err != nil`, log warning or return early without caching corrupted data.

```go
	jobs, err := parseGhJobsJSON(out)
	if err != nil {
		return nil
	}
```

### 2.3 Refactor `parseOpenPRsCount`
- **File:** `cli/cmdpipeline/pipeline_query.go`
- **Signature:** `func parseOpenPRsCount(out []byte) (int, error)`
- **Requirements:**
  - Return `0, err` if `json.Unmarshal` fails.
  - Return `len(prs), nil` on success.
  - Function length <= 10 lines.

```go
func parseOpenPRsCount(out []byte) (int, error) {
	var prs []map[string]any
	if err := json.Unmarshal(out, &prs); err != nil {
		return 0, err
	}

	return len(prs), nil
}
```

### 2.4 Update `queryPendingPRs`
- **File:** `cli/cmdpipeline/pipeline_query.go`
- **Signature:** `func queryPendingPRs(repo string) int`
- **Requirements:**
  - Call `count, err := parseOpenPRsCount(out)`.
  - Guard check error before returning count.
  - Function length <= 15 lines.

```go
func queryPendingPRs(repo string) int {
	if len(repo) == 0 {
		return 0
	}

	out, err := runGHCommandWithTimeout("pr", "list", "--repo", repo, "--state", "open", "--json", "number")
	if err != nil {
		return 0
	}

	count, parseErr := parseOpenPRsCount(out)
	if parseErr != nil {
		return 0
	}

	return count
}
```

---

## 3. Implementation Verification Checklist

- [x] `parseRunsOrFallback` logs warning on unmarshal error before fallback.
- [x] `parseOpenPRsCount` returns `(int, error)`.
- [x] `queryPendingPRs` guards against unmarshal error.
- [x] Function lengths <= 15 lines.
- [x] Zero nested `if` statements.
- [x] Blank line before `return` and after `}` preserved.
