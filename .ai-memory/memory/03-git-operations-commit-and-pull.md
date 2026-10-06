# 03 — Git Operations: Flat Commit Suite, Concurrency & Pull Resilience

- **Domain:** Git Operations, Flat Commit Mechanics & Pull Worker Pools
- **Authoritative Specification:** [03-git-operations-and-pull](../../02-spec/21-app/03-git-operations-and-pull/01-architecture-spec.md)
- **Status:** Active & Ratified

---

## 1. Flat Commit Suite (`gitmap commit`, `gitmap c`, `cm`)

GitMap streamlines Git workflows through high-level semantic commit commands:

- **Flat Aliases:** Direct shortcuts `c` and `cm` map to `gitmap commit`.
- **Automatic Staging:** Detects modified and untracked files and stages them automatically when `--all` (`-a`) is specified.
- **Interactive Remediation:** If the working directory has conflicts, detached HEAD, or unstaged deletions, the commit engine presents surgical remediation prompts rather than aborting abruptly.
- **Normalized Trailers:** Automatically appends semantic trailers (e.g., `Signed-off-by`, `Co-authored-by`, `Task-ID`) based on project configuration.

---

## 2. Commit-In and Commit-Right Engine

- **Path Resolution:** Operates on normalized relative Git paths, resolving paths relative to the Git repository root regardless of the active working directory.
- **Append-Only History Invariant:** Normal commits append directly to the current branch tip. Commit history rewriting is strictly prohibited during operational flows.
- **Delta Extraction:** Commits record file path modifications and change sizes in SQLite telemetry for fast change audits without reading Git packfiles.
- **Zero Raw Blobs in DB:** Database tracking stores only relative path strings and metadata; blob retrieval is delegated to `git cat-file` on demand.

---

## 3. PAS Worker Concurrency & Pull Pool Architecture

For multi-repository operations (e.g., `gitmap pull-all`, `gitmap status-all`), GitMap employs the Parallel Async Scheduler (PAS):

- **Concurrency Sizing Formula:**
  $$\text{Workers} = \min(\text{Logical Cores}, 8)$$
  *(Falls back to $\lfloor\text{Cores} / 2\rfloor$ under low-memory environments).*
- **Adaptive Worker Pool:** Distributes repositories evenly across workers, preventing disk I/O thrashing and network bandwidth starvation.
- **Fast-Forward Merge Policy:** Pull operations enforce `--ff-only` by default. Diverged branches are flagged as dirty and presented in a summary table for manual resolution rather than triggering unexpected merge commits.

---

## 4. Pull-All Self-Healing & OMZ Ignore Strategy

- **Oh My Zsh (OMZ) Ignore Rule:** Automated pull scripts explicitly bypass third-party framework repositories like `.oh-my-zsh` and custom plugins unless targeted directly, preventing upstream update conflicts.
- **Push Self-Healing:** When remote reject errors occur due to non-fast-forward state, GitMap inspects commit ancestry, fetches remote tags, and offers safe rebase or fast-forward solutions.
- **Dirty Repo Handling:** Repositories with uncommitted local modifications are bypassed safely during batch pulls with clear status warnings in the terminal UI.
