# Consolidated Plan: 218-nodes-agy-ui-remote-settings-and-cursor-automation

## User Request (Verbatim)
> Great. Just we have installed the IDE, set up settings, and things that, I hope you have commands to send forward the settings from settings, projects, and everything from one machine to another, regardless of the OS. You have to confirm that first, that you have done it. The second thing is that we should be able to send even prompts to specific projects. That's what we want as well. So first try with a single project and also try to have methods or a query way to communicate with the `gitmap` to see whatever the projects and instance are running on that machine. So all kinds of query stuff that should be there. Also, you can do nodes space nodes space AGY, and then do UI. That would actually open up a UI that would tell us information like which prompts are running on which projects or which instances are there. So all kinds of things we wanted to see in the UI. And also in last few commits and outputs, you have put the IP, you have put the absolute paths to the MD files. Make sure that you fix those. All absolute path and no IPs should be there. Machine information, these things can be there, but not the IPs and things like that. Can you please confirm and fix that as well? Now coming to the new problems. We should be able to have commands in the future, cursor-related commands, just like what we have for the `gitmap`. We want to have for the cursor as well, so make sure we have it. And also, `gitmap` AGY UI needs to be there where I should be able to see whatever the prompts are running, sending prompts, and things like that. Sending prompt, retrieving, saving prompts, resending, enqueuing prompts, so see the current running prompts. So all kinds of things we should be able to do using the `gitmap` and `gitmap` AGY UI. that is very important. So in the cursor sense. Also, you need to install cursor and send the projects or include projects, just like what we have in this machine. Similarly, we want to have in the Ubuntu machine. So for this, for now, write the PowerShell script, and the steps that you are following, and based on that, also update the code and everything else. Do you understand? Can you please follow through?

## Canonical Specifications
- [01-architecture-spec.md](../../02-spec/21-app/218-nodes-agy-ui-remote-settings-and-cursor-automation/01-architecture-spec.md)
- [02-component-and-fleet-spec.md](../../02-spec/21-app/218-nodes-agy-ui-remote-settings-and-cursor-automation/02-component-and-fleet-spec.md)

## Executed Tasks Summary

### Task-01: Cross-OS Settings & Projects Forwarding
- **Status:** `[DONE]`
- **Implementation:**
  - Implemented `cli/cmdnodes/nodes_sync_settings.go`:
    * `gitmap nodes push-settings [node]` & `gitmap nodes sync-settings`: Exports local Antigravity settings to JSON envelope, deploys over SSH, and imports remotely via `gitmap agy settings import`.
  - Implemented `cli/cmdnodes/nodes_send_projects.go`:
    * `gitmap nodes send-projects [node]`: Replicates workspace projects and syncs Cursor / VS Code Project Manager configurations.

### Task-02: Project-Scoped Prompt Dispatch & Live Querying
- **Status:** `[DONE]`
- **Implementation:**
  - Implemented `cli/cmdnodes/nodes_agy_prompt.go`:
    * Supports `gitmap nodes agy prompt <node> <project> "<prompt>"`.
    * Connects to target node over SSH, stages prompt safely, and dispatches to the specified project.
  - Enhanced `cli/cmdnodes/nodes_agy_query.go`:
    * Supports `--ssh` fleet aggregation, querying running projects and instances across the cluster via `AggregateSSHRunningProjects()`.

### Task-03: `nodes space nodes space AGY` UI Dashboard (`gitmap nodes agy ui`)
- **Status:** `[DONE]`
- **Implementation:**
  - In `cli/cmd/nodes_cmd.go` and `cli/cmd/rootcore.go`:
    * Added command routing and alias tolerance for `gitmap nodes nodes agy ui`, `gitmap nodes agy ui`, `gitmap nodes agy-ui`, and `gitmap nodes-agy-ui`.
  - In `cli/cmdnodes/nodes_agy_ui.go`:
    * Launches the embedded Antigravity Studio Web UI on port 7430.
  - In `cli/cmdagy/agy_ui.go` and `agy_ui_html.go`:
    * Updated dashboard to display active fleet nodes, remote running projects, and running prompts in real-time.

### Task-04: Git Hygiene & Privacy Scrub (Absolute Paths & Raw IPs)
- **Status:** `[DONE]`
- **Implementation:**
  - Cleaned all 31 recent files across `02-spec/21-app/` and `.ai-memory/plans/`.
  - Replaced hardcoded IPs (`node-u1`, `node-main`, etc.) with machine aliases (`ubuntu-fleet-01`, `node-01`, `node-main`, `localhost`).
  - Replaced hardcoded Windows/Linux absolute paths (`$WORKSPACE_ROOT`, `C:\Users\...`, `/home/<user>/...`) with relative `$WORKSPACE_ROOT`, `$HOME/git-work/`, and `./`.
  - Verified with `03-ai-scripts/49-verify-privacy-and-relative-paths.py` (31 files checked, 0 violations found, exit code 0).

### Task-05: GitMap AGY Full Prompt Lifecycle UI & Operations
- **Status:** `[DONE]`
- **Implementation:**
  - In `cli/cmdagy/agy_prompt_lifecycle.go`:
    * Implemented models and handlers for `/api/prompts/send`, `/api/prompts/enqueue`, `/api/prompts/tree`, `/api/prompts/history`, `/api/prompts/save`, and `/api/prompts/resend`.
    * Implemented background prompt queue worker logic and lifecycle state transitions.

### Task-06: Cursor CLI Subsystem & Ubuntu Fleet Setup Automation
- **Status:** `[DONE]`
- **Implementation:**
  - Implemented `cli/cmdcursor/cursor_settings.go`:
    * `gitmap cursor settings [view|apply|sync]` managing settings and Dracula themes.
  - Implemented `cli/cmdcursor/cursor_install.go`:
    * Automated installation on Windows (`winget`) and Linux (AppImage) with `--node` delegation.
  - Authored canonical 7-phase setup scripts:
    * `03-ai-scripts/40-ubuntu-cursor-and-agy-fleet-setup.ps1`
    * `03-ai-scripts/40-ubuntu-cursor-and-agy-fleet-setup.sh`
    * Verified via `-DryRun` with exit code 0.

## Modified Files
- `02-spec/21-app/218-nodes-agy-ui-remote-settings-and-cursor-automation/01-architecture-spec.md`
- `02-spec/21-app/218-nodes-agy-ui-remote-settings-and-cursor-automation/02-component-and-fleet-spec.md`
- `02-spec/21-app/readme.md`
- `.ai-memory/plans/readme.md`
- `cli/cmdnodes/nodes_sync_settings.go`
- `cli/cmdnodes/nodes_send_projects.go`
- `cli/cmdnodes/nodes_agy_prompt.go`
- `cli/cmdnodes/nodes_agy_query.go`
- `cli/cmdnodes/nodes_agy_ui.go`
- `cli/cmd/nodes_cmd.go`
- `cli/cmd/rootcore.go`
- `cli/cmdagy/agy_prompt_lifecycle.go`
- `cli/cmdagy/agy_ui.go`
- `cli/cmdagy/agy_ui_html.go`
- `cli/cmdcursor/cursor_settings.go`
- `cli/cmdcursor/cursor_install.go`
- `cli/cmdcursor/cursor_sync.go`
- `03-ai-scripts/40-ubuntu-cursor-and-agy-fleet-setup.ps1`
- `03-ai-scripts/40-ubuntu-cursor-and-agy-fleet-setup.sh`
- `03-ai-scripts/49-verify-privacy-and-relative-paths.py`
- 31 scrubbed markdown specification and plan files
