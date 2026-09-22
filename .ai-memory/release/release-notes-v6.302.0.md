## Quick Install v6.302.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.302.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.302.0/install.sh | bash
```

## Changelog v6.302.0

- AGY Workspace & Pin Management: Added `gitmap agy pins` (`ls`, `add`, `rm`, `edit`, `help`), `gitmap agy rm-rejoin-read` (`rrr`), `gitmap agy rm-rejoin-pin-read` (`rrpr`), supporting sequential IDs, project slugs, and shell completion.
- Raw Error Propagation: Audited codebase error management across CLI commands to guarantee zero swallowed errors; wrapped underlying SQL and filesystem errors with structured `apperror.WrapSimple` / `apperror.Wrap` preserving root causes.
- Antigravity Compilation & Import Hygiene: Added missing `database/sql` import in `cli/cmdagy/agy_history_cmd.go`, fixed `loadAllAgyProjects` signature mismatch in `agy_pin_projects.go` and `agy_pins_edit.go`, and cleaned unused imports in `agy_projects.go`.
- Quad Runner Verification: Verified all 13 modified Go packages pass green via parallel quad runner in local CI/CD.
