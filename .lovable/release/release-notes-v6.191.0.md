## Quick Install v6.191.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.191.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.191.0/install.sh | bash
```

## Changelog v6.191.0

- Synchronized release deployment and version bump prompt in 01-prompts/01-release.md to Prompt Version 2.1.0
- Aligned SSoT manifests and verified version propagation across version.json, package.json, constants.go, and readme.md
- Enforced strict zero-tag policy delegating git tags to automated CI release orchestrators
- Passed 100% test verification across lazyregex, regexnew, pipelinedb, and constants packages
- Pinned active repository version to v6.191.0 in root readme.md and .lovable/user-preferences
