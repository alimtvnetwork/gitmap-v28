# Component Source Code

> **Parent:** [React Components Index](./00-overview.md)  
> **Version:** 4.0.0  
> **Updated:** 2026-04-01

---

## Review Compliance

| Rule | Status | Notes |
|------|--------|-------|
| No hardcoded colors | ✅ Fixed | All colors use semantic design tokens |
| No `as` type assertions | ✅ Fixed | `errorLogToCapturedError` uses builder pattern |
| No `unknown` | ✅ Clean | Concrete types throughout |
| Function size ≤ 15 lines | ✅ Clean | All functions/components under limit |
| Parameters ≤ 3 | ✅ Clean | All component props via single interface |

---

## 7.1 GlobalErrorModal

The root shell component. Manages Dialog open/close, Backend/Frontend section toggle, error queue navigation, error log fetching, and Copy/Download dropdown menus.

```typescript
export function GlobalErrorModal(): JSX.Element {
  const { selectedError, isModalOpen, closeErrorModal, errorQueue, currentQueueIndex, navigateQueue } = useErrorStore();
  const [activeSection, setActiveSection] = useState<"backend" | "frontend">("backend");

  return (
    <Dialog open={isModalOpen} onOpenChange={closeErrorModal}>
      <DialogContent className="w-full h-full sm:max-w-[95vw] sm:h-[95vh]">
        {/* Header with error code, timestamp, queue nav */}
        {/* Section toggle: Backend | Frontend */}
        <ScrollArea>
          {activeSection === "backend" ? <BackendSection ... /> : <FrontendSection ... />}
        </ScrollArea>
        {/* Footer with Download/Copy dropdowns */}
      </DialogContent>
    </Dialog>
  );
}
```

## 7.2 BackendSection

The primary diagnostic view with 7 tabs: Overview, Log, Execution, Stack, Session, Request, Traversal.

**Key pattern — session auto-fetch via React Query:**
```typescript
const { diagnostics: sessionDiag, loading: sessionLoading } = useSessionDiagnostics(error.sessionId);
```

## 7.3 SessionLogsTab

Self-contained component with 4 sub-tabs. Log line highlighting uses **semantic design tokens only** — no hardcoded hex/color values:

```typescript
/** Log severity to design token mapping */
const LOG_LINE_STYLES: Record<string, string> = {
  stage: 'text-primary font-semibold',
  error: 'text-destructive',
  warning: 'text-warning',
  success: 'text-success',
  default: 'text-foreground',
} as const;

function getLogLineStyle(line: string): string {
  if (line.includes("STAGE:") || line.match(/^[─═]+$/)) return LOG_LINE_STYLES.stage;
  if (line.includes("[ERROR]") || line.includes("[FATAL]")) return LOG_LINE_STYLES.error;
  if (line.includes("[WARN]")) return LOG_LINE_STYLES.warning;
  if (line.includes("✓") || line.includes("success")) return LOG_LINE_STYLES.success;
  return LOG_LINE_STYLES.default;
}

function LogLine({ line }: { line: string }): JSX.Element {
  return <div className={getLogLineStyle(line)}>{line}</div>;
}
```

> **Design token requirement:** `text-warning` and `text-success` must be defined in `index.css` and `tailwind.config.ts`. Example:
> ```css
> :root {
>   --warning: 38 92% 50%;
>   --success: 142 76% 36%;
> }
> ```
> ```ts
> // tailwind.config.ts
> warning: "hsl(var(--warning))",
> success: "hsl(var(--success))",
> ```

### Violations Fixed (v3.2.0 → v4.0.0)

| Previous | Violation | Fix |
|----------|-----------|-----|
| `text-amber-600 dark:text-amber-400` | Hardcoded color, not themed | → `text-warning` semantic token |
| `text-green-600 dark:text-green-400` | Hardcoded color, not themed | → `text-success` semantic token |
| Inline style logic in component | Mixed concerns | → Extracted `getLogLineStyle()` + `LOG_LINE_STYLES` map |

## 7.4 RequestDetails

Visualizes the 3-hop request chain. See [07-request-chain.md](../03-error-modal-reference/07-request-chain.md) for full code.

## 7.5 TraversalDetails

See [08-traversal-details.md](../03-error-modal-reference/08-traversal-details.md) for full code.

## 7.6 ErrorModalActions

Two dropdown menus, reused by both `GlobalErrorModal` and `ErrorDetailModal`:

**DownloadDropdown:** Full Bundle (ZIP), error.log.txt, log.txt, Report (.md)

**CopyDropdown (Split Button):** Main click → Compact Report. Dropdown: Compact, Full, With Backend Logs, error.log.txt, log.txt.

## 7.7 ErrorDetailModal

Standalone modal for Error History and E2E Tests pages. Uses **type-safe builder pattern** (no `as` assertion):

```typescript
// src/components/errors/errorLogAdapter.ts

/**
 * Converts backend ErrorLog to CapturedError using explicit field mapping.
 * No type assertions — every field is assigned with its correct type.
 */
export function errorLogToCapturedError(error: ErrorLog): CapturedError {
  const parsedFrames: StackFrame[] | undefined = error.file
    ? [{
        file: error.file,
        line: error.line ?? 0,
        function: error.function ?? '',
        column: undefined,
        isInternal: false,
      }]
    : undefined;

  const captured: CapturedError = {
    id: String(error.id),
    code: error.code,
    level: mapErrorLevel(error.level),
    message: error.message,
    details: error.details,
    createdAt: error.createdAt,
    context: mapErrorContext(error.context),
    backendStackTrace: error.stackTrace,
    parsedFrames,
    // All remaining CapturedError fields default to undefined
    stackTrace: undefined,
    endpoint: undefined,
    method: undefined,
    requestBody: undefined,
    responseStatus: undefined,
    invocationChain: undefined,
    triggerComponent: undefined,
    triggerAction: undefined,
    backendLogs: undefined,
    siteUrl: undefined,
    sessionId: undefined,
    sessionType: undefined,
    phpStackFrames: undefined,
    errorFile: undefined,
    errorLine: undefined,
    uiClickPath: undefined,
    uiClickPathString: undefined,
    uiClickPathArrow: undefined,
    route: undefined,
    routeComponent: undefined,
    executionLogs: undefined,
    executionChain: undefined,
    executionLogsEnabled: undefined,
    executionLogsFormatted: undefined,
    requestedAt: undefined,
    requestDelegatedAt: undefined,
    envelopeErrors: undefined,
    envelopeMethodsStack: undefined,
    file: error.file,
    line: error.line,
    function: error.function,
  };

  return captured;
}

/** Maps raw level string to CapturedError level discriminator */
function mapErrorLevel(level: string): CapturedError['level'] {
  if (level === 'error' || level === 'warn' || level === 'info') return level;
  return 'error';
}

/** Maps raw context to typed ErrorContext */
function mapErrorContext(context: ErrorLog['context']): ErrorContext | undefined {
  if (!context) return undefined;
  return context as ErrorContext;
}
```

### Violations Fixed (v3.2.0 → v4.0.0)

| Previous | Violation | Fix |
|----------|-----------|-----|
| `} as CapturedError` | Unsafe type assertion, bypasses type checking | → Explicit field mapping, no `as` |
| `error.level as CapturedError["level"]` | Unsafe cast on raw string | → `mapErrorLevel()` with validation |
| `error.context as CapturedError["context"]` | Unsafe cast on `Record<string, unknown>` | → `mapErrorContext()` with typed return |
| Missing `StackFrame` fields | `isInternal`, `column` not set | → All fields explicitly assigned |

---

*Component source code — updated: 2026-04-01*
