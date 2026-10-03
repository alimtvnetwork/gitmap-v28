# Rust FFI & Platform Abstraction (AI Execution Prompt)

> **/goal** Architect robust, cross-platform Rust abstractions, conditional compilation architectures, and strictly encapsulated FFI boundaries across Windows, Linux, and macOS.
> **/learn** Master platform isolation principles: trait-based `PlatformApi` contracts, OS-specific module segregation via `#[cfg(target_os)]`, safe RAII wrappers around native OS handles, and vetting platform binding crates.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Encapsulate all platform-specific logic behind unified, platform-agnostic traits (`PlatformApi`).
- [ ] `/learn` Organize OS-specific code into discrete submodules (`windows.rs`, `linux.rs`, `macos.rs`) gated by `#[cfg(target_os)]`.
- [ ] `/goal` Wrap raw FFI pointers, window handles, and system hooks in safe RAII types that implement `Drop`.
- [ ] `/learn` Prioritize well-maintained, official binding crates (e.g., Microsoft `windows`, `x11rb`, `objc2`) over unmaintained bindings.
- [ ] `/learn` Verify zero guideline violations via `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Overview

Guidelines for Foreign Function Interface (FFI) calls and cross-platform code organization. The Time Log CLI uses conditional compilation to support Windows, Linux, and macOS from a single codebase.

---

## Platform Abstraction Pattern

### Public Interface (Platform-Agnostic)

```rust
// src/platform/mod.rs

#[cfg(target_os = "windows")]
mod windows;
#[cfg(target_os = "linux")]
mod linux;
#[cfg(target_os = "macos")]
mod macos;

// Re-export platform-specific implementations under a unified API
#[cfg(target_os = "windows")]
pub use windows::*;
#[cfg(target_os = "linux")]
pub use linux::*;
#[cfg(target_os = "macos")]
pub use macos::*;

/// Platform-agnostic window information
#[derive(Debug, Clone)]
pub struct WindowInfo {
    pub app_name: String,
    pub window_title: String,
    pub process_id: u32,
}

/// Every platform module must implement these functions
pub trait PlatformApi: Send + Sync {
    fn get_active_window(&self) -> Result<WindowInfo, OsError>;
    fn get_idle_seconds(&self) -> Result<u64, OsError>;
    fn capture_screen(&self, mode: CaptureMode) -> Result<RawBuffer, CaptureError>;
    fn register_click_hook(&self, callback: ClickCallback) -> Result<HookHandle, OsError>;
    fn register_autostart(&self) -> Result<(), OsError>;
    fn unregister_autostart(&self) -> Result<(), OsError>;
}
```

### Platform Implementation Example

```rust
// src/platform/windows.rs
#[cfg(target_os = "windows")]

use windows::Win32::UI::WindowsAndMessaging::*;

pub struct WindowsPlatform;

impl PlatformApi for WindowsPlatform {
    fn get_active_window(&self) -> Result<WindowInfo, OsError> {
        // SAFETY: GetForegroundWindow is always safe to call.
        // Returns HWND(0) when no window has focus.
        let handle = unsafe { GetForegroundWindow() };
        if handle.0 == 0 {
            return Err(OsError::NoFocusedWindow);
        }

        let title = self.get_window_title(handle)?;
        let process_id = self.get_window_process_id(handle)?;
        let app_name = self.get_process_name(process_id)?;

        Ok(WindowInfo { app_name, window_title: title, process_id })
    }

    // ... other trait methods
}
```

---

## Conditional Compilation Rules

### Rule 1: Use `cfg` at module level, not scattered inline

```rust
// ✅ REQUIRED — entire module is platform-specific
#[cfg(target_os = "windows")]
mod windows_hooks {
    pub fn register_mouse_hook() -> Result<HookHandle, OsError> { ... }
}

// ❌ FORBIDDEN — scattered cfg attributes throughout a function
pub fn register_hook() -> Result<HookHandle, OsError> {
    #[cfg(target_os = "windows")]
    { return windows_register(); }
    #[cfg(target_os = "linux")]
    { return linux_register(); }
    #[cfg(target_os = "macos")]
    { return macos_register(); }
}
```

### Rule 2: Feature flags for optional capabilities

```toml

# Cargo.toml

[features]
default = ["x11"]
x11 = ["x11rb"]
wayland = ["wayland-client", "pipewire"]
browser-extension = []  # Enable native messaging host
```

```rust
// Feature-gated code
#[cfg(feature = "wayland")]
mod wayland_capture {
    pub fn capture_via_pipewire() -> Result<RawBuffer, CaptureError> { ... }
}
```

### Rule 3: Compile-time platform validation

```rust
// Ensure at least one display server is available on Linux
#[cfg(all(target_os = "linux", not(any(feature = "x11", feature = "wayland"))))]
compile_error!("Linux builds require either 'x11' or 'wayland' feature");
```

---

## FFI Safety Checklist

For every `unsafe` FFI call, verify:

| Check | Description |
|-------|-------------|
| ✅ Null check | All returned pointers are checked for null before dereference |
| ✅ Lifetime | Borrowed data outlives the FFI call |
| ✅ Thread safety | FFI call is safe to invoke from the calling thread |
| ✅ Error handling | OS error codes are checked and converted to Rust `Result` |
| ✅ Resource cleanup | All acquired resources (handles, hooks) are released on drop |
| ✅ SAFETY comment | `// SAFETY:` comment documents invariants |

### RAII for OS Resources

```rust
/// Hook handle that unregisters on drop
pub struct HookHandle {
    #[cfg(target_os = "windows")]
    handle: HHOOK,
}

impl Drop for HookHandle {
    fn drop(&mut self) {
        #[cfg(target_os = "windows")]
        {
            // SAFETY: UnhookWindowsHookEx is safe to call with a valid HHOOK.
            // The handle was obtained from SetWindowsHookExW and has not been
            // previously unhooked.
            unsafe { UnhookWindowsHookEx(self.handle) };
        }
    }
}
```

---

## Crate Selection for Platform APIs

| Platform | Crate | Purpose |
|----------|-------|---------|
| Windows | `windows` (official Microsoft) | Win32 API bindings |
| Linux (X11) | `x11rb` | X11 protocol, pure Rust |
| Linux (Wayland) | `wayland-client` | Wayland protocol |
| macOS | `core-foundation`, `core-graphics` | System framework bindings |
| macOS | `objc2` | Objective-C runtime interop |
| Cross-platform | `raw-window-handle` | Window handle abstraction |

### Crate Evaluation Criteria

1. **Official or well-maintained** — Prefer Microsoft's `windows` crate over `winapi`
2. **Safe wrappers available** — Prefer crates that provide safe Rust APIs
3. **Minimal dependencies** — Avoid pulling in unnecessary transitive deps
4. **Active maintenance** — Last release within 6 months

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Memory Safety (unsafe policy) | `./05-memory-safety.md` |
| Cross-Language Guidelines | `../01-cross-language/readme.md` |

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/05-rust/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-RUST-007: Rust FFI Safety Boundaries and Platform Isolation

**Given** Rust codebases integrating with external C libraries, Win32 APIs, or Unix system calls.
**When** Audited against FFI safety guidelines and platform isolation architecture.
**Then** All external C interfaces are isolated behind safe `PlatformApi` trait boundaries, every OS handle implements deterministic RAII drop cleanup, conditional compilation is confined to module boundaries, and zero violations are detected with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only
```
**Expected:** exit 0. Zero violations.
