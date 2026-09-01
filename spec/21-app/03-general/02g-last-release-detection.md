# Last Release Detection Script

Part of [PowerShell Build & Deploy Patterns](02-powershell-build-deploy.md).

## Purpose

A standalone PowerShell script (`scripts/Get-LastRelease.ps1`) resolves
and displays the latest released version. Keeping this in a separate
file avoids bloating `run.ps1` and allows reuse from any context
(manual invocation, CI, update scripts).

## Resolution Order

The script uses a three-tier fallback strategy:

| Priority | Source | Method |
|----------|--------|--------|
| 1 | Binary | `toolname list-versions --limit 1` — parses first `vX.Y.Z` from output |
| 2 | JSON | `.gitmap/release/latest.json` — reads `tag` or `version` field |
| 3 | Git tag | `git tag --list "v*" --sort=-version:refname` — first stable `vX.Y.Z` |

## Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `-BinaryPath` | string | `""` | Path to CLI binary; falls back to `Get-Command` |
| `-RepoRoot` | string | `""` | Repo root for `latest.json` lookup; falls back to CWD |
| `-Label` | string | `"Last release"` | Display label prefix |

## Output

```
  Last release:    v2.24.0 (binary)
```

The parenthetical suffix indicates which source resolved the version:
`binary`, `latest.json`, or `git tag`. If all sources fail, prints
`unknown`.

## Integration Points

**`run.ps1`** — called after the final "All done!" message:

```powershell
$lastReleaseScript = Join-Path (Join-Path $RepoRoot "scripts") "Get-LastRelease.ps1"
if (Test-Path $lastReleaseScript) {
    & $lastReleaseScript -BinaryPath $binaryPath -RepoRoot $RepoRoot
}
```

**Update script** — embedded in the version-verify section between
version lines and the active/deployed match check:

```powershell
$lastReleaseScript = Join-Path (Join-Path "<repoPath>" "scripts") "Get-LastRelease.ps1"
if (Test-Path $lastReleaseScript) {
    & $lastReleaseScript -BinaryPath $activeBinary -RepoRoot "<repoPath>"
}
```

## Design Rules

- **No error exits** — the script always succeeds; missing data shows
  `unknown`.
- **Three-tier fallback** ensures a result even when the binary is
  unavailable (fresh clone) or `.gitmap/release/` metadata hasn't been
  generated yet.
- **Separate file** keeps `run.ps1` lean and allows reuse from any
  context.

## Cross-References

- Generic spec: [04-build-scripts.md](../08-generic-update/04-build-scripts.md) §Validation
