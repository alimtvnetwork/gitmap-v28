# Spec 316: Repository Management and Testing Instructions

> **Spec ID:** `SPEC-316`
> **Subsystem:** Cross-repo program — Antigravity-Manager, `cli/` (gitmap-v28), yt-manage, img-pdf-v2
> **Status:** In Progress
> **Date:** 2026-10-24

---

## 1. Overview & Objective

The user ordered one coordinated verification-and-release program across the four
repositories in the active working scope:

1. **Antigravity tools (Antigravity-Manager):** test that auto-switching works and that
   fast-forward follows the discussed logic, verify as much as possible, finalize CI/CD
   with zero failures, then release with a minor version bump and tags (the `gitmap
   release` command may be used to create the tags).
2. **gitmap (gitmap-v28):** review the recent commands that have been implemented,
   integrated, and created; improve them where warranted; resolve any CI/CD failure;
   verify the fast-forward follows the modern logic defined in the spec and recent
   committed/completed tasks; test thoroughly; then release with a minor version bump
   so the user can test it in the morning.
3. **img-pdf-v2:** clone via `gitmap cfr` into the workspace repos (done — local folder
   `img-pdf`, remote verified as `alimtvnetwork/img-pdf-v2`), create its dedicated chat
   (done), investigate why the right-click button does not work, and why the YouTube
   functionalities (e.g. YouTube thumbnail download) are not done yet.
4. **yt-manage / gitmap YT manage:** check whether it works.
5. **Rerun rule (applies to Antigravity tools):** after a fast-forward or a switch, the
   system must rerun the running prompt — a single one.

## 2. Scope Lock

Per the user's rule of 2026-10-24, active work is limited to exactly four repositories:
`gitmap-v28`, `Antigravity-Manager`, `yt-manage`, `img-pdf-v2`. No other repository is
touched unless the user names it. The full GitHub account inventory is never a clone
list (see the hard scope lock added to the master prompt in coding-guidelines-v24,
Step 3 item 5, commit `8a2cfca8`).

## 3. Track A — Antigravity-Manager

- Locate the auto-switch implementation and its trigger conditions; verify the
  selection logic matches the discussed behavior.
- Locate the fast-forward implementation; verify it follows the modern logic
  (fast-forward only when the branch can move cleanly; no silent merge/rebase of
  published history).
- Trace the post-switch / post-fast-forward path and verify exactly **one** running
  prompt is re-dispatched (restore/rerun semantics: unscoped restore is FIFO, explicit
  instance scope is strict; `agm rerun` takes the most recent N and actually re-sends —
  fixed in v4.190.0 and must not regress).
- This VM is headless: the Tauri GUI and the Antigravity IDE cannot run here. Verify via
  the `agm` CLI headless where possible, plus code-level verification and CI.
- CI/CD gate: no failing workflow on main at release time. Pre-existing failures must
  be distinguished from regressions introduced by this program; regressions are fixed,
  pre-existing failures are reported with evidence.
- Release: minor bump from v4.253.0 → **v4.254.0**, tag created (via `gitmap release`
  where it fits), release attributed to Antigravity-Manager on its own v4.x line.

## 4. Track B — gitmap-v28

- Inventory the commands added in the recent wave (specs 271–315 and the CLI surface)
  and review them for correctness and improvement candidates.
- Verify the fast-forward implementation in the pull engine against the governing
  specs: `02-spec/21-app/03-git-operations-and-pull/`,
  `02-spec/21-app/243-pull-forward-slash-paths-pre-pull-and-interactive-remediation-prompt/`,
  and `02-spec/21-app/270-gitmap-nodes-pa-fleet-pull-command/`.
- Known observation to fold into the review (2026-10-24): `gitmap cpf` in a clone with
  no git identity configured surfaces a stack trace on first failure instead of a clean
  actionable error; the underlying cause (missing `user.name`/`user.email`) only becomes
  clear on a second run.
- Resolve CI failures on main that are attributable to the recent command work.
- Release: minor bump from v6.599.0 → **v6.600.0**, tag created, release attributed to
  gitmap on its own v6.x line, so the user can test in the morning.

## 5. Track C — img-pdf-v2

- Repository: `alimtvnetwork/img-pdf-v2`, cloned via `gitmap cfr`; local folder is
  `img-pdf` (cfr strips the version suffix). TypeScript/Vite application.
- Right-click: locate the context-menu / right-click handling (registration, event
  wiring, platform gating) and produce a root cause with file-level evidence.
- YouTube thumbnail download: determine whether the feature is absent, stubbed, or
  broken; cite the relevant modules or their absence, and state what finishing it
  requires.
- Deliverable is a root-cause report first; code fixes follow once the cause is proven.

## 6. Track D — yt-manage

- Verify the yt-manage CLI: entry points, command surface, auth model (OAuth), and
  whether the shipped functionality works as documented.
- Check any gitmap-side YT integration the user referred to as "gitmap YT manage".
- Headless limits (no interactive OAuth in this VM) must be stated explicitly; verify
  everything that can be verified without user credentials.

## 7. Cross-Cutting Rules

- Pull before any work in each repo (`gitmap pull`, fallback `git pull`); trees clean
  before starting. (Satisfied for this program on 2026-10-24.)
- Codebase search exclusively via GitMap (`gitmap aum search`, `gitmap find`,
  `gitmap cat`, `gitmap ps`, `gitmap py`). Total ban on rg, ripgrep, grep, git grep,
  Select-String.
- Strictly relative Git paths in all committed material, including the release page
  and release notes.
- Subagents are read-only researchers: no git commands, no file writes (Worker Git
  Ban). The lead performs all git operations.
- Verify via CI; do not run test suites locally. Fast formatters/linters allowed.
- Atomic commits via `gitmap cpf` / `gitmap cpb` in hyphen format; commit and push
  without asking, within scope.
- Never delete repositories or files unless explicitly asked. No force-push, no
  rebasing published history.

## 8. Verification Gates & Definition of Done

| Track | Done when |
|---|---|
| A | Auto-switch, fast-forward, and single-prompt rerun traced and verified (headless CLI + code + CI); gaps fixed or reported with evidence; CI green on main; v4.254.0 released with tag |
| B | Recent-command review delivered with improvement list; fast-forward verified against spec; CI failures resolved or proven pre-existing; v6.600.0 released with tag |
| C | Right-click root cause proven with file evidence; YouTube thumbnail gap explained with evidence and a finishing plan |
| D | yt-manage functional verdict delivered: what works, what does not, what needs user credentials |
