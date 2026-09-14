---
name: macro-resilience-trees-and-direct-dispatch
description: Autonomously implement and verify macro run-until error tolerance, dedicated failure logging, no-terminal suppression, execution summary trees, no-tree toggles, and direct root CLI macro dispatch.
---

# Macro Resilience, Execution Summary Trees & Direct Dispatch

This skill guides the autonomous implementation and verification of:
1. **Macro Execution Resilience (`--run-until`, `run-until`)**:
   - Continuing execution when individual macro steps fail.
   - Recording step failure details (command, exit code, elapsed time, stderr diagnostics, timestamp) into dedicated macro log files.
   - Reporting the log file path and failure reason to the user.
2. **Terminal Output Suppression (`--no-terminal`)**:
   - Running macros without polluting the terminal stdout/stderr, while retaining full internal output and logging capabilities.
3. **Execution Summary Trees (`--summary`, `--no-tree`, default tree)**:
   - Always displaying a visual execution tree summary at the end of macro runs by default.
   - Supporting `--summary` to focus solely on the execution tree summary.
   - Supporting `--no-tree` to omit the tree summary when requested.
4. **Direct Top-Level Macro Dispatch (`gitmap <macro-name>`, `gitmap run <macro-name>`)**:
   - Seamlessly resolving saved macros at the root CLI level so running `gitmap <macro-name>` or `gitmap run <macro-name>` executes the macro rather than returning `E1001: Unknown command`.
