# 329.1 — Version Update Downgrade Guard & Sibling Probe Specification

> **Status:** Drafted  
> **Subsystem:** `cli/cmdupdate/`  
> **Parent Ledger:** `02-spec/21-app/329-gitmap-context-menu-and-version-fix/00-master-audit-ledger.md`  
> **Rule Compliance:** Leaf Error Sovereignty, Positive Booleans (`is*`/`has*`), Functions <= 15 lines.

---

## 1. Problem Statement & Root Cause Analysis

In recent builds (e.g. `v6.613.0`), invoking `gitmap update` resulted in an unexpected downgrade to `v6.610.0` (as captured in `assets/screenshots/gitmap-context-menu-and-version-fix-01.png`):

```text
Checking for updates...
Discovered release via sibling-probe: v6.610.0
Downgrading / updating gitmap: v6.613.0 -> v6.610.0
```

### Root Cause
1. **Equality-Only Version Check:** In `cli/cmdupdate/update_version_target.go`, `IsAlreadyUpdated(currentVersion, targetVersion, isForce)` only evaluates string equality (`normCurrent == normTarget`). It does NOT perform semantic version comparison (`semver.Compare(target, current) <= 0`). Therefore, an older target version (e.g. `6.610.0`) is considered "not yet updated" rather than "older than current", causing the update loop to proceed with downloading and replacing the running binary.
2. **Unvalidated Sibling Repository Discovery:** In `cli/cmdupdate/updateprobe.go`, `resolveLatestRepoSlug` probes sibling repositories (`gitmap-v29`, `gitmap-v30`, etc.) by sending HTTP HEAD requests to repository root endpoints. If a sibling repository exists on GitHub, it is immediately returned as the resolved target repository slug without verifying whether that repository actually contains releases newer than the current version. If a sibling repository was created earlier or only contains experimental/older tags, `gitmap update` queries its release feed and resolves a release that is older than the current repository's latest release.

---

## 2. Architectural Design & Contract

### 2.1 Strict Monotonic SemVer Guard
1. **Comparison Helper:** Leverage SemVer comparison logic (`isNewer(tagA, tagB string) bool` or `CompareSemVer(v1, v2 string) int`) in `cli/cmdupdate/update_version_target.go`.
2. **Downgrade Rejection Rule:** An update candidate is skipped if `CompareSemVer(cand.Version, currentVersion) <= 0`, UNLESS:
   - The user explicitly specified `--allow-downgrade`, OR
   - The user explicitly requested an exact target version via `--version <tag>` with `--force`.
3. **User Communication:** When skipped because the running version is newer or equal, the CLI outputs a clear, informative message:
   ```text
   Already up to date (current v6.613.0 >= target v6.610.0). No update needed.
   Use --allow-downgrade to install an older release.
   ```

### 2.2 Sibling Probe Monotonicity Verification
1. When probing sibling repositories (`gitmap-v(N+1)` through `gitmap-v(N+max)`), candidate sibling repositories must only be adopted if their latest release tag is strictly newer than `constants.Version` (and newer than the current repository's latest release).
2. If all sibling repositories have older releases or no releases with assets for the current OS/architecture, the probe falls back cleanly to `constants.UpdateCurrentRepoSlug`.

### 2.3 CLI Option Expansion
Update `UpdateOptions` struct in `cli/cmdupdate/update_options.go` / `update_version_target.go`:
```go
type UpdateOptions struct {
    TargetVersion   string
    IsForce         bool
    IsDryRun        bool
    IsJSON          bool
    IsQuiet         bool
    IsAllowDowngrade bool
    MaxFallbackTags int
}
```
CLI flags `--allow-downgrade` and `-d` set `IsAllowDowngrade = true`.

---

## 3. Implementation Checklist & Files

| File Path | Description of Changes |
| :--- | :--- |
| `cli/cmdupdate/update_version_target.go` | Implement `CompareSemVer`, update `IsAlreadyUpdated` to reject downgrades without `IsAllowDowngrade`, add `--allow-downgrade` flag parsing. |
| `cli/cmdupdate/update_binary_swap.go` | Update `isUpdateSkipped` to display downgrade prevention notices. |
| `cli/cmdupdate/updateprobe.go` | Ensure sibling probe resolution falls back to current slug if sibling releases are older than current version. |
| `cli/cmdupdate/update_version_target_test.go` | Unit tests for version comparison, downgrade prevention, and flag parsing. |
