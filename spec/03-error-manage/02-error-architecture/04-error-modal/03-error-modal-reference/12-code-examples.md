# React Code Examples

> **Parent:** [Error Modal Reference](./00-overview.md)  
> **Version:** 2.2.0  
> **Updated:** 2026-03-31

---

## Example: Capturing an API Error with Full Context

```typescript
import { useErrorStore } from '@/stores/errorStore';
import { api } from '@/lib/api';

async function enablePlugin(slug: string, siteId: number) {
  try {
    const response = await api.post('/api/v1/plugins/enable', {
      slug,
      siteId,
    });
    
    if (!response.success) {
      useErrorStore.getState().captureError(response.error!, {
        endpoint: '/api/v1/plugins/enable',
        method: HttpMethod.Post,
        requestBody: { slug, siteId },
        responseStatus: response.error?.statusCode,
        siteUrl: 'https://example.com',
        sessionId: response.error?.sessionId,
        sessionType: 'plugin_lifecycle',
        context: {
          source: 'PluginActions.enablePlugin',
          triggerComponent: 'PluginCard',
          triggerAction: 'enable_clicked',
        },
      });
    }

  } catch (error) {
    useErrorStore.getState().captureException(error, {
      source: 'PluginActions.enablePlugin',
      triggerComponent: 'PluginCard',
      triggerAction: 'enable_clicked',
      endpoint: '/api/v1/plugins/enable',
      method: HttpMethod.Post,
      requestBody: { slug, siteId },
    });
  }
}
```

---

## Example: Opening the Error Modal Manually

```typescript
import { useErrorStore } from '@/stores/errorStore';

function MyComponent() {
  const { openErrorModal, captureException } = useErrorStore();
  
  const handleDangerousAction = async () => {
    try {
      await riskyOperation();
    } catch (error) {
      const captured = captureException(error, {
        source: 'MyComponent.handleDangerousAction',
        triggerComponent: 'MyComponent',
        triggerAction: 'dangerous_action_clicked',
      });
      // Automatically opens the modal:
      openErrorModal(captured);
    }

  };
  
  return <Button onClick={handleDangerousAction}>Do Something</Button>;
}
```

---

## Example: Using the Error Boundary

```tsx
import { AppErrorBoundary } from '@/components/errors/AppErrorBoundary';

function App() {
  return (
    <AppErrorBoundary>
      <Router>
        <Routes>
          {/* All routes are wrapped — any unhandled error opens the modal */}
        </Routes>
      </Router>
    </AppErrorBoundary>
  );
}
```

---

## Example: Consuming Envelope Errors in a Hook

```typescript
import { useApiQuery } from '@/hooks/useApiQuery';

function usePlugins() {
  return useApiQuery<Plugin[]>('/plugins', ['plugins']);
  // If the query fails:
  // 1. useApiQuery calls requireSuccess()
  // 2. requireSuccess() detects envelope failure
  // 3. Error is auto-captured into errorStore with full envelope data
  // 4. Modal opens with Backend → Overview → "Backend Error" red banner
}
```

---

## Example: Rendering PHP Stack Frames

```tsx
// Extracting PHP frames from error context:
const phpStackFrames: PHPStackFrame[] = (() => {
  const ctx = selectedError?.context as Record<string, unknown> | undefined;
  if (!ctx) return [];
  if (Array.isArray(ctx.stackTraceFrames)) return ctx.stackTraceFrames;
  const errorDetails = ctx.errorDetails as Record<string, unknown> | undefined;
  if (errorDetails && Array.isArray(errorDetails.stackTraceFrames))
    return errorDetails.stackTraceFrames;
  return [];
})();

// Rendering:
{phpStackFrames.map((frame, i) => (
  <tr key={i}>
    <td className="p-2 font-mono text-muted-foreground">{i}</td>
    <td className="p-2 font-mono text-orange-500">
      {frame.class ? `${frame.class}::${frame.function}()` : `${frame.function}()`}
    </td>
    <td className="p-2 font-mono text-muted-foreground">{frame.fileBase || frame.file}</td>
    <td className="p-2 font-mono text-right">{frame.line || '?'}</td>
  </tr>
))}
```

---

*Code examples — updated: 2026-03-31*
