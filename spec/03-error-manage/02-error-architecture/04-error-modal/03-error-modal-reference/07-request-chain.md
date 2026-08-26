# Request Chain Visualization

> **Parent:** [Error Modal Reference](./00-overview.md)  
> **Version:** 2.2.0  
> **Updated:** 2026-03-31

---

The `RequestDetails` component renders the full React → Go → Delegated Server chain:

```tsx
// Key data sources:
const hasRequestChain = error.requestedAt || error.requestDelegatedAt;
const hasDelegation = error.requestDelegatedAt || sessionDiagnostics?.response?.requestUrl;
const hasDelegatedServer = error.delegatedRequestServer || error.envelopeErrors?.DelegatedRequestServer;

// Node 1: React → Go
// - endpoint, method, responseStatus from CapturedError
// - requestBody (collapsible JSON)

// Node 2: Go → Delegated (shown if delegated)
// - delegatedEndpoint from DelegatedRequestServer or requestDelegatedAt
// - method, statusCode from DelegatedRequestServer
// - requestBody from DelegatedRequestServer (if POST/PUT)

// Node 3: Delegated Response (shown if DelegatedRequestServer present, NEW v2.0.0)
// - Response JSON from DelegatedRequestServer.Response
// - StackTrace from DelegatedRequestServer.StackTrace
// - AdditionalMessages from DelegatedRequestServer.AdditionalMessages
```

---

## Sample Code: Rendering the 3-Hop Request Chain

```tsx
{/* Node 1: React → Go */}
<div className="border rounded-t-md bg-muted/50 p-3">
  <div className="flex items-center gap-2 mb-1">
    <div className="w-2 h-2 rounded-full bg-blue-500" />
    <Badge variant="outline" className="font-mono bg-blue-500/10 border-blue-500/30">
      React → Go
    </Badge>
    <Badge variant="outline" className="font-mono">{error.method}</Badge>
    <Badge variant={error.responseStatus >= 400 ? "destructive" : "secondary"}>
      {error.responseStatus}
    </Badge>
  </div>
  <p className="ml-4 text-xs font-mono text-muted-foreground break-all">
    {error.requestedAt || error.endpoint}
  </p>
</div>

{/* Connector */}
<div className="flex items-center pl-4">
  <div className="w-0.5 h-4 bg-border ml-[3px]" />
</div>

{/* Node 2: Go → Delegated */}
<div className={cn("border bg-muted/50 p-3", hasDelegatedServer ? "" : "rounded-b-md")}>
  <div className="flex items-center gap-2 mb-1">
    <div className="w-2 h-2 rounded-full bg-orange-500" />
    <Badge variant="outline" className="font-mono bg-orange-500/10 border-orange-500/30 text-orange-600">
      Go → Delegated
    </Badge>
    {delegatedServer && (
      <>
        <Badge variant="outline" className="font-mono">{delegatedServer.Method}</Badge>
        <Badge variant={delegatedServer.StatusCode >= 400 ? "destructive" : "secondary"}>
          {delegatedServer.StatusCode}
        </Badge>
      </>
    )}
  </div>
  <p className="ml-4 text-xs font-mono text-orange-600 break-all">
    {delegatedServer?.DelegatedEndpoint || error.requestDelegatedAt}
  </p>
</div>

{/* Node 3: Delegated Response (NEW v2.0.0) */}
{hasDelegatedServer && (
  <>
    <div className="flex items-center pl-4">
      <div className="w-0.5 h-4 bg-border ml-[3px]" />
    </div>
    <div className="border rounded-b-md bg-muted/50 p-3">
      <div className="flex items-center gap-2 mb-1">
        <div className="w-2 h-2 rounded-full bg-purple-500" />
        <Badge variant="outline" className="font-mono bg-purple-500/10 border-purple-500/30 text-purple-600">
          Delegated Response
        </Badge>
      </div>
      {delegatedServer?.Response && (
        <Collapsible>
          <CollapsibleTrigger className="text-xs text-muted-foreground">
            ▸ Response JSON
          </CollapsibleTrigger>
          <CollapsibleContent>
            <pre className="text-xs bg-muted p-2 rounded mt-1 overflow-auto max-h-32">
              {JSON.stringify(delegatedServer.Response, null, 2)}
            </pre>
          </CollapsibleContent>
        </Collapsible>
      )}
      {delegatedServer?.StackTrace?.length > 0 && (
        <div className="mt-2">
          <p className="text-xs text-muted-foreground mb-1">Stack Trace:</p>
          {delegatedServer.StackTrace.map((line, i) => (
            <p key={i} className="text-xs font-mono text-purple-600 ml-2">{line}</p>
          ))}
        </div>
      )}
      {delegatedServer?.AdditionalMessages && (
        <p className="mt-2 text-xs text-purple-600">
          ℹ {delegatedServer.AdditionalMessages}
        </p>
      )}
    </div>
  </>
)}
```

---

*Request chain visualization — updated: 2026-03-31*
