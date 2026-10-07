# GitMap v6.506.0

## What's Changed in v6.506.0

- **Cursor Settings Guard Validation**: Added `requireCursorSettingsFile` in `cli/cmdcursor/cursor_settings.go` with descriptive `E1037` error envelopes prompting users to run `gitmap cursor settings apply`.
- **Views AI Settings UI Modernization**: Refactored `src/pages/Settings.tsx` to invert `/api/instances` fetch validation into clean guard clauses, eliminating nested `if` statements.
- **Process Discovery Spacing Hygiene**: Enforced vertical blank line spacing rules in `cli/cmdagy/agy_instance_discovery.go` and `cli/cmdcursor/cursor_settings_test.go`.
- **Quality Gates Verification**: Verified zero violations across `check-nested-ifs.py`, `check-enum-and-boolean.py`, and `check-relative-paths.py`.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.506.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.506.0/install.sh | sh
```
