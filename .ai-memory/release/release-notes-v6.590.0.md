# GitMap v6.590.0

## What's Changed in v6.590.0

- **AppError migration completed** — bare fmt.Errorf converted to typed AppError returns across all production packages
- **Coding-guideline execute suite finished** — all 39 prompts applied, partially applied with documented remainders, or closed with evidence
- **Deduplicated helpers** — filesystem and string helpers consolidated into shared packages, 51 local copies removed
- **Dispatch and matching efficiency** — map-based dispatch, case-insensitive match helpers, per-call regex compilation eliminated
- **Enum compliance** — linter-verified, 9 violations fixed; heavy test suite tagged and wired into CI
- **Build repaired** — deduplicated constants and store types, Result-envelope callers aligned; full build and vet green
- **Plan backlog drained** — 24 completed plans filed away

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.590.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.590.0/install.sh | sh
```
