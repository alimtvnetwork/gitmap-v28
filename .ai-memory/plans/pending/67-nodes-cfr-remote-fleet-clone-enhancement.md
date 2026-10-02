# Execution Plan: 67-nodes-cfr-remote-fleet-clone-enhancement

## User Request (Verbatim)
```text
Okay. So you can see the nodes CFR has issues. First of all, the UI is very crappy. You need to fix the terminal UI. And then, whatever is happening, we have no clue. First, it should show how many machines are active, how many are online, offline, and where the requests are sent. Okay? And also at the same time, it should actually communicate with the other, let's say, Git maps into the other machines, and also the Git map versions from those other machines needs to come and show up into the report first. So this is how the communication should start and the delegation of that. So current folder that we are in, or it will by default, if we are not inside the work directory, then it will try to clone inside the work directory. Okay? If we are inside some relative path inside the work directory, we will try to respect that and try to clone inside that work directory. Okay? So there will be by default path that user can provide, which is optional, where they could actually give the actual path where they wanted to clone. So this needs to be there. And also nodes help should have all this clone CFR help with example of clone, accept, query machine, things like that. So they should always show up at the end as a suggestion as well. Okay, so improve the performance as well. So send these to the machines as an async. Do not just wait, send one after another. Send it async so that the performance is faster and the response from this Git map should come as JSON. So make sure that you use the JSON flags so that you can display it nicely how you want inside the terminal. Do you understand the problem? Can you please follow through, make it improve, and make a minor bump, and release it?
```

## Architectural Analysis & Scope Breakdown

### 1. Pre-Flight Machine Status & Remote GitMap Version Probing
- Before dispatching clone operations, probe all registered fleet nodes concurrently with a fast timeout (2–3 seconds).
- Query remote GitMap version using `gitmap version --json` over SSH (with fallback to parsing `gitmap version`).
- Ensure `gitmap version --json` on this machine (and updated remotes) outputs valid JSON `{"version": "...", "commit": "...", "os": "...", "arch": "..."}`.
- Parse JSON response and collect telemetry: Total machines, Online, Offline, and remote GitMap version per node.
- Render pre-flight report table first, showcasing machine status, remote GitMap versions, and reachability.

### 2. Modernized Terminal UI & Progress Telemetry
- Overhaul `renderFleetStartBanner` and `renderFleetResultsTable` in `cli/cmdnodes/nodes_clone_table.go`.
- Replace the static 3-line ASCII banner with an informative, styled dashboard:
  - Header with run parameters, target repository/file, and target directory.
  - Machine status breakdown: Total machines, Active / Online, Offline / Unreachable.
  - Destination mapping showing where requests are dispatched (`<alias> -> <ip>:<remote-dest-path>`).
  - Styled phase indicators and clean status badges (`ONLINE`, `OFFLINE`, `CLONED`, `FAILED`).
  - Next-step suggestions in the footer (`gitmap nodes ping`, `gitmap ssh <alias>`, `gitmap nodes cfr <repo>`).

### 3. Asynchronous Concurrent Multi-Node Dispatch
- Eliminate sequential blocking or unbounded hangs during remote execution.
- Implement bounded worker pool / goroutine fan-out with timeout context (`context.WithTimeout`, 3 minutes).
- Send requests to machines asynchronously (concurrently) passing `--json` to `gitmap cfr <args> --json`.
- Parse remote JSON response (`DirectCloneJSONResponse` or `jsonenvelope`) to extract success, duration, and error details for terminal rendering.
- Run local clone concurrently or clearly demarcate local and remote execution phases.

### 4. Intelligent Work Directory & Relative Path Resolution
- Check if CWD is inside the default workspace (`fsutil.IsInsideWorkDir`).
- If CWD is a relative subfolder inside the work directory (e.g. `d:\work\presentations-repos`), compute relative path (`filepath.Rel`).
- Replicate this relative path on remote machines: `<remoteWorkDir>/<relSubfolder>` (e.g. `D:\work\presentations-repos` on Windows, `~/work/presentations-repos` on Linux).
- If CWD is outside the work directory, default to the canonical work root (`D:\work` / `~/work`).
- Support an optional explicit user-provided destination path argument or flag (`-d`, `--dir`, `--dest`, `--target-dir`).
- Ensure remote directory creation (`mkdir -p` / `New-Item -ItemType Directory`) prior to clone execution.

### 5. Nodes Help Menu & Actionable Suggestions Parity
- Author `cli/helptext/nodes.md` documenting `gitmap nodes`, `nodes ping`, `nodes clone`, `nodes cfr`, `nodes cfrp`, flags, and usage examples.
- Update `printUnifiedNodesHelp` in `cli/cmd/nodes_cmd.go` to detail `cfr`, optional `[dest]`, accept keys, and query machine commands.
- Ensure footer suggestions display after `nodes`, `nodes ping`, and `nodes cfr` commands.

### 6. Verification & SemVer Minor Release Ceremony
- Run file-scoped coding guideline checks (`03-ai-scripts/05-guideline-autofixer.py`).
- Perform minor version bump to `v6.454.0` in `version.json` and sync documentation.
- Commit atomically via `gitmap cpf "nodes - modernize cfr ui remote version probe and async dispatch"`.

## Subtasks Decomposition
- **Task-01**: Terminal UI Overhaul & Pre-Flight Status Dashboard (`cli/cmdnodes/nodes_clone_table.go`, `cli/cmdnodes/nodes_clone_help.go`)
- **Task-02**: Remote GitMap Version Probing via JSON Protocol (`cli/cmdnodes/nodes_clone_remote.go`, `cli/cmd/rootutility.go`, `cli/cmdssh/ssh_node_version.go`)
- **Task-03**: Asynchronous Multi-Node Dispatch & Concurrency Performance (`cli/cmdnodes/nodes_clone.go`, `cli/cmdnodes/nodes_clone_remote.go`)
- **Task-04**: Intelligent Work Directory & Relative Path Resolution (`cli/cmdnodes/nodes_clone_file.go`, `cli/cmdnodes/nodes_clone.go`, `cli/cmdnodes/nodes_clone_types.go`)
- **Task-05**: Nodes Help Menu, Examples & Command Suggestions Alignment (`cli/helptext/nodes.md`, `cli/cmd/nodes_cmd.go`, `cli/cmdnodes/nodes_clone_help.go`)
- **Task-06**: Verification, Minor Bump & Release Ceremony (`version.json`, root `readme.md`, `02-spec/21-app/`)
