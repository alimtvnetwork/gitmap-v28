# Subtask 07 — Windows VirTool:MSIL/Subti.C RCA document (program 273, WS5)

## File box

- ONE file: `.ai-memory/virus-rca-virtool-msil-subtic.md` (≤300 lines).
  Documentation ONLY — no code changes. (The fix already shipped in
  commit 385aba38e; this wave RECORDS the analysis.)

## Context

Spec: `02-spec/21-app/gitmap-help-text-and-ui-enhancement/07-windows-virus-rca.md`.
Ground truth: commit 385aba38e (2026-10-15), 13 files, +241/−230. Heuristic
triggers removed: `reg add HKCU\...\Run` self-persistence
(`cli/osutil/startup.go`, −54 lines), auto-login credential handling
(`cli/cmdos/os_autologin_windows.go`), VSS writer interop
(`cli/cmdchromeprofile/chromeprofile_vss_windows.go`, −78 lines), file-lock
probing (`cli/fspath/lockcheck_windows.go`, −70 lines). Root cause: Defender
scores behavior combinations, not intent. Prevention: EV code signing
(~USD 300–500/yr), Microsoft Defender vendor submission per release,
reproducible builds, CI AV gate, avoidance list (never re-add registry Run
writes, VSS interop, auto-login code, raw lock probing). Write for a
continuing AI: keep the open questions (engine version, per-asset signing
state, `gitmap doctor` AV check as WS4 candidate). Rules: gitmap search
tools only; relative paths.

## Checklist

- [ ] Triggers table: file → removed behavior → heuristic reason.
- [ ] Prevention section with costs and next steps.
- [ ] Avoidance list explicit (never re-add items).
- [ ] Open questions for the continuing AI.
- [ ] Commit: `gitmap cpf "rca - windows virtool msil subtic analysis"` + push
  (verify `git ls-remote`).
