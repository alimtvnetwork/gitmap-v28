# Subtask 75.4: U1 Codebase Pull, Build & Bi-Directional Communication

- **Parent Plan:** [75-gitmap-u1-ubuntu-agm-fleet-integration.md](../../pending/75-gitmap-u1-ubuntu-agm-fleet-integration.md)
- **Spec Reference:** [02-spec/21-app/205-gitmap-u1-ubuntu-agm-fleet-integration/01-architecture-spec.md](../../../../02-spec/21-app/205-gitmap-u1-ubuntu-agm-fleet-integration/01-architecture-spec.md)
- **Status:** Completed
- **Target Area:** Remote Node `U1` (`/home/a/git-work/gitmap`), `cli/cmdssh`

## Objective
Establish bi-directional development synchronization: pull updated GitMap code on node `U1`, build the native Linux binary, verify local execution in `/home/a/git-work/`, and validate communication between Windows and Ubuntu.

## Implementation Details
1. Commit and push local Windows GitMap enhancements.
2. Connect to `U1` and execute `git pull` inside `/home/a/git-work/gitmap`.
3. Compile native Linux binary: `go build -o gitmap ./cli` and link to binary path.
4. Verify execution of `gitmap version`, `gitmap tree`, and scan visualizer directly on the Ubuntu terminal.
