# 4-Part RCA 74: Go Interface Naming and Enum Type Suffix Compliance in cmdos

- Date: 2026-09-22
- Quality Gates: Enum Guidelines Linter & Interface Naming Check
- Affected Files:
  - `cli/cmdos/os_theme_types.go`
  - `cli/cmdos/os_autologin_types.go`
  - `cli/cmdos/os_clean_sys_types.go`
  - `cli/cmdos/os_dns_types.go`
  - `cli/cmdos/os_tweak_types.go`

---

## 1. Why It Happened
During rapid implementation of the OS subsystems (`os autologin`, `os clean sys`, `os dns`, `os theme`, `os tweak`), internal orchestrator interfaces were named using the `Engine` suffix (`AutoLoginEngine`, `SystemCleanEngine`, `DNSEngine`, `ThemeEngine`, `TweakEngine`) and the desktop theme enum was named `ThemeMode`.

## 2. How It Happened
The linter quality gates:
1. `linter-scripts/check-enum-guidelines.py` enforces that all Go enum type definitions end with the mandatory `Type` suffix (e.g. `ThemeModeType`).
2. `linter-scripts/check-interface-naming.py` strictly enforces standard Go interface naming where interface definitions must end in `er` or `or` (e.g., `Reader`, `Writer`, `Operator`). Names ending in `Engine` violate this naming standard.

## 3. Root Cause
Lack of compliance with repository-wide interface naming standards (`*er` / `*or`) and enum type naming standards (`*Type`) in newly authored `cmdos` type definitions.

## 4. Code Fix
1. Rename enum `ThemeMode` to `ThemeModeType` in `cli/cmdos/os_theme_types.go`, `os_theme_windows.go`, `os_theme_linux.go`, `os_theme_other.go`, and `os_theme_cmd.go`.
2. Rename interface `AutoLoginEngine` to `AutoLoginOperator` in `os_autologin_types.go` and engine references.
3. Rename interface `SystemCleanEngine` to `SystemCleanOperator` in `os_clean_sys_types.go` and references.
4. Rename interface `DNSEngine` to `DNSOperator` in `os_dns_types.go`, `os_dns_windows.go`, `os_dns_linux.go`, `os_dns_other.go`, and `os_dns_cmd.go`.
5. Rename interface `ThemeEngine` to `ThemeOperator` in `os_theme_types.go`, `os_theme_windows.go`, `os_theme_linux.go`, `os_theme_other.go`, and `os_theme_cmd.go`.
6. Rename interface `TweakEngine` to `TweakOperator` in `os_tweak_types.go`, `os_tweak_windows.go`, `os_tweak_other.go`, `os_tweak_cmd.go`, and `os_tweak_mock_test.go`.
7. Verify with `linter-scripts/check-enum-guidelines.py` and `linter-scripts/check-interface-naming.py`.
