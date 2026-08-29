# Subtask 18.03: Refactor Explicit Booleans in Scripts

## Target
- Files: `install.ps1`, `install-quick.ps1`, `linter-scripts/tests/`
- Current Score: 90.0 / 100

## Violations to Fix
- [ ] [CG-BOOL-001] Replace `$true -eq $true` explicit comparisons with implicit `$isSuccess`.

## Acceptance Criteria
- [ ] Implicit boolean evaluation only.
- [ ] Zero behavioral regression in installer scripts.
