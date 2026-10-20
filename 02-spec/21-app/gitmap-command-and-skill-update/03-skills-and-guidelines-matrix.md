# Skills & Coding Guidelines Matrix: Screenshot Scraper Replacement

## 1. Non-Negotiable Command Replacement Matrix Addition

In `.agents/skills/gitmap/SKILL.md` (Section 2) and `02-spec/02-coding-guidelines/`:

| Operation | ❌ Prohibited Legacy / Ad-Hoc Python Command | ✅ Mandatory Native GitMap Command |
|---|---|---|
| **Screenshot Asset Download** | `python -c "import urllib.request..."` or ad-hoc scrapers | `gitmap prnt <url> [dest]` (alias `gitmap download-prnt`) |
| **Batch Screenshot Fetch** | `for u in list: urllib.request.urlopen(u)...` | `gitmap prnt <url1> <url2> <url3>...` |
| **Image Link Extraction** | `python -c "...re.findall(r'https://image...', h)..."` | `gitmap prnt <url1> <url2>... --extract-only` |
| **Structured Scraper Telemetry**| Custom JSON printing from Python | `gitmap prnt <url1> <url2>... --json` |

---

## 2. Common AI Mistakes Update

In `02-spec/02-coding-guidelines/06-ai-optimization/04-common-ai-mistakes.md`:

### Mistake #22: Ad-Hoc Python Web Scrapers for Screenshots
- **Symptom:** The AI writes a python one-liner using `urllib.request` or creates a temporary `fetch_urls.py` file to parse and extract image URLs from `prnt.sc` / Lightshot.
- **Root Cause:** Ignorance of GitMap's built-in `gitmap prnt` command family.
- **Remediation:**
  - Execute `gitmap prnt <url1> <url2>...` directly.
  - If only direct URLs are needed, pass `-e` or `--extract-only`.
  - Never write ad-hoc Python network scrapers inside repository workspaces.

---

## 3. Operational Guardrails

1. **Zero External Dependencies:** No Python, pip packages, curl pipelines, or puppeteer/playwright scripts for downloading Lightshot or screenshot URLs.
2. **Automatic Antigravity Synchronization:** All downloaded screenshots are dual-saved to `assets/screenshots/` and `.gemini/antigravity/brain/<active>/` for instant agent inspection.
3. **Graceful Error Isolation:** In batch mode, if one URL 404s or times out, sibling downloads must continue uninterrupted.
