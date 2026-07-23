# Future Features — Pending Discussion

## Overview

This document tracks planned features that require further discussion
before specification and implementation.

---

## 1. REST API / Cloud Sync

**Status**: Pending discussion

A remote repository or REST API endpoint where gitmap-v27 can sync
repo metadata, task definitions, and configuration. This enables:

- Multi-machine synchronization of scan results and groups.
- Team-shared repo registries.
- Remote backup of gitmap-v27 database.
- Cloud-based task management and monitoring.

### Open Questions
- Self-hosted vs managed service?
- Authentication model (API key, OAuth)?
- Conflict resolution strategy for concurrent edits?
- Data scope: full DB sync or selective (repos, groups, tasks)?

---

## 2. Install Command — Remote Manifests

**Status**: Pending discussion

Extend `gitmap-v27 install` with REST URL support for remote install plans:

```bash
gitmap-v27 install --from https://company.com/dev-setup.json
```

### Capabilities
- Fetch a JSON manifest listing tools and versions.
- Execute batch installation from the manifest.
- Version pinning and update tracking.
- Organization-specific tool configurations.

### Open Questions
- Manifest schema and versioning?
- Signature verification for security?
- Offline caching of manifests?

---

## 3. Collaborative Watch Tasks

**Status**: Future consideration

Allow watch tasks to sync across machines via the REST API,
enabling distributed file synchronization workflows.

---

## See Also

- task — file sync automation (spec 79)
- env — environment variable management (spec 80)
- install — developer tool installer (spec 81)
