# 04 Verification Gates

## 1. Acceptance Criteria: Split-DB Cache
- **Speed**: `gitmap scan <group>` must be highly performant, eliminating the 140-sec latency issue.
- **Persistence**: Evaluate ignore patterns ONCE and cache the boolean result (true/false for ignore bit) in Split-DB (`repodb.db`).
- **Validation**: Ignore evaluation should NEVER re-run unless the `.gitmapignore` hash changes.
- **Architecture Validation**:
  - `gitmap.db` is used for local metadata.
  - `repodb.db` is used for cached repo data.
  - `pipelinedb.db` is used for pipeline executions.
  - `nodes.db` is used for SSH nodes history and status.

## 2. Acceptance Criteria: Ignoring Binaries
- The cache engine must automatically identify binary files and cache them as ignored.
- Binary evaluation must be done once and not repeatedly evaluated on subsequent reads.

## 3. Acceptance Criteria: See Commands Output
- `gitmap see all fast [group]` must display all fast paths correctly.
- `gitmap see all repo-cache [group]` must display repo cache output correctly.
- `gitmap see ssh-node all` must list all globally or locally registered SSH nodes.
- `gitmap see ssh-node <id>` must display detailed info (id, os, type, ip, identityfile, alive status, load/cpu/ram, latency).
- `gitmap see pending-tasks all <group>` must list all pending-tasks for a group in ai-memory.

## 4. Acceptance Criteria: Bug Fix (Screenshot)
- The issue where `gitmap scan <group>` evaluates ignore regexes (which are heavy) on EVERY file via `WalkDir` must be resolved.
- The GitMap cache table must store both ignored and non-ignored paths (true/false) and filter them on **read** rather than **write**.
