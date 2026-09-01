# Build Patterns

Part of [PowerShell Build & Deploy Patterns](02-powershell-build-deploy.md).

## Build with Embedded Variables

Use Go's `-ldflags` to embed values at compile time:

```powershell
$absRepoRoot = (Resolve-Path $RepoRoot).Path
$ldflags = "-X 'pkg/constants.RepoPath=$absRepoRoot'"
go build -ldflags $ldflags -o $outPath .
```

## Version Verification

After building, immediately run the binary with `version` to confirm:

```powershell
$versionOutput = & $binaryPath version 2>&1
Write-Info "Version: $versionOutput"
```

This catches build issues early — if the version doesn't match
expectations, the build is suspect.

## Data Folder Copy

If the binary needs companion data files, copy them alongside:

```powershell
if ($Config.copyData) {
    Copy-Item $dataSource $dataDest -Recurse
}
```

## Cross-References

- Generic spec: [04-build-scripts.md](../08-generic-update/04-build-scripts.md) §Build
- Icon embedding: [04-windows-icon-embedding.md](04-windows-icon-embedding.md)
