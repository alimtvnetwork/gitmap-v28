# 01 — Health Endpoint Format Mismatch

> **Created:** 2026-02-04  
**Version:** 3.2.0  
> **Time Wasted:** ~1 hour  
> **Severity:** High (blocked all development)

---

## Symptoms

1. Frontend showed "Backend disconnected" banner even when backend was running
2. Error modal displayed E9005 error code
3. User reported API endpoints "not working" despite backend logs showing successful requests
4. Diagnostics showed `VITE_API_URL: (not set)` which confused the issue further

---

## Root Cause

### Problem 1: Response Format Mismatch

**Backend returned:**
```json
{"status":"healthy","timestamp":"2026-02-04T11:00:00Z"}
```

**Frontend expected (in BackendStatus.tsx):**
```typescript
const data = JSON.parse(raw) as { success?: boolean; status?: string };
const connected = data.success === true || data.status === "ok";
```

The frontend checked for `success === true` OR `status === "ok"`, but backend returned `status: "healthy"`.

**Result:** Every health check failed, showing "Backend disconnected" even when backend was perfectly healthy.

### Problem 2: Confusing Diagnostics

The error modal displayed:
- `VITE_API_URL: (not set)` — but this was actually the **resolved** value, not the raw env var
- No distinction between raw environment variables and resolved/effective URLs
- Missing UI origin (host + port) to help diagnose port mismatches

### Problem 3: Missing API Index Route

Accessing `http://localhost:8080/api/v1` returned 404, making users think the API wasn't running at all.

---

## Solution

### 1. Standardized Health Response (Backend)

**File:** `backend/internal/api/handlers/handlers.go`

```go
// Before (non-standard format)
func Health(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(`{"status":"healthy",...}`))
}

// After (standard envelope with typed response)
type HealthResponse struct {
    Status    string
    Timestamp string
}

func Health(w http.ResponseWriter, r *http.Request) {
    respondSuccess(w, HealthResponse{
        Status:    "ok",
        Timestamp: time.Now().Format(time.RFC3339),
    })
}
// Returns: {"success":true,"data":{"status":"ok","timestamp":"..."}}
```

### 2. Added API Index Route (Backend)

**File:** `backend/internal/api/router.go`

```go
api.HandleFunc("", handlers.APIIndex).Methods(httpmethod.Get.String())
api.HandleFunc("/", handlers.APIIndex).Methods(httpmethod.Get.String())
```

Now `GET /api/v1` returns:
```json
{"success":true,"data":{"name":"WP Plugin Publish API","version":"v1","health":"/api/v1/health","ws":"/ws"}}
```

### 3. Fixed Detection Logic (Frontend)

**File:** `src/components/shared/BackendStatus.tsx`

```typescript
// Before: Checked for specific field values
const connected = data.success === true || data.status === "ok";

// After: Any 2xx JSON response = connected
if (response.ok) {
    setIsConnected(true);
}
```

### 4. Enhanced Diagnostics (Frontend)

**Files:** `src/lib/diagnostics.ts`, `src/lib/api.ts`, `GlobalErrorModal.tsx`

Now clearly shows:
- **Raw env vars:** `VITE_API_URL (raw): http://localhost:8080`
- **Resolved values:** `Resolved API Origin: http://localhost:8080`
- **UI Origin:** `http://localhost:8080`
- **API Base (absolute):** `http://localhost:8080/api/v1`

---

## Prevention

### For AI Agents

1. **Always verify both ends** before claiming an endpoint works:
   ```bash
   # Verify backend
   curl http://localhost:8080/api/v1/health
   
   # Verify frontend expectation
   grep -A5 "const connected" src/components/shared/BackendStatus.tsx
   ```

2. **Check response format** in the handler code, not just the spec

3. **Test the actual URL** the frontend will call, not just the path

### For Developers

1. All API endpoints MUST return the standard envelope: `{success:true, data:{...}}`
2. Health checks should use HTTP status codes as the primary indicator
3. Diagnostics should always show both raw inputs and resolved outputs

---

## Related Files

- `backend/internal/api/handlers/handlers.go` — Health handler
- `backend/internal/api/router.go` — Route registration
- `src/components/shared/BackendStatus.tsx` — Detection logic
- `src/lib/diagnostics.ts` — Diagnostics utility
- `src/lib/api.ts` — API client with error context
- `src/components/errors/GlobalErrorModal.tsx` — Error display

---

## Lessons Learned

1. **Response format consistency is critical** — One non-conforming endpoint can break the entire UX
2. **Diagnostics must distinguish raw vs resolved** — Users need to know what they configured vs what's actually being used
3. **404 on base URLs confuses everyone** — Always register an index/info route
4. **AI agents make confident mistakes** — Always verify assumptions against actual code
