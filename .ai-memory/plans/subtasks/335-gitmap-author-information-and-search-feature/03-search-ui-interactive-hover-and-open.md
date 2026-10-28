# Subtask 03: Search UI Interactive Hover and File Opener
**Sequence ID:** 335  
**Parent Plan:** [335-gitmap-author-information-and-search-feature.md](.ai-memory/plans/pending/335-gitmap-author-information-and-search-feature.md)  
**Status:** Pending  

---

## Actions
1. Implement JS logic in `cli/cmdui/ui_assets_js_fleet.go` to execute `/api/search/repo`.
2. Add interactive hover preview displaying contextual lines around matched content.
3. Wire click-to-open to switch to `tab-editor` or call `/api/project/open-antigravity`.
4. Support `gitmap search --ui` and `gitmap repo-search --ui` flags in `cli/cmdsearch/search.go` / `search_entry.go`.
