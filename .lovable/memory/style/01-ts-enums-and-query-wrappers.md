# TypeScript Enums and Query Wrappers

- **Date**: 2026-08-09
- **Topic**: Code Quality, TypeScript strictly-typed structures, Python/PHP unified error logging.

## The Problem

TypeScript code previously used magic strings and union string types like `"pass" | "fail" | "fallback"`.
Additionally, booleans were being explicitly inverted (`!response.isSuccess`) instead of using positive assertions (`response.isFail`).
Queries across multiple languages were doing scattershot logging which violated error management guidelines.

## The Solution

1. **Enums Only**: TypeScript union string types are strictly prohibited for control flow or statuses. You must use Enums.
2. **Suffix Requirement**: Every single Enum must end with the suffix `Type` (e.g., `StatusType`, `ResultType`).
3. **Explicit Booleans**: Always use explicit boolean state checks like `response.isFail`. NEVER use inverted success booleans (`!response.isSuccess`).
4. **Query Wrappers**: Go, Python, PHP, and TS must utilize a `queryWrapper` that automatically intercepts and logs query failures according to `.lovable/error-manage/guidelines.md`. The return object must contain explicit `isSuccess` and `isFailure` booleans.

## Enforcement

This is a strict code quality baseline. Any AI agent modifying or creating TS/Python/PHP code must ensure these rules are followed, else the PR will be rejected.
