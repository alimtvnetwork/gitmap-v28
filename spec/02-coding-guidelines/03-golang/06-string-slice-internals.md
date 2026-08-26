# Go String & Slice Internals

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
