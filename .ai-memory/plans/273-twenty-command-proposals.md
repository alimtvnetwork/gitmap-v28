# Program 273 WS4 — 20 New Command Proposals (proposals ONLY, zero code)

> Spec: `02-spec/21-app/gitmap-help-text-and-ui-enhancement/06-twenty-new-command-proposals.md`.
> Audit gap IDs: n/a — subtask-04 audit report not yet produced; proposals seed
> from spec-06 directions (noted per proposal).
> Naming note: spec-06 seeds a `pr list/view/merge` family, but `pr` is taken
> by the merge-commit replay engine, so the GitHub PR family is proposed under
> `mr` (the WS3 `merge write` PR-creation command). `theme` is a global flag
> only — no `theme` command verb exists today.

### P-01 — `gitmap mr list`

- PROBLEM: After `mr` creates a PR, there is no way to list open PRs without
  leaving the terminal for github.com; reviewing what is awaiting merge across
  626 repos means clicking through the web UI repo by repo.
- INVOCATION: `gitmap mr list --state open`
- EXPECTED OUTPUT: `#  repo            age   checks  reviews / alimtvnetwork/go-email-reader  2d  ✓ 1/1 / alimtvnetwork/story-telling-v1  5h  ✗ 0/1`
- ALIASES: `mrl`
- HOST/TARGET: github / PR surface (migration surface, seed dir 1)
- EFFORT: S · VALUE: H
- GAIN: Replaces ~2 min of web-UI clicking per repo with one sub-second table.
- DEPENDS ON: WS3 `mr` (merge write) implementation; `gitmap login` token auth.

### P-02 — `gitmap mr view <n>`

- PROBLEM: Triaging a PR today requires opening the browser for the diff,
  check runs, and review state; there is no terminal view of "is this PR
  actually mergeable right now".
- INVOCATION: `gitmap mr view 42`
- EXPECTED OUTPUT: `PR #42 "release - bump to v1.1.0" · checks: 3/3 ✓ · reviews: 1/1 ✓ · +180/-24 · mergeable: YES`
- ALIASES: `mrv`
- HOST/TARGET: github / PR surface (seed dir 1)
- EFFORT: S · VALUE: H
- GAIN: One command answers "can I merge this" in seconds instead of 3-4 page loads.
- DEPENDS ON: WS3 `mr`; GitHub API check-run/review endpoints.

### P-03 — `gitmap mr merge <n>`

- PROBLEM: Merging an approved PR still means a browser round-trip (or raw
  `gh pr merge`); there is no gitmap-native merge that respects the house
  merge strategy and records the merge in the release flow.
- INVOCATION: `gitmap mr merge 42 --strategy squash`
- EXPECTED OUTPUT: `merged #42 (squash) → abc1234 · branch fix/foo deleted · remote updated`
- ALIASES: `mrm`
- HOST/TARGET: github / PR surface (seed dir 1)
- EFFORT: M · VALUE: H
- GAIN: Closes the PR loop entirely inside gitmap; ~1-2 min saved per merge.
- DEPENDS ON: WS3 `mr`; P-02 mergeability check logic.

### P-04 — `gitmap mr checkout <n>`

- PROBLEM: Reviewing someone else's PR locally means hand-typing
  `git fetch origin pull/42/head:pr-42 && git checkout pr-42` — slow and
  error-prone when done repeatedly across the fleet.
- INVOCATION: `gitmap mr checkout 42`
- EXPECTED OUTPUT: `fetched PR #42 → branch pr/42-user-slug · checked out · 3 commits behind main`
- ALIASES: `mrco`
- HOST/TARGET: github / PR surface (seed dir 1)
- EFFORT: S · VALUE: M
- GAIN: Turns a 30-second fiddly fetch into one command; zero typos in refspecs.
- DEPENDS ON: WS3 `mr`.

### P-05 — `gitmap repo export`

- PROBLEM: Moving a repo (Lovable→GitHub, machine→machine) today is a manual
  copy of clone URL, branches, remotes, and gitmap registration; nothing
  portable is produced, so migrations drift and get redone.
- INVOCATION: `gitmap repo export story-telling-v1 --out ./story-telling-v1.gitmap.json`
- EXPECTED OUTPUT: `exported 1 repo · 2 remotes · 3 branches · written ./story-telling-v1.gitmap.json (4.2 KB)`
- ALIASES: `repox`
- HOST/TARGET: scanner / repo registry (migration surface, seed dir 2)
- EFFORT: M · VALUE: H
- GAIN: Makes repo migration a file you can version, diff, and hand to anyone.
- DEPENDS ON: `scan export` manifest format (reuse, don't reinvent).

### P-06 — `gitmap repo import`

- PROBLEM: `repo export` output has no consumer: re-registering a repo on a
  new machine still means clone + `scan` + manual remote fix-ups as separate
  steps that each fail differently.
- INVOCATION: `gitmap repo import ./story-telling-v1.gitmap.json --into ~/workspace/repos`
- EXPECTED OUTPUT: `cloned → ~/workspace/repos/story-telling-v1 · registered in gitmap.db · remotes restored (2) · scan updated`
- ALIASES: `repoi`
- HOST/TARGET: scanner / clone pipeline (seed dir 2)
- EFFORT: M · VALUE: H
- GAIN: Full export→import round-trip turns a 10-minute migration into one command.
- DEPENDS ON: P-05 manifest format; `clone` + `scan` internals.

### P-07 — `gitmap help topic`

- PROBLEM: `gitmap help` lists 392 commands but teaches no learning path; a
  new operator (or agent) cannot discover "the 8 commands for doing a release"
  without already knowing their names.
- INVOCATION: `gitmap help topic releases`
- EXPECTED OUTPUT: `TOPIC releases: 1) gitmap summary 2) gitmap release-dry 3) gitmap release-notes 4) gitmap release vX.Y.Z 5) gitmap pe -t …`
- ALIASES: `ht`
- HOST/TARGET: help / discoverability (seed dir 3)
- EFFORT: S · VALUE: H
- GAIN: Cuts command-discovery time from minutes of guessing to one lookup.
- DEPENDS ON: Program 273 help-registry + example registry (spec 03).

### P-08 — `gitmap howto <task>`

- PROBLEM: Users think in tasks ("ship a version", "fix a red pipeline"), not
  verbs; today that intent dies at the help screen and costs a support round.
- INVOCATION: `gitmap howto "ship a new version"`
- EXPECTED OUTPUT: `RECIPE: 1) gitmap summary 2) gitmap release-dry 3) gitmap release vX.Y.Z 4) gitmap pe -t — see: help topic releases`
- ALIASES: `howto`
- HOST/TARGET: help / discoverability (seed dir 3)
- EFFORT: S · VALUE: M
- GAIN: Task→command mapping in one shot; fewer "how do I…?" interruptions.
- DEPENDS ON: P-07 topic index; 273 suggestion-engine data.

### P-09 — `gitmap theme list`

- PROBLEM: Three palette modes exist (`bright|standard|monochrome`) but are
  only discoverable via `helpdoc/theme.md`; there is no command that shows
  what is available and which is active.
- INVOCATION: `gitmap theme list`
- EXPECTED OUTPUT: `● bright (active) · ○ standard · ○ monochrome — set via --theme or GITMAP_THEME`
- ALIASES: `themes`
- HOST/TARGET: terminal UI / theming (seed dir 4)
- EFFORT: S · VALUE: M
- GAIN: Discoverable theming; ends "how do I change the colors" questions.
- DEPENDS ON: 273 theme-integration artifacts (subtask 02).

### P-10 — `gitmap theme preview <name>`

- PROBLEM: Switching themes today is blind trial-and-error across real
  commands; you cannot see what `monochrome` looks like before committing to it.
- INVOCATION: `gitmap theme preview monochrome`
- EXPECTED OUTPUT: Renders a sample status table, checkmarks, and error banner in the `monochrome` palette without changing the active theme.
- ALIASES: `themep`
- HOST/TARGET: terminal UI / theming (seed dir 4)
- EFFORT: M · VALUE: M
- GAIN: Try-before-you-buy theming; no more theme roulette on live output.
- DEPENDS ON: P-09; 273 theme-integration artifacts.

### P-11 — `gitmap nodes jump <node>`

- PROBLEM: Reaching a fleet node means remembering its hostname, SSH alias,
  and workspace path; `ssh` subcommands exist but no one-shot "take me to
  that machine's repo dir".
- INVOCATION: `gitmap nodes jump build-ubuntu-02`
- EXPECTED OUTPUT: `connected → build-ubuntu-02 · cwd ~/workspace/repos/gitmap-v28 (interactive shell)`
- ALIASES: `nj`
- HOST/TARGET: fleet / ssh ergonomics (seed dir 5)
- EFFORT: S · VALUE: M
- GAIN: One command replaces hostname lookup + ssh + cd on every fleet hop.
- DEPENDS ON: Existing `ssh` connection pool + node registry.

### P-12 — `gitmap nodes version`

- PROBLEM: After every release, fleet nodes silently run different gitmap
  versions; there is no command that shows version skew, so "works on my
  machine" fleet bugs waste hours.
- INVOCATION: `gitmap nodes version`
- EXPECTED OUTPUT: `build-ubuntu-02  v6.553.0 ✓ · mac-mini  v6.550.0 ⚠ behind · win-dev  unreachable ✗`
- ALIASES: `nv`
- HOST/TARGET: fleet / release verification (seed dir 5)
- EFFORT: S · VALUE: H
- GAIN: Instant fleet-wide "are we all on the new version" — the release saga's missing check.
- DEPENDS ON: `nodes ping` reachability; two-phase handshake pattern.

### P-13 — `gitmap pe triage`

- PROBLEM: `pe all` shows failures but the next step is manual: open the repo,
  re-run `pe`, check the branch, check recent commits — a 5-minute ritual per
  red repo before the RCA can even start.
- INVOCATION: `gitmap pe triage`
- EXPECTED OUTPUT: `worst: go-email-reader (exit 1, job build, step go-test) · branch main · last commit 2h ago · 25-line trace attached · suggested: gitmap pe -t go-email-reader`
- ALIASES: `pet`
- HOST/TARGET: pipeline / diagnostics (seed dir 6)
- EFFORT: M · VALUE: H
- GAIN: Collapses the 5-minute per-failure triage ritual into one ranked bundle.
- DEPENDS ON: `pe`/`pe all` extractors; `summary` commit data.

### P-14 — `gitmap pe rerun <run>`

- PROBLEM: A flaky CI failure forces a browser trip to click "Re-run jobs";
  from the terminal there is no re-run, so transient reds cost context switches.
- INVOCATION: `gitmap pe rerun 11823445501`
- EXPECTED OUTPUT: `re-run requested for run 11823445501 (failed jobs only) · watch: gitmap pipeline-ai status -t <eta>`
- ALIASES: `perr`
- HOST/TARGET: pipeline / CI actions (seed dir 6)
- EFFORT: M · VALUE: M
- GAIN: Flaky-failure recovery without leaving the terminal (~1-2 min per incident).
- DEPENDS ON: GitHub Actions re-run API; `pipeline-ai status` waiter.

### P-15 — `gitmap release train`

- PROBLEM: Releasing N repos (AGM + gitmap dual-release saga) is a manual
  sequence — version bump, CI gate, tag, release notes — repeated per repo with
  ordering done in the operator's head.
- INVOCATION: `gitmap release train agm gitmap-v28 --order agm,gitmap-v28`
- EXPECTED OUTPUT: `[1/2] agm v4.200.1 ✓ 19 assets · [2/2] gitmap-v28 v6.553.0 ✓ 16 assets · notes aggregated`
- ALIASES: `rtrain`
- HOST/TARGET: release automation (cross-repo workflow)
- EFFORT: L · VALUE: H
- GAIN: Turns a multi-hour, error-prone dual-release night into one ordered, gated run.
- DEPENDS ON: `release` internals; `pe` CI gate; `summary` notes.

### P-16 — `gitmap errors last`

- PROBLEM: The no-silent-errors rule tracks failures in the error module, but
  there is no command to surface the LAST tracked error — the real message and
  its fallback are buried where nobody looks.
- INVOCATION: `gitmap errors last`
- EXPECTED OUTPUT: `14:02:11 push failed: remote rejected (hook declined) · fallback: fix available via gitmap fix-auth · module: push-recovery`
- ALIASES: `errl`
- HOST/TARGET: diagnostics / error module
- EFFORT: S · VALUE: H
- GAIN: Makes the standing "no silent errors" rule actually visible and actionable.
- DEPENDS ON: Error-module storage format (wherever failures are tracked).

### P-17 — `gitmap scan diff`

- PROBLEM: After a `scan`, nobody knows what changed — new clones, deleted
  dirs, moved repos are invisible until something breaks downstream.
- INVOCATION: `gitmap scan diff`
- EXPECTED OUTPUT: `+ 2 repos (~/work/new-svc, ~/work/exp) · − 1 (~/tmp/scratch gone) · ~ 1 moved (old-path → new-path)`
- ALIASES: `scand`
- HOST/TARGET: scanner / repo registry
- EFFORT: S · VALUE: M
- GAIN: Instant workspace-change awareness; catches accidental deletions early.
- DEPENDS ON: `scan` snapshot persistence (last-scan DB state).

### P-18 — `gitmap db doctor`

- PROBLEM: Three SQLite tiers + per-repo caches grow silently; corruption or
  bloat is only discovered when a command mysteriously slows down or fails.
- INVOCATION: `gitmap db doctor`
- EXPECTED OUTPUT: `gitmap.db 12 MB ✓ · installation.db 3 MB ✓ · pipeline.db 48 MB ⚠ stale>30d · repodb/ 210 MB ⚠ — run with --repair`
- ALIASES: `dbdoc`
- HOST/TARGET: database / split-db engine
- EFFORT: M · VALUE: M
- GAIN: Proactive cache health; prevents the "why is gitmap slow today" mystery.
- DEPENDS ON: Split-DB inventory (db names/locations registry).

### P-19 — `gitmap hooks install --all`

- PROBLEM: House rules (no colons in commit messages, lowercase file names)
  are enforced by memory and review, not by machinery — so `gitmap cpf`
  still receives colon-laden messages that the tool then has to work around.
- INVOCATION: `gitmap hooks install --all`
- EXPECTED OUTPUT: `installed commit-msg guard in 41/43 repos · 2 skipped (no .git) · --dry-run to preview`
- ALIASES: `hkin`
- HOST/TARGET: git operations / policy enforcement (bulk operation)
- EFFORT: M · VALUE: H
- GAIN: Shifts convention enforcement from human memory to git hooks; kills a whole class of bad commits fleet-wide.
- DEPENDS ON: Policy-hook template set (commit-msg colon guard, filename case check).

### P-20 — `gitmap health`

- PROBLEM: The morning "is everything OK" check is 4 commands (`pc`, `pe all`,
  `nodes version`, `db doctor`) stitched together by hand, every single day.
- INVOCATION: `gitmap health`
- EXPECTED OUTPUT: `repos: 43 (2 dirty) · pipelines: 1 failing · fleet: 1 node behind · dbs: 1 stale · worst item: pe triage → go-email-reader`
- ALIASES: `hlth`
- HOST/TARGET: diagnostics / daily dashboard (cross-subsystem)
- EFFORT: M · VALUE: H
- GAIN: One daily command replaces a 5-minute multi-command ritual and points at the worst problem first.
- DEPENDS ON: `pc`, `pe all`, P-12, P-18 (aggregation only, no new data sources).
