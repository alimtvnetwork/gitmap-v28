## Quick Install v6.267.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.267.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.267.0/install.sh | bash
```

## Changelog v6.267.0

- Improve `gitmap pipeline errors` visual hierarchy with top and bottom newline padding and 2-space indentation on the `Reading pipeline logs...` progress message
- Add leading newline gap before `● Recent Commits Pipeline Summary` cyan table header for cleaner visual spacing
- Display actionable pipeline database cleanup guidance (`• Cleanup: gitmap pipeline clear -y`) directly below database metadata in both clean and failure reports
- Update `printEmptyCachedFailures` to include pipeline cleanup command guidance
- Add unit test `TestPrintRecentCommitsHeaderPadding` to verify consistent header padding
