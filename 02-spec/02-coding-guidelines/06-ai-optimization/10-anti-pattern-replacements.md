# Anti-Pattern Replacements (AI Execution Prompt)

> **/goal** Eliminate obsolete shell commands, raw HTTP invocations, manual filesystem grepping, single-stream download bottlenecks, and absolute paths by substituting them with native, ultra-low-latency GitMap commands.
> **/learn** Enforce the canonical command substitution matrix: use `gitmap url read` for web assets, 10-split segmented downloading for large payloads, `gitmap heatmap` for past decision discovery, `gitmap aum search` for file scanning, and strictly relative workspace paths.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Replace all raw HTTP fetchers (`curl`, `wget`, `Invoke-WebRequest`, Python urllib) with `gitmap url read <url>`.
- [ ] `/learn` Accelerate large payload downloads (>= 1–2 MB) using GitMap's automatic 10-split segmented downloader and `aria2c` integration.
- [ ] `/goal` Inspect past agent decisions, file touch history, and active leases using `gitmap heatmap` or `gitmap rc heatmap --repo <alias>`.
- [ ] `/learn` Strictly ban `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, and `findstr` in favor of `gitmap aum search`, `gitmap find`, and `gitmap cat`.
- [ ] `/goal` Prohibit absolute filesystem paths (`C:\`, `D:\`, `file:///`) across all scripts, specs, and documentation, enforcing relative paths from repository root.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 1.0.0
**Updated:** 2026-10-23
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## 1. Executive Summary & Purpose

Autonomous AI agents operating within GitMap repositories frequently fall back to legacy shell utilities, ad-hoc network scripts, and brute-force filesystem scans. In multi-agent swarm environments, these anti-patterns cause severe friction:
1. **Network Redundancy & Latency:** Repeatedly invoking `curl` or `Invoke-WebRequest` across 10+ subagents burns external rate limits and introduces seconds of blocking network I/O.
2. **Download Bottlenecks:** Single-stream downloads for files $\ge 1\text{--}2\text{ MB}$ stall execution pipelines.
3. **Telemetry Blind Spots:** Grepping source trees or guessing past agent edits ignores existing cognitive traces in SQLite heatmap databases.
4. **Context Window Exhaustion:** Unfiltered shell output from `grep` or `Select-String` floods agent context windows with megabytes of unstructured text.
5. **Portability Failures:** Hardcoded drive letters and absolute paths break cross-platform portability across Windows, Linux, and macOS environments.

This specification mandates 5 primary anti-pattern replacements with native GitMap CLI tools.

---

## 2. Anti-Pattern Substitution Matrix

| Prohibited Anti-Pattern | Mandatory GitMap Command | Core Rationale & Technical Advantage |
| :--- | :--- | :--- |
| ❌ Raw `curl`, `wget`, `Invoke-WebRequest`, custom Python/Go HTTP fetchers | ✅ `gitmap url read <url>`<br>*(aliases: `gitmap url`, `gitmap get`)* | Dedicated SQLite Split-DB (`repodb/url_cache.db`), WAL mode, sub-millisecond retrieval (<0.5ms), SHA-256 caching, unbuffered `output.Raw()` streaming. |
| ❌ Single-stream downloads for files $\ge 1\text{--}2\text{ MB}$ | ✅ GitMap 10-split segmented engine<br>*(Host `aria2c` or native Go Range workers)* | Automatic threshold detection; 10 parallel connections with 5KB chunks saturate bandwidth and prevent socket timeouts. |
| ❌ Grepping filesystem or guessing past agent decisions | ✅ `gitmap heatmap`<br>`gitmap rc heatmap --repo <alias>` | 3-tier SQLite discovery, alias resolution via `db.ListAliasesWithRepo()`, recent tasks across 2–5 repos, dynamic suggestion pills. |
| ❌ `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, `findstr` | ✅ `gitmap aum search "<pattern>" [dir]`<br>`gitmap find "<pattern>" [dir]`<br>`gitmap cat <filepath>` | Multi-core streaming live search (<15ms), binary null-byte probe, 500 KB file cap, DH2D SQLite hot cache, zero buffer bloat. |
| ❌ Raw `$env:LIBCLANG_PATH` exports & ad-hoc Cargo scripts | ✅ `gitmap cargo [args...]`<br>`gitmap rust status` | Automated LLVM/libclang path discovery on Windows & Linux, injecting `LIBCLANG_PATH` & `PATH` transparently. |
| ❌ Raw `taskkill /F /PID <pid>` or `kill -9 <pid>` | ✅ `gitmap kill <pid\|name\|:port>`<br>*(alias: `gitmap terminate`)* | Cross-platform process termination with safety guards protecting system-critical processes, IDEs, and shells. |
| ❌ Third brackets or lower-case commits (`feat([module]): ...`) | ✅ `gitmap cpf "(module): summary"`<br>`gitmap cbf "(module): summary"`<br>`gitmap cfr "(module): summary"` | Enforces uppercase `TYPE(MODULE): summary` with parentheses; square brackets `[...]` are strictly prohibited. |
| ❌ Manual per-folder `~/.ssh/config` host aliases | ✅ `gitmap wg add <name> <dir> --ssh-key <key>`<br>`gitmap wg bind <dir> <key>` | SQLite-backed folder & ancestor SSH key inheritance for zero-prompt remote Git operations. |
| ❌ Absolute filesystem paths (`C:\`, `D:\`, `file:///`) | ✅ Relative Git-root paths<br>*(e.g. `02-spec/...`, `cli/...`)* | Ensures zero path leakage across host machines, container boundaries, and multi-node SSH cluster fleets. |

---

## 3. Deep-Dive Anti-Pattern Specifications

### Anti-Pattern 1: Raw Web Fetchers vs `gitmap url read`

#### Context & Forbidden Pattern
AI agents frequently invoke raw shell commands or author inline scripts to fetch schemas, documentation, release binaries, or test payloads. These network calls bypass local caching and risk terminal stream corruption.

```bash
<!-- ❌ FORBIDDEN: Raw shell fetchers and ad-hoc scripts -->
curl -s https://example.com/schema.json > schema.json
wget -qO- https://example.com/data.txt
Invoke-WebRequest -Uri "https://example.com/api" -UseBasicParsing
python -c "import urllib.request; print(urllib.request.urlopen('https://example.com').read())"
```

#### Required Replacement
All remote HTTP/HTTPS resources must be read via `gitmap url read <url>` (or shorthand aliases `gitmap url <url>` and `gitmap get <url>`).

```bash
<!-- ✅ REQUIRED: GitMap URL split-db cached streaming -->
# Byte-faithful stream directly to stdout (<0.5ms on cache hit)
gitmap url read https://example.com/schema.json

# Shorthand aliases for agent ergonomics
gitmap url https://example.com/schema.json
gitmap get https://example.com/schema.json

# Force fresh network refresh when upstream is known to have updated
gitmap url read https://example.com/schema.json --fresh

# Inspect cached payload metadata, ETag, and TTL
gitmap url info https://example.com/schema.json
```

#### Technical Guarantees
- **Split-DB Isolation:** Stored in `repodb/url_cache.db` with WAL mode and `PRAGMA mmap_size = 268435456` (256MB memory mapping).
- **Sub-Millisecond Retrieval:** Repeated reads return in $<0.5\text{ms}$ without touching the network.
- **Byte-Faithful Streaming:** Emits payloads via `output.Raw()` so binary, JSON, and text remain bit-identical. Diagnostic telemetry is isolated to `output.UIErr()`.
- **Request Coalescing:** Concurrent agent requests for the same URL share a single in-flight fetch via `singleflight.Group`.

---

### Anti-Pattern 2: Single-Stream Large Downloads vs 10-Split Segmented Engine

#### Context & Forbidden Pattern
Downloading large assets ($\ge 1\text{--}2\text{ MB}$), archives, language runtimes, or database seeds using single-stream HTTP connections leads to slow transfers, dropped connections, and pipeline timeouts.

```bash
<!-- ❌ FORBIDDEN: Single-stream downloads for files >= 1-2 MB -->
curl -O https://example.com/large-archive.tar.gz
Invoke-WebRequest -Uri "https://example.com/large-seed.db" -OutFile "seed.db"
```

#### Required Replacement
GitMap automatically activates its 10-split segmented download pipeline whenever target payloads meet or exceed the 1–2 MB threshold and the server signals `Accept-Ranges: bytes`.

```bash
<!-- ✅ REQUIRED: 10-split segmented downloading -->
# GitMap automatically probes HEAD Content-Length and selects Tier 1 (aria2c) or Tier 2 (Go Range workers)
gitmap url read https://example.com/large-archive.tar.gz > large-archive.tar.gz

# In Go source code, use the downloader engine directly:
res, err := downloader.ExecuteSegmentedDownload(ctx, downloader.DownloadOptions{
    URL:             "https://example.com/large-archive.tar.gz",
    DestinationPath: "large-archive.tar.gz",
    ChunkCount:      10,
    MinChunkSize:    5120, // 5KB minimum chunk
})
```

#### Architecture Tiers
1. **Tier 1 (Host aria2c):** If `aria2c` is present in system `PATH`, GitMap delegates to `aria2c -s 10 -x 10 -k 5K --file-allocation=none` for maximum hardware throughput.
2. **Tier 2 (Native Go 10-Split Range Downloader):** If `aria2c` is absent, GitMap creates pre-allocated destination storage and launches 10 concurrent goroutines writing into target byte offsets using `io.NewOffsetWriter`.
3. **Tier 3 (Single-Stream Fallback):** For servers omitting Range headers or payloads $< 1\text{ MB}$, falls back cleanly to streaming HTTP.

---

### Anti-Pattern 3: Filesystem Grepping for Past Decisions vs `gitmap heatmap`

#### Context & Forbidden Pattern
When newly joined or dynamically respawned subagents need to understand what previous agents did, they often resort to grepping source files or guessing which files were touched.

```bash
<!-- ❌ FORBIDDEN: Grepping filesystem or guessing past agent actions -->
git log --stat -n 5
grep -rn "TODO: task 314" .
git diff HEAD~3
```

#### Required Replacement
Agents must query GitMap's cognitive heatmap discovery catalog to inspect past decisions, touched files, reasoning states, and active file claims.

```bash
<!-- ✅ REQUIRED: Heatmap discovery and task inspection -->
# Discover all heatmap databases, alias mappings, recent tasks, and suggestion pills
gitmap heatmap

# Filter heatmap telemetry by repository alias or slug
gitmap heatmap --repo gitmap
gitmap rc heatmap --repo gitmap

# Hierarchical causal execution tree for a specific task
gitmap heatmap tree --task "314-gitmap-heatmap-and-url-handling-instructions"

# Non-blocking SQL query against cognitive action logs
gitmap heatmap query "SELECT FilePath, ReasoningState, ExpectedOutcome, ActualOutcome FROM HeatmapActionLog WHERE IsExpectationMet = 1 ORDER BY HeatmapActionLogId DESC LIMIT 10" --json
```

#### Technical Guarantees
- **3-Tier Introspection:** Surveys Local Root (`<EffectiveDataDir>/heatmap.db`), Local Per-Repo Split DBs (`<EffectiveDataDir>/heatmap/<repo>/sql.db`), and Multi-AI Repo-Cache DBs (`repo-cache/.heatmap/<seq>-<ai>/`).
- **Alias Resolution:** Automatically enriches database entries with registered aliases from `gitmap.db` via `db.ListAliasesWithRepo()`.
- **Multi-Repo Task Aggregation:** Aggregates top 5 to 10 recent tasks across 2 to 5 repositories with completion rates and action counts.
- **Dynamic AI Suggestion Pills:** Emits high-leverage command suggestions directly into terminal output.

---

### Anti-Pattern 4: Legacy Search Tools vs `gitmap aum search`

#### Context & Forbidden Pattern
Using slow, unbuffered external utilities like `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, or `findstr` triggers high CPU consumption, unbuffered text flooding, and terminal freezing.

```powershell
<!-- ❌ FORBIDDEN: Slow external search utilities -->
Select-String -Path "cli/*.go" -Pattern "Heatmap"
Get-ChildItem -Recurse -Filter "*.go" | Select-String "RunRCHeatmap"
grep -rn "url_cache" cli/
rg "DownloadOptions" cli/
findstr /s /i "DispatchHeatmap" *.go
```

#### Required Replacement
Always use `gitmap aum search` (or `gitmap search`, `gitmap find`, and `gitmap cat`).

```bash
<!-- ✅ REQUIRED: High-speed GitMap AUM search and inspection -->
# Multi-core streaming live search (<15ms) scoped by directory and extension
gitmap aum search "HeatmapActionLog" cli -e .go

# Fast symbol and keyword search backed by DH2D hot cache
gitmap search "UrlCacheEntry"

# Sub-millisecond file pattern discovery
gitmap find "*heatmap*.go" cli

# Zero-disk direct stream of file content into process stdout
gitmap cat cli/cmdurl/cmd_url.go

# Flexible path parsing: handles bracket notation, markdown embeds, repo-cache paths, bullets, and comments
gitmap cat @[<drive>:\path\to\media_1792894009084_45b4893f.png]
gitmap cat "![Attached Image 1](file:///<drive>:/path/to/media_1792894009084_45b4893f.png)"
gitmap cat "1. @[assets/media_1792894009084_45b4893f.png] # consider the file is relative from the repo"
gitmap cat "1. @[../repo-cache/assets/media_1792894009084_45b4893f.png] # consider the file is relative from the repo outside of the repo"

# Safe image asset inspection: renders structured TokyoNight cards (dimensions, size, status) instead of binary noise
gitmap cat assets/screenshots/media_1792894009084_45b4893f.png
```

---

### Anti-Pattern 5: Absolute Paths vs Relative Git-Root Paths

#### Context & Forbidden Pattern
Hardcoding drive letters (`C:\`, `D:\`), machine-specific home folders, or `file:///` URIs breaks portability across machines, containerized runners, and multi-agent swarm environments.

```markdown
<!-- ❌ FORBIDDEN: Absolute filesystem paths and drive letters -->
[Spec File](file:///repo-root/02-spec/21-app/314-overview.md)
Path: C:\Users\Administrator\.gitmap\heatmap.db
Include: /var/lib/gitmap/repodb/url_cache.db
```

#### Required Replacement
All paths in specifications, plans, code references, and documentation must be relative to the Git repository root.

```markdown
<!-- ✅ REQUIRED: Strictly relative paths from workspace root -->
[Spec File](02-spec/21-app/314-gitmap-heatmap-and-url-handling-instructions/00-overview.md)
Database: repodb/url_cache.db
Source: cli/cmdurl/cmd_url.go
Plan: .ai-memory/plans/subtasks/314-gitmap-heatmap-and-url-handling-instructions/07-subtask-coding-guidelines-and-skills.md
```

---

## 4. Coding Conventions & Invariants

When implementing or invoking GitMap replacements, adhere to these non-negotiable coding conventions:

1. **Positive Boolean Naming:** All booleans must use `is*` or `has*` prefixes exclusively. Negative booleans (`isNotCached`, `disableAria2`) are strictly forbidden.
   - ✅ `isCacheHit`, `isFresh`, `hasAria2c`, `isRangeSupported`, `isActive`
   - ❌ `notCached`, `noRange`, `canDownload`, `shouldCache`

2. **15-Line Function Limit:** Functions must remain concise ($\le 8$ lines preferred, $\le 15$ lines maximum). Break complex pipelines into discrete composable helpers.

3. **Zero Nested Ifs:** Flatten conditional branches using early returns and guard clauses.

4. **Error Management Sovereignty:** Return structured errors using `apperror.WrapSimple(err, "context")`. Never call `os.Exit` or exit handlers inside leaf functions.

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/06-ai-optimization/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-AI-010: Anti-Pattern Command Replacements Enforcement

**Given** AI agents operating within the repository and codebase guidelines.
**When** Audited against this anti-pattern replacement specification.
**Then** Zero prohibited commands, raw HTTP invocations, filesystem grep traversals, or absolute paths are detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only
```
**Expected:** exit 0. Zero violations.
