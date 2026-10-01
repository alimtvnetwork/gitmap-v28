# Subtask 02: Concurrency Worker & Hands Flags and paswh Command

> **Parent Plan:** [61-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry](../../pending/61-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md)  
> **Tracking Spec:** [198-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md](../../../../02-spec/21-app/198-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md)  
> **Primary File Targets:** `cli/cmdpull/pull.go`, `cli/cmdssh/ssh_pas_fleet.go`, `cli/cmd/rootcore.go`, `cli/cloneconcurrency/resolve.go`  

---

## 1. Objective

Implement granular worker pool and hand concurrency constraints across local and SSH fleet pulls:
1. Flags `--w <N>` / `--workers <N>` to constrain worker concurrency.
2. Flags `--hand <Y>` / `--h <Y>` to constrain concurrent operations per worker.
3. Shorthand `--wwoh` (`--worker-with-one-hand`) to enforce $N=1, Y=1$.
4. Root command `gitmap paswh [N] [Y]` expanding to `gitmap pull-all --ssh --w N --hand Y` defaulting to $1 \times 1$.

---

## 2. Implementation Scope

- **`cli/cloneconcurrency/resolve.go`:**
  - Define `ResolvePullConcurrency(workerFlag, handFlag int, isWwoh bool) (int, int)`.
  - Enforce bounds: $N \ge 1$, $Y \ge 1$.
- **`cli/cmdpull/pull.go`:**
  - Parse `--w`, `--worker`, `--workers`, `--hand`, `--h`, `--hands`, `--wwoh` flags.
  - Pass resolved worker and hand counts to the parallel pull executor.
- **`cli/cmdssh/ssh_pas_fleet.go`:**
  - Support forwarding concurrency flags to remote SSH nodes when delegating pull execution.
- **`cli/cmd/rootcore.go`:**
  - Register root dispatch for `paswh` command:
    - `gitmap paswh`: defaults to $N=1, Y=1$.
    - `gitmap paswh 2`: sets $N=2, Y=1$.
    - `gitmap paswh 2 2`: sets $N=2, Y=2$.

---

## 3. Verification

- Verify `gitmap paswh` command routes correctly to `gitmap pull-all --ssh --w 1 --hand 1`.
- Verify `gitmap pull-all --w 2 --hand 1` accepts flags without syntax error.
- Verify compilation with `go build ./...` in `cli/`.
