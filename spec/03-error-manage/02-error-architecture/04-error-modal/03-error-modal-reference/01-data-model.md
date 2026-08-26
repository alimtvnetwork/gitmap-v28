# Data Model: CapturedError

> **Parent:** [Error Modal Reference](./00-overview.md)  
> **Version:** 2.2.0  
> **Updated:** 2026-03-31

---

The `CapturedError` interface is the **central data structure** for all error diagnostics. Every field must be populated correctly for the modal to render complete information.

```typescript
export interface CapturedError {
  // === Identity ===
  id: string;                              // Unique UUID
  code: string;                            // Error code (e.g., "E5001", "E9005")
  level: 'error' | 'warn' | 'info';       // Severity level
  message: string;                         // Primary error message
  details?: string;                        // Extended description
  createdAt: string;                       // ISO 8601 timestamp

  // === Frontend Location ===
  file?: string;                           // Source file where error originated
  line?: number;                           // Line number
  function?: string;                       // Function name
  stackTrace?: string;                     // Raw JS stack trace string
  parsedFrames?: StackFrame[];             // Parsed stack frames (file, line, function, isInternal)
  context?: Record<string, unknown>;       // Arbitrary context data

  // === API Request Context ===
  endpoint?: string;                       // API endpoint path
  method?: string;                         // HTTP method (GET, POST, etc.)
  requestBody?: unknown;                   // JSON request body sent
  responseStatus?: number;                 // HTTP response status code

  // === Trigger Context (Component → Action) ===
  triggerComponent?: string;               // React component name (e.g., "PluginCard")
  triggerAction?: string;                  // User action (e.g., "enable_clicked")
  invocationChain?: string[];              // Call chain: ["PluginsPage", "usePluginActions.enable", "api.post"]

  // === Backend Execution Logs ===
  backendLogs?: BackendLogEntry[];         // Session execution logs from Go
  backendStackTrace?: string;              // Go stack trace string
  siteUrl?: string;                        // Target WordPress site URL
  sessionId?: string;                      // Session ID for fetching deep diagnostics
  sessionType?: string;                    // "publish", "sync", "connection_test", etc.

  // === PHP/WordPress Error Details ===
  phpStackFrames?: PHPStackFrame[];        // Structured PHP frames (file, line, function, class)
  errorFile?: string;                      // PHP file where error occurred
  errorLine?: number;                      // PHP line number

  // === User Interaction Tracking ===
  uiClickPath?: ClickEvent[];             // Last N user clicks before error
  uiClickPathString?: string;             // Formatted click path for copy/display

  // === React Execution Logger ===
  executionLogs?: ExecutionLogEntry[];     // Function/component/effect/handler calls
  executionChain?: CallChain | null;       // Structured call chain
  executionLogsEnabled?: boolean;          // Whether debug mode was on
  executionLogsFormatted?: string;         // Pre-formatted chain for display

  // === Universal Response Envelope Fields ===
  requestedAt?: string;                    // Attributes.RequestedAt (Go endpoint path)
  requestDelegatedAt?: string;             // Attributes.RequestDelegatedAt (PHP endpoint URL)
  envelopeErrors?: EnvelopeErrors;         // Errors block from envelope
  envelopeMethodsStack?: EnvelopeMethodsStack; // MethodsStack block from envelope

  // NOTE: DelegatedRequestServer data is accessed via envelopeErrors.DelegatedRequestServer
  // There is no top-level delegatedRequestServer field — it lives inside the envelope Errors block.
}
```

---

## Supporting Types

```typescript
export interface StackFrame {
  function: string;
  file: string;
  line: number;
  column?: number;
  isInternal: boolean;    // true if from node_modules or browser internals
}

export interface BackendLogEntry {
  timestamp: string;
  level: 'debug' | 'info' | 'warn' | 'error';
  message: string;
  step?: string;          // Pipeline step (e.g., "backup", "upload", "activate")
  details?: Record<string, unknown>;
}

export interface PHPStackFrame {
  file?: string;
  fileBase?: string;      // Basename for compact display
  line?: number;
  function?: string;
  class?: string;
}

// From the Universal Response Envelope
export interface EnvelopeErrors {
  BackendMessage: string;
  DelegatedServiceErrorStack?: string[];  // PHP error lines (legacy)
  Backend?: string[];                      // Go stack trace lines
  Frontend?: string[];                     // Reserved for frontend injection
  DelegatedRequestServer?: DelegatedRequestServer; // NEW v2.0.0
}

// NEW v2.0.0 — Structured delegated server error details
export interface DelegatedRequestServer {
  DelegatedEndpoint: string;    // Exact URL of downstream endpoint
  Method: string;               // HTTP method (GET, POST, etc.)
  StatusCode: number;           // HTTP status code from downstream
  RequestBody?: unknown;        // What was sent to downstream (null for GET)
  Response?: unknown;           // Full response body from downstream
  StackTrace?: string[];        // Stack trace lines from downstream server
  AdditionalMessages?: string;  // Human-readable error context
}

export interface EnvelopeMethodsStack {
  Backend: Array<{
    Method: string;       // Go function name
    File: string;         // Go source file
    LineNumber: number;
  }>;
  Frontend: Array<{
    Method: string;
    File: string;
    LineNumber: number;
  }>;
}
```

---

*Data model — updated: 2026-03-31*
