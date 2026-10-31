# Specification: AI Heatmap File Touch & Analysis Notification

## 1. Overview
Multi-agent concurrency requires that when an AI agent touches, reads, or analyzes a file to diagnose an issue, other agents and developers know that the file is currently under investigation. GitMap Heatmap tracks these active touches in Split-DB (`repodb/heatmap.db`).

## 2. CLI Contract
- **Command:** `gitmap heatmap touch <filepath> [--reason "<why>"] [--agent <name>]`
- **Aliases:** `gitmap ht touch`, `gitmap touch-log`
- **Behavior:**
  - Records a touch event in `heatmap.db` with timestamp, relative filepath, caller agent persona, and reason.
  - Automatically flags active inspection in `gitmap heatmap live` and `gitmap heatmap status`.
- **Auto-touch with `gitmap cat`:**
  - `gitmap cat <filepath> --touch` automatically emits a heatmap touch event before streaming file contents.

## 3. Mandatory AI Workflow
Whenever an AI agent begins analyzing a target file:
```bash
gitmap heatmap touch <target-path> --reason "Investigating SQL query execution and database aliases"
```
This guarantees zero file collision and complete cognitive auditability.
