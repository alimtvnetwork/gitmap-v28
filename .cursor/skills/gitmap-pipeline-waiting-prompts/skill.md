---
name: gitmap-pipeline-waiting-prompts
description: Update execute, RCA, and CI/CD fix prompts and skills to use GitMap pipeline-ai with dynamic timeout waiting.
---

# GitMap Pipeline-AI & Dynamic Waiting Protocol

This skill encapsulates the guidelines and workflow updates for CI/CD pipeline monitoring and diagnostic fetching across all orchestration, RCA, and execution prompts.

## Core Protocols

1. **GitMap Pipeline-AI Section Integration:**
   - Use `gitmap pipeline-ai status --json` (or alias `gitmap pl-ai status --json`) to query remote CI/CD pipeline status.
   - Parse `is_running`, `etaSeconds`, and `nextAiCommand`.
   - Use GitMap's automated error extraction to isolate actionable failure lines (`##[error]`, `FAIL:`, compile errors) for 4-part RCA without reading noisy passing step logs.

2. **Adaptive Waiting / Sleep (Anti-Credit-Waste Mandate):**
   - NEVER loop rapidly or busy-poll (`gh run view` in tight loops). This burns user credits and API token quota.
   - When a pipeline is running, use `gitmap pipeline-ai status -t <etaSeconds>` to sleep for the estimated completion duration before querying again.
   - Alternatively, schedule wait intervals or sleep based on the computed ETA:
     - ETA > 120s: wait 20s-30s
     - 60s < ETA <= 120s: wait 10s-20s
     - ETA <= 60s: wait 5s-10s
