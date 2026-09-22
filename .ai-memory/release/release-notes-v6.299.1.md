## Quick Install v6.299.1

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.299.1/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.299.1/install.sh | bash
```

---

## What's Changed in v6.299.1

- **CI/CD & Gofmt Formatting:** Formatted `cli/cmd/rootcore.go` and `cli/cmdschedule/helpers.go` using canonical `gofmt -w` to eliminate extra blank lines and ensure 100% compliance with `go-format-check.py`.
- **Linter Unit Tests:** Verified all 18/18 test cases pass green in `.github/scripts/tests/test_ci_scripts.py`.
- **4-Part RCA 73:** Recorded `.ai-memory/cicd-issues/73-gofmt-whitespace-drift-in-rootcore-and-helpers-rca.md` and registered in `.ai-memory/cicd-index.md`.
- **Cross-Platform Quality Gates:** All 4 policy linters (nested ifs, boolean guidelines, relative paths, enum & boolean) and cross-platform `go vet` (Linux, Darwin, Windows) pass with 0 errors.
