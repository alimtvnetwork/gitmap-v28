# GitMap v6.575.0

## What's Changed in v6.575.0

- **GoEmailSender Connection & Discovery**: Resolved cross-platform executable detection on Windows (`isExecutableFile` without POSIX permission mode flags, supporting `.exe`, `.cmd`, `.bat`), scanning across PATH, `GOPATH/bin`, and `~go/bin` for `go-email-reader` and `GoEmailSender`.
- **Multi-Login Profile Fallback**: Added `--profile` / `-p` support to `gitmap clone`, `reclone`, `cfr`, and `cfrp`. Probes primary/default profile first; on HTTP 401/404 or auth failure, executes concurrent background probing across secondary profiles in `02-git-profiles.json`, automatically binding projects to the winning identity.
- **Interactive Email Drafting & Scheduling Studio**: Added `gitmap space email ui` (and `gitmap space send-email --ui`) for terminal drafting, recipient and template review, account status display, and automated trigger registration.
- **Relative DB Path Reporting & Repo-Secrets URL**: Standardized `gitmap heatmap --json` and `gitmap agent task status --json` to report clean relative database paths (`.ai-memory/...` or traversal `../repo-cache/...`) along with `repoSecretsUrl`.
- **Fleet Nodes SSH Dispatch & Peer Delegation**: Added `gitmap nodes run [target] "<command>"` and peer GitMap delegation protocol (`PeerDelegationPayload`), routing `space`, `heatmap`, and `agent` subcommands cleanly over remote SSH hosts.
- **Isolated E2E Companion Tests**: Added end-to-end integration test suite guarded by `//go:build e2e` and `GITMAP_RUN_E2E=1`.

### Related Specifications & Plans
- Spec: `02-spec/21-app/goemailsender-connection-and-testing/01-architecture-spec.md`
- Plan: `.ai-memory/plans/goemailsender-connection-and-testing.md`

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.575.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.575.0/install.sh | sh
```
