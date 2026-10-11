# GitMap v6.544.0 Release Notes

## Summary of Changes
GitMap v6.544.0 cleans up unused test helpers in `cli/cmd/schemaregistry_assert_test.go` to satisfy the `Lint Baseline Diff` analyzer in CI/CD pipelines, alongside comprehensive updates for repo-cache manifest injection, two-tier subtree remediation rendering, antivirus mitigations, and native Go path query acceleration.

---

### 1. Test Suite Cleanup & Lint Baseline Diff Compliance
- **Purged Unused Test Assertions:**
  - Removed unused `assertSchemaKeysArray`, `assertSchemaKeysSlice`, and `readEveryObjectKeys` from `cli/cmd/schemaregistry_assert_test.go`.
  - Retained active helper `assertSchemaKeysFirstObject` (exercised by `cli/cmd/sshlistjson_contract_test.go`).
  - Satisfies `Lint Baseline Diff` with zero new analyzer findings.
- **RCA Documentation:**
  - Added `.ai-memory/cicd-issues/121-unused-schemaregistry-assert-helpers-rca.md` and indexed in `.ai-memory/cicd-issues/index.md`.

---

### 2. Verified Capabilities
- **Companion Repo-Cache JSON Injection:** Probing default work directory and split-DB for repo-cache and adding companion record to `gitmap.json` (`cli/cmdscan/scanoutput.go`).
- **Two-Tier Subtree Solutions:** Structured subtree formatting with native GitMap remediation commands (`cli/cmdpull/pull_efficient_render.go`).
- **Antivirus Hardening:** Mitigated false-positive heuristics in Winlogon registry writes and privacy tweaks (`cli/cmdos/os_autologin_windows.go`, `cli/cmdos/os_tweak_windows_privacy.go`).
- **Native Go Path Querying (<5ms):** High-speed lookups via `gitmap rc --path`, `gitmap rs --path`, and `gitmap path <repo>` (`cli/cmd/special_repos_cli.go`, `cli/cmd/path_cli.go`).
- **File Size Gate Governance:** Decomposed `cli/cmd/rootutility.go` and `cli/cmdpipeline/pipeline_db_ops.go` strictly under 300 lines (spec 243.4).

---

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.544.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.544.0/install.sh | sh
```
