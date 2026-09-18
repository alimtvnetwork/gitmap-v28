# Spec 85 — Winget Package Distribution

## Status: Research / Future

## Overview

Publish `gitmap-v28` to the Windows Package Manager (winget) so users can install via:

```powershell
winget install AliMTVNetwork.GitMap
```

## How Winget Packages Work

Winget uses YAML manifest files stored in the `microsoft/winget-pkgs` GitHub repository. There is no binary hosting — manifests point to existing download URLs (our GitHub release assets).

A package submission requires **three manifest files** (multi-file format):

1. **Version manifest** (`AliMTVNetwork.GitMap.yaml`)
2. **Default locale manifest** (`AliMTVNetwork.GitMap.locale.en-US.yaml`)
3. **Installer manifest** (`AliMTVNetwork.GitMap.installer.yaml`)

## Required Manifest Files

### `AliMTVNetwork.GitMap.yaml` (version)

```yaml
PackageIdentifier: AliMTVNetwork.GitMap
PackageVersion: 2.49.1
DefaultLocale: en-US
ManifestType: version
ManifestVersion: 1.6.0
```

### `AliMTVNetwork.GitMap.locale.en-US.yaml` (locale)

```yaml
PackageIdentifier: AliMTVNetwork.GitMap
PackageVersion: 2.49.1
PackageLocale: en-US
Publisher: AliMTVNetworkSolutions
PublisherUrl: https://github.com/alimtvnetwork
PackageName: GitMap
PackageUrl: https://github.com/alimtvnetwork/gitmap-v28
License: MIT
LicenseUrl: https://github.com/alimtvnetwork/gitmap-v28/blob/main/LICENSE
ShortDescription: Git repository scanner, manager, and navigator CLI tool.
Description: GitMap scans, catalogs, and manages Git repositories across your machine. It provides cloning, grouping, aliasing, release management, and more.
Tags:
  - git
  - cli
  - devtools
  - repository
  - manager
ManifestType: defaultLocale
ManifestVersion: 1.6.0
```

### `AliMTVNetwork.GitMap.installer.yaml` (installer)

```yaml
PackageIdentifier: AliMTVNetwork.GitMap
PackageVersion: 2.49.1
InstallerType: zip
NestedInstallerType: portable
NestedInstallerFiles:
  - RelativeFilePath: gitmap.exe
    PortableCommandAlias: gitmap-v28
Installers:
  - Architecture: x64
    InstallerUrl: https://github.com/alimtvnetwork/gitmap-v28/releases/download/v2.49.1/gitmap-v28.49.1-windows-amd64.zip
    InstallerSha256: <SHA256_OF_ZIP>
  - Architecture: arm64
    InstallerUrl: https://github.com/alimtvnetwork/gitmap-v28/releases/download/v2.49.1/gitmap-v28.49.1-windows-arm64.zip
    InstallerSha256: <SHA256_OF_ZIP>
ManifestType: installer
ManifestVersion: 1.6.0
```

## Submission Steps

### First-Time Submission

1. **Install wingetcreate:**
   ```powershell
   winget install wingetcreate
   ```

2. **Generate manifests (interactive):**
   ```powershell
   wingetcreate new https://github.com/alimtvnetwork/gitmap-v28/releases/download/v2.49.1/gitmap-v28.49.1-windows-amd64.zip
   ```

3. **Validate manifests:**
   ```powershell
   winget validate --manifest <manifest-dir>
   ```

4. **Test locally:**
   ```powershell
   winget install --manifest <manifest-dir>
   ```

5. **Submit PR to `microsoft/winget-pkgs`:**
   ```powershell
   wingetcreate submit <manifest-dir>
   ```
   This creates a PR automatically. Requires a GitHub PAT with `public_repo` scope.

### Version Updates

```powershell
wingetcreate update AliMTVNetwork.GitMap --version 2.49.1 \
  --urls https://github.com/alimtvnetwork/gitmap-v28/releases/download/v2.49.1/gitmap-v28.49.1-windows-amd64.zip \
  --submit --token <GITHUB_PAT>
```

## CI/CD Automation

Add to `.github/workflows/release.yml`:

```yaml
- name: Publish to Winget
  if: startsWith(github.ref, 'refs/tags/v')
  env:
    WINGET_PAT: ${{ secrets.WINGET_PAT }}
  run: |
    wingetcreate update AliMTVNetwork.GitMap \
      --version ${{ github.ref_name }} \
      --urls "https://github.com/alimtvnetwork/gitmap-v28/releases/download/${{ github.ref_name }}/gitmap-${{ github.ref_name }}-windows-amd64.zip" \
      --submit --token $WINGET_PAT
```

> **Note:** `wingetcreate` runs on Windows only. The release job needs `runs-on: windows-latest`.

## Requirements Before Implementation

- [ ] Choose a stable `PackageIdentifier` (proposed: `AliMTVNetwork.GitMap`)
- [ ] Create a GitHub PAT with `public_repo` scope and store as `WINGET_PAT` secret
- [ ] Decide if ARM64 should be included from day one
- [ ] Ensure release assets follow the exact naming convention the manifest references
- [ ] First PR must pass Microsoft's automated validation bot + community review

## Risks

- **Review time**: First submission typically takes 3-5 days for Microsoft bot + human review
- **Naming**: The `PackageIdentifier` is permanent — choose carefully
- **Windows-only CI**: `wingetcreate` only runs on Windows, so the release job must use `windows-latest`
- **ZIP + portable**: Winget's `zip` + `portable` installer type is newer — verify it works with `NestedInstallerFiles`

## Comparison: Chocolatey vs Winget

| Aspect | Chocolatey | Winget |
|--------|-----------|--------|
| Package format | `.nupkg` (NuGet) | YAML manifests |
| Hosting | Push to chocolatey.org | PR to microsoft/winget-pkgs |
| Review | Human moderation (1-7 days) | Bot + human (3-5 days) |
| CI tool | `choco push` | `wingetcreate update --submit` |
| User base | Developers, sysadmins | All Windows 10/11 users |
| Auto-update | Via `choco upgrade` | Via `winget upgrade` |

## References

- https://learn.microsoft.com/en-us/windows/package-manager/package/manifest
- https://learn.microsoft.com/en-us/windows/package-manager/package/
- https://github.com/microsoft/winget-pkgs
- https://github.com/microsoft/winget-create
