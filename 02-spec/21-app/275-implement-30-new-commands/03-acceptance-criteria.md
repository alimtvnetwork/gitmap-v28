# 03 — Acceptance criteria

**Spec:** 275 · **Updated:** 2026-10-17

A command is DONE only when all of these hold:

1. **Registered** — invocable via its primary name AND every alias in the
   proposal (`gitmap <name>`, `gitmap <alias>`).
2. **`--help` renders** — description, usage, flags, example; no panic,
   no garbled output.
3. **Help-example entry** — `gitmap help example <cmd>` shows the 3-part
   format (what-it-does + invocation + expected output).
4. **Happy path works** — the command does what the proposal's EXPECTED
   OUTPUT describes (against real local state where read-only; dry-run
   where side effects exist).
5. **Clean errors** — missing required flag → usage + flag name, non-zero
   exit. Invalid flag → clean error, not a stack dump.
6. **No live side effects** — no real PRs, merges, releases, hook installs,
   or checkouts performed during implementation or E2E.
7. **No regressions** — `go build ./...` exit 0 (lead-verified); existing
   commands unaffected.
8. **File hygiene** — every new file ≤300 lines, lowercase names, grouped
   by concern; no binaries committed; relative paths only.

## E2E matrix (WS3)

For each of the 30 commands, record in `04-e2e-results.md`:
`<cmd> | --help | dry-run/happy-path | invalid-flag | PASS/FAIL | notes`

## Release gate (WS5)

- All 30 commands PASS the E2E matrix.
- `gitmap pe` 100% green on main.
- Then: minor bump, tag, push, `gh release create`, verify green on the
  release commit.
