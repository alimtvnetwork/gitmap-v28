# Prompt: Antigravity-Manager Linux Headless & Keyring Account Switching Remediation

## Objective
Refactor `Antigravity-Manager` (`src-tauri/src/modules/integration.rs`) to support reliable, non-blocking account switching on Linux across desktop, headless VM, SSH, and containerized environments.

---

## 1. Root Cause Analysis (RCA)

### Issue 1: Hard Failure on Keyring Timeout Despite File-Based Credential Success
- **Location:** `src-tauri/src/modules/integration.rs` lines 616–710 (`write_to_system_keyring`) and lines 145–180 (`apply_account_credentials`).
- **Mechanism:** On Linux, `write_to_system_keyring` attempts to store credentials into the freedesktop Secret Service using `secret-tool store --collection=login` and `secret-tool store (default)`.
- **Failure:** In headless VM or SSH sessions without an unlocked desktop keyring daemon, `secret-tool` blocks until reaching the 10-second timeout.
- **Bug:** `apply_account_credentials` already successfully writes file-based credentials to `~/.gemini/oauth_creds.json` and `google_accounts.json`. However, line 706 of `integration.rs` executes:
  ```rust
  if login_res.is_err() && default_res.is_err() {
      return Err(login_res.unwrap_err());
  }
  ```
  This causes the entire account switch function to return an `Err("Keyring write timed out (10s)")`, discarding the successful token synchronization.

### Issue 2: Self-Recursive Launcher Script on Linux (`Argument list too long`)
- **Location:** Linux launcher generation and deployment scripts.
- **Mechanism:** When the launcher script contains:
  ```bash
  DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  EXEC="$DIR/antigravity"
  exec "$EXEC" --no-sandbox "$@"
  ```
  If this script is copied or symlinked to `/usr/local/bin/antigravity` or `~/.local/bin/antigravity`, `$DIR/antigravity` resolves back to the script itself.
- **Failure:** Spawning `antigravity` executes an infinite recursive loop until Linux aborts with `E2BIG: Argument list too long`.

### Issue 3: Unhandled Headless Environment During IDE Relaunch
- **Mechanism:** In Step 4/5 of account switching, AGM executes the IDE executable without checking whether an active display server (`DISPLAY` or `WAYLAND_DISPLAY`) is available.
- **Failure:** Electron aborts with `Missing X server or $DISPLAY`, producing an unnecessary crash log in non-GUI sessions.

---

## 2. Required Surgical Code Changes

### Change A: Graceful Headless Fallback in `integration.rs`
In `apply_account_credentials`:
```rust
// Attempt system keyring write, but do NOT abort if file-based sync succeeded
let keyring_res = write_to_system_keyring("antigravity", &email, &token);
if let Err(e) = keyring_res {
    log::warn!(
        "[Desktop] Keyring write failed ({}), but file-based credentials in ~/.gemini were successfully synced. Continuing gracefully.",
        e
    );
}
```

In `write_to_system_keyring`:
1. Check if `DBUS_SESSION_BUS_ADDRESS` or an active Secret Service is reachable before spawning `secret-tool`.
2. If `secret-tool` is missing or fails due to timeout/D-Bus unreachable, return a specialized `Err("Keyring unavailable in headless session")` that is treated as non-fatal when file sync succeeds.

### Change B: Robust Path Resolution in Launcher Scripts
Ensure launcher scripts distinguish between the launcher entrypoint in `bin/` and the application bundle in `share/`:
```bash
if [ -x "${HOME}/.local/share/antigravity-ide/antigravity" ]; then
    DIR="${HOME}/.local/share/antigravity-ide"
elif [ -d "/opt/antigravity" ]; then
    DIR="/opt/antigravity"
else
    DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
fi
EXEC="${DIR}/antigravity"
```

### Change C: Headless Display Guard Before Spawning GUI IDE
In Step 4/5 (`relaunch_antigravity_ide`):
```rust
let has_display = std::env::var("DISPLAY").is_ok() || std::env::var("WAYLAND_DISPLAY").is_ok();
if !has_display {
    log::info!("[Desktop] Headless session detected (no DISPLAY or WAYLAND_DISPLAY); skipping GUI IDE launch.");
    return Ok(());
}
```

---

## 3. Verification Criteria
1. `agm switch <account>` on a headless Linux SSH session exits with status 0.
2. `~/.gemini/oauth_creds.json` is updated with active tokens.
3. No infinite recursive loops or `Argument list too long` errors occur when launching the binary.
