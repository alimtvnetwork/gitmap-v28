# Program 274 WS3 — 20 Codebase Improvements (proposals ONLY, zero code)

> Wave 1, workstream 3. Seeded from the last-100-commit review
> (`/tmp/274-last100.txt`): recurring commit shapes are treated as evidence
> of a missing mechanism. Each item: what it is, why it matters, expected
> benefit. Numbered 1–20.
>
> HARD EXCLUSIONS honored (272 decision record + program brief): no pipeline
> DB retention policy, no staging-repo philosophy changes, no error-message
> overhaul, no release-tag management changes. Already-shipped items are not
> re-proposed (deprecation enforcement, coverage floor gate, `nc` codegen).

### 1. Auto-derived lint baseline cache

- WHAT: Replace the manually-bumped lint baseline cache version (v3 → v4
  reseeds in the commit log) with a baseline derived from a hash of the lint
  config plus tree state, so the cache invalidates itself exactly when the
  rules or the code change.
- WHY: "Bump lint baseline cache version" is a human-remembered ritual; a
  forgotten bump silently compares against a stale baseline and either hides
  new violations or fails clean code.
- BENEFIT: Deletes a whole commit category; baseline staleness becomes
  impossible instead of policed.

### 2. Shift the 300-line file gate left with split suggestions

- WHAT: Make the file-size cap fail locally (fastgate / pre-commit) with a
  concrete suggested split plan — proposed file boundaries by concern —
  instead of CI discovering the breach after push.
- WHY: 4+ commits in the last 100 are "decompose X for file size gate"
  (pipeline_db_ops, rootutility, completion/install, nodes tasks). The fix is
  always the same shape; only its discovery is late and expensive.
- BENEFIT: Turns a push→CI-fail→fix→push cycle into a single local decision.

### 3. Gosec G115-safe integer conversion helper plus one codemod

- WHAT: Add one house helper for integer conversions that satisfies gosec
  G115, then run a single codemod rewriting the recurring `int(x)` patterns
  in probe/telemetry code — instead of fixing file by file.
- WHY: Two separate commits in the last 100 eliminate G115 in cpu/hw/nodes
  probes; the pattern will keep reappearing wherever telemetry code is
  written because the safe idiom isn't the default idiom.
- BENEFIT: Kills the vulnerability class at the source; future G115 findings
  in this shape drop to zero.

### 4. Hermetic self-update end-to-end test

- WHAT: A test that spins a fake release server, runs the updater under a
  temp HOME, and asserts the full contract: binary swap lands the new
  version, launcher wrapper survives, symlinks are preserved, no deploy
  contamination, console-safe handoff works.
- WHY: 5+ update commits in the last 100 (compiler errors in update UI,
  binary swap, deploy contamination, writability, terminal suggestions). The
  update path is the least-tested most-feared path in the codebase.
- BENEFIT: Update regressions get caught before release instead of by users
  reporting broken installs.

### 5. Release asset-count and draft-visibility gate in CI

- WHAT: The release workflow fails unless the published asset count equals
  the build matrix expectation and the release is verifiably non-draft after
  publish. (Verification only — no change to how tags are managed.)
- WHY: v4.184.0 shipped with 5 of 19 assets and v6.102.0 sat as an invisible
  draft; both were discovered by humans noticing something missing, hours or
  days later.
- BENEFIT: Short or invisible releases become structurally impossible; the
  release-race incidents stop recurring.

### 6. Defender heuristic trigger regression scan in CI

- WHAT: Scan every release binary for the 7 known Windows Defender heuristic
  trigger patterns identified in 385aba38e, and fail the release before
  publish on a match. (Pairs with proposed `release winscan`.)
- WHY: `VirTool:MSIL/Subti.C` false positives are currently found by users'
  antivirus after release; each incident costs a hotfix release and erodes
  trust in the Windows builds.
- BENEFIT: A reintroduced trigger pattern is caught in CI minutes after the
  offending commit, not after it ships to users.

### 7. Generate the helpdoc example registry from dispatcher metadata

- WHAT: Extend the program-272 `nc` codegen approach to the 9
  hand-maintained `cli/helpdoc/examples*.go` files: the 3-part example format
  from program 273 becomes a build-time artifact, and a command cannot ship
  without a registered example.
- WHY: The example registry is already drifting by hand (273 added the format
  and seeded entries manually); every new command otherwise needs a human to
  remember three files.
- BENEFIT: `gitmap help example <cmd>` stays complete forever; "missing
  example" stops being a review finding.

### 8. Per-package delta inside the coverage floor gate

- WHAT: Extend the program-272 global 2% coverage gate: fail if any single
  package drops more than 5%, even when the global number still holds.
- WHY: A global average masks a single package's coverage collapsing —
  exactly the packages (split constants, markers, telemetry) that the commit
  log shows breaking most often.
- BENEFIT: Coverage rot gets pinpointed to the package that rotted, before
  it becomes the next CI-red fix commit.

### 9. Single source of truth for the command list across all consumers

- WHAT: Derive the command list from dispatcher registration for every
  consumer — the coverage tests and the `gitmap:cmd` top-level marker
  validation that today parse split constants files by hand.
- WHY: Commits show these hand parsers desyncing on every file split
  ("restore gitmap:cmd marker on split commands_extra block", "coverage test
  parses split constants files"). Program 272 fixed this for `nc` only.
- BENEFIT: File splits stop breaking command-list consumers; one registration
  feeds nc, help, coverage, and marker checks alike.

### 10. Flaky-test quarantine with auto-retry-once

- WHAT: CI retries a failed test once; tests that fail-then-pass go to a
  quarantine list (visible, triaged weekly) instead of failing the pipeline.
- WHY: The log shows "prevent false failures on in-progress runs" and heavy
  test stabilization work (tempdir cleanup); flaky reds currently block every
  release and train everyone to distrust CI.
- BENEFIT: Real regressions still fail fast; transient flakes stop vetoing
  releases and stop generating noise commits.

### 11. Coalesce automated sync commits

- WHAT: Batch the "automated repository sync and working tree commit"
  commits (6+ in the last 100) per program window instead of per trigger, so
  one sync commit covers one logical sync.
- WHY: The current per-trigger commits bury meaningful history under noise;
  `git log` for a program reads as sync-spam with the real work hidden
  between entries.
- BENEFIT: Quieter, more reviewable history; bisect and archaeology get
  faster without losing any sync information.

### 12. SQLite WAL mode plus busy_timeout on all split DBs

- WHAT: Enable WAL journaling and a sane busy_timeout on gitmap.db,
  installation.db, pipeline.db, and the per-repo tier — a config change, not
  a schema change.
- WHY: Agents and commands hit the three tiers concurrently (fleet pulls,
  pipeline polling, agent telemetry); "database is locked" flakes are the
  classic symptom and the commit log shows ongoing DB hardening work.
- BENEFIT: Concurrent readers/writers stop stepping on each other; a whole
  class of transient failures disappears. (Locking behavior only — not a
  retention policy change.)

### 13. Slowest-test report and per-test timeout budget

- WHAT: CI reports the 10 slowest tests every run and fails any new test
  exceeding a time budget; the suite is heavy and getting heavier
  (tempdir-cleanup stabilization commits are the symptom).
- WHY: Without a budget, each new test adds seconds nobody notices until the
  suite takes 20 minutes and developers avoid running it — which is when
  regressions start escaping.
- BENEFIT: The suite stays runnable locally, which is the precondition for
  developers actually running it before push.

### 14. Secret-scan in CI tuned to house secret shapes

- WHAT: A pre-commit/CI scan (e.g. gitleaks with house rules) keyed to the
  secret shapes the repo actually handles: GCM env values, RSA-encrypted
  credential blobs, app-password formats, PATs in remotes.
- WHY: The codebase manages real credentials (email app-passwords,
  RSA-OAEP stores, hardened GCM env); AGENTS.md already documents which
  files must never leak, but nothing mechanical enforces it.
- BENEFIT: A mistyped `git add` can never publish a credential; the standing
  no-plaintext rule gets a machine enforcer.

### 15. Windows-first smoke job in CI

- WHAT: A dedicated windows-latest smoke job covering exactly the paths that
  keep breaking on user machines: install, self-update, selfinstall, and the
  Defender-sensitive binary layout.
- WHY: Windows-only breakage keeps escaping (POSIX-only test skips in the
  log, the user's Mac/Windows install failures, the Defender heuristic
  saga) because CI's heart beats on Linux.
- BENEFIT: The platform the owner actually runs on gets first-class
  verification; "works on CI, broken on my machine" incidents shrink.

### 16. `--json` output contract tests for agent-consumed commands

- WHAT: Golden JSON schemas for the ~20 commands agents parse most
  (help-json exists as the pattern); CI fails on undeclared schema drift.
- WHY: Agents are now first-class consumers of gitmap output, and a renamed
  JSON field breaks every automation silently — discovered at 3am during a
  delegated night, not in review.
- BENEFIT: Output becomes a versioned contract; agent fleets stop breaking on
  cosmetic refactors.

### 17. Versioned llm-docs artifact published per release

- WHAT: Publish `llm-docs-vX.Y.Z.md` (from the existing `llm-docs` command)
  as a release asset, so agents pin to the exact command surface of the
  version they run.
- WHY: Agents operate from stale training data about gitmap's 500+ commands;
  the matrix generator exists but its output isn't versioned or attached to
  anything.
- BENEFIT: "Which commands exist in v6.556.2" becomes an answered question
  instead of a hallucination risk; onboarding new agents gets deterministic.

### 18. Release-notes preview as a PR comment on version-bump PRs

- WHAT: When a version-bump PR opens, CI posts the would-be release notes
  (changelog entries + expected asset list) as a PR comment for review
  before the tag exists.
- WHY: Release notes are currently assembled at publish time ("reuse release
  notes" was itself a fix commit); errors in notes are found after the
  release is public.
- BENEFIT: Notes get reviewed like code, while they're still cheap to fix.

### 19. Automated contrast check for the design-system themes

- WHAT: A CI check asserting minimum contrast ratios for status colors
  (red/green/amber, checkmarks, error banners) across all 4 program-273
  palettes including monochrome.
- WHY: 273 shipped four themes; a palette that looks fine in a screenshot
  can still render red-on-dark unreadable in a real terminal, and nobody
  re-checks after each theme tweak.
- BENEFIT: Theme work stays accessible by construction; "I can't read the
  error output in monochrome" never becomes a support thread.

### 20. Incremental mode for `gitmap pc`

- WHAT: mtime/checkpoint-based short-circuit in pending-commits so clean
  repos are skipped on repeat runs; full walk only on demand or when the
  checkpoint is stale.
- WHY: `pc` is the most-run fleet command in the daily ritual (and inside
  273's proposed `health` aggregation); walking every repo every time is the
  dominant cost of the morning check.
- BENEFIT: The daily fleet-health loop drops from minutes to seconds, which
  is what makes running it daily actually happen.
