# Issue: Update Flow Spec Alignment

**Status**: ✅ Resolved

## Problem

Repeated mismatch between general specs (`spec/03-general/`) and app-specific specs (`spec/01-app/`) for the update flow, causing implementation drift — lock/retry loops, async/foreground oscillation, missing version checks.

## Root Cause

1. General and app-specific specs described different update sequences with divergent terminology (`temp script` vs `handoff copy`, `update --from-copy` vs `update-runner`)
2. Missing explicit prohibitions (e.g., never use `cmd.Start()` in `runUpdate()`)
3. Post-update validation requirements were not strict acceptance criteria
4. Specs did not mandate rename-first for PATH sync — only mentioned rename as a "fallback"

## Solution

### Canonical Two-Phase Flow

**Phase 1 — Handoff and foreground execution:**
1. `gitmap-v28 update` creates handoff copy in active binary directory (`gitmap-update-<pid>.exe`, fallback `%TEMP%`)
2. Launch handoff copy with hidden `update-runner` command using `cmd.Run()` (foreground/blocking)
3. Parent waits for worker to complete. Terminal stays attached. **Never async detach.**

**Phase 2 — Update pipeline and validation:**
1. Handoff copy resolves repo path
2. Run `run.ps1 -Update` (full pipeline: pull, build, deploy)
3. PATH sync uses **rename-first** in update mode, then copy-retry as fallback
4. Print executable-derived version comparison (before and after)
5. Run `gitmap-v28 changelog --latest` from the updated binary
6. Run `gitmap-v28 update-cleanup` to remove handoff and `.old` artifacts

## Acceptance Criteria

- Active PATH binary equals deployed binary version after update
- If versions differ, update exits with clear failure output
- Changelog output executed via updated binary
- Cleanup runs after successful update
- Zero lock-retry loops during normal update

## Prevention — Do Not Repeat

- Any update-flow change must update ALL of: `spec/03-general/02-powershell-build-deploy.md`, `spec/03-general/03-self-update-mechanism.md`, `spec/01-app/09-build-deploy.md`, `spec/01-app/02-cli-interface.md`, and `spec/02-app-issues/`
- Keep one source-of-truth sequence and mirror verbatim across specs
- Explicit prohibitions must be documented (e.g., "never use `cmd.Run()`", "never add `Read-Host`")
