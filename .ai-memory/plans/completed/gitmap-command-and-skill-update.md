# Execution Plan: GitMap Command and Skill Update

**Slug:** `gitmap-command-and-skill-update`  
**Task Name:** GitMap Command and Skill Update  
**Status:** DONE  
**Lifecycle Budget:** N = 400 Steps (Phase 1: 1-100, Phase 2: 101-250, Phase 3: 251-400)  
**Concurrency:** A = 3, H = 2, Capacity = 6  

---

## 1. User Request (Verbatim)

```text
python -c "import urllib.request, re;
for u in ['https://prnt.sc/FnF7IMTncnd5', 'https://prnt.sc/oxgHWyUQkiWW', 'https://prnt.sc/FdL_kq6hZEDW', 'https://prnt.sc/uXKknIPwaS3Z', 'https://prnt.sc/6LEiq8Qfq_bV']:
    req = urllib.request.Request(u, headers={'User-Agent': 'Mozilla/5.0'})
    try:
        h = urllib.request.urlopen(req, timeout=5).read().decode('utf-8')
        m = re.findall(r'https://image\.prntscr\.com/[^\"]+|https://api\.prntscr\.com/[^\"]+|https://i\.imgur\.com/[^\"]+', h)
        print(u, m[:2])
    except Exception as e:
        print(u, e)
"


# GitMap Command and Skill Update: high priority instruction, non-negotiable task

For this type of code, I think we already have done a `gitmap` command. Please confirm if you have it. If you don't, if you make it. So make that `gitmap` command and also circulate in the `gitmap` skills with forbidden this and the `gitmap` one so that we can follow through correctly and download and install following the `gitmap` skill. And you update those skill in the coding guideline, not here, but in the coding guideline. And also you sync up with all the repositories using coding guideline principles and coding guideline knowledge. Then make sure the some parts should be also mentioned in the execute prompt V6. So what they should use, a one-liner just to refer as a `gitmap` skill. Does this make sense? So you have to update the prompt in the coding guideline, you have to update the skill, and then from there, you have to update all the repositories, including the `gitmap`. Does this make sense?
```

---

## 2. Architecture & Design Blueprint

1. **CLI Enhancement (`cli/cmdasset/`):**
   - Confirm `gitmap prnt` / `gitmap download-prnt` exists.
   - Extend `cli/cmdasset/asset_prnt.go` to support:
     - Multi-URL batch input (`gitmap prnt <url1> <url2>...`).
     - `--extract-only` (`-e`, `-p`) flag: print direct URLs without downloading.
     - `--json` (`-j`) flag: output Universal JSON response envelope.
     - Enhanced regex for `api.prntscr.com`, `i.imgur.com`, `image.prntscr.com`, `img.lightshot.app`.
     - Concurrent worker pool downloading with error isolation.
2. **Skill Update (`.agents/skills/gitmap/SKILL.md`):**
   - Ban `python -c "import urllib.request, re..."` and ad-hoc scrapers in Section 2 Replacement Matrix.
   - Document `gitmap prnt` multi-URL syntax and `--extract-only`.
3. **Coding Guidelines Update (`d:/work/coding-guidelines`):**
   - Add Mistake #22 to `02-spec/02-coding-guidelines/06-ai-optimization/04-common-ai-mistakes.md`.
4. **V6 Prompt Update (`01-prompts/v4/14-execute/02-execute-parent-task-with-n-steps-v6.md`):**
   - Add Screenshot Asset row in Section 3 Tool Matrix.
   - Add Total Ban on raw Python scrapers with one-liner reference to `@[.agents/skills/gitmap]`.
5. **Fleet Synchronization:**
   - Synchronize skills, prompts, and coding guidelines across repositories.

---

## 3. Subtask Breakdown Matrix

| Subtask ID | Title | Disjoint Files Owned | Status |
|---|---|---|---|
| **Subtask-01** | Architecture Specs & `gitmap prnt` Batch Engine | `02-spec/21-app/gitmap-command-and-skill-update/`, `cli/cmdasset/` | DONE |
| **Subtask-02** | GitMap Skill & Coding Guidelines Update | `.agents/skills/gitmap/SKILL.md`, `02-spec/02-coding-guidelines/` | DONE |
| **Subtask-03** | Canonical V6 Execution Prompt & Upstream Guidelines | `d:/work/coding-guidelines/`, `01-prompts/` | DONE |
| **Subtask-04** | Unit Verification, Quality Gates & Release | `cli/cmdasset/*_test.go`, local tests | DONE |
