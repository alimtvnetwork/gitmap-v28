# Program 274 WS3 — 10 New Command Proposals (proposals ONLY, zero code)

> Wave 1, workstream 3. Seeded from the last-100-commit review (`/tmp/274-last100.txt`).
> Duplicate-diligence: checked against `gitmap nc` catalog (105 commands),
> `cli/helpdoc/examples*.go` registry (9 files), the 273 P-01..P-20 proposals,
> and live `--help` for `update`, `release`, `ssh`, `prompts`, `agent`, `agm`,
> `fix`, `scan`, `ide`, `scaffold`, `selfinstall`, `update-ide`.
> Rejected ideas from the 272 decision record (pipeline DB retention,
> staging-repo philosophy, error-message overhaul, release-tag management
> changes) are not proposed here.
> Dropped during diligence (already exist): `update --dry-run`, `update all`
> (fleet), `ssh troubleshoot (doctor)`, `agm update-all`, `ide`, `scaffold`.

### C-01 — `gitmap install verify`

- PURPOSE: Post-install self-check that the gitmap installation on this
  machine is actually healthy.
- INVOCATION: `gitmap install verify`
- EXPECTED OUTPUT: `binary v6.556.2 ✓ · launcher wrapper intact ✓ · PATH registered ✓ · symlinks preserved (3/3) ✓ · no deploy contamination ✓ · install source: github release`
- PAIN IT KILLS: Install breakage is per-OS and diagnosed by hand every time:
  tar-extract-dir bugs, `run.sh` deploy contamination, lost launcher wrappers,
  broken symlinks, and the user's still-open "install fails on my Mac" report.
  Today the only signal is "it doesn't work"; this turns it into a checklist.
- NOT A DUPLICATE OF: `update --dry-run` (previews a download, never audits
  installed state); `selfinstall` (installer entry point, no verification).

### C-02 — `gitmap install urls <tag>`

- PURPOSE: Extract every download URL from the release's installer scripts
  (`install.ps1`, `install.sh`, Homebrew cask) and HEAD-check each one.
- INVOCATION: `gitmap install urls v6.556.2`
- EXPECTED OUTPUT: `12 urls checked · 12 × 200 ✓ (0.8s) — or: ✗ https://…/gitmap_6.556.2_macOS_arm64.dmg → 404 (cask stanza)`
- PAIN IT KILLS: The 2026-10-15 AGM install-fix wave was "all verified HTTP
  200" done by a human clicking links (wrong DMG filename, AppImage-vs-.deb
  mixups, stale cask URLs). Broken installer URLs are currently discovered by
  users, not by the release process.
- NOT A DUPLICATE OF: C-03 (that audits the GitHub release object; this audits
  the installer scripts' outbound links).

### C-03 — `gitmap release verify <repo> <tag>`

- PURPOSE: Post-release integrity audit: published asset count vs the build
  matrix expectation, draft visibility, and tag/version parity.
- INVOCATION: `gitmap release verify gitmap-v28 v6.556.2`
- EXPECTED OUTPUT: `v6.556.2: 16/16 assets ✓ · draft=false ✓ · tag→abc0f2faa = HEAD ✓ · checksums present ✓`
- PAIN IT KILLS: The release-wave races — v4.184.0 shipped with 5 of 19
  assets, v6.102.0 sat as an invisible draft, and each incident was found by a
  human noticing something missing. This is the "did the release actually land
  whole" check, runnable by anyone in seconds.
- NOT A DUPLICATE OF: 273 P-15 `release train` (orchestrates multi-repo
  releases; this audits one finished release); `release --dry-run` (previews
  before, never verifies after). Does not change tag management (272-excluded).

### C-04 — `gitmap release winscan <tag>`

- PURPOSE: Pre-publish scan of the Windows release assets against the 7 known
  Windows Defender heuristic trigger patterns identified in 385aba38e.
- INVOCATION: `gitmap release winscan v6.556.2`
- EXPECTED OUTPUT: `gitmap_6.556.2_windows_amd64.zip: 0/7 trigger patterns ✓ — or: ✗ pattern #4 (credential-string cluster) matched in gitmap.exe+0x1a2f30`
- PAIN IT KILLS: Every release cycle risks another `VirTool:MSIL/Subti.C`
  false positive on user machines, and remediation today is reactive
  (ship → flagged → hotfix → re-ship). This moves the check to before publish.
- NOT A DUPLICATE OF: `fix` subcommands (source-code fixers; this scans built
  release binaries); anything in the `nc` catalog (no AV/defender command
  exists).

### C-05 — `gitmap prompts drift`

- PURPOSE: Read-only drift report between the canonical prompt home
  (`coding-guidelines-v24` `01-prompts/27-muse-prompts/`) and every registered
  copy, by content hash.
- INVOCATION: `gitmap prompts drift`
- EXPECTED OUTPUT: `27-muse-prompts/execute-in-a-step.md: canonical a1b2… · gitmap-v28 copy a1b2… ✓ · agm copy 9f8e… ✗ DRIFT (12 lines)`
- PAIN IT KILLS: "Sync canonical prompts, skills, gitmap agent standards" is
  currently a manual chore commit; drift between the canonical home and copies
  is discovered by accident, weeks later, when an agent follows a stale prompt.
- NOT A DUPLICATE OF: `prompt list`, `prompt show --copy`, `prompt-template`
  (all operate on local prompt files; none compare canonical vs copies).

### C-06 — `gitmap agent report --since <duration>`

- PURPOSE: Shift report from the agent task SQLite DB: per-agent tasks
  claimed/completed, durations, and open blockers over the window.
- INVOCATION: `gitmap agent report --since 12h`
- EXPECTED OUTPUT: `worker-01: 4 claimed · 4 done (2h11m) · worker-02: 3 claimed · 2 done · 1 blocked (ci red, see pe triage) · crashed: 0`
- PAIN IT KILLS: The user delegates whole nights to agent fleets and the
  morning report is hand-assembled from coordinator messages; a dead or stuck
  worker is found by noticing silence. This makes fleet output queryable.
- NOT A DUPLICATE OF: `agent task`/`subtask`/`log` (lifecycle management, no
  rollup); `agent crashed`/`diagnose` (forensics on failures, no shift
  summary); 273 P-20 `health` (workstation/pipeline dashboard, not agent
  telemetry).

### C-07 — `gitmap ci open [repo]`

- PURPOSE: Open the repo's latest GitHub Actions run in Chrome, deep-linked
  to the failing job when the run is red.
- INVOCATION: `gitmap ci open go-email-reader`
- EXPECTED OUTPUT: `opening https://github.com/alimtvnetwork/go-email-reader/actions/runs/11823445501/job/… (failed: go-test) in Chrome`
- PAIN IT KILLS: ~10 of the last 100 commits are CI-red fix cycles, and every
  one starts with 2–3 minutes of browser navigation: Actions tab → find the
  run → find the failing job. The terminal already knows the answer (`pe`,
  `pipeline-ai`); this hands the exact deep link to the browser in one step.
- NOT A DUPLICATE OF: 273 P-14 `pe rerun` (re-runs jobs; this opens the run
  for log inspection); `pipeline-ai status` (waits on state, never opens a
  browser). Chrome angle per owner suggestion.

### C-08 — `gitmap scan verify`

- PURPOSE: Two-way index↔filesystem consistency check: DB entries pointing
  at deleted/moved dirs AND on-disk repos missing from the index.
- INVOCATION: `gitmap scan verify`
- EXPECTED OUTPUT: `index 43 repos · 41 on disk ✓ · 2 stale entries (~/tmp/scratch gone, ~/work/old renamed) · 0 unindexed · run with --repair to prune`
- PAIN IT KILLS: `orphans` covers only one direction (dirs not in the DB);
  the reverse — stale index entries for deleted directories — surfaces later
  as mysterious downstream failures in `pa`, `pc`, and `sync`. Nobody runs a
  consistency check because none exists.
- NOT A DUPLICATE OF: `orphans` (one direction only); 273 P-17 `scan diff`
  (what changed since last scan; this is index truth vs filesystem truth).

### C-09 — `gitmap spec progress --spec <path>`

- PURPOSE: Parse a program spec's parts and ledger checkboxes into a
  per-workstream completion overview.
- INVOCATION: `gitmap spec progress --spec 02-spec/21-app/gitmap-help-text-and-ui-enhancement`
- EXPECTED OUTPUT: `WS1 ████████ 8/8 done · WS2 ████░░░░ 4/8 · WS3 ░░░░░░░░ 0/8 · next: WS2 subtask 5 (helpdoc registry)`
- PAIN IT KILLS: Programs 271–274 each hand-write ledgers, and "mark waves
  complete" is a manual commit; program state is only visible by reading the
  whole ledger file. The lead's 100-commit log is full of ledger bookkeeping
  commits that a parser could summarize instantly.
- NOT A DUPLICATE OF: `spec next` (next sequence number); `spec-audit`
  (conformance/quality audit of spec content; this is completion state, not
  quality judgment).

### C-10 — `gitmap agm queue audit`

- PURPOSE: Cross-check Antigravity Manager's prompt queue against actual
  dispatch receipts and flag phantom dispatches (queued ≠ sent).
- INVOCATION: `gitmap agm queue audit`
- EXPECTED OUTPUT: `queue: 6 pending · 41 dispatched ✓ · 2 PHANTOM (marked dispatched, no receipt: prompt #118, #121)`
- PAIN IT KILLS: The genuine defect found 2026-10-16 — a failed queue
  dispatch still recorded as "dispatched" — plus the user's broken Send/Queue
  buttons. Today the only evidence is log lines the user has to paste by
  hand; this makes queue truth a command.
- NOT A DUPLICATE OF: `agm status`/`install`/`update`/`update-all`
  (lifecycle, no queue inspection); `rerun` (re-runs failed Antigravity
  prompts; this audits whether "dispatched" means sent).
