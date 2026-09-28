# Completed Plan: 184-agy-rp-running-projects-and-recreate-safety.md

## Status
Completed (Verified 100% Green, Tested, Rebuilt Binary)

## Spec Reference
[02-spec/21-app/89-agy-rp-running-projects-and-recreate-safety.md](../../../02-spec/21-app/89-agy-rp-running-projects-and-recreate-safety.md)

## Overview
Remapped `gitmap agy rp` strictly to `running-projects` (alias for `running-projects ls`), removed `rp` alias from `recreate-project` (retaining `recreate`, `rcp`, `rec`), restored top-level `gitmap rp` to `release-pending`, and established strict safety guards preventing `recreate-project` from converting system directories, user home profiles (`C:\Users\Administrator`), drive roots, or non-git folders into ad-hoc Antigravity projects.

## Target Subtasks
1. `01-agy-rp-running-projects-remap-and-alias-hygiene.md` — Remap `agy rp` to `running-projects`, sanitize recreate aliases to `rcp`/`rec`, restore top-level `rp` to `release-pending`, and update help documentation.
2. `02-recreate-restricted-and-non-git-safety-guards.md` — Implement restricted directory detection (`isRestrictedSystemOrHomeDir`), non-git directory guard (`isGitRepo`), update target resolution and recreate ops, and add comprehensive unit tests.

## Verification Checklist
- [x] `gitmap agy rp` executes `running-projects` listing active or queued prompts.
- [x] `gitmap agy recreate-project` / `recreate` / `rcp` / `rec` works for valid project targets.
- [x] Top-level `gitmap rp` invokes `release-pending`, not `recreate-project`.
- [x] Running `gitmap agy recreate` from `C:\Users\Administrator` aborts immediately with `E9101` without wiping or registering.
- [x] Running `gitmap agy recreate C:\` or other drive roots aborts immediately with `E9101`.
- [x] Running `gitmap agy recreate <non-git-folder>` aborts immediately with `E9102`.
- [x] Linter `check-nested-ifs.py --changed-only` passes with zero violations.
- [x] Unit tests in `cli/cmdagy/agy_recreate_test.go` pass.
