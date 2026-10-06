# 15 — Fast File Reader and Testing Hygiene

- **Subsystem:** File I/O Engine & Test Isolation
- **Status:** Authoritative Reference

## 1. 8KB Binary Probe Guard
- Scans the initial 8KB chunk of files for null bytes (`\x00`) to safely identify and skip binary assets.
- Automatically strips UTF-8 Byte Order Marks (`\xef\xbb\xbf`) and normalizes CRLF to LF (`\n`).

## 2. Hermetic Test Isolation
- Unit tests run against temporary mock directories (`t.TempDir()`) without touching user configuration or live Git trees.
- Destructive system commands (power-off, shutdown, disk formatting) are strictly intercepted by mock executors.
