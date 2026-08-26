# Plan: Coding Guidelines Multi-Repo, Version Status & Status Dirty Filter

## Context
Full implementation of:
1. Coding Guidelines version & status inspection from `version.json`.
2. Multi-repo parent directory auto-discovery for CG install/update, pull, and push.
3. Flexible target specifiers (paths, aliases, IDs, URLs) for CG.
4. Selective update enforcement (only updating repos with existing CG).
5. `gitmap status --dirty` filter.

Inputs:
- .lovable/spec/commands/03-cg-multirepo-and-status-dirty.md
