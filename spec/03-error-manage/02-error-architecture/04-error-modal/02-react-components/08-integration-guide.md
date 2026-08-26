# Integration Guide

> **Parent:** [React Components Index](./00-overview.md)  
> **Version:** 3.2.0  
> **Updated:** 2026-03-31

---

## Minimal Setup

1. **Install dependencies:** `zustand`, `@radix-ui/react-dialog`, `@radix-ui/react-tabs`, `lucide-react`, `sonner`
2. **Copy files:** All files from `src/components/errors/`, `src/stores/errorStore.ts`, `src/hooks/useSessionDiagnostics.ts`
3. **Mount the modal:**

```tsx
import { GlobalErrorModal } from '@/components/errors/GlobalErrorModal';

function App() {
  return (
    <>
      <RouterProvider router={router} />
      <GlobalErrorModal />
    </>
  );
}
```

4. **Capture errors in React Query:**

```tsx
const queryClient = new QueryClient({
  defaultOptions: { queries: { retry: false } },
});

queryClient.getQueryCache().subscribe((event) => {
  const isQueryError = event.type === 'updated' && event.query.state.status === 'error';
  if (!isQueryError) return;

  const isSuppressed = event.query.meta?.suppressGlobalError === true;
  if (isSuppressed) return;

  const error = event.query.state.error;
  useErrorStore.getState().captureException(error, {
    source: 'App.showGlobalError',
    triggerComponent: 'QueryClient',
    triggerAction: 'async_operation',
  });
});
```

---

## Required Utility Functions

```typescript
export function formatDateTimeUtc(isoString: string): string {
  return new Date(isoString).toLocaleString('en-US', { timeZone: 'UTC', ... });
}

export function toClipboardText(text: string): string {
  return text.replace(/\r\n/g, '\n').replace(/[ \t]+$/gm, '').trim();
}

export function unescapeEmbeddedNewlines(text: string): string {
  return text.replace(/\\n/g, '\n').replace(/\\t/g, '\t');
}
```

---

## Adapting for Non-WordPress Projects

- **Remove PHP-specific code:** `phpStackFrames`, `phpStackTraceLog`, PHP-related UI
- **Rename "Delegated Server":** Change labels to your architecture (e.g., "API → Microservice")
- **Simplify if no delegation:** Remove `RequestDetails` 3-hop chain, keep simple request display
- **Keep the core:** `CapturedError` model, `errorStore`, `GlobalErrorModal` shell, `FrontendSection`, `errorReportGenerator` are fully generic

---

*Integration guide — updated: 2026-03-31*
