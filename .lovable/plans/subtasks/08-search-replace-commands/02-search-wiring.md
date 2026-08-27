# Subtask 2: Search Wiring

1. Rewrite `gitmap/cmd/search_entry.go` to be fully featured.
2. Use the same DB connection block (`getFindDB`) as `find_entry.go` (extract it to `db_resolver.go` to avoid duplication).
3. Connect `runSearch` and `runRepoSearch` to `searcher.SearchRepoDB`.
4. Connect `runRepoRegex` to `searcher.SearchRepoDBRegex`.
5. Connect JSON variants to output stringified results instead of `pterm`.
6. For `runSearchReplaceAll` (reset), add logic to delete `.db` files in `repo_search/` or clear cache tables.

- [ ] Task incomplete
