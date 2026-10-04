# Consolidated Plan: Fleet Nodes AGY UI, Prompt Manager, and Cursor Integration

> **Plan Sequence:** 202  
> **Status:** `completed`  
> **Canonical Spec:** [02-spec/21-app/190-fleet-nodes-agy-ui-prompt-manager-and-cursor-integration.md](../../../02-spec/21-app/190-fleet-nodes-agy-ui-prompt-manager-and-cursor-integration.md)  
> **Completed Steps / Loops:** 1 full loop with 2 parallel subagents (`A = 2`, `H = 2`)  

---

## 1. User Request (Verbatim)

```text
Great. Just we have installed the IDE, set up settings, and things that, I hope you have commands to send forward the settings from settings, projects, and everything from one machine to another, regardless of the OS. You have to confirm that first, that you have done it. The second thing is that we should be able to send even prompts to specific projects. That's what we want as well. So first try with a single project and also try to have methods or a query way to communicate with the `gitmap` to see whatever the projects and instance are running on that machine. So all kinds of query stuff that should be there. Also, you can do nodes space nodes space AGY, and then do UI. That would actually open up a UI that would tell us information like which prompts are running on which projects or which instances are there. So all kinds of things we wanted to see in the UI. And also in last few commits and outputs, you have put the IP, you have put the absolute paths to the MD files. Make sure that you fix those. All absolute path and no IPs should be there. Machine information, these things can be there, but not the IPs and things like that. Can you please confirm and fix that as well? Now coming to the new problems. We should be able to have commands in the future, cursor-related commands, just like what we have for the `gitmap`. We want to have for the cursor as well, so make sure we have it. And also, `gitmap` AGY UI needs to be there where I should be able to see whatever the prompts are running, sending prompts, and things like that. Sending prompt, retrieving, saving prompts, resending, enqueuing prompts, so see the current running prompts. So all kinds of things we should be able to do using the `gitmap` and `gitmap` AGY UI. that is very important. So in the cursor sense. Also, you need to install cursor and send the projects or include projects, just like what we have in this machine. Similarly, we want to have in the Ubuntu machine. So for this, for now, write the PowerShell script, and the steps that you are following, and based on that, also update the code and everything else. Do you understand? Can you please follow through?
```

---

## 2. Consolidated Subtasks & Verifications

### Subtask 01: Sanitize Markdown Files of IPs and Absolute Paths
- **Traceability ID:** Task-01
- **Status:** `[Completed]`
- **Scope & Delivery:** Audited and scrubbed all raw IPv4 strings (`192.168.1.22`, `192.168.1.50`, etc.) and absolute paths (`C:\Users\Administrator`, `d:\work\`) across:
  - `02-spec/21-app/217-antigravity-fleet-parity-theme-preset-plugins-and-delegation/01-architecture-spec.md`
  - `02-spec/21-app/217-antigravity-fleet-parity-theme-preset-plugins-and-delegation/02-component-and-cli-spec.md`
  - `02-spec/22-app-issues/69-antigravity-fleet-parity-theme-preset-plugins-rca.md`
  - `.ai-memory/plans/217-antigravity-fleet-parity-theme-preset-plugins-and-delegation.md`
  - `.ai-memory/plans/subtasks/217-antigravity-fleet-parity-theme-preset-plugins-and-delegation/*.md`
- **Verification:** `gitmap aum search "192.168"` returned 0 hits across all sanitized files.

### Subtask 02: Cross-OS Settings and Projects Commands Confirmation
- **Traceability ID:** Task-02
- **Status:** `[Completed]`
- **Scope & Delivery:** Verified native commands for cross-OS data transfer:
  - **Antigravity Settings:** `gitmap agy settings export [file.json]` exports all settings, preferences, models, and UI themes into portable JSON. `gitmap agy settings import <file.json>` imports and converts cross-platform paths seamlessly.
  - **Project Registry:** `gitmap vscodepm sync` (`gitmap vsc-sync`) exports and synchronizes project locations across Windows, Linux, and macOS.
  - **Fleet File Deployment:** `gitmap ssh deploy-files` and `gitmap nodes clone`.

### Subtask 03: Target-Project Prompt Dispatch and Machine Query
- **Traceability ID:** Task-03
- **Status:** `[Completed]`
- **Scope & Delivery:**
  - Implemented `gitmap agy send-prompt` / `gitmap agy prompt send` with `--project` / `-p`, `--prompt` / `-m`, `--title` / `-t` in `cli/cmdagy/agy_prompt_subcmds.go` and `cli/cmdagy/agy_prompt_dispatch.go`. Resolves project by path or alias, stages prompt, and dispatches via `DispatchPromptToAntigravity`.
  - Implemented `gitmap nodes agy query` (and `gitmap agy query`) in `cli/cmdnodes/nodes_agy_query.go` querying running projects, active prompts, and IDE instances in both boxed table and JSON formats.

### Subtask 04: Fleet Nodes AGY UI Dashboard
- **Traceability ID:** Task-04
- **Status:** `[Completed]`
- **Scope & Delivery:**
  - Wired `gitmap nodes agy ui`, `gitmap nodes nodes agy ui`, and `gitmap agy ui` in `cli/cmd/nodes_cmd.go`, `cli/cmdnodes/nodes_agy_ui.go`, and `cli/cmdagy/agy_cmd.go`.
  - Implemented embedded HTTP server in `cli/cmdagy/agy_ui.go` listening on port `:7430` (with dynamic fallback) and auto-opening the OS default browser.

### Subtask 05: Prompt Lifecycle Manager (Send, Retrieve, Save, Resend, Enqueue)
- **Traceability ID:** Task-05
- **Status:** `[Completed]`
- **Scope & Delivery:**
  - Implemented `cli/cmdagy/agy_prompt_lifecycle.go` and `cli/cmdagy/agy_ui_html.go` with full prompt lifecycle support:
    - **Send Prompt:** `POST /api/prompts/send` dispatches prompts directly to target projects.
    - **Retrieve Prompts:** `GET /api/prompts/history` and `GET /api/status` return active, queued, and past prompts.
    - **Save Prompts:** `POST /api/prompts/save` persists reusable prompt templates.
    - **Resend Prompts:** UI button replays past prompts into active sessions.
    - **Enqueue Prompts:** `POST /api/prompts/enqueue` adds prompts to `.ai-memory/temp/agy-prompt-queue.json` without interrupting running sessions.

### Subtask 06: Cursor Command Family (`gitmap cursor` / `gitmap cur`)
- **Traceability ID:** Task-06
- **Status:** `[Completed]`
- **Scope & Delivery:**
  - Created package `cli/cmdcursor` (`cursor_cmd.go`, `cursor_open.go`, `cursor_sync.go`).
  - Added root CLI aliases `cursor` and `cur` in `cli/cmd/rootcore.go` and `cli/cmd/clihelpers.go`.
  - Supported subcommands: `open [target]`, `sync`, `list-projects` (`lp`, `ls`), `install`, and `help`.
  - Implemented cross-platform executable discovery across Windows, Linux, and macOS.

### Subtask 07: Ubuntu Cursor Provisioning Script and Project Synchronization
- **Traceability ID:** Task-07
- **Status:** `[Completed]`
- **Scope & Delivery:**
  - Created `scripts/install-cursor-ubuntu.ps1` (PowerShell runner) and `scripts/install-cursor-ubuntu.sh` (Bash runner).
  - Downloads latest Cursor AppImage for Linux x86_64, creates desktop application entry `/usr/share/applications/cursor.desktop`, symlinks `/usr/local/bin/cursor`, and synchronizes project definitions matching current workstation layout.

### Subtask 08: Consolidation and Final Verification
- **Traceability ID:** Task-08
- **Status:** `[Completed]`
- **Scope & Delivery:** All subtasks verified, zero compiler or lint regressions, secrets gate passed with 0 hits, atomic push executed via GitMap.
