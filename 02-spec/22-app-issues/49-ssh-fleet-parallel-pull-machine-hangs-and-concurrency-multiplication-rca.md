# Issue 49 RCA: SSH Fleet Parallel Pull Machine Hangs and Concurrency Multiplication

## 1. Reproduction

1. **Triggering Fleet Pull-All:**
   - Execute `gitmap pa --ssh` (or `gitmap pae --ssh`) from the host machine.
   - Fleet contains multiple registered nodes (e.g. `w1` on node-w1, `w2` on node-w2, `w3` on node-w3, and the local host machine).
   - The fleet dispatch launches `executeLocalVMPull` locally and `executeRemoteNodePull` concurrently across all online nodes.

2. **Observed System Freeze:**
   - On the local machine and across the virtual machines (VMs), CPU utilization abruptly spikes to 100%.
   - Disk I/O queue lengths saturate physical storage bandwidth.
   - SSH terminal sessions freeze and become unresponsive; timeouts trigger on SSH channels.
   - Host machine and virtual machines become completely hung until processes are manually killed or timeout.

## 2. Root Cause Analysis

1. **Multiplication of Concurrency Across Shared Physical Subsystems:**
   - In `cli/cmdpull/pull.go`, the default parallel worker count is `0` (`h.pFlag = fs.Int("parallel", 0, ...)`).
   - When `parallel == 0`, `cloneconcurrency.Resolve(0)` sets the worker count to `runtime.NumCPU()`.
   - On an 8-core or 16-core system, each instance of `gitmap pa` launches 8 to 16 concurrent goroutines.
   - Each goroutine invokes git subprocesses (`git fetch`, `git rev-parse`, `git status`, `git pull`).
   - When the host runs `gitmap pa` locally AND simultaneously commands `w1`, `w2`, and `w3` to run `gitmap pa --json`, all nodes launch full-concurrency worker pools simultaneously.
   - When VMs run on the same physical host (e.g., VMware Workstation, Hyper-V) or share network storage, total concurrent git processes spike to 32–64+, completely exhausting CPU thread pools and disk I/O channels.

2. **Lack of Priority Throttling on SSH-Delegated Requests ("Half Priority"):**
   - In `cli/cmdssh/ssh_pull_fleet.go`:
     - Line 173 executed: `subArgs := []string{"pa", "--json"}` on the local host.
     - Line 230 executed: `crypto.RunCommand(client, "gitmap pa --json", "")` on remote nodes.
   - Neither invocation passed a throttled worker limit (`--parallel 2` or `--half-priority`).
   - Remote machines receiving commands via SSH treated the request with normal, unconstrained priority rather than recognizing it as a delegated background task that should run at reduced concurrency (half priority).

3. **Absence of CPU & Resource Awareness ("CPU Culture"):**
   - `cloneconcurrency.Resolve` purely checked `runtime.NumCPU()` without considering available system resources, virtualized environments, or whether the execution was triggered remotely over SSH.
   - In virtual machines where virtual CPU cores are oversubscribed or limited to 2–4 vCPUs, spawning 4–8 concurrent git pull operations locks up the guest OS scheduler.

## 3. Corrective Implementation

1. **Adaptive SSH & Resource-Aware Concurrency Resolver (`cli/cloneconcurrency/`):**
   - Introduce `ResolveWithPriority(requested int, isSSHRequest bool) (int, bool)`.
   - When `isSSHRequest` is true (or when environment variables `SSH_CLIENT` / `SSH_CONNECTION` are set, or `--half-priority` is passed):
     - Calculate half priority: `max(1, runtime.NumCPU() / 2)`.
     - Cap maximum concurrency for any SSH-delegated task at 2 workers (or max 3 on high-core machines).
     - Prevent thread exhaustion on shared VM hosts.

2. **Explicit Throttled Fleet Dispatch (`cli/cmdssh/ssh_pull_fleet.go`):**
   - In `executeLocalVMPull`: invoke `gitmap pa --json --parallel 2` so local execution during fleet orchestration does not starve host resources needed to maintain SSH sessions.
   - In `runRemotePullJSON`: dispatch `gitmap pa --json --parallel 2` to all remote nodes.
   - If a remote node fails due to a stale pending task, extract the task ID and automatically recover.

3. **Universal Application Across All SSH Fleet Commands:**
   - Apply half-priority worker limits across `pa --ssh`, `pae --ssh`, `update --ssh`, and cluster runners.

## 4. Prevention

1. **Unit Test Verification:**
   - Add unit tests verifying `ResolveWithPriority` correctly clamps concurrency under simulated SSH conditions.
2. **Deterministic Defaults:**
   - Never allow unconstrained `NumCPU` worker pools when delegating operations across machine boundaries.
