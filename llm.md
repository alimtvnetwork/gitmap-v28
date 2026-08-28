# Gitmap LLM Specification

Gitmap is a powerful CLI designed for autonomous agents, LLMs, and developers to navigate, analyze, and modify codebases efficiently. 

## Capabilities

1. **Repository Discovery & Scanning**:
   - gitmap scan: Automatically locate and index Git repositories.
   - gitmap rescan: Refresh indexed local repositories.

2. **Repository Cloning & Synchronization**:
   - gitmap clone, gitmap clone-next: Manage bulk cloning.

3. **Intelligent File Search & Indexing (Split DB)**:
   - Gitmap utilizes a Split SQLite Database architecture for rapid indexing.
   - gitmap search <query>: Execute exact searches quickly on-the-fly.
   - gitmap repo-search <query>: Perform analytical, cache-backed searches across the repo.
   - gitmap repo-search-json <query>: Retrieve search results as structured JSON (ideal for LLM tool consumption).
   - *Note*: Gitmap natively skips .git, 
ode_modules, and handles files >300KB using lazy regex to maintain high performance.

4. **Regex & Replacement**:
   - gitmap replace <query> <replacement>: Find and replace exact phrases.
   - gitmap replace-regex <regex> <replacement>: Regex replacement with automatic history tracking.
   - gitmap replace history: View replacement operations.

5. **Releases & Commits**:
   - gitmap release: Automate semantic version bumping and changelog generation.
   - gitmap commit-in: Perform robust, orchestrator-driven batch commits.

## Instructions for LLMs
- Always prefer JSON output commands (sj, epo-search-json) when parsing data programmatically.
- Avoid modifying the root SQLite DB manually; rely on the CLI commands.
- Before running heavy regex operations across a large repo, use gitmap search to verify matches.


## AI File Search Patterns
When searching codebases, LLMs can use native gitmap commands OR standard terminal tools.
Here are equivalent alternative command samples for LLM search operations:

- **Find a specific struct definition**:
  - `gitmap file-search . "type SearchResult struct"` 
  - `Get-ChildItem -Path gitmap -Recurse -File | Select-String "type SearchResult struct"` 

- **Search functions with Regex context**:
  - `gitmap file-search cmd/ "func dispatch[A-Z]" 0 10` 
  - `Get-ChildItem -Path gitmap/cmd -Filter *.go | Select-String "func dispatch[A-Z]"` 

- **Find specific function contexts**:
  - `gitmap file-search cmd/root.go "func finishCommandAudit" 0 10` 
  - `cat gitmap/cmd/root.go | Select-String "func finishCommandAudit" -Context 0,10` 



## STRICT Commit Guidelines for AI Agents
AI Agents MUST NEVER use raw git commands for committing or pushing (e.g., git commit, git push).
Instead, you MUST use the native Gitmap commands for grouping and committing:
- gitmap feature <name>: Start or group a feature.
- gitmap bar: Use the bar command to handle specific groups.
- gitmap release: Use this to finalize and release grouped commits.
- gitmap commit-in: Commit changes inside grouped nodes.

Always group your commits. DO NOT commit every single file as a separate big node. Use Gitmap's orchestration.
