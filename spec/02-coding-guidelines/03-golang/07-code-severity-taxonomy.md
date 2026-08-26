# Code Severity Taxonomy

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**Source:** Consolidated from `01-pre-code-review-guides/03-golang-code-review-guides.md`

---

## 1. Purpose

Classification system for code issues found during review. Helps prioritize fixes and communicate severity.

---

## 2. Code Red 🔴

Issues that are **critical** and must be fixed immediately:

- [ ] Code has serious issues, tends to be buggy in future
- [ ] Requires lots of investigation in future
- [ ] Fluctuating / non-deterministic results
- [ ] May have locking-related issues (race conditions)
- [ ] Looks alright but contains hard-to-detect bugs

**Examples:**
- Calling methods on unchecked return values
- Missing nil checks on pointers
- Multiple defers creating unclear execution order
- Mutation of shared state without locks
- File/path error logged without exact file path or failure reason ([rule](../../03-error-manage/01-error-resolution/app-issues/error-management-file-path-and-missing-file-code-red-rule.md))

---

## 3. Dangerous ⚠️

Issues that **complement Code Red** and pose significant risk:

- [ ] Code could throw or panic at uncertain times
- [ ] Code may cause UX or UI issues
- [ ] May not trace/log properly — swallows errors silently
- [ ] Not optimized in BigO — could be easily improved
- [ ] Lazy evaluation not done properly (eager when should be lazy)
- [ ] Code may have unknown loopholes
- [ ] Mutex not used properly (e.g., mutex created inside function instead of as struct field)

---

## 4. Usage in Reviews

When flagging issues in code review, use these labels:

```
🔴 [Code Red] — Pointer dereference without nil check on line 42
⚠️ [Dangerous] — Regex compiled inside loop on line 78
```

---

## 5. Cross-References

- [Null Pointer Safety](../01-cross-language/19-null-pointer-safety.md) — Common Code Red issues
- [Lazy Evaluation](../01-cross-language/16-lazy-evaluation-patterns.md) — Dangerous: improper lazy
- [Code Mutation Avoidance](../01-cross-language/18-code-mutation-avoidance.md) — Code Red: unlocked mutation

---

*Code severity taxonomy — consolidated from pre-code review guides.*
