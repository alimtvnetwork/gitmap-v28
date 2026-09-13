## Quick Install v6.224.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.224.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.224.0/install.sh | bash
```

## Changelog v6.224.0

- Added dynamic CPU freeness detection to scale test workers dynamically from 16 to 32+ based on CPU idle headroom
- Enforced wipe-before-write pre-build and pre-test temporary storage cleanup across all build, test, and packaging gates
- Added ClearRepoTestTempDir, ClearRepoSandboxTempDir, and ClearAllRepoTempDirs to cli/tempdir package
- Relocated e2e smoke test worker sandboxes to OS test temp with automatic post-run sweeping
- Automated old CI/CD run pruning (capping at 5 runs) and loose coverage profile cleanup recovering 455+ MB of storage
- Enabled --allow-serial-runners across concurrent golangci-lint invocations to eliminate race conditions
