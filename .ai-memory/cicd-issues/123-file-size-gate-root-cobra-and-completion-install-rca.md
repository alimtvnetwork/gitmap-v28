# RCA 123: File Size Gate Failure on root_cobra_completion.go and completion/install.go

## 1. Symptom & Error Telemetry
- **CI Run**: Workflow Run `#38161301385`
- **Failing Step**: File Size Gate / Code Hygiene
- **Failure Output**:
  ```text
  ✖ File size check failed:
    cli/cmd/root_cobra_completion.go: 652 lines (exceeds 300-line limit)
    cli/completion/install.go: 434 lines (exceeds 300-line limit)
  ```
- **Blast Radius**: Commit `14d064d` introduced modifications to shell completions and Cobra command registration, growing both files beyond the 300-line ceiling enforced by `03-ai-scripts/52-file-size-check.py` and Spec 243.4.

## 2. Root Cause Analysis
- `cli/cmd/root_cobra_completion.go` accumulated top-level subcommands and flag documentation across multiple domains (Antigravity commands, tooling commands, app commands, shell completions), reaching 652 lines.
- `cli/completion/install.go` accumulated platform-specific shell registration handlers, PowerShell profile reconciliation, Unix profile reconciliation, legacy block stripping, and string normalization helpers, reaching 434 lines.
- When commit `14d064d` modified both files, the CI file size gate (`03-ai-scripts/52-file-size-check.py --diff HEAD~1 --limit 300`) flagged both files as exceeding the strict 300-line cap.

## 3. Surgical Resolution & Code Mitigations
1. **Decomposed `cli/cmd/root_cobra_completion.go`**:
   - `cli/cmd/root_cobra_completion.go` (~160 lines): Retained core Cobra completion generator and root completion wiring.
   - `cli/cmd/root_cobra_commands_agy.go` (~85 lines): Extracted Antigravity-specific subcommands and flags.
   - `cli/cmd/root_cobra_commands_tools.go` (~220 lines): Extracted developer tooling subcommands and flags.
   - `cli/cmd/root_cobra_commands_apps.go` (~170 lines): Extracted application management subcommands and flags.
2. **Decomposed `cli/completion/install.go`**:
   - `cli/completion/install.go` (~195 lines): Retained core profile resolution, generation dispatch, and top-level installer entrypoints.
   - `cli/completion/install_powershell.go` (~140 lines): Extracted PowerShell-specific profile paths, snippet generation, and script installation.
   - `cli/completion/install_reconcile.go` (~110 lines): Extracted Unix profile reconciliation, marker block replacements, and native EOL normalizers.
3. Formatted all decomposed files with `gofmt -w` and validated with `go -C cli build -p 2 .`.

## 4. Verification & Preventive Measures
- Ran `python 03-ai-scripts/52-file-size-check.py --files cli/cmd/root_cobra_completion.go cli/cmd/root_cobra_commands_agy.go cli/cmd/root_cobra_commands_tools.go cli/cmd/root_cobra_commands_apps.go cli/completion/install.go cli/completion/install_powershell.go cli/completion/install_reconcile.go` confirming 7/7 files pass strictly $\le 300$ lines.
- Ran `python .github/scripts/go-format-check.py --check-only` confirming 0 unformatted Go files across all 4,504 files.
- Ran `python linter-scripts/check-relative-paths.py` and `python linter-scripts/check-forbidden-strings.py` confirming 0 violations.
- Verified local build passes with `go -C cli build -p 2 .`.
