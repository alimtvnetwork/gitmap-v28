# Subtask 02: Muse Help Text 5 Semantic Clusters & ~4,300 LOC Code Reduction

> **Subtask ID:** Subtask-02  
> **Parent Plan:** `.ai-memory/plans/241-cicd-enhancements-help-text-grouping-and-suggestion-engine-design.md`  
> **Target Path:** `.ai-memory/plans/subtasks/241-cicd-enhancements-help-text-grouping-and-suggestion-engine-design/02-muse-help-text-clustering-and-code-reduction.md`  
> **Status:** PENDING  

---

## 1. Objectives & Scope
1. Consolidate GitMap's 24 disparate command groups into 5 structured semantic clusters based on Muse's architectural guidance.
2. Design a centralized, data-driven Markdown help parser `termhelp.FromMarkdown(mdBytes []byte) (*HelpMenu, error)`.
3. Eliminate 36+ redundant `*_help_menu.go` and `*_help_sections.go` files across `cli/`.
4. Quantify and verify codebase reduction of ~4,300 lines of code.

---

## 2. Detailed Technical Plan

### 2.1 The 5 Semantic Clusters Mapping
- **Cluster 1: Core & Repository Operations**
  * Commands: `scan`, `list`, `clone`, `pull`, `status`, `reconcile`, `stash`, `wip`, `group`, `cd`.
  * Rationale: Foundational day-to-day Git and repository state management.
- **Cluster 2: Release & Commit Engineering**
  * Commands: `release`, `changelog`, `list-versions`, `commit`, `cpf`, `cpb`, `cpr`, `cpar`.
  * Rationale: Version bumping, changelog authoring, and automated commit/push operations.
- **Cluster 3: Fleet & Remote SSH Cluster**
  * Commands: `ssh`, `nodes`, `cluster`, `sc`, `deploy`, `vhost`, `nginx`.
  * Rationale: Multi-node distributed operations, server delegation, and remote execution.
- **Cluster 4: AI & Automation Intelligence**
  * Commands: `agy`, `agm`, `aum`, `ai`, `macro`, `pipeline` (`pl`, `pe`), `rerun`, `sug`.
  * Rationale: Agent workflows, Model Context Protocol (MCP) integrations, CI self-healing, and terminal macros.
- **Cluster 5: System, OS & Developer Tooling**
  * Commands: `os`, `apps`, `install`, `uninstall`, `storage`, `clean`, `vscode`, `sync`.
  * Rationale: Host machine environment management, developer IDE synchronization, and system caches.

### 2.2 Data-Driven Markdown Help Engine
- In `cli/termhelp/markdown.go`:
  - Implement AST parser reading Markdown tokens:
    * `# Heading` -> Menu title and header banner.
    * `## Usage` -> Command syntax.
    * `## Flags` / list items -> Flag descriptions.
    * ````code blocks```` -> Example usage commands.
  - Generates Catppuccin-bordered output dynamically.
- Deprecate hardcoded Go struct initializations across `cli/cmd*/*_help_menu.go`.

---

## 3. Verification Criteria
- [ ] All 5 clusters correctly categorize 100% of registered GitMap commands.
- [ ] `termhelp.FromMarkdown` successfully converts standard Markdown into Catppuccin terminal boxes.
- [ ] File removal checklist accounts for ~3,560 LOC of Go boilerplate.
- [ ] Total codebase net line reduction target of ~4,300 LOC verified.
