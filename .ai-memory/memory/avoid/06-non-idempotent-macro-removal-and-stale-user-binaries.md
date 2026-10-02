# Avoid: Non-Idempotent Windows PowerShell Removal in Macro Steps & Stale User-Profile Binaries

**Status:** 🚫 Blocked — Cross-Platform Idempotent Removal & Multi-Path Binary Deployment Mandated
**Recorded:** 2026-10-01

---

## 1. Rule

🔴 **NEVER permit macro execution on Windows PowerShell to invoke raw, unshimmed removal commands (`rm`, `rmdir`, `del`, `Remove-Item`, `erase`) that abort terminatingly with `ItemNotFoundException` (`exit status 1`) when encountering absent targets, and NEVER deploy or update `gitmap.exe` to only a single system path while neglecting user-specific installation targets (`C:\Users\<user>\AppData\Local\gitmap\gitmap.exe` and `AppData\Local\gitmap-cli\gitmap.exe`).**

---

## 2. Forbidden Actions

- ❌ Allowing macro execution engines to dispatch raw `rm <target>` directly to Windows PowerShell without platform adaptation shims.
- ❌ Relying on PowerShell's default `rm` alias (`Remove-Item`) without idempotent existence guards (`Test-Path`), which triggers `ItemNotFoundException` whenever the target file or directory is already absent.
- ❌ Compiling or installing `gitmap.exe` exclusively to `%LOCALAPPDATA%\gitmap-cli\gitmap.exe` (`Administrator`) while leaving stale binaries in other active user profiles (`C:\Users\Alim\AppData\Local\gitmap\gitmap.exe` and `C:\Users\Alim\AppData\Local\gitmap-cli\gitmap.exe`).
- ❌ Leaving user shell sessions invoking legacy binaries that lack recent platform adaptation and safe removal enhancements.

---

## 3. Permitted & Mandatory Actions

- ✅ **Platform-Adaptive Removal:** Always route macro execution through `macro.AdaptCommandForPlatform(cmdText)` (`cli/macro/safe_rm.go`), which dynamically transforms Windows removal commands into safe, idempotent PowerShell loops:
  ```powershell
  foreach ($__target in @('<target>')) {
      if (Test-Path -LiteralPath $__target) {
          Remove-Item -Recurse -Force -LiteralPath $__target
      } elseif (Test-Path -Path $__target) {
          Remove-Item -Recurse -Force -Path $__target
      }
  }
  ```
- ✅ **Native CLI Safe Removal:** Use `gitmap safe-rm <path...> [--force]` (alias `gitmap rm-safe`), which guarantees exit code 0 when targets are already absent.
- ✅ **Multi-Target Binary Synchronization:** Whenever compiling or deploying runtime binaries, ensure synchronization across all active binary targets on the host:
  1. `./gitmap.exe` and `./bin/gitmap.exe`
  2. `%LOCALAPPDATA%\gitmap-cli\gitmap.exe` (`C:\Users\Administrator\AppData\Local\gitmap-cli\gitmap.exe`)
  3. `C:\Users\Alim\AppData\Local\gitmap-cli\gitmap.exe`
  4. `C:\Users\Alim\AppData\Local\gitmap\gitmap.exe`
- ✅ **Macro Editor Step Persistence:** Ensure `gitmap macro edit` clearly distinguishes between persistent macro steps and in-builder diagnostic commands so user steps are reliably recorded.

---

## 4. Root Cause

Executing `gitmap alim1` failed terminatingly with `exit status 1` (`ItemNotFoundException` from `Remove-Item`) on missing target directory `C:\Users\Alim\test` because user `Alim` was executing a stale runtime binary (`v6.442.0` at `C:\Users\Alim\AppData\Local\gitmap\gitmap.exe`) that preceded the introduction of `safe-rm` and `AdaptCommandForPlatform` platform adaptation shims, combined with an interactive edit session that did not append the user's setup steps.
