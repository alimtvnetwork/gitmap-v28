# Subtask 05: Antigravity E2E Verification, Skills/Prompts & Minor Release

- **Subtask ID**: `subtask-05-skills-prompts-and-minor-bump`
- **Parent Task**: `command-and-database-integration-instructions`
- **Specification Reference**:
  - `02-spec/21-app/command-and-database-integration-instructions/05-antigravity-verification-and-testing.md`
- **Owned Files**:
  - `.agents/skills/gitmap/SKILL.md`
  - `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`
  - `02-spec/02-coding-guidelines/`
  - `VERSION`
  - `package.json`
- **Assigned Worker**: Worker 02
- **Status**: PENDING

---

## 1. Objectives

1. Execute End-to-End Verification Protocol against live Antigravity databases:
   - Target `%USERPROFILE%\.gemini\antigravity\conversation_summaries.db` (or `~/.gemini/antigravity/...`).
   - Verify table enumeration (`ls`), column introspection (`table conversation_summaries`), and query execution (`SELECT conversation_id, title FROM conversation_summaries`).
   - Ensure read-only guard (`PRAGMA query_only = ON;`) prevents modifications to live IDE databases.
2. Update GitMap Skill Documentation (`.agents/skills/gitmap/SKILL.md`):
   - Document the complete `db-*` command ecosystem (`db-read`, `db-add`, `db-edit`, `db-find`, `db-view`, `db-listing`) with copy-paste examples.
   - Expand `heatmap` telemetry documentation with practical command patterns (`gitmap rc heatmap status`, `tree`, `live`, `sync`).
3. Update Coding Guidelines & Execute V6 Prompt:
   - Amend `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md` and related coding guidelines to state that native `gitmap db-*` commands are the authorized, preferred method for database introspection and operations, replacing external Python one-liners.
4. Execute Version Bump & Release Ceremony:
   - Bump minor version to `v6.580.0` across repository manifests (`VERSION`, `package.json`, root `README.md`).
   - Document all new commands, bug fixes (deadlock resolution), and capabilities in release notes.
   - Verify all CI/CD gates pass before tagging and pushing.

---

## 2. Documentation Additions (`.agents/skills/gitmap/SKILL.md`)

```markdown
### SQLite Database Operations (`gitmap db-*`)

Replace slow, fragile Python one-liners with native, CGO-free Go SQLite commands:

- `gitmap db-read <alias|path> "SELECT ..."`: Run safe, read-only SQL queries.
- `gitmap db-read <alias|path> ls`: List all user tables, views, and row counts.
- `gitmap db-read <alias|path> table <tbl> --columns --json`: Inspect table schema and datatypes.
- `gitmap db-add <alias> <table> csv "v1", "v2"`: Insert rows with transaction safety.
- `gitmap db-edit <alias> -table "tbl" -id "id" csv ...`: Update record by primary key.
- `gitmap db-find <alias> -table "tbl" --find "term*"`: Wildcard or regex search across columns.
- `gitmap db-view <alias> -t "t1 as a, t2 as b" --viewname "v1"`: Auto-join views.
- `gitmap db-listing [--json]`: View all registered database aliases, file paths, and sizes.
```

---

## 3. Verification Criteria

1. Live Antigravity database queries run cleanly without throwing SQL errors or locking files.
2. All documentation references in `.agents/skills/gitmap/SKILL.md` and coding guidelines are relative and valid.
3. `gitmap --version` reports `v6.580.0`.
4. Continuous integration suite (`go test ./...`) passes cleanly with zero failures.
