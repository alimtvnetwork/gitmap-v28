## Quick Install v6.189.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.189.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.189.0/install.sh | bash
```

## Changelog v6.189.0

- Eliminated redundant secondary regexMap lookups in gitmap/lazyregex and pkg/regexnew, transitioning to single-map pattern deduplication
- Added isCompiled boolean flag, self-contained compiled *regexp.Regexp state, and mutex synchronization to LazyRegexp and LazyRegex
- Introduced Count, IsFound, GroupBy (named capture group extraction), and FindAllGroups methods to LazyRegexp and LazyRegex
- Implemented CompileAppError and CompileBuilder methods returning structured diagnostic AppErrors and AppBuilders on compilation failure
- Updated root readme.md pinned version to v6.189.0 and synchronized all SSoT manifests
