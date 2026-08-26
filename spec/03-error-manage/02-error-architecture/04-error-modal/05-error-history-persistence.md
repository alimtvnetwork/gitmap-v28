# Error Modal — Error History Persistence

> **Version:** 1.1.0  
> **Updated:** 2026-04-01
> **Status:** Active  
> **Purpose:** Specifies the error history system that persists captured errors to the backend and provides browsing, export, and multi-error queue functionality.

---

## 1. Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                     Error History Flow                                │
│                                                                       │
│  Error Captured ──▸ errorStore (Zustand, in-memory)                   │
│       │                    │                                          │
│       │                    ▼                                          │
│       │          useErrorHistorySync (auto-save)                      │
│       │                    │                                          │
│       │                    ▼                                          │
│       │          POST /api/v1/error-history  ──▸ Backend DB           │
│       │                                                               │
│       ▼                                                               │
│  GlobalErrorModal (live view)                                         │
│                                                                       │
│  ErrorQueueBadge ──▸ ErrorHistoryDrawer ──▸ GET /api/v1/error-history│
│       │                    │                                          │
│       │                    ├── View in GlobalErrorModal (single)       │
│       │                    ├── View in queue (multi-select)            │
│       │                    ├── Export (POST /api/v1/error-history/bulk)│
│       │                    ├── Delete (DELETE /api/v1/error-history/:id)│
│       │                    └── Clear (DELETE /api/v1/error-history)    │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 2. Hooks

### 2.1 `useErrorHistory` — CRUD Operations

Located in `src/hooks/useErrorHistory.ts`.

```typescript
export function useErrorHistory() {
  // Query: GET /api/v1/error-history?limit=100
  const { data: response, isLoading, error, refetch } = useQuery({
    queryKey: ["error-history"],
    queryFn: () => api.listErrorHistory({ limit: 100 }),
    staleTime: 30000,
  });

  // Mutations — all use meta: { suppressGlobalError: true }
  const saveMutation = useMutation({
    meta: { suppressGlobalError: true },
    mutationFn: (input: ErrorHistoryInput) => api.saveErrorHistory(input),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["error-history"] }),
  });

  const deleteMutation = useMutation({
    meta: { suppressGlobalError: true },
    mutationFn: (id: number) => api.deleteErrorHistory(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["error-history"] }),
  });

  const clearMutation = useMutation({
    meta: { suppressGlobalError: true },
    mutationFn: () => api.clearErrorHistory(),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["error-history"] }),
  });

  const exportMutation = useMutation({
    meta: { suppressGlobalError: true },
    mutationFn: (ids: number[]) => api.bulkExportErrorHistory(ids),
  });

  return {
    errors: data?.errors || [],
    total: data?.total || 0,
    isLoading, error, refetch,
    saveError, deleteError, clearErrors, exportErrors,
    isExporting: exportMutation.isPending,
  };
}
```

### 2.2 `useErrorHistorySync` — Auto-Save on Capture

Mounted in `App.tsx` — automatically persists new errors from the Zustand store to the backend.

```typescript
export function useErrorHistorySync() {
  const { getPendingSyncErrors, markErrorSynced } = useErrorStore();
  const syncingRef = useRef<Set<string>>(new Set());

  const saveMutation = useMutation({
    meta: { suppressGlobalError: true },
    mutationFn: (input: ErrorHistoryInput) => api.saveErrorHistory(input),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["error-history"] }),
  });

  useEffect(() => {
    const pendingErrors = getPendingSyncErrors();

    for (const error of pendingErrors) {
      if (syncingRef.current.has(error.id)) continue;
      syncingRef.current.add(error.id);
      saveMutation.mutate(capturedToInput(error), {
        onSuccess: () => { markErrorSynced(error.id); syncingRef.current.delete(error.id); },
        onError: () => { syncingRef.current.delete(error.id); /* retry next render */ },
      });
    }

  }, [getPendingSyncErrors, markErrorSynced, capturedToInput, saveMutation]);
}
```

### 2.3 `recordToCapturedError` — Backend → Frontend Conversion

Converts `ErrorHistoryRecord` (backend API shape) back to `CapturedError` for display in the GlobalErrorModal.

> **v1.1.0:** Removed unsafe `as` cast on `level`; uses validated mapping function instead.

```typescript
/** Maps raw level string to CapturedError level discriminator */
function mapErrorLevel(level: string): CapturedError['level'] {
  if (level === 'error' || level === 'warn' || level === 'info') return level;
  return 'error';
}

export function recordToCapturedError(record: ErrorHistoryRecord): CapturedError {
  return {
    id: record.errorId,
    code: record.code,
    level: mapErrorLevel(record.level),
    message: record.message,
    details: record.details,
    context: record.context,
    stackTrace: record.stackTrace,
    endpoint: record.endpoint,
    method: record.method,
    requestBody: record.requestBody,
    responseStatus: record.responseStatus,
    sessionId: record.sessionId,
    sessionType: record.sessionType,
    phpStackFrames: record.phpStackFrames,
    backendStackTrace: record.backendStackTrace,
    siteUrl: record.siteUrl,
    triggerComponent: record.triggerComponent,
    triggerAction: record.triggerAction,
    invocationChain: record.invocationChain,
    createdAt: record.createdAt,
  };
}
```

---

## 3. Backend API Endpoints

| Method | Endpoint | Purpose |
|--------|----------|---------|
| `GET` | `/api/v1/error-history?limit=100` | List recent errors |
| `POST` | `/api/v1/error-history` | Save a new error |
| `DELETE` | `/api/v1/error-history/:id` | Delete a single error |
| `DELETE` | `/api/v1/error-history` | Clear all errors |
| `POST` | `/api/v1/error-history/bulk-export` | Export selected errors as Markdown report |

### ErrorHistoryInput (POST body)

```typescript
/** Typed context for error history persistence */
interface ErrorHistoryContext {
  [key: string]: string | number | boolean | null;
}

interface ErrorHistoryInput {
  errorId: string;
  code: string;
  level: string;
  message: string;
  details?: string;
  context?: ErrorHistoryContext;
  stackTrace?: string;
  endpoint?: string;
  method?: string;
  requestBody?: RequestPayload;
  responseStatus?: number;
  sessionId?: string;
  sessionType?: string;
  phpStackFrames?: PHPStackFrame[];
  backendLogs?: string[];           // Formatted as "[timestamp] [level] message"
  backendStackTrace?: string;
  siteUrl?: string;
  triggerComponent?: string;
  triggerAction?: string;
  invocationChain?: string[];
}
```

---

## 4. UI Components

### 4.1 ErrorQueueBadge

**File:** `src/components/errors/ErrorQueueBadge.tsx`  
**Location:** Application header/navigation bar  
**Visibility:** Only shown when `displayCount > 0`

```tsx
export function ErrorQueueBadge() {
  const { recentErrors } = useErrorStore();
  const { total } = useErrorHistory();
  const displayCount = recentErrors.length || total;

  if (displayCount === 0) return null;

  return (
    <Button
      variant="ghost"
      className="relative h-8 px-2 text-destructive hover:text-destructive hover:bg-destructive/10"
      onClick={() => setDrawerOpen(true)}
    >
      <AlertCircle className="h-4 w-4" />
      <Badge variant="destructive" className="absolute -top-1 -right-1 h-5 min-w-5 px-1 text-xs">
        {displayCount > 99 ? "99+" : displayCount}
      </Badge>
    </Button>
  );
}
```

### 4.2 ErrorHistoryDrawer

**File:** `src/components/errors/ErrorHistoryDrawer.tsx`  
**Component:** `Sheet` (side drawer, `sm:max-w-lg`)

**Features:**
- Search/filter by message, code, endpoint
- Multi-select with checkboxes
- Bulk actions: View (in queue), Copy (export), Delete
- Single-click opens error in GlobalErrorModal
- Select All / Deselect All toggle
- Clear all with confirmation dialog
- Refresh button

**Item card states:**
```typescript
selected: "bg-accent border-primary"
default:  "bg-card hover:bg-accent/50"
```

### 4.3 ErrorDetailModal

**File:** `src/components/errors/ErrorDetailModal.tsx`  
**Used on:** Errors page, E2E Tests page  
**Receives:** `ErrorLog` (backend API type), adapts via `errorLogToCapturedError()`

**Tabs:** Overview, Stack Trace, Request/Response, Suggested Fixes  
**Footer:** `DownloadDropdown` + Close + Split Copy Button (compact/full)

---

## 5. Data Flow: Error Lifecycle

```
1. Error occurs (API failure, render crash, etc.)
       │
2. errorStore.captureError() or captureException()
       │ ── Captures: stack trace, click path, execution logs, envelope data
       │ ── Adds to recentErrors[] with pendingSync flag
       │
3. useErrorHistorySync (useEffect) detects pending error
       │ ── POST /api/v1/error-history
       │ ── On success: markErrorSynced(id)
       │ ── On failure: retry on next render cycle
       │
4. User views errors via:
       ├── GlobalErrorModal (immediate, from capture)
       ├── ErrorQueueBadge → ErrorHistoryDrawer (persisted list)
       └── Errors page → ErrorDetailModal (backend ErrorLog format)
```

---

## 6. AppErrorBoundary

**File:** `src/components/errors/AppErrorBoundary.tsx`  
**Type:** React class component (Error Boundary)  
**Purpose:** Prevents full white-screen crashes by catching render errors

```tsx
interface Props {
  children: React.ReactNode;
}

interface State {
  captured: CapturedError | null;
}

export class AppErrorBoundary extends React.Component<Props, State> {
  state: State = { captured: null };

  componentDidCatch(error: Error, info: React.ErrorInfo): void {
    const { captureException, openErrorModal } = useErrorStore.getState();
    const captured = captureException(error, {
      source: 'AppErrorBoundary.componentDidCatch',
      triggerComponent: 'ErrorBoundary',
      triggerAction: 'render_crash',
    });
    openErrorModal(captured);
    this.setState({ captured });
  }

  render() {
    if (!this.state.captured) return this.props.children;
    return (
      <div className="flex min-h-[60vh] items-center justify-center px-6 py-12">
        <div className="w-full max-w-lg space-y-4 rounded-lg border bg-background p-6">
          <h1 className="text-lg font-semibold">Something went wrong</h1>
          <p className="text-sm text-muted-foreground">
            A UI render error was caught. You can open the error details modal...
          </p>
          <Button variant="outline" onClick={() => openErrorModal(this.state.captured!)}>
            View error details
          </Button>
          <Button onClick={() => window.location.reload()}>Reload</Button>
        </div>
      </div>
    );
  }
}
```

**Integration:**
```tsx
// App.tsx
<AppErrorBoundary>
  <RouterProvider router={router} />
</AppErrorBoundary>
<GlobalErrorModal />
```

---

## Cross-References

- [Error Modal Overview](./00-overview.md) — Modal architecture overview
- [Error Modal Reference](./03-error-modal-reference.md) — GlobalErrorModal structure and data model
- [suppressGlobalError Pattern](./06-suppress-global-error.md) — Why all history mutations use the flag
- [Error Handling Cross-Stack Spec](../01-error-handling-reference.md) — Tier 3 frontend capture pipeline
- [React Components Reference](./02-react-components.md) — Component code

---

*Error History Persistence v1.0.0 — created: 2026-03-03*
