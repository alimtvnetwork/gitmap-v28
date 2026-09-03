# Specification: Repository Creation, Multi-Account Profiles, and Git-Backed Cloud Backup Suite

---
**Version:** 1.0.0  
**Updated:** 2026-09-03  
**Status:** Active  
**AI Confidence Score:** Production-Ready  
**Ambiguity Score:** None  
**Module Health Score:** 100/100  
---

## Purpose & Scope

This specification defines the architecture, data contracts, and CLI interfaces for:
1. **Multi-Account Git Profiles (`gitmap profiles`)**: Managing multiple user accounts and organizations across GitHub and GitLab with auto-discovery, usage tracking, and numeric sequence picking (`1, 2, 3...`).
2. **Repository Creation Engine (`gitmap create`)**: End-to-end local repository provisioning and remote cloud repository publishing bound to the active default profile.
3. **Git-Backed Cloud Backup & Disaster Recovery (`gitmap backup`)**: Transactional backup and restoration of GitMap SQLite databases, schemas, macros, and catalog metadata into a dedicated private cloud repository.

---

## File Inventory

| File | Description |
|------|-------------|
| [`00-overview.md`](./00-overview.md) | Module entry point, metadata, and architecture overview |
| [`01-git-profiles-management.md`](./01-git-profiles-management.md) | Multi-account profiles schema, discovery engine, and sequence picker |
| [`02-repo-creation.md`](./02-repo-creation.md) | Local repo initialization, README/.gitignore generation, and remote push |
| [`03-git-backed-cloud-backup.md`](./03-git-backed-cloud-backup.md) | Cloud Git backup synchronization, manifest hashing, and restoration |
| [`97-acceptance-criteria.md`](./97-acceptance-criteria.md) | Test suites, Given/When/Then scenarios, and validation matrix |
| [`99-consistency-report.md`](./99-consistency-report.md) | Structural, cross-reference, and consistency audit report |
