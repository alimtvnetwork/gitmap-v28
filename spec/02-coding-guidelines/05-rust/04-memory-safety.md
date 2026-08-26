# Rust Memory Safety

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

Ownership idioms, lifetime guidelines, and strict `unsafe` policy for Rust projects.

---

## Ownership Rules

### Rule 1: Prefer borrowing over cloning

```rust
// ✅ Correct — borrow when ownership isn't needed
fn classify_url(url: &str, categories: &[UrlCategory]) -> CategoryMatch {
    categories.iter().find(|category| category.matches(url))
}

// ❌ Avoid — unnecessary clone
fn classify_url(url: String, categories: Vec<UrlCategory>) -> CategoryMatch {
    categories.iter().find(|category| category.matches(&url))
}
```

### Rule 2: Use `Arc` for shared ownership across tasks

```rust
// ✅ Correct — config shared across multiple collector tasks
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
// ✅ Correct — accepts both &String and &str
fn parse_browser_title(title: &str) -> Option<TabInfo> { ... }

// ❌ Avoid — unnecessarily restrictive
fn parse_browser_title(title: &String) -> Option<TabInfo> { ... }
```

---

## Lifetime Guidelines

### Keep lifetimes simple — avoid naming when elision works

```rust
// ✅ Correct — elision handles this
fn get_name(&self) -> &str { &self.name }

// ❌ Unnecessary — explicit lifetime adds noise
fn get_name<'a>(&'a self) -> &'a str { &self.name }
```

### Name lifetimes descriptively when multiple are needed

```rust
// ✅ Clear what each lifetime represents
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
// ✅ Correct — justified FFI call
// SAFETY: GetForegroundWindow returns a valid HWND or null.
// Null is checked immediately after the call.
let handle = unsafe { GetForegroundWindow() };
if handle.is_invalid() {
    return Err(OsError::WindowInfoFailed);
}

// ❌ Forbidden — no safety comment
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
// ✅ Correct — unsafe FFI wrapped in safe public API
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

// ❌ Forbidden — exposing unsafe to callers
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
// ✅ Correct — pre-allocate when size is known
let mut events = Vec::with_capacity(batch_size);

// ❌ Avoid — repeated reallocations
let mut events = Vec::new();
for _ in 0..1000 {
    events.push(event); // May reallocate multiple times
}
```

---

## Cross-References

| Reference | Location |
|-----------|----------|
| FFI & Platform Abstraction | `./06-ffi-platform.md` |
| Cross-Language Guidelines | `../01-cross-language/00-overview.md` |
