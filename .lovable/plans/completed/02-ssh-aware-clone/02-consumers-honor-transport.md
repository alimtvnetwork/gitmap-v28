# Subtask 02 — All consumers honor identified transport

**Parent:** 02-ssh-aware-clone
**Status:** pending
**Created:** 2026-06-07

## Work
- `gitmap/formatter/clonescript.go` `cloneURL`: pick URL matching `IdentifiedTransport`; fall back to the other only if empty, and log.
- `direct-clone.ps1` / `direct-clone-ssh.ps1` generators: keep both for convenience, but the primary `clone.ps1` and the terminal "command:" line MUST use identified transport per repo.
- `gitmap/probe/*` and any pull/safe-pull path: select URL via identified transport.
- Terminal report (`scanoutput.go`): per-repo block must show `transport:`, `https:`, `ssh:`, and `command:` in the identified transport.
- Tests: golden update for terminal block; clone script generation tests for mixed-transport manifest.
