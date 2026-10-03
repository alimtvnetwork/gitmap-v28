# Go String & Slice Internals (AI Execution Prompt)

> **/goal** Master Go string and slice memory representation, eliminate redundant pointers (`*string`, `*[]T`), prevent unintentional mutations across sub-slices, and enforce preallocation and zero-allocation idioms.
> **/learn** Internal headers (`StringHeader` 16 bytes, `SliceHeader` 24 bytes), value-passing semantics with underlying reference sharing, re-allocation mechanics on `append()`, sub-slice mutation hazards, and compiler copy optimizations.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Pass strings and slices by value for read-only operations; never accept `*string` or `*[]T` unless actively mutating the slice header.
- [ ] `/learn` Understand that sub-slicing (`s[1:3]`) references the original backing array; copy explicitly when sub-slices must not retain large allocations or mutate shared memory.
- [ ] `/goal` Preallocate slice capacity via `make([]T, 0, cap)` when the element count is known to avoid incremental heap reallocations.
- [ ] `/learn` Verify that strings are treated as immutable values and string concatenations in tight loops use `strings.Builder`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 3.2.0
**Updated:** 2026-04-16
**Source:** Consolidated from `01-pre-code-review-guides/03-golang-code-review-guides.md`

---

## 1. Purpose

Understanding Go's internal representation of strings and slices helps avoid common performance and correctness mistakes.

---

## 2. String Internals

Strings in Go are **passed by value**, but the underlying data is **passed by reference**. The string header struct gets copied, but the byte data does not.

```go
// Go's internal string representation
type StringHeader struct {
    Data uintptr  // pointer to byte data (shared, not copied)
    Len  int      // length of string
}
```

**Implications:**

- Passing a string to a function copies the header (16 bytes) but **not** the data
- No need to pass `*string` for read-only use — it's already efficient
- Strings are immutable — any modification creates a new allocation

```go
// ❌ FORBIDDEN: Passing *string or *[]T for read-only operations
func renderSummary(content *string, items *[]string) string {
    return fmt.Sprintf("%s: %d", *content, len(*items))
}

// ✅ REQUIRED: Value semantics for read-only strings and slices
func renderSummary(content string, items []string) string {
    return fmt.Sprintf("%s: %d", content, len(items))
}
```

**References:**

- [What is the point of passing a pointer to strings in Go?](https://stackoverflow.com/questions/24642311)
- [Is string passed by value or reference?](https://groups.google.com/g/golang-nuts/c/ZRKSJ3GPkLw)

---

## 3. Slice Internals

```go
// Go's internal slice representation
type SliceHeader struct {
    Data uintptr  // pointer to array data
    Len  int      // current length
    Cap  int      // capacity
}
```

**Implications:**

- Passing a slice copies the header (24 bytes) but shares the underlying array
- `append()` may create a new array if capacity is exceeded
- Slicing (`s[1:3]`) shares the same underlying array — mutations affect both

**References:**

- [Go Slice Tricks](https://github.com/golang/go/wiki/SliceTricks)
- [Go Slice Tricks Cheat Sheet](https://ueokande.github.io/go-slice-tricks/)

---

## 4. Cross-References

- [Null Pointer Safety](../01-cross-language/19-null-pointer-safety.md) — Nil checks for slices
- [Code Mutation Avoidance](../01-cross-language/18-code-mutation-avoidance.md) — Immutability principles

---

*Go string & slice internals — consolidated from pre-code review guides.*

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-GO-006: Go String and Slice Preallocation & Memory Efficiency Standards

**Given** Go source code across packages and CLI modules.
**When** Guideline linters audit the codebase.
**Then** String and slice parameters adhere to value semantics without redundant pointers, slice allocations are preallocated where capacities are known, and zero memory corruption occurs.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only
```
**Expected:** exit 0. Zero violations detected.
