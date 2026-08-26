# Traversal Details

> **Parent:** [Error Modal Reference](./00-overview.md)  
> **Version:** 2.2.0  
> **Updated:** 2026-03-31

---

The `TraversalDetails` component provides a unified view of the request lifecycle:

```tsx
interface TraversalDetailsProps {
  error: CapturedError;
  copySection: (label: string, content: string) => void;
}

// Data sources:
const hasEndpoints = error.requestedAt || error.requestDelegatedAt;
const hasMethodsStack = error.envelopeMethodsStack?.Backend?.length > 0;
const hasDelegatedStack = error.envelopeErrors?.DelegatedServiceErrorStack?.length > 0;
const hasBackendTrace = error.envelopeErrors?.Backend?.length > 0;
const hasDelegatedServer = error.delegatedRequestServer || error.envelopeErrors?.DelegatedRequestServer;
```

---

## Sections Rendered

1. **Endpoint Flow** — 3-hop badges: React → Go → Delegated (with method + status for delegated hop)
2. **Methods Stack Table** — Go call chain from envelope
3. **Delegated Server Details** — Purple-themed block (NEW v2.0.0): endpoint, method, status, stack trace, response JSON, additional messages
4. **Delegated Service Error Stack** — Orange-themed PHP error lines (legacy, ScrollArea)
5. **Backend Trace** — Go stack trace lines (ScrollArea)

---

## Sample Code: Rendering Delegated Server Details in Traversal

```tsx
{/* Delegated Server Details (NEW v2.0.0) */}
{hasDelegatedServer && (() => {
  const ds = error.delegatedRequestServer || error.envelopeErrors?.DelegatedRequestServer;
  if (!ds) return null;
  return (
    <div className="space-y-2">
      <div className="flex items-center justify-between">
        <h4 className="text-sm font-medium flex items-center gap-2">
          <div className="w-2 h-2 rounded-full bg-purple-500" />
          Delegated Server Details
        </h4>
        <Button variant="ghost" size="sm" onClick={() => copySection("Delegated Server", 
          `Endpoint: ${ds.DelegatedEndpoint}\nMethod: ${ds.Method}\nStatus: ${ds.StatusCode}\n` +
          (ds.StackTrace?.join('\n') || '') + '\n' + (ds.AdditionalMessages || '')
        )}>
          <Copy className="h-3 w-3" />
        </Button>
      </div>
      <div className="bg-purple-500/5 border border-purple-500/20 rounded p-3 space-y-2">
        <div className="flex items-center gap-2 text-xs">
          <Badge variant="outline" className="font-mono">{ds.Method}</Badge>
          <Badge variant={ds.StatusCode >= 400 ? "destructive" : "secondary"}>{ds.StatusCode}</Badge>
        </div>
        <p className="text-xs font-mono text-purple-600 break-all">{ds.DelegatedEndpoint}</p>
        {ds.StackTrace?.length > 0 && (
          <ScrollArea className="max-h-32">
            {ds.StackTrace.map((line, i) => (
              <p key={i} className="text-xs font-mono text-purple-500">{line}</p>
            ))}
          </ScrollArea>
        )}
        {ds.AdditionalMessages && (
          <p className="text-xs text-purple-600 border-t border-purple-500/20 pt-2">
            ℹ {ds.AdditionalMessages}
          </p>
        )}
      </div>
    </div>
  );
})()}
```

---

*Traversal details — updated: 2026-03-31*
