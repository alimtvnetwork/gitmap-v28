# CI/CD Task: Fix Go Interface Naming and Enum Type Suffix Compliance

## Source
- Runner job: `Enum Guidelines Linter` and `Interface Naming Check`
- Error type: FAIL
- Detected at: 2026-09-22T11:20:15+08:00

## Error Summary
- `cli/cmdos/os_theme_types.go:4`: Go enum `ThemeMode` missing mandatory `Type` suffix.
- `cli/cmdos/os_autologin_types.go:21`: Interface `AutoLoginEngine` must have 'er' or 'or' as suffix.
- `cli/cmdos/os_clean_sys_types.go:4`: Interface `SystemCleanEngine` must have 'er' or 'or' as suffix.
- `cli/cmdos/os_dns_types.go:43`: Interface `DNSEngine` must have 'er' or 'or' as suffix.
- `cli/cmdos/os_theme_types.go:12`: Interface `ThemeEngine` must have 'er' or 'or' as suffix.
- `cli/cmdos/os_tweak_types.go:28`: Interface `TweakEngine` must have 'er' or 'or' as suffix.

## Required Fix
Rename enum `ThemeMode` to `ThemeModeType` and rename all five interfaces from `*Engine` to `*Operator` or `*Manager` to satisfy mandatory suffix conventions.

## Acceptance Criteria
- [ ] `check-enum-guidelines.py` reports 0 violations
- [ ] `check-interface-naming.py` reports 0 violations
- [ ] `03-ai-scripts/06-cicd-local-runner.py` reports ✅ PASS for all quality gates
- [ ] No regression in unit tests or cross-compilation

## Status
- [x] resolved
