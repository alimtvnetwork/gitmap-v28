# Error Store (Zustand)

> **Parent:** [React Components Index](./00-overview.md)  
> **Version:** 4.0.0  
> **Updated:** 2026-04-01

---

## Review Compliance

| Rule | Status | Notes |
|------|--------|-------|
| No `unknown` | ✅ Fixed | `captureException` param typed via union |
| Parameters ≤ 3 | ✅ Clean | Max 2 per method |
| Typed error context | ✅ Fixed | Uses `ErrorContext` interface |

---

## Store Interface

```typescript
import type { CapturedError, ErrorContext } from './ErrorModalTypes';

/** Metadata for error capture enrichment */
interface CaptureErrorMeta {
  source?: string;
  triggerComponent?: string;
  triggerAction?: string;
}

interface ErrorStore {
  // --- State ---
  selectedError: CapturedError | null;
  isModalOpen: boolean;
  recentErrors: CapturedError[];
  errorQueue: CapturedError[];
  currentQueueIndex: number;
  pendingSync: Set<string>;

  // --- Actions ---
  captureError: (error: ApiError, meta?: CaptureErrorMeta) => CapturedError;
  captureException: (error: Error | string, context?: ErrorContext) => CapturedError;
  openErrorModal: (error: CapturedError) => void;
  openErrorQueue: (errors: CapturedError[], startIndex?: number) => void;
  navigateQueue: (direction: 'prev' | 'next') => void;
  closeErrorModal: () => void;
  clearRecentErrors: () => void;
  markErrorSynced: (errorId: string) => void;
  getPendingSyncErrors: () => CapturedError[];
  getQueuedErrorsMarkdown: () => string;
}
```

### Violations Fixed (v3.2.0 → v4.0.0)

| Previous | Violation | Fix |
|----------|-----------|-----|
| `captureException: (error: unknown, ...)` | Banned `unknown` param | → `Error \| string` union type |
| `context?: ErrorContext` (was unnamed) | Relied on external import without showing type | → Explicit `ErrorContext` reference |

---

## Key Behaviors

- **`captureError`**: Converts an API error into `CapturedError`, auto-capturing UI click path, execution logs, route info, and envelope diagnostic fields.
- **`captureException`**: Converts `Error` or string into `CapturedError`. If the value is neither, callers must convert to `Error` first.
- **`openErrorQueue`**: Opens the modal with multiple errors and queue navigation.
- **`navigateQueue`**: Cycles through queued errors (wraps around).
- **Envelope extraction**: Automatically extracts `requestedAt`, `requestDelegatedAt`, `envelopeErrors`, and `envelopeMethodsStack` from context.

## Stack Trace Parser

```typescript
export function parseFullStackTrace(stack: string): {
  frames: StackFrame[];
  primaryFrame: StackFrame | null;
  invocationChain: string[];
  rawStack: string;
} {
  // Handles patterns:
  // "at functionName (file:line:col)"
  // "at file:line:col" (anonymous)
  // "at async functionName (file:line:col)"
  // Webpack/Vite URLs: "at fn (http://localhost:5173/src/file.tsx:123:45)"
  
  // Internal frame detection:
  // - node_modules, chrome-extension://, @tanstack, react-dom, react., scheduler.
  // - Object.*, Module, <anonymous>
  
  // Returns first non-internal frame as primaryFrame
  // Limits invocation chain to 8 levels
}
```

---

*Error store — updated: 2026-04-01*
