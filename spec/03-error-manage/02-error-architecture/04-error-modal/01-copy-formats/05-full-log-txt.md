# log.txt (Backend Full Log)

> **Parent:** [Copy Formats Index](./00-overview.md)  
> **Version:** 3.2.0  
> **Updated:** 2026-03-31  
> **Purpose:** Raw text file from `GET /api/v1/logs/full`. Contains ALL log entries (info, warn, error).

---

## Sample Entry

```
[2026-02-12 00:53:30] [INFO] GET /api/v1/sites
  Duration: 12ms
  Status: 200

[2026-02-12 00:53:30] [INFO] GET /api/v1/sites/1/mappings
  Duration: 8ms
  Status: 200

[2026-02-12 00:53:34] [ERROR] GET /api/v1/sites/1/snapshots/settings
  Duration: 4190ms
  Status: 500
  Error: [E3001] failed to fetch snapshot settings...
```

---

*log.txt format — updated: 2026-03-31*
