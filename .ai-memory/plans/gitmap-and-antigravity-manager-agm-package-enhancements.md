# Master Plan: GitMap Single-Repo Pull UX Simplification, Progress Bar Fix, Shell Wrapping RCA, and AGM Integration

**Slug:** `gitmap-and-antigravity-manager-agm-package-enhancements`  
**Parent Task ID:** `task-20261018034151-gitmap-and-antigravity-manager-agm-package-enhancements`  
**Budget:** 400 steps  
**Target Version:** v6.572.0  
**Specification:** `02-spec/21-app/gitmap-and-antigravity-manager-agm-package-enhancements/01-architecture-spec.md`  

---

## User Request (Verbatim)

```markdown
# GitMap and `Antigravity Manager (AGM)` Package Enhancements: high priority instruction, non-negotiable task

For the `gitmap` pool, the single repo, it needs to be very simple. It should not show that much of a 20%, 30% progress bar. It should just say progressing and at the end, the previous UI. Very simple. It has updated from this commit to this. Probably this many lines has been changed. Everything is clean and active, things like that. If the state is dirty, we could mention the dirty. If it executed auto merge, then it would also mention auto merge is executed, things like that. Very simple UI. And also there is a little bit bug into the main test UI, which I do like if it could show multiple progress bar when you are downloading or pull all for the repos. I do think that it could get too much buggy, so try to keep the single progress bar which we have in the past. It is a bug I think. I would consider this as a bug. And also in some cases, the AI tried to run PowerShell command using `gitmap`, which it could have executed already using the `gitmap` directly. So why it does that? Can you find the root cause and try to fix it? Can you please start working on this as we have discussed?

Also, the recent changes in the `Antigravity Manager (`Antigravity Manager (AGM)`)` package, check that new commands and new functionalities with the prompt rerun and prompt review. Let's check if we have from our end, `gitmap`, AGY commands there to backup, restore, take the whole prompt with images. Images becomes a very big problem. You can suggest me some ways how I could create the image as a standalone single URL or things like that. What could be the easiest way to deal with it? You shared that thought and currently started fixing this. That should be enough for now. Also, please check the RC agent, which is similar to the `heatmap` RC command, which I requested from another AI here. Try to see if you could synchronize with that AI. If not, then skip that part. Try to do a git pull before you start working. Is it clear?

at the do minor bump and release

## task: GitMap and `Antigravity Manager (AGM)` Package Enhancements
## slug: gitmap-and-`Antigravity Manager (AGM)`-package-enhancements

# Actionable Items Must Follow Non-Negotiable

1. Git pull BEFORE starting any work — first try `gitmap pull <repo path>` (or run `gitmap pull` in repository folder); if that fails, fall back to regular `git pull` and resolve any merge conflicts so the repository is completely up to date and clean before starting.
2. Write spec under 02-spec/21-app/gitmap-and-`Antigravity Manager (AGM)`-package-enhancements/ and enqueue plan task in .ai-memory/plans/gitmap-and-`Antigravity Manager (AGM)`-package-enhancements.md (subtasks in .ai-memory/plans/subtasks/gitmap-and-`Antigravity Manager (AGM)`-package-enhancements/) first
3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
5. Simplify the `gitmap` pool UI to avoid showing detailed progress bars; instead, display a simple "progressing" message and summarize changes at the end.
6. Investigate and fix the bug in the main test UI related to multiple progress bars; ensure a single progress bar is maintained.
7. Identify and resolve the issue where AI runs PowerShell commands unnecessarily when `gitmap` commands suffice.
8. Review recent changes in the `Antigravity Manager (`Antigravity Manager (AGM)`)` package, focusing on new commands and functionalities.
9. Explore solutions for handling images as standalone URLs in the `Antigravity Manager (AGM)` package.
10. Check the RC agent for synchronization with the `heatmap` RC command; skip if synchronization is not possible.
11. Perform a minor bump and release after completing the tasks.
```

---

## Subtask Breakdown & Execution Roadmap

| Subtask ID | Title | File Box Ownership | Status |
|---|---|---|---|
| **Subtask-01** | GitMap Single-Repo Pull UI Simplification | `cli/cmdpull/pull_cwd.go`, `cli/cmdpull/pull_worker.go`, `cli/cloner/safe_pull.go`, `cli/cmdpull/pull_summary_single.go` | `QUEUED` |
| **Subtask-02** | Fix Progress Bar Column Width & Waterfall Wrapping Bug | `cli/cmdpull/pull_progress_bar_render.go`, `cli/cmdpull/pull_progress_bar_visual.go` | `QUEUED` |
| **Subtask-03** | Shell Wrapping RCA, Guidelines & Native Git Commands | `01-prompts/`, `.agents/skills/gitmap/SKILL.md`, `.ai-memory/strictly-avoid.md`, `cli/cmdgit/` | `QUEUED` |
| **Subtask-04** | AGM Package Updates Review & AGY Subsystem Parity | `cli/cmdagy/agy_rerun*.go`, `cli/cmdagy/agy_prompt_subcmds.go`, `cli/cmdagy/agy_running_prompts_backup.go` | `QUEUED` |
| **Subtask-05** | Standalone Image URL & Lightweight Prompt Architecture | `02-spec/21-app/gitmap-and-antigravity-manager-agm-package-enhancements/03-image-url-architecture.md`, `cli/cmdagy/agy_image_assets.go` | `QUEUED` |
| **Subtask-06** | RC Agent & Heatmap RC Continuous Synchronization Verification | `cli/heatmap/heatmap_rc_cmd.go`, `cli/heatmap/heatmap_conflict.go` | `QUEUED` |
| **Subtask-07** | Targeted Verification, Spec/Docs Update & Minor Release Ceremony | `version.json`, `changelog.md`, `readme.md`, `cli/version/version.go` | `QUEUED` |

---

## Concurrency Topology & Wave Allocations
- **Wave 1 (Core Engine Edits & Fixes)**:
  - Worker 01: Subtask-01 (Single-repo pull simplification) & Subtask-02 (Progress bar runewidth and wrap fix)
  - Worker 02: Subtask-03 (Shell wrapping prompt fixes & native git commands)
  - Worker 03: Subtask-04 (AGM & AGY alignment) & Subtask-05 (Image asset architecture)
- **Wave 2 (Verification, Heatmap RC Sync & Release)**:
  - Lead / Worker: Subtask-06 (Heatmap RC sync verification) & Subtask-07 (Targeted tests, version bump to `v6.572.0`, atomic push)
