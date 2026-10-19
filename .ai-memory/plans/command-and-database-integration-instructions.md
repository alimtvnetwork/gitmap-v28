# Execution Plan: Command and Database Integration Instructions

- **Task Slug**: `command-and-database-integration-instructions`
- **Parent Task ID**: `task-20261020000428-command-and-database-integration-instructions`
- **Specifications Suite**:
  - `02-spec/21-app/command-and-database-integration-instructions/01-architecture-overview.md`
  - `02-spec/21-app/command-and-database-integration-instructions/02-cli-command-contracts.md`
  - `02-spec/21-app/command-and-database-integration-instructions/03-alias-registry-and-storage.md`
  - `02-spec/21-app/command-and-database-integration-instructions/04-view-and-join-engine.md`
  - `02-spec/21-app/command-and-database-integration-instructions/05-antigravity-verification-and-testing.md`
- **Status**: PENDING
- **Budget**: 400 Steps (Phase 1: 1..100, Phase 2: 101..250, Phase 3: 251..400)
- **Capacity**: A = 3 agents, H = 2 hands, Concurrency = 6

---

## 1. Verbatim User Request

```text
Try to have similar functionalites in gitmap like

gitmap db-read "path" "sql"
gitmap db-read "path" ls # should show tables list from sqlite db
gitmap db-read "path" -as "alias name unique saves into it's sqlite db as indeitifier"

gitmap db-read aliasName "sql ..." --limit/l 20/100
gitmap db-read aliasName table "table1" "table 2" --columns/c --json # list out columns in json or list view if no jsonf flag
gitmap db-read aliasName table "table1" "table 2" --json # show table view with datatypes
gitmap db-add aliasName "table" # prompt for each column
gitmap db-add aliasName "table" csv "col1", "col2", "col3"  # prompt for each column, etc
gitmap db-edit aliasName -table"table" -id "id-first" csv "col1", "col2", "col3"  # prompt for each column, etc
gitmap db-find aliasName -table"table, table" --find/f "text" / "text *" /  "*text"
gitmap db-find aliasName -table"table, table" --find/f "text" / "text *" /  "*text"
gitmap db-find aliasName -table "table, table" --grep/r "regex search in db"

gitmap db-find aliasName "table, table" --columns "col1 can be from table 1 or 2", "col2 can be from table 2 or 1"  --find/f "text" / "text *" /  "*text" --json --limit/l 100 # show results in json by default limit is 100

gitmap db-listing # 

# Command and Database Integration Instructions: high priority instruction, non-negotiable task

Hi there. Instead of having this type of command running by AI, make sure that we have these type of commands, and also test these commands with the database from the Antigravity as well. We should be able to explore, read, add things to the database with the path. Also, we should have the examples in the `gitmap` skills properly, for the `heatmap` as well in details, in the `gitmap` skills, in the coding guideline, and also in the `gitmap` book. The same way we should update, and also inside the execute V6 prompt in the coding guideline, we should update these commands as an exception to not using the Python one, but this one. And also the help should explore everything for Database (Database (Db)) read help, Database (Database (Db)) add help, Database (Database (Db)) edit help, something like this. Also, if you can think of more commands, go ahead and define that. The first thing would be defining everything. Make sure that you understood the requirement. You write the spec for the requirement. You create the help text for the terminal and also UI, and also command docs with examples. Make sure of that. And try to make sure the performance is very faster, so that the queries can run faster and things like that. Also, we could be able to see `gitmap`. Database (Database (Db)) listing should show us how many DBs we have searched, and also the alias for those DBs with the path information and what is their size. So this needs to be displayed very nicely. And also the JSON version with the flag needs to be there. Yeah, so make sure that the raw query section is available. Also, we should be able to create view in the Git Database (Database (Db)). We view and we can give alias means the database alias name that we have set. If we don't set, then we have to give the path at the first parameter anyway. And then we can say the... all kinds of example needs to be there so that AI can understand and use it effectively. Also, you test it with end-to-end so that you understand that these are working in the database and test with some sample data, do some queries, check the verifications of the data so that everything is correct. I think this type of end-to-end testing you can keep inside your system, so that can run very effectively, very fast in the CI/CD as well. You create views based on some table one, table two, and that would create a join automatically for the view. We can also give the view name. If the view name is already there, then it will just perform a query. And also we can give a find operation and give each column and where we want to find. The view find you could use the aliasing where, and we can say t1.fill equal search text. This type of search can be operated with the view and also with the view we can always do the table to see all the tables that is there in the database, get table names, and we can see the table names easily from the database. There should be a short form, which I probably did not find it or didn't give it. So you can give the table format, like how many columns are there. This type of commands needs to be there. So try to understand the examples that I've given. If you can think of new examples, go ahead and create those. And then finally, do a minor bump and do a release. Do you understand?
```

---

## 2. Five Architectural Pillars

1. **Pillar 1: High-Performance Native Go SQLite Engine**
   - Pure Go CGO-free driver (`modernc.org/sqlite`) eliminating host Python dependencies.
   - Enforce hardware/engine read-only safety via `PRAGMA query_only = ON` for read/find/view operations.
   - Centralized connection management (`SetMaxOpenConns(1)`) preventing concurrency locks and WAL mode enablement.

2. **Pillar 2: Persistent Alias Registry & Database Listing**
   - Persistent `DatabaseAlias` SQLite table in `gitmap.db` mapping human aliases to physical database paths.
   - Live filesystem telemetry tracking file byte sizes and cumulative search query frequency.
   - CLI command `gitmap db-listing` (alias `dbl`) with styled terminal output and structured `--json` envelopes.

3. **Pillar 3: Schema Introspection & Data Mutation**
   - `gitmap db-read` (alias `dbr`) supporting raw SQL execution, `ls` table summary, and `table <tbl>` column datatype introspection.
   - Two-phase table extraction resolving connection pool deadlocks in `listTablesAndViews` under `MaxOpenConns=1`.
   - `gitmap db-add` (alias `dba`) with interactive prompt and CSV insertion; `gitmap db-edit` (alias `dbe`) targeting records by ID.

4. **Pillar 4: Advanced Search, Virtual Views & Cross-Table Auto-Joins**
   - `gitmap db-find` executing wildcard prefix/suffix (`*text`, `text*`) and regex search (`--grep/-r`) across table columns.
   - `gitmap db-view` creating persistent virtual views with automatic foreign-key / naming-convention join inference.
   - Aliased column projections (`t1 as a, t2 as b`) and declarative filter expressions.

5. **Pillar 5: Antigravity E2E Testing, Skills/Prompts Ecosystem & Release Ceremony**
   - End-to-end verification against live Antigravity databases (`conversation_summaries.db`, trajectories) and hermetic mock test suites.
   - Update `.agents/skills/gitmap/` with exhaustive `heatmap` telemetry and native database command documentation.
   - Update coding guidelines and execute-v6 prompt establishing native `gitmap db-*` commands as the authorized exception replacing Python scripts.
   - Minor SemVer version bump ceremony (`v6.580.0`).

---

## 3. Discrete Subtask Allocations & Milestones

| Subtask ID | Title | Owned Spec & Files | Assigned Owner | Status | Target Deliverables |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Subtask-01** | Specification Suite & Database Registry Engine | `.ai-memory/plans/subtasks/command-and-database-integration-instructions/01-spec-and-registry-engine.md` | Worker 01 | PENDING | Complete 5-part specs in `02-spec/21-app/command-and-database-integration-instructions/`; `DatabaseAlias` schema; `gitmap db-listing` |
| **Subtask-02** | `db-read` Engine, Schema Introspection & Deadlock Regression Fix | `.ai-memory/plans/subtasks/command-and-database-integration-instructions/02-db-read-listing-and-schema.md` | Worker 02 | PENDING | `gitmap db-read` (SQL, ls, table, --columns, --json); deadlock fix for `listTablesAndViews` under `MaxOpenConns=1` |
| **Subtask-03** | `db-add` & `db-edit` Mutation Engine with Single-Writer Locks | `.ai-memory/plans/subtasks/command-and-database-integration-instructions/03-db-add-edit-and-mutation.md` | Worker 03 | PENDING | `gitmap db-add` (interactive + CSV); `gitmap db-edit` (-table, -id, csv); transactional integrity and input validation |
| **Subtask-04** | `db-find` & `db-view` Virtual Auto-Join Query Engine | `.ai-memory/plans/subtasks/command-and-database-integration-instructions/04-db-find-and-view-joins.md` | Worker 01 | PENDING | `gitmap db-find` (wildcards, regex --grep); `gitmap db-view` (declarative auto-join, view creation, filter pushdown) |
| **Subtask-05** | Antigravity E2E Verification, Skills/Prompts & Minor Release | `.ai-memory/plans/subtasks/command-and-database-integration-instructions/05-skills-prompts-and-minor-bump.md` | Worker 02 | PENDING | Live Antigravity E2E tests; `.agents/skills/gitmap` docs; coding guideline v6 prompt update; minor version bump (`v6.580.0`) |

---

## 4. Verification Contracts

- **Subtask-01:** `DatabaseAlias` migration applies cleanly; `gitmap db-listing --json` outputs array of registered databases with valid sizes and query counts; unit tests pass in `cli/cmddb`.
- **Subtask-02:** `gitmap db-read <alias> ls` correctly inspects tables; `gitmap db-read <alias> "SELECT ..."` executes with `PRAGMA query_only=ON`; `TestListTablesAndViews_DeadlockRegression_MaxOpenConns1` passes without hanging.
- **Subtask-03:** `gitmap db-add <alias> <table...> csv ...` appends rows transactionally; `gitmap db-edit` mutates targeted row by ID without data loss; write operations on read-only databases fail with descriptive errors.
- **Subtask-04:** `gitmap db-find` returns matching records for wildcard and regex patterns; `gitmap db-view` automatically links joined tables via common keys or explicit `-t` specifications.
- **Subtask-05:** `go test ./cli/cmddb/... -v` passes cleanly in CI; live Antigravity test verifies queries against `conversation_summaries.db`; skill and prompt files updated with zero broken links; minor release tagged.
