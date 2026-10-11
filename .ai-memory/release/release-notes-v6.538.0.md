# GitMap v6.538.0 Release Notes

## Summary of Changes
GitMap v6.538.0 resolves a `go vet` redundant newline error in `cli/cmd/rootutility_update_args.go` following the modular decomposition of utility command handlers, ensuring full compliance across all platform build checks and race detectors.

---

### 1. Fix Details
- **Redundant Newline Elimination:**
  - `cli/cmd/rootutility_update_args.go`: Switched `fmt.Println(updateUsageHelpText)` to `fmt.Print(updateUsageHelpText)` to prevent redundant newline warnings flagged by `go vet` and CI race tests.
- **RCA Documentation:**
  - Added `.ai-memory/cicd-issues/116-redundant-newline-in-rootutility-update-args-rca.md` and indexed in `.ai-memory/cicd-issues/index.md`.

---

### 2. Verified Capabilities
- **Repo-Cache JSON Injection:** Automated companion repo-cache discovery and inclusion in `gitmap.json` (`cli/cmdscan/scanoutput.go`).
- **Two-Tier Subtree Solutions:** Structured remediation subtrees and native GitMap command alternatives (`cli/cmdpull/pull_efficient_render.go`).
- **Antivirus Hardening:** Mitigated false-positive heuristics in Winlogon registry writes and privacy tweaks (`cli/cmdos/os_autologin_windows.go`, `cli/cmdos/os_tweak_windows_privacy.go`).
- **Native Go Path Querying:** High-speed `<5ms` path lookups via `gitmap rc --path` and `gitmap path <repo>` (`cli/cmd/special_repos_cli.go`, `cli/cmd/path_cli.go`).
- **File Size Gate Adherence:** Modular decomposition of `cli/cmd/rootutility.go`, `cli/cmd/rootutility_update.go`, and `cli/cmd/rootutility_update_args.go`, all under the 300-line cap (spec 243.4).

---

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.538.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.538.0/install.sh | sh
```
