# Error Capture Pipeline

> **Parent:** [Error Modal Reference](./00-overview.md)  
> **Version:** 2.3.0  
> **Updated:** 2026-04-01

---

## Step 1: API Client Detects Failure

When an API call fails, the client extracts envelope data:

```typescript
// src/lib/api/envelope.ts

/** Typed shape of a raw Universal Response Envelope from the backend */
interface RawEnvelope {
  Status: { IsSuccess: boolean; Code: number; Message: string };
  Attributes: {
    RequestedAt?: string;
    RequestDelegatedAt?: string;
    SessionId?: string;
    HasAnyErrors?: boolean;
  };
  Results: unknown[];
  Errors?: EnvelopeErrors;
  MethodsStack?: EnvelopeMethodsStack;
}

function isRawEnvelope(value: object): value is RawEnvelope {
  return 'Status' in value && 'Attributes' in value && 'Results' in value;
}

export function parseEnvelope(response: object | null): ParsedEnvelope | null {
  if (!response) return null;

  if (!isRawEnvelope(response)) return null;

  return {
    isSuccess: response.Status.IsSuccess,
    code: response.Status.Code,
    message: response.Status.Message,
    requestedAt: response.Attributes.RequestedAt,
    requestDelegatedAt: response.Attributes.RequestDelegatedAt,
    sessionId: response.Attributes.SessionId,
    hasErrors: response.Attributes.HasAnyErrors ?? false,
    errors: response.Errors ?? null,
    methodsStack: response.MethodsStack ?? null,
    results: response.Results,
  };
}
```

---

## Step 2: Error Store Captures & Enriches

```typescript
// src/stores/errorStore.ts
captureError: (apiError, meta) => {
  const captured = buildCapturedError(apiError, meta);
  
  // Enrich with UI click path (last 10 clicks)
  const clickPath = getClickPathForError();
  captured.uiClickPath = clickPath.events;
  captured.uiClickPathString = clickPath.formatted;
  
  // Enrich with React execution logs (if debug mode enabled)
  const execLogs = getExecutionLogsForError();
  captured.executionLogs = execLogs.entries;
  captured.executionChain = execLogs.chain;
  captured.executionLogsEnabled = execLogs.enabled;
  captured.executionLogsFormatted = execLogs.formatted;
  
  // Extract envelope diagnostic fields
  if (apiError.envelope) {
    captured.requestedAt = apiError.envelope.requestedAt;
    captured.requestDelegatedAt = apiError.envelope.requestDelegatedAt;
    captured.sessionId = apiError.envelope.sessionId || meta?.sessionId;
    captured.envelopeErrors = apiError.envelope.errors;
    captured.envelopeMethodsStack = apiError.envelope.methodsStack;
  }
  
  commitErrorToStore(captured);

  return captured;
},
```

---

## Step 3: Modal Opens

```typescript
// Anywhere in the app:
const { openErrorModal } = useErrorStore();
openErrorModal(capturedError);
```

---

*Capture pipeline — updated: 2026-03-31*
