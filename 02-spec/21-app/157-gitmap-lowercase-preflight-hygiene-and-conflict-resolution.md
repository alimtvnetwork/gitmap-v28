# Spec 156: GitMap Lowercase Pre-Flight Hygiene, Working Tree Conflict Resolution & Automated Push

> **Spec ID:** `SPEC-156`  
> **Version:** `v6.338.0`  
> **Status:** Active  
> **Date:** 2026-09-25  

---

## 0. User Request (Verbatim)

```text
https://prnt.sc/n_O_kOGNEJqZ


PS D:\work\presentations-repos\hiltrax> gitmap lowercase "*.md"

⚡ GitMap Lowercase File Renamer
  ● Mode:        Git Repository (2-step git mv)
  ● Working Dir: D:\work\presentations-repos\hiltrax
  ● Filter:      *.md
  ● Scanned:     1043 files in directory
  ● Matched:     18 uppercase file(s)


════════════════════════════════════════════════════════════════
⚡ Pre-Flight Verification: Git Lowercase Renamer
════════════════════════════════════════════════════════════════
  ● Found: 18 uppercase file(s) to rename:
    • SKILL.md -> skill.md
    • SKILL.md -> skill.md
    • SKILL.md -> skill.md
    • SKILL.md -> skill.md
    • AGENTS.md -> agents.md
    • LLM.md -> llm.md
    • AdvanceStepSlide.md -> advancestepslide.md
    • CapsuleListSlide.md -> capsulelistslide.md
    ... and 10 more file(s)

  ● Planned Change: Safe 2-step atomic git mv
    1. Move to temp:   git mv <file> <file>.tmp-lcf (breaks case-collision)
    2. Move to target: git mv <file>.tmp-lcf <file_lowercase>

  ● After-Effects:
    • Working tree and Git index synchronized (git add -A)
    • Automatically committed to branch: main
────────────────────────────────────────────────────────────────
Type 'confirm' or 'yes' (or 'y') to proceed: y
  [1/18] SKILL.md -> skill.md (filter [*.md])
        Step 1: git mv SKILL.md -> SKILL.md.tmp-lcf (OK)
        Step 2: git mv SKILL.md.tmp-lcf -> skill.md (OK)
  [2/18] SKILL.md -> skill.md (filter [*.md])
        Step 1: git mv SKILL.md -> SKILL.md.tmp-lcf (OK)
        Step 2: git mv SKILL.md.tmp-lcf -> skill.md (OK)
  [3/18] SKILL.md -> skill.md (filter [*.md])
        Step 1: git mv SKILL.md -> SKILL.md.tmp-lcf (OK)
        Step 2: git mv SKILL.md.tmp-lcf -> skill.md (OK)
  [4/18] SKILL.md -> skill.md (filter [*.md])
        Step 1: git mv SKILL.md -> SKILL.md.tmp-lcf (OK)
        Step 2: git mv SKILL.md.tmp-lcf -> skill.md (OK)
  [5/18] AGENTS.md -> agents.md (filter [*.md])
        Step 1: git mv AGENTS.md -> AGENTS.md.tmp-lcf (OK)
        Step 2: git mv AGENTS.md.tmp-lcf -> agents.md (OK)
  [6/18] LLM.md -> llm.md (filter [*.md])
        Step 1: git mv LLM.md -> LLM.md.tmp-lcf (OK)
        Step 2: git mv LLM.md.tmp-lcf -> llm.md (OK)
  [7/18] AdvanceStepSlide.md -> advancestepslide.md (filter [*.md])
        Step 1: git mv AdvanceStepSlide.md -> AdvanceStepSlide.md.tmp-lcf (OK)
        Step 2: git mv AdvanceStepSlide.md.tmp-lcf -> advancestepslide.md (OK)
  [8/18] CapsuleListSlide.md -> capsulelistslide.md (filter [*.md])
        Step 1: git mv CapsuleListSlide.md -> CapsuleListSlide.md.tmp-lcf (OK)
        Step 2: git mv CapsuleListSlide.md.tmp-lcf -> capsulelistslide.md (OK)
  [9/18] ERDiagramSlide.md -> erdiagramslide.md (filter [*.md])
        Step 1: git mv ERDiagramSlide.md -> ERDiagramSlide.md.tmp-lcf (OK)
        Step 2: git mv ERDiagramSlide.md.tmp-lcf -> erdiagramslide.md (OK)
  [10/18] FocusTimelineSlide.md -> focustimelineslide.md (filter [*.md])
        Step 1: git mv FocusTimelineSlide.md -> FocusTimelineSlide.md.tmp-lcf (OK)
        Step 2: git mv FocusTimelineSlide.md.tmp-lcf -> focustimelineslide.md (OK)
  [11/18] ImageSlide.md -> imageslide.md (filter [*.md])
        Step 1: git mv ImageSlide.md -> ImageSlide.md.tmp-lcf (OK)
        Step 2: git mv ImageSlide.md.tmp-lcf -> imageslide.md (OK)
  [12/18] KeywordSlide.md -> keywordslide.md (filter [*.md])
        Step 1: git mv KeywordSlide.md -> KeywordSlide.md.tmp-lcf (OK)
        Step 2: git mv KeywordSlide.md.tmp-lcf -> keywordslide.md (OK)
  [13/18] MetricGridSlide.md -> metricgridslide.md (filter [*.md])
        Step 1: git mv MetricGridSlide.md -> MetricGridSlide.md.tmp-lcf (OK)
        Step 2: git mv MetricGridSlide.md.tmp-lcf -> metricgridslide.md (OK)
  [14/18] MiddleTitleSlide.md -> middletitleslide.md (filter [*.md])
        Step 1: git mv MiddleTitleSlide.md -> MiddleTitleSlide.md.tmp-lcf (OK)
        Step 2: git mv MiddleTitleSlide.md.tmp-lcf -> middletitleslide.md (OK)
  [15/18] QrMeetingSlide.md -> qrmeetingslide.md (filter [*.md])
        Step 1: git mv QrMeetingSlide.md -> QrMeetingSlide.md.tmp-lcf (OK)
        Step 2: git mv QrMeetingSlide.md.tmp-lcf -> qrmeetingslide.md (OK)
  [16/18] SectionDividerSlide.md -> sectiondividerslide.md (filter [*.md])
        Step 1: git mv SectionDividerSlide.md -> SectionDividerSlide.md.tmp-lcf (OK)
        Step 2: git mv SectionDividerSlide.md.tmp-lcf -> sectiondividerslide.md (OK)
  [17/18] StepTimelineSlide.md -> steptimelineslide.md (filter [*.md])
        Step 1: git mv StepTimelineSlide.md -> StepTimelineSlide.md.tmp-lcf (OK)
        Step 2: git mv StepTimelineSlide.md.tmp-lcf -> steptimelineslide.md (OK)
  [18/18] TitleSlide.md -> titleslide.md (filter [*.md])
        Step 1: git mv TitleSlide.md -> TitleSlide.md.tmp-lcf (OK)
        Step 2: git mv TitleSlide.md.tmp-lcf -> titleslide.md (OK)

✔ Committed 18 lowercase file rename(s): "chore: rename 18 files to lowercase across repository" (9582c388)


════════════════════════════════════════════════════════════════
✔ Lowercase Rename Summary:
  ● Total Files Scanned: 1043
  ● Files Matched:       18
  ● Files Renamed:       18
  ● Git Status:          Committed (9582c388)

  ● Manipulation Steps Performed:
    1. Step 1 (Safe Temp Move):  git mv <file> <file>.tmp-lcf
       Avoids silent case-collision / no-op on case-insensitive filesystems (Windows/macOS)
    2. Step 2 (Target Rename):   git mv <file>.tmp-lcf <file_lowercase>
       Registers true case rename in Git index tree
    3. Step 3 (Index Sync):      git add -A
    4. Step 4 (Atomic Commit):   git commit -m "chore: rename ..."
════════════════════════════════════════════════════════════════

PS D:\work\presentations-repos\hiltrax>


Okay. So here the big problem is that even if you fix those uppercase to lowercase, the first thing you created is the Git conflict, and you didn't resolve the Git conflict. Okay, and it looks like that there are too many Git merge happen, conflict resolved. Why? Because you are doing this change, you should make sure that it is done properly. Okay, if there is something pending, you make sure that you discard it, you don't care it, and you confirm with the user. You list out all the files that you are dealing with, make sure there is no hidden stuff. And then when you start, you make sure you commit and resolve and push to the Git, which is missing from your task. Is it understood? Can you please fix it and bump the minor version and release it
```

---

## 1. Visual Evidence & Reproduction

![GitMap Lowercase Conflict Evidence](assets/screenshots/gitmap-lowercase-conflict-resolution-01.png)

### Key Observations from User Run:
1. `gitmap lowercase "*.md"` found 18 uppercase files.
2. In the pre-flight verification summary, only 8 files were listed, followed by `... and 10 more file(s)`, obscuring the full scope of modified files.
3. The command proceeded without checking whether the repository had pending uncommitted changes, staged changes, or unresolved merge conflicts.
4. When `git add -A` was run, it swept up any pre-existing dirty working tree state into the rename commit, resulting in unexpected merge conflicts when pulling/syncing.
5. The command committed locally (`9582c388`) but never pushed to the remote repository (`git push`), leaving the local and remote branches desynchronized.
6. The user explicitly requests:
   - Full disclosure: list out all files being dealt with (zero hidden files).
   - Pre-flight working tree hygiene: if something is pending, disclose it, allow discarding with confirmation, or clean conflicts.
   - Complete lifecycle: commit, verify clean merge state, and push to Git.
   - Minor version bump and release.

---

## 2. Technical Architecture & Protocols

### 2.1 Complete File Disclosure Protocol (No Hidden Files)
- The pre-flight box must list **every single matched file** without truncation or ellipsis (`... and N more files`), unless explicitly paginated or configured with a high ceiling.
- Each file entry displays the full relative path from repository root (e.g. `docs/AdvanceStepSlide.md → docs/advancestepslide.md`) rather than just the base filename, so users see exact locations.

### 2.2 Working Tree Pre-Flight Hygiene & Conflict Audit
Before executing any file moves:
1. Run `git status --porcelain` to detect:
   - Unmerged/conflicted files (`UU`, `AA`, `UD`, `DU`, `DD`, `AU`, `UA`).
   - Staged changes (`M `, `A `, `D `, `R `).
   - Unstaged modified files (` M`, ` D`).
   - Untracked files (`??`).
2. If unresolved merge conflicts exist:
   - Abort immediately with an explicit error detailing the conflicted files.
   - Do NOT attempt to run `git mv` or `git add -A` over a conflicted index.
3. If uncommitted dirty/pending changes exist:
   - Output a dedicated `⚠️ Warning: Working tree has pending uncommitted changes:` section listing all dirty files.
   - Offer the user an explicit option to discard pending changes (`--discard-pending` or prompt: `Discard pending changes to proceed with clean renames? [y/N]`).
   - If confirmed, execute `git reset --hard HEAD` and `git clean -fd` (or discard pending files) to establish a clean slate before renaming.
   - If not confirmed and `--force` is not passed, cancel safely without touching any files.

### 2.3 Automated Git Push Synchronization
1. After staging (`git add -A`) and committing (`git commit -m ...`), `gitmap lowercase` must automatically execute `git push` to the remote branch.
2. Add `--no-push` flag for users who explicitly wish to keep commits local.
3. If `git push` encounters upstream diverged commits:
   - Pull with rebase or report upstream status clearly so the user is never left with silent merge conflicts.

---

## 3. Verification & Acceptance Criteria

1. **AC-SPEC156-01 (Full Disclosure):** All matched files are listed individually with their full relative paths in the pre-flight verification box.
2. **AC-SPEC156-02 (Pre-Flight Conflict Gate):** If the repository has unresolved merge conflicts, the renamer halts with a descriptive error.
3. **AC-SPEC156-03 (Pending Discard Support):** If uncommitted changes exist, the renamer surfaces them and provides safe discard confirmation.
4. **AC-SPEC156-04 (Automated Push):** The command commits and pushes to the current tracking branch unless `--no-push` or `--no-commit` is specified.
5. **AC-SPEC156-05 (Minor Version Bump):** Minor version is bumped and release orchestrated as requested.
