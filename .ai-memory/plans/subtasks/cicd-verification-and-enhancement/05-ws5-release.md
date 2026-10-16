# Subtask 05 — WS5: minor bump + release until green

**Version:** 1.0.0
**Updated:** 2026-10-17
**Owner:** Lead
**Status:** BLOCKED until WS2 reports 100% green

---

## Goal

Publish a release from a **minor version bump** and keep CI green on the
release commit — repeating bump+release until green. This is the program's
definition of done.

## Inputs

- WS2 completion report: "WS2 COMPLETE — CI 100% green on <sha>; WS5 unblocked."
- `01-overview.md` §5 — definition of done (normative).
- `version.json` — single source of truth (6.557.0 at program start).

## Gate (hard)

Do NOT start until the WS2 completion report lands with an explicit green
HEAD sha. Starting a bump or release on a red tree is a program failure.

## Steps

1. Verify the gate: re-run `gitmap pe` personally on the reported green sha
   (standing rule: never trust a self-reported "green").
2. Minor bump: update `version.json` via the standard bump flow (the owner's
   "minor bump" = a small increment from 6.557.0, not a major); the bump
   script is the sole writer for propagated copies — never hand-edit version
   strings elsewhere.
3. Commit the bump atomically and push.
4. Run the release ceremony (release workflow / release command per repo
   practice); record tag timing and verify the published asset count.
5. Watch CI on the release commit. If red: diagnose (WS2 loop rules apply),
   fix, bump again, re-release. Repeat until CI is green ON the release
   commit.
6. NEVER cut a new bump/release while a Release run is in flight — wait for
   the in-flight run to finish first.

## Acceptance criteria

- [ ] WS2's green report received and personally re-verified before starting.
- [ ] A release is published from a minor bump of `version.json`.
- [ ] CI is green on the release commit itself (not just an earlier commit).
- [ ] Release evidence recorded: tag, timing, version the owner runs,
      verified asset count.
- [ ] No release was cut while another Release run was in flight.

## Out of scope

New features; folding in the standing open items (Send/Queue buttons, PR #6,
release-pipeline durable fix) without owner approval.
