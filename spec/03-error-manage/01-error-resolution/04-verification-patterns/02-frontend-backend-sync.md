# Frontend-Backend Sync Verification Pattern

> **Version:** 1.0.0
> **Created:** 2026-02-04
> **Status:** MANDATORY

---

## Overview

This document defines the mandatory verification pattern for ensuring frontend and backend are properly synchronized. This pattern MUST be followed before claiming any API integration works.

---

## The Three-Step Verification

### Step 1: Backend Verification

Verify the backend endpoint exists and returns the expected format.

```bash

# Test health endpoint

curl -s http://localhost:8080/api/v1/health | jq .

# Expected output:

{
  "success": true,
  "data": {
    "status": "ok",
    "timestamp": "2026-02-04T12:00:00Z"
  }
}
```

**Check for:**

- ✅ HTTP 200 status code
- ✅ Standard envelope `{success: true, data: {...}}`
- ✅ Correct Content-Type header (`application/json`)
- ✅ CORS headers present (if cross-origin)

### Step 2: Frontend Verification

Verify the frontend code correctly handles the backend response.

```bash

# Find detection logic

grep -rn "connected\|isConnected\|setConnected" src/

# Check what the code expects

cat src/components/shared/BackendStatus.tsx
```

**Check for:**

- ✅ Uses HTTP status code (2xx) as primary indicator
- ✅ Handles response envelope correctly
- ✅ Has proper error handling
- ✅ Timeout configured appropriately

### Step 3: Integration Verification

Test the actual integration in the browser.

```javascript
// Open browser console and test
fetch('/api/v1/health')
  .then(r => r.json())
  .then(console.log)
  .catch(console.error);
```

**Check for:**

- ✅ Network request visible in DevTools
- ✅ Response matches expected format
- ✅ No CORS errors
- ✅ UI updates correctly

---

## Endpoint Existence Checklist

Before implementing frontend code for an endpoint:

| Check | How | Status |
|-------|-----|--------|
| Route registered | Check router file | ☐ |
| Handler implemented | Check handlers file | ☐ |
| Response format standard | Return `{success, data}` | ☐ |
| Error response standard | Return `{success: false, error}` | ☐ |
| OpenAPI documented | Check openapi.yaml | ☐ |

---

## Response Format Verification

### Standard Success Response

```json
{
  "success": true,
  "data": {
    "key": "value"
  }
}
```

### Standard Error Response

```json
{
  "success": false,
  "error": {
    "code": "ERR_NOT_FOUND",
    "message": "Resource not found"
  }
}
```

### Verify with curl

```bash

# Success case

curl -s http://localhost:8080/api/v1/resource/123 | jq '.success'

# Should return: true

# Error case

curl -s http://localhost:8080/api/v1/resource/invalid | jq '.success'

# Should return: false

```

---

## Environment Variable Verification

### Check Raw vs Resolved

```typescript
// In diagnostics code, always show both:
const rawApiUrl = import.meta.env.VITE_API_URL ?? '(not set)';
const resolvedApiUrl = getResolvedApiUrl(); // After any defaults/fallbacks

console.log('Raw VITE_API_URL:', rawApiUrl);
console.log('Resolved API URL:', resolvedApiUrl);
```

### Common Misconfigurations

| Symptom | Likely Cause |
|---------|--------------|
| `VITE_API_URL: (not set)` | Missing .env file or wrong variable name |
| `VITE_API_URL: /api` | Relative URL, may need full origin |
| `VITE_API_URL: http://localhost:8080` but errors | Port mismatch or backend not running |

---

## WebSocket Verification

For CLI frontends using WebSocket:

```javascript
// Test WebSocket connection
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = () => console.log('WS Connected');
ws.onmessage = (e) => console.log('WS Message:', e.data);
ws.onerror = (e) => console.error('WS Error:', e);
ws.onclose = (e) => console.log('WS Closed:', e.code, e.reason);
```

**Check for:**

- ✅ Connection established (onopen fired)
- ✅ Messages received in expected format
- ✅ Ping/pong working (if implemented)
- ✅ Graceful close handling

---

## Automated Verification Script

Create a verification script for your project:

```bash
#!/bin/bash

# verify-api.sh

API_BASE="${API_BASE:-http://localhost:8080}"

echo "=== API Verification ==="
echo "Base URL: $API_BASE"
echo ""

# Check health

echo "1. Health Check:"
HEALTH=$(curl -s -w "\n%{http_code}" "$API_BASE/api/v1/health")
HTTP_CODE=$(echo "$HEALTH" | tail -n1)
BODY=$(echo "$HEALTH" | head -n1)

if [ "$HTTP_CODE" = "200" ]; then
    echo "   ✅ HTTP 200"
    if echo "$BODY" | jq -e '.success == true' > /dev/null 2>&1; then
        echo "   ✅ success: true"
    else
        echo "   ❌ Missing success: true"
    fi
else
    echo "   ❌ HTTP $HTTP_CODE"
fi

# Check API index

echo ""
echo "2. API Index:"
INDEX=$(curl -s -w "\n%{http_code}" "$API_BASE/api/v1")
HTTP_CODE=$(echo "$INDEX" | tail -n1)

if [ "$HTTP_CODE" = "200" ]; then
    echo "   ✅ HTTP 200"
else
    echo "   ❌ HTTP $HTTP_CODE (missing index route?)"
fi

echo ""
echo "=== Verification Complete ==="
```

---

## Integration with CI/CD

Add verification to your CI pipeline:

```yaml

# .github/workflows/verify-api.yml

name: API Verification

on: [push, pull_request]

jobs:
  verify:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Start backend
        run: |
          cd backend
          go run . daemon &
          sleep 5

      - name: Verify endpoints
        run: |
          ./linter-scripts/verify-api.sh
```

---

## Cross-Reference

- [Error Resolution Overview](../../01-index.md)
- [Health Endpoint Mismatch Retrospective](../03-retrospectives/02-health-endpoint-mismatch.md)
- Shared CLI Frontend *(external spec)*
