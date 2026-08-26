# 12 — Code Signing

**Version:** 1.0.0  
**Updated:** 2026-04-17

---

## Purpose

Define how compiled CLI binaries and installer scripts are
cryptographically signed so end users (and Windows SmartScreen /
macOS Gatekeeper) can verify authenticity. Covers Authenticode
(Windows), `codesign` (macOS), and detached GPG signatures
(cross-platform release artifacts).

---

## Why Sign

| Without signing | With signing |
|-----------------|--------------|
| SmartScreen warning: "Unrecognized app" | Publisher name shown, no warning at trusted level |
| Gatekeeper blocks: "from an unidentified developer" | Runs without prompt |
| Tampering undetectable | Mismatch breaks the signature → user warned |
| No publisher attribution | Properties → Digital Signatures shows issuer |

---

## What to Sign

| Artifact | Required? | Mechanism |
|----------|-----------|-----------|
| `<binary>.exe` (Windows) | ✅ Required | Authenticode (`signtool`) |
| `<binary>` (macOS) | ✅ Required | `codesign` + notarization |
| `<binary>` (Linux) | ⚠️ Recommended | Detached GPG `.sig` |
| `install.ps1`, `install.sh` | ✅ Required | Detached GPG `.sig` next to script |
| `<binary>-updater.exe` | ✅ Required | Same as main binary |
| `latest.json`, `checksums.txt` | ✅ Required | Detached GPG `.sig` |

---

## Windows — Authenticode

### Prerequisites

- A code-signing certificate from a trusted CA (DigiCert, Sectigo,
  SSL.com). EV certs avoid SmartScreen warm-up entirely; OV certs
  warm up after enough downloads.
- `signtool.exe` from the Windows SDK.

### Signing command

```powershell
signtool sign `
    /tr   http://timestamp.digicert.com `
    /td   sha256 `
    /fd   sha256 `
    /a    `
    /d    "<binary> CLI" `
    /du   "https://riseup-asia.com" `
    "<binary>.exe"
```

### Verification

```powershell
signtool verify /pa /v "<binary>.exe"
```

A valid result MUST show:

```
Successfully verified: <binary>.exe
```

### Required flags

| Flag | Purpose |
|------|---------|
| `/tr` | Timestamp server URL (RFC 3161). MANDATORY — without it the signature expires when the cert expires. |
| `/td sha256` | Timestamp digest |
| `/fd sha256` | File digest (SHA-1 is forbidden) |
| `/a` | Auto-select best cert from store |
| `/d` | Description shown in UAC and Properties |
| `/du` | URL shown in UAC dialog |

---

## macOS — codesign + notarization

### Sign

```bash
codesign --sign "Developer ID Application: Riseup Asia LLC (TEAMID)" \
         --options runtime \
         --timestamp \
         --force \
         "<binary>"
```

### Notarize (required for Gatekeeper)

```bash
xcrun notarytool submit "<binary>.zip" \
    --apple-id    "$APPLE_ID" \
    --password    "$APP_SPECIFIC_PASSWORD" \
    --team-id     "$TEAM_ID" \
    --wait
```

### Staple

```bash
xcrun stapler staple "<binary>"
```

Stapling embeds the notarization ticket so Gatekeeper works offline.

---

## Detached GPG Signatures (cross-platform)

For install scripts, JSON manifests, and Linux binaries:

```bash
gpg --detach-sign --armor --local-user "$RELEASE_KEY_ID" install.ps1

# Produces: install.ps1.asc

```

Verification (user-side):

```bash
gpg --verify install.ps1.asc install.ps1
```

The public key MUST be published at a stable URL referenced from
the project README (e.g., `https://riseup-asia.com/release-key.asc`).

---

## CI Integration

Signing keys MUST live in CI secret storage, never in the repo:

| Secret | Format | Used by |
|--------|--------|---------|
| `WINDOWS_CERT_PFX` | base64-encoded `.pfx` | `signtool` |
| `WINDOWS_CERT_PASSWORD` | string | `signtool` |
| `APPLE_ID`, `APP_SPECIFIC_PASSWORD`, `TEAM_ID` | string | `notarytool` |
| `GPG_PRIVATE_KEY` | armored ASCII | `gpg --import` |
| `GPG_PASSPHRASE` | string | `gpg --batch` |

CI workflow MUST:

1. Import secrets to a temporary keystore.
2. Sign all release artifacts.
3. **Wipe the keystore** before any other step that might log files.

---

## Constraints

- Timestamping (`/tr` on Windows, `--timestamp` on macOS) is
  MANDATORY. Without it, signatures become invalid the moment the
  cert expires, even on binaries shipped years earlier.
- SHA-1 is forbidden. Always use SHA-256 for both file and
  timestamp digests.
- The `update` command MUST verify the signature on a downloaded
  new binary before deploying it (binary-based update path). See
  [`../14-update/14-checksums-verification.md`](../14-update/14-checksums-verification.md).
- Never commit signing keys, certs, or passphrases — even in
  encrypted form. Use CI secret storage exclusively.
- Public verification keys MUST be published at a stable HTTPS URL
  with a 1-year minimum cache lifetime.

---

## Cross-References

- [04-build-scripts.md](04-build-scripts.md) §Build Step — sign as last step before deploy
- [11-windows-icon-embedding.md](11-windows-icon-embedding.md) — sign AFTER `.syso` is linked, never before
- [`../14-update/14-checksums-verification.md`](../14-update/14-checksums-verification.md) — checksum + signature flow for binary-based update
- [`../14-update/17-release-pipeline.md`](../14-update/17-release-pipeline.md) — CI pipeline that orchestrates signing

---

*Code signing — v1.0.0 — 2026-04-17*
