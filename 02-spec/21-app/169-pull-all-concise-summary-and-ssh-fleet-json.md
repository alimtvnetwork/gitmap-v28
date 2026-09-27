# 169. Pull-All Concise Summary Filtering, Noise Elimination & SSH Fleet JSON Orchestration

## 1. Executive Summary & Architectural Overview

This specification establishes two critical enhancements to GitMap's repository orchestration and SSH fleet synchronization:
1. **Pull-All Concise Summary Filtering & Zero Log Clutter:**
   - In `gitmap pa` / `gitmap pull-all`, suppress all verbose `remote.origin.url` diagnostic lines during pre-flight discovery.
   - Filter the post-pull summary table to strictly display repositories that have active changes (`+X/-Y`), failures, or dirty worktrees. Repositories with status `up-to-date` are omitted from the summary table to eliminate terminal bloat across large 60+ repository workspaces.
2. **SSH Fleet Synchronization with JSON Communication (`gitmap pa --ssh`):**
   - When `--ssh` is passed to `gitmap pa`, `gitmap pull all`, or `gitmap agm update`, GitMap dispatches tasks across all registered cluster fleet nodes.
   - The primary node displays all remote worker IPs where tasks have been enqueued, explicitly shows that the local VM is running concurrently, and sends execution commands formatted to output JSON (`gitmap pa --json`).
   - Responses returned from remote nodes via SSH are parsed as JSON and aggregated into a unified, concise status table on the terminal.

```mermaid
flowchart TD
    subgraph LocalPull [1. Local Pull-All (Concise)]
        A["gitmap pa (cwd: workspace)"] --> B["Quiet Pre-flight (Suppress remote.origin.url)"]
        B --> C["Concurrent Pull Progress Bar"]
        C --> D["Filter Summary: Hide up-to-date"]
        D --> E["Render Only: Updated (+X/-Y), Dirty, Failed"]
    end
    subgraph SSHFleet [2. SSH Fleet Pull-All (--ssh)]
        F["gitmap pa --ssh"] --> G["Discover Cluster Fleet Nodes (cluster.LoadFleet)"]
        G --> H["Print Fleet Enqueue Banner: Node IPs + Local VM"]
        H --> I["Async SSH Execution: gitmap pa --json"]
        H --> J["Concurrent Local Execution: gitmap pa --json"]
        I & J --> K["Receive JSON Payloads"]
        K --> L["Aggregate & Render Unified Clean Summary Table"]
    end
```

---

## 2. User Request (Verbatim)

```text
PS [workspace]> gitmap pa --ssh

  → gitmap pull-all (cwd: [workspace])
    → resolved 64 repo(s) to pull

→ remote.origin.url: git@github.com:alimtvnetwork/ai-empathy-prompt-tuner-v1 → git@github.com:alimtvnetwork/ai-empathy-prompt-tuner-v1.git
→ remote.origin.url: https://github.com/alimtvnetwork/alim-cv-v8 → git@github.com:alimtvnetwork/alim-cv-v8.git
→ remote.origin.url: git@github.com:alimtvnetwork/alim-karim-profile-v2 → git@github.com:alimtvnetwork/alim-karim-profile-v2.git
... [64 verbose remote.origin.url lines] ...
  ⠴ [████████████████████████] 100% (64/64 repos) | ↓ riseup-asia-website-project-v6: pulling remote ch...

    • ai-empathy-prompt-tuner-v1             up-to-date
... [60 up-to-date lines] ...
    • wp-exam-v2                             +1266/-131 (10)
    • pwp-mobile                             failed
    • riseup-asia-website-project-v6         failed

  ✓ Pull all complete: 64 pulled (32.5s)
  ── Pull Remediation Options ──
...

I think first thing we should appreciate that the SSH Git clone is working. That's really nice. Now let's come to the issue, what I do feel like is very wrong. The first of all, if you run the Git map PA in a local directory, it will only show at the end whichever that actually has issues, failed, and the updated ones. Other ones, just ignore from the last summary. That's the first thing. Okay. Now, coming to the point, when we do the SSH, I mentioned this before, the communication needs to be done using JSON. Okay? The first thing it'll do is that definitely it would check that the Git map is installed and up to date. If not, then it would do that. So it did nicely how it represents the Git map SSH execute IP, just like how we do it. The way that it shows the commands that has been executed, the similar way it would do it here. But first thing it would say, "Yes, here we have enqueued the task into this other worker," and the current VM, it is going to run itself to be here. And the current VM is the one that it is this. So we are not sending the request, we're running it here. It mentions this. So all these IPs would show up where the request has been sent first. Okay? And then it would show the current one is also running. So once those async requests send it to those IPs, and using SSH would come back as a response. When it comes back on the async process, then we combine it and finally give a nice output, just like the current one from the JSON that we have received to the terminal. Okay, so the current response that you have done, it's trashy, and I'm not sure why there is this GitHub origin URL that we have to add, why we have to add this and display it everywhere. I don't know, so correct me if I'm wrong. Okay, so make sure that you do not do this type of output in the SSH. Whenever the SSH is mentioned. So all of these like update all or Git map AGM update hyphen hyphen SSH, then it indirectly send the request through these machines as a SSH, but it would receive the response as JSON. Remember that. And when it receives the JSON, it will just update what is the status of it nicely. Do you understand? Can you please help me with this?
```

---

## 3. Data Contracts & Architecture Requirements

### 3.1 Pull-All Summary Filtering
1. **Omit `up-to-date` items:**
   In `cli/cmdpull/pullall.go` and `cli/cmdpull/pull.go`:
   - A repository is included in the terminal summary if and only if:
     - `state.Status == PullStatusFailed` (or error is non-nil)
     - `state.Status == PullStatusDirty`
     - `state.Changes != ""` and `state.Changes != "up-to-date"`
   - If 60 of 64 repos are `up-to-date`, the final summary displays only the 4 changed/failed/dirty repos, accompanied by a count:
     `✓ Pull all complete: 64 pulled (4 changed/failed, 60 up-to-date) (32.5s)`.
2. **Suppress `remote.origin.url` Log Flooding:**
   - Remove or guard the `fmt.Printf("→ remote.origin.url: ...")` printing so routine pulls remain completely clean.

### 3.2 SSH Fleet JSON Protocol
1. **Command:** `gitmap pa --ssh` (alias `gitmap pull-all --ssh`, `gitmap pull all --ssh`)
2. **Node Enqueue Display:**
   ```text
   Enqueuing 'pull-all' across SSH fleet:
     • Remote Node [alpha-win] (10.20.0.11): Enqueued (async)
     • Remote Node [beta-linux] (10.20.0.12): Enqueued (async)
     • Remote Node [gamma-mac] (10.20.0.13): Enqueued (async)
     • Local VM (127.0.0.1 - localhost): Running locally
   ```
3. **Execution & JSON Parsing:**
   - Remote command dispatched via SSH: `gitmap pa --json`
   - Remote stdout parsed into `PullBatchJSONSummary`:
     ```json
     {
       "total": 64,
       "pulledCount": 64,
       "successCount": 62,
       "failedCount": 2,
       "durationMs": 32500,
       "states": [
         { "repoName": "wp-exam-v2", "status": "updated", "changes": "+1266/-131 (10)" },
         { "repoName": "pwp-mobile", "status": "failed", "changes": "error: merge conflict" }
       ]
     }
     ```
   - Terminal formats a consolidated multi-node table displaying only active changes and failures across all nodes.
