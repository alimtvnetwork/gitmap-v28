# Subtask 75.5: Antigravity Manager (AGM) Linux Account Switching RCA & Migration

- **Parent Plan:** [75-gitmap-u1-ubuntu-agm-fleet-integration.md](../../pending/75-gitmap-u1-ubuntu-agm-fleet-integration.md)
- **Spec Reference:** [02-spec/21-app/205-gitmap-u1-ubuntu-agm-fleet-integration/01-architecture-spec.md](../../../../02-spec/21-app/205-gitmap-u1-ubuntu-agm-fleet-integration/01-architecture-spec.md)
- **Status:** Completed
- **Target Area:** `Antigravity-Manager` on `U1`, AGM Settings Exporter/Importer

## Objective
Diagnose why Antigravity Manager cannot switch accounts on Ubuntu Linux, export settings from Windows, transfer them via SCP to `U1`, and establish account switching functionality.

## Implementation Details
1. Inspect AGM codebase at `/home/a/git-work/Antigravity-Manager`.
2. Analyze root cause for Linux account switching failure:
   - Linux Secret Service / DBus keyring dependencies in headless or VM sessions.
   - Profile locking files (`SingletonLock`) under `~/.config/antigravity`.
   - Windows path backslash assumptions vs Linux forward-slash XDG paths.
3. Export Windows AGM profiles/tokens into a secure bundle.
4. SCP archive to `U1`, unpack into `~/.config/antigravity`, and verify AGM CLI functionality.
5. Provide a headless keyring fallback or environment-based token switcher.
