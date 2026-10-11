# GitMap v6.543.0 Release Notes

## Summary of Changes
GitMap v6.543.0 purges legacy unused status formatting files from `cli/cmd/` to ensure clean baseline diff validation in CI/CD pipelines, alongside cumulative enhancements for repo-cache manifest injection, two-tier subtree solutions, antivirus mitigations, and Go native path acceleration.

---

### 1. Code Cleanup & Linter Baseline Compliance
- **Purged Unused Status Formatters:**
  - Removed `cli/cmd/statusformat.go`, `cli/cmd/statusmissing.go`, and `cli/cmd/statusprint.go` (legacy status rendering replaced by dedicated `cli/cmdstatus/` package).
  - Removed unreferenced `statusSummary` struct from `cli/cmd/dbhelpers.go`.
  - Satisfies `Lint Baseline Diff (unused)` gate with zero findings.
- **RCA Documentation:**
  - Added `.ai-memory/cicd-issues/120-unused-statusformat-and-statusprint-rca.md` and indexed in `.ai-memory/cicd-issues/index.md`.

---

### 2. Core Feature Highlights
- **Companion Repo-Cache JSON Injection:** Automated companion repo-cache discovery in default work directory and split-DB, synthesizing `model.ScanRecord` with notes `companion:repo-cache` in `gitmap.json` (`cli/cmdscan/scanoutput.go`).
- **Two-Tier Subtree Solutions:** Formatted pull and pending commit errors as structured two-tier subtrees with native GitMap remediation commands (`cli/cmdpull/pull_efficient_render.go`).
- **Antivirus & EDR Hardening:** Removed Winlogon plaintext password registry writes and service tampering commands (`cli/cmdos/os_autologin_windows.go`, `cli/cmdos/os_tweak_windows_privacy.go`).
- **Native Go Path Acceleration (<5ms):** Instant query commands `gitmap rc --path`, `gitmap rs --path`, and `gitmap path <repo>` replacing slow Python subprocess wrappers (`cli/cmd/special_repos_cli.go`, `cli/cmd/path_cli.go`).
- **File Size Gate Governance:** Decomposed `cli/cmd/rootutility.go` and `cli/cmdpipeline/pipeline_db_ops.go` to strictly adhere to the repository 300-line limit (spec 243.4).

---

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.543.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.543.0/install.sh | sh
```
