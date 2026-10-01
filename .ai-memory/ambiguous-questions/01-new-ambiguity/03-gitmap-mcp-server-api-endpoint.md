# Ambiguity 03: GitMap as an MCP Server API Endpoint for Autonomous AI Agents

- **Slug:** `gitmap-mcp-server-api-endpoint`
- **Status:** `open`
- **Raised At:** 2026-10-01
- **Blocking:** No (Deferred to future MCP phase per user instructions in `02-spec/21-app/pas-fix-parts/04-part.md`)
- **Spec Reference:** [02-spec/21-app/197-gitmap-pas-fix-and-repo-cache-commands.md](../../../02-spec/21-app/197-gitmap-pas-fix-and-repo-cache-commands.md)

---

## 1. Question / Unresolved Decision

The user mandated exposing GitMap operations via an API endpoint and Model Context Protocol (MCP) server so that AI agents (such as Google Antigravity, Claude Code, and external tooling) can query GitMap natively:
> *"API endpoint I wanted to discuss because we want to Use that Gitmap as an MCP server for the AI in the future. Put a question mark and ambiguity question in your ambiguity folder so we can discuss this later. Put it as a pending task. We will discuss it later on, how it's going to work. Okay."*

Specific open questions:
1. **Transport Mechanism:** Should the MCP server operate over standard stdio JSON-RPC (standard for local agents like Claude Code / Antigravity), or run an embedded HTTP/SSE daemon (`gitmap mcp serve --port 8765`), or both?
2. **Tool Scope:** Which GitMap CLI verbs should be exposed as first-class MCP tools?
   - Candidate tools: `gitmap_find`, `gitmap_cat`, `gitmap_search`, `gitmap_status`, `gitmap_pull_all`, `gitmap_cache_search`, `gitmap_nodes_status`.
3. **Authentication & Multi-Node Cluster Access:** Should the MCP server expose fleet node operations across SSH directly through MCP tool invocations, and how should host authentication keys be scoped?

---

## 2. Options Considered

- **Option A (Stdio-first MCP):** Implement `gitmap mcp` command implementing the MCP stdio protocol using JSON-RPC 2.0. Lightweight, zero network port binding, instant integration into agent config files.
- **Option B (Dual Stdio + SSE HTTP Server):** Support both `gitmap mcp` (stdio) and `gitmap mcp --http :8765` for remote web agents and cluster UI integration.
- **Option C (Embedded within Main-Worker Service):** Delegate MCP server endpoints to the central service daemon.

---

## 3. Impact if Guessed Wrong

Prematurely implementing a custom HTTP server without aligned MCP JSON schemas risks API churn and potential security exposure of local cluster credentials. Storing this ambiguity blocks premature guessing while capturing the user's intent losslessly.
