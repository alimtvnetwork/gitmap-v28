# Backend error.log.txt with Delegated Server Info

> **Parent:** [Copy Formats Index](./00-overview.md)  
> **Version:** 3.2.0  
> **Updated:** 2026-03-31  
> **Purpose:** Enhanced error.log.txt format that includes the Delegated Server Info section when the Go backend proxied a request to a downstream server.

---

## Complete Sample

```
[2026-02-12 00:53:34] HTTP 500 GET FAILED
  Requested To: GET http://localhost:8080/api/v1/sites/1/snapshots/settings
  Duration: 4.1904552s
  Error Code: 500
  Error Message: [E3001] failed to fetch snapshot settings: get snapshot settings (GET https://demoat.attoproperty.com.au/riseup-asia-uploader/v1/snapshots/settings): status 403
  Backend Error: [E3025] [E3001] failed to fetch snapshot settings: get snapshot settings (GET https://demoat.attoproperty.com.au/riseup-asia-uploader/v1/snapshots/settings): status 403
  Go Backend Stack:
    D:/wp-work/riseup-asia/wp-onboarding-v5/backend/internal/api/handlers/handler_factory.go:107 handlers.init.handleSiteActionById.func63
    D:/wp-work/riseup-asia/wp-onboarding-v5/backend/internal/api/middleware/session_logging.go:107 api.NewServer.SessionLogging.func3.1
    D:/wp-work/riseup-asia/wp-onboarding-v5/backend/internal/api/middleware/middleware.go:245 api.NewServer.Recovery.func2.1
    D:/wp-work/riseup-asia/wp-onboarding-v5/backend/internal/api/middleware/middleware.go:66 api.NewServer.Logging.func1.1
    D:/wp-work/riseup-asia/wp-onboarding-v5/backend/internal/api/middleware/middleware.go:45 middleware.CORS.func1
  Go Methods Stack:
    #0 handlers.init.handleSiteActionById.func63 at D:/wp-work/riseup-asia/wp-onboarding-v5/backend/internal/api/handlers/handler_factory.go:107
    #1 api.NewServer.SessionLogging.func3.1 at D:/wp-work/riseup-asia/wp-onboarding-v5/backend/internal/api/middleware/session_logging.go:107
    #2 api.NewServer.Recovery.func2.1 at D:/wp-work/riseup-asia/wp-onboarding-v5/backend/internal/api/middleware/middleware.go:245
    #3 api.NewServer.Logging.func1.1 at D:/wp-work/riseup-asia/wp-onboarding-v5/backend/internal/api/middleware/middleware.go:66
    #4 middleware.CORS.func1 at D:/wp-work/riseup-asia/wp-onboarding-v5/backend/internal/api/middleware/middleware.go:45
  Delegated Server Info:
    Endpoint: "https://demoat.attoproperty.com.au/wp-json/riseup-asia-uploader/v1/snapshots/settings"
    Method: "GET"
    Status: 403
    Stacktrace:
        #0 riseup-asia-uploader.php(1098): FileLogger->error()
        #1 class-wp-hook.php(341): Plugin->enrichErrorResponse()
        #2 plugin.php(205): WP_Hook->apply_filters()
        #3 class-wp-rest-server.php(462): apply_filters()
        #4 rest-api.php(467): WP_REST_Server->serve_request()
    RequestBody:
        (none — GET request)
    Additional Message:
        Endpoint 'snapshots' is not enabled in plugin settings. Navigate to WordPress Admin → Settings → Riseup Asia → Endpoints to enable it.
  Response Body:
    {
  "Status": {
    "IsSuccess": false,
    "IsFailed": true,
    "Code": 500,
    "Message": "[E3001] failed to fetch snapshot settings: get snapshot settings (GET https://demoat.attoproperty.com.au/wp-json/riseup-asia-uploader/v1/snapshots/settings): status 403",
    "Timestamp": "2026-02-11T16:53:35Z"
  },
  "Attributes": {
    "HasAnyErrors": true,
    "IsSingle": false,
    "IsMultiple": false
  },
  "Results": [],
  "Errors": {
    "BackendMessage": "[E3025] [E3001] failed to fetch snapshot settings: get snapshot settings (GET https://demoat.attoproperty.com.au/riseup-asia-uploader/v1/snapshots/settings): status 403",
    "Backend": [
      "handler_factory.go:107 handlers.init.handleSiteActionById.func63",
      "session_logging.go:107 api.NewServer.SessionLogging.func3.1",
      "middleware.go:245 api.NewServer.Recovery.func2.1",
      "middleware.go:66 api.NewServer.Logging.func1.1",
      "middleware.go:45 middleware.CORS.func1"
    ],
    "DelegatedRequestServer": {
      "DelegatedEndpoint": "https://demoat.attoproperty.com.au/wp-json/riseup-asia-uploader/v1/snapshots/settings",
      "Method": "GET",
      "StatusCode": 403,
      "RequestBody": null,
      "Response": {
        "code": "rest_forbidden",
        "message": "This endpoint is disabled",
        "data": {
          "status": 403,
          "plugin_version": "1.54.0",
          "timestamp": "2026-02-11T11:50:58Z",
          "log_hint": "Check WP Settings → Riseup Asia → Endpoints to enable snapshots"
        }
      },
      "StackTrace": [
        "#0 riseup-asia-uploader.php(1098): FileLogger->error()",
        "#1 class-wp-hook.php(341): Plugin->enrichErrorResponse()",
        "#2 plugin.php(205): WP_Hook->apply_filters()",
        "#3 class-wp-rest-server.php(462): apply_filters()",
        "#4 rest-api.php(467): WP_REST_Server->serve_request()"
      ],
      "AdditionalMessages": "Endpoint 'snapshots' is not enabled in plugin settings."
    }
  },
  "MethodsStack": {
    "Backend": [
      {
        "Method": "handlers.init.handleSiteActionById.func63",
        "File": "handler_factory.go",
        "LineNumber": 107
      },
      {
        "Method": "api.NewServer.SessionLogging.func3.1",
        "File": "session_logging.go",
        "LineNumber": 107
      }
    ],
    "Frontend": []
  }
}
```

---

## Delegated Server Info Format Rules

| Line | Content | When Present |
|------|---------|-------------|
| `Delegated Server Info:` | Section header | Always when delegated request failed |
| `Endpoint:` | Exact downstream URL (quoted) | Always |
| `Method:` | HTTP method (quoted) | Always |
| `Status:` | HTTP status code | Always |
| `Stacktrace:` | Indented stack trace lines from delegated server | If delegated server returned stack trace |
| `RequestBody:` | JSON body sent to delegated server | If POST/PUT/PATCH, otherwise "(none — GET request)" |
| `Additional Message:` | Human-readable error context | If available |

---

*error.log.txt with delegated info — updated: 2026-03-31*
