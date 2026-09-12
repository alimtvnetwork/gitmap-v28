# Execute Pending Tasks (Autonomous Execution)

## 1. Initial State: Audit and Re-Sequence Pending Tasks

Before starting execution, you must ensure the queue of pending tasks is properly ordered and named.

- Read the `.lovable/plans/pending/` directory and `.lovable/plans/index.md`.
- Count exactly how many pending tasks exist.
- **Naming Correction**: Check if the pending task files are correctly sequenced with a 2-digit numerical prefix (e.g., `01-<slug>.md`, `02-<slug>.md`).
- If the naming is incorrect or missing prefixes, **fix it immediately**. Rename the files to follow the sequential `01-`, `02-`, `03-` format and update `.lovable/plans/index.md` to match the new filenames in the same operation.

## 2. Uninterrupted Autonomous Execution (Self-Looping)

You are the sole orchestrator. Your job is to complete ALL pending tasks without stopping.

- **Make a Great Plan**: Analyze all pending tasks and devise a comprehensive execution plan.
- **Do NOT Ask Questions**: Do not stop to ask the user for permission. Do not stop to ask clarifying questions.
- **Self-Loop**: You must self-loop continuously until every single pending task in the queue is verifiably completed.
- Ensure everything gets done properly and deeply. If you are not going deep, you are not doing the job.

## 3. High-Stakes Code Standards & Root Cause Analysis

While executing the pending tasks, you must adhere strictly to the project's code standards and root cause protocols:

- **Root Cause First**: Find the root cause of every problem before applying any fix. Record the root cause into the .lovable memory before touching code.

You MUST follow the project's strict coding guidelines. These files are located in the `spec/` folder.

## 4. Runner In-Flight ETA Wait Protocol & Sleep Discipline

When running tests, linters, or CI/CD pipelines in the background:
- **25-Second Heartbeat Cadence**: Runner progress heartbeats occur strictly every 25 seconds or more. Never spam or query status quickly.
- **1-Minute AI Sleep Discipline**: When inspecting or waiting for background runner execution, agents MUST sleep/wait for **1 minute (60 seconds) each time**, or dynamically sleep based on the remaining ETA read from `.lovable/temp/runner-eta.json` (or based on previous total approximate delay) instead of busy-polling.
- **Reactive Notification**: Once launched, allow background tasks to run and await the system's reactive message notification instead of looping on status calls.

