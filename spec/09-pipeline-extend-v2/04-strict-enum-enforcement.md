# Strict Enum Enforcement

This codebase rejects string union types (also known as pipes) in TypeScript. All string sub-items must be rigorously extracted into strict `Enum` declarations to guarantee refactoring safety and eliminate magic strings.

## The Rule

Any instance of:
```typescript
type Status = "pass" | "fail" | "fallback";
```
MUST be converted into an Enum that explicitly ends with the suffix `Type`:
```typescript
export enum StatusType {
  Pass = "pass",
  Fail = "fail",
  Fallback = "fallback"
}
```

## UI and Spec Implications

This applies heavily to `src/components/ui/` (e.g., `chart.tsx`, `sidebar.tsx`, `carousel.tsx`) and `src/components/spec/`. Component properties that were previously evaluated against magic strings:
```typescript
if (props.indicator === "dashed") { ... }
```
MUST be rewritten to evaluate against the Enum specifically:
```typescript
if (props.indicator === ChartIndicatorType.Dashed) { ... }
```

No AI Agent may commit code containing string union types.
