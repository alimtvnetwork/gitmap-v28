# Plan 311: Workspace-Wide Default Pull & Accelerated SSH Fleet Remote Execution

## User Request (Verbatim)
```text
PS C:\Users\Administrator> gitmap pa --ssh

  Enqueuing 'pull-all' across SSH fleet:
    • Remote Node [w3] (192.168.1.12): Offline (skipped, no task enqueued)
    • Remote Node [w1] (192.168.1.3) [GitMap v6.589.0]: Online → Enqueued (async)
    • Remote Node [w2] (192.168.1.7): Local machine (covered by local execution, skipped)
    • Remote Node [w4] (192.168.1.13): Offline (skipped, no task enqueued)
    • Remote Node [u1] (192.168.1.22): Offline (skipped, no task enqueued)
    • Remote Node [main] (192.168.1.20): Offline (skipped, no task enqueued)
    • Current Machine [AI-Main (127.0.0.1)] [GitMap v6.590.0]: Running locally (direct execution, not enqueued)

  • Local VM (127.0.0.1 - localhost): 68 pulled (42 active, 26 up-to-date)
    ...
  ▶ Node [w1] (IP: 192.168.1.3): failed
      (invalid JSON: The filename, directory name, or volume label syntax is incorrect.)

gitmap pa: execute failed: [E9000:EXECUTION] SSH fleet pull-all failed: [E_SSH_FLEET_FAIL:EXECUTION] SSH fleet pull-all had 1 node failure(s): (at=cmdssh/ssh_pull_fleet.go:63) (at=cmdpull/pull_dispatch.go:102)

PS D:\> gitmap pa
  → gitmap pull-all (cwd: D:\)
PS D:\> gitmap cd $work
D:\work
PS D:\work> gitmap pa
  → gitmap pull-all (cwd: D:\work)

FIX ALL these bugs, 
gitmap pull will run on all workspaces not on the running dir, fix the code and fix the bugs 
and remote `gitmap pa --ssh ` is super slow now figure out the root cause and fix it and release please
```

## 1. Disjoint Subtask Breakdown
| Task-ID | Subtask | Owner | Owned Files | Status | Evidence |
|---|---|---|---|---|---|
| Task-01 | SSH Fleet Windows/Linux Shell & Concurrency Fix | Worker 01 | `cli/cmdssh/ssh_pull_fleet_dispatch.go` | COMPLETED | Passes shellType to secrets.RunCommand, --parallel 4 default, resolveTargetNodeOSFast |
| Task-02 | SSH Fleet Liveness & Version Check Acceleration | Worker 02 | `cli/cmdssh/ssh_pull_fleet_liveness.go`, `cli/cmdssh/ssh_health.go` | COMPLETED | 500ms timeout without retry, stored BuildVersion first, 300ms REST timeout |
| Task-03 | Workspace-Wide Pull & CWD Isolation Removal | Worker 03 | `cli/cmdpull/pull_cwd.go`, `cli/cmdpull/pull_dispatch.go`, `cli/cmdpull/pull_dispatch_routes.go` | COMPLETED | opts.all=true default, isPullCWDEnabled requires . or cwd, (all workspaces) banner |

## 2. Invariants
- Functions <= 15 lines.
- Blank lines before return / after }.
- Positive booleans (is* or has* only).
- Strict relative paths.
- FastGate pass (exit 0).
- Minor version release: `6.590.0` -> `6.591.0`.
