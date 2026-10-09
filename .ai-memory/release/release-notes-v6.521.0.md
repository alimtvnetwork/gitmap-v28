# GitMap v6.521.0

## What's Changed in v6.521.0

- **Test-suite repair** — fixed ~64 test files broken by the wave-2 cmd package split (tests referencing old `package cmd` locations/symbols). Tests moved to owner subpackages, exported symbols capitalized, cross-package references qualified. Test-only changes; no production code modified. `go build ./...` and `go vet ./...` clean.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.521.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.521.0/install.sh | sh
```
