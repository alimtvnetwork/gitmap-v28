---
name: zero-storage-actions-purge
description: Autonomously discover, purge, and govern GitHub Actions storage, deleting stored artifacts and caches to maintain 0.0 GB usage across repositories.
---

# Zero-Storage Actions Purge & Storage Governance

Autonomously inspects repositories, identifies bloated GitHub Actions storage (artifacts and caches), eliminates excessive caching configurations in workflows, and executes multi-threaded API deletion of remote artifacts and caches to maintain zero storage footprint.

## Core Capabilities

- **Storage Auditing:** Inspects active caches (`/actions/caches`) and stored artifacts (`/actions/artifacts`) via GitHub REST API across repositories.
- **Workflow Optimization:**
  - Adds `retention-days: 1` to all temporary build handoff `upload-artifact` steps.
  - Removes bloated compiler object caches (`cache-all-crates: "true"`, multi-gigabyte target directory caching).
  - Adds post-release artifact deletion hooks to release workflows.
- **Scheduled Automated Purge:** Deploys `.github/workflows/purge-actions-artifacts.yml` running on nightly cron (`0 2 * * *`) and `workflow_dispatch`.
- **High-Speed Autonomous Deletion:** Runs `03-ai-scripts/34-purge-github-actions-artifacts.py` with multi-threaded concurrent deletion (`ThreadPoolExecutor`) to instantly reclaim gigabytes of storage across repositories.

## Workflow Execution Steps

1. **Defensive Git Pull:** Run `git pull` on the target repository before modifying any workflows or scripts.
2. **Audit & Fix Workflows:**
   - Scan `.github/workflows/` for missing `retention-days: 1` on `actions/upload-artifact`.
   - Remove redundant or massive object caches.
   - Install or update `.github/workflows/purge-actions-artifacts.yml`.
3. **Deploy Purge Automation:** Ensure `03-ai-scripts/34-purge-github-actions-artifacts.py` exists with dual artifact and cache purge support.
4. **Execute Remote Purge:** Run:
   ```bash
   python 03-ai-scripts/34-purge-github-actions-artifacts.py --repo <owner/repo>
   ```
5. **Verify & Test:** Validate YAML integrity of workflows and verify storage has dropped to 0.0 MB.
6. **Atomic Commit & Push:** Commit all workflow changes and scripts atomically and immediately push to `origin`.
