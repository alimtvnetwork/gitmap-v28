# Subtask 04: Enum Type Suffix & Semantic Naming Compliance

> **Parent Plan:** `17-boolean-and-naming-audit.md`  
> **Scope:** TypeScript, PHP, and Go enum and type declarations

## Objective

Audit and enforce that all enums across the codebase end in the `Type` suffix (e.g. `UserRoleType`, `SyncStatusType`).

## Action Steps

1. Verify all TS enum definitions end with `Type`.
2. Verify all Go typed constants/enums follow semantic naming.
3. Ensure no magic string comparisons against enum fields.
