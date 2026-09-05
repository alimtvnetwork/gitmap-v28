## Quick Install v6.190.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.190.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.190.0/install.sh | bash
```

## Changelog v6.190.0

- Refactored lazy regex engine with strict thread-safe compile locking across gitmap/lazyregex and pkg/regexnew
- Introduced CompileResult envelope wrapping compiled regexp, structured AppError, and fluent AppBuilder diagnostics
- Added dedicated GroupMap data type with rich query, mutation, cloning, and serialization methods
- Added dedicated GroupList data type with bounds-safe indexing, key deduplication, and predicate filtering
- Updated root readme.md pinned version to v6.190.0 and synchronized all SSoT manifests
