---
name: create-run-ps1-file
description: >-
  Autonomously architect, construct, and verify cross-platform execution runners (run.ps1 and run.sh) driven by run.config.json and standalone installation scripts (local-install.ps1 and local-install.sh) with GitMap integration and automatic rerun.
---

# Create Run & Install Scripts Skill (`create-run-ps1-file`)

This skill autonomously architects, constructs, and verifies cross-platform execution runners (`run.ps1` and `run.sh`) driven by a central `run.config.json` manifest, along with accompanying standalone environment installation scripts (`local-install.ps1` and `local-install.sh`). It implements robust dependency pre-flight checks, intelligent fallback to GitMap (`gitmap aum install`) or `local-install.ps1`/`local-install.sh` with automatic rerun, and clean process cleanup on exit.

## Architecture & Lifecyle

```mermaid
flowchart TD
    A["run.ps1 / run.sh Execution"] --> B["Parse run.config.json"]
    B --> C{"Toolchains & Dependencies Installed?"}
    C -- "Yes" --> D["Execute Target Command / Start Services"]
    C -- "No" --> E{"Is GitMap CLI Available?"}
    E -- "Yes" --> F["Run gitmap aum install"]
    E -- "No" --> G["Run .\\local-install.ps1 / ./local-install.sh"]
    F --> H["Environment Ready"]
    G --> H
    H --> I["Self-Rerun: run.ps1 / run.sh with Original Args"]
    I --> D
    D --> J["Process Cleanup on Exit (Trap/Finally)"]
```

## Implementation Standards

1. **Manifest Driven (`run.config.json`):**
   - Centralize all ports, paths, and commands.
   - Support token replacement (e.g. `{fePort}`, `{bePort}`).
2. **Self-Healing Fallback:**
   - Detect missing dependencies at runtime.
   - Check for GitMap environment manager first (`gitmap aum install`).
   - Fall back to `local-install.ps1` / `local-install.sh`.
   - Re-verify and automatically re-execute `run.ps1`.
3. **CI Mode Support (`-CI` / `--ci`):**
   - Non-interactive execution of build, lint, and test suites.
   - Immediate exit 1 on any failure.
4. **Clean Process Management:**
   - Track spawned PIDs and terminate in `finally` / `trap` blocks.
