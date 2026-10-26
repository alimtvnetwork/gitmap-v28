# Forensic Audit Ledger: GitMap Author Information and Search Feature
**Sequence ID:** 335  
**Task Slug:** gitmap-author-information-and-search-feature  
**Timestamp:** 2026-10-27  

---

## 1. Executive Summary & Root Cause Analysis (RCA)
- **Problem Statement 1 (Commit Author Information):**
  When performing commits and pushes (`gitmap cp`, `gitmap cpr`, `gitmap cpf`, `gitmap cbf`, `gitmap cfr`, `gitmap cpc`), the commit summary terminal card (`renderCommitPushSummaryCard`) renders the commit SHA, branch, status, and message, but omits the author identity. Furthermore, users cannot specify a custom author (`--author` / `-a`), nor does the CLI clearly indicate which local Git author (`user.name <user.email>`) signed the commit.
- **Problem Statement 2 (Interactive Repository Search UI):**
  While GitMap has an ultra-fast internal multi-threaded search engine via AUM (`cmdautomation.RunSearch`), users currently lack a dedicated, interactive web search UI. The UI needs to allow repository-scoped search with rich hover previews (matching line numbers, snippet content, context) and direct 1-click open capabilities (opening the target file in the built-in editor or the system editor / Antigravity IDE).

---

## 2. Forensic Requirement Mapping
| ID | Requirement | Component / File Target | Verification Standard |
|---|---|---|---|
| REQ-01 | Add `--author` / `-a` flag and author resolution | `cli/cmdcommit/commit_flags.go`, `commit_types.go` | Parsed into `CommitMetadata.Author`; falls back to `git config user.name <email>`. |
| REQ-02 | Execute commit with author attribution | `cli/cmdcommit/commit_push.go` | Appends `--author="Name <email>"` to `git commit` if set. |
| REQ-03 | Render Author in Terminal Summary Card | `cli/cmdcommit/commit_ui.go` | Displays `Author: Name <email>` in Lipgloss/term card. |
| REQ-04 | Add `/api/search/repo` Endpoint | `cli/cmdui/ui_server.go` | Calls `cmdautomation.RunSearch`, returns structured JSON results with matches and line context. |
| REQ-05 | Add Search UI Tab to Fleet Web UI | `cli/cmdui/ui_assets_markup_misc.go`, `ui_assets_js_settings.go` | Added to sidebar navigation (`NAV_ITEMS`); responsive search view with filter options. |
| REQ-06 | Interactive Search Hover & File Opener | `cli/cmdui/ui_assets_js_fleet.go` | Hover over search results shows line preview; clicking opens in Editor or calls `/api/project/open-antigravity`. |
| REQ-07 | CLI Command `gitmap search --ui` / `repo-search --ui` | `cli/cmdsearch/search.go`, `search_entry.go` | Launches Web UI directly into Search mode. |

---

## 3. Strict Coding Guideline Enforcement
- All functions <= 15 lines max.
- Zero nested if statements (guard clauses only).
- Positive booleans (`is*`, `has*` only).
- Dedicated semantic commits: `gitmap cpc`, `gitmap cpf`, `gitmap cbf`, `gitmap cfr`, `gitmap cpr`.
