# GitMap v6.536.0

## What's Changed in v6.536.0

- **CI/CD Build Fix**: Removed unused `cmdsummary` import in `cli/cmd/rootutility.go` that triggered compiler failure across all target platforms (Ubuntu, Windows, macOS) during release binary compilation and `go test -race`.
- **RCA Documentation**: Root cause analysis and mitigation details documented in `.ai-memory/cicd-issues/114-unused-import-cmdsummary-in-rootutility-rca.md`.
- **Full Module Verification**: Verified complete module compilation (`go build ./...`) across all packages with zero compilation errors.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.536.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.536.0/install.sh | sh
```

