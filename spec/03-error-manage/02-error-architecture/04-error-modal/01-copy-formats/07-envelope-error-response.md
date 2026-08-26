# Envelope Error Response (JSON)

> **Parent:** [Copy Formats Index](./00-overview.md)  
> **Version:** 3.2.0  
> **Updated:** 2026-03-31  
> **Purpose:** Raw JSON response structure from the Go backend when an error occurs.

---

See [envelope-error.json](../../05-response-envelope/envelope-error.json) for the canonical full sample.

---

## Key Fields for Error Display

```json
{
  "Status": {
    "Code": 500,
    "Message": "[E3001] failed to fetch snapshot settings..."
  },
  "Attributes": {
    "RequestedAt": "http://localhost:8080/api/v1/sites/1/snapshots/settings",
    "RequestDelegatedAt": "https://example.com/wp-json/.../snapshots/settings",
    "SessionId": "a1b2c3d4-...",
    "HasAnyErrors": true
  },
  "Errors": {
    "BackendMessage": "...",
    "DelegatedServiceErrorStack": ["..."],
    "Backend": ["..."],
    "DelegatedRequestServer": {
      "DelegatedEndpoint": "https://...",
      "Method": "GET",
      "StatusCode": 403,
      "RequestBody": null,
      "Response": { "...": "..." },
      "StackTrace": ["..."],
      "AdditionalMessages": "..."
    }
  },
  "MethodsStack": {
    "Backend": [{ "Method": "...", "File": "...", "LineNumber": 0 }],
    "Frontend": []
  }
}
```

---

## Frontend Extraction Map

| Envelope Field | → CapturedError Field | Used In |
|---------------|----------------------|---------|
| `Status.Message` | `message` | All views |
| `Attributes.RequestedAt` | `requestedAt` | Traversal tab |
| `Attributes.RequestDelegatedAt` | `requestDelegatedAt` | Traversal tab |
| `Attributes.SessionId` | `sessionId` | Session tab auto-fetch |
| `Errors.BackendMessage` | `backendMessage` | Overview tab banner |
| `Errors.Backend` | `backendStackTrace` | Stack tab |
| `Errors.DelegatedServiceErrorStack` | `phpStackFrames` (parsed) | Stack tab |
| `Errors.DelegatedRequestServer` | `delegatedRequestServer` | Stack + Request + Traversal |
| `MethodsStack.Backend` | `envelopeMethodsStack` | Execution + Traversal |

---

*Envelope error response — updated: 2026-03-31*
