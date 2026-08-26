# Code Signing

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Purpose

Defines the pattern for integrating code signing into the release pipeline. Windows binaries SHOULD be signed to avoid SmartScreen warnings and to establish publisher trust. SHOULD (not MUST) is intentional: code signing requires a paid certificate that may be unavailable in early-stage or fork builds. When the certificate is unavailable, the release pipeline MUST (a) skip signing without failing the build, and (b) emit a clearly-labelled `signing-skipped` banner in the GitHub Step Summary so reviewers see the unsigned status. This spec covers SignPath integration as the primary signing provider, with a feature-flag gating pattern for gradual rollout.

---

## Feature-Flag Gating

Code signing is enabled/disabled via a GitHub repository variable (not a secret):

```yaml
- name: Sign Windows binaries
  if: vars.SIGNPATH_SIGNING_ENABLED == 'true'
  uses: signpath/github-action-submit-signing-request@v1
```

| Variable | Value | Effect |
|----------|-------|--------|
| `SIGNPATH_SIGNING_ENABLED` | `true` | Signing step executes |
| `SIGNPATH_SIGNING_ENABLED` | `false` or unset | Signing step skipped silently |

This allows signing to be disabled during development or if the signing service is unavailable, without modifying the workflow file.

---

## SignPath Integration

### Required Secrets

| Secret | Purpose |
|--------|---------|
| `SIGNPATH_API_TOKEN` | Authentication token for SignPath API |
| `SIGNPATH_ORGANIZATION_ID` | Organization identifier |
| `SIGNPATH_PROJECT_SLUG` | Project slug in SignPath |
| `SIGNPATH_SIGNING_POLICY_SLUG` | Signing policy (e.g., `release-signing`) |

### Workflow Step

```yaml
- name: Sign Windows binaries
  if: vars.SIGNPATH_SIGNING_ENABLED == 'true'
  uses: signpath/github-action-submit-signing-request@v1
  with:
    api-token: ${{ secrets.SIGNPATH_API_TOKEN }}
    organization-id: ${{ secrets.SIGNPATH_ORGANIZATION_ID }}
    project-slug: ${{ secrets.SIGNPATH_PROJECT_SLUG }}
    signing-policy-slug: ${{ secrets.SIGNPATH_SIGNING_POLICY_SLUG }}
    artifact-configuration-slug: "exe"
    input-artifact-path: "dist"
    output-artifact-path: "dist"
    wait-for-completion: true
    wait-for-completion-timeout-in-seconds: 600
    parameters: |
      include: "*.exe"
```

### Key Configuration

| Parameter | Value | Purpose |
|-----------|-------|---------|
| `artifact-configuration-slug` | `"exe"` | Tells SignPath to sign PE executables |
| `input-artifact-path` | `"dist"` | Directory containing unsigned binaries |
| `output-artifact-path` | `"dist"` | Overwrite in-place with signed binaries |
| `wait-for-completion` | `true` | Block until signing is done |
| `wait-for-completion-timeout-in-seconds` | `600` | 10 minute timeout |
| `parameters.include` | `"*.exe"` | Only sign `.exe` files, skip archives |

---

## Signature Verification

After signing, verify that Windows binaries were processed. There are two levels of verification:

### Level 1: CI Smoke Check (Linux Runner)

This confirms the signing action completed and `.exe` files exist. It does **not** verify Authenticode signatures (which requires Windows tooling):

```yaml
- name: Verify signatures (smoke check)
  if: vars.SIGNPATH_SIGNING_ENABLED == 'true'
  run: |
    echo "Checking signed Windows binaries..."
    COUNT=0
    for exe in dist/*.exe; do
      [ -f "$exe" ] || continue
      SIZE=$(stat -c%s "$exe")
      echo "  ✓ Present: $(basename "$exe") ($SIZE bytes)"
      COUNT=$((COUNT + 1))
    done
    if [ "$COUNT" -eq 0 ]; then
      echo "::error::No .exe files found after signing step"
      exit 1
    fi
    echo "All $COUNT Windows binaries present after signing."
```

### Level 2: Full Authenticode Verification (Windows Runner — Optional)

For true signature verification, use a Windows runner with `signtool`:

```yaml
- name: Verify Authenticode signatures
  if: vars.SIGNPATH_SIGNING_ENABLED == 'true'
  runs-on: windows-latest
  run: |
    $unsigned = 0
    Get-ChildItem dist/*.exe | ForEach-Object {
      $result = & signtool verify /pa $_.FullName 2>&1
      if ($LASTEXITCODE -ne 0) {
        Write-Error "UNSIGNED: $($_.Name)"
        $unsigned++
      } else {
        Write-Host "  ✓ Verified: $($_.Name)"
      }
    }
    if ($unsigned -gt 0) { exit 1 }
```

> **Known limitation:** Level 2 requires a Windows runner, adding ~2 minutes to the pipeline. Most projects use Level 1 only, relying on SignPath's own verification. Level 2 is recommended for production releases where signature integrity is critical.

---

## Pipeline Placement

Code signing MUST occur **after** binary building and **before** compression/checksums:

```
1. Build binaries (all platforms)
2. → Sign Windows binaries (conditional)
3. → Verify signatures (conditional)
4. Compress and checksum (zip/tar.gz + checksums.txt)
5. Generate install scripts
6. Create GitHub Release
```

If signing is disabled, steps 2–3 are skipped and unsigned binaries proceed directly to compression.

---

## Constraints

- **Sign before compress** — signing operates on raw `.exe` files, not archives
- **In-place replacement** — signed binaries overwrite unsigned ones in `dist/`
- **Feature flag, not conditional logic** — use `vars.SIGNPATH_SIGNING_ENABLED`, never hardcode `if: false`
- **Timeout** — set a generous timeout (600s) to handle queue delays
- **Scope** — only `.exe` files are signed; Linux and macOS binaries are not code-signed

---

## Alternative Signing Providers

| Provider | Action | Notes |
|----------|--------|-------|
| SignPath | `signpath/github-action-submit-signing-request@v1` | Cloud-hosted, free tier available |
| Azure Trusted Signing | `azure/trusted-signing-action@v0.5.0` | Requires Azure subscription |
| DigiCert | Manual `signtool` invocation | Requires Windows runner |

---

## Cross-References

- [Go Binary Release Pipeline](./02-go-binary-deploy/02-release-pipeline.md) — Where signing is used
- [Shared Conventions](./01-shared-conventions.md) — Platform, runner, action version rules
- [Install Script Generation](./04-install-script-generation.md) — Scripts that download signed binaries
- [Self-Update Mechanism](./06-self-update-mechanism.md) — Update flow that distributes signed binaries
- [Self-Update & App Update (Full Specs)](../14-update/00-overview.md) — Client-side update implementation

---

*Code signing — v3.2.0 — 2026-04-10*
