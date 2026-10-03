# Rust Memory Safety (AI Execution Prompt)

> **/goal** Enforce strict Rust ownership idioms, borrow semantics, lifetime hygiene, and mandatory safety justification comments for all unsafe blocks.
> **/learn** Master memory safety principles: preferring borrowing over cloning, using `Arc` for cross-task shared state, `&str` over `&String`, lifetime elision idioms, and zero unannotated `unsafe` code.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Prefer borrowing (`&T`, `&str`, `&[T]`) over cloning (`.clone()`) and taking owned types unnecessarily.
- [ ] `/learn` Utilize `Arc<T>` for thread-safe shared ownership across concurrent async tasks without unnecessary deep copies.
- [ ] `/goal` Document every `unsafe` block with a clear, mandatory `// SAFETY:` rationale explaining pointer validity and invariant preservation.
- [ ] `/learn` Provide collection capacity hints via `Vec::with_capacity()` when allocation size is known ahead of time.
- [ ] `/learn` Verify zero guideline violations via `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Overview

Ownership idioms, lifetime guidelines, and strict `unsafe` policy for Rust projects.

---

## Ownership Rules

### Rule 1: Prefer borrowing over cloning

```rust
// ✅ REQUIRED — borrow when ownership isn't needed
fn classify_url(url: &str, categories: &[UrlCategory]) -> CategoryMatch {
    categories.iter().find(|category| category.matches(url))
}

// ❌ FORBIDDEN — unnecessary clone
fn classify_url(url: String, categories: Vec<UrlCategory>) -> CategoryMatch {
    categories.iter().find(|category| category.matches(&url))
}
```

### Rule 2: Use `Arc` for shared ownership across tasks

```rust
// ✅ REQUIRED — config shared across multiple collector tasks
let config = Arc::new(config);

for collector in collectors {
    let config = Arc::clone(&config);
    tokio::spawn(async move {
        collector.run(config).await
    });
}
```

### Rule 3: Prefer `&str` over `&String` in function parameters

```rust
// ✅ REQUIRED — accepts both &String and &str
fn parse_browser_title(title: &str) -> Option<TabInfo> { ... }

// ❌ FORBIDDEN — unnecessarily restrictive
fn parse_browser_title(title: &String) -> Option<TabInfo> { ... }
```

---

## Lifetime Guidelines

### Keep lifetimes simple — avoid naming when elision works

```rust
// ✅ REQUIRED — elision handles this
fn get_name(&self) -> &str { &self.name }

// ❌ FORBIDDEN — explicit lifetime adds noise
fn get_name<'a>(&'a self) -> &'a str { &self.name }
```

### Name lifetimes descriptively when multiple are needed

```rust
// ✅ REQUIRED — clear what each lifetime represents
fn merge_activities<'session, 'filter>(
    session: &'session Session,
    filter: &'filter ActivityFilter,
) -> Vec<&'session AppActivity> { ... }
```

---

## `unsafe` Policy

### Strict Rule: `unsafe` requires justification comment and review

Every `unsafe` block **must** include:

1. A `// SAFETY:` comment explaining why it's sound
2. The invariants being upheld

```rust
// ✅ REQUIRED — justified FFI call
// SAFETY: GetForegroundWindow returns a valid HWND or null.
// Null is checked immediately after the call.
let handle = unsafe { GetForegroundWindow() };
if handle.is_invalid() {
    return Err(OsError::WindowInfoFailed);
}

// ❌ FORBIDDEN — no safety comment
let handle = unsafe { GetForegroundWindow() };
```

### Where `unsafe` is permitted

| Context | Allowed | Notes |
|---------|---------|-------|
| FFI calls (Win32, X11, macOS) | ✅ Yes | Must have `// SAFETY:` comment |
| Platform-specific hooks | ✅ Yes | Wrapped in safe abstraction |
| Performance-critical hot paths | ⚠️ Rarely | Must prove safe alternative is too slow |
| Convenience / avoiding borrow checker | ❌ Never | Restructure the code instead |

### Wrap `unsafe` in safe abstractions

```rust
// ✅ REQUIRED — unsafe FFI wrapped in safe public API
pub fn get_active_window() -> Result<WindowInfo, OsError> {
    // SAFETY: GetForegroundWindow is safe to call and returns
    // HWND(0) when no window has focus, which we handle below.
    let handle = unsafe { GetForegroundWindow() };
    if handle.0 == 0 {
        return Err(OsError::NoFocusedWindow);
    }
    // ... safe code to extract window info
    Ok(info)
}

// ❌ FORBIDDEN — exposing unsafe to callers
pub unsafe fn get_active_window_raw() -> HWND {
    GetForegroundWindow()
}
```

---

## String Handling

```rust
// ✅ Use String for owned data, &str for borrowed
pub struct AppActivity {
    pub app_name: String,      // Owned — stored in struct
    pub window_title: String,  // Owned — stored in struct
}

// ✅ Accept &str in functions that don't need ownership
pub fn matches_pattern(text: &str, pattern: &str) -> bool { ... }

// ✅ Use Cow<str> when sometimes owned, sometimes borrowed
use std::borrow::Cow;
pub fn normalize_app_name(name: &str) -> Cow<'_, str> {
    if name.contains(".exe") {
        Cow::Owned(name.replace(".exe", ""))
    } else {
        Cow::Borrowed(name)
    }
}
```

---

## Collection Guidelines

| Need | Type | Example |
|------|------|---------|
| Ordered, growable | `Vec<T>` | Event buffer |
| Key-value lookup | `HashMap<K, V>` | URL category cache |
| Ordered key-value | `BTreeMap<K, V>` | Sorted reports |
| Unique set | `HashSet<T>` | Excluded domains |
| Fixed size | `[T; N]` | Heatmap grid cells |
| Optional value | `Option<T>` | Nullable URL field |

### Capacity hints

```rust
// ✅ REQUIRED — pre-allocate when size is known
let mut events = Vec::with_capacity(batch_size);

// ❌ FORBIDDEN — repeated reallocations
let mut events = Vec::new();
for _ in 0..1000 {
    events.push(event); // May reallocate multiple times
}
```

---

## Cross-References

| Reference | Location |
|-----------|----------|
| FFI & Platform Abstraction | `./07-ffi-platform.md` |
| Cross-Language Guidelines | `../01-cross-language/readme.md` |

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/05-rust/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-RUST-005: Rust Zero Unsafe and RAII Lifetime Safety

**Given** Rust codebases handling memory allocations, lifetimes, references, and unsafe blocks.
**When** Audited against memory safety standards and lifetime hygiene rules.
**Then** Borrowing is preferred over cloning, every `unsafe` block includes a mandatory `// SAFETY:` comment justifying pointer validity and sound invariants, all unsafe calls are encapsulated in safe RAII wrappers, and zero violations are detected with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only
```
**Expected:** exit 0. Zero violations.
