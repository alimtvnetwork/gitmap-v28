---
name: Inject/Open idempotency
description: Per-tool stamps on Repo (LastInjectedDesktopAt/VSCodeAt) gate Desktop+VSCode side effects in `inject` and `open`; --force/-f bypasses. Schema v25.
type: feature
---
Schema v25 adds two TEXT DEFAULT '' columns to `Repo`:
- `LastInjectedDesktopAt` — stamped after `registerSingleDesktop` runs.
- `LastInjectedVSCodeAt`  — stamped after `openInVSCode` runs.

Both `gitmap inject` and `gitmap open` consult them per-tool: empty == do
the action; non-empty == skip with one-line notice ("already injected
(<ts>) — pass --force to re-register"). `--force` (`-f`) zeros both
gates and re-stamps to CURRENT_TIMESTAMP after the side effects run.

Helpers live in `gitmap/cmd/inject_idempotency.go` (parseInjectForceFlag,
loadInjectStamps, markInjected, shouldRunDesktop, shouldRunVSCode) and
`gitmap/store/inject_idempotency.go` (GetInjectTimestamps, MarkInjected
+ typed `constants.InjectKind`). All best-effort: DB hiccups degrade to
"never injected" so the user-visible Desktop/VS Code calls always run.

Bumped `SchemaVersionCurrent` 24 → 25; migration is a pair of additive
ALTERs in `migrateRepoInjectTimestamps()`.

Specs: `spec/04-generic-cli/31-open.md`,
`spec/02-app-issues/26-cfrp-no-version-suffix-hard-error.md`.
