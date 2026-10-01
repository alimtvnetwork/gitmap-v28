---
status: PENDING
---

# Subtask: Cache Engine

## Objective
Implement the Split-DB cache engine to store and retrieve repository and ignore states efficiently.

## Requirements
- Implement the Split-DB architecture: `gitmap.db` (local metadata), `repodb.db` (cached repo data), `pipelinedb.db` (pipeline executions), and `nodes.db` (SSH nodes).
- The cache table must store both ignored and non-ignored paths (true/false for the ignore bit).
- Filtering of ignored paths must occur on **read** rather than on **write**.
- The cache engine must store the hash of the `.gitmapignore` file and use it to invalidate the cache.

## Next Steps
1. Review Split-DB implementation in `cli/store/` and `cli/repodb/`.
2. Update the cache schema to include an `ignored` boolean flag for paths.
3. Modify the write logic to insert all paths regardless of ignore status.
4. Modify the read logic to filter out paths where `ignored = true` when required.
