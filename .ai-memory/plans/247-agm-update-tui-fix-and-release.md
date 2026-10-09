# Plan 247 — agm update TUI fix + release 6.517.0

## Problem
`gitmap agm update` (`cli/cmdinstall/agm_update.go`) has no proper UI: on Linux/macOS
`runAgmUpdateLinuxQuiet` swallows ALL installer output via `CombinedOutput()` (minutes
of silence, then one line or a full dump + stack trace on failure); on Windows/Unix
dispatch (`installagmanager.go:163`, `:173-175`) the installer inherits the terminal
raw (`cmd.Stdout = os.Stdout` …), so its escape sequences corrupt/crash the screen
with zero coordination. Root cause: no renderer owns the screen.

## Fix (per spec 01)
One mutex-guarded pterm-based renderer (GitLab-style: spinner + live log tail +
stage checkmarks) fed by CAPTURED (piped, never inherited) child output.
Structured failure panel instead of raw dump + stack trace. Plain-line fallback
when stdout is not a TTY. `cli/glyphs/` untouched (UI renders through the existing
pipeline, no byteFaithful exemption). Reuses in-repo precedents
(`cli/cloner/batchprogress.go`, `cli/cluster/pool.go`).

Test seam: env override for the installer source (default unchanged) so e2e can
point at a mock installer without network.

## File boxes (disjoint)
- W1 (implement): `cli/cmdinstall/agm_update.go`, `cli/cmdinstall/installagmanager.go`,
  + ONE new helper file/package under `cli/cmdinstall/`. NOTHING else.
- W2 (verify): `~/workspace/gitmap-e2e-247/` ONLY (mock installer + e2e script + logs).
- Lead: spec/plan/readme, build verification, release (version bump + tag + push).
- DO NOT TOUCH: `cli/cmdagent/agent_subtask.go` (246 workstream active), `cli/cmdupdate/`,
  `cli/cmddownload/`, `cli/glyphs/`, version files (lead-owned release step).

## Execution
1. W1 implements per spec; `go build ./...` (build ONLY — no `go test`).
2. Lead personally verifies the build (process rule).
3. W2 Phase A authors mock installer + `e2e-247-agm-update.sh`; Phase B runs them
   against the lead-verified binary + `gitmap aum search` code assertions (U02/U04).
4. Secrets gate + relative-path linter; targeted `git add` ONLY (never `git add -A` —
   246's modified files must not be swept in).
5. Release 6.517.0: bump (`37-bump-version.py` or manual), changelog entry, commit,
   annotated tag `v6.517.0` message exactly `Release v6.517.0`, push main + tag.
   Documented deviation: orchestrator test gates skipped per standing rule.

## Acceptance (spec 02)
U01 live progress during update; U02 zero stdio-inheritance in agm update paths;
U03 structured failure panel, no stack trace; U04 no concurrent-writer corruption;
U05 non-TTY plain fallback; U06 help/dispatch intact; U07 version.json=6.517.0,
tag v6.517.0 on origin, changelog entry, `gitmap version` = 6.517.0.
Failure policy: any FAIL → rejected, root-caused, re-verified; tag cut ONLY after
U01–U06 pass.

## Specs
- `02-spec/21-app/247-agm-update-tui-fix-and-release/01-rca-scope-and-plan.md`
- `02-spec/21-app/247-agm-update-tui-fix-and-release/02-acceptance-criteria-and-evidence.md`
