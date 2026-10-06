# 11 — Macro Automation and Execution Idempotency

- **Subsystem:** Interactive Macro Automation & VMware Integration
- **Status:** Authoritative Reference

## 1. Idempotent Macro Playback
- Macro steps define explicit preconditions (`exists`, `absent`, `matches`, `process_running`).
- Safely skips redundant operations on repeated runs without triggering errors.

## 2. Child Process Lifecycle Management
- Tracks and terminates lingering background processes spawned during macro sequences.
- Prevents stale development servers or build watchers from blocking subsequent tasks.

## 3. VMware Guest Automation
- Manages shared folder mounts (`/mnt/hgfs`), guest tools status, and snapshot synchronization.
