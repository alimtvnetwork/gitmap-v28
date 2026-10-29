# 348-03: Antigravity Skill Integration, Onboarding Guide, Verification Matrix, and Release Procedure

## Status: APPROVED
- **Document Identifier:** `02-spec/21-app/348-recommit-and-skill-creation-instructions/03-migration-and-verification-spec.md`
- **Parent Plan:** `.ai-memory/plans/pending/348-recommit-and-skill-creation-instructions.md`
- **Related Specs:**
  - `02-spec/21-app/348-recommit-and-skill-creation-instructions/01-architecture-spec.md`
  - `02-spec/21-app/348-recommit-and-skill-creation-instructions/02-component-spec.md`
  - `02-spec/21-app/348-recommit-and-skill-creation-instructions/readme.md`
- **Execution Subtask:** `.ai-memory/plans/subtasks/348-recommit-and-skill-creation-instructions/03-antigravity-integration-and-release.md`

---

## 1. Executive Summary & Scope

This specification establishes the definitive Antigravity skill onboarding guide, comprehensive verification matrix, and atomic commit release procedure for Plan 348 (`recommit-and-skill-creation-instructions`).

This specification covers:
1. **Antigravity Skill Onboarding Protocols:** Standard operating procedures for integrating agent skills via:
   - Workspace repository directory (`.agents/skills/<name>/skill.md`).
   - User global directory (`~/.gemini/antigravity/skills/<name>/SKILL.md`).
   - Programmatic Python SDK configuration (`LocalAgentConfig(skills_paths=[...])`).
2. **Progressive Disclosure & Lifecycle Management:** How skills are discovered, indexed by metadata, and progressively loaded into the context window.
3. **Verification Matrix & Quality Gates:** Mechanical checks ensuring strict conformance with repository coding standards (relative paths, positive booleans, leaf error sovereignty, function line limits).
4. **GitMap Atomic Release Ceremony:** Step-by-step semantic release process utilizing dedicated GitMap commit verbs (`gitmap cpf`, `gitmap cpr`).

---

## 2. Antigravity Skill Discovery & Loading Architecture

### 2.1 Skill Structure & Specification Standards

An Antigravity skill is a standardized, modular package containing operational knowledge, procedures, and runbooks following the Agent Skills specification. Every skill directory contains:

```text
.agents/skills/<skill-name>/
├── skill.md (or SKILL.md)    # Required: Main instruction document with YAML frontmatter
├── scripts/                  # Optional: Deterministic automation or validation scripts
├── examples/                 # Optional: Reference implementations and payloads
├── resources/                # Optional: Reusable templates and assets
└── references/               # Optional: In-depth reference manuals loaded on demand
```

### 2.2 YAML Frontmatter Contract

The primary instruction file (`skill.md` or `SKILL.md`) must open with a valid YAML frontmatter block:

```markdown
---
name: skill-identifier
description: >-
  Actionable third-person description stating exactly what the skill accomplishes
  and the exact triggers for when the model should activate it.
---
```

Frontmatter validation rules:
- `name`: Must be lowercase, alphanumeric with hyphens (kebab-case). Matches the directory name.
- `description`: Must be concise and written in the third person. It serves as the trigger index in the primary agent context window.

### 2.3 Progressive Disclosure Mechanism

To maximize available context and prevent token saturation:
1. **Index Phase (Lightweight):** On agent startup, Antigravity parses only the YAML frontmatter (`name` and `description`) across all discovered skill folders.
2. **Context Injection Phase (On-Demand):** When a user prompt matches the activation criteria or an orchestrating prompt explicitly requests a skill, Antigravity loads the full content of `skill.md` into the agent's active context.
3. **Deep Reference Resolution:** Bulky guidelines or extended documentation in `references/` remain unloaded until the agent explicitly traverses them via file reading tools.

### 2.4 Precedence Hierarchy

When skills share identical names across multiple discovery locations, Antigravity resolves them using the following priority order (highest to lowest):

| Priority Tier | Discovery Location | Scope | Team Sharing |
| :--- | :--- | :--- | :--- |
| **Tier 1 (Highest)** | Workspace Directory (`.agents/skills/<name>/`) | Project Repository | Yes (committed to VCS) |
| **Tier 2** | Workspace Config Manifest (`.agents/skills.json`) | Project Repository | Yes (committed to VCS) |
| **Tier 3** | User Global Directory (`~/.gemini/antigravity/skills/<name>/`) | Machine User | No (local to developer) |
| **Tier 4** | Global Customization Config (`~/.gemini/config/skills/`) | Machine User | No (local to developer) |
| **Tier 5 (Lowest)** | Built-in Application Skills | Runtime Environment | Bundled with binary |

---

## 3. Complete Skill Onboarding Guide

### 3.1 Method 1: Workspace Repository Registration (Recommended for Teams)

Workspace registration binds the skill directly to the git repository, ensuring that every engineer and autonomous agent operating on the codebase shares identical procedural knowledge.

#### Step-by-Step Procedure:
1. Navigate to the repository root.
2. Create the dedicated skill directory inside `.agents/skills/`:
   ```bash
   mkdir -p .agents/skills/gitmap-time-restriction
   mkdir -p .agents/skills/gitmap-automation-bot
   ```
3. Author the instruction file as `skill.md` (or `SKILL.md`) with valid YAML frontmatter:
   ```markdown
   ---
   name: gitmap-time-restriction
   description: >-
     Configures time restrictions, operational windows, and recommit enforcement
     using GitMap CLI commands. Use when the user requests time-gated execution or recommit.
   ---

   # GitMap Time Restriction Skill
   ...
   ```
4. Verify the folder structure using relative Git paths:
   - `.agents/skills/gitmap-time-restriction/skill.md`
   - `.agents/skills/gitmap-automation-bot/skill.md`
5. Commit the files to version control using GitMap semantic commit verbs:
   ```bash
   gitmap cpf "SKILLS - add gitmap time restriction and automation bot skills"
   ```

### 3.2 Method 2: User Global Directory Registration (Machine-Wide)

Global registration provides skills across all repositories and workspaces on the host machine without modifying individual project trees.

#### Step-by-Step Procedure:
1. Locate the user Antigravity configuration directory in the home folder:
   - Linux/macOS: `~/.gemini/antigravity/skills/`
   - Windows: `~/.gemini/antigravity/skills/` (resolves to `%USERPROFILE%\.gemini\antigravity\skills\`)
2. Create the target skill folder:
   ```bash
   mkdir -p ~/.gemini/antigravity/skills/<skill-name>
   ```
3. Create `SKILL.md` inside that directory with the required frontmatter:
   ```markdown
   ---
   name: <skill-name>
   description: >-
     User global capability for cross-project workflows.
   ---

   # User Global Skill
   ...
   ```
4. On the next session initialization, Antigravity automatically indexes the new skill and lists it in available skills.

### 3.3 Method 3: Programmatic API Registration via Python SDK

Developers building custom autonomous agent harnesses or running batch evaluations can dynamically register skills using the Antigravity Python SDK.

#### Configuration Mechanics:
The `LocalAgentConfig` class accepts `skills_paths` as a parameter. Paths can point to:
- A directory containing multiple skill subfolders (batch discovery).
- A specific skill directory containing `SKILL.md` (isolated inclusion).

#### Implementation Example (Adhering to Coding Standards):

```python
"""Antigravity agent initialization with programmatic skill paths."""

from typing import List
from google.antigravity import Agent, LocalAgentConfig


class AgentSkillManager:
    """Manages skill discovery paths with positive boolean validation."""

    def __init__(self, skills_paths: List[str]) -> None:
        self.skills_paths = skills_paths
        self.has_custom_paths = len(skills_paths) > 0
        self.is_initialized = False

    def get_agent_config(self) -> LocalAgentConfig:
        """Constructs LocalAgentConfig with validated skill paths."""
        if self.has_custom_paths:
            return LocalAgentConfig(skills_paths=self.skills_paths)

        return LocalAgentConfig()

    async def run_agent_session(self, prompt: str) -> str:
        """Executes an agent turn with configured skills."""
        config = self.get_agent_config()
        self.is_initialized = True

        async with Agent(config) as agent:
            response = await agent.chat(prompt)
            return await response.text()
```

#### Go Leaf Error Sovereignty & Helper Pattern:

When bridging configuration from Go CLI orchestrators into agent runners, enforce leaf error sovereignty and <= 15 line functions:

```go
// SkillConfig defines agent skill registration parameters.
type SkillConfig struct {
	SkillPaths []string
	IsEnabled  bool
	HasCustom  bool
}

// BuildSkillConfig constructs validated skill configuration.
func BuildSkillConfig(paths []string) (SkillConfig, error) {
	if len(paths) == 0 {
		return SkillConfig{}, apperror.New("skill paths cannot be empty")
	}

	return SkillConfig{
		SkillPaths: paths,
		IsEnabled:  true,
		HasCustom:  len(paths) > 1,
	}, nil
}
```

---

## 4. Verification Matrix & Quality Gates

Before any skill, prompt, or specification is released, it must pass the following mechanical quality gates.

### 4.1 Mechanical Verification Matrix

| Verification ID | Category | Target Deliverable | Verification Command / Method | Acceptance Criteria | Conformance Status |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **VR-01** | Relative Paths | All authored specs & skills | `gitmap aum search ":\\\\" 02-spec/21-app/348*` | 0 occurrences of host absolute paths. | PASS |
| **VR-02** | Positive Booleans | Code & spec attributes | Scan identifiers for banned negative prefixes | All booleans use `is*` or `has*` prefix. Ban on `can`, `should`, `was`, `will`. | PASS |
| **VR-03** | Function Sizing | Go & Python helpers | AST / line counter scan | Every function is <= 15 lines (<= 8 lines preferred). | PASS |
| **VR-04** | Error Sovereignty | Go leaf functions | Return signature validation | Leaf functions return `(T, error)` directly. No `cliexit.HandleError` in leaves. | PASS |
| **VR-05** | Frontmatter Schema | `.agents/skills/*/skill.md` | YAML parser validation | Valid `name` matching dir, non-empty `description` in 3rd person. | PASS |
| **VR-06** | Prompt Formatting | `01-prompts/22-letterly/12*` | GitMap streaming inspection | Conforms to Letterly template with explicit on-call skill invocations. | PASS |
| **VR-07** | Prompt Formatting | `01-prompts/14-execute/09*` | GitMap streaming inspection | Conforms to V7 execution prompt standard; does not modify legacy V6. | PASS |
| **VR-08** | Working Tree | Git status | `gitmap ps` or `git status` | Clean working tree; zero untracked artifacts. | PASS |

### 4.2 Quality Gate Checklist
- [x] **Gate 1 (Path Hygiene):** Strict relative Git paths only (`02-spec/...`, `.ai-memory/...`, `.agents/skills/...`).
- [x] **Gate 2 (Boolean Nomenclature):** Positive booleans only (`isVerified`, `hasCustomPaths`, `isEnabled`).
- [x] **Gate 3 (Error Architecture):** Leaf functions return errors directly; orchestrators handle exit codes.
- [x] **Gate 4 (Complexity Control):** Zero nested if statements; all guards flattened with early returns.
- [x] **Gate 5 (Preservation):** V6 execution prompt (`01-prompts/14-execute/08-*`) remains intact without modifications.

---

## 5. GitMap Atomic Release Ceremony

The release ceremony follows the strict GitMap semantic commit discipline.

### 5.1 Pre-Flight Synchronization
Always execute `gitmap pull` before finalizing code changes:
```bash
gitmap pull
```
*Fallback:* If `gitmap pull` encounters remote divergence, run `git pull` and resolve merges before proceeding.

### 5.2 Commit Discipline Rules
- Total ban on raw `git commit`, `git add -A`, or arbitrary commit strings.
- Total ban on square brackets `[...]` or colons `:` inside commit messages.
- Total ban on nested conventional prefixes like `chore(release):` or `feat:`.
- Semantic commit verbs:
  - `gitmap cpf`: Features and new capabilities.
  - `gitmap cfr`: Refactoring existing modules.
  - `gitmap cbf`: Bug fixes.
  - `gitmap cpc`: Routine chores and maintenance.
  - `gitmap cpr`: Version bump releases.

### 5.3 Staged Commit Sequence for Plan 348

#### Commit 1: Architecture & Component Specifications
```bash
gitmap cpf "SPEC - author 348 recommit and skill creation specifications"
```

#### Commit 2: Dedicated Antigravity Skills
```bash
gitmap cpf "SKILLS - add gitmap time restriction and automation bot skills"
```

#### Commit 3: Letterly (12) and Execution Prompt (09 V7)
```bash
gitmap cpf "PROMPTS - add letterly 12 and v7 execution prompt with time restriction"
```

#### Commit 4: Subtask Plan & Antigravity Guide Integration
```bash
gitmap cpc "CHORE - document antigravity skill onboarding and verification subtask"
```

#### Commit 5: Semantic Release Ceremony (When Tagging)
```bash
gitmap cpr "v6.614.0 add gitmap time restriction skills and v7 execution prompt"
```
*(Format strictly as `"v<version> <summary>"`; total ban on redundant `VERSION -` or `RELEASE -` prefixes).*

---

## 6. Maintenance & Fleet Synchronization

Once committed, the newly authored skills and prompts can be synchronized across all connected workspaces using GitMap's synchronization toolchain:
```bash
gitmap aum sync --skills
```
This ensures uniform agent capabilities across all active project repositories.
