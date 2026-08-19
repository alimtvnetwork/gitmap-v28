---
name: Execute Pending Tasks
description: Autonomous execution engine for `.lovable/plans/pending/` tasks.
---

# Execute Pending Tasks (Autonomous Execution)

**Goal:** Execute all pending tasks without stopping or asking the user for permission.

## 1. Initial State: Audit and Re-Sequence Pending Tasks
- Read the `.lovable/plans/pending/` directory and `.lovable/plans/index.md`.
- Count exactly how many pending tasks exist.
- **Naming Correction**: Check if the pending task files are correctly sequenced with a 2-digit numerical prefix (e.g., `01-<slug>.md`, `02-<slug>.md`). 
- If the naming is incorrect or missing prefixes, **fix it immediately**. Rename the files to follow the sequential `01-`, `02-`, `03-` format and update `.lovable/plans/index.md` to match the new filenames in the same operation.

## 2. Uninterrupted Autonomous Execution (Self-Looping)
- **Make a Great Plan**: Analyze all pending tasks and devise a comprehensive execution plan. 
- **Do NOT Ask Questions**: Do not stop to ask the user for permission. Do not stop to ask clarifying questions. 
- **Self-Loop**: You must self-loop continuously until every single pending task in the queue is verifiably completed.
- Ensure everything gets done properly and deeply. If you are not going deep, you are not doing the job.

## 3. High-Stakes Code Standards & Root Cause Analysis
- **Root Cause First**: Find the root cause of every problem before applying any fix. Record the root cause into `.lovable/memory/` before touching code.
- You MUST follow the project's strict coding guidelines. These files are located in the `spec/` folder.
