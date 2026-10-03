# Git Reconciliation & Mechanical Conflict Resolution — Workflow (must follow)

> **Prompt Version:** 1.0.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

[/goal](slashCommand;goal) Safely reconcile diverged Git branches, mechanically resolve merge conflicts across known file patterns, and push cleanly to the remote tracking branch without rewriting published Git history.

---

## Architectural Motivation & Why This Is Very Important

When autonomous agents, multiple developers, or CI/CD automated release bots commit concurrently, Git branches diverge. In environments connected to platforms like Lovable or live deployment pipelines, **rewriting published git history (force-pushing, rebasing published commits, or squashing pushed history) is strictly forbidden** because it corrupts external state, orphans client sessions, and invalidates build history.

To achieve reliable, deterministic, and autonomous synchronization:
1. **Divergence Must Be Non-Destructively Reconciled:** Always use standard merge semantics (`git merge --no-ff`) instead of rebasing or force-pushing.
2. **Conflict Resolution Must Be Mechanical & Deterministic:** Instead of guessing or hallucinating arbitrary changes, conflicts in structured files (such as `.gitignore`, `version.json`, `changelog.md`, dependency manifests, markdown task checklists, and code import blocks) follow rigorous mathematical and structural rules:
   - Sets and ignore files obey **deduplicated union**.
   - SemVer manifests obey **monotonic version dominance** (highest version wins).
   - Changelogs and release notes obey **chronological section preservation**.
   - Checklists obey **completion dominance** (`- [x]` overrides `- [ ]`).
   - Declarative imports obey **lexicographical union**.
3. **Pre-Commit Syntax & Zero-Marker Invariants:** Every resolved file MUST be strictly audited to ensure **zero conflict markers** (`<<<<<<<`, `=======`, `>>>>>>>`) remain, and syntax validation (e.g. AST parsing for Python, JSON schema verification) MUST pass before staging.

---

## Step-by-Step Instructions for Mechanical Conflict Resolution

Follow these structured steps sequentially to execute or mechanically automate merge conflict reconciliation:

### Step 1: Remote Tracking Fetch & Working Tree Safety Guard
- Run `git fetch <remote>` to retrieve the latest remote tracking refs without modifying the working tree.
- Inspect `git status --porcelain` to identify untracked, modified, or conflicted files.
- **Safety Invariant:** If untracked or uncommitted changes conflict with incoming files, stash or stage them before proceeding to prevent accidental data loss.

### Step 2: Divergence & Commit Delta Analysis
- Determine current branch and remote counterpart: `git rev-parse --abbrev-ref HEAD` and `<remote>/<branch>`.
- Count ahead and behind commits:
  - `git rev-list --count <remote>/<branch>..HEAD` (ahead count)
  - `git rev-list --count HEAD..<remote>/<branch>` (behind count)
  - `git merge-base HEAD <remote>/<branch>` (common ancestor commit)
- Branch State Decision Tree:
  - If ahead == 0 and behind == 0: Branch is already synchronized. Proceed directly to verification.
  - If ahead == 0 and behind > 0: Fast-forward cleanly using `git merge --ff-only <remote>/<branch>`.
  - If ahead > 0 and behind == 0: Local commits pending; proceed to push.
  - If ahead > 0 and behind > 0: **Diverged branches detected**. Standard non-destructive merge required.

### Step 3: Non-Destructive Merge Execution
- Initiate standard non-fast-forward merge:
  ```bash
  git merge --no-ff <remote>/<branch>
  ```
- If the merge exits with status 0, the merge resolved cleanly without conflict. Proceed directly to Step 6.
- If the merge exits with a non-zero code, Git has flagged one or more conflicted files. Proceed immediately to Step 4.

### Step 4: Mechanical Conflict Resolution by Domain Handlers
Query conflicted files using `git diff --name-only --diff-filter=U`. For each conflicted file, route through the appropriate deterministic solver:

1. **Ignore & Set Files (`.gitignore`, `.gitattributes`, `.npmignore`, `.dockerignore`):**
   - Extract `ours` and `theirs` blocks between `<<<<<<<`, `=======`, and `>>>>>>>`.
   - Compute the deduplicated line union, preserving comments and section headers.
2. **Version Manifests (`version.json`, `package.json`, `prompt-version.template.json`):**
   - Extract version strings from `ours` and `theirs`.
   - Parse SemVer tuples `(major, minor, patch)`.
   - Adopt the block with the higher SemVer constraint to prevent version rollbacks.
3. **Changelogs & Release Notes (`changelog.md`, `release-notes.md`):**
   - Preserve release headers and entries from both branches chronologically.
   - Separate distinct release entries cleanly with blank lines.
4. **Markdown Lists & Checklists (`.md`, `.markdown`):**
   - Parse list items (`- `, `* `, `+ `, `\d+\. `).
   - If checklist items (`- [ ]` vs `- [x]`), adopt completion state (`- [x]`) for identical item text.
   - For markdown tables (`| col | col |`), deduplicate table rows and preserve single divider row.
5. **Dependency Manifests (`requirements*.txt`, `*-requirements.txt`):**
   - Union distinct package names.
   - For duplicate packages with differing version constraints, adopt the higher SemVer constraint.
6. **Code Import Blocks (`.py`, `.ts`, `.js`, `.go`, `.rs`):**
   - If the conflict block consists purely of declarative import/use statements, union all distinct import lines and deduplicate.
7. **Generic Fallback (`union`, `ours`, or `theirs`):**
   - For general prose or non-colliding additions, compute clean union.

### Step 5: Post-Resolution Marker Audit & Syntax Validation
- **Marker Audit:** Scan file contents for `<<<<<<<`, `=======`, or `>>>>>>>`. If any marker exists, resolution has failed; abort commit.
- **Syntax Validation:**
  - For Python files (`.py`), run `ast.parse()` to ensure valid Python syntax.
  - For JSON files (`.json`), run `json.loads()` to ensure valid JSON structure.
- **Stage Resolved File:** Run `git add <file>` once verified.

### Step 6: Merge Commit Creation
- Once all conflicted files are resolved and staged:
  ```bash
  git commit -m "Merge remote-tracking branch '<remote>/<branch>' into <branch> (mechanically resolved conflicts)"
  ```

### Step 7: Push & Final Synchronization Verification
- Push synchronized branch to remote tracking repository:
  ```bash
  git push <remote> <branch>
  ```
- Run post-flight verification: `git rev-list --count <remote>/<branch>..HEAD` must equal 0.
- Verify working tree is clean: `git status --porcelain` must be empty.

---

## Automated Execution via Python AI Script

Use the canonical Python automation script in `03-ai-scripts/47-git-reconcile-and-resolve-conflict.py` for one-command reconciliation, conflict resolution, and self-testing:

```bash
# 1. Run self-tests on synthetic conflict test cases
python 03-ai-scripts/47-git-reconcile-and-resolve-conflict.py --self-test

# 2. Check divergence without modifying working tree
python 03-ai-scripts/47-git-reconcile-and-resolve-conflict.py --check

# 3. Reconcile, mechanically resolve conflicts, and push to origin
python 03-ai-scripts/47-git-reconcile-and-resolve-conflict.py --push

# 4. Mechanically resolve active conflict markers in working tree
python 03-ai-scripts/47-git-reconcile-and-resolve-conflict.py resolve-conflicts --strategy smart
```

---

## Action Items — Must Follow (Non-Negotiable)

- [ ] Fetch remote tracking refs and inspect ahead/behind divergence counts before modifying working tree.
- [ ] Strictly avoid rebasing, force-pushing, or squashing published Git history (preserving Lovable and CI/CD integrity).
- [ ] Execute non-destructive standard merge (`git merge --no-ff`) for diverged branches.
- [ ] Route conflicted files through domain-specific mechanical solvers (gitignore, SemVer JSON, changelogs, markdown checklists, requirements, imports).
- [ ] Audit every resolved file for zero remaining conflict markers (`<<<<<<<`, `=======`, `>>>>>>>`).
- [ ] Validate AST and syntax integrity for all modified code and JSON files prior to staging.
- [ ] Stage verified files with `git add` and complete the structured merge commit.
- [ ] Push to remote tracking branch and verify zero ahead/behind delta post-push.

---

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

## Anti-Hallucination, Micro-Tasking, & Self-Looping

> [!CAUTION]
> **CRITICAL RULE: DO NOT ATTEMPT TO READ, PLAN, AND EXECUTE EVERYTHING AT ONCE.**
> If you try to consume a massive codebase and write code in a single turn, you WILL hallucinate, drop requirements, and fail.

To survive massive checklists and complex codebases, you MUST operate using these three principles:

1. **Phase 1: Read & Understand (Isolated Loop):** Your very first action must be purely exploratory. Do NOT write code. Break down the task, read the specific files, trace the dependencies, and understand the architectural boundary. Once you understand the scope, end your turn and self-loop to begin execution.
2. **Phase 2: Bounded Micro-Tasking (Sequential Self-Looping):** Never attempt to execute the entire checklist in one response. Treat each checklist section or file as a strict, isolated boundary. Execute *only* the first small portion, verify it, end your turn, and self-loop to process the next portion.
3. **Phase 3: Multi-Agent Parallelization:** If tasks are independent, you MUST spawn dedicated sub-agents to handle them concurrently. Give each sub-agent an extremely small, strictly defined bounding box (e.g., "Only edit File X"). Never give a sub-agent a generic or multi-file task.
