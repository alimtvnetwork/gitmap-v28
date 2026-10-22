# Subtask 16: Zero-Allocation Terminal Padding Helper (PadRightMax)

> **Subtask ID:** Subtask-16  
> **Parent Plan:** `.ai-memory/plans/289-improvements-planning-and-execution.md`  
> **Component Spec:** `02-spec/21-app/289-improvements-planning-and-execution/02-component-spec.md`  
> **Target Path:** `cli/termout/termpad.go`  
> **Status:** PENDING  

---

## 1. Objectives & Scope
1. Implement high-throughput padding helper `PadRightMax` in `cli/termout/termpad.go`.
2. Eliminate heap allocations caused by repeated `strings.Repeat(" ", padLen)` calls during table, help menu, and tree rendering.
3. Leverage a static pre-allocated space pool up to 256 characters.
4. Provide fallback to standard string repetition for oversized target widths.

---

## 2. Detailed Technical Plan

### 2.1 Pre-allocated Space Buffer Pool
- Define static space buffer in `cli/termout/termpad.go`:
  ```go
  const spacePool = "                                                                                                                                                                                                                                                                " // 256 spaces
  ```

### 2.2 Helper Implementation: `PadRightMax`
- Implement zero-allocation padding function:
  ```go
  // PadRightMax pads s on the right up to targetWidth using a static space pool,
  // avoiding dynamic heap allocations for typical terminal column widths.
  func PadRightMax(s string, targetWidth int) string {
  	curLen := len(s)
  	if curLen >= targetWidth {
  		return s
  	}

  	needed := targetWidth - curLen
  	if needed <= len(spacePool) {
  		return s + spacePool[:needed]
  	}

  	return s + strings.Repeat(" ", needed)
  }
  ```

### 2.3 Additional Padding Utilities
- Implement `SpacePool(n int) string` helper to allow zero-allocation space slice retrieval across `table.go` and `render.go`:
  ```go
  // SpacePool returns a string of n spaces from the static pool without heap allocation.
  func SpacePool(n int) string {
  	if n <= 0 {
  		return ""
  	}
  	if n <= len(spacePool) {
  		return spacePool[:n]
  	}

  	return strings.Repeat(" ", n)
  }
  ```

---

## 3. Strict Coding Guidelines Compliance
- **Function Limit:** Functions are 11 lines and 8 lines respectively (both <= 15 lines).
- **Guard Clauses:** Early returns on `curLen >= targetWidth` and `n <= 0`.
- **Zero Nested Ifs:** Linear execution flow with early exit.
- **Return Formatting:** Single blank line before return statements and after `}` blocks.

---

## 4. Verification Criteria
- [ ] `PadRightMax("abc", 5)` produces `"abc  "` without heap allocations.
- [ ] `PadRightMax("longer", 4)` returns `"longer"` unmodified.
- [ ] Requesting width beyond 256 gracefully falls back to `strings.Repeat`.
- [ ] Benchmarks demonstrate zero allocations for `n <= 256`.
