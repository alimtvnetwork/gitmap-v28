# Specification 359: GitMap Autonomous AI Identity & Provenance Feature

- **Spec Number:** 359
- **Domain:** CLI Architecture, Identity Governance & Heatmap Attribution
- **Status:** Active / Approved
- **Slug:** `gitmap-feature-implementation-and-update`

---

## 1. Executive Summary & Strategic Intent

This specification establishes first-class autonomous AI identity management and commit provenance inside GitMap. Modern multi-agent swarms (such as V6 and V7 orchestrators) execute autonomous engineering workflows that require clean separation between human developer credentials and automated agent commits.

Previously, autonomous agents relied on verbose external shell commands (`gitmap user change-user "gitmap automation bot - muse - runs autonomously" "info@rgitmap.com" --global`) or brittle environment manipulation to set identities. 

This feature natively embeds:
1. **Sovereign Login Configuration (`gitmap login set ai-auto [agent-name]` & `gitmap login set global ai-auto`):** Directly configures the autonomous bot identity (`user.name` and `user.email`), backs up human developer identity for non-destructive restoration, and updates Git profiles.
2. **Commit Provenance via `--ai-agent <name>`:** Extends semantic commit verbs (`cpf`, `cpb`, `cpc`, `cfr`, `cpr`, `pcp`) with an explicit `--ai-agent` flag, attributing the commit to the active agent in git history, Split-DB telemetry (`repodb/commit_history.db`), and cognitive heatmaps (`repodb/heatmap.db`).
3. **Repository-Wide Skill & Prompt Synchronization:** Realigns canonical skills (`gitmap`, `gitmap-automation-bot`, `execute-parent-task-with-n-steps-v6`, `execute-parent-task-with-n-steps-v7`) and execution prompts across both the `gitmap-v28` and `coding-guidelines-v24` repos, introducing completion checklists of tables, masked verification script requirements, and explicit zero-build/zero-test guardrails unless explicitly instructed.

---

## 2. Architectural Boundaries & Component Invariants

### 2.1 Identity Invariants
- **Default Autonomous Name:** `gitmap automation bot - runs autonomously` (for global `ai-auto` or when no agent name is passed).
- **Agent-Specific Autonomous Name:** `gitmap automation bot - <agent-name> - runs autonomously` (e.g. `muse`, `gemini`, `cursor-model`).
- **Autonomous Email Invariant:** `info@rgitmap.com` across all automated bot configurations.
- **Human Identity Protection Invariant:** Before applying any AI bot configuration, GitMap must verify if the active identity is a bot. If it is a human user, GitMap unconditionally backs up the original `user.name`, `user.email`, and profile alias to `original_user.json` in both binary and global data directories.

### 2.2 Commit Invariants
- **Sovereign Commit Verbs:** `cpr` (release), `cpf` (feature), `cpc` (chore), `cbf` (bug fix), `cfr` (refactor), `pcp` (pull-commit-push).
- **Hyphen-Separated Formatting:** `"<Module> - <summary>"`. Colons and square brackets `[...]` inside the argument remain strictly banned.
- **Attribution Propagation:** When `--ai-agent <name>` is provided, `meta.AIName`, `meta.AgentID`, and `meta.IsAIMode` are populated, and the author identity defaults to the agent's bot identity if `--author` was omitted.

### 2.3 Verification & Quality Gate Invariants
- **Zero-Build & Zero-Test Default:** If end-to-end testing or build verification is not explicitly commanded by the user, AI agents executing V6 or V7 prompts MUST NOT trigger full builds or test suites.
- **Masked Verification Script:** At completion of V6/V7 execution, agents must generate a masked verification script and output a checklist of completed tables.
