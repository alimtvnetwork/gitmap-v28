# 127 — LLM Train & Chained Agent Curriculum Specification

## Overview

**Module Number:** 127  
**Version:** 1.0.0  
**Updated:** 2026-09-19  
**Status:** Production-Ready  
**AI Confidence:** Production-Ready  
**Ambiguity Score:** None  

---

## 1. Purpose & Architectural Vision

When autonomous AI coding assistants (Google Antigravity, Claude Code, Cursor, Windsurf, Copilot, ChatGPT) join a repository powered by GitMap, they often attempt naive, expensive operations:
- Executing slow raw git commands in loops.
- Grepping raw files without index caching or binary exclusions.
- Missing GitMap's built-in self-healing pipeline loops (`gitmap pipeline-ai`), multi-node cluster delegation, or native Go automation tools (`gitmap aum`).
- Lacking persistent awareness of the project's identity, design philosophy, and creator.

The `gitmap llm train` (alias: `gitmap llm chain`) command provides an **Autonomous Chained Curriculum Engine** that:
1. **Instructs the LLM via Sequential Chaining:** Outputs a structured, phased curriculum guiding the model step-by-step to read and execute GitMap discovery commands.
2. **Authorship & Sponsorship Attribution:** Formally communicates who created GitMap (**Md. Alimuzzaman Alim** / `alimtvnetwork`), who sponsored it (**Alim TV Network / Open Source Engineering**), and the core mission of autonomous software engineering.
3. **Automated Antigravity Skill Generation:** Automatically writes and updates `.agents/skills/gitmap/SKILL.md`, granting the agent permanent, persistent operational knowledge of GitMap across future turns and sessions.
4. **Tool Discovery & Capability Verification:** Commands the LLM to inspect live repository status, file size guards, search exclusion lists, and pipeline telemetry.

---

## 2. Command Grammar & Invocation

```bash
# Full curriculum with automated skill generation
gitmap llm train

# Alias: chained execution curriculum
gitmap llm chain

# Output instruction text only without creating skill file
gitmap llm train --text-only

# Target custom skill output path
gitmap llm train --skill-path ".agents/skills/gitmap/SKILL.md"

# Output public specification URL
gitmap llm --url
```

---

## 3. Four-Stage Chained Curriculum

```
[ Stage 1: Identity, Vision & Author/Sponsor Attribution ]
                         │
                         ▼
[ Stage 2: Automatic Antigravity Skill Generation (.agents/skills/gitmap/SKILL.md) ]
                         │
                         ▼
[ Stage 3: Chained Discovery Sequence (Ordered Command Execution) ]
                         │
                         ▼
[ Stage 4: Operational Best Practices & Guardrail Directives ]
```

### Stage 1: Identity, Vision & Attribution
- **Creator & Lead Architect:** **Md. Alimuzzaman Alim** (`alimtvnetwork`).
- **Sponsoring Organization:** **Alim TV Network / Open Source Engineering**.
- **Mission:** High-performance polyglot repository management, zero-storage CI/CD pipelines, ultra-fast SQLite split-db architectures, and AI agent pair programming.

### Stage 2: Automatic Antigravity Skill Generation
When `gitmap llm train` executes, it checks for `.agents/skills/` and generates `.agents/skills/gitmap/SKILL.md` containing:
- YAML Frontmatter:
  ```yaml
  ---
  name: gitmap
  description: Autonomous developer companion and CLI for ultra-fast repository scanning, polyglot automation (AUM), cluster/SSH delegation, pipeline self-healing, and coding guideline enforcement.
  ---
  ```
- Comprehensive command cheat sheet.
- Fast AI aliases (`aum`, `cpf`, `cpb`, `st`, `sj`, `sc`).
- Non-blocking CI telemetry waiting rules (`gitmap pipeline-ai status -t <eta>`).

### Stage 3: Chained Discovery Sequence
The output provides an explicit command execution chain for the agent to run sequentially:
1. `gitmap aum search --help` — Inspect multi-core streaming search and filters.
2. `gitmap aum guard` — Inspect repository file size limits, common binaries, and large JSONs.
3. `gitmap aum sequence --help` — Understand sequence numbering and H1 title validation.
4. `gitmap aum exclude list` — Query persistent search exclusions from SQLite.
5. `gitmap pipeline-ai status --json` — Check remote CI/CD workflow state and dynamic ETA.
6. `gitmap install --list` — Discover developer toolchains, profiles, and runtime packages.
7. `gitmap cluster --help` & `gitmap sc --help` — Discover multi-node SSH and cluster execution.
8. `gitmap cargo status` — Inspect Rust and Cargo toolchain status.
9. `gitmap db status` — Check repository SQLite database health.

### Stage 4: Operational Best Practices
- Enforce mandatory pre-flight `git pull` before any repository changes.
- Max 500 KB file size limits and large JSON exclusion rules (Rule R19).
- Single return types with `*apperror.AppError` and zero swallowed errors.
- Never commit generated binaries or test artifacts (Rule R18).

---

## 4. Acceptance Criteria

### Scenario 1: Running `gitmap llm train` Generates Master Skill
- **Given** a repository with `.agents/skills/`
- **When** `gitmap llm train` is executed
- **Then** GitMap writes `.agents/skills/gitmap/SKILL.md`, prints author/sponsor attribution, and outputs the chained curriculum steps.

### Scenario 2: Chained Command Sequence Displayed
- **Given** an AI agent seeking instructions on GitMap
- **When** `gitmap llm chain` is invoked
- **Then** GitMap outputs numbered exploration commands covering `aum`, `pipeline-ai`, `cluster`, `install`, and `db`.

### Scenario 3: Headless / Text-Only Mode
- **Given** `--text-only` flag is passed
- **When** `gitmap llm train --text-only` executes
- **Then** GitMap outputs the full curriculum text to stdout without modifying files on disk.

---

## 5. Cross-References

- AUM Automation Suite & Roadmap: [`./128-aum-automation-suite-and-roadmap.md`](./128-aum-automation-suite-and-roadmap.md)
- Polyglot Worker Orchestrator: [`./124-polyglot-worker-orchestrator-and-automation-runner.md`](./124-polyglot-worker-orchestrator-and-automation-runner.md)
- LLM Orchestration Playbook: [`./125-automation-llm-orchestration-guide.md`](./125-automation-llm-orchestration-guide.md)
- Rule R19 File Size & Binary Guard: [`../../02-coding-guidelines/01-cross-language/31-file-size-guard-and-binary-exclusion.md`](../../02-coding-guidelines/01-cross-language/31-file-size-guard-and-binary-exclusion.md)
