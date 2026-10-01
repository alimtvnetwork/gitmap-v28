# Subtask 02: Fix Ignore All

## Description
Implement and refine the `gitmap fix-ignores-all-ssh (fias)` and local `fix-ignore-all (fia)` commands to use the centralized ignore management suite and PAS formula.

## Objectives
- Implement `fias` to follow the PAS Formula for remote fleet execution.
- Add support for `-y` (auto-confirm) to bypass interactive prompts.
- When `-y` is absent, show interactive resolution proposals grouped by repository.
- Ensure that the remediation process deduplicates `.gitignore` patterns and applies standardized ignore rules (e.g., ignoring `.gitmap/`).

## Acceptance Criteria
- `gitmap fia -y` fixes all ignore issues locally without prompts.
- `gitmap fias` uses the PAS formula and updates the `TaskQueue` accurately.
- Interactive mode presents a clean summary and options for batch remediation.
