# Specification: 173 — SSH Fleet Liveness, Error Remediation, and Rerun Help

## Status
`active`

## Domain
`cli / ssh fleet / preflight liveness / failure diagnostics / actionable remediation / agy rerun help / help routing`

## User Request (Verbatim)
```text
is it done with release ca n you please check the e2e and gitmap pe for the release?


PS D:\work\gitmap> gitmap pa --ssh

  Enqueuing 'pull-all' across SSH fleet:
    • Remote Node [alpha-win] (10.20.0.11): Enqueued (async)
    • Remote Node [beta-linux] (10.20.0.12): Enqueued (async)
    • Remote Node [gamma-mac] (10.20.0.13): Enqueued (async)
    • Remote Node [w1] (192.168.1.3): Enqueued (async)
    • Remote Node [w2] (192.168.1.7): Enqueued (async)
    • Remote Node [w3] (192.168.1.12): Enqueued (async)
    • Local VM (127.0.0.1 - localhost): Running locally

  [alpha-win|10.20.0.11] Offline: [E9000:EXECUTION] execution: node [alpha-win|10.20.0.11] is unreachable: connection timed out (at=cmdssh/ssh_target_nodes.go:28)
  Stack Trace:
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.checkRemoteNodeOnline (cmdssh/ssh_target_nodes.go:28)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.executeRemoteNodePull (cmdssh/ssh_pull_fleet.go:115)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.executeFleetPullAll.func2 (cmdssh/ssh_pull_fleet.go:89)
  [gamma-mac|10.20.0.13] Offline: [E9000:EXECUTION] execution: node [gamma-mac|10.20.0.13] is unreachable: connection timed out (at=cmdssh/ssh_target_nodes.go:28)
  Stack Trace:
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.checkRemoteNodeOnline (cmdssh/ssh_target_nodes.go:28)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.executeRemoteNodePull (cmdssh/ssh_pull_fleet.go:115)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.executeFleetPullAll.func2 (cmdssh/ssh_pull_fleet.go:89)
  [beta-linux|10.20.0.12] Offline: [E9000:EXECUTION] execution: node [beta-linux|10.20.0.12] is unreachable: connection timed out (at=cmdssh/ssh_target_nodes.go:28)
  Stack Trace:
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.checkRemoteNodeOnline (cmdssh/ssh_target_nodes.go:28)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.executeRemoteNodePull (cmdssh/ssh_pull_fleet.go:115)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.executeFleetPullAll.func2 (cmdssh/ssh_pull_fleet.go:89)

  ▶ Local VM (127.0.0.1 - localhost): 64 pulled (5 active, 59 up-to-date)
      • Antigravity-Manager             dirty
      • coding-guidelines-v24           +7/-7 (3)
      • gitmap-v28                      dirty
      • pwp-mobile                      failed
      • riseup-asia-website-project-v6  failed
  ▶ Node [alpha-win] (IP: 10.20.0.11): failed
      (offline: node unreachable)
  ▶ Node [beta-linux] (IP: 10.20.0.12): failed
      (offline: node unreachable)
  ▶ Node [gamma-mac] (IP: 10.20.0.13): failed
      (offline: node unreachable)

  ▶ [w1] (IP: 192.168.1.3): 61 pulled (16 active, 45 up-to-date)
      • alim-cv-v8                             dirty
      • bright-buddy-block                     failed
      • bsrm-presentation-hiltrax-v4           failed
      • digital-name-card                      failed
      • flat-slide-show                        failed
      • gitlogger-new-v2                       failed
      • gitmap-v28                             +55/-31 (3)
      • hiltrax-v1                             dirty
      • ki-health-ppt-v5                       failed
      • kita-social-media-content-calender-v2  failed
      • lara-publishing-v1                     failed
      • maid-app-spec-presentation-v1          failed
      • remix-of-presentation-riseup-asia-v7   failed
      • seo-packages-v2                        failed
      • slides-spec                            failed
      • wp-exam-v2                             dirty

  ▶ [w2] (IP: 192.168.1.7): 64 pulled (7 active, 57 up-to-date)
      • Antigravity-Manager             dirty
      • coding-guidelines-v24           +7/-7 (3)
      • gitmap-v28                      dirty
      • global-ppt-v1                   failed
      • hiltrax-v1                      failed
      • pwp-mobile                      failed
      • riseup-asia-website-project-v6  failed
  ▶ Node [w3] (IP: 192.168.1.12): failed
      (invalid JSON: -> gitmap pull-all (cwd: C:\Users\Administrator)
pending task already exists for pa at C:\Users\Administrator (Id 234)
flag provided but not defined: -json
...
)

PS [REPO_ROOT]>

Okay. A couple of things are very wrong. When you enqueue the task using SSH, the first thing is that you should check the machine is on and off. That's your first verification. After you verify that these machines are on, then you start the task. And also you verify which machine is this machine, where it is running from. So you write in the first summary of the status that this machine is this machine. It's not going to be enqueued. It's going to run here, and all the other machines are enqueued as async. So this is what your status should be. You did it wrong. And then, if any machine is off, there is no need to put a stack trace. That's a known thing. It's off. There's no need to put the error. It's a known error. You should handle it. Okay? Now, coming to the point, lots of, let's say, Git repo seems like failed. Now, this is where I wanted to know why it failed. So the machines which are running this Git map, that should actually send the error logs and the stack trace along with the JSON. That is what is missing. Okay? And also the next steps, how these machines can fix it, that is very, very important. And the terminal output from those machines are not very good. Okay. Try to make it better, if possible. If machines are shut down, do not just do this stuff writing nicely, like these are shut down or offline. Okay, we only are working with these machines enqueued, and then you display the current machines results. Once those are back, you reply those machines back as well. So that is the idea. This is for any SSH task, not only this one. All the SSH enqueued task, PAE update, all these things. Okay, I want you to look into all the SSH commands like this and fix it just like what I'm saying. Do you understand? Do you have any question and confusion? Let me know
```

## Architectural Design & Invariants

### 1. SSH Pre-Flight Machine Liveness & Distinct Identity
- Probe machine TCP/SSH liveness before queuing fleet tasks.
- Offline machines are labeled `Offline (skipped, no task enqueued)` without throwing or printing any stack trace.
- The machine executing the command is explicitly recognized as:
  `Current Machine [<name> (<ip>)]: Running locally (direct execution, not enqueued)`
- Only truly reachable remote nodes are enqueued asynchronously.

### 2. Rich Diagnostic Failure Reason & Actionable Next Steps
- When a repo pull/sync fails, extract the exact error category (`divergent branches`, `uncommitted conflicts`, `authentication required`, `untracked file collision`, `network timeout`).
- Transmit `errorDetails` and `remediationHint` in the JSON batch summary.
- The terminal renderer presents:
  ```text
  • <repo>  failed
      ↳ Reason: <exact git error / conflict>
      ↳ Next Step: <concrete git command to resolve>
  ```
- Self-heal stale remote tasks (e.g. `pending task already exists for pa (Id 234)`) automatically.

### 3. `gitmap rerun help` & Comprehensive Routing
- When `help`, `-h`, `--help`, or `help` in any case or position is supplied to `gitmap rerun`, `gitmap rr`, or `gitmap agy rerun`, immediately display comprehensive help guidance without attempting project matching or prompt replay.
- Guard `findProjectByFlexibleTarget` so keywords (`help`, `info`, `man`) never substring-match against project names (e.g. `strhelper`).
- Register embedded help in `cli/helptext/rerun.md`.
