# Subtask 01 — Canonical help-example format + data registry (program 273, WS2)

## File box

- `cli/helpdoc/examples.go` — NEW (≤300 lines): `HelpExample` / `Example`
  types + the command→examples registry.
- `cli/cmd/help.go` — surgical: `help example <cmd>` subcommand + append
  "Examples" section to `help <cmd>` from the registry.
- `03-ai-scripts/51-helptext-generator.py` — point at the registry so
  `cli/helpdoc/*.md` regenerates (never hand-edit generated md).
- Touch nothing else.

## Context

Spec: `02-spec/21-app/gitmap-help-text-and-ui-enhancement/03-help-example-format.md`.
Every command gets what-it-does + concrete invocation + expected output,
normalized placeholders (`<sha>`, `<date>`, `<machine>`, `<count>`, `<path>`).
Single source of truth: one Go registry; terminal help, `help example`, and
the docs UI all render from it. Rules: gitmap search tools only (TOTAL BAN on
grep/rg); relative paths; ≤300-line files; `go build ./...` from `cli/` to
verify — NEVER `go test`.

## Checklist

- [ ] Types: `Example{Invocation, ExpectedOutput string}`,
  `HelpExample{WhatItDoes string, Examples []Example}`.
- [ ] Registry: `map[string]HelpExample` keyed by verb; every existing verb
  AND alias present (reuse the subtask-04 audit if it lands first).
- [ ] `gitmap help example <cmd>` renders the canonical 3-part layout.
- [ ] `gitmap help <cmd>` appends the Examples section from the registry.
- [ ] Generator script updated; regenerate `cli/helpdoc/*.md`.
- [ ] Pipe + `--theme mono` render clean (no ANSI/control chars).
- [ ] Commit: `gitmap cpf "help - canonical help-example registry"` + push
  (verify `git ls-remote`).
