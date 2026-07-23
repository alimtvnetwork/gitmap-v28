# Replace Command — Phased Plan

Spec: `spec/04-generic-cli/15-replace-command.md`.

## Phase 1 — Constants & dispatch (DONE this session)
- Add `CmdReplace`, `CmdReplaceAlias` to `constants_cli.go`.
- Add `replace` flag/message constants in new `constants_replace.go`.
- Wire `replace` into `roottooling.go` dispatch table.

## Phase 2 — Engine (DONE this session)
- `gitmap/cmd/replace.go` — entrypoint, mode detection, delegates.
- `gitmap/cmd/replaceflags.go` — flag parsing (`--yes`, `--dry-run`, `--quiet`).
- `gitmap/cmd/replacewalk.go` — repo walk, dir exclusions, binary detection.
- `gitmap/cmd/replaceapply.go` — atomic temp+rename rewrite, summary print.
- `gitmap/cmd/replaceversion.go` — remote URL parse → `(base, K)`.
- `gitmap/cmd/replaceaudit.go` — report-only scan, line-level output.

## Phase 3 — Bump version + memory (DONE this session)
- `Version` 3.95.0 → 3.96.0 (minor bump per project policy).
- Update `mem://index.md` with replace-command entry.

## Remaining / Follow-ups
- Generate Go unit tests (`replace_test.go`) covering: literal happy
  path, binary skip, `.gitmap/release-assets` exclusion, version-suffix
  detection, `-N` clamp, `all` on v1 no-op, audit-no-write.
- Add `gitmap/helptext/replace.md` simulation block (per Command Help
  System standard, 3-8 line realistic example).
- Add a marker-comment `// gitmap:cmd top-level` is already on the
  const block; CI generate-check should catch the new constant — verify
  on next CI pass.
- Consider integrating with `runReleaseSelf` to auto-suggest
  `gitmap replace -1` after version bumps.
