# 364 — GitMap Command Enhancements: Multi-Node Deploy, Copy (CP), Move (MV) & Software Updates

> **Specification Sequence:** `364`  
> **Canonical Path:** `02-spec/21-app/364-gitmap-command-enhancements/00-overview.md`  
> **Parent Plan:** `.ai-memory/plans/pending/364-gitmap-command-enhancements.md`  
> **Author:** GitMap Architecture & Swarm Team  
> **Date:** 31-Oct-26  
> **Status:** `active`  

---

## 1. Executive Summary

This specification module defines the architecture, CLI contracts, and execution mechanics for GitMap's next-generation fleet command enhancements:
1. **Reclaiming `cp` for Copy & `mv` for Move:** Reassigning the universal two-letter alias `cp` exclusively to `gitmap copy` (eliminating legacy collision with `commit-push`), while ensuring `mv` exclusively represents `gitmap move`. All discrete 3-letter semantic commit verbs (`cpf`, `cpb`/`cbf`, `cpr`, `cpc`, `cfr`) remain strictly preserved.
2. **Multi-Node File Transport (`copy` / `cp` & `move` / `mv`):** Supporting 8 comprehensive syntax variants spanning local-to-single-node, local-to-all-nodes (fleet broadcast), remote-to-local (pull), remote-to-remote, recursive directory copy, and local clipboard disambiguation.
3. **OS-Aware Binary Deployment (`deploy`):** Enhancing `gitmap deploy` with binary magic byte inspection (PE `.exe` vs ELF), automatic OS incompatibility skipping (e.g. automatically skipping Windows binaries on Ubuntu nodes unless overridden with `--force` or `--force-copy`), offline node bypass with bounded TCP probes, and rich post-execution suggestions.
4. **Cross-Node Fleet Software Updates (`update`):** Centralized fleet update commands for Cursor IDE, Visual Studio Code, Antigravity Manager (AGM), and Antigravity Tools (AGY).

---

## 2. User Request (Verbatim)

```text
The `gitmap`, I wanted to have the deploy, copy, and move command, which I could not find in the help text. The deploy command would have copy binaries and also with the installation. For example, Windows, we should be able to copy binaries to just all the machines using copy command. And there is one more thing that I thought of we are doing a little bit wrong because usually the copy command is CP, and we are doing CP for commit pull or commit push. I think that needs to be corrected everywhere, along with the commands, examples, help text, UI help text and everything. So CP should be reserved for copy because people know CP.copy, and MV for move file from one machine to another. So create separate syntaxes. There should be four to eight different syntax for deploying a copy a file for all nodes, single node, a specific node from this current folder file. Also move file from one node to another or move file from here to all the nodes, things like that. So all types of commands, try to share that at the end that you have implemented. Make CP as copy command for the `gitmap` deployed and also with the deploy, we should be able to install the binary, regardless of the OS. If it is a EXE on the Windows and we are deploying, then we can also say skip Unix or skip Ubuntu. It would automatically skip that. So by default, EXE files can be skipped from Ubuntu machines when we do the deploy copy. And it can have a flag like force copy, then it would automatically do that, and also it would provide those suggestions after it runs, and it would mention what it is doing and why it is keeping the Ubuntu machine, things like that, so that it is efficient. And only the open machines or connected machines, those are running that can copy to remember that. And installation, if it is a EXE, so for example, the cursor update, it does not update automatically. So we should have a command to update cursor. In all the machines, we should have a command for the answer nodes, update cursor or VS code or `Antigravity Manager or Antigravity Tools`, things like that. So try to make sure that you work on it, you implement that, and also do a minor bump and release.
```

---

## 3. Modular Specification Map

| Document | Primary Scope | Primary Owner |
| :--- | :--- | :--- |
| [00-overview.md](00-overview.md) | High-level architecture, user request, invariant rules | Root Orchestrator |
| [01-alias-and-copy-move-spec.md](01-alias-and-copy-move-spec.md) | Alias remapping (`cp`->copy, `mv`->move), 8 copy syntaxes, move syntax | Subagent 1 |
| [02-deploy-and-os-compat-spec.md](02-deploy-and-os-compat-spec.md) | `deploy` binary command, magic byte check, OS skip logic, feedback docket | Subagent 2 |
| [03-fleet-software-update-spec.md](03-fleet-software-update-spec.md) | Cross-node updates for Cursor, VS Code, AGM, AGY tools | Subagent 3 |

---

## 4. Architectural Invariants & Non-Negotiable Standards

1. **Leaf Error Sovereignty:** Leaf functions must return `error` or `*apperror.AppError` directly upward. No calling exit handlers in leaf functions.
2. **Positive Polarity Only:** All booleans must use `is*`, `has*`, `can*`, `should*` prefixes (e.g. `isOnline`, `isCompatible`, `hasTarget`). Zero inverted checks (`!isSuccess` is banned; use `isFail` or `IsFail()`).
3. **100% Relative Git Paths:** Zero hardcoded absolute paths (`C:\...`, `/home/...`, `file:///...`) in documentation, code, or release notes.
4. **Zero-Build & Zero-Test During Maintenance:** Targeted file-level verification only. Full builds and test suites are not executed during routine turns.
5. **Semantic Commit Verbs Preservation:** `cpf`, `cpb`/`cbf`, `cpr`, `cpc`, `cfr` must remain completely intact.
