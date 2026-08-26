# Error Modal — suppressGlobalError Meta Pattern

> **Version:** 1.2.0  
> **Updated:** 2026-04-01
> **Status:** Active  
> **Purpose:** Defines the `suppressGlobalError` React Query meta pattern that prevents duplicate error displays when queries/mutations handle errors locally.

---

## 1. Problem Statement

The application has a **global error handler** on React Query's `QueryCache` and `MutationCache` that captures all unhandled errors and displays them in the Global Error Modal. However, many queries and mutations already provide **local error feedback** via:

- `onError` callbacks with `toast.error()`
- `isError` state with inline retry UI
- Try-catch blocks with local toast notifications

Without suppression, the user sees **both** the local toast/inline error **and** the global error modal — a confusing double-notification.

---

## 2. Solution: `meta.suppressGlobalError`

Any React Query `useQuery` or `useMutation` that provides its own error feedback **MUST** include `meta: { suppressGlobalError: true }` to opt out of the global handler.

### 2.1 Global Handler Implementation

```typescript
// src/App.tsx — QueryClient configuration
import { QueryClient, QueryCache, MutationCache } from "@tanstack/react-query";

// showGlobalError() handles both ApiClientError and generic errors:
// - ApiClientError → captures via captureError, shows toast with "View Details" action
//   (E9005 errors open the modal immediately)
// - Generic errors → captures via captureException, shows toast with "View Details" action
function showGlobalError(error: Error | ApiClientError, context?: { endpoint?: string; method?: string }): void {
  const { captureError, captureException, openErrorModal } = useErrorStore.getState();

  if (isApiClientError(error)) {
    const captured = captureError(error.apiError, { /* request metadata */ });

    if (error.apiError.code === "E9005") { openErrorModal(captured); return; }
    toast.error(error.apiError.message, {
      action: { label: "View Details", onClick: () => openErrorModal(captured) },
    });

    return;
  }

  const captured = captureException(error, { source: "App.showGlobalError", /* ... */ });
  toast.error(`Request failed: ${context?.endpoint}`, {
    action: { label: "View Details", onClick: () => openErrorModal(captured) },
  });
}

const queryClient = new QueryClient({
  queryCache: new QueryCache({
    onError: (error, query) => {
      if (query.meta?.suppressGlobalError) return;
      showGlobalError(error, { endpoint: String(query.queryKey?.[0] ?? "query") });
    },
  }),
  mutationCache: new MutationCache({
    onError: (error, _variables, _context, mutation) => {
      if (mutation.meta?.suppressGlobalError) return;
      showGlobalError(error, { endpoint: String(mutation.options.mutationKey?.[0] ?? "mutation") });
    },
  }),
  defaultOptions: {
    queries: { staleTime: 5 * 60 * 1000, retry: false, refetchOnWindowFocus: false },
  },
});
```

### 2.2 Usage in Queries

```typescript
// Query with local error handling via isError + inline retry UI
const { data, isError, refetch } = useQuery({
  queryKey: ['remote-plugins', siteId],
  queryFn: async () => { /* ... */ },
  retry: false,
  meta: { suppressGlobalError: true },  // ← REQUIRED
});

// In JSX:
{isError && (
  <div className="text-destructive">
    <p>Failed to load</p>
    <Button onClick={() => refetch()}>Try again</Button>
  </div>
)}
```

### 2.3 Usage in Mutations

```typescript
// Mutation with local error handling via onError + toast
const deleteMutation = useMutation({
  meta: { suppressGlobalError: true },  // ← REQUIRED
  mutationFn: async (id: string) => {
    const response = await api.deleteSession(id);

    return requireSuccess(response, { endpoint: `/sessions/${id}`, method: "DELETE" });
  },
  onSuccess: () => {
    toast.success("Deleted successfully");
    queryClient.invalidateQueries({ queryKey: ["sessions"] });
  },
  onError: () => {
    toast.error("Failed to delete");
  },
});
```

---

## 3. Decision Matrix: When to Use

| Error Handling Pattern | Uses suppressGlobalError? | Example |
|----------------------|--------------------------|---------|
| `onError` callback with `toast.error()` | ✅ Yes | Session delete, clear mutations |
| `isError` state with inline retry button | ✅ Yes | Remote plugins panel |
| Try-catch with `toast.error()` in `mutationFn` | ✅ Yes | Test run start/rerun |
| No local error handling (relies on global modal) | ❌ No | Sites list query, settings query |
| Intentionally uses global modal (comment: "Ensure errors surface in GlobalErrorModal") | ❌ No | Theme setting mutation |
| Error silently swallowed (`console.warn` only) | ✅ Yes | Error history sync (background) |

**Rule of thumb:** If the user will see **any** local error feedback (toast, inline message, retry button), add the flag.

---

## 4. Complete Registry

All queries and mutations with `suppressGlobalError: true`:

| File | Queries/Mutations |
|------|-------------------|
| `src/components/sites/RemotePluginsPanel.tsx` | remote-plugins query, forceSyncMutation, toggleMutation, deleteMutation |
| `src/components/sites/RemotePluginFileBrowser.tsx` | plugin files query |
| `src/components/sites/SiteCard.tsx` | snapshots query, cron jobs query |
| `src/hooks/useRemoteSnapshots.ts` | All queries (snapshots, settings, providers) and all mutations (create, delete, restore, updateSettings, fullBackup, incrementalBackup, import, cleanup) |
| `src/hooks/useSiteHealth.ts` | useCheckAllSitesHealth mutation |
| `src/hooks/useErrorHistory.ts` | All mutations (save, delete, clear, export) in both useErrorHistory and useErrorHistorySync |
| `src/components/settings/SnapshotSettingsTab.tsx` | Snapshot cron queries (×2), snapshots query |
| `src/pages/Sessions.tsx` | deleteMutation |
| `src/pages/RequestSessions.tsx` | deleteMutation, clearMutation |
| `src/pages/Tests.tsx` | startRun, rerunCase |
| `src/hooks/useSettings.ts` | `useSaveSettings` mutation (called with inline `onError` toasts throughout Settings.tsx) |
| `src/hooks/useRemoteDebugRoutes.ts` | debug-routes query (inline error UI with retry in DebugRoutesPanel) |

### Intentionally NOT Suppressed

| File | Query/Mutation | Reason |
|------|---------------|--------|
| `src/hooks/useTheme.ts` | `updateSettingMutation` | Comment says "Ensure errors surface in GlobalErrorModal" |
| `src/hooks/useDashboardStats.ts` | dashboard-stats query | No local error handling |
| `src/pages/PublishHistory.tsx` | publish history queries | No local error handling |
| Most `useApiQuery` consumers | Various | `useApiQuery` uses `requireSuccess` which throws → caught by global handler |

---

## 5. Adding New Queries/Mutations — Checklist

When creating a new `useQuery` or `useMutation`:

1. **Does it have an `onError` callback?** → Add `meta: { suppressGlobalError: true }`
2. **Does it use `isError` for inline UI?** → Add `meta: { suppressGlobalError: true }`
3. **Does it catch errors and show toasts?** → Add `meta: { suppressGlobalError: true }`
4. **Does it have no local error handling at all?** → Do NOT add the flag (let global handler catch it)
5. **Update this registry** after adding the flag

---

## Cross-References

- [Error Handling Cross-Stack Spec](../01-error-handling-reference.md) — Tier 3: Frontend Error Handling
- [Error Modal Reference](./03-error-modal-reference.md) — §13 React Code Examples
- [Error History Persistence](./05-error-history-persistence.md) — Why all history mutations use the flag

---

*suppressGlobalError Pattern v1.1.0 — updated: 2026-03-03*
