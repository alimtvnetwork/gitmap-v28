# Specification: React Execution Logger

**Version:** 3.2.0  
**Created:** 2026-03-09  
**Status:** Implemented

---

## 1. Executive Summary

The React Execution Logger provides deep visibility into frontend execution by tracking function calls, component renders, effects, and handlers. When an error occurs, this creates a complete "call chain" showing exactly what code path led to the failure.

---

## 2. Problem Statement

JavaScript stack traces show where an error occurred but not the full context of how we got there. React's component model and hooks pattern obscure the actual execution flow. Developers need:

1. Component render sequence
2. Effect trigger chain
3. Handler invocation context
4. API call timing
5. Parent-child function relationships

---

## 3. Requirements

### 3.1 Functional Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| F1 | Track function executions with arguments | MUST |
| F2 | Track component renders with props | MUST |
| F3 | Track useEffect triggers with dependencies | MUST |
| F4 | Track event handlers with component context | MUST |
| F5 | Track API calls with method/endpoint | MUST |
| F6 | Build parent-child call relationships | MUST |
| F7 | Format as human-readable chain | MUST |
| F8 | Toggle via debug mode setting | MUST |
| F9 | Zero overhead when disabled | MUST |
| F10 | Integrate with error capture | MUST |

### 3.2 Non-Functional Requirements

| ID | Requirement | Target |
|----|-------------|--------|
| NF1 | Memory overhead (enabled) | < 1MB |
| NF2 | Memory overhead (disabled) | 0 |
| NF3 | CPU overhead per log | < 0.1ms |
| NF4 | Max entries retained | 100 |

---

## 4. Architecture

### 4.1 Component Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                     React Application                            │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐         │
│  │Component │  │Component │  │ Effect   │  │ Handler  │         │
│  │ Render   │  │ Render   │  │ Trigger  │  │ Invoke   │         │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘         │
│       │             │             │             │                │
│       └─────────────┴─────────────┴─────────────┘                │
│                              │                                   │
│                              ▼                                   │
│       ┌──────────────────────────────────────────────────────┐  │
│       │              useExecutionLogger Hook                  │  │
│       │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  │  │
│       │  │logFunc  │  │logComp  │  │logEffect│  │logHandler│ │  │
│       │  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘  │  │
│       │       └────────────┴────────────┴────────────┘       │  │
│       │                         │                             │  │
│       │                         ▼                             │  │
│       │  ┌───────────────────────────────────────────────┐   │  │
│       │  │            Zustand Store                       │   │  │
│       │  │  entries: ExecutionLogEntry[]                  │   │  │
│       │  │  callStack: string[]                           │   │  │
│       │  └───────────────────────────────────────────────┘   │  │
│       └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Error Store                                 │
│  captureError() ──► includes executionChain from logger         │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Global Error Modal                             │
│  Stack Tab ► Frontend Sub-tab ► React Execution Chain           │
└─────────────────────────────────────────────────────────────────┘
```

### 4.2 Call Stack Tracking

```
Time ───────────────────────────────────────────────────────────►

logComponent('PluginList')
    │
    ├── logEffect('loadPlugins')
    │       │
    │       ├── logApiCall('GET', '/api/v1/plugins')
    │       │
    │       └── logFunction('processResponse')
    │
    └── logHandler('onClick', 'RefreshButton')
            │
            └── logFunction('refreshData')
                    │
                    └── [ERROR OCCURS HERE]

Result Chain:
1. [component] PluginList rendered
2. [effect] loadPlugins triggered
3. [api] GET /api/v1/plugins
4. [function] processResponse called
5. [handler] onClick in RefreshButton
6. [function] refreshData called
7. [ERROR] TypeError: Cannot read property 'x' of undefined
```

---

## 5. Data Model

### 5.1 ExecutionLogEntry

```typescript
interface ExecutionLogEntry {
  // Identity
  id: string;           // Unique ID (nanoid)
  parentId?: string;    // Parent entry for call chain
  
  // Classification
  type: 'function' | 'component' | 'effect' | 'handler' | 'api';
  name: string;         // Function/component name
  context?: string;     // Additional context (e.g., component for handler)
  
  // Data
  args?: unknown[];     // Function arguments (optional)
  result?: unknown;     // Return value (optional)
  error?: string;       // Error if threw
  
  // Timing
  timestamp: number;    // Date.now() at entry
  duration?: number;    // Execution time in ms
}
```

### 5.2 Store State

```typescript
interface ExecutionLoggerState {
  entries: ExecutionLogEntry[];
  enabled: boolean;
  callStack: string[];  // Stack of entry IDs
  maxEntries: number;   // Rolling buffer limit
  
  // Actions
  logFunction: (name: string, args?: unknown[]) => string;
  logComponent: (name: string, props?: object) => void;
  logEffect: (name: string, deps?: unknown[]) => void;
  logHandler: (name: string, context?: string) => void;
  logApiCall: (method: string, endpoint: string) => void;
  endFunction: (id: string, result?: unknown, error?: string) => void;
  
  // Getters
  getFormattedChain: () => string;
  getRecentEntries: (count: number) => ExecutionLogEntry[];
  clear: () => void;
  setEnabled: (enabled: boolean) => void;
}
```

---

## 6. API Reference

### 6.1 logFunction

Track function execution with optional arguments.

```typescript
function processItems(items: Item[]): ProcessedItem[] {
  const id = logFunction('processItems', [items.length]);
  try {
    const result = items.map(transform);
    endFunction(id, result);

    return result;
  } catch (e) {
    endFunction(id, undefined, e.message);

    throw e;
  }
}
```

### 6.2 logComponent

Track component render.

```typescript
function UserCard({ user }: Props) {
  logComponent('UserCard', { userId: user.id });

  return <div>{user.name}</div>;
}
```

### 6.3 logEffect

Track effect trigger with dependencies.

```typescript
useEffect(() => {
  logEffect('fetchUserData', [userId]);
  fetchUser(userId);
}, [userId]);
```

### 6.4 logHandler

Track event handler invocation.

```typescript
<Button onClick={() => {
  logHandler('onClick', 'DeleteButton');
  handleDelete();
}}>
  Delete
</Button>
```

### 6.5 logApiCall

Track API request.

```typescript
async function getPlugins() {
  logApiCall('GET', '/api/v1/plugins');

  return await api.get('/api/v1/plugins');
}
```

### 6.6 getFormattedChain

Get human-readable execution chain.

```typescript
const chain = useExecutionLogger.getState().getFormattedChain();
console.log(chain);
// Output:
// 1. [component] PluginList rendered
// 2. [effect] loadPlugins triggered (deps: [1])
// 3. [api] GET /api/v1/plugins
// 4. [function] handleResponse called with [Object]
```

---

## 7. Integration Points

### 7.1 Error Store Integration

```typescript
// In errorStore.ts
captureError: (input) => {
  const logger = useExecutionLogger.getState();
  const executionChain = logger.enabled 
    ? logger.getFormattedChain() 
    : undefined;
  
  const error: CapturedError = {
    ...input,
    executionChain,
    executionEntries: logger.getRecentEntries(20),
  };
  
  // Store error...
}
```

### 7.2 Global Error Modal

```tsx
// In GlobalErrorModal.tsx - Stack Tab - Frontend Sub-tab
{error.executionChain ? (
  <div className="space-y-2">
    <div className="flex justify-between">
      <h4>React Execution Chain</h4>
      <CopyButton text={error.executionChain} />
    </div>
    <pre className="bg-muted p-3 rounded text-xs">
      {error.executionChain}
    </pre>
  </div>
) : (
  <p className="text-muted-foreground">
    Enable Debug Mode in settings to capture execution chain.
  </p>
)}
```

---

## 8. Configuration

### 8.1 Enable/Disable

```typescript
// Via settings
const debugMode = useSettings(s => s.debugMode);
useEffect(() => {
  useExecutionLogger.getState().setEnabled(debugMode);
}, [debugMode]);
```

### 8.2 Config Reference

```json
{
  "logging": {
    "FrontendDebugMode": false
  }
}
```

---

## 9. Performance

### 9.1 When Disabled

All log methods are immediate no-ops:

```typescript
logFunction: (name, args) => {
  if (!get().enabled) return '';
  // ... actual logging
}
```

### 9.2 When Enabled

- Rolling buffer limits memory (default 100 entries)
- No deep cloning of arguments
- Lazy formatting (only on getFormattedChain)
- ID generation uses fast nanoid

### 9.3 Benchmarks

| Operation | Time (enabled) | Time (disabled) |
|-----------|----------------|-----------------|
| logFunction | 0.02ms | 0.001ms |
| logComponent | 0.01ms | 0.001ms |
| getFormattedChain | 0.5ms | 0.001ms |

---

## 10. Best Practices

### 10.1 DO

```typescript
// Log at function boundaries
function processData(items) {
  const id = logFunction('processData', [items.length]);
  // ...
  endFunction(id, result);
}

// Log handlers with context
<Button onClick={() => {
  logHandler('onClick', 'SaveDraftButton');
  saveDraft();
}}>

// Log effects with meaningful deps
useEffect(() => {
  logEffect('syncWithServer', [lastModified]);
}, [lastModified]);
```

### 10.2 DON'T

```typescript
// Don't log in tight loops
items.forEach(item => {
  logFunction('processItem'); // Too many entries!
});

// Don't log sensitive data
logFunction('login', [username, password]); // Exposes password!

// Don't log in render body (use logComponent instead)
function MyComp() {
  logFunction('MyComp'); // Wrong!
  logComponent('MyComp'); // Correct!
}
```

---

## 11. Files

| File | Purpose |
|------|---------|
| `src/hooks/useExecutionLogger.ts` | Logger implementation |
| `src/stores/errorStore.ts` | Integration with error capture |
| `src/components/errors/GlobalErrorModal.tsx` | UI display |

---

## 12. Future Enhancements

1. **Async Tracking** - Track Promise chains and async/await
2. **Redux Integration** - Log action dispatches
3. **React Query Integration** - Track query lifecycle
4. **Source Maps** - Map to original source locations
5. **Flamegraph Visualization** - Visual call stack

---

*Specification created: 2026-02-06*
