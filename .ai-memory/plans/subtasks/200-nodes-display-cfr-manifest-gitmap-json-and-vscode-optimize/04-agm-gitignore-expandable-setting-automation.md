# Subtask 04: AGM Gitignore Expandable Setting and Scan Automation

> **Parent Plan:** `200-nodes-display-cfr-manifest-gitmap-json-and-vscode-optimize.md`
> **Status:** `PENDING`
> **Target Files:**
> - `cli/config/config.go`
> - `cli/model/record.go`
> - `cli/gitignoreagm/cli_prompt.go`
> - `cli/cmdscan/scan.go`

---

## Technical Specification

1. **Setting Definition**:
   - `autoGitignoreAgm` toggle.
   - Default: false (prompts user when interactive, or logs notice when non-interactive).
   - When enabled (true): automatically deletes from git, commits deletion, updates `.gitignore`, and commits `.gitignore` during `gitmap scan` without user prompt!
   - Can be set via:
     - `gitmap gitignore agm enable-scan-auto`
     - `gitmap gitignore agm disable-scan-auto`
     - `gitmap gitignore agm status`
     - Or in `store.DB` (`Settings` table: `autoGitignoreAgm` = `"true"` / `"false"`).
     - Or CLI flags on `scan`: `gitmap scan --auto-gitignore-agm` / `--no-auto-gitignore-agm`.

2. **Integration in `cli/cmdscan/scan.go`**:
   - In `checkAgmResumeTaskOnScan`:
     - Check if setting `autoGitignoreAgm` is enabled in DB or config.
     - Pass `isAutoYes = isAutoYes || isConfigAutoAgmEnabled`.

3. **Help & Examples**:
   - Update `gitmap gitignore agm help` to showcase all automated scan remediation workflows.
