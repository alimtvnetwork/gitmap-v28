# Go Backend Fix: Missing Delegation Fields


**Version:** 3.2.0  
**Last Updated:** 2026-04-16  

> **Problem:** When the Go backend proxies a request to a downstream service (e.g., WordPress PHP plugin) and the downstream call fails, the response envelope is missing `Attributes.RequestDelegatedAt` and `Errors.DelegatedRequestServer`. This makes it impossible for the frontend error modal to show the full 3-hop request chain or the third-party's error details.

## What Must Change

When a handler in the Go backend makes an HTTP request to a downstream/delegated service and receives a **non-2xx response**, the envelope builder **MUST** populate:

### 1. `Attributes.RequestDelegatedAt` (string)

The **full URL** of the downstream endpoint that was called.

```go
// Example: in the handler that proxies to WordPress
envelope.Attributes.RequestDelegatedAt = delegatedUrl
// e.g. "https://demoat.attoproperty.com.au/wp-json/riseup-asia-uploader/v1/snapshots/providers"
```

**Rule:** This field MUST be set for **every** delegated request — even successful ones. It tells the frontend that a 3rd-party hop occurred.

---

### 2. `Errors.DelegatedRequestServer` (object)

Structured error details from the downstream server. **Required when `IsFailed=true` AND `RequestDelegatedAt` is non-empty.**

```go
type DelegatedRequestServer struct {
    DelegatedEndpoint  string                               // Full URL
    Method             string                               // HTTP method used
    StatusCode         int                                  // Response status code
    RequestBody        json.RawMessage `json:",omitempty"`   // Request body sent (null for GET)
    Response           json.RawMessage `json:",omitempty"`   // Full response body (parsed JSON if possible)
    StackTrace         []string        `json:",omitempty"`   // PHP/Node/etc stack trace lines
    AdditionalMessages string          `json:",omitempty"`   // Human-readable diagnostic hint
}
```

---

## Implementation Pattern

In the handler or service layer that makes the delegated HTTP call:

```go
func (s *Service) fetchFromDelegatedServer(
	context stdctx.Context,
	site *models.Site,
	path string,
) apperror.Result[*Envelope] {
    delegatedUrl := fmt.Sprintf("%s/wp-json/%s", site.Url, path)

    result := s.executeDelegatedRequest(context, delegatedUrl)

    if result.HasError() {
        return apperror.Fail[*Envelope](result.Error())
    }

    resp, bodyBytes := result.Value().Response, result.Value().Body
    defer resp.Body.Close()

    return s.buildDelegatedEnvelope(delegatedUrl, path, resp.StatusCode, bodyBytes)
}

func (s *Service) buildDelegatedEnvelope(
    delegatedUrl, path string,
    statusCode int,
    bodyBytes []byte,
) apperror.Result[*Envelope] {
    envelope := NewEnvelope()
    envelope.Attributes.RequestDelegatedAt = delegatedUrl

    if statusCode >= 400 {
        return s.buildDelegatedErrorEnvelope(envelope, delegatedUrl, path, statusCode, bodyBytes)
    }

    // Success path — RequestDelegatedAt is still set so frontend knows a hop occurred
    // ... parse Results ...
    return apperror.Succeed(envelope)
}

func (s *Service) executeDelegatedRequest(context stdctx.Context, url string) apperror.Result[DelegatedRequestResult] {
    req, err := http.NewRequestWithContext(context, http.MethodGet, url, nil)

    if err != nil {
        return apperror.Fail[DelegatedRequestResult](
            apperror.Wrap(err, apperror.ErrWpConnect, "failed to build delegated request"),
        )
    }

    resp, err := s.httpClient.Do(req)

    if err != nil {
        return apperror.Fail[DelegatedRequestResult](
            apperror.Wrap(err, apperror.ErrWpConnect, "failed to reach delegated server").
                WithEndpoint(url),
        )
    }

    bodyBytes, _ := io.ReadAll(resp.Body)

    return apperror.Succeed(DelegatedRequestResult{
        Response: resp,
        Body:     bodyBytes,
    })
}

func (s *Service) buildDelegatedErrorEnvelope(
    envelope *Envelope,
    delegatedUrl, path string,
    statusCode int,
    bodyBytes []byte,
) apperror.Result[*Envelope] {
    var parsed DelegatedResponseBody
    _ = json.Unmarshal(bodyBytes, &parsed)

    envelope.Errors = &EnvelopeErrors{
        BackendMessage:         s.formatDelegatedError(path, statusCode),
        DelegatedRequestServer: s.newDelegatedServer(delegatedUrl, statusCode, bodyBytes, &parsed),
    }

    return apperror.Fail[*Envelope](
        apperror.New(
        "E3001", fmt.Sprintf("delegated request failed with status %d", statusCode),
    )
}

func (s *Service) formatDelegatedError(path string, statusCode int) string {
    return fmt.Sprintf("[E3001] failed to fetch %s: %s (GET %s): status %d",
        path, path, path, statusCode)
}

func (s *Service) newDelegatedServer(
    delegatedUrl string,
    statusCode int,
    bodyBytes []byte,
    parsed *DelegatedResponseBody,
) *DelegatedRequestServer {
    return &DelegatedRequestServer{
        DelegatedEndpoint:  delegatedUrl,
        Method:             http.MethodGet,
        StatusCode:         statusCode,
        RequestBody:        nil,
        Response:           json.RawMessage(bodyBytes),
        StackTrace:         parsed.Data.StackTrace,
        AdditionalMessages: parsed.Data.LogHint,
    }
}
```

---

## Typed Response Body for Delegation Parsing

Instead of using `interface{}` and type assertions, define a concrete struct matching the WordPress REST API error format:

```go
// DelegatedResponseBody represents the parsed JSON structure from a WordPress-style error response.
type DelegatedResponseBody struct {
    Data DelegatedResponseData `json:"data"` // EXEMPTED: WordPress REST API response format
}

// DelegatedResponseData holds the structured error details from a delegated server.
type DelegatedResponseData struct {
    StackTrace []string `json:"stack_trace"` // EXEMPTED: WordPress REST API response format
    LogHint    string   `json:"log_hint"`    // EXEMPTED: WordPress REST API response format
}
```

This replaces the legacy `extractPHPStackTrace(body interface{})` and `extractLogHint(body interface{})` helper functions with direct struct field access (`parsed.Data.StackTrace`, `parsed.Data.LogHint`), eliminating all `interface{}` type assertions.

---

## Checklist

- [ ] `Attributes.RequestDelegatedAt` is set for **all** delegated calls (success and failure)
- [ ] `Errors.DelegatedRequestServer` is populated for **all** failed delegated calls (status ≥ 400)
- [ ] `DelegatedRequestServer.Response` includes the full parsed response body
- [ ] `DelegatedRequestServer.StackTrace` includes PHP stack trace when available
- [ ] `DelegatedRequestServer.AdditionalMessages` includes the `log_hint` from WordPress error responses
- [ ] `DelegatedRequestServer.Method` matches the actual HTTP method used
- [ ] `DelegatedRequestServer.RequestBody` is populated for POST/PUT/DELETE (null for GET)

---

## Frontend Detection

The React frontend now shows a **"Missing Delegation Data"** warning (amber banner) in the error modal Overview tab when:
1. The error message contains a third-party endpoint pattern (e.g., `(GET https://demoat.attoproperty.com.au/wp-json/riseup-asia-uploader/v1/...)`)
2. But `RequestDelegatedAt` is empty/missing
3. And `DelegatedRequestServer` is null/missing

This warning explicitly states it's a backend bug to help developers identify the issue quickly.
