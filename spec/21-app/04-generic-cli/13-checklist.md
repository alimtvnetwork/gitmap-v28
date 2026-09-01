# Implementation Checklist

> **Related specs:**
> - [01-overview.md](01-overview.md) — design philosophy guiding each phase
> - [02-project-structure.md](02-project-structure.md) — scaffold phase package layout
> - [11-build-deploy.md](11-build-deploy.md) — build and deploy phase details

## Instructions for AI

This is a sequenced implementation plan. Execute each phase in order.
Reference the numbered spec files for detailed patterns.
All constraints from `08-code-style.md` apply to every file you write.
- After any Go refactor or file split, run `go test ./<affected-package>` immediately.
- Do not leave unused imports or stale symbols for a later cleanup pass.

---

## Phase 1: Scaffold (Do First)

- [ ] Run `go mod init <module-path>`
- [ ] Create `main.go` with minimal entry point calling `cmd.Run()`
- [ ] Create `constants/constants.go` with `Version`, tool name
- [ ] Create `constants/constants_cli.go` with all command names + aliases
- [ ] Create `constants/constants_messages.go` with error messages
- [ ] Create `constants/constants_terminal.go` with ANSI color codes
- [ ] Create `cmd/root.go` with `Run()` and `dispatch()`
- [ ] Create `cmd/rootflags.go` with flag registration helpers
- [ ] Create `cmd/rootusage.go` with `printUsage()` help text
- [ ] Implement `version` command (print version, exit 0)
- [ ] Implement `help` command (print usage, exit 0)

**Verify:** `go build && ./toolname version && ./toolname help`

---

## Phase 2: Configuration

- [ ] Create `model/` package with core data structs
- [ ] Create `config/config.go` with three-layer merge logic
- [ ] Create `data/config.json` with default settings
- [ ] Wire config loading into first real command

**Verify:** Tool reads config, flag overrides work

---

## Phase 3: Core Command

- [ ] Create `scanner/scanner.go` (or domain-specific logic package)
- [ ] Create `mapper/mapper.go` for data transformation
- [ ] Implement first real command (e.g., `scan`) with flag parsing

**Verify:** `./toolname scan <input>` produces correct data

---

## Phase 4: Output Formatting

- [ ] Create `formatter/terminal.go` with colored banner + item list
- [ ] Create `formatter/csv.go` with header row + CSV writer
- [ ] Create `formatter/json.go` with 2-space indented JSON writer
- [ ] Create `formatter/structure.go` with Markdown tree visualization
- [ ] Create `formatter/template.go` with shared template loading
- [ ] Create `formatter/templates/` directory with embedded `.tmpl` files
- [ ] Create output directory structure (`toolname-output/`)
- [ ] Add `--output` flag to main command
- [ ] Create date formatting utility function (see `14-date-formatting.md`)

**Verify:** Each format produces correct output; dates display consistently

---

## Phase 5: Database

- [ ] Create `store/store.go` with DB init, open, close, and migration
- [ ] Create `store/repo.go` (or domain CRUD file) with upsert logic
- [ ] Create `constants/constants_store.go` with SQL statements, DB paths, table names
- [ ] Wire DB upsert into main command's output flow
- [ ] Implement `db-reset` command

**Verify:** Data persists across runs, `db-reset --confirm` clears it

---

## Phase 6: Additional Commands

- [ ] Implement each remaining command in its own file
- [ ] Add flag parsing function per command
- [ ] Wire all commands into dispatch (split into multi-layer if 15+ cases)

**Verify:** Each command executes correctly with valid input

---

## Phase 7: Help System

- [ ] Create `helptext/print.go` with `go:embed` and `Print()` function
- [ ] Create one `.md` file per command (see `09-help-system.md` for format)
- [ ] Create `cmd/helpcheck.go` with `checkHelp()` function
- [ ] Add `checkHelp` call as the **first line** of every command handler
- [ ] Verify `--help` and `-h` both work on every command

**Verify:** `./toolname scan --help` prints help and exits 0; `./toolname cd -h` works

---

## Phase 8: Build & Deploy

- [ ] Create build script (`run.ps1` and/or `Makefile`)
- [ ] Add `-ldflags` for compile-time embedded variables (e.g., repo path)
- [ ] Add deploy step with nested directory structure
- [ ] Add retry-on-lock logic for Windows deploy
- [ ] Add version verification after build (run binary with `version`)
- [ ] For installer scripts, print a post-install summary with version, binary path, install directory, and PATH target/status
- [ ] Implement `update` command with copy-and-handoff self-update
- [ ] Implement `update-cleanup` command for artifact removal
- [ ] Implement post-install shell wrapper per [21-post-install-shell-activation.md](21-post-install-shell-activation.md):
  - [ ] Add `<TOOL>_WRAPPER` constant + `ShellWrapperMarkerPrefix`/`Suffix` constants
  - [ ] Create `setup/wrapper.go` with `DetectShell`, `ResolveProfilePath`, `InjectSnippet`, `RemoveSnippet`, `TryInSessionActivate`, `PrintReloadInstruction`
  - [ ] Wire wrapper injection into `setup` (and installer); attempt in-session activation, fall back to printed reload one-liner
  - [ ] Add `doctor` check returning `LOADED` / `INSTALLED_BUT_NOT_LOADED` / `NOT_INSTALLED`
  - [ ] Emit stderr warning from every shell-dependent subcommand when `<TOOL>_WRAPPER` is unset
  - [ ] Tests: fresh injection, re-injection (no duplicates), marker-based removal, all three `doctor` states

**Verify:** `./run.ps1` builds, deploys, and prints correct version; `./toolname update` works; `./toolname setup` activates the wrapper in the current shell (or prints the exact reload command); `./toolname doctor` reports `LOADED`

---

## Phase 9: Testing

- [ ] Add unit tests for `mapper` (table-driven, input/output pairs)
- [ ] Add unit tests for `config` (merge priority verification)
- [ ] Add unit tests for `formatter` (capture output via `io.Writer`)
- [ ] Add unit tests for `store` (in-memory SQLite)
- [ ] Add integration tests under `tests/` for command flag parsing
- [ ] Add wrapper tests: snippet injection (fresh), re-injection (idempotent), marker-based removal, `doctor` status detection
- [ ] Verify all tests pass: `go test ./...`

**Verify:** `go test ./...` — zero failures

---

## Phase 10: Polish

- [ ] Update `README.md` with grouped command reference + examples
- [ ] Verify all files ≤ 200 lines (split if exceeded)
- [ ] Verify all functions ≤ 15 lines (extract helpers if exceeded)
- [ ] Verify every edited Go file has zero unused imports after refactors
- [ ] Verify no magic strings (all in `constants`)
- [ ] Verify positive conditionals only (no `!`, no `!=`)
- [ ] Verify blank line before every `return`
- [ ] Verify no circular imports between packages
- [ ] Final version bump

**Verify:** Full `go build && go vet ./... && go test ./...` passes clean

---

## Quick Reference — File Counts

| Phase | Files Created |
|-------|--------------|
| Scaffold | ~8 |
| Configuration | ~3 |
| Core Command | ~3 |
| Formatting | ~7 |
| Database | ~4 |
| Commands | 1 per command |
| Help | 1 per command + `print.go` + `helpcheck.go` |
| Build | 1–2 scripts |
| Tests | 1 per testable package |

---

## Dependency Graph

```
Phase 1 (Scaffold)
  └─► Phase 2 (Config)
       └─► Phase 3 (Core Command)
            ├─► Phase 4 (Formatting)
            │    └─► Phase 5 (Database)
            │         └─► Phase 6 (More Commands)
            │              └─► Phase 7 (Help System)
            └──────────────────────► Phase 8 (Build & Deploy)
                                         └─► Phase 9 (Testing)
                                              └─► Phase 10 (Polish)
```

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
