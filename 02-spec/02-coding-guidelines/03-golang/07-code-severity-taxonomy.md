# Code Severity Taxonomy (AI Execution Prompt)

> **/goal** Establish a strict, deterministic code severity classification system (Code Red 🔴 and Dangerous ⚠️) to prioritize immediate review remediations, eradicate silent panics, prevent race conditions, and guarantee structured fault logging.
> **/learn** Master Code Red defect triggers (unchecked return values, nil dereferences, multi-defer ambiguity, unlocked shared state, ungrounded path logging), understand Dangerous risk vectors (unbounded Big-O, swallowed errors, loop-nested regex, improper mutex scoping), and apply structured review tags.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Halt execution immediately on any Code Red violation: unchecked returns, nil pointer dereferences, or unlocked concurrent state mutations.
- [ ] `/learn` Never swallow errors silently or log file system failures without exact paths, operation context, and structured fault metadata.
- [ ] `/goal` Enforce proper mutex scoping as struct fields rather than local function variables to prevent race conditions.
- [ ] `/learn` Hoist expensive operations (regex compilation, reflection, eager loading) out of hot paths and loops to eliminate Dangerous Big-O bottlenecks.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

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
- File/path error logged without exact file path or failure reason ([rule](../../03-error-manage/01-error-resolution/app-issues/02-error-management-file-path-and-missing-file-code-red-rule.md))

```go
// ❌ FORBIDDEN: Unchecked return value causing potential nil panic (Code Red)
user, _ := findUser(userId)
fmt.Println(user.Name)

// ✅ REQUIRED: Strict error check before accessing return value
user, appErr := findUser(userId)
if appErr != nil {
    return appfault.Wrap(appErr, "findUser")
}

fmt.Println(user.Name)
```

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

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-GO-007: Code Severity Taxonomy, Error Classification, and Fault Logging

**Given** Go source code across packages and CLI modules.
**When** Guideline linters audit the codebase.
**Then** Code review issues are classified deterministically as Code Red or Dangerous, zero unchecked errors or ungrounded paths exist, and all faults log complete contextual metadata.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only
```
**Expected:** exit 0. Zero violations detected.
