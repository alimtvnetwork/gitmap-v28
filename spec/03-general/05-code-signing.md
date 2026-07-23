# Code Signing with SignPath.io

## Overview

Code signing ensures Windows users see a trusted publisher name
instead of "Unknown publisher" warnings. This project uses
[SignPath.io](https://signpath.io), which provides **free code
signing certificates for open-source projects**.

## Why Code Signing Matters

| Without Signing | With Signing |
|----------------|--------------|
| SmartScreen: "Windows protected your PC" | No SmartScreen warning (after reputation) |
| UAC shows "Unknown publisher" | Shows "Riseup Asia LLC" (or org name) |
| Antivirus false positives common | Reduced false positive rate |
| Users must click "Run anyway" | Clean, trusted first-run experience |
| Downloads flagged by browsers | Recognized as legitimate software |

## SignPath.io — How It Works

### Eligibility

SignPath offers free certificates to OSS projects that meet:
- Public GitHub repository
- Active project with meaningful commits
- OSI-approved license

### Architecture

```
GitHub Actions (release.yml)
  └─ Build binaries (go build)
      └─ Submit .exe files to SignPath API
          └─ SignPath signs with org certificate
              └─ Download signed binaries
                  └─ Attach to GitHub Release
```

The private key **never leaves SignPath's infrastructure** — binaries
are uploaded, signed server-side, and returned. No key management
or hardware tokens needed on your end.

### Certificate Type

SignPath issues **standard (OV) code signing certificates** for OSS:
- Organization Validation — shows your org name
- Timestamped — signature remains valid after certificate expires
- SHA-256 digest — meets modern Windows requirements

## Setup Steps

### 1. Register with SignPath

1. Go to [signpath.io](https://signpath.io) and sign up
2. Apply for the **OSS program** (free tier)
3. Link your GitHub repository
4. SignPath reviews the project (typically 1–5 business days)

### 2. Configure Signing Policy

In the SignPath dashboard:

1. Create a **Project** named `gitmap-v27`
2. Create a **Signing Policy** (e.g., `release-signing`)
3. Define an **Artifact Configuration**:

```xml
<artifact-configuration xmlns="http://signpath.io/artifact-configuration/v1">
  <pe-file>
    <authenticode-sign hash-algorithm="sha256" />
  </pe-file>
</artifact-configuration>
```

### 3. Create GitHub Integration

SignPath provides a GitHub Action. Add these secrets to the repo:

| Secret | Source | Description |
|--------|--------|-------------|
| `SIGNPATH_API_TOKEN` | SignPath dashboard → API Tokens | Authentication token |
| `SIGNPATH_ORGANIZATION_ID` | SignPath dashboard → Organization | Org identifier |
| `SIGNPATH_PROJECT_SLUG` | Project settings | e.g., `gitmap-v27` |
| `SIGNPATH_SIGNING_POLICY_SLUG` | Policy settings | e.g., `release-signing` |

### 4. CI Pipeline Integration

Add the signing step to `.github/workflows/release.yml` **after**
building Windows binaries but **before** compression and checksum
generation:

```yaml
- name: Sign Windows binaries
  if: runner.os == 'Linux'  # SignPath action runs on any OS
  uses: signpath/github-action-submit-signing-request@v1
  with:
    api-token: ${{ secrets.SIGNPATH_API_TOKEN }}
    organization-id: ${{ secrets.SIGNPATH_ORGANIZATION_ID }}
    project-slug: ${{ secrets.SIGNPATH_PROJECT_SLUG }}
    signing-policy-slug: ${{ secrets.SIGNPATH_SIGNING_POLICY_SLUG }}
    artifact-configuration-slug: "exe"
    input-artifact-path: "dist/gitmap-windows-amd64.exe"
    output-artifact-path: "dist/gitmap-windows-amd64.exe"
    wait-for-completion: true
    wait-for-completion-timeout-in-seconds: 300
```

Repeat for each Windows target (amd64 and arm64).

### 5. Signing Both Windows Targets

```yaml
- name: Sign Windows binaries
  run: |
    for binary in dist/gitmap-windows-amd64.exe dist/gitmap-windows-arm64.exe; do
      if [ -f "$binary" ]; then
        echo "Signing: $binary"
        # Submit to SignPath and wait for signed binary
      fi
    done
```

### 6. Updater Binary

The `gitmap-updater` Windows binaries should also be signed:

```yaml
# Sign both gitmap-v27 and gitmap-updater Windows builds
WINDOWS_BINARIES=(
  "dist/gitmap-windows-amd64.exe"
  "dist/gitmap-windows-arm64.exe"
  "dist/gitmap-updater-windows-amd64.exe"
  "dist/gitmap-updater-windows-arm64.exe"
)
```

## Pipeline Order (Updated)

The signing step inserts between build and compress:

| Step | Action |
|------|--------|
| 1 | Checkout + Setup Go |
| 2 | Resolve version from tag |
| 3 | Build all targets (go build) |
| 4 | Generate Windows resources (go-winres) |
| **5** | **Sign Windows binaries (SignPath)** |
| 6 | Compress (zip/tar.gz) |
| 7 | Generate checksums |
| 8 | Generate install scripts |
| 9 | Publish GitHub Release |

**Critical**: Signing MUST happen before compression and checksums.
The checksum must reflect the signed binary, not the unsigned one.

## Verification

### Manual Verification

After a signed release, verify on Windows:

```powershell
# Check digital signature
Get-AuthenticodeSignature .\gitmap.exe

# Expected output:
# SignerCertificate: [Thumbprint]  Riseup Asia LLC
# Status:           Valid
# StatusMessage:    Signature verified.
```

### CI Verification

Add a post-sign validation step:

```yaml
- name: Verify signature
  if: runner.os == 'Windows'
  run: |
    $sig = Get-AuthenticodeSignature dist/gitmap-windows-amd64.exe
    if ($sig.Status -ne "Valid") {
      Write-Error "Signature verification failed: $($sig.StatusMessage)"
      exit 1
    }
    Write-Host "Signature valid: $($sig.SignerCertificate.Subject)"
```

## Timestamping

All signatures must include a timestamp from a trusted TSA
(Timestamp Authority). This ensures the signature remains valid
even after the signing certificate expires:

```
Timestamp server: http://timestamp.digicert.com
Digest algorithm: SHA-256
```

SignPath handles timestamping automatically — no manual configuration
needed.

## SmartScreen Reputation

Even with a valid certificate, Windows SmartScreen uses a
reputation system:

| Stage | SmartScreen Behavior |
|-------|---------------------|
| First signed release | May still show warning (building reputation) |
| After ~100-500 installs | Warning disappears for most users |
| EV certificate | Immediate reputation (but costs $300+/year) |

SignPath's OV certificate builds reputation over time. Each signed
release with real downloads increases trust.

## Constraints

- **Sign before compress** — checksums must cover the signed binary
- **Never store private keys in CI** — SignPath handles key management
- **Always timestamp** — prevents expiry-related failures
- **Sign all Windows binaries** — gitmap-v27 + gitmap-updater
- **Verify after signing** — CI must validate the signature
- **Build-once rule** — sign the same binary that was built, never rebuild

## Costs

| Item | Cost |
|------|------|
| SignPath OSS program | Free |
| Certificate issuance | Free (included) |
| Signing operations | Free (OSS tier) |
| Renewal | Free (as long as project remains active OSS) |

## Alternative: Azure Trusted Signing

If SignPath is unavailable or the project goes private:

| Feature | SignPath (OSS) | Azure Trusted Signing |
|---------|---------------|----------------------|
| Cost | Free | ~$10/month |
| Certificate type | OV | Microsoft-managed |
| SmartScreen | Builds over time | Immediate |
| Setup | GitHub integration | Azure subscription |
| Key management | Server-side | Cloud HSM |

## Cross-References

| Topic | Spec |
|-------|------|
| Release pipeline | [02-release-pipeline.md](../07-generic-release/02-release-pipeline.md) |
| CI architecture | [release.yml](../../.github/workflows/release.yml) |
| Icon embedding | [04-windows-icon-embedding.md](04-windows-icon-embedding.md) |
| Build & deploy | [09-build-deploy.md](../01-app/09-build-deploy.md) |
| Chocolatey package | [84-chocolatey-package.md](../01-app/84-chocolatey-package.md) |
| Winget package | [85-winget-package.md](../01-app/85-winget-package.md) |

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
