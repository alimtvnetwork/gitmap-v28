# 03-git-operations-and-pull: Pull Worker Components, Commit Suite & Push-Fix Specification

- **Spec ID:** `03-git-operations-and-pull/02-component-spec.md`
- **Status:** `APPROVED`
- **Version:** `1.0.0`
- **Subsystem:** Pull Worker Pool, Semantic Commit CLI, Push-Fix Handler, Failure Tree Formatter
- **Dependencies:** `cli/cmdpull`, `cli/cmdpushfix`, `cli/committransfer`, `cli/repodb`
- **Version Baseline:** `v6.498.0`
- **Target Version:** `v6.499.0`

---

## 1. Component Breakdown

The Git Operations and Pull cluster comprises four primary operational components:

```
03-git-operations-and-pull/
├── 01-architecture-spec.md
└── 02-component-spec.md
```

| Component | Target Source Path | Primary Role |
| :--- | :--- | :--- |
| **Pull-All Orchestrator** | `cli/cmdpull/pull_all.go`, `cli/cmdpull/pull_worker.go` | Concurrent repository pulling, worker dispatch, and semaphore management. |
| **Commit Suite** | `cli/cmd/commit.go`, `cli/committransfer/commit.go` | Semantic flat commit execution, auto-staging, Conventional Commit generation. |
| **Push-Fix Engine** | `cli/cmdpushfix/push_fix.go` | Automated recovery for non-fast-forward pushes, rebase handling, auth retry. |
| **Failure Subtree Formatter** | `cli/cmdpull/pull_tree.go` | Tree visualization of failed repository pulls and diagnostic root causes. |

---

## 2. Pull Worker Pool & Semaphore Controller

### 2.1 Concurrency Implementation
The pull worker pool manages concurrent repository operations via Go channels:

```go
package cmdpull

import (
    "context"
    "sync"
)

type PullPool struct {
    maxConcurrency int
    semaphore      chan struct{}
    wg             sync.WaitGroup
}

func NewPullPool(concurrency int) *PullPool {
    if concurrency <= 0 || concurrency > 8 {
        concurrency = 8
    }
    return &PullPool{
        maxConcurrency: concurrency,
        semaphore:      make(chan struct{}, concurrency),
    }
}

func (p *PullPool) Dispatch(ctx context.Context, fn func()) {
    p.semaphore <- struct{}{}
    p.wg.Add(1)
    go func() {
        defer func() {
            <-p.semaphore
            p.wg.Done()
        }()
        fn()
    }()
}
```

---

## 3. Semantic Flat Commit Suite

### 3.1 Flat Commit CLI Syntax
- `gitmap commit <type> <scope> <hyphen-separated-description>`:
  Auto-stages working tree changes (`git add -A`) and produces a standardized commit message:
  `<type>(<scope>): <hyphen-separated-description>`
- Aliases: `gitmap cm`, `gitmap c`, `gitmap ca`, `gitmap commit-all`.

### 3.2 Auto-Staging Safety Filters
Before staging, the commit engine inspects the working tree for accidental binary additions or credentials (`.env`, private keys) and aborts execution if sensitive patterns are discovered.

---

## 4. Push-Fix Command Suite & Auth Recovery

### 4.1 Push Recovery Lifecycle
1. Execute `git push origin <branch>`.
2. On non-fast-forward rejection (`[rejected] (fetch first)`):
   - Invoke `git fetch origin <branch>`.
   - Run `git rebase origin/<branch>` with autostash enabled.
   - If clean: retry push.
   - If merge conflict arises: halt and output conflict details to console.
3. On authentication failure (`Permission denied (publickey)`):
   - Refresh SSH credentials via `cli/cluster/vault.go`.
   - Retry operation up to 2 times before logging error to `repodb/pipeline.db`.

---

## 5. Verification & Acceptance Criteria

```yaml
verificationGates:
  isPoolConcurrencyCappedAtEight: true
  isFlatCommitFormatEnforced: true
  isPushFixRecoveryAutomated: true
  isPositiveBooleansUsed: true
```

- [x] Concurrency cap prevents CPU thrashing and socket starvation.
- [x] Conventional Commit format strictly generated.
- [x] Push-fix restores branch synchronization automatically.
