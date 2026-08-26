# TypeScript Interfaces

> **Parent:** [React Components Index](./00-overview.md)  
> **Version:** 4.0.0  
> **Updated:** 2026-04-01

---

## Review Compliance

| Rule | Status | Notes |
|------|--------|-------|
| No `unknown` | ✅ Fixed | Replaced with typed alternatives |
| No `Record<string, unknown>` | ✅ Fixed | Replaced with named interfaces |
| No `any` | ✅ Clean | None present |
| PascalCase enums | ✅ Clean | Level uses string literal union (acceptable for 3-value discriminator) |
| Function size ≤ 15 lines | N/A | Interfaces only |

---

## 2.1 Core Error Model — `CapturedError`

This is the central type stored in the error store and consumed by all modal components. See [01-data-model.md](../03-error-modal-reference/01-data-model.md) for full documentation with field descriptions.

```typescript
/** Parsed stack frame with file, line, column info */
export interface StackFrame {
  function: string;
  file: string;
  line: number;
  column?: number;
  isInternal: boolean;
}

/** Structured key-value details attached to a log entry */
export interface LogEntryDetails {
  [key: string]: string | number | boolean | null;
}

export interface BackendLogEntry {
  timestamp: string;
  level: 'debug' | 'info' | 'warn' | 'error';
  message: string;
  step?: string;
  details?: LogEntryDetails;
}

export interface PHPStackFrame {
  file?: string;
  fileBase?: string;
  line?: number;
  function?: string;
  class?: string;
}

export interface ClickEvent {
  id: string;
  element: string;
  text?: string;
  action: string;
  componentName?: string;
  route?: string;
  timestamp: number;
}

export interface ExecutionLogEntry {
  timestamp: number;
  source: string;
  action: string;
  details?: LogEntryDetails;
}

export interface CallChain {
  chain: string[];
  formatted: string;
}

/** Delegated server error details from the response envelope */
export interface DelegatedRequestServerError {
  DelegatedEndpoint: string;
  Method: string;
  StatusCode: number;
  RequestBody?: string;
  Response?: string;
  StackTrace?: string[];
  AdditionalMessages?: string;
}

export interface EnvelopeErrors {
  BackendMessage: string;
  DelegatedServiceErrorStack?: string[];
  Backend?: string[];
  Frontend?: string[];
  DelegatedRequestServer?: DelegatedRequestServerError;
}

export interface EnvelopeMethodFrame {
  Method: string;
  File: string;
  LineNumber: number;
}

export interface EnvelopeMethodsStack {
  Backend: EnvelopeMethodFrame[];
  Frontend: EnvelopeMethodFrame[];
}

/** Structured context attached to a captured error */
export interface ErrorContext {
  source?: string;
  triggerComponent?: string;
  triggerAction?: string;
  [key: string]: string | number | boolean | null | undefined;
}

/** The full captured error object */
export interface CapturedError {
  id: string;
  code: string;
  level: 'error' | 'warn' | 'info';
  message: string;
  details?: string;
  context?: ErrorContext;
  file?: string;
  line?: number;
  function?: string;
  stackTrace?: string;
  createdAt: string;
  endpoint?: string;
  method?: string;
  requestBody?: string;
  responseStatus?: number;
  invocationChain?: string[];
  parsedFrames?: StackFrame[];
  triggerComponent?: string;
  triggerAction?: string;
  backendLogs?: BackendLogEntry[];
  backendStackTrace?: string;
  siteUrl?: string;
  sessionId?: string;
  sessionType?: string;
  phpStackFrames?: PHPStackFrame[];
  errorFile?: string;
  errorLine?: number;
  uiClickPath?: ClickEvent[];
  uiClickPathString?: string;
  uiClickPathArrow?: string;
  route?: string;
  routeComponent?: string;
  executionLogs?: ExecutionLogEntry[];
  executionChain?: CallChain | null;
  executionLogsEnabled?: boolean;
  executionLogsFormatted?: string;
  requestedAt?: string;
  requestDelegatedAt?: string;
  envelopeErrors?: EnvelopeErrors;
  envelopeMethodsStack?: EnvelopeMethodsStack;
}
```

### Violations Fixed (v3.2.0 → v4.0.0)

| Previous | Violation | Fix |
|----------|-----------|-----|
| `details?: Record<string, unknown>` | Banned `Record<string, unknown>` | → `LogEntryDetails` (typed map) |
| `context?: Record<string, unknown>` | Banned `Record<string, unknown>` | → `ErrorContext` (named interface) |
| `requestBody?: unknown` | Banned `unknown` | → `string` (serialized JSON) |
| `Response?: unknown` in `DelegatedRequestServer` | Banned `unknown` | → `string` (serialized response) |
| Inline `DelegatedRequestServer` type | Large inline object | → `DelegatedRequestServerError` interface |

---

## 2.2 Session Diagnostics

```typescript
export interface SessionStackFrame {
  function: string;
  file?: string;
  line?: number;
  class?: string;
}

/** Typed HTTP headers map */
export interface HttpHeaders {
  [header: string]: string;
}

/** Typed request body for session diagnostics */
export interface SessionRequestBody {
  [key: string]: string | number | boolean | null;
}

export interface SessionRequest {
  url: string;
  method: string;
  headers?: HttpHeaders;
  body?: SessionRequestBody;
}

export interface SessionResponse {
  requestUrl: string;
  responseUrl: string;
  statusCode: number;
  headers?: HttpHeaders;
  body?: string;
}

export interface SessionStackTraces {
  golang?: SessionStackFrame[];
  php?: SessionStackFrame[];
}

export interface SessionDiagnostics {
  request?: SessionRequest;
  response?: SessionResponse;
  stackTrace?: SessionStackTraces;
  phpStackTraceLog?: string;
}
```

### Violations Fixed (v3.2.0 → v4.0.0)

| Previous | Violation | Fix |
|----------|-----------|-----|
| `headers?: Record<string, string>` | Generic record | → `HttpHeaders` (named type) |
| `body?: Record<string, unknown>` | Banned `Record<string, unknown>` | → `SessionRequestBody` (typed map) |
| `body?: unknown` on response | Banned `unknown` | → `string` (serialized) |
| Inline nested objects | Hard to test individually | → Named interfaces per object |

---

## 2.3 Shared Component Props

```typescript
export interface SectionCommonProps {
  copySection: (label: string, content: string) => void;
  formatTs: (ts: string) => string;
}

export interface AppInfo {
  appName: string;
  appVersion: string;
  gitCommit?: string;
  buildTime?: string;
}
```

---

*TypeScript interfaces — updated: 2026-04-01*
