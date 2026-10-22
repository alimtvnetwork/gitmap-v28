# Subtask 28: Remote Transport Guard Clause & Boolean Enforcement

> **Subtask ID:** Subtask-28  
> **Parent Plan:** `.ai-memory/plans/pending/289-improvements-planning-and-execution.md`  
> **Component Spec:** `02-spec/21-app/289-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/cmdclone/remotetransport.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Audit and refactor `cli/cmdclone/remotetransport.go` to strictly enforce zero nested if statements using early guard returns.
2. Standardize all boolean parameters and local variables to use `is*` or `has*` positive prefixes exclusively (`isSSH`, `isHTTPS` replacing `useSSH`, `useHTTPS`).
3. Preserve existing transport rewrite behavior for `--ssh` and `--https` flags without breaking CLI compatibility.
4. Ensure all functions remain within the <= 15 line limit with strict vertical whitespace styling.

---

## 2. Detailed Technical Plan

### 2.1 Boolean Parameter Normalization
- Update `ApplyTransportFlag` signature to eliminate banned `use*` boolean prefixes:
  ```go
  func ApplyTransportFlag(dir string, isSSH, isHTTPS bool) (bool, string, string, error) {
  	isRequestedSSH, isActionable := resolveRequestedTransport(isSSH, isHTTPS)
  	if !isActionable {
  		return false, "", "", nil
  	}

  	old, err := currentOriginURL(dir)
  	if err != nil {
  		return false, "", "", err
  	}

  	return rewriteTransportURL(dir, old, isRequestedSSH)
  }
  ```

### 2.2 Guard Clause Verification in `resolveRequestedTransport`
- Verify flat guard structure:
  ```go
  func resolveRequestedTransport(isSSH, isHTTPS bool) (bool, bool) {
  	if !isSSH && !isHTTPS {
  		return false, false
  	}

  	if isSSH && isHTTPS {
  		fmt.Fprintln(os.Stderr, "⚠ both --ssh and --https set; --ssh wins")

  		return true, true
  	}

  	return isSSH, true
  }
  ```

### 2.3 Early Guard Inversion in `applyConvertedURL`
- Ensure error and equality checks return early:
  ```go
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

## 3. Strict Coding Guidelines Compliance
- **Rule 3 (Zero Nested Ifs):** Completely flat conditional structure; zero nested if statements.
- **Rule 4 (Boolean Naming):** Replaced `useSSH`/`useHTTPS` with positive `isSSH`/`isHTTPS`.
- **Rule 2 (Function Limit):** All functions are 6 to 13 lines (<= 15 line max).
- **Rule 6 (Vertical Whitespace):** Single blank line after each `}` closure and before return statements.

---

## 4. Verification Criteria
- [x] `remotetransport.go` contains zero nested if statements.
- [x] All booleans adhere to `is*` or `has*` positive prefixes.
- [x] Error return sovereignty preserved across all functions.
- [x] Zero compiler builds and zero tests executed during planning.
