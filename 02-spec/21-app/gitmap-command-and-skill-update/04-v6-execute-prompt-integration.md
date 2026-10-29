# V6 Execution Prompt Integration & Canonical Reference

## 1. Executive Reference Requirement
The user prompt mandates:
> "Then make sure the some parts should be also mentioned in the execute prompt V6. So what they should use, a one-liner just to refer as a `gitmap` skill. Does this make sense? So you have to update the prompt in the coding guideline, you have to update the skill, and then from there, you have to update all the repositories, including the `gitmap`."

---

## 2. One-Liner Reference Specification

In `01-prompts/v4/14-execute/02-execute-parent-task-with-n-steps-v6.md` and `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`:

### Section 3: Non-Negotiable Tool Replacement Table
Add row:
```markdown
| **Screenshot Asset & URL Extraction** | `python -c "import urllib.request..."` or ad-hoc scrapers | `gitmap prnt <url1> [url2...] [--extract-only] [--json]` (or `gitmap download-prnt`) |
```

### TOTAL BAN ON AD-HOC PYTHON NETWORK SCRAPERS
In Section 3 (TOTAL BAN rules):
```markdown
- **TOTAL BAN ON RAW PYTHON SCRAPERS & URL FETCHER SCRIPTS:** NEVER execute `python -c "import urllib.request..."` or write ad-hoc Python scripts (`fetch_urls.py`, etc.) to scrape web pages, extract images, or download screenshot URLs (e.g., `prnt.sc`, `image.prntscr.com`, `i.imgur.com`). You MUST use native compiled `gitmap prnt <url1> [url2...] [--extract-only] [--json]` (alias `gitmap download-prnt`), which automatically handles User-Agent headers, multi-tier CDN resolution, and bidirectional Antigravity Brain synchronization. See `@[\1](file;.agents/skills/\1)` for full command reference.
```

---

## 3. Upstream & Cross-Repository Sync Strategy
1. Apply changes to canonical prompt in canonical coding-guidelines repository (`01-prompts/v4/14-execute/02-execute-parent-task-with-n-steps-v6.md`).
2. Apply changes to canonical skill in canonical coding-guidelines repository (`.agents/skills/gitmap/SKILL.md`).
3. Apply changes to local repository root.
4. Run `03-ai-scripts/38-sync-prompts-skills-scripts.py` to synchronize across connected fleet repositories.
