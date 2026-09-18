---
name: pipeline-errors-agy-fix-injection-and-multi-project-batching
description: Autonomously implement and verify Antigravity (AGY) CI/CD pipeline error fix injection, embedding full failing logs, full file path terminal display, direct prompt dispatch into Antigravity IDE/CLI, queue management, and multi-project parallel batching with configurable project limits.
---

# Pipeline Errors AGY Fix Injection and Multi-Project Batching

## Overview
Autonomously diagnose, repair, and verify `gitmap pipeline errors agy fix` and `gitmap agy fix-pipeline` across single-repo and multi-project contexts:
1. Ensure the generated fix prompt file contains the complete, failing CI/CD error logs alongside the RCA fix instructions.
2. In terminal outputs, display strictly full absolute paths for all generated payload, queue, and ledger files.
3. Automatically inject/dispatch the prompt payload into Antigravity (`agy` CLI or IDE session), not just copy to clipboard.
4. Support multi-project parallel pipeline error scanning, collection, and AGY task enqueueing (default limit: 3 projects per batch, subsequent batches for remaining projects).
5. Update terminal help text, documentation (`cli/helptext/pipeline.md`, `cli/helptext/agy-fix-pipeline.md`), Web UI documentation, and datasets.
