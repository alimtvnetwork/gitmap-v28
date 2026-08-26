# error.log.txt (Backend Error Log)

> **Parent:** [Copy Formats Index](./00-overview.md)  
> **Version:** 3.2.0  
> **Updated:** 2026-03-31  
> **Purpose:** Raw text file from `GET /api/v1/logs/error`. Contains only error-level entries.

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
```

---

## Format Rules

| Line | Content |
|------|---------|
| `[timestamp]` | ISO-like timestamp `YYYY-MM-DD HH:MM:SS` |
| `HTTP {code} {method} FAILED` | Request summary |
| `Requested To:` | Full URL of Go backend endpoint |
| `Duration:` | Request duration |
| `Error Code:` | HTTP status code |
| `Error Message:` | Primary error with error code prefix |
| `Backend Error:` | Wrapped error with additional context code |
| `Go Backend Stack:` | File:line + function for each Go frame |
| `Go Methods Stack:` | Numbered, formatted method references |

---

*error.log.txt format — updated: 2026-03-31*
