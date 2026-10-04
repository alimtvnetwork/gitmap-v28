# Subtask 75.1: U1 SSH Node Discovery, Telemetry Probe & Local Caching

- **Parent Plan:** [75-gitmap-u1-ubuntu-agm-fleet-integration.md](../../pending/75-gitmap-u1-ubuntu-agm-fleet-integration.md)
- **Spec Reference:** [02-spec/21-app/205-gitmap-u1-ubuntu-agm-fleet-integration/01-architecture-spec.md](../../../../02-spec/21-app/205-gitmap-u1-ubuntu-agm-fleet-integration/01-architecture-spec.md)
- **Status:** Completed
- **Target Area:** `cli/cmdnode`, `cli/cmdssh`, node cache store

## Objective
Connect to Ubuntu machine `U1` via SSH, probe remote host environment metadata (OS name, version, architecture, GitMap version, default workspace `/home/a/git-work`), and persist the telemetry into the local GitMap cache on Windows.

## Implementation Details
1. Discover connection settings for node `U1` (host, port, user `a`, SSH keys/credentials).
2. Execute remote telemetry probe retrieving `/etc/os-release`, `uname -m`, and `gitmap version`.
3. Verify default workspace `/home/a/git-work`.
4. Store machine profile in local GitMap cache (`~/.gitmap/nodes_cache.json` / SQLite `installation.db`) to enable zero-latency offline OS lookups.
