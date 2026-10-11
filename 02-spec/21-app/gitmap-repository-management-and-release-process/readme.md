# GitMap Repository Management and Release Process Specification Suite

- **Directory:** `02-spec/21-app/gitmap-repository-management-and-release-process/`
- **Slug:** `gitmap-repository-management-and-release-process`
- **Specification Version:** 1.0.0
- **Status:** APPROVED

## Overview

This specification suite defines the architecture, component interaction, and implementation requirements for automated companion repository discovery (`repo-cache`), Split-DB integration, GitMap JSON manifest injection, and two-tier hierarchical subtree remediation rendering using native GitMap commands.

## Documents

1. [01-architecture-spec.md](01-architecture-spec.md) — System architecture, companion `repo-cache` resolution hierarchy, Split-DB integration, manifest injection pipeline, and two-tier subtree solutions rendering engine with native GitMap commands.
2. [02-repo-cache-json-injection.md](../../../.ai-memory/plans/subtasks/gitmap-repository-management-and-release-process/02-repo-cache-json-injection.md) — Implementation subtask plan for repository discovery cache probing, model tagging, and portable `gitmap.json` manifest synchronization.
3. [03-subtree-solutions-formatting.md](../../../.ai-memory/plans/subtasks/gitmap-repository-management-and-release-process/03-subtree-solutions-formatting.md) — Implementation subtask plan for two-tier hierarchical subtree formatting and eradication of raw git commands in favor of native GitMap CLI primitives.

## Visual References

- Screenshot 01: `assets/screenshots/gitmap-repository-management-and-release-process-01.png` — Shows dirty pull status output where remediation hints are rendered flat/inline with raw `git -C` commands, motivating the two-tier subtree restructuring.
- Screenshot 02: `assets/screenshots/gitmap-repository-management-and-release-process-02.png` — Illustrates Python subprocess overhead for repository queries, driving native Go path acceleration.
