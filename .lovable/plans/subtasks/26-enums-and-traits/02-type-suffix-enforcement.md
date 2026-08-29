# Subtask 26-02: Type Suffix (*Type) Enforcement

## Goal

Ensure all Enum type definitions end with `Type` (e.g., `ConflictModeType`, `ActionKindType`, `ValidationVerdictType`, `TaskStatusType`).

## Status: DONE

- Refactored `gitmap/cmd/commitin/enums.go` with explicit `*Type` type aliases.
