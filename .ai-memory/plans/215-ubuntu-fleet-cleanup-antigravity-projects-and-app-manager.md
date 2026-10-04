# Plan: 215-ubuntu-fleet-cleanup-antigravity-projects-and-app-manager

## User Request (Verbatim)
```text
Hi there. I really liked what you have done. I really like the background now, what we have right now. So that is all right. So far your work is good, but I have some feedbacks and improvements that we have to work on. So let's get into one by one. First, inside the Git work, there is something like a tilde folder that I wanted to remove. You can get in and try to remove that folder properly. That's one big issue I think. The next one inside the home section or home folder, we have a bunch of downloads. Let me highlight this. These downloads usually should be done using a temp directory, and after the download and installation, this should be removed. But it was not removed. And repo cache and repo secret should be created inside the git work, where is our work directory is, but it was created on the root level. Try to move, remove, and make sure there is only one. So gitmap needs to be synced as well. So all these things, these are new steps. Make sure that you add it to the next steps and also make sure that these are fixed. After that, now you have the logs. Now you know how to install. How to do things? I want to make sure that the project that we have inside the gitmap, it's not being loaded into the Antigravity. I wanted to load into the Antigravity so that all the projects that we have can be listed in the Antigravity and from the Antigravity. What we have done earlier, we have loaded all the conversations and all these things. Those are not connected to the projects, and I want those to be connected to the projects properly. So everything should be visible properly. Then, in the Ubuntu, to uninstall an app is really hard. Can you make something to uninstall apps that we have installed? Or maybe you can see what is installed and from there we can uninstall easily. Like the yellow icon, Antigravity tools. Can you remove it from the system? We don't want to have it. I want to test this. And also in the Windows, uninstall clot. So these are the test cases for both platforms. Then we don't want to pass very long commands that you are writing. It's really hard to run. Why not make it a `gitmap` command so that we can easily pass it to the other machine with a json response? If we are going to install multiple tools, we can pass comma separated tools like: `gitmap install --tools antigravity,chrome,vscode,flameshot` and it should install all of them. Also we need a way to migrate the whole project and settings and conversations from one machine to another in 1 single command. If we can run a command with `--json`, it should return json. So all these things, can you make it as a help text or documentation or something in the terminal so that anyone can know what commands to run?
```

## Architecture & Objectives
1. **Node `u1` Deep Filesystem Hygiene**:
   - Safely remove `/home/a/git-work/~` (root-owned Oh My Zsh artifact).
   - Purge loose downloads, credential JSONs, and scripts from `/home/a`.
   - Remove literal Windows directory `/home/a/C:\Users\Administrator\AppData\Roaming\Antigravity`.
   - Remove empty `/home/a/repo-secrets` and `/home/a/repo-cache`; guarantee canonical copies live exclusively at `/home/a/git-work/repo-secrets` and `/home/a/git-work/repo-cache`.
2. **Antigravity Projects & Conversations Mapping Engine**:
   - Generate project JSON definitions under `/home/a/.gemini/config/projects/<project-id>.json` for all 74 workspace repositories in `/home/a/git-work/` with matching GUIDs and Linux folder URIs.
   - Ensure `conversation_summaries.db` links conversation IDs to project IDs, eliminating "No Project" state.
3. **App Management Subsystem (`gitmap apps`)**:
   - Implement `gitmap apps list` (parsing XDG `.desktop` entries on Linux and Winget/Registry on Windows).
   - Implement `gitmap apps uninstall <app>` with `--purge` and `--json`.
   - Test by purging legacy yellow `antigravity-tools` v4.7.2 on `u1` and `clot` on Windows.
4. **Modular Installation & Cross-Node Delegation**:
   - Add `--tools` and comma-separated tool install support to `gitmap install`.
   - Add `--json` flag support to fleet and migration commands for machine-to-machine task handoff.
   - Comprehensive terminal help text and command catalog.
5. **Verification Scorecard & Zero-Build Commit**:
   - Run multi-point health check across `u1` and Windows.
   - Commit cleanly via `gitmap cpf "apps - implement cross-platform app manager and fleet governance"`.

## Discrete Subtasks
- **Subtask 01**: `01-ubuntu-node-u1-deep-filesystem-hygiene.md` — Eradicate `/home/a/git-work/~`, purge home root artifacts, clone/init `/home/a/git-work/repo-cache`.
- **Subtask 02**: `02-antigravity-projects-and-workspaces-registration.md` — Generate and stream 74 project JSONs to `/home/a/.gemini/config/projects/`, update conversation DB mappings.
- **Subtask 03**: `03-app-uninstaller-subsystem-and-legacy-purge.md` — Implement `gitmap apps list/uninstall`, test purging `antigravity-tools` on `u1` and `clot` on Windows.
- **Subtask 04**: `04-compact-gitmap-commands-and-delegation-protocol.md` — Modular `--tools` installer in GitMap, `--json` cross-node delegation, terminal help.
- **Subtask 05**: `05-live-verification-scorecard-and-atomic-push.md` — Consolidated audit ledger, scorecard validation, atomic GitMap push.
