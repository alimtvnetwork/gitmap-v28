# Spec 84 — Chocolatey Package Distribution

## Status: Research / Future

## Overview

Publish `gitmap-v27` as a Chocolatey package so users can install via:

```powershell
choco install gitmap-v27
```

## How Chocolatey Packages Work

Chocolatey uses the NuGet infrastructure. A package is a `.nupkg` file containing:

1. **`.nuspec`** — XML metadata (id, version, description, authors, project URL, license)
2. **`tools/chocolateyInstall.ps1`** — Download + install script
3. **`tools/chocolateyUninstall.ps1`** — Cleanup script (optional but recommended)

The package does NOT embed the binary. Instead, `chocolateyInstall.ps1` downloads the `.zip` from our GitHub release and extracts it.

## Required Files

### `gitmap.nuspec`

```xml
<?xml version="1.0" encoding="utf-8"?>
<package xmlns="http://schemas.chocolatey.org/2012/06/nuspec">
  <metadata>
    <id>gitmap-v27</id>
    <version>2.49.1</version>
    <title>GitMap</title>
    <authors>AliMTVNetworkSolutions</authors>
    <projectUrl>https://github.com/alimtvnetwork/gitmap-v27</projectUrl>
    <licenseUrl>https://github.com/alimtvnetwork/gitmap-v27/blob/main/LICENSE</licenseUrl>
    <requireLicenseAcceptance>false</requireLicenseAcceptance>
    <description>Git repository scanner, manager, and navigator CLI tool.</description>
    <tags>git cli devtools repository manager</tags>
    <projectSourceUrl>https://github.com/alimtvnetwork/gitmap-v27</projectSourceUrl>
    <packageSourceUrl>https://github.com/alimtvnetwork/gitmap-v27</packageSourceUrl>
    <releaseNotes>https://github.com/alimtvnetwork/gitmap-v27/releases</releaseNotes>
  </metadata>
  <files>
    <file src="tools/**" target="tools" />
  </files>
</package>
```

### `tools/chocolateyInstall.ps1`

```powershell
$ErrorActionPreference = 'Stop'

$packageArgs = @{
  packageName    = 'gitmap-v27'
  url64bit       = 'https://github.com/alimtvnetwork/gitmap-v27/releases/download/v2.49.1/gitmap-v27.49.1-windows-amd64.zip'
  checksum64     = '<SHA256_HASH>'
  checksumType64 = 'sha256'
  unzipLocation  = "$(Split-Path -Parent $MyInvocation.MyCommand.Definition)"
}

Install-ChocolateyZipPackage @packageArgs
```

### `tools/chocolateyUninstall.ps1`

```powershell
$toolsDir = "$(Split-Path -Parent $MyInvocation.MyCommand.Definition)"
Remove-Item "$toolsDir\gitmap.exe" -Force -ErrorAction SilentlyContinue
```

## Publishing Steps

1. **Create an account** at https://community.chocolatey.org/account/Register
2. **Get an API key** from https://community.chocolatey.org/account
3. **Set API key locally:**
   ```
   choco apikey --key <YOUR_KEY> --source https://push.chocolatey.org/
   ```
4. **Build the package:**
   ```
   choco pack gitmap.nuspec
   ```
5. **Test locally:**
   ```
   choco install gitmap-v27 --debug --verbose --source .
   ```
6. **Push to community feed:**
   ```
   choco push gitmap.2.49.1.nupkg --source https://push.chocolatey.org/
   ```
7. **Wait for moderation** — Chocolatey community packages go through human review (can take 1-7 days).

## CI/CD Automation

Add a step to `.github/workflows/release.yml` after asset upload:

```yaml
- name: Publish to Chocolatey
  if: startsWith(github.ref, 'refs/tags/v')
  env:
    CHOCO_API_KEY: ${{ secrets.CHOCO_API_KEY }}
  run: |
    choco apikey --key $CHOCO_API_KEY --source https://push.chocolatey.org/
    choco pack choco/gitmap.nuspec
    choco push gitmap.*.nupkg --source https://push.chocolatey.org/
```

## Requirements Before Implementation

- [ ] Create Chocolatey community account
- [ ] Obtain and store `CHOCO_API_KEY` as a GitHub Actions secret
- [ ] Decide on directory layout: `choco/gitmap.nuspec` + `choco/tools/`
- [ ] Automate version + checksum substitution in `chocolateyInstall.ps1` during release
- [ ] Consider whether to also support ARM64 (Chocolatey's `url` vs `url64bit`)

## Risks

- **Moderation delays**: First submission can take up to a week for review
- **Version sync**: Must update the `.nuspec` version and download URL on every release
- **Checksum automation**: The SHA256 must be injected after the release assets are uploaded

## References

- https://docs.chocolatey.org/en-us/create/create-packages/
- https://community.chocolatey.org/courses/creating-chocolatey-packages/
- https://docs.chocolatey.org/en-us/create/helpers/install-chocolateyzippackage
