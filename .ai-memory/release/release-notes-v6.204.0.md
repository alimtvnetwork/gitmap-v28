## Quick Install v6.204.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.204.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.204.0/install.sh | bash
```

## Changelog v6.204.0

- Fix Linux PNPM, Yarn, and Bun installation via npm global and standalone script fallbacks
- Add upfront tool tree hierarchy preview before profile execution and --tree flag support
- Remove Ollama and ubuntu-dev-ai profile, isolating Ollama strictly to standalone ai profile
- Add explicit profile installation examples across CLI help, documentation, and install ls
