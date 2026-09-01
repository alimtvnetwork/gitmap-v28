# Specification 16: Scan Nested WorkDirs, Terminal UI Polish, Interactive Macro Recording, and Cluster Join Daemon

## 1. Executive Summary

This specification defines major enhancements to the GitMap ecosystem across seven key pillars:
1. **Nested Git Scanning & Work Directory Grouping**: Deep recursive scanning into nested directories with two presentation modes (flat table vs. folder/work-directory-grouped table), automatic default work directory tracking, and interactive missing repository resolution (`gitmap rm` / relocate path).
2. **Terminal UI & Glyph Robustness**: Eradication of broken ASCII bracket fallbacks (e.g. `[dna]`, `[]`, `v clean`) in favor of universal UTF-8 terminal symbols (`✔`, `●`, `✖`, `➜`, `📂`, `⚡`) verified across PowerShell 7, Windows Terminal, and Linux/macOS consoles.
3. **Parallel Clone Visualization**: Fixing line-overlapping in `gitmap clone gitmap.json` with dedicated multi-worker status lines, clear distinctions for existing repositories (safe-pull vs skip), and per-repo elapsed timers.
4. **Help System & DAC / Fix-Repo Redesign**: Re-architecting cluttered utility help pages (`fix-repo`, `doctor`, `setup`), highlighting cluster commands (`serve`, `servers-clients`, `clients`, `cluster nodes`, `clients ls`, `cluster history`), enabling multi-keyword search discovery (`join`, `servers`, `cluster`), and expanding the Web Help Dashboard (`gitmap hd`).
5. **Multi-Commit Replay & Auto-Repo Creation**: Verifying and extending `commit-in`, `commit-left`, and `commit-right` to support true multi-commit sequential replays, automatic target repository scaffolding/creation when missing, and enhanced SEO commit generation.
6. **Interactive Mode Macro Recording & Execution**: Real-time interactive session recording (`gitmap interactive --record <name>`), logging terminal inputs and execution outcomes into SQLite, macro replay execution (`gitmap execute <name>`), and recurring cron/interval scheduling.
7. **Cluster Join Daemon & Lightweight Node Agent**: Designing a standalone lightweight CLI (`gitmap-node-join` / `gitmap-agent`) configured to execute on Windows startup, automatically join central orchestrators, and provide clear diagnostic warnings for unreachable cluster nodes.

## 2. Referenced Visual Artifacts

The following user-provided screenshots illustrate current production deficiencies and anchor the requirements in this specification:

- `images/01-parallel-clone-terminal-overlap.png`: Demonstrates parallel clone stdout corruption where `[ 1/14] Cloning repo1 ...[ 2/14] Cloning repo2 ...` concatenated without newlines and missing glyphs.
- `images/02-gitmap-rm-prompt-architect.png`: Demonstrates the repository untrack command (`gitmap rm <path>`) used in missing repo remediation.
- `images/03-gitmap-status-v-clean-fallback.png`: Demonstrates status table output with fallback `v clean` markers.
- `images/04-gitmap-status-missing-repos.png`: Demonstrates `26 repos · 24 clean · 2 missing` with `[] not found` missing repo markers.
- `images/05-gitmap-scan-bracket-emojis.png`: Demonstrates `gitmap scan` artifact summary displaying bracketed text strings (`[chart]`, `[dna]`, `[tree]`, `[wand]`, `[db]`) instead of rich visual glyphs.

## 3. Strict Architectural & Code Constraints

All implementations derived from this specification MUST strictly follow:
1. **Zero-Swallow Error Policy**: Every non-nil error MUST be logged to `os.Stderr` or structured logger.
2. **File Size Budget**: Maximum 200 lines per source file. Modularize components into distinct sub-files.
3. **Function Size Budget**: Maximum 15 lines per function.
4. **Positive Logic Flow**: Conditionals MUST use positive boolean expressions (`hasX`, `isY`), no raw negations (`!`), and early returns.
5. **No Magic Strings**: All CLI flags, prompts, table headers, and status strings MUST be centralized in `constants/`.
6. **No Raw Direct File Manipulation in `.gitmap/release/`**: Automated release tooling MUST be used for versioning ceremonies.
