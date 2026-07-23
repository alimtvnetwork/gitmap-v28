# Plan: Update-awareness commands + pending visibility + release

Scope is spec-only (this is the docs/spec site repo). Deliverables are new spec pages, updates to existing spec pages, a JSON schema entry per command, and a release bump. No Go code is written here; the Go implementation lands in the `gitmap` binary repo from these specs.

## 1. New commands to specify

### 1.1 `gitmap list --update` (alias: `lu`)

List every scanned repo that has an available upgrade (local tag < latest remote tag, or ahead-of-origin main). Columns: `repo`, `current`, `latest`, `behind`, `source` (release|import|git-tag), `path`.

Flags:

- `--json` structured output
- `--source release|import|git-tag` filter
- `--limit N`
- `--stale-days N` only repos whose last scan is older than N days (forces rescan first)

Exit 0 always; exit 2 if scan cache is missing.

### 1.2 `gitmap stats` additions (alias: `ss`)

Extend existing `spec/01-app/.../stats.md` output with an **Upgrades** block:

```
UPGRADES
  repos scanned:        142
  up-to-date:           118
  upgradable:            22   (gitmap list --update)
  unknown / no tags:      2
  last full scan:  2026-07-22 10:14
```

JSON adds `upgrades: { scanned, upToDate, upgradable, unknown, lastScanAt }`.

### 1.3 `gitmap hd` (help-dashboard, alias for `help --dashboard`)

One-screen dashboard showing: pending tasks count, upgradable repos count, last scan time, current gitmap version vs latest, and 5 most-recent completed tasks. Pure read command.

### 1.4 `gitmap update apply [repo...]` and `gitmap update all`

- `update apply <repo>` upgrades a single scanned repo to its latest release tag (git fetch + checkout tag; for source-linked repos, run their release/self-update path).
- `update all` iterates every upgradable repo from `list --update`.
- Both accept `-y` / `--yes` to skip the confirmation prompt.
- `--dry-run` prints the plan without touching anything.
- `--only release|import|git-tag` limits by source.
- On per-repo failure, enqueue a `PendingTask` with type `Upgrade` (new TaskType seed) and continue.

Exit codes: 0 all succeeded, 1 partial, 2 nothing upgradable.

### 1.5 `gitmap update list` (alias of `list --update`)

Same output, discoverable under the `update` verb.

## 2. Pending task surfacing

Update `gitmap/helptext/pending.md` and `stats.md` so the pending count is visible from three entry points: `pending`, `stats`, `hd`.

Add new TaskType seed `Upgrade` in `constants/constants_pending_task.go` spec section (spec/01-app pending-task doc) so retry via `do-pending` replays `gitmap update apply <repo>`.

## 3. Spec files to add / change

New:

- `spec/01-app/NN-list-update.md` — command signature, flags, JSON schema, examples
- `spec/01-app/NN-update-apply.md`
- `spec/01-app/NN-update-all.md`
- `spec/01-app/NN-hd-dashboard.md`
- `spec/08-json-schemas/list-update.schema.json`
- `spec/08-json-schemas/update-apply.schema.json`
- `spec/08-json-schemas/hd.schema.json`

Change:

- `spec/01-app/19-list-versions.md` — cross-link to `list --update`
- existing `stats` spec — add Upgrades block + JSON field
- `gitmap/helptext/stats.md`, `pending.md`, `do-pending.md` — mention Upgrade task type
- `constants/constants_pending_task.go` spec doc — seed `Upgrade`, add `SQLSeedTaskTypes` note
- `plan.md` and `.lovable/memory/reports/20260723-rejog-reliability.md` — mark the P0/P1 items this touches
- `src/types/helpJson.ts` — add the four new command entries so the docs site index renders them
- `src/pages/FlagReference.tsx` — surface the new flags

## 4. Release

- Bump `.gitmap/release/latest.json` to next minor (v5.32.0) with `branch: release/v5.32.0`.
- Add `.gitmap/release/v5.32.0.json` release note listing the new commands.
- Update `.github/coverage.floor` only if the spec-site tests change; otherwise leave.
- Site version-sync test (`src/test/version-sync.test.ts`) will pick up the bump; add fixtures if needed.

## 5. Acceptance criteria

1. `list --update` spec defines exact columns, flags, exit codes, JSON schema, and 2 worked examples (text + JSON).
2. `stats` spec shows the Upgrades block in both text and JSON, matching a fixture.
3. `hd` spec renders a single-screen dashboard with pending + upgrades + version + last-scan.
4. `update apply` / `update all` specs cover: dry-run, `-y`, per-repo failure -> PendingTask(Upgrade), partial-success exit code, source filter.
5. `pending` / `do-pending` specs replay `Upgrade` tasks correctly.
6. Site docs index (`helpJson.ts`) lists all four new commands and their aliases (`lu`, `hd`).
7. Release v5.32.0 file exists and `latest.json` points at it; version-sync test passes.

## 6. Next task selection

Pick one to implement first once you switch me to build mode:

- **A.** All of the above in one pass (spec + site index + release bump).
- **B.** Just the `list --update` + `stats` upgrades block + site index (smallest useful slice, no `update apply` yet).
- **C.** Spec only, no release bump this turn.
- **D.** Release-focused: skip `hd`, ship `list --update` + `update apply/all` + release.

Which slice do you want, A / B / C / D?  
  


Add the help text. When we do the scan, it should also show that also in the help text and also in the UI. Okay? Uh, you do all these things all together. Okay? Not one, but all together as, uh, detailed as possible. And also make sure that there are other commands, and commands needs to be in the help section. Okay? So make sure of that. Um, and also I think there are some small issues you can check and fix those. And these are non-negotiable, okay? Must apply