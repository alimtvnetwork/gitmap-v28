# Error Modal — Frontend Specification

> **Version:** 2.1.0  
> **Updated:** 2026-03-09  
> **Status:** Active  
> **Location:** `src/components/errors/`  
> **Purpose:** Comprehensive specification for the Global Error Modal — how errors are captured, enriched, displayed, and exported across the React → Go → Delegated Server request chain.

---

## Table of Contents

1. [Architecture Overview](#1-architecture-overview)
2. [Data Model: CapturedError](#2-data-model-capturederror)
3. [Error Capture Pipeline](#3-error-capture-pipeline)
4. [Envelope Parsing & Enrichment](#4-envelope-parsing--enrichment)
5. [Modal Structure & Components](#5-modal-structure--components)
6. [Backend Section (Tabs)](#6-backend-section-tabs)
7. [Frontend Section (Tabs)](#7-frontend-section-tabs)
8. [Request Chain Visualization](#8-request-chain-visualization)
9. [Traversal Details](#9-traversal-details)
10. [Session Diagnostics Auto-Fetch](#10-session-diagnostics-auto-fetch)
11. [Error Report Generation](#11-error-report-generation)
12. [Error Queue Navigation](#12-error-queue-navigation)
13. [React Code Examples](#13-react-code-examples)
14. [File Reference](#14-file-reference)

---

## 1. Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         Error Capture Flow                              │
│                                                                         │
│  User Action ──▸ API Call ──▸ Go Backend ──▸ PHP (WordPress)            │
│       │              │              │              │                     │
│       │              │              │              ▼                     │
│       │              │              │     ┌──────────────────┐          │
│       │              │              │     │ PHP Error         │          │
│       │              │              │     │ - stackTrace      │          │
│       │              │              │     │ - stackTraceFrames│          │
│       │              │              │     │ - fatal-errors.log│          │
│       │              │              │     └────────┬─────────┘          │
│       │              │              │              │                     │
│       │              │              ▼              │                     │
│       │              │     ┌──────────────────┐    │                     │
│       │              │     │ Go Error Handler  │◀──┘                     │
│       │              │     │ - apperror.Wrap() │                         │
│       │              │     │ - session logger  │                         │
│       │              │     │ - envelope builder│                         │
│       │              │     └────────┬─────────┘                         │
│       │              │              │                                    │
│       │              ▼              │ Universal Response Envelope        │
│       │     ┌──────────────────┐    │                                    │
│       │     │ API Client       │◀───┘                                    │
│       │     │ - parseEnvelope()│                                         │
│       │     │ - extract errors │                                         │
│       │     └────────┬─────────┘                                        │
│       │              │                                                   │
│       ▼              ▼                                                   │
│  ┌──────────────────────────────────┐                                   │
│  │        Error Store (Zustand)      │                                   │
│  │  - buildCapturedError()           │                                   │
│  │  - commitErrorToStore()           │                                   │
│  │  - enrich with click path         │                                   │
│  │  - enrich with execution logs     │                                   │
│  └────────────────┬─────────────────┘                                   │
│                   │                                                      │
│                   ▼                                                      │
│  ┌──────────────────────────────────┐                                   │
│  │       Global Error Modal          │                                   │
│  │  ┌─────────┐ ┌─────────────────┐ │                                   │
│  │  │ Backend │ │    Frontend     │ │                                   │
│  │  │ Section │ │    Section      │ │                                   │
│  │  └─────────┘ └─────────────────┘ │                                   │
│  │  ┌─────────────────────────────┐ │                                   │
│  │  │ Download/Copy Actions       │ │                                   │
│  │  └─────────────────────────────┘ │                                   │
│  └──────────────────────────────────┘                                   │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Data Model: CapturedError

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
  endpoint?: string;                       // API endpoint path (e.g., "http://localhost:8080/api/v1/plugins/enable")
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

### Supporting Types

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

## 3. Error Capture Pipeline

### Step 1: API Client Detects Failure

When an API call fails, the client extracts envelope data:

```typescript
// src/lib/api/envelope.ts
export function parseEnvelope(response: unknown): ParsedEnvelope | null {
  if (!response || typeof response !== 'object') return null;
  const r = response as Record<string, unknown>;
  
  // Detect Universal Response Envelope (PascalCase)
  if (r.Status && r.Attributes && r.Results) {
    const status = r.Status as { IsSuccess: boolean; Code: number; Message: string };
    const attrs = r.Attributes as {
      RequestedAt?: string;
      RequestDelegatedAt?: string;
      SessionId?: string;
      HasAnyErrors?: boolean;
    };
    
    return {
      isSuccess: status.IsSuccess,
      code: status.Code,
      message: status.Message,
      requestedAt: attrs.RequestedAt,
      requestDelegatedAt: attrs.RequestDelegatedAt,
      sessionId: attrs.SessionId,
      hasErrors: attrs.HasAnyErrors ?? false,
      errors: r.Errors as EnvelopeErrors | null,
      methodsStack: r.MethodsStack as EnvelopeMethodsStack | null,
      results: r.Results as unknown[],
    };
  }

  return null;
}
```

### Step 2: Error Store Captures & Enriches

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

### Step 3: Modal Opens

```typescript
// Anywhere in the app:
const { openErrorModal } = useErrorStore();
openErrorModal(capturedError);
```

---

## 4. Envelope Parsing & Enrichment

The Universal Response Envelope provides six top-level blocks. The error modal consumes three of them for diagnostics:

### Errors Block → Backend Section (Overview + Stack tabs)

```json
{
  "Errors": {
    "BackendMessage": "Failed to fetch plugin details from remote site",
    "DelegatedServiceErrorStack": [
      "PHP Fatal error: Class 'PDO' not found in plugin-manager.php on line 42",
      "#0 endpoints.php(15): PluginManager->connect()",
      "#1 {main}"
    ],
    "Backend": [
      "site_handlers.go:327 handlers.EnableRemotePlugin",
      "service.go:1245 site.(*Service).EnableRemotePlugin"
    ],
    "Frontend": [],
    "DelegatedRequestServer": {
      "DelegatedEndpoint": "https://example.com/wp-json/riseup.../v1/snapshots/settings",
      "Method": "GET",
      "StatusCode": 403,
      "RequestBody": null,
      "Response": {
        "code": "rest_forbidden",
        "message": "This endpoint is disabled",
        "data": { "status": 403, "plugin_version": "1.54.0" }
      },
      "StackTrace": [
        "#0 riseup-asia-uploader.php(1098): FileLogger->error()",
        "#1 class-wp-hook.php(341): Plugin->enrichErrorResponse()"
      ],
      "AdditionalMessages": "Endpoint 'snapshots' is not enabled in plugin settings."
    }
  }
}
```

**Mapping:**
- `BackendMessage` → Overview tab red banner
- `DelegatedServiceErrorStack` → Stack tab (orange-themed PHP trace) + Traversal tab (legacy)
- `DelegatedRequestServer` → Stack tab (purple-themed delegated section) + Request tab (3rd hop) + Traversal tab (NEW v2.0.0)
- `DelegatedRequestServer.AdditionalMessages` → Overview tab info banner
- `DelegatedRequestServer.Response` → Request tab delegated response JSON viewer
- `DelegatedRequestServer.StackTrace` → Stack tab delegated server stack trace
- `Backend` → Stack tab (Go trace)

### MethodsStack Block → Execution + Traversal tabs

```json
{
  "MethodsStack": {
    "Backend": [
      { "Method": "handlers.EnableRemotePlugin", "File": "site_handlers.go", "LineNumber": 327 },
      { "Method": "site.(*Service).EnableRemotePlugin", "File": "service.go", "LineNumber": 1245 },
      { "Method": "wordpress.(*Client).doRequest", "File": "uploader.go", "LineNumber": 350 }
    ],
    "Frontend": []
  }
}
```

**Mapping:** Rendered as a sortable table with `#`, `Method`, `File`, `Line` columns.

### Attributes Block → Session ID propagation

```json
{
  "Attributes": {
    "RequestedAt": "http://localhost:8080/api/v1/plugins/enable",
    "RequestDelegatedAt": "https://example.com/wp-json/riseup-asia-uploader/v1/enable",
    "SessionId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "HasAnyErrors": true
  }
}
```

**Mapping:**
- `RequestedAt` + `RequestDelegatedAt` → Traversal tab endpoint flow
- `SessionId` → Session tab auto-fetch trigger

---

## 5. Modal Structure & Components

### Component Hierarchy

```
GlobalErrorModal.tsx
├── DialogHeader (error code, timestamp, queue navigation)
├── Section Toggle (Backend / Frontend buttons)
├── ScrollArea
│   ├── BackendSection.tsx (when activeSection === "backend")
│   │   ├── Tabs
│   │   │   ├── Overview (error message, request, timing, badges)
│   │   │   ├── Log (error.log.txt content)
│   │   │   ├── Execution (Go call chain table + session logs)
│   │   │   ├── Stack (Go + PHP stack traces, session diagnostics)
│   │   │   ├── Session (SessionLogsTab — logs, request, response, stack trace)
│   │   │   ├── Request (RequestDetails — request chain visualization)
│   │   │   └── Traversal (TraversalDetails — endpoint flow + methods stack)
│   │   └── (Internal sub-components: OverviewContent, ErrorLogContent, etc.)
│   │
│   └── FrontendSection.tsx (when activeSection === "frontend")
│       ├── Tabs
│       │   ├── Overview (trigger context, message, call chain, click path)
│       │   ├── Stack (parsed/raw JS stack, React execution chain)
│       │   ├── Context (full error context JSON)
│       │   └── Fixes (suggested fixes by error code)
│       └── (Internal sub-components)
│
├── DialogFooter
│   ├── DownloadDropdown (ZIP bundle, error.log.txt, log.txt, report.md)
│   ├── Close button
│   └── CopyDropdown (Split Button: main click = compact report, chevron dropdown = full report, with backend logs, error.log.txt, log.txt)
│
└── ErrorModalTypes.ts (shared types: PHPStackFrame, AppInfo, SectionCommonProps)
```

### Visual Layout Diagrams

#### Full Modal Layout (Desktop: 95vw × 95vh)

```
┌─────────────────────────────────────────────────────────────────────┐
│ ┌─ DialogHeader ──────────────────────────────────────────────────┐ │
│ │ [E5001]  Failed to enable plugin   2026-02-09 14:32:01         │ │
│ │                                          [◀ 1/3 ▶] [Copy All] │ │
│ └─────────────────────────────────────────────────────────────────┘ │
│ ┌─ Section Toggle ────────────────────────────────────────────────┐ │
│ │  [ ● Backend ]  [ ○ Frontend ]                                 │ │
│ └─────────────────────────────────────────────────────────────────┘ │
│ ┌─ ScrollArea (flex-1) ──────────────────────────────────────────┐ │
│ │ ┌─ Tab Bar ──────────────────────────────────────────────────┐ │ │
│ │ │ Overview │ Log │ Execution │ Stack │ Session │ Request │ Traversal │
│ │ └────────────────────────────────────────────────────────────┘ │ │
│ │                                                                │ │
│ │  ┌─ Active Tab Content ─────────────────────────────────────┐  │ │
│ │  │                                                          │  │ │
│ │  │  (Tab-specific content rendered here)                    │  │ │
│ │  │                                                          │  │ │
│ │  └──────────────────────────────────────────────────────────┘  │ │
│ └────────────────────────────────────────────────────────────────┘ │
│ ┌─ DialogFooter ─────────────────────────────────────────────────┐ │
│ │  [▼ Download]                              [Close] [▼ Copy]   │ │
│ └─────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
```

#### Backend Section — Overview Tab

```
┌──────────────────────────────────────────────────────────────────┐
│  ┌─ Error Banner (red) ────────────────────────────────────────┐ │
│  │ ⚠ Backend Error: Failed to fetch plugin details from site  │ │
│  └─────────────────────────────────────────────────────────────┘ │
│                                                                  │
│  ┌─ Delegated Info Banner (blue, NEW v2.0.0) ──────────────────┐ │
│  │ ℹ Endpoint 'snapshots' is not enabled in plugin settings.   │ │
│  │   (from DelegatedRequestServer.AdditionalMessages)          │ │
│  └─────────────────────────────────────────────────────────────┘ │
│                                                                  │
│  ┌─ Request Info ──────────────────────────────────────────────┐ │
│  │  Method: POST   Endpoint: /api/v1/plugins/enable            │ │
│  │  Status: 500    Site: https://example.com                   │ │
│  └─────────────────────────────────────────────────────────────┘ │
│                                                                  │
│  ┌─ Timing ───────────────────────────────────────────────────┐  │
│  │  Requested At:           /api/v1/plugins/enable             │  │
│  │  Delegated At:           https://site.com/wp-json/...       │  │
│  └─────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ Availability Badges ──────────────────────────────────────┐  │
│  │  [✓ Session] [✓ Stack Traces] [✓ Delegated Info] [✓ Exec] │  │
│  └─────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────┘
```

#### Backend Section — Stack Tab

```
┌──────────────────────────────────────────────────────────────────┐
│  ┌─ Go Backend Stack (blue-themed) ───────────────────────────┐  │
│  │  site_handlers.go:327  handlers.EnableRemotePlugin          │  │
│  │  service.go:1245       site.(*Service).EnableRemotePlugin   │  │
│  └─────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ Delegated Server Stack (purple-themed, NEW v2.0.0) ───────┐  │
│  │  ┌─ Header ─────────────────────────────────────────────┐   │  │
│  │  │ 🟣 Delegated Server  GET  403                        │   │  │
│  │  │ https://example.com/wp-json/riseup.../snapshots/...  │   │  │
│  │  └──────────────────────────────────────────────────────┘   │  │
│  │  Stack Trace:                                               │  │
│  │  #0 riseup-asia-uploader.php(1098): Logger->error()         │  │
│  │  #1 class-wp-hook.php(341): Plugin->enrichError()          │  │
│  │  #2 plugin.php(205): WP_Hook->apply_filters()               │  │
│  │  Response:                                                  │  │
│  │  ▸ { "code": "rest_forbidden", "message": "...", ... }     │  │
│  └─────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ PHP Delegated Stack (orange-themed, legacy) ───────────────┐  │
│  │  PHP Fatal error: Class 'PDO' not found in plugin-mgr.php  │  │
│  │  #0 endpoints.php(15): PluginManager->connect()             │  │
│  │  #1 {main}                                                  │  │
│  └─────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ PHP Structured Frames (table) ────────────────────────────┐  │
│  │  #  │ Function                    │ File              │ Line│  │
│  │  0  │ PluginManager::connect()    │ plugin-mgr.php    │ 42  │  │
│  │  1  │ handle_enable()             │ endpoints.php     │ 15  │  │
│  └─────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ Session Diagnostics (auto-fetched) ───────────────────────┐  │
│  │  Go frames: 3 │ PHP frames: 2 │ stacktrace.txt: available  │  │
│  └─────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────┘
```

#### Backend Section — Request Tab (Chain Visualization)

```
┌──────────────────────────────────────────────────────────────────┐
│  ┌─ Node 1: React → Go ──────────────────────────────────────┐  │
│  │  🔵 [React → Go]  [POST]  [500]                           │  │
│  │  /api/v1/plugins/enable                                    │  │
│  │  ▸ Request Body: { "slug": "my-plugin", "SiteId": 1 }     │  │
│  └──────────┬─────────────────────────────────────────────────┘  │
│             │ (vertical connector line)                           │
│  ┌──────────┴─────────────────────────────────────────────────┐  │
│  │  🟠 [Go → Delegated]  [GET]  [403]                        │  │
│  │  https://example.com/wp-json/riseup.../v1/snapshots/...    │  │
│  │  ▸ Request Body: (none — GET)                              │  │
│  └──────────┬─────────────────────────────────────────────────┘  │
│             │ (vertical connector line)                           │
│  ┌──────────┴─────────────────────────────────────────────────┐  │
│  │  🟣 [Delegated Response]  (NEW v2.0.0)                     │  │
│  │  ▸ Response: { "code": "rest_forbidden", "message": ... }  │  │
│  │  ▸ Stack Trace: #0 riseup-asia-uploader.php(1098)...       │  │
│  │  ▸ Additional: Endpoint 'snapshots' is not enabled...      │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ Environment ──────────────────────────────────────────────┐  │
│  │  API Base: http://localhost:8080   VITE_API_URL: ...        │  │
│  └────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────┘
```

#### Backend Section — Traversal Tab

```
┌──────────────────────────────────────────────────────────────────┐
│  ┌─ Endpoint Flow (3-hop, NEW v2.0.0) ─────────────────────────┐ │
│  │  [React] http://localhost:8080                               │ │
│  │    ──▸                                                      │ │
│  │  [Go] /api/v1/sites/1/snapshots/settings                    │ │
│  │    ──▸                                                      │ │
│  │  [Delegated] https://site.com/wp-json/riseup.../settings    │ │
│  │              GET → 403                                      │ │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ Methods Stack (table) ────────────────────────────────────┐  │
│  │  #  │ Method                          │ File            │ Ln│  │
│  │  1  │ handlers.handleSiteActionById   │ handler_factory │107│  │
│  │  2  │ api.SessionLogging              │ session_log     │107│  │
│  │  3  │ api.Recovery                    │ middleware      │245│  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ Delegated Server Details (purple, NEW v2.0.0) ─────────────┐ │
│  │  Endpoint: https://site.com/wp-json/riseup.../settings      │ │
│  │  Method: GET │ Status: 403                                  │ │
│  │  Stack Trace:                                                │ │
│  │    #0 riseup-asia-uploader.php(1098): Logger->error()       │ │
│  │    #1 class-wp-hook.php(341): enrichErrorResponse()       │ │
│  │  Additional: Endpoint not enabled in plugin settings        │ │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ Delegated Service Error Stack (orange, legacy) ────────────┐ │
│  │  PHP Fatal error: Class 'PDO' not found...                  │ │
│  │  #0 endpoints.php(15): PluginManager->connect()             │ │
│  └────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────┘
```

#### Frontend Section — Overview Tab

```
┌──────────────────────────────────────────────────────────────────┐
│  ┌─ Trigger Context ──────────────────────────────────────────┐  │
│  │  Component: PluginCard  →  Action: enable_clicked           │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ Message ──────────────────────────────────────────────────┐  │
│  │  Failed to enable plugin "my-plugin" on site                │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ Call Chain ───────────────────────────────────────────────┐  │
│  │  PluginsPage                                                │  │
│  │    └─ usePluginActions.enable                               │  │
│  │        └─ api.post("http://localhost:8080/api/v1/plugins/enable")                │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ User Interaction Path (last 10 clicks) ───────────────────┐  │
│  │  14:31:55  PluginsPage     "Plugins" tab       click  /    │  │
│  │  14:31:58  PluginCard      "Enable" button     click  /    │  │
│  │  14:32:01  PluginCard      "Confirm" button    click  /    │  │
│  └────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────┘
```

#### Frontend Section — Stack Tab

```
┌──────────────────────────────────────────────────────────────────┐
│  [● Parsed] [○ Raw]                    [□ Show internal frames] │
│                                                                  │
│  ┌─ Parsed Stack Frames (table) ──────────────────────────────┐  │
│  │  #  │ Function              │ File                │ Line   │  │
│  │  0  │ enablePlugin          │ usePluginActions.ts  │ 45    │  │
│  │  1  │ handleClick           │ PluginCard.tsx       │ 112   │  │
│  │  2  │ callCallback          │ react-dom.js         │ 3942  │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ React Execution Chain ────────────────────────────────────┐  │
│  │  [render] PluginsPage                           14:31:50   │  │
│  │  [effect] usePluginActions                      14:31:51   │  │
│  │  [handler] enablePlugin                         14:32:01   │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  Error Location: usePluginActions.ts:45 in enablePlugin()        │
└──────────────────────────────────────────────────────────────────┘
```

#### DialogFooter — Action Menus

The **Copy** button uses a **Split Button** pattern: the main button area copies the **Compact Report** instantly (no API call), while the chevron arrow opens a dropdown with all copy options.

```
┌──────────────────────────────────────────────────────────────────┐
│  [▼ Download]                       [Close]  [ Copy ][▼]        │
│  ┌─────────────────┐                    ┌──────────────────────┐│
│  │ Full Bundle (ZIP)│        main click: │ (copies compact)     ││
│  │ error.log.txt    │        chevron ▼:  │ Compact Report       ││
│  │ log.txt          │                    │ Full Report          ││
│  │ Report (.md)     │                    │ With Backend Logs    ││
│  └─────────────────┘                    │ error.log.txt        ││
│                                         │ log.txt              ││
│                                         └──────────────────────┘│
└──────────────────────────────────────────────────────────────────┘
```

The **ErrorDetailModal** (standalone viewer used on E2E Tests page) uses the same Split Button pattern via an `errorLogAdapter.ts` bridge that maps the `ErrorLog` API shape to a `CapturedError`-compatible object for report generation.

### Full-Screen Layout

The modal uses a **full-screen** layout on mobile and **95vw × 95vh** on desktop:

```tsx
<DialogContent className={cn(
  "flex flex-col p-0 gap-0 overflow-hidden",
  "w-full h-full max-w-full max-h-full rounded-none",          // Mobile: full screen
  "sm:max-w-[95vw] sm:w-[95vw] sm:max-h-[95vh] sm:h-[95vh] sm:rounded-lg",  // Desktop
  "lg:max-w-6xl"
)}>
```

---

## 6. Backend Section (Tabs)

### Tab: Overview

Displays the primary error information:
- Error message with code badge and timestamp
- Target site URL (if WordPress operation)
- API request method and endpoint
- **Backend Error** banner (from `envelopeErrors.BackendMessage`) — red-themed
- Timing section (`requestedAt`, `requestDelegatedAt`)
- Availability badges (session, stack traces, execution logs)

### Tab: Log

Fetches and displays `error.log.txt` from the backend:
- Auto-fetched when the modal opens on the Backend section
- Refresh, Copy, and Download buttons
- `ScrollArea` with monospace font, 400px height

### Tab: Execution

Two sections:
1. **Go Call Chain** — Table from `envelopeMethodsStack.Backend`
2. **Session Execution Logs** — Timeline from `backendLogs[]` with level-based coloring and step labels

### Tab: Stack

Multi-source stack traces:
1. **Go Backend** — From `envelopeErrors.Backend` (blue-themed)
2. **PHP Delegated** — From `envelopeErrors.DelegatedServiceErrorStack` (orange-themed)
3. **PHP Structured Frames** — Table from `phpStackFrames[]` (file, line, class::function)
4. **Session Diagnostics** — Auto-fetched Go and PHP frames from session API
5. **PHP stacktrace.txt** — Raw backtrace from session diagnostics

### Tab: Session

Full session-level diagnostics (only shown when `sessionId` exists):
- Sub-tabs: **Logs**, **Request**, **Response**, **Stack Trace**
- Fetches from `GET /api/v1/sessions/{id}/logs` and `GET /api/v1/sessions/{id}/diagnostics`
- Log rendering with color-coded levels and stage headers

### Tab: Request

**Request Chain Visualization** — A two-node vertical chain:

```
┌──────────────────────────────────┐
│ 🔵 React → Go                   │
│ POST /api/v1/plugins/enable      │
│ Status: 500                      │
│ Request Body: { ... }            │
└──────────┬───────────────────────┘
           │
┌──────────┴───────────────────────┐
│ 🟠 Go → PHP                     │
│ https://site.com/wp-json/...     │
│ Status: 500                      │
│ PHP Response Body: { ... }       │
│ PHP Error Stack: ...             │
└──────────────────────────────────┘
```

Plus environment diagnostics: API base, VITE_API_URL, resolved origin, UI origin.

### Tab: Traversal

- **Endpoint Flow** — Go endpoint → PHP endpoint with badges
- **Methods Stack** — Table: #, Method, File, Line
- **Delegated Server Details** — Purple-themed section with endpoint, method, status, stack trace, response JSON (NEW v2.0.0)
- **Delegated Service Error Stack** — Orange-themed PHP error lines (legacy)
- **Backend Trace** — Go stack trace lines

### Sections Rendered (Priority Order):

1. **Endpoint Flow** — React → Go → Delegated badges with full URLs (3-hop when DelegatedRequestServer present)
2. **Methods Stack Table** — Go call chain from envelope
3. **Delegated Server Details** — Purple-themed block with endpoint, method, status, stack trace, response (NEW v2.0.0)
4. **Delegated Service Error Stack** — Orange-themed PHP error lines (legacy, ScrollArea)
5. **Backend Trace** — Go stack trace lines (ScrollArea)

---

## 7. Frontend Section (Tabs)

### Tab: Overview

- **Trigger Context** — Component → Action badge (e.g., `PluginCard → enable_clicked`)
- **Message** and **Details**
- **Call Chain** — Indented tree visualization of `invocationChain[]`
- **User Interaction Path** — Last 10 clicks with component name, text, action type, and route

### Tab: Stack

- **Parsed/Raw toggle** — Switch between parsed frame table and raw stack string
- **Show internal frames** — Toggle for node_modules frames
- **React Execution Chain** — From `useExecutionLogger` (component renders, effect triggers, handler calls)
- **Error Location** — File, line, function

### Tab: Context

- Full JSON context (`error.context`) with syntax highlighting via `JsonHighlighter`

### Tab: Fixes

- Suggested fixes keyed by error code (e.g., E1001 → check backend port, E5001 → check plugin activation)

---

## 8. Request Chain Visualization

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

### Sample Code: Rendering the 3-Hop Request Chain

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

## 9. Traversal Details

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
const hasDelegatedServer = error.delegatedRequestServer || error.envelopeErrors?.DelegatedRequestServer; // NEW v2.0.0
```

### Sections Rendered:

1. **Endpoint Flow** — 3-hop badges: React → Go → Delegated (with method + status for delegated hop)
2. **Methods Stack Table** — Go call chain from envelope
3. **Delegated Server Details** — Purple-themed block (NEW v2.0.0): endpoint, method, status, stack trace, response JSON, additional messages
4. **Delegated Service Error Stack** — Orange-themed PHP error lines (legacy, ScrollArea)
5. **Backend Trace** — Go stack trace lines (ScrollArea)

### Sample Code: Rendering Delegated Server Details in Traversal

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

## 10. Session Diagnostics Auto-Fetch

When `sessionId` is present, the `BackendSection` automatically fetches deep diagnostics:

```tsx
// src/hooks/useSessionDiagnostics.ts
export function useSessionDiagnostics(sessionId?: string) {
  const [diagnostics, setDiagnostics] = useState<SessionDiagnostics | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!sessionId) return;
    setLoading(true);
    api.getSessionDiagnostics(sessionId)
      .then(resp => {
        if (resp.success && resp.data) setDiagnostics(resp.data);
      })
      .finally(() => setLoading(false));
  }, [sessionId]);

  return { diagnostics, loading };
}
```

### SessionDiagnostics Shape

```typescript
interface SessionDiagnostics {
  request?: {
    method: string;
    url: string;
    body?: Record<string, unknown>;
  };
  response?: {
    statusCode: number;
    requestUrl: string;   // The full PHP endpoint that processed the request
    body?: unknown;
  };
  stackTrace?: {
    golang?: SessionStackFrame[];
    php?: SessionStackFrame[];
  };
  phpStackTraceLog?: string;  // Raw stacktrace.txt content
}

interface SessionStackFrame {
  function: string;
  file?: string;
  line?: number;
  class?: string;
}
```

---

## 11. Error Report Generation

The `errorReportGenerator.ts` produces two report formats from any `CapturedError`:

### Compact Report (Default)

The **Compact Error Report** is the default copy format. It is designed for high-signal diagnostic sharing — stripped of noise, built entirely from client-side memory (no API calls). The main "Copy" button uses this format.

```typescript
export function generateCompactReport(error: CapturedError, app?: AppInfo): string {
  // Sections generated:
  // 1. App metadata (name, version)
  // 2. Error code + level
  // 3. Page (route + component)
  // 4. User Interaction (arrow-joined click path)
  // 5. Trigger Context (component, action, source)
  // 6. Message
  // 7. Request (method, endpoint, status)
  // 8. Frontend Execution Chain (stripped: no timestamps, no base URLs)
  // 9. Context JSON
  // 10. Frontend Stack Trace
  // 11. Backend error.log.txt (reconstructed from CapturedError memory)
  return markdownReport;
}
```

**Stripping rules for execution chain:**
- Timestamps (e.g., `[12:58:22 AM] ⬡`) are removed
- Base API URLs (e.g., `http://localhost:8080/api/v1`) are stripped to relative paths (e.g., `/sites`)
- Result: clean, scannable list like `GET /sites`, `POST /error-history`

### Full Report

```typescript
export function generateErrorReport(error: CapturedError, app?: AppInfo): string {
  // 17 sections including full URLs, backend logs, PHP traces, etc.
  return markdownReport;
}
```

### ErrorDetailModal Adapter

The standalone `ErrorDetailModal` (used on E2E Tests page) receives `ErrorLog` objects from the backend API, not `CapturedError`. The `errorLogAdapter.ts` bridges this gap:

```typescript
// src/components/errors/errorLogAdapter.ts
export function errorLogToCapturedError(error: ErrorLog): CapturedError;
```

This allows both modals to use `generateCompactReport` and `generateErrorReport` with the same interface.

### Download/Copy Options

**Copy (Split Button):**

| Action | Trigger | Content |
|--------|---------|---------|
| **Copy (main click)** | Primary button area | Compact Report (instant, from memory) |
| **Copy Compact Report** | Chevron dropdown | Same as main click |
| **Copy Full Report** | Chevron dropdown | Full 17-section Markdown report |
| **Copy with Backend Logs** | Chevron dropdown | Full report + error.log.txt (API fetch) |
| **Copy error.log.txt** | Chevron dropdown | Raw error log content (API fetch) |
| **Copy log.txt** | Chevron dropdown | Raw full log content (API fetch) |

**Download:**

| Action | Content |
|--------|---------|
| **Full Bundle (ZIP)** | Markdown report + logs, sent to `/api/v1/errors/bundle` |
| **error.log.txt** | Backend error log file |
| **log.txt** | Backend full diagnostic log |
| **Report (.md)** | Markdown error report |

---

## 12. Error Queue Navigation

Multiple concurrent errors are queued and navigable:

```tsx
// In ErrorStore:
errorQueue: CapturedError[];
currentQueueIndex: number;
navigateQueue: (direction: 'prev' | 'next') => void;
getQueuedErrorsMarkdown: () => string;  // All errors as one Markdown doc

// In GlobalErrorModal header:
{hasMultipleErrors && (
  <div className="flex items-center gap-1">
    <Button onClick={() => navigateQueue('prev')}>
      <ChevronLeft />
    </Button>
    <Badge>{currentQueueIndex + 1}/{errorQueue.length}</Badge>
    <Button onClick={() => navigateQueue('next')}>
      <ChevronRight />
    </Button>
    <Button onClick={copyAllErrors}>
      <CopyPlus /> All
    </Button>
  </div>
)}
```

---

## 13. React Code Examples

### Example: Capturing an API Error with Full Context

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

### Example: Opening the Error Modal Manually

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

### Example: Using the Error Boundary

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

### Example: Consuming Envelope Errors in a Hook

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

### Example: Rendering PHP Stack Frames

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

## 14. File Reference

| File | Lines | Purpose |
|------|-------|---------|
| `src/stores/errorStore.ts` | 674 | Central error store (Zustand): CapturedError type, capture pipeline, queue management |
| `src/components/errors/GlobalErrorModal.tsx` | 239 | Modal shell: header, section toggle, footer actions |
| `src/components/errors/BackendSection.tsx` | 721 | Backend tabs: Overview, Log, Execution, Stack, Session, Request, Traversal |
| `src/components/errors/FrontendSection.tsx` | 328 | Frontend tabs: Overview, Stack, Context, Fixes |
| `src/components/errors/RequestDetails.tsx` | 176 | Request chain visualization (React → Go → Delegated, 3-hop) |
| `src/components/errors/TraversalDetails.tsx` | 149 | Envelope traversal: endpoint flow, methods stack, delegated server details, error stacks |
| `src/components/errors/SessionLogsTab.tsx` | 443 | Session diagnostics: logs, request, response, stack traces |
| `src/components/errors/ErrorModalActions.tsx` | 194 | Download and Copy dropdown menus |
| `src/components/errors/errorReportGenerator.ts` | 185 | Pure function: Markdown report generation (compact + full) + suggested fixes |
| `src/components/errors/errorLogAdapter.ts` | 22 | Maps backend `ErrorLog` → `CapturedError` for ErrorDetailModal report generation |
| `src/components/errors/ErrorModalTypes.ts` | 26 | Shared types: PHPStackFrame, AppInfo, SectionCommonProps, DelegatedRequestServer |
| `src/components/errors/ErrorDetailModal.tsx` | — | Standalone error detail viewer (Split Button copy + DownloadDropdown) |
| `src/components/errors/ErrorHistoryDrawer.tsx` | — | Error history browser drawer |
| `src/components/errors/ErrorQueueBadge.tsx` | — | Error queue indicator badge |
| `src/components/errors/AppErrorBoundary.tsx` | — | React error boundary wrapping the app |
| `src/hooks/useSessionDiagnostics.ts` | — | Hook for auto-fetching session-level diagnostics |
| `src/hooks/useClickTracker.ts` | — | Click path tracking for error context |
| `src/hooks/useExecutionLogger.ts` | — | React execution logger (debug mode) |
| `src/lib/api/envelope.ts` | — | Envelope parsing and error extraction (incl. DelegatedRequestServer) |

---

## Cross-References

- [Error Handling Cross-Stack Spec](../01-error-handling-reference.md) — PHP, Go, frontend error chain + DelegatedRequestServer flow
- [Copy Format Samples](./01-copy-formats.md) — Complete samples for all copy/export formats
- [React Components Reference](./02-react-components.md) — Portable React code for rebuilding the modal
- [Response Envelope Schema](../05-response-envelope/envelope.schema.json) — JSON Schema for envelope (incl. DelegatedRequestServer)
- [Envelope Configurability](../05-response-envelope/01-adr.md) — DelegatedRequestServer presence rules
- [Session-Based Logging](../07-logging-and-diagnostics/02-session-based-logging.md) — Backend session system
- [React Execution Logger](../07-logging-and-diagnostics/01-react-execution-logger.md) — Frontend debug logger
- [TypeScript Standards](../../../02-coding-guidelines/02-typescript/08-typescript-standards-reference.md) — Type safety rules

---

*Error Modal specification v2.1.0 — updated: 2026-02-17*
