# Subtask 03 — Help-example data for every command (program 273, WS2)

## File box

- `cli/helpdoc/examples.go` — ADD entries only (split into
  `examples_<domain>.go` files if the registry nears 300 lines).
- No renderer changes (subtask 01 owns those). No dispatch changes.

## Context

Spec: `02-spec/21-app/gitmap-help-text-and-ui-enhancement/03-help-example-format.md`.
Subtask 01 built the registry + renderers; this subtask fills the data: one
`HelpExample` per command verb AND per alias, each with what-it-does +
runnable invocation + normalized expected output (≤25 lines). Domains to
cover: scan, clone, pull, commit suite (`c`/`cm`/`cpf`/`cpb`/`cpr`/`pcp`),
nodes, ssh, agm, update, prompt, llm train, task, schedule, fix, db,
pipeline (`pe`/`pea`/`te`), tree, summary, merge (after subtask 05 lands —
coordinate), plus UI help screens. Rules: gitmap search tools only; relative
paths; every data file ≤300 lines; NEVER `go test`.

## Checklist

- [ ] Coverage: registry keys ⊇ dispatcher verbs + aliases (diff vs the
  dispatch tables; flag gaps to subtask 04).
- [ ] Every invocation verified runnable as shown (run it, normalize
  unstable values to `<sha>`/`<date>`/`<machine>`/`<count>`/`<path>`).
- [ ] Real-looking arguments (no `<foo>` placeholders in invocations).
- [ ] UI help screens (docs viewer) show the same examples (single source).
- [ ] Commit: `gitmap cpf "help - help-example data for all commands"` + push
  (verify `git ls-remote`).
