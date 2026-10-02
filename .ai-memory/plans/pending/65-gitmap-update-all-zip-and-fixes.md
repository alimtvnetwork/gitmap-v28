# Execution Plan: 65-gitmap-update-all-zip-and-fixes

## User Request (Verbatim)
```text
gitmap update all --include-others
gitmap update all zip --include-others
gitmap update-all-zip --include-others (uaz)

So I think we need to fix some of the Git map issues. First of all, if we do the Git map nodes, it shows this table. I think the top table that we see, that should be in the bottom part. That's the first thing. Second is that the aliasing worker and things like that, we don't need the worker role in the second table. That's the first thing. Second subsystem is fine. Okay. All right. So you can actually, if you reduce the roles column, the table looks nice. Okay, that's the first observation. The second thing is that in the Git map, there is no way to update all the ADM exe. That's another issue I think you need to work on. Okay? Another command that I want you to introduce is the update all zip. Okay? Update all zip. You can also use include others with the update zip, or let's say Git map update all. If we use the include others, then it will try to download the AGM package on all other machines as well.




So how this update all in Zipl will work? The first time, so we need to have a system that I think you need to make updates in other installer with some flags. Way that it's going to work is that it's going to update first time in the current machine and zip the uploaded package and deliver to all these other machines using SCP, okay? And put it to a temp directory and then run installer script, depending on the OS. It could be PowerShell, it could be Shell, and it would install an update in all these machine together. So that is powerful because that also saves the network and other stuff, to reduce traffic, and that is quite powerful. Can we do that? That's the first thing. Second is that can you design this system, and can you make sure that it works? You can also test it in other machines, other worker machines, okay? 



Also for the Git map, I want you to test and understand how the cache is removed and let's say, updated during the AUM search. And also I want you to check these searches that originally works are now. So I want you to find out similar, so what you can find a test case in our system, it's not like these images needs to be even I request you not to put like Amazon, Prize of Asia, do not put these names. Just put different names that you get the idea, and I will make sure that it will work. Now, for the prompt section. Now for the prompt section, I'm coding that into the V6 prompt. As I mean, execute the current task with the n8n steps V6, that prompt. I think that already has the Git map to run the command through Git map. So I think this type of code example needs to be there. So rather than writing like this, one needs to write, or AI needs to write like using Git map so that everything is available. But one more thing I do see here is that why do we have the same package repeated twice?




If let's say package is repeated, then the best way to deal with it is to use its two different paths so that I can understand. That's the first thing. Second, if a package failed, and you mentioned that it failed, next step, you do a sub-command-- I mean, sub-node, and then provide the solution if you have two different solutions. Okay? Try to do that. That would be effective and helpful. One more thing. The .git map should not be ignored. Inside the .git map, the backup folder should be ignored. Okay? But user can choose to ignore if they like. Does this make sense?




So, a little bit you need to work on the AGM package as well, because AGM, this represents Anti-Gravity Manager, so that you can see that the Anti-Gravity Manager is also exported. First download it and then export it using SCP, and then download on all these machines. You can test it from one machine to the another, which is absolutely fine. So let me know your thoughts. What do you think? And can you implement this? Firstly, start all these tasks as I have mentioned. Write your spec first and then also 3D plan modes. This is very important. And do you understand
```

## Architecture Strategy & Discovery Synthesis

### 1. GitMap `nodes` Command Table Layout & Role Column Pruning
- **Current Layout**: Primary status table (ALIAS, ROLE, HOST, USER, STATUS, ENROLLED) is rendered first, followed by Commands Matrix (NODE, ROLE, SUBSYSTEMS, SUPPORTED COMMANDS & CLUSTERS).
- **Target Layout**:
  - Invert rendering order in `renderUnifiedNodesTable` (`cli/cmd/nodes_cmd.go`): render Commands Matrix first, followed by primary node status table with `Total: X registered node(s)` at the bottom.
  - Remove redundant `ROLE` column from Commands Matrix (`NODE (ALIAS) | SUBSYSTEMS | SUPPORTED COMMANDS & CLUSTERS`). Expand width of `SUPPORTED COMMANDS & CLUSTERS` column.
  - Update `cli/cmd/nodes_cmd_test.go` string assertions to verify inverted sequence and omitted role column.

### 2. Multi-Node Update Engine & Zip SCP Distribution (`update all zip`, `uaz`, `--include-others`)
- **Commands**:
  - `gitmap update all --include-others`
  - `gitmap update all zip --include-others`
  - `gitmap update-all-zip --include-others` (alias: `uaz`)
- **Architecture**:
  - Update `cli/cmd/rootutility.go` to route `uaz`, `update-all-zip`, and positional `zip` target cleanly.
  - Update `FleetUpdateOptions` in `cli/cmdupdate/update_fleet.go` with `IsZip: bool` and `IncludeOthers: bool`.
  - When `--include-others` is specified: expand targets beyond SSH connections to include cluster hosts and nodes (`collectUnifiedNodes` / `dbpkg.ListClusterNodes`).
  - When `IsZip` is specified:
    - Package/zip the updated binary or installer locally.
    - Transfer zip to remote nodes via SSH stream (`StreamFileToRemote`) into temp directory (`C:\Windows\Temp` or `/tmp`).
    - Run remote PowerShell script (Windows) or shell script (Unix) to unpack and install.
  - AGM (Anti-Gravity Manager) Integration: Support AGM asset resolution, packaging, and SCP distribution to worker nodes.

### 3. AUM Search Cache Lifecycle & Safe Test Fixtures
- Search cache operates on two levels:
  - In-memory `MemoryCache` via `gitmap aum cache clear / warm`.
  - SQLite Split-DB `SearchHotCache` in `.gitmap/data/search/sql.db` via `ComputeDH2D` and auto-promotion to `HOT_MEMORY_CACHE` on 2nd hit.
- Wire `gitmap search clean` to explicitly purge `SearchHotCache` and vacuum the search database.
- Create tests in `cli/tests/e2e/ssh_agy_nodes_agm_aum_tempe2e_test.go` using safe synthetic names (`sample_service_node`, `generic_alpha_repo`, `test_fixture_beta`), strictly avoiding proprietary names (no Amazon, no Prize of Asia).

### 4. Canonical Prompts: Replace Raw `rg` / `grep` with GitMap Commands
- Update `01-prompts/14-execute/13-execute-parent-task-with-n-steps-v6.md` and `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`.
- Enforce explicit ban on raw `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, and `findstr`.
- Provide concrete GitMap replacement examples (`gitmap aum search "<symbol>" [dir] [-e <.ext>]`, `gitmap search`, `gitmap lf`, `gitmap cat`).

### 5. Repository Deduplication & Sub-Node Diagnostics with Solutions
- In `cli/cmdpull/pull_efficient.go`, deduplicate repository records by clean lowercase path (`filepath.Clean(strings.ToLower(r.AbsolutePath))`) so Windows case variation doesn't cause duplicate lines.
- In `cli/cmdpull/pull_efficient_render.go`, when repositories have distinct paths, render the path alongside the name.
- When pull fails, output structured sub-nodes with up to two alternative solutions:
  - Diverged branch: Option 1 (Preserve Local / Rebase) vs Option 2 (Discard Local / Hard Reset).
  - Dirty tree: Option 1 (Commit WIP) vs Option 2 (Stash).

### 6. `.gitignore` Policy: Retain `.gitmap/`, Ignore Only `.gitmap/backup/`
- Update `cli/cmdpull/pull.go`: Remove `.gitmap/` from `buildDefaultTrackedPathsToCheck`. Only check `.gitmap/backup/` and resume JSON files.
- Update `cli/cmdignore/ignore_groups.go`: Update `ImmutableDefaultRules` to enforce only `.gitmap/backup/`, allowing legitimate project files inside `.gitmap/`.
- Deduplicate ignore patterns by stripping leading/trailing slashes during normalization.

---

## Discrete Subtask Definitions

### Task-01: Nodes Table Reordering and Role Column Pruning
- **Owner:** Worker 01
- **Files:** `cli/cmd/nodes_cmd.go`, `cli/cmd/nodes_cmd_test.go`

### Task-02: Multi-Node Update Engine & Zip SCP Distribution
- **Owner:** Worker 02
- **Files:** `cli/cmd/rootutility.go`, `cli/cmdupdate/update_fleet.go`, `cli/cmdupdate/update_fleet_test.go`

### Task-03: AUM Search Cache Lifecycle & Safe Test Fixtures
- **Owner:** Worker 01
- **Files:** `cli/cmd/search.go`, `cli/tests/e2e/ssh_agy_nodes_agm_aum_tempe2e_test.go`

### Task-04: Canonical Prompts: Replace Raw `rg` with GitMap Commands
- **Owner:** Worker 02
- **Files:** `01-prompts/14-execute/13-execute-parent-task-with-n-steps-v6.md`, `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`

### Task-05: Repository Deduplication & Sub-Node Solutions
- **Owner:** Worker 01
- **Files:** `cli/cmdpull/pull_efficient.go`, `cli/cmdpull/pull_efficient_render.go`, `cli/cmdpull/pull_remediation_hint.go`

### Task-06: `.gitignore` Policy: Retain `.gitmap/`, Exclude Only `.gitmap/backup/`
- **Owner:** Worker 02
- **Files:** `cli/cmdpull/pull.go`, `cli/cmdignore/ignore_groups.go`, `cli/cmdignore/fix_ignore.go`
