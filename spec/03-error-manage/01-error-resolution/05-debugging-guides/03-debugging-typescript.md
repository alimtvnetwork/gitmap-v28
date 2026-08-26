# TypeScript/React Debugging Guide

> **Version:** 1.0.0  
> **Created:** 2026-02-04  
> **Applies To:** All React Frontend Applications, Shared CLI Frontend

---

## Overview

This guide covers debugging patterns for TypeScript/React applications, particularly frontends connecting to Go backends. It includes API integration verification, state management debugging, and common troubleshooting scenarios.

---

## API Integration Verification (CRITICAL)

### Endpoint Existence Check

Before implementing any API call, ALWAYS verify:

```typescript
// ❌ WRONG: Assuming endpoint exists
const response = await fetch('/api/v1/health');

// ✅ CORRECT: Verify endpoint exists in backend first
// 1. Check backend router: Is /api/v1/health registered?
// 2. Check handler: Does the handler return expected format?
// 3. Then implement frontend
```

### Response Format Verification

```typescript
// Standard response envelope (ALL backends must use this)
interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: {
    code: number;
    message: string;
    details?: string;
  };
}

// Correct detection logic
async function fetchWithValidation<T>(url: string): Promise<T> {
  const response = await fetch(url);
  
  // Primary indicator: HTTP status code (not response body!)
  if (!response.ok) {
    const body = await response.json() as ApiResponse<never>;

    throw new ApiError(
      body.error?.code ?? response.status,
      body.error?.message ?? `HTTP ${response.status}`
    );
  }
  
  const body = await response.json() as ApiResponse<T>;
  
  // Secondary check: success field
  if (!body.success) {
    throw new ApiError(
      body.error?.code ?? 0,
      body.error?.message ?? 'Unknown error'
    );
  }
  
  return body.data!;
}
```

### Connection Status Detection

```typescript
// Health check implementation
interface HealthStatus {
  connected: boolean;
  version?: string;
  latency?: number;
  error?: string;
}

async function checkHealth(baseUrl: string): Promise<HealthStatus> {
  const start = Date.now();
  
  try {
    const response = await fetch(`${baseUrl}/api/v1/health`, {
      method: HttpMethod.Get,
      headers: { 'Content-Type': 'application/json' },
      signal: AbortSignal.timeout(5000), // 5s timeout
    });
    
    // Use HTTP status as PRIMARY indicator
    if (!response.ok) {
      return {
        connected: false,
        error: `HTTP ${response.status}`,
      };
    }
    
    const data = await response.json();
    
    return {
      connected: true,
      version: data.data?.version,
      latency: Date.now() - start,
    };
  } catch (error) {
    return {
      connected: false,
      error: error instanceof Error ? error.message : 'Unknown error',
    };
  }
}
```

---

## Environment Variable Diagnostics

### Show Raw vs Resolved Values

```typescript
// DiagnosticsPanel component
interface DiagnosticsData {
  raw: Record<string, string | undefined>;
  resolved: Record<string, string>;
  origin: string;
  timestamp: string;
}

function getDiagnostics(): DiagnosticsData {
  return {
    raw: {
      VITE_API_BASE_URL: import.meta.env.VITE_API_BASE_URL,
      VITE_API_PORT: import.meta.env.VITE_API_PORT,
      NODE_ENV: import.meta.env.MODE,
    },
    resolved: {
      apiBaseUrl: getApiBaseUrl(), // After fallback logic
      wsUrl: getWebSocketUrl(),
      origin: window.location.origin,
    },
    origin: window.location.origin,
    timestamp: new Date().toISOString(),
  };
}

// Display in error modals
function ErrorModal({ error, diagnostics }: ErrorModalProps) {
  return (
    <Dialog>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Connection Error</DialogTitle>
        </DialogHeader>
        <div className="space-y-4">
          <Alert variant="destructive">
            <AlertDescription>{error.message}</AlertDescription>
          </Alert>
          
          <Collapsible>
            <CollapsibleTrigger>Show Diagnostics</CollapsibleTrigger>
            <CollapsibleContent>
              <pre className="text-xs bg-muted p-2 rounded">
                {JSON.stringify(diagnostics, null, 2)}
              </pre>
            </CollapsibleContent>
          </Collapsible>
        </div>
      </DialogContent>
    </Dialog>
  );
}
```

---

## Common Issues and Solutions

### Issue: "Backend disconnected" but server is running

**Symptoms:**
- Go server logs show it's running
- Frontend shows "disconnected" or "connection failed"

**Check:**

1. **Is the API base URL correct?**
   ```typescript
   // Debug: Log the actual URL being called
   console.log('API URL:', getApiBaseUrl());
   ```

2. **Is there a CORS issue?**
   ```typescript
   // Browser console will show CORS errors
   // Check backend for proper CORS headers
   ```

3. **Is the response format correct?**
   ```typescript
   // Check raw response
   fetch('/api/v1/health')
     .then(r => r.text())
     .then(text => console.log('Raw response:', text));
   ```

4. **Is the frontend checking the right field?**
   ```typescript
   // ❌ Wrong: Checking body field for connection status
   const connected = data.status === 'ok';
   
   // ✅ Correct: Use HTTP status as primary indicator
   const connected = response.ok; // status 200-299
   ```

### Issue: API calls work in development but fail in production

**Symptoms:**
- Works on localhost:3000
- Fails when deployed

**Check:**

1. **Environment variables set correctly?**
   ```typescript
   // Add to production build process
   console.log('Build-time env:', {
     VITE_API_BASE_URL: import.meta.env.VITE_API_BASE_URL,
     MODE: import.meta.env.MODE,
   });
   ```

2. **Relative vs absolute URLs?**
   ```typescript
   // ❌ May break in production
   fetch('http://localhost:8080/api/v1/health');
   
   // ✅ Better: Use relative URL or env var
   fetch(`${import.meta.env.VITE_API_BASE_URL}/api/v1/health`);
   ```

3. **HTTPS in production, HTTP in dev?**
   ```typescript
   // Ensure WebSocket URL uses correct protocol
   const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
   ```

### Issue: State not updating after API call

**Symptoms:**
- API returns data successfully
- UI doesn't reflect new data

**Check:**

1. **Is state being set correctly?**
   ```typescript
   // Use React Query for automatic cache invalidation
   const queryClient = useQueryClient();
   
   const mutation = useMutation({
     mutationFn: updateData,
     onSuccess: () => {
       // Invalidate and refetch
       queryClient.invalidateQueries({ queryKey: ['data'] });
     },
   });
   ```

2. **Is the component re-rendering?**
   ```typescript
   // Debug: Add useEffect to track renders
   useEffect(() => {
     console.log('Component rendered with:', data);
   }, [data]);
   ```

### Issue: TypeScript errors with API responses

**Symptoms:**
- Type errors when accessing response data
- `Property 'x' does not exist on type`

**Solution:**

```typescript
// Define response types explicitly
interface HealthResponse {
  status: string;
  version: string;
  timestamp: string;
}

// Use type assertion with validation
function parseHealthResponse(data: unknown): HealthResponse {
  // Validate at runtime
  if (
    typeof data !== 'object' ||
    data === null ||
    !('status' in data) ||
    !('version' in data)
  ) {
    throw new Error('Invalid health response format');
  }
  
  return data as HealthResponse;
}

// Or use Zod for schema validation
import { z } from 'zod';

const HealthResponseSchema = z.object({
  status: z.string(),
  version: z.string(),
  timestamp: z.string(),
});

type HealthResponse = z.infer<typeof HealthResponseSchema>;

async function getHealth(): Promise<HealthResponse> {
  const response = await fetch('/api/v1/health');
  const data = await response.json();

  return HealthResponseSchema.parse(data.data);
}
```

---

## React Query Debugging

### Enable DevTools

```typescript
// main.tsx
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Routes />
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  );
}
```

### Debug Query State

```typescript
// Check query state
const { data, error, status, fetchStatus, isLoading, isFetching } = useQuery({
  queryKey: ['health'],
  queryFn: fetchHealth,
});

console.log('Query state:', {
  status,        // 'pending' | 'error' | 'success'
  fetchStatus,   // 'fetching' | 'paused' | 'idle'
  isLoading,     // First load (no cached data)
  isFetching,    // Any fetch (including background refetch)
});
```

---

## WebSocket Debugging

### Connection State Tracking

```typescript
function useWebSocket(url: string) {
  const [state, setState] = useState<'connecting' | 'connected' | 'disconnected'>('connecting');
  const wsRef = useRef<WebSocket | null>(null);
  
  useEffect(() => {
    console.log(`WebSocket connecting to: ${url}`);
    
    const ws = new WebSocket(url);
    wsRef.current = ws;
    
    ws.onopen = () => {
      console.log('WebSocket connected');
      setState('connected');
    };
    
    ws.onclose = (event) => {
      console.log('WebSocket closed:', {
        code: event.code,
        reason: event.reason,
        wasClean: event.wasClean,
      });
      setState('disconnected');
    };
    
    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
    
    ws.onmessage = (event) => {
      console.log('WebSocket message:', event.data);
    };
    
    return () => {
      ws.close();
    };
  }, [url]);
  
  return { state, ws: wsRef.current };
}
```

---

## Console Logging Best Practices

### Structured Logging

```typescript
// Logger utility
const logger = {
  debug: (message: string, data?: unknown) => {
    if (import.meta.env.DEV) {
      console.log(`[DEBUG] ${message}`, data ?? '');
    }

  },
  
  info: (message: string, data?: unknown) => {
    console.log(`[INFO] ${message}`, data ?? '');
  },
  
  warn: (message: string, data?: unknown) => {
    console.warn(`[WARN] ${message}`, data ?? '');
  },
  
  error: (message: string, error?: unknown) => {
    console.error(`[ERROR] ${message}`, error ?? '');
  },
  
  api: (method: string, url: string, status: number, duration: number) => {
    const emoji = status >= 400 ? '❌' : '✅';
    console.log(`${emoji} [API] ${method} ${url} → ${status} (${duration}ms)`);
  },
};

// Usage
logger.api('GET', '/api/v1/health', 200, 45);
logger.error('Failed to fetch settings', error);
```

### Request/Response Logging

```typescript
// API client with logging
async function apiRequest<T>(
  method: string,
  url: string,
  body?: unknown
): Promise<T> {
  const start = Date.now();
  
  logger.debug(`Request: ${method} ${url}`, body);
  
  try {
    const response = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: body ? JSON.stringify(body) : undefined,
    });
    
    const data = await response.json();
    const duration = Date.now() - start;
    
    logger.api(method, url, response.status, duration);
    logger.debug(`Response:`, data);
    
    if (!response.ok) {
      throw new ApiError(data.error?.code ?? response.status, data.error?.message);
    }
    
    return data.data;
  } catch (error) {
    const duration = Date.now() - start;
    logger.error(`Request failed: ${method} ${url} (${duration}ms)`, error);

    throw error;
  }
}
```

---

## Browser DevTools Tips

### Network Tab Filtering

```
// Filter by URL pattern
/api/v1/

// Filter by status
status-code:500

// Filter by method
method:POST
```

### Console Commands

```javascript
// Clear console
clear()

// Monitor function calls
monitor(functionName)

// Time operations
console.time('operation');
// ... do work
console.timeEnd('operation');

// Group related logs
console.group('API Call');
console.log('URL:', url);
console.log('Response:', data);
console.groupEnd();

// Table display
console.table([{ id: 1, name: 'Test' }]);
```

---

## Error Boundary Pattern

```typescript
// ErrorBoundary.tsx
import { Component, ReactNode } from 'react';

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
}

interface State {
  hasError: boolean;
  error?: Error;
}

class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('ErrorBoundary caught:', {
      error: error.message,
      stack: error.stack,
      componentStack: errorInfo.componentStack,
    });
    
    // Report to error tracking service
    // reportError(error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return this.props.fallback ?? (
        <div className="p-4 border border-destructive rounded">
          <h2 className="text-lg font-semibold text-destructive">Something went wrong</h2>
          <pre className="mt-2 text-sm text-muted-foreground">
            {this.state.error?.message}
          </pre>
        </div>
      );
    }

    return this.props.children;
  }
}

export { ErrorBoundary };
```

---

## Performance Debugging

### React Profiler

```typescript
import { Profiler, ProfilerOnRenderCallback } from 'react';

const onRender: ProfilerOnRenderCallback = (
  id,
  phase,
  actualDuration,
  baseDuration,
  startTime,
  commitTime
) => {
  console.log(`[Profiler] ${id}:`, {
    phase,          // 'mount' | 'update'
    actualDuration, // Time spent rendering
    baseDuration,   // Estimated time without memoization
    startTime,
    commitTime,
  });
};

function App() {
  return (
    <Profiler id="App" onRender={onRender}>
      <MainContent />
    </Profiler>
  );
}
```

### why-did-you-render

```typescript
// wdyr.ts (import before React)
import React from 'react';

if (import.meta.env.DEV) {
  const whyDidYouRender = await import('@welldone-software/why-did-you-render');
  whyDidYouRender.default(React, {
    trackAllPureComponents: true,
    logOnDifferentValues: true,
  });
}
```

---

## Cross-Reference

- [Error Resolution Overview](../../00-overview.md)
- [Frontend-Backend Sync Verification](../04-verification-patterns/01-frontend-backend-sync.md)
- [Go Debugging Guide](./02-debugging-go.md)
- React Coding Guidelines *(external spec)*
- Shared CLI Frontend *(external spec)*
