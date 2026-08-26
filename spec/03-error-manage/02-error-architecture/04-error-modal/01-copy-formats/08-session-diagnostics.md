# Session Diagnostics (JSON)

> **Parent:** [Copy Formats Index](./00-overview.md)  
> **Version:** 3.2.0  
> **Updated:** 2026-03-31  
> **Purpose:** Session-linked request/response data fetched when a sessionId is available on the error.

---

Fetched via:
- `GET /api/v1/sessions/{sessionId}/diagnostics` (for diagnostics)
- `GET /api/v1/sessions/{sessionId}/logs` (for logs)
- `GET /api/v1/request-sessions/{sessionId}` (raw session data)

---

## Complete Sample

```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "method": "GET",
  "path": "http://localhost:8080/api/v1/sites/1/snapshots/settings",
  "QueryString": "",
  "RequestHeaders": {
    "content-type": "application/json",
    "authorization": "[REDACTED]",
    "user-agent": "Mozilla/5.0 ..."
  },
  "RequestBody": "",
  "ResponseStatus": 500,
  "ResponseBody": "{\"Status\":{\"IsSuccess\":false,...},\"Errors\":{\"BackendMessage\":\"...\",\"DelegatedRequestServer\":{...}}}",
  "StartTime": "2026-02-11T16:58:26.000Z",
  "EndTime": "2026-02-11T16:58:32.730Z",
  "DurationMs": 6730,
  "error": "[E3001] failed to fetch snapshot settings: get snapshot settings (GET https://demoat.attoproperty.com.au/wp-json/riseup-asia-uploader/v1/snapshots/settings): status 403"
}
```

---

## Session ↔ Error Linkage

```
Error Modal opens
  │
  ├─ CapturedError.sessionId present?
  │     │
  │     ├─ YES → Auto-fetch GET /api/v1/request-sessions/{id}
  │     │         → Merge into Session tab (logs, request, response, stack trace)
  │     │         → Parse responseBody for DelegatedRequestServer
  │     │
  │     └─ NO → Session tab shows "No session data available"
  │
  └─ CapturedError.envelopeErrors.DelegatedRequestServer present?
        │
        ├─ YES → Show in Stack tab (delegated server section)
        │         → Show in Request tab (3-hop chain: React → Go → Delegated)
        │         → Show in Traversal tab
        │
        └─ NO → Standard 2-hop display (React → Go)
```

---

*Session diagnostics — updated: 2026-03-31*
