# New-Command Checklist

Mandatory steps for adding a `gitmap` subcommand. Updated for the spec-03
`HelpDisplay` technique (spec 243.3); the old hand-written help-builder
approach is retired.

## The spec-03 HelpDisplay technique (mandatory)

Every new command uses the `HelpDisplay` system from day one:

1. **Struct-first.** Model the command's help as data, not strings:
   `helpdisplay.HelpDisplay` → `CommandHelpGroup` → `CommandHelper`
   (command, description, example, url, subHelpers). The structs are the
   source of truth — never hand-format help output.
2. **Generated `helptext/*.md`.** The `helptext/<command-id>.md` topic is
   GENERATED from the structs (generator: `03-ai-scripts/` or a gitmap
   subcommand — see the spec-03 plan). Do not hand-write or hand-edit the
   generated file; change the struct and regenerate.
3. **Single-touch registration.** The command is registered in exactly one
   dispatch table (`root*.go` `dispatchEntry`). Help wiring
   (`--help` → build `HelpDisplay` → print) is automatic via the
   displayer; no per-command help plumbing.

Reference: `02-spec/21-app/243-cli-help-displayer-and-backup-branch/03-cli-help-displayer.md`.

## Checklist

- [ ] `Cmd<Name> = "<id>"` constant added in `cli/constants/constants_cli.go`
      (aliases as `Cmd<Name>Alias` / `Cmd<Name>Alias2` — aliases are a kept feature).
- [ ] Implementation lives in its owning `cli/cmd<name>/` package (small
      files; `cli/cmd/` keeps only the thin `dispatchEntry` line).
- [ ] Single `dispatchEntry` added to the correct `root*.go` table —
      one line, no extra wiring.
- [ ] Help modeled struct-first: `CommandHelper`/`CommandHelpGroup`
      populated for the command (description, example, subcommands).
- [ ] `helptext/<id>.md` GENERATED from the structs (never hand-written).
- [ ] `TestEveryCmdIDHasHelpFile` passes (CI enforces: every primary
      command ID must have its `helptext/<id>.md`).
- [ ] `--help` and `--help <sub>` render via the displayer (theme
      inheritance: item → group → display → theme → global default).
- [ ] File-size discipline: new files ~100-line median; split past ~300.
- [ ] `gofmt` clean; positive booleans; zero nesting.

## Migration note

Commands predating spec 243 keep their hand-built help until ported
(spec 03 §Migration: incremental, duplicated/outdated builders first).
New or rewritten commands MUST use the displayer from day one.
