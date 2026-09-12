# Interactive Macro Builder PWD Header, In-Builder LS & Helper Commands

**Updated:** 2026-09-09
**Specification:** `spec/21-macro-builder-enhancements/`
**Reference Plan:** `.lovable/plans/pending/83-interactive-macro-builder-pwd-ls-search.md`

---

## 1. Context & Architecture

During interactive macro creation (`gitmap macro add <name>`), users enter a sequence of commands to be executed sequentially or recorded into automated recipes.

### Key UX Deficiencies Identified:
1. **Missing PWD Awareness:** In interactive prompt mode (`Step N> `), users are often navigating or composing commands without visibility into their active working directory (`PWD`).
2. **Interactive `ls` Ambiguity:** When users type `ls` in the prompt, they frequently intend to inspect files in the current directory before choosing which command to execute next. The builder must support executing directory listing in-place while asking whether to record it or continue composition.
3. **In-Builder Helper Commands:** Composition of macros benefits from fast in-prompt utilities:
   - `find <pattern>`: scans directory tree for matching files.
   - `search <query>`: greps text contents of repository/workspace files.
   - `replace <old> <new> [glob]`: performs surgical string replacement in files.
   - `pwd` / `pwd on` / `pwd off`: toggles current directory header display above prompt line.

---

## 2. PWD Header Specification

Above each `Step N> ` input line, render:
```text
  [PWD: /path/to/current/workdir]
  Step 1>
```
- Toggleable via command line flag `--pwd` / `--no-pwd`.
- Toggleable in-session via `pwd on` and `pwd off`.
- State saved to user preferences in `uipref/` (`IsMacroPwdVisible`).

---

## 3. Offline CI/CD Release Installer Architecture

In local CI/CD execution environments:
- Remote GitHub releases may not be published immediately when version bumps occur locally (e.g. `v6.199.0` while `v6.198.0` is published).
- `smoke-installer.py` automatically checks if the target release asset exists on GitHub via HEAD probe.
- When absent on GitHub, it automatically creates a temporary local release archive (`.zip` on Windows, `.tar.gz` on Linux) with sha256 checksums, spins up a Python `http.server.ThreadingHTTPServer` on localhost, and directs `install.ps1` or `install.sh` to download and install via `GITMAP_DOWNLOAD_URL`.
- This guarantees 100% test coverage of the installer download, sha256 verification, and archive extraction logic locally without network failure or retry stalls.

---

## 4. Multi-Agent Pipeline Section Concurrency

`03-ai-scripts/06-cicd-local-runner.py` executes quality gates across 3 autonomous section agents:
1. **Agent 1:** Linters, AST Checks & Static Analyzers (22 gates across thread workers).
2. **Agent 2:** Go Compilation, Frontend Build & E2E Smoke Suite.
3. **Agent 3:** Go Smart Incremental Unit Tests, Coverage Profiling & Race Detection.

Each agent operates concurrently across all logical CPU cores (`os.cpu_count() or 16`), with results synchronized through thread locks (`DISK_WRITE_LOCK`).
