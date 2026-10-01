# Subtask 03: Commit Push All

## Description
Develop the `gitmap commit push all` functionality to batch commit and push operations across all tracked repositories, with robust review mechanisms.

## Objectives
- Scan all tracked repositories for modified, untracked, or deleted files.
- Add an interactive `--review` (`-r`) flag to summarize file changes as a folder tree before execution.
- Implement automated fallback commit prefixes (e.g., `Feature:`, `Bug:`, `Release:`) when requested.
- Track all resulting commit and push steps in the `TaskQueue` for full historical visibility.

## Acceptance Criteria
- `gitmap commit push all --review` prompts the user with a detailed change summary before committing.
- Commit failures (e.g., conflict or hook failure) are gracefully caught and recorded in the GitMap errors DB.
- Operations respect task servers integration, enqueueing tasks prior to execution.
