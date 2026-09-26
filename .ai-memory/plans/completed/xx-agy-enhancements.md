# Subtask 01: `create-repo` Enhancements & `common` Command
Traceability ID: Task-01
Spec Reference: [02-spec/21-app/xx-agy-enhancements/01-overview.md](../../../02-spec/21-app/xx-agy-enhancements/01-overview.md)
Target Files: `cli/cmd/sync.go`, `cli/constants/constants.go`, `cli/cmd/repo_create_init.go`, `cli/cmd/repo_create_params.go`, `cli/cmd/repo_cmd_dispatch.go`, `cli/cmd/recreate_repo.go` (new)
Action: 
1. Rename all instances of `gitmap sync` to `gitmap common` in `cli/cmd/sync.go` and `cli/constants/constants.go`. The behavior stays exactly the same, only the command changes.
2. In `cli/cmd/repo_create_params.go`, add `IsCommon bool` and `IsCG bool` for `--common` and `--cg`.
3. In `cli/cmd/repo_create_init.go`, if `p.LocalDir` exists, do not error; simply ensure it's a git repository by running `git init` (if `.git` is missing). If `--common` is true, run `gitmap common all` programmatically or via exec, then commit it. If `--cg` is true, backup current branch (`git checkout -b backup-cg-sync`), copy `02-spec/02-coding-guidelines/` into the repo, and commit.
4. Create `cli/cmd/recreate_repo.go` which provides `gitmap recreate-repo <name>`, taking an existing git repo, duplicating it cleanly (or stripping remote and assigning new remote) onto the parent GitHub account using `gh repo create`.
Acceptance Criteria: 
- `gitmap common` replaces `sync`.
- `gitmap create-repo <existing> --common --cg` works safely on existing folders.
- `gitmap recreate-repo` clones/duplicates to a new remote.
Targeted Verification: `python 03-ai-scripts/05-guideline-autofixer.py cli/cmd/sync.go cli/cmd/repo_create_init.go`
# Subtask 02: `agy ls` and Advanced Queue Listing
Traceability ID: Task-03
Spec Reference: [02-spec/21-app/xx-agy-enhancements/01-overview.md](../../../02-spec/21-app/xx-agy-enhancements/01-overview.md)
Target Files: `cli/cmdagy/agy_ls.go`, `cli/cmdagy/agy_ls_table.go`, `cli/cmdagy/agy_most_conv.go` (new)
Action: 
1. Expand `gitmap agy ls N` and `gitmap agy conv ls N` (default 8) to parse an optional `-f` or `--file` JSON export flag.
2. In `cli/cmdagy/agy_ls.go` (or `agy_active_prompts.go`), update the querying logic to order conversations by: 1) Currently running, 2) Most recent. If invoked inside a project repo, limit the list to conversations rooted in this project path, and parse the prompt logs to show a 200-character prefix of the pending/running prompt in the table or JSON.
3. Introduce `gitmap agy most-conv ls N` which groups and sorts repositories by the highest count of conversations, outputting to a table or JSON file (e.g. `most-repo-conversations.json`).
4. Apply 1-minute caching for sequence IDs.
Acceptance Criteria: 
- `gitmap agy ls 5` shows 5 conversations sorted by running then recency.
- In-project `ls` shows 200-character prompt excerpts.
- `gitmap agy most-conv ls 10 -f` outputs top 10 projects by conversation volume to JSON.
Targeted Verification: `python 03-ai-scripts/05-guideline-autofixer.py cli/cmdagy/agy_ls.go`
# Subtask 03: `watch-prompts` and Cross-Node SSH (`lapp`/`wapp`)
Traceability ID: Task-04
Spec Reference: [02-spec/21-app/xx-agy-enhancements/01-overview.md](../../../02-spec/21-app/xx-agy-enhancements/01-overview.md)
Target Files: `cli/cmdagy/agy_watch_prompts.go` (new), `cli/cmdagy/agy_look_prompts.go` (new), `cli/cmdagy/agy_lapp.go` (new)
Action: 
1. Implement `gitmap agy watch-prompts`: Loops every 30 seconds blocking the terminal, displaying running/pending prompts for the current repo with 200 chars. Supports `--json`, `-f`, `--compact <words>`, `--count <N>`, `--all`.
2. Implement `gitmap agy look-prompts`: A one-off snapshot equivalent of `watch-prompts`.
3. Add `--ssh` flag to `watch-prompts` and `look-prompts` which uses `ssh` to retrieve the JSON output from all configured nodes in parallel and aggregate them.
4. Implement `gitmap agy look-all-projects-prompts` (`lapp`) and `wapp` which bypass the current project filter and summarize prompts across *all* projects globally (also supporting `--ssh`).
Acceptance Criteria: 
- `watch-prompts` blocks and refreshes every 30s.
- `lapp --ssh` aggregates prompt JSON from remote nodes.
- Flags like `--compact` cleanly truncate the string in output.
Targeted Verification: `python 03-ai-scripts/05-guideline-autofixer.py cli/cmdagy/agy_watch_prompts.go`
# Subtask 04: `inject-prompts` and `enhance-prompts` Engine
Traceability ID: Task-05
Spec Reference: [02-spec/21-app/xx-agy-enhancements/01-overview.md](../../../02-spec/21-app/xx-agy-enhancements/01-overview.md)
Target Files: `cli/cmdagy/agy_inject_prompts.go` (new), `cli/cmdagy/agy_enhance_prompts.go` (new), `cli/cmdagy/agy_prompt_subcmds.go`
Action: 
1. `inject-prompts` (`ip`, `ipt`, `ipr`): Read text or folders (sorting sequentially `01-..`, `02-..`). Prefix/Suffix templating engine (load from `01-migration-templates.md` or a config map). Add loop mechanism for `--rerun Y`. Add `--watch` flag to wait until the prompt completes (querying state) before injecting the next one.
2. `enhance-prompts` (`ep`): Take a text prompt or file, invoke AGY (no project mode or a dummy project), and break the prompt down into `01-index.md`, `02-verbatim-instructions.md`, and subsequent task `.md` files. Store in the requested folder or `$variable` (supporting exported env vars).
3. `iprs` / `ip ssh`: Support cross-node targeting using `--ssh-only-node <nodeid>`.
Acceptance Criteria: 
- `gitmap agy ip <folder> --prefix default --rerun 2` loops over files, prefixes templates, and queues them.
- `gitmap agy ep "text"` generates sequenced Markdown files.
Targeted Verification: `python 03-ai-scripts/05-guideline-autofixer.py cli/cmdagy/agy_inject_prompts.go`
# Subtask 05: Verify E2E and Test Setup
Traceability ID: Task-01
Spec Reference: [02-spec/21-app/xx-agy-enhancements/01-overview.md](../../../02-spec/21-app/xx-agy-enhancements/01-overview.md)
Target Files: `cli/tests/heavy_test/agy_enhancements_e2e_test.go` (new)
Action: 
1. Write a temporary E2E test `TestAGYEnhancements_E2E` inside `cli/tests/heavy_test/` utilizing the `os/exec` pattern to execute the newly created commands:
   - Create a dummy existing folder, run `gitmap create-repo <folder> --common` and verify it initializes correctly.
   - Run `gitmap agy conv ls 5 -f out.json` and verify the JSON exists.
   - Run `gitmap agy ipt "Hello world" -p "test"`.
2. Skip by default (`if os.Getenv("HEAVY_E2E") == "" { t.Skip() }`).
Acceptance Criteria: 
- Temporary integration test proves the commands don't crash.
Targeted Verification: `python 03-ai-scripts/05-guideline-autofixer.py cli/tests/heavy_test/agy_enhancements_e2e_test.go`
