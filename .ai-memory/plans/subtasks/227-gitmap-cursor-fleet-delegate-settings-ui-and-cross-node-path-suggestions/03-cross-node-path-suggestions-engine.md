# Subtask 227.03: Cross-Node Path Suggestions Engine & Installer Candidate Probing

**Subtask ID:** `227.03`  
**Parent Plan:** [227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions.md](../../227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions.md)  
**Spec References:**
- [01-architecture-spec.md](../../../../02-spec/21-app/227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions/01-architecture-spec.md)
- [02-component-and-ui-spec.md](../../../../02-spec/21-app/227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions/02-component-and-ui-spec.md)
**Status:** Ready  
**Target Area:** `install-quick.ps1`, `install-quick.sh`, `cli/completion`, `cli/cmdpipeline`  

---

## 1. Objective

Deliver robust, proactive directory and repository suggestions across both installer scripts and CLI workflows:
1. Provide candidate install path discovery in `install-quick.ps1` and `install-quick.sh` with interactive numbered selection menus.
2. Implement cross-node repository suggestions in Go that combine local indexed repositories with remote fleet node entries (e.g. `u1:<repo>`).
3. Enrich CLI diagnostics when target repositories cannot be resolved locally, pointing operators to remote fleet alternatives.

---

## 2. Implementation Steps

### Step 1: Windows Installer Probing (`install-quick.ps1`)
1. Implement `Get-CandidateInstallDirs`:
   - Inspect existing `powershell.json` deployed at standard locations (`D:\gitmap`, `C:\gitmap`, `E:\gitmap`, `$HOME\gitmap`).
   - Query `$env:PATH` for existing `gitmap.exe` installation directories.
   - Enumerate fixed storage drives with > 5 GB free disk space.
   - Deduplicate candidates and sort with the current default (`D:\gitmap` or active deployment) at position `[1]`.
2. Refactor `Read-InstallDir`:
   - Print candidate list with numbered options `[1..N]` and detailed annotations (e.g., `(Recommended - Existing Installation)`, `(System Drive)`).
   - Prompt user for input: allow Enter for default `[1]`, single-digit choice `1..N`, or an arbitrary custom path.
   - Validate and resolve path before returning.
3. Preserve non-interactive mode:
   - When `-Interactive` is omitted, default automatically to candidate `[1]` with an informational message.

### Step 2: Linux / macOS Installer Probing (`install-quick.sh`)
1. Implement `detect_candidate_install_dirs`:
   - Check user privilege: if root (`id -u == 0`), probe `/usr/local/bin` and `/opt/gitmap`. If non-root, probe `${HOME}/.local/bin` and `${HOME}/gitmap`.
   - Inspect `which gitmap` for existing binary locations.
   - Check existing `powershell.json` deploy paths in candidate roots.
2. Refactor `prompt_dir`:
   - Output numbered options menu to `/dev/tty` or stderr.
   - Read user choice safely: press Enter for default `[1]`, enter `1..N`, or provide custom path.
   - Sanitize path using `sanitize_install_dir`.
3. Guard eval and pipe execution:
   - Ensure `curl ... | bash` or `eval` without terminal tty does not stall, falling back immediately to candidate `[1]`.

### Step 3: Go Cross-Node Suggestion Aggregator (`cli/completion/cross_node_suggest.go`)
1. Create `cli/completion/cross_node_suggest.go`:
   - Define data types:
     ```go
     type CrossNodeCandidate struct {
         NodeAlias string
         RepoSlug  string
         FullPath  string
         IsRemote  bool
     }
     ```
   - Implement `CollectCrossNodeSuggestions(prefix string) []string`:
     - Retrieve local candidates from `store.DB` (`Repo` and `WorkDir` tables).
     - Retrieve registered remote nodes from `cmdssh.FetchAllSSHConnections()`.
     - Synthesize remote suggestions using `<nodeAlias>:<repo>` convention (e.g. `u1:gitmap-v28`, `u1:/home/a/git-work/gitmap-v28`).
     - Filter by prefix and return sorted unique slice.
2. Maintain strict package boundaries and keep functions under 15 lines.

### Step 4: Dynamic Completion Integration (`cli/completion/dynamic.go`)
1. Update `isRepoPathCmd` in `cli/completion/dynamic.go` to include `cursor`, `cur`, `nodes`, and `ssh`.
2. Update `Dynamic` completion handler to query `CollectCrossNodeSuggestions` when prefix matches or contains a colon (`:`).
3. Ensure performance stays O(few ms) by using cached fleet inventory and avoiding blocking remote SSH network calls during tab completion.

### Step 5: Enhanced CLI Diagnostic (`cli/cmdpipeline/pipeline_repo_suggest.go`)
1. Extend `PrintRepoTargetNotFoundDiagnostic`:
   - When a repository is not located in local SQLite indexes, query cross-node fleet nodes.
   - Print both local suggestions and remote fleet execution commands:
     ```
     ✖ Repository not found: "repo-name"

     Did you mean (Local):
       • gitmap cd repo-name

     Did you mean (Remote Fleet Node u1):
       • gitmap nodes run u1 -- gitmap cd repo-name
       • gitmap cur delegate --node u1
     ```

### Step 6: Verification & Automated Tests
1. Create unit tests in `cli/completion/cross_node_suggest_test.go`:
   - Test empty prefix, partial prefix matching, and node-scoped prefix matching (e.g., `u1:`).
   - Test deduplication and sorting.
2. Test PowerShell installer script via simulated interactive inputs (`1`, `2`, custom path, Enter).
3. Test Bash installer script in child shell.

---

## 3. Evidence Checklist
- [ ] `install-quick.ps1` discovers candidates and accepts numbered inputs.
- [ ] `install-quick.sh` presents numbered menu and preserves non-interactive pipelining.
- [ ] `cli/completion/cross_node_suggest.go` aggregates local and remote candidates without network lag.
- [ ] `pipeline_repo_suggest.go` prints formatted remote alternatives.
- [ ] All unit tests pass with 100% success.
