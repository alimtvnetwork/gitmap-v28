# Subtask 06: Build, Minor Bump v6.639.0, Release & Fleet Verification

- **Task Code:** `TASK-06`
- **Slug:** `gitmap-installation-suggestion-issue - Task 06`
- **Parent Plan:** [365-gitmap-installation-suggestion-issue.md](../pending/365-gitmap-installation-suggestion-issue.md)
- **Status:** `pending`

---

## Scope & Deliverables
1. Compile `gitmap.exe` with `go build -p 1 -o ../bin/gitmap.exe .` and deploy to `$env:LOCALAPPDATA\gitmap-cli\`.
2. Update `cli/changelog.md` and bump version to `v6.639.0`.
3. Verify terminal suggestion behavior across `install.ps1`, `local-install.ps1`, `gitmap install gitmap`, and `gitmap completion powershell`.
4. Commit release using `gitmap cpr "v6.639.0 persistent terminal suggestions in installer and cli"`.
5. Push to remote origin.
