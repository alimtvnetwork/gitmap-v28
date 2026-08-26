# 09 — Version Verification (Three-Branch)

**Version:** 1.0.0  
**Updated:** 2026-04-17

---

## Purpose

Define the post-update verification logic that confirms the update
actually took effect, and produces clear diagnostic output when it
did not. This solves the common failure mode where a deploy succeeds
but the PATH binary still resolves to the *old* file.

---

## Two Binaries, Two Truths

After deploy, there are up to **two distinct binaries** to verify:

| Name | Where | How resolved |
|------|-------|--------------|
| **Active binary** | The one on `PATH` | `Get-Command <binary>` (PowerShell) / `command -v <binary>` (Bash) |
| **Deployed binary** | The deploy target | `<deployPath>/<binary>/<binary>.exe` from config |

These can disagree when:

- The PATH binary is in a different directory than the deploy target.
- A previous PATH-sync step failed silently.
- The user has multiple installs and PATH resolves to the wrong one.

---

## Three-Branch Decision Table

The verification step compares both versions and chooses one of three
outcomes:

| Active version | Deployed version | Outcome | Exit |
|----------------|------------------|---------|------|
| `<new>` | `<new>` | ✅ Success — both match | 0 |
| `<new>` | `unknown` | ⚠️ Warning — active binary is correct, deployed path is misconfigured | 0 |
| `unknown` OR active ≠ deployed | any | ❌ Failure — PATH still serves the old binary | 1 |

**Why warning, not failure, when only deployed is unknown?** A
misconfigured `deployPath` in the config file should not block an
otherwise-successful update. The user is still on the new version —
they just need to fix their config.

---

## Required Trace Output

When verification produces a warning or failure, the worker MUST emit
trace-level diagnostic lines so the user can self-diagnose without
reading source code.

```
  Version before:   <binary> v1.2.0
  Version active:   <binary> v1.3.0
  Version deployed: unknown
  Active binary:    C:\Users\user\bin\<binary>.exe
  Deployed binary:  (not resolved)

  [WARN] Deployed binary could not be verified.
  [TRACE] activeAfter=v1.3.0  deployedAfter=unknown
  [HINT] Check that <config-file> 'deployPath' points to the correct
         directory and that the binary exists at: <expected-path>
  [OK]   Active PATH binary updated successfully: v1.3.0
```

### Required trace points

| Trace line | Emitted when |
|------------|--------------|
| `deployedBinary: not resolved` | `deployPath` is unset or config missing |
| `deployedBinary: path not found: <path>` | Config resolved but file does not exist |
| `Get-Command <binary>: <path>` | Always — shows what PATH actually resolves to |
| `activeAfter=<v> deployedAfter=<v>` | Always — single-line machine-grep summary |

---

## Implementation

```
func verifyUpdate(deployedPath, expectedVersion string) (int, error) {
    activePath := resolveActiveBinary()           // Get-Command / command -v
    activeVer  := runVersion(activePath)
    deployedVer := runVersion(deployedPath)        // "unknown" if path nil

    printTrace(activePath, activeVer, deployedPath, deployedVer)

    switch {
    case activeVer == "unknown":
        return 1, errors.New("active binary version unreadable")
    case activeVer != expectedVersion:
        return 1, fmt.Errorf("active still at %s, expected %s",
            activeVer, expectedVersion)
    case deployedVer == "unknown":
        printWarn("deployed path misconfigured — active is correct")
        return 0, nil
    case activeVer != deployedVer:
        return 1, fmt.Errorf("active=%s != deployed=%s",
            activeVer, deployedVer)
    default:
        return 0, nil
    }
}
```

---

## Constraints

- Verification MUST run before `update-cleanup`. The `.old` backup
  is the rollback safety net — do not delete it until success is
  confirmed.
- The version string MUST come from `<binary> version`, not from a
  static constant in the worker. Stale constants would mask real
  failures.
- Trace output MUST be on by default during update — this is not
  debug output, it is the user's only visibility into a fail mode.
- Never short-circuit verification on a "successful" deploy. Deploy
  succeeding only proves the file was written; it does not prove the
  user's PATH now serves it.

---

## Cross-References

- [01-self-update-overview.md](01-self-update-overview.md) §Version Comparison
- [03-rename-first-deploy.md](03-rename-first-deploy.md) §PATH Sync
- [05-handoff-mechanism.md](05-handoff-mechanism.md) §Two-Phase Summary
- [06-cleanup.md](06-cleanup.md) §Automatic Cleanup (must run after verify)

---

*Version verification — v1.0.0 — 2026-04-17*
