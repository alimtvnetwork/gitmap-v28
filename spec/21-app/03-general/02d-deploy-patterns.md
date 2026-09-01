# Deploy Patterns

Part of [PowerShell Build & Deploy Patterns](02-powershell-build-deploy.md).

## Retry-on-Lock

When deploying to a target that may be in use (especially on Windows),
wrap `Copy-Item` in a retry loop:

```powershell
$maxAttempts = 20
$attempt = 1
while ($true) {
    try {
        Copy-Item $BinaryPath $destFile -Force -ErrorAction Stop
        break
    } catch {
        if ($attempt -ge $maxAttempts) { throw }
        Write-Warn "Target is in use; retrying ($attempt/$maxAttempts)..."
        Start-Sleep -Milliseconds 500
        $attempt++
    }
}
```

## Nested Deploy Structure

Deploy the binary into a named subfolder within the target directory.
This keeps the deploy target organized when multiple tools share the
same parent directory:

```
deploy-target/
└── toolname/
    ├── toolname.exe
    └── data/
        └── config.json
```

The subfolder (not the parent) should be added to the system `PATH`.

## Deploy Target on PATH

The deploy directory should be on the system `PATH` so the tool can
be run from any terminal without specifying the full path.

## Cross-References

- Generic spec: [03-rename-first-deploy.md](../08-generic-update/03-rename-first-deploy.md) §Implementation
- Deploy path resolution: [02-deploy-path-resolution.md](../08-generic-update/02-deploy-path-resolution.md)
