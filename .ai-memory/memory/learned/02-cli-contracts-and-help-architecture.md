# 02 — CLI Contracts and Help Architecture

- **Subsystem:** CLI Framework & Universal Help Interceptor
- **Status:** Authoritative Reference

## 1. Command Routing & Interception
- **Cobra Foundation:** Root command dispatch configured in `cli/cmd/root.go` routes execution through subcommand packages (`cmdcommit`, `cmdnodes`, `cmdpas`, etc.).
- **Universal Help Interceptor (`tryInterceptCommandHelp`):**
  - Intercepts any invocation ending in `help`, `-h`, or `--help`.
  - Dispatches to rich Markdown documentation in `cli/helptext/` when available.
  - Dynamically synthesizes a modern framed box help screen (`RenderDynamicCommandHelp`) when no static document exists.
  - Guarantees 100% interception with zero unintended side effects or live command execution.

## 2. Levenshtein Typo Suggestion Engine
- In-memory prefix trie evaluates edit distances across 350+ commands and aliases in <1ms.
- Employs deterministic suggestions: `Did you mean: gitmap <command>?`.

## 3. Terminal Presentation & Boxes
- Standardized double-line framed box banners (`termhelp`).
- Two-column aligned flags and descriptions with smart word wrapping.
