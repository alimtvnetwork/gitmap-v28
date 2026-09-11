## Quick Install v6.213.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.213.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.213.0/install.sh | bash
```

## Changelog v6.213.0

- Add live execution feedback and --exec toggle during interactive macro creation
- Add universal path expansion for %TEMP%, //temp, /temp, , and ~ across macro engine and cd tracking
- Fix CI/CD race in-memory DB collision and anchor lifecycle in store
- Fix process isolation for concurrent E2E smoke test workers on Windows
- Ensure readme platform badges and single-line install section compliance
