# Gitmap Feature Upgrades & Issues Implementation Plan

This plan details the resolution of the reported bugs, the UI/UX improvements to the CLI, and the architecture for the new Multi-Machine Join feature. It breaks the requirements down into granular tasks for execution.

## Open Questions

> [!WARNING]
> **Multi-VM Join Security & Architecture:**
> 1. Should the connection between machines be secured with TLS, or is plain TCP/HTTP acceptable for internal networks?
> 2. How should authentication work? A randomly generated token that the server prints and the client uses to join? (e.g., `gitmap join 192.168.1.50 --token ABC123`)
> 3. Does the server machine orchestrate the commands (e.g. server says "pull all", and clients pull their assigned repos), or do they just sync their repository state/file changes?

## Proposed Changes

### Phase 1: Bug Fixes & Version Normalization

*Fixing the double 'v' issue and ensuring semantic versioning is strictly followed during releases.*

- [ ] 1.1 Locate the version parsing logic in the `release` command.
- [ ] 1.2 Implement a robust version normalizer function (e.g. `NormalizeSemver`).
- [ ] 1.3 Add regex or string manipulation to strip all leading `v` or `V` characters from user input.
- [ ] 1.4 Ensure the normalized version is prefixed with exactly one `v` before committing the release.
- [ ] 1.5 Handle edge cases (empty input, pure numbers, `vvv1.0.0`).
- [ ] 1.6 Update the CI/CD scripts (e.g., PowerShell deployment scripts) to use the normalized version.
- [ ] 1.7 Add unit tests for the version normalizer covering `v1.0.0`, `vv1.0.0`, `1.0.0`, `Vv1.0.0`.
- [ ] 1.8 Verify the changelog generator correctly identifies the normalized version.
- [ ] 1.9 Ensure tag generation uses the correct normalized tag without duplicates.
- [ ] 1.10 End-to-end test a dummy release with a malformed version number.

### Phase 2: Gitmap Status UI Enhancements

*Improving the alignment, clarity, and readability of the `gitmap status` command.*

- [ ] 2.1 Audit the current `status` command terminal output logic.
- [ ] 2.2 Identify the column width calculation for `REPO`, `STATUS`, `SYNC`, `BRANCH`, `STASH`, `FILES`.
- [ ] 2.3 Implement dynamic column padding based on the longest string in each column.
- [ ] 2.4 Add minimum width constraints to the `BRANCH` column so short branch names (e.g., `[pkg] 1`) do not break alignment.
- [ ] 2.5 Increase the visual padding (distance) between columns for a cleaner look.
- [ ] 2.6 Format the `[pkg] X` branch indicators to stand out without misaligning the row.
- [ ] 2.7 Update header styling to clearly distinguish headers from row data.
- [ ] 2.8 Review colors: ensure branch names, dirty states, and clean states have accessible contrast.
- [ ] 2.9 Test the status UI with deeply nested/long repository names.
- [ ] 2.10 Test the status UI with exceptionally long branch names.

### Phase 3: Gitmap Scan & JSON Export Flow

*Making repo discovery portable and self-documenting.*

- [ ] 3.1 Update `gitmap scan` to automatically export a portable JSON file (`gitmap.json`) to the current working directory or a specified path.
- [ ] 3.2 Ensure the JSON payload contains relative paths or easily adaptable metadata for cross-machine portability.
- [ ] 3.3 Add a final summary block to the `gitmap scan` output explaining what files were generated.
- [ ] 3.4 Print explicit, copy-pasteable commands in the scan output (e.g., `To clone these repositories on another machine, run: gitmap clone gitmap.json`).
- [ ] 3.5 Document the available flags for the `clone` command in the scan summary (e.g., `--clean`, `--missing-only`).
- [ ] 3.6 Implement a `--compact` flag for `gitmap scan` that generates a minified, essential-only JSON for easier transfer.
- [ ] 3.7 Add instructions on how to securely transfer the JSON file between machines via SSH or network.

### Phase 4: Gitmap Clone Upgrades & Parallel UI

*Creating a 30-step UI and logic upgrade for parallel cloning/pulling.*

**Logic Upgrades:**
- [x] 4.1 Parse the `gitmap.json` file in the `clone` command.
- [x] 4.2 Check for existing directories before attempting a `git clone`.
- [x] 4.3 Implement default safe mode: if repo exists, perform a `git pull` instead of cloning.
- [x] 4.4 Implement `--clean` flag: forcefully delete the local folder and re-clone.
- [x] 4.5 Implement `--missing-only` flag: skip existing directories entirely.
- [x] 4.6 Implement interactive prompting if `--clean` is NOT passed and conflicts are detected.
- [x] 4.7 Catch `git pull` failures in existing repos and queue them into an error summary.
- [x] 4.8 Ensure failed pulls do not crash the entire batch operation.

**UI Upgrades (30-Step Plan for Terminal Output):**
- [x] 4.9 Hide the default stdout/stderr of the raw `git clone` / `git pull` commands.
- [x] 4.10 Initialize a terminal UI library (e.g., `pterm`, `bubbletea`, or simple ANSI escape codes).
- [x] 4.11 Render a fixed header: `[gitmap] Processing X repositories...`
- [x] 4.12 Create a dynamic layout area for active workers.
- [x] 4.13 Configure a bounded worker pool (e.g., 5-10 concurrent workers) to prevent network/disk thrashing.
- [x] 4.14 Assign a UI line for each active worker.
- [x] 4.15 Display a spinner for active operations (e.g., `⠋ Cloning repo-A...`).
- [x] 4.16 Display a progress percentage if determinable, or a time elapsed counter.
- [x] 4.17 Once a worker finishes, replace the spinner with a green checkmark `✓ repo-A cloned in 2.1s`.
- [x] 4.18 If a pull is executed instead of a clone, show `✓ repo-B updated (pull)`.
- [x] 4.19 If an error occurs, show a red cross `✗ repo-C failed (merge conflict)`.
- [x] 4.20 Move completed items to a scrolling "completed" list above the active workers.
- [x] 4.21 Keep the active workers pinned to the bottom of the terminal output.
- [x] 4.22 Ensure line-clearing ANSI codes work correctly on both Windows (cmd/pwsh) and Unix terminals.
- [x] 4.23 Upon total completion, clear the active worker area.
- [x] 4.24 Print a summary box: `Completed: 12, Skipped: 1, Failed: 1`.
- [x] 4.25 List specific errors below the summary box with actionable next steps.
- [x] 4.26 Ensure the terminal cursor is restored and visible upon exit.
- [x] 4.27 Handle `Ctrl+C` gracefully, canceling all workers and cleaning up the terminal UI.
- [x] 4.28 Add visual padding (empty lines) around the summary for aesthetic breathing room.
- [x] 4.29 Use distinct colors for `Clone`, `Pull`, `Skip`, and `Fail` events.
- [x] 4.30 Test the UI with large JSON files (50+ repos) to ensure smooth scrolling.
- [x] 4.31 Verify UI performance when operations complete instantaneously.
- [x] 4.32 Trigger `gitmap status` automatically on the target directory upon successful completion of the clone/pull batch.

### Phase 5: Parallel Git Push/Pull All UI

*Extending the clone UI to standard push/pull all commands.*

- [x] 5.1 Adapt the bounded worker pool from Phase 4 for `gitmap pull --all`.
- [x] 5.2 Adapt the bounded worker pool for `gitmap push --all`.
- [x] 5.3 Implement the same dynamic spinner UI for push/pull operations.
- [x] 5.4 Ensure credentials prompts (if SSH agent / credential manager isn't loaded) don't break the parallel UI.
- [x] 5.5 If a repo is fully up to date, show a gray `~ repo-A (up to date)`.
- [x] 5.6 Automatically trigger `gitmap status` on the directory once all parallel pushes/pulls finish.

### Phase 6: Multi-Machine Join (Networked Gitmap)

*Research and implementation plan for the Kubernetes-like VM joining feature.*

**Research & Architecture:**
- [x] 6.1 Evaluate communication protocols: gRPC vs. HTTP/REST vs. WebSockets for CLI-to-CLI communication.
- [x] 6.2 Design the Node Roles: `Server` (Orchestrator) and `Client` (Worker).
- [x] 6.3 Determine Network Discovery: Can clients find the server via mDNS broadcasting on the LAN, or is explicit IP binding required?
- [x] 6.4 Design Auth/Security: Implement a generated Join Token to prevent unauthorized machines from joining the cluster.
- [x] 6.5 Define the Workload: How are tasks distributed? Does the server trigger a scan, split the repos, and assign them to clients to clone/build?

**Implementation Steps:**
- [x] 6.6 Implement `gitmap serve` command to start the orchestrator daemon.
- [x] 6.7 Bind the server to a network interface and display the IP/Port.
- [x] 6.8 Generate and display a secure `Join Token`.
- [x] 6.9 Print the exact join command: `gitmap join <IP>:<PORT> --token <TOKEN>`.
- [x] 6.10 Implement `gitmap join` command.
- [x] 6.11 Establish connection and handshake between Client and Server.
- [x] 6.12 Build the connection heartbeat/ping mechanism to detect dropped nodes.
- [x] 6.13 Create the server-side registry of connected nodes (VMs).
- [x] 6.14 Expose a `gitmap cluster status` command to view connected machines.
- [x] 6.15 Implement workload distribution algorithm (e.g. Server has 100 repos, connected to 4 VMs -> assigns 25 repos to each VM to parallelize heavy workloads).
- [x] 6.16 Enable command broadcasting: Server sends a `PullAll` event, and all clients execute it on their local copies.
- [x] 6.17 Stream client logs/progress back to the server.
- [x] 6.18 Display aggregated cluster progress on the server's terminal UI.
- [x] 6.19 Implement graceful node disconnects.
- [x] 6.20 Handle node failures (reassign workloads if a node drops during execution).

## Verification Plan

### Automated Tests

- Test version normalization logic covering all permutations of `v`.
- Test `gitmap.json` deserialization and path resolving logic.

### Manual Verification

- Simulate a multi-repo directory, run `gitmap status` and verify alignment with varying branch name lengths.
- Run `gitmap clone gitmap.json` with a mix of existing repos, missing repos, and clean flags to verify interactive prompts and safe-pulls.
- Test the parallel UI by throttling network speeds to ensure spinners and progress bars render correctly without artifacts.
- Deploy a server and two client instances (using `localhost` different ports or local Docker/VMs) to verify the cluster join handshake.
