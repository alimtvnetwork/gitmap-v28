# Completed Plan: 69-nodes-cfr-ui-async-delegation

## User Request (Verbatim)
```text
Okay. So you can see the nodes CFR has issues. First of all, the UI is very crappy. You need to fix the terminal UI. And then, whatever is happening, we have no clue. First, it should show how many machines are active, how many are online, offline, and where the requests are sent. Okay? And also at the same time, it should actually communicate with the other, let's say, Git maps into the other machines, and also the Git map versions from those other machines needs to come and show up into the report first. So this is how the communication should start and the delegation of that. So current folder that we are in, or it will by default, if we are not inside the work directory, then it will try to clone inside the work directory. Okay? If we are inside some relative path inside the work directory, we will try to respect that and try to clone inside that work directory. Okay? So there will be by default path that user can provide, which is optional, where they could actually give the actual path where they wanted to clone. So this needs to be there. And also nodes help should have all this clone CFR help with example of clone, accept, query machine, things like that. So they should always show up at the end as a suggestion as well. Okay, so improve the performance as well. So send these to the machines as an async. Do not just wait, send one after another. Send it async so that the performance is faster and the response from this Git map should come as JSON. So make sure that you use the JSON flags so that you can display it nicely how you want inside the terminal. Do you understand the problem? Can you please follow through, make it improve, and make a minor bump, and release it?
```

## Summary of Completed Work
1. **Pre-Flight Fleet Probing & Metrics Handshake (`nodes_clone_probe.go`):**
   - Implemented concurrent pre-flight probing over SSH querying `gitmap version --json` (with fallback to `gitmap --version`) with a 3-second timeout ceiling.
   - Populated machine status metrics: Total, Online, Offline, and Auth Failed counts.
   - Cached remote GitMap versions in memory via `recordNodeVersion` for display in pre-flight and results tables.
2. **Target Directory & Workdir Resolution (`nodes_clone.go`, `nodes_clone_types.go`, `nodes_clone_remote.go`):**
   - Implemented local CWD detection relative to workdir (`D:\work` on Windows, `~/work` on Linux).
   - Preserves relative subdirectories when inside work directory (e.g. `./internal\tooling` replicates to `<remote-work-root>/internal/tooling`).
   - Defaults to root workdir when outside work directory.
   - Supports explicit destination flags (`-d`, `--dest`, `--dir`, `--target-dir`) and positional second argument `[dest]`.
   - Pre-creates remote target directories automatically (`mkdir -p` / `New-Item -ItemType Directory -Force`).
3. **Asynchronous Parallel Fleet Delegation (`nodes_clone_remote.go`):**
   - Concurrent dispatch across remote workers via goroutines and wait groups.
   - Enforces `--json` flag on remote GitMap executions, parsing JSON responses with fallback.
4. **Terminal UI Modernization (`nodes_clone_table.go`):**
   - Polished pre-flight readiness table with ANSI-padded columns (`NODE (ALIAS)`, `HOST`, `OS`, `VERSION`, `DESTINATION`, `STATUS`) and status badges (`● online`, `▲ auth_failed`, `○ offline`, `✗ unreachable`).
   - Startup banner with title, metadata, and route mappings (`• Dispatch: <alias> -> <host>:<dest>`).
   - Results table displaying remote node version, execution duration, and details.
   - Context-aware footer suggestions displaying next-step commands.
5. **Nodes Help & Catalog Registration (`nodes_clone_help.go`, `catalog.go`):**
   - Comprehensive CLI help with rich examples for `clone`, `cfr`, `cfrp`, `except-self`, `accept`, and query machine (`gitmap machine --ssh`, `gitmap nodes ping`).
   - Registered `nodes-cfr`, `nodes-clone`, `nodes-cfrp`, `cfr`, `cfrp` in `cli/helptext/catalog.go`.
6. **Minor Version Bump & Release Ceremony:**
   - Bumped minor version from `v6.453.0` to `v6.454.0` via `03-ai-scripts/37-bump-version.py`.
