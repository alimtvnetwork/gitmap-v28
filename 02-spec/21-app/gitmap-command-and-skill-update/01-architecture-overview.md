# Architecture Overview: GitMap Screenshot & Asset Downloader Engine

## 1. Problem Statement & Anti-Pattern Analysis
In AI-assisted software engineering, autonomous agents frequently encounter user bug reports containing screenshot links hosted on services such as Lightshot (`prnt.sc`), Imgur (`i.imgur.com`), or image CDNs (`image.prntscr.com`, `api.prntscr.com`).

A recurring anti-pattern observed in AI interactions is the execution of ad-hoc raw Python one-liners:
```python
# ❌ STRICTLY FORBIDDEN: Ad-hoc Python scraping script
python -c "import urllib.request, re;
for u in ['https://prnt.sc/...']:
    req = urllib.request.Request(u, headers={'User-Agent': 'Mozilla/5.0'})
    h = urllib.request.urlopen(req).read().decode('utf-8')
    m = re.findall(r'https://image\.prntscr\.com/[^\"]+', h)
    ...
"
```

### Why This Anti-Pattern is Prohibited:
1. **Security & Sandbox Escape Risks:** Ad-hoc scripts bypass repository policy checks, lack bounded timeouts, and risk certificate or socket leaks.
2. **Brittle HTML Scraping:** Hardcoded ad-hoc regexes fail on variations like OpenGraph `og:image`, redirect headers, or Lightshot layout changes.
3. **Missing Artifact Lifecycle:** Images downloaded into temporary or working trees litter git status, fail to synchronize with the Antigravity Brain, and lack structured JSON envelopes for downstream AI processing.
4. **Zero Cache & High Latency:** Sequential Python scraping does not take advantage of concurrent worker pools, connection reuse, or localized caching.

---

## 2. Native Solution: `gitmap prnt` Engine
GitMap provides native, compiled, CGO-free screenshot asset extraction and parallel downloading via `gitmap prnt` (aliases: `gitmap download-prnt`, `gitmap asset download-prnt`).

### Subsystem Architecture:
```
User/AI Input:
  gitmap prnt <url1> <url2>... [--extract-only] [--json] [-o <dir>]
                   │
                   ▼
  Target Parser & URL Normalizer (cli/cmdasset/)
  ├── Validates URLs (prnt.sc, imgur, lightshot direct, CDN)
  ├── Detects direct vs indirect image targets
  ├── Configures flags: ExtractOnly, JsonFormat, Concurrency
  │
  ├── Concurrency Dispatcher (Worker Pool: goroutines + sync.WaitGroup)
  │   ├── Headless HTTP Client (User-Agent emulation, 15s timeout)
  │   ├── Multi-Tier Regex Parser:
  │   │   ├── Tier 1: Direct CDN regex (image.prntscr.com, i.imgur.com)
  │   │   ├── Tier 2: OpenGraph og:image meta tag
  │   │   ├── Tier 3: DOM image elements (#screenshot-image, .screenshot__image)
  │   │   └── Tier 4: Fallback URL matcher
  │   └── Parallel Image Fetcher (io.ReadAll with 10MB safety cap)
  │
  ├── Multi-Destination Writer:
  │   ├── Local Repository Asset: assets/screenshots/<slug>_<timestamp>.png
  │   └── Antigravity Brain Sync: ~/.gemini/antigravity/brain/<active>/screenshot_<slug>.png
  │
  └── Output Formatter:
      ├── Terminal Table: Formatted colored summary table with direct URLs
      └── Universal JSON Envelope: { Status, Attributes, Results }
```

---

## 3. Key Architecture Guarantees
- **O(1) Dependencies:** Pure Go standard library (`net/http`, `regexp`, `os`, `sync`). Zero external Python or third-party scraper dependencies.
- **Concurrent Batch Fetching:** Proactively resolves and downloads multiple URLs in parallel using goroutines.
- **Antigravity Brain Interoperability:** Automatically detects active Antigravity session directories and duplicates screenshot artifacts into the brain workspace so agents can inspect them immediately.
- **Non-Destructive Extraction:** Supports `--extract-only` to display resolved direct URLs without writing disk files.
