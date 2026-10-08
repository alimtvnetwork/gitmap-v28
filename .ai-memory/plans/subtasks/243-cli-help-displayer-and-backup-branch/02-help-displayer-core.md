# Subtask 02 — Help displayer core (`cli/helpdisplay/`)

## Objective
Build the new package per `02-spec/21-app/243-cli-help-displayer-and-backup-branch/03-cli-help-displayer.md`:
`Displayer` interface, `HelpDisplay` struct, `CommandHelper` (command/description/example/url/subHelpers),
`CommandHelpGroup` (header/commands/hints), `Suggestion` (text + variables),
`RenderContext` (cached values), `Theme` (header/command/description/hint colors + parent inheritance).

## Files (disjoint box)
- `cli/helpdisplay/` (new package — one small file per type)

## Done criteria
- [ ] Binding pattern verified against 2 sibling packages first (cite file:line in the ledger).
- [ ] `NewHelpDisplay(...) Displayer` constructor binds struct → interface.
- [ ] Colors reuse `cli/theme/` + `cli/glyphs/` — no new color engine.
- [ ] Inheritance chain works: item > group > display > theme > global default; explicit override wins.
- [ ] `Suggestion` variables interpolate from `RenderContext` at render time (demo with a cached IP value).
- [ ] `--help` auto-display hook point identified and documented (wiring comes in subtask 03).
- [ ] `go build` + `go vet` clean; every file small (median ~100 lines).
- [ ] Atomic commit via `gitmap cpf "helpdisplay - core displayer package"` + immediate push.
- [ ] Evidence: a small demo program output (before/after rendering sample) in the ledger.
