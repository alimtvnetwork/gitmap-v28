---
name: agy-prompts-templates-and-rerun-suite
description: Autonomously implement and verify Antigravity (AGY) prompt replay, prefix templates management, prompt listing with multi-project VS Code opening, and remote delegation across SSH, Cluster, SC, and local execution.
---

# AGY Prompts Templates and Rerun Suite

## Overview
Autonomously manage, template, and rerun Google Antigravity (AGY) prompts across local and remote nodes (SSH, Cluster, SC), with prompt history inspection, prefix templating, and non-admin VS Code workspace opening.

## Core Capabilities
1. `gitmap agy rerun last [N] [-p <template>]`:
   - Reruns the last N prompts with optional prefix prompt template.
   - Built-in default template: `is-done` ("Is it done properly? Can we check properly the missing items from the task that is mentioned below? Please check it carefully. Do not make any mistakes.").
2. `gitmap agy list-prompts [N] [--all-projects] [--projects <N>] [--project <name>]`:
   - List historical prompts per project or across all projects.
   - Opens prompt diffs/summaries in a new non-admin VS Code window from temporary files without replacing existing windows.
3. `gitmap agy scan`:
   - Scans projects and recent prompt activity counts.
4. `gitmap prompts-template [add|edit|ls|rm|import|export|import-all|export-all]`:
   - JSON-backed template registry for custom prompt prefixes and reusable verification prompts.
5. Remote Triad Delegation:
   - `gitmap ssh exec agy ...`
   - `gitmap sc exec agy ...`
   - `gitmap cluster exec agy ...`
   - `gitmap exec agy ...`
6. Visual & Storage Hygiene:
   - Storage restore DB command (`gitmap storage restore-db`).
   - Show repo DBs in `gitmap storage ls`.
   - Conditional terminal padding with intelligent detection.
