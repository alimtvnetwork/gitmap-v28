## Quick Install v6.301.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.301.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.301.0/install.sh | bash
```

---

## What's Changed in v6.301.0

- **Interface Naming Compliance:** Renamed all internal `cmdos` orchestrator interfaces (`AutoLoginEngine` -> `AutoLoginOperator`, `SystemCleanEngine` -> `SystemCleanOperator`, `DNSEngine` -> `DNSOperator`, `ThemeEngine` -> `ThemeOperator`, `TweakEngine` -> `TweakOperator`) to strictly adhere to Go naming guidelines (`*er` / `*or` suffix) enforced by `check-interface-naming.py`.
- **Enum Naming Compliance:** Renamed `ThemeMode` to `ThemeModeType` in `cli/cmdos/os_theme_types.go` and cross-platform handlers to satisfy the mandatory `*Type` suffix rule enforced by `check-enum-guidelines.py`.
- **Quality Gate Self-Healing & Verification:** All 43 CI/CD quality gates in `03-ai-scripts/06-cicd-local-runner.py` executed and passed 100% green (including E2E smoke tests, history pin/purge, and race detector).
- **4-Part RCA 74:** Recorded root-cause analysis in `.ai-memory/cicd-issues/74-interface-naming-and-enum-suffix-compliance-rca.md` and registered in `.ai-memory/cicd-index.md`.
