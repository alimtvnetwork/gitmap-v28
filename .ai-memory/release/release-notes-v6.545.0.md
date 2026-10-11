# GitMap v6.545.0 Release Notes

## Summary of Changes
GitMap v6.545.0 enforces canonical `gofmt` compliance in `cli/cmd/schemaregistry_assert_test.go` to ensure clean execution of CI lint-script unit tests, alongside completed implementations for companion repo-cache discovery, two-tier subtree remediation rendering, antivirus mitigations, and native Go path query acceleration.

---

### 1. Formatting & Lint Script Test Compliance
- **Gofmt Normalization:**
  - Formatted `cli/cmd/schemaregistry_assert_test.go` with canonical `gofmt -w`.
  - Verified with `python .github/scripts/tests/test_ci_scripts.py TestGoFormatCheck` and `python .github/scripts/go-format-check.py --check-only` (4499 Go files scanned, 0 unformatted).
- **RCA Documentation:**
  - Added `.ai-memory/cicd-issues/122-unformatted-gofmt-in-schemaregistry-assert-test-rca.md` and indexed in `.ai-memory/cicd-issues/index.md`.

---

### 2. Core Feature Highlights
- **Companion Repo-Cache JSON Injection:** Automated repo-cache detection in default work directory and split-DB, synthesizing companion records with notes `companion:repo-cache` into `gitmap.json` (`cli/cmdscan/scanoutput.go`).
- **Two-Tier Subtree Solutions:** Formatted remediation solutions into structured subtrees with native GitMap remediation commands (`cli/cmdpull/pull_efficient_render.go`).
- **Antivirus & EDR Hardening:** Removed Winlogon plaintext password registry writes and telemetry service tampering (`cli/cmdos/os_autologin_windows.go`, `cli/cmdos/os_tweak_windows_privacy.go`).
- **Native Go Path Querying (<5ms):** Instant query commands `gitmap rc --path`, `gitmap rs --path`, and `gitmap path <repo>` replacing slow Python subprocess wrappers (`cli/cmd/special_repos_cli.go`, `cli/cmd/path_cli.go`).
- **File Size Gate Governance:** Decomposed `cli/cmd/rootutility.go` and `cli/cmdpipeline/pipeline_db_ops.go` strictly under 300 lines (spec 243.4).

---

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.545.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.545.0/install.sh | sh
```
