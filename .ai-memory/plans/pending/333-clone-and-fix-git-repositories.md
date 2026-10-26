# Plan: 333 - Clone and Fix Git Repositories

## User Request (Verbatim)
# Clone and Fix Git Repositories: high priority instruction, non-negotiable task

Usually here we should be able to clone by giving the repository name, and there is no suggestion here. There should be a suggestion on the repositories when we start with the cloning CFR, CRFB, or clone or reclone. Any one of the commands should have a suggestion. So by default, the SQLite will keep as a cache. If not, then it will try to search also. There is no suggestion showing up. That is a serious issue. We need to work on it. So fix that cloning. So try to clone two repositories from my Git repo using `gitmap`, fix it if required to. It should be just by name or by suggestion also came up. There is nothing. So clone the yt-manage V1 and also the go email sender. These two, clone it right now and fix the code in the `gitmap` as much as possible. And also include these repositories into the Antigravity IDE current exe automatically. Can you please do that for me?

## task: Clone and Fix Git Repositories
## slug: clone-and-fix-git-repositories

# Actionable Items Must Follow Non-Negotiable

1. Git pull BEFORE starting any work — first try `gitmap pull <repo path>` (or run `gitmap pull` in repository folder); if that fails, fall back to regular `git pull` and resolve any merge conflicts so the repository is completely up to date and clean before starting.
2. Write spec under 02-spec/21-app/XX-clone-and-fix-git-repositories/ and enqueue plan task in .ai-memory/plans/pending/XX-clone-and-fix-git-repositories.md (subtasks in .ai-memory/plans/subtasks/XX-clone-and-fix-git-repositories/) first, (xx should be figured out by the folder structure and running AI tools, xx is sequence 01, 02)
3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
5. Fix the cloning process to ensure repository suggestions appear when starting with cloning commands like CFR, CRFB, clone, or reclone.
6. Clone the `yt-manage V1` and `go email sender` repositories from the Git repo using `gitmap`.
7. Fix any issues in the `gitmap` code related to cloning by name or suggestion.
8. Include the cloned repositories into the Antigravity IDE current execution automatically.

---

## Strategic Architecture & Problem Analysis

### 1. Root Cause 1: `gitmap clone <name>` File vs Slug Confusion
In `cli/cmdclone/clone.go:236-252`, `dispatchCloneExecution` checks `isDirectURL(cf.Source)`. Because a bare repo name or slug (e.g. `yt-manage-v1`) does not start with `https://`, `ssh://`, or `git@`, it routes to `executeManifestCloneFlow(cf)`. `executeManifestCloneFlow` calls `cloner.LoadRecords(cf.Source)` which executes `os.Open(cf.Source)`. When the file does not exist locally, it fails with:
`Error: clone failed for source file yt-manage-v1: [E9000:EXECUTION] open clone source "yt-manage-v1": open yt-manage-v1: The system cannot find the file specified. (at=cloner/cloner.go:94)`
It never checked `fsutil.FileExists()`, never invoked `ResolveRepoSlug()`, and never queried SQLite `store.DB` or GitHub for repository suggestions.

### 2. Root Cause 2: `gitmap cfr <name>` Alias Collision
In `cli/cmd/roottooling_util.go:177`, `cfr` is mapped to `runCfrDisambiguated(args)`. `isCfrCloneRequest(args)` checks only if `args[0]` is a URL or begins with specific flags. If given a bare repo name like `yt-manage-v1`, it returns `false`, defaulting to `cmdcommit.RunCommitPushRefactor`. Because `D:\work` is not a git repository (`cmdpull.IsGitRepoCWD() == false`), it errors with `ERROR Not a git repository (or any of the parent directories).`
Furthermore, `crfb` does not currently exist as an alias in GitMap.

### 3. Root Cause 3: Antigravity IDE Post-Clone Sync Omission
In `cli/cmdclone/clone_direct.go:207-221`, `postCloneRegistration` registers with GitHub Desktop and VS Code Project Manager, but omits `workspacesync.SyncAntigravity(params.TargetDir, repoName)`.

---

## Multi-Agent Subtask Breakdown

| Task-ID | Subtask | Owner | Owned Files | Status | Evidence Contract |
|---|---|---|---|---|---|
| Task-01 | Spec Authoring & Architecture Design | Spec Authors 01-03 | `02-spec/21-app/333-clone-and-fix-git-repositories/*`, `.ai-memory/plans/subtasks/333-clone-and-fix-git-repositories/*` | IN_PROGRESS | Specs and subtask files authored |
| Task-02 | Clone Name Resolution, Suggestion Engine & Alias Routing | Worker 01 | `cli/cmdclone/clone.go`, `cli/cmdclone/clone_manifest.go`, `cli/cmdclone/gh_resolver.go`, `cli/cmdclone/clonefixrepo.go`, `cli/cmd/roottooling_util.go`, `cli/cmd/rootcore.go` | QUEUED | `gitmap clone <name>` resolves or suggests, `cfr`/`crfb` route to clone |
| Task-03 | Antigravity IDE Integration & Target Repos Cloning | Worker 02 | `cli/cmdclone/clone_direct.go`, Antigravity project JSONs, cloned repos | QUEUED | Repositories cloned, Antigravity project config updated |
| Task-04 | Verification, Pre-Commit Hygiene & GitMap Atomic Push | Lead Orchestrator | `cli/*`, `.ai-memory/*` | QUEUED | Clean git tree, atomic `gitmap cbf` |

---

## Coding Guidelines Compliance Matrix
- **Error Sovereignty:** Leaf functions return error directly, no exit calls.
- **Positive Booleans:** `is*`/`has*` only.
- **Function Sizing:** Functions <= 15 lines max. Zero nested ifs.
- **Relative Paths:** Only relative paths throughout all documentation, code, and plans.
