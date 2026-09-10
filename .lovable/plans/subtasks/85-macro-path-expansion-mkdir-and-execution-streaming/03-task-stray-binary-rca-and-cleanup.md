# Subtask 03: Stray Binary Root Cause Analysis, Elimination & Build Target Guard

## 1. Objective
Author a detailed Root Cause Analysis in `.lovable/issues/08-stray-parent-binary-rca.md`, remove the stray binary `D:\wp-work\riseup-asia\gitmap.exe`, align `Makefile` build targets to avoid root/parent pollution, and update `.lovable/strictly-avoid.md`.

## 2. Target Files
- `.lovable/issues/08-stray-parent-binary-rca.md`
- `.lovable/strictly-avoid.md`
- `Makefile`
- Parent directory executable cleanup (`D:\wp-work\riseup-asia\gitmap.exe`)

## 3. Requirements
- Document exact reasons why `D:\wp-work\riseup-asia\gitmap.exe` was created:
  - Prescribed `cd gitmap && go build -o ../gitmap .` build command executing from repo root.
  - Workspace parent scanning defaults in `run.ps1`.
  - Misinterpreted "root `gitmap.exe`" instruction.
- Safely delete `D:\wp-work\riseup-asia\gitmap.exe`.
- Update `Makefile` target `build` to write to `bin/gitmap.exe`.
- Clarify `.lovable/strictly-avoid.md` with an explicit Total Ban on binaries in parent workspaces or outside `bin/` and user AppData directories.
