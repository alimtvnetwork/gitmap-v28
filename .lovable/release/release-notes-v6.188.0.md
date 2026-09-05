## Quick Install v6.188.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.188.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.188.0/install.sh | bash
```

## Changelog v6.188.0

- Refactored lazy regular expression engine in gitmap/lazyregex with thread-safe global maps and compiled regex pattern caching
- Implemented reusable pkg/regexnew in coding guidelines codebase with New Creator pattern, batch registration, and nil-safe predicates
- Refactored pipelinedb models and code generator with dedicated enum subpackages, gofmt tab-alignment, and type aliases
- Documented Rule 5 (Lazy Regex & Global Map Deduplication) in cross-language regex usage guidelines
- Updated root readme.md pinned version to v6.188.0 and synchronized all SSoT manifests
