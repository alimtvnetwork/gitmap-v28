# 05 — Antigravity (AGY) Agent System & IDE Ecosystem

- **Domain:** AI Agents, Antigravity SDK, Multi-IDE Parity & Workstation Governance
- **Authoritative Specification:** [05-antigravity-and-ide](../../02-spec/21-app/05-antigravity-and-ide/01-architecture-spec.md)
- **Status:** Active & Ratified

---

## 1. Antigravity (AGY) Agent Architecture

GitMap deeply integrates with the Google Antigravity (AGY) SDK to orchestrate autonomous development:

- **Agent Types:**
  - `Lead Orchestrator`: Handles planning, file scoping, subtask dispatch, git transactions, and release tagging.
  - `Worker Subagents`: Pure execution agents with write tools constrained to designated file boundaries; forbidden from executing git commands.
  - `Research Subagents`: Read-only discovery agents inspecting codebases, specs, and external references.
- **Decision Logs & Transcripts:** Subagent actions and tool trajectories are logged in `.system_generated/logs/transcript.jsonl` for full auditability.
- **Prompt Freeze Remediation:** Prompts use structured JSON response contracts and strict boundaries to prevent agent stalling or endless conversational loops.

---

## 2. Multi-IDE Synchronization & Workspace Parity

GitMap synchronizes developer environments across modern IDE ecosystems:

- **Supported IDEs:** Visual Studio Code, Cursor IDE, and Google Antigravity Desktop IDE.
- **Settings & Extension Sync:**
  - Centralized manifests in `repo-secrets/09-antigravity-backup/vault/` track installed extensions, themes, and keybindings.
  - GitMap CLI synchronization commands (`gitmap agy sync-ide`) mirror settings between host workstations and fleet nodes.
- **Skills & Plugins Management:** Antigravity plugins and custom `.agents/skills/` are synced from canonical repositories to local IDE workspaces.

---

## 3. Autonomous CI/CD Pipeline Fix Injection

When CI/CD pipelines fail, GitMap bridges error logs directly to the AI agent:

- **Fix Extraction:** GitMap extracts the failing test frame, bounded stack traces (5 leading + 20 trailing lines), and relevant file paths.
- **Prompt Synthesis:** Formats an autonomous remediation prompt embedding the exact error logs, coding guideline rules, and target source files.
- **Queue Injection (`gitmap agy inject`):** Enqueues the generated prompt directly into the Antigravity task queue, allowing the agent to begin immediate surgical remediation.
- **Deduplication Flag (`--force`):** Prevents duplicate prompt injections for identical pipeline failures.

---

## 4. Ubuntu Workstation Governance & Customization

- **Clean Developer Profiles:** Configures unified shell environments (`zsh`, `oh-my-zsh`), terminal prompts, and PATH exports across Ubuntu fleet nodes.
- **Window Management & Dock Alignment:** Governs GNOME / desktop dock positions, custom launchers, and default window geometry for multi-monitor setups.
- **Automated Dependency Scaffolding:** Ensures Node/pnpm, Go toolchains, Python virtual environments, and Docker runtimes are aligned across all developer workstations.
