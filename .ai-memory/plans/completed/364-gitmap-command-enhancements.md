# Plan 364: GitMap Command Enhancements — Multi-Node Deploy, Copy (CP), Move (MV) & Fleet Software Updates

- **Plan Slug:** `364-gitmap-command-enhancements`
- **Specification:** [02-spec/21-app/364-gitmap-command-enhancements/00-overview.md](../../../02-spec/21-app/364-gitmap-command-enhancements/00-overview.md)
- **Status:** `completed`
- **Creation Date:** 31-Oct-26
- **Architecture Lead:** Root Orchestrator & Swarm (A=3, H=2, Capacity=6)

---

## 1. User Request (Verbatim)

```text
The `gitmap`, I wanted to have the deploy, copy, and move command, which I could not find in the help text. The deploy command would have copy binaries and also with the installation. For example, Windows, we should be able to copy binaries to just all the machines using copy command. And there is one more thing that I thought of we are doing a little bit wrong because usually the copy command is CP, and we are doing CP for commit pull or commit push. I think that needs to be corrected everywhere, along with the commands, examples, help text, UI help text and everything. So CP should be reserved for copy because people know CP.copy, and MV for move file from one machine to another. So create separate syntaxes. There should be four to eight different syntax for deploying a copy a file for all nodes, single node, a specific node from this current folder file. Also move file from one node to another or move file from here to all the nodes, things like that. So all types of commands, try to share that at the end that you have implemented. Make CP as copy command for the `gitmap` deployed and also with the deploy, we should be able to install the binary, regardless of the OS. If it is a EXE on the Windows and we are deploying, then we can also say skip Unix or skip Ubuntu. It would automatically skip that. So by default, EXE files can be skipped from Ubuntu machines when we do the deploy copy. And it can have a flag like force copy, then it would automatically do that, and also it would provide those suggestions after it runs, and it would mention what it is doing and why it is keeping the Ubuntu machine, things like that, so that it is efficient. And only the open machines or connected machines, those are running that can copy to remember that. And installation, if it is a EXE, so for example, the cursor update, it does not update automatically. So we should have a command to update cursor. In all the machines, we should have a command for the answer nodes, update cursor or VS code or `Antigravity Manager or Antigravity Tools`, things like that. So try to make sure that you work on it, you implement that, and also do a minor bump and release.
```

---

## 2. Core Execution Phases & Subtasks

### Subtask 01: Alias Disambiguation & Multi-Node File Copy & Move Engine
- **Spec:** [01-alias-and-copy-move-spec.md](../../../02-spec/21-app/364-gitmap-command-enhancements/01-alias-and-copy-move-spec.md)
- **Subtask Plan:** [subtasks/364-gitmap-command-enhancements/01-alias-and-copy-move.md](../subtasks/364-gitmap-command-enhancements/01-alias-and-copy-move.md)
- **Scope:**
  - Reassign `cp` to `copy` across CLI, constants, dispatchers, completions, help documents, and UI (`src/data/commands.ts`).
  - Reserve `mv` exclusively for `move`.
  - Preserve all semantic commit verbs (`cpf`, `cpb`/`cbf`, `cpr`, `cpc`, `cfr`).
  - Implement 8 syntax variants for `gitmap copy` / `gitmap cp`.
  - Implement extended syntax for `gitmap move` / `gitmap mv` with node detection and atomic remote cleanup.

### Subtask 02: OS-Aware Binary Deploy & Installation Engine
- **Spec:** [02-deploy-and-os-compat-spec.md](../../../02-spec/21-app/364-gitmap-command-enhancements/02-deploy-and-os-compat-spec.md)
- **Subtask Plan:** [subtasks/364-gitmap-command-enhancements/02-deploy-and-os-compat.md](../subtasks/364-gitmap-command-enhancements/02-deploy-and-os-compat.md)
- **Scope:**
  - Enhance `gitmap deploy` with binary magic byte inspection (`MZ` for PE, `ELF` for Linux).
  - Filter out incompatible OS targets (e.g. skip `.exe` on Ubuntu) unless overridden with `--force` or `--force-copy`.
  - Fast concurrent TCP liveness check to skip offline nodes without delay.
  - Rich post-execution diagnostic docket with next-step suggestions.

### Subtask 03: Cross-Node Fleet Software Updates & UI/Doc Synchronization
- **Spec:** [03-fleet-software-update-spec.md](../../../02-spec/21-app/364-gitmap-command-enhancements/03-fleet-software-update-spec.md)
- **Subtask Plan:** [subtasks/364-gitmap-command-enhancements/03-fleet-software-update.md](../subtasks/364-gitmap-command-enhancements/03-fleet-software-update.md)
- **Scope:**
  - Implement fleet update commands for Cursor IDE, VS Code, AGM, and AGY tools across nodes.
  - Update `src/data/commands.ts` and documentation pages.
  - Regenerate shell completions.

### Subtask 04: Minor Version Bump & Release Ceremony
- **Scope:**
  - Bump minor version to `v6.639.0`.
  - Update `version.json`, `package.json`, `changelog.md`, and `.ai-memory/release/release-notes-v6.639.0.md`.
  - Create semantic atomic commit.

