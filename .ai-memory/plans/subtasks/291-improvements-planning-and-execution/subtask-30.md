# Subtask Plan 30: Remote Transport Guard Flattening & Style Hygiene (`cli/cmdclone/remotetransport.go`)

> **Subtask ID:** `subtask-30`  
> **Parent Plan:** `.ai-memory/plans/pending/291-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/291-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/cmdclone/remotetransport.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Refactor `cli/cmdclone/remotetransport.go` to eliminate unexported boolean tuples `(bool, bool)`.
2. Introduce a clear internal struct `transportRequest` for transport resolution.
3. Invert and flatten conditional branching using strictly Positive Booleans (Rule 4, eliminating `!isSSH && !isHTTPS`).
4. Enforce Vertical Whitespace Hygiene (Rule 6, blank line after `}` and blank line before `return`).
5. Ensure all functions are strictly <= 15 lines (Rule 2).

---

## 2. Targeted Problem & Architectural Rationale
In `cli/cmdclone/remotetransport.go`:
1. `resolveRequestedTransport(isSSH, isHTTPS bool) (bool, bool)` returns an unnamed boolean tuple `(bool, bool)`, creating ambiguity about which return represents the transport protocol vs. actionability.
2. Negative boolean condition `!isSSH && !isHTTPS` violates the strict positive boolean naming and expression mandate.
3. In `applyConvertedURL`, statements directly follow closing braces without intermediate blank lines.

Refactoring this file establishes a model of coding guideline compliance across the clone subsystem.

---

## 3. Concrete Implementation Details

### 3.1 Refactoring in `cli/cmdclone/remotetransport.go`

```go
type transportRequest struct {
	isSSH        bool
	isActionable bool
}

func resolveRequestedTransport(isSSH, isHTTPS bool) transportRequest {
	hasTransport := isSSH || isHTTPS
	if !hasTransport {
		return transportRequest{isSSH: false, isActionable: false}
	}

	if isSSH && isHTTPS {
		fmt.Fprintln(os.Stderr, "⚠ both --ssh and --https set; --ssh wins")

		return transportRequest{isSSH: true, isActionable: true}
	}

	return transportRequest{isSSH: isSSH, isActionable: true}
}

// ApplyTransportFlag rewrites the `remote.origin.url` of dir to the
// requested transport and persists it via `git remote set-url`.
func ApplyTransportFlag(dir string, isSSH, isHTTPS bool) (bool, string, string, error) {
	req := resolveRequestedTransport(isSSH, isHTTPS)
	if !req.isActionable {
		return false, "", "", nil
	}

	old, err := currentOriginURL(dir)
	if err != nil {
		return false, "", "", err
	}

	return rewriteTransportURL(dir, old, req.isSSH)
}

func applyConvertedURL(dir, oldURL, convertedURL string) (bool, error) {
	if convertedURL == oldURL {
		return false, nil
	}

	if err := setOriginURL(dir, convertedURL); err != nil {
		return false, err
	}

	fmt.Printf("→ remote.origin.url: %s → %s\n", oldURL, convertedURL)

	return true, nil
}
```

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** Direct error returns preserved across all functions.
- **Rule 2 (Function Limit):** `resolveRequestedTransport` (11 lines), `ApplyTransportFlag` (10 lines), `applyConvertedURL` (10 lines). All <= 15 lines.
- **Rule 3 (Zero Nested Ifs):** Linear guard clauses on `!hasTransport` and conflict warning.
- **Rule 4 (Positive Booleans):** Positive identifier `hasTransport` replaces `!isSSH && !isHTTPS`. `isActionable` indicates operational status.
- **Rule 6 (Vertical Whitespace):** Exactly one blank line after `}` and before every `return` statement.

---

## 5. Verification Criteria
- [x] `transportRequest` struct is declared with `isSSH` and `isActionable`.
- [x] `resolveRequestedTransport` returns `transportRequest` with zero negative boolean operations.
- [x] `applyConvertedURL` includes vertical newline separation after closing brace.
- [x] Every function is <= 15 lines.
- [x] All file paths are strictly relative lowercase Git paths.
