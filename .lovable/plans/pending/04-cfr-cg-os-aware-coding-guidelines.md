# cfr / cfrp `cg` modifier: OS-aware Coding Guidelines v24 integration

Slug: cfr-cg-os-aware-coding-guidelines
Steps: 10
Status: pending
Created: 2026-07-16

## Context

Extend `gitmap cfr` and `gitmap cfrp` with a `cg` (Coding Guidelines) modifier that, after clone + fix, runs the OS-appropriate coding-guidelines-v24 installer (PowerShell on Windows, bash on Linux/macOS), then stages, commits, and pushes the resulting files. Must not disturb the existing `install clean-code` (v15) alias.

- Captured command: `.lovable/spec/commands/06-cfr-cg-os-aware-coding-guidelines.md`
- Existing installer reference: `gitmap/cmd/installcleancode.go`, `gitmap/constants/constants_cleancode.go`
- Existing cfr entry: `gitmap/cmd/clonefixrepo.go`, `gitmap/cmd/clonefixrepo_escape.go`

## Steps

1. Add v24 installer constants (`DefaultCodingGuidelinesURLWindows`, `DefaultCodingGuidelinesURLUnix`, commit message, log messages) to a new `gitmap/constants/constants_codingguidelines.go`; do NOT modify the v15 `constants_cleancode.go`.
2. Introduce a shared runner `gitmap/cmd/codingguidelines.go` exposing `RunCodingGuidelinesInstall(opts)` that detects `runtime.GOOS`, dispatches to PowerShell (`powershell` then `pwsh`) or bash (`curl | bash`), streams stdio, and returns a typed error. See `./subtasks/04-cfr-cg-os-aware-coding-guidelines/SS-01-runner-design.md`.
3. Add `CfrModifierFlags` parsing helper (`gitmap/cmd/clonefixrepo_modifiers.go`) that recognizes positional tokens `cg`, `p`, in any order, before or after the URL; returns `{WantCodingGuidelines, WantPublic, RepoURL, PassThroughArgs}`.
4. Wire the modifier parser into both `cfr` and `cfrp` dispatch in `gitmap/cmd/clonefixrepo.go`; `cfrp` implies `WantPublic=true`, `cfr p` also sets it; `cg` sets `WantCodingGuidelines=true`.
5. After the existing cfr clone + escape-nested + fix sequence completes successfully, invoke `RunCodingGuidelinesInstall` inside the cloned working directory when `WantCodingGuidelines` is true.
6. Add auto commit + push stage: `git add -A`, commit `chore: install coding guidelines (v24)` (skip if working tree clean), push to tracked upstream only. Support `--no-commit` and `--no-push` flags surfaced through `CfrModifierFlags`.
7. Add help text: `gitmap/helptext/cfr.md` and `gitmap/helptext/cfrp.md` updated with a `## Coding Guidelines (cg)` section plus fenced `## Examples` block covering `cfr cg <url>`, `cfr p cg <url>`, `cfrp cg <url>`, and `--no-push` / `--no-commit`.
8. Register the new UI command entries in `src/data/commands.ts` (modifier documentation on existing cfr/cfrp entries, not a new top-level command) with examples matching the helptext.
9. Add unit tests: `gitmap/cmd/clonefixrepo_modifiers_test.go` covering token order permutations, and `gitmap/cmd/codingguidelines_test.go` covering OS dispatch selection via injectable `execRunner` seam (no network calls in tests). See `./subtasks/04-cfr-cg-os-aware-coding-guidelines/SS-02-test-matrix.md`.
10. Bump MINOR version across unified sites (Go constants, `version.json`, root README pins), add CHANGELOG entry under a new `## Added` block naming `cfr cg` / `cfrp cg`, and register the release JSON via existing release flow.

## Verification

- `go build ./...` and `go test ./gitmap/...` green (including new modifier + runner tests).
- Helptext golden test passes (both new sections include fenced `## Examples`).
- Manual smoke on Windows: `gitmap cfr cg <test-repo>` clones, runs PowerShell installer, commits, pushes.
- Manual smoke on Linux: same invocation runs bash installer path.
- `cmd_constants_test.go` still passes (no new top-level command IDs introduced).
- Version-sync test (`src/test/version-sync.test.ts`) green after MINOR bump.

## Appended from prior pending tasks

- `01-bulk-visibility-mapub-mapri.md` (still pending)
- `02-ssh-aware-clone.md` (still pending)
- `03-reclone-transport-and-vscode-open.md` (still pending)
